---
status: approved
granted: 2026-10-13
action: report a child process that a signal ended by the signal's name, number and core state in the Delivery Queue's item-binary notice and step errors and in every test helper that runs a child Roundfix binary and returns its exit code, make Verification Capacity read cancellation after an immediate grant so a cancelled Run starts no Verification command or pre-work probe, and describe the Signal Exit in the roundfix skill's deliver reference, the deliver user guide and the glossary with the recorded version raise
consuming: 0269-a-killed-child-names-its-signal-and-a-cancelled-run-starts-nothing
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - internal/cli/baseline_release_gate_test.go
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0269

On 2026-09-30 the maintainer asked for unattended work through every release
of the program, and on skills: "considere autorizado a ajustar todas as skills
se necessário". On 2026-10-11, answering "Grants dos caminhos governados" for
cycle AE, the maintainer wrote "Concedo" for the Governed Paths each Spec
declares. On 2026-10-13 the maintainer approved cycle AG (AskUserQuestion at
about 17:20 BRT: six Specs, order B, A, C, D, E, F, and for F6 "Só nomear o
sinal"), then answered "Pode iniciar" at about 17:25 BRT to start them under the
same operating rules as cycles AE and AF, which carry the same operations and
the same grant for the Governed Paths each Spec declares. The Makefile and the
CI workflows are not covered in AG. Cycle AG item D is this Spec.

The governed set was measured on the authoring branch at `9963ffd3` with
`speccheck.GovernedPath`, through a `go test -overlay` probe in
`internal/speccheck` that wrote nothing to the repository, against every path
task_01 to task_04 declare and every planning file. Four are governed.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/deliver.md` covers the `deliver`
  Behavior Surfaces whose source `internal/cli/deliver_workflow.go` changes;
  ADR-0256's authoring rule and `docs/agents/cli.md` require the skill to
  change with the command, and it gains the Signal Exit form of the fallback
  notice and of a step failure.
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  owned skill's version fields, which the record command raises with any
  content change.
- `internal/cli/baseline_release_gate_test.go` holds `runBaselineReleaseCLI`,
  a helper that runs a child Roundfix binary and returns
  `ExitError.ExitCode()`, which is `-1` with no signal named for a child a
  signal ended; task_03 moves it onto `testfixture.ChildExitCode`. Only that
  helper changes; no release-gate assertion changes.

## What is not governed

These paths are ordinary: `skills/roundfix/references/deliver.md`,
`skills/testdata/owned-skill-versions.json`,
`docs/user-guide/commands/deliver.md`,
`docs/references/behavior-surfaces.json`, `CONTEXT.md`,
`docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md`, the new
`internal/childexit` package, `internal/testfixture/child_exit.go`,
`internal/cli/deliver_workflow.go`, the other `internal/cli` test files task_02
and task_03 declare, `internal/daemon/task_engine.go` and
`internal/daemon/verification_capacity_cancellation_test.go`.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill version and
the Behavior Surface Record after the approved edits and add no source path.

```yaml
command: make skills-sync
```

`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
rewrites the version lines of `.agents/skills/roundfix/SKILL.md` and its mirror
and `skills/testdata/owned-skill-versions.json`, and
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
re-records the fingerprint of the `deliver` user-guide surface in
`docs/references/behavior-surfaces.json`. No Baseline module changes, so
neither `make baseline-digests` nor the module version record command is
declared. No record or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`,
  `.roundfixrc.yml`, lint, formatter or test-runner configuration, any
  Baseline module, or the `### QA settlement` section of any skill.
- No `sudo`, no process trace, and no signal sent by a test to a process it
  did not start; the fixtures signal only themselves.
- No test, Verification command or QA row opens a network connection, reads
  or writes the real `~/.roundfix`, or reads or writes an adopter repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
