package main

import (
	"flag"
	"log"
	"os"
	"oslib/oslib"
)

func main() {
	showIssues := flag.Bool("issues", false, "Fetch only good first issues")
	showMonthlyReport := flag.Bool("monthlyreport", false, "Fetch only good first issues")
	flag.Parse()

	configPath := "config.json"
	config, err := oslib.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("Environment variable GITHUB_TOKEN is required")
	}
	if *showIssues {
		oslib.GenerateIssuesReport(config.Orgs, githubToken, config.Labels)
	} else if *showMonthlyReport {
		oslib.GenerateTeamAchievements(config.Users, githubToken)
	} else {
		oslib.GenerateReport(config.Users, githubToken)
	}
}
