package options

import "github.com/hydn-co/mesh-sdk/pkg/connectorutil"

type TerraformOptionsCore struct {
	Hostname     string `json:"hostname"     title:"Hostname"     description:"HCP Terraform or Terraform Enterprise hostname." binding:"required"`
	Organization string `json:"organization" title:"Organization" description:"Terraform organization name."                    binding:"required"`
}

func (o *TerraformOptionsCore) GetHostname() string {
	if o == nil {
		return ""
	}
	return o.Hostname
}

func (o *TerraformOptionsCore) GetOrganization() string {
	if o == nil {
		return ""
	}
	return o.Organization
}

func validateTerraformOptionsCore(o *TerraformOptionsCore) error {
	return connectorutil.RequireStrings(
		"feature options",
		connectorutil.RequiredString{Name: "hostname", Value: o.GetHostname()},
		connectorutil.RequiredString{Name: "organization", Value: o.GetOrganization()},
	)
}

type VaultOptionsCore struct {
	Address   string `json:"address"             title:"Address"   description:"Vault base URL."           binding:"required"`
	Namespace string `json:"namespace,omitempty" title:"Namespace" description:"Optional Vault namespace."`
}

func (o *VaultOptionsCore) GetAddress() string {
	if o == nil {
		return ""
	}
	return o.Address
}

func (o *VaultOptionsCore) GetNamespace() string {
	if o == nil {
		return ""
	}
	return o.Namespace
}

func validateVaultOptionsCore(o *VaultOptionsCore) error {
	return connectorutil.RequireStrings(
		"feature options",
		connectorutil.RequiredString{Name: "address", Value: o.GetAddress()},
	)
}
