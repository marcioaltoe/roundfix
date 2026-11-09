---
status: approved
granted: 2026-11-09
action: add the shared check that refuses Borrowed Evidence, apply it in qa-report accept, the Settle Command and the Daemon's QA settlement with one repair turn, name the Run Worktree, the audited head and the evidence directory in the QA prompt, and document it in the qa-gate skill, the qa-report, settle and implement guides, the Roundfix settle and implement references and CONTEXT.md, as ADR-0376 decides
consuming: 0366-a-qa-report-relies-only-on-evidence-its-own-run-produced
paths:
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/implement.md
  - .agents/skills/roundfix/references/settle.md
  - skills/qa-gate/SKILL.md
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

# Approved authority for Spec 0366

On 2026-11-06 quaestor reported a terminal QA gate that failed its Spec on
evidence an earlier Run of the same Spec left on the item branch, including
an executable inside that earlier Run's Run Worktree
(`references/2026-11-07-a-later-runs-qa-gate-reused-an-earlier-runs-evidence.md`).
The maintainer asked no separate question for this measured defect with the
reporter's own suggestion; the cycle BA rules answer covers it. The answers
this Spec rests on are, verbatim:

- Rules (AskUserQuestion, 2026-11-09 ~08:35–08:45 BRT), question "Que
  autorizações você concede de saída para eu levar o ciclo BA até o fim sem
  parar (implementações, merges e releases inclusos como no AZ)?", selected:
  - "Mesmas do AZ (Recommended)" — "Release patch assim que cada Spec ficar
    pronta, teto de 2 corretivas por Spec, merge pela fila, avisos ao
    Jarvis, sem notas de release às caixas."
  - "Scripts/Makefile/CI" — "Mexer em ~/.roundfix-operator, Makefile e
    .github/workflows se a migração ou a Spec B precisarem."
  - "go.mod/lint" — "Mexer em go.mod/go.sum e na configuração de
    lint/formatter se uma Task precisar."
  - "Avisar Quaestor/Fiscus" — "Mensagem ai-memory curta a cada um dizendo
    qual Spec adotou o finding (não é nota de release)."

"Mesmas do AZ" keeps the cycle AZ rules, and through them the cycle AJ grant
of 2026-10-19 "Mexer em Governed Paths declarados" (option text: "'Concedo'
para os caminhos que cada Spec declarar (módulos Baseline, registros, skills,
testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere
autorizado a ajustar todas as skills se necessário"). The design needs no
Makefile, workflow, Go module, lint or formatter change.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it right after the merge.

The governed set was measured on 2026-11-09 on the branch `docs/plan-0366`,
based on `2e59aebe`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference and
  ADR-0376.
- The six paths above are governed: the `qa-gate` and Roundfix skills'
  `SKILL.md` pairs and the Roundfix skill's canonical implement and settle
  references, which tasks 01 to 03 change.

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
  - .agents/skills/qa-gate/SKILL.md
  - skills/qa-gate/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

## What is not governed

- Every Go source and test file tasks 01 and 02 create or change under
  `internal/spec`, `internal/cli` and `internal/daemon`.
- The mirrors `skills/roundfix/references/settle.md` and
  `skills/roundfix/references/implement.md`, and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/qa-report.md`,
  `docs/user-guide/commands/settle.md`,
  `docs/user-guide/commands/implement.md`, `CONTEXT.md`,
  `docs/references/behavior-surfaces.json` and ADR-0376.
- Every planning file of this Spec and its adopted reference.
