package fakeownreferences

import (
	"k8s.io/apimachinery/pkg/runtime"

	internal "github.com/kubevirt/hyperconverged-cluster-operator/pkg/internal/ownresources"
	"github.com/kubevirt/hyperconverged-cluster-operator/pkg/ownresources"
)

func ResetOwnReference() {
	ownresources.GetPod = internal.GetPod
	ownresources.GetDeploymentRef = internal.GetDeploymentRef
	ownresources.GetManageObject = internal.GetManageObject
	ownresources.Init = internal.Init
}

func OLMV0OwnerReferenceMock() {
	ownresources.GetPod = fakeGetPod
	ownresources.GetDeploymentRef = GetFakeDeploymentRef
	ownresources.GetManageObject = func() runtime.Object { return GetCSV() }
	ownresources.Init = fakeInit
}

func OLMV1OwnerReferenceMock() {
	ownresources.GetPod = fakeGetPod
	ownresources.GetDeploymentRef = GetFakeDeploymentRef
	ownresources.GetManageObject = func() runtime.Object { return clusterExtension.DeepCopy() }
	ownresources.Init = fakeInit
}

func NoOLMOwnerReferenceMock() {
	ownresources.GetPod = fakeGetPod
	ownresources.GetDeploymentRef = GetFakeDeploymentRef
	ownresources.GetManageObject = func() runtime.Object { return nil }
	ownresources.Init = fakeInit
}
