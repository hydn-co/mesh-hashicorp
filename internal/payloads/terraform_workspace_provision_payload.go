package payloads

import (
	"fmt"
	"strings"
)

type TerraformWorkspaceProvisionPayload struct {
	Name        string `json:"name"                  binding:"required"`
	Description string `json:"description,omitempty"`
	ProjectID   string `json:"project_id,omitempty"`
}

func (p *TerraformWorkspaceProvisionPayload) GetDiscriminator() string {
	return "mesh://hashicorp/actions/terraform_workspace_provision_payload"
}

func (p *TerraformWorkspaceProvisionPayload) Validate() error {
	if p == nil {
		return fmt.Errorf("payload is required")
	}
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}
