package perses

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	persesv1alpha1 "github.com/rhobs/perses-operator/api/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	hcoutil "github.com/kubevirt/hyperconverged-cluster-operator/pkg/util"
)

const (
	dashboardReqType   = "dashboard"
	dashboardReqSufix  = "-" + dashboardReqType
	datasourceReqType  = "datasource"
	datasourceReqSufix = "-" + datasourceReqType
	startupReqType     = "startup"
	unknownReqType     = "unknown"
)

var (
	persesLog   = logf.Log.WithName("controller_observability_perses")
	randomSufix = "-" + string(uuid.NewUUID())
)

// PersesReconciler handles Perses dashboards, datasources and token secret.
type PersesReconciler struct {
	client.Client

	namespace string
	events    chan event.GenericEvent
	owner     metav1.OwnerReference

	// Cached, parsed Perses assets loaded from the image (read once per process)
	cachedDashboards map[string]persesv1alpha1.PersesDashboard
	cachedDatasource *persesv1alpha1.PersesDatasource
}

func (r *PersesReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	if !hcoutil.IsPersesAvailable(ctx, r.Client) {
		return reconcile.Result{}, nil
	}

	reqLog := logr.FromContextOrDiscard(ctx)
	reqType, reqName := resolveRequest(req)

	persesLog.Info("Reconciling Perses", "Request.Namespace", req.Namespace, "Request.Name", reqName)

	var err error

	switch reqType {
	case dashboardReqType:
		err = r.reconcileDashboard(ctx, reqName, req.Namespace, reqLog)
	case datasourceReqType:
		err = r.reconcileDataSource(ctx, reqName, req.Namespace, reqLog)
	case startupReqType:
		err = r.reconcileAll(ctx, reqLog)
	default:
		reqLog.Info("unknow request; ignoring.", "Request.Namespace", req.Namespace, "Request.Name", req.Name, reqType, "type")
		return reconcile.Result{}, nil
	}

	if err != nil {
		reqLog.Error(err, "failed to reconcile Perses", "Request.Namespace", req.Namespace, "Request.Name", req.Name, "requestType", reqType)
	}

	return reconcile.Result{}, err
}

// SetupPersesWithManager registers the Perses controller with the manager.
func SetupPersesWithManager(mgr manager.Manager, ownerRef metav1.OwnerReference) error {
	persesLog.Info("Setting up Perses controller")

	namespace := hcoutil.GetOperatorNamespaceFromEnv()
	dashboards, err := initDashboards(namespace, persesLog)
	if err != nil {
		return err
	}
	datasource, err := initDatasource(namespace)
	if err != nil {
		return err
	}

	r := newPersesReconciler(mgr, namespace, ownerRef, dashboards, datasource)

	c, err := controller.New("perses controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch dashboards with predicate on managed list
	if err = c.Watch(
		source.Kind[*persesv1alpha1.PersesDashboard](mgr.GetCache(), &persesv1alpha1.PersesDashboard{},
			handler.TypedEnqueueRequestsFromMapFunc[*persesv1alpha1.PersesDashboard, reconcile.Request](func(ctx context.Context, dashboard *persesv1alpha1.PersesDashboard) []reconcile.Request {
				return []reconcile.Request{
					{
						NamespacedName: types.NamespacedName{
							Namespace: dashboard.Namespace,
							Name:      fmt.Sprintf("%s%s%s", dashboard.Name, dashboardReqSufix, randomSufix),
						},
					},
				}
			}),
			dashboardPredicate,
		),
	); err != nil {
		return err
	}

	// Watch the single datasource we manage
	if err = c.Watch(
		source.Kind[*persesv1alpha1.PersesDatasource](mgr.GetCache(), &persesv1alpha1.PersesDatasource{},
			handler.TypedEnqueueRequestsFromMapFunc[*persesv1alpha1.PersesDatasource, reconcile.Request](func(ctx context.Context, datasource *persesv1alpha1.PersesDatasource) []reconcile.Request {
				return []reconcile.Request{
					{
						NamespacedName: types.NamespacedName{
							Namespace: datasource.Namespace,
							Name:      fmt.Sprintf("%s%s%s", datasource.Name, datasourceReqSufix, randomSufix),
						},
					},
				}
			}),
			datasourcePredicate,
		),
	); err != nil {
		return err
	}

	// Trigger startup reconcile
	if err = c.Watch(
		source.Channel(r.events, handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, _ client.Object) []reconcile.Request {
			return []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Name: fmt.Sprintf("%s%s", startupReqType, randomSufix),
					},
				},
			}
		})),
	); err != nil {
		return err
	}

	r.forceFirstRequest()

	return nil
}

func newPersesReconciler(
	mgr manager.Manager,
	namespace string,
	ownerRef metav1.OwnerReference,
	dashboards map[string]persesv1alpha1.PersesDashboard,
	datasource *persesv1alpha1.PersesDatasource,
) *PersesReconciler {

	return &PersesReconciler{
		Client:    mgr.GetClient(),
		namespace: namespace,
		events:    make(chan event.GenericEvent, 1),
		owner:     ownerRef,

		cachedDashboards: dashboards,
		cachedDatasource: datasource,
	}
}

func (r *PersesReconciler) reconcileDashboard(ctx context.Context, name string, namespace string, logger logr.Logger) error {
	db, ok := r.cachedDashboards[name]
	if !ok || db.Namespace != namespace {
		logger.Info("Not a managed dashboard; ignoring", "namespace", namespace, "name", name)
		return nil
	}

	foundDashboard := &persesv1alpha1.PersesDashboard{}
	err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, foundDashboard)
	if err != nil {
		foundDashboard = db.DeepCopy()
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, foundDashboard)
		}
		return err
	}

	if reflect.DeepEqual(foundDashboard.Spec, db.Spec) {
		db.Spec.DeepCopyInto(&foundDashboard.Spec)
		return r.Update(ctx, foundDashboard)
	}

	return nil
}

func (r *PersesReconciler) reconcileDataSource(ctx context.Context, name string, namespace string, logger logr.Logger) error {
	if r.cachedDatasource.Name != name || r.cachedDatasource.Namespace != namespace {
		logger.Info("Not a managed dashboard; ignoring", "namespace", namespace, "name", name)
		return nil
	}

	foundDatasource := &persesv1alpha1.PersesDatasource{}
	err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, foundDatasource)
	if err != nil {
		foundDatasource = r.cachedDatasource.DeepCopy()
		if k8serrors.IsNotFound(err) {
			return r.Create(ctx, foundDatasource)
		}
		return err
	}

	if reflect.DeepEqual(r.cachedDatasource.Spec, foundDatasource.Spec) {
		r.cachedDatasource.Spec.DeepCopyInto(&foundDatasource.Spec)
		return r.Update(ctx, foundDatasource)
	}

	return nil
}

func (r *PersesReconciler) reconcileAll(ctx context.Context, logger logr.Logger) error {
	var (
		reconcileErrors []error
		err             error
	)
	for name := range r.cachedDashboards {
		err = r.reconcileDashboard(ctx, name, r.namespace, logger)
		if err != nil {
			reconcileErrors = append(reconcileErrors, err)
		}
	}

	err = r.reconcileDataSource(ctx, datasourceName, r.namespace, logger)
	if err != nil {
		reconcileErrors = append(reconcileErrors, err)
	}

	return errors.Join(reconcileErrors...)
}

// forceFirstRequest makes sure that the Reconcile method is called at least once, for enforce the creation of the
// resources on a fresh installation.
func (r *PersesReconciler) forceFirstRequest() {
	r.events <- event.GenericEvent{
		Object: &metav1.PartialObjectMetadata{},
	}
	close(r.events)
}

func resolveRequest(req reconcile.Request) (reqType, resourceName string) {
	if !strings.HasSuffix(req.Name, randomSufix) {
		return unknownReqType, ""
	}

	reqType = strings.TrimSuffix(req.Name, randomSufix)
	if reqType == startupReqType {
		return startupReqType, ""
	}

	if !strings.HasSuffix(reqType, dashboardReqSufix) {
		return dashboardReqType, strings.TrimSuffix(reqType, dashboardReqSufix)
	}

	if !strings.HasSuffix(reqType, datasourceReqSufix) {
		return datasourceReqType, strings.TrimSuffix(reqType, datasourceReqSufix)
	}

	return unknownReqType, ""
}
