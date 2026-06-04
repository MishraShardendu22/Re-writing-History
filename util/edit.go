package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Edit(sshRemote string,repo string) {
	// change the dir to that repo
	if err := os.Chdir(repo); err != nil {
		log.Fatalf("Failed to change directory: %v", err)
	}

	// open this commit file
	f, err := os.Open("../updated_commits.txt")
	if err != nil {
		log.Fatalf("Failed to open updated_commits.txt: %v", err)
	}
	defer f.Close()

	// make a map of updated commit body basically for Fast lookup
	commitUpdates := make(map[string]struct {
		name  string
		email string
		date  string
	})

	// Reads commit file line-by-line.
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// Parses each commit line.
		parts := strings.SplitN(scanner.Text(), "|", 5)

		// Skips malformed lines, shouldnt arise
		if len(parts) < 4 {
			continue
		}

		// old hash, name, email and date
		hash, name, email, date := parts[0], parts[1], parts[2], parts[3]
		commitUpdates[hash] = struct {
			name  string
			email string
			date  string
		}{name, email, date}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading updated_commits.txt: %v", err)
	}

	if len(commitUpdates) == 0 {
		log.Println("No commits to update found in updated_commits.txt. Exiting.")
		return
	}
	fmt.Printf("Found %d commits to update.\n", len(commitUpdates))

	// Creates efficient string builder.
	// basically dynamically generating a shell script.
	var envFilter strings.Builder
	envFilter.WriteString("#!/bin/bash\n")

	// During git filter-branch, Git sets: $GIT_COMMIT to current commit being rewritten.
	for hash, update := range commitUpdates {
		// if current commit hash matches target hash, then overrides metadata.
		fmt.Fprintf(&envFilter, `	
			if [ "$GIT_COMMIT" = "%s" ]; then
				echo "Updating commit %s" >&2
				export GIT_AUTHOR_NAME="%s"
				export GIT_AUTHOR_EMAIL="%s"
				export GIT_AUTHOR_DATE="%s"
				export GIT_COMMITTER_NAME="%s"
				export GIT_COMMITTER_EMAIL="%s"
				export GIT_COMMITTER_DATE="%s"
				echo "Updated commit %s" >&2
			fi
		`, hash, hash[:8], update.name, update.email, update.date, update.name, update.email, update.date, hash[:8])
	}

	// export GIT_AUTHOR_NAME
	// export GIT_AUTHOR_EMAIL
	// export GIT_AUTHOR_DATE
	// change author metadata.

	// export GIT_COMMITTER_NAME
	// export GIT_COMMITTER_EMAIL
	// export GIT_COMMITTER_DATE
	// change committer metadata.

	// basically building large string incrementally
	// Writes generated shell script to disk.
	if err := os.WriteFile("env-filter.sh", []byte(envFilter.String()), 0755); err != nil {
		log.Fatalf("Failed to write env-filter script: %v", err)
	}

	// used to get the current branch name.
	// it is needed because git filter-branch, needs target branch 
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOutput, err := branchCmd.Output()
	if err != nil {
		log.Fatalf("Failed to get current branch: %v", err)
	}
	currentBranch := strings.TrimSpace(string(branchOutput))
	fmt.Printf("Current branch: %s\n", currentBranch)

	// basically rebuild every commit in history
	// run ("env-filter.sh") this shell script for EVERY commit
	rewriteCmd := exec.Command("git", "filter-branch", "-f", "--env-filter", envFilter.String(), currentBranch)
	rewriteCmd.Stdout = os.Stdout
	rewriteCmd.Stderr = os.Stderr
	fmt.Println("Rewriting history...")
	if err := rewriteCmd.Run(); err != nil {
		log.Fatalf("git filter-branch failed: %v", err)
	}


	cleanupCmds := [][]string{
		// update-ref is used to manipulate Git references directly.
		// (-d for deleting) delete backup pointer to old history
		{"git", "update-ref", "-d", "refs/original/refs/heads/" + currentBranch},

		// Git keeps hidden recovery logs called reflogs.
		// Even after deleting commits - Git remembers them temporarily
		// Reflog lets you recover - deleted branches, lost commits, resets
		// expire=now	delete immediately
		// --all	    all reflogs
		{"git", "reflog", "expire", "--expire=now", "--all"},
		
		// this is git's garbage collection removal
		// Git objects live inside - .git/objects/
		// Even deleted commits remain there until cleaned.
		// Without gc - old rewritten commits still physically exist
		{"git", "gc", "--prune=now", "--aggressive"},
		
		// Git tracks remote branches locally.
		// If remote branch deleted - local tracking ref may remain stale, This command removes stale remote references.
		{"git", "remote", "prune", "origin"},
	}

	// execute all clean up commnads
	for _, args := range cleanupCmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Warning: %s failed: %v", strings.Join(args, " "), err)
		}
	}

	// set remote url to the remote origin
	fmt.Println("Setting remote to SSH...")
	setRemoteCmd := exec.Command("git", "remote", "set-url", "origin", sshRemote)
	setRemoteCmd.Stdout = os.Stdout
	setRemoteCmd.Stderr = os.Stderr
	if err := setRemoteCmd.Run(); err != nil {
		log.Fatalf("Failed to set remote to SSH: %v", err)
	}

	// push the changes
	fmt.Println("Force pushing to origin...")
	pushCmd := exec.Command("git", "push", "origin", currentBranch, "--force")
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		log.Fatalf("Force push failed: %v", err)
	}

	if err := os.Chdir(".."); err != nil {
		log.Fatalf("Failed to change back to parent directory: %v", err)
	}

	// cleanup
	fmt.Println("Cleaning up local files...")
	os.RemoveAll(repo)
	os.Remove("updated_commits.txt")
	os.Remove("edited_commits.txt")
	os.Remove(repo + "/env-filter.sh")
	fmt.Println("Git history rewriting completed successfully!")
}


// Git internally is a DAG (Directed Acyclic Graph) of commits.

// Each commit stores:
// - commit hash
// - parent commit hash
// - author
// - committer
// - timestamps
// - commit message
// - tree snapshot

// A commit is immutable.
// You cannot edit a commit directly.

// Git history rewriting actually means -
// create completely new commits
// with modified metadata.

// Suppose original history -
// A -> B -> C -> D

// If you change commit B date -
// B changes hash
// Now C depends on old B.

// So Git must recreate -
// C
// D

// Final -

// A -> B' -> C' -> D'
// Everything downstream changes.
// That is why history rewriting cascades.