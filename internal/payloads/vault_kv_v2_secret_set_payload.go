package payloads

import (
	"encoding/json"
	"fmt"
)

type VaultKVV2SecretSetPayload struct {
	MountPath  string          `json:"mount_path"    binding:"required" title:"Mount Path"            description:"The Vault KV v2 mount path that stores the secret."`
	SecretPath string          `json:"secret_path"   binding:"required" title:"Secret Path"           description:"The path of the secret relative to the KV v2 mount."`
	Data       json.RawMessage `json:"data"          binding:"required" title:"Secret Data"           description:"A JSON object containing the key-value pairs to write to the secret."                                 additionalProperties:"true"`
	CAS        *int            `json:"cas,omitempty"                    title:"Check-And-Set Version" description:"Optional optimistic concurrency version. Set to 0 to require that the secret does not already exist."                             minimum:"0"`
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
