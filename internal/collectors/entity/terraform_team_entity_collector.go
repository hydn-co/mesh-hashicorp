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
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewTerraformTeamEntityCollector(
	ctx *connector.TypedFeatureContext[*options.TerraformTeamEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &TerraformTeamEntityCollector{TypedFeatureContext: ctx}
}

type TerraformTeamEntityCollector struct {
	initialized bool
	token       string
	*connector.TypedFeatureContext[*options.TerraformTeamEntityCollectorOptions, *connector.NoPayload]
}

func (c *TerraformTeamEntityCollector) Init(_ context.Context) error {
	opts := c.GetOptions()
	if err := options.ValidateTerraformOptions(opts); err != nil {
		return err
	}
	token, err := credentials.ExtractToken(c.GetCredentials())
	if err != nil {
		return fmt.Errorf("parse api key credentials: %w", err)
	}
	c.token = token
	c.initialized = true

	return nil
}

func (c *TerraformTeamEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !c.initialized {
		return fmt.Errorf("terraform team entity collector not initialized")
	}
	collectors.LogCollector(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Starting HCP Terraform team entity collector",
	)
	opts := c.GetOptions()

	client, err := api.NewClient(http.DefaultClient, "terraform", opts.GetHostname(), c.token)
	if err != nil {
		return fmt.Errorf("build terraform client: %w", err)
	}

	teamEnum := client.TeamEnumerator(ctx, opts.GetOrganization())
	if err := enumerators.ForEach(teamEnum, func(result api.TerraformTeamResult) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if result.Team.ID == "" {
			collectors.LogCollector(
				ctx,
				c.TypedFeatureContext,
				slog.LevelWarn,
				"terraform team result returned empty team id",
				"team_name",
				result.Team.Attributes.Name,
			)
			return fmt.Errorf("terraform team result returned empty team id")
		}

		group := collectors.NewTerraformGroup(result.Team)
		if err := c.Emit(ctx, group); err != nil {
			return fmt.Errorf("emit group %s: %w", group.GroupRef, err)
		}

		if err := enumerators.ForEach(
			enumerators.Slice(result.Team.Relationships.Users.Data),
			func(member api.TerraformResourceIdentifier) error {
				if member.ID == "" {
					collectors.LogCollector(
						ctx,
						c.TypedFeatureContext,
						slog.LevelWarn,
						"terraform team relationship returned empty user id",
						"team_id",
						result.Team.ID,
					)
					return fmt.Errorf(
						"terraform team %s relationship returned empty user id",
						result.Team.ID,
					)
				}

				user, ok := result.UsersByID[member.ID]
				if !ok {
					user = api.TerraformUser{ID: member.ID}
				}

				groupMember := collectors.NewTerraformGroupMember(result.Team.ID, member.ID, user)
				if err := c.Emit(ctx, groupMember); err != nil {
					return fmt.Errorf(
						"emit group member %s:%s: %w",
						groupMember.GroupRef,
						groupMember.AccountRef,
						err,
					)
				}

				return nil
			},
		); err != nil {
			return fmt.Errorf("enumerate team members for %s: %w", group.GroupRef, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("enumerate teams: %w", err)
	}

	collectors.LogCollector(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Finished HCP Terraform team entity collector",
	)
	return nil
}

func (c *TerraformTeamEntityCollector) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !c.initialized {
		return fmt.Errorf("terraform team entity collector not initialized")
	}
	c.initialized = false
	c.token = ""
	return nil
}
