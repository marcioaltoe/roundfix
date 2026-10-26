---
type: feat
status: promoted
created: 2026-10-22
spec: 0309-a-python-baseline-profile-for-uv-pytest-and-ruff
reason: null
---

# A Baseline Profile for Python repositories

## Opportunity

Quaestor, a Python 3.12 command-line tool of the maintainer, adopted the
Context-Driven Baseline on 2026-10-19 through a repository-owned profile
derived from `rust-cli` without the `rust` module, because the embedded
catalog has profiles only for Go, Rust and TypeScript. Its agents therefore
read no Baseline rule about uv, ruff or pytest; those rules live in its own
`docs/agents/specific-repository.md`, and its tooling files were written by
hand. In an ai-memory message to `default/marcioaltoe-roundfix` of
2026-10-21 at 19:02 BRT ("Pedido de Marcio"), Quaestor asked Roundfix to ship
a Baseline Profile for Python repositories in the shape Quaestor uses today
(read at `marcioaltoe/quaestor@e6b390d`):

- tools: uv, Python 3.12, the hatchling build backend, pytest and ruff, with
  the development dependency group `dev`;
- gate: `make verify`, whose targets are `fmt-check` (`uv run ruff format
  --check src tests`), `lint` (`uv run ruff check src tests`) and `test`
  (`uv run pytest -q`), and `make verify-incremental`, equal to `verify`;
- Project Config: `defaults.verification: make verify`, `worktree.bootstrap:
  uv sync --group dev --frozen` (under 2 s, measured on 2026-10-19),
  `worktree.bootstrap_timeout: 5m`, `worktree.caches: [.venv, .pytest_cache,
  .ruff_cache]` and `verification.tools: [uv, make, git, rtk]`;
- CI: `.github/workflows/verify.yml` runs `make verify` on pull requests and
  pushes to `main`;
- tests that read real client data skip when the data is not reachable, so the
  gate stays hermetic.

## Value

A Python repository adopts the Baseline with a built-in profile instead of a
hand-made one, its agents read the uv, ruff and pytest rules as managed
guidance that later refreshes keep current, the skills it restores are the
language-neutral ones it needs and no Go or Rust skill, and a new repository
starts with a working gate instead of writing it by hand.

## Shape

Expected, as stated by the request and the maintainer's decisions of
2026-10-22 ("Preset de baseline Python" in cycle AP, delivered last; "Perfil
completo e pode, se necessário, solicitar ou requerer novas skills
recomendadas ao repo skills para que tenhamos a curadoria e ter o melhor
resultado no uso do python, uv e linters e outros toolings especificos do
setup."):

1. A built-in `python` Baseline Profile with a `python` module whose guide,
   `docs/agents/python.md`, states the uv and pytest commands and rules.
2. The Python equivalent of the Go guide's cross-platform build rule: how a
   change that branches on the operating system is exercised.
3. The pytest Verification rule: a Task Verification selects a new test by its
   node id alone, the rule `roundfix spec check` already enforces through its
   overlapping-selector finding.
4. Generating the Makefile, the Project Config and the CI workflow in the
   shape above.
5. The profile works with the skills that exist today; the recommended Python
   skills the upstream skills repository lacks are requested through a
   Secondbrain Inbox Entry under `inbox/skills/`, filed by the operator.

On 2026-10-24 the maintainer named the profile `python` instead of
`python-cli` ("mude o nome do profile de python-cli para python (vai ser para
cli, backend e workers)") and chose "Nenhuma superfície fixa (Recomendado)":
the profile carries the language and the workflow modules only, serves a
command-line tool, a backend service or a worker alike, and leaves the rules
of a specific surface to the repository's own extension.
