package actions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-hashicorp/internal/credentials"
	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-hashicorp/internal/payloads"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewTerraformWorkspaceProvisionAction(
	ctx *connector.TypedFeatureContext[*options.TerraformWorkspaceProvisionActionOptions, *payloads.TerraformWorkspaceProvisionPayload],
) runner.Feature {
	return &TerraformWorkspaceProvisionAction{TypedFeatureContext: ctx}
}

type TerraformWorkspaceProvisionAction struct {
	token string
	*connector.TypedFeatureContext[*options.TerraformWorkspaceProvisionActionOptions, *payloads.TerraformWorkspaceProvisionPayload]
}

func (a *TerraformWorkspaceProvisionAction) Init(ctx context.Context) error {
	if err := options.ValidateTerraformOptions(a.GetOptions()); err != nil {
		return err
	}
	if payload := a.GetPayload(); payload == nil {
		return fmt.Errorf("terraform workspace provision payload is required")
	} else if err := payload.Validate(); err != nil {
		return fmt.Errorf("invalid terraform workspace provision payload: %w", err)
	}
	token, err := credentials.ExtractToken(a.GetCredentials())
	if err != nil {
		return fmt.Errorf("parse api key credentials: %w", err)
	}
	a.token = token
	logAction(ctx, a.TypedFeatureContext, slog.LevelInfo, "Initialized HCP Terraform workspace provision action")
	return nil
}

func (a *TerraformWorkspaceProvisionAction) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logAction(ctx, a.TypedFeatureContext, slog.LevelInfo, "Starting HCP Terraform workspace provision action")
	return fmt.Errorf("terraform workspace provision action not implemented")
}

func (a *TerraformWorkspaceProvisionAction) Stop(context.Context) error { return nil }
