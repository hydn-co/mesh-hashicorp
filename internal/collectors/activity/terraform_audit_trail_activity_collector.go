package activity

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

func NewTerraformAuditTrailActivityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformAuditTrailActivityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformAuditTrailActivityCollector{TypedFeatureContext: ctx}
}

type TerraformAuditTrailActivityCollector struct {
	initialized bool
	token       string
	*connector.TypedFeatureContext[*options.TerraformAuditTrailActivityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformAuditTrailActivityCollector) Init(_ context.Context) error {
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

func (c *TerraformAuditTrailActivityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !c.initialized {
		return fmt.Errorf("terraform audit trail activity collector not initialized")
	}
	collectors.LogCollector(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Starting HCP Terraform audit trail activity collector",
	)
	return fmt.Errorf("terraform audit trail activity collector not implemented")
}

func (c *TerraformAuditTrailActivityCollector) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !c.initialized {
		return fmt.Errorf("terraform audit trail activity collector not initialized")
	}
	c.initialized = false
	c.token = ""
	return nil
}
