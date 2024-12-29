package oslib

import (
	"bytes"
	"log"
	"os"
	"text/template"
	"time"
)

func GenerateReport(users []string, token string) {
	data := make(map[string]map[string][]Issue)
	for i, user := range users {
		data[user] = map[string][]Issue{
			"assigned_issues": FetchAssignedIssues(user, token),
			"created_issues":  FetchCreatedIssues(user, token),
			"open_prs":        FetchOpenPRs(user, token),
		}
		if i < len(users)-1 {
			log.Printf("Sleeping for 30 seconds to avoid rate-limiting")
			time.Sleep(10 * time.Second)
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
					</div>
				</div>
			</div>
			{{ end }}
		</div>
	</body>
	</html>
	`))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatalf("Error rendering template: %v", err)
	}

	if err := os.WriteFile("docs/index.html", buf.Bytes(), 0644); err != nil {
		log.Fatalf("Error saving HTML file: %v", err)
	}

	log.Println("HTML report generated and saved to dashboard.html")
}
