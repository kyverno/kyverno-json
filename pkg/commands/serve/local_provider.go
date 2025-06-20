package serve

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"sigs.k8s.io/yaml"
)


type localProvider struct {
	policies []v1alpha1.ValidatingPolicy
}

func (p *localProvider) Get() ([]v1alpha1.ValidatingPolicy, error) {
	return p.policies, nil
}

func PolicylocalProvider(policies []v1alpha1.ValidatingPolicy) *localProvider {
	return &localProvider{
		policies: policies,
	}
}

func LoadPoliciesFromDir(dir string) ([]v1alpha1.ValidatingPolicy, error) {
	policies := []v1alpha1.ValidatingPolicy{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(strings.ToLower(path), ".yaml") || strings.HasSuffix(strings.ToLower(path), ".yml") {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", path, err)
			}

			var policy v1alpha1.ValidatingPolicy
			if err := yaml.Unmarshal(content, &policy); err != nil {
				return fmt.Errorf("error parsing YAML in file %s: %w", path, err)
			}

			policies = append(policies, policy)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", dir, err)
	}

	return policies, nil
}
