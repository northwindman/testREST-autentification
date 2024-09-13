package refresh

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew_Success(t *testing.T) {
	length := 16

	token, err := New(length)

	require.NoError(t, err)

	assert.Equal(t, length, len(token))
}

func TestNew_ZeroLength(t *testing.T) {
	length := 0

	token, err := New(length)

	require.NoError(t, err)

	assert.Empty(t, token)
}

func TestNew_NegativeLength(t *testing.T) {
	length := -5

	token, err := New(length)

	require.Error(t, err)

	assert.Empty(t, token)
}
