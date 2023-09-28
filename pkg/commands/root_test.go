package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	cmd := NewRootCommand()
	assert.NotNil(t, cmd)
	cmd.SetArgs([]string{
		"--plan",
		"../../tf.plan.json",
		"--policy",
		"../../policy.yaml",
	})
	err := cmd.Execute()
	assert.NoError(t, err)
}
