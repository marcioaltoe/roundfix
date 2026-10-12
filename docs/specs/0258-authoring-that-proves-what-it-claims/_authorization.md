---
status: approved
granted: 2026-10-11
action: add four structural proof gaps to the Spec Consistency Check for every active Spec, state the four proof rules in the write-techspec and write-tasks skills, and document the codes in the roundfix command reference and the user guide
consuming: 0258-authoring-that-proves-what-it-claims
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-tasks/references/task-template.md
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/write-techspec/references/techspec-template.md
  - internal/docscontract/testdata/corpus-golden.json
  - internal/spec/archive_layout_characterization_test.go
  - internal/speccheck/coherence.go
  - skills/baseline_skill_contract_test.go
  - skills/roundfix/SKILL.md
  - skills/write-tasks/SKILL.md
  - skills/write-techspec/SKILL.md
  - skills/write-techspec/references/techspec-template.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0258

On 2026-09-30 the maintainer asked for unattended work through every release
of the program. Of the skills the maintainer said "considere autorizado a
ajustar todas as skills se necessário". Of Baseline sources, guides and
`.roundfixrc.yml` the maintainer said "Autorizar os dois". On 2026-10-11 the
maintainer answered "Grants dos caminhos governados" with "Concedo" for the
Governed Paths each Spec of cycle AE declares.

On 2026-10-11 the maintainer chose "Skills + gap no spec check" for this
Spec's source: the `write-tasks` and `write-techspec` skills gain the four
rules, and `roundfix spec check --strict` reports a missing section or line
structurally. On 2026-10-12 the maintainer decided that the new gaps hold
every active Spec with no exemption by date, that archived Specs stay out of
scope, and that the adopter notes name the sections an active Spec completes
before its next Run.

No live provider call is authorized by this record. Tests, Verification and QA
use temporary repositories and this repository's own fixtures. No other
repository is read or written by any Task.

The governed set was measured with `GovernedPath` on the authoring branch at
`7a8b7f70`. The probe was a `go test -overlay` test in `internal/speccheck`
that wrote nothing to the repository. It ran against every file the Tasks
declare, and the 14 paths above are governed.

## Why each governed path is unavoidable

- `internal/speccheck/coherence.go` registers the stage of each new code and
  calls the TechSpec-stage detector.
- `internal/docscontract/testdata/corpus-golden.json` and
  `internal/spec/archive_layout_characterization_test.go` pin the active
  corpus code set; each new code joins it at 0.
- `.agents/skills/write-techspec/SKILL.md`,
  `.agents/skills/write-techspec/references/techspec-template.md`,
  `.agents/skills/write-tasks/SKILL.md` and
  `.agents/skills/write-tasks/references/task-template.md`, with the governed
  mirrors `skills/write-techspec/SKILL.md`,
  `skills/write-techspec/references/techspec-template.md` and
  `skills/write-tasks/SKILL.md`, carry the four rules and the raised versions.
- `skills/baseline_skill_contract_test.go` gains the test that reads the four
  rules.
- `.agents/skills/roundfix/references/spec.md`, `.agents/skills/roundfix/SKILL.md`
  and `skills/roundfix/SKILL.md` document the four codes and raise the roundfix
  skill version.

## What is not governed

These paths are ordinary: the two new detector files and their tests under
`internal/speccheck`, `internal/speccheck/citations.go`, the fixtures under
`internal/speccheck/testdata/repo/docs/specs/`,
`internal/docscontract/corpus_test.go`, `internal/cli/spec_check_test.go`,
`internal/cli/deliver_plan_test.go`, `skills/write-tasks/references/task-template.md`,
`skills/roundfix/references/spec.md`, `skills/testdata/owned-skill-versions.json`,
`docs/user-guide/commands/spec.md`, `docs/references/behavior-surfaces.json`,
`CONTEXT.md` and
`docs/adr/0265-a-spec-proves-what-it-claims-and-every-active-spec-is-held-to-it.md`.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors and the Behavior
Surface Record. This declaration records the regeneration that follows the
approved edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

No Baseline module changes, so `make baseline-digests` reports no change; an
authoring prototype ran it after editing the three owned skills and it rewrote
nothing.

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml` or the CI
  workflows.
- No change to any Baseline module, guide under `docs/agents/`, Source
  Baseline asset or production Go file outside `internal/speccheck`.
- No change to the `### QA settlement` section of any skill, to another Spec's
  folder, to archived Specs, existing Archive Records, history entries or
  `CHANGELOG.md`.
- No test, Verification command or QA row opens a network connection, reads
  or writes the real `~/.roundfix`, or reads another repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
