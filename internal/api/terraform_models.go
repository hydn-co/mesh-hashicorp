package api

type TerraformResourceIdentifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type TerraformResourceIdentifiers struct {
	Data []TerraformResourceIdentifier `json:"data,omitempty"`
}

type TerraformResourceIdentifierLink struct {
	Data TerraformResourceIdentifier `json:"data,omitzero"`
}

type TerraformUserAttributes struct {
	Username         *string `json:"username,omitempty"`
	IsServiceAccount bool    `json:"is-service-account,omitempty"`
	AuthMethod       string  `json:"auth-method,omitempty"`
	Email            string  `json:"email,omitempty"`
}

type TerraformUser struct {
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Attributes TerraformUserAttributes `json:"attributes,omitzero"`
}

func (u TerraformUser) DisplayName() string {
	if u.Attributes.Username != nil && *u.Attributes.Username != "" {
		return *u.Attributes.Username
	}
	if u.Attributes.Email != "" {
		return u.Attributes.Email
	}
	return u.ID
}

type TerraformOrganizationMembershipAttributes struct {
	Status string `json:"status"`
}

type TerraformOrganizationMembershipRelationships struct {
	Teams TerraformResourceIdentifiers    `json:"teams,omitzero"`
	User  TerraformResourceIdentifierLink `json:"user,omitzero"`
}

type TerraformOrganizationMembership struct {
	ID            string                                       `json:"id"`
	Type          string                                       `json:"type"`
	Attributes    TerraformOrganizationMembershipAttributes    `json:"attributes,omitzero"`
	Relationships TerraformOrganizationMembershipRelationships `json:"relationships,omitzero"`
}

type TerraformTeamAttributes struct {
	Name       string `json:"name"`
	Visibility string `json:"visibility,omitempty"`
	UsersCount int    `json:"users-count,omitempty"`
}

type TerraformTeamRelationships struct {
	Users TerraformResourceIdentifiers `json:"users,omitzero"`
}

type TerraformTeam struct {
	ID            string                     `json:"id"`
	Type          string                     `json:"type"`
	Attributes    TerraformTeamAttributes    `json:"attributes,omitzero"`
	Relationships TerraformTeamRelationships `json:"relationships,omitzero"`
}

type TerraformOrganizationMembershipResult struct {
	Membership TerraformOrganizationMembership `json:"membership,omitzero"`
	User       TerraformUser                   `json:"user,omitzero"`
}

type TerraformTeamResult struct {
	Team      TerraformTeam            `json:"team,omitzero"`
	UsersByID map[string]TerraformUser `json:"-"`
}

type terraformLinks struct {
	Next string `json:"next,omitempty"`
}

type terraformPaginationInfo struct {
	CurrentPage int  `json:"current-page,omitempty"`
	NextPage    *int `json:"next-page,omitempty"`
}

type terraformMeta struct {
	Pagination terraformPaginationInfo `json:"pagination,omitzero"`
}

type terraformMembershipsResponse struct {
	Data     []TerraformOrganizationMembership `json:"data,omitempty"`
	Included []TerraformUser                   `json:"included,omitempty"`
	Links    terraformLinks                    `json:"links,omitzero"`
	Meta     terraformMeta                     `json:"meta,omitzero"`
}

type terraformTeamsResponse struct {
	Data     []TerraformTeam `json:"data,omitempty"`
	Included []TerraformUser `json:"included,omitempty"`
	Links    terraformLinks  `json:"links,omitzero"`
	Meta     terraformMeta   `json:"meta,omitzero"`
}
