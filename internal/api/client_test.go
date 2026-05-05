package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldRejectClientWhenTokenInvalid(t *testing.T) {
	// Arrange

	// Act
	client, err := NewClient(http.DefaultClient, "terraform", "app.terraform.io", "")

	// Assert
	require.Error(t, err)
	assert.Nil(t, client)
	assert.EqualError(t, err, "validate terraform token: token is required")
}

func TestShouldRejectClientWhenHTTPClientMissing(t *testing.T) {
	// Arrange

	// Act
	client, err := NewClient(nil, "terraform", "app.terraform.io", "abc123")

	// Assert
	require.Error(t, err)
	assert.Nil(t, client)
	assert.EqualError(t, err, "http client is required")
}
