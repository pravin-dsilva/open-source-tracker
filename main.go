package main

import (
	"log"
	"os"
	"oslib/oslib"
)

func main() {

	configPath := "config.json"
	config, error := oslib.LoadConfig(configPath)

	if error != nil {
		log.Fatalf("Error loading config: %v", error)
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("Environment variable GITHUB_TOKEN is required")
	}

	oslib.GenerateReport(config.Users, githubToken)
}
