package entity

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-hashicorp/internal/credentials"
	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewTerraformTeamAccessEntityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformTeamAccessEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformTeamAccessEntityCollector{TypedFeatureContext: ctx}
}

type TerraformTeamAccessEntityCollector struct {
	state connectorutil.FeatureState
	token string
	*connector.TypedFeatureContext[*options.TerraformTeamAccessEntityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformTeamAccessEntityCollector) Init(_ context.Context) error {
	if err := connectorutil.Validate(c.GetOptions(), "feature options"); err != nil {
		return err
	}
	token, err := credentials.ExtractToken(c.GetCredentials())
	if err != nil {
		return fmt.Errorf("parse api key credentials: %w", err)
	}
	c.token = token
	c.state.MarkReady()

	return nil
}

func (c *TerraformTeamAccessEntityCollector) Start(ctx context.Context) error {
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
		"Starting HCP Terraform team access entity collector",
	)
	return fmt.Errorf("terraform team access entity collector not implemented")
}

func (c *TerraformTeamAccessEntityCollector) Stop(ctx context.Context) error {
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
