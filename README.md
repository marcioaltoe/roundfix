# Roundfix

Roundfix is a local-first Go CLI for resolving pull request review feedback and
executing Spec Task Graphs with local coding agents. It fetches unresolved
CodeRabbit findings, stores them as local markdown Review Issue artifacts,
assigns bounded Batches or Tasks to a local Agent runtime, verifies Agent
changes, creates Daemon-owned commits, and pushes only at configured clean
boundaries.

Roundfix is not a general workflow engine, CI healer, or task orchestration
system. The MVP focuses on one review-resolution loop for an Open Pull Request.

## Requirements

- Go 1.26 or newer.
- `make`.
- GitHub CLI `gh` authenticated for the target repository.
- Node.js 22.13 or newer with npm/npx.
- acpx `0.12.0` on `PATH`. Run `roundfix setup` after installing Roundfix to
  verify the exact version; when acpx is missing or mismatched, setup offers the
  pinned install command and runs it only after confirmation or `--yes`.
- A supported ACP Runtime selected through acpx: `codex`, `claude`, or
  `opencode`.
- `rtk` is optional. The `Makefile` uses it when available and falls back to the
  plain Go toolchain when it is not installed.

For latency-sensitive setups, configure direct adapter binaries in acpx config
so default adapters do not launch through `npx -y` on first use.

Known constraint: acpx `0.12.0` has a hard 10 MiB queue-owner per-message
buffer with no CLI, config, or environment override found in the pinned
package. Large docs-task payloads, especially turns that print or return large
skill/docs file content, can trigger `-32603 Message buffer exceeded 10485760
bytes`. Roundfix proceeds to verification when acpx delivered a parsed prompt
result before that exit; if completed work is preserved in the tree but a Task
stays failed, review the tree and run the Settle Command.

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

Bootstrap a machine for Roundfix Runs:

```bash
roundfix setup
```

Use `roundfix setup --yes` to accept every offered install or file change, or
`roundfix setup --no-input` to report missing pieces without prompting.

Update an installed Roundfix binary:

```bash
roundfix upgrade
```

Use `roundfix upgrade --check` to report the latest release outcome without
installing it. Operational commands also emit one best-effort stderr note per
day when the installed version is behind, using the shape
`roundfix 1.0.0 is behind latest 1.1.0; run roundfix upgrade`.

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

Verify and prepare this machine for Roundfix Runs:

```bash
go run ./cmd/roundfix setup
```

Upgrade Roundfix or check the release channel:

```bash
go run ./cmd/roundfix upgrade
go run ./cmd/roundfix upgrade --check
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

Execute a Spec's Task Graph:

```bash
go run ./cmd/roundfix implement --spec <slug> --agent codex
```

Settle one failed Spec Task whose completed work is already in the tree:

```bash
go run ./cmd/roundfix settle --spec <slug> --task <task_id>
```

Stop a live Run gracefully, or force-stop a dead or runaway Run:

```bash
go run ./cmd/roundfix stop <run-id>
go run ./cmd/roundfix stop --force <run-id>
```

Validate or install the shipped Roundfix agent skill:

```bash
go run ./cmd/roundfix skills check
go run ./cmd/roundfix skills install
```

By default, `skills install` writes the shipped skill to
`<repo>/.agents/skills/roundfix`. Use `--target codex`, `--target claude`,
`--target opencode`, or `--target all` for user-scoped Agent skill directories.
If the project already has `.claude/skills`, Roundfix asks whether to create
`.claude/skills/roundfix` as a symlink to the project-local skill.

Supported Agent names are `codex`, `claude`, and `opencode`. Supported Review
Source is `coderabbit`.

Preflight and run messages use color automatically in interactive terminals.
Set `ROUNDFIX_COLOR=always` to force color, `ROUNDFIX_COLOR=never` to disable
it, or set `NO_COLOR` to suppress color.

## Command Boundaries

- `setup` verifies Node.js, the pinned acpx version, the configured Agent
  probe, acpx local adapter overrides, User Config, and Project Config. Each
  check prints one deterministic report line such as `node: ok`,
  `acpx: installed`, or `User Config: skipped`. `--yes` accepts offers;
  `--no-input` skips offers instead of prompting.
- `upgrade` resolves the latest Roundfix release through the GitHub CLI.
  Successful stdout outcomes are `upgraded 1.0.0 → 1.1.0`,
  `already current 1.0.0`, `no releases published`, and, with `--check`,
  `upgrade available 1.0.0 → 1.1.0`. Failures leave the current binary
  untouched and print a manual fallback on stderr.
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
- `implement` executes a Spec's Task Graph on the current branch as one Run.
  Tasks run in dependency order, each Task's Verification commands gate one
  commit, and `implement.auto_push: true` makes a Clean spec Run push its
  branch upstream and append `pushed <remote>/<branch>` to stdout. Unresolved,
  Failed, Stopped, and failing-QA Runs never push.
- `settle` targets one failed Task, re-runs its Verification commands, changes
  nothing when verification fails, and on pass settles it `completed`, stages
  all current worktree changes plus the task file, creates the standard Task
  commit, creates no Run, writes no Run Event Journal entries, and never
  pushes.
- `stop` is graceful by default. It records a Stop Request in the Run Database
  and reports `Stop Request recorded; the Run stops after the current Work Item
  settles.` Use `--force` only for a dead, stuck, or runaway Run; it cancels
  the Agent Session best-effort, completes the Run Stopped immediately, and
  releases its Active Run locks.
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

implement:
  auto_push: false

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
make skills-install
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
skills/                          shipped Roundfix agent skill
docs/                            product docs and architecture decisions
```

Start with:

- [Product brief](docs/product-brief.md)
- [Project glossary](CONTEXT.md)
- [Architecture decisions](docs/adr/)

## License

MIT. See [LICENSE](LICENSE).
