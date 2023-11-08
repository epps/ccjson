package main

import "testing"

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:     "Empty Object",
			input:    "{}",
			expected: []Token{{Type: BeginObject, Literal: "{"}, {Type: EndObject, Literal: "}"}, {Type: EOF, Literal: ""}},
		},
		{
			name:     "Empty Array",
			input:    "[]",
			expected: []Token{{Type: BeginArray, Literal: "["}, {Type: EndArray, Literal: "]"}, {Type: EOF, Literal: ""}},
		},
		{
			name:  "Object with String Value",
			input: `{ "key": "value" }`,
			expected: []Token{
				{Type: BeginObject, Literal: "{"},
				{Type: String, Literal: "key"},
				{Type: NameSeparator, Literal: ":"},
				{Type: String, Literal: "value"},
				{Type: EndObject, Literal: "}"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "Object with Number Value (Carriage Returns)",
			input: `{\r"key": 42\r}`,
			expected: []Token{
				{Type: BeginObject, Literal: "{"},
				{Type: String, Literal: "key"},
				{Type: NameSeparator, Literal: ":"},
				{Type: Number, Literal: "42"},
				{Type: EndObject, Literal: "}"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "Object with Boolean Value (New Line Characters)",
			input: `{\n"a": true,\n"b": false,\n"c": null\n}`,
			expected: []Token{
				{Type: BeginObject, Literal: "{"},
				{Type: String, Literal: "a"},
				{Type: NameSeparator, Literal: ":"},
				{Type: True, Literal: "true"},
				{Type: ValueSeparator, Literal: ","},
				{Type: String, Literal: "b"},
				{Type: NameSeparator, Literal: ":"},
				{Type: False, Literal: "false"},
				{Type: ValueSeparator, Literal: ","},
				{Type: String, Literal: "c"},
				{Type: NameSeparator, Literal: ":"},
				{Type: Null, Literal: "null"},
				{Type: EndObject, Literal: "}"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "Object with Boolean, Null, and Number Values (Decimal and Negative)",
			input: `{ "key1": true, "key2": false, "key3": null, "key4": "value", "key5": 101.2, "key6": -42 }`,
			expected: []Token{
				{Type: BeginObject, Literal: "{"},
				{Type: String, Literal: "key1"},
				{Type: NameSeparator, Literal: ":"},
				{Type: True, Literal: "true"},
				{Type: ValueSeparator, Literal: ","},
				{Type: String, Literal: "key2"},
				{Type: NameSeparator, Literal: ":"},
				{Type: False, Literal: "false"},
				{Type: ValueSeparator, Literal: ","},
				{Type: String, Literal: "key3"},
				{Type: NameSeparator, Literal: ":"},
				{Type: Null, Literal: "null"},
				{Type: ValueSeparator, Literal: ","},
				{Type: String, Literal: "key4"},
				{Type: NameSeparator, Literal: ":"},
				{Type: String, Literal: "value"},
				{Type: ValueSeparator, Literal: ","},
				{Type: String, Literal: "key5"},
				{Type: NameSeparator, Literal: ":"},
				{Type: Number, Literal: "101.2"},
				{Type: ValueSeparator, Literal: ","},
				{Type: String, Literal: "key6"},
				{Type: NameSeparator, Literal: ":"},
				{Type: Number, Literal: "-42"},
				{Type: EndObject, Literal: "}"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "Object with Empty Nested Object and Array Values",
			input: `{ "key": "value", "key-n": 101, "key-o": {}, "key-l": [] }`,
			expected: []Token{
				{Type: BeginObject, Literal: "{"},
				{Type: String, Literal: "key"},
				{Type: NameSeparator, Literal: ":"},
				{Type: String, Literal: "value"},
				{Type: ValueSeparator, Literal: ","},
				{Type: String, Literal: "key-n"},
				{Type: NameSeparator, Literal: ":"},
				{Type: Number, Literal: "101"},
				{Type: ValueSeparator, Literal: ","},
				{Type: String, Literal: "key-o"},
				{Type: NameSeparator, Literal: ":"},
				{Type: BeginObject, Literal: "{"},
				{Type: EndObject, Literal: "}"},
				{Type: ValueSeparator, Literal: ","},
				{Type: String, Literal: "key-l"},
				{Type: NameSeparator, Literal: ":"},
				{Type: BeginArray, Literal: "["},
				{Type: EndArray, Literal: "]"},
				{Type: EndObject, Literal: "}"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "Object with Nested Object and Array Values",
			input: `{ "key": "value", "key-n": 101, "key-o": { "inner key": "inner value" }, "key-l": [ "list value" ] }`,
			expected: []Token{
				{Type: BeginObject, Literal: "{"},
				{Type: String, Literal: "key"},
				{Type: NameSeparator, Literal: ":"},
				{Type: String, Literal: "value"},
				{Type: ValueSeparator, Literal: ","},
				{Type: String, Literal: "key-n"},
				{Type: NameSeparator, Literal: ":"},
				{Type: Number, Literal: "101"},
				{Type: ValueSeparator, Literal: ","},
				{Type: String, Literal: "key-o"},
				{Type: NameSeparator, Literal: ":"},
				{Type: BeginObject, Literal: "{"},
				{Type: String, Literal: "inner key"},
				{Type: NameSeparator, Literal: ":"},
				{Type: String, Literal: "inner value"},
				{Type: EndObject, Literal: "}"},
				{Type: ValueSeparator, Literal: ","},
				{Type: String, Literal: "key-l"},
				{Type: NameSeparator, Literal: ":"},
				{Type: BeginArray, Literal: "["},
				{Type: String, Literal: "list value"},
				{Type: EndArray, Literal: "]"},
				{Type: EndObject, Literal: "}"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:     "Empty Input",
			input:    "",
			expected: []Token{{Type: EOF, Literal: ""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			for _, expected := range tt.expected {
				actual := l.NextToken()
				if actual != expected {
					t.Errorf("expected %v, got %v", expected, actual)
				}
			}
		})
	}
}
