package collectors

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-hashicorp/internal/credentials"
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
	token string
	*connector.TypedFeatureContext[*options.TerraformPolicyEntityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformPolicyEntityCollector) Init(_ context.Context) error {
	if err := options.ValidateTerraformOptions(c.GetOptions()); err != nil {
		return err
	}
	token, err := credentials.ExtractToken(c.GetCredentials())
	if err != nil {
		return fmt.Errorf("parse api key credentials: %w", err)
	}
	c.token = token

	return nil
}

func (c *TerraformPolicyEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logCollector(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting HCP Terraform policy entity collector")
	return fmt.Errorf("terraform policy entity collector not implemented")
}

func (c *TerraformPolicyEntityCollector) Stop(context.Context) error { return nil }
