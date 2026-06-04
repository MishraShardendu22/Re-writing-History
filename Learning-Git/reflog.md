# reflog

`git reflog expire` removes reflog entries.

Reflog is basically Git's:

```text id="lk1h1u"
local movement history
```

It tracks where refs previously pointed.

---

Example:

```bash id="kk1gt0"
git checkout main
git commit
git reset --hard HEAD~1
```

Even after deleting the commit from branch history:

```text id="kvqodq"
commit may still be recoverable
```

because reflog recorded:

```text id="q6ojpu"
HEAD moved from A to B
HEAD moved from B to C
HEAD moved from C to B
```

---

Reflog is local only.

Not shared with remote.

Stored under:

```text id="b5e7vl"
.git/logs/
```

Examples:

```text id="e40gx1"
.git/logs/HEAD
.git/logs/refs/heads/main
```

---

Your command:

```bash id="yzhax2"
git reflog expire --expire=now --all
```

Meaning:

| Part            | Meaning                          |
| --------------- | -------------------------------- |
| `reflog expire` | remove reflog entries            |
| `--expire=now`  | mark entries immediately expired |
| `--all`         | apply to every reflog            |

---

Internal effect:

Before:

```text id="0caw90"
HEAD reflog:
A -> B
B -> C
C -> D
```

After:

```text id="s8bf3x"
reflog entries deleted
```

So Git loses recovery references.

---

Why this matters after `filter-branch`:

Even if:

```bash id="bn1nyj"
git update-ref -d refs/original/...
```

removes backup refs,

reflog may STILL contain:

```text id="9vq6zu"
old commit hashes
```

meaning old history is still reachable.

---

Git garbage collection rule:

```text id="k48m0u"
reachable object = cannot delete
```

Reachable through:

* branches
* tags
* reflogs
* special refs

---

So cleanup sequence is usually:

```bash id="b6vnm6"
git filter-branch ...
git update-ref -d refs/original/refs/heads/main
git reflog expire --expire=now --all
git gc --prune=now
```

---

Flow:

```text id="slqlj6"
rewrite commits
↓
delete backup refs
↓
delete reflog references
↓
garbage collect unreachable objects
```

Only after all this can old commits disappear physically.

---

Important:

Without `git gc`:

```text id="l9lt8w"
objects may still exist in .git/objects
```

just unreachable.

---

Recovery example:

```bash id="m3d0xq"
git reflog
```

might show:

```text id="07du8s"
54ab12c HEAD@{1}: reset: moving to HEAD~1
```

Then:

```bash id="h2jlwm"
git checkout 54ab12c
```

recovers supposedly "deleted" commit.

That is exactly what `reflog expire` prevents.

---