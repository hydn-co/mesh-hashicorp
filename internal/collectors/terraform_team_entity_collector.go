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

func NewTerraformTeamEntityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformTeamEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformTeamEntityCollector{TypedFeatureContext: ctx}
}

type TerraformTeamEntityCollector struct {
	token string
	*connector.TypedFeatureContext[*options.TerraformTeamEntityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformTeamEntityCollector) Init(_ context.Context) error {
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

func (c *TerraformTeamEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	logCollector(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting HCP Terraform team entity collector")

	client, err := newTerraformClient(c.GetOptions().GetHostname(), c.token)
	if err != nil {
		return fmt.Errorf("build terraform client: %w", err)
	}

	teamEnum := client.TeamEnumerator(ctx, c.GetOptions().GetOrganization())
	if err := enumerators.ForEach(teamEnum, func(result api.TerraformTeamResult) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		group := newTerraformGroup(result.Team)
		if err := c.Emit(ctx, group); err != nil {
			return fmt.Errorf("emit group %s: %w", group.GroupRef, err)
		}

		for _, member := range result.Team.Relationships.Users.Data {
			if member.ID == "" {
				continue
			}

			user, ok := result.UsersByID[member.ID]
			if !ok {
				user = api.TerraformUser{ID: member.ID}
			}

			groupMember := newTerraformGroupMember(result.Team.ID, member.ID, user)
			if err := c.Emit(ctx, groupMember); err != nil {
				return fmt.Errorf(
					"emit group member %s:%s: %w",
					groupMember.GroupRef,
					groupMember.AccountRef,
					err,
				)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("enumerate teams: %w", err)
	}

	logCollector(ctx, c.TypedFeatureContext, slog.LevelInfo, "Finished HCP Terraform team entity collector")
	return nil
}

func (c *TerraformTeamEntityCollector) Stop(context.Context) error { return nil }
