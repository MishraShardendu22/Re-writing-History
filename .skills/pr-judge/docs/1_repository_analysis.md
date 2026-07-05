# Repository Analysis

Always begin by inspecting the repository itself. Do not judge a PR until you know how this codebase is organized and what patterns it already uses.

Infer these from evidence in the repository:

- Project architecture and module boundaries.
- Abstraction level and layering.
- Dependency structure and ownership of data, validation, orchestration, and I/O.
- Coding conventions, naming conventions, and formatting habits.
- Testing philosophy and the level of coverage expected.
- Logging, configuration, and error-handling style.
- Dependency injection style, if any.
- Engineering philosophy: simplicity, extensibility, framework-driven design, or performance-driven design.
- Maintainers' actual style, not generic best practices.

Use nearby files, tests, scripts, and docs as evidence. If the repository has no strong convention in a category, say so explicitly instead of inventing one.