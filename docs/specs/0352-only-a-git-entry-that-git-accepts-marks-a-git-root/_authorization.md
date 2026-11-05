---
status: approved
granted: 2026-11-05
action: make the Git root walk of configuration resolution and of the skills project scope accept only a .git entry that Git accepts, give the tests one helper that writes the smallest .git directory Git accepts, and document the rule in the Roundfix skill, the configuration and runs guides, ADR-0362 and ADR-0357's amendment notice
consuming: 0352-only-a-git-entry-that-git-accepts-marks-a-git-root
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/review-runs.md
  - .agents/skills/roundfix/references/runs.md
  - internal/cli/cli_test.go
  - internal/gittest/gittest.go
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

# Approved authority for Spec 0352

On 2026-11-05 at about 07:00 BRT the maintainer answered the operator's cycle
AZ question with AskUserQuestion, after his MacBook showed that a `~/.git`
holding only `gk/` and `info/` made Roundfix take the home directory for a Git
root (`references/2026-11-05-a-stray-home-git-directory-makes-the-home-a-git-root.md`).
The answer this Spec rests on is, verbatim:

- "Spec agora + v0.132.1 (Recommended)" — option text: "findGitRoot passa a
  aceitar só um .git que o Git reconhece: pasta com HEAD e objects/, ou
  arquivo 'gitdir:'. Um teste cobre um ~/.git falso. Mesmas regras do AY
  (merge pela fila, notas aos 8 adotantes)."

"Mesmas regras do AY" keeps the cycle AY rules of 2026-11-04 (Rules "Mesmas
do AX (Recommended)", option text: "Sem supervisão; merge pela fila; corrigir
e continuar nas paradas; Makefile, CI, go.mod e lint liberados; mudanças de
Baseline e skills que chegam aos adotantes liberadas."), and through them the
cycle AX, AW, AV, AU, AT and AS answers, the cycle AJ grant of 2026-10-19
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os
caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere autorizado
a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The answer asks for v0.132.1 after this Spec merges, so the Release
Ceiling is `patch` and the release may change only the three version-only
files in `release_paths`: the changelog, the npm package manifest and the
application version. The release is a hotfix and stays outside the limit of
one non-hotfix release per day. No Task performs it; the operator cuts it
after the merge.

The governed set was measured on 2026-11-05 on the branch `docs/plan-0352`,
based on `38cc5854`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference and
  ADR-0362.
- The six paths above are governed: the Roundfix skill's `SKILL.md` pair and
  its canonical `review-runs` and `runs` references, the CLI test file
  `internal/cli/cli_test.go`, and the `gittest` helper
  `internal/gittest/gittest.go`.

## What is not governed

- `internal/config/config.go`, `internal/config/profile_config.go`,
  `internal/cli/cli.go`, the new test files
  `internal/config/git_root_test.go` and
  `internal/cli/rejected_git_directory_test.go`, and every other test file
  task_01 changes under `internal/config`, `internal/cli` and
  `internal/judge`.
- The mirrors `skills/roundfix/references/review-runs.md` and
  `skills/roundfix/references/runs.md`, and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/configuration.md`, `docs/user-guide/commands/runs.md`,
  `docs/references/behavior-surfaces.json`, ADR-0357 and ADR-0362.
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
  `docs/agents/`, `docs/references/skill-coverage.json`,
  `docs/references/coverage-record.json`, `CONTEXT.md`, any Archive Record,
  any skill other than `roundfix`, any Roundfix reference other than
  `review-runs` and `runs`, or the `### QA settlement` section of any skill.
- No change to `findBaselineRepositoryRoot`, `internal/store`,
  `internal/worktree`, `internal/spec`, `internal/speccheck` or
  `internal/baseline`.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
