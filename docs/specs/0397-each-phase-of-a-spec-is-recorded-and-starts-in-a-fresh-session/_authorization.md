---
status: approved
granted: 2026-11-14
action: add the Phase Record under Roundfix Home with the roundfix phase start command and its Boundary Notice, record each prompt of a Run and of a pre-PR review in it, remove old day files in the Run Retention Sweep, state the session boundary rule and the recording step in the four authoring skills, the Roundfix skill and the autonomous-work guide of the Baseline, install a published authoring skill without the Declared Skill Header, keep a fix or a refactor in the session, and declare this repository's authoring models in its Project Config with the generated subagents
consuming: 0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session
paths:
  - .roundfixrc.yml
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/modes.md
  - .agents/skills/roundfix/references/phase.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/roundfix/references/storage.md
  - .agents/skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - docs/agents/autonomous-work.md
  - docs/agents/setup-context.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md
  - internal/baseline/assets/modules/autonomous-work.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - skills/roundfix/SKILL.md
  - skills/write-idea/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-tasks/SKILL.md
  - skills/write-techspec/SKILL.md
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

# Approved authority for Spec 0397

The Spec adopts one ai-memory message, reported by jarvis on 2026-11-14, and the maintainer approved its Idea on 2026-11-16. The maintainer's answers are recorded whole in [`references/2026-11-15-maintainer-answers.md`](references/2026-11-15-maintainer-answers.md). The answers this record rests on are, verbatim:

- Authorization (AskUserQuestion, 2026-11-14 ~11:35 BRT). Question "Autorizações para a Spec 0396 (Indexus) e para a 0397 (limite do Claude): valem as mesmas do ciclo BE?" → "Mesmas do BE (Recommended)" — option text: "Push, PR, merge, release depois de cada Spec, até 3 correções; Makefile/CI/go.mod/lint só se a correção exigir. Com o Quaestor ocioso, 1 vaga."
- The route and the scope (2026-11-14). The message from jarvis → "Medidas já, sem Spec" and "Spec com Ideia após o reset"; adoption of the per-phase models in this repository → "Adotar na Spec 0397"; automatic compaction refused: "passa para determinístico e pode perder informação importante". The scope as the operator recorded it: "o Run grava o id e a fase num lugar fixo para o ai_usage.py; o Supervisor fecha a sessão em cada fronteira de fase".
- The Idea stage (2026-11-15). "Manter a Ideia". The research subagent's model: "Sonnet, effort medium (Recommended)".
- The shape (the night of 2026-11-15 to 2026-11-16, about 00:00 BRT):
  - "Cobertura do registro de fase" → "Run + sessões Supervisor (Recommended)" — "Registra sessão, Spec e fase também fora do Run (autoria, entrega, release, triagem) e, dentro do Run, a Task ou o QA em curso."
  - "Quem relata o consumo por fase" → "Roundfix só expõe (Recommended)" — "Lugar e formato documentados e estáveis; o ai_usage.py cruza com os transcripts; o Roundfix não lê transcript de harness nem calcula preço."
  - "Fronteira de sessão do Supervisor" → "Instrução e aviso (Recommended)" — "Skills e guia mandam gravar handoff, dizer o que ficou armado e encerrar; o Roundfix avisa, sem bloquear, quando a mesma sessão registra duas fases; você abre a próxima com /roundfix follow-up."
  - "Adoção dos modelos por fase no Roundfix" → "Só subagentes gerados (Recommended)" — "Na 0397: autores por fase e pesquisador em Sonnet em .claude/agents/, sem cabeçalho nas skills deste repositório; uma skill digitada segue no modelo da sessão."
  - "Correção pequena × declaração por fase" → "Exceção pela rota (Recommended)" — "Na rota de correção ou refactor a skill escreve na sessão e só delega quando você pedir; feature sempre delega; muda o texto das skills para todos os adotantes."
  - "Direção da Ideia" → "Original + respostas (Recommended)" — "Registro de fase, fronteira de sessão e adoção dos modelos numa Spec; o custo de nascer de cada sessão vira Backlog Entry."
- The success criterion (2026-11-16, in the session): "Adicione ao critério de sucesso: com queda de performance nula ou quase nula e com aumento de qualidade de artefatos para implementações."
- The approval of the Idea (AskUserQuestion, 2026-11-16): "Aprovar e seguir (Recommended)".

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the answers of the cycles before it. Those release to the Runs the owned skills and the changes of skills that reach the adopters: the 2026-09-30 grant of every skill ("considere autorizado a ajustar todas as skills se necessário"). They also release the Baseline source with its derived files, the guides under `docs/agents/` and this repository's Project Config ("Autorizar os dois"), and the cycle AJ grant of 2026-10-19 "Mexer em Governed Paths declarados" (option text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica).").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted in this cycle only for a correction that requires one. This Spec's design needs none of them.

This repository's Project Config is in this Spec's scope: the answer "Adotar na Spec 0397" and the answer "Só subagentes gerados (Recommended)" decide that this Spec adds the Authoring Model Declaration to `.roundfixrc.yml`. The declaration uses only keys the released binary already knows, so the installed binary of the Delivery Queue keeps reading the file after the merge.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release`. The Spec is a `type: feat`, so the Release Ceiling is `minor` and the release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge.

The governed set was measured on 2026-11-16 on the branch `docs/plan-0397`, based on `2caa1d17`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote nothing to the repository: `git status --porcelain` showed only this Spec's planning files and ADR-0407 and ADR-0408 afterwards, and the probe file never existed in the tree.
- It covered eighty-six paths: every file the TechSpec's Build Order names, the Go files it expects the Tasks to create, the generated subagent files, the planning files of this Spec, the two ADRs and three controls.
- The twenty paths above are governed: this repository's Project Config; the `SKILL.md` pairs of the Roundfix skill and of the four authoring skills; the Roundfix skill's canonical `modes`, `setup` and `storage` references and its new `phase` reference; the autonomous-work module of the Baseline, the profile that pins its digest and the formatter golden of its guide; and the generated guide `docs/agents/autonomous-work.md` with the Setup Manifest.
- The three controls, `internal/cli/cli_test.go`, `docs/references/coverage-record.json` and the `Makefile`, are governed, and no Task declares or changes them; none is granted.

## What is not governed

- Every Go source and test file under `internal/phaserecord`, `internal/cli`, `internal/daemon`, `internal/baseline` and `skills` that the Tasks declare.
- The five generated files under `.claude/agents/`.
- The mirrors under `skills/roundfix/references/`, `skills/testdata/owned-skill-versions.json`, `internal/baseline/module-versions.json` and the catalog records under `internal/baseline/testdata/`.
- `docs/user-guide/**`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`, ADR-0407, ADR-0408 and the amendment notice in ADR-0401.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

The repository-owned commands write the skill mirrors, the generated subagents and the generated records. This declaration records the regeneration that follows the approved edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/write-idea/SKILL.md
  - skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - skills/write-techspec/SKILL.md
```

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/assets/modules/autonomous-work.json
  - internal/baseline/module-versions.json
```

```yaml
command: make baseline-digests
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/autonomous-work.md
  - docs/agents/setup-context.json
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix skills install --target project
outputs:
  - .claude/agents/roundfix-idea-author.md
  - .claude/agents/roundfix-prd-author.md
  - .claude/agents/roundfix-researcher.md
  - .claude/agents/roundfix-tasks-author.md
  - .claude/agents/roundfix-techspec-author.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

No mirror, version field, digest, generated subagent or record is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint, formatter or test-runner configuration, `skills-lock.json`, any Vendored Skill under `.agents/skills`, any Baseline module other than `autonomous-work`, any template or Setup Snapshot, any guide under `docs/agents/` other than `autonomous-work.md` and the Setup Manifest, `docs/references/coverage-record.json`, any Archive Record, any owned skill other than `roundfix`, `write-idea`, `write-prd`, `write-techspec` and `write-tasks`, or the `### QA settlement` section of any skill.
- `.roundfixrc.yml` gains only the `skills.authoring.claude` section, with keys the released binary already reads. No other key of Project Config changes.
- No shipped `SKILL.md`, canonical or mirrored, gains a `model` or `effort` key.
- No command, Task or QA row writes to the real `~/.roundfix`, to a user-level skill directory of the host or to the operator's usage script, which lives outside this repository.
- No test, Verification command or QA row:
  - starts Claude Code, another Agent harness, an Agent session or an ACP adapter;
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads a harness transcript, or reads or writes the real home directory;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value other than a session identifier the test itself set, or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
