---
status: approved
granted: 2026-11-10
action: make a commit that a review record or the disposition ledger names, and that a complete repository proves absent, restart the Reviewer Lineage of roundfix review with a reason naming it and be refused by name in review dispose and the Delivery Retry's correction proof, keep every other Git failure blocking, and document it in the review and deliver guides and the Roundfix skill's review and deliver references
consuming: 0372-a-review-whose-reviewed-head-is-gone-opens-a-new-lineage
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

# Approved authority for Spec 0372

On 2026-11-10 Conexus reported that `roundfix review` blocks every later
review of a checkout whose recorded reviewed head the repository no longer
holds
(`references/2026-11-10-conexus-review-blocks-when-the-reviewed-head-is-gone.md`,
reported by conexus). The report, ai-memory message
`01a1c53f-c3be-78a0-851d-cda2f4dd43a8` of 2026-11-10, is a measured defect
with the reporter's own suggestion, quoted verbatim:

- "Sugestão: tratar o resultado \"not a valid commit\" de `git merge-base
  --is-ancestor` (ou um `git cat-file -e` prévio do head anterior) como
  reinício de linhagem, com um motivo como \"the reviewed head <prior head>
  is no longer in the repository\"."

No further question was asked for it. The maintainer's answer this Spec
rests on (AskUserQuestion via `/roundfix --follow-up`, 2026-11-10 ~15:35
BRT) is, verbatim:

- Question "Com que autorizações eu toco sem supervisão as três Specs deste
  ciclo: 0372 (review do Conexus), 0371 (relatório de QA, que você já decidiu
  às 06:20 pelo Jarvis) e 0373 (regra de resposta)?" → "Mesmas do BA, teto
  minor (Recommended)" — option text: "Escrever, entregar pela fila, merge e
  release de cada Spec assim que pronta (patch nas duas correções; a regra
  pode sair como minor), teto de 2 Tasks corretivas, eventos ao Jarvis, sem
  notas às caixas, e Makefile/CI/go.mod/lint só onde a Spec exigir."

"Mesmas do BA" keeps the cycle BA rules of 2026-11-09 ("Mesmas do AZ
(Recommended)" — "Release patch assim que cada Spec ficar pronta, teto de 2
corretivas por Spec, merge pela fila, avisos ao Jarvis, sem notas de release
às caixas."), and through them the cycle AJ grant of 2026-10-19 "Mexer em
Governed Paths declarados" (option text: "'Concedo' para os caminhos que cada
Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and the
2026-09-30 grant of every skill ("considere autorizado a ajustar todas as
skills se necessário"). The design needs no Makefile, workflow, Go module,
lint or formatter change.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`.
No Task performs it; the operator cuts the patch right after the merge.

The governed set was measured on 2026-11-10 on the branch `docs/plan-0372`,
based on `9fda0b65`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec and its adopted
  reference.
- The four paths above are governed: the Roundfix `SKILL.md` pair whose
  version fields task_01 raises and the canonical `review` and `deliver`
  references task_01 edits. Their mirrors, the guides, the records and every
  `internal/cli` file are not governed.

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
  `docs/agents/`, `docs/references/skill-coverage.json`, `CONTEXT.md`, any
  ADR, any Archive Record, any skill other than `roundfix`, any Roundfix
  reference other than `review` and `deliver`, or the `### QA settlement`
  section of any skill.
- No change to `internal/preflight/**` (the detached-HEAD refusal stays),
  `internal/spec/qa*.go`, `internal/speccheck/**`, `internal/baseline/**` or
  `internal/cli/follow_up*.go`, which other Specs of cycle BB own.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - creates files in `/tmp`, `/private/tmp` or `$HOME` outside its own
    temporary directory;
  - starts an Agent session or an ACP adapter;
  - reaches the network, or fetches from any remote other than a local
    `file://` source it created;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
