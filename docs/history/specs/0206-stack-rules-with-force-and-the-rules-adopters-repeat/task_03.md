---
task: task_03
spec: 0206-stack-rules-with-force-and-the-rules-adopters-repeat
status: completed
type: backend
complexity: medium
---

# Task 03: Rust rules carry force, with the typed-error policy

## Overview

The Rust module ships two paragraphs with no force and no error policy. This Task splits them into clauses with force, adds the rule the maintainer decided (a thin binary over the library, typed errors in library and domain code, a type-erased error only at a binary's entry point, and no panic on a user-reachable path), requires that rule in the Rust CLI profile, and gives the Rust guide its scope sentence. It is verified through a Rust CLI plan rendered in a temporary repository.

This is an authorized tooling Task. It may change only the files in its Context, the derived files the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST replace the rules of `internal/baseline/assets/modules/rust.json` with the clauses and the new rule the TechSpec's "Exact texts" gives for task_03, keeping both existing rule identifiers and the module's skill lists, and MUST move its `schemaVersion` to `setup-context-driven/module-v3` with `"repositoryExtensions": []`.
2. MUST list `rule.rust.library-and-entry-point` in `guide.rust` and in the Rust CLI profile's `requiredRules`, and MUST insert the TechSpec's scope sentence into the Rust guide template between its heading and its rules token.
3. MUST raise by one, from the value on the starting main, the versions the TechSpec's "Version changes" lists for task_03; the new rule starts at version 1.
4. MUST add the nine Rust clauses to the force record.
5. MUST create `internal/baseline/stack_force_rust_test.go` with the four tests the TechSpec's Testing Approach 3 names, reusing the helpers task_02 created.
6. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`;
   a second refresh MUST report `File changes: 0`. MUST NOT hand-edit a snapshot or a golden.
7. MUST NOT edit any skill, including the vendored Rust skills, and MUST change no exported function signature and rename or remove no top-level test.

## Subtasks

- [ ] Split the Rust rules into clauses and add the error rule.
- [ ] Require the rule in the profile and add the guide's scope sentence.
- [ ] Raise versions, update the force record and create the new test file.
- [ ] Regenerate and refresh twice.

## Acceptance Criteria

- [ ] No Rust rule carries rule-level guidance.
- [ ] A Rust CLI plan renders the scope sentence, the typed-error, entry-point and no-panic clauses with their force, and no longer the manifest paragraph it replaced.
- [ ] The Rust CLI profile requires the new rule.
- [ ] A second Managed Refresh is a no-op.

## Context

- instruction: `docs/adr/0202-a-stack-rule-carries-force-and-a-preference-names-what-licenses-the-exception.md`
- interface: `internal/baseline/assets/modules/rust.json`
- interface: `internal/baseline/assets/profiles/rust-cli.json`
- interface: `internal/baseline/assets/templates/guides/rust.md`
- interface: `internal/baseline/assets/templates/index.json`
- interface: `internal/baseline/clause_characterization_test.go`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/stack_force_rust_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestEveryRustRuleCarriesForce|TestTheRustGuideStatesTheErrorPolicy|TestTheRustClausesCarryTheirForce|TestTheRustProfileRequiresTheErrorRule|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestShippedGuidanceCitesNoRepositoryRecord|TestCatalogCompatibility|TestBaselinePlanCharacterization)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryRustRuleCarriesForce TestTheRustGuideStatesTheErrorPolicy TestTheRustClausesCarryTheirForce TestTheRustProfileRequiresTheErrorRule TestBaselineClauseForceIsCharacterized TestCatalogCompatibility; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && template=internal/baseline/assets/templates/guides/rust.md && tr -s '[:space:]' ' ' < "$template" | grep -qF -- "These rules govern the repository's Rust crates" || { printf 'missing scope sentence in %s\n' "$template" >&2; exit 1; }; go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the four new tests do not exist and the template has no scope sentence, so the command fails.

## References

- `_prd.md` → Goal 1; Goal 3; Story 3; Core Feature 2; Success Metric 1; Success Metric 3
- `_techspec.md` → Candidate rules; Exact texts (task_03); Version changes; API Contract 3; Testing Approach 3; Testing Approach 5; Build Order 3
- ADR-0059, ADR-0081, ADR-0149, ADR-0186, ADR-0190, ADR-0202

## Result

The first handoff stopped over the nine-clause count discrepancy. After the
same-session Verification Feedback requested repair, implementation followed
all eight clauses explicitly defined in the TechSpec's Exact texts, plus its
separate scope sentence. No ninth obligation was invented. Requirement 4 and
Testing Approach 3 still say nine; the Spec owner can correct that wording as
a follow-up outside this Task's implementation diff.

The inspected diagnostic artifact was
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261002T051041Z_3b0b74eb6eaab110/verification/batch-003-attempt-1.log`.
It identified the missing Rust test, consistent with the first handoff having
made no source changes. Its body is not embedded here.

### Implementation and acceptance evidence

- No Rust rule carries rule-level guidance: the module now uses module-v3,
  empty `repositoryExtensions`, and clauses for both retained rule identities
  and the new error rule. `TestEveryRustRuleCarriesForce` passes.
- A temporary Rust CLI plan renders the scope sentence and all eight exact
  clause texts with force, including typed library errors, thin binaries,
  entry-point-only type-erased reports and the prohibition of panics on
  user-reachable paths. It excludes the replaced manifest paragraph and has
  no unlabelled rule bullet. `TestTheRustGuideStatesTheErrorPolicy` passes.
- The force record gains the eight defined Rust clauses with their specified
  enforcement. `TestTheRustClausesCarryTheirForce` passes using the existing
  clause-force helper; the rule-level and wording helpers are reused too.
- The guide and Rust CLI profile both list
  `rule.rust.library-and-entry-point`. Its version starts at 1;
  `TestTheRustProfileRequiresTheErrorRule` passes. Rust module version rises
  3 to 4, guide and existing rules 2 to 3, and Rust template 1 to 2.
- The second Managed Refresh reports `File changes: 0`, identical current and
  prior catalog digests, and verified idempotence. The first successful
  refresh changed only the Setup Manifest; this repository selects no Rust
  guide.

### Focused checks and regeneration

- The four new Rust tests were written before the asset changes. A focused
  run failed for both old guidance paragraphs, all eight missing clauses,
  the absent scope sentence, the retained manifest paragraph and the missing
  profile requirement. After the asset changes the same focused run exited 0.
- `rtk make baseline-digests` exited 0 with `ok:true` and `changed:true`.
  It rewrote only the catalog digest, normalized catalog and four declared
  plan-characterization goldens. No generated artifact was hand-edited.
- `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk proxy go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
  exited 0 for both successful refreshes, first with one file change, then
  zero. The initial default-cache attempt was sandbox-denied; the first
  task-cache attempt could not write the Git-private transaction lock.
  Elevated execution allowed the authorized refresh, and the no-op refresh
  then succeeded without elevation. Both successful refreshes reported
  nested-carrier warnings for the existing formatter fixture and Source
  Baseline corpus, whose bytes were preserved.
- `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk proxy go test -count=1 -v -run '^Test(EveryRustRuleCarriesForce|TheRustGuideStatesTheErrorPolicy|TheRustClausesCarryTheirForce|TheRustProfileRequiresTheErrorRule)$' ./internal/baseline`
  exited 0 after regeneration and both refreshes, with all four named tests
  reporting PASS. A preceding default-cache invocation was sandbox-denied.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0. Changed-path
  inspection found only this Task's declared sources, test, force record,
  sanctioned derived files and Task Result. Initial status already listed
  this Task file as modified; status remains Daemon-owned.

The declared Verification sequence was not rerun. No Task status or checkbox,
Task Graph, other Task file, skill, exported signature or existing top-level
test name was changed. No commit, push or pull request was created.
