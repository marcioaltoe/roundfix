# Configuration

Every Roundfix config key, where it lives, and the local state it produces.
For per-command behavior see the [command reference](commands.md); for the
task-oriented walkthrough see the [operational guide](usage.md).

## Sources and precedence

Roundfix reads YAML config in this order, later sources overriding earlier
ones per key:

1. Built-in defaults.
2. User Config at `~/.roundfix/config.yml`.
3. Project Config at `<repo>/.roundfixrc.yml`.
4. CLI flags.

Create config with `roundfix init` (Project Config) or
`roundfix init --scope user`; `--force` overwrites an existing file.

Removed keys that Roundfix registers as deprecated never break an existing
config — they are ignored with one stderr warning naming the replacement:

- `resolve.concurrent` → `worktree.concurrency`
- `defaults.model` → `runtimes.<runtime>.model`

Unknown keys that are not registered as deprecated fail strict validation.

## Full example

```yaml
defaults:
  agent: codex
  verification: make verify
  # Empty uses Roundfix Home artifacts/<repo-id>; set a path to override.
  artifact_dir: ""
  auto_commit: true

runtimes:
  codex:
    model: gpt-5.6-sol
    # Empty reasoning_effort means the Agent Model manages reasoning.
    reasoning_effort: ""
  claude:
    model: opus
    reasoning_effort: ""
  opencode:
    model: ""
    reasoning_effort: ""

review_source:
  name: coderabbit
  include_nitpicks: true

watch:
  until_clean: true
  max_rounds: 6
  poll_interval: 30s
  review_timeout: 30m
  quiet_period: 30s
  # How long watch polls for the Review Source check on the pushed HEAD
  # before ending CleanUnverified instead of Clean.
  check_grace_period: 5m
  auto_push: true

implement:
  auto_push: false

notify:
  # Send one notification when resolve, watch, or implement reaches a terminal outcome.
  enabled: true
  # Empty uses the native desktop notifier when available.
  command: ""

specs:
  # Directory holding Spec folders; relative values resolve against the repo root.
  root: "docs/specs"

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

## Agent selection

Roundfix owns the Agent Model and Default Reasoning Effort for every Agent
Session. It resolves each value independently — built-in defaults, User
Config, Project Config, then one-Run flags — and never reads or mutates
runtime-owned model configuration.

Built-in selections:

- Codex: `model: gpt-5.5`, `reasoning_effort: xhigh`.
- Claude: `model: opus`, `reasoning_effort: ""` (model-managed).
- OpenCode: no built-in model. Provide one through config, its one-Run flag,
  or Interactive Input.

An empty `reasoning_effort` is a deliberate model-managed selection: Roundfix
assigns the Agent Model and skips the runtime-specific reasoning option, and
Run headers render `Default Reasoning Effort: model-managed`. This is the
required shape for the codex `gpt-5.6` family, whose models manage reasoning
and reject every `reasoning_effort` value exposed by codex-acp.

Interactive Input asks for Agent, Agent Model, then Default Reasoning Effort.
The Codex Model Catalog is ordered `gpt-5.6-sol`, `gpt-5.6-terra`,
`gpt-5.6-luna`, `gpt-5.5`, `gpt-5.4`, `gpt-5.4-mini`, `gpt-5.3-codex-spark`;
the Claude catalog is `Default`, `Opus`, `Fable`, `Sonnet`, `Haiku`, with
`Default` showing the concrete configured model. Catalogs are picker data, not
allowlists: non-interactive flags and typed values can pass a custom pair,
which the installed ACP adapter validates during Preflight Validation. A Batch
that later fails because the Agent Session rejects the model reports
`Agent Model "<model>" not advertised by runtime "<runtime>"; advertised: <list>`
— and `roundfix doctor` shows the same probe as its `model:` line.

Override both values for one Run without touching config:

```bash
roundfix resolve --pr 123 --agent codex --model gpt-5.6-sol --reasoning-effort "" --no-input
roundfix implement --spec example-spec --agent claude --model opus --reasoning-effort "" --qa --detach
```

An explicit empty `--reasoning-effort ""` is the one-Run model-managed
override; an explicit empty `--model ""` is invalid and exits `2`.

### Fallback Selection

If a runtime rejects the selected Agent Model or a non-empty reasoning value,
Roundfix probes that runtime's catalog newest-first and its reasoning
vocabulary highest-first — never crossing to another runtime and never
re-proposing the failed model. Interactively, it presents the proven Fallback
Selection and asks `Use this Fallback Selection for this Run? [y/N]:`; the
confirmation applies to that Run only. With `--no-input`, `--detach`, or
non-interactive stderr, it exits `2` before creating a Run and prints one
concrete `Re-run:` line with explicit `--model` and `--reasoning-effort`
flags. There is no flag or key that pre-authorizes a fallback.

Run startup and inspection surfaces show the stored selection
(`Agent Model: …`, `Default Reasoning Effort: …`); `attach` reads those values
from the Run row, so later config changes never rewrite a historical Run.

## Spec Root

`specs.root` is the directory holding Spec folders; default `docs/specs`.
Relative values resolve against the repository root; absolute values are used
as-is. Roundfix resolves the root once at command start and carries the
absolute path into Run and Task Worktrees, so every surface reads the same
Spec artifacts. When the resolved root is outside the repository working tree
after symlink evaluation, Spec artifacts are external and are never staged
into code-repository commits.

## Worktrees and bootstrap

Spec Runs execute in Run Worktrees; concurrently executing Tasks get sibling
Task Worktrees (`<location>/<repo-slug>/<run-id>.<task_id>`). Review Runs run
in the user's checkout and create no worktree. `worktree.location` sets only
the parent directory — the final path segments are fixed.

A new worktree starts from committed Git state: untracked files are absent
unless listed in `worktree.copy` (repository-relative, inside the repository,
and already gitignored). `worktree.bootstrap` runs once in each new worktree
after copy and before Agent work, bounded by `worktree.bootstrap_timeout`
(default `10m`); a start failure, non-zero exit, or timeout fails the owning
Run or settles the owning Task failed with
`worktree bootstrap failed: <command>: <reason>`. Dependency installation,
migrations, seeding, and caching belong in the configured command. For a
stateful monorepo with one shared database, keep execution sequential so
bootstrap runs once:

```yaml
worktree:
  concurrency: 1
  copy: [".env", "packages/backend/.env"]
  bootstrap: "bun install && bun run db:migrate && bun run db:seed"
  bootstrap_timeout: 10m
```

`worktree.concurrency: 2` (the default) can run two Verification commands at
once — with heavy gates like `make verify`, expect matching CPU and cache
load.

## Notifications

`notify.enabled` defaults to `true`. With the default empty `notify.command`,
Roundfix uses the native desktop path (`osascript` on macOS, `notify-send` on
Linux, silent no-op elsewhere). A non-empty command replaces the native path
and runs through the shell with a 30s timeout, receiving `ROUNDFIX_RUN_ID`,
`ROUNDFIX_OUTCOME`, `ROUNDFIX_KIND`, and `ROUNDFIX_TARGET`. Only `resolve`,
`watch`, and `implement` notify, once, at the terminal outcome; Detached Runs
notify from the detached child. Notification failures warn on stderr and never
change the Run outcome.

## Logs and retention

`logs.agent` defaults to `false`: no per-Batch Agent log files are written,
while the Run Event Journal still records every Agent payload losslessly. Set
it to `true` to write `<artifact-dir>/runs/<run-id>/agent/batch-<nnn>.log`
files for debugging. The Detached Run console log is always written.

`store.journal_retention` defaults to `336h` (14 days); `0` keeps everything.
Terminal Runs older than the window become eligible for journal and run
artifact pruning — by `roundfix gc` on demand and by a best-effort sweep at
`implement`/`resolve`/`watch` startup. Active Runs, `runs` rows, active-run
locks, and Review artifacts under the Spec Root are never pruned.

## Local state

- Run Database: `~/.roundfix/roundfix.db`
- Run Worktrees: `<worktree.location>/<repo-slug>/<run-id>` (spec Runs only)
- Task Worktrees: `<worktree.location>/<repo-slug>/<run-id>.<task_id>`
- Spec Root: `specs.root`, default `<repo>/docs/specs`
- Default Artifact Directory: Roundfix Home `artifacts/<repo-id>`
- Review Issue artifacts (ADR-0029 resolver):
  - Explicit `--artifact-dir` or `defaults.artifact_dir`:
    `<artifact-dir>/reviews/pr-<number>/round-<nnn>/issue_<nnn>.md`
  - PR associated with a Spec:
    `<specs.root>/<slug>/reviews/round-<nnn>/issue_<nnn>.md`
  - No Spec association:
    `<specs.root>/_reviews/pr-<number>/round-<nnn>/issue_<nnn>.md`
- Per-Batch Agent logs (only `logs.agent: true`):
  `<artifact-dir>/runs/<run-id>/agent/batch-<nnn>.log`
- Detached Run console log: `<artifact-dir>/runs/<run-id>/console.log`
- Failed Verification diagnostics:
  `<artifact-dir>/runs/<run-id>/verification/batch-<nnn>-attempt-<1|2>.log`

For review commands, explicit `--spec <slug>` wins over trailer discovery;
otherwise Roundfix uses the newest `Roundfix-Spec: <slug>` trailer on the PR
head commit when that Spec folder exists, falling back to the `_reviews` path.

## How Runs use context

Roundfix keeps Agent context bounded by design:

- Agent prompts receive the assigned Work Item contract and bounded guidance.
  Successful Verification output never enters Agent context; on an attempt-1
  failure the Daemon sends exactly one Verification Feedback prompt with the
  failed command, wrapped failure, and diagnostic artifact path, then reruns
  the full Verification sequence and settles from that verdict.
- Spec Task prompts include one complete assigned Task plus a Spec Context
  Bundle: paths for standard Spec artifacts, Task-authored `## Context`
  entries (capped at 50 unique paths), and sorted files changed by prior
  integrated Tasks — the whole manifest capped at 200 paths. The bundle embeds
  no full PRD, TechSpec, skill document, source file, or prior diff.
- Console Logs and the Live Run View render file reads and edits as compact
  lines (`edit internal/daemon/task_engine.go (+8/-3)`); the Run Event Journal
  remains the lossless evidence boundary.

## Known constraint: acpx message buffer

acpx `0.12.0` has a hard 10 MiB queue-owner per-message buffer with no
override. Turns that print or return very large file content can trigger
`-32603 Message buffer exceeded 10485760 bytes`. Roundfix proceeds to
verification when acpx delivered a parsed prompt result before that exit
(ADR-0020); if completed work is preserved but a Task stays failed, recover it
with the Settle Command. For latency-sensitive setups, configure direct
adapter binaries in acpx config so default adapters do not launch through
`npx -y` on first use.
