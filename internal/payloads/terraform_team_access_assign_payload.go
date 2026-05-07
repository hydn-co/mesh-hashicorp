package payloads

import (
	"fmt"
	"strings"
)

type TerraformTeamAccessAssignPayload struct {
	TeamID      string `json:"team_id"      binding:"required" title:"Team ID"      description:"The HCP Terraform team ID receiving the workspace access."`
	WorkspaceID string `json:"workspace_id" binding:"required" title:"Workspace ID" description:"The HCP Terraform workspace ID where access will be assigned."`
	Access      string `json:"access"       binding:"required" title:"Access Level" description:"The Terraform access level to grant, such as read, plan, write, or admin."`
}

func (p *TerraformTeamAccessAssignPayload) GetDiscriminator() string {
	return "mesh://hashicorp/actions/terraform_team_access_assign_payload"
}

func (p *TerraformTeamAccessAssignPayload) Validate() error {
	if p == nil {
		return fmt.Errorf("payload is required")
	}
	if strings.TrimSpace(p.TeamID) == "" {
		return fmt.Errorf("team_id is required")
	}
	if strings.TrimSpace(p.WorkspaceID) == "" {
		return fmt.Errorf("workspace_id is required")
	}
	if strings.TrimSpace(p.Access) == "" {
		return fmt.Errorf("access is required")
	}
	return nil
}
