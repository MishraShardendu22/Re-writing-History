# gc - garbage collection

`git gc` = Git Garbage Collection.

It is Git's repository maintenance and cleanup system.

---

Git stores everything as objects:

```text id="pd6lc5"
commit
tree
blob
tag
```

Inside:

```text id="ygt6qv"
.git/objects/
```

Objects are immutable.

Even after deleting branches or rewriting history:

```text id="lln4si"
objects usually still physically exist
```

until garbage collection removes them.

---

Your command:

```bash id="u1nm0g"
git gc --prune=now --aggressive
```

Meaning:

| Part           | Meaning                                |
| -------------- | -------------------------------------- |
| `gc`           | run repository cleanup                 |
| `--prune=now`  | immediately delete unreachable objects |
| `--aggressive` | deeper compression optimization        |

---

Internal responsibilities of `git gc`:

### 1. Remove unreachable objects

Example:

```text id="rkh9gt"
A -> B -> C
```

Then:

```bash id="7yjl3x"
git reset --hard A
```

Now:

```text id="0cl4vf"
B and C unreachable
```

No refs point to them.

`git gc --prune=now` can physically delete them.

---

### 2. Pack loose objects

Git initially stores objects separately:

```text id="1u1nvh"
.git/objects/ab/c123...
.git/objects/ef/9123...
```

Thousands of tiny files become inefficient.

So Git creates:

```text id="j4mx1j"
.git/objects/pack/
```

Containing:

```text id="0m9w9x"
.pack
.idx
```

files.

This is called:

```text id="k1bg88"
packfiles
```

Much faster and smaller.

---

### 3. Delta compression

Git stores object differences.

Instead of:

```text id="k73o5m"
file_v1
file_v2
file_v3
```

fully duplicated,

Git may store:

```text id="5c7g4u"
base version
+
binary diffs
```

Very similar to compression dictionaries.

---

### 4. Cleanup temporary metadata

Removes:

* stale reflogs
* temporary packs
* orphan metadata

---

`--prune=now`

Normally Git keeps unreachable objects for safety.

Default grace period:

```text id="p70xsz"
about 2 weeks
```

Why?

To allow recovery from mistakes.

`--prune=now` skips waiting.

Meaning:

```text id="d5s9up"
delete unreachable objects immediately
```

---

Important distinction:

Unreachable ≠ deleted yet.

Example:

```bash id="ajgq3n"
git branch -D feature
```

does NOT instantly remove commits.

It only removes refs.

Commits become:

```text id="20j77d"
dangling/unreachable
```

Then GC later removes them.

---

`--aggressive`

More expensive optimization pass.

Git:

* recomputes deltas deeper
* tries stronger compression
* repacks more thoroughly

Effects:

| Effect    | Result      |
| --------- | ----------- |
| CPU       | much higher |
| RAM       | higher      |
| Speed     | slower      |
| Repo size | smaller     |

---

Usually used:

* large repos
* CI mirrors
* monorepos
* after massive history rewrites

Not needed regularly.

---

Very important sequence after history rewrite:

```bash id="65lgg6"
git filter-branch ...
git update-ref -d refs/original/refs/heads/main
git reflog expire --expire=now --all
git gc --prune=now --aggressive
```

Flow:

```text id="hbm0k9"
rewrite history
↓
remove backup refs
↓
remove reflog recovery pointers
↓
garbage collect unreachable objects
↓
physically erase old commits
```

---

Without GC:

Old commits may still physically exist:

```text id="j0c7k5"
inside .git/objects/
```

even though branch history no longer references them.

---

Real-world production usage:

1. Remove leaked secrets permanently
2. Shrink bloated repositories
3. Optimize CI clone performance
4. Cleanup after monorepo surgery
5. Git server maintenance
6. Large enterprise repo maintenance

---