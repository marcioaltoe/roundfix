---
status: approved
granted: 2026-11-01
action: make every QA settlement path, the Daemon's QA settlement, roundfix qa-report accept and the Settle Command for a qa Task, refuse a QA Report whose local link does not resolve in the committed tree, naming the link, its line and the reason ignored or missing, give the Daemon's QA gate one repair turn before that refusal, share the link extraction with SC-HISTORY-RETAINED, and teach the qa-gate skill to keep evidence under names the repository does not ignore
consuming: 0340-qa-evidence-links-resolve-in-the-committed-tree
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
---

# Approved authority for Spec 0340

On 2026-11-01, around 22:20 BRT, the maintainer answered the operator's
cycle AU questions with AskUserQuestion, in the `/roundfix --unattended`
mode after the cycle AT follow-up. The answers this Spec rests on are,
verbatim:

- Scope: "3 Specs (Recommended)" — option text: "0338 release: tipo do PR
  segue a Spec + teto de release no _authorization.md (os dois mexem no
  release plan e no write-*); 0339 apply aprovado apaga as aposentadas (nova
  ADR, context7-cli/exa-web-search ficam); 0340 links de evidência da QA
  existem no commit. Fila nessa ordem."
- QA evidence: "Recusar no settle (Recommended)" — option text: "Antes do
  settle da QA, cada link local de evidência precisa existir na árvore que
  vai ser commitada; link ignorado ou ausente falha com o caminho e o motivo,
  e o Agente tem um turno de reparo. A skill qa-gate manda salvar evidência
  com extensão que o repositório não ignora (.txt)."
- Release of the cycle: "Uma no fim, v0.129.0 (Recommended)" — option text:
  "Depois das três Specs, com aviso aos adotantes; o apply que apaga
  aposentadas e o teto de release chegam juntos pela Baseline."
- Authorization: "Mesmas do AT (Recommended)" — option text: "Push e PR dos
  planos, entrega e merge pela fila, Makefile/workflows/go.mod/lint quando o
  desenho pedir, 'Consertar e seguir' em falha técnica, mudança da Baseline
  (agent-instructions, release, spec-routing, skills de autoria) que chega
  aos adotantes, release v0.129.0 no fim e aviso aos adotantes."
- Start: "Já (Recommended)".

Those answers keep the cycle AT rules of 2026-11-01 (Authorization "Mesmas
do AS (Recommended)"), and through them the cycle AS answers: of 2026-10-31,
"Consertar e seguir (Recommended)" (option text: "Diagnostico, emendo o
plano ou crio uma Spec bloqueadora e entrego, sem teto fixo de corretivas.
Só pulo a Spec se o conserto exigir uma decisão de produto ou arquitetura
que contradiga o PRD aprovado; aí ela sai do ciclo e a release vai com as
demais, ou espera se for a 0333.") and the Tooling answer "Makefile,
.github/workflows, go.mod / go.sum, Config de lint/formatter"; of
2026-10-30, "Mesmas regras do AR (Recommended)" (option text: "Push de
branches, PR dos planos, entrega e merge pela Delivery Queue, eu escolho o
ponto de release e aviso os adotantes por ai-memory."). They keep the cycle
AQ answers of 2026-10-26 ("Sim, pode fazer a higiene." and "Vamos seguir a
sua sugestão de sequencia para a spec e os demais itens. Talvez seja
importante usar o jev para validar."), the cycle AJ grant of 2026-10-19
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os
caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere
autorizado a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the release is the operator's, v0.129.0 after
the three Specs of the cycle, as the Release answer states.

The governed set was measured on 2026-11-01 on the branch `docs/plan-0340`,
based on `dd3b2b62`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain internal` showed no
  change from it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference and
  ADR-0351.
- The six paths above are governed: the canonical `qa-gate` skill, the
  canonical Roundfix `implement` and `settle` references, and the `SKILL.md`
  pairs of both skills whose version fields rise. The canonical
  `domain-modeling` skill is governed too, but task_01 only reads it as an
  instruction, so it is not listed.

## What is not governed

- `internal/spec/record_links.go`, `internal/spec/qa_evidence_links.go`,
  `internal/spec/qa_evidence_links_test.go`, `internal/cli/qa_report.go`,
  `internal/cli/qa_report_test.go`, `internal/cli/settle.go`,
  `internal/daemon/task_engine.go`, `internal/daemon/qa_evidence_links_test.go`.
- The mirrors `skills/roundfix/references/implement.md` and
  `skills/roundfix/references/settle.md`, and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/qa-report.md`,
  `docs/user-guide/commands/settle.md`,
  `docs/user-guide/commands/implement.md`,
  `docs/references/behavior-surfaces.json`, `CONTEXT.md`,
  `docs/adr/0351-a-qa-report-settles-only-when-its-local-links-resolve-in-the-committed-tree.md`.
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
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - skills/qa-gate/SKILL.md
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
  or the `### QA settlement` section of any skill.
- No change to `spec.QAReportEligibility`, `roundfix archive`, the
  `qa-report accept` usage text, the Secret Scan or the non-image binary
  check.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
