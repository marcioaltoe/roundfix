# The root statement: the maintainer's decision and its measurement (2026-11-18)

Recorded by the author on 2026-11-18, amending the plan of Spec 0405 before
its Run.

## The maintainer's decision

AskUserQuestion, 2026-11-18 about 02:40 BRT, verbatim:

> "Não quero que mude nunca o ~/.claude/CLAUDE.md para que a baseline do
> roundfix funcione. Tudo tem que funcionar por configurações no repositório.
> Ou seja, se for necessário a informação de que o cabeçalho só vale com
> /skill isso deve ser adicionado ao agents.md e claude.md do repositório. E
> quem deve fazer essa alteração/ajuste no baseline é o roundfix."

- Where the change goes: "Juntar na 0405 (Recommended)".
- Authorization: "Mesmas do BE (Recommended)": implement, commit, push,
  pull_request, merge and release; up to five corrections; the Makefile, CI,
  `go.mod` and lint only if a correction requires one.

## The root files' standard (about 06:30 BRT)

AskUserQuestion and the session, 2026-11-18 about 06:30 BRT, verbatim:

> "Temos que ter sempre um padrão igual para todos os projetos que usam o baseline do roundfix. Então se for necessário instruções diferentes no AGENTS.md e CLAUDE.md, vamos ter sempre dois arquivos e sem symlink já sem nem dar opção ao usuário, é regra de uso do baseline. Se não for necessário e formos manter arquivo AGENTS.md e o CLAUDE.md com symlink para o AGENTS.md, vai ser sempre assim. Sem opção ao usuário de mudar o baseline no projeto. O único caminho para personalização do baseline é o arquivo que contem regras especificas do repositorio, hoje chamado de specific-repository.md"

> "eu gosto da ideia de usar o CLAUDE.md como symlink para o AGENTS.md, porém… o trecho que fala sobre o uso da skill afeta de alguma forma … o uso das skills no Codex ou Open Code? Se não afeta, a gente vai ter dessa forma: um AGENTS.md, que é o arquivo, e o CLAUDE.md sempre o symlink… Mas se houver necessidade de separar… essa instrução, que é específica para o Claude, ficaria no CLAUDE.md."

The Supervisor answered that the statement opens with "In Claude Code," and
says it does not apply to Codex or OpenCode, whose models come from the Agent
Selection Profiles; so the instruction does not differ per harness, and the
standard is one `AGENTS.md` with `CLAUDE.md` a symbolic link to it. The rule
that makes the link mandatory → "Spec 0406 depois da 0405 (Recommended)";
authorization → "Mesmas do BE (Recommended)".

## Where the root files get the statement

Measured on a scratch worktree of `eb888135`, removed afterwards.

- The Baseline writes its root blocks into `AGENTS.md`. Every built-in profile
  (`go-cli-tui`, `go-cli-typescript-monorepo`, `python`, `rust-cli`,
  `standard-typescript-monorepo`) includes the `spec-workflow` module, whose
  root block `root.spec-workflow` (template `template.root.spec-workflow`,
  `internal/baseline/assets/templates/root/spec-workflow.md`) holds one bullet
  and no clause (`"rules": []`).
- A greenfield adoption with the built binary (`roundfix baseline plan
  --profile go-cli-tui` with every required decision, then `roundfix baseline
  apply`) in a temporary repository whose `CLAUDE.md` was a regular file
  created `AGENTS.md` with the statement and left `CLAUDE.md` without any
  managed block, writing a backup copy of it. The Baseline writes no root block
  into a regular `CLAUDE.md` today; Spec 0406 makes `CLAUDE.md` always a link to `AGENTS.md`.
- Every adopter checkout under `~/dev` that holds a Setup Manifest (argus,
  quaestor, roundfix) keeps `CLAUDE.md` as a symbolic link to `AGENTS.md`
  (`ls -l`, read only). In each, the statement in `AGENTS.md` is the statement
  `CLAUDE.md` shows.

## The prototype

The prototype added one bullet to the template (the wording before the 06:30 decision, without the Claude Code scope; the scope changes no file the regeneration touches), raised the template's version
(1 to 2) and the root block's version (2 to 3), then ran, in this order:

1. `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`, which raised `spec-workflow` from 26 to 27 and recorded it in `internal/baseline/module-versions.json`;
2. `make baseline-digests`;
3. `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`.

Files the three commands and the edit changed:

- the edit: `internal/baseline/assets/templates/root/spec-workflow.md`,
  `internal/baseline/assets/templates/index.json`,
  `internal/baseline/assets/modules/spec-workflow.json`;
- the module record: `internal/baseline/module-versions.json` and the
  module's version line;
- `make baseline-digests`: the Standard TypeScript formatter golden
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/AGENTS.md`,
  its pin in `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/testdata/catalog.diagnostics.golden.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json` and four plan goldens
  under `internal/baseline/testdata/plan-characterization/`
  (`advisory-only-divergences`, `clean-adoption`,
  `idempotent-replan-after-verified-apply`,
  `same-baseline-changed-profile-and-catalog-digests`);
- the Managed Refresh: `AGENTS.md` (the new bullet) and
  `docs/agents/setup-context.json` (the catalog digest and the root block's
  digest).

Unchanged: the Source Baseline corpus and manifests, every Setup Snapshot,
the other four profiles, every guide under `docs/agents/` and the parity
corpus.

## Governed Paths

Measured with `speccheck.GovernedPath` through a `go test -overlay` probe in
`internal/speccheck` that wrote nothing to the repository.

- **Governed and changed:** `internal/baseline/assets/templates/root/spec-workflow.md`,
  `internal/baseline/assets/templates/index.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/AGENTS.md`
  and `docs/agents/setup-context.json`.
- **Governed and not changed:** the Source Baseline corpus `AGENTS.md` and
  `internal/baseline/assets/source-baselines/index.json`,
  `internal/baseline/assets/setups/go.json`,
  `internal/baseline/assets/profiles/go-cli-tui.json`,
  `internal/baseline/assets/modules/core.json`, `docs/agents/spec-routing.md`
  and `internal/baseline/plan_test.go`.
- **Not governed:** `AGENTS.md`, `CLAUDE.md`,
  `internal/baseline/module-versions.json`, the catalog snapshots and plan
  goldens under `internal/baseline/testdata/`,
  `internal/cli/baseline_update_test.go`, a new test file under
  `internal/baseline/`, ADR-0431 and this Spec's files.

## Existing tests

With the prototype, these packages passed with `-count=1`, so the statement
moves no existing test:

- `./internal/baseline` (1,820 tests);
- `./internal/cli` (235 s);
- `go test -tags docscontract ./internal/docscontract`;
- `./skills`, `./internal/speccheck`, `./internal/spec`;
- `./internal/baselineacp`, the other package that imports `internal/baseline`.

Specs 0401 to 0404 change no Baseline module, template or profile, so the
Task's base moves none of these files.
