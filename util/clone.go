package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Clone(url string) {
	repoName := strings.TrimSuffix(filepath.Base(url), ".git")
	cloneDir := "./" + repoName

	CloneRepoIfNeeded(url, cloneDir)
	commitLog := GetGitLog(cloneDir)

	CreateAndWriteToTheFile("edited_commits.txt", commitLog)
}

func CloneRepoIfNeeded(url, dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Println("Cloning repo...")
		cmd := exec.Command("git", "clone", url, dir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("git clone failed: %v", err)
		}
	}
}

func GetGitLog(dir string) []byte {
	cmd := exec.Command("git", "log", "--pretty=format:%H|%an|%ae|%ad|%s", "--date=iso")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("git log failed: %v", err)
	}
	return out
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
		log.Fatal(err)
	}

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
