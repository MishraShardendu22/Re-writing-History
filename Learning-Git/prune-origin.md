# `git remote prune origin`

## Purpose
Remove stale remote-tracking references.

## Command
```bash
git remote prune origin
```

## Example

Remote previously had:
```text
origin/dev
origin/test
```

If `test` deleted remotely:
```text
local stale reference remains
```

This command removes it.

## Internal Concept
Git caches remote branch references locally.

This command synchronizes remote-tracking refs.

---