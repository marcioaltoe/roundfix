### doctor

```bash
roundfix doctor
```

Read-only readiness report; mutates nothing and exits nonzero when any check
fails. One stdout line per check with `ok`, `failed`, or `skipped`; `residue:`
and `storage:` also report `found` or `partial`. Failure lines include
`next: <action>` when a remediation is known. The checks:

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
adapter: ok (claude: command="npx -y @agentclientprotocol/claude-agent-acp@0.63.0"; package=@agentclientprotocol/claude-agent-acp; version=0.63.0 | codex: command="npx -y @agentclientprotocol/codex-acp@1.1.5"; package=@agentclientprotocol/codex-acp; version=1.1.5)
profiles: ok (3 distinct tuples; 10 category references)
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

