package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldRejectVaultClientWhenTokenInvalid(t *testing.T) {
	// Arrange

	// Act
	client, err := NewVaultClient(http.DefaultClient, "vault.example.com", "", "")

	// Assert
	require.Error(t, err)
	assert.Nil(t, client)
	assert.EqualError(t, err, "validate vault token: token is required")
}

func TestShouldRejectVaultClientWhenHTTPClientMissing(t *testing.T) {
	// Arrange

	// Act
	client, err := NewVaultClient(nil, "vault.example.com", "", "abc123")

	// Assert
	require.Error(t, err)
	assert.Nil(t, client)
	assert.EqualError(t, err, "http client is required")
}

func TestShouldListVaultEntityIDsUsingGetListQueryAndHeaders(t *testing.T) {
	t.Parallel()

	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/identity/entity/id", r.URL.Path)
		require.Equal(t, "true", r.URL.Query().Get("list"))
		require.Equal(t, "test-token", r.Header.Get(vaultTokenHeader))
		require.Equal(t, "admin", r.Header.Get(vaultNamespaceHeader))
		require.Equal(t, "true", r.Header.Get(vaultRequestHeader))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"keys":["entity-b","entity-a"]}}`))
	}))
	defer server.Close()

	client, err := NewVaultClient(server.Client(), server.URL, "admin", "test-token")
	require.NoError(t, err)

	// Act
	keys, err := client.ListIdentityEntityIDs(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"entity-b", "entity-a"}, keys)
}

func TestShouldTreatEmptyVaultListAsNoResults(t *testing.T) {
	t.Parallel()

	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/identity/group/id", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"errors":[]}`))
	}))
	defer server.Close()

	client, err := NewVaultClient(server.Client(), server.URL, "", "test-token")
	require.NoError(t, err)

	// Act
	keys, err := client.ListIdentityGroupIDs(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Nil(t, keys)
}
