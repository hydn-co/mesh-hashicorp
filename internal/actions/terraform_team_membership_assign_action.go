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

func NewTerraformTeamMembershipAssignAction(
	ctx *connector.TypedFeatureContext[*options.TerraformTeamMembershipAssignActionOptions, *payloads.TerraformTeamMembershipAssignPayload],
) runner.Feature {
	return &TerraformTeamMembershipAssignAction{TypedFeatureContext: ctx}
}

type TerraformTeamMembershipAssignAction struct {
	*connector.TypedFeatureContext[*options.TerraformTeamMembershipAssignActionOptions, *payloads.TerraformTeamMembershipAssignPayload]
}

func (a *TerraformTeamMembershipAssignAction) Init(ctx context.Context) error {
	if err := options.ValidateTerraformOptions(a.GetOptions()); err != nil {
		return err
	}
	if payload := a.GetPayload(); payload == nil {
		return fmt.Errorf("terraform team membership assign payload is required")
	} else if err := payload.Validate(); err != nil {
		return fmt.Errorf("invalid terraform team membership assign payload: %w", err)
	}
	logAction(ctx, a.TypedFeatureContext, slog.LevelInfo, "Initialized HCP Terraform team membership assign action")
	return nil
}

func (a *TerraformTeamMembershipAssignAction) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logAction(ctx, a.TypedFeatureContext, slog.LevelInfo, "Starting HCP Terraform team membership assign action")
	return fmt.Errorf("terraform team membership assign action not implemented")
}

func (a *TerraformTeamMembershipAssignAction) Stop(context.Context) error { return nil }
