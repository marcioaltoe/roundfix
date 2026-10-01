## Setup, doctor, and upgrade

Use `roundfix setup [--yes] [--no-input]` to take a machine from fresh to
Run-ready. It checks Node.js and minimum-supported acpx, proves Adapter
Readiness for every distinct runtime referenced by the effective required
profiles, builds the generated Agent Selection Profiles in memory, exact-proves
every distinct tuple, and only then offers acpx local adapter overrides, User
Config, and Project Config writes. Each check prints one
deterministic report line with status `ok`, `installed`, `skipped`,
`offered: declined`, or `failed`. Tested report lines include:

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

