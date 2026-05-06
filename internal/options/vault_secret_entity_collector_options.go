package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type VaultSecretEntityCollectorOptions struct {
	VaultOptionsCore `json:",inline"`
}

func (o *VaultSecretEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://hashicorp/collectors/vault_secret_entity_collector_options"
}

func (o *VaultSecretEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Secrets}
}

func (o *VaultSecretEntityCollectorOptions) GetRequirements() []string {
	return []string{"hashicorp", "vault"}
}
