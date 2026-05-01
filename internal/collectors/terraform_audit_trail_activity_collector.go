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

func NewTerraformAuditTrailActivityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformAuditTrailActivityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformAuditTrailActivityCollector{TypedFeatureContext: ctx}
}

type TerraformAuditTrailActivityCollector struct {
	token string
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

	return nil
}

func (c *TerraformAuditTrailActivityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logCollector(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting HCP Terraform audit trail activity collector")
	return fmt.Errorf("terraform audit trail activity collector not implemented")
}

func (c *TerraformAuditTrailActivityCollector) Stop(context.Context) error { return nil }
