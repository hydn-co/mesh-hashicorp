package entity

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-hashicorp/internal/collectors"
	"github.com/hydn-co/mesh-hashicorp/internal/credentials"
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
	initialized bool
	token       string
	*connector.TypedFeatureContext[*options.TerraformWorkspaceEntityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformWorkspaceEntityCollector) Init(_ context.Context) error {
	if err := options.ValidateTerraformOptions(c.GetOptions()); err != nil {
		return err
	}
	token, err := credentials.ExtractToken(c.GetCredentials())
	if err != nil {
		return fmt.Errorf("parse api key credentials: %w", err)
	}
	c.token = token
	c.initialized = true

	return nil
}

func (c *TerraformWorkspaceEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !c.initialized {
		return fmt.Errorf("terraform workspace entity collector not initialized")
	}
	collectors.LogCollector(
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
	if !c.initialized {
		return fmt.Errorf("terraform workspace entity collector not initialized")
	}
	c.initialized = false
	c.token = ""
	return nil
}
