---
status: approved
granted: 2026-11-16
action: add the Spec Queue Table to the status view of the Follow-up Command, with a Queue Forecast from the repository's delivered Specs and a named condition where a wait cannot be dated, list the Specs being authored from the Phase Record and the open Backlog Entries, and describe the table in the Roundfix skill, the user guide and the glossary
consuming: 0398-the-status-view-shows-the-spec-queue-with-forecasts
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/follow-up.md
  - .agents/skills/roundfix/references/modes.md
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

# Approved authority for Spec 0398

The Spec adopts one Backlog Entry, recorded by the operator from the maintainer's request of 2026-11-15, and the maintainer approved its Idea on 2026-11-16. The maintainer's answers are recorded whole in [`references/2026-11-16-maintainer-answers.md`](references/2026-11-16-maintainer-answers.md). The answers this record rests on are, verbatim:

- The request (2026-11-15, 08:52 BRT): "É possível que no comando /roundfix --status seja informados além do que já temos uma tabela com a fila de specs programada, da ordem de implementacao mais antiga para a mais nova (incluindo possíveis mudanças, ou seja, a fila corrente), com o status e breve observação e para as que ainda não iniciaram, previsão de inicio e fim?"
- The route (2026-11-15): "Backlog agora, Spec depois (Recommended)"; when to author it (2026-11-15, ~20:15 BRT): "Logo após o plano da 0397 (Recommended)".
- The Idea stage (AskUserQuestion, 2026-11-16 ~08:00 BRT): "Rodar a Ideia" — "write-idea em Fable xhigh com perguntas antes do PRD, como na 0397."
- Authorization (AskUserQuestion, 2026-11-16 ~08:00 BRT): "Mesmas do BE (Recommended)" — "Push, PR, merge, release minor depois da Spec, até 3 correções; Makefile/CI/go.mod/lint só se a correção exigir; com o Quaestor ocioso, entregue depois da 0397."
- The shape (AskUserQuestion, 2026-11-16 ~09:15 BRT):
  - "Tamanho da V1" → "Tabela com previsão (Recommended)" — "O pedido original, neste repositório: status, observação e previsão pelo histórico de Runs do próprio repositório."
  - "Linhas além da fila" → "Autoria e backlog (Recommended)" — "Entram Specs em Ideia, PRD ou TechSpec (pelo registro de fase da 0397) e backlog aberto, marcados 'ainda não é Spec', com a espera declarada e sem previsão de fim."
  - "Fim previsto" → "Merge no main (Recommended)" — "Fim é a entrega integrada (Run, revisão, checks e merge); a release segue na linha 'próxima release' do status."
  - "Formato" → "Provável e limite (Recommended)" — "Dois horários por célula, ex.: '~10:30, até 12:10' (mediana e percentil 85), com a base no rodapé."
  - "Esperas" → "Condição e duração (Recommended)" — "Início 'depois de \<condição\>' e fim como duração após o início; as de trás encadeiam igual, sem horário inventado."
  - "Validação" → "Backtest antes (Recommended)" — "No PRD, o método é repetido sobre ~100 Specs dos últimos 30 dias e só é aceito se o fim real cair dentro do limite em pelo menos 80%."
- A Spec being authored (AskUserQuestion, 2026-11-16 ~09:50 BRT): "Mostra previsão (Recommended)" — "Uma Spec em autoria recebe o provável e o limite pelo histórico de plano-até-merge; só o backlog fica sem fim."
- The approval of the Idea (AskUserQuestion, 2026-11-16): "Aprovar e seguir (Recommended)".
- The basis of the forecast of a Spec being authored (AskUserQuestion, 2026-11-16 ~11:40 BRT, asked with the backtest's result, plan to merge at 77.1%): "Aceitar a troca (Recommended)" — "Spec em autoria mostra 'depois do merge do plano' + duração provável e limite da entrega; sua regra de 80% é respeitada."
- The standing condition on performance (2026-11-16, made for Spec 0397 and carried to this Spec by the operator's brief): "com queda de performance nula ou quase nula".

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the answers of the cycles before it. Those release to the Runs the owned skills and the changes of skills that reach the adopters: the 2026-09-30 grant of every skill ("considere autorizado a ajustar todas as skills se necessário"). They also keep the cycle AJ grant of 2026-10-19 "Mexer em Governed Paths declarados" (option text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica).").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted in this cycle only for a correction that requires one. This Spec's design needs none of them, and it changes no Baseline module, no guide under `docs/agents/` and no Project Config.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release`. The Spec is a `type: feat`, so the Release Ceiling is `minor` and the release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge.

One decision of the PRD departed from the wording of an answer, and the maintainer accepted it. A Spec being authored shows a forecast, as he decided, from the time from start to merge. The plan-to-merge history his earlier answer names held its limit for 77.1% in the backtest, under the 80% his answer on validation requires, and his answer "Aceitar a troca (Recommended)" settles it.

The governed set was measured on 2026-11-16 on the branch `docs/plan-0398`, based on `f4580437`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote nothing to the repository: `git status --porcelain` showed only this Spec's planning files, ADR-0409 and ADR-0410 afterwards, and the probe file never existed in the tree.
- It covered forty-three paths: every file the TechSpec's Build Order names, the Go files it expects the Tasks to create, the planning files of this Spec, the four ADRs it writes or amends and three controls.
- The four paths above are governed: the `SKILL.md` pair of the Roundfix skill and its canonical `follow-up` and `modes` references.
- The three controls, `internal/cli/cli_test.go`, `docs/references/coverage-record.json` and the `Makefile`, are governed, and no Task declares or changes them; none is granted.

## What is not governed

- Every Go source and test file under `internal/queueforecast`, `internal/phaserecord` and `internal/cli` that the Tasks declare, and the fixture under `internal/queueforecast/testdata`.
- The mirrors under `skills/roundfix/references/` and `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/follow-up.md`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`, ADR-0409, ADR-0410 and the amendment notices in ADR-0353 and ADR-0407.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

The repository-owned commands write the skill mirrors and the generated records. This declaration records the regeneration that follows the approved edits and adds no source path.

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

No mirror, version field or record is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint, formatter or test-runner configuration, `skills-lock.json`, any Vendored Skill under `.agents/skills`, any Baseline module, template or Setup Snapshot, any guide under `docs/agents/`, `.roundfixrc.yml`, `docs/references/coverage-record.json`, any Archive Record or `CHANGELOG.md`.
- The status view writes nothing. No command, Task or QA row writes to the real `~/.roundfix` or reads the real Run Database, the real Phase Record or the maintainer's checkout.
- No test, Verification command or QA row:
  - starts an Agent harness, an Agent session or an ACP adapter;
  - reaches the network, calls ai-memory or runs the real `gh`;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository, except a temporary one the test itself created;
  - asserts a wall-clock duration;
  - signals a process it did not start, or ends a process by name.
- The fixture of the replay holds only what the Git history and the Run listing of this repository already show: Spec numbers and times.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
