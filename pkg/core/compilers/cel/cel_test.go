package cel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_compiler_Compile(t *testing.T) {
	c := NewCompiler(DefaultEnv)
	_, err := c.Compile("object.?spec")
	assert.NoError(t, err)
}
