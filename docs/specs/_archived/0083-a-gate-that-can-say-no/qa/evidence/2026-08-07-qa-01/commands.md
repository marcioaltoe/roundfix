# QA 0083 command evidence

Build: `0cebcce5`.

## Authoritative gate

- `rtk make verify` — exit 0 after the sandbox-only GitHub denial was rerun in
  the full-access gate environment. The gate ran `go test -parallel 16 ./...`,
  the corpus budget, skill synchronization and contract checks, build, and
  `bin/roundfix spec check`; 26 Go packages passed and Spec 0083 had no
  consistency findings.
- `rtk go test ./... -count=1` — exit 0; `Go test: 3654 passed in 26 packages`.

## Gate says no

- `/tmp/roundfix-qa-0083-short-CDAGGy/repo`: injected
  `TestQAInducedShortFailure`; `rtk make verify` exited 2, named
  `roundfix/internal/spec`, and printed `QA-03 deliberate short failure`.
- `/tmp/roundfix-qa-0083-high-volume-xyKjHo/repo`: injected a test that logged
  300 lines and failed; `rtk make verify` exited 2. The retained full output
  named `QA-04 deliberate high-volume failure` and
  `FAIL roundfix/internal/spec`.
- `TestAuthoritativeGateReportsFailure` — exit 0. Its high-volume failure,
  short failure, passing package, and wrapper-masking mutation subtests all
  passed.

## Retired noisy gates

- `/tmp/roundfix-qa-0083-authoring-H7u2zp/repo`: added disposable ADR 9999 and
  archived Spec 9999; `rtk make verify` exited 0.
- The copied tree's verbose `TestCheckCorpusGolden` exited 0 and reported
  `SC-CONSTRAINT-MISSING: 324` for the archived corpus. The current tree's same
  command exited 0 and reported `SC-CONSTRAINT-MISSING: 320`, proving the
  historical count changed without gating.
- With twelve `openssl speed -seconds 45 -multi 12 sha256` workers active,
  `TestCheckCorpusBudget -count=3 -parallel=1 -v` exited 0. It logged
  1.15226875s, 1.089449292s, and 1.164038708s while preserving 82 Check
  operations across 82 Specs.
- `/tmp/roundfix-qa-0083-inefficiency-NSvaUy/repo`: duplicated the real
  `checkCorpusSpec` operation; `rtk make verify` exited 2 and reported 164 Check
  operations across 82 Specs, want at most 82.

## Retained gates

- `/tmp/roundfix-qa-0083-spec-check-qEgRte/repo`: added same-subject `MUST` and
  `MUST NOT` requirements. `roundfix spec check` exited 1 with
  `SC-REQUIREMENT-CONTRADICTORY`; `rtk make verify` exited 2 and the active
  corpus plus active-error check both named that code.
- `/tmp/roundfix-qa-0083-published-example-C31xj0/repo`: changed the README
  example to `roundfix baseline --qa-invalid-flag`. The focused parser test
  exited 1 and `rtk make verify` exited 2, both naming the invalid flag and
  `README.md`.
- `/tmp/roundfix-qa-0083-coverage-ib6rUz/repo`: renamed the recorded
  `TestCompareCoverageRecordsReportsMissingTest`. The focused invariant exited
  1 and `rtk make verify` exited 2, reporting both the renamed addition and the
  missing recorded identity.
- Current `TestCoverageEquivalence` — exit 0. `rtk rg --files docs | rtk rg
  'coverage-record\\.json$'` returned only
  `docs/references/coverage-record.json`. The exact archived-tree diff against
  Task 02's preimage, excluding the authorized moved record, exited 0.

## Stabilized tests

All three commands below ran while the same twelve-worker CPU load was active:

- Cancellation test `-count=20 -v` — exit 0; 20 parent tests and all 40
  subtests passed in 7.801s.
- Capacity and Daemon-status integrated test `-count=20 -v` — exit 0; all 20
  tests passed in 12.226s.
- Corpus budget `-count=3 -v` — exit 0 with the elapsed observations above.

The adjacent cancellation race run and the capacity race run both exited 0.
Source inspection found the changed waits at `<-promptStarted.done` and
`waitImplementAgentStarts`; neither target assertion races an Agent start
against a real-time deadline.

## Authorization and cleanup

- Authorization commit `b1670ed6` predates Task commits `20a162bc`, `9bf89139`,
  `f1744333`, `81ed8156`, `6ca5ee76`, and `0cebcce5`.
- `git diff-tree --no-commit-id --name-status -r` on each Task commit showed
  only its expressly bounded tooling paths plus its own Task file. No Task
  commit contains an authorization, prerequisite fix, consequent fix, derived
  pin, or unrelated path.
- Removed copies: `short-CDAGGy`, `high-volume-xyKjHo`, `authoring-H7u2zp`,
  `inefficiency-NSvaUy`, `spec-check-qEgRte`,
  `published-example-C31xj0`, and `coverage-ib6rUz`.
- Final `find /tmp -maxdepth 1 -name 'roundfix-qa-0083-*' -print` returned no
  path.
