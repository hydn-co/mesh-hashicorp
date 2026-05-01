package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldRejectClientWhenTokenInvalid(t *testing.T) {
	client, err := NewClient(http.DefaultClient, "app.terraform.io", "")

	require.Error(t, err)
	assert.Nil(t, client)
	assert.EqualError(t, err, "validate terraform token: token is required")
}
