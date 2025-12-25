package perses

import (
	"embed"
	"io/fs"
	"maps"
	"path"
	"slices"

	"github.com/go-logr/logr"
	persesv1alpha1 "github.com/rhobs/perses-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	hcoutil "github.com/kubevirt/hyperconverged-cluster-operator/pkg/util"
)

// PersesDashboardsFS contains the embedded Perses dashboard YAML files.
//
//go:embed resources/dashboards/*.yaml
var persesDashboardsFS embed.FS

// PersesDatasourcesBytes contains the embedded Perses datasource YAML file.
//
//go:embed resources/datasources/00-datasource-default.yaml
var persesDatasourcesBytes []byte

var (
	datasourceName    string
	managedDashboards []string
)

func initDashboards(namespace string, logger logr.Logger) (map[string]persesv1alpha1.PersesDashboard, error) {
	dashboards, err := parseDashboards(persesDashboardsFS, namespace, logger)
	if err != nil {
		return nil, err
	}

	managedDashboards = slices.Collect(maps.Keys(dashboards))

	return dashboards, nil
}

func initDatasource(namespace string) (*persesv1alpha1.PersesDatasource, error) {
	ds := &persesv1alpha1.PersesDatasource{}
	err := yaml.Unmarshal(persesDatasourcesBytes, ds)
	if err != nil {
		return nil, err
	}

	ds.Namespace = namespace

	datasourceName = ds.Name

	return ds, nil
}

func parseDashboards(root fs.FS, namespace string, logger logr.Logger) (map[string]persesv1alpha1.PersesDashboard, error) {
	dashboards := make(map[string]persesv1alpha1.PersesDashboard)

	err := fs.WalkDir(root, ".", func(filePath string, d fs.DirEntry, internalErr error) error {
		if internalErr != nil {
			return internalErr
		}

		if d.IsDir() || path.Ext(d.Name()) != ".yaml" {
			return nil
		}

		file, err := root.Open(filePath)
		if err != nil {
			logger.Error(internalErr, "Can't open the perses dashboard yaml file", "file name", filePath)
			return err
		}
		defer file.Close()

		dec := yaml.NewYAMLToJSONDecoder(file)

		dashboard := persesv1alpha1.PersesDashboard{}
		err = dec.Decode(&dashboard)
		if err != nil {
			return err
		}

		dashboard.Namespace = namespace

		dashboards[dashboard.Name] = dashboard
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dashboards, nil
}

var (
	dashboardPredicate = predicate.NewTypedPredicateFuncs[*persesv1alpha1.PersesDashboard](func(d *persesv1alpha1.PersesDashboard) bool {
		// Only reconcile dashboards in operator namespace and that are in the managed list
		if d.Namespace != hcoutil.GetOperatorNamespaceFromEnv() {
			return false
		}
		return slices.Contains(managedDashboards, d.Name)
	})

	datasourcePredicate = predicate.NewTypedPredicateFuncs[*persesv1alpha1.PersesDatasource](func(ds *persesv1alpha1.PersesDatasource) bool {
		// Only reconcile the single datasource we manage, in operator namespace
		return ds.Namespace == hcoutil.GetOperatorNamespaceFromEnv() && ds.Name == datasourceName
	})
)
