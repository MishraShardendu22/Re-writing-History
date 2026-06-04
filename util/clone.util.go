package util

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func Clone(url string, repo string) {
	Logger.LogAttrs(context.Background(), slog.LevelInfo, "clone phase started",
		slog.String("repo", repo),
		slog.String("url", url),
	)

	repoName := strings.TrimSuffix(filepath.Base(url), ".git")
	cloneDir := "./" + repoName

	CloneRepoIfNeeded(url, cloneDir)
	commitLog := GetGitLog(cloneDir)

	CreateAndWriteToTheFile("edited_commits.txt", commitLog)

	Logger.LogAttrs(context.Background(), slog.LevelInfo, "clone phase completed",
		slog.String("repo", repo),
		slog.String("clone_dir", cloneDir),
	)
}

func CloneRepoIfNeeded(url, dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		Logger.LogAttrs(context.Background(), slog.LevelInfo, "cloning repository",
			slog.String("url", url),
			slog.String("dir", dir),
		)

		gc := NewGitCommand("git", "clone", url, dir)
		if err := gc.RunWithStdout(); err != nil {
			Logger.LogAttrs(context.Background(), slog.LevelError, "git clone failed",
				slog.String("url", url),
				slog.String("dir", dir),
				slog.String("error", err.Error()),
			)
			fmt.Fprintf(os.Stderr, "git clone failed: %v\n", err)
			os.Exit(1)
		}
	} else {
		Logger.LogAttrs(context.Background(), slog.LevelInfo, "repository already exists, skipping clone",
			slog.String("dir", dir),
		)
	}
}

func GetGitLog(dir string) []byte {
	Logger.LogAttrs(context.Background(), slog.LevelInfo, "getting git log",
		slog.String("dir", dir),
	)

	gc := NewGitCommandInDir(dir, "git", "log", "--pretty=format:%H|%an|%ae|%ad|%s", "--date=iso")
	out, err := gc.Run()
	if err != nil {
		Logger.LogAttrs(context.Background(), slog.LevelError, "git log failed",
			slog.String("dir", dir),
			slog.String("error", err.Error()),
		)
		fmt.Fprintf(os.Stderr, "git log failed: %v\n", err)
		os.Exit(1)
	}

	Logger.LogAttrs(context.Background(), slog.LevelInfo, "git log retrieved",
		slog.String("dir", dir),
		slog.Int("bytes", len(out)),
	)
	return []byte(out)
}

// --pretty=format:...
// Customizes git log output.

// %H	Full commit hash
// %an	Author name
// %ae	Author email
// %ad	Author date
// %s	Commit message

// "--date=iso"
// Makes dates readable ISO format.

func CreateAndWriteToTheFile(filename string, data []byte) {
	// better approach
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		Logger.LogAttrs(context.Background(), slog.LevelError, "failed to write file",
			slog.String("filename", filename),
			slog.String("error", err.Error()),
		)
		fmt.Fprintf(os.Stderr, "failed to write file %s: %v\n", filename, err)
		os.Exit(1)
	}

	Logger.LogAttrs(context.Background(), slog.LevelInfo, "file written successfully",
		slog.String("filename", filename),
		slog.Int("bytes", len(data)),
	)

	// another way of doing this
	// // create the file
	// f, err := os.Create(filename)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()

	// // write to that file
	// w := bufio.NewWriter(f)

	// // Writes data into buffer, not necessarily into file immediately
	// _, err = w.WriteString(string(data))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Without Flush(), some data may remain in memory and never reach the file.
	// // Forces buffered data to actually be written into the file
	// w.Flush()
}

// 0 644
// │ └── actual permissions
// └──── octal notation

// Writing HTTP logs continuously -> bufio.Writer
// Saving generated JSON once -> os.WriteFile
