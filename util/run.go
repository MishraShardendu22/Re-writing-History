package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func AIRun(startD, endD, apiKey string) {
	in, err := ioutil.ReadFile("edited_commits.txt")
	if err != nil {
		log.Fatal(err)
	}

	systemPrompt := `You are a strict assistant for editing Git commit history.

Instructions:
- Input is a list of commits in the format:
  <hash>|<author>|<email>|<timestamp>|<message>

- Your task is to:
  1. Update the timestamp (4th field) to a new realistic value between ` + startD + ` and ` + endD + `.
  2. Rewrite commit messages into professional production-quality Git commit messages.
  3. Randomize commit authors across commits while preserving valid author-email pairing consistency.
  4. Ensure all authors end up with approximately the same number of commits.
  5. Make commit timing patterns appear fully organic and human-like:
     - irregular gaps between commits
     - realistic work/rest periods
     - varying commit densities
     - avoid predictable spacing patterns
     - include occasional clustered commits and long inactive gaps

- Commit message format rules:
  - Follow conventional commit style.
  - Use this structure:
    <type>(optional-scope): short lowercase summary

- Allowed commit types:
  - feat
  - fix
  - refactor
  - docs
  - chore

- Example formats:
  - feat(auth): add oauth login support
  - fix(api): prevent duplicate user creation
  - refactor(cache): simplify redis client setup
  - docs(readme): update installation steps
  - chore(deps): upgrade grpc version

- Commit message requirements:
  - Use imperative tense
  - Keep messages concise
  - No trailing period
  - Lowercase preferred
  - Scope should be realistic and technical when applicable

- Example professional scopes:
  - auth
  - api
  - ui
  - db
  - cache
  - docker
  - ci
  - parser

- Preserve the exact order of commits — do not reorder them.

- Do NOT modify:
  - commit hash
  - overall line structure
  - field separator format

- Output format must remain exactly:
  <hash>|<author>|<email>|<new timestamp>|<professional message>

- Do NOT add explanations, comments, markdown, headers, or extra lines.
- Output ONLY the transformed commit entries.`

	reqBody := chatRequest{
		Model: "openai/gpt-oss-120b:free",
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: string(in)},
		},
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var cr chatResponse
	if err := json.Unmarshal(respBody, &cr); err != nil {
		log.Fatal(err)
	}

	if len(cr.Choices) == 0 {
		log.Fatalf("empty response from OpenRouter: %s", string(respBody))
	}

	lines := strings.Split(cr.Choices[0].Message.Content, "\n")
	var cleaned []string
	for _, l := range lines {
		if strings.Count(l, "|") == 4 {
			cleaned = append(cleaned, l)
		}
	}
	output := strings.Join(cleaned, "\n")

	if err := ioutil.WriteFile("updated_commits.txt", []byte(output), 0644); err != nil {
		log.Fatal(err)
	}
}

var start string
var end string

func Run(startD string, endD string) {
	start = startD
	end = endD

	inFile := "edited_commits.txt"
	f, err := os.Open(inFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	outFile := "updated_commits.txt"
	out, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	
	dates := generateDates(lenLines(inFile))

	// Creates line reader for input file, Reads line-by-line.
	scanner := bufio.NewScanner(f)
	f.Seek(0, 0)
	
	idx := 0

	// writting in the updated file
	w := bufio.NewWriter(out)
	scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// split file messages in parts by Pipe '|'
		parts := strings.SplitN(line, "|", 5)
		
		// use the dates from time slice we generated
		newDate := dates[idx].Format("2006-01-02 15:04:05 +0530")

		// improve message
		msg := parts[4]
		newMsg := fmt.Sprintf("update: %s", strings.ToLower(msg))

		// new file lines generated
		resultOutline := fmt.Sprintf("%s|%s|%s|%s|%s\n", parts[0], parts[1], parts[2], newDate, newMsg)
		w.WriteString(resultOutline)
		
		idx++
	}

	w.Flush()
}

func lenLines(path string) int {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// Creates scanner for reading file line-by-line.
	s := bufio.NewScanner(file)
	count := 0

	// Reads next line repeatedly.
	// returns true if next line exist, else false 
	for s.Scan() {
		count++
	}

	return count
}

// generates n timestamps, basically the number of commits
// returns slice (vector array) of time.Time
func generateDates(n int) []time.Time {
	start, _ := time.Parse("2006-01-02 15:04:05 -0700", start)
	end, _ := time.Parse("2006-01-02 15:04:05 -0700", end)

	// Subtract and get the total time difference bw the start and the end dates
	durations := end.Sub(start)

	// Why n-1?
	// Because - 5 dates create 4 gaps
	gaps := durations / time.Duration(n-1)

	// make empty time slice of that size
	dates := make([]time.Time, n)
	for i := 0; i < n; i++ {
		// from start date add, time.Duration(i) * gaps time to the start
		dates[i] = start.Add(time.Duration(i) * gaps)
	}

	return dates
}
