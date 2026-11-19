---
status: approved
granted: 2026-11-19
action: add the Flow Command List with its Named Exclusions and Guard Entries, the `roundfix permissions` command that shows it for a harness and writes or removes it only against a previewed digest, one check of the Doctor Command and one field of the Update Check; add the project permission files of the three harnesses to the governed set; describe the command in the Roundfix skill's setup and modes references, in the user guide, in the glossary and in the decision records
consuming: 0413-a-session-runs-the-flows-own-commands-without-a-permission-prompt
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/modes.md
  - .agents/skills/roundfix/references/setup.md
  - internal/speccheck/governed.go
  - internal/speccheck/governed_repocontract_test.go
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

# Approved authority for Spec 0413

The maintainer's request and answers are recorded whole in [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md). The answers this record rests on are, verbatim (2026-11-19):

- the request, in his own words (about 00:15 BRT): "Outra coisa que normalmente trava o fluxo é o uso de comandos bash, se houver algo que possa ser feito para resolver isso tb, uma inclusão de autorização permanente seria muito bom. Mas depende do harness.";
- the split (about 01:10 BRT): "Duas Specs (Recommended)", whose option text reads "(B) autorização permanente de comandos de shell vira a 0413, com Ideia própria partindo desta pesquisa, entregue depois da 0412";
- the target and the form (the same call): "Sessão Supervisor e autores (Recommended)" and "Mostra a lista e grava com seu sim (Recommended)", the second with the option text "comando mostra a lista fechada (comandos roundfix e leituras git/gh, sem curinga de shell nem interpretador) e grava no arquivo de permissões do projeto só com confirmação; push e merge continuam atrás do _authorization.md";
- the approval of the Idea (about 02:20 BRT): "Aprovar e seguir ao PRD (Recommended)", whose option text reads "inclui os arquivos de permissão do harness como Caminhos Governados e nenhuma chave nova no .roundfixrc.yml", with "Os que as skills nomeiam (Recommended)", "Todas as leituras nomeadas" and ".claude/settings.json (Recommended)";
- the open questions of this cycle (about 00:20 BRT): "Seguir o padrão do autor (Recommended)";
- the authorization of Specs 0412 and 0413 (later the same day): "Mesmas do BE para as duas (Recommended)", whose option text reads "implement, commit, push, PR, merge e release minor depois de cada uma; até 5 correções; Makefile/CI/go.mod/lint só se a correção exigir".

## Where the grant of the operations comes from

The grant is the maintainer's answer of 2026-11-19 to the authorization of Specs 0412 and 0413: "Mesmas do BE para as duas (Recommended)", with the option text "implement, commit, push, PR, merge e release minor depois de cada uma; até 5 correções; Makefile/CI/go.mod/lint só se a correção exigir". It grants this Spec `implement`, `commit`, `push`, `pull_request`, `merge` and `release`, a `minor` Release Ceiling and up to five corrections, the same authorization "Mesmas do BE (Recommended)" that he gave on 2026-11-18 for the Specs then planned in this cycle. The operator recorded the answer and told the author of this plan; the author read it in the record before writing this file.

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the grants of the cycles before it:

- the grant of 2026-09-30 for every skill: "considere autorizado a ajustar todas as skills se necessário";
- the cycle AJ grant of 2026-10-19, "Mexer em Governed Paths declarados", whose option text reads: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica)."

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted only for a correction that requires one. This Spec's design needs none of them.

The Spec is a `type: feat`. Its Release Ceiling is `minor` because it adds one command, one check of the Doctor Command and one field of the Update Check; no command, flag or exit code is removed or renamed. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge, corrections included.

## What the grant does not cover

This record authorizes the delivery of the product change. It is not a yes to the Flow Command List for this repository. No Task runs `roundfix permissions` with a digest in this repository, and no Task creates or changes `.claude/settings.json`, `.claude/settings.local.json`, a file under `.codex/` or an OpenCode configuration file here. Whether this repository applies the list is the maintainer's own act after the release.

## How the governed set was measured

The author measured the governed set on 2026-11-19 on the planning clone of the branch `docs/plan-0413`, based on `abe2ebe1`, with `speccheck.GovernedPath` in a `go test -overlay` probe of `internal/speccheck` that wrote nothing to the repository. The probe covered every file the TechSpec names and the controls below.

- The six paths above are the governed ones this Spec changes:
  - the canonical and mirrored `SKILL.md` of the Roundfix skill, which names the new command and carries the version fields;
  - the canonical `setup` and `modes` references of the Roundfix skill;
  - `internal/speccheck/governed.go`, which gains the pattern for the harness permission files, and `internal/speccheck/governed_repocontract_test.go`, which pins that the governed set only grows.
- These governed paths are not granted, and no Task may change them:
  - `.roundfixrc.yml`: this Spec adds no Project Config key;
  - every other owned skill and every other file of the Roundfix skill, except the mirrors that `make skills-sync` writes from the granted files;
  - every Vendored Skill under `.agents/skills`;
  - everything under `internal/baseline/assets/`, and every guide under `docs/agents/`, `agent-instructions.md` and `setup-context.json` included: the tooling-authority clause keeps its text, and no Baseline module changes;
  - `docs/references/coverage-record.json`: the design removes and renames no test;
  - `.gitignore` and every other ignore file;
  - the controls: the `Makefile`, `.github/workflows/ci-verify.yml` and `go.mod`.
- From the release on, the governed set also holds the project permission files of the three harnesses. This repository tracks none of them, and this Spec creates none.

## What is not governed

- The new package of the Flow Command List under `internal/`, and under `internal/cli/` the new command file, `cli.go`, `health.go`, `doctor.go`, `update_check.go` and their tests.
- `internal/speccheck/governed_test.go` and `internal/docscontract/command_documentation_test.go`.
- The mirrors `skills/roundfix/references/setup.md` and `skills/roundfix/references/modes.md`, `skills/testdata/owned-skill-versions.json` and the skill tests.
- `CONTEXT.md`, the user guide pages `commands.md`, `commands/permissions.md`, `commands/doctor.md`, `commands/update.md`, `configuration.md` and `usage.md`, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`, ADR-0444, ADR-0445 and ADR-0446.
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

No mirror, version field or record is hand-edited.

## Limits

- No change to any of these:
  - the Makefile, the CI workflows, `go.mod` or `go.sum`, except for a correction that requires one;
  - the lint, formatter or test-runner configuration;
  - `.roundfixrc.yml`, `skills-lock.json` or any ignore file;
  - any Baseline module, template, profile, Source Baseline or Setup Snapshot, any guide under `docs/agents/` and this repository's `AGENTS.md`;
  - any harness permission file of this repository: `.claude/settings.json`, `.claude/settings.local.json`, anything under `.codex/`, `opencode.json` and `opencode.jsonc`;
  - the Author Subagent files under `.claude/agents/` and what `roundfix skills install` and `roundfix baseline update` write;
  - how a Run answers a permission request, the sandbox of a Run and Full-Access Mode;
  - any existing command, flag or exit code, and any check code of `roundfix spec check`;
  - `docs/references/coverage-record.json`;
  - any other skill, any Vendored Skill and any Archive Record;
  - the `### QA settlement` section of any skill;
  - any accepted ADR other than ADR-0444, ADR-0445 and ADR-0446, which this Spec adds.
- No change to any file of another repository, and no read of one.
- No test or Task Verification:
  - reaches the network, calls ai-memory or runs the real `gh`, `codex`, `claude` or `opencode`;
  - reads or writes the real `~/.roundfix`, the real home directory or a harness file under it;
  - starts an Agent session, an Author Subagent or an ACP adapter, or asks a person a question;
  - writes a harness permission file anywhere but a temporary repository the test created;
  - reads a file under the Spec root;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
