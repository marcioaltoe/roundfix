## Setup, doctor, and upgrade

Use `roundfix setup [--yes] [--no-input]` to take a machine from fresh to
Run-ready. It checks Node.js and minimum-supported acpx, proves Adapter
Readiness for every distinct runtime referenced by the effective required
profiles, builds the generated Agent Selection Profiles in memory, exact-proves
every distinct tuple, and only then offers acpx local adapter overrides, User
Config, and Project Config writes. Each check prints one
deterministic report line with status `ok`, `installed`, `skipped`,
`offered: declined`, `warn`, or `failed`. Tested report lines include:

```text
node: ok
acpx: installed
adapter: ok (claude: command="npx -y @agentclientprotocol/claude-agent-acp@0.84.0"; package=@agentclientprotocol/claude-agent-acp; version=0.84.0 | codex: command="npx -y @agentclientprotocol/codex-acp@2.0.1"; package=@agentclientprotocol/codex-acp; version=2.0.1)
profile readiness: passed
acpx agents override: installed
User Config: installed
Project Config: installed
```

`--yes` accepts every offered install or file change. `--no-input` skips
offers instead of prompting and writes nothing. When acpx is missing or older
than `0.12.0`, setup offers `npm install -g acpx@0.12.0`. Version `0.12.0` and
newer versions are accepted; Setup never downgrades a newer installation.

After `acpx` and before adapter work, Setup prints Doctor's five readiness
lines in order: `gh`, `git`, `remote`, `toolchain`, and `environment`. They
also appear when acpx is unavailable. Each uses `<name>: <status> (<detail>)`;
`failed` and `warn` findings include a stable `DR-` code and `; next: <action>`.
Setup offers no install or change for these lines, including with `--yes`.
Only `failed` makes Setup exit `1` at the end; `warn` alone does not.

- `gh` checks GitHub CLI version, login for the repository's forge, and write
  permission. `DR-GH-UNAUTHENTICATED` points to `gh auth login --hostname <host>`.
- `git` checks Git version and repository `user.name` and `user.email`.
- `remote` checks the delivery remote (`watch.push_remote`, otherwise `origin`)
  and whether it is a reachable GitHub forge remote.
- `toolchain` finds executables named by configured Verification, bootstrap,
  regeneration commands, and `verification.tools`, without running them.
  `DR-TOOL-MISSING` names the missing executable and the command that needs it.
- `environment` names missing `NODE_OPTIONS` preload files and reports optional
  `ROUNDFIX_` key variables as set or not set, never their values.

Forge reads use the user's `gh` and Git credentials, no standard input, disabled
Git terminal prompts, and a ten-second cancellation limit per read. A timeout
or unreachable forge reports `warn`; Roundfix does not log in or change Git or
shell configuration. Run `roundfix doctor` again after taking the next action.

| Code | Status | Next action |
| --- | --- | --- |
| `DR-GH-MISSING` | failed | `install GitHub CLI 2.81.0 or newer from https://cli.github.com` |
| `DR-GH-VERSION` | failed | `upgrade GitHub CLI to 2.81.0 or newer` |
| `DR-GH-UNAUTHENTICATED` | failed | `gh auth login --hostname <host>` |
| `DR-GH-TOKEN-REJECTED` | failed | `gh auth refresh --hostname <host>` |
| `DR-GH-UNREACHABLE` | warn | `re-run roundfix doctor when <host> is reachable` |
| `DR-GH-PERMISSION` | failed | `ask for write access to <owner/repo>, or gh auth switch --hostname <host>` |
| `DR-GH-PERMISSION-UNVERIFIED` | warn | `re-run roundfix doctor when <host> is reachable` |
| `DR-GIT-MISSING` | failed | `install Git 2.23.0 or newer` |
| `DR-GIT-VERSION` | failed | `upgrade Git to 2.23.0 or newer` |
| `DR-GIT-IDENTITY` | failed | `git config user.name <name>` or `git config user.email <address>` |
| `DR-REMOTE-MISSING` | failed | `git remote add <remote> <url>` |
| `DR-REMOTE-FORGE` | failed | `point <remote> at the repository's GitHub URL, or set watch.push_remote` |
| `DR-REMOTE-UNREACHABLE` | warn | `re-run roundfix doctor when <host> is reachable` |
| `DR-TOOL-MISSING` | failed | `install <tool>, or change the command that names it` |
| `DR-TOOL-UNREAD` | warn | `list the tools this command needs under verification.tools in Project Config` |
| `DR-NODE-PRELOAD-MISSING` | warn | `remove the preload from NODE_OPTIONS where your shell sets it` |
| `DR-LANE-DOCKER-MISSING` | warn | Docker is absent; install Docker, or set `delivery.linux_lane: false`. |
| `DR-LANE-DOCKER-UNREACHABLE` | warn | `docker info` failed; start Docker, then re-run `roundfix doctor`. |
| `DR-LANE-IMAGE-MISSING` | warn | The image is not local; `docker pull <image>`. |
| `DR-LANE-IMAGE-UNKNOWN` | warn | The image cannot be resolved from `go.mod`; set `delivery.linux_lane_image`. |
| `DR-LIGHT-ACCESS-MODE` | warn | No light runtime offers a full-access mode; set `openrouter.light_models: []` in User Config, or `defaults.agent_full_access: false`. |
| `DR-DISK-LOW` | warn | Less than 10,240 MiB free where worktrees, a Shared Build Cache or the temporary directory live; `roundfix gc`, then free space on that filesystem. |
| `DR-DISK-UNREAD` | warn | The free space could not be read; check the path, then re-run `roundfix doctor`. |

The supported adapters are official `@agentclientprotocol/codex-acp` version
`2.0.1` or newer and official
`@agentclientprotocol/claude-agent-acp` version `0.84.0` or newer. When Setup
needs explicit commands, it proposes
`npx -y @agentclientprotocol/codex-acp@2.0.1` and
`npx -y @agentclientprotocol/claude-agent-acp@0.84.0`. A bare or stale
override can resolve to a package outside the official lineage. Setup proposes
migration from the failed lineage proof rather than from recognizing a
superseded package by name, proves the replacement, and asks before writing. The
official install actions are
`npm install -g @agentclientprotocol/codex-acp@2.0.1` and
`npm install -g @agentclientprotocol/claude-agent-acp@0.84.0`. Decline,
`--no-input`, failed exact proof, or a later write failure preserves every
unauthorized target. A rejected Sol/high proof never becomes an offer to use
model-managed reasoning.

Use `roundfix doctor` when you only need a read-only readiness report. It runs
the Node.js, minimum-supported acpx, effective adapter check for every distinct
runtime referenced by the effective required profiles, required Agent
Selection Profiles, Repository Skill Set, process residue, storage, the
Machine Verification Capacity and the
conditional `light-tier`, `disk-space` and `linux-lane` checks, and codex runtime hygiene checks and
exits nonzero if any check fails. Runtime entries on the aggregate `adapter:`
line are deduplicated and sorted. Adapter failures name the effective command,
package classification, and official install action. Profile failure names the
exact runtime/model/reasoning tuple, every affected category, bounded adapter
evidence, and the next `roundfix profiles configure` or
`roundfix profiles validate` action. A rejected explicit `high` does not
recommend model-managed reasoning. The command has no flags and mutates
nothing.

A terminal Run whose recorded owner PID is alive with a different start identity
adds nothing: its owner is gone and the PID was reused, so `partial` is left for
a process table Doctor cannot read.

- `light-tier:` — printed after `recommendations:` only when
  `openrouter.light_models` is not empty and `defaults.agent_full_access` is
  `true`; `warn` with `DR-LIGHT-ACCESS-MODE` when no light runtime offers a
  full-access mode, otherwise `ok`. It opens no Agent Session and never changes
  the exit code.
- `disk-space:` — printed after `storage:` only when the filesystem holding
  the worktree location (`worktree.location`), a declared Shared Build Cache
  or the temporary directory has less than the Free-Space Floor of 10 GiB
  (10,240 MiB) free, or its free space cannot be read.
  Each location is named on its own and read at its nearest existing
  directory: `DR-DISK-LOW` gives its free MiB with `next: roundfix gc`, and
  `DR-DISK-UNREAD` gives the read error. Findings are `warn` and never change
  the exit code; with enough free space everywhere the line is absent.
- `verification-capacity:` — the read-only Machine Verification Capacity
  check, printed after the `disk-space:` results and before `linux-lane:`.
  It reports `ok (runs.max_verifying <n>; <held> held, <waiting> waiting)`,
  or `ok (off; <held> held, <waiting> waiting)` at `0`. The detail adds
  `; holding: <holder description>, <step>` per holder, then
  `; waiting <n>: <holder description>, <step>` per waiter numbered from 1,
  each in ticket order. An unreadable entry changes it to `warn`, adds
  `; unreadable entry <path>` and the next action
  `inspect <path> and remove it if no Roundfix process wrote it`.
  A directory read error reports `warn (read Machine Verification Capacity: <error>)`.
  Warnings keep the exit code, and the check writes nothing.
- `verification-cores:` — the read-only Verification Core Budget check, printed
  after `verification-capacity:`. It names the budget, its source and the count
  of active Implement Runs across repositories (0 without a Run Database).
  Its three forms are:
  `ok (<cores> per Verification: <source>; <k> active Implement Runs)`;
  `warn (<cores> per Verification: <source>; <k> active Implement Runs; <c> at once may take <total> of <cpus> logical CPUs; next: lower GOMAXPROCS in the environment or set runs.max_verifying in User Config)`;
  and `warn (<cores> per Verification: <source>; active Implement Runs unknown: <error>; next: inspect the Run Database with roundfix runs list)`.
  The source describes the positive whole-number `GOMAXPROCS` override or half
  the logical CPUs divided by `runs.max_verifying`, then `runs.max_active`
  while capacity is off, or undivided while both are 0; the budget is at least 1.
  `<c>` is the positive `runs.max_verifying`, otherwise positive `runs.max_active`,
  otherwise the active Run count, at least 1. Oversubscription warns only when
  `<total>` (`<c>` times `<cores>`) exceeds the logical CPUs. The check writes
  nothing and never changes the exit code.
- `linux-lane:` — printed after `verification-cores:` and before `codex:` only when the
  checkout's effective `delivery.linux_lane` is `true`. It checks Docker on
  the command environment's PATH with `docker info`, then
  `docker image inspect <image>`, each with a ten-second limit that cancels
  the child. The image is `delivery.linux_lane_image` when configured,
  otherwise `golang:<go directive>` from `go.mod` at `HEAD`, ignoring dirty
  changes. Without a Git repository or root `go.mod`, the line is `skipped`.
  A missing Go version reports `DR-LANE-IMAGE-UNKNOWN`. Findings are `warn`
  and never change the exit code; a healthy setup reports
  `ok (docker reachable; image <image> present)`. No pull, container start,
  removal, or next action is executed, and no environment value is printed.
  Doctor reads this checkout's configuration; the queue reads the fetched
  default's configuration instead.

```text
node: ok
acpx: ok
adapter: ok (claude: command="npx -y @agentclientprotocol/claude-agent-acp@0.84.0"; package=@agentclientprotocol/claude-agent-acp; version=0.84.0 | codex: command="npx -y @agentclientprotocol/codex-acp@2.0.1"; package=@agentclientprotocol/codex-acp; version=2.0.1)
profiles: ok (4 distinct tuples; 10 category references)
skills: ok (<total> required: <owned> Roundfix-owned, <external> external)
codex: ok
```

Roundfix checks a required external skill only for presence: an entry in
`skills-lock.json` and a folder under `.agents/skills`. It reads neither the
folder's files nor the entry's `computedHash`, never compares the skill with
its Setup Snapshot, and never replaces a present copy unasked;
`roundfix baseline skills restore --profile <id> --skill <name>` replaces one
deliberately.

`roundfix skills check` validates the shipped bundle and, inside an adopted
Git repository, fails with exit `1` for each Missing Required Skill, naming it
with `roundfix baseline skills restore --profile <id> --skill <name>`.

`roundfix skills check`, `roundfix doctor` and `roundfix update check` also name each Stale User Skill Copy: they read the five user-level directories in the `skills check` order, follow symbolic links and name a folder once under the first directory that reaches it. `skills check` warns on stderr and keeps its exit code, Doctor adds `user-skills: warn` after `skills:`, and the Update Check adds `User skills:` and `userSkills`; each names the refresh command. None of them changes a copy; run a refresh only when the user asks, because it writes outside the repository and through the link the report names.

The binary ships twelve Roundfix-owned skills: the operational `roundfix`
skill and eleven authorial workflow skills.

After the bundle install, an existing `skills-lock.json` entry naming an owned
skill receives the installed folder's `computedHash` and loses its `ref`; every
other field is kept. No entry or lock is created when none exists, and the
names are reported as `Updated skills-lock.json computedHash: <name>` followed
by `Removed skills-lock.json ref: <name>` for each entry whose `ref` was
removed.

With a Setup Manifest, Doctor appends
`DR-SKILL-LOCK-STALE: skills update would replace these Roundfix-owned skills
with the copy their skills-lock.json entry pins: <names>` when an installed
owned skill has a `ref` or a mismatched `computedHash`. It changes `ok` to
`warn`, keeps `failed` and `unversioned`, never changes Doctor's exit code, and
adds `roundfix skills install --target project` as the next action. A bundle
install rewrites the installed owned entries with the installed hash and no
`ref`, clearing the warning.

## Migration check

The synopsis is `roundfix migrate [--check]`. With `--check`, the migration
check reads the database and never migrates or writes it. A current database
prints `Run Database is at schema version <n>, the version this binary
supports: <path>` on stdout and exits `0`. An older database prints
`roundfix: migrate check: <error>` on stderr and exits `2`; its error names
`roundfix migrate` as the remedy. A newer database prints the same error form
and exits `2`; its error names `roundfix upgrade`. An absent database prints
`No Run Database at <path>; nothing to migrate` on stdout, exits `0`, and
creates nothing. Any other read failure prints
`roundfix: migrate check failed: <error>` on stderr and exits `1`.

The migration check uses these surface transcripts:

```text
$ roundfix migrate --check
stdout:
Run Database is at schema version <n>, the version this binary supports: <path>
stderr:
exit: 0
```

```text
$ roundfix migrate --check
stdout:
stderr:
roundfix: migrate check: Run Database "<path>" has schema version 21, older than the schema version <n> this binary supports; run 'roundfix migrate' to upgrade it
exit: 2
```

```text
$ roundfix migrate --check
stdout:
No Run Database at <path>; nothing to migrate
stderr:
exit: 0
```

## Config initialization

```bash
roundfix init [--scope <project|user>] [--force]
```

Creates Project Config at `<repo>/.roundfixrc.yml` or, with `--scope user`,
User Config at `~/.roundfix/config.yml`. With no scope, the command asks and
defaults to Project Config. It checks whether the destination already exists;
`--force` permits overwriting it.

The User Config template written by Setup and `roundfix init --scope user`
includes the Machine Verification Capacity key:

```yaml
runs:
  max_verifying: 1
```

It defaults to `1`; `0` turns the bound off. Only User Config sets it.

## Config compatibility

Roundfix treats registered removed config keys as migrations, not Preflight
Validation failures. The current deprecated keys are `resolve.concurrent` and
`defaults.model`; each is ignored and prints exactly once per User Config or
Project Config load on stderr:

```text
config: resolve.concurrent is deprecated and ignored; use worktree.concurrency
config: defaults.model is deprecated and ignored; use profiles.<category>.preferred.model
```

Unknown keys that are not in the deprecation registry still fail strict
validation.

## Worktree Build Caches

When `worktree.caches` declares rebuildable directories, `storage:` adds
`; Worktree Build Caches reclaimable: <n> (<bytes> bytes)` before `; next:`.
Doctor measures only declared caches in ended worktrees under the configured
worktree location, never follows a cache symlink or a resolved escape, and
writes nothing. The byte walk stops after two seconds total; an unfinished
walk reports `at least <bytes> bytes` and continues no background measurement.
A cache count greater than zero reports `found` and names `roundfix gc`.
Read errors report `partial`; storage never reports `failed` or changes the
exit code. With no declared cache, the storage line stays byte-identical.

```text
storage: found (Runs reclaimable: 0; Run Database free bytes: 0; Worktree Build Caches reclaimable: 1 (11 bytes); next: roundfix gc)
```

## Update check

`roundfix update check [--refresh] [--format <text|json>]` is the read-only Update Check a session runs before its first task, as the Baseline's agent instructions say. It compares the installed Roundfix with the latest published release, lists the `### ` titles of each newer release's notes, newest first, and summarizes this repository's `roundfix baseline update` preview with the installed binary. The release listing is one `gh api` call cached for 24 hours in `~/.roundfix/update-check.json` and shared by every repository on the machine; `--refresh` reads it now, and a failed lookup is one line ending `the session continues.` and is cached for the same 24 hours. Every report exits `0`; a usage error or an invalid configuration exits `2`. `update.mode` (Project Config over User Config) selects `notify`, the default, or `auto`; the `Next:` line names the `--update` Skill Mode when a release is newer or the Baseline has a plan. JSON uses `roundfix/update-check/v1` with `state` (`update_available`, `current`, `unknown` or `development`), `next` (`apply` or `none`), `releases`, `lookup` and `baseline`.

When a newer Roundfix wrote the repository's Baseline, the check warns on stderr and reports `Baseline: failed:` with the refusal; run `roundfix upgrade` before anything else in that repository. Doctor reports the refusal on its `baseline-version: failed` line. The Run readers and `roundfix stop` warn and continue without Project Config.

## Pending History after upgrade

The upgrade notice names every pending history unit, including a Reduced
History Entry left by an earlier version. Run `roundfix baseline update` to
review removal of retired families and uncited superseded records, conversion
of Legacy Archive Folders, missing successor lineage, per-path ADR refusals,
citation rewrites, `record-links` and `live-links` repairs. The latter prints
`live-links: rewrites <l> link(s) in <n> file(s)` in history text and the same
line prefixed with `- ` in baseline text.
No new reduced entry is written; retired text stays in Git.
