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

func NewTerraformTeamAccessAssignAction(
	ctx *connector.TypedFeatureContext[*options.TerraformTeamAccessAssignActionOptions, *payloads.TerraformTeamAccessAssignPayload],
) runner.Feature {
	return &TerraformTeamAccessAssignAction{TypedFeatureContext: ctx}
}

type TerraformTeamAccessAssignAction struct {
	*connector.TypedFeatureContext[*options.TerraformTeamAccessAssignActionOptions, *payloads.TerraformTeamAccessAssignPayload]
}

func (a *TerraformTeamAccessAssignAction) Init(ctx context.Context) error {
	if err := options.ValidateTerraformOptions(a.GetOptions()); err != nil {
		return err
	}
	if payload := a.GetPayload(); payload == nil {
		return fmt.Errorf("terraform team access assign payload is required")
	} else if err := payload.Validate(); err != nil {
		return fmt.Errorf("invalid terraform team access assign payload: %w", err)
	}
	logAction(ctx, a.TypedFeatureContext, slog.LevelInfo, "Initialized HCP Terraform team access assign action")
	return nil
}

func (a *TerraformTeamAccessAssignAction) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logAction(ctx, a.TypedFeatureContext, slog.LevelInfo, "Starting HCP Terraform team access assign action")
	return fmt.Errorf("terraform team access assign action not implemented")
}

func (a *TerraformTeamAccessAssignAction) Stop(context.Context) error { return nil }
