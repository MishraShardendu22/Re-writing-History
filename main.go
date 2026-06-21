package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/MishraShardendu22/util"
	"github.com/joho/godotenv"
)

var repo = "Dhvani-Commit-Tester"
var userName = "ShardenduMishra22"
var sshRemote = "git@github.com-learning:" + userName + "/" + repo + ".git"

// var repo = "Employee-Frontend"
// var repo = "Employee-Leave-Management"
// var userName = "singhvanshiki"
// var sshRemote = "git@github.com-didi:" + userName + "/" + repo + ".git"

var start = "2026-06-17 00:00:00 +0530"
var end = "2026-06-19 23:59:59 +0530"

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		util.Logger.LogAttrs(context.Background(), slog.LevelError, "failed to load .env file",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		util.Logger.LogAttrs(context.Background(), slog.LevelError, "API_KEY environment variable is not set")
		os.Exit(1)
	}

	overallStart := time.Now()

	// 1. Clone the repo (no ssh) and generate edited_commits.txt
	cloneStart := time.Now()
	util.LogPhaseStart("clone", slog.String("repo", repo))
	util.Clone("https://github.com/" + userName + "/" + repo + ".git", repo)
	util.LogPhaseEnd("clone", time.Since(cloneStart), slog.String("repo", repo))

	// This step generally, I made optional.
	// I could not do this manually edit commits and push as well

	// 2.1 Run AI to generate updated_commits.txt
	aiRunStart := time.Now()
	util.LogPhaseStart("ai-run", slog.String("repo", repo))
	util.AIRun(start, end, apiKey)
	util.LogPhaseEnd("ai-run", time.Since(aiRunStart), slog.String("repo", repo))

	// 2.2 Run a non AI automated sequential time script (very limited)
	// util.Run(start,end)

	// 2.3 Do it manually and push

	// 3. Edit commit history and push
	editStart := time.Now()
	util.LogPhaseStart("rewrite-push", slog.String("repo", repo))
	util.Edit(sshRemote, repo)
	util.LogPhaseEnd("rewrite-push", time.Since(editStart), slog.String("repo", repo))

	overallDuration := time.Since(overallStart)
	util.LogPhaseEnd("overall", overallDuration,
		slog.String("repo", repo),
		slog.Int64("active_jobs", util.ActiveJobs()),
	)
}
