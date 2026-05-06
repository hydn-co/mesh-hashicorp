package activity

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewTerraformAuditTrailActivityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformAuditTrailActivityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformAuditTrailActivityCollector{TypedFeatureContext: ctx}
}

type TerraformAuditTrailActivityCollector struct {
	state connectorutil.FeatureState
	token string
	*connector.TypedFeatureContext[*options.TerraformAuditTrailActivityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformAuditTrailActivityCollector) Init(_ context.Context) error {
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

func (c *TerraformAuditTrailActivityCollector) Start(ctx context.Context) error {
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
		"Starting HCP Terraform audit trail activity collector",
	)
	return fmt.Errorf("terraform audit trail activity collector not implemented")
}

func (c *TerraformAuditTrailActivityCollector) Stop(ctx context.Context) error {
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
