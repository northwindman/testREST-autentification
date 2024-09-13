package format

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashStringEmpty(t *testing.T) {
	input := ""
	hashed, err := HashString(input)
	assert.Error(t, err)
	assert.Empty(t, hashed)
}

func TestVerifyStringWithCorruptedHash(t *testing.T) {
	validString := "test_string"
	corruptedHash := "corrupted_hash"

	ok := VerifyString(validString, corruptedHash)
	assert.False(t, ok)
}

func TestInBase64Empty(t *testing.T) {
	input := ""
	expected := ""
	encoded := InBase64(input)
	assert.Equal(t, expected, encoded)
}

func TestFromBase64Empty(t *testing.T) {
	encoded := ""
	decoded, err := FromBase64(encoded)
	assert.Error(t, err)
	assert.Empty(t, decoded)
}
