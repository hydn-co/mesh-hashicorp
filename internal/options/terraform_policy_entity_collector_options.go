package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type TerraformPolicyEntityCollectorOptions struct {
	TerraformOptionsCore `json:",inline"`
}

func (o *TerraformPolicyEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://hashicorp/collectors/terraform_policy_entity_collector_options"
}

func (o *TerraformPolicyEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Policies, spaces.SecurityConfigurations}
}

func (o *TerraformPolicyEntityCollectorOptions) GetRequirements() []string {
	return []string{"hashicorp", "terraform"}
}
