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

func NewTerraformWorkspaceEntityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformWorkspaceEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformWorkspaceEntityCollector{TypedFeatureContext: ctx}
}

type TerraformWorkspaceEntityCollector struct {
	state connectorutil.FeatureState
	token string
	*connector.TypedFeatureContext[*options.TerraformWorkspaceEntityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformWorkspaceEntityCollector) Init(_ context.Context) error {
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

func (c *TerraformWorkspaceEntityCollector) Start(ctx context.Context) error {
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
		"Starting HCP Terraform workspace entity collector",
	)
	return fmt.Errorf("terraform workspace entity collector not implemented")
}

func (c *TerraformWorkspaceEntityCollector) Stop(ctx context.Context) error {
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
