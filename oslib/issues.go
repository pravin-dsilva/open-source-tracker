package oslib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Issue struct {
	Title     string `json:"title"`
	URL       string `json:"html_url"`
	Repo      string `json:"repository_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func FetchAssignedIssues(username, token string) []Issue {
	return fetchGitHubData(fmt.Sprintf("https://api.github.com/search/issues?q=assignee:%s+is:issue+is:open", username), token)
}

func FetchCreatedIssues(username, token string) []Issue {
	return fetchGitHubData(fmt.Sprintf("https://api.github.com/search/issues?q=author:%s+is:issue+is:open", username), token)
}

func FetchOpenPRs(username, token string) []Issue {
	return fetchGitHubData(fmt.Sprintf("https://api.github.com/search/issues?q=author:%s+is:pr+is:open", username), token)
}

func FetchClosedPRs(username, token string) []Issue {
	oneYearAgo := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	apiURL := fmt.Sprintf(
		"https://api.github.com/search/issues?q=author:%s+is:pr+is:closed+closed:>=%s",
		username,
		oneYearAgo,
	)
	return fetchGitHubData(apiURL, token)
}

func fetchGitHubData(apiURL, token string) []Issue {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "token "+token)

	for {
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to fetch data: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("Failed to read response body: %v", err)
			}

			var result struct {
				Items []Issue `json:"items"`
			}
			if err := json.Unmarshal(body, &result); err != nil {
				log.Fatalf("Failed to unmarshal response: %v", err)
			}

			// Process repository field to extract only the repo name
			for i := range result.Items {
				repoParts := strings.Split(result.Items[i].Repo, "/")
				if len(repoParts) >= 2 {
					result.Items[i].Repo = repoParts[len(repoParts)-1]
				}
			}
			//Sorting
			sort.Slice(result.Items, func(i, j int) bool {
				return result.Items[i].CreatedAt > result.Items[j].CreatedAt
			})
			return result.Items
		}

		if resp.StatusCode == 403 {
			log.Println("Rate limit exceeded. Retrying after 60 seconds...")
			time.Sleep(60 * time.Second)
			continue
		}

		// Log and exit for non-retryable status codes
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("GitHub API returned status: %d. Body: %s", resp.StatusCode, body)
	}

	return nil
}

func FetchIssues(org, token string, label string) []Issue {
	apiURL := fmt.Sprintf("https://api.github.com/search/issues?q=org:%s+label:%q+is:issue+is:open", org, label)
	return fetchGitHubData(apiURL, token)
}
