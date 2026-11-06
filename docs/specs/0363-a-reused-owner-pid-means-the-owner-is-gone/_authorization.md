---
status: approved
granted: 2026-11-06
action: add the store question OwnerGone, make the Doctor residue check read a proven start-identity mismatch at a terminal Run's owner PID as an owner that is gone, make the Active Run lock reclaim, the Delivery Queue Run slots and deliver start use the same rule, and document it in the doctor and deliver guides and the Roundfix skill's setup and deliver references, as ADR-0373 decides
consuming: 0363-a-reused-owner-pid-means-the-owner-is-gone
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/setup.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
  - release
release_ceiling: patch
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0363

On 2026-11-06 argus reported that `roundfix doctor` keeps `residue: partial`
forever once a dead owner's PID is reused
(`references/2026-11-06-doctor-residue-stays-partial-when-an-owner-pid-is-reused.md`).
The report arrived after the operator's cycle AZ batch 4 questions; the
maintainer asked no separate question for this clear defect, and the batch's
rules answer covers it as a new Spec of the batch. The answers this Spec
rests on are, verbatim:

- Rules (AskUserQuestion, 2026-11-06 ~12:00 BRT): "Mesmas do AZ
  (Recommended)" — option text: "Sem supervisão, merge pela fila, conserto e
  retry quando parar. Release assim que cada Spec ficar pronta; resposta ao
  OnionCry, ao sql-library e ao Jarvis."
- The maintainer's standing instruction of 2026-11-05: "Quero que faça os
  releases hoje assim que estiverem prontos. E se existe no baseline regra
  para um release por dia por algum motivo deve ser removida."

"Mesmas do AZ" keeps the cycle AZ rules of 2026-11-05, and through them the
cycle AY, AX, AW, AV, AU, AT and AS answers, the cycle AJ grant of 2026-10-19
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os
caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere autorizado
a ajustar todas as skills se necessário"). The design needs no Makefile,
workflow, Go module, lint or formatter change.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it right after the merge.

The governed set was measured on 2026-11-06 on the branch `docs/plan-0363`,
based on `9a285d77`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference and
  ADR-0373.
- The four paths above are governed: the Roundfix skill's `SKILL.md` pair
  and its canonical setup and deliver references, which task_01 changes.

## Sanctioned regeneration

The repository-owned commands write the skill mirrors and the generated
records. No mirror, version field, digest or record is hand-edited.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

## What is not governed

- Every Go source and test file tasks 02 and 03 create or change under
  `internal/store` and `internal/cli`.
- The mirrors `skills/roundfix/references/setup.md` and
  `skills/roundfix/references/deliver.md`, and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/doctor.md`, `docs/user-guide/commands/deliver.md`,
  `docs/references/behavior-surfaces.json` and ADR-0373.
- Every planning file of this Spec and its adopted reference.
