# Agent instructions – mesh-hashicorp

This file orients AI agents and automated tooling working in the `mesh-hashicorp` repository.

## Before any task

- Read this file first.
- Read the feature you are modifying, or the closest existing feature, before changing code.
- If you touch capability registration or schemas, read `cmd/main.go`, `internal/options/core.go`, `internal/options/register.go`, and `internal/payloads/register.go` first.
- Prefer existing repo and `mesh-sdk` APIs over one-off wrappers.

## Project summary

- `mesh-hashicorp`: A Mesh platform connector for HashiCorp products.
- Entry point: `cmd/main.go`
- Framework: `mesh-sdk`
- Language: Go (see `go.mod`)

## Non-negotiable rules

1. All collectors and actions embed `*connector.TypedFeatureContext[...]` and implement `Init`, `Start`, and `Stop`.
2. Collectors are split by feature type under `internal/collectors/entity/` and `internal/collectors/activity/`. Shared collector helpers stay in `internal/collectors/` only when used by both.
3. All option types live in `internal/options/` with one type per file plus shared `core.go`, `register.go`, and validation helpers.
4. All payload types live in `internal/payloads/` with one type per file plus `register.go`.
5. Use only `net/http` for provider API calls.
6. Wrap errors with `fmt.Errorf("context: %w", err)`.
7. Validate options in `Init()` and validate action payloads in action `Init()`.
8. Use `testkit.TestPolymorphicRegistrations()` for option and payload registration coverage.
9. Use behavioral test names: `TestShould{Expectation}When{Condition}`.

## Primary sources

| Need | Source |
|------|--------|
| Manifest and feature registration | `cmd/main.go` |
| Options | `internal/options/` |
| Payloads | `internal/payloads/` |
| Collectors | `internal/collectors/entity/`, `internal/collectors/activity/`, and shared helpers in `internal/collectors/` |
| Actions | `internal/actions/` |
| Credentials | `internal/credentials/` |
| SDK framework | `mesh-sdk/pkg/runner`, `mesh-sdk/pkg/connector`, `mesh-sdk/pkg/testkit` |