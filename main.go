package main

import (
	"log"
	"os"

	"github.com/MishraShardendu22/util"
	"github.com/joho/godotenv"
)

var repo = "Dhvani-Commit-Tester"
var sshRemote = "git@github.com-learning:ShardenduMishra22/" + repo + ".git"

var	start = "2026-06-01 00:00:00 +0530"
var	end = "2026-06-05 23:59:59 +0530"

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}


func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is not set. Please set it in your environment or .env file.")
	}

	// 1. Clone the repo (no ssh) and generate edited_commits.txt
	util.Clone("https://github.com/ShardenduMishra22/" + repo + ".git")

	// This step generally, I made optional. 
	// I could not do this manually edit commits and push as well 
	
	// 2.1 Run AI to generate updated_commits.txt
	util.AIRun(start, end, apiKey)

	// 2.2 Run a non AI automated sequential time script (very limited)
	// util.Run(start,end)

	// 2.3 Do it manually and push

	// 3. Edit commit history and push
	util.Edit(sshRemote,repo)
}
