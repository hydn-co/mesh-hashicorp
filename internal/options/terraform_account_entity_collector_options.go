package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type TerraformAccountEntityCollectorOptions struct {
	TerraformOptionsCore `json:",inline"`
}

func (o *TerraformAccountEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://hashicorp/collectors/terraform_account_entity_collector_options"
}

func (o *TerraformAccountEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Accounts}
}

func (o *TerraformAccountEntityCollectorOptions) GetRequirements() []string {
	return []string{"hashicorp", "terraform"}
}
