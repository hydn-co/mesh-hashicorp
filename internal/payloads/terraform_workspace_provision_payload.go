package payloads

import (
	"fmt"
	"strings"
)

type TerraformWorkspaceProvisionPayload struct {
	Name        string `json:"name"                  binding:"required" title:"Workspace Name" description:"The HCP Terraform workspace name to create."`
	Description string `json:"description,omitempty"                    title:"Description"    description:"An optional description for the workspace."`
	ProjectID   string `json:"project_id,omitempty"                     title:"Project ID"     description:"The Terraform project ID to associate with the workspace when provided."`
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
