package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_GetValue(t *testing.T) {
	tests := []struct {
		name     string
		command  Command
		expected string
	}{
		{
			name:     "Given command with single equals sign, When GetValue is called, Then it should return the value after the equals sign",
			command:  Command("/stock=AAPL"),
			expected: "AAPL",
		},
		{
			name:     "Given a command with multiple equals signs, When GetValue is called, Then it should return the value after the first equals sign",
			command:  Command("/stock=AAPL=test"),
			expected: "AAPL=test",
		},
		{
			name:     "Given a command without an equals sign, When GetValue is called, Then it should return an empty string",
			command:  Command("/stock"),
			expected: "",
		},
		{
			name:     "Given an empty command, When GetValue is called, Then it should return an empty string",
			command:  Command(""),
			expected: "",
		},
		{
			name:     "Given command with equals sign but no value, When GetValue is called, Then it should return an empty string",
			command:  Command("/stock="),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.command.GetValue()
			assert.Equal(t, tt.expected, result)
		})
	}
}
