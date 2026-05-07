package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type VaultIdentityGroupEntityCollectorOptions struct {
	VaultOptionsCore `json:",inline"`
}

func (o *VaultIdentityGroupEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://hashicorp/collectors/vault_identity_group_entity_collector_options"
}

func (o *VaultIdentityGroupEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Groups, spaces.GroupMembers}
}

func (o *VaultIdentityGroupEntityCollectorOptions) GetRequirements() []string {
	return []string{"hashicorp", "vault"}
}
