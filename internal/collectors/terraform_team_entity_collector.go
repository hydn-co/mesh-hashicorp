package collectors

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewTerraformTeamEntityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformTeamEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformTeamEntityCollector{TypedFeatureContext: ctx}
}

type TerraformTeamEntityCollector struct {
	*connector.TypedFeatureContext[*options.TerraformTeamEntityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformTeamEntityCollector) Init(_ context.Context) error {
	return options.ValidateTerraformOptions(c.GetOptions())
}

func (c *TerraformTeamEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logCollector(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting HCP Terraform team entity collector")
	return fmt.Errorf("terraform team entity collector not implemented")
}

func (c *TerraformTeamEntityCollector) Stop(context.Context) error { return nil }
