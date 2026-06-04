package util

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

func Edit(sshRemote string, repo string) {
	// cleanup, has to happen even if scrpit is interuppted
	defer func() {
		Logger.LogAttrs(context.Background(), slog.LevelInfo, "cleaning up local files")

		os.Chdir("..")

		os.RemoveAll(repo)
		os.Remove("updated_commits.txt")
		os.Remove("edited_commits.txt")
	}()

	Logger.LogAttrs(context.Background(), slog.LevelInfo, "edit phase started",
		slog.String("repo", repo),
		slog.String("ssh_remote", sshRemote),
	)

	// change the dir to that repo
	if err := os.Chdir(repo); err != nil {
		Logger.LogAttrs(context.Background(), slog.LevelError, "failed to change directory to repo",
			slog.String("repo", repo),
			slog.String("error", err.Error()),
		)
		fmt.Fprintf(os.Stderr, "Failed to change directory: %v\n", err)
		os.Exit(1)
	}

	// open this commit file
	f, err := os.Open("../updated_commits.txt")
	if err != nil {
		Logger.LogAttrs(context.Background(), slog.LevelError, "failed to open updated_commits.txt",
			slog.String("error", err.Error()),
		)
		fmt.Fprintf(os.Stderr, "Failed to open updated_commits.txt: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// make a map of updated commit body basically for Fast lookup
	commitUpdates := make(map[string]struct {
		name    string
		email   string
		date    string
		message string
	})

	// Reads commit file line-by-line.
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// Parses each commit line.
		parts := strings.SplitN(scanner.Text(), "|", 5)

		// Skips malformed lines, shouldnt arise
		if len(parts) < 5 {
			continue
		}

		// old hash, name, email and date
		hash, name, email, date, msg := parts[0], parts[1], parts[2], parts[3], parts[4]
		commitUpdates[hash] = struct {
			name    string
			email   string
			date    string
			message string
		}{name, email, date, msg}
	}

	if err := scanner.Err(); err != nil {
		Logger.LogAttrs(context.Background(), slog.LevelError, "error reading updated_commits.txt",
			slog.String("error", err.Error()),
		)
		fmt.Fprintf(os.Stderr, "Error reading updated_commits.txt: %v\n", err)
		os.Exit(1)
	}

	if len(commitUpdates) == 0 {
		Logger.LogAttrs(context.Background(), slog.LevelWarn, "no commits to update, exiting")
		return
	}
	Logger.LogAttrs(context.Background(), slog.LevelInfo, "found commits to update",
		slog.Int("commit_count", len(commitUpdates)),
	)

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
		Logger.LogAttrs(context.Background(), slog.LevelError, "failed to write env-filter script",
			slog.String("error", err.Error()),
		)
		fmt.Fprintf(os.Stderr, "Failed to write env-filter script: %v\n", err)
		os.Exit(1)
	}
	Logger.LogAttrs(context.Background(), slog.LevelInfo, "env-filter script written")

	// generates commit message rewrite script
	var msgFilter strings.Builder
	msgFilter.WriteString("#!/bin/bash\n")

	first := true
	for hash, update := range commitUpdates {
		if first {
			fmt.Fprintf(&msgFilter,
				`if [ "$GIT_COMMIT" = "%s" ]; then
					echo "%s"
				`, hash, update.message)

			first = false
		} else {
			fmt.Fprintf(&msgFilter,
				`elif [ "$GIT_COMMIT" = "%s" ]; then
					echo "%s"
				`, hash, update.message)
		}
	}

	msgFilter.WriteString(`
		else
			cat
		fi
	`)

	if err := os.WriteFile("msg-filter.sh", []byte(msgFilter.String()), 0755); err != nil {
		Logger.LogAttrs(context.Background(), slog.LevelError, "failed to write msg-filter script",
			slog.String("error", err.Error()),
		)
		fmt.Fprintf(os.Stderr, "Failed to write msg-filter script: %v\n", err)
		os.Exit(1)
	}

	Logger.LogAttrs(context.Background(), slog.LevelInfo, "msg-filter script written")

	// used to get the current branch name.
	// it is needed because git filter-branch, needs target branch
	branchOut, err := NewGitCommand("git", "rev-parse", "--abbrev-ref", "HEAD").Run()
	if err != nil {
		Logger.LogAttrs(context.Background(), slog.LevelError, "failed to get current branch",
			slog.String("error", err.Error()),
		)
		fmt.Fprintf(os.Stderr, "Failed to get current branch: %v\n", err)
		os.Exit(1)
	}

	currentBranch := strings.TrimSpace(branchOut)

	Logger.LogAttrs(context.Background(), slog.LevelInfo, "current branch identified",
		slog.String("branch", currentBranch),
	)

	// basically rebuild every commit in history
	// run ("env-filter.sh") this shell script for EVERY commit
	Logger.LogAttrs(context.Background(), slog.LevelInfo, "rewriting history via filter-branch",
		slog.String("branch", currentBranch),
	)

	// the need for these sh files is that,
	// filter branch accepts only sh files in it as arguments
	rewriteCmd := exec.Command(
		"git",
		"filter-branch",
		"-f",
		"--env-filter",
		envFilter.String(),
		"--msg-filter",
		msgFilter.String(),
		currentBranch,
	)

	rewriteCmd.Stdout = os.Stdout
	rewriteCmd.Stderr = os.Stderr

	if err := rewriteCmd.Run(); err != nil {
		Logger.LogAttrs(context.Background(), slog.LevelError, "git filter-branch failed",
			slog.String("branch", currentBranch),
			slog.String("error", err.Error()),
		)
		fmt.Fprintf(os.Stderr, "git filter-branch failed: %v\n", err)
		os.Exit(1)
	}

	cleanupCmds := []*GitCommand{
		// update-ref is used to manipulate Git references directly.
		// (-d for deleting) delete backup pointer to old history
		NewGitCommand("git", "update-ref", "-d", "refs/original/refs/heads/"+currentBranch),

		// Git keeps hidden recovery logs called reflogs.
		// Even after deleting commits - Git remembers them temporarily
		// Reflog lets you recover - deleted branches, lost commits, resets
		// expire=now	delete immediately
		// --all	    all reflogs
		NewGitCommand("git", "reflog", "expire", "--expire=now", "--all"),

		// this is git's garbage collection removal
		// Git objects live inside - .git/objects/
		// Even deleted commits remain there until cleaned.
		// Without gc - old rewritten commits still physically exist
		NewGitCommand("git", "gc", "--prune=now", "--aggressive"),

		// Git tracks remote branches locally.
		// If remote branch deleted - local tracking ref may remain stale, This command removes stale remote references.
		NewGitCommand("git", "remote", "prune", "origin"),
	}

	// execute all clean up commnads
	for _, gc := range cleanupCmds {
		if err := gc.RunWithStdout(); err != nil {
			Logger.LogAttrs(context.Background(), slog.LevelWarn, "cleanup command produced warning",
				slog.String("command", gc.Name),
				slog.String("args", strings.Join(gc.Args, " ")),
				slog.String("error", err.Error()),
			)
		}
	}

	// set remote url to the remote origin
	Logger.LogAttrs(context.Background(), slog.LevelInfo, "setting remote to SSH",
		slog.String("ssh_remote", sshRemote),
	)
	setRemoteCmd := exec.Command("git", "remote", "set-url", "origin", sshRemote)
	setRemoteCmd.Stdout = os.Stdout
	setRemoteCmd.Stderr = os.Stderr
	if err := setRemoteCmd.Run(); err != nil {
		Logger.LogAttrs(context.Background(), slog.LevelError, "failed to set remote to SSH",
			slog.String("error", err.Error()),
		)
		fmt.Fprintf(os.Stderr, "Failed to set remote to SSH: %v\n", err)
		os.Exit(1)
	}

	// push the changes
	Logger.LogAttrs(context.Background(), slog.LevelInfo, "force pushing to origin",
		slog.String("branch", currentBranch),
	)
	pushCmd := exec.Command("git", "push", "origin", currentBranch, "--force")
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		Logger.LogAttrs(context.Background(), slog.LevelError, "force push failed",
			slog.String("branch", currentBranch),
			slog.String("error", err.Error()),
		)
		fmt.Fprintf(os.Stderr, "Force push failed: %v\n", err)
		os.Exit(1)
	}

	if err := os.Chdir(".."); err != nil {
		Logger.LogAttrs(context.Background(), slog.LevelError, "failed to change back to parent directory",
			slog.String("error", err.Error()),
		)
		fmt.Fprintf(os.Stderr, "Failed to change back to parent directory: %v\n", err)
		os.Exit(1)
	}
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
