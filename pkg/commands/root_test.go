package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TfPlan(t *testing.T) {
	cmd := NewRootCommand()
	assert.NotNil(t, cmd)
	cmd.SetArgs([]string{
		"--payload",
		"../../testdata/tf-plan/tf.plan.json",
		"--pre-process",
		"planned_values.root_module.resources",
		"--policy",
		"../../testdata/tf-plan/policy.yaml",
	})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func Test_PayloadYaml(t *testing.T) {
	cmd := NewRootCommand()
	assert.NotNil(t, cmd)
	cmd.SetArgs([]string{
		"--payload",
		"../../testdata/payload-yaml/payload.yaml",
		"--pre-process",
		"planned_values.root_module.resources",
		"--policy",
		"../../testdata/payload-yaml/policy.yaml",
	})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func Test_FooBar(t *testing.T) {
	cmd := NewRootCommand()
	assert.NotNil(t, cmd)
	cmd.SetArgs([]string{
		"--payload",
		"../../testdata/foo-bar/payload.yaml",
		"--policy",
		"../../testdata/foo-bar/policy.yaml",
	})
	err := cmd.Execute()
	assert.NoError(t, err)
}
