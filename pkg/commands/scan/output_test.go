package scan

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
	"github.com/kyverno/kyverno-json/pkg/matching"
	"gotest.tools/assert"
)

var resp = jsonengine.Response{
	Policies: []jsonengine.PolicyResponse{
		{
			Rules: []jsonengine.RuleResponse{
				{
					Identifier: "test-identifier",
					Error:      errors.New("test-error"),
					Violations: matching.Results{
						{
							Message: "test-message",
						},
					},
				},
			},
		},
	},
}

func Test_OutputJSON(t *testing.T) {
	var buff bytes.Buffer
	out := newOutput(&buff, "json")

	out.responses(resp)
	output := buff.String()
	assert.Assert(t, strings.Contains(output, "test-error"))
	assert.Assert(t, strings.Contains(output, "test-message"))
	assert.Assert(t, strings.Contains(output, "test-identifier"))
}

func Test_OutputText(t *testing.T) {
	var buff bytes.Buffer
	out := newOutput(&buff, "text")

	out.println(resp)
	output := buff.String()
	assert.Assert(t, strings.Contains(output, "test-error"))
	assert.Assert(t, strings.Contains(output, "test-message"))
	assert.Assert(t, strings.Contains(output, "test-identifier"))
}
