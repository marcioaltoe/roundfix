---
status: approved
granted: 2026-11-09
action: make the strict Spec Consistency Check report, as SC-UNREACHABLE-UNREAD, text in a PRD's Unreachable Acceptance section that no fielded declaration reads, through one reader the QA gate shares; make SC-ADR-RELATED seed its candidates only with the listed ADRs whose TechSpec ADR Alternatives line does not begin with not applicable, as ADR-0377 decides; and document both in ADR-0093, ADR-0168, CONTEXT.md, the spec command and context-driven development guides, the Roundfix spec reference and the write-techspec and write-tasks skills
consuming: 0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/write-techspec/references/techspec-template.md
  - internal/speccheck/coherence.go
  - internal/speccheck/constraints.go
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
release_ceiling: patch
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0367

On 2026-11-07 the maintainer, through jarvis, asked for the Backlog Entry
"The strict check names an Unreachable Acceptance entry it cannot read"
(`references/2026-11-07-the-strict-check-names-an-unreachable-acceptance-it-cannot-read.md`),
and Fiscus reported that `SC-ADR-RELATED` grows through ADRs a Spec lists as
not applying
(`references/2026-11-07-fiscus-sc-adr-related-grows-through-not-applicable-adrs.md`,
reported by fiscus). The maintainer's answers this Spec rests on
(AskUserQuestion, 2026-11-09 ~08:35–08:45 BRT) are, verbatim:

- Question "Finding do Fiscus: a regra SC-ADR-RELATED (ADRs 0093/0168
  aceitas) faz uma ADR marcada “não se aplica” puxar todas as ADRs que a
  citam (72 ADRs numa Spec). Emendo a ADR?" → "Emendar (Recommended)" —
  option text: "Só ADRs aplicadas ou alteradas pela Spec entram no fecho; uma
  ADR “não se aplica” deixa de propagar. Task com “Amended by” na ADR
  aceita."
- Question "Os dois itens da checagem estrita (backlog do Unreachable
  Acceptance ilegível + cascata do SC-ADR-RELATED do Fiscus) vão numa Spec só
  ou em duas?" → "Uma Spec (Recommended)" — option text: "Mesmo pacote
  (internal/speccheck), ~2 Tasks + QA, um release patch. Jev: coerente 0,85;
  dividir 0,44."
- Question "Que autorizações você concede de saída para eu levar o ciclo BA
  até o fim sem parar (implementações, merges e releases inclusos como no
  AZ)?" → selected "Mesmas do AZ (Recommended)" — "Release patch assim que
  cada Spec ficar pronta, teto de 2 corretivas por Spec, merge pela fila,
  avisos ao Jarvis, sem notas de release às caixas."; "Scripts/Makefile/CI" —
  "Mexer em ~/.roundfix-operator, Makefile e .github/workflows se a migração
  ou a Spec B precisarem."; "go.mod/lint" — "Mexer em go.mod/go.sum e na
  configuração de lint/formatter se uma Task precisar."; "Avisar
  Quaestor/Fiscus" — "Mensagem ai-memory curta a cada um dizendo qual Spec
  adotou o finding (não é nota de release)."

"Mesmas do AZ" keeps the cycle AZ rules, and through them the cycle AJ grant
of 2026-10-19 "Mexer em Governed Paths declarados" (option text: "'Concedo'
para os caminhos que cada Spec declarar (módulos Baseline, registros,
skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem
autorização específica).") and the 2026-09-30 grant of every skill
("considere autorizado a ajustar todas as skills se necessário"). The first
answer's option text grants the amendment of the accepted ADRs through a
Task. The design needs no Makefile, workflow, Go module, lint or formatter
change, although the answers would allow them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`.
No Task performs it; the operator cuts the patch right after the merge.

The governed set was measured on 2026-11-09 on the branch `docs/plan-0367`,
based on `2e59aebe`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted references,
  the removed Backlog Entry path and ADR-0377.
- The eleven paths above are governed: the canonical and mirrored owned
  skill files task_03 changes, and `internal/speccheck/coherence.go` and
  `internal/speccheck/constraints.go`, which an earlier authorization bounded
  and task_01 changes.

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
  - .agents/skills/write-techspec/SKILL.md
  - skills/write-techspec/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
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
  any skill other than `roundfix`, `write-techspec` and `write-tasks`, any
  Roundfix reference other than `spec`, `internal/cli/cli_test.go`, or the
  `### QA settlement` section of any skill.
- No change to `internal/spec/qa*.go`, `internal/spec/archive.go` or
  `internal/agent/**`, which other Specs of cycle BA own; `Unreachable` keeps
  the declarations and errors it returns.
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
