package main

import (
	"flag"
	"log"
	"os"
	"oslib/oslib"
)

func main() {
	showGoodFirstIssues := flag.Bool("good-first-issues", false, "Fetch only good first issues")
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
	if *showGoodFirstIssues {
		oslib.GenerateGoodFirstIssuesReport(config.Orgs, githubToken)
	} else {
		oslib.GenerateReport(config.Users, githubToken)
	}
}
