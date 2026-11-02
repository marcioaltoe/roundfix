---
status: approved
granted: 2026-11-02
action: make the Spec Consistency Check read a Glossary Definition written as a top-level list item or with a spaced em dash, name the accepted forms in the SC-GLOSSARY-MISSING fix and print the line-anchored glossary Verification in the SC-GLOSSARY-UNPLANNED fix, and teach the write-prd glossary guide, the write-tasks Glossary requirement and the Roundfix spec reference the same forms and command
consuming: 0341-a-glossary-definition-written-as-a-list-item-counts
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-prd/references/glossary.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-tasks/SKILL.md
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

# Approved authority for Spec 0341

On 2026-11-02 the maintainer answered the operator's cycle AV questions with
AskUserQuestion, after Pantheon reported that the glossary check reads only a
definition line that starts with `**Term**:`. The answers this Spec rests on
are, verbatim:

- Scope: "Aceitar lista + mensagem (Recommended)" — option text: "Reconhecer
  também `- **Termo**:` e `- **Termo** —` (e `**Termo** —`) como definição; a
  mensagem do achado e os guias write-prd/write-tasks dizem os formatos
  aceitos; a Verification gerada pelo write-tasks passa a casar com qualquer
  um deles."
- Release: "v0.129.1 logo após o merge (Recommended)" — option text: "Patch
  único ao fim do ciclo; desbloqueia o Pantheon. Aviso aos 8 adotantes, com
  resposta ao Pantheon."
- Rules: "Mesmas do AU (Recommended)" — option text: "Sem supervisão: push,
  PR do plano, entrega e merge pela fila; Makefile/CI/go.mod/lint liberados
  se o desenho pedir; falha técnica = consertar e seguir; findings do review
  corrigidos ou dispensados com evidência; mensagem do Pantheon consumida
  quando a Spec estiver no main."

Those answers keep the cycle AU rules of 2026-11-01 (Authorization "Mesmas do
AT (Recommended)", option text: "Push e PR dos planos, entrega e merge pela
fila, Makefile/workflows/go.mod/lint quando o desenho pedir, 'Consertar e
seguir' em falha técnica, mudança da Baseline (agent-instructions, release,
spec-routing, skills de autoria) que chega aos adotantes, release v0.129.0 no
fim e aviso aos adotantes."), and through them the cycle AS answers of
2026-10-31 ("Consertar e seguir (Recommended)" and the Tooling answer
"Makefile, .github/workflows, go.mod / go.sum, Config de lint/formatter"), the
cycle AJ grant of 2026-10-19 "Mexer em Governed Paths declarados" (option
text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and the 2026-09-30 grant of every skill
("considere autorizado a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Release answer asks for the patch v0.129.1 right after the
merge, so the Release Ceiling is `patch` and the release may change only the
three version-only files in `release_paths`: the changelog, the npm package
manifest and the application version. No Task performs the release.

The governed set was measured on 2026-11-02 on the branch `docs/plan-0341`,
based on `900e13f8`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain internal` showed no
  change from it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference and
  the three release paths.
- The eight paths above are governed: the canonical Roundfix `spec`
  reference, the canonical `write-prd` glossary guide, the canonical
  `write-tasks` skill, and the `SKILL.md` pairs of the three skills whose
  version fields rise. The canonical `domain-modeling` skill is governed too,
  but task_01 only reads it as an instruction, so it is not listed.
  `dist/npm/roundfix/package.json` is governed too; it is a release path, not
  a Task path.

## What is not governed

- `internal/speccheck/glossary.go`, `internal/speccheck/glossary_test.go`,
  `skills/glossary_definition_test.go`.
- The mirrors `skills/roundfix/references/spec.md` and
  `skills/write-prd/references/glossary.md`, and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/spec.md`, `docs/references/behavior-surfaces.json`,
  `CONTEXT.md`.
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
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-tasks/SKILL.md
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
  the `domain-modeling` skill, or the `### QA settlement` section of any skill.
- No change to the Glossary Declaration grammar, the bold candidate rule, the
  glossary horizon, the binding rule (a Verification naming the exact bolded
  term) or the glossary file set.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
