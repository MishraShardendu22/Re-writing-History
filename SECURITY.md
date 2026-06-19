# Security Policy

## Reporting a Vulnerability

We take the security of this project seriously. If you discover a security vulnerability, please report it privately.

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, open a GitHub Issue with the label `security` and the maintainer will follow up privately.

### What to Include

When reporting a vulnerability, please include:
- A clear description of the issue
- Steps to reproduce the vulnerability
- Any potential impact or exploit scenarios
- Your contact information for follow-up

### What to Expect

- **Acknowledgment**: You will receive an acknowledgment of your report within 48 hours.
- **Status Updates**: We will provide updates on the investigation and remediation progress.

## Scope

The following areas are in scope for security review:

- **API Key Handling**: Secure storage and usage of the OpenRouter API key.
- **Git Operations**: Safety of `git filter-branch`, `git push --force`, and history rewriting operations.
- **Configuration**: Environment variable handling, `.env` file management, secret exposure.

## Best Practices

- **Never commit API keys** -- Use `.env` files (excluded via `.gitignore`) or environment variables.
- **Review rewritten history** -- Verify the output before force-pushing to remote repositories.
- **Use a test repository first** -- Test the script on a fork or backup before running on production repositories.

## Supported Versions

| Version       | Supported          |
|---------------|--------------------|
| Latest        | Active support     |

Only the latest release of this project is actively supported with security updates.

## Contact

For security-related matters, please open a GitHub Issue with the `security` label.


## Why build this ?

I built a small project to understand Git beyond the normal workflow. Mostly we only use high-level commands like commit, push, rebase, and merge. 

I wanted to understand how Git actually stores commits, authorship metadata, timestamps, and object relationships internally, I wanted to understand plumbing commands and not just porcelian commands. 

I experimented with changing commit metadata such as author information and commit dates using Git plumbing commands and history rewriting techniques. The goal wasn't to misrepresent work but to understand how Git's object model, commit hashes, history rewriting, and repository integrity work under the hood.

The project taught me about Git objects, SHA-based commit identity, DAG structures, reflogs, filter-repo/filter-branch style history rewriting, author vs committer metadata, and the implications of rewriting shared history.