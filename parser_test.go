package main

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedValue interface{}
		expectError   bool
	}{
		{
			name:          "Empty Input",
			input:         "",
			expectedValue: nil,
			expectError:   true,
		},
		{
			name:          "String",
			input:         `"hello"`,
			expectedValue: "hello",
			expectError:   false,
		},
		{
			name:          "Number",
			input:         `42`,
			expectedValue: float64(42),
			expectError:   false,
		},
		{
			name:          "Boolean",
			input:         `true`,
			expectedValue: true,
			expectError:   false,
		},
		{
			name:          "Empty Object",
			input:         "{}",
			expectedValue: make(map[string]interface{}),
			expectError:   false,
		},
		{
			name:          "Simple Object with String Value",
			input:         `{"key":"value"}`,
			expectedValue: map[string]interface{}{"key": "value"},
			expectError:   false,
		},
		{
			name:          "Simple Object with Numeric, Boolean and Null Values",
			input:         `{ "keyA": "value", "keyB": 42.5, "keyC": true, "keyD": null }`,
			expectedValue: map[string]interface{}{"keyA": "value", "keyB": float64(42.5), "keyC": true, "keyD": nil},
			expectError:   false,
		},
		{
			name:          "Simple Object Containing Nested Objects",
			input:         `{"key": {"key2": "value"}, "key3": []}`,
			expectedValue: map[string]interface{}{"key": map[string]interface{}{"key2": "value"}, "key3": []interface{}{}},
			expectError:   false,
		},
		{
			name:          "Empty Array",
			input:         "[]",
			expectedValue: make([]interface{}, 0),
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			actualValue, actualError := p.Parse()
			if actualError != nil && !tt.expectError {
				t.Fatalf("Expected no error but received: %v", actualError)
			}
			if !reflect.DeepEqual(actualValue, tt.expectedValue) {
				t.Fatalf("Expected %T but received %T", tt.expectedValue, actualValue)
			}
		})
	}
}
