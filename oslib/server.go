package oslib

import (
	"bytes"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"
)

func GenerateIssuesReport(orgs []string, token string, labels []string) {
	for _, label := range labels {
		var Issues []Issue

		for _, org := range orgs {
			Issues = append(Issues, FetchIssues(org, token, label)...)
		}

		sort.Slice(Issues, func(i, j int) bool {
			return Issues[i].CreatedAt > Issues[j].CreatedAt
		})

		tmpl := template.Must(template.New("goodFirstIssues").Parse(`
    <!DOCTYPE html>
    <html>
    <head>
        <title>Issues</title>
        <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/css/bootstrap.min.css" rel="stylesheet">
    </head>
    <body class="container mt-5">
        <h1 class="mb-4">Issues</h1>
        {{ if . }}
        <table class="table table-striped">
            <thead>
                <tr>
                    <th>Title</th>
                    <th>Repository</th>
                    <th>URL</th>
                    <th>Created At</th>
                </tr>
            </thead>
            <tbody>
                {{ range . }}
                <tr>
                    <td>{{ .Title }}</td>
                    <td>{{ .Repo }}</td>
                    <td><a href="{{ .URL }}" target="_blank">{{ .URL }}</a></td>
                    <td>{{ .CreatedAt }}</td>
                </tr>
                {{ end }}
            </tbody>
        </table>
        {{ else }}
        <p>No issues found</p>
        {{ end }}
    </body>
    </html>
    `))

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, Issues); err != nil {
			log.Fatalf("Error rendering template: %v", err)
		}

		var outputFile string
		if label == "good+first+issue" {
			outputFile = "docs/good_first_issues.html"
		} else if label == "help+wanted" {
			outputFile = "docs/help_wanted.html"
		}

		if err := os.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
			log.Fatalf("Error saving HTML file (%s): %v", outputFile, err)
		}

		log.Println("HTML report is generated")
	}
}
func GenerateReport(users []string, token string) {
	data := make(map[string]map[string][]Issue)
	for i, user := range users {
		data[user] = map[string][]Issue{
			"assigned_issues": FetchAssignedIssues(user, token),
			"created_issues":  FetchCreatedIssues(user, token),
			"open_prs":        FetchOpenPRs(user, token),
			"closed_prs":      FetchClosedPRs(user, token),
		}
		if i < len(users)-1 {
			log.Printf("Sleeping for 30 seconds to avoid rate-limiting")
			time.Sleep(30 * time.Second)
		}

	}

	tmpl := template.Must(template.New("dashboard").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>GitHub Dashboard</title>
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/css/bootstrap.min.css" rel="stylesheet">
		<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/js/bootstrap.bundle.min.js"></script>
		<style>
			h3 { font-size: 1.25rem; font-weight: bold; background-color: #f8f9fa; padding: 0.5rem; border-radius: 0.25rem; }
			.table td, .table th { width: auto; text-align: left; }
		</style>
	</head>
	<body class="container mt-5">
		<h1 class="mb-4">GitHub Dashboard</h1>
		<div class="accordion" id="usersAccordion">
			{{ range $user, $data := . }}
			<div class="accordion-item">
				<h2 class="accordion-header" id="heading-{{ $user }}">
					<button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapse-{{ $user }}" aria-expanded="false" aria-controls="collapse-{{ $user }}">
						<span style="font-size: 1.5rem; font-weight: bold;">{{ $user }}</span>
					</button>
				</h2>
				<div id="collapse-{{ $user }}" class="accordion-collapse collapse" aria-labelledby="heading-{{ $user }}" data-bs-parent="#usersAccordion">
					<div class="accordion-body">
						<h3 style="background-color: #d1e7dd;">Assigned Issues</h3>
						{{ if $data.assigned_issues }}
						<table class="table table-striped">
							<thead>
								<tr><th>Title</th><th>Repository</th><th>URL</th><th>Updated At</th></tr>
							</thead>
							<tbody>
								{{ range $issue := $data.assigned_issues }}
								<tr>
									<td>{{ $issue.Title }}</td>
									<td>{{ $issue.Repo }}</td>
									<td><a href="{{ $issue.URL }}" target="_blank">{{ $issue.URL }}</a></td>
									<td>{{ $issue.UpdatedAt }}</td>
								</tr>
								{{ end }}
							</tbody>
						</table>
						{{ else }}<p>No assigned issues</p>{{ end }}

						<h3 style="background-color: #ffeeba;">Created Issues</h3>
						{{ if $data.created_issues }}
						<table class="table table-striped">
							<thead>
								<tr><th>Title</th><th>Repository</th><th>URL</th><th>Updated At</th></tr>
							</thead>
							<tbody>
								{{ range $issue := $data.created_issues }}
								<tr>
									<td>{{ $issue.Title }}</td>
									<td>{{ $issue.Repo }}</td>
									<td><a href="{{ $issue.URL }}" target="_blank">{{ $issue.URL }}</a></td>
									<td>{{ $issue.UpdatedAt }}</td>
								</tr>
								{{ end }}
							</tbody>
						</table>
						{{ else }}<p>No created issues</p>{{ end }}

						<h3 style="background-color: #f8d7da;">Open PRs</h3>
						{{ if $data.open_prs }}
						<table class="table table-striped">
							<thead>
								<tr><th>Title</th><th>Repository</th><th>URL</th><th>Updated At</th></tr>
							</thead>
							<tbody>
								{{ range $issue := $data.open_prs }}
								<tr>
									<td>{{ $issue.Title }}</td>
									<td>{{ $issue.Repo }}</td>
									<td><a href="{{ $issue.URL }}" target="_blank">{{ $issue.URL }}</a></td>
									<td>{{ $issue.UpdatedAt }}</td>
								</tr>
								{{ end }}
							</tbody>
						</table>
						{{ else }}<p>No open PRs</p>{{ end }}
						<h3 style="background-color: #f8d7da;">Closed PRs (past 1 year)</h3>
						{{ if $data.closed_prs }}
						<table class="table table-striped">
							<thead>
								<tr><th>Title</th><th>Repository</th><th>URL</th><th>Updated At</th></tr>
							</thead>
							<tbody>
								{{ range $issue := $data.closed_prs }}
								<tr>
									<td>{{ $issue.Title }}</td>
									<td>{{ $issue.Repo }}</td>
									<td><a href="{{ $issue.URL }}" target="_blank">{{ $issue.URL }}</a></td>
									<td>{{ $issue.UpdatedAt }}</td>
								</tr>
								{{ end }}
							</tbody>
						</table>
						{{ else }}<p>No closed PRs</p>{{ end }}
					</div>
				</div>
			</div>
			{{ end }}
		</div>

		</tbody>
	</table>
</div>
	</body>
	</html>
	`))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatalf("Error rendering template: %v", err)
	}

	if err := os.WriteFile("docs/user_dashboard.html", buf.Bytes(), 0644); err != nil {
		log.Fatalf("Error saving HTML file: %v", err)
	}

	log.Println("HTML report generated and saved to dashboard.html")
}

func GenerateTeamAchievements(users []string, token string) {
	activityMap := FetchMonthlyActivity(users, token)
	groupedData, months := GroupMonthlyActivity(activityMap)

	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("Jan 2")
		},
		"badgeClass": func(action string) string {
			switch action {
			case "created_issue_open":
				return "primary"
			case "created_issue_closed":
				return "info"
			case "opened_pr":
				return "warning"
			case "closed_pr":
				return "success"
			default:
				return "secondary"
			}
		},
		"actionLabel": func(action string) string {
			switch action {
			case "created_issue_open":
				return "Created Issue (Open)"
			case "created_issue_closed":
				return "Created Issue (Closed)"
			case "opened_pr":
				return "Opened PR"
			case "closed_pr":
				return "Closed PR"
			default:
				return action
			}
		},
		"filterByAction": func(acts []Activity, action string) []Activity {
			var filtered []Activity
			for _, a := range acts {
				if a.Action == action {
					filtered = append(filtered, a)
				}
			}
			return filtered
		},
		"escapeID": func(s string) string {
			safe := strings.ReplaceAll(s, "@", "_")
			safe = strings.ReplaceAll(safe, ".", "_")
			safe = strings.ReplaceAll(safe, "/", "_")
			return safe
		},
	}

	tmpl := template.Must(template.New("teamAchievements").Funcs(funcMap).Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Team Achievements</title>
	<meta charset="utf-8">
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/css/bootstrap.min.css" rel="stylesheet">
	<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/js/bootstrap.bundle.min.js"></script>
</head>
<body class="container mt-5">
	<h1 class="mb-4">Team Achievements by Month</h1>

	{{ range .Months }}
		{{ $month := . }}
		<h2 class="mt-4">{{ $month }}</h2>

		{{ range $user, $activities := index $.Data $month }}
			<div class="card mb-2">
				<div class="card-header">
					<h5 class="mb-0">
						<button class="btn btn-link text-decoration-none" data-bs-toggle="collapse" data-bs-target="#collapse-{{ $month }}-{{ $user | escapeID }}" aria-expanded="false" aria-controls="collapse-{{ $month }}-{{ $user | escapeID }}">
							{{ $user }}
						</button>
					</h5>
					<div class="mt-2">
						<span class="badge bg-primary">Open Issues: {{ len (filterByAction $activities "created_issue_open") }}</span>
						<span class="badge bg-info text-dark">Closed Issues: {{ len (filterByAction $activities "created_issue_closed") }}</span>
						<span class="badge bg-warning text-dark">Opened PRs: {{ len (filterByAction $activities "opened_pr") }}</span>
						<span class="badge bg-success">Closed PRs: {{ len (filterByAction $activities "closed_pr") }}</span>
					</div>
				</div>
				<div id="collapse-{{ $month }}-{{ $user | escapeID }}" class="collapse">
					<ul class="list-group list-group-flush">
						{{ range $a := $activities }}
						<li class="list-group-item">
							<span class="badge bg-{{ badgeClass $a.Action }}">{{ actionLabel $a.Action }}</span>
							<a href="{{ $a.URL }}" target="_blank">{{ $a.Title }}</a>
							<span class="text-muted">in {{ $a.Repo }} on {{ formatDate $a.Timestamp }}</span>
						</li>
						{{ end }}
					</ul>
				</div>
			</div>
		{{ end }}

	{{ end }}
</body>
</html>
`))

	data := struct {
		Data   map[string]map[string][]Activity
		Months []string
	}{
		Data:   groupedData,
		Months: months,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatalf("Error rendering template: %v", err)
	}

	if err := os.WriteFile("docs/team_achievements.html", buf.Bytes(), 0644); err != nil {
		log.Fatalf("Error saving HTML file: %v", err)
	}

	log.Println("Team achievements dashboard generated: docs/team_achievements.html")
}

func GenerateKubernetesContributions(users []string, token string) {
	data := FetchKubernetesPRs(users, token)

	tmpl := template.Must(template.New("k8sTable").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Kubernetes PR Contributions</title>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/css/bootstrap.min.css" rel="stylesheet">
	<style>
		body { font-size: 0.95rem; }
		table td, table th { vertical-align: top; }
		.pr-link { font-size: 0.8rem; margin-right: 6px; display: inline-block; text-decoration: none; }
		.pr-link:hover { text-decoration: underline; }
	</style>
</head>
<body class="container mt-5">
	<h1 class="mb-4">Kubernetes PR Contributions</h1>
	<div class="mb-3">
		<a href="user_dashboard.html" class="btn btn-secondary">← Back to Dashboard</a>
	</div>

	<table class="table table-bordered table-sm align-middle">
		<thead class="table-light">
			<tr>
				<th style="width: 20%;">User</th>
				<th>Pull Requests</th>
			</tr>
		</thead>
		<tbody>
			{{ range $user, $prs := . }}
			<tr>
				<td><strong>{{$user}}</strong></td>
				<td>
					{{ if $prs }}
						{{ range $i, $pr := $prs }}
							<a href="{{ $pr.URL }}" target="_blank" class="pr-link">{{ printf "[%s]" $pr.Repo }}</a>
						{{ end }}
					{{ else }}
						<em>No PRs</em>
					{{ end }}
				</td>
			</tr>
			{{ end }}
		</tbody>
	</table>
</body>
</html>
`))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatalf("Error rendering Kubernetes contributions UI: %v", err)
	}

	if err := os.WriteFile("docs/kubernetes_contributions.html", buf.Bytes(), 0644); err != nil {
		log.Fatalf("Error writing HTML file: %v", err)
	}

	log.Println("Generated compact Kubernetes PR table: docs/kubernetes_contributions.html")
}
