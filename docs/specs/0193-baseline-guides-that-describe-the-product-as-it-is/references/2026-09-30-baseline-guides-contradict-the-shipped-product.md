---
type: fix
status: promoted
created: 2026-09-30
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
reason: null
---

# Baseline guides contradict the product they ship with

## Symptom

An audit of the Baseline modules against `main` on 2026-09-30 found clauses that tell an Agent something the product no longer does. Every adopter receives them as mandatory guidance on its next `roundfix baseline update`.

- **Loop order.** `clause.autonomous.loop-01-qa-once` says "archive and commit the candidate, apply the configured pre-PR review policy". The Delivery Queue reviews first and archives second. The clause never names the Delivery Queue, `roundfix reopen` or Task Carry-Forward.
- **Tooling Task scope.** `clause.spec.project-constraints-03-bounded-execution` says an authorized tooling Task "may mutate only its bounded repository-relative files and its own Task file" and fails on any other path. The audit judges only Governed Paths (ADR-0130), and the Daemon records an ordinary undeclared path instead of refusing it (ADR-0166).
- **Absent record.** `clause.core.tooling-commit-choreography` says an absent record "does not by itself refuse implementation, commit, or push". `roundfix deliver start` refuses a Spec whose record does not grant five operations.
- **A control that does not exist.** `clause.core.verification-two-tiers` says an untrusted source "requires an `execution_approvals` entry". No Go source reads that field.
- **Empty grant.** `clause.spec.project-constraints-02-tooling-authorization` does not say that `paths: []` grants operations only (ADR-0179).
- **Outside evidence.** `clause.spec.project-constraints-06-outside-evidence` says the row "never blocks the Spec". ADR-0104 holds Pull Request preparation at the gate until the row is satisfied.
- **ADR lifecycle.** `clause.context.adr-02-active-status` says "Only `accepted` is active", and the next clause says a legacy ADR without frontmatter is active.
- **Backlog status.** The Backlog contract has no `deferred` status. Twenty entries under `docs/history/backlog/` use it.
- **Repository-only records.** Four clauses cite this repository's own numbers to every adopter: `Spec 0078`, `ADR-0014`, `ADR-0080`, `ADR-0091`, `ADR-0092` and `ADR-0163`. An adopter's `ADR-0014` is a different decision.
- **Repository guide.** `docs/agents/specific-repository.md` still says nothing under a `_archived` tree is validated. The archive root is `docs/history/`.

## Where

`internal/baseline/assets/modules/autonomous-work.json`, `core.json`, `spec-workflow.json` and `context-workflow.json`; the guides rendered from them under `docs/agents/`; `docs/agents/specific-repository.md`.

## Expected

Each clause states what the shipped product does, in words an adopter can act on, and cites no record that exists only in this repository. A clause that describes a control the product does not enforce says so. Mechanical checks keep the loop order equal to the Delivery Queue's order and keep repository-only citations out.

## Evidence

- The Delivery Queue's stage order in `internal/delivery/engine.go`, and the action order its tests pin: `run`, `policy`, `review`, `archive`, `gate`, `publication`, `push`, `pull-request`, `checks`, `merge`.
- `git grep -l execution_approvals -- '*.go'` on `9e439dbb` returns nothing.
- The Secondbrain mirrors of two adopters, read on 2026-09-30: `projects/oraculum/mirror/docs/agents/autonomous-work.md` and `projects/vortex/mirror/docs/agents/autonomous-work.md` carry the old order, and all five adopter mirrors carry `Spec 0078`.
