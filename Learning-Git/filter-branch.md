# filter-branch

`git filter-branch` is basically:

> "Replay all commits, modify something during replay, create a new history."

Core idea:

Git commits are immutable objects.

A commit contains:

* tree snapshot
* parent commit hash
* author metadata
* committer metadata
* commit message

Because parent hashes are embedded inside commits:

If one commit changes:

```text
ALL child commits must also change
```

That is why history rewriting cascades.

---

Example history:

```text
A -> B -> C
```

Suppose you change author email in `B`.

Git cannot modify `B`.

Instead:

```text
B becomes B'
```

But `C` points to old `B`.

So Git must recreate `C` too:

```text
A -> B' -> C'
```

---

What `filter-branch` does internally:

For every commit:

```text
checkout commit state
↓
apply your filters
↓
create entirely new commit object
↓
reconnect parent pointers
↓
move branch ref to new history
```

Very similar to:

```text
map(old_commit -> new_commit)
```

---

Your command:

```bash
git filter-branch -f --env-filter '...' main
```

Meaning:

| Part            | Meaning                            |
| --------------- | ---------------------------------- |
| `filter-branch` | rewrite commit history             |
| `-f`            | overwrite previous rewrite backups |
| `--env-filter`  | modify commit metadata env vars    |
| `main`          | rewrite only main branch           |

---

`--env-filter`

This is literally shell code executed once per commit.

Git temporarily exports commit metadata:

```bash
GIT_AUTHOR_NAME
GIT_AUTHOR_EMAIL
GIT_AUTHOR_DATE
GIT_COMMITTER_NAME
GIT_COMMITTER_EMAIL
GIT_COMMITTER_DATE
GIT_COMMIT
```

Then your shell script can modify them.

Example:

```bash
git filter-branch --env-filter '
if [ "$GIT_AUTHOR_EMAIL" = "old@gmail.com" ]; then
    export GIT_AUTHOR_EMAIL="new@gmail.com"
fi
' main
```

Meaning:

```text
for every commit:
    if author email matches:
        replace it
    create new commit
```

---

Important distinction:

Author:

```text
person who originally wrote code
```

Committer:

```text
person who actually created/applied commit
```

They can differ during:

* rebases
* cherry-picks
* patches
* maintainer workflows

---

`refs/original`

After rewrite:

Git keeps backups:

```text
refs/original/refs/heads/main
```

This points to OLD history.

Why?

Safety.

So you can recover if rewrite destroys history.

---

Problem:

Even after rewrite:

```text
old commits still reachable
```

because backup refs still point to them.

So garbage collection cannot remove them.

---

That is why:

```bash
git update-ref -d refs/original/refs/heads/main
```

is used.

Meaning:

```text
delete backup pointer
```

Now old commits may become unreachable.

---

Then usually:

```bash
git reflog expire --expire=now --all
git gc --prune=now
```

Meaning:

1. remove reflog references
2. run garbage collection
3. physically remove unreachable commits

---

Production warning:

`git filter-branch` is old and slow.

Modern replacement:

```bash
git filter-repo
```

Much faster and safer.

Git itself officially discourages heavy `filter-branch` usage now.

---

Real-world usage:

1. Remove secrets accidentally committed
2. Rewrite author emails across org migration
3. Remove huge files from history
4. GDPR/legal cleanup
5. Monorepo extraction
6. OSS contribution cleanup

---