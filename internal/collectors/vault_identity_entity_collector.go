package collectors

import (
	"context"
	"fmt"
	"log/slog"

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
	*connector.TypedFeatureContext[*options.VaultIdentityEntityCollectorOptions, *connector.NoPayload]
}

func (c *VaultIdentityEntityCollector) Init(_ context.Context) error {
	return options.ValidateVaultOptions(c.GetOptions())
}

func (c *VaultIdentityEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logCollector(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting Vault identity entity collector")
	return fmt.Errorf("vault identity entity collector not implemented")
}

func (c *VaultIdentityEntityCollector) Stop(context.Context) error { return nil }
