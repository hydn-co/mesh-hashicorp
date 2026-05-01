package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/fgrzl/enumerators"
)

const terraformPageSize = 100

func (c *Client) OrganizationMembershipEnumerator(
	ctx context.Context,
	organization string,
) enumerators.Enumerator[TerraformOrganizationMembershipResult] {
	page := 1

	return enumerators.PageItemEnumerator(func() ([]TerraformOrganizationMembershipResult, bool, error) {
		if err := ctx.Err(); err != nil {
			return nil, false, err
		}

		response := terraformMembershipsResponse{}
		query := url.Values{}
		query.Set("include", "user,teams")
		query.Set("page[number]", strconv.Itoa(page))
		query.Set("page[size]", strconv.Itoa(terraformPageSize))

		path := fmt.Sprintf("/api/v2/organizations/%s/organization-memberships", url.PathEscape(organization))
		if err := c.get(ctx, path, query, &response); err != nil {
			return nil, false, err
		}

		usersByID := make(map[string]TerraformUser, len(response.Included))
		for _, user := range response.Included {
			usersByID[user.ID] = user
		}

		results := make([]TerraformOrganizationMembershipResult, 0, len(response.Data))
		for _, membership := range response.Data {
			result := TerraformOrganizationMembershipResult{Membership: membership}
			userID := membership.Relationships.User.Data.ID
			if user, ok := usersByID[userID]; ok {
				result.User = user
			} else if userID != "" {
				result.User = TerraformUser{ID: userID}
			}
			results = append(results, result)
		}

		next, ok := nextPage(page, response.Links, response.Meta.Pagination)
		if !ok {
			return results, false, nil
		}
		page = next

		return results, true, nil
	})
}

func (c *Client) TeamEnumerator(
	ctx context.Context,
	organization string,
) enumerators.Enumerator[TerraformTeamResult] {
	page := 1

	return enumerators.PageItemEnumerator(func() ([]TerraformTeamResult, bool, error) {
		if err := ctx.Err(); err != nil {
			return nil, false, err
		}

		response := terraformTeamsResponse{}
		query := url.Values{}
		query.Set("include", "users")
		query.Set("page[number]", strconv.Itoa(page))
		query.Set("page[size]", strconv.Itoa(terraformPageSize))

		path := fmt.Sprintf("/api/v2/organizations/%s/teams", url.PathEscape(organization))
		if err := c.get(ctx, path, query, &response); err != nil {
			return nil, false, err
		}

		usersByID := make(map[string]TerraformUser, len(response.Included))
		for _, user := range response.Included {
			usersByID[user.ID] = user
		}

		results := make([]TerraformTeamResult, 0, len(response.Data))
		for _, team := range response.Data {
			results = append(results, TerraformTeamResult{
				Team:      team,
				UsersByID: usersByID,
			})
		}

		next, ok := nextPage(page, response.Links, response.Meta.Pagination)
		if !ok {
			return results, false, nil
		}
		page = next

		return results, true, nil
	})
}
