package credentials

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldBuildBearerAuthorizationHeaderWhenTokenValid(t *testing.T) {
	// Arrange

	// Act
	headerValue, err := GetBearerAuthorizationHeaderValue("abc123")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Bearer abc123", headerValue)
}

func TestShouldRejectBearerAuthorizationHeaderWhenTokenEmpty(t *testing.T) {
	// Arrange

	// Act
	_, err := GetBearerAuthorizationHeaderValue("")

	// Assert
	require.Error(t, err)
	assert.EqualError(t, err, "token is required")
}

func TestShouldTrimWhitespaceWhenBuildingBearerAuthorizationHeader(t *testing.T) {
	// Arrange

	// Act
	headerValue, err := GetBearerAuthorizationHeaderValue(" abc123 ")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Bearer abc123", headerValue)
}

func TestShouldNormalizeTokenWhenTokenValid(t *testing.T) {
	// Arrange

	// Act
	token, err := NormalizeToken(" abc123 ")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "abc123", token)
}

func TestShouldRejectNormalizeTokenWhenTokenEmpty(t *testing.T) {
	// Arrange

	// Act
	_, err := NormalizeToken("")

	// Assert
	require.Error(t, err)
	assert.EqualError(t, err, "token is required")
}

func TestShouldExtractTokenWhenAPIKeyProvided(t *testing.T) {
	// Arrange

	// Act
	token, err := ExtractToken(json.RawMessage(`{"api_key":"abc123"}`))

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "abc123", token)
}

func TestShouldTrimWhitespaceWhenExtractingToken(t *testing.T) {
	// Arrange

	// Act
	token, err := ExtractToken(json.RawMessage(`{"api_key":" abc123 "}`))

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "abc123", token)
}

func TestShouldRejectCredentialsWhenAPIKeyMissing(t *testing.T) {
	// Arrange

	// Act
	_, err := ExtractToken(json.RawMessage(`{"foo":"bar"}`))

	// Assert
	require.Error(t, err)
	assert.EqualError(t, err, "api_key is required")
}

func TestShouldRejectCredentialsWhenJSONInvalid(t *testing.T) {
	// Arrange

	// Act
	_, err := ExtractToken(json.RawMessage(`{"api_key":`))

	// Assert
	require.Error(t, err)
	assert.ErrorContains(t, err, "decode api key credentials")
}
