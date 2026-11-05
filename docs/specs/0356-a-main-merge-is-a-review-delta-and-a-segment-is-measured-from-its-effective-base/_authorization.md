---
status: approved
granted: 2026-11-05
action: measure every pre-PR review diff inside the candidate from its Effective Base, keep the Reviewer Lineage across a conflict-free merge of the default branch and name each restart, print that restart from deliver status and deliver retry, and document the rules in the review and deliver guides, the Roundfix review and deliver references, CONTEXT.md, ADR-0366 and the amendment notices of ADR-0197, ADR-0282, ADR-0286 and ADR-0352
consuming: 0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/review.md
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

# Approved authority for Spec 0356

On 2026-11-05 at about 08:20 BRT the maintainer answered the operator's cycle
AZ batch 2 questions with AskUserQuestion, after Fiscus reported a Delivery
Queue of about 21 Specs that had run for two days
(`references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-main-merge-restarts-the-review.md`).
The answers this Spec rests on are, verbatim:

- Question: "Relato do Fiscus (fila de 21 Specs há 2 dias). Quais pontos
  viram Spec agora? Os que ficarem de fora viram Backlog Entry." — all four
  options selected; the two this Spec answers:
  - "2. Merge do main (Recommended)" — option text: "Um merge do main sem
    conflito, em que o único pai novo é o main, conta como delta e não
    reinicia o review na rodada 1. Quando o review for reiniciar, o retry
    avisa antes."
  - "4. Trecho do review (Recommended)" — option text: "Cada trecho do review
    passa a ser medido contra o merge-base com o main atual. Um commit antigo
    que só parece reverter o main deixa de bloquear o review por tamanho."
- Rules (all): "Mesmas do AZ (Recommended)" — option text: "Sem supervisão,
  com merge pela fila e conserto quando a entrega parar. Makefile, CI, go.mod
  e lint liberados, assim como mudanças de Baseline e skills que chegam aos
  adotantes. Release assim que cada Spec ficar pronta, com notas aos
  adotantes, aviso ao Fluxus e ao Fiscus."
- Standing instruction of 2026-11-05: "Quero que faça os releases hoje assim
  que estiverem prontos. E se existe no baseline regra para um release por dia
  por algum motivo deve ser removida." The repository has no such rule; this
  Spec changes nothing for it.

"Mesmas do AZ" keeps the cycle AZ rules of 2026-11-05, and through them the
cycle AY, AX, AW, AV, AU, AT and AS answers, the cycle AJ grant of 2026-10-19
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os
caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere autorizado
a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The answers ask for a release as soon as this Spec is ready; it is
a `type: fix` Spec, so the Release Ceiling is `patch`, and the release may
change only the three version-only files in `release_paths`: the changelog,
the npm package manifest and the application version. No Task performs it;
the operator cuts it after the merge.

The governed set was measured on 2026-11-05 on the branch `docs/plan-0356`,
based on `ad77a797`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference and
  ADR-0366.
- The four paths above are governed: the Roundfix skill's `SKILL.md` pair and
  its canonical `deliver` and `review` references.

## What is not governed

- `internal/cli/review_effective_base.go`, `internal/cli/review_segments.go`,
  `internal/cli/review_lineage.go`, `internal/cli/review.go`,
  `internal/cli/review_lineage_notice.go`, `internal/cli/deliver.go` and the
  three new test files under `internal/cli`.
- The mirrors `skills/roundfix/references/deliver.md` and
  `skills/roundfix/references/review.md`, and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/review.md`, `docs/user-guide/commands/deliver.md`,
  `docs/references/behavior-surfaces.json`, `CONTEXT.md`, ADR-0197, ADR-0282,
  ADR-0286, ADR-0352 and ADR-0366.
- Every planning file of this Spec and its adopted reference.

## Sanctioned regeneration

The repository-owned commands write the skill mirrors and the generated
records. This declaration records the regeneration that follows the approved
edits and adds no source path.

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

No mirror, version field or record is hand-edited. No top-level test is
removed or renamed, so the coverage record is not re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module, profile, template or derived Baseline file, any guide under
  `docs/agents/`, `docs/references/skill-coverage.json`,
  `docs/references/coverage-record.json`, any Archive Record, any skill other
  than `roundfix`, any Roundfix reference other than `deliver` and `review`,
  or the `### QA settlement` section of any skill.
- No change to `internal/delivery`, `internal/store`, `internal/config` or
  any ADR other than ADR-0197, ADR-0282, ADR-0286 and ADR-0352.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
