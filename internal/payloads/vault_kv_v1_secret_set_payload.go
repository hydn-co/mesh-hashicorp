package payloads

import "fmt"

type VaultKVV1SecretSetPayload struct {
	MountPath  string         `json:"mount_path"  binding:"required"`
	SecretPath string         `json:"secret_path" binding:"required"`
	Data       map[string]any `json:"data"        binding:"required"`
}

func (p *VaultKVV1SecretSetPayload) GetDiscriminator() string {
	return "mesh://hashicorp/actions/vault_kv_v1_secret_set_payload"
}

func (p *VaultKVV1SecretSetPayload) Validate() error {
	if p == nil {
		return fmt.Errorf("payload is required")
	}

	return validateVaultSecretSetPayloadFields(p.MountPath, p.SecretPath, p.Data)
}
