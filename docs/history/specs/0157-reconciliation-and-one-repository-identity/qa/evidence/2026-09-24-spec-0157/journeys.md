# Spec 0157 QA journey evidence

Build: `972bbe8f42ebc60d78f22e2cb975448b0a090260`

## R2 — refuse an unproven cleanup candidate

- Fixture: `/private/tmp/roundfix-0157-r2.2U297C`
- Command: `HOME=... /private/tmp/roundfix-0157-qa reconcile --apply --format json`
- Exit: `1`, without a panic.
- Result: `action: preserve`; `refusalReason: apply Run Branch candidate
  "roundfix/run-run_20260924T131549Z_aad3266ff0df9bf6": worktree revalidation
  returned no evidence`.
- Fresh Git confirmation: the Run Worktree remained listed at `b4798c1`, and
  `refs/heads/roundfix/run-run_20260924T131549Z_aad3266ff0df9bf6` still resolved
  to `b4798c14c90abd6bdcd3867fc926c4879c35a365`.

## R3 — archived copy keeps report identity

- Fixture: `/private/tmp/roundfix-0157-r3.FHe5hb`
- Command: `HOME=... /private/tmp/roundfix-0157-qa reconcile
  run_20260924T131550Z_a11a0946c2ac10f0 --format json`.
- Exit: `0` on two fresh runs.
- Result: `classification: superseded`; `supersedingReport:
  docs/history/specs/qa-0157-fixture/qa/qa-report-2026-07-28.md`.
- Fresh Git confirmation: `git show` read the Run-side active report as `fail`
  and the target-side archived report of the same name as `pass`.

## R4 and R5 — one identity and earlier artifacts

- Fixture: `/private/tmp/roundfix-0157-r45.EmtCfm`.
- Native Git evidence: `git worktree list --porcelain` reported main and linked
  at `4db52b4`; from linked, `git rev-parse --git-common-dir` returned
  `/private/tmp/roundfix-0157-r45.EmtCfm/main/.git`.
- Identity result: main, linked and independently computed pre-change main
  Artifact Directories all resolved to `.../artifacts/db203fe68217fa91`.
- Main-checkout `roundfix runs list --state all --limit 0` listed current Run
  `run_20260924T131550Z_cb47c44565fbe4bb` and earlier Run
  `run_20260924T131550Z_fdb404ffe4395049`.
- A fresh store read kept the earlier Run at Artifact Directory
  `.../artifacts/9e452eeeea9249b6` and read its marker as `earlier worktree
  artifact\n`.

## R6 — promise rule

`GOCACHE=/private/tmp/roundfix-0157-qa-gocache go run ./cmd/roundfix spec check
0157-reconciliation-and-one-repository-identity --strict` exited `0` on three
current-source runs and printed `No findings`.

## Focused regression checks

The focused `go test -count=1` selections for `internal/worktree`,
`internal/config` and `internal/store` all exited `0`:

```text
ok roundfix/internal/worktree 1.288s
ok roundfix/internal/config 0.326s
ok roundfix/internal/store 0.705s
```
