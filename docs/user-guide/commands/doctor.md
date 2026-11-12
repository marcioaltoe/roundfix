### doctor

```bash
roundfix doctor
```

Read-only readiness report; mutates nothing and exits nonzero when any check
fails. One stdout line per check with `ok`, `warn`, `failed`, or `skipped`; `residue:`
and `storage:` also report `found` or `partial`. Failure and warning lines include
`next: <action>` when a remediation is known. The checks:

When `NODE_OPTIONS` names a preload whose file no longer exists, Roundfix drops
that preload from the environment it gives an agent process. It checks absolute
preload paths from the last `NODE_OPTIONS` entry, keeps existing paths and
package names, and reports each dropped path once per process on standard error:

```text
roundfix: notice: NODE_OPTIONS preload "<path>" does not exist; Roundfix left it out of the agent environment
```

- `node:` — Node.js meets the minimum version.
- `acpx:` — the installed acpx version is at least the minimum supported
  version. Newer versions are accepted and are not downgraded.
- `adapter:` — every distinct runtime referenced by the effective required
  profiles proves its adapter package lineage and supported version. The
  runtime entries are deduplicated and sorted; legacy, unknown, old, and
  missing adapters fail the aggregate line with that runtime's official
  install action. A Cursor runtime is checked with `cursor-agent status`; no
  login fails the line with `cursor_login_required` and
  `next: run cursor-agent login in a terminal yourself`.
- `profiles:` — the required Agent Selection Profiles pass exact proof through
  disposable ACP Sessions. Success names distinct tuples and category
  references. Failure names the exact tuple, affected categories,
  classification, bounded adapter evidence, and the next
  `roundfix profiles configure` or `roundfix profiles validate` action. A
  rejected explicit `high` does not recommend model-managed reasoning.
- `recommendations:` — immediately after `profiles:`, compares the loaded
  configuration with the shipped snapshot. `ok` means no category differs
  (current and pinned profiles both qualify); `found` names the counts and
  suggests `roundfix profiles check`. If comparison cannot run, `skipped`
  carries the reason. This line opens no Agent Session and never fails Doctor.
- `light-tier:` — printed after `recommendations:` only when
  `openrouter.light_models` is not empty and `defaults.agent_full_access` is
  `true`. It builds each light model's runtime as a Run would and reports
  `DR-LIGHT-ACCESS-MODE` (`warn`) when none offers a full-access mode, because
  light Tasks then run on their category profile; otherwise it reports `ok`.
  It opens no Agent Session, reads no key, and never changes the exit code.
- `gh:` — GitHub CLI is on `PATH`, at least 2.81.0, authenticated for the
  delivery forge and able to write to the repository. The active account login
  is reported; tokens and scope lists are never printed.
- `git:` — Git is on `PATH`, at least 2.23.0, with `user.name` and `user.email`
  set for the repository. Missing keys are named; identity values are never
  printed.
- `remote:` — the delivery remote exists, identifies a GitHub repository and
  answers `git ls-remote <remote> HEAD`. The remote is `watch.push_remote`, or
  `origin` when unset. Outside Git, it is `skipped` with
  `requires a Git repository`.
- `toolchain:` — after `remote:`, finds the tools required by
  `defaults.verification`, `worktree.bootstrap`, each
  `delivery.derived_paths[].regenerate`, the Setup Manifest's
  `verification.gate` and `verification.incremental` decisions, and
  `verification.tools`. Discovery checks executable files and symlink targets
  on `PATH` without running tools. `DR-TOOL-MISSING` fails the line and names
  every source requiring the missing tool. Only the first word after leading
  `NAME=value` words is read; a shell keyword, builtin or non-bare name yields
  `DR-TOOL-UNREAD` (`warn`) with the command. List its tools under
  `verification.tools` to make their requirements explicit.
- `environment:` — after `toolchain:`, reports one
  `DR-NODE-PRELOAD-MISSING` (`warn`) per missing absolute preload path in the
  last `NODE_OPTIONS` value, using the agent environment's parser. Existing
  paths and package names produce no finding. It lists `spec judge keys` in
  preference order: `ROUNDFIX_OPENROUTER_JUDGE_API_KEY`,
  `ROUNDFIX_OPENROUTER_API_KEY`, `ROUNDFIX_TYPESAFE_API_KEY`; then
  `implementation keys`: `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY`,
  `ROUNDFIX_OPENROUTER_API_KEY`. Each entry is `set` or `not set`, never a
  key value. Both stages fall back to the shared `ROUNDFIX_OPENROUTER_API_KEY`
  when their stage key is unset; the judge can then use TypeSafe directly.
  Empty values count as unset, and the last entry for a name wins. The generic
  `OPENROUTER_API_KEY` is never read. Remove a missing preload where your shell
  sets `NODE_OPTIONS`; unset optional keys do not fail the line.
- `skills:` — the required Repository Skill Set matches its local
  authorities. The running binary's embedded artifacts are authoritative for
  the 13 Roundfix-owned skills, including the Roundfix Skill. Each required
  external skill must be present: named by `skills-lock.json` and installed under
  `.agents/skills`. The minimum version of an owned skill is the version
  of that skill the running binary carries, so a copy installed by an older
  binary fails this line until it is refreshed.
  `DR-SKILL-LOCK-STALE` warns when an installed owned skill has a `ref` or
  a `computedHash` that differs from its installed tree in `skills-lock.json`.
- `residue:` — live processes from terminal Run lineages, or a partial result
  when Doctor cannot inspect every lineage. A terminal Run whose recorded owner
  PID is alive with a different start identity adds nothing: its owner is gone
  and the PID was reused, so `partial` is left for a process table Doctor
  cannot read.
- `storage:` — whether Run storage is reclaimable from terminal Runs or free
  Run Database pages, with the applicable `roundfix gc` next action. It is
  read-only and never reports failed, so it does not change Doctor's exit
  code.
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
- `codex:` — macOS-only runtime hygiene: inspects `com.apple.quarantine` (the
  real XProtect trigger) and code-signature validity, resolving `CODEX_PATH`
  first and then `codex` on `PATH`. It does not use `spctl --assess`, which
  rejects any signed CLI that is not a notarized app. A quarantined or
  improperly-signed codex fails with the next action to reinstall codex with
  the official curl installer into `~/.local/bin` and set `CODEX_PATH`.
  Skipped on non-Darwin platforms.

Doctor has no separate `agent:` or `model:` authority. The aggregate
`profiles:` result is the Agent Selection Profile Readiness contract.
Repository Skill Set readiness runs after it as an independent check, even
when profile proof fails, and appears before `codex:`.
Outside Git, profile proof uses the process working directory while Repository
Skill Set inspection does not run. The `skills:` line reports
`Repository Skill Set readiness requires a Git repository` and keeps the
run-from-Git next action.

```text
node: ok
acpx: ok
adapter: ok (claude: command="npx -y @agentclientprotocol/claude-agent-acp@0.84.0"; package=@agentclientprotocol/claude-agent-acp; version=0.84.0 | codex: command="npx -y @agentclientprotocol/codex-acp@2.0.1"; package=@agentclientprotocol/codex-acp; version=2.0.1)
profiles: ok (4 distinct tuples; 10 category references)
recommendations: ok (snapshot 2026-09-30; 5 current, 0 differ, 0 pinned)
pre-pr-review: ok (provider=codex; source=default)
gh: ok (gh <version>; github.com login=<login>; <owner/repo> permission=WRITE)
git: ok (<version> >= 2.23.0)
remote: ok (origin: github.com/<owner/repo>; reachable)
toolchain: ok (<tools> found)
environment: ok (spec judge keys: ROUNDFIX_OPENROUTER_JUDGE_API_KEY not set, ROUNDFIX_OPENROUTER_API_KEY not set, ROUNDFIX_TYPESAFE_API_KEY not set; implementation keys: ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY not set, ROUNDFIX_OPENROUTER_API_KEY not set)
skills: ok (<required> required: <owned> Roundfix-owned, <external> external)
residue: ok (no process residue found)
storage: ok (nothing to reclaim; Runs reclaimable: 0; Run Database free bytes: 0)
codex: ok
```

The skills counts come from the repository's Repository Skill Set.

Roundfix checks a required external skill only for presence: an entry in
`skills-lock.json` and a folder under `.agents/skills`. It reads neither the
folder's files nor the entry's `computedHash`, never compares the skill with
its Setup Snapshot, and never replaces a present copy unasked; `roundfix
baseline skills restore --profile <id> --skill <name>` replaces one
deliberately.

Roundfix reads a `skills-lock.json` of version 1 or 2, the version the skills CLI 1.8 writes; every `computedHash` it writes into the lock or compares with an installed folder follows the hash algorithm of the version the file declares, and it never changes that version. Any other lock version fails the `skills:` line with the version named and `roundfix upgrade` as its only next action, because the skills command would write that version again.

With a Setup Manifest, Doctor also checks installed Roundfix-owned lock entries
and appends the exact suffix
`; DR-SKILL-LOCK-STALE: skills update would replace these Roundfix-owned skills with the copy their skills-lock.json entry pins: <names>`,
where `<names>` are sorted and joined by `, `:

```text
skills: warn (<counts>; DR-SKILL-LOCK-STALE: skills update would replace these Roundfix-owned skills with the copy their skills-lock.json entry pins: <names>; next: roundfix skills install --target project)
```

An `ok` line becomes `warn`; `failed` and `unversioned` keep their status.
The warning never changes Doctor's exit code and never names external entries.
It adds `roundfix skills install --target project` once to the next action,
joined to an existing action by ` && `. One bundle install rewrites the
installed owned entries with the installed tree's hash and no `ref`, clearing
the warning. If the scan cannot run, Doctor appends
`; skills-lock.json comparison unavailable: <error>` and keeps the status and
next action. The release plan's `skills` check reads this same line.

When a **Stale User Skill Copy** exists, Doctor prints this line right after `skills:`; it never changes Doctor's exit code, prints no such line otherwise, and the release plan does not read it.

```text
user-skills: warn (user-level Roundfix skill copies older than this binary: <skill> <found> (binary <version>) at <path>[ -> <link target>], ...; next: <refresh> && <refresh>)
```

In the catalog's sync window right after a release, `skills update` over an
unpinned entry can still install the previous copy. Doctor's version check
reports that older installed copy, and one bundle install restores it.

A missing required skill, or an invalid required lock declaration,
prints one sorted blocking line and makes Doctor exit `1`. Doctor still prints
every other readiness result:

```text
skills: failed (missing: agentic-cli-design; next: bunx skills add marcioaltoe/skills@agentic-cli-design)
```

For mixed ownership, Doctor joins the Roundfix-owned restore and external
update actions with `&&`, so external remediation runs only after the owned
restore succeeds. A failure owned only by the external lock or skill set prints
only the external action.
Doctor never runs either command, never deletes skills, and never updates
`skills-lock.json`. The check is offline and read-only: it reads only local
embedded artifacts, `.agents/skills`, and `skills-lock.json`. Unrelated extra
installed skills and lock entries are ignored and are not removed or flagged.

Doctor prints `gh`, `git` and `remote` after `pre-pr-review` and before
`skills`. Only the `gh` and `remote` lines contact the repository's forge.
They make at most three network reads: `gh auth status --active --hostname
<host> --json hosts`, `gh repo view <owner/repo> --json viewerPermission`, and
`git ls-remote <remote> HEAD`. Every `gh` and `git` child has empty standard
input, `GIT_TERMINAL_PROMPT=0`, and a ten-second limit that cancels the child.
The permission read uses the resolved forge host, including GitHub Enterprise.
Doctor remains read-only and mutates nothing; it never runs the next actions.

The accepted remote URL forms are `https://host/owner/repo(.git)`,
`ssh://user@host/owner/repo(.git)`, and `user@host:owner/repo(.git)`.
`github.com` is a forge host; another host must appear in the authentication
JSON. Local paths, `file://` URLs and other forms fail the forge check.

A definite failure exits `1`. A read that times out or cannot connect yields
`warn` with a next action; warnings alone leave the exit code at `0`.
`ADMIN`, `MAINTAIN` and `WRITE` permissions qualify as write access.
Authentication JSON with no host entry fails; an `error` entry naming HTTP 401
fails as a rejected token. Timeout and other authentication errors warn.
Child diagnostics are not forwarded, protecting credentials and identity values.

| Code | Status | Meaning and next action |
| --- | --- | --- |
| `DR-GH-MISSING` | failed | `gh` is absent; install GitHub CLI 2.81.0 or newer from https://cli.github.com. |
| `DR-GH-VERSION` | failed | Version below the floor or unreadable; upgrade GitHub CLI to 2.81.0 or newer. |
| `DR-GH-UNAUTHENTICATED` | failed | No account for the host; `gh auth login --hostname <host>`. |
| `DR-GH-TOKEN-REJECTED` | failed | The forge returned HTTP 401; `gh auth refresh --hostname <host>`. |
| `DR-GH-UNREACHABLE` | warn | Login could not be confirmed; re-run roundfix doctor when `<host>` is reachable. |
| `DR-GH-PERMISSION` | failed | Permission is below write; ask for write access to `<owner/repo>`, or `gh auth switch --hostname <host>`. |
| `DR-GH-PERMISSION-UNVERIFIED` | warn | Permission could not be read; re-run roundfix doctor when `<host>` is reachable. |
| `DR-GIT-MISSING` | failed | Git is absent; install Git 2.23.0 or newer. |
| `DR-GIT-VERSION` | failed | Version below the floor or unreadable; upgrade Git to 2.23.0 or newer. |
| `DR-GIT-IDENTITY` | failed | An identity key is empty; `git config user.name <name>` or `git config user.email <address>`. |
| `DR-REMOTE-MISSING` | failed | The delivery remote is absent; `git remote add <remote> <url>`. |
| `DR-REMOTE-FORGE` | failed | Remote form or host is not a GitHub forge; point `<remote>` at the repository's GitHub URL, or set `watch.push_remote`. |
| `DR-REMOTE-UNREACHABLE` | warn | Remote read did not answer; re-run roundfix doctor when `<host>` is reachable. |
| `DR-LANE-DOCKER-MISSING` | warn | Docker is absent; install Docker, or set `delivery.linux_lane: false`. |
| `DR-LANE-DOCKER-UNREACHABLE` | warn | `docker info` failed; start Docker, then re-run `roundfix doctor`. |
| `DR-LANE-IMAGE-MISSING` | warn | The image is not local; `docker pull <image>`. |
| `DR-LANE-IMAGE-UNKNOWN` | warn | The image cannot be resolved from `go.mod`; set `delivery.linux_lane_image`. |
| `DR-LIGHT-ACCESS-MODE` | warn | No light runtime offers a full-access mode; set `openrouter.light_models: []` in User Config, or `defaults.agent_full_access: false`. |
| `DR-DISK-LOW` | warn | Less than 10,240 MiB free where worktrees, a Shared Build Cache or the temporary directory live; `roundfix gc`, then free space on that filesystem. |
| `DR-DISK-UNREAD` | warn | The free space could not be read; check the path, then re-run `roundfix doctor`. |

An offline machine with otherwise ready checks prints warnings and exits `0`:

```text
gh: warn (gh <version>; DR-GH-UNREACHABLE: could not confirm the github.com login: read timed out; next: re-run roundfix doctor when github.com is reachable)
git: ok (<version> >= 2.23.0)
remote: warn (origin: github.com/<owner/repo>; DR-REMOTE-UNREACHABLE: git ls-remote origin failed: read did not answer; next: re-run roundfix doctor when github.com is reachable)
```

When a newer Roundfix wrote this repository's Baseline, Doctor prints this line first, runs every other check without Project Config, and exits `1`:

```text
baseline-version: failed (Setup Manifest "<path>" records Roundfix <recorded>, newer than this binary's <running>; a newer roundfix wrote this repository's Baseline, so use that binary or run 'roundfix upgrade'; Project Config not read; next: roundfix upgrade)
```

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
