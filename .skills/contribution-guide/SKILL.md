---
name: contribution-guide
description: Use when the user asks how to contribute, submit code, create pull requests, follow contribution guidelines, become a contributor, or asks about the contribution workflow.
license: MIT
metadata:
  author: mishrashardendu22
  version: "1.0"
---

# Contribution Guide Skill
Purpose: help an agent satisfy repository contribution requirements before opening or updating a PR.

## When to Use
Use when a user needs help satisfying repository contribution rules or preparing commits, branches, signatures, or PR metadata.

## Activation Criteria
Trigger on contributing, opening a PR, pushing changes, or merge readiness. Inspect the repository before assuming any contribution requirement.

## Repository Analysis Workflow
Read docs in order: `1_repository_analysis.md`, `2_contribution_rules.md`, `3_commit_signing_and_messages.md`, `4_contribution_workflow.md`, `5_examples.md`; infer whether the repo uses DCO, CLA, signed commits, branch rules, or special CI gates; match guidance to the repo's actual workflow and terminology; highlight only requirements evidenced by the repository.

## Validation Checklist
Infer requirements from repository evidence, separate required checks from optional best practices, and ensure commit, signing, and branch advice match the workflow.