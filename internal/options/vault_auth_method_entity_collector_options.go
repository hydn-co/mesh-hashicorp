package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type VaultAuthMethodEntityCollectorOptions struct {
	VaultOptionsCore `json:",inline"`
}

func (o *VaultAuthMethodEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://hashicorp/collectors/vault_auth_method_entity_collector_options"
}

func (o *VaultAuthMethodEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Applications}
}

func (o *VaultAuthMethodEntityCollectorOptions) GetRequirements() []string {
	return []string{"hashicorp", "vault"}
}
