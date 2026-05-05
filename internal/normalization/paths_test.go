package normalization

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldNormalizeVaultMountPathByTrimmingWhitespaceAndSlashes(t *testing.T) {
	// Arrange

	// Act
	normalizedValue, err := NormalizeVaultMountPath("  /secret/child/  ")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "secret/child", normalizedValue)
}

func TestShouldRejectNormalizeVaultSecretPathWhenEmpty(t *testing.T) {
	// Arrange

	// Act
	normalizedValue, err := NormalizeVaultSecretPath("  ///  ")

	// Assert
	require.Error(t, err)
	assert.Empty(t, normalizedValue)
	assert.EqualError(t, err, "vault secret path is required")
}
