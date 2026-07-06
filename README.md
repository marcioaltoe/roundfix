# Roundfix

Roundfix is a local-first Go CLI for resolving pull request review feedback and
executing Spec Task Graphs with local coding agents. It fetches unresolved
CodeRabbit findings, stores them as local markdown Review Issue artifacts,
assigns bounded Batches or Tasks to a local Agent runtime, verifies Agent
changes, creates Daemon-owned commits, and pushes only at configured clean
boundaries.

Roundfix is not a general workflow engine, CI healer, or task orchestration
system. It runs two local-first loops — resolving CodeRabbit review feedback on
an Open Pull Request, and executing a Spec's Task Graph — plus first-class
supporting commands (`setup`, `doctor`, `upgrade`, `gc`, `settle`, `archive`,
`stop`) that keep a machine Run-ready and manage Run lifecycle and storage.

## Install

Roundfix ships through npm as a `roundfix` launcher package with a per-platform
binary package for each target (darwin, linux, and windows on both `x64` and
`arm64`, except windows which is `x64` only). The launcher runs the binary for
your platform, so `npx`, `bunx`, and global installs all behave identically to a
locally built binary — same stdout, stderr, and exit codes.

Run it once without installing:

```bash
npx roundfix --version
bunx roundfix --version
```

Install it globally to put `roundfix` on `PATH`:

```bash
npm install -g roundfix
# or: bun add -g roundfix
```

After installing, make the machine Run-ready with `roundfix setup` and verify
readiness with `roundfix doctor`. Node.js is a hard prerequisite either way,
since the ACP Agent layer runs through acpx. Building from source
(`make build`, below) stays supported and produces the same CLI.

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

Run `roundfix doctor` at any time to diagnose Run readiness without installing,
writing config, or otherwise mutating the machine. Doctor checks Node.js, the
pinned acpx version, the configured Agent probe, and codex runtime hygiene. On
macOS, codex hygiene resolves `CODEX_PATH` first and then `codex` on `PATH`,
checks the `com.apple.quarantine` attribute (the real XProtect trigger), and
verifies the binary's code signature. It does not use `spctl --assess`, which
rejects any signed CLI that is not a notarized app — codex is never
Apple-notarized. A quarantined or improperly-signed codex fails with
the next action to reinstall codex with the official curl installer into
`~/.local/bin`, then set `CODEX_PATH` to that binary. On non-Darwin platforms
the codex check is skipped and does not fail the command.

When Roundfix launches codex through `codex-acp` on macOS, it uses the same
configured-path-then-`PATH` resolution and passes a verified-clean codex through
`CODEX_PATH`. If no clean codex is available, Roundfix surfaces the hygiene
risk instead of silently spawning a known unsafe binary.

For latency-sensitive setups, configure direct adapter binaries in acpx config
so default adapters do not launch through `npx -y` on first use.

Known constraint: acpx `0.12.0` has a hard 10 MiB queue-owner per-message
buffer with no CLI, config, or environment override found in the pinned
package. Large docs-task payloads, especially turns that print or return large
skill/docs file content, can trigger `-32603 Message buffer exceeded 10485760
bytes`. Roundfix proceeds to verification when acpx delivered a parsed prompt
result before that exit; if completed work is preserved but a Task stays
failed, review the kept Task Worktree or Run Worktree and run the Settle
Command.

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

Once the binary is on `PATH`, use `roundfix setup` to make the machine
Run-ready, `roundfix doctor` to diagnose readiness without changing anything,
and `roundfix upgrade` to update the binary. The Commands and Command
Boundaries sections below cover their flags and outcomes.

## GitHub Access

Roundfix uses the GitHub CLI (`gh`) from the local machine. It does not ask for
or store a GitHub token directly. Authenticate `gh` for the target repository
before running operational commands:

```bash
gh auth status
```

## Commands

These examples call the installed `roundfix` binary. From a source checkout
without installing, substitute `go run ./cmd/roundfix` for `roundfix`.

Show help:

```bash
roundfix --help
```

Show version:

```bash
roundfix --version
roundfix -v
```

Create a Project Config in the current repository:

```bash
roundfix init
```

Create a User Config instead:

```bash
roundfix init --scope user
```

Verify and prepare this machine for Roundfix Runs:

```bash
roundfix setup
```

Diagnose this machine for Roundfix Runs without installing or writing config:

```bash
roundfix doctor
```

Upgrade Roundfix or check the release channel:

```bash
roundfix upgrade
roundfix upgrade --check
```

Fetch unresolved CodeRabbit Review Issues into local Round artifacts:

```bash
roundfix fetch --source coderabbit --pr <number> [--spec <slug>]
```

Resolve downloaded Compatible Artifacts with a selected Agent:

```bash
roundfix resolve --pr <number> --agent codex [--spec <slug>]
```

Run the watched review-resolution loop:

```bash
roundfix watch --source coderabbit --pr <number> --agent codex [--spec <slug>] --until-clean
```

Execute a Spec's Task Graph:

```bash
roundfix implement --spec <slug> --agent codex
```

Start a Detached Run for scripts or CI. The `--detach` flag is available on
`resolve`, `watch`, and `implement`:

```bash
roundfix implement --spec <slug> --agent codex --detach
```

Detached Runs print exactly four stdout lines — the Run ID, the console log
path, and the `attach`/`stop` follow-up commands:

```text
Run detached: <run-id>
Console log: <path>
Follow: roundfix attach <run-id>
Stop: roundfix stop <run-id>
```

See Command Boundaries for the full detach contract.

Settle one failed Spec Task whose completed work is already in a kept Task
Worktree, kept Run Worktree, or the current repository:

```bash
roundfix settle --spec <slug> --task <task_id>
```

Archive a completed Spec after all Tasks are completed and the newest QA
Report has `verdict: pass`:

```bash
roundfix archive <slug>
```

Preview or reclaim old terminal Run journal and run artifact storage:

```bash
roundfix gc --dry-run
roundfix gc
```

Stop a live Run gracefully, or force-stop a dead or runaway Run:

```bash
roundfix stop <run-id>
roundfix stop --force <run-id>
```

List, validate, or install the bundled Roundfix skills:

```bash
roundfix skills list
roundfix skills check
roundfix skills install
```

The binary ships 14 Roundfix-owned skills: the operational `roundfix` skill plus
the authorial workflow skills (`write-idea`, `write-prd`, `write-techspec`,
`write-tasks`, `setup-workflow`, `implement-task`, `implement-spec`,
`brainstorming`, `council`, `business-analyst`, `archive-spec`, `qa-gate`,
`evidence-gate`). `skills list` also prints the recommended external skills,
which install through your own skills tooling and are never shipped.

By default, `skills install` writes all bundled skills to
`<repo>/.agents/skills`. Use `--target codex`, `--target claude`,
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
- `doctor` runs the read-only readiness checks for Node.js, pinned acpx, the
  configured Agent probe, and codex runtime hygiene. It prints one stdout line
  per check with `ok`, `failed`, or `skipped`; failure lines include
  `next: <action>` when a remediation is known. It mutates nothing and exits
  nonzero when any check fails. The codex check is macOS-only: it inspects
  `com.apple.quarantine` and code-signature validity, reports the curl reinstall
  into `~/.local/bin` as the next action for a quarantined or improperly-signed codex,
  and is skipped on non-Darwin platforms.
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
  assigned issues, integrates the Run Worktree, and runs Final Push only when
  no Unresolved Review Issues remain.
- `watch` waits for CodeRabbit status on the current PR HEAD, observes the
  configured quiet period, fetches unresolved issues, resolves Batches, and
  repeats until `Clean`, `MaxRoundsReached`, `BudgetExceeded`, `TimedOut`,
  `Failed`, or `Stopped`.
- `implement` executes a Spec's Task Graph in a Run Worktree as one Run. The
  scheduler executes the current Wave up to `worktree.concurrency` at a time
  (default `2`; `1` keeps sequential behavior), with concurrently running
  Tasks in Task Worktrees and one commit per completed Task on the Run Branch.
  `implement.auto_push: true` makes a Clean spec Run push its branch upstream
  and append `pushed <remote>/<branch>` to stdout. Integration Pending,
  Unresolved Outcome, Failed, Stopped, and failing-QA Runs never push.
- `resolve`, `watch`, and `implement` accept `--detach` to start a Detached
  Run. The foreground command prints the four-line report, exits `0`, and
  leaves follow-up control to `roundfix attach <run-id>` and
  `roundfix stop <run-id>`. Detached startup failures before the handshake
  relay the child's stderr and exit code verbatim, with no stdout report.
- `settle` targets one failed Task by resolving its kept Task Worktree first,
  then its kept Run Worktree, then the current repository. It re-runs the
  Task's Verification commands in the selected tree, changes nothing when
  verification fails, and on pass settles it `completed`, creates the standard
  Task commit, creates no Run, writes no Run Event Journal entries, and never
  pushes. Task Worktree settlements integrate onto the Run Branch before the
  existing Run-level integration runs.
- `archive` is non-interactive, creates no Run, and never pushes. Before
  touching the filesystem, it verifies every Task in the Spec's Task Graph is
  `completed` and that the newest QA Report has `verdict: pass`. On pass, it
  stamps `_prd.md` with `status: archived`, `archived`, and `source_slug`, then
  moves `docs/specs/<slug>/` to `docs/specs/_archived/<slug>/`. Refusals exit
  `2`, write the Preflight Validation failure to stderr, and leave the folder
  in place.
- `gc` is non-interactive. It resolves `store.journal_retention`, computes the
  cutoff, prunes eligible terminal Runs' Run Event Journal rows and
  `<artifact-dir>/runs/<run-id>` directories, removes orphaned `runs/<id>`
  directories under the resolved run artifact root, and reports Runs, journal
  rows, and artifact bytes reclaimed on stdout. `--dry-run` lists the same
  eligible set without deleting anything. `journal_retention: 0` skips pruning
  and reports that no pruning was performed. Retention never deletes Active
  Runs, `runs` rows, active-run locks, or Review artifacts under `docs/specs/`.
- `stop` is graceful by default. It records a Stop Request in the Run Database
  and reports `Stop Request recorded; the Run stops after the current Work Item
  settles.` Use `--force` only for a dead, stuck, or runaway Run; it cancels
  the Agent Session best-effort, completes the Run Stopped immediately, and
  releases its Active Run locks. It also reaps kept Run or Task Worktrees and
  branches for terminal Runs whose branch has no commits beyond its base,
  reporting each removed pair on stderr as
  `roundfix: reaped terminal Worktree path=<path> branch=<branch>`.
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

Removed keys that Roundfix registers as deprecated never break an existing
config: Roundfix ignores them and prints one stderr warning naming the
replacement. The current deprecated key is `resolve.concurrent`, which prints
`config: resolve.concurrent is deprecated and ignored; use worktree.concurrency`
and then continues. Unknown keys that are not registered as deprecated still
fail strict validation.

Example:

```yaml
defaults:
  agent: codex
  verification: make verify
  # Empty uses Roundfix Home artifacts/<repo-id>; set a path to override.
  artifact_dir: ""
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

worktree:
  # Parent directory; Roundfix always appends <repo-slug>/<run-id>[.<task_id>].
  location: "~/.roundfix/worktrees"
  # Maximum concurrent Task Worktrees for spec Runs; 1 keeps sequential behavior.
  concurrency: 2
  # Repository-relative untracked files copied into each new Run or Task Worktree.
  # List only files that are already gitignored.
  copy: []
  # Command run in each new worktree after copy and before Agent work; empty disables it.
  bootstrap: ""
  # Maximum time before the bootstrap command fails the owning Run or Task.
  bootstrap_timeout: 10m

store:
  # Terminal Run journals older than this duration are eligible for pruning; 0 keeps everything.
  journal_retention: 336h

budget:
  enabled: true
  max_run_duration: 2h

resolve:
  batch_size: 3
```

Per-Batch Agent log files are development opt-in. `logs.agent` defaults to
`false` in built-in, User, and Project Config; when it is false, Roundfix still
records every Agent payload in the Run Event Journal. Set it to `true` in User
or Project Config to write the per-Batch files again:

```yaml
logs:
  agent: true
```

`store.journal_retention` defaults to `336h` (14 days). It accepts Go duration
strings, and `0` keeps every Run Event Journal and run artifact directory. A
non-zero value makes terminal Runs older than the window eligible for journal
and run artifact pruning. Active Runs are never eligible, and retention never
deletes `runs` rows or active-run locks. The `implement`, `resolve`, and
`watch` preflight sweep runs the same prune best-effort and reports one stderr
summary when it frees storage; failures are warnings and do not block the Run.
Use `roundfix gc --dry-run` to preview the same terminal Run set, and
`roundfix gc` to reclaim on demand.

## Local State

- Run Database: `~/.roundfix/roundfix.db`
- Run Worktrees: `<worktree.location>/<repo-slug>/<run-id>`
- Task Worktrees: `<worktree.location>/<repo-slug>/<run-id>.<task_id>`
- Default Artifact Directory: Roundfix Home `artifacts/<repo-id>`, used for
  Artifact Directory-backed Run files and as the base when an explicit
  Artifact Directory is configured.
- Review Issue artifacts, with the ADR-0029 resolver:
  - Explicit `--artifact-dir` or `defaults.artifact_dir` preserves the legacy
    layout: `<artifact-dir>/reviews/pr-<number>/round-<nnn>/issue_<nnn>.md`.
  - Otherwise, a PR associated with a Spec writes
    `docs/specs/<slug>/reviews/round-<nnn>/issue_<nnn>.md`.
  - Without a valid Spec association, review artifacts write to
    `docs/specs/_reviews/pr-<number>/round-<nnn>/issue_<nnn>.md`.
- Per-Batch Agent logs, only when `logs.agent: true`:
  `<artifact-dir>/runs/<run-id>/agent/batch-<nnn>.log`
- Detached Run console log, always written for Detached Runs:
  `<artifact-dir>/runs/<run-id>/console.log`
- Journal Retention prunes only terminal Run Event Journal rows and
  `<artifact-dir>/runs/<run-id>` directories. Review artifacts under
  `docs/specs/<slug>/reviews/` or `docs/specs/_reviews/` are outside retention
  scope.

For review commands, explicit `--spec <slug>` wins over trailer discovery. When
`--spec` is absent, Roundfix uses the newest `Roundfix-Spec: <slug>` trailer on
the PR head commit only if `docs/specs/<slug>/` exists; an unknown or invalid
slug falls back to the spec-less `_reviews` path. Roundfix never commits or
gitignores review artifacts, so versioning them stays a repository decision.
ADR-0030 keeps per-Batch Agent log files off by default because the Run Event
Journal already stores the raw payloads; the Detached Run console log remains
unconditional for the ADR-0028 detach contract.

With automatic Round selection, `fetch` reuses an existing matching Round when
the same HEAD already has the same Review Issue fingerprints. If the fetched
payload is new, Roundfix writes the next Round directory. Roundfix does not
overwrite existing Round artifacts. Repeated findings across different payloads
are deduplicated later during `resolve` by Review Issue Fingerprint, preferring
`source_ref` such as `thread:<id>,comment:<id>` and falling back to
`review_hash`.

Operational Runs that start an Agent work in Run Worktrees. `worktree.location`
sets only the parent directory with Project Config > User Config > built-in
default precedence. Roundfix always appends the readable repo slug and Run ID;
for concurrently executing Tasks, Task Worktrees are siblings of the Run
Worktree and append `.<task_id>`. Those final path segments are fixed and not
configurable.

A new Run Worktree starts from committed Git state, so untracked files in the
user's checkout are absent unless they are listed under `worktree.copy`; add
repository-relative paths there when Verification or local tooling needs
untracked files. Copied environment files must already be gitignored; Roundfix
does not add ignore rules for arbitrary copied files, and Batch or Task commits
can include changed repository files. Dirty user checkout behavior is
command-specific: `implement` no longer blocks on it and instead prints a note
that overlapping local changes end the Run in Integration Pending. Other
operational commands retain their existing preflight rules, including the local
Project Config allowances for `fetch`, `resolve`, and `watch` at
`.roundfixrc.yml`. Batch commits exclude that config file so local setup
changes do not mix with review fixes. Terminal Run outcomes release the Active
Run lock for the PR Head Branch.

Worktree Bootstrap prepares each newly created Run or Task Worktree after
`worktree.copy` placement and before Agent work and Verification.
`worktree.bootstrap` is a shell command run in the worktree root; an empty value
skips the step. `worktree.bootstrap_timeout` defaults to `10m` and bounds the
command. Bootstrap output streams to stderr and the Run Event Journal. A start
failure, non-zero exit, or timeout fails the owning Run for a Run Worktree or
settles only the owning Task failed for a Task Worktree, with a message shaped
as `worktree bootstrap failed: <command>: <reason>`.

Roundfix owns running and timing the bootstrap command. Dependency installation,
database provisioning, migrations, seeding, and cache strategy belong in the
configured command. For a stateful monorepo that uses one shared database, keep
Task execution sequential so bootstrap runs once on the reused Run Worktree:

```yaml
worktree:
  concurrency: 1
  copy: [".env", "packages/backend/.env"]
  bootstrap: "bun install && bun run db:migrate && bun run db:seed"
  bootstrap_timeout: 10m
```

The files listed in `worktree.copy` must be repository-relative and must stay
inside the repository. They must also be gitignored before they are copied.

Spec Runs schedule Tasks by Wave. `worktree.concurrency` defaults to `2`,
which can run two Verification commands at once; if those commands are heavy
(`make verify`, for example), expect matching local CPU and cache load. Set
`worktree.concurrency: 1` for sequential execution.

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
skills/                          shipped Roundfix skill bundle (synced from .agents/skills)
docs/                            product docs and architecture decisions
```

Start with:

- [Operational guide](docs/usage.md)
- [CONTEXT-driven development](docs/context-driven-development.md)
- [Project glossary](CONTEXT.md)
- [Architecture decisions](docs/adr/)
- [Release runbook](docs/release-runbook.md)

## License

MIT. See [LICENSE](LICENSE).
