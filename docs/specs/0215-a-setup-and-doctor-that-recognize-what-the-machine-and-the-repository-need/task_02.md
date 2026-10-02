---
task: task_02
spec: 0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need
status: pending
type: backend
complexity: high
---

# Task 02: Doctor reports the toolchain and the environment, and this repository declares what it needs

## Overview

A repository's Verification needs tools that nothing checks, a dead
`NODE_OPTIONS` preload hides behind `node: ok`, and the optional features'
keys exist in one shell and not another. This Task adds the `toolchain` and
`environment` lines of the TechSpec's "The toolchain and environment lines",
the Project Config key `verification.tools`, and this repository's own
declarations of its tools and derived paths. It is verifiable alone through
Doctor with a fake resolver and environment, and through this repository's
Project Config.

## Requirements

1. MUST add `verification.tools` to User and Project Config as "Data Models"
   states: default empty, Project Config replacing the User Config list, each
   entry a bare executable name, and a config error naming
   `verification.tools` for an invalid or duplicate entry.
2. MUST export `ResolveExecutable` from the capability-discovery resolver
   without changing how it resolves, and MUST find every tool through it,
   never by running the tool (ADR-0087).
3. MUST implement the `toolchain` line as the TechSpec states: the configured
   commands it reads, the first-word rule with the shell keyword and builtin
   list, `DR-TOOL-MISSING` naming the tool and every source that needs it,
   and `DR-TOOL-UNREAD` as `warn`.
4. MUST add the exported `MissingNodePreloads` in the new file
   `internal/agent/missing_node_preloads.go`, built on Spec 0211's
   `agentNodeOptions` in the same package, and call it from the `environment`
   line, never a second `NODE_OPTIONS` parser,
   with one `DR-NODE-PRELOAD-MISSING` `warn` per missing path.
5. MUST report each `ROUNDFIX_` key variable the judge's transports name as
   `set` or `not set`, MUST NOT print, store or compare a key value beyond a
   non-empty test, and MUST NOT read `OPENROUTER_API_KEY`,
   `TYPESAFE_API_KEY` or any variable the transports do not name.
6. MUST print `toolchain` and `environment` after `remote` in Doctor, as
   Surface Transcript 1 states.
7. MUST add to `.roundfixrc.yml`: `verification.tools` listing `go`, `gofmt`,
   `git` and `make`; and `delivery.derived_paths` declaring
   `internal/baseline/testdata/catalog.digest`,
   `internal/baseline/testdata/catalog.normalized.json`,
   `internal/baseline/testdata/catalog.diagnostics.golden.json` and
   `internal/baseline/testdata/plan-characterization/*.golden.json` with
   `make baseline-digests`; `docs/agents/setup-context.json` with
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`;
   and `skills/testdata/owned-skill-versions.json` with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   It MUST NOT change any other key of that file.
8. MUST describe `verification.tools` in `docs/user-guide/configuration.md`,
   and the two lines with their codes in `docs/user-guide/commands/doctor.md`.
9. MUST add the tests named in Verification: config validation and the load
   of this repository's `.roundfixrc.yml` in the new file
   `internal/config/verification_tools_test.go`; the two lines, the key
   sentinel and Surface Transcript 1 through `runCLI` in the new file
   `internal/cli/readiness_toolchain_test.go`; the wrapper's reuse of the
   agent splitter in the new file
   `internal/agent/missing_node_preloads_test.go`; and the exported resolver
   never executing a candidate in the new file
   `internal/baseline/resolve_executable_test.go`.

## Subtasks

- [ ] Add and validate `verification.tools`.
- [ ] Export the resolver and the preload wrapper.
- [ ] Implement the `toolchain` and `environment` lines and wire them into Doctor.
- [ ] Declare this repository's tools and derived paths.
- [ ] Describe the key and the lines in the guides.
- [ ] Add the tests.

## Acceptance Criteria

- [ ] A tool missing from `PATH` that a configured command or
      `verification.tools` names fails the line with `DR-TOOL-MISSING` and
      every source; a command starting with `cd` warns with `DR-TOOL-UNREAD`.
- [ ] `NODE_OPTIONS` naming a missing `--require` path yields
      `DR-NODE-PRELOAD-MISSING` with that path; an existing preload and a
      package name yield nothing.
- [ ] With a key variable set to a sentinel and `OPENROUTER_API_KEY` set, the
      output reports the key as `set`, never contains the sentinel, and never
      names `OPENROUTER_API_KEY`.
- [ ] This repository's Project Config loads with the four tools and the three
      derived-path declarations.

## Context

- creates: `internal/cli/readiness_toolchain.go`
- creates: `internal/cli/readiness_toolchain_test.go`
- creates: `internal/config/verification_tools_test.go`
- creates: `internal/agent/missing_node_preloads_test.go`
- creates: `internal/baseline/resolve_executable_test.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/doctor_test.go`
- interface: `internal/config/config.go`
- interface: `internal/baseline/profile_alignment.go`
- creates: `internal/agent/missing_node_preloads.go`
- interface: `.roundfixrc.yml`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/configuration.md`
- instruction: `internal/config/delivery.go`
- instruction: `internal/judge/questions.go`
- instruction: `internal/cli/spec_judge.go`
- instruction: `docs/adr/0211-the-agent-environment-drops-a-node-preload-whose-file-is-missing.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestToolchainReadinessNamesEachMissingToolAndItsSource|TestEnvironmentReadinessNamesAMissingPreloadAndNeverAKeyValue|TestDoctorPrintsTheFiveReadinessLines)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestToolchainReadinessNamesEachMissingToolAndItsSource TestEnvironmentReadinessNamesAMissingPreloadAndNeverAKeyValue TestDoctorPrintsTheFiveReadinessLines; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestVerificationToolsAreBareExecutableNames|TestThisRepositoryDeclaresItsToolsAndDerivedPaths)$" ./internal/config 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestVerificationToolsAreBareExecutableNames TestThisRepositoryDeclaresItsToolsAndDerivedPaths; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestMissingNodePreloadsReusesTheAgentSplitter)$" ./internal/agent 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestMissingNodePreloadsReusesTheAgentSplitter" || { printf 'missing pass: %s\n' TestMissingNodePreloadsReusesTheAgentSplitter >&2; exit 1; }; out="$(go test -count=1 -v -run "^(TestResolveExecutableNeverRunsTheCandidate)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestResolveExecutableNeverRunsTheCandidate" || { printf 'missing pass: %s\n' TestResolveExecutableNeverRunsTheCandidate >&2; exit 1; }` — expected: exit 0; before this Task neither test exists, so the command fails.
- `for phrase in "verification.tools" "bare executable name"; do tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' docs/user-guide/configuration.md "$phrase" >&2; exit 1; }; done; for phrase in "DR-TOOL-MISSING" "DR-TOOL-UNREAD" "DR-NODE-PRELOAD-MISSING"; do tr -s '[:space:]' ' ' < docs/user-guide/commands/doctor.md | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/doctor.md "$phrase" >&2; exit 1; }; done` — expected: exit 0; before this Task neither guide carries these phrases, so the command fails.

## References

- `_prd.md` → User Stories 4, 6, 7; Core Features 4, 5, 8; Success Metrics 1, 4, 6
- `_techspec.md` → Interfaces; Invariant 3; The toolchain and environment lines; Finding codes; Data Models; API Contract 1; API Contract 5; Surface Transcript 1; Testing Approach 1, 3, 6; Build Order 2
- ADR-0220; ADR-0087; ADR-0211; ADR-0192
