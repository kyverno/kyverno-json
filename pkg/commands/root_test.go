package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	cmd := NewRootCommand()
	assert.NotNil(t, cmd)
	cmd.SetArgs([]string{
		"--payload",
		"../../tf.plan.json",
		"--pre-process",
		"planned_values.root_module.resources",
		"--policy",
		"../../policy.yaml",
	})
	err := cmd.Execute()
	assert.NoError(t, err)
}
