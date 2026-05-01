package collectors

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewTerraformTeamAccessEntityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformTeamAccessEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformTeamAccessEntityCollector{TypedFeatureContext: ctx}
}

type TerraformTeamAccessEntityCollector struct {
	*connector.TypedFeatureContext[*options.TerraformTeamAccessEntityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformTeamAccessEntityCollector) Init(_ context.Context) error {
	return options.ValidateTerraformOptions(c.GetOptions())
}

func (c *TerraformTeamAccessEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logCollector(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting HCP Terraform team access entity collector")
	return fmt.Errorf("terraform team access entity collector not implemented")
}

func (c *TerraformTeamAccessEntityCollector) Stop(context.Context) error { return nil }
