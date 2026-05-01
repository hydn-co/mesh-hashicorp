package collectors

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewTerraformWorkspaceEntityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformWorkspaceEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformWorkspaceEntityCollector{TypedFeatureContext: ctx}
}

type TerraformWorkspaceEntityCollector struct {
	*connector.TypedFeatureContext[*options.TerraformWorkspaceEntityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformWorkspaceEntityCollector) Init(_ context.Context) error {
	return options.ValidateTerraformOptions(c.GetOptions())
}

func (c *TerraformWorkspaceEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logCollector(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting HCP Terraform workspace entity collector")
	return fmt.Errorf("terraform workspace entity collector not implemented")
}

func (c *TerraformWorkspaceEntityCollector) Stop(context.Context) error { return nil }
