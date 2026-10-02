---
task: task_04
spec: 0207-a-go-cli-and-a-typescript-monorepo-in-one-baseline
status: completed
type: backend
complexity: medium
---

# Task 04: The root gate reaches each toolchain's Verification

## Overview

The composed profile's repository gate is one root `make verify` that must run both the Go and the Bun workspace Verification. task_03 marked those two expectations as parts of the gate. This Task validates the marker and has profile alignment report, without blocking, each part a root Make gate does not reach through its prerequisites or recipes.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST leave `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json` as task_03 left it; its `verification.go` and `verification.workspace` entries carry `"partOfGate": true`.
2. MUST report `catalog.profile.verification.invalid` from `internal/baseline/catalog_validate.go` for a profile verification entry whose `partOfGate` is present and not a boolean.
3. MUST create `internal/baseline/verification_gate_parts.go` with `makeTargetReach` and `gatePartReached` as the TechSpec's Interfaces give, and call them from `resolveVerificationProjection` in `internal/baseline/profile_alignment.go` only when the selected `verification.gate` is `make <target>` declared in a Makefile, appending one non-blocking, recommended `verification.gate.part.missing` divergence per unreached part with the message and next action the TechSpec gives.
4. MUST create `internal/baseline/verification_gate_parts_test.go` with the six tests the TechSpec's Testing Approach 4 names, using temporary repositories only. The skip test MUST show one reported divergence naming the unreached part; the recipe test MUST reach a part only through `$(MAKE) <target>`; the cycle test MUST terminate on two targets that name each other.
5. MUST update the Profiles section of `docs/user-guide/context-driven-development.md` as Fixed texts gives for task_04.
6. MUST run `make baseline-digests`, then `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` twice; the second MUST report `File changes: 0`.
7. MUST NOT change the Standard TypeScript Monorepo Profile's verification entries or any Verification role mapping behavior, and MUST NOT rename or remove a top-level test or an exported function.
8. MUST keep `TestGuidanceCompositionJourney` passing unchanged. Since task_03 it runs a journey for the composed profile whose fixture `verify` target reaches neither part, so that plan now carries two non-blocking `verification.gate.part.missing` divergences and still applies and converges.

## Subtasks

- [ ] Validate the gate-part marker.
- [ ] Read the root Make target's reach and report unreached parts.
- [ ] Write the six tests and update the public guide.
- [ ] Regenerate and refresh this repository twice.

## Acceptance Criteria

- [ ] A root `verify` target whose prerequisites or recipes reach both `make verify-go` and `bun run verify` yields no part divergence.
- [ ] A root gate that reaches only one part yields exactly one non-blocking `verification.gate.part.missing` naming the other.
- [ ] A gate that is not a Make target is not checked, and a cyclic Makefile terminates.
- [ ] A non-boolean `partOfGate` is refused.

## Context

- instruction: `docs/adr/0204-a-composed-profile-takes-a-setup-composed-from-upstream-setups-by-name.md`
- interface: `internal/baseline/catalog_validate.go`
- interface: `internal/baseline/profile_alignment.go`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/agents/setup-context.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/unsatisfied-blocking-capabilities.golden.json`
- creates: `internal/baseline/verification_gate_parts.go`
- creates: `internal/baseline/verification_gate_parts_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAGateThatReachesBothPartsReportsNoDivergence|TestAGateThatSkipsAPartReportsItOnce|TestAPartReachedThroughARecipeInvocationCounts|TestAGateThatIsNotAMakeTargetIsNotChecked|TestMakeTargetReachStopsAtACycle|TestANonBooleanPartOfGateIsRefused|TestCatalogCompatibility)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAGateThatReachesBothPartsReportsNoDivergence TestAGateThatSkipsAPartReportsItOnce TestAPartReachedThroughARecipeInvocationCounts TestAGateThatIsNotAMakeTargetIsNotChecked TestMakeTargetReachStopsAtACycle TestANonBooleanPartOfGateIsRefused TestCatalogCompatibility; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && grep -qF 'verification.gate.part.missing' docs/user-guide/context-driven-development.md && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the six new tests and the guide paragraph do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Interfaces; Data Models; Fixed texts; Testing Approach 4; Testing Approach 5; Build Order 4; Risks & Considerations
- `_prd.md` → Goal 3; Story 4; Core Feature 4; Success Metric 4; Success Metric 7
- `_techspec.md` → API Contract 6
- ADR-0081, ADR-0149, ADR-0204

## Result

Implemented the bounded gate-part check. Catalog validation refuses a present
non-boolean `partOfGate`. Profile alignment follows literal prerequisites and
recursive Make recipes in the selected gate's declared Makefile, visits each
target once, and reports each unreached marked part as a recommended,
non-blocking `verification.gate.part.missing` divergence. It checks resolved
projection commands, preserving existing Verification role mappings. The
Profiles guide now states both expected commands and the advisory diagnostic.

Acceptance evidence from focused checks:

| Criterion | Evidence |
| --- | --- |
| Gate reaches both toolchains | `TestAGateThatReachesBothPartsReportsNoDivergence` passes for direct and transitive prerequisites, an order-only prerequisite, an inline recipe, and RTK-prefixed recipes. `TestAPartReachedThroughARecipeInvocationCounts` passes for `$(MAKE)`, `${MAKE}`, literal `make`, chained invocations, and a mapped Go role. |
| Gate skips exactly one part | `TestAGateThatSkipsAPartReportsItOnce` passes, asserting exactly one recommended, non-blocking diagnostic naming `verification.workspace`, with the authored message and next action. An unrelated target's recipe does not satisfy the part. The recipe test also proves that an echo or a non-MAKE variable does not reach the Go part. |
| Non-Make gates skipped; cycles terminate | `TestAGateThatIsNotAMakeTargetIsNotChecked` passes for a Bun script and an undeclared Make target. `TestMakeTargetReachStopsAtACycle` passes on two targets naming each other, returning exactly two targets and two recipes. |
| Non-boolean marker refused | `TestANonBooleanPartOfGateIsRefused` passes for a string, null, number, and object, each producing `catalog.profile.verification.invalid`. |

Focused commands and outcomes:

- Initial focused build failed on the two absent helper functions, establishing
  the pre-implementation signal.
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk proxy go test ./internal/baseline -count=1 -v -run 'TestAGate|TestAPartReached|TestMakeTargetReach|TestANonBoolean|TestTheComposedProfilePlanConverges'`
  — exit 0; all six authored tests and composed-profile convergence pass.
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk go test ./internal/baseline -run 'TestPortableVerificationRoleMapping|TestExecutableVerificationCommandRequiresLocalDeclaration|TestIncrementalVerificationDecisionProjectsLikeTheGate' -count=1`
  — exit 0; existing role-mapping and command-declaration checks pass.
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk go test ./internal/cli -run '^TestGuidanceCompositionJourney$' -count=1`
  — exit 0 when run alone; the journey test remains unchanged. An earlier
  invocation encountered host Go-cache permissions; the first cache-local run
  overlapped digest regeneration, so its assertions passed but the repository
  guard refused the concurrent fixture writes. The isolated rerun passed both.
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk make baseline-digests`
  — exit 0, `ok: true`, `changed: false`; derived artifacts already match.
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk proxy go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
  — run twice sequentially; both exit 0 and report `File changes: 0`.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0.

Scope review: only this Task file, the public guide, catalog validation,
profile alignment, and the two new gate-part files changed. Neither profile
asset, the frozen corpora, the Task Graph, nor another Task file changed. No
top-level test or exported function was renamed or removed. The declared
Verification command was not run; status and settlement remain Daemon-owned.
