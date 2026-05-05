package actions

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hydn-co/mesh-hashicorp/internal/api"
	"github.com/hydn-co/mesh-hashicorp/internal/credentials"
	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-hashicorp/internal/payloads"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewVaultKVV1SecretSetAction(
	ctx *connector.TypedFeatureContext[*options.VaultKVV1SecretSetActionOptions, *payloads.VaultKVV1SecretSetPayload],
) runner.Feature {
	return &VaultKVV1SecretSetAction{TypedFeatureContext: ctx}
}

type VaultKVV1SecretSetAction struct {
	initialized bool
	token       string
	opts        *options.VaultKVV1SecretSetActionOptions
	payload     *payloads.VaultKVV1SecretSetPayload
	*connector.TypedFeatureContext[*options.VaultKVV1SecretSetActionOptions, *payloads.VaultKVV1SecretSetPayload]
}

func (a *VaultKVV1SecretSetAction) Init(ctx context.Context) error {
	opts := a.GetOptions()
	if err := options.ValidateVaultOptions(opts); err != nil {
		return err
	}
	payload := a.GetPayload()
	if err := payloads.ValidatePayload(payload, "vault kv v1 secret set payload"); err != nil {
		return err
	}
	token, err := credentials.ExtractToken(a.GetCredentials())
	if err != nil {
		return fmt.Errorf("parse api key credentials: %w", err)
	}

	a.opts = opts
	a.payload = payload
	a.token = token
	a.initialized = true

	logAction(
		ctx,
		a.TypedFeatureContext,
		slog.LevelInfo,
		"Initialized Vault KV v1 secret set action",
		"mount_path", payload.MountPath,
		"secret_path", payload.SecretPath,
	)
	return nil
}

func (a *VaultKVV1SecretSetAction) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !a.initialized {
		return fmt.Errorf("vault kv v1 secret set action not initialized")
	}

	logAction(
		ctx,
		a.TypedFeatureContext,
		slog.LevelInfo,
		"Starting Vault KV v1 secret set action",
		"mount_path", a.payload.MountPath,
		"secret_path", a.payload.SecretPath,
	)

	client, err := api.NewVaultClient(
		http.DefaultClient,
		a.opts.GetAddress(),
		a.opts.GetNamespace(),
		a.token,
	)
	if err != nil {
		return fmt.Errorf("build vault client: %w", err)
	}

	if err := client.SetKVV1Secret(
		ctx,
		a.payload.MountPath,
		a.payload.SecretPath,
		a.payload.Data,
	); err != nil {
		return fmt.Errorf("set vault kv v1 secret: %w", err)
	}

	logAction(
		ctx,
		a.TypedFeatureContext,
		slog.LevelInfo,
		"Finished Vault KV v1 secret set action",
		"mount_path", a.payload.MountPath,
		"secret_path", a.payload.SecretPath,
		"kv_version", "1",
	)
	return nil
}

func (a *VaultKVV1SecretSetAction) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !a.initialized {
		return fmt.Errorf("vault kv v1 secret set action not initialized")
	}

	logAction(
		ctx,
		a.TypedFeatureContext,
		slog.LevelInfo,
		"Stopping Vault KV v1 secret set action",
	)
	a.initialized = false
	a.token = ""
	a.opts = nil
	a.payload = nil
	return nil
}
