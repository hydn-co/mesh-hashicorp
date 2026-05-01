package collectors

import (
	"strings"

	"github.com/hydn-co/mesh-hashicorp/internal/api"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/entities"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/types"
)

func newTerraformAccount(accountRef string, user api.TerraformUser, status string) *entities.Account {
	displayName := user.DisplayName()
	description := "HCP Terraform organization membership"
	if authMethod := strings.TrimSpace(user.Attributes.AuthMethod); authMethod != "" {
		description += " via " + authMethod
	}
	if trimmedStatus := strings.TrimSpace(status); trimmedStatus != "" && !strings.EqualFold(trimmedStatus, "active") {
		description += " (" + trimmedStatus + ")"
	}

	account := entities.NewAccount()
	account.AccountRef = accountRef
	account.AccountType = types.AccountTypeUser
	if user.Attributes.IsServiceAccount {
		account.AccountType = types.AccountTypeServicePrincipal
	}
	account.Name = displayName
	account.DisplayName = displayName
	account.Description = description
	account.Enabled = strings.EqualFold(strings.TrimSpace(status), "active")
	if email := strings.TrimSpace(user.Attributes.Email); email != "" {
		account.PrimaryEmail = &types.Email{Address: email}
	}
	return account
}

func newTerraformGroup(team api.TerraformTeam) *entities.Group {
	group := entities.NewGroup()
	group.GroupRef = team.ID
	group.Name = strings.TrimSpace(team.Attributes.Name)
	if group.Name == "" {
		group.Name = team.ID
	}
	if visibility := strings.TrimSpace(team.Attributes.Visibility); visibility != "" {
		group.Description = "Visibility: " + visibility
	}
	return group
}

func newTerraformGroupMember(groupRef, accountRef string, user api.TerraformUser) *entities.GroupMember {
	groupMember := entities.NewGroupMember()
	groupMember.GroupRef = groupRef
	groupMember.AccountRef = accountRef
	if authMethod := strings.TrimSpace(user.Attributes.AuthMethod); authMethod != "" {
		groupMember.RoleRef = authMethod
	}
	return groupMember
}
