module github.com/kubevirt/hyperconverged-cluster-operator

go 1.21

toolchain go1.21.7

require (
	dario.cat/mergo v1.0.0
	github.com/blang/semver/v4 v4.0.0
	github.com/evanphx/json-patch/v5 v5.9.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-logr/logr v1.4.1
	github.com/google/uuid v1.6.0
	github.com/kubevirt/cluster-network-addons-operator v0.91.0
	github.com/kubevirt/monitoring/pkg/metrics/parser v0.0.0-20240125201600-b689e9c89409
	github.com/machadovilaca/operator-observability v0.0.11
	github.com/onsi/ginkgo/v2 v2.15.0
	github.com/onsi/gomega v1.31.1
	github.com/openshift/api v3.9.1-0.20190517100836-d5b34b957e91+incompatible
	github.com/openshift/custom-resource-status v1.1.2
	github.com/openshift/library-go v0.0.0-20240126142105-b5320774c6c4
	github.com/operator-framework/api v0.21.0
	github.com/operator-framework/operator-lib v0.12.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.71.2
	github.com/prometheus/client_model v0.5.0
	github.com/samber/lo v1.39.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/sync v0.6.0
	golang.org/x/tools v0.17.0
	gomodules.xyz/jsonpatch/v2 v2.4.0
	k8s.io/api v0.29.1
	k8s.io/apiextensions-apiserver v0.29.1
	k8s.io/apimachinery v0.29.1
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20240126223410-2919ad4fcfec
	k8s.io/utils v0.0.0-20240102154912-e7106e64919e
	kubevirt.io/api v1.2.0-beta.0
	kubevirt.io/application-aware-quota v1.1.4
	kubevirt.io/containerized-data-importer-api v1.58.1
	kubevirt.io/controller-lifecycle-operator-sdk/api v0.2.4
	kubevirt.io/managed-tenant-quota v1.2.0
	kubevirt.io/ssp-operator/api v0.19.0
	sigs.k8s.io/controller-runtime v0.17.0
	sigs.k8s.io/controller-tools v0.14.0
)

// TODO: consume v0.12.0 as soon as available
replace github.com/operator-framework/operator-lib => github.com/operator-framework/operator-lib v0.0.0-20230717184314-6efbe3a22f6f

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/emicklei/go-restful/v3 v3.11.2 // indirect
	github.com/evanphx/json-patch v5.9.0+incompatible // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-logr/zapr v1.3.0 // indirect
	github.com/go-openapi/jsonpointer v0.20.2 // indirect
	github.com/go-openapi/jsonreference v0.20.4 // indirect
	github.com/go-openapi/swag v0.22.9 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/gobuffalo/flect v1.0.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20240125082051-42cd04596328 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.18.0 // indirect
	github.com/prometheus/common v0.46.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/exp v0.0.0-20240119083558-1b970713d09a // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/oauth2 v0.16.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/term v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiserver v0.29.1 // indirect
	k8s.io/component-base v0.29.1 // indirect
	k8s.io/klog/v2 v2.120.1 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

exclude k8s.io/cluster-bootstrap v0.0.0

exclude k8s.io/api v0.0.0

exclude k8s.io/apiextensions-apiserver v0.0.0

exclude k8s.io/apimachinery v0.0.0

exclude k8s.io/apiserver v0.0.0

exclude k8s.io/code-generator v0.0.0

exclude k8s.io/component-base v0.0.0

exclude k8s.io/kube-aggregator v0.0.0

exclude k8s.io/cli-runtime v0.0.0

exclude k8s.io/kubectl v0.0.0

exclude k8s.io/client-go v2.0.0-alpha.0.0.20181121191925-a47917edff34+incompatible

exclude k8s.io/client-go v0.0.0

exclude k8s.io/cloud-provider v0.0.0

exclude k8s.io/cri-api v0.0.0

exclude k8s.io/csi-translation-lib v0.0.0

exclude k8s.io/kube-controller-manager v0.0.0

exclude k8s.io/kube-proxy v0.0.0

exclude k8s.io/kube-scheduler v0.0.0

exclude k8s.io/kubelet v0.0.0

exclude k8s.io/legacy-cloud-providers v0.0.0

exclude k8s.io/metrics v0.0.0

exclude k8s.io/sample-apiserver v0.0.0

// Pinned to v0.29.1
replace (
	k8s.io/api => k8s.io/api v0.29.1
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.29.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.29.1
	k8s.io/apiserver => k8s.io/apiserver v0.29.1
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.29.1
	k8s.io/client-go => k8s.io/client-go v0.29.1
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.29.1
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.29.1
	k8s.io/code-generator => k8s.io/code-generator v0.29.1
	k8s.io/component-base => k8s.io/component-base v0.29.1
	k8s.io/cri-api => k8s.io/cri-api v0.29.1
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.29.1
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.29.1
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.29.1
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.29.1
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.29.1
	k8s.io/kubectl => k8s.io/kubectl v0.29.1
	k8s.io/kubelet => k8s.io/kubelet v0.29.1
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.29.1
	k8s.io/metrics => k8s.io/metrics v0.29.1
	k8s.io/node-api => k8s.io/node-api v0.29.1
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.29.1
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.29.1
	k8s.io/sample-controller => k8s.io/sample-controller v0.29.1
)

replace (
	github.com/appscode/jsonpatch => github.com/appscode/jsonpatch v1.0.1
	github.com/go-kit/kit => github.com/go-kit/kit v0.12.0
	github.com/openshift/machine-api-operator => github.com/openshift/machine-api-operator v0.2.1-0.20230329185430-d3973b45c2b6
)

replace sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v1.0.2

replace vbom.ml/util => github.com/fvbommel/util v0.0.0-20180919145318-efcd4e0f9787

replace bitbucket.org/ww/goautoneg => github.com/munnerz/goautoneg v0.0.0-20120707110453-a547fc61f48d

replace github.com/openshift/api => github.com/openshift/api v0.0.0-20230503133300-8bbcb7ca7183

// Fixes various security issues forcing newer versions of affected dependencies,
// prune the list once not explicitly required
replace (
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.5.0
)

// FIX: Unhandled exception in gopkg.in/yaml.v3
replace gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1

// FIX: CVE-2023-44487
replace golang.org/x/net => golang.org/x/net v0.17.0
