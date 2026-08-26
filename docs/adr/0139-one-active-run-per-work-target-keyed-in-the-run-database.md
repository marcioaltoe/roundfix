---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-26T00:00:00Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# One active Run per work target, keyed in the Run Database

A Run targets either an Open Pull Request's Review Issues or a Spec's Tasks, so the Run Database rejects a second Active Run for the same work target: `active_run_locks` is keyed by `(target_kind, target_key)` — `pr` with `<head_repository>#<pr_head_branch>`, `spec` with `<git_root>#<spec_slug>` — and `runs` carries an `implement` Kind with a nullable `spec_slug`, keeping one uniform table because Attach, stop, and the Run Event Journal treat all Runs uniformly by run id. A same-working-tree rejection guards the checkout while concurrent Runs would mutate it. Supersedes ADR-0005 through its members.

Consolidates ADR-0012 and ADR-0016 (2026-08-26); both are archived under docs/history/adr/.
