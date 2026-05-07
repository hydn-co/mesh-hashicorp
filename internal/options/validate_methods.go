package options

func (o *TerraformAccountEntityCollectorOptions) Validate() error {
	if o == nil {
		return validateTerraformOptionsCore(nil)
	}
	return validateTerraformOptionsCore(&o.TerraformOptionsCore)
}

func (o *TerraformAuditTrailActivityCollectorOptions) Validate() error {
	if o == nil {
		return validateTerraformOptionsCore(nil)
	}
	return validateTerraformOptionsCore(&o.TerraformOptionsCore)
}

func (o *TerraformPolicyEntityCollectorOptions) Validate() error {
	if o == nil {
		return validateTerraformOptionsCore(nil)
	}
	return validateTerraformOptionsCore(&o.TerraformOptionsCore)
}

func (o *TerraformTeamAccessAssignActionOptions) Validate() error {
	if o == nil {
		return validateTerraformOptionsCore(nil)
	}
	return validateTerraformOptionsCore(&o.TerraformOptionsCore)
}

func (o *TerraformTeamAccessEntityCollectorOptions) Validate() error {
	if o == nil {
		return validateTerraformOptionsCore(nil)
	}
	return validateTerraformOptionsCore(&o.TerraformOptionsCore)
}

func (o *TerraformTeamEntityCollectorOptions) Validate() error {
	if o == nil {
		return validateTerraformOptionsCore(nil)
	}
	return validateTerraformOptionsCore(&o.TerraformOptionsCore)
}

func (o *TerraformTeamMembershipAssignActionOptions) Validate() error {
	if o == nil {
		return validateTerraformOptionsCore(nil)
	}
	return validateTerraformOptionsCore(&o.TerraformOptionsCore)
}

func (o *TerraformTeamProvisionActionOptions) Validate() error {
	if o == nil {
		return validateTerraformOptionsCore(nil)
	}
	return validateTerraformOptionsCore(&o.TerraformOptionsCore)
}

func (o *TerraformWorkspaceEntityCollectorOptions) Validate() error {
	if o == nil {
		return validateTerraformOptionsCore(nil)
	}
	return validateTerraformOptionsCore(&o.TerraformOptionsCore)
}

func (o *TerraformWorkspaceProvisionActionOptions) Validate() error {
	if o == nil {
		return validateTerraformOptionsCore(nil)
	}
	return validateTerraformOptionsCore(&o.TerraformOptionsCore)
}

func (o *VaultAuthMethodEntityCollectorOptions) Validate() error {
	if o == nil {
		return validateVaultOptionsCore(nil)
	}
	return validateVaultOptionsCore(&o.VaultOptionsCore)
}

func (o *VaultIdentityAccountEntityCollectorOptions) Validate() error {
	if o == nil {
		return validateVaultOptionsCore(nil)
	}
	return validateVaultOptionsCore(&o.VaultOptionsCore)
}

func (o *VaultIdentityGroupEntityCollectorOptions) Validate() error {
	if o == nil {
		return validateVaultOptionsCore(nil)
	}
	return validateVaultOptionsCore(&o.VaultOptionsCore)
}

func (o *VaultKVV1SecretSetActionOptions) Validate() error {
	if o == nil {
		return validateVaultOptionsCore(nil)
	}
	return validateVaultOptionsCore(&o.VaultOptionsCore)
}

func (o *VaultKVV2SecretSetActionOptions) Validate() error {
	if o == nil {
		return validateVaultOptionsCore(nil)
	}
	return validateVaultOptionsCore(&o.VaultOptionsCore)
}

func (o *VaultPolicyEntityCollectorOptions) Validate() error {
	if o == nil {
		return validateVaultOptionsCore(nil)
	}
	return validateVaultOptionsCore(&o.VaultOptionsCore)
}

func (o *VaultSecretEntityCollectorOptions) Validate() error {
	if o == nil {
		return validateVaultOptionsCore(nil)
	}
	return validateVaultOptionsCore(&o.VaultOptionsCore)
}
