package payloads

import (
	"encoding/json"
	"fmt"
)

type VaultKVV1SecretSetPayload struct {
	MountPath  string          `json:"mount_path"  binding:"required" title:"Mount Path"  description:"The Vault KV v1 mount path that stores the secret."`
	SecretPath string          `json:"secret_path" binding:"required" title:"Secret Path" description:"The path of the secret relative to the KV v1 mount."`
	Data       json.RawMessage `json:"data"        binding:"required" title:"Secret Data" description:"A JSON object containing the key-value pairs to write to the secret." additionalProperties:"true"`
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
