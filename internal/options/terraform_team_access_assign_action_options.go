package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type TerraformTeamAccessAssignActionOptions struct {
	TerraformOptionsCore `json:",inline"`
}

func (o *TerraformTeamAccessAssignActionOptions) GetDiscriminator() string {
	return "mesh://hashicorp/actions/terraform_team_access_assign_action_options"
}

func (o *TerraformTeamAccessAssignActionOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Permissions, spaces.GroupPermissions}
}

func (o *TerraformTeamAccessAssignActionOptions) GetRequirements() []string {
	return []string{"hashicorp", "terraform"}
}
