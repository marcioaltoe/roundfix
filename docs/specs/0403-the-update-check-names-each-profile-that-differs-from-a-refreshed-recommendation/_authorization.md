---
status: approved
granted: 2026-11-17
action: refresh the Recommended Profile to the 2026-11-17 snapshot the maintainer approved, make the Model Catalog follow the installed adapters, make the Update Check and the --update Skill Mode name each configured profile that differs with the command that fixes it, align this repository's Project Config, and describe all of it in the Roundfix skill, the guides, the glossary and the amended ADR-0181
consuming: 0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation
paths:
  - .roundfixrc.yml
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/modes.md
  - .agents/skills/roundfix/references/profiles.md
  - .agents/skills/roundfix/references/setup.md
  - internal/docscontract/publicdocs_test.go
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

# Approved authority for Spec 0403

The maintainer's answers are recorded whole in [`references/2026-11-17-maintainer-answers.md`](references/2026-11-17-maintainer-answers.md). The answers this record rests on are, verbatim (AskUserQuestion, 2026-11-17 ~12:40 BRT):

- Pediu "revisão e sugestão/atualização dos modelos já configurados para implementações e revisões. Já peguei configuração antiga sem atualizar." → "Spec 0403 separada (Recommended)": atualiza o Recommended Profile (pesquisa; tabela proposta passa pelo Marcio antes de aplicar) e o update check / --update mostra cada categoria que difere com o comando para corrigir; emenda ADR-0181; depois da 0402. […] Autorização 0403 → "Mesmas do BE (Recommended)".

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the grants of the cycles before it. These include the 2026-09-30 grant of every skill ("considere autorizado a ajustar todas as skills se necessário") and the cycle AJ grant of 2026-10-19 "Mexer em Governed Paths declarados", whose option text reads: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica)." The operator's brief of 2026-11-17 adds `.roundfixrc.yml` to this Spec's paths: its `review`, `docs` and `chore` profiles follow the approved table, and its comments drop the stale `gpt-5.6-sol` measurement and Fable 5.1 notes.

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted only for a correction that requires one. This Spec's design needs none of them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release`. The Spec is a `type: fix`. Its Release Ceiling is `minor`, not `patch`, because the design adds a public output: the `recommendations` member of the `roundfix/update-check/v1` JSON report. The release may change only the three version-only files in `release_paths`. No Task performs it; the operator cuts it after the merge, up to three corrections included.

## How the governed set was measured

The author measured the governed set on 2026-11-17 in a clone on the branch `docs/plan-0403`, based on `63a9f493`, with `speccheck.GovernedPath`.

- The measurement used a `go test -overlay` probe in `internal/speccheck` that wrote nothing to the repository. The probe file never existed in the tree.
- It covered 39 paths: every source, test, guide, skill, record, configuration and document file the TechSpec names or expects the Tasks to create, ADR-0181, ADR-0416, `internal/cli/cli_test.go`, and three controls.
- The seven paths above are the governed ones this Spec's design changes.
- These governed paths are not granted, and no Task may change them:
  - `internal/cli/cli_test.go`; new tests go in new files;
  - the three controls: the `Makefile`, `.github/workflows/ci-verify.yml` and `go.mod`.

The author also applied the new snapshot and catalog to a scratch worktree of `63a9f493`, removed afterwards, and ran `./internal/config`, `./internal/agent`, `./internal/cli`, `./internal/tui`, `./skills` and `-tags docscontract ./internal/docscontract`. The tests that failed are the ones the Tasks declare; `internal/cli/cli_test.go` passed. With a sentence added to each changed guide and skill reference, the sanctioned regeneration below changed exactly the `SKILL.md` pair (version `0.1.248` to `0.1.249`), the three mirrored references, `skills/testdata/owned-skill-versions.json` and three guide entries of `docs/references/behavior-surfaces.json`. No governed path outside `paths` changed.

## What is not governed

- Every Go source and test file under `internal/config`, `internal/agent` and `internal/cli` that the Tasks declare, `internal/cli/cli.go` included.
- The mirrors under `skills/roundfix/references/`, `skills/testdata/owned-skill-versions.json` and the tests under `skills/`.
- `docs/user-guide/usage.md`, `docs/user-guide/configuration.md`, `docs/user-guide/commands/update.md`, `docs/references/model-selection.md`, `docs/references/behavior-surfaces.json`, `CONTEXT.md`, ADR-0181 and ADR-0416.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

Repository-owned commands write the skill mirrors and the derived records. This declaration records the regeneration that follows the approved edits and adds no source path.

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

- No change to any of these:
  - the Makefile, the CI workflows, `go.mod` or `go.sum`;
  - the lint, formatter or test-runner configuration;
  - any `.roundfixrc.yml` key outside `profiles` and the comments above it;
  - `skills-lock.json`;
  - any Baseline module, template, profile, Source Baseline or Setup Snapshot, and any file under `docs/agents/`;
  - any other skill, Vendored Skill or Archive Record;
  - the Baseline semantic analysis selections and the light tier;
  - the `### QA settlement` section of any skill.
- No change to any file of another repository, and no read of one.
- No test or Task Verification:
  - reaches the network, calls ai-memory or runs the real `gh`, `codex` or `claude`;
  - reads or writes the real `~/.roundfix`, `~/.codex`, `~/.claude` or the real home directory;
  - starts an Agent session or an ACP adapter, or runs `roundfix profiles validate` or `roundfix profiles check --apply` against real runtimes;
  - sends a request to a provider, the judge or a reviewer;
  - routes an OpenAI or Anthropic model through OpenRouter;
  - signals a process it did not start, or ends a process by name;
  - reads the clock to decide a snapshot, or waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
