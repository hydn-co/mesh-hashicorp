package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type VaultIdentityEntityCollectorOptions struct {
	VaultOptionsCore `json:",inline"`
}

func (o *VaultIdentityEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://hashicorp/collectors/vault_identity_entity_collector_options"
}

func (o *VaultIdentityEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Accounts, spaces.Groups, spaces.GroupMembers}
}

func (o *VaultIdentityEntityCollectorOptions) GetRequirements() []string {
	return []string{"hashicorp", "vault"}
}
