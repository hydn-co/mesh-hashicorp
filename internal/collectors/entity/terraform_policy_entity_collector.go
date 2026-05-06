package entity

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewTerraformPolicyEntityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformPolicyEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformPolicyEntityCollector{TypedFeatureContext: ctx}
}

type TerraformPolicyEntityCollector struct {
	state connectorutil.FeatureState
	token string
	*connector.TypedFeatureContext[*options.TerraformPolicyEntityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformPolicyEntityCollector) Init(_ context.Context) error {
	if err := connectorutil.Validate(c.GetOptions(), "feature options"); err != nil {
		return err
	}
	token, err := connectorutil.ExtractAPIKey(c.GetCredentials())
	if err != nil {
		return fmt.Errorf("parse api key credentials: %w", err)
	}
	c.token = token
	c.state.MarkReady()

	return nil
}

func (c *TerraformPolicyEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := c.state.RequireReady(); err != nil {
		return err
	}
	connectorutil.LogFeature(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Starting HCP Terraform policy entity collector",
	)
	return fmt.Errorf("terraform policy entity collector not implemented")
}

func (c *TerraformPolicyEntityCollector) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := c.state.RequireReady(); err != nil {
		return err
	}
	c.state.Reset()
	c.token = ""
	return nil
}
