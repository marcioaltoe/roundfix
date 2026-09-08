---
status: pending
created_at: 2026-09-08
updated_at: 2026-09-08
kind: finding
---

# Run naming — internal branches ignore the repository's selected prefix (2026-09-08)

The Baseline tells Agents to use the repository's selected branch prefix,
but Run and Task branches use a fixed namespace. A prefixed initial branch
cannot prevent the executor from creating additional nonconforming branches.

Source: Secondbrain `inbox/roundfix/_triaged/2026-09-08-branches-internas-do-run-ignoram-o-prefixo-do-repositorio.md`.
The original capture and its dated evidence remain intact there.

## 1. Observed behavior

- Symptom / evidence: `internal/store/store.go` defines `RunBranchPrefix = "roundfix/run-"` and validates persisted Run branches against it. `internal/worktree/worktree.go` derives Run/Task names from the same constant. Baseline `branch.prefix` renders guidance without configuring those names.
- Root cause: The Baseline decision and executor namespace are separate contracts with no shared effective naming policy.
- Action / suggestion: Route to provisional P6, repository identity and Run branch policy. Preserve recognition of existing Run resources when honoring the selected prefix. Pantheon's missing adoption manifest is a separate readiness issue.

The implementation group is provisional; no Spec has been assigned and no
implementation or terminal verification is claimed by this triage.
