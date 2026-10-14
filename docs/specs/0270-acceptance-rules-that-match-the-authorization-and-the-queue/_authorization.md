---
status: approved
granted: 2026-10-13
action: accept an outside-evidence row the Spec's own authorization keeps offline toward a qualifying partial in the shared QA policy, report a QA requirement that pins a byte invariant to a commit as a strict gap, and state both rules in the Baseline guides, the owned skills and the user guide with the recorded version raises and regenerated derived files
consuming: 0270-acceptance-rules-that-match-the-authorization-and-the-queue
paths:
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/settle.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-tasks/SKILL.md
  - docs/agents/autonomous-work.md
  - docs/agents/setup-context.json
  - docs/agents/spec-routing.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - internal/baseline/assets/modules/autonomous-work.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/docscontract/testdata/corpus-golden.json
  - internal/spec/archive_layout_characterization_test.go
  - internal/speccheck/coherence.go
  - skills/archive-spec/SKILL.md
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-tasks/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0270

On 2026-09-30 the maintainer asked for unattended work through every release
of the program. Of the skills the maintainer said "considere autorizado a
ajustar todas as skills se necessário". Of Baseline sources, guides and
`.roundfixrc.yml` the maintainer said "Autorizar os dois". On 2026-10-11 the
maintainer answered "Grants dos caminhos governados" with "Concedo" for the
Governed Paths each Spec declares. On 2026-10-13, through a structured
question of about 17:20 BRT (20:20 UTC), the maintainer approved cycle AG with
six Specs in the order B, A, C, D, E, F, and answered "Pode iniciar" at about
17:25 BRT (20:25 UTC) to start them under the operating rules of cycles AE and
AF: the same operations and the same "Concedo" for the Governed Paths each
Spec declares. This Spec is cycle AG item E. The maintainer decided its
policy: an outside-evidence row blocked because the Spec's own authorization
forbids network access counts toward a qualifying partial without the QA
Archive Override, in one canonical wording that the QA skill writes, and
settlement and archive share the policy.

No live provider call is authorized by this record. Tests, Verification and QA
use temporary repositories, temporary homes and this repository's own
fixtures. No test, Verification command or QA row opens a network
connection, and no other repository is read or written by any Task.

The governed set was measured on the authoring branch at `9963ffd3` with
`speccheck.GovernedPath`, through a `go test -overlay` probe in
`internal/speccheck` that wrote nothing to the repository. It ran against
every file the Tasks declare, and the 22 paths above are governed.

## Why each governed path is unavoidable

- `internal/speccheck/coherence.go` registers the new code's stage.
  `internal/docscontract/testdata/corpus-golden.json` and
  `internal/spec/archive_layout_characterization_test.go` pin the active
  corpus code set, which the new code joins at 0.
- `internal/baseline/assets/modules/spec-workflow.json` and
  `internal/baseline/assets/modules/autonomous-work.json` carry the two
  Baseline clauses that name the exempt outside-evidence rows. Each gains one
  sentence and the next module version.
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, the
  two formatter goldens, `docs/agents/spec-routing.md`,
  `docs/agents/autonomous-work.md` and `docs/agents/setup-context.json` are
  derived from those modules and change only through `make baseline-digests`
  and the Managed Refresh.
- `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/archive-spec/SKILL.md`
  and `.agents/skills/roundfix/SKILL.md` carry the QA settlement table this
  Spec extends under its own heading, and their version fields. The roundfix
  references `archive.md`, `settle.md` and `spec.md` describe the archive,
  settlement and spec check commands whose policy and codes change.
  `.agents/skills/write-tasks/SKILL.md` gains the byte-invariant rule. The
  four governed mirrors under `skills/` carry the same bytes after
  `make skills-sync`.

## What is not governed

These paths are ordinary: `internal/spec/qa.go`,
`internal/spec/qa_partial_policy_test.go`,
`internal/cli/qa_partial_policy_test.go`,
`internal/daemon/qa_partial_policy_test.go`,
`internal/speccheck/pinned_bytes.go`, `internal/speccheck/pinned_bytes_test.go`,
`internal/speccheck/citations.go`, `internal/docscontract/corpus_test.go`,
`internal/baseline/module-versions.json`, the catalog and plan
characterization goldens under `internal/baseline/testdata/`,
`internal/baseline/network_prohibited_row_clause_test.go`,
`skills/roundfix/references/archive.md`, `skills/roundfix/references/settle.md`,
`skills/roundfix/references/spec.md`,
`skills/testdata/owned-skill-versions.json`,
`skills/acceptance_rules_contract_test.go`, the user guide pages,
`docs/references/behavior-surfaces.json`, `CONTEXT.md` and
`docs/adr/0276-a-row-the-authorization-keeps-offline-qualifies-and-a-byte-invariant-compares-with-the-delivery-base.md`.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files, the skill
mirrors and the generated records. This declaration records the regeneration
that follows the approved edits and adds no source path.

```yaml
command: make baseline-digests
```

The Module Version Record step writes each changed module's version line and
the record; suiteguard accepts those writes only for this declared command:

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/autonomous-work.json
  - internal/baseline/assets/modules/spec-workflow.json
```

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
```

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
renders this repository's guides and Setup Manifest, and a second refresh
reports no file change. No digest pin, golden, generated guide or record is
hand-edited. An authoring prototype at `9963ffd3` ran these commands after the
approved edits and changed exactly the paths the Tasks declare.

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml` or the CI
  workflows.
- No change to any Baseline clause other than
  `clause.spec.project-constraints-06-outside-evidence` and
  `clause.autonomous.loop-01-qa-once`, no clause removed or renamed, and no
  Source Baseline corpus, retention or snapshot change.
- No change to the `### QA settlement` section of any skill, to another
  Spec's folder, to archived Specs, existing Archive Records, history entries
  or `CHANGELOG.md`.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
