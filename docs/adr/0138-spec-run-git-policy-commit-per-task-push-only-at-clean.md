---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-26T00:00:00Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# Spec Run git policy: one commit per verified Task, push only at Clean by opt-in

The Daemon creates one commit per successfully verified Task — the code changes plus the updated task file — with a Conventional Commits message derived from the task frontmatter, on the user's current non-default branch. A repository can opt a Spec Run into pushing that branch when, and only when, the Run ends Clean, through a Project Config key that defaults to off; every other outcome never pushes, and opening pull requests stays permanently outside Roundfix's scope. The original blanket "never push" protected early dogfooding; a per-repository opt-in keeps the decision in the repo's own config while the commit-per-Task ownership split of ADR-0001 stands unchanged.

Consolidates ADR-0013 and ADR-0021 (2026-08-26); both are archived under docs/history/adr/.
