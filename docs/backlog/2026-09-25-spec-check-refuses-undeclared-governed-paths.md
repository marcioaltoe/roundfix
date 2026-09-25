---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# `spec check` passes a Task that will touch a governed path nobody declared

## Symptom

A Task edits a governed path (Makefile, `.roundfixrc.yml`, `internal/cli/cli_test.go`, `internal/speccheck/mechanical_test.go`) that its `_authorization.md` and bounded-file rows never name. The daemon's governed audit refuses the commit only after the agent finished the work, and the Run ends Unresolved. Five of the sixteen Unresolved Runs in Specs 0155–0169 had this cause.

## Where

`internal/speccheck` (authoring gate) against `internal/speccheck/governed.go` (`GovernedPath`); Task `paths:` and Verification commands in `task_NN.md`.

## Expected

`spec check` resolves every path a Task declares or its Verification names through `GovernedPath`, and refuses with a stable code when that governed path is missing from `_authorization.md` and the bounded-file rows. The refusal happens before dispatch.

## Evidence

Runs of 0155 (Makefile, `.roundfixrc.yml`, `ci-verify.yml`), 0159 (`archive_layout_characterization_test.go`), 0162 and 0167 (`cli_test.go`), 0163 (`mechanical_test.go`); efficiency diagnosis in secondbrain `raw/roundfix/2026-09-25-sequencia-de-eficiencia.md`.
