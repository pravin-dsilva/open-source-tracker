package oslib

import "fmt"

func FetchKubernetesPRs(users []string, token string) map[string][]Issue {
	orgs := []string{"kubernetes", "kubernetes-sigs"}
	result := make(map[string][]Issue)

	for _, user := range users {
		var userPRs []Issue
		for _, org := range orgs {
			apiURL := fmt.Sprintf("https://api.github.com/search/issues?q=org:%s+author:%s+is:pr", org, user)
			prs := fetchGitHubData(apiURL, token)
			userPRs = append(userPRs, prs...)
		}
		result[user] = userPRs
	}

	return result
}
