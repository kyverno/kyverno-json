package policy

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/data"
	fileinfo "github.com/kyverno/kyverno/ext/file-info"
	"github.com/kyverno/kyverno/ext/resource/convert"
	"github.com/kyverno/kyverno/ext/resource/loader"
	yamlutils "github.com/kyverno/kyverno/ext/yaml"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

var (
	gv_v1alpha1               = schema.GroupVersion{Group: "json.kyverno.io", Version: "v1alpha1"}
	validatingPolicy_v1alpha1 = gv_v1alpha1.WithKind("ValidatingPolicy")
)

func Load(path ...string) ([]*v1alpha1.ValidatingPolicy, error) {
	var policies []*v1alpha1.ValidatingPolicy
	for _, path := range path {
		p, err := load(path)
		if err != nil {
			return nil, err
		}
		policies = append(policies, p...)
	}
	return policies, nil
}

func load(path string) ([]*v1alpha1.ValidatingPolicy, error) {
	var files []string
	err := filepath.Walk(path, func(file string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fileinfo.IsYaml(info) {
			files = append(files, file)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var policies []*v1alpha1.ValidatingPolicy
	for _, path := range files {
		content, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return nil, err
		}
		p, err := Parse(content)
		if err != nil {
			return nil, err
		}
		policies = append(policies, p...)
	}
	return policies, nil
}

func Parse(content []byte) ([]*v1alpha1.ValidatingPolicy, error) {
	documents, err := yamlutils.SplitDocuments(content)
	if err != nil {
		return nil, err
	}
	var policies []*v1alpha1.ValidatingPolicy
	// TODO: no need to allocate a validator every time
	loader, err := loader.New(openapiclient.NewLocalCRDFiles(data.Crds(), data.CrdsFolder))
	if err != nil {
		return nil, err
	}
	for _, document := range documents {
		gvk, untyped, err := loader.Load(document)
		if err != nil {
			return nil, err
		}
		switch gvk {
		case validatingPolicy_v1alpha1:
			policy, err := convert.To[v1alpha1.ValidatingPolicy](untyped)
			if err != nil {
				return nil, err
			}
			policies = append(policies, policy)
		default:
			return nil, fmt.Errorf("policy type not supported %s", gvk)
		}
	}
	return policies, nil
}
