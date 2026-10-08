---
task: task_02
spec: 0251-skills-keep-up-with-the-behavior-they-describe
status: completed
type: backend
complexity: high
---

# Task 02: The repository keeps a Skill Coverage Map and a current Behavior Surface Record

## Overview

Gives this repository its Skill Coverage Map and its generated Behavior
Surface Record, and the repository contract that computes every Behavior
Surface and keeps both current. The record becomes a declared derived path,
so two items that change a surface merge by regeneration. The slice is
verifiable on its own: the contract passes on the repository and fails, in
fixture cases, on a stale fingerprint and on an unmapped surface.

## Requirements

1. MUST create `internal/docscontract/skill_coverage_test.go` with build tag
   `docscontract`, a `//verify:always` header directive and a
   `-record-skill-coverage` test flag. It MUST compute every Behavior Surface
   and its fingerprint bytes exactly as `_techspec.md` → Data Model 3 states,
   reusing the package's `commandPaths` reader and `skillcoverage.Fingerprint`,
   with `cli.Run` in process and no spawned process.
2. MUST implement `TestTheSkillCoverageMapIsCurrent` exactly as
   `_techspec.md` → API Contract 2 states. With the flag it writes the record
   with `skillcoverage.EncodeRecord` to `skillcoverage.RecordPath` before
   checking; without it, it never writes. Every failure message MUST name each
   offending surface id.
3. MUST create `docs/references/skill-coverage.json`, a map that satisfies
   the contract and `_techspec.md` → Data Model 4: every command surface lists
   the reference the Roundfix skill's reference index assigns to its first
   word; every guide surface lists its own path among its sources and the
   owned skill files that describe it, or an `uncovered` reason; `config keys`
   and `exit codes` are covered by the Roundfix skill files that state them.
   No `command` entry lists `internal/cli/cli.go` as a source.
4. MUST create `docs/references/behavior-surfaces.json` only by running
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
5. MUST add these fixture tests to the same file, each feeding in-memory
   surfaces, records and maps to the helpers the contract uses:
   - `TestAStaleBehaviorSurfaceIsReported`: a record with one changed
     fingerprint fails naming that id, and a record missing one id fails
     naming it.
   - `TestAnUnmappedBehaviorSurfaceIsReported`: a computed surface without an
     entry, and an entry without a computed surface, each fail naming the id.
   - `TestCommandCoverageFollowsTheRoundfixReferenceIndex`: a command entry
     that omits its indexed reference fails naming the command and the
     reference.
6. MUST implement `_techspec.md` → API Contract 3: append the derived-path
   declaration of the record to `.roundfixrc.yml` and add both paths to the
   `//verify:relevant` directive of `internal/config/regeneration_declared_test.go`,
   changing no other line of either file.
7. MUST write the glossary terms **Behavior Surface**, **Skill Coverage Map**,
   **Behavior Surface Record** and **Coverage Review** into `CONTEXT.md`
   through the `domain-modeling` skill, each citing ADR-0256.

## Subtasks

- [ ] Write the contract, its record flag and its three fixture tests.
- [ ] Author the map and record the record with the record command.
- [ ] Declare the record as a derived path and widen the regeneration directive.
- [ ] Add the four glossary terms.

## Acceptance Criteria

- [ ] The contract passes on the repository and fails naming the surface on a stale fingerprint, a missing entry, an extra entry and a missing indexed reference.
- [ ] The record is written only by the record command, and re-running it changes no byte.
- [ ] The record's derived declaration selects `TestRegenerationIsDeclared`.
- [ ] `CONTEXT.md` defines the four terms.

## Context

- creates: `internal/docscontract/skill_coverage_test.go`
- creates: `docs/references/skill-coverage.json`
- creates: `docs/references/behavior-surfaces.json`
- interface: `.roundfixrc.yml`
- interface: `internal/config/regeneration_declared_test.go`
- interface: `CONTEXT.md`
- instruction: `internal/docscontract/command_documentation_test.go`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/testdata/owned-skill-versions.json`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md`
- instruction: `docs/adr/0253-every-repository-contract-test-runs-on-main-and-before-a-release.md`

## Verification

- `out="$(go test -count=1 -tags docscontract -v -run '^(TestTheSkillCoverageMapIsCurrent|TestAStaleBehaviorSurfaceIsReported|TestAnUnmappedBehaviorSurfaceIsReported|TestCommandCoverageFollowsTheRoundfixReferenceIndex)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestTheSkillCoverageMapIsCurrent TestAStaleBehaviorSurfaceIsReported TestAnUnmappedBehaviorSurfaceIsReported TestCommandCoverageFollowsTheRoundfixReferenceIndex; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done; sel="$(go test -count=1 -v -run '^TestRegenerationContractRunsWhenADerivedInputChanges$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$sel"; exit 1; }; printf '%s\n' "$sel" | grep -q -- '--- PASS: TestRegenerationContractRunsWhenADerivedInputChanges/declaration_[0-9]*/paths/docs/references/behavior-surfaces.json' || { printf 'the record is not a declared derived path that selects the regeneration contract\n%s\n' "$sel" >&2; exit 1; }; norm="$(tr -s '[:space:]' ' ' < CONTEXT.md)"; printf '%s\n' "$norm" | grep -qF -- '**Behavior Surface**:' || { printf 'missing glossary term Behavior Surface\n' >&2; exit 1; }; printf '%s\n' "$norm" | grep -qF -- '**Skill Coverage Map**:' || { printf 'missing glossary term Skill Coverage Map\n' >&2; exit 1; }; printf '%s\n' "$norm" | grep -qF -- '**Behavior Surface Record**:' || { printf 'missing glossary term Behavior Surface Record\n' >&2; exit 1; }; printf '%s\n' "$norm" | grep -qF -- '**Coverage Review**:' || { printf 'missing glossary term Coverage Review\n' >&2; exit 1; }` — expected: exit 0. Before this Task the four tests, the declaration and the terms do not exist, so the first check fails. After it, the contract and its fixtures pass, the record's declaration selects the regeneration contract, and the glossary defines the four terms.

## References

- `_prd.md` → Core Feature 1; Core Feature 2; Core Feature 3; User Story 4; Success Metric 4; Glossary
- `_techspec.md` → Data Model 3; Data Model 4; API Contract 2; API Contract 3; Invariant 1; Invariant 2; Build Order 2
- ADR-0256; ADR-0253; ADR-0252; ADR-0192; ADR-0187

## Result

Implemented the Task 02 slice for Daemon Verification. Task status and all
Subtask and Acceptance Criteria checkboxes remain Daemon-owned; no terminal
verdict is recorded here.

The new `docscontract` contract computes 96 Behavior Surfaces: 56 command
paths, 38 guide pages, configuration keys and exit codes. Command help runs
in process and falls back to the nearest successful ancestor; YAML tags are
parsed from direct non-test configuration sources, exit constants are evaluated
from their declarations, and guide fingerprints retain the original bytes.
The authored map links commands to the Roundfix reference index, guides to
owned skill files and their own paths, and configuration and exits to the
Roundfix references that describe them. Commands whose handlers remain in the
shared `internal/cli/cli.go` use their guide source without declaring that
shared file, as required by the Data Model.

Acceptance evidence:

- Repository contract and diagnostic fixtures: the required record command
  exited 0 after validating the repository map and generated record. The
  in-memory fixtures passed for a stale fingerprint, a missing recorded id,
  an extra recorded id, a missing map entry, an extra map entry and a missing
  indexed reference, with exact diagnostics naming the offending id and,
  where applicable, the reference. Additional focused fixtures passed for
  owned-file restrictions, unmatched sources, YAML tag extraction, evaluated
  constant values, exact guide bytes and ancestor-help fallback.
- Record ownership and determinism: the only record write in the contract is
  inside the `-record-skill-coverage` branch, before validation, and uses
  `skillcoverage.EncodeRecord`. The record was created and regenerated only
  by the required command. Comparing the saved record with regenerated bytes
  using `cmp` exited 0.
- Derived-path selection: the focused
  `declaration_6/paths/docs/references/behavior-surfaces.json` subtest passed
  and proved selection of `TestRegenerationIsDeclared`. The configuration diff
  adds only the two-line declaration; the regeneration test diff changes only
  its relevance directive to add both paths.
- Glossary: `CONTEXT.md` now defines Behavior Surface, Skill Coverage Map,
  Behavior Surface Record and Coverage Review using the repository glossary
  format, with ADR-0256 linked from every definition.

Focused commands and outcomes (each Go command used
`GOCACHE=/private/tmp/roundfix-task02-gocache`):

- `rtk proxy go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  — passed. An earlier run exposed an incorrect storage-guide source in the
  authored map; the source was corrected to `commands/storage-report.md` and
  the actual `internal/cli/gc.go` handler. A later run exposed unrelated typed
  constants when evaluating all CLI constant declarations; evaluation now
  selects the exit-code groups, and the final record run passed.
- `rtk proxy go test -count=1 -tags docscontract ./internal/docscontract -run '^Test(SkillCoveragePaths|BehaviorSurfaceBytes|AStaleBehaviorSurface|AnUnmappedBehaviorSurface|CommandCoverageFollows)' -v`
  — passed all five selected tests and their fixture cases after the final
  contract edit.
- `rtk proxy go test -count=1 -v ./internal/verifyselect -run '^TestRegenerationContractRunsWhenADerivedInputChanges$/^declaration_6$/^paths$/^docs$/^references$/^behavior-surfaces.json$'`
  — passed the new derived-path selection case.
- `rtk proxy cmp /private/tmp/roundfix-task02-record-before.json docs/references/behavior-surfaces.json`
  — exited 0; regenerated bytes are identical.

Final scope inspection found exactly the seven Task 02 paths, map/record
agreement on all 96 ids, all four ADR-linked glossary definitions and the
preserved status. `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

The authored `## Verification` command was not run. No commit, push, pull
request, Task Graph edit or other Task edit was performed. No follow-up work
was added to this slice.

### Attempt 1 Verification Feedback repair

Inspected the Daemon diagnostic artifact for `make verify-changed`. Two
repository expectation tests still described the pre-Task-02 contract:
the derived declaration list omitted the Behavior Surface Record, and the
contract discovery test expected the old relevance list and eight
`docscontract` files instead of nine.

Updated only the expectations affected by this Task in
`internal/config/verification_tools_test.go` and
`internal/verifyselect/contracts_repository_test.go`. The derived declaration
assertion now includes the record and its exact regeneration command; the
relevance assertion includes both coverage paths, and the file-count assertion
includes the new contract file. Path expectations use `skillcoverage.RecordPath`
and `skillcoverage.MapPath`. Exact equality assertions remain in place.
These two ordinary test paths are additional Task 02 repair paths.

Focused repair evidence:

- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk proxy go test -count=1 ./internal/config ./internal/verifyselect -run '^(TestThisRepositoryDeclaresItsToolsAndDerivedPaths|TestRepositoryContractTestsDeclareTheirRelevance)$' -v`
  — exited 0; both previously failing tests and all three selection scenarios
  passed.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exited 0.
- Reviewed the repair diff: only the two tests' imports and affected
  expectations changed. The repair changes no computed Behavior Surface, so
  the generated record needs no regeneration.

Task status is preserved. Neither `make verify-changed` nor the authored
Verification command was rerun. No commit, push, pull request, Task Graph edit
or other Task edit was performed. The implementation is handed back for the
Daemon's configured Verification retry, with no terminal Task verdict.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/config/verification_tools_test.go`
- `internal/verifyselect/contracts_repository_test.go`
