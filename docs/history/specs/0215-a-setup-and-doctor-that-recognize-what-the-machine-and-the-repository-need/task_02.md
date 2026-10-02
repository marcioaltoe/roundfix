---
task: task_02
spec: 0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need
status: completed
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

## Result

Implemented this Task's slice for Daemon Verification. Status remains
Daemon-owned; no declared Verification command was run, and no commit, push
or Pull Request was created.

- `verification.tools` defaults empty, validates bare names and duplicates
  with a key-specific error in both config scopes, and replaces rather than
  appends the User Config list. Generated config examples expose the empty
  list. The configuration guide describes the rule and precedence.
- Toolchain discovery uses the exported Baseline resolver, retaining its
  existing symlink, permission and bounded-chain behavior without executing
  candidates. It combines all configured command sources and explicit tools,
  reports one missing-tool finding per executable with every source, and
  warns on unreadable first words. Doctor now appends `toolchain` and
  `environment` after `remote`; the command guide documents their codes.
- Environment diagnosis reuses `agentNodeOptions` through
  `MissingNodePreloads`, reports unique missing preload paths, and uses only
  presence booleans for transport-declared `ROUNDFIX_` keys. Generic provider
  keys and undeclared keys are never selected for inspection or reporting.
- `.roundfixrc.yml` adds exactly the four tools and three derived-path
  declarations from this Task. All previously present keys are unchanged;
  no derived artifact was regenerated.

### Acceptance evidence

| Criterion | Implementation and focused-check evidence |
| --- | --- |
| Missing tools name every source; `cd` warns | `TestToolchainReadinessNamesEachMissingToolAndItsSource` checks aggregation from effective Verification, regeneration, both Setup Manifest decisions and explicit tools, one resolution per tool, and the complete authored builtin list. `cd` yields `DR-TOOL-UNREAD` and the line remains a warning when no tool is missing. |
| Missing preload is named; existing paths and packages yield nothing | `TestEnvironmentReadinessNamesAMissingPreloadAndNeverAKeyValue` checks the missing path, last `NODE_OPTIONS` entry, accepted existing path and package, and unique findings for repeated missing paths. `TestMissingNodePreloadsReusesTheAgentSplitter` compares the wrapper with the agent parser for quotes, file URLs, duplicates, relative paths and an unbalanced quote, plus independently asserts expected missing paths. |
| Sentinel key is reported only as set; generic key is absent | `TestEnvironmentReadinessNamesAMissingPreloadAndNeverAKeyValue` exercises Doctor through `runCLI`, checks the transport-declared key states, and rejects sentinel values, generic key names and undeclared key names from output. It also checks that a final empty key entry overrides an earlier non-empty one. |
| This Project Config loads with four tools and three declarations | `TestThisRepositoryDeclaresItsToolsAndDerivedPaths` loads the actual `.roundfixrc.yml` through `ResolveConfigProposal` and compares the full tool list, derived path lists and regeneration commands. `TestVerificationToolsAreBareExecutableNames` checks invalid entries, duplicates, replacement and explicit empty replacement. |

`TestDoctorPrintsTheFiveReadinessLines` exercises Surface Transcript 1 through
`runCLI`, with no GitHub account, missing Git email, missing `rtk` and a dead
preload. It checks all five lines in order, coded findings, next actions,
stdout/stderr placement and exit `1`. The standalone resolver test compiles an
executable through `testfixture.FixtureBinary` with a side-effect marker and
proves discovery of both it and its
symlink leaves that marker absent, then rejects a non-executable target.

### Focused checks

- Red starting point:
  `GOCACHE=/tmp/roundfix-task02-cache rtk proxy go test ./internal/config ./internal/agent ./internal/baseline -run 'Test(VerificationTools|MissingNodePreloads|ResolveExecutable)'`
  exited `1` because `Verification.Tools`, `MissingNodePreloads` and
  `ResolveExecutable` were not yet defined.
- At the initial handback:
  `GOCACHE=/tmp/roundfix-task02-cache rtk proxy go test ./internal/config ./internal/agent ./internal/baseline ./internal/cli -run 'Test(Verification|ThisRepository|MissingNode|ResolveExecutable|Toolchain|EnvironmentReadiness|Doctor|RunDoctor|Forge|Readiness|Init)'`
  exited `0` for all four packages. This includes the new acceptance tests,
  existing Doctor and forge checks, config Verification checks and config
  initialization checks. All external reads use local fixtures or injected
  runners; no live network probe was run.
- `gofmt` was applied to all changed Go files. `git diff --check` exited `0`.
  Diff review confirms only `.roundfixrc.yml` changes among Governed Paths;
  `_tasks.md`, other Task files and generated artifacts are untouched.

Declared Verification and Task settlement remain for the Daemon. No follow-up
work was added to this slice.

### Verification Feedback repair — attempt 1

Inspected the Daemon diagnostic artifact
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261002T124121Z_c56fb26225ec8dc1/verification/batch-002-attempt-1.log`.
The recorded `make verify-changed` attempt reported two failures caused by
this slice: the generated config inserted the new tools list before the
existing Verification Capacity block, violating its tested layout; and the
resolver test directly wrote an executable instead of using the repository's
compiled-fixture mechanism.

Moved the generated tools list after the existing Verification settings,
preserving that block's contract, and added an assertion in this Task's config
test that generated config still exposes `tools: []`. Replaced the resolver
test's direct executable write with `testfixture.FixtureBinary`; the test
retains its marker assertion, symlink resolution and non-executable-target
rejection. Existing contract tests and the frozen executable inventory were
left unchanged.

- Before repair,
  `GOCACHE=/tmp/roundfix-task02-cache rtk proxy go test ./internal/config ./internal/testfixture -run 'Test(DefaultConfigYAMLVerificationCapacity|NoTestWritesAnExecutableOutsideTheResidue)'`
  exited `1`, reproducing both reported failures in isolation.
- After repair,
  `GOCACHE=/tmp/roundfix-task02-cache rtk proxy go test -count=1 ./internal/config ./internal/baseline ./internal/testfixture -run 'Test(DefaultConfigYAML|VerificationTools|ThisRepositoryDeclares|ResolveExecutable|NoTestWritesAnExecutableOutsideTheResidue)'`
  exited `0` for all three packages. This exercises generated-config contracts,
  tools validation and repository declarations, the non-executing resolver
  probe, and the repository-wide executable-write guard.
- Applied `gofmt` to the three repaired Go files; `git diff --check` exited
  `0`. Repairs remain inside this Task's paths.

Task status remains Daemon-owned. Neither `make verify-changed` nor any
declared Task Verification command was rerun by this Agent; no commit, push
or Pull Request was created. The Daemon owns the full Verification rerun.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/readiness_forge.go`

## Carry-forward provenance

- Source Run: `run_20261002T124121Z_c56fb26225ec8dc1`
- Source commit: `c1c1b68d0df774f9408e6d0f6e2e02e47ee6c177`
