package entity

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-hashicorp/internal/collectors"
	"github.com/hydn-co/mesh-hashicorp/internal/credentials"
	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewVaultAuthMethodEntityCollector(
	ctx *connector.TypedFeatureContext[*options.VaultAuthMethodEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &VaultAuthMethodEntityCollector{TypedFeatureContext: ctx}
}

type VaultAuthMethodEntityCollector struct {
	token string
	*connector.TypedFeatureContext[*options.VaultAuthMethodEntityCollectorOptions, *connector.NoPayload]
}

func (c *VaultAuthMethodEntityCollector) Init(_ context.Context) error {
	if err := options.ValidateVaultOptions(c.GetOptions()); err != nil {
		return err
	}
	token, err := credentials.ExtractToken(c.GetCredentials())
	if err != nil {
		return fmt.Errorf("parse api key credentials: %w", err)
	}
	c.token = token

	return nil
}

func (c *VaultAuthMethodEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	collectors.LogCollector(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Starting Vault auth method entity collector",
	)
	return fmt.Errorf("vault auth method entity collector not implemented")
}

func (c *VaultAuthMethodEntityCollector) Stop(context.Context) error { return nil }
