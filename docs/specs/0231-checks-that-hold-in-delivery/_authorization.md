---
status: approved
granted: 2026-10-05
action: let the detach fixture tests count a process-group member that is already exiting as ended, let CI's pull request job test the head merged with the default branch tip it fetches and record that tip, and let the Delivery Queue re-run a failed check that tested an older default branch instead of parking it
consuming: 0231-checks-that-hold-in-delivery
paths:
  - .github/workflows/ci-verify.yml
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0231

On 2026-10-05 the maintainer answered the cycle's questions through
AskUserQuestion. For the scope the answer was "0230 + S/T": deliver Spec 0230
and author this Spec from sources S and T. For tooling the maintainer gave a
named grant, "Concedo", covering the governed test files this Spec declares
(for example `internal/cli/implement_test.go` and the detach tests) and
`.github/workflows`, limited to the paths the Spec declares. The advisory
`roundfix spec judge` was allowed for this repository's Spec artifacts only,
run once on the finished Spec, its suggestions reported and never used as a
gate.

The governed set was measured with `GovernedPath` on the authoring branch at
`c5b2f27e`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare. Of them, only
`.github/workflows/ci-verify.yml` is a Governed Path; the detach tests,
`internal/cli/implement_test.go` and every `internal/delivery` file are
ordinary, so the test part of the grant bounds no path here.

## Why each governed path is unavoidable

- `.github/workflows/ci-verify.yml` — a re-run of a `pull_request` job keeps
  the original event's merge commit, so no queue behavior alone can make CI
  test the current default branch. The job must check out the head, merge the
  default branch tip it fetches, verify against that tip and record it as the
  `tested-base` annotation the queue reads (ADR-0236). task_02 changes this
  file and nothing else, so the tooling change is its own commit.

## What is not governed

`internal/cli/detach_fixture_group_darwin_test.go`,
`internal/delivery/check_rerun.go`, `internal/delivery/engine.go`, their
tests, `internal/delivery/stale_check_test.go` and
`docs/user-guide/commands/deliver.md` are ordinary.

## Limits

- No other workflow changes; the `push` path of `ci-verify.yml` keeps its full
  Verification, its budget and `make verify-docs`.
- No change to the Makefile, `go.mod`, the lint, formatter or test-runner
  configuration, any skill or any agent guide.
- No production change to detach, `internal/store` process inspection or the
  Run Database schema.
- No test or Verification command reaches GitHub; the queue's GitHub calls are
  exercised through scripted command runners.
- No change to archived Specs, existing QA Reports, `CONTEXT.md`,
  `CHANGELOG.md` or the `### QA settlement` section of any skill.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
