package api

import (
	"context"
	"encoding/json"
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

func TestShouldListVaultMountsUsingSysMountsEndpoint(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/sys/mounts", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"secret/":{"type":"kv","options":{"version":"2"}},"auth/":{"type":"system"}}}`))
	}))
	defer server.Close()

	client, err := NewVaultClient(server.Client(), server.URL, "", "test-token")
	require.NoError(t, err)

	mounts, err := client.ListMounts(context.Background())

	require.NoError(t, err)
	assert.Equal(t, map[string]VaultMount{
		"auth":   {Type: "system", Options: nil},
		"secret": {Type: "kv", Options: map[string]string{"version": "2"}},
	}, mounts)
}

func TestShouldListVaultKVV1SecretsRecursivelyWithoutReadingValues(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/sys/mounts/secret":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"type":"kv"}`))
		case r.Method == "LIST" && r.URL.Path == "/v1/secret":
			require.Empty(t, r.URL.RawQuery)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"keys":["app/","root"]}}`))
		case r.Method == "LIST" && r.URL.Path == "/v1/secret/app":
			require.Empty(t, r.URL.RawQuery)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"keys":["config"]}}`))
		default:
			t.Fatalf("unexpected request: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
		}
	}))
	defer server.Close()

	client, err := NewVaultClient(server.Client(), server.URL, "", "test-token")
	require.NoError(t, err)

	secrets, err := client.ListKVSecrets(context.Background(), "secret")

	require.NoError(t, err)
	assert.Equal(t, []VaultSecret{
		{Ref: "secret/app/config", Name: "config", Provider: "HashiCorp Vault", Path: "app/config", Type: "kv-v1"},
		{Ref: "secret/root", Name: "root", Provider: "HashiCorp Vault", Path: "root", Type: "kv-v1"},
	}, secrets)
}

func TestShouldListVaultKVV2SecretsRecursivelyUsingMetadataPaths(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/sys/mounts/secret":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"type":"kv","options":{"version":"2"}}}`))
		case r.Method == "LIST" && r.URL.Path == "/v1/secret/metadata":
			require.Empty(t, r.URL.RawQuery)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"keys":["app/","root"]}}`))
		case r.Method == "LIST" && r.URL.Path == "/v1/secret/metadata/app":
			require.Empty(t, r.URL.RawQuery)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"keys":["config"]}}`))
		default:
			t.Fatalf("unexpected request: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
		}
	}))
	defer server.Close()

	client, err := NewVaultClient(server.Client(), server.URL, "", "test-token")
	require.NoError(t, err)

	secrets, err := client.ListKVSecrets(context.Background(), "secret")

	require.NoError(t, err)
	assert.Equal(t, []VaultSecret{
		{Ref: "secret/app/config", Name: "config", Provider: "HashiCorp Vault", Path: "app/config", Type: "kv-v2"},
		{Ref: "secret/root", Name: "root", Provider: "HashiCorp Vault", Path: "root", Type: "kv-v2"},
	}, secrets)
}

func TestShouldWriteVaultKVV1SecretUsingMountMetadata(t *testing.T) {
	t.Parallel()

	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/sys/mounts/secret":
			require.Equal(t, "write-token", r.Header.Get(vaultTokenHeader))
			require.Equal(t, "admin", r.Header.Get(vaultNamespaceHeader))
			require.Equal(t, "true", r.Header.Get(vaultRequestHeader))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"type":"kv","options":null}`))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/secret/app/config":
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			var body map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, map[string]any{"foo": "bar"}, body)
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client, err := NewVaultClient(server.Client(), server.URL, "admin", "write-token")
	require.NoError(t, err)

	// Act
	err = client.SetKVV1Secret(
		context.Background(),
		"secret",
		"app/config",
		json.RawMessage(`{"foo":"bar"}`),
	)

	// Assert
	require.NoError(t, err)
}

func TestShouldWriteVaultKVV2SecretUsingMountMetadata(t *testing.T) {
	t.Parallel()

	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/sys/mounts/secret":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"type":"kv","options":{"version":"2"}}}`))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/secret/data/app/config":
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			var body struct {
				Options struct {
					CAS int `json:"cas"`
				} `json:"options"`
				Data map[string]any `json:"data"`
			}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, 3, body.Options.CAS)
			assert.Equal(t, map[string]any{"foo": "bar"}, body.Data)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"version":4}}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client, err := NewVaultClient(server.Client(), server.URL, "", "write-token")
	require.NoError(t, err)
	cas := 3

	// Act
	err = client.SetKVV2Secret(
		context.Background(),
		"secret/",
		"/app/config/",
		json.RawMessage(`{"foo":"bar"}`),
		&cas,
	)

	// Assert
	require.NoError(t, err)
}

func TestShouldRejectWritingVaultKVV1SecretWhenMountIsKVV2(t *testing.T) {
	t.Parallel()

	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/sys/mounts/secret", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"type":"kv","options":{"version":"2"}}}`))
	}))
	defer server.Close()

	client, err := NewVaultClient(server.Client(), server.URL, "", "write-token")
	require.NoError(t, err)

	// Act
	err = client.SetKVV1Secret(
		context.Background(),
		"secret",
		"app/config",
		json.RawMessage(`{"foo":"bar"}`),
	)

	// Assert
	require.Error(t, err)
	assert.EqualError(t, err, "vault mount secret is kv v2, not kv v1")
}

func TestShouldRejectWritingVaultKVV2SecretWhenMountIsKVV1(t *testing.T) {
	t.Parallel()

	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/sys/mounts/secret", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"type":"kv"}`))
	}))
	defer server.Close()

	client, err := NewVaultClient(server.Client(), server.URL, "", "write-token")
	require.NoError(t, err)

	// Act
	err = client.SetKVV2Secret(
		context.Background(),
		"secret",
		"app/config",
		json.RawMessage(`{"foo":"bar"}`),
		nil,
	)

	// Assert
	require.Error(t, err)
	assert.EqualError(t, err, "vault mount secret is kv v1, not kv v2")
}
