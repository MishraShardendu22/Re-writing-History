# Review Process

1. Identify the PR's goal and whether the change is actually needed.
2. Inspect the touched code and the nearest surrounding implementation.
3. Check whether similar behavior already exists somewhere else in the repository.
4. Compare the change against the repository's architecture and abstraction level.
5. Check maintainability, complexity, readability, tests, and documentation.
6. Classify issues by severity and separate required fixes from suggestions.

Questions the review must answer:

- What problem is this PR solving?
- Is the change necessary, or is it duplicating existing behavior?
- Does it fit the repository architecture and module boundaries?
- Does it match the local coding style and abstraction level?
- Does it add complexity or technical debt?
- Is the testing strategy consistent with the repository?
- Are the comments and documentation appropriate for this project?
- Would maintainers likely accept this implementation as-is?