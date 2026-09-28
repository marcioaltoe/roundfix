---
status: approved
granted: 2026-09-28
action: make a schema-version refusal name the action that resolves it and ship roundfix migrate, let a help token mean help only where it is an argument, bound the QA Report front matter as the derived QA Verification does, and govern every bounded path while running the repository-contract tests in a repository gate
consuming: 0174-operator-surfaces-that-tell-the-truth
paths:
  - internal/cli/cli_test.go
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - skills/qa-gate/SKILL.md
  - internal/speccheck/governed.go
  - internal/speccheck/governed_repocontract_test.go
  - Makefile
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0174

The maintainer approved Onda 2 of the efficiency sequence in chat on 2026-09-28
("Siga para onda 2 como sugerido até o release"); this Spec is part of it. The
Go sources and Go tests ride the standing grant of 2026-09-21 for governed
source a slice genuinely needs; the `Makefile` change that makes `repo-test` run
the `repocontract` tests was expressly authorized by the maintainer in chat on
2026-09-28 (the standing grant excludes Verification configuration); the skill files ride the standing grant
of 2026-09-18 for keeping the shipped skills true to the CLI. The set was
measured with `GovernedPath` on `0160f70a`.

## Why each governed path is unavoidable

- `internal/cli/cli_test.go` — the root and per-command help tests live there:
  task_01 adds the `migrate` help row and the root usage entry, and task_02 adds
  the help-token tests beside them.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the Run
  discovery section tells a Supervisor what to do when `runs list` refuses a
  Run Database of another schema version (task_01); `.agents/skills/` is
  canonical and `skills/` its mirror.
- `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md` — the qa-gate
  skill tells the QA Agent to keep the seeded front matter the report's only
  one, and that a report whose front matter is empty or duplicated never
  settles the gate (task_03).
- `internal/speccheck/governed.go` — `GovernedPath` must match
  `skills/_ownership.yml`, which archived Spec 0163 bounds (task_04).
- `internal/speccheck/governed_repocontract_test.go` — `TestGovernedSetOnlyGrows`
  lists the newly governed path (task_04).
- `Makefile` — `repo-test` must run every `repocontract` test in the repository
  so `make verify-docs` enforces them (task_04).

## What is not governed

`internal/store/store.go`, `internal/store/store_test.go` and the new
`internal/store/migrate_test.go`; `internal/cli/cli.go`, the new
`internal/cli/migrate.go` and `internal/cli/migrate_test.go`,
`internal/cli/qa_report_test.go` and `internal/cli/archive_test.go`;
`internal/spec/qa.go` and the new `internal/spec/qa_frontmatter_test.go`;
`internal/daemon/task_engine.go`, `internal/daemon/task_engine_test.go` and the
new `internal/daemon/qa_frontmatter_test.go`; `internal/speccheck/governed_test.go`;
`internal/suiteguardcontract/contract.go` and the new
`internal/suiteguardcontract/repository_gate_test.go`; the new `main_test.go`
and `suiteguard_repocontract_test.go` files in `internal/authorization` and
`internal/verifyselect`; `docs/user-guide/commands.md` and `CONTEXT.md` are
ordinary.

## Sanctioned regeneration

The repository-owned commands resolve their generated outputs; these
declarations record the regeneration that follows the approved skill edits and
add no source paths.

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

## Limits

- No change to the derived QA Verification command, to archived Specs, to
  archived QA Reports, to `skills/_ownership.yml` itself or to the CI workflow.
- No change to the signatures of `store.Open`, `store.OpenReader` or
  `store.OpenStorageReader`, and no renamed or removed top-level test.
- The live Run Database under `~/.roundfix` is never opened by a Task or the
  gate.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
