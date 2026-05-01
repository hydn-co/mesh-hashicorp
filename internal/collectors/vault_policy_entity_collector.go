package collectors

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewVaultPolicyEntityCollector(
	ctx *connector.TypedFeatureContext[*options.VaultPolicyEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &VaultPolicyEntityCollector{TypedFeatureContext: ctx}
}

type VaultPolicyEntityCollector struct {
	*connector.TypedFeatureContext[*options.VaultPolicyEntityCollectorOptions, *connector.NoPayload]
}

func (c *VaultPolicyEntityCollector) Init(_ context.Context) error {
	return options.ValidateVaultOptions(c.GetOptions())
}

func (c *VaultPolicyEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logCollector(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting Vault policy entity collector")
	return fmt.Errorf("vault policy entity collector not implemented")
}

func (c *VaultPolicyEntityCollector) Stop(context.Context) error { return nil }
