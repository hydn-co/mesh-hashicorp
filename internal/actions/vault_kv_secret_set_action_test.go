package actions

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/fgrzl/json/polymorphic"
	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-hashicorp/internal/payloads"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newVaultKVV1SecretSetActionForTest(
	opts *options.VaultKVV1SecretSetActionOptions,
	payload *payloads.VaultKVV1SecretSetPayload,
	credentials json.RawMessage,
) *VaultKVV1SecretSetAction {
	config := &connector.Configuration{
		Options:     polymorphic.NewEnvelope(opts),
		Credentials: credentials,
	}
	if payload != nil {
		config.Payload = polymorphic.NewEnvelope(payload)
	}
	ctx := connector.NewTypedFeatureContext[
		*options.VaultKVV1SecretSetActionOptions,
		*payloads.VaultKVV1SecretSetPayload,
	](connector.NewFeatureContext(connector.WithConfiguration(config)))
	return &VaultKVV1SecretSetAction{TypedFeatureContext: ctx}
}

func newVaultKVV2SecretSetActionForTest(
	opts *options.VaultKVV2SecretSetActionOptions,
	payload *payloads.VaultKVV2SecretSetPayload,
	credentials json.RawMessage,
) *VaultKVV2SecretSetAction {
	config := &connector.Configuration{
		Options:     polymorphic.NewEnvelope(opts),
		Credentials: credentials,
	}
	if payload != nil {
		config.Payload = polymorphic.NewEnvelope(payload)
	}
	ctx := connector.NewTypedFeatureContext[
		*options.VaultKVV2SecretSetActionOptions,
		*payloads.VaultKVV2SecretSetPayload,
	](connector.NewFeatureContext(connector.WithConfiguration(config)))
	return &VaultKVV2SecretSetAction{TypedFeatureContext: ctx}
}

func TestShouldRejectVaultKVV1SecretSetInitWhenPayloadMissing(t *testing.T) {
	// Arrange
	action := newVaultKVV1SecretSetActionForTest(
		&options.VaultKVV1SecretSetActionOptions{
			VaultOptionsCore: options.VaultOptionsCore{Address: "https://vault.example.com"},
		},
		nil,
		nil,
	)

	// Act
	err := action.Init(context.Background())

	// Assert
	require.Error(t, err)
	assert.EqualError(t, err, "vault kv v1 secret set payload required but not provided")
}

func TestShouldRejectVaultKVV1SecretSetStartWhenNotInitialized(t *testing.T) {
	// Arrange
	action := newVaultKVV1SecretSetActionForTest(
		&options.VaultKVV1SecretSetActionOptions{
			VaultOptionsCore: options.VaultOptionsCore{Address: "https://vault.example.com"},
		},
		&payloads.VaultKVV1SecretSetPayload{
			MountPath:  "secret",
			SecretPath: "app/config",
			Data:       json.RawMessage(`{"foo":"bar"}`),
		},
		json.RawMessage(`{"api_key":"token"}`),
	)

	// Act
	err := action.Start(context.Background())

	// Assert
	require.Error(t, err)
	assert.EqualError(t, err, "feature not ready")
}

func TestShouldRejectVaultKVV2SecretSetStopWhenNotInitialized(t *testing.T) {
	// Arrange
	action := newVaultKVV2SecretSetActionForTest(
		&options.VaultKVV2SecretSetActionOptions{
			VaultOptionsCore: options.VaultOptionsCore{Address: "https://vault.example.com"},
		},
		&payloads.VaultKVV2SecretSetPayload{
			MountPath:  "secret",
			SecretPath: "app/config",
			Data:       json.RawMessage(`{"foo":"bar"}`),
		},
		json.RawMessage(`{"api_key":"token"}`),
	)

	// Act
	err := action.Stop(context.Background())

	// Assert
	require.Error(t, err)
	assert.EqualError(t, err, "feature not ready")
}

func TestShouldRejectVaultKVV2SecretSetInitWhenPayloadInvalid(t *testing.T) {
	// Arrange
	action := newVaultKVV2SecretSetActionForTest(
		&options.VaultKVV2SecretSetActionOptions{
			VaultOptionsCore: options.VaultOptionsCore{Address: "https://vault.example.com"},
		},
		&payloads.VaultKVV2SecretSetPayload{SecretPath: "app/config", Data: json.RawMessage(`{"foo":"bar"}`)},
		nil,
	)

	// Act
	err := action.Init(context.Background())

	// Assert
	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid vault kv v2 secret set payload")
	assert.ErrorContains(t, err, "mount_path is required")
}
