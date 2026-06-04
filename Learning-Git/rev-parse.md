# git rev-parse

`git rev-parse` is an internal Git reference parser and resolver.
It converts human-readable Git references into actual internal object values.

Examples of refs:

* `HEAD`
* branch names
* tags
* commit hashes
* `HEAD~1`
* `origin/main`

Git internally stores references like:

```bash
refs/heads/main
refs/tags/v1.0
```

`rev-parse` resolves and normalizes them.

Example:

```bash
git rev-parse HEAD
```

Returns:

```bash
e7a1c9f1c2...
```

Meaning:

* Resolve `HEAD`
* Output the full commit SHA-1/SHA-256 hash

---

Your example:

```bash
git rev-parse --abbrev-ref HEAD
```

Flow:

1. `HEAD`

   * special pointer to current checked-out commit/branch

2. Normally Git resolves it to:

   ```bash
   refs/heads/main
   ```

3. `--abbrev-ref`

   * strips the full ref path
   * gives human-readable short ref

Result:

```bash
main
```

instead of:

```bash
refs/heads/main
```

---

Important behavior:

If detached HEAD:

```bash
git checkout e7a1c9f
```

then:

```bash
git rev-parse --abbrev-ref HEAD
```

returns:

```bash
HEAD
```

because you are no longer on a branch ref.

---

Real-world usage:

Shell scripts:

```bash
branch=$(git rev-parse --abbrev-ref HEAD)
```

CI/CD:

* detect current branch
* compare refs
* validate deployment branch

Git hooks:

* pre-push
* pre-commit
* release automation

---

Other useful variants:

```bash
git rev-parse --show-toplevel
```

Returns repo root path.

---

```bash
git rev-parse --git-dir
```

Returns `.git` directory location.

---

```bash
git rev-parse HEAD~1
```

Resolves previous commit hash.

---

```bash
git rev-parse origin/main
```

Resolves remote branch commit hash.

---