package options

import "github.com/fgrzl/json/polymorphic"

func init() {
	polymorphic.RegisterType[TerraformAccountEntityCollectorOptions]()
	polymorphic.RegisterType[TerraformTeamEntityCollectorOptions]()
	polymorphic.RegisterType[TerraformWorkspaceEntityCollectorOptions]()
	polymorphic.RegisterType[TerraformPolicyEntityCollectorOptions]()
	polymorphic.RegisterType[TerraformTeamAccessEntityCollectorOptions]()
	polymorphic.RegisterType[TerraformAuditTrailActivityCollectorOptions]()
	polymorphic.RegisterType[VaultIdentityAccountEntityCollectorOptions]()
	polymorphic.RegisterType[VaultIdentityGroupEntityCollectorOptions]()
	polymorphic.RegisterType[VaultPolicyEntityCollectorOptions]()
	polymorphic.RegisterType[VaultAuthMethodEntityCollectorOptions]()
	polymorphic.RegisterType[VaultSecretEntityCollectorOptions]()
	polymorphic.RegisterType[TerraformTeamProvisionActionOptions]()
	polymorphic.RegisterType[TerraformWorkspaceProvisionActionOptions]()
	polymorphic.RegisterType[TerraformTeamMembershipAssignActionOptions]()
	polymorphic.RegisterType[TerraformTeamAccessAssignActionOptions]()
	polymorphic.RegisterType[VaultKVV1SecretSetActionOptions]()
	polymorphic.RegisterType[VaultKVV2SecretSetActionOptions]()
}
