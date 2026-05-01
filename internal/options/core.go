package options

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

type VaultOptionsCore struct {
	Address   string `json:"address"             title:"Address"   description:"Vault base URL."                binding:"required"`
	Namespace string `json:"namespace,omitempty" title:"Namespace" description:"Vault namespace when required."`
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
