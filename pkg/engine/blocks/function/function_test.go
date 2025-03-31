package function

import (
	"context"
	"strconv"
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

func Test_Function(t *testing.T) {
	// Define a function that transforms a request into a response
	processingFunc := func(ctx context.Context, request TestRequest) TestResponse {
		return TestResponse{
			Output: "Processed: " + request.Input,
		}
	}

	// Create the function engine
	engine := New[TestRequest, TestResponse](processingFunc)

	// Test cases
	testCases := []struct {
		name     string
		request  TestRequest
		expected TestResponse
	}{
		{
			name:     "empty input",
			request:  TestRequest{Input: ""},
			expected: TestResponse{Output: "Processed: "},
		},
		{
			name:     "non-empty input",
			request:  TestRequest{Input: "hello"},
			expected: TestResponse{Output: "Processed: hello"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Run with context and request
			response := engine.Run(context.Background(), tc.request)

			// Verify the response
			assert.Equal(t, tc.expected.Output, response.Output)
		})
	}
}

// Test with primitive types
func Test_FunctionWithPrimitives(t *testing.T) {
	// Test with int to string conversion
	intToStringFunc := func(ctx context.Context, n int) string {
		return "Number: " + strconv.Itoa(n)
	}
	intToStringEngine := New[int, string](intToStringFunc)
	result := intToStringEngine.Run(context.Background(), 42)
	assert.Equal(t, "Number: 42", result)

	// Test with string to int conversion
	stringToIntFunc := func(ctx context.Context, s string) int {
		val, _ := strconv.Atoi(s)
		return val * 2
	}
	stringToIntEngine := New[string, int](stringToIntFunc)
	resultInt := stringToIntEngine.Run(context.Background(), "21")
	assert.Equal(t, 42, resultInt)

	// Test with context usage
	contextAwareFunc := func(ctx context.Context, s string) bool {
		value := ctx.Value("key")
		return value == s
	}

	contextAwareEngine := New[string, bool](contextAwareFunc)

	// Create context with value
	ctxWithValue := context.WithValue(context.Background(), "key", "match")

	// Should match
	resultTrue := contextAwareEngine.Run(ctxWithValue, "match")
	assert.Equal(t, true, resultTrue)

	// Should not match
	resultFalse := contextAwareEngine.Run(ctxWithValue, "no match")
	assert.Equal(t, false, resultFalse)
}

// Test with complex data transformations
func Test_FunctionWithComplexTransformations(t *testing.T) {
	// Map transformation function
	mapTransformFunc := func(ctx context.Context, input map[string]interface{}) map[string]interface{} {
		result := make(map[string]interface{})
		for k, v := range input {
			// Transform each value by appending "-transformed"
			if strVal, ok := v.(string); ok {
				result[k] = strVal + "-transformed"
			} else {
				result[k] = v
			}
		}
		return result
	}

	mapEngine := New[map[string]interface{}, map[string]interface{}](mapTransformFunc)

	inputMap := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": 123,
	}

	resultMap := mapEngine.Run(context.Background(), inputMap)

	assert.Equal(t, "value1-transformed", resultMap["key1"].(string))
	assert.Equal(t, "value2-transformed", resultMap["key2"].(string))
	assert.Equal(t, 123, resultMap["key3"].(int))
}
