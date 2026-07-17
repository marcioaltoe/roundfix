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
- `defaults.model` → `profiles.<category>.preferred.model`

Unknown keys that are not registered as deprecated fail strict validation.

## Full example

```yaml
defaults:
  verification: make verify
  # Empty uses Roundfix Home artifacts/<repo-id>; set a path to override.
  artifact_dir: ""
  auto_commit: true

profiles:
  general:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.6-terra
        reasoning_effort: max
  backend:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.6-terra
        reasoning_effort: max
  frontend:
    preferred:
      runtime: claude
      model: claude-fable-5
      reasoning_effort: medium
    fallbacks:
      - runtime: codex
        model: gpt-5.6-sol
        reasoning_effort: high
  qa:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.6-terra
        reasoning_effort: max
  review:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.6-terra
        reasoning_effort: max

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

## Agent selection profiles

Roundfix owns Agent Selection Profiles in Project Config, User Config, and
built-ins. Each profile is one atomic object: a complete Preferred Selection
plus a non-empty ordered Fallback Chain. Project Config replaces User Config,
User Config replaces built-ins, and no tuple field or fallback entry merges
across scopes.

Required profiles are `general`, `backend`, `frontend`, `qa`, and `review`.
Optional Task Type profiles `data`, `infra`, `docs`, `test`, and `chore`
inherit the effective `general` profile when absent. `roundfix profiles show`
labels that inherited recommendation source instead of duplicating stored
config.

Built-in required profiles use these official identifiers:

- `general`, `backend`, `qa`, `review`: preferred `codex / gpt-5.6-sol / high`;
  fallback `codex / gpt-5.6-terra / max`.
- `frontend`: preferred `claude / claude-fable-5 / medium`; fallback
  `codex / gpt-5.6-sol / high`.

Custom model strings are accepted verbatim for forward-compatible ACP proof.
They are not added to an allowlist, and the installed runtime adapter proves
them before any Run mutation. An explicit empty `reasoning_effort: ""` means
model-managed reasoning; omitted `reasoning_effort` is invalid because the
runtime would not receive a complete tuple.

Manage profiles with:

```bash
roundfix profiles show --category backend
roundfix profiles show --all --json
roundfix profiles configure --scope project --file profiles.yml --dry-run --json
roundfix profiles validate --category review --json
```

`profiles configure` writes only the `profiles` schema after strict validation
and confirmation. `--dry-run` leaves config bytes unchanged and reports
`changed: false` in JSON. `profiles validate` proves the selected tuples
through disposable ACP sessions, deduplicates exact tuples, and closes every
session without creating a Run, worktree, prompt, credential, or runtime-owned
setting.

One-Run flags such as `--agent`, `--model`, and `--reasoning-effort` replace
only the Preferred Selection for each relevant category and preserve that
category's Fallback Chain. When one override spans multiple Task or QA
categories, text and JSON output include a cross-category warning.

Recommendations shown by `profiles show` are a dated advisory snapshot
(`2026-07-16`) with five entries per category, benchmark/result/cost evidence,
rationale, and `category_specific: false`. Recommendation rank never changes
configuration, proof order, Preferred Selection, Fallback Chain, or routing.

Automatic fallback is pre-prompt only. Roundfix records and shows the fallback
notification before activating the next configured tuple. Once
`agent_work_started` is recorded, no prompt, tool, verification, cancellation,
rate-limit, or session-loss failure can start a replacement session.

Legacy configs with `defaults.agent` or top-level `runtimes` still load during
the compatibility window only when the same file has no `profiles` section.
Do not edit runtime-owned settings or credentials for migration; write complete
profiles with `roundfix profiles configure --scope user|project`.

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
