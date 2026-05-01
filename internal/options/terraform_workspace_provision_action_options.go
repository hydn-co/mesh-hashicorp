package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type TerraformWorkspaceProvisionActionOptions struct {
	TerraformOptionsCore `json:",inline"`
}

func (o *TerraformWorkspaceProvisionActionOptions) GetDiscriminator() string {
	return "mesh://hashicorp/actions/terraform_workspace_provision_action_options"
}

func (o *TerraformWorkspaceProvisionActionOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Applications}
}

func (o *TerraformWorkspaceProvisionActionOptions) GetRequirements() []string {
	return []string{"hashicorp", "terraform"}
}
