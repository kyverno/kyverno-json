package scan

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Execute(t *testing.T) {
	tests := []struct {
		name          string
		payload       string
		preprocessors []string
		policies      []string
		wantErr       bool
		out           string
	}{{
		name:     "foo-bar",
		payload:  "../../../test/foo-bar/payload.yaml",
		policies: []string{"../../../test/foo-bar/policy.yaml"},
		out:      "../../../test/foo-bar/out.txt",
		wantErr:  false,
	}, {
		name:     "wildcard",
		payload:  "../../../test/wildcard/payload.json",
		policies: []string{"../../../test/wildcard/policy.yaml"},
		out:      "../../../test/wildcard/out.txt",
		wantErr:  false,
	}, {
		name:     "pod-no-latest",
		payload:  "../../../test/pod-no-latest/payload.yaml",
		policies: []string{"../../../test/pod-no-latest/policy.yaml"},
		out:      "../../../test/pod-no-latest/out.txt",
		wantErr:  false,
	}, {
		name:     "pod-all-latest",
		payload:  "../../../test/pod-all-latest/payload.yaml",
		policies: []string{"../../../test/pod-all-latest/policy.yaml"},
		out:      "../../../test/pod-all-latest/out.txt",
		wantErr:  false,
	}, {
		name:     "scripted",
		payload:  "../../../test/scripted/payload.yaml",
		policies: []string{"../../../test/scripted/policy.yaml"},
		out:      "../../../test/scripted/out.txt",
		wantErr:  false,
	}, {
		name:          "payload-yaml",
		payload:       "../../../test/payload-yaml/payload.yaml",
		preprocessors: []string{"planned_values.root_module.resources"},
		policies:      []string{"../../../test/payload-yaml/policy.yaml"},
		out:           "../../../test/payload-yaml/out.txt",
		wantErr:       false,
	}, {
		name:          "tf-plan",
		payload:       "../../../test/tf-plan/payload.json",
		preprocessors: []string{"planned_values.root_module.resources"},
		policies:      []string{"../../../test/tf-plan/policy.yaml"},
		out:           "../../../test/tf-plan/out.txt",
		wantErr:       false,
	}, {
		name:     "escaped",
		payload:  "../../../test/escaped/payload.yaml",
		policies: []string{"../../../test/escaped/policy.yaml"},
		out:      "../../../test/escaped/out.txt",
		wantErr:  false,
	}, {
		name:     "dockerfile",
		payload:  "../../../test/dockerfile/payload.json",
		policies: []string{"../../../test/dockerfile/policy.yaml"},
		out:      "../../../test/dockerfile/out.txt",
		wantErr:  false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Command()
			assert.NotNil(t, cmd)
			var args []string
			args = append(args, "--payload", tt.payload)
			for _, preprocessor := range tt.preprocessors {
				args = append(args, "--pre-process", preprocessor)
			}
			for _, policy := range tt.policies {
				args = append(args, "--policy", policy)
			}
			args = append(args, "--payload", tt.payload)
			cmd.SetArgs(args)
			out := bytes.NewBufferString("")
			cmd.SetOut(out)
			if err := cmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("command.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			actual, err := io.ReadAll(out)
			assert.NoError(t, err)
			if tt.out != "" {
				expected, err := os.ReadFile(tt.out)
				assert.NoError(t, err)
				assert.Equal(t, string(expected), string(actual))
			}
		})
	}
}
