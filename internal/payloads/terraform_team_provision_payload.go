package payloads

import (
	"fmt"
	"strings"
)

type TerraformTeamProvisionPayload struct {
	Name string `json:"name" binding:"required" title:"Team Name" description:"The HCP Terraform team name to create."`
}

func (p *TerraformTeamProvisionPayload) GetDiscriminator() string {
	return "mesh://hashicorp/actions/terraform_team_provision_payload"
}

func (p *TerraformTeamProvisionPayload) Validate() error {
	if p == nil {
		return fmt.Errorf("payload is required")
	}
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}
