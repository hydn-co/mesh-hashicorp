package collectors

import (
	"context"
	"fmt"
	"log/slog"

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
	*connector.TypedFeatureContext[*options.TerraformAuditTrailActivityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformAuditTrailActivityCollector) Init(_ context.Context) error {
	return options.ValidateTerraformOptions(c.GetOptions())
}

func (c *TerraformAuditTrailActivityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logCollector(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting HCP Terraform audit trail activity collector")
	return fmt.Errorf("terraform audit trail activity collector not implemented")
}

func (c *TerraformAuditTrailActivityCollector) Stop(context.Context) error { return nil }
