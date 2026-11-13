---
status: approved
granted: 2026-11-13
action: inventory a nested AGENTS.md or CLAUDE.md only when Git lists it as a file of the repository, so an ignored copy and a carrier inside a linked worktree, a nested repository or a submodule raise no nested-carrier-conflict, and document it in the baseline guide, the Roundfix skill's baseline reference and the ADR-0070 amendment notice
consuming: 0390-a-nested-carrier-git-does-not-list-raises-no-conflict
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - internal/baseline/repository_test.go
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

# Approved authority for Spec 0390

On 2026-11-13 quaestor reported, in ai-memory message
`01a1d4de-8d39-7f60-b575-5703959df8da`, that `roundfix baseline update`
warned `baseline.inventory.nested-carrier-conflict` for the `AGENTS.md` and
`CLAUDE.md` of a linked worktree kept under a directory its `.gitignore`
lists. The message is copied verbatim in this Spec's `references/`
(`references/2026-11-13-quaestor-a-nested-carrier-conflict-inside-a-git-ignored-worktree.md`).
The maintainer then answered the operator with AskUserQuestion via
`/roundfix --follow-up` (2026-11-13 ~10:30 BRT). The answer this Spec rests on
is, verbatim:

- Question "O Quaestor relatou que o `roundfix baseline update` acusa
  conflito nos AGENTS.md e CLAUDE.md de um worktree guardado numa pasta que o
  Git ignora (.storage/). Esses arquivos não são do repositório. Viro isso na
  Spec 0390 e entrego depois da 0389?" → "Sim, mesmas regras do BE
  (Recommended)" — option text: "O inventário passa a pular caminhos
  ignorados pelo Git e raízes de worktrees ligados. Push, merge e patch logo
  após o merge; até 3 correções; Makefile, CI, go.mod e lint liberados se a
  correção exigir."

"Mesmas regras do BE" carries the cycle BE answers of 2026-11-13 ~08:20 BRT:
a patch release after each Spec ("Patch após cada Spec (Recommended)") and up
to three corrections per Spec with Makefile, CI workflow, `go.mod` and lint
changes granted when a correction requires one ("Até 3 por Spec + Makefile/CI
(Recommended)"). The cycle BE rules keep the earlier cycles' grants of owned
skills and of changes to skills that reach the adopters, the cycle AJ grant of
2026-10-19 "Mexer em Governed Paths declarados" (option text: "'Concedo' para
os caminhos que cada Spec declarar (módulos Baseline, registros, skills,
testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere autorizado
a ajustar todas as skills se necessário").

This Spec's design needs no Makefile, `.github/workflows`, `go.mod`, `go.sum`,
lint or formatter configuration change, and no Baseline module, template or
profile change. The answer would grant a Makefile, workflow, `go.mod` or lint
change if a correction required one; none is listed here, and any such
change must first be added to `paths` on the default branch with its reason.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it after the merge.

The governed set was measured on 2026-11-13 on the branch `docs/plan-0390`,
based on `2e7a2a7a`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  planning files, its adopted reference and ADR-0399 afterwards, and the probe
  file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every file the scratch prototype of tasks 01 and 02 changed
  or created, every planning file of this Spec, its adopted reference, ADR-0399,
  `CONTEXT.md` and `docs/agents/setup-context.json`.
- The four paths above are governed: the Roundfix skill's `SKILL.md` pair and
  its canonical `baseline` reference, which task_02 changes, and
  `internal/baseline/repository_test.go`, which task_01 changes. The probe
  also named `.agents/skills/domain-modeling/SKILL.md` and
  `docs/agents/setup-context.json`, which no Task changes; neither is
  granted.

## What is not governed

- `internal/baseline/repository.go`, `internal/baseline/plan.go`,
  `internal/baseline/apply.go` and `internal/baseline/transaction.go`, which
  task_01 changes.
- The mirror `skills/roundfix/references/baseline.md` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/baseline.md`, `docs/references/behavior-surfaces.json`,
  ADR-0070 and ADR-0399.
- Every planning file of this Spec and its adopted reference.

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

No mirror, version field, digest or record is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`,
  `skills-lock.json`, any Vendored Skill under `.agents/skills`, any Baseline
  module, template, profile, Setup Snapshot or plan golden, any guide under
  `docs/agents/`, `docs/agents/setup-context.json`, `CONTEXT.md`,
  `docs/references/coverage-record.json`, any Archive Record, any skill other
  than `roundfix`, any Roundfix reference other than `baseline`, or the
  `### QA settlement` section of any skill.
- No test, Verification command or QA row:
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message or a Spec Event Announcement; the
  operator does.
