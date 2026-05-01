package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type TerraformTeamProvisionActionOptions struct {
	TerraformOptionsCore `json:",inline"`
}

func (o *TerraformTeamProvisionActionOptions) GetDiscriminator() string {
	return "mesh://hashicorp/actions/terraform_team_provision_action_options"
}

func (o *TerraformTeamProvisionActionOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Groups}
}

func (o *TerraformTeamProvisionActionOptions) GetRequirements() []string {
	return []string{"hashicorp", "terraform"}
}
