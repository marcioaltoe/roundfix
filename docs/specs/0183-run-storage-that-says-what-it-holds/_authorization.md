---
status: approved
granted: 2026-09-29
action: make a retention prune, a GC report and the operational sweep count only Runs that still hold something to reclaim, add a read-only Doctor storage check, let gc sanitize recognize a pre-key default Artifact Root, and keep runs list silent about Runs whose recorded checkout is gone
consuming: 0183-run-storage-that-says-what-it-holds
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0183

The maintainer approved the Onda 5 plan in chat on 2026-09-29. Answering a
structured question with "Duas filas", they chose to deliver Specs 0181 and
0182 through one `roundfix deliver start`, then rebuild the binary and deliver
Specs 0183 and 0184 through a second queue. This Spec is part of that second
wave. The skill files ride the standing grant of 2026-09-18 for keeping the
shipped skills true to the CLI. The set was measured with `GovernedPath` on the
authoring branch at `fe517d2b`, through a `go test -overlay` probe that wrote
nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the Roundfix
  skill describes what `gc` reports and the operational sweep's stderr line
  (task_01), and the Doctor checks (task_03). `.agents/skills/` is canonical
  and `skills/` its mirror.

## What is not governed

`internal/store/journal.go`, `internal/cli/gc.go`, `internal/cli/doctor.go`,
`internal/cli/health.go`, `internal/config/config.go`,
`internal/worktree/worktree.go`, the new `internal/cli/doctor_storage.go` and
`internal/store/free_bytes.go`, the
existing tests `internal/store/journal_test.go`, `internal/cli/gc_test.go`,
`internal/worktree/worktree_test.go` and `internal/cli/doctor_test.go`, the new
test files `internal/store/prune_reclaimable_test.go`,
`internal/cli/gc_reclaimable_test.go`, `internal/cli/doctor_storage_test.go`,
`internal/cli/gc_sanitize_pre_key_root_test.go`,
`internal/config/artifact_directory_for_path_test.go`,
`internal/worktree/retained_vanished_checkout_test.go` and
`internal/cli/runs_list_vanished_checkout_test.go`,
`docs/user-guide/commands.md` and `CONTEXT.md` are ordinary. `internal/cli/cli_test.go` and
`docs/references/coverage-record.json` are governed and are not touched: no
help text changes and no existing top-level test is renamed or removed.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No change to archived Specs, the Run Database schema, Journal Retention
  eligibility, `roundfix storage`, `gc compact` or the Doctor and GC help text.
- The live Run Database and Artifact Roots under `~/.roundfix` are never opened
  for writing by a Task or the gate; tests use temporary Roundfix Homes.
- No paid API use, provider call from a test, release, tag, deployment or
  branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
