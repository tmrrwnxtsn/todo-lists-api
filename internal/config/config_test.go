package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_getEnv(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		defaultValue  string
		existingValue string
		expectedValue string
	}{
		{
			name:          "get existing env variable",
			key:           "TOKEN",
			defaultValue:  "12345",
			existingValue: "54321",
			expectedValue: "54321",
		},
		{
			name:          "get default env variable",
			key:           "TOKEN",
			defaultValue:  "12345",
			existingValue: "",
			expectedValue: "12345",
		},
		{
			name:          "empty variable key",
			key:           "",
			defaultValue:  "12345",
			existingValue: "",
			expectedValue: "12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.key != "" && tt.existingValue != "" {
				err := os.Setenv(tt.key, tt.existingValue)
				assert.NoError(t, err)
				defer func() {
					err := os.Unsetenv(tt.key)
					assert.NoError(t, err)
				}()
			}

			variable := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, variable, tt.expectedValue)
		})
	}
}
