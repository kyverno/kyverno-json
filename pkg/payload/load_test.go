package payload

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	basePath := "../../test/commands/scan"
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Valid JSON File",
			path:    filepath.Join(basePath, "/tf-ec2/payload.json"),
			wantErr: false,
		},
		{
			name:    "Valid YAML File",
			path:    filepath.Join(basePath, "/pod-no-latest/payload.yaml"),
			wantErr: false,
		},
		{
			name:    "Not a YAML or a JSON File",
			path:    filepath.Join(basePath, "/dockerfile/out.txt"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(tt.path)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
