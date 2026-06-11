# Roundfix

Roundfix is a local-first Go CLI for resolving pull request review feedback with
local coding agents. It fetches unresolved CodeRabbit findings, stores them as
local markdown Review Issue artifacts, assigns bounded Batches to a local Agent
runtime, verifies Agent changes, creates Batch commits, and runs the Final Push
only after no Unresolved Review Issues remain.

Roundfix is not a general workflow engine, CI healer, or task orchestration
system. The MVP focuses on one review-resolution loop for an Open Pull Request.

## Requirements

- Go 1.26 or newer.
- `make`.
- GitHub CLI `gh` authenticated for the target repository.
- A local Agent runtime command for the selected Agent:
  - `codex-acp` for Codex, with `npx --yes @zed-industries/codex-acp` as a fallback
  - `claude-agent-acp`, with `npx --yes @agentclientprotocol/claude-agent-acp` as a fallback
  - `opencode acp`
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

Install the CLI into your Go bin directory for local testing:

```bash
make install
```

Make sure your Go bin directory, usually `~/go/bin`, is on `PATH` before
running `roundfix` directly.

## GitHub Access

Roundfix uses the GitHub CLI (`gh`) from the local machine. It does not ask for
or store a GitHub token directly. Authenticate `gh` for the target repository
before running operational commands:

```bash
gh auth status
```

## Commands

Show help:

```bash
go run ./cmd/roundfix --help
```

Show version:

```bash
go run ./cmd/roundfix --version
go run ./cmd/roundfix -v
```

Create a Project Config in the current repository:

```bash
go run ./cmd/roundfix init
```

Create a User Config instead:

```bash
go run ./cmd/roundfix init --scope user
```

Fetch unresolved CodeRabbit Review Issues into local Round artifacts:

```bash
go run ./cmd/roundfix fetch --source coderabbit --pr <number>
```

Resolve downloaded Compatible Artifacts with a selected Agent:

```bash
go run ./cmd/roundfix resolve --pr <number> --agent codex
```

Run the watched review-resolution loop:

```bash
go run ./cmd/roundfix watch --source coderabbit --pr <number> --agent codex --until-clean
```

Validate or install the shipped Roundfix agent skills:

```bash
go run ./cmd/roundfix skills check
go run ./cmd/roundfix skills install --target codex
```

Supported Agent names are `codex`, `claude`, and `opencode`. Supported Review
Source is `coderabbit`.

Preflight and run messages use color automatically in interactive terminals.
Set `ROUNDFIX_COLOR=always` to force color, `ROUNDFIX_COLOR=never` to disable
it, or set `NO_COLOR` to suppress color.

## Command Boundaries

- `fetch` validates local state, creates a Fetch Run, fetches unresolved
  CodeRabbit review threads, writes markdown Round artifacts, and stops at the
  `Fetched` terminal outcome. It never starts an Agent, commits, pushes, or
  resolves Review Source threads.
- `resolve` works only over downloaded Compatible Artifacts. It does not fetch
  Review Source issues. It assigns a bounded Batch, runs the selected Agent
  runtime, verifies terminal assigned issues, commits successful Batches when
  auto-commit is enabled, resolves source threads for `resolved` and `invalid`
  assigned issues, and runs Final Push only when no Unresolved Review Issues
  remain.
- `watch` waits for CodeRabbit status on the current PR HEAD, observes the
  configured quiet period, fetches unresolved issues, resolves Batches, and
  repeats until `Clean`, `MaxRoundsReached`, `BudgetExceeded`, `TimedOut`,
  `Failed`, or `Stopped`.
- Agents own only assigned issue files, triage, code edits, tests,
  verification commands, and assigned Review Issue status updates. They must
  not commit, push, resolve Review Source threads, edit unassigned issue files,
  or mark issues as `duplicated`.

## Config

Roundfix reads YAML config in this order:

1. Built-in defaults.
2. User Config at `~/.roundfix/config.yml`.
3. Project Config at `<repo>/.roundfixrc.yml`.
4. CLI flags.

Use `roundfix init` to create config. When `--scope` is omitted, Roundfix asks
where to write the file and defaults to Project Config when you press Enter.
Use `--force` to overwrite an existing config file.

Example:

```yaml
defaults:
  agent: codex
  verification: make verify
  artifact_dir: .roundfix
  auto_commit: true

review_source:
  name: coderabbit
  include_nitpicks: true

watch:
  until_clean: true
  max_rounds: 6
  poll_interval: 30s
  review_timeout: 30m
  quiet_period: 30s
  auto_push: true

budget:
  enabled: true
  max_run_duration: 2h

resolve:
  batch_size: 3
  concurrent: 1
```

## Local State

- Run Database: `~/.roundfix/roundfix.db`
- Default Artifact Directory: `<repo>/.roundfix/`
- Review Issue artifacts:
  `<artifact-dir>/reviews/pr-<number>/round-<nnn>/issue_<nnn>.md`
- Agent logs:
  `<artifact-dir>/runs/<run-id>/agent/batch-<nnn>.log`

With automatic Round selection, `fetch` reuses an existing matching Round when
the same HEAD already has the same Review Issue fingerprints. If the fetched
payload is new, Roundfix writes the next Round directory. Roundfix does not
overwrite existing Round artifacts. Repeated findings across different payloads
are deduplicated later during `resolve` by Review Issue Fingerprint, preferring
`source_ref` such as `thread:<id>,comment:<id>` and falling back to
`review_hash`.

Roundfix rejects dirty worktree changes outside the Artifact Directory before
starting operational work. `fetch` allows a local Project Config change at
`.roundfixrc.yml` because it never starts an Agent, commits, or pushes.
`resolve` and `watch` also allow `.roundfixrc.yml`, but Batch commits exclude it
so local setup changes do not mix with review fixes. Terminal Run outcomes
release the Active Run lock for the PR Head Branch.

The current CodeRabbit fetch imports unresolved inline review threads.
CodeRabbit review-body summaries and outside-diff comments are not converted
into Review Issue artifacts yet.

## Development

Run the local verification gate:

```bash
make verify
```

The `verify` target runs:

```text
fmt-check -> test -> skills-check -> build
```

Useful targets:

```bash
make fmt
make fmt-check
make test
make test-race
make build
make install
make deps
make skills-check
make skills-install TARGET=codex
```

## Project Structure

```text
cmd/roundfix/                    CLI entry point
internal/agent/                  ACP Agent runtime execution
internal/cli/                    command parsing, output, and exit codes
internal/config/                 YAML config loading and validation
internal/daemon/                 verification, commits, source resolution, push
internal/preflight/              git, PR, worktree, and push safety checks
internal/reviewsource/           Review Source boundary
internal/reviewsource/coderabbit/ CodeRabbit implementation
internal/tui/                    Interactive Input and ACP Live Run View
internal/rounds/                 Round artifacts, issue parsing, batching
internal/store/                  central Run Database
internal/watch/                  watch state machine
skills/                          shipped Roundfix agent skills
docs/                            product docs and architecture decisions
```

Start with:

- [Product brief](docs/product-brief.md)
- [Project glossary](CONTEXT.md)
- [Architecture decisions](docs/adr/)

## License

MIT. See [LICENSE](LICENSE).
