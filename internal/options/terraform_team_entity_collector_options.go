package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type TerraformTeamEntityCollectorOptions struct {
	TerraformOptionsCore `json:",inline"`
}

func (o *TerraformTeamEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://hashicorp/collectors/terraform_team_entity_collector_options"
}

func (o *TerraformTeamEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Groups, spaces.GroupMembers}
}

func (o *TerraformTeamEntityCollectorOptions) GetRequirements() []string {
	return []string{"hashicorp", "terraform"}
}
