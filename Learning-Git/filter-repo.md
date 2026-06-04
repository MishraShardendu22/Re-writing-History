# filter-repo

`git filter-branch` is considered old, slow, and dangerous for large rewrites.

Git itself warns against heavy usage now.

Modern replacement:

```bash id="e0bdc7"
git filter-repo
```

Official repo:

[git-filter-repo](https://github.com/newren/git-filter-repo?utm_source=chatgpt.com)

---

Why `filter-branch` is considered bad:

### 1. Extremely slow

`filter-branch`:

* checks out commits repeatedly
* uses shell filters
* recreates commits inefficiently

Large repos can take:

* minutes
* hours

---

### 2. Error-prone

Easy to accidentally:

* corrupt refs
* miss tags
* leave backup refs
* keep secrets reachable

---

### 3. Complex cleanup

Usually requires:

```bash id="yvf74u"
git update-ref -d refs/original/...
git reflog expire --expire=now --all
git gc --prune=now
```

---

### 4. Poor UX

Shell-based filters are awkward:

```bash id="0f8z4h"
--tree-filter
--index-filter
--env-filter
--msg-filter
```

---

`git filter-repo` solves this.

Advantages:

| Feature                 | filter-branch  | filter-repo  |
| ----------------------- | -------------- | ------------ |
| Speed                   | slow           | very fast    |
| Safety                  | risky          | safer        |
| Cleanup                 | manual         | automatic    |
| Large repos             | painful        | designed for |
| UX                      | shell-heavy    | cleaner      |
| Official recommendation | deprecated-ish | recommended  |

---

Example:

Remove file from history.

Old:

```bash id="pb6hzt"
git filter-branch --tree-filter 'rm -f secret.txt' -- --all
```

Modern:

```bash id="l8sdfu"
git filter-repo --path secret.txt --invert-paths
```

Much simpler.

---

Rename email globally:

```bash id="0s4o9m"
git filter-repo --mailmap my-mailmap.txt
```

---

Install on Fedora:

```bash id="3k2trm"
sudo dnf install git-filter-repo
```

or:

```bash id="6j8stt"
pip install git-filter-repo
```

---

Important:

`git filter-repo` is NOT bundled with Git by default on many systems.

It is an external tool maintained by Git contributors.

---

Git documentation itself basically says:

```text id="npppp0"
use filter-repo instead of filter-branch
```

for most rewrite operations.

---

Key Terminology:

* History rewriting
* Repository surgery
* Packfile rewriting
* Commit graph rewrite
* DAG rewrite
* Reachability cleanup
* Immutable objects
* Object rewriting
* Repository migration
* Metadata rewrite
