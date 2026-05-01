package api

import (
	"net/url"
	"strconv"
	"strings"
)

func nextPage(current int, links terraformLinks, meta terraformPaginationInfo) (int, bool) {
	if meta.NextPage != nil && *meta.NextPage > current {
		return *meta.NextPage, true
	}
	if links.Next == "" {
		return 0, false
	}
	parsed, err := url.Parse(links.Next)
	if err != nil {
		return 0, false
	}
	page := strings.TrimSpace(parsed.Query().Get("page[number]"))
	if page == "" {
		return 0, false
	}
	next, err := strconv.Atoi(page)
	if err != nil || next <= current {
		return 0, false
	}
	return next, true
}
