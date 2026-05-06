package entity

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sort"

	"github.com/fgrzl/enumerators"
	"github.com/hydn-co/mesh-hashicorp/internal/api"
	"github.com/hydn-co/mesh-hashicorp/internal/collectors"
	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewVaultSecretEntityCollector(
	ctx *connector.TypedFeatureContext[*options.VaultSecretEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &VaultSecretEntityCollector{TypedFeatureContext: ctx}
}

type VaultSecretEntityCollector struct {
	state connectorutil.FeatureState
	token string
	*connector.TypedFeatureContext[*options.VaultSecretEntityCollectorOptions, *connector.NoPayload]
}

func (c *VaultSecretEntityCollector) Init(_ context.Context) error {
	opts := c.GetOptions()
	if err := connectorutil.Validate(opts, "feature options"); err != nil {
		return err
	}
	token, err := connectorutil.ExtractAPIKey(c.GetCredentials())
	if err != nil {
		return fmt.Errorf("parse api key credentials: %w", err)
	}
	c.token = token
	c.state.MarkReady()

	return nil
}

func (c *VaultSecretEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := c.state.RequireReady(); err != nil {
		return err
	}
	connectorutil.LogFeature(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Starting Vault secret entity collector",
	)
	opts := c.GetOptions()

	client, err := api.NewVaultClient(
		http.DefaultClient,
		opts.GetAddress(),
		opts.GetNamespace(),
		c.token,
	)
	if err != nil {
		return fmt.Errorf("build vault client: %w", err)
	}

	mounts, err := client.ListMounts(ctx)
	if err != nil {
		return fmt.Errorf("list vault mounts: %w", err)
	}

	mountPaths := make([]string, 0, len(mounts))
	for mountPath, mount := range mounts {
		if mount.Type == "kv" {
			mountPaths = append(mountPaths, mountPath)
		}
	}
	sort.Strings(mountPaths)

	if err := enumerators.ForEach(enumerators.Slice(mountPaths), func(mountPath string) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		secrets, err := client.ListKVSecrets(ctx, mountPath)
		if err != nil {
			return fmt.Errorf("list vault secrets for mount %s: %w", mountPath, err)
		}

		for _, secret := range secrets {
			secretEntity, err := collectors.NewVaultSecret(secret)
			if err != nil {
				return fmt.Errorf("map vault secret %s: %w", secret.Ref, err)
			}
			if err := c.Emit(ctx, secretEntity); err != nil {
				return fmt.Errorf("emit secret %s: %w", secretEntity.SecretRef, err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("enumerate vault secrets: %w", err)
	}

	connectorutil.LogFeature(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Finished Vault secret entity collector",
	)
	return nil
}

func (c *VaultSecretEntityCollector) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := c.state.RequireReady(); err != nil {
		return err
	}
	connectorutil.LogFeature(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Stopping Vault secret entity collector",
	)
	c.state.Reset()
	c.token = ""
	return nil
}
