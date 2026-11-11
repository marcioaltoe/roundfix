---
status: approved
granted: 2026-11-11
action: measure a pre-PR review across a conflicting merge of the default branch from its Conflicted Carry, keep the Reviewer Lineage across that merge, add the Conflicted Carry lines to the round-2, Correction Review and segment prompts and the absence-claim sentence to every review prompt, and document it in the review and deliver guides, the Roundfix skill's review and deliver references, the glossary, ADR-0391 and the amendment notices of ADR-0366 and ADR-0197
consuming: 0381-a-conflicting-merge-of-the-default-branch-is-a-review-delta
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

# Approved authority for Spec 0381

On 2026-11-11 Conexus reported two false compile findings from the pre-PR
review right after a merge of `main` into an item branch
(`references/2026-11-11-conexus-review-false-compile-findings-after-a-main-merge.md`,
reported by conexus). The maintainer answered the operator's cycle BC
questions with AskUserQuestion (via `/roundfix --follow-up`, 2026-11-11
~14:45–14:50 BRT). The answer this Spec rests on is, verbatim:

- Question "Specs novas D (mensagens: pergunta não bloqueante recebe
  resposta; precedência de 'bloqueado' limitada) e E (review pré-PR revisa só
  o delta e com contexto depois de merge da main): como autorizar?" → "Após a
  0379 (Recommended)" — option text: "Mesmas autorizações do ciclo BB (push,
  merge, release, teto minor); entregues depois da fila 0377→0379, D antes de
  E."

This Spec is E. "Mesmas autorizações do ciclo BB" keeps the cycle BB rules
of 2026-11-10 ("Mesmas do BA, teto minor (Recommended)" and "Mesmas do BB,
teto minor (Recommended)"), and through them the cycle BA, AZ, AY, AX, AW, AV,
AU, AT and AS answers, which release to the Runs the owned skills and changes
of skills that reach the adopters, the cycle AJ grant of 2026-10-19 "Mexer em
Governed Paths declarados" (option text: "'Concedo' para os caminhos que cada
Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and the
2026-09-30 grant of every skill ("considere autorizado a ajustar todas as
skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable in this cycle only where a design
needs them. This Spec needs none of them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it after the merge.

The governed set was measured on 2026-11-11 on the branch `docs/plan-0381`,
based on `013af3b4`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  folder and ADR-0391 afterwards, and the probe file never existed in the
  tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every file the scratch prototype of tasks 01 and 02
  changed or created, every planning file of this Spec, its adopted
  reference and ADR-0391.
- The four paths above are governed: the Roundfix skill's `SKILL.md` pair
  and its canonical `review` and `deliver` references, which task_01
  changes.

## What is not governed

- Every Go source and test file task_02 creates or changes under
  `internal/cli/`.
- The mirrors under `skills/roundfix/references/` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/**`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`,
  ADR-0366, ADR-0197 and ADR-0391.
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

No mirror, version field, digest or record is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`,
  `skills-lock.json`, any Vendored Skill under `.agents/skills`, any Baseline
  module, template, profile or Setup Snapshot, any guide under `docs/agents/`,
  `docs/references/coverage-record.json`, any Archive Record, any skill other
  than `roundfix`, any Roundfix reference other than `review` and `deliver`,
  `internal/preflight/**`, `internal/delivery/**`, or the
  `### QA settlement` section of any skill.
- No test, Verification command or QA row:
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task sends an ai-memory message or a Spec Event Announcement; the
  operator does.
