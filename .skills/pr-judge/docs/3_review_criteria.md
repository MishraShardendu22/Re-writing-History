# Review Criteria

Judge every category relative to the repository, not against a universal checklist.

Group related concerns together when scoring or explaining findings:

- Correctness and regression risk.
- Necessity and duplication.
- Architecture, module boundaries, and dependency direction.
- Maintainability, readability, and code cleanliness.
- Best practices, including KISS, SOLID, DRY, and YAGNI, but only when they match the repository's style.
- Pattern usage and repository consistency.
- Testing quality and coverage.
- Documentation and comments.

Use local evidence to decide whether a concern is real. A change that would be textbook-correct in another project can still be wrong here if it breaks the repository's existing conventions.