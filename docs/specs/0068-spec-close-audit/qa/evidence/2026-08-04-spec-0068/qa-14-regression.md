# QA-14 — Existing command regressions

Status: pass.

- Built `bin/roundfix spec check 0068-spec-close-audit` exited 0 with no
  finding; its two skips are optional-section presence skips.
- The fresh seven-test Spec Audit CLI selection includes established clean
  `spec check` behavior.
- `rtk go test ./internal/worktree ./internal/cli -count=1` exited 0 with
  1,041 tests, covering dirty, ambiguous, stale-evidence, unintegrated,
  dry-run, mixed-result, and idempotent reconcile behavior.
- Full repository Verification passed on the same build.
