# Examples

## Repository Requiring DCO

The repository's CI rejects commits without a `Signed-off-by` line. The correct advice is to amend the commit with `git commit -s --amend`, not to suggest unrelated workflow changes.

## Repository Requiring CLA

The PR is blocked by a CLA bot. The repository's docs say contributors must sign the agreement in the portal before review. The right response is to complete the CLA flow and wait for the bot to clear the status.

## Repository Requiring Signed Commits

The branch protection rules require signed commits. The contributor should use the signing method the repository documents, such as `git commit -S`, and verify the signature status before opening the PR.

## Example Signed-off-by Trailer

```text
Signed-off-by: Dev Contributor <dev@example.com>
```

## Example Signed Commit

```bash
git commit -S -s -m "Fix request parsing in handler"
```

## Example Contribution Checklist

- Branch name matches the repository convention.
- Commit message matches recent history.
- DCO sign-off is present if required.
- Commit signature is present if required.
- CLA is complete if required.
- CI checks pass.

## Example PR Readiness Review

The repository uses direct commits, DCO sign-off, and a minimal PR template. This change is ready only after the contributor adds the sign-off line and confirms CI passes. Recommending a new branch model or extra signing workflow would be inconsistent with the repository's existing process.