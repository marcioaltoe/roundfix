---
status: approved
granted: 2026-11-09
action: give roundfix reconcile a digest-confirmed --discard-dirty switch that records a Branch Disposition with the discarded paths before it removes an ended Run's Run Worktree and Run Branch whose Spec is archived on the default branch and whose commit proof holds, classify a recorded Run Worktree that Git does not register and that holds only ignored files safe, and remove empty repository directories under the Worktree Root on a full-scan --apply, as ADR-0379 decides; and document it in ADR-0053, ADR-0289, ADR-0290, CONTEXT.md, the reconcile user guide and the Roundfix skill's reconcile reference
consuming: 0369-reconcile-discards-confirmed-dirty-leftovers-and-ignored-only-worktrees
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/reconcile.md
  - internal/cli/cli_test.go
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

# Approved authority for Spec 0369

On 2026-11-09 Fiscus reported two gaps of `roundfix reconcile`
(`references/2026-11-09-fiscus-reconcile-leaves-abandoned-run-leftovers-without-a-path.md`,
reported by fiscus). The maintainer's answers this Spec rests on
(AskUserQuestion via `/roundfix --follow-up`, 2026-11-09 ~11:20 BRT) are,
verbatim:

- Question "Fiscus relatou duas lacunas do `roundfix reconcile`: (1) sobra
  suja de Run terminal (Failed/Unresolved) de Spec já arquivada fica 'dirty'
  sem caminho suportado; (2) diretório de Run Worktree não registrado no Git,
  só com artefatos ignorados de runtime, fica 'unknown' para sempre (e
  diretórios de grupo vazios ficam). O que fazer?" → "Spec E agora
  (Recommended)" — option text: "Spec fix: modo explícito com confirmação
  (ex. --discard-dirty) só para Run terminal cuja Spec está arquivada no
  default branch, gravando Branch Disposition com os paths descartados;
  diretório não registrado só com artefatos ignorados vira liberável; grupos
  vazios removidos. Patch após o merge."
- Question "Para as novas Specs (D testes de disco cheio, E e F se
  aprovadas), valem as mesmas regras do ciclo BA e o corte da meia-noite?" →
  "Mesmas do BA + corte (Recommended)" — option text: "Sem supervisão, merge
  pela fila, teto de 2 corretivas, release patch assim que cada Spec ficar
  pronta, Makefile/CI/go.mod/lint liberados, Jarvis avisado. Nenhum item
  começa depois das 21:00 BRT; a fila para antes da meia-noite e o resto fica
  para amanhã."

"Mesmas do BA" keeps the cycle BA rules, and through them the cycle AJ grant
of 2026-10-19 "Mexer em Governed Paths declarados" (option text: "'Concedo'
para os caminhos que cada Spec declarar (módulos Baseline, registros,
skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem
autorização específica).") and the 2026-09-30 grant of every skill
("considere autorizado a ajustar todas as skills se necessário"). The design
needs no Makefile, workflow, Go module, lint or formatter change, although
the answers would allow them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`.
No Task performs it; the operator cuts the patch right after the merge.

The governed set was measured on 2026-11-09 on the branch `docs/plan-0369`,
based on `cfc08683`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference and
  ADR-0379.
- The four paths above are governed: the canonical `reconcile` reference
  and the Roundfix `SKILL.md` pair whose version fields task_03 raises, and
  `internal/cli/cli_test.go`, whose `reconcile` help row task_01 updates to
  the new synopsis. The mirror `skills/roundfix/references/reconcile.md` and
  the read-only `domain-modeling` skill are not governed.

## Sanctioned regeneration

The repository-owned commands write the skill mirror and the generated
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
  `reconcile`, or the `### QA settlement` section of any skill.
- No change to `internal/cli/deliver_disk_full*.go`,
  `internal/cli/doctor_disk_space.go`, `internal/spec/**`,
  `internal/speccheck/**` or `internal/baseline/**`, which other Specs of
  cycle BA own.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`, or runs `roundfix reconcile`
    with `--apply` or `--discard-dirty` outside a temporary home and
    repository;
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
