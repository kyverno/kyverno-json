package policy

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/data"
	fileinfo "github.com/kyverno/kyverno-json/pkg/utils/file-info"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource/loader"
	yamlutils "github.com/kyverno/kyverno/pkg/utils/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

var (
	gv_v1alpha1     = schema.GroupVersion{Group: "json.kyverno.io", Version: "v1alpha1"}
	policy_v1alpha1 = gv_v1alpha1.WithKind("Policy")
)

func Load(path ...string) ([]*v1alpha1.Policy, error) {
	var policies []*v1alpha1.Policy
	for _, path := range path {
		p, err := load(path)
		if err != nil {
			return nil, err
		}
		policies = append(policies, p...)
	}
	return policies, nil
}

func load(path string) ([]*v1alpha1.Policy, error) {
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
	var policies []*v1alpha1.Policy
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

func Parse(content []byte) ([]*v1alpha1.Policy, error) {
	documents, err := yamlutils.SplitDocuments(content)
	if err != nil {
		return nil, err
	}
	var policies []*v1alpha1.Policy
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
		case policy_v1alpha1:
			// TODO: don't use kyverno's convert for now to workaround the bug in api machinery code
			// https://kubernetes.slack.com/archives/C0EG7JC6T/p1696331287543159
			policy, err := To[v1alpha1.Policy](untyped)
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

func Into[T any](untyped unstructured.Unstructured, result *T) error {
	return runtime.DefaultUnstructuredConverter.FromUnstructured(untyped.UnstructuredContent(), result)
}

func To[T any](untyped unstructured.Unstructured) (*T, error) {
	var result T
	if err := Into(untyped, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
