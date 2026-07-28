# Roundfix

Roundfix is a local-first Go CLI that resolves pull request review feedback and
executes Spec Task Graphs with local coding agents. It fetches unresolved
CodeRabbit findings as local markdown artifacts, assigns bounded Batches or
Tasks to an ACP Agent runtime, verifies every change with your repository's own
gate, creates Daemon-owned commits, and pushes only when nothing unresolved
remains.

It is not a workflow engine, CI healer, or task orchestrator. It runs two
loops — review resolution on an Open Pull Request, and Spec Task Graph
execution — plus the support commands that keep a machine Run-ready (`setup`,
`doctor`, `upgrade`, `gc`, `settle`, `archive`, `stop`).

The spec pipeline it runs is CONTEXT-driven development, adapted from Matt
Pocock's **skills** work — see
[source and attribution](docs/user-guide/context-driven-development.md#source-and-attribution).

## Install

Roundfix ships through npm as a `roundfix` launcher with per-platform binary
packages (darwin/linux `x64`+`arm64`, windows `x64`):

```bash
npx roundfix --version        # run once without installing
npm install -g roundfix       # or: bun add -g roundfix
```

Then make the machine Run-ready and check it:

```bash
roundfix setup     # proves adapters and generated profiles before writes
roundfix doctor    # read-only profile and Repository Skill Set readiness
```

## Requirements

- Node.js 22.13+ with npm/npx (the ACP Agent layer requires acpx `0.12.0` or
  newer; `setup` installs the minimum tested version only when acpx is missing
  or older).
- A supported ACP Runtime: `codex`, `claude`, or `opencode`.
- GitHub CLI `gh` authenticated for the target repository (review loops only).
- Building from source additionally needs Go 1.26+ and `make`.

`roundfix doctor` diagnoses all of it — including official Codex adapter
identity, exact Agent Selection Profile proof, Repository Skill Set readiness,
and macOS codex hygiene — with one line per check and a `next:` action on
failures. The independent Repository Skill Set result follows `profiles:`:

```text
profiles: ok (3 distinct tuples; 10 category references)
skills: ok (39 required: 14 Roundfix-owned, 25 external)
```

The running binary's embedded bundle is authoritative for Roundfix-owned
skills, including the Roundfix Skill. Each required external skill must match
its `computedHash` in the repository's `skills-lock.json`. Missing, outdated,
or invalid required skill state prints `skills: failed (...)`, still runs the
other Doctor checks, and makes Doctor exit `1`.

Doctor checks these authorities locally, without network access or writes. It
ignores unrelated extra skill directories and lock entries, never deletes or
updates skills automatically, and only prints the applicable command:

```bash
roundfix skills install --target project
bunx skills experimental_install && bunx skills update -p -y
```

Use `roundfix setup` to provision the minimum supported acpx or official
adapter, or migrate a stale legacy override after authorization.

## Quickstart

Resolve a pull request's CodeRabbit feedback until it is merge-ready:

```bash
roundfix watch --source coderabbit --pr <number> --until-clean
```

Execute a Spec's Task Graph (specs live under `docs/specs/<slug>/`), with the
QA gate at the end:

```bash
roundfix implement --spec <slug> --qa
```

Adopt or update the repository's Context-Driven Baseline through the public,
confirmation-gated workflow:

```bash
roundfix baseline
```

For scripts, CI, or agents, detach the Run and monitor it without owning it:

```bash
roundfix implement --spec <slug> --detach
roundfix events <run-id> --follow --filter outcome  # terminal outcome for Supervisors
roundfix attach <run-id>                             # read-only Live Run View for humans
```

Both loops keep one contract: stdout carries only the deterministic report,
diagnostics go to stderr, and exit codes are stable (`0` clean, `1`
unresolved/failed, `2` preflight, `3` Clean Unverified or Review Skipped watch,
`130` Ctrl-C) — which is what lets an agent drive Roundfix reliably.

## Documentation

- [Operational guide](docs/user-guide/usage.md) — the two loops end to end:
  run, monitor, read outcomes, recover failed Tasks.
- [Command reference](docs/user-guide/commands.md) — every command's flags,
  output shapes, exit codes, and the boundaries it never crosses.
- [Configuration](docs/user-guide/configuration.md) — config precedence, agent
  and model selection, worktrees and bootstrap, retention, local state paths.
- [CONTEXT-driven development](docs/user-guide/context-driven-development.md)
  — the method behind the spec pipeline and the complete Baseline adoption,
  automation, migration, recovery, and security contract.
- [Project glossary](CONTEXT.md) and
  [architecture decisions](docs/adr/) — the vocabulary and decision log.
- [Release runbook](docs/user-guide/release-runbook.md).

## Build from source

```bash
make build        # writes bin/roundfix
make install      # installs into your Go bin directory
```

The built CLI is identical to the npm-launched binary — same stdout, stderr,
and exit codes.

## Development

```bash
make verify       # fmt-check -> test -> skills-sync-check -> skills-check -> build
```

Project layout:

```text
cmd/roundfix/                     CLI entry point
internal/agent/                   ACP Agent runtime execution
internal/cli/                     command parsing, output, and exit codes
internal/config/                  YAML config loading and validation
internal/daemon/                  verification, commits, source resolution, push
internal/preflight/               git, PR, worktree, and push safety checks
internal/reviewsource/coderabbit/ CodeRabbit implementation
internal/rounds/                  Round artifacts, issue parsing, batching
internal/store/                   central Run Database
internal/tui/                     Interactive Input and Live Run View
internal/watch/                   watch state machine
skills/                           shipped skill bundle (synced from .agents/skills)
docs/                             product docs and architecture decisions
```

## License

MIT. See [LICENSE](LICENSE).
