---
name: pr-judge
description: Use when the user asks to review, evaluate, critique, or judge a pull request, code review, patch, commit, or software changes.
license: MIT
metadata:
  author: mishrashardendu22
  version: "1.0"
---

# PR Judge Skill

Purpose: review pull requests as a repository-aware engineering judge, not a generic linter.

## When to Use
Use when a user asks whether a PR is necessary, safe, maintainable, or ready to merge.

## Activation Criteria
Trigger on a request mentioning a pull request, patch, diff, review, or approval. Inspect the repository before forming any opinion.

## Repository Analysis Workflow
Read docs in order: `1_repository_analysis.md`, `2_review_process.md`, `3_review_criteria.md`, `4_scoring_and_output.md`, `5_examples.md`; infer the repository's architecture, conventions, testing style, and maintainers' preferences from evidence; compare the change against the repo's abstraction level, complexity, and style; validate correctness, tests, docs, and regression risk before judging.

## Validation Checklist
Ground every judgment in repository evidence, separate required changes from suggestions, and ensure the final score matches the documented rubric.