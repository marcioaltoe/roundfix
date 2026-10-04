---
status: approved
granted: 2026-10-04
action: make branch.prefix an optional Baseline decision whose unrecorded state renders the commit-type rule, state the re-execution rule for a QA row without inputs in the qa-gate skill, let roundfix release plan report the skills and baseline checks read-only, and describe this in the skills and guides
consuming: 0223-adjustments-the-adopters-asked-for
paths:
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/release.md
  - docs/agents/setup-context.json
  - internal/baseline/assets/decisions.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/templates/guides/agent-instructions.md
  - internal/cli/baseline_plan_test.go
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0223

On 2026-10-04 the maintainer, asked through a structured question, decided to
deliver this cycle's item L, the adjustments the adopters asked for, and
answered "Autorizar os dois" for the Baseline source and the agent guides. On
2026-09-30 the maintainer had said of the skills "considere autorizado a
ajustar todas as skills se necessário". The same 2026-10-04 answer granted
"Concedo" for `internal/spec/archive.go`; that path belongs to Spec 0226,
which owns the archive's relative links, and this record does not use it.

The governed set was measured with `GovernedPath` on the authoring branch at
`6ea9e1e9`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/qa-gate/SKILL.md` and its mirror `skills/qa-gate/SKILL.md`
  — the re-execution rule for a row without inputs lives in the skill, whose
  content changes only with its version (task_01).
- `.agents/skills/roundfix/SKILL.md`, its mirror `skills/roundfix/SKILL.md`
  and `.agents/skills/roundfix/references/release.md` — the release plan's two
  new lines, and the version raise the skill-sync rule requires for a change
  of CLI behavior (task_01).
- `internal/baseline/assets/decisions.json`,
  `internal/baseline/assets/modules/core.json` and
  `internal/baseline/assets/templates/guides/agent-instructions.md` — the
  decision's optional flag, the core module's required list, and the template
  whose prefix sentence moves into the decision's renderer (task_02).
- `docs/agents/setup-context.json` — this repository's Setup Manifest; the
  public Managed Refresh rewrites its catalog digest after the catalog changes
  (task_02).
- `internal/cli/baseline_plan_test.go` — the characterization case
  `decisions-absent-names-every-required-decision` asserts that every decision
  of the Go CLI/TUI profile is reported missing. An optional decision is no
  longer missing by design, so the case changes with the contract it
  characterizes (task_02). Spec 0207 bounded the governed Baseline tests that
  its optional decision invalidated under the same "Autorizar os dois" answer.
  On 2026-10-04 the maintainer, asked through a structured question, answered
  "Confirmo" that this test is bounded under "Autorizar os dois" following
  Spec 0207's precedent.

## What is not governed

`internal/baseline/project_decision_render.go`, the catalog snapshots and the
four plan goldens under `internal/baseline/testdata/`, the two new Baseline
test files, `internal/cli/doctor.go`, `internal/cli/releaseplan_command.go`,
`internal/cli/cli.go`, the new `internal/cli/releaseplan_checks.go` and its
test file, `skills/roundfix/references/release.md`,
`skills/testdata/owned-skill-versions.json`,
`docs/user-guide/release-runbook.md`,
`docs/user-guide/context-driven-development.md`, ADR-0228 and this Spec's own
files are ordinary.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the derived Baseline
artifacts and this repository's Setup Manifest. This declaration records the
regeneration that follows the approved edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

`make baseline-digests` regenerates the catalog snapshots and plan goldens,
and `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
renders this repository's Setup Manifest; a second refresh reports no file
change. No digest pin, golden or generated guide is hand-edited.

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, the CI workflows or `.roundfixrc.yml`.
- No change to archived Specs, existing QA Reports, `CONTEXT.md`, the
  `### QA settlement` section of any skill, or the human first-adoption
  questions.
- No test, Verification command or QA row reaches GitHub, a provider or the
  network, or writes under the real `~/.roundfix`.
- No release, tag or deployment.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
