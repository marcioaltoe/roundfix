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
  install action.
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
- `skills:` — the required Repository Skill Set matches its local
  authorities. The running binary's embedded artifacts are authoritative for
  the 14 Roundfix-owned skills, including the Roundfix Skill. Each of the 25
  required external skills must hash to its `computedHash` in
  `skills-lock.json`. The minimum version of an owned skill is the version
  of that skill the running binary carries, so a copy installed by an older
  binary fails this line until it is refreshed.
- `residue:` — live processes from terminal Run lineages, or a partial result
  when Doctor cannot inspect every lineage.
- `storage:` — whether Run storage is reclaimable from terminal Runs or free
  Run Database pages, with the applicable `roundfix gc` next action. It is
  read-only and never reports failed, so it does not change Doctor's exit
  code.
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
skills: ok (<required> required: <owned> Roundfix-owned, <external> external)
residue: ok (no process residue found)
storage: ok (nothing to reclaim; Runs reclaimable: 0; Run Database free bytes: 0)
codex: ok
```

The skills counts come from the repository's Repository Skill Set.

A missing or outdated required skill, or an invalid required lock declaration,
prints one sorted blocking line and makes Doctor exit `1`. Doctor still prints
every other readiness result:

```text
skills: failed (missing: handoff; outdated: roundfix; next: roundfix skills install --target project && bunx skills experimental_install && bunx skills update -p -y)
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

An offline machine with otherwise ready checks prints warnings and exits `0`:

```text
gh: warn (gh <version>; DR-GH-UNREACHABLE: could not confirm the github.com login: read timed out; next: re-run roundfix doctor when github.com is reachable)
git: ok (<version> >= 2.23.0)
remote: warn (origin: github.com/<owner/repo>; DR-REMOTE-UNREACHABLE: git ls-remote origin failed: read did not answer; next: re-run roundfix doctor when github.com is reachable)
```
