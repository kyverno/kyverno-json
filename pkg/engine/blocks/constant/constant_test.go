package constant

import (
	"context"
	"testing"

	"gotest.tools/assert"
)

// Test types for generic engine
type TestRequest struct {
	Input string
}

type TestResponse struct {
	Output string
}

func Test_Constant(t *testing.T) {
	// Test creating a constant engine with a predefined response
	expectedResponse := TestResponse{Output: "test output"}
	engine := New[TestRequest, TestResponse](expectedResponse)

	// Test that Run returns the constant response regardless of the input
	testCases := []struct {
		name    string
		request TestRequest
	}{
		{
			name:    "empty request",
			request: TestRequest{},
		},
		{
			name:    "non-empty request",
			request: TestRequest{Input: "test input"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Run with different requests should always return the same response
			response := engine.Run(context.Background(), tc.request)

			// Verify the response matches the expected constant value
			assert.Equal(t, expectedResponse.Output, response.Output)
		})
	}
}

// Test with primitive types
func Test_ConstantWithPrimitives(t *testing.T) {
	// Test with string
	stringEngine := New[int, string]("constant string")
	stringResponse := stringEngine.Run(context.Background(), 42)
	assert.Equal(t, "constant string", stringResponse)

	// Test with int
	intEngine := New[string, int](100)
	intResponse := intEngine.Run(context.Background(), "any input")
	assert.Equal(t, 100, intResponse)

	// Test with bool
	boolEngine := New[float64, bool](true)
	boolResponse := boolEngine.Run(context.Background(), 3.14)
	assert.Equal(t, true, boolResponse)
}

// Test with nil value
func Test_ConstantWithNil(t *testing.T) {
	// Create engine with nil map
	var nilMap map[string]interface{}
	mapEngine := New[string, map[string]interface{}](nilMap)
	mapResponse := mapEngine.Run(context.Background(), "any input")
	assert.Assert(t, mapResponse == nil)
}
