# mesh-hashicorp

Mesh connector for HashiCorp platforms.

Current bootstrap scope:

- HCP Terraform entity collectors
- HCP Terraform audit activity collector
- HCP Terraform provisioning actions
- Vault entity collectors
- Vault audit deferred

The repository is bootstrapped to match the current modular connector structure used by `mesh-azure`:

- one option type per file under `internal/options/`
- one payload type per file under `internal/payloads/`
- one feature per file under `internal/collectors/` and `internal/actions/`

## Build

```bash
go build ./...
go test ./...
```

## Manifest

```bash
go run ./cmd/... -describe
go run ./cmd/... -list
```