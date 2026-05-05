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

func NewTerraformTeamAccessEntityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformTeamAccessEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformTeamAccessEntityCollector{TypedFeatureContext: ctx}
}

type TerraformTeamAccessEntityCollector struct {
	initialized bool
	token       string
	*connector.TypedFeatureContext[*options.TerraformTeamAccessEntityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformTeamAccessEntityCollector) Init(_ context.Context) error {
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

func (c *TerraformTeamAccessEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !c.initialized {
		return fmt.Errorf("terraform team access entity collector not initialized")
	}
	collectors.LogCollector(
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
	if !c.initialized {
		return fmt.Errorf("terraform team access entity collector not initialized")
	}
	c.initialized = false
	c.token = ""
	return nil
}
