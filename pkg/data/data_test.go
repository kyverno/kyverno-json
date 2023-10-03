package data

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrds(t *testing.T) {
	data := Crds()
	{
		file, err := fs.Stat(data, "crds/json.kyverno.io_policies.yaml")
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
