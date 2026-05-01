package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldDecodeTerraformUserAttributesWhenProviderUsesKebabCase(t *testing.T) {
	var user TerraformUser
	err := json.Unmarshal([]byte(`{
		"id":"user-1",
		"type":"users",
		"attributes":{
			"username":"alice",
			"is-service-account":true,
			"auth-method":"hcp_sso",
			"email":"alice@example.com"
		}
	}`), &user)

	require.NoError(t, err)
	assert.Equal(t, "alice", user.DisplayName())
	assert.True(t, user.Attributes.IsServiceAccount)
	assert.Equal(t, "hcp_sso", user.Attributes.AuthMethod)
	assert.Equal(t, "alice@example.com", user.Attributes.Email)
}

func TestShouldDecodeTerraformTeamAttributesWhenProviderUsesKebabCase(t *testing.T) {
	var team TerraformTeam
	err := json.Unmarshal([]byte(`{
		"id":"team-1",
		"type":"teams",
		"attributes":{
			"name":"owners",
			"visibility":"organization",
			"users-count":3
		}
	}`), &team)

	require.NoError(t, err)
	assert.Equal(t, "owners", team.Attributes.Name)
	assert.Equal(t, "organization", team.Attributes.Visibility)
	assert.Equal(t, 3, team.Attributes.UsersCount)
}

func TestShouldDecodePaginationWhenProviderUsesKebabCase(t *testing.T) {
	var response terraformTeamsResponse
	err := json.Unmarshal([]byte(`{
		"meta":{
			"pagination":{
				"current-page":1,
				"next-page":2
			}
		}
	}`), &response)

	require.NoError(t, err)
	if assert.NotNil(t, response.Meta.Pagination.NextPage) {
		assert.Equal(t, 2, *response.Meta.Pagination.NextPage)
	}
	assert.Equal(t, 1, response.Meta.Pagination.CurrentPage)
}
