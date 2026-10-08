---
status: approved
granted: 2026-10-08
action: remove terminal Runs past a User Config Run Retention window from the machine-wide Run Database with their dependent rows and artifact directories, keep Active, queue-referenced and worktree-holding Runs, run the Run Retention Sweep at most once a day at Run and Delivery Queue start under a budget and through roundfix gc, compact the database incrementally, report the sweep in gc --dry-run and gc, and describe it in the guides, the glossary and the Roundfix skill
consuming: 0250-a-run-database-that-keeps-only-recent-runs
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/events.md
  - .agents/skills/roundfix/references/runs.md
  - .agents/skills/roundfix/references/storage.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0250

On 2026-10-08 the maintainer decided how the Run Database is bounded. The
question asked how Roundfix should keep the 1.1 GB machine-wide database
small, and the answer was "Automática + gc (Recommended)": keep only the data
of Runs still running, and of any Run a Delivery Queue still references, plus
the last N days, N one of 30, 15 or 7 and 30 by default; run it automatically
at the start of each Run and of the Delivery Queue at most once a day, and
through `roundfix gc`; remove the terminal Runs older than N days with their
rows, events, artifacts and dependent rows; and compact the database. Active
Runs and queue-referenced Runs never go.

The same day the maintainer granted:

- the Governed Paths this Spec declares;
- a standing `qa_override` for an environment-only `partial`;
- after this Spec ships, the operator applies the 30-day window to this
  machine's database. That is not part of any Task or QA row.

On 2026-09-30 the maintainer had granted the skills: "considere autorizado a
ajustar todas as skills se necessário".

The governed set was measured with `GovernedPath` on the authoring branch at
`c73a92e0`, through a `go test -overlay` probe in `internal/speccheck` that
wrote nothing to the repository, against every file the Tasks declare. The
skill's derived files were measured in a disposable clone: an edit to
`.agents/skills/roundfix/references/storage.md`, then `make skills-sync`, the
skill version record and `make baseline-digests`, changed only the two
`storage.md` copies, the two `SKILL.md` copies and
`skills/testdata/owned-skill-versions.json`.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/storage.md` describes `roundfix gc`, the
  Journal Retention prune at Run start and the compaction; it must describe
  Run Retention, the automatic sweep and the new report, or the skill no
  longer matches the shipped CLI.
- `.agents/skills/roundfix/references/runs.md` and
  `.agents/skills/roundfix/references/events.md` describe `runs show` and
  `events`; each gains the Run Retention hint for an unknown Run.
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  skill's two version fields, raised with its content.

## What is not governed

These paths are ordinary:

- the new and changed Go sources and tests under `internal/store`,
  `internal/config` and `internal/cli` that the Tasks declare;
- the reference mirrors `skills/roundfix/references/storage.md`,
  `skills/roundfix/references/runs.md` and
  `skills/roundfix/references/events.md`, and
  `skills/testdata/owned-skill-versions.json`;
- `CONTEXT.md`, the user guides under `docs/user-guide/`, and
  `docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md`.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors and the skill
versions. This declaration records the regeneration that follows the approved
edits and adds no source path.

```yaml
command: make skills-sync
```

`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
records the raised version, which is the next version at record time. No
Baseline module changes, so no module version is recorded and
`make baseline-digests` rewrites nothing. No record or mirror is hand-edited.

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, the CI workflows or
  `.roundfixrc.yml`.
- No change to any Baseline module, to the `### QA settlement` section of any
  skill, or to `internal/cli/cli_test.go`.
- No test, Verification command or QA row reads or writes the real
  `~/.roundfix`; every one uses a temporary home and a fixture Run Database.
  No row reaches the network.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
