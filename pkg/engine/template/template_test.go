package template

import (
	"context"
	"encoding/json"
	"testing"

	jpbinding "github.com/jmespath-community/go-jmespath/pkg/binding"
	"gotest.tools/assert"
)

var (
	resourceRaw = `
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
				},
				{
					"name": "webserver-4",
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

func Test_Execute(t *testing.T) {
	bindings := jpbinding.NewBindings()
	bindings = bindings.Register("$allowedContainerNames", jpbinding.NewBinding(`["webserver-1", "webserver-2", "webserver-3"]`))
	bindings = bindings.Register("$tag", jpbinding.NewBinding(`:latest`))

	var payload map[string]interface{}
	err := json.Unmarshal([]byte(resourceRaw), &payload)
	assert.NilError(t, err)

	val, err := Execute(context.Background(), "contains($allowedContainerNames, spec.containers[*].name | [0])", payload, bindings)
	assert.NilError(t, err)
	assert.Equal(t, val.(bool), true)

	val, err = Execute(context.Background(), "contains($allowedContainerNames, 'bad')", payload, bindings)
	assert.NilError(t, err)
	assert.Equal(t, val.(bool), false)

	val, err = Execute(context.Background(), "ends_with(spec.containers[*].image | [1], $tag)", payload, bindings)
	assert.NilError(t, err)
	assert.Equal(t, val.(bool), true)
}
