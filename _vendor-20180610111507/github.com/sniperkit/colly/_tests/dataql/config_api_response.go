package main

var (

	//////////////////////////////////////////
	// Github API - Results vars
	//////////////////////////////////////////

	// GITHUB_API_HEADERS_STARRED_DEFAULT sets the list of headers to match
	// results returned by a github's api specific endpoint.
	// The test endpoint returns all starred repositories by a specifc user.
	// Important: It DOESN'T return headers for nested objects like for the owner information.
	GITHUB_API_HEADERS_STARRED_DEFAULT = []string{"html_url", "keys_url", "pulls_url", "milestones_url", "subscription_url", "compare_url", "has_downloads", "id", "git_refs_url", "statuses_url", "stargazers_url", "git_url", "default_branch", "branches_url", "pushed_at", "watchers_count", "teams_url", "notifications_url", "labels_url", "name", "fork", "commits_url", "comments_url", "full_name", "issue_comment_url", "merges_url", "size", "license", "clone_url", "language", "owner", "private", "events_url", "languages_url", "contributors_url", "contents_url", "homepage", "has_pages", "description", "tags_url", "blobs_url", "git_tags_url", "issues_url", "stargazers_count", "has_wiki", "forks_count", "forks", "url", "releases_url", "created_at", "collaborators_url", "has_issues", "node_id", "forks_url", "subscribers_url", "downloads_url", "deployments_url", "svn_url", "mirror_url", "updated_at", "ssh_url", "has_projects", "hooks_url", "archived", "open_issues_count", "watchers", "issue_events_url", "assignees_url", "trees_url", "git_commits_url", "archive_url", "open_issues"}

	// GITHUB_API_HEADERS_STARRED_NESTED sets the list of headers to match
	// results returned by a github's api specific endpoint.
	// The test endpoint returns all starred repositories by a specifc user.
	// Important: It DOES return headers for nested objects like for the owner information with dot notation of paths.
	GITHUB_API_HEADERS_STARRED_NESTED = []string{
		"html_url",
		"keys_url",
		"pulls_url",
		"milestones_url",
		"subscription_url",
		"compare_url",
		"has_downloads",
		"id",
		"git_refs_url",
		"statuses_url",
		"stargazers_url",
		"git_url",
		"default_branch",
		"branches_url",
		"pushed_at",
		"watchers_count",
		"teams_url",
		"notifications_url",
		"labels_url",
		"name",
		"fork",
		"commits_url",
		"comments_url",
		"full_name",
		"issue_comment_url",
		"merges_url",
		"size",
		"license",
		"clone_url",
		"language",
		"owner",
		"private",
		"events_url",
		"languages_url",
		"contributors_url",
		"contents_url",
		"homepage",
		"has_pages",
		"description",
		"tags_url",
		"blobs_url",
		"git_tags_url",
		"issues_url",
		"stargazers_count",
		"has_wiki",
		"forks_count",
		"forks",
		"url",
		"releases_url",
		"created_at",
		"collaborators_url",
		"has_issues",
		"node_id",
		"forks_url",
		"subscribers_url",
		"downloads_url",
		"deployments_url",
		"svn_url",
		"mirror_url",
		"updated_at",
		"ssh_url",
		"has_projects",
		"hooks_url",
		"archived",
		"open_issues_count",
		"watchers",
		"issue_events_url",
		"assignees_url",
		"trees_url",
		"git_commits_url",
		"archive_url",
		"open_issues",
	}

	//-- End
)
