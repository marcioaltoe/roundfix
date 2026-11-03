---
status: approved
granted: 2026-11-03
action: give the write-prd, write-techspec, write-tasks and archive-spec skills and the Roundfix skill's implement reference a generic Spec Event Announcement step that exists only when the user's own instructions define its destination and format, add the Spec Event terms to the glossary and the context-driven development guide, and add a skill-contract test that holds the five skills to the generic wording and keeps every owned skill free of a personal assistant's or tracker's name
consuming: 0344-the-spec-skills-announce-the-spec-events-the-user-defines
paths:
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/implement.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - skills/archive-spec/SKILL.md
  - skills/roundfix/SKILL.md
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

# Approved authority for Spec 0344

On 2026-11-03 the maintainer answered the operator's cycle AW questions with
AskUserQuestion, in the `/roundfix --follow-up` mode, after Jarvis asked the
Roundfix-distributed Spec skills to send its Spec board one block per Spec
event. The answers this Spec rests on are, verbatim:

- 0344: "Passo genérico nas skills (Recommended)" — option text: "Spec 0344:
  as skills ganham um 'aviso de evento da Spec' nos cinco pontos, que só
  existe quando as instruções do próprio usuário o definem (mesmo padrão do
  passo de anúncio de release). Nenhuma skill cita Jarvis nem Todoist. O
  formato do bloco e o destino default/jarvis vão para o seu CLAUDE.md
  global; eu proponho o texto."
- Rules (all three): "Mesmas do AV (Recommended)" — option text: "Push e PR
  dos planos, entrega e merge pela Delivery Queue, consertar e seguir em
  falha técnica (pular só se o conserto contradisser o PRD). Makefile,
  .github/workflows, go.mod/go.sum e config de lint/formatter liberados
  quando o desenho pedir."
- Release (all three): "Uma v0.130.0 no fim (Recommended)" — option text:
  "Minor único depois das três Specs, porque a 0342 cria a chave
  pre_pr_review.omit_paths e a 0343 cria o --all. Notas aos 8 adotantes,
  resposta ao Fluxus com os campos waiting_on/since, aviso ao Vortex da chave
  nova e aviso ao Jarvis de que a versão saiu."

Those answers keep the cycle AV rules of 2026-11-02 (Rules "Mesmas do AU
(Recommended)", option text: "Sem supervisão: push, PR do plano, entrega e
merge pela fila; Makefile/CI/go.mod/lint liberados se o desenho pedir; falha
técnica = consertar e seguir; findings do review corrigidos ou dispensados com
evidência; mensagem do Pantheon consumida quando a Spec estiver no main."),
the cycle AU rules of 2026-11-01 (Authorization "Mesmas do AT
(Recommended)"), and through them the cycle AS answers of 2026-10-31
("Consertar e seguir (Recommended)" and the Tooling answer "Makefile,
.github/workflows, go.mod / go.sum, Config de lint/formatter"), the cycle AJ
grant of 2026-10-19 "Mexer em Governed Paths declarados" (option text:
"'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and the 2026-09-30 grant of every skill
("considere autorizado a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Release answer asks for one minor, v0.130.0, after Specs 0342,
0343 and 0344, so the Release Ceiling is `minor` and the release may change
only the three version-only files in `release_paths`: the changelog, the npm
package manifest and the application version. No Task performs the release;
the operator cuts it after the three Specs merge.

The governed set was measured on 2026-11-03 on the branch `docs/plan-0344`,
based on `4e61bf6d`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference,
  ADR-0354 and the three release paths.
- The eleven paths above are governed: the five canonical skill files that
  gain the step (the `write-prd`, `write-techspec`, `write-tasks` and
  `archive-spec` `SKILL.md` files and the Roundfix implement reference), the
  canonical Roundfix `SKILL.md` whose version fields rise, and the five
  `SKILL.md` mirrors. The canonical `domain-modeling` skill is governed too,
  but task_01 only reads it as an instruction, so it is not listed.
  `dist/npm/roundfix/package.json` is governed too; it is a release path, not
  a Task path.

## What is not governed

- `skills/spec_event_announcement_test.go`.
- The mirror `skills/roundfix/references/implement.md` and
  `skills/testdata/owned-skill-versions.json`.
- `CONTEXT.md`, `docs/user-guide/context-driven-development.md`,
  `docs/references/behavior-surfaces.json` and ADR-0354.
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
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - skills/archive-spec/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-tasks/SKILL.md
  - skills/write-techspec/SKILL.md
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
  the `domain-modeling` skill, any Go file outside
  `skills/spec_event_announcement_test.go`, or the `### QA settlement`
  section of any skill.
- No change to the Roundfix skill's `references/modes.md` or
  `references/follow-up.md`, which Spec 0343 owns.
- No owned skill or guide names a personal assistant, a tracker, a channel or
  the reported block, and only the skill-contract test lists those names, as
  the words it forbids; the user's own instructions own them, and this Spec
  never edits them.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session;
  - reaches the network;
  - sends a message, or a request to a provider, the judge or a reviewer;
  - reads another repository or the user's own instructions;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
