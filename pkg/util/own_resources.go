package util

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/go-logr/logr"
	csvv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	operatorcontv1 "github.com/operator-framework/operator-controller/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// OwnResources holds the running POd, Deployment and CSV, if exist
type OwnResources struct {
	pod        *corev1.Pod
	deployment *appsv1.Deployment
	manageObj  client.Object
}

// GetPod returns the running pod, or nil if not exists
func (or *OwnResources) GetPod() *corev1.Pod {
	if or == nil {
		return nil
	}
	return or.pod.DeepCopy()
}

// GetDeployment returns the deployment that controls the running pod, or nil if not exists
func (or *OwnResources) GetDeployment() *appsv1.Deployment {
	if or == nil {
		return nil
	}
	return or.deployment
}

// GetManageObject returns the object that defines the application, or nil if not exists
func (or *OwnResources) GetManageObject() client.Object {
	if or == nil {
		return nil
	}

	return or.manageObj
}

func findOwnResources(ctx context.Context, cl client.Reader, logger logr.Logger) *OwnResources {
	or := &OwnResources{}

	if !GetClusterInfo().IsRunningLocally() {
		var err error

		or.pod, err = getPod(ctx, cl, logger)
		if err != nil {
			or.pod = nil
			logger.Error(err, "Can't get self pod")
			return or
		}

		operatorNs := GetOperatorNamespaceFromEnv()
		or.deployment, err = getDeploymentFromPod(ctx, or.pod, cl, operatorNs, logger)
		if err != nil {
			logger.Error(err, "Can't get deployment")
			return or
		}

		or.getOwnerReference(ctx, cl, logger, operatorNs)

	} else {
		deployment := &appsv1.Deployment{}
		err := cl.Get(ctx, client.ObjectKey{
			Namespace: os.Getenv("OPERATOR_NAMESPACE"),
			Name:      "hyperconverged-cluster-operator",
		}, deployment)
		if err != nil {
			logger.Error(err, "Can't get deployment")
			return or
		}
		or.deployment = deployment
		or.pod = nil
		or.manageObj = nil
	}

	return or
}

func getPod(ctx context.Context, c client.Reader, logger logr.Logger) (*corev1.Pod, error) {
	ci := GetClusterInfo()
	operatorNs := GetOperatorNamespaceFromEnv()

	// This is taken from k8sutil.getPod. This method only receives client. But the client is not always ready. We'll
	// use --- instead
	if ci.IsRunningLocally() {
		return nil, nil
	}
	podName := os.Getenv(PodNameEnvVar)
	if podName == "" {
		return nil, fmt.Errorf("required env %s not set, please configure downward API", PodNameEnvVar)
	}

	pod := &corev1.Pod{}
	key := client.ObjectKey{Namespace: operatorNs, Name: podName}
	err := c.Get(ctx, key, pod)
	if err != nil {
		logger.Error(err, "Failed to get Pod", "Pod.Namespace", operatorNs, "Pod.Name", podName)
		return nil, err
	}

	// .Get() clears the APIVersion and Kind,
	// so we need to set them before returning the object.
	pod.APIVersion = "v1"
	pod.Kind = "Pod"

	logger.Info("Found Pod", "Pod.Namespace", operatorNs, "Pod.Name", pod.Name)

	return pod, nil
}

func getDeploymentFromPod(ctx context.Context, pod *corev1.Pod, c client.Reader, operatorNs string, logger logr.Logger) (*appsv1.Deployment, error) {
	if pod == nil {
		return nil, nil
	}
	rsReference := metav1.GetControllerOf(pod)
	if rsReference == nil || rsReference.Kind != "ReplicaSet" {
		err := errors.New("failed getting HCO replicaSet reference")
		logger.Error(err, "Failed getting HCO replicaSet reference")
		return nil, err
	}
	rs := &appsv1.ReplicaSet{}
	err := c.Get(context.TODO(), client.ObjectKey{
		Namespace: operatorNs,
		Name:      rsReference.Name,
	}, rs)
	if err != nil {
		logger.Error(err, "Failed to get HCO ReplicaSet")
		return nil, err
	}

	dReference := metav1.GetControllerOf(rs)
	if dReference == nil || dReference.Kind != "Deployment" {
		err = errors.New("failed getting HCO deployment reference")
		logger.Error(err, "Failed getting HCO deployment reference")
		return nil, err
	}
	deployment := &appsv1.Deployment{}
	err = c.Get(ctx, client.ObjectKey{
		Namespace: operatorNs,
		Name:      dReference.Name,
	}, deployment)
	if err != nil {
		logger.Error(err, "Failed to get HCO Deployment")
		return nil, err
	}

	return deployment, nil
}

func (or *OwnResources) getOwnerReference(ctx context.Context, cl client.Reader, logger logr.Logger, operatorNs string) {
	if or.deployment == nil {
		return
	}

	if GetClusterInfo().IsManagedByOLMV0() {
		csv, err := getCSVFromDeployment(ctx, or.deployment, cl, operatorNs, logger)
		if err != nil {
			logger.Error(err, "Can't get CSV")
			or.manageObj = nil
		} else if csv != nil {
			or.manageObj = csv
		}
	} else {
		clusterExt, err := getClusterExtensionFromDeployment(ctx, or.deployment, cl, logger)
		if err != nil {
			logger.Error(err, "Can't get ClusterExtension")
			or.manageObj = nil
		} else if clusterExt != nil {
			or.manageObj = clusterExt
		}
	}
}

func getCSVFromDeployment(ctx context.Context, d *appsv1.Deployment, c client.Reader, operatorNs string, logger logr.Logger) (*csvv1alpha1.ClusterServiceVersion, error) {
	idx := slices.IndexFunc(d.GetOwnerReferences(), func(owner metav1.OwnerReference) bool {
		return owner.Kind == csvv1alpha1.ClusterServiceVersionKind
	})

	if idx < 0 {
		err := errors.New("no CSV owner reference reference")
		return nil, err
	}

	csvReference := d.GetOwnerReferences()[idx]

	csv := &csvv1alpha1.ClusterServiceVersion{}
	err := c.Get(ctx, client.ObjectKey{
		Namespace: operatorNs,
		Name:      csvReference.Name,
	}, csv)

	if err != nil {
		logger.Error(err, "Failed to get HCO CSV")
		return nil, err
	}

	return csv, nil
}

func getClusterExtensionFromDeployment(ctx context.Context, d *appsv1.Deployment, c client.Reader, logger logr.Logger) (*operatorcontv1.ClusterExtension, error) {
	idx := slices.IndexFunc(d.GetOwnerReferences(), func(owner metav1.OwnerReference) bool {
		return owner.Kind == operatorcontv1.ClusterExtensionKind
	})

	if idx < 0 {
		logger.Info("No known owner reference")
		return nil, nil
	}

	clusterExtReference := d.GetOwnerReferences()[idx]
	clusterExt := &operatorcontv1.ClusterExtension{}
	err := c.Get(ctx, client.ObjectKey{
		Name: clusterExtReference.Name,
	}, clusterExt)

	if err != nil {
		return nil, fmt.Errorf("failed to get HCO ClusterExtension; %w", err)
	}

	return clusterExt, nil
}
