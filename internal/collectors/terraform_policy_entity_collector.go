package collectors

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewTerraformPolicyEntityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformPolicyEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformPolicyEntityCollector{TypedFeatureContext: ctx}
}

type TerraformPolicyEntityCollector struct {
	*connector.TypedFeatureContext[*options.TerraformPolicyEntityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformPolicyEntityCollector) Init(_ context.Context) error {
	return options.ValidateTerraformOptions(c.GetOptions())
}

func (c *TerraformPolicyEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logCollector(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting HCP Terraform policy entity collector")
	return fmt.Errorf("terraform policy entity collector not implemented")
}

func (c *TerraformPolicyEntityCollector) Stop(context.Context) error { return nil }
