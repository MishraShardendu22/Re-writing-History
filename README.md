# Git Commit Metadata Rewriter

## Overview

The Git Commit Metadata Rewriter is a specialized utility engineered in Go to retroactively modify Git commit metadata, including author names, email addresses, and timestamps, for cloned GitHub repositories. The system integrates with an AI model to generate realistic, non-uniform commit distributions, subsequently rewriting the repository's history to reflect the updated timelines.

## Live Resources

- **Video Demonstration**: [YouTube](https://www.youtube.com/watch?v=jHUGKcj5OwE)
- **Commit Testing Sample**: [Test Repository](https://github.com/ShardenduMishra22/Dhvani-Commit-Tester)

## Core Capabilities

- **Automated Repository Management**: Seamlessly clones public GitHub repositories for localized processing.
- **AI-Driven Timestamp Generation**: Leverages the OpenRouter API to programmatically generate realistic commit date distributions.
- **History Modification**: Utilizes Git filter-branch operations to systematically rewrite commit history with the newly generated metadata.
- **Automated Synchronization**: Executes force-pushes to synchronize the rewritten history with the remote repository origin.
- **State Cleanup**: Automatically purges temporary files and local repository clones upon process completion to maintain system hygiene.

## Operational Workflow

1. **Initialization**: Clones the specified target GitHub repository to the local environment.
2. **Extraction**: Exports the existing commit log to `edited_commits.txt`.
3. **Generation**: Interfaces with the AI model to establish new, realistic commit dates within a defined timeframe.
4. **Data Mapping**: Compiles the updated commit metadata into `updated_commits.txt`.
5. **Execution**: Rewrites the Git history applying the new metadata parameters.
6. **Synchronization**: Force-pushes the modified commit history to the remote source.
7. **Cleanup**: Removes all local artifacts and cloned data.

## Installation and Configuration

### 1. Repository Setup

Clone the repository and resolve dependencies:

```bash
git clone <repository-url>
cd Hackathon-Time-Script
go mod tidy
```

### 2. API Authentication

An OpenRouter API key is required for timestamp generation.
Create a `.env` file in the project root directory:

```env
API_KEY=your_openrouter_api_key_here
```

### 3. Target Configuration

Modify `main.go` to define the target repository and the desired date range parameters:

- `repo`: Set to the target repository identifier (e.g., `MishraShardendu22`).
- `start` / `end`: Define the boundaries for the newly generated commit dates.

### 4. Execution

```bash
go run main.go
```

## AI Integration Details

The system relies on the OpenRouter API to produce plausible commit timelines. The underlying prompt is strictly constrained to alter timestamps while preserving original author attributions and commit messages. Failure to provide a valid API key will result in immediate process termination.

## Troubleshooting

- **Missing API_KEY**: Verify the `.env` file is present and properly formatted.
- **Authentication Failures**: Ensure SSH keys are correctly configured with GitHub and that the executing environment holds sufficient push permissions for the target repository.
- **Empty AI Responses**: Confirm API key validity and network connectivity. An index out of range panic typically indicates a malformed or empty API response.
- **Protected Branches**: Force-push operations will fail on protected branches. Branch protection rules on the remote must be temporarily disabled prior to execution.

## System Architecture

```text
Hackathon-Time-Script/
├── main.go                # Primary execution script
├── util/
│   ├── clone.go           # Repository cloning and log extraction logic
│   ├── run.go             # AI model integration and timestamp generation
│   └── edit.go            # History rewriting and synchronization procedures
├── edited_commits.txt     # Extracted original commit log (Generated)
├── updated_commits.txt    # AI-modified commit log (Generated)
├── .env                   # Environment configuration (Ignored in version control)
├── go.mod                 # Go module dependencies
└── go.sum                 # Go module checksums
```

## Contributing

Review `CONTRIBUTING.md` for detailed contributor guidelines and procedures. Adherence to the `CODE_OF_CONDUCT.md` is strictly enforced.

## Security

For instructions on disclosing security vulnerabilities, please refer to `SECURITY.md`.

## License

Distributed under the MIT License. See `LICENSE` for further details. This project was developed to observe GitHub's visualization of commit histories and the systemic effects of rewritten repositories.
