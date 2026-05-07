package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type VaultIdentityAccountEntityCollectorOptions struct {
	VaultOptionsCore `json:",inline"`
}

func (o *VaultIdentityAccountEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://hashicorp/collectors/vault_identity_account_entity_collector_options"
}

func (o *VaultIdentityAccountEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Accounts}
}

func (o *VaultIdentityAccountEntityCollectorOptions) GetRequirements() []string {
	return []string{"hashicorp", "vault"}
}
