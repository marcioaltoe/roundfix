---
status: approved
granted: 2026-09-30
action: record each passing QA row's evidence snapshot when a QA pass closes, hand a failed pass's QA Report to the next pass, carry rows whose recorded inputs are byte-identical, keep rows that read the repository Verification, the Pull Request row or Task commits observed on every pass, and teach the carry in the qa-gate skill
consuming: 0202-a-qa-gate-that-reruns-only-stale-rows
paths:
  - .agents/skills/qa-gate/SKILL.md
  - skills/qa-gate/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0202

On 2026-09-30 the maintainer approved the program of Waves 7 to 9, which
includes the QA gate that reruns only stale rows. The same day the maintainer
wrote "considere autorizado a ajustar todas as skills se necessário", which
authorizes edits to every Roundfix-owned skill. The maintainer also answered
"Autorizar os dois" for the Baseline source and guides and for
`.roundfixrc.yml`, but this Spec uses only the skill authorization. The set
was measured with `GovernedPath` on the authoring branch at `30e8504f`,
through a `go test -overlay` probe that wrote nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md` — the qa-gate
  skill tells the gate Agent how to plan rows, declare inputs and treat a
  rerun. It must state that a carried row is kept and counts as passed, that
  every executed row declares its inputs including the new `commit_range`
  kind, and that the Daemon alone writes `evidence_snapshots` (task_04).
  `.agents/skills/` is canonical and `skills/` its mirror.

## What is not governed

`internal/speccheck/mechanical.go`, `internal/speccheck/report.go`,
`internal/daemon/task_engine.go`, `internal/daemon/task_context.go`,
`internal/agent/spec_prompt.go`, the new files
`internal/speccheck/evidence_record.go` and `internal/daemon/qa_prior_pass.go`,
the new test files this Spec's Tasks name under `internal/speccheck/`,
`internal/daemon/` and `internal/agent/`,
`skills/testdata/owned-skill-versions.json` and
`docs/user-guide/context-driven-development.md` are ordinary.
`internal/speccheck/mechanical_test.go` is governed and is not touched: its
carry-forward tests must pass unedited.

## Sanctioned regeneration

The repository-owned commands resolve their generated outputs. These
declarations record the regeneration that follows the approved skill edit and
add no source paths. A probe on a scratch copy of `30e8504f` showed that
`make baseline-digests` rewrites nothing for a qa-gate text and version
change.

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

## Limits

- No change to verdict rules, typed blocked counts, report naming, the
  `### QA settlement` section, archive eligibility or Task Carry-Forward.
- No edit to the Makefile, to lint, formatter or test-runner configuration,
  to a CI workflow or to `go.mod`.
- The live Run Database under `~/.roundfix` is never opened for writing by a
  Task or the gate, and no test reaches GitHub or a provider.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
