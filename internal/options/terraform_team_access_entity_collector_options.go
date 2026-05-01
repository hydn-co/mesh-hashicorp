package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type TerraformTeamAccessEntityCollectorOptions struct {
	TerraformOptionsCore `json:",inline"`
}

func (o *TerraformTeamAccessEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://hashicorp/collectors/terraform_team_access_entity_collector_options"
}

func (o *TerraformTeamAccessEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Permissions, spaces.GroupPermissions}
}

func (o *TerraformTeamAccessEntityCollectorOptions) GetRequirements() []string {
	return []string{"hashicorp", "terraform"}
}
