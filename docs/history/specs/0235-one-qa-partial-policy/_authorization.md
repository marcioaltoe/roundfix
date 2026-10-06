---
status: approved
granted: 2026-10-06
action: make Daemon settlement, roundfix settle, roundfix archive and roundfix qa-report accept apply one QA partial policy that exempts the pre-PR Pull Request row and network-denied outside-evidence rows, make the Delivery Queue park only partials that need the override, and describe the policy in the qa-gate, archive-spec and Roundfix skills and in the Spec routing and autonomous-work Baseline guides
consuming: 0235-one-qa-partial-policy
paths:
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/settle.md
  - docs/agents/autonomous-work.md
  - docs/agents/setup-context.json
  - docs/agents/spec-routing.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - internal/baseline/assets/modules/autonomous-work.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - skills/archive-spec/SKILL.md
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0235

On 2026-09-30 the maintainer asked for unattended work through the program
and said of the skills "considere autorizado a ajustar todas as skills se
necessário", and of the Baseline source and the guides "Autorizar os dois".
On 2026-10-06 the maintainer answered this cycle's questions through
AskUserQuestion. The scope was "Y1, Y2 e Y3", and this Spec is Y2, delivered
first. For rows a Run sandbox cannot reach, the answer was "Igual à linha do
PR": an outside-evidence row blocked only because the sandbox denied network
access never decides a qualifying partial, the report records that it was
not reached, and whoever needs that proof declares Unreachable Acceptance.
For the Governed Paths this Spec declares, the answer was "Concedo".

The governed set was measured with `GovernedPath` on the authoring branch at
`bde616d1`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare. The derived Baseline files
were measured by applying task_04's edit in a disposable clone and running
`make baseline-digests` and the Managed Refresh there.

## Why each governed path is unavoidable

- `.agents/skills/qa-gate/SKILL.md` and its mirror — the gate must write the
  network-denied marker exactly, and its outside-evidence and sandbox
  paragraphs state the old rule.
- `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`
  and their mirrors — the `### QA settlement` table, identical in the three
  skills by contract, states which unmet rows a qualifying partial may hold.
  Each owned skill's version is raised with its content.
- `.agents/skills/roundfix/references/archive.md` and `settle.md` — they state
  the policy for the Archive Command and settle, and the repository's
  skill-sync rule ships them with the behavior change.
- `internal/baseline/assets/modules/spec-workflow.json` and
  `autonomous-work.json` — their clauses say that a blocked outside-evidence
  row holds Pull Request preparation and that only the Pull Request row is
  exempt.
- `docs/agents/spec-routing.md`, `docs/agents/autonomous-work.md`,
  `docs/agents/setup-context.json`, the two formatter goldens and
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json` — the
  deterministic fallout of those module edits, written by `make
  baseline-digests` and the Managed Refresh, never by hand.

## What is not governed

The Go sources and tests under `internal/spec`, `internal/cli`,
`internal/daemon` and `internal/baseline` that the Tasks declare, the
Baseline test data under `internal/baseline/testdata/`, the skill reference
mirrors `skills/roundfix/references/archive.md` and
`skills/roundfix/references/settle.md`,
`skills/testdata/owned-skill-versions.json`, the user guides, and the ADRs
are ordinary.

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

`make baseline-digests` regenerates the catalog snapshots, the plan goldens,
the formatter goldens and the profile digest, and
`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
renders this repository's guides and Setup Manifest; a second refresh reports
no file change. No digest pin, golden or generated guide is hand-edited.

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration or the CI workflows.
- No change to the `pass` rule, the mechanical stage's count checks or the
  frontmatter count fields, and no exemption for any other environment row.
- No change to archived Specs, existing QA Reports, `CONTEXT.md` or
  `CHANGELOG.md`.
- No test, Verification command or QA row reaches the network or reads or
  writes the real `~/.roundfix`.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
