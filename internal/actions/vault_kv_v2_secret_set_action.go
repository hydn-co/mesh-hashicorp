package actions

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hydn-co/mesh-hashicorp/internal/api"
	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-hashicorp/internal/payloads"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewVaultKVV2SecretSetAction(
	ctx *connector.TypedFeatureContext[*options.VaultKVV2SecretSetActionOptions, *payloads.VaultKVV2SecretSetPayload],
) runner.Feature {
	return &VaultKVV2SecretSetAction{TypedFeatureContext: ctx}
}

type VaultKVV2SecretSetAction struct {
	state   connectorutil.FeatureState
	token   string
	opts    *options.VaultKVV2SecretSetActionOptions
	payload *payloads.VaultKVV2SecretSetPayload
	*connector.TypedFeatureContext[*options.VaultKVV2SecretSetActionOptions, *payloads.VaultKVV2SecretSetPayload]
}

func (a *VaultKVV2SecretSetAction) Init(ctx context.Context) error {
	opts := a.GetOptions()
	if err := connectorutil.Validate(opts, "feature options"); err != nil {
		return err
	}
	payload := a.GetPayload()
	if err := connectorutil.Validate(payload, "vault kv v2 secret set payload"); err != nil {
		return err
	}
	token, err := connectorutil.ExtractAPIKey(a.GetCredentials())
	if err != nil {
		return fmt.Errorf("parse api key credentials: %w", err)
	}

	a.opts = opts
	a.payload = payload
	a.token = token
	a.state.MarkReady()

	connectorutil.LogFeature(
		ctx,
		a.TypedFeatureContext,
		slog.LevelInfo,
		"Initialized Vault KV v2 secret set action",
		"mount_path", payload.MountPath,
		"secret_path", payload.SecretPath,
		"cas_provided", payload.CAS != nil,
	)
	return nil
}

func (a *VaultKVV2SecretSetAction) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := a.state.RequireReady(); err != nil {
		return err
	}

	connectorutil.LogFeature(
		ctx,
		a.TypedFeatureContext,
		slog.LevelInfo,
		"Starting Vault KV v2 secret set action",
		"mount_path", a.payload.MountPath,
		"secret_path", a.payload.SecretPath,
		"cas_provided", a.payload.CAS != nil,
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

	if err := client.SetKVV2Secret(
		ctx,
		a.payload.MountPath,
		a.payload.SecretPath,
		a.payload.Data,
		a.payload.CAS,
	); err != nil {
		return fmt.Errorf("set vault kv v2 secret: %w", err)
	}

	connectorutil.LogFeature(
		ctx,
		a.TypedFeatureContext,
		slog.LevelInfo,
		"Finished Vault KV v2 secret set action",
		"mount_path", a.payload.MountPath,
		"secret_path", a.payload.SecretPath,
		"kv_version", "2",
		"cas_provided", a.payload.CAS != nil,
	)
	return nil
}

func (a *VaultKVV2SecretSetAction) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := a.state.RequireReady(); err != nil {
		return err
	}

	connectorutil.LogFeature(
		ctx,
		a.TypedFeatureContext,
		slog.LevelInfo,
		"Stopping Vault KV v2 secret set action",
	)
	a.state.Reset()
	a.token = ""
	a.opts = nil
	a.payload = nil
	return nil
}
