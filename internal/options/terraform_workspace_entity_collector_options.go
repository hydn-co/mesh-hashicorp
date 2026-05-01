package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type TerraformWorkspaceEntityCollectorOptions struct {
	TerraformOptionsCore `json:",inline"`
}

func (o *TerraformWorkspaceEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://hashicorp/collectors/terraform_workspace_entity_collector_options"
}

func (o *TerraformWorkspaceEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Applications}
}

func (o *TerraformWorkspaceEntityCollectorOptions) GetRequirements() []string {
	return []string{"hashicorp", "terraform"}
}
