package perses

import (
	"context"
	"maps"
	"slices"

	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	persesv1alpha1 "github.com/rhobs/perses-operator/api/v1alpha1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kubevirt/hyperconverged-cluster-operator/controllers/commontestutils"
)

var _ = Describe("Perses controller", func() {
	var (
		dashboards map[string]persesv1alpha1.PersesDashboard
		datasource *persesv1alpha1.PersesDatasource
		s          *runtime.Scheme
		startupReq reconcile.Request
	)

	BeforeEach(func() {
		s = scheme.Scheme
		Expect(apiextensionsv1.AddToScheme(s)).To(Succeed())
		Expect(persesv1alpha1.AddToScheme(s)).To(Succeed())

		var err error
		dashboards, err = initDashboards(commontestutils.Namespace, GinkgoLogr)
		Expect(err).ToNot(HaveOccurred())
		Expect(dashboards).ToNot(BeEmpty())

		datasource, err = initDatasource(commontestutils.Namespace)
		Expect(err).ToNot(HaveOccurred())

		startupReq = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name: startupReqType + randomSufix,
			},
		}
	})

	Context("PersesReconciler end-to-end behavior", func() {
		makeCRD := func(name string) *apiextensionsv1.CustomResourceDefinition {
			return &apiextensionsv1.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
				},
			}
		}

		It("should apply datasource and dashboard when CRDs are available", func(ctx context.Context) {
			ctx = logr.NewContext(ctx, GinkgoLogr)
			tracker := testing.NewObjectTracker(s, serializer.NewCodecFactory(s).UniversalDecoder())
			Expect(tracker.Add(makeCRD("persesdashboards.perses.dev"))).To(Succeed())
			Expect(tracker.Add(makeCRD("persesdatasources.perses.dev"))).To(Succeed())

			cl := fake.NewClientBuilder().
				WithScheme(s).
				WithObjectTracker(tracker).
				Build()

			r := &PersesReconciler{
				Client:           cl,
				namespace:        commontestutils.Namespace,
				cachedDatasource: datasource,
				cachedDashboards: dashboards,
			}
			_, err := r.Reconcile(ctx, startupReq)
			Expect(err).ToNot(HaveOccurred())

			ds := &persesv1alpha1.PersesDatasource{}
			Expect(cl.Get(ctx, client.ObjectKey{Namespace: commontestutils.Namespace, Name: datasource.Name}, ds)).To(Succeed())

			dbNames := slices.Collect(maps.Keys(dashboards))
			Expect(dbNames).ToNot(BeEmpty())
			Expect(dbNames).To(ContainElement("perses-dashboard-node-memory-overview"))

			for name := range dashboards {
				db := &persesv1alpha1.PersesDashboard{}
				Expect(cl.Get(ctx, client.ObjectKey{Namespace: commontestutils.Namespace, Name: name}, db)).To(Succeed())
			}
		})

		It("should skip reconcile when Perses CRDs are missing", func(ctx context.Context) {
			ctx = logr.NewContext(ctx, GinkgoLogr)
			tracker := testing.NewObjectTracker(s, serializer.NewCodecFactory(s).UniversalDecoder())
			cl := fake.NewClientBuilder().
				WithScheme(s).
				WithObjectTracker(tracker).
				Build()

			r := &PersesReconciler{
				Client:           cl,
				namespace:        commontestutils.Namespace,
				cachedDatasource: datasource,
				cachedDashboards: dashboards,
			}
			_, err := r.Reconcile(ctx, startupReq)
			Expect(err).ToNot(HaveOccurred())

			dsList := &persesv1alpha1.PersesDatasourceList{}
			Expect(cl.List(ctx, dsList)).To(Succeed())
			Expect(dsList.Items).To(BeEmpty())

			dbList := &persesv1alpha1.PersesDashboardList{}
			Expect(cl.List(ctx, dbList)).To(Succeed())
			Expect(dbList.Items).To(BeEmpty())
		})
	})

	Context("SetupPersesWithManager guard", func() {
		It("should skip controller registration when Perses CRDs are not available", func() {
			old := checkPersesAvailable
			checkPersesAvailable = func(_ context.Context, _ client.Client) bool { return false }
			defer func() { checkPersesAvailable = old }()

			// Build a no-CRD fake client and a lightweight manager mock
			cl := fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
			mgr, err := commontestutils.NewManagerMock(nil, manager.Options{Scheme: scheme.Scheme}, cl, GinkgoLogr)
			Expect(err).ToNot(HaveOccurred())

			err = SetupPersesWithManager(mgr, metav1.OwnerReference{})
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
