package clusterinfo

import (
	"context"

	"github.com/go-logr/logr"
	openshiftconfigv1 "github.com/openshift/api/config/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	hcoutil "github.com/kubevirt/hyperconverged-cluster-operator/pkg/util"
)

type ClusterInfoMock struct {
	isOpenshift      bool
	runningLocally   bool
	isManagedByOLMV0 bool
}

// make sure ClusterInfoMock implements hcoutil.Clus`terInfo
var _ hcoutil.ClusterInfo = &ClusterInfoMock{}

func NewClusterInfoMock(opts ...MockOption) *ClusterInfoMock {
	mock := &ClusterInfoMock{}

	for _, op := range opts {
		op(mock)
	}

	return mock
}

func (ClusterInfoMock) Init(_ context.Context, _ client.Client, _ logr.Logger) error {
	return nil
}

func (c ClusterInfoMock) IsOpenshift() bool {
	return c.isOpenshift
}

func (c ClusterInfoMock) IsRunningLocally() bool {
	return c.runningLocally
}

func (ClusterInfoMock) GetBaseDomain() string {
	return ""
}

func (c ClusterInfoMock) IsManagedByOLMV0() bool {
	return c.isManagedByOLMV0
}

func (ClusterInfoMock) IsConsolePluginImageProvided() bool {
	return true
}

func (ClusterInfoMock) IsMonitoringAvailable() bool {
	return true
}

func (ClusterInfoMock) IsDeschedulerAvailable() bool {
	return true
}

func (ClusterInfoMock) IsNADAvailable() bool {
	return true
}

func (ClusterInfoMock) IsDeschedulerCRDDeployed(ctx context.Context, cl client.Client) bool {
	return true
}

func (ClusterInfoMock) IsSingleStackIPv6() bool {
	return true
}

func (ClusterInfoMock) GetTLSSecurityProfile(_ *openshiftconfigv1.TLSSecurityProfile) *openshiftconfigv1.TLSSecurityProfile {
	return nil
}

func (ClusterInfoMock) RefreshAPIServerCR(_ context.Context, _ client.Client) error {
	return nil
}

func (ClusterInfoMock) GetPod() *corev1.Pod                 { return nil }
func (c ClusterInfoMock) GetDeployment() *appsv1.Deployment { return nil }
func (c ClusterInfoMock) GetManageObject() runtime.Object   { return nil }
