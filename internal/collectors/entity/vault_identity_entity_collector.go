package entity

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sort"

	"github.com/fgrzl/enumerators"
	"github.com/hydn-co/mesh-hashicorp/internal/api"
	"github.com/hydn-co/mesh-hashicorp/internal/collectors"
	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewVaultIdentityEntityCollector(
	ctx *connector.TypedFeatureContext[*options.VaultIdentityEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &VaultIdentityEntityCollector{TypedFeatureContext: ctx}
}

type VaultIdentityEntityCollector struct {
	state connectorutil.FeatureState
	token string
	*connector.TypedFeatureContext[*options.VaultIdentityEntityCollectorOptions, *connector.NoPayload]
}

func (c *VaultIdentityEntityCollector) Init(_ context.Context) error {
	opts := c.GetOptions()
	if err := connectorutil.Validate(opts, "feature options"); err != nil {
		return err
	}
	token, err := connectorutil.ExtractAPIKey(c.GetCredentials())
	if err != nil {
		return fmt.Errorf("parse api key credentials: %w", err)
	}
	c.token = token
	c.state.MarkReady()

	return nil
}

func (c *VaultIdentityEntityCollector) Start(ctx context.Context) error {
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
		"Starting Vault identity entity collector",
	)
	opts := c.GetOptions()

	client, err := api.NewVaultClient(
		http.DefaultClient,
		opts.GetAddress(),
		opts.GetNamespace(),
		c.token,
	)
	if err != nil {
		return fmt.Errorf("build vault client: %w", err)
	}

	entityIDs, err := client.ListIdentityEntityIDs(ctx)
	if err != nil {
		return fmt.Errorf("list vault identity entities: %w", err)
	}
	sort.Strings(entityIDs)

	if err := enumerators.ForEach(enumerators.Slice(entityIDs), func(entityID string) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if entityID == "" {
			connectorutil.LogFeature(
				ctx,
				c.TypedFeatureContext,
				slog.LevelWarn,
				"vault identity entity list returned empty entity id",
			)
			return fmt.Errorf("vault identity entity list returned empty entity id")
		}

		entity, err := client.GetIdentityEntity(ctx, entityID)
		if err != nil {
			return fmt.Errorf("read vault identity entity %s: %w", entityID, err)
		}

		account, err := collectors.NewVaultAccount(entity)
		if err != nil {
			return fmt.Errorf("map vault identity entity %s: %w", entityID, err)
		}
		if err := c.Emit(ctx, account); err != nil {
			return fmt.Errorf("emit account %s: %w", account.AccountRef, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("enumerate vault identity entities: %w", err)
	}

	groupIDs, err := client.ListIdentityGroupIDs(ctx)
	if err != nil {
		return fmt.Errorf("list vault identity groups: %w", err)
	}
	sort.Strings(groupIDs)

	if err := enumerators.ForEach(enumerators.Slice(groupIDs), func(groupID string) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if groupID == "" {
			connectorutil.LogFeature(
				ctx,
				c.TypedFeatureContext,
				slog.LevelWarn,
				"vault identity group list returned empty group id",
			)
			return fmt.Errorf("vault identity group list returned empty group id")
		}

		group, err := client.GetIdentityGroup(ctx, groupID)
		if err != nil {
			return fmt.Errorf("read vault identity group %s: %w", groupID, err)
		}

		groupEntity, err := collectors.NewVaultGroup(group)
		if err != nil {
			return fmt.Errorf("map vault identity group %s: %w", groupID, err)
		}
		if err := c.Emit(ctx, groupEntity); err != nil {
			return fmt.Errorf("emit group %s: %w", groupEntity.GroupRef, err)
		}

		if err := enumerators.ForEach(enumerators.Slice(group.MemberEntityIDs), func(memberEntityID string) error {
			if memberEntityID == "" {
				connectorutil.LogFeature(
					ctx,
					c.TypedFeatureContext,
					slog.LevelWarn,
					"vault identity group returned empty member entity id",
					"group_id",
					groupEntity.GroupRef,
				)
				return fmt.Errorf(
					"vault identity group %s returned empty member entity id",
					groupEntity.GroupRef,
				)
			}

			groupMember, err := collectors.NewVaultGroupMember(groupEntity.GroupRef, memberEntityID)
			if err != nil {
				return fmt.Errorf("map vault group member %s:%s: %w", groupEntity.GroupRef, memberEntityID, err)
			}
			if err := c.Emit(ctx, groupMember); err != nil {
				return fmt.Errorf(
					"emit group member %s:%s: %w",
					groupMember.GroupRef,
					groupMember.AccountRef,
					err,
				)
			}

			return nil
		}); err != nil {
			return fmt.Errorf("enumerate vault group members for %s: %w", groupEntity.GroupRef, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("enumerate vault identity groups: %w", err)
	}

	connectorutil.LogFeature(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Finished Vault identity entity collector",
	)
	return nil
}

func (c *VaultIdentityEntityCollector) Stop(ctx context.Context) error {
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
		"Stopping Vault identity entity collector",
	)
	c.state.Reset()
	c.token = ""
	return nil
}
