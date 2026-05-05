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
	"github.com/hydn-co/mesh-hashicorp/internal/credentials"
	"github.com/hydn-co/mesh-hashicorp/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func NewVaultAuthMethodEntityCollector(
	ctx *connector.TypedFeatureContext[*options.VaultAuthMethodEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &VaultAuthMethodEntityCollector{TypedFeatureContext: ctx}
}

type VaultAuthMethodEntityCollector struct {
	initialized bool
	token       string
	*connector.TypedFeatureContext[*options.VaultAuthMethodEntityCollectorOptions, *connector.NoPayload]
}

func (c *VaultAuthMethodEntityCollector) Init(_ context.Context) error {
	opts := c.GetOptions()
	if err := options.ValidateVaultOptions(opts); err != nil {
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

func (c *VaultAuthMethodEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !c.initialized {
		return fmt.Errorf("vault auth method entity collector not initialized")
	}
	collectors.LogCollector(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Starting Vault auth method entity collector",
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

	authMethods, err := client.ListAuthMethods(ctx)
	if err != nil {
		return fmt.Errorf("list vault auth methods: %w", err)
	}

	paths := make([]string, 0, len(authMethods))
	for path := range authMethods {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	if err := enumerators.ForEach(enumerators.Slice(paths), func(path string) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if path == "" {
			collectors.LogCollector(
				ctx,
				c.TypedFeatureContext,
				slog.LevelWarn,
				"vault auth method list returned empty mount path",
			)
			return fmt.Errorf("vault auth method list returned empty mount path")
		}

		application, err := collectors.NewVaultApplication(path, authMethods[path])
		if err != nil {
			return fmt.Errorf("map vault auth method %s: %w", path, err)
		}
		if err := c.Emit(ctx, application); err != nil {
			return fmt.Errorf("emit application %s: %w", application.ApplicationRef, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("enumerate vault auth methods: %w", err)
	}

	collectors.LogCollector(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Finished Vault auth method entity collector",
	)
	return nil
}

func (c *VaultAuthMethodEntityCollector) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !c.initialized {
		return fmt.Errorf("vault auth method entity collector not initialized")
	}
	collectors.LogCollector(
		ctx,
		c.TypedFeatureContext,
		slog.LevelInfo,
		"Stopping Vault auth method entity collector",
	)
	c.initialized = false
	c.token = ""
	return nil
}
