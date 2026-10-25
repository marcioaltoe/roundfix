---
status: accepted
created_at: 2026-10-23T12:00:00Z
updated_at: 2026-10-24T17:30:00Z
deprecated_at: null
superseded_by: null
---

# A Python profile imports the published upstream Python preset

Roundfix shipped no Baseline Profile for a Python repository, so Quaestor, a
Python 3.12 project run with uv, pytest and ruff, adopted the Baseline on
2026-10-19 through a repository-owned profile that only removes the `rust`
module from a built-in one and carries no Python guidance. On 2026-10-22 the
maintainer chose a full Python profile and allowed it to request new
recommended skills from the upstream skills repository. The first version of
this decision, of 2026-10-23, derived a `python` setup from the skills the Go
and Rust lists share, because upstream had no Python list.

On 2026-10-24 upstream published one: `setups/python.txt` at the immutable
commit `32e4143e1b018005a65069afa408cec2c6e73a83` of `marcioaltoe/skills`
(pull request #130), 47 skills. The same day upstream replaced `handoff` with
six `ai-memory-*` skills in every preset at commit
`65d5f365cace27c4b468db1723213723b10c870b` (pull request #131), leaving the
Python list 52 skills. Its first 36 entries are the sequence the Go and Rust
lists share at that commit, the Baseline skills, `find-skills` and the six
ai-memory skills; then uv, ruff, ty,
pytest (vendored from Microsoft's debugpy repository), Typer, FastAPI,
Python services and file pipelines, and the backend, document and container
skills. The same day the maintainer decided: "Emendar a 0309 agora
(Recomendado)" — "Antes de a 0309 entrar na fila: importa o preset publicado
(commit 32e4143e) e o módulo python passa a despachar uv, ruff e pytest para
as skills novas; ty, typer, fastapi, pipelines e streamlit ficam só quando
escolhidos. Mantém a v0.123."

## Decision

**The `python` setup is the upstream preset, imported by name.** It is an
ordinary Setup Snapshot that mirrors `setups/python.txt`, written only by the
asset sync (ADR-0191) from a clean checkout of upstream at
`65d5f365cace27c4b468db1723213723b10c870b`. A new snapshot enters the sync as
a seed: a byte copy of the `rust` snapshot with only its `id` and
`source.path` changed, which the sync then rewrites, the path Spec 0208
used for its renamed snapshots.

**Every upstream snapshot names the same commit.** The asset sync rebuilds
every upstream snapshot from one checkout, and a repository-owned profile
takes each skill's contract from all embedded snapshots and refuses one whose
contracts disagree (ADR-0241, `restore.snapshot-conflict`). Skills the Python
preset shares with the `go`, `rust` and `typescript` lists changed bytes
upstream, so a `python` snapshot at a commit the others do not name would
make every repository-owned profile fail. On 2026-10-24 the maintainer
ordered the Baseline's move to ai-memory before the Python profile ("Antes da
0307 (Recomendado)" — "A 0317 entra logo após a 0316. A 0307 é emendada para
usar só o ai-memory (sem triagem de inbox), e a 0317 também tira a fonte
secondbrain-inbox do comando follow-up. O setup Python vai uma versão para
frente (v0.127)."): Spec 0317 moves `go`, `rust`, `typescript` and the
composed `go-cli-typescript-bun` to `65d5f36` and restores this repository's
changed skills, and the Python import then writes only the `python` snapshot
at that same commit.

**The `python` profile.** A built-in profile selects the universal, context,
Spec, autonomous-work, external-triage, Secondbrain and repository-extension
modules and a new `python` module, on the `python` setup. It selects no
application-surface module: on 2026-10-24 the maintainer named it `python`,
not `python-cli`, because it serves a command-line tool, a backend service or
a worker alike, and a repository states the rules of its own surface in its
repository-owned extension. The module's guide, `docs/agents/python.md`,
states with force that uv owns dependencies and the project environment, that
ruff's findings fail the gate and are never hidden, that a Task Verification
selects a new test by its pytest node id alone, that tests stay hermetic and
data from outside the repository skips by name, that every operating-system
branch is exercised, and that every entry point stays thin and is tested
through its observable contract. Every Python gate also checks the lockfile
with `uv lock --check --offline` and type-checks the package and its tests
with `uv run ty check --error-on-warning`: on 2026-10-24 the maintainer chose
"ty e lock check sempre" ("Todo repositório Python recebe ty no verify e `uv
lock --check`, igual ao Quaestor."). uv decides whether a lock is current by
comparing it with the project metadata, so the offline check needs neither
the network nor a cache for a current lock and fails for a stale one; CI
keeps `uv sync --locked`.

**What the module dispatches.** The module requires the nine Python skills of
the preset and dispatches each through one trigger. `uv`, `ruff`, `ty`,
`pytest` and `python-services` are dispatched whenever their work arises.
`typer`, `fastapi`, `python-file-pipelines` and `developing-with-streamlit`
are dispatched only when the repository has chosen that tool: each trigger
begins "Only when the repository", the way existing triggers already name the
tool they serve (the `backend` module's `drizzle-orm` or `better-auth`). The
preset's other skills (`backend-docs`, `external-api-adapters`,
`integration-contract-testing`, `security-best-practices`, `docker`, `pdf`,
`pdftk-server` and `agentic-cli-design`) are installed with the setup and
dispatched by no module of this profile, because the catalog lets one module
dispatch a skill and the `backend` module already dispatches four of them.

## Considered Options

- A Common Setup Snapshot, the part of `go` and `rust` both carry, was the
  first version of this decision; it was dropped on 2026-10-24 because
  upstream published a curated Python list, which a derived intersection
  cannot follow.
- Importing only the `python` snapshot at a commit the others do not name was
  rejected: shared skills would carry two contracts and ADR-0241's refusal
  would fail every repository-owned profile. Relaxing ADR-0241 to accept
  disagreeing contracts was rejected with it, because a repository profile
  would then restore a skill whose bytes depend on which snapshot it read
  first. Importing it at `32e4143e` while refreshing the others in the same
  Spec was the first version of this amendment; the batch 6 order gave that
  refresh to Spec 0317 at `65d5f36`.
- A Baseline decision per optional tool (one for Typer, FastAPI, file
  pipelines and Streamlit) was rejected: a decision can activate a module or
  select a template but cannot dispatch a skill, so it would need four new
  modules; the conditional trigger states the same choice in the dispatch
  guide an Agent already reads.
- Leaving ty to the repositories that choose it, and the lock check to CI,
  was the first version of this amendment; the maintainer rejected it on
  2026-10-24 ("ty e lock check sempre"). An online `uv lock --check` inside a
  Run was rejected with it, because the gate of a Run has no network.
- Composing `go` and `rust`, or taking the `go` setup, was rejected: either
  would restore Go, TUI or Rust skills into a Python repository.
- A `python-cli` profile with the CLI surface was rejected by the maintainer
  on 2026-10-24 ("Nenhuma superfície fixa"): one profile serves command-line
  tools, backend services and workers.

## Consequences

- Adopters of the other built-in profiles see no setup change from the
  Python import; the changed skills and the ai-memory skills reach them with
  Spec 0317's refresh, where a skill they modified is kept and reported,
  never overwritten.
- A later upstream change to `setups/python.txt` reaches the profile through
  the next asset sync, like every other preset.
- The profile itself writes no Makefile, Project Config or CI workflow; its
  guidance never grants the authority to create repository tooling, which the
  core tooling clause reserves to the maintainer.
