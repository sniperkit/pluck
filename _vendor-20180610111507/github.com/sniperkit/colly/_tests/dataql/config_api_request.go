package main

import (
	"fmt"
	"strings"
)

var (

	//////////////////////////////////////////
	// Github API - Request vars
	//////////////////////////////////////////

	// GITHUB_API_DOMAIN specifies the Github's API v3 domain
	GITHUB_API_DOMAIN = "https://api.github.com"

	// GITHUB_API_ACCOUNT sets the github username to target for github's api requests.
	GITHUB_API_ACCOUNT = "roscopecoltran"

	// GITHUB_API_PAGINATION_OFFSET specifies the start page for the request
	GITHUB_API_PAGINATION_OFFSET = 1

	// GITHUB_API_PAGINATION_OFFSET specifies the max number of results to return
	GITHUB_API_PAGINATION_PER_PAGE = 10

	// GITHUB_API_PAGINATION_SORT_KEY specifies the sort key to use to order returned results
	GITHUB_API_PAGINATION_SORT_KEY = "updated"

	// GITHUB_API_PAGINATION_SORT_ORDER specifies the sort order. available: desc or asc
	GITHUB_API_PAGINATION_SORT_ORDER = "desc"

	// GITHUB_API_ENDPOINT_PARAMS sets the parameters list for the pagination options.
	GITHUB_API_ENDPOINT_PARAMS_LIST = []string{
		fmt.Sprintf("page=%d", GITHUB_API_PAGINATION_OFFSET),
		fmt.Sprintf("per_page=%d", GITHUB_API_PAGINATION_PER_PAGE),
		fmt.Sprintf("sort=%s", GITHUB_API_PAGINATION_SORT_KEY),
		fmt.Sprintf("direction=%s", GITHUB_API_PAGINATION_SORT_ORDER),
	}

	// GITHUB_API_ENDPOINT_QUERY_STRING sets the query string
	GITHUB_API_ENDPOINT_QUERY_STRING = fmt.Sprintf("?%s", strings.Join(GITHUB_API_ENDPOINT_PARAMS_LIST, "&"))

	// GITHUB_API_ENDPOINT_URI sets the request uri
	GITHUB_API_ENDPOINT_URI = fmt.Sprintf("/users/%s/starred", GITHUB_API_ACCOUNT)

	// GITHUB_API_ENDPOINT_REQUEST sets the full request url
	GITHUB_API_ENDPOINT_REQUEST = fmt.Sprintf("%s%s%s", GITHUB_API_DOMAIN, GITHUB_API_ENDPOINT_URI, GITHUB_API_ENDPOINT_QUERY_STRING)

	//-- End
)
