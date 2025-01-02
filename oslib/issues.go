package oslib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

type Issue struct {
	Title string `json:"title"`
	URL   string `json:"html_url"`
	Repo  string `json:"repository_url"`
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


func fetchGitHubData(apiURL, token string) []Issue {
    client := &http.Client{}
    req, err := http.NewRequest("GET", apiURL, nil)
    if err != nil {
        log.Fatalf("Failed to create request: %v", err)
    }
    req.Header.Set("Authorization", "token "+token)

    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("Failed to fetch data: %v", err)
    }
    defer resp.Body.Close()

    log.Printf("GitHub API URL: %s, Status Code: %d", apiURL, resp.StatusCode)
    log.Printf("Rate Limit Remaining: %s, Rate Limit Reset: %s",
        resp.Header.Get("X-RateLimit-Remaining"),
        resp.Header.Get("X-RateLimit-Reset"),
    )

    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        log.Fatalf("GitHub API returned status: %s. Body: %s", resp.Status, string(body))
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Failed to read response body: %v", err)
    }
    //log.Printf("GitHub API Response: %s", string(body))

    var result struct {
        Items []Issue `json:"items"`
    }
    if err := json.Unmarshal(body, &result); err != nil {
        log.Fatalf("Failed to unmarshal response: %v. Body: %s", err, string(body))
    }

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

func FetchGoodFirstIssues(org, token string) []Issue {
	apiURL := fmt.Sprintf("https://api.github.com/search/issues?q=org:%s+label:%q+is:issue+is:open", org, "good+first+issue")
	return fetchGitHubData(apiURL, token)
}
