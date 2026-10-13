---
status: approved
granted: 2026-10-12
action: make every test that waits on a child process, a goroutine or Git bound by real time only the budget it asserts, add the Real-Time Budget rule with a ceiling per test file, and state the rule in the repository guide and the glossary
consuming: 0265-tests-that-wait-on-a-child-do-not-race-the-machine
paths:
  - docs/agents/specific-repository.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0265

On 2026-10-12 the maintainer approved cycle AF by answering "Sim, como no AE"
to "As regras de operação do AE valem para o AF?". The cycle AE rules include
the grant of 2026-10-11, the answer "Concedo" to "Grants dos caminhos
governados", for the Governed Paths each Spec declares, and the operations
implement, commit, push, pull_request and merge. The Makefile and the CI
workflows are not covered in cycle AF, and this Spec changes neither.

The operator minted this Spec as the blocking Spec that the unattended operating
rules call for once Spec 0259's two Corrective Tasks were spent on tests that do
not survive machine load.

The governed set was measured on the authoring branch at `94fbf85b` with
`GovernedPath`, through a `go test -overlay` probe in `internal/speccheck` that
wrote nothing to the repository, against every file the Tasks declare. One of
them is governed.

## Why the governed path is unavoidable

- `docs/agents/specific-repository.md` gains, outside every setup marker, one
  repository rule (task_04): a test bounds by real time only the budget it
  asserts, and the Real-Time Budget rule in `internal/testfixture` holds each
  test file to its ceiling. Whoever writes the next test meets it there, beside
  the Parallel Test Package rule.

## What is not governed

The test files and the two production files the Tasks declare under
`internal/cli`, `internal/agent`, `internal/daemon`, `internal/worktree`,
`internal/store`, `internal/spec` and `internal/testfixture`; `CONTEXT.md`; and
`docs/adr/0271-a-test-bounds-by-real-time-only-what-it-asserts.md`, which this
planning change adds.

## Sanctioned regeneration

None. No skill, Baseline module, derived digest or record changes, so no
regeneration command runs.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any skill, any Baseline module, or the
  `### QA settlement` section of any skill.
- No production budget value changes; the two new dependencies default to the
  values used today.
- No test, Verification command or QA row reads or writes the real
  `~/.roundfix`, signals a process it did not start, uses `pkill` or `killall`,
  or reaches the network. The QA gate's synthetic load runs niced, records each
  burner's PID, and stops every burner by that PID when its window ends.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
