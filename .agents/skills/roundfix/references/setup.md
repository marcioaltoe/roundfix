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
Selection Profiles, Repository Skill Set, and codex runtime hygiene checks and
exits nonzero if any check fails. Runtime entries on the aggregate `adapter:`
line are deduplicated and sorted. Adapter failures name the effective command,
package classification, and official install action. Profile failure names the
exact runtime/model/reasoning tuple, every affected category, bounded adapter
evidence, and the next `roundfix profiles configure` or
`roundfix profiles validate` action. A rejected explicit `high` does not
recommend model-managed reasoning. The command has no flags and mutates
nothing.

```text
node: ok
acpx: ok
adapter: ok (claude: command="npx -y @agentclientprotocol/claude-agent-acp@0.84.0"; package=@agentclientprotocol/claude-agent-acp; version=0.84.0 | codex: command="npx -y @agentclientprotocol/codex-acp@2.0.1"; package=@agentclientprotocol/codex-acp; version=2.0.1)
profiles: ok (4 distinct tuples; 10 category references)
skills: ok (<total> required: <owned> Roundfix-owned, <external> external)
codex: ok
```

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
