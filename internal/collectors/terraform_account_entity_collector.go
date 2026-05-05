package collectors

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fgrzl/enumerators"
	"github.com/hydn-co/mesh-hashicorp/internal/api"
	"github.com/hydn-co/mesh-hashicorp/internal/credentials"
	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewTerraformAccountEntityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformAccountEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformAccountEntityCollector{TypedFeatureContext: ctx}
}

type TerraformAccountEntityCollector struct {
	token string
	*connector.TypedFeatureContext[*options.TerraformAccountEntityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformAccountEntityCollector) Init(_ context.Context) error {
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

func (c *TerraformAccountEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logCollector(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting HCP Terraform account entity collector")

	client, err := newTerraformClient(c.GetOptions().GetHostname(), c.token)
	if err != nil {
		return fmt.Errorf("build terraform client: %w", err)
	}

	membershipEnum := client.OrganizationMembershipEnumerator(ctx, c.GetOptions().GetOrganization())
	if err := enumerators.ForEach(membershipEnum, func(result api.TerraformOrganizationMembershipResult) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		userID := result.Membership.Relationships.User.Data.ID
		if userID == "" {
			return nil
		}

		user := result.User
		if user.ID == "" {
			user = api.TerraformUser{ID: userID}
		}

		account := newTerraformAccount(userID, user, result.Membership.Attributes.Status)
		if err := c.Emit(ctx, account); err != nil {
			return fmt.Errorf("emit account %s: %w", account.AccountRef, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("enumerate organization memberships: %w", err)
	}

	logCollector(ctx, c.TypedFeatureContext, slog.LevelInfo, "Finished HCP Terraform account entity collector")
	return nil
}

func (c *TerraformAccountEntityCollector) Stop(context.Context) error { return nil }
