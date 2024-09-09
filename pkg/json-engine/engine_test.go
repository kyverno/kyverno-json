package jsonengine

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"gotest.tools/assert"
)

var (
	policyRaw = `
{
  "apiVersion": "json.kyverno.io/v1alpha1",
  "kind": "ValidatingPolicy",
  "metadata": {
    "name": "test"
  },
  "spec": {
    "rules": [
      {
        "name": "pod-no-latest",
        "context": [
          {
            "name": "tag",
            "variable": ":latest"
          }
        ],
        "match": {
          "any": [
            {
              "apiVersion": "v1",
              "kind": "Pod"
            }
          ]
        },
        "identifier": "metadata.name",
        "assert": {
          "all": [
            {
              "check": {
                "spec": {
                  "~foo.containers->foos": {
                    "(at($foos, $foo).image)->foo": {
                      "(contains($foo, ':'))": true,
                      "(ends_with($foo, $tag))": false
                    }
                  }
                }
              }
            },
            {
              "check": {
                "spec": {
                  "~.containers->foo": {
                    "image": {
                      "(contains(@, ':'))": true,
                      "(ends_with(@, ':latest'))": false
                    }
                  }
                }
              }
            },
            {
              "check": {
                "~index.(spec.containers[*].image)->images": {
                  "(contains(@, ':'))": true,
                  "(ends_with(@, ':latest'))": false
                }
              }
            }
          ]
        }
      }
    ]
  }
}
`

	payloadRaw = `
{
	"apiVersion": "v1",
	"kind": "Pod",
	"metadata": {
		"name": "webserver"
	},
	"spec": {
		"containers": [
			{
				"name": "webserver-3",
				"image": "nginx:latest",
				"ports": [
					{
						"containerPort": 80
					}
				]
			}
		]
	}
}
`
)

func Test_JSONEngine(t *testing.T) {
	var res interface{}
	err := json.Unmarshal([]byte(payloadRaw), &res)
	assert.NilError(t, err)

	var pol v1alpha1.ValidatingPolicy
	err = json.Unmarshal([]byte(policyRaw), &pol)
	assert.NilError(t, err)

	e := New()
	resp := e.Run(context.Background(), Request{
		Resource: res,
		Policies: []*v1alpha1.ValidatingPolicy{
			&pol,
		},
	})
	assert.Equal(t, len(resp.Policies[0].Rules[0].Violations), 3)
}
