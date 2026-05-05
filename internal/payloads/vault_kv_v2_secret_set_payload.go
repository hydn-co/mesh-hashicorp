package payloads

import "fmt"

type VaultKVV2SecretSetPayload struct {
	MountPath  string         `json:"mount_path"    binding:"required"`
	SecretPath string         `json:"secret_path"   binding:"required"`
	Data       map[string]any `json:"data"          binding:"required"`
	CAS        *int           `json:"cas,omitempty"`
}

func (p *VaultKVV2SecretSetPayload) GetDiscriminator() string {
	return "mesh://hashicorp/actions/vault_kv_v2_secret_set_payload"
}

func (p *VaultKVV2SecretSetPayload) Validate() error {
	if p == nil {
		return fmt.Errorf("payload is required")
	}
	if err := validateVaultSecretSetPayloadFields(p.MountPath, p.SecretPath, p.Data); err != nil {
		return err
	}
	if p.CAS != nil && *p.CAS < 0 {
		return fmt.Errorf("cas cannot be negative")
	}

	return nil
}
