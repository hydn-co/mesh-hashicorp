package credentials

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldBuildBearerAuthorizationHeaderWhenTokenValid(t *testing.T) {
	headerValue, err := GetBearerAuthorizationHeaderValue("abc123")

	require.NoError(t, err)
	assert.Equal(t, "Bearer abc123", headerValue)
}

func TestShouldRejectBearerAuthorizationHeaderWhenTokenEmpty(t *testing.T) {
	_, err := GetBearerAuthorizationHeaderValue("")

	require.Error(t, err)
	assert.EqualError(t, err, "token is required")
}

func TestShouldTrimWhitespaceWhenBuildingBearerAuthorizationHeader(t *testing.T) {
	headerValue, err := GetBearerAuthorizationHeaderValue(" abc123 ")

	require.NoError(t, err)
	assert.Equal(t, "Bearer abc123", headerValue)
}

func TestShouldExtractTokenWhenAPIKeyProvided(t *testing.T) {
	token, err := ExtractToken(json.RawMessage(`{"api_key":"abc123"}`))

	require.NoError(t, err)
	assert.Equal(t, "abc123", token)
}

func TestShouldTrimWhitespaceWhenExtractingToken(t *testing.T) {
	token, err := ExtractToken(json.RawMessage(`{"api_key":" abc123 "}`))

	require.NoError(t, err)
	assert.Equal(t, "abc123", token)
}

func TestShouldRejectCredentialsWhenAPIKeyMissing(t *testing.T) {
	_, err := ExtractToken(json.RawMessage(`{"foo":"bar"}`))

	require.Error(t, err)
	assert.EqualError(t, err, "api_key is required")
}

func TestShouldRejectCredentialsWhenJSONInvalid(t *testing.T) {
	_, err := ExtractToken(json.RawMessage(`{"api_key":`))

	require.Error(t, err)
	assert.ErrorContains(t, err, "decode api key credentials")
}
