package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type VaultKVV2SecretSetActionOptions struct {
	VaultOptionsCore `json:",inline"`
}

func (o *VaultKVV2SecretSetActionOptions) GetDiscriminator() string {
	return "mesh://hashicorp/actions/vault_kv_v2_secret_set_action_options"
}

func (o *VaultKVV2SecretSetActionOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.SecurityConfigurations}
}

func (o *VaultKVV2SecretSetActionOptions) GetRequirements() []string {
	return []string{"hashicorp", "vault"}
}
