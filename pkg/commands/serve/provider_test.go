package serve

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadPolicy(t *testing.T) {
	testFile := filepath.Join("testdata", "dockerfile-policy.yaml")

	policies, err := LoadPoliciesFromDir(testFile)

	assert.NoError(t, err)
	assert.NotNil(t, policies)

	for _, p := range policies {
		assert.NotEmpty(t, p.Name, "Policy name should not be empty")
		assert.NotEmpty(t, p.Spec.Rules, "Policy rules should not be empty")
	}

}
