package collectors

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewTerraformAccountEntityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformAccountEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformAccountEntityCollector{TypedFeatureContext: ctx}
}

type TerraformAccountEntityCollector struct {
	*connector.TypedFeatureContext[*options.TerraformAccountEntityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformAccountEntityCollector) Init(_ context.Context) error {
	return options.ValidateTerraformOptions(c.GetOptions())
}

func (c *TerraformAccountEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logCollector(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting HCP Terraform account entity collector")
	return fmt.Errorf("terraform account entity collector not implemented")
}

func (c *TerraformAccountEntityCollector) Stop(context.Context) error { return nil }
