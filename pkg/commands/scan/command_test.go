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
		payload:  "../../../test/commands/scan/foo-bar/payload.yaml",
		policies: []string{"../../../test/commands/scan/foo-bar/policy.yaml"},
		out:      "../../../test/commands/scan/foo-bar/out.txt",
		wantErr:  false,
	}, {
		name:     "wildcard",
		payload:  "../../../test/commands/scan/wildcard/payload.json",
		policies: []string{"../../../test/commands/scan/wildcard/policy.yaml"},
		out:      "../../../test/commands/scan/wildcard/out.txt",
		wantErr:  false,
	}, {
		name:     "pod-no-latest",
		payload:  "../../../test/commands/scan/pod-no-latest/payload.yaml",
		policies: []string{"../../../test/commands/scan/pod-no-latest/policy.yaml"},
		out:      "../../../test/commands/scan/pod-no-latest/out.txt",
		wantErr:  false,
	}, {
		name:     "pod-all-latest",
		payload:  "../../../test/commands/scan/pod-all-latest/payload.yaml",
		policies: []string{"../../../test/commands/scan/pod-all-latest/policy.yaml"},
		out:      "../../../test/commands/scan/pod-all-latest/out.txt",
		wantErr:  false,
	}, {
		name:     "scripted",
		payload:  "../../../test/commands/scan/scripted/payload.yaml",
		policies: []string{"../../../test/commands/scan/scripted/policy.yaml"},
		out:      "../../../test/commands/scan/scripted/out.txt",
		wantErr:  false,
	}, {
		name:          "payload-yaml",
		payload:       "../../../test/commands/scan/payload-yaml/payload.yaml",
		preprocessors: []string{"planned_values.root_module.resources"},
		policies:      []string{"../../../test/commands/scan/payload-yaml/policy.yaml"},
		out:           "../../../test/commands/scan/payload-yaml/out.txt",
		wantErr:       false,
	}, {
		name:          "tf-plan",
		payload:       "../../../test/commands/scan/tf-plan/payload.json",
		preprocessors: []string{"planned_values.root_module.resources"},
		policies:      []string{"../../../test/commands/scan/tf-plan/policy.yaml"},
		out:           "../../../test/commands/scan/tf-plan/out.txt",
		wantErr:       false,
	}, {
		name:     "escaped",
		payload:  "../../../test/commands/scan/escaped/payload.yaml",
		policies: []string{"../../../test/commands/scan/escaped/policy.yaml"},
		out:      "../../../test/commands/scan/escaped/out.txt",
		wantErr:  false,
	}, {
		name:     "dockerfile",
		payload:  "../../../test/commands/scan/dockerfile/payload.json",
		policies: []string{"../../../test/commands/scan/dockerfile/policy.yaml"},
		out:      "../../../test/commands/scan/dockerfile/out.txt",
		wantErr:  false,
	}, {
		name:     "tf-s3",
		payload:  "../../../test/commands/scan/tf-s3/payload.json",
		policies: []string{"../../../test/commands/scan/tf-s3/policy.yaml"},
		out:      "../../../test/commands/scan/tf-s3/out.txt",
		wantErr:  false,
	}, {
		name:          "tf-ec2",
		payload:       "../../../test/commands/scan/tf-ec2/payload.json",
		preprocessors: []string{"planned_values.root_module.resources"},
		policies:      []string{"../../../test/commands/scan/tf-ec2/policy.yaml"},
		out:           "../../../test/commands/scan/tf-ec2/out.txt",
		wantErr:       false,
	}, {
		name:          "tf-ecs-cluster-1",
		payload:       "../../../test/commands/scan/tf-ecs-cluster/payload.json",
		preprocessors: []string{"planned_values.root_module.resources"},
		policies:      []string{"../../../test/commands/scan/tf-ecs-cluster/01-policy.yaml"},
		out:           "../../../test/commands/scan/tf-ecs-cluster/01-out.txt",
		wantErr:       false,
	}, {
		name:          "tf-ecs-cluster-2",
		payload:       "../../../test/commands/scan/tf-ecs-cluster/payload.json",
		preprocessors: []string{"planned_values.root_module.resources"},
		policies:      []string{"../../../test/commands/scan/tf-ecs-cluster/02-policy.yaml"},
		out:           "../../../test/commands/scan/tf-ecs-cluster/02-out.txt",
		wantErr:       false,
	}, {
		name:          "tf-ecs-service-1",
		payload:       "../../../test/commands/scan/tf-ecs-service/payload.json",
		preprocessors: []string{"planned_values.root_module.resources"},
		policies:      []string{"../../../test/commands/scan/tf-ecs-service/01-policy.yaml"},
		out:           "../../../test/commands/scan/tf-ecs-service/01-out.txt",
		wantErr:       false,
	}, {
		name:          "tf-ecs-service-2",
		payload:       "../../../test/commands/scan/tf-ecs-service/payload.json",
		preprocessors: []string{"planned_values.root_module.resources"},
		policies:      []string{"../../../test/commands/scan/tf-ecs-service/02-policy.yaml"},
		out:           "../../../test/commands/scan/tf-ecs-service/02-out.txt",
		wantErr:       false,
	}, {
		name:          "tf-ecs-task-definition",
		payload:       "../../../test/commands/scan/tf-ecs-task-definition/payload.json",
		preprocessors: []string{"planned_values.root_module.resources"},
		policies:      []string{"../../../test/commands/scan/tf-ecs-task-definition/policy.yaml"},
		out:           "../../../test/commands/scan/tf-ecs-task-definition/out.txt",
		wantErr:       false,
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
