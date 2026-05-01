package collectors

import (
	"context"
	"fmt"
	"log/slog"

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
	*connector.TypedFeatureContext[*options.VaultAuthMethodEntityCollectorOptions, *connector.NoPayload]
}

func (c *VaultAuthMethodEntityCollector) Init(_ context.Context) error {
	return options.ValidateVaultOptions(c.GetOptions())
}

func (c *VaultAuthMethodEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logCollector(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting Vault auth method entity collector")
	return fmt.Errorf("vault auth method entity collector not implemented")
}

func (c *VaultAuthMethodEntityCollector) Stop(context.Context) error { return nil }
