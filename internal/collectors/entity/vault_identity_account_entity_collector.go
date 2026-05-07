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

func NewVaultIdentityAccountEntityCollector(
	ctx *connector.TypedFeatureContext[*options.VaultIdentityAccountEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &VaultIdentityAccountEntityCollector{TypedFeatureContext: ctx}
}

type VaultIdentityAccountEntityCollector struct {
	state connectorutil.FeatureState
	token string
	*connector.TypedFeatureContext[*options.VaultIdentityAccountEntityCollectorOptions, *connector.NoPayload]
}

func (c *VaultIdentityAccountEntityCollector) Init(_ context.Context) error {
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

func (c *VaultIdentityAccountEntityCollector) Start(ctx context.Context) error {
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
		"Starting Vault identity account entity collector",
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

	connectorutil.LogFeature(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Finished Vault identity account entity collector",
	)
	return nil
}

func (c *VaultIdentityAccountEntityCollector) Stop(ctx context.Context) error {
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
		"Stopping Vault identity account entity collector",
	)
	c.state.Reset()
	c.token = ""
	return nil
}
