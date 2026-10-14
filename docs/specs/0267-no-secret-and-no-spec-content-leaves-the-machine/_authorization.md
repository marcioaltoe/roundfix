---
status: approved
granted: 2026-10-13
action: add a Secret Scan that refuses a QA Report commit and every Roundfix push carrying a Secret Match, withhold the Judge Key Variables from every Agent Session and Verification command, and describe both in the qa-gate and roundfix skills, the QA prompt, the user guides and the glossary, with the recorded skill version raises
consuming: 0267-no-secret-and-no-spec-content-leaves-the-machine
paths:
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/implement.md
  - .agents/skills/roundfix/references/runtime.md
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0267

On 2026-10-13 the maintainer approved cycle AG through AskUserQuestion at
about 17:20 BRT (six Specs, in the order B, A, C, D, E, F), then answered
"Pode iniciar" at about 17:25 BRT to start them under the operating rules of
cycles AE and AF. Those rules carry the 2026-10-11 answer "Concedo" for the
Governed Paths each Spec declares and the operations listed above; the
Makefile and the CI workflows are not covered. On 2026-09-30 the maintainer
had granted the skills ("considere autorizado a ajustar todas as skills se
necessário"). Cycle AG item A is this Spec: QA evidence carries no secret and
no Spec content leaves the machine.

The governed set was measured on the authoring branch at `9963ffd3` with
`speccheck.GovernedPath`, through a `go test -overlay` probe in
`internal/speccheck` that wrote nothing to the repository, against every path
the Tasks declare and every planning file. Seven are governed.

## Why each governed path is unavoidable

- `.agents/skills/qa-gate/SKILL.md` holds the QA contract the QA Agent
  follows; task_04 adds the rule that evidence records environment variable
  names, never values, under its own heading.
- `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/implement.md` and
  `.agents/skills/roundfix/references/runtime.md` describe the Delivery Queue
  park, the QA Report commit and the Agent Session environment, which
  task_01, task_02 and task_03 change.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` and
  `skills/qa-gate/SKILL.md` carry the skills' version fields, raised with
  their content.

## What is not governed

The Go sources and tests under `internal/secretscan`, `internal/judge`,
`internal/agent`, `internal/daemon` and `internal/delivery`; the reference
mirrors under `skills/roundfix/references/`;
`skills/testdata/owned-skill-versions.json`;
`docs/references/behavior-surfaces.json`; the user guides
`docs/user-guide/commands/deliver.md`, `docs/user-guide/commands/implement.md`
and `docs/user-guide/configuration.md`; `CONTEXT.md`; and
`docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md`.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill versions
and the Behavior Surface Record. This declaration records the regeneration
that follows the approved edits and adds no source path.

```yaml
command: make skills-sync
```

`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
records the raised skill versions, and
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
re-records the fingerprints of the three changed user-guide surfaces. No
Baseline module changes, so no module version is recorded; on a scratch clone
at `9963ffd3`, `make baseline-digests` rewrote nothing after an edit of the
`qa-gate` and `roundfix` skills on 2026-10-13. No record or mirror is
hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`,
  `.roundfixrc.yml`, the lint, formatter or test-runner configuration, any
  Baseline module, or the `### QA settlement` section of any skill.
- No change to `internal/cli/review.go`, which Spec 0266 owns.
- No test, fixture, document or QA evidence holds a string that matches a
  Secret Pattern; tests build every sample at runtime.
- No test, Verification command or QA row reads a real credential, reads or
  writes the real `~/.roundfix`, or reaches the network, and no QA evidence
  records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
