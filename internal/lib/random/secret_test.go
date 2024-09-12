package random

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSecret(t *testing.T) {
	tests := []struct {
		name           string
		Length         int
		ExpectedLength int
	}{
		{
			"5 symbols",
			5,
			5,
		},
		{
			"10 symbols",
			10,
			10,
		},
		{
			"12 symbols",
			12,
			12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewSecret(tt.Length)
			require.NoError(t, err)
			assert.Equal(t, tt.ExpectedLength, len(token))
		})
	}
}
