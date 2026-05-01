package actions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-hashicorp/internal/payloads"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewTerraformTeamProvisionAction(
	ctx *connector.TypedFeatureContext[*options.TerraformTeamProvisionActionOptions, *payloads.TerraformTeamProvisionPayload],
) runner.Feature {
	return &TerraformTeamProvisionAction{TypedFeatureContext: ctx}
}

type TerraformTeamProvisionAction struct {
	*connector.TypedFeatureContext[*options.TerraformTeamProvisionActionOptions, *payloads.TerraformTeamProvisionPayload]
}

func (a *TerraformTeamProvisionAction) Init(ctx context.Context) error {
	if err := options.ValidateTerraformOptions(a.GetOptions()); err != nil {
		return err
	}
	if payload := a.GetPayload(); payload == nil {
		return fmt.Errorf("terraform team provision payload is required")
	} else if err := payload.Validate(); err != nil {
		return fmt.Errorf("invalid terraform team provision payload: %w", err)
	}
	logAction(ctx, a.TypedFeatureContext, slog.LevelInfo, "Initialized HCP Terraform team provision action")
	return nil
}

func (a *TerraformTeamProvisionAction) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logAction(ctx, a.TypedFeatureContext, slog.LevelInfo, "Starting HCP Terraform team provision action")
	return fmt.Errorf("terraform team provision action not implemented")
}

func (a *TerraformTeamProvisionAction) Stop(context.Context) error { return nil }
