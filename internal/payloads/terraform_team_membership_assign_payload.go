package payloads

import (
	"fmt"
	"strings"
)

type TerraformTeamMembershipAssignPayload struct {
	TeamID string `json:"team_id" binding:"required" title:"Team ID" description:"The HCP Terraform team ID that will receive the membership."`
	UserID string `json:"user_id" binding:"required" title:"User ID" description:"The HCP Terraform user ID to assign to the team."`
}

func (p *TerraformTeamMembershipAssignPayload) GetDiscriminator() string {
	return "mesh://hashicorp/actions/terraform_team_membership_assign_payload"
}

func (p *TerraformTeamMembershipAssignPayload) Validate() error {
	if p == nil {
		return fmt.Errorf("payload is required")
	}
	if strings.TrimSpace(p.TeamID) == "" {
		return fmt.Errorf("team_id is required")
	}
	if strings.TrimSpace(p.UserID) == "" {
		return fmt.Errorf("user_id is required")
	}
	return nil
}
