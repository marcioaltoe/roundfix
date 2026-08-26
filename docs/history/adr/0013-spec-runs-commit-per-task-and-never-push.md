---
status: superseded # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-07-05T22:17:04Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null
superseded_by: ADR-0138
---

# Spec runs commit per task and never push

Spec Runs keep ADR 0001's ownership split: the Daemon creates one commit per successfully verified Task — the code changes plus the updated task file — with a Conventional Commits message derived from the task frontmatter type and title. Spec Runs work on the user's current branch (Preflight Validation rejects the repository default branch), never push, and never open pull requests: Final Push remains a review-Run concept, and handing the branch to a pull request is the user's explicit decision after the Run.
