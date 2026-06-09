# Roundfix

Roundfix is a local-first Go CLI for resolving pull request review feedback with
local coding agents. The planned daemon fetches unresolved CodeRabbit findings,
assigns bounded batches to an Agent, verifies the changes, commits each
successful batch, and pushes only after no Unresolved Review Issues remain.

Current status: early scaffold. The CLI prints help and version output. The
operational commands `fetch`, `resolve`, and `watch` are reserved but not
implemented yet.

## What Roundfix is for

Roundfix targets the review loop that starts after a pull request is open:

1. Validate the local repo, pull request, config, artifact directory, selected
   Agent, and GitHub access before any long-running work starts.
2. Fetch unresolved Review Source issues for the current pull request.
3. Persist review issues as markdown artifacts.
4. Start the selected local ACP Runtime, such as Codex, Claude Code, or
   OpenCode, through the user's existing local setup.
5. Verify each successful batch.
6. Create one local commit per successful batch.
7. Run the Final Push only when no Unresolved Review Issues remain.

Roundfix is not a general workflow engine, CI healer, or task orchestration
system.

## Requirements

- Go 1.26 or newer.
- `make`.
- `rtk` is optional. The `Makefile` uses it when available and falls back to the
  plain Go toolchain when it is not installed.

## Build

```bash
make build
```

The binary is written to:

```text
bin/roundfix
```

Remove build artifacts with:

```bash
make clean
```

## Usage

Show help:

```bash
make run ARGS=--help
```

Print the version:

```bash
make version
```

Current CLI output:

```text
Usage:
  roundfix --help
  roundfix --version
  roundfix fetch
  roundfix resolve
  roundfix watch
```

The planned command contract is:

```bash
roundfix fetch --source coderabbit --pr <number>
roundfix resolve --pr <number> --agent codex
roundfix watch --source coderabbit --pr <number> --agent codex --until-clean
```

These commands are documented in the product brief, but they do not run the
review loop yet.

## Development

Run the local verification gate:

```bash
make verify
```

Useful targets:

```bash
make fmt
make fmt-check
make test
make test-race
make build
make deps
```

The current `verify` target runs formatting checks, tests, and build:

```text
fmt-check -> test -> build
```

## Project structure

```text
cmd/roundfix/       CLI entry point
internal/app/       application metadata
internal/cli/       command parsing, stdout/stderr, and exit codes
docs/               product docs and architecture decisions
```

Start with:

- [Product brief](docs/product-brief.md)
- [Project glossary](CONTEXT.md)
- [Architecture decisions](docs/adr/)

## Planned product shape

Roundfix has three operational commands:

- `fetch`: download unresolved Review Source issues for an Open Pull Request and
  persist markdown artifacts.
- `resolve`: process downloaded Unresolved Review Issues for an Open Pull
  Request with the selected Agent.
- `watch`: automate `fetch` and `resolve` across Review Source rounds until the
  pull request is clean, `max_rounds` is reached, or a Run Budget stops the run.

The daemon owns Review Source access, Run state, verification, commit policy,
and Final Push. The child Agent owns only triage, code edits, tests, and issue
artifact updates for its assigned Batch.

## License

MIT. See [LICENSE](LICENSE).
