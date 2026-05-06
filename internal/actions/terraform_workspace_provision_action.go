package actions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-hashicorp/internal/credentials"
	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-hashicorp/internal/payloads"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewTerraformWorkspaceProvisionAction(
	ctx *connector.TypedFeatureContext[*options.TerraformWorkspaceProvisionActionOptions, *payloads.TerraformWorkspaceProvisionPayload],
) runner.Feature {
	return &TerraformWorkspaceProvisionAction{TypedFeatureContext: ctx}
}

type TerraformWorkspaceProvisionAction struct {
	state connectorutil.FeatureState
	token string
	*connector.TypedFeatureContext[*options.TerraformWorkspaceProvisionActionOptions, *payloads.TerraformWorkspaceProvisionPayload]
}

func (a *TerraformWorkspaceProvisionAction) Init(ctx context.Context) error {
	opts := a.GetOptions()
	if err := connectorutil.Validate(opts, "feature options"); err != nil {
		return err
	}
	payload := a.GetPayload()
	if err := connectorutil.Validate(payload, "terraform workspace provision payload"); err != nil {
		return err
	}
	token, err := credentials.ExtractToken(a.GetCredentials())
	if err != nil {
		return fmt.Errorf("parse api key credentials: %w", err)
	}
	a.token = token
	a.state.MarkReady()
	connectorutil.LogFeature(
		ctx,
		a.TypedFeatureContext,
		slog.LevelInfo,
		"Initialized HCP Terraform workspace provision action",
	)
	return nil
}

func (a *TerraformWorkspaceProvisionAction) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := a.state.RequireReady(); err != nil {
		return err
	}
	connectorutil.LogFeature(
		ctx,
		a.TypedFeatureContext,
		slog.LevelInfo,
		"Starting HCP Terraform workspace provision action",
	)
	return fmt.Errorf("terraform workspace provision action not implemented")
}

func (a *TerraformWorkspaceProvisionAction) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := a.state.RequireReady(); err != nil {
		return err
	}
	a.state.Reset()
	a.token = ""
	return nil
}
