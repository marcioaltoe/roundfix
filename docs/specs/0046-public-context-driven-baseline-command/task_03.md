---
task: task_03
spec: 0046-public-context-driven-baseline-command
status: completed
type: backend
complexity: medium
---

# Task 03: Manage repository-owned Baseline Profiles

## Overview

Expose repository-owned Baseline Profile authoring, inspection, and validation
through the public CLI. A maintainer can create one valid declarative profile
and resolve it against the embedded catalog without confusing it with Agent
Selection Profiles.

## Requirements

1. MUST implement `roundfix baseline profile init`, `show`, and `validate`
   with the TechSpec flag and output contracts.
2. MUST discover custom profiles only under
   `.roundfix/baseline/profiles/<id>.json`.
3. MUST allow custom profiles to compose only embedded modules, decisions,
   Repository Capabilities, and templates.
4. MUST reject remote references, custom assets, executable content, unknown
   IDs, profile composition, unsafe paths, and invalid schemas.
5. MUST keep stdout, stderr, JSON schemas, help, and exit categories aligned
   with the public CLI contract.

## Subtasks

- [x] Add Baseline command dispatch and profile subcommand parsing.
- [x] Implement repository profile discovery and strict loading.
- [x] Implement init, text/JSON show, and validation results.
- [x] Add CLI and filesystem tests for valid and rejected profiles.
- [x] Add help and schema contract tests.

## Acceptance Criteria

- [x] `profile init` creates a valid repository-owned profile from an allowed built-in source.
- [x] `profile show` resolves the same normalized profile in text and JSON.
- [x] `profile validate` accepts an ID or explicit profile path and reports stable failures.
- [x] User-scoped profiles and repository-provided assets never participate.
- [x] The command family does not alter Agent Selection Profile behavior.

## Context

- instruction: `docs/adr/0067-custom-baseline-profiles-are-repository-owned.md`
- interface: `internal/cli/profiles.go`
- interface: `internal/config/profiles.go`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestCustomProfile|TestBaselineProfileCommand|TestBaselineProfileHelp'` — expected: authoring, discovery, validation, output, and rejection cases pass.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline profile --help` — expected: help lists only init, show, and validate with the approved public grammar.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Story 4; Core Features 1, 3–4, 16, 18; Non-Goals / Out of Scope.
- `_techspec.md` → Data Models: CustomProfile; API Contracts: profile operations; Build Order 2 and 7.
- ADR-0067 → repository ownership and embedded-only composition.

## Result

Roundfix now exposes `baseline profile init`, `show`, and `validate` as a
separate command family from Agent Selection Profiles. Repository-owned
profiles use the strict `roundfix/custom-baseline-profile/v1` declaration
under `.roundfix/baseline/profiles/<id>.json`, bind to the embedded catalog
schema, and resolve to deterministic `roundfix/baseline-profile/v1` content
and a domain-separated digest.

The loader accepts only direct regular JSON files in the repository profile
directory. It rejects duplicate or unknown fields, invalid schema and IDs,
built-in ID collisions, unknown or misordered modules, conflicts, unknown
decisions, Repository Capabilities, or templates, invalid decision values,
profile composition, custom assets, executable or remote fields, symlinks,
escaping paths, and filename/ID mismatches.

Acceptance evidence:

- `TestCustomProfileInitAndLoad` created `team-go` from the embedded
  `go-cli-tui` profile, loaded the exact repository file, reproduced its
  normalized digest, and discovered it lexically from the repository-only
  directory.
- `TestBaselineProfileCommandInitShowAndValidate` proved `init`, text and JSON
  `show`, ID/path validation, `roundfix/baseline-result/v1`, stdout/stderr
  separation, and matching normalized digests across output formats.
- `TestCustomProfileRejectsInvalidDeclarations` covered unknown catalog IDs,
  remote references, custom assets, executable content, profile composition,
  duplicate keys, and invalid typed decision values.
- `TestCustomProfileRejectsUnsafePathsAndNonRepositorySources` rejected an
  outside explicit path and an escaping symlink. The CLI rejection suite also
  proved a user-home profile does not participate in discovery.
- `TestBaselineProfileHelpContract` proved the public grammar advertises only
  `init`, `show`, and `validate`, while the existing `roundfix profiles` help
  remains the Agent Selection Profile contract. The full existing CLI suite
  passed unchanged.

Verification:

- Pre-change: the focused Go test failed to compile because
  `InitCustomProfile`, repository discovery/loading, strict parsing, and the
  Baseline Profile types did not exist.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-go-cache go test -count=1
  ./internal/baseline ./internal/cli -run
  'TestCustomProfile|TestBaselineProfileCommand|TestBaselineProfileHelp'`:
  passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-go-cache go run
  -buildvcs=false ./cmd/roundfix baseline profile --help`: passed and listed
  only `init`, `show`, and `validate`.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-go-cache go test -count=1
  ./internal/baseline ./internal/cli`: passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-go-cache make verify`: passed
  1,762 Go tests, both 256-test setup suites, embedded asset validation,
  shipped-skill checks, and the Roundfix build.

The isolated `GOCACHE` keeps Go build artifacts inside the Task Worktree
sandbox. The Daemon remains responsible for the task's verbatim authoritative
Verification commands.
