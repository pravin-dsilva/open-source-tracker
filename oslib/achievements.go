package oslib

import (
	"sort"
	"time"
)

type Activity struct {
	Title     string
	URL       string
	Repo      string
	Timestamp time.Time
	Action    string
}

// FetchMonthlyActivity gathers all PR and Issue activities for a list of users
func FetchMonthlyActivity(users []string, token string) map[string][]Activity {
	activityByUser := make(map[string][]Activity)

	for _, user := range users {
		var activities []Activity

		// 1. Open PRs
		for _, pr := range FetchOpenPRs(user, token) {
			t, _ := time.Parse(time.RFC3339, pr.CreatedAt)
			activities = append(activities, Activity{
				Title:     pr.Title,
				URL:       pr.URL,
				Repo:      pr.Repo,
				Timestamp: t,
				Action:    "opened_pr",
			})
		}

		// 2. Closed PRs
		for _, pr := range FetchClosedPRs(user, token) {
			t, _ := time.Parse(time.RFC3339, pr.CreatedAt)
			activities = append(activities, Activity{
				Title:     pr.Title,
				URL:       pr.URL,
				Repo:      pr.Repo,
				Timestamp: t,
				Action:    "closed_pr",
			})
		}

		// 3. Open Created Issues
		for _, issue := range FetchCreatedIssues(user, token) {
			t, _ := time.Parse(time.RFC3339, issue.CreatedAt)
			activities = append(activities, Activity{
				Title:     issue.Title,
				URL:       issue.URL,
				Repo:      issue.Repo,
				Timestamp: t,
				Action:    "created_issue_open",
			})
		}

		// 4. Closed Created Issues
		for _, issue := range FetchClosedIssues(user, token) {
			t, _ := time.Parse(time.RFC3339, issue.CreatedAt)
			activities = append(activities, Activity{
				Title:     issue.Title,
				URL:       issue.URL,
				Repo:      issue.Repo,
				Timestamp: t,
				Action:    "created_issue_closed",
			})
		}

		activityByUser[user] = activities
	}

	return activityByUser
}

type MonthlyUserActivity struct {
	User       string
	Month      string // "2025-07"
	Activities []Activity
}

func GroupMonthlyActivity(activityByUser map[string][]Activity) (map[string]map[string][]Activity, []string) {
	grouped := make(map[string]map[string][]Activity)
	monthSet := make(map[string]bool)

	for user, acts := range activityByUser {
		for _, act := range acts {
			month := act.Timestamp.Format("2006-01")

			if _, ok := grouped[month]; !ok {
				grouped[month] = make(map[string][]Activity)
			}
			grouped[month][user] = append(grouped[month][user], act)
			monthSet[month] = true
		}
	}

	// Sort months descending
	var months []string
	for m := range monthSet {
		months = append(months, m)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(months)))

	return grouped, months
}
