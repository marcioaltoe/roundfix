---
status: approved
granted: 2026-11-10
action: make the related-ADR seed rule of roundfix spec check read a closed list of not-applicable tokens, the English not applicable and the Portuguese não se aplica with its unaccented and decomposed spellings, name both in the SC-ADR-RELATED fix text, write the ADR-0389 amendment notice into ADR-0377, and align the glossary, the write-techspec and write-tasks skills, the Roundfix spec reference and the spec and Context-Driven Development guides
consuming: 0379-an-adr-marked-as-not-applying-in-portuguese-draws-in-no-citing-adr
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/write-techspec/references/techspec-template.md
  - skills/roundfix/SKILL.md
  - skills/write-tasks/SKILL.md
  - skills/write-techspec/SKILL.md
  - skills/write-techspec/references/techspec-template.md
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

# Approved authority for Spec 0379

On 2026-11-10 the operator recorded that an ADR marked as not applying in
Portuguese still seeds the related ADRs
(`references/2026-11-10-an-adr-marked-as-not-applying-in-portuguese-still-seeds-the-related-adrs.md`,
recorded by the operator from Fiscus's message whose id begins `01a1c323`).
This Spec adopts the entry and amends ADR-0377 with ADR-0389.

The maintainer's answers this Spec rests on (AskUserQuestion via
`/roundfix --follow-up`, 2026-11-10 ~20:55–21:00 BRT) are, verbatim:

- Question "Quais Backlog Entries viram Spec agora, depois da 0374? As que
  você não marcar ficam abertas no backlog." → selected, among four,
  "Marcador \"não se aplica\"" — option text: "Aceitar o marcador em
  português nas ADR Alternatives. O Fiscus não pediu e não está travado."
- Question "Com que autorizações eu toco sem supervisão as seis Specs novas:
  0374 (capacidade da máquina), 0375 (orçamento de núcleos), 0376 (skills de
  usuário velhas), 0377 (portão do CI), 0378 (digest do bump de skills) e
  0379 (marcador \"não se aplica\")?" → "Mesmas do BB, teto minor
  (Recommended)" — option text: "Escrever, entregar pela fila, merge e
  release de cada Spec assim que pronta (patch nas correções, minor nas
  features; 0374 e 0375 saem juntas), teto de 2 Tasks corretivas, CI na 0377
  e Makefile na 0378 só nos arquivos que elas exigirem, eventos ao Jarvis e
  resposta ao Vortex quando a 0374 e a 0375 saírem."

"Mesmas do BB" keeps the cycle BB rules of 2026-11-10 ("Mesmas do BA, teto
minor (Recommended)" — "Escrever, entregar pela fila, merge e release de
cada Spec assim que pronta (patch nas duas correções; a regra pode sair como
minor), teto de 2 Tasks corretivas, eventos ao Jarvis, sem notas às caixas,
e Makefile/CI/go.mod/lint só onde a Spec exigir."), and through them the
cycle AJ grant of 2026-10-19 "Mexer em Governed Paths declarados" (option
text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and the 2026-09-30 grant of every skill
("considere autorizado a ajustar todas as skills se necessário"). The design
needs no Makefile, workflow, Go module, lint or formatter change.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: feat`, so the Release Ceiling is `minor` and
the release may change only the three version-only files in `release_paths`.
No Task performs it; the operator cuts the minor right after the merge.

The governed set was measured on 2026-11-10 on the branch `docs/plan-0379`,
based on `dfb4dbac` and rebased on `b6cd0f14`, which changes none of the
probed files, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, ADR-0389, its adopted
  reference and the Backlog Entry's former path.
- The nine paths above are governed: the canonical and mirror `SKILL.md`
  files of `roundfix`, `write-tasks` and `write-techspec` whose version
  fields task_01 raises, the canonical `spec` reference and the canonical and
  mirror TechSpec template task_01 edits. The `spec` reference's mirror, the
  version record, the guides, `CONTEXT.md`, the ADRs, the Behavior Surface
  Record and every `internal/speccheck` file are not governed.

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
  - .agents/skills/write-tasks/SKILL.md
  - skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
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
  formatter or test-runner configuration, `.roundfixrc.yml`,
  `skills-lock.json`, any Baseline module, profile, template or derived
  Baseline file, any guide under `docs/agents/`,
  `docs/references/skill-coverage.json`, any ADR other than ADR-0377, any
  Archive Record, any skill other than `roundfix`, `write-tasks` and
  `write-techspec`, any Roundfix reference other than `spec`, or the
  `### QA settlement` section of any skill.
- No change to `internal/speccheck/constraints.go` (the Project Constraints
  parser) or to any finding code, severity or location.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - creates files in `/tmp`, `/private/tmp` or `$HOME` outside its own
    temporary directory;
  - starts an Agent session or an ACP adapter;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
