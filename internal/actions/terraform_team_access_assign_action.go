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

func NewTerraformTeamAccessAssignAction(
	ctx *connector.TypedFeatureContext[*options.TerraformTeamAccessAssignActionOptions, *payloads.TerraformTeamAccessAssignPayload],
) runner.Feature {
	return &TerraformTeamAccessAssignAction{TypedFeatureContext: ctx}
}

type TerraformTeamAccessAssignAction struct {
	initialized bool
	token       string
	*connector.TypedFeatureContext[*options.TerraformTeamAccessAssignActionOptions, *payloads.TerraformTeamAccessAssignPayload]
}

func (a *TerraformTeamAccessAssignAction) Init(ctx context.Context) error {
	opts := a.GetOptions()
	if err := options.ValidateTerraformOptions(opts); err != nil {
		return err
	}
	payload := a.GetPayload()
	if err := payloads.ValidatePayload(payload, "terraform team access assign payload"); err != nil {
		return err
	}
	token, err := credentials.ExtractToken(a.GetCredentials())
	if err != nil {
		return fmt.Errorf("parse api key credentials: %w", err)
	}
	a.token = token
	a.initialized = true
	logAction(ctx, a.TypedFeatureContext, slog.LevelInfo, "Initialized HCP Terraform team access assign action")
	return nil
}

func (a *TerraformTeamAccessAssignAction) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !a.initialized {
		return fmt.Errorf("terraform team access assign action not initialized")
	}
	logAction(ctx, a.TypedFeatureContext, slog.LevelInfo, "Starting HCP Terraform team access assign action")
	return fmt.Errorf("terraform team access assign action not implemented")
}

func (a *TerraformTeamAccessAssignAction) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !a.initialized {
		return fmt.Errorf("terraform team access assign action not initialized")
	}
	a.initialized = false
	a.token = ""
	return nil
}
