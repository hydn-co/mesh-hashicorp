# mesh-hashicorp Implementation Plan

## Objective

Bootstrap `mesh-hashicorp` from an empty repository into a production-ready Mesh connector that follows the bootstrap skill, matches current `mesh-azure` and `mesh-slack` conventions, and uses a single-source-of-truth implementation style.

Because HashiCorp is a vendor umbrella rather than a single API surface, the connector should be structured as one repository with product-scoped features. The initial implementation should avoid pretending that Vault, HCP Terraform, Boundary, and HCP generic services can be delivered in one pass.

## Scope Decision

### Repository Scope

- Repository name remains `mesh-hashicorp`.
- Feature names must be product-scoped, for example `hashicorp_terraform_*` and `hashicorp_vault_*`.
- Shared code lives once under `internal/api`, `internal/options`, `internal/credentials`, and shared logging helpers.
- Product-specific behavior stays isolated in distinct client helpers, options cores, mappings, and collector/action files.

### Phase Strategy

Phase 1 should ship one vertically complete product slice rather than shallow support across multiple HashiCorp products.

Recommended sequencing:

1. Phase 0: bootstrap repository skeleton and shared plumbing.
2. Phase 1: HCP Terraform / Terraform Cloud entity, activity, and action coverage plus Vault entity coverage.
3. Phase 2: Highest-value Vault write actions.
4. Phase 3: Vault activity plus Boundary or additional HCP services once the first two slices are stable.

## Client Strategy

### Decision

Use raw `net/http` for provider API interactions.

### Why

This matches the bootstrap skill's explicit rule, keeps the connector consistent with `mesh-azure` and `mesh-slack`, and avoids carrying multiple client abstractions for unrelated HashiCorp product families.

### SDK Evaluation

- `github.com/hashicorp/vault/api` latest resolved version: `v1.23.0`
- `github.com/hashicorp/vault/api` measured non-stdlib dependency count: `37`
- `github.com/hashicorp/go-tfe` latest resolved version: `v1.103.0`
- `github.com/hashicorp/go-tfe` measured non-stdlib dependency count: `14`
- `github.com/hashicorp/hcp-sdk-go` latest resolved version: `v0.172.0`

### Recommendation

- Reject `vault/api` for connector runtime use. Its dependency footprint is too large for a small connector binary, and it hides behavior such as Vault-specific headers that we can implement directly.
- Reject `go-tfe` for phase 1 even though it is lighter. HCP Terraform uses a predictable JSON:API surface with straightforward pagination, so raw HTTP is simpler, keeps behavior explicit, and avoids diverging from the bootstrap skill.
- Reject `hcp-sdk-go` for initial work. It is an umbrella SDK for a broad service estate and is the wrong abstraction for a narrowly scoped connector.

## Product Assessment

### HCP Terraform / Terraform Cloud

This is the best Phase 1 target because it has the strongest combination of:

- identity and access resources
- auditable events through audit trails
- write APIs suitable for action features
- a single bearer-token HTTP model
- stable JSON:API pagination and relationship semantics

Phase 1 should treat HCP Terraform as the first complete vertical slice while also shipping Vault entity coverage in the same milestone.

### Vault

Vault entity coverage belongs in the first milestone alongside HCP Terraform, because identity, group, policy, and auth method inventory are high-value and available through the HTTP API.

Vault activity should not be treated as blocked. The key distinction is only that Vault does not expose a first-party paginated event list comparable to HCP Terraform audit trails. Instead, it writes line-delimited JSON audit entries to configured audit devices and sinks. That makes Vault activity a source-adapter problem, not a collector-pattern problem.

### Boundary

Boundary is a reasonable later candidate, but it should be deferred until the shared HashiCorp repo skeleton, options model, and HTTP client helpers are proven with HCP Terraform and Vault.

## Phase 0: Repository Bootstrap

Create the full skeleton expected by the bootstrap skill, but align package structure with the current modular `mesh-azure` layout rather than the older monolithic `options.go` and `payloads.go` examples described in the prompt.

### Files To Create

- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`
- `.github/workflows/version.yml`
- `.github/dependabot.yml`
- `.gitignore`
- `.golangci.yml`
- `AGENTS.md`
- `GitVersion.yml`
- `README.md`
- `lefthook.yml`
- `cmd/main.go`
- `internal/api/`
- `internal/collectors/logging.go`
- `internal/actions/logging.go`
- `internal/credentials/credentials.go`
- `internal/options/core.go`
- `internal/options/register.go`
- `internal/options/*_options.go`
- `internal/options/options_test.go`
- `internal/payloads/register.go`
- `internal/payloads/*_payload.go`
- `internal/payloads/payloads_test.go`
- `internal/mappings/context.go`
- `internal/mappings/helpers_*.go`
- `internal/mappings/*_mappings.go`
- `internal/mappings/*_mappings_test.go`

### Structural Alignment With `mesh-azure`

Investigation of the refactored `mesh-azure` shows three patterns worth copying directly:

- `internal/options/` uses `core.go`, `register.go`, and one option type per single-purpose file such as `user_entity_collector_options.go`.
- `internal/mappings/` is split by activity domain and concern, with shared context and helper files instead of one large mapper.
- `internal/payloads/` now uses `register.go`, one payload type per single-purpose file such as `user_provision_payload.go`, and centralized `payloads_test.go` coverage.

Use that as the repo standard.

### Bootstrap Rules

- Follow the exact manifest registration shape used in `mesh-azure` and `mesh-slack`.
- Keep shared option cores in `internal/options/core.go`.
- Keep polymorphic option registration in `internal/options/register.go`.
- Keep exactly one option type per `internal/options/*_options.go` file.
- Keep payload registration in `internal/payloads/register.go`.
- Keep exactly one action payload type per `internal/payloads/*_payload.go` file.
- Keep activity mapping logic in `internal/mappings/` from the start, split by domain and concern.
- Use `runner.APIKeyCredential` for both HCP Terraform and Vault token-based features.
- Add shared options cores instead of duplicating base URL, organization, namespace, or token-adjacent configuration.

### Shared Option Cores

Plan for two shared cores in `internal/options/core.go`:

```go
type TerraformOptionsCore struct {
	Hostname     string `json:"hostname" title:"Hostname" description:"HCP Terraform or Terraform Enterprise hostname." binding:"required"`
	Organization string `json:"organization" title:"Organization" description:"Terraform organization name." binding:"required"`
}

type VaultOptionsCore struct {
	Address   string `json:"address" title:"Address" description:"Vault base URL." binding:"required"`
	Namespace string `json:"namespace,omitempty" title:"Namespace" description:"Vault namespace when required."`
}
```

### Shared API Helpers

Bootstrap shared HTTP concerns once:

- request construction
- auth header injection
- user agent
- pagination helpers
- rate-limit and retry handling
- JSON decode helpers
- provider error envelope parsing

### Package Layout Rules

Use modular, single-purpose files throughout the connector:

- `internal/options/`
- `core.go` for shared option cores only.
- `register.go` for polymorphic registrations only.
- one `*_options.go` file per collector or action option type.
- `options_test.go` for centralized polymorphic and requirements coverage.

- `internal/payloads/`
- `register.go` for polymorphic registrations only.
- one `*_payload.go` file per action payload type.
- `payloads_test.go` for centralized polymorphic registration and validation coverage.

- `internal/mappings/`
- `context.go` for mapper context only.
- `helpers_*.go` for shared extraction and normalization helpers.
- one `*_mappings.go` file per provider activity domain or concern.
- one focused `*_mappings_test.go` file per mapping area.
- mapper context should create deterministic event metadata via `types.NewEventMetadataForSource(...)` using typed event instances, matching the current `mesh-azure` activity mapping pattern.

- `internal/collectors/` and `internal/actions/`
- one feature implementation per file.
- shared logging helpers stay in `logging.go` only.

## Phase 1: HCP Terraform Vertical Slice

### Goal

Deliver the first release-worthy milestone with full HCP Terraform entity, activity, and action coverage plus Vault entity coverage.

### Authentication Model

- Use bearer token authentication.
- Use organization tokens for entity collectors and actions.
- Use audit trails tokens for audit collectors.
- Make token expectations explicit in option descriptions and README examples.

### Planned Entity Collectors

#### `hashicorp_terraform_account_entity_collector`

Collect organization users and emit:

- `Account`

Notes:

- Do not force `Person` or `Employee` until the API exposes enough stable profile data to justify it.
- Prefer included relationship expansion when it reduces round trips without obscuring logic.

#### `hashicorp_terraform_team_entity_collector`

Collect teams and memberships and emit:

- `Group`
- `GroupMember`

#### `hashicorp_terraform_workspace_entity_collector`

Collect workspaces and emit:

- `Application`

Model Terraform workspaces as applications.

#### `hashicorp_terraform_policy_entity_collector`

Collect policy sets and policies and emit:

- `Policy`
- `SecurityConfiguration` where a policy-set style record clearly behaves like a policy container

#### `hashicorp_terraform_team_access_entity_collector`

Collect team-to-workspace access and emit:

- `Permission`
- `GroupPermission`

Model workspace access as permission assignment rather than role assignment.

### Planned Activity Collector

#### `hashicorp_terraform_audit_trail_activity_collector`

Collect audit trail events with `FeatureResumeBehaviorLastActivity`.

Initial event targets:

- team lifecycle and membership changes
- workspace lifecycle and permission changes
- policy changes
- token or organization administrative actions where the API exposes them

Resume strategy:

- primary watermark: audit event timestamp
- boundary tiebreaker: audit event ID
- skip boundary event on resume to avoid duplication

Mapping placement:

- keep collector focused on enumeration and emission only
- start with a modular mapping package rather than inline translation
- use `internal/mappings/context.go` plus Terraform-specific mapping files such as `terraform_audit_team_mappings.go`, `terraform_audit_workspace_mappings.go`, and shared helpers when concerns split naturally

### Planned Actions

#### `hashicorp_terraform_team_provision_action`

Create a team in an organization.

#### `hashicorp_terraform_workspace_provision_action`

Create a workspace in an organization.

#### `hashicorp_terraform_team_membership_assign_action`

Assign an existing user to a team.

#### `hashicorp_terraform_team_access_assign_action`

Grant a team access to a workspace.

### Phase 1 File Shape

Use flat packages with product-prefixed files, for example:

- `internal/api/terraform_client.go`
- `internal/api/terraform_pagination.go`
- `internal/api/terraform_models.go`
- `internal/options/core.go`
- `internal/options/register.go`
- `internal/options/terraform_account_entity_collector_options.go`
- `internal/options/terraform_team_entity_collector_options.go`
- `internal/options/terraform_workspace_entity_collector_options.go`
- `internal/options/terraform_policy_entity_collector_options.go`
- `internal/options/terraform_team_access_entity_collector_options.go`
- `internal/options/terraform_audit_trail_activity_collector_options.go`
- `internal/options/terraform_team_provision_action_options.go`
- `internal/options/terraform_workspace_provision_action_options.go`
- `internal/options/terraform_team_membership_assign_action_options.go`
- `internal/options/terraform_team_access_assign_action_options.go`
- `internal/collectors/terraform_account_entity_collector.go`
- `internal/collectors/terraform_team_entity_collector.go`
- `internal/collectors/terraform_workspace_entity_collector.go`
- `internal/collectors/terraform_policy_entity_collector.go`
- `internal/collectors/terraform_team_access_entity_collector.go`
- `internal/collectors/terraform_audit_trail_activity_collector.go`
- `internal/payloads/register.go`
- `internal/payloads/terraform_team_provision_payload.go`
- `internal/payloads/terraform_workspace_provision_payload.go`
- `internal/payloads/terraform_team_membership_assign_payload.go`
- `internal/payloads/terraform_team_access_assign_payload.go`
- `internal/actions/terraform_team_provision_action.go`
- `internal/actions/terraform_workspace_provision_action.go`
- `internal/actions/terraform_team_membership_assign_action.go`
- `internal/actions/terraform_team_access_assign_action.go`
- `internal/mappings/context.go`
- `internal/mappings/helpers_jsonapi.go`
- `internal/mappings/terraform_audit_team_mappings.go`
- `internal/mappings/terraform_audit_workspace_mappings.go`
- `internal/mappings/terraform_audit_policy_mappings.go`

### Phase 1 Validation

- `go mod tidy`
- `go build ./...`
- `go vet ./...`
- `go test ./...`
- `golangci-lint run`
- `go run ./cmd/... describe`

## Phase 2: Vault Vertical Slice

### Goal

Add the highest-value Vault write actions on top of the Vault entity coverage already delivered in the first milestone.

### Planned Entity Collectors

#### `hashicorp_vault_identity_account_entity_collector`

Collect Vault identity entities and emit:

- `Account`

#### `hashicorp_vault_identity_group_entity_collector`

Collect Vault identity groups and memberships and emit:

- `Group`
- `GroupMember`

#### `hashicorp_vault_policy_entity_collector`

Collect ACL and related policies and emit:

- `Policy`

#### `hashicorp_vault_auth_method_entity_collector`

Collect enabled auth methods and emit:

- `Application`

Model Vault auth methods as applications.

### Planned Actions

#### `hashicorp_vault_group_provision_action`

Create an identity group.

#### `hashicorp_vault_group_membership_assign_action`

Add an entity to a group.

#### `hashicorp_vault_policy_write_action`

Write or update a policy only if the Mesh use case needs a provisioning-style response action.

### Vault Activity Decision

Vault activity is intentionally deferred until after the initial implementation slices. The collector pattern remains valid, but the first shipped implementation should prioritize HCP Terraform coverage, Vault entity coverage, and Vault write actions before adding an audit-source adapter for Vault file, socket, or syslog audit outputs.

## Phase 3: Optional Expansion

After HCP Terraform and Vault stabilize, evaluate:

- Boundary users, groups, roles, targets, and sessions
- HCP platform organization-level IAM data only if there is a coherent API slice worth owning

Do not add additional product families until the repo already has:

- one proven shared HTTP stack
- stable option registration tests
- a clear pattern for product-scoped feature naming

## Implementation Order

1. Create repository skeleton and CI files.
2. Add `AGENTS.md`, `README.md`, and manifest bootstrap.
3. Add credentials and shared options cores.
4. Build shared HTTP helpers.
5. Implement HCP Terraform entity collectors.
6. Implement HCP Terraform audit collector and mappings.
7. Implement HCP Terraform actions.
8. Implement Vault entity collectors.
9. Validate the first milestone: HCP Terraform plus Vault entities.
10. Implement Vault write actions in a second pass.
11. Return to Vault audit after the core product slices are stable.

## Open Questions To Resolve Early

1. Which Vault write action should be prioritized immediately after entity coverage: group provisioning, membership assignment, or policy write?
2. Which Vault audit ingestion mode should ship first when the deferred activity slice starts: mounted file tailing or a socket-based reader?

## Exit Criteria For The First Real Milestone

The first milestone is complete when the repository is bootstrapped and the Phase 1 scope compiles, validates, and passes the standard connector checks:

- HCP Terraform entity collectors
- HCP Terraform activity collector
- HCP Terraform actions
- Vault entity collectors

- `go mod tidy` produces no changes
- `go build ./...` succeeds
- `go vet ./...` succeeds
- `go test ./...` succeeds
- `golangci-lint run` succeeds
- `manifest.Validate()` succeeds
- `go run ./cmd/... describe` emits a valid manifest