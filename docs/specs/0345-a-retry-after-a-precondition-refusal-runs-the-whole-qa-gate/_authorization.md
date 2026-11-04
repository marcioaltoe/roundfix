---
status: approved
granted: 2026-11-04
action: make the QA Report Import refuse a candidate QA Report commit that adds or modifies a report recording a Precondition Refusal, so a retry after a precondition refusal runs the whole QA gate, pin with a test that the Daemon-written refusal report goes through the configured Format Command before its commit, and document the import in the Roundfix skill's implement reference, the configuration guide, the glossary and a new ADR
consuming: 0345-a-retry-after-a-precondition-refusal-runs-the-whole-qa-gate
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/implement.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
  - release
release_ceiling: minor
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0345

On 2026-11-04 the maintainer answered the operator's cycle AX questions with
AskUserQuestion, in the `/roundfix --follow-up` mode, after Pantheon reported
that `roundfix deliver retry` re-imported a QA Report of a precondition
refusal and never ran the QA gate again. The answers this Spec rests on are,
verbatim:

- 0345: "Rodar QA de novo + format (Recommended)" — option text: "Retry após
  recusa de precondição descarta o report de recusa e roda o QA gate inteiro
  (precondição, depois agente); o report de precondição passa pelo
  verification.format antes do commit; implement.md documenta a importação
  de QA Report no retry."
- Rules (all four): "Mesmas do AW (Recommended)" — option text: "Sem
  supervisão: push e PR dos planos, entrega e merge pela fila, consertar e
  seguir em falha técnica; Makefile, .github/workflows, go.mod/go.sum e lint
  liberados quando o desenho pedir; mudança de Baseline que chega aos
  adotantes liberada."
- Release (all four): "Uma v0.131.0 no fim (Recommended)" — option text:
  "Uma minor (a mudança da Baseline e o runs list fora de repo são feat);
  notas aos 8 adotantes, resposta ao Pantheon com o defeito corrigido e o
  contorno que podem largar."

Those answers keep the cycle AW rules of 2026-11-03 (Rules "Mesmas do AV
(Recommended)", option text: "Push e PR dos planos, entrega e merge pela
Delivery Queue, consertar e seguir em falha técnica (pular só se o conserto
contradisser o PRD). Makefile, .github/workflows, go.mod/go.sum e config de
lint/formatter liberados quando o desenho pedir."), and through them the cycle
AV, AU, AT and AS answers, the cycle AJ grant of 2026-10-19 "Mexer em Governed
Paths declarados" (option text: "'Concedo' para os caminhos que cada Spec
declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and the
2026-09-30 grant of every skill ("considere autorizado a ajustar todas as
skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Release answer asks for one minor, v0.131.0, after Specs 0345,
0346, 0347 and 0348 merge, so the Release Ceiling is `minor` and the release
may change only the three version-only files in `release_paths`: the
changelog, the npm package manifest and the application version. No Task
performs the release; the operator cuts it.

The governed set was measured on 2026-11-04 on the branch `docs/plan-0345`,
based on `6a0ed904`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec and its adopted reference.
- The three paths above are governed: the Roundfix skill's `SKILL.md` pair
  and its canonical `implement` reference. The canonical `domain-modeling`
  skill is governed too, but task_02 only reads it as an instruction, so it
  is not listed.

## What is not governed

- `internal/daemon/qa_prior_pass.go`, `internal/daemon/qa_prior_pass_test.go`,
  `internal/daemon/qa_format_test.go`.
- The mirror `skills/roundfix/references/implement.md` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/configuration.md`, `docs/references/behavior-surfaces.json`,
  `CONTEXT.md` and ADR-0355.
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
  `docs/agents/`, `docs/references/skill-coverage.json`, any Archive Record,
  any skill other than `roundfix`, any Roundfix reference other than
  `implement`, or the `### QA settlement` section of any skill.
- No change to `internal/daemon/task_engine.go`,
  `internal/daemon/qa_format.go`, `internal/spec`, `internal/speccheck`,
  `internal/cli` or the carry-forward code.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session or an ACP adapter;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
