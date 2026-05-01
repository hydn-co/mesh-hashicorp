package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type TerraformTeamMembershipAssignActionOptions struct {
	TerraformOptionsCore `json:",inline"`
}

func (o *TerraformTeamMembershipAssignActionOptions) GetDiscriminator() string {
	return "mesh://hashicorp/actions/terraform_team_membership_assign_action_options"
}

func (o *TerraformTeamMembershipAssignActionOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.GroupMembers}
}

func (o *TerraformTeamMembershipAssignActionOptions) GetRequirements() []string {
	return []string{"hashicorp", "terraform"}
}
