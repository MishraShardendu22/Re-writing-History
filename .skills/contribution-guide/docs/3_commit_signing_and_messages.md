# Commit Signing and Messages

Treat commit signing, sign-off, and message format as separate checks and only combine them when the repository requires it.

Topics that often belong together:

- Signed-off-by trailer format.
- `git commit -s` for DCO sign-off.
- `git commit -S` for cryptographic signing.
- `git commit -S -s` when both are required.
- Commit subject style and body expectations.
- GPG or SSH signing, depending on the repository.

Infer the message style from recent history. If the repository uses short imperative subjects, follow that. If it uses Conventional Commits or issue-key prefixes, follow that instead.