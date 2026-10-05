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

### Omitted and empty values

An omitted key inherits the value from the next lower-precedence source. For
example, a Project Config that omits `notify.command` keeps the User Config
value when one exists, then falls back to the built-in default.

An explicit empty string (`""`) or empty list (`[]`) overrides the lower
source. Empty values are valid only where this page defines their behavior:

| Explicit empty value | Behavior |
| --- | --- |
| `defaults.artifact_dir: ""` | Uses `~/.roundfix/artifacts/<repo-id>` |
| `worktree.copy: []` | Copies no ignored files into Run or Task Worktrees |
| `worktree.bootstrap: ""` | Disables Worktree Bootstrap |
| `watch.push_remote: ""` and `watch.push_branch: ""` | Uses the upstream remote and branch detected by Preflight Validation |
| `notify.command: ""` | Uses the native desktop notifier when available |
| `profiles.<category>.*.reasoning_effort: ""` | Lets the model manage reasoning effort |

Do not use a bare YAML value such as `command:` to mean an empty string. Use
the explicit value shown above or omit the key to inherit it.

Setting `verification.repository_at_settlement: false` turns off only the
repository Verification appended at settlement for non-QA Tasks of a graph
with a QA gate Task; declared Task commands, entry preconditions, and the QA
gate keep their existing Verification.

`verification.tools` lists additional tools Doctor checks on `PATH` without
running them, such as tools used inside a Makefile. It defaults to `[]` and is
accepted in both User and Project Config; a Project Config list replaces the
User Config list, including an explicit empty list. Each entry must be a bare
executable name matching `^[A-Za-z0-9][A-Za-z0-9._+-]*$`, for example `go` or
`gofmt`. Paths, arguments, empty names and duplicate entries are config errors
naming `verification.tools`. Doctor also checks the tools named at the start
of the effective Verification, bootstrap and regeneration commands and the
Setup Manifest's Verification decisions.

## Task and Verification capacities

Task Capacity and Verification Capacity are independent, config-only limits
for one Implement Run. Use the built-in Task Capacity `2` and Verification
Capacity `1` when two implementation-ready Tasks can overlap but the
repository gate must run sequentially:

```yaml
worktree:
  concurrency: 2

verification:
  concurrency: 1
```

`worktree.concurrency` limits concurrent Task Worktree lifecycles, including
Agent work. `verification.concurrency` limits concurrent Task Verification
attempts inside that Run. The built-in values are `2` and `1`, respectively;
User Config overrides each built-in value, and Project Config overrides the
User Config value for the same key. Neither capacity has a command-line flag.

Both values must be positive integers. Zero, negative, fractional, and
otherwise invalid values fail strict configuration validation before the Run
is created. Omitting `verification.concurrency` uses `1`; it never inherits
`worktree.concurrency`.

Verification Capacity is not a machine-wide lock. It does not coordinate
other Roundfix Runs, CI jobs, manually started commands, or other processes.
Projects must still choose a safe Task Capacity for Worktree Bootstrap, Agent
work, and resources used outside Daemon Verification.

## Active Implement Run ceiling

`runs.max_active` limits Active Implement Runs across every repository in the
Run Database. The built-in value is `3`; set it in User Config to the measured
capacity of the machine:

```yaml
runs:
  max_active: 3
```

When the number of Active Implement Runs reaches the ceiling, `implement`
refuses before creating another Run and lists each holding Run's id,
repository, and Spec. Set `max_active: 0` to disable the bound. Negative values
fail configuration validation. Only User Config can set this machine-wide
ceiling. A `runs.max_active` value in Project Config is ignored with a warning,
so a repository cannot raise or disable the user's limit.

## Jev monthly ceiling

`jev.monthly_ceiling_usd` is a finite number of US dollars greater than zero.
Set it only in User Config at `~/.roundfix/config.yml`; when unset, the
effective ceiling is US$5. The judge, the Jev Router gate, and its key-limit
check all read the same value, and the key's monthly credit limit must be at
most that value:

```yaml
jev:
  monthly_ceiling_usd: 50
```

A `jev.monthly_ceiling_usd` value in Project Config is ignored with this
warning on standard error:
`config: jev.monthly_ceiling_usd in Project Config is ignored; set jev.monthly_ceiling_usd in User Config`.
An invalid User Config value fails with:
`jev.monthly_ceiling_usd must be a finite number greater than 0`.

An older Roundfix binary refuses a User Config containing this key as an
unknown key. Set it only after every Roundfix binary on the machine includes
this change.

## Context-Driven Baseline state

User Config and Project Config are operational Roundfix state. They do not
select, version, or authorize a Context-Driven Baseline. The
`setup-context-driven` workflow stores its Setup Manifest and resolved setup
decisions in `docs/agents/setup-context.json`; structured setup and Baseline
Readoption answers use a separate
`setup-context-driven/decisions/0.0.1` decision file supplied with
`--decision-file`.

The `0.0.1` reset does not recreate, migrate, or renumber User Config, Project
Config, Runs, or Run Database state. See
[CONTEXT-driven development](context-driven-development.md#configure-or-audit-the-baseline)
for the audit, decision, preview, apply, formatter, Verification, audit, and
reapply workflow.

## Full example

```yaml
defaults:
  # false keeps each ACP Runtime's normal sandbox or permission mode.
  agent_full_access: false
  # Verification command for review Batches; Spec Tasks use their task file commands.
  verification: make verify
  # Empty uses Roundfix Home artifacts/<repo-id>; set a path to override.
  artifact_dir: ""
  auto_commit: true

profiles:
  general:
    preferred:
      runtime: codex
      model: gpt-6.1-sol
      reasoning_effort: "high"
    fallbacks:
      - runtime: claude
        model: opus
        reasoning_effort: "high"
  backend:
    preferred:
      runtime: codex
      model: gpt-6.1-sol
      reasoning_effort: "high"
    fallbacks:
      - runtime: claude
        model: opus
        reasoning_effort: "high"
  frontend:
    preferred:
      runtime: claude
      model: opus
      reasoning_effort: "high"
    fallbacks:
      - runtime: codex
        model: gpt-6.1-sol
        reasoning_effort: "xhigh"
  qa:
    preferred:
      runtime: codex
      model: gpt-6.1-sol
      reasoning_effort: "high"
    fallbacks:
      - runtime: claude
        model: opus
        reasoning_effort: "high"
  review:
    preferred:
      runtime: codex
      model: gpt-5.6-luna
      reasoning_effort: "max"
    fallbacks:
      - runtime: codex
        model: gpt-6.1-sol
        reasoning_effort: "high"

# Legacy PR-feedback Review Source, read only by fetch, watch and resolve; it never selects or requests a pre-PR reviewer.
review_source:
  name: coderabbit
  # false excludes CodeRabbit findings whose severity is nitpick.
  include_nitpicks: false

watch:
  until_clean: true
  max_rounds: 6
  poll_interval: 30s
  review_timeout: 30m
  quiet_period: 30s
  # How long watch polls for the Review Source check on the pushed HEAD
  # before ending CleanUnverified instead of Clean.
  check_grace_period: 5m
  # Final Push uses the detected upstream when both values are empty.
  auto_push: true
  push_remote: ""
  push_branch: ""

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

runs:
  # Maximum Active Implement Runs across repositories; 0 disables the bound.
  max_active: 3

worktree:
  # Parent directory; Roundfix always appends <repo-slug>/<run-id>[.<task_id>].
  location: "~/.roundfix/worktrees"
  # Maximum concurrent Task Worktrees for spec Runs; 1 keeps sequential behavior.
  concurrency: 2
  # Repository-relative untracked files copied into each new Run or Task Worktree.
  # List only files that are already gitignored.
  # Empty copies no files.
  copy: []
  # Command run in each new worktree after copy and before Agent work; empty disables it.
  bootstrap: ""
  # Maximum time before the bootstrap command fails the owning Run or Task.
  bootstrap_timeout: 10m

verification:
  # Maximum concurrent Task Verification attempts within one Implement Run.
  concurrency: 1
  # Tools Doctor checks on PATH without running them.
  tools: [go, gofmt, git, make]
  # Append repository Verification when a non-QA Task of a gated graph settles.
  repository_at_settlement: true

store:
  # Terminal Run journals older than this duration are eligible for pruning; 0 keeps everything.
  journal_retention: 336h

logs:
  # false keeps Agent payloads only in the lossless Run Event Journal.
  agent: false

budget:
  enabled: true
  # budget.max_run_duration bounds how long a Run may run; the Run Window bounds when one may start.
  # An Implement Run's allowance renews at each Task settlement.
  max_run_duration: 2h

resolve:
  batch_size: 3
```

## Key reference

The defaults below apply when neither User Config nor Project Config sets the
key. Duration values use Go duration syntax such as `30s`, `10m`, and `2h`.

### General settings

| Key | Built-in default | Effect |
| --- | --- | --- |
| `defaults.agent_full_access` | `false` | Keeps the ACP Runtime's normal sandbox or permission mode. `true` requests its full-access mode before Agent work. |
| `defaults.verification` | `make verify` | Verification command for review Batches. Spec Tasks use the commands in each task file. |
| `defaults.artifact_dir` | `""` | Uses `~/.roundfix/artifacts/<repo-id>`. An absolute or repository-relative value overrides it. |
| `defaults.auto_commit` | `true` | Enables Daemon-owned review commits. `watch.auto_push: true` requires it. |
| `specs.root` | `docs/specs` | Locates active and archived Specs. Relative paths resolve from the repository root. |

### Review and Run settings

| Key | Built-in default | Effect |
| --- | --- | --- |
| `pre_pr_review.provider` | `codex` | Selects the independent reviewer for a candidate before a Pull Request exists. Supported values are `codex`, `claude`, `coderabbit`, and `none`. Project Config takes precedence over User Config, which takes precedence over this built-in default; an absent key inherits and never means `none`. Declaring a provider invokes nothing by itself; `roundfix review` runs the selected policy over the current candidate. Publication and merge gating remain separate. |
| `review_source.name` | `coderabbit` | Selects the legacy PR-feedback Review Source, read only by `fetch`, `watch`, and `resolve`; it never selects or requests a pre-PR reviewer. CodeRabbit is the only supported value. |
| `review_source.include_nitpicks` | `false` | Excludes CodeRabbit findings whose severity is `nitpick`. Set `true` to include them. |
| `watch.until_clean` | `true` | Continues the watch cycle until its clean-outcome contract or another bound ends the Run. |
| `watch.max_rounds` | `6` | Limits the number of review Rounds in one watch Run. |
| `watch.poll_interval` | `30s` | Sets the delay between Review Source polls. |
| `watch.review_timeout` | `30m` | Bounds the review wait across the watch Run. |
| `watch.check_grace_period` | `5m` | Bounds the final wait for Review Source evidence on the pushed head. |
| `watch.quiet_period` | `30s` | Requires this quiet interval before Roundfix fetches a settled review. |
| `watch.auto_push` | `true` | Runs Final Push after no Unresolved Review Issues remain. |
| `watch.push_remote` | `""` | Uses the upstream remote detected by Preflight Validation. |
| `watch.push_branch` | `""` | Uses the upstream branch detected by Preflight Validation. |
| `implement.auto_push` | `false` | Leaves a Clean Spec Run local. `true` pushes its upstream branch but never opens a pull request. |
| `runs.max_active` | `3` | Limits Active Implement Runs across every repository in the Run Database. Only User Config can set it; Project Config is ignored with a warning. `0` disables the bound. |
| `jev.monthly_ceiling_usd` | `5` | Sets the shared monthly Jev ceiling in User Config, in US dollars. Project Config is ignored with a warning. |
| `notify.enabled` | `true` | Sends one terminal outcome notification for `resolve`, `watch`, and `implement`. |
| `notify.command` | `""` | Uses the native desktop notifier. A non-empty shell command replaces it. |
| `budget.enabled` | `true` | Enforces the configured Run duration budget. |
| `budget.max_run_duration` | `2h` | Sets one Run Budget allowance. An Implement Run's allowance renews at each Task settlement; the Run Window bounds when a Run may start. |
| `resolve.batch_size` | `3` | Limits how many Review Issues one Agent Batch receives. |

### Worktree, log, and retention settings

| Key | Built-in default | Effect |
| --- | --- | --- |
| `worktree.location` | `~/.roundfix/worktrees` | Sets the parent directory for Run and Task Worktrees. |
| `worktree.concurrency` | `2` | Limits concurrent Task Worktrees. `1` keeps Task execution sequential. |
| `verification.concurrency` | `1` | Limits concurrent Task Verification attempts within one Implement Run, independently from Task Capacity. |
| `verification.tools` | `[]` | Declares tools Doctor must find on `PATH` without executing them. Project Config replaces the User Config list. |
| `verification.repository_at_settlement` | `true` | Appends `defaults.verification` after declared Verification for non-QA Tasks of a graph with a QA gate Task. `false` turns off only this appended command. |
| `worktree.copy` | `[]` | Copies no ignored files. Entries must be repository-relative and already ignored by Git. |
| `worktree.bootstrap` | `""` | Disables Worktree Bootstrap. A non-empty command runs after copy and before Agent work. |
| `worktree.bootstrap_timeout` | `10m` | Bounds each Worktree Bootstrap command. |
| `logs.agent` | `false` | Writes no per-Batch Agent log files; the Run Event Journal remains lossless. |
| `store.journal_retention` | `336h` | Keeps terminal Run journals and Run artifacts for 14 days. `0` disables pruning. |

## Agent selection profiles

Roundfix owns Agent Selection Profiles in Project Config, User Config, and
built-ins. Each profile is one atomic object: a complete Preferred Selection
plus a non-empty ordered Fallback Chain. Project Config replaces User Config,
User Config replaces built-ins, and no tuple field or fallback entry merges
across scopes.

### Jev Router

The Jev Router is an experimental OpenCode selection that only Project Config
may select. Set its effort to `""`: the router chooses the model and effort.
User Config and a one-Run override cannot name it. For example, opt in for
`docs` in your Project Config, with the current default as its fallback:

```yaml
profiles:
  docs:
    preferred:
      runtime: opencode
      model: roundfix-openrouter/typesafe/jev-router
      reasoning_effort: ""
    fallbacks:
      - runtime: codex
        model: gpt-6.1-sol
        reasoning_effort: high
```

Export `ROUNDFIX_OPENROUTER_API_KEY` in the environment that starts Roundfix.
Roundfix uses this key name and never `OPENROUTER_API_KEY`; without the key,
the selection is refused with `jev_router_key_missing`. The key value stays
in the environment. Roundfix gives OpenCode the provider through
`OPENCODE_CONFIG_CONTENT`, with `{env:ROUNDFIX_OPENROUTER_API_KEY}` as its key
placeholder, replacing an inherited value for routed sessions only.

Every routed prompt runs under the configured monthly Jev ceiling shared with
`roundfix spec judge`; `jev.monthly_ceiling_usd` is US$5 when unset. The router
changes no built-in or Recommended Profile.

Before a routed prompt starts, the key must also report a numeric `limit`
no greater than the ceiling and `limit_reset: monthly`. Set a monthly credit
limit of at most the configured ceiling on the key at OpenRouter; an unlimited
key, a lifetime limit, or a monthly limit above the ceiling is refused with
`jev_router_key_unbounded`. OpenRouter enforces this limit while a prompt is
running. A numeric `limit_remaining` at or below zero is refused with
`jev_ceiling_reached`. These refusals activate the configured fallback before
Agent work begins and fail the Work Item after work begins.

### Access policy readiness

When `defaults.agent_full_access: true`, profile readiness applies the
requested full-access mode during the disposable proof for the Preferred
Selection and every Fallback Selection. The proof records the effective access
policy. If a runtime has no mode for that policy or its adapter refuses the
mode, readiness refuses before a Run is created and names the failed predicate
and the remedy: disable full access or select a runtime that supports it.

When full access is not requested, readiness keeps the runtime's normal access
policy and performs no additional access-mode proof.

Required profiles are `general`, `backend`, `frontend`, `qa`, and `review`.
Optional Task Type profiles `data`, `infra`, `docs`, `test`, and `chore`
inherit the effective `general` profile when absent. `roundfix profiles show`
labels that inherited effective profile source instead of duplicating stored
config; its Recommended Profile remains specific to the requested category.

Built-in required profiles use these official identifiers:

- `general`, `backend`, and `qa`: preferred
  `codex / gpt-6.1-sol / high`, fallback `claude / opus / high`.
- `frontend`: preferred `claude / opus / high`, fallback
  `codex / gpt-6.1-sol / xhigh`.
- `review`: preferred `codex / gpt-5.6-luna / max`, fallback
  `codex / gpt-6.1-sol / high`.

The Model Catalog recognizes `gpt-6.1-sol`, `gpt-6-astra`, `gpt-6-sol`,
`gpt-6-luna`, `gpt-5.6-sol`, `gpt-5.6-terra`, and `gpt-5.6-luna` as official Codex identifiers, plus `opus`, `sonnet`,
`claude-fable-5-1`, `haiku`, and `default` as Claude
identifiers. Those are the values the
Claude adapter advertises, with the bracketed context suffix removed as the
capability parser removes it — the adapter advertises Opus 5.5 as `opus[1m]`.
Identifier validity does not prove operational availability. Recommendations are advisory
profiles only; the effective adapter in each environment must complete Exact
Agent Selection Proof before Roundfix can use a tuple. This exact proof is the
operational readiness authority. Custom model strings remain accepted verbatim
for forward-compatible proof and are never added to an allowlist. The dated
Recommended Profile supplies each required built-in Preferred Selection and
Fallback Chain; User Config and Project Config can replace them.

When an adapter advertises an independent reasoning control, Roundfix treats
every advertised Agent Model identifier as opaque. A bracketed identifier such
as `opus[1m]` is selectable exactly as printed, and its canonical prefix
`opus` remains selectable with a separate reasoning effort. The `[1m]` suffix
is a context-window annotation, not a reasoning effort. See
[ADR-0079](../adr/0079-independent-reasoning-controls-make-model-identifiers-opaque.md).

For Codex, Adapter Readiness requires the official
`@agentclientprotocol/codex-acp` package at version `2.0.1` or newer. For
Claude, it requires official `@agentclientprotocol/claude-agent-acp` at
version `0.84.0` or newer. The deterministic install actions are
`npm install -g @agentclientprotocol/codex-acp@2.0.1` and
`npm install -g @agentclientprotocol/claude-agent-acp@0.84.0`.

Setup writes `npx -y @agentclientprotocol/codex-acp@2.0.1` or
`npx -y @agentclientprotocol/claude-agent-acp@0.84.0` when an explicit
override needs migration. Migration follows from failed official lineage proof
rather than from recognizing any particular superseded package, so it covers a bare
override that resolves to a differently scoped package as well as an earlier
explicit pin such as `1.1.4`. Setup and Doctor report the effective command,
the required official package, and the applicable install action.

An explicit empty `reasoning_effort: ""` means model-managed reasoning;
omitted `reasoning_effort` is invalid because the runtime would not receive a
complete tuple. Roundfix never changes an explicit `high` request to empty
reasoning when proof fails.

Manage profiles with:

```bash
roundfix profiles show --category backend
roundfix profiles show --all --json
roundfix profiles configure --scope project --file profiles.yml --dry-run --json
roundfix profiles validate --category review --json
```

`profiles configure` prepares the candidate in memory, validates it, and
proves every distinct Preferred Selection and fallback before confirmation or
write. `--dry-run` performs the same proof, leaves config bytes unchanged, and
reports `changed: false` in JSON. Proof failure, cleanup failure, declined
confirmation, or output failure preserves the target bytes. `profiles
validate` performs the same read-only Exact Agent Selection Proof over the
effective profiles, deduplicates exact tuples, and closes every disposable ACP
Session without creating a Run, worktree, prompt, credential, or
runtime-owned setting.

Agent-starting commands have two supported selection forms: omit `--agent`,
`--model`, and `--reasoning-effort` to use the effective profiles, or provide
all three together for one complete one-Run override. A complete override
replaces only the Preferred Selection for each relevant category and preserves
that category's Fallback Chain. Any partial subset exits `2` before config
load, proof, Session creation, or Run mutation. An explicitly empty
`--reasoning-effort ""` counts as present. When one complete override spans
multiple Task or QA categories, text and JSON output include a cross-category
warning.

The Recommended Profile shown by `profiles show` is dated `2026-09-30`.
Each of the ten categories has a Preferred Selection and Fallback Chain with
a rationale for each selection. JSON schema `roundfix/profiles/v2` lists
these rows in order with `preferred` and `fallback` roles, the selection,
source date, and rationale. Interactive configure prints the same advisory
rows. The Recommended Profile never selects, routes, or writes configuration.

A configured Agent Selection Profile can carry a **Profile Deviation**: a
dated, reasoned record that its difference from the Recommended Profile is
deliberate for the snapshot it names. Put it after `fallbacks`:

```yaml
profiles:
  backend:
    preferred: {runtime: codex, model: gpt-5.6-sol, reasoning_effort: high}
    fallbacks:
      - {runtime: codex, model: gpt-5.5, reasoning_effort: xhigh}
    deviation:
      from: 2026-09-30
      reason: Keep the model validated for this repository
```

- `deviation` must be a mapping containing exactly `from` and `reason`; both
  are required.
- `from` must be a calendar date written `YYYY-MM-DD`, and `reason` must be a
  non-empty string after trimming.

`profiles configure` writes a deviation its fragment carries. Replacing a
profile with a fragment without a deviation removes the old record.
A Roundfix older than this release refuses a configuration that uses it.

Automatic fallback is pre-prompt only. Roundfix records and shows the fallback
notification before activating the next configured tuple. Once
`agent_work_started` is recorded, no prompt, tool, verification, cancellation,
rate-limit, or session-loss failure can start a replacement session.

Legacy configs with `defaults.agent` or top-level `runtimes` still load during
the compatibility window only when the same file has no `profiles` section.
Do not edit runtime-owned settings or credentials for migration; write complete
profiles with `roundfix profiles configure --scope user|project`.

Setup treats configuration as one fail-before-mutation proposal. It resolves
the effective adapters, proposes migration of stale Codex or Claude overrides,
builds User Config and Project Config bytes in memory, and proves every
generated selection before asking to write. `--yes` authorizes offered
changes; `--no-input` reports diagnosis and writes nothing. Declining an
adapter migration or any later write leaves every unauthorized target
unchanged.

The durable contract is [ADR-0055](../history/adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md).
Its environment evidence is preserved in the
[archived Spec 0041 validation record](../history/specs/0041-agent-selection-runtime-readiness/references/validation.md).
[Spec 0036](../history/specs/0036-doctor-skill-readiness/_prd.md) remains ordered after
this profile-aware Doctor readiness work so it can append Repository Skill Set
readiness without duplicating exact proof.

### Cursor runtime

`cursor` is opt-in and runs as `cursor-agent acp` through acpx. A Cursor
selection names the advertised model value verbatim, including its brackets:

```yaml
{runtime: cursor, model: "grok-4-20[thinking=true]", reasoning_effort: ""}
```

A non-empty `reasoning_effort` is refused: Cursor states effort, thinking,
context and speed inside the model value. The value states speed, so choosing
`fast=false` avoids Fast billing when Cursor offers that variant. Copy the
advertised value rather than composing one. `cursor` is not supported in the
legacy `runtimes:` section or `defaults.agent`; name it in profiles.

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
See [Run Database lifecycle](run-database-lifecycle.md) for the owner and
retention rule of every durable table.

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
- Exclusive retry diagnostics:
  `<artifact-dir>/runs/<run-id>/verification/batch-<nnn>-attempt-<1|2>-retry-1.log`

For review commands, explicit `--spec <slug>` wins over trailer discovery;
otherwise Roundfix uses the newest `Roundfix-Spec: <slug>` trailer on the PR
head commit when that Spec folder exists, falling back to the `_reviews` path.

## How Runs use context

Roundfix keeps Agent context bounded by design:

- Agent prompts receive the assigned Work Item contract and bounded guidance.
  Successful Verification output never enters Agent context. On a
  deterministic first-attempt failure, the Daemon releases Verification
  Capacity and sends exactly one Verification Feedback prompt with the failed
  command, wrapped failure, and diagnostic artifact path, then queues and runs
  the final complete Verification sequence. A Temporary Verification Failure
  uses the separate exclusive retry protocol and sends no Agent feedback.
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

## Delivery conflict recovery

`delivery.derived_paths` declares generated paths and the command that rebuilds
each group. It defaults to `[]`. Project Config replaces the User Config list;
an explicit empty list disables those declarations.

```yaml
delivery:
  derived_paths:
    - paths: [generated/catalog.json, generated/cache/]
      regenerate: make generate
```

Each entry needs a non-empty `paths` list and a non-empty `regenerate` command.
Paths must be clean repository-relative files, directories ending in `/`, or
Go `path.Match` patterns over the whole path. Absolute paths and `..` segments
are refused with a `delivery.derived_paths` config error.

An entry can also declare optional `lines` for files that contain source and
regenerated lines:

```yaml
delivery:
  derived_paths:
    - paths: [skills/testdata/owned-skill-versions.json]
      lines:
        paths: [.agents/skills/*/SKILL.md, skills/*/SKILL.md]
        match: '^ *version: '
      regenerate: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
```

`lines` needs a non-empty `paths` list with the same path rules and a non-empty
`match` that compiles as a Go regular expression. An incomplete or invalid
value is refused with a `delivery.derived_paths[<n>].lines` config error.
Whole-file `paths` takes precedence over `lines.paths`.

For a conflicted `lines.paths` file, the owner takes the default branch's side of
each conflict hunk only when its sides have the same number of lines and differ
only on lines matching `match` on both sides; identical lines between them are
kept. Other hunks keep the file a source conflict and park the item. The resolver
uses two-sided merge hunks, regardless of the user's conflict style.
Regeneration may change matching
lines only when they match before and after, with the file's line count kept;
any other change aborts with `regenerated <path> outside delivery.derived_paths`.

`delivery.item_binary` declares the Roundfix binary the Delivery Queue builds
and can run for an item's `implement`, `archive` and `review` steps. It has two
fields, `build` and `path`; Project Config replaces User Config for this value.
For example:

```yaml
delivery:
  item_binary:
    build: make build
    path: bin/roundfix
```

`delivery.item_binary requires build and path` when either field is empty.
`delivery.item_binary has unsafe path "<path>"` when `path` is absolute, uses a
backslash, contains a `..` segment, is `.`, or is not already its clean form.
An unknown nested key is refused with
`delivery.item_binary.<key> is not a supported config key`.

The Delivery Queue reads these declarations from the fetched default-branch
commit, with User Config beneath it. A declaration added or changed only on an
item branch takes effect after that branch merges. During a derived conflict,
the owner takes the default branch's conflicted files and runs each matched
command once in declaration order through the Verification executor. Logs are
under the artifact directory's `delivery/<slug>/derived-regeneration-<n>.log`.
Declare only files that can be regenerated; a regeneration that changes an
undeclared path aborts the merge and parks the item.
