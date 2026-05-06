package entity

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fgrzl/enumerators"
	"github.com/hydn-co/mesh-hashicorp/internal/api"
	"github.com/hydn-co/mesh-hashicorp/internal/collectors"
	"github.com/hydn-co/mesh-hashicorp/internal/credentials"
	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewTerraformAccountEntityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformAccountEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformAccountEntityCollector{TypedFeatureContext: ctx}
}

type TerraformAccountEntityCollector struct {
	state connectorutil.FeatureState
	token string
	*connector.TypedFeatureContext[*options.TerraformAccountEntityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformAccountEntityCollector) Init(_ context.Context) error {
	opts := c.GetOptions()
	if err := connectorutil.Validate(opts, "feature options"); err != nil {
		return err
	}
	token, err := credentials.ExtractToken(c.GetCredentials())
	if err != nil {
		return fmt.Errorf("parse api key credentials: %w", err)
	}
	c.token = token
	c.state.MarkReady()

	return nil
}

func (c *TerraformAccountEntityCollector) Start(ctx context.Context) error {
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
		"Starting HCP Terraform account entity collector",
	)
	opts := c.GetOptions()

	client, err := api.NewClient(http.DefaultClient, "terraform", opts.GetHostname(), c.token)
	if err != nil {
		return fmt.Errorf("build terraform client: %w", err)
	}

	membershipEnum := client.OrganizationMembershipEnumerator(ctx, opts.GetOrganization())
	if err := enumerators.ForEach(membershipEnum, func(result api.TerraformOrganizationMembershipResult) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		userID := result.Membership.Relationships.User.Data.ID
		if userID == "" {
			connectorutil.LogFeature(
				ctx,
				c.TypedFeatureContext,
				slog.LevelWarn,
				"terraform organization membership returned empty user id",
				"membership_id",
				result.Membership.ID,
			)
			return fmt.Errorf(
				"terraform organization membership %s returned empty user id",
				result.Membership.ID,
			)
		}

		user := result.User
		if user.ID == "" {
			user = api.TerraformUser{ID: userID}
		}

		account := collectors.NewTerraformAccount(userID, user, result.Membership.Attributes.Status)
		if err := c.Emit(ctx, account); err != nil {
			return fmt.Errorf("emit account %s: %w", account.AccountRef, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("enumerate organization memberships: %w", err)
	}

	connectorutil.LogFeature(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Finished HCP Terraform account entity collector",
	)
	return nil
}

func (c *TerraformAccountEntityCollector) Stop(ctx context.Context) error {
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
