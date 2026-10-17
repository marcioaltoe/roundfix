---
status: approved
granted: 2026-10-16
action: make the history cleanup that roundfix history sanitize and roundfix baseline update plan write a superseded ADR's missing lineage line into its successor, refuse only an ADR without an accepted successor, and give Live Documents and user guide pages the Git Locators of missing files; make the pre-PR review omit deleted history and declared generated files and review a candidate still above the bound in at most four Review Segments, and reword Delivery Convention C2 for QA Tasks; make roundfix spec check report a Task that edits an owned skill without its version record; give roundfix archive --plan an --offline form the qa-gate skill uses; make the Delivery Queue run a fresh copy of each item binary build; and describe each change in the user guides, the Roundfix, write-tasks and qa-gate skills and the glossary
consuming: 0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec
paths:
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/review.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-tasks/SKILL.md
  - docs/references/coverage-record.json
  - internal/docscontract/testdata/corpus-golden.json
  - internal/spec/archive_layout_characterization_test.go
  - internal/speccheck/coherence.go
  - internal/speccheck/constraints.go
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

# Approved authority for Spec 0276

On 2026-10-16 the maintainer approved cycle AH through AskUserQuestion: the
authorization "Governed Paths declarados" (the option text: "\"Concedo\" como
em AE/AF/AG para os Governed Paths que cada Spec declarar: módulos da
Baseline, registros de versão, make baseline-digests, skills e testes
governados. Makefile, workflows de CI, go.mod e configuração de lint ficam
fora."), "Merge e releases", and the failure rules "Regras de AE–AG
(Recomendado)"; the cycle's leftover sources went to "Spec E no fim
(Recomendado)". On 2026-10-17 the maintainer chose, for this Spec, item E:
"Escrever e entregar (Recomendado)"; for the adopters, "Não faça a edição, só
envie por inbox o que temos a trabalhar em cada um dos repos para que no
melhor momento seja atualizado. Alguns estão com runs ativas e abrir PR
enquanto isso pode ser um problema."; for the review bound, "Sim, como na 0272
(Recomendado)"; and for retries, "Sim, uma vez por item (Recomendado)". The
2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
necessário") also stands.

The governed set was measured on the authoring branch `docs/plan-0276`, based
on `e5da154a`, with `speccheck.GovernedPath`, through a `go test -overlay`
probe in `internal/speccheck` that wrote nothing to the repository (`git
status` was identical before and after, and the probe file never existed in
the tree), against every file the Tasks declare and every planning file.
Sixteen of them are governed; `docs/references/coverage-record.json` is
granted as well, for the sanctioned re-record a Task needs only if it must
rename or remove a test.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01 to
  task_04).
- `.agents/skills/roundfix/references/archive.md` covers the `history
  sanitize` and `archive` surfaces (task_01, task_04);
  `.agents/skills/roundfix/references/baseline.md` covers `baseline update`
  and `.agents/skills/roundfix/references/setup.md` names the history units
  (task_01); `.agents/skills/roundfix/references/review.md` covers `review`
  (task_02); `.agents/skills/roundfix/references/spec.md` lists the check
  codes (task_03); `.agents/skills/roundfix/references/deliver.md` covers the
  Delivery Queue's item binary (task_04).
- `internal/speccheck/constraints.go` holds `Check`, the full sweep that must
  call the new detector, and `internal/speccheck/coherence.go` registers it at
  the Tasks stage (task_03).
- `internal/docscontract/testdata/corpus-golden.json` and
  `internal/spec/archive_layout_characterization_test.go` pin the count of
  every finding code; the new code joins them (task_03).
- `.agents/skills/write-tasks/SKILL.md` and `skills/write-tasks/SKILL.md`
  carry the authoring rule for owned skill Tasks and their version fields
  (task_03).
- `.agents/skills/qa-gate/SKILL.md` and `skills/qa-gate/SKILL.md` carry the
  QA gate's archive plan instruction and their version fields (task_04).
- `docs/references/coverage-record.json` is re-recorded only by its
  sanctioned command, and only by a Task that renames or removes a test; the
  Tasks are written so that none needs to.

## What is not governed

The production and test sources the Tasks change in `internal/spec`,
`internal/cli` and `internal/speccheck` other than the two files above; the
skill reference mirrors under `skills/roundfix/references/`;
`skills/testdata/owned-skill-versions.json`;
`docs/references/behavior-surfaces.json`; `CONTEXT.md`; the user guides under
`docs/user-guide/commands/`; the ADRs `docs/adr/0278-*`, `docs/adr/0281-*` and
`docs/adr/0282-*`; and every planning file of this Spec.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill versions,
the Behavior Surface Record and, if needed, the coverage record. This
declaration records the regeneration that follows the approved edits and adds
no source path.

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

```yaml
command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
outputs:
  - docs/references/coverage-record.json
```

No Baseline module changes, so no module version is recorded; an owned skill
edit leaves the derived Baseline artifacts byte-identical
(`TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical`), so `make
baseline-digests` has nothing to rewrite. No record or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module, `.roundfixrc.yml`,
  or the `### QA settlement` section of any skill.
- No change to `SC-HISTORY-RETAINED`, the review bound's value, the Claude
  review budget, any judge call other than the archive plan's, or an existing
  finding code's behavior or count.
- No edit to any adopter repository; no test, Verification command or QA row
  reads or writes the real `~/.roundfix`, reaches the network, sends a request
  to a provider or the judge, or reads another repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
