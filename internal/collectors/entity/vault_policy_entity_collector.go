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

func NewVaultPolicyEntityCollector(
	ctx *connector.TypedFeatureContext[*options.VaultPolicyEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &VaultPolicyEntityCollector{TypedFeatureContext: ctx}
}

type VaultPolicyEntityCollector struct {
	state connectorutil.FeatureState
	token string
	*connector.TypedFeatureContext[*options.VaultPolicyEntityCollectorOptions, *connector.NoPayload]
}

func (c *VaultPolicyEntityCollector) Init(_ context.Context) error {
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

func (c *VaultPolicyEntityCollector) Start(ctx context.Context) error {
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
		"Starting Vault policy entity collector",
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

	policyNames, err := client.ListPolicyNames(ctx)
	if err != nil {
		return fmt.Errorf("list vault policies: %w", err)
	}
	sort.Strings(policyNames)

	if err := enumerators.ForEach(enumerators.Slice(policyNames), func(policyName string) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if policyName == "" {
			connectorutil.LogFeature(
				ctx,
				c.TypedFeatureContext,
				slog.LevelWarn,
				"vault policy list returned empty policy name",
			)
			return fmt.Errorf("vault policy list returned empty policy name")
		}

		policyEntity, err := collectors.NewVaultPolicy(policyName)
		if err != nil {
			return fmt.Errorf("map vault policy %s: %w", policyName, err)
		}
		if err := c.Emit(ctx, policyEntity); err != nil {
			return fmt.Errorf("emit policy %s: %w", policyEntity.PolicyRef, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("enumerate vault policies: %w", err)
	}

	connectorutil.LogFeature(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Finished Vault policy entity collector",
	)
	return nil
}

func (c *VaultPolicyEntityCollector) Stop(ctx context.Context) error {
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
		"Stopping Vault policy entity collector",
	)
	c.state.Reset()
	c.token = ""
	return nil
}
