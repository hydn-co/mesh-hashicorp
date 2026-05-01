package payloads

import "github.com/fgrzl/json/polymorphic"

func init() {
	polymorphic.RegisterType[TerraformTeamProvisionPayload]()
	polymorphic.RegisterType[TerraformWorkspaceProvisionPayload]()
	polymorphic.RegisterType[TerraformTeamMembershipAssignPayload]()
	polymorphic.RegisterType[TerraformTeamAccessAssignPayload]()
}
