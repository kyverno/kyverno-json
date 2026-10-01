package serve_test

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno-json/pkg/commands/serve"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CommandFlags(t *testing.T) {
	cmd := serve.Command()
	require.NotNil(t, cmd, fmt.Sprintf("%s command should not be nil", cmd.Use))

	tests := []struct {
		flagName     string
		defaultValue string
		userValue    string
	}{
		{"server-host", "0.0.0.0", "127.0.0.1"},
		{"server-port", "8080", "9090"},
		{"gin-mode", gin.ReleaseMode, gin.DebugMode},
		{"gin-log", "true", "false"},
		{"gin-cors", "true", "false"},
		{"gin-max-body-size", "2097152" /* 2MB, 2 * 1024 * 1024 */, "1048576" /* 1MB, 1 * 1024 * 1024 */},
	}

	for _, tt := range tests {
		t.Run(tt.flagName, func(t *testing.T) {
			flag := cmd.Flags().Lookup(tt.flagName)
			require.NotNil(t, flag, "Flag should be defined")
			assert.Equal(t, tt.defaultValue, flag.DefValue, "Default value should match")

			/* Set the flag to the user-defined value */
			err := cmd.Flags().Set(tt.flagName, tt.userValue)
			require.NoError(t, err, "Setting flag should not return an error")

			/* Verify the flag is set to the user-defined value */
			assert.Equal(t, tt.userValue, flag.Value.String(), "User-defined value should match")
		})
	}
}
