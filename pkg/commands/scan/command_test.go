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
		payload:  "../../../testdata/foo-bar/payload.yaml",
		policies: []string{"../../../testdata/foo-bar/policy.yaml"},
		out:      "../../../testdata/foo-bar/out.txt",
		wantErr:  false,
	}, {
		name:     "jim",
		payload:  "../../../testdata/jim/payload.json",
		policies: []string{"../../../testdata/jim/policy.yaml"},
		out:      "../../../testdata/jim/out.txt",
		wantErr:  false,
	}, {
		name:     "pod-no-latest",
		payload:  "../../../testdata/pod-no-latest/payload.yaml",
		policies: []string{"../../../testdata/pod-no-latest/policy.yaml"},
		out:      "../../../testdata/pod-no-latest/out.txt",
		wantErr:  false,
	}, {
		name:     "pod-all-latest",
		payload:  "../../../testdata/pod-all-latest/payload.yaml",
		policies: []string{"../../../testdata/pod-all-latest/policy.yaml"},
		out:      "../../../testdata/pod-all-latest/out.txt",
		wantErr:  false,
	}, {
		name:     "scripted",
		payload:  "../../../testdata/scripted/payload.yaml",
		policies: []string{"../../../testdata/scripted/policy.yaml"},
		out:      "../../../testdata/scripted/out.txt",
		wantErr:  false,
	}, {
		name:          "payload-yaml",
		payload:       "../../../testdata/payload-yaml/payload.yaml",
		preprocessors: []string{"planned_values.root_module.resources"},
		policies:      []string{"../../../testdata/payload-yaml/policy.yaml"},
		out:           "../../../testdata/payload-yaml/out.txt",
		wantErr:       false,
	}, {
		name:          "tf-plan",
		payload:       "../../../testdata/tf-plan/tf.plan.json",
		preprocessors: []string{"planned_values.root_module.resources"},
		policies:      []string{"../../../testdata/tf-plan/policy.yaml"},
		out:           "../../../testdata/tf-plan/out.txt",
		wantErr:       false,
	}, {
		name:     "escaped",
		payload:  "../../../testdata/escaped/payload.yaml",
		policies: []string{"../../../testdata/escaped/policy.yaml"},
		out:      "../../../testdata/escaped/out.txt",
		wantErr:  false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Command(nil)
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
