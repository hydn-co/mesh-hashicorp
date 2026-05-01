package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type TerraformAuditTrailActivityCollectorOptions struct {
	TerraformOptionsCore `json:",inline"`
}

func (o *TerraformAuditTrailActivityCollectorOptions) GetDiscriminator() string {
	return "mesh://hashicorp/collectors/terraform_audit_trail_activity_collector_options"
}

func (o *TerraformAuditTrailActivityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Activity}
}

func (o *TerraformAuditTrailActivityCollectorOptions) GetRequirements() []string {
	return []string{"hashicorp", "terraform"}
}
