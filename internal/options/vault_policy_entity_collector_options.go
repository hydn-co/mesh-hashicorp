package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type VaultPolicyEntityCollectorOptions struct {
	VaultOptionsCore `json:",inline"`
}

func (o *VaultPolicyEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://hashicorp/collectors/vault_policy_entity_collector_options"
}

func (o *VaultPolicyEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Policies}
}

func (o *VaultPolicyEntityCollectorOptions) GetRequirements() []string {
	return []string{"hashicorp", "vault"}
}
