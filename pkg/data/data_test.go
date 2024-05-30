package data

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrds(t *testing.T) {
	data, err := Crds()
	assert.NoError(t, err)
	{
		file, err := fs.Stat(data, "crds/json.kyverno.io_validatingpolicies.yaml")
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.False(t, file.IsDir())
	}
	{
		file, err := fs.Stat(data, "crds")
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.True(t, file.IsDir())
	}
}
