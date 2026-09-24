# Spec 0157 QA rerun evidence

Build: `934c775ff03f720f0e170d95c70afffcf798992a`

Current-tree CLI: `roundfix 0.16.0 (934c775f-dirty, built 2026-09-24
16:33:32 -0300)` at `/private/tmp/roundfix-0157-rerun-qa`.

## R2 — refuse an unproven cleanup candidate

- Fixture: `/private/tmp/roundfix-0157-rerun-r2.rO8dzO`.
- Command: `HOME=... /private/tmp/roundfix-0157-rerun-qa reconcile --apply
  --format json`.
- Exit: `1`, without a panic.
- Result: the Run Branch candidate recorded `action: preserve` and
  `refusalReason: apply Run Branch candidate
  "roundfix/run-run_20260924T193358Z_68ec49712d0a41d5": worktree
  revalidation returned no evidence`.
- Fresh Git confirmation: `git worktree list --porcelain` still listed the Run
  Worktree at `4a2438dbec58fe1d99f0e4ba0d0ed66de03c9d5a`, and `git
  show-ref --verify
  refs/heads/roundfix/run-run_20260924T193358Z_68ec49712d0a41d5` resolved
  the Run Branch to the same head.

## R3 — same name requires the same content

- Fixture: `/private/tmp/roundfix-0157-rerun-r3.ocjo9p`.
- Different-content blobs: active
  `0374cb8c68d83d9c7155d9af8f73f386b561c8c8`, archived
  `af54a1011967d076aa0e8c8a972509fc92348f44`.
- Two fresh read-only reconciliation commands classified the Run
  `unintegrated`, returned no `supersedingReport`, and preserved it.
- After the archived report was committed with the active report's exact
  bytes, both Git blob reads returned
  `0374cb8c68d83d9c7155d9af8f73f386b561c8c8`.
- Two more fresh reconciliation commands classified the Run `superseded` and
  named
  `docs/history/specs/qa-0157-fixture/qa/qa-report-2026-07-28.md` as the
  superseding report.

## R4 and R5 — shared identity, checkout-local Run root and earlier artifacts

- Fixture: `/private/tmp/roundfix-0157-rerun-r45.b2AJaJ`.
- Main, linked and independently computed pre-change main Artifact Directories
  all resolved to
  `/private/tmp/roundfix-0157-rerun-r45.b2AJaJ/home/.roundfix/artifacts/0a9c7690aaba1cef`.
- A fresh Run-store reader listed current Run
  `run_20260924T193547Z_00f22f8b8bdf2abe` with `gitRoot` equal to the linked
  checkout and its Artifact Directory equal to the shared main identity.
- Two fresh `runs list --state all --limit 0` commands from main listed the
  current and legacy Runs. Two more from linked still listed both after the
  fixture gained an unreadable `.git/worktrees/unreadable/gitdir` entry.
- Native Git listed main and linked at
  `2c6fcc76f5e18385cc1fe0db885728d3ed2d5456`; from linked, `git
  rev-parse --git-common-dir` returned
  `/private/tmp/roundfix-0157-rerun-r45.b2AJaJ/main/.git`.
- The legacy Run remained bound to Artifact Directory
  `/private/tmp/roundfix-0157-rerun-r45.b2AJaJ/home/.roundfix/artifacts/bcbb7e7c496ac383`.
  Reading `earlier-run.txt` from that directory returned exactly `earlier
  worktree artifact\n`.

## R6 — promise rule

`/private/tmp/roundfix-0157-rerun-qa spec check
0157-reconciliation-and-one-repository-identity --strict` exited 0 from the
audited source and printed `No findings`. It explicitly reported that authored
Verification commands were not executed.

## Focused regression selection

The named `go test -count=1` selection for `internal/worktree`,
`internal/config` and `internal/store` exited 0:

```text
ok  roundfix/internal/worktree  1.759s
ok  roundfix/internal/config    0.463s
ok  roundfix/internal/store     1.359s
```
