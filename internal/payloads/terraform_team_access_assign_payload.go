package payloads

import (
	"fmt"
	"strings"
)

type TerraformTeamAccessAssignPayload struct {
	TeamID      string `json:"team_id"      binding:"required"`
	WorkspaceID string `json:"workspace_id" binding:"required"`
	Access      string `json:"access"       binding:"required"`
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
