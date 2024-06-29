package binding

import (
	"encoding/json"
	"testing"

	jpbinding "github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"gotest.tools/assert"
)

var (
	context = `
	[
		{
			"name": "tag",
			"variable": "latest"
		},
		{
			"name": "tag",
			"variable": "(concat(':', $tag))"
		}
	]
	`

	contextFromResource = `
	[
		{
			"name": "containerName",
			"variable": "(spec.containers[*].name | [0])"
		}
	]
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

func Test_NewContextBindings(t *testing.T) {
	var ctxEntry []v1alpha1.ContextEntry
	err := json.Unmarshal([]byte(context), &ctxEntry)
	assert.NilError(t, err)

	bindings := jpbinding.NewBindings()
	bindings = NewContextBindings(bindings, nil, ctxEntry...)

	b, err := bindings.Get("$tag")
	assert.NilError(t, err)
	val, err := b.Value()
	assert.NilError(t, err)
	assert.Equal(t, val.(string), ":latest")

	var res interface{}
	err = json.Unmarshal([]byte(resource), &res)
	assert.NilError(t, err)
	err = json.Unmarshal([]byte(contextFromResource), &ctxEntry)
	assert.NilError(t, err)

	bindings = NewContextBindings(bindings, res, ctxEntry...)

	b, err = bindings.Get("$containerName")
	assert.NilError(t, err)
	containerName, err := b.Value()
	assert.NilError(t, err)
	assert.Equal(t, containerName.(string), "webserver-3")
}
