---
status: approved
granted: 2026-09-30
action: prove Claim Receipts and report attributions that lack one, declare and trace Surface Transcripts, start both rules at the commit that added the concrete-contract guide, and teach both forms in the write-prd, write-techspec, write-tasks, qa-gate and Roundfix skills
consuming: 0191-claims-with-receipts-and-contracts-as-they-ship
paths:
  - internal/speccheck/coherence.go
  - internal/docscontract/testdata/corpus-golden.json
  - internal/spec/archive_layout_characterization_test.go
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/write-techspec/references/techspec-template.md
  - .agents/skills/write-techspec/references/concrete-contracts.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-prd/references/prd-template.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - skills/write-techspec/SKILL.md
  - skills/write-techspec/references/techspec-template.md
  - skills/write-prd/SKILL.md
  - skills/write-prd/references/prd-template.md
  - skills/write-tasks/SKILL.md
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0191

On 2026-09-30 the maintainer approved the program "Três ondas", whose Wave 7
includes this Spec (blocks A2 and A3: concrete contracts and proved
receipts). The same day the maintainer wrote "considere autorizado a ajustar
todas as skills se necessário", which authorizes edits to every
Roundfix-owned skill, and answered "Autorizar os dois" to a structured
question about the Baseline source and guides and `.roundfixrc.yml`. This Spec
edits no Baseline asset and does not touch `.roundfixrc.yml`. The three
governed Go paths ride the standing grant of 2026-09-21 for governed sources a
slice genuinely needs. The set was measured with `GovernedPath` on the
authoring branch at `9e439dbb`, through a `go test -overlay` probe that wrote
nothing to the repository.

## Why each governed path is unavoidable

- `internal/speccheck/coherence.go` — the stage table there lists every
  diagnostic code with its authoring stage, and the stage-scoped check calls
  each detector. The five new codes and the two detectors are registered and
  called there (task_01, task_02, task_03).
- `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go` — the corpus golden
  pins the active count of every code, and the characterization test pins the
  golden's bytes. The five new codes join both at `0` (task_01).
- `.agents/skills/write-techspec/SKILL.md`,
  `.agents/skills/write-techspec/references/techspec-template.md`,
  `.agents/skills/write-techspec/references/concrete-contracts.md` and the
  mirrors `skills/write-techspec/SKILL.md` and
  `skills/write-techspec/references/techspec-template.md` — the skill ships the
  guide whose commit starts the contract, and its template gains the Surface
  Transcripts section (task_04).
- `.agents/skills/write-prd/SKILL.md`,
  `.agents/skills/write-prd/references/prd-template.md` and their mirrors — the
  skill and template ask for a receipt on every attribution (task_04).
- `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/qa-gate/SKILL.md` and
  their mirrors — Tasks name each transcript, and the gate reproduces it
  (task_04).
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the Spec
  Consistency Check section lists the stable identifiers (task_02, task_03).

`.agents/skills/` is canonical and `skills/` is its mirror.

## What is not governed

`internal/speccheck/citations.go`, the new files
`internal/speccheck/authoring_horizon.go`, `internal/speccheck/receipts.go`
and `internal/speccheck/transcripts.go`, their new tests, the new
characterization test and its fixtures under
`internal/speccheck/testdata/receipt-characterization/`, the existing tests
`internal/speccheck/coherence_test.go` and
`internal/docscontract/corpus_test.go`, the new test
`skills/concrete_contract_guide_test.go`, the mirror
`skills/write-techspec/references/concrete-contracts.md`, `CONTEXT.md` and
`docs/user-guide/context-driven-development.md` are ordinary.
`internal/speccheck/constraints.go`,
`internal/speccheck/constraints_characterization_test.go`,
`internal/speccheck/mechanical_test.go`, `internal/cli/cli_test.go` and
`docs/references/coverage-record.json` are governed and are not touched: no
help text changes and no existing top-level test is renamed or removed.

## Sanctioned regeneration

The repository-owned commands resolve their generated outputs. These
declarations record the regeneration that follows the approved skill edits and
add no source paths.

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

## Limits

- No change to the Makefile, a lint, formatter or test-runner configuration, a
  CI workflow, `go.mod`, a Baseline asset or `.roundfixrc.yml`.
- No change to archived Specs, to another active Spec, or to the
  `roundfix-speccheck/v1` schema.
- No text is added inside the `### QA settlement` section of any skill.
- No model or network call is added to the Spec Consistency Check.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
