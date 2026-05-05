package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type VaultKVV1SecretSetActionOptions struct {
	VaultOptionsCore `json:",inline"`
}

func (o *VaultKVV1SecretSetActionOptions) GetDiscriminator() string {
	return "mesh://hashicorp/actions/vault_kv_v1_secret_set_action_options"
}

func (o *VaultKVV1SecretSetActionOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.SecurityConfigurations}
}

func (o *VaultKVV1SecretSetActionOptions) GetRequirements() []string {
	return []string{"hashicorp", "vault"}
}
