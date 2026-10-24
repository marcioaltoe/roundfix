---
spec: 0310-the-python-cli-profile-writes-its-starter-files
status: active
created: 2026-10-23
surfaces: [backend, cli, docs]
---

# The Python CLI profile writes its Starter Files

A maintainer starting a Python repository on the `python-cli` Baseline Profile
still writes by hand the Makefile that declares `make verify`, the Project
Config that bootstraps each worktree with uv, and the CI workflow, and the
profile's alignment refuses to plan until the Makefile exists. After this Spec
the profile's module offers three Starter Files in the shape Quaestor uses, each a **Starter File** that Roundfix writes once and never manages,
and a new confirmation-gated command, `roundfix baseline starters`, writes the
ones whose paths do not exist after the maintainer confirms the exact preview,
never touching a file the repository already has.

## Sources

This Spec is a secondary Spec of one source. Spec 0309, which delivers the
`python-cli` profile, owns the source and adopted it as a promoted Backlog
Entry; this Spec links the owner's copy and adopts nothing:

- [`../0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff/references/2026-10-22-a-baseline-profile-for-python-repositories.md`](../0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff/references/2026-10-22-a-baseline-profile-for-python-repositories.md),
  type `feat`: Quaestor's request of 2026-10-21 and the maintainer's
  decisions of 2026-10-22. This Spec delivers its expected item 4: generating
  the Makefile, the Project Config and the CI workflow.

At the authoring base `dde9d5f5` `docs/backlog/` and `docs/findings/` held no
entry; nothing else shares this context.

## Prerequisites

`requires: [0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff]`: this
Spec's Starter Files belong to the `python` module that Spec 0309 creates, and
its Verification commands run tests over that module and the `python-cli`
profile. The operator delivers it right after Spec 0309. Spec 0307 changes the
Roundfix skill's version and the Baseline digests before it; this Spec's
Tasks raise the skill version over what the tree holds and regenerate every
derived file through the sanctioned commands, naming no version or digest
value.

## Project Constraints

- Identifier strategy: not applicable — no project-owned identifier or record
  schema is minted; catalog identifiers follow the dotted naming and the JSON
  schema name follows the `roundfix/baseline-*` payloads. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files only; tests use
  temporary repositories and homes. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0323, "A profile Starter File is
  written only where absent and only on a confirmed preview", which this Spec
  records; ADR-0068, "Baseline Command uses one confirmation-gated workflow",
  whose digest confirmation the new command follows ("apply only that document
  and its exact approved digest"); ADR-0075, "Profile divergence uses
  confirmed repository-owned adaptation", whose divergence next action this
  Spec extends for a profile with Starter Files; ADR-0219, "A profile draft
  binds to the closest built-in profile it adapts", and ADR-0241, "A
  repository profile takes its skill contracts from the embedded Setup
  Snapshots", under which a repository-owned draft of `python-cli` keeps the
  module's Starter Files; and ADR-0322, "A Python profile takes the
  skills the Go and Rust setups share", whose module this Spec extends.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 through AskUserQuestion, scope "Preset de baseline
  Python", authorization "same as cycle AO", which carries the cycle AJ grant
  "Mexer em Governed Paths declarados" ("'Concedo' para os caminhos que cada
  Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
  Makefile, CI, go.mod, lint, que pedem autorização específica)."), and the
  preset's reach, the maintainer's own answer of the same day. The Makefile,
  Project Config and workflow this Spec adds are Baseline template assets
  rendered into adopters, not this repository's own tooling. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0310-the-python-cli-profile-writes-its-starter-files/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `docs/agents/setup-context.json`,
  `internal/baseline/assets/modules/python.json`,
  `internal/baseline/assets/templates/index.json`,
  `internal/baseline/assets/templates/starters/python/Makefile`,
  `internal/baseline/assets/templates/starters/python/roundfixrc.yml`,
  `internal/baseline/assets/templates/starters/python/verify.yml`,
  `internal/cli/baseline_documentation_contract_test.go`,
  `internal/docscontract/publicdocs_test.go`,
  `skills/roundfix/SKILL.md`. No change to this repository's Makefile,
  workflows, Go module, lint, formatter or Project Config.

## Goals

- A maintainer previews, for a profile, which Starter Files would be created
  and which kept, and creates the absent ones only by confirming that exact
  preview.
- The `python-cli` profile offers a Makefile whose `verify` runs ruff's format
  check, ruff's linter and pytest through uv, a Project Config that bootstraps
  each worktree with a frozen uv sync, and a CI workflow that refuses a stale
  lockfile and runs `make verify`.
- A file that exists is never overwritten, compared or removed, and a failed
  write leaves no partial set.
- Profile alignment tells a repository without a declared gate how to get one
  when its profile has Starter Files.

## User Stories

1. As a maintainer starting a Python repository, I want Roundfix to write the
   Makefile, Project Config and CI workflow the profile expects, so that my
   first plan is ready without copying another repository's files.
2. As a maintainer of an existing repository, I want the command to keep every
   file I already have, so that my own tooling decisions stay mine.
3. As an Agent, I want the command to write nothing without the maintainer's
   confirmed digest, so that the tooling clause holds.
4. As a maintainer whose plan is blocked by an undeclared gate, I want the
   next action to name the command, so that I find the way to a working gate.

## Core Features

1. Starter File declarations: a module field `starterFiles` and a template
   kind `starter-file`, validated when the catalog loads.
2. The `python` module's three Starter Files: `Makefile`, `.roundfixrc.yml`
   and `.github/workflows/verify.yml`.
3. The preview and the confirmed write: `create` only where nothing exists at
   the path, `keep` otherwise, a Plan Digest, exclusive creation, refusal of a
   digest that no longer matches, and removal of what was created when a
   later write fails.
4. Profile alignment's next action for an undeclared Verification command
   names the command when the profile has Starter Files.
5. The command `roundfix baseline starters --profile <id> [--confirm-plan
   <digest>] [--repo <path>] [--format <text|json>]`, its help, guide, Roundfix
   skill reference, Skill Coverage Map entry, and the glossary term **Starter
   File**.

## User Experience

```text
$ roundfix baseline starters --profile python-cli
Baseline starters: python-cli
- create Makefile
- create .roundfixrc.yml
- create .github/workflows/verify.yml
State: preview
Plan Digest: sha256:<digest>
Next: roundfix baseline starters --profile python-cli --confirm-plan sha256:<digest>
```

The preview exits 3. Running the printed next command writes the three files,
prints the same lines with `State: applied` and no `Next:` line, and exits 0.
A later preview lists each as `keep`, prints `State: current` and exits 0.

## Non-Goals / Out of Scope

- Starter Files for any profile other than `python-cli`.
- Managing, refreshing, comparing or removing a Starter File after it is
  written.
- A `pyproject.toml`, `uv.lock` or source tree starter: the repository's
  package is its own.
- Running uv, make, pytest, ruff or a CI workflow from Roundfix or its tests.
- Any change to this repository's own Makefile, workflows or Project Config.

## Success Metrics

1. Success Metric: the embedded catalog returns, for `python-cli`, the Starter
   Files `Makefile`, `.roundfixrc.yml` and `.github/workflows/verify.yml` in
   that order with `sha256:` content identities, and none for the other four
   profiles; a declaration with an unsafe or duplicate path, an unknown
   template or a template of another kind is refused with its own code.
2. Success Metric: the Makefile's `verify` runs `fmt-check`, `lint` and `test`
   through `uv run --frozen` and `verify-incremental` runs `verify`; the
   Project Config written by the command loads through Roundfix's Project
   Config loader with `make verify`, the frozen bootstrap, a five-minute
   bootstrap timeout, the three caches and the tools `uv`, `make` and `git`;
   the workflow runs `uv sync --group dev --locked` then `make verify` on pull
   requests and pushes to `main`.
3. Success Metric: a preview writes nothing; a confirmed apply creates only
   absent paths and leaves an existing file byte-identical; a digest that no
   longer matches is refused with nothing written; a write failure removes
   every file the run created; a second preview is `current` and exits 0.
4. Success Metric: for `python-cli` without a declared gate, the
   `verification.command.undeclared` next action names `roundfix baseline
   starters --profile python-cli`; for `rust-cli` it keeps its bytes.
5. Success Metric: the command's preview and applied transcripts match line by
   line, its help, guide and Roundfix reference describe it, the skill's
   version is raised and recorded, and the glossary defines **Starter File**.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- Quaestor's tooling at `marcioaltoe/quaestor@e6b390d`, read through `gh api`
  at authoring without cloning: its `Makefile` ("verify: fmt-check lint test",
  "verify-incremental: verify"), `.roundfixrc.yml` ("bootstrap: uv sync
  --group dev --frozen", "bootstrap_timeout: 5m", "caches: [.venv,
  .pytest_cache, .ruff_cache]") and `.github/workflows/verify.yml` ("uses:
  astral-sh/setup-uv@v10.3.0", "run: make verify"). The QA compares the
  starters with these quoted lines; it never reads that repository.
- uv's documentation, "Locking and syncing"
  (https://docs.astral.sh/uv/concepts/projects/sync/), quoted: "To use the
  lockfile without checking if it is up-to-date, use the `--frozen` option"
  and "If the lockfile is not up-to-date, uv will raise an error instead of
  updating the lockfile." — why the worktree bootstrap and the Makefile use
  `--frozen` and the CI workflow uses `--locked`.

## Unreachable Acceptance

- criterion: a real Python repository runs the written Makefile with uv, ruff and pytest, and GitHub Actions runs the written workflow green
  reason: it needs uv, Python tooling, the network and a GitHub repository, which a Run sandbox and this Spec's authorization exclude
  satisfied-by: the first adopter that writes the Starter Files after the release, which the operator announces to `default/marcioaltoe-quaestor`, recorded as this Spec's unproven step

## Validation

The decisions were checked on 2026-10-23 against these sources.

- Characterization at `dde9d5f5`: the Baseline writes only managed guidance
  and the Setup Manifest; the tooling clause of
  `docs/agents/agent-instructions.md` forbids creating repository tooling
  without express authorization and says a profile does not grant it; the
  template validator does not refuse an unknown template kind; `roundfix
  baseline skills restore` is the model of a separate confirmation-gated
  Baseline command (non-empty preview exits 3, `--confirm-plan` applies the
  exact preview). A prototype on top of Spec 0309's prototype added the three
  templates, the module field, the catalog validation, the resolver and the
  confirmed write, ran the module record command, `make baseline-digests` and
  the Managed Refresh, passed the `internal/baseline` suite, and failed each
  named test under the sabotages the TechSpec lists.
- Secondbrain, read only. Read `wiki/index.md` and ran `qmd query "Python
  Baseline profile uv pytest ruff Quaestor" --all --files --min-score 0.3`,
  which returned Quaestor's decision record
  `projects/quaestor/mirror/docs/adr/0004-roundfix-e-baseline.md`: "Configuração
  do projeto em `.roundfixrc.yml`: cada worktree de Run ou Task roda `uv sync
  --group dev --frozen` antes do agente, e `.venv`, `.pytest_cache` e
  `.ruff_cache` são caches descartáveis. O CI (`.github/workflows/verify.yml`)
  roda `make verify` numa execução limpa a cada pull request e push na
  `main`." — the shape the starters reproduce.
- Exa: uv's CLI reference (https://docs.astral.sh/uv/reference/cli/) and
  locking page above. Challenging: Quaestor runs `uv run` without a flag in
  its Makefile, which may update `uv.lock` during the gate and may need the
  network to resolve; the starters use `uv run --frozen` in the Makefile (no
  lock rewrite and no resolution inside a Run) and `uv sync --group dev
  --locked` in CI, which has the network and must refuse a stale lockfile.
  Quaestor's `verification.tools` also lists `rtk`, a tool of its host; the
  starter lists only `uv`, `make` and `git`.

## Glossary

- adds: **Starter File**

## Decisions

- 2026-10-22: Scope of cycle AP (AskUserQuestion, questions of about 15:25 BRT
  answered by about 18:30 BRT) — "Preset de baseline Python", delivered last.
- 2026-10-22: Reach of the preset — "Perfil completo e pode, se necessário,
  solicitar ou requerer novas skills recomendadas ao repo skills para que
  tenhamos a curadoria e ter o melhor resultado no uso do python, uv e linters
  e outros toolings especificos do setup." (the maintainer's own answer).

Autonomous route: no person present; decisions taken from the operator's
cycle AP batch 3 briefing and the characterization above, with no further
interview:

- A separate confirmation-gated command writes Starter Files only where
  absent, and the repository owns them afterwards; see ADR-0323.
- The Makefile uses `uv run --frozen`, the Project Config bootstraps with
  `--frozen`, and CI syncs with `--locked`.
- The command's exit codes follow the Baseline commands: 0 for a current or
  applied result, 3 for a preview or a stale confirmation, 2 for invalid
  input, 1 for a write failure.

## Open Questions

None.
