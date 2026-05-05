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

func NewVaultIdentityEntityCollector(
	ctx *connector.TypedFeatureContext[*options.VaultIdentityEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &VaultIdentityEntityCollector{TypedFeatureContext: ctx}
}

type VaultIdentityEntityCollector struct {
	token string
	*connector.TypedFeatureContext[*options.VaultIdentityEntityCollectorOptions, *connector.NoPayload]
}

func (c *VaultIdentityEntityCollector) Init(_ context.Context) error {
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

func (c *VaultIdentityEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	collectors.LogCollector(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Starting Vault identity entity collector",
	)
	return fmt.Errorf("vault identity entity collector not implemented")
}

func (c *VaultIdentityEntityCollector) Stop(context.Context) error { return nil }
