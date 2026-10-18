---
task: task_04
spec: 0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers
status: completed
type: test
complexity: low
---

# Task 04: Every Task settlement reports a stale Setup Manifest

## Overview

On 2026-10-18 the v0.88.0 release checks found the Setup Manifest stale,
because Spec 0277 had changed one Makefile line and both `declarationDigest`
values hash the whole Makefile (`_prd.md` → Sources 4). The same drift
followed three earlier Makefile edits. This Task appends one test to an
existing always-relevant documentation contract file. The test builds the
release check's Managed Refresh plan and fails while the plan proposes any
change, so the Task that moves a digest also regenerates the manifest. It is
verifiable on its own by running the test on this tree.

## Requirements

1. MUST answer the operator observation of 2026-10-18 (v0.88.0, release pull
   request #573) cited in `_prd.md` → Sources 4 by implementing API Contract 6
   of `_techspec.md` (ADR-0287).
2. MUST append `TestTheSetupManifestIsCurrent` to
   `internal/docscontract/publicdocs_test.go`. It calls `t.Parallel()` first
   and builds the plan from the repository root through these calls:
   - `baseline.LoadEmbeddedCatalog`;
   - `baseline.ResolveManifestInput`;
   - `baseline.BuildPlan` with `baseline.PreservationModeManagedRefresh` and
     the directories of `PATH`.

   It fails when any call fails, when the input needs new decisions, when no
   plan is returned (naming its state and message), or when the plan holds any
   file change or history move. The data source is the repository's own files
   and the embedded Baseline catalog. The test writes nothing and reaches no
   network.
3. MUST make the failure message name every changed path and the command
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes`.
4. MUST NOT add a new file to `internal/docscontract`, since nine
   always-relevant files are pinned. It MUST NOT add a Repository Contract
   Test, edit the Makefile or a CI workflow, start a process from the test, or
   change any other test.

## Subtasks

- [ ] Append the test to the existing always-relevant file.
- [ ] Run it on this tree and confirm the pin of always-relevant files holds.

## Acceptance Criteria

- [ ] On this tree the test passes, and the nine-file pin still passes.
- [ ] A Managed Refresh plan with any change fails the test with the regeneration command.

## Consumers

None. The test reads the repository and the embedded catalog, and no other
package reads it.

## Context

- instruction: `docs/adr/0287-a-plan-is-checked-for-undeclared-verification-paths-and-every-settlement-checks-the-setup-manifest.md`
- instruction: `docs/adr/0103-an-applied-refresh-republishes-the-setup-manifest.md`
- instruction: `internal/cli/releaseplan_checks.go`
- interface: `internal/docscontract/publicdocs_test.go`

## Verification

- `out="$(go test -count=1 -v -tags docscontract -run '^TestTheSetupManifestIsCurrent$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheSetupManifestIsCurrent ' || { printf 'missing pass: TestTheSetupManifestIsCurrent\n' >&2; exit 1; }; grep -qF -- 'baseline update --repo . --no-skills --yes' internal/docscontract/publicdocs_test.go || { printf 'missing regeneration command in the failure message\n' >&2; exit 1; }; go test -count=1 -run '^TestRepositoryContractTestsDeclareTheirRelevance$' ./internal/verifyselect || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the test does not exist, so the command fails at its pass check. Afterwards the test passes on this tree, its message names the regeneration command, and the nine-file pin and the whole documentation contract package pass.

## References

- `_prd.md` → Sources; Goals; User Story 3; Core Feature 4; Success Metric 4
- `_techspec.md` → API Contract 6; API Contract 7; Invariant 4; Current behavior; Testing Approach; Build Order 3
- ADR-0287; ADR-0103; ADR-0252; ADR-0253; ADR-0259


## Result

Appended `TestTheSetupManifestIsCurrent` to the existing always-relevant
`internal/docscontract/publicdocs_test.go`. Its first statement is
`t.Parallel()`. It loads the embedded catalog, resolves the repository's
manifest input, and builds the Managed Refresh plan with the directories of
`PATH`. Catalog, input and planning errors, new decisions, and a missing plan
all fail with the regeneration command; a missing plan also reports its
state and message. Every file change and both paths of every history move
appear in the stale-manifest diagnostic. The test adds no writes, network
calls or process-starting code.

### Focused evidence

- Before the edit, `rg -n 'TestTheSetupManifestIsCurrent|baseline.BuildPlan'
  internal/docscontract/publicdocs_test.go` returned no matches (exit 1).
- `go test -c -tags docscontract -o
  /private/tmp/roundfix-task04-docscontract.test ./internal/docscontract`:
  exit 0; compiled the package with the appended test.
- From `internal/docscontract`,
  `/private/tmp/roundfix-task04-docscontract.test -test.run
  '^TestTheSetupManifestIsCurrent$' -test.v -test.count=1`: exit 0;
  `TestTheSetupManifestIsCurrent` passed in 0.10 s on this tree. This supports
  the current-tree portion of acceptance criterion 1.
- A Python inspection counted exactly nine `docscontract`-tagged Go files
  and asserted that every one, including `publicdocs_test.go`, still carries
  `//verify:always`: exit 0. No file was added and no relevance directive was
  changed. The existing pin test remains for Daemon Verification.
- A disposable local clone containing this working tree ran the same focused
  binary after appending one newline byte to its Makefile. The test exited 1
  in 0.10 s with `Setup Manifest is stale`,
  `docs/agents/setup-context.json`, and
  `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes`.
  The probe asserted those diagnostics and exited 0, supporting acceptance
  criterion 2. The real Makefile and manifest were untouched. Initial scratch
  probes exposed missing Git metadata and then missing HEAD; copying the
  repository's existing history supplied those planning preconditions without
  creating a commit.

The source diff appends only this test; existing tests are unchanged. The
Task's pre-existing bytes, status, checkboxes and Verification commands are
preserved. Declared Verification, the nine-file pin test and the full
package gate remain Daemon-owned and were not run. No follow-up scope was
identified.
