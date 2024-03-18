package matching

import (
	"context"
	"encoding/json"
	"testing"

	jpbinding "github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"gotest.tools/assert"
)

var (
	policy = `
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
            "variable": "latest"
          },
          {
            "name": "tag",
            "variable": "(concat(':', $tag))"
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
                "~.(spec.containers[*].image)": {
                  "(contains(@, ':'))": true,
                  "(ends_with(@, ':latest'))": true
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
	match1 = `
{
	"all": [
		{
			"apiVersion": "v1",
			"kind": "Pod"
		},
		{
			"metadata": {
				"name": "webserver"
			}
		}
	]
}
`

	match2 = `
{
	"any": [
		{
			"apiVersion": "v1",
			"kind": "Pod"
		},
		{
			"apiVersion": "v1",
			"kind": "Deployment"
		}
	]
}
`
	resource = `
{
  "apiVersion": "v1",
  "kind": "Pod",
  "metadata": {
    "name": "webserver"
  },
  "spec": {
    "containers": [
      {
        "name": "webserver-1",
        "image": "nginx:latest",
        "ports": [
          {
            "containerPort": 80
          }
        ]
      },
      {
        "name": "webserver-2",
        "image": "nginx:latest",
        "ports": [
          {
            "containerPort": 80
          }
        ]
      },
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

	resourcebad = `
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
        "image": "nginx:1.2",
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

func Test_Match(t *testing.T) {
	var res interface{}
	err := json.Unmarshal([]byte(resource), &res)
	assert.NilError(t, err)

	var m1, m2 v1alpha1.Match
	err = json.Unmarshal([]byte(match1), &m1)
	assert.NilError(t, err)

	err = json.Unmarshal([]byte(match2), &m2)
	assert.NilError(t, err)

	errs, err := Match(context.Background(), nil, &m1, res, nil)
	assert.NilError(t, err)
	assert.Equal(t, len(errs), 0)

	errs, err = Match(context.Background(), nil, &m2, res, nil)
	assert.NilError(t, err)
	assert.Equal(t, len(errs), 0)

	m2.All = m2.Any
	m2.Any = nil
	errs, err = Match(context.Background(), nil, &m2, res, nil)
	assert.NilError(t, err)
	assert.Equal(t, len(errs), 0)
}

func Test_MatchAssert(t *testing.T) {
	var res interface{}
	err := json.Unmarshal([]byte(resource), &res)
	assert.NilError(t, err)

	var pol v1alpha1.ValidatingPolicy
	err = json.Unmarshal([]byte(policy), &pol)
	assert.NilError(t, err)

	bindings := jpbinding.NewBindings().Register("$payload", jpbinding.NewBinding(res))
	bindings = bindings.Register("$policy", jpbinding.NewBinding(pol))
	bindings = bindings.Register("$rule", jpbinding.NewBinding(pol.Spec.Rules[0]))

	errs, err := MatchAssert(context.Background(), nil, pol.Spec.Rules[0].Assert, res, bindings)
	assert.NilError(t, err)
	assert.Equal(t, len(errs), 0)

	var resBad interface{}
	err = json.Unmarshal([]byte(resourcebad), &resBad)
	assert.NilError(t, err)

	errs, err = MatchAssert(context.Background(), nil, pol.Spec.Rules[0].Assert, resBad, bindings)
	assert.NilError(t, err)
	assert.Equal(t, len(errs), 1)
}
