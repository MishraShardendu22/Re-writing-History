# fsck

Filesystem Check.
Low-level Git integrity verification command.

Checks entire Git object database for:

* corrupted objects
* broken commit chains
* missing blobs
* dangling commits
* unreachable objects
* invalid refs

Equivalent idea:

```text
fsck in Linux filesystem
```

but for Git internals.

---

Basic usage:

```bash
git fsck
```

Example output:

```text
dangling commit a1b2c3
dangling blob d4e5f6
```

or:

```text
missing blob 91f31f...
broken link from commit ...
```

---

What Git checks internally

Git stores objects in:

```text
.git/objects/
```

Types:

* blob
* tree
* commit
* tag

`git fsck` validates:

* SHA integrity
* object existence
* object graph connectivity
* ref validity

---

Common object types

Blob:

```text
file content
```

Tree:

```text
directory structure
```

Commit:

```text
snapshot + metadata + parent pointers
```

Tag:

```text
named pointer
```

---

Very important outputs

`dangling commit`

```text
commit exists but no branch/tag points to it
```

Usually after:

* rebase
* reset
* filter-branch
* amend

Not necessarily bad.

Example:

```bash
git reset --hard HEAD~1
```

Old commit becomes dangling.

Still recoverable until GC.

---

`unreachable object`

```text
object cannot be reached from any ref
```

Usually garbage candidate.

---

`missing blob`

```text
file object missing/corrupted
```

Very bad.

Repo may be damaged.

---

`broken link`

```text
commit/tree references object that does not exist
```

Repo corruption.

---

Useful flags

`--full`

```bash
git fsck --full
```

Deep verification.

Modern Git often does this by default.

---

`--unreachable`

```bash
git fsck --unreachable
```

Show unreachable objects.

---

`--dangling`

```bash
git fsck --dangling
```

Show dangling objects.

---

`--no-reflogs`

Normally reflogs are considered reachable references.

This ignores reflogs.

Useful after history rewrite cleanup.

Example:

```bash
git fsck --no-reflogs
```

This is important in your script context.

Because after:

```bash
git reflog expire
git gc --prune=now
```

you want to ensure:

```text
old rewritten commits are truly gone
```

---

Production use cases

1. Repository corruption detection
2. Recovery after disk damage
3. CI validation
4. After history rewrite
5. Before migration
6. Git server maintenance

---

In your rewrite pipeline

You could add:

```bash
git fsck --full
```

after:

```bash
git gc --prune=now
```

to verify:

```text
repository is still structurally valid
```

and:

```bash
git fsck --unreachable
```

to confirm:

```text
old commits are actually unreachable
```

---

Example flow

Before rewrite:

```text
A -> B -> C
```

After rewrite:

```text
A' -> B' -> C'
```

Old commits:

```text
A B C
```

become dangling/unreachable.

`git fsck` can reveal them until:

```bash
git gc --prune=now
```

removes them permanently.

---

Important distinction

`git fsck`

```text
checks integrity
```

`git gc`

```text
cleans storage
```

`git reflog`

```text
tracks ref history
```

`git update-ref`

```text
manipulates refs directly
```

---

Key Terminology

* Git Object Database
* Blob Object
* Tree Object
* Commit Object
* Reachability
* Dangling Commit
* Unreachable Object
* SHA-1/SHA-256 Integrity
* Object Graph
* Repository Corruption
* Ref Validation
* Garbage Collection
* Reflog Reachability
