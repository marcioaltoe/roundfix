---
task: task_01
spec: 0199-stack-rules-that-say-what-they-mean
status: completed
type: backend
complexity: medium
---

# Task 01: The Bun and TypeScript rules say what the profile runs

## Overview

Four clauses of the `bun` and `typescript` modules read wrongly when taken literally. One sends tests to "Bun-owned commands", which reads as the bare Bun runner, while the profile's Verification runs the package's `test` script. One forbids "another package manager" without naming any, so it binds a Go or Rust toolchain in the same repository. One states two obligations, and one names a profile setting that does not exist. This Task replaces the four texts with the ones the TechSpec gives, adds one scope sentence to each part of the TypeScript and Bun guide, and creates the wording check the later Tasks reuse.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST apply, in `internal/baseline/assets/modules/bun.json` and `internal/baseline/assets/modules/typescript.json`, the four replacements the TechSpec's "Exact texts" gives for task_01. Every clause keeps its `id` and its `enforcement`, and every other clause stays byte-identical.
2. MUST insert, in `internal/baseline/assets/templates/guides/typescript-bun.md` and `internal/baseline/assets/templates/guides/bun.md`, the scope paragraph the TechSpec gives for each, between the heading and the rules token.
3. MUST raise by one, from the value on the starting main, the versions the TechSpec's "Version changes" lists for task_01, including `template.guide.typescript-bun` and `template.guide.bun` in `internal/baseline/assets/templates/index.json`. MUST leave the versions of `template.guide.backend` and `template.guide.frontend` as they are.
4. MUST create `internal/baseline/stack_rule_wording_test.go` with `stackWordingFindings`, `clauseForce` and the three tests the TechSpec's Testing Approach 1 names. The negative test MUST feed the check the four replaced sentences as a literal. The file MUST compile before the regeneration runs, because the regeneration compiles the package's tests.
5. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a
   golden or a generated guide, and MUST keep each module file's existing
   formatting by replacing strings and version numbers in place.
6. MUST NOT add or remove a clause, and MUST NOT edit `internal/baseline/plan_test.go` or any other existing test file. The new tests call the existing helpers from the new file.
7. MUST change no exported function signature and rename or remove no top-level test.

## Subtasks

- [ ] Replace the four clause texts and raise the versions.
- [ ] Add the two scope paragraphs and raise the two template versions.
- [ ] Create the wording check, its negative test and the force test.
- [ ] Regenerate the pins, the golden and the Setup Manifest.

## Acceptance Criteria

- [ ] The rendered TypeScript and Bun guide carries the test script rule, the named package managers, the single type-error obligation, the Verification-based warnings condition and both scope sentences.
- [ ] The rendered guide carries none of the four replaced sentences, and the wording check reports each of them when it is fed the old text.
- [ ] Each of the four reworded clauses keeps its identifier and enforcement level.
- [ ] A second Managed Refresh is a no-op.

## Context

- instruction: `docs/adr/0190-a-stack-rule-names-the-language-or-workspace-it-governs.md`
- interface: `internal/baseline/assets/modules/bun.json`
- interface: `internal/baseline/assets/modules/typescript.json`
- interface: `internal/baseline/assets/templates/guides/typescript-bun.md`
- interface: `internal/baseline/assets/templates/guides/bun.md`
- interface: `internal/baseline/assets/templates/index.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/typescript-bun.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/stack_rule_wording_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheTypeScriptAndBunGuideSaysWhatItGoverns|TestStackWordingCheckReportsTheWordingItReplaced|TestTheRewordedBunAndTypeScriptClausesKeepTheirForce|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestStandardTypeScriptStructuralClauseRetention)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTheTypeScriptAndBunGuideSaysWhatItGoverns TestStackWordingCheckReportsTheWordingItReplaced TestTheRewordedBunAndTypeScriptClausesKeepTheirForce TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization TestStandardTypeScriptStructuralClauseRetention; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && golden=internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/typescript-bun.md && for phrase in "These rules govern the repository's TypeScript sources and tests." "These rules govern the Bun workspace" "never through a bare runner such as" "(npm, pnpm, yarn, npx)" "never hide one to make Verification pass." "When the repository's Verification treats warnings as errors"; do tr -s '[:space:]' ' ' < "$golden" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$golden" "$phrase" >&2; exit 1; }; done && for phrase in "scripts, tests, and lockfile updates" "preserve the repository's package and lockfile workflow" "Do not substitute another package manager" "selected TypeScript/Bun profile"; do if tr -s '[:space:]' ' ' < "$golden" | grep -qF -- "$phrase"; then printf 'stale phrase in %s: %s\n' "$golden" "$phrase" >&2; exit 1; fi; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task none of the three new named tests exists and the golden still carries the four replaced sentences, so the command fails.

## References

- [_techspec.md](_techspec.md) — Exact texts; Version changes; Testing Approach 1; Testing Approach 5; Build Order 1
- `_prd.md` → Goal 1; Goal 2; Goal 5; Core Feature 1; Success Metric 1; Success Metric 7; Success Metric 8
- `_techspec.md` → API Contract 1
- ADR-0058, ADR-0059, ADR-0060, ADR-0081, ADR-0149, ADR-0186, ADR-0190

## Result

Implemented the Task 01 wording slice. The four clauses use the TechSpec's
exact replacements, both guide templates carry their scope paragraphs, and
all nine named module/rule/guide/template versions increased by one. Backend
and frontend template versions remain unchanged. Source JSON formatting was
preserved through in-place string and version replacement.

Acceptance evidence:

- Rendered guide wording: `TestTheTypeScriptAndBunGuideSaysWhatItGoverns`
  passes against a Standard TypeScript Monorepo Plan postimage. It requires
  the test-script rule, named JavaScript package managers, TypeScript type-error
  obligation, Verification-based warning condition and both scope paragraphs.
- Replaced wording: `TestStackWordingCheckReportsTheWordingItReplaced`
  passes with four independent literal old sentences, asserting seven missing
  requirements and four replaced-sentence findings. Its absent-guide and
  wrapped-whitespace cases also pass. The rendered-guide test rejects all four
  old sentences. Before the source edits, the rendered-guide test failed with
  precisely those seven missing and four replaced-sentence findings.
- Clause force: `TestTheRewordedBunAndTypeScriptClausesKeepTheirForce` passes
  for all four original identifiers and literal enforcement expectations,
  including absent-identifier checks. A separate comparison with HEAD confirmed
  unchanged clause sets/enforcement and byte-identical untouched clauses.
- Managed Refresh: the sanctioned update applied one Setup Manifest change;
  repeating the same command exited 0 with `File changes: 0`,
  `Idempotence: verified` and `approved Baseline Plan is already applied`.

Focused checks and regeneration:

- `GOCACHE=/tmp/roundfix-0199-task01-gocache rtk proxy go test -count=1 -v -run '^Test(TheTypeScriptAndBunGuideSaysWhatItGoverns|StackWordingCheckReportsTheWordingItReplaced|TheRewordedBunAndTypeScriptClausesKeepTheirForce)$' ./internal/baseline`
  — exit 0; all three new tests and their subtests passed.
- `GOCACHE=/tmp/roundfix-0199-task01-gocache rtk proxy make baseline-digests`
  — exit 0; sanctioned formatter golden, profile digest, catalog snapshots and
  four plan-characterization goldens regenerated; strict catalog validation
  passed. No generated file was hand-edited.
- `GOCACHE=/tmp/roundfix-0199-task01-gocache rtk proxy go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
  — successful apply and zero-change repeat both exited 0. The first sandboxed
  attempt could not create the Git-private transaction directory; rerunning
  with the required filesystem access succeeded. The command reports the two
  nested fixture/source-carrier warnings and skips skills as requested.
- `GOCACHE=/tmp/roundfix-0199-task01-gocache rtk make verify-incremental`
  — exit 0 on the sequential rerun with required process-table access;
  formatting, vet, all package tests, skill checks and build passed. The first
  run exited 2: owner-process integration tests lacked process-table access,
  and suite guards detected this Agent's concurrent Managed Refresh and Result
  writes. The rerun kept the repository unchanged while checking it.
- Changed-file postflight — all 17 changed paths are in this Task's Context,
  including this Task file and the new test; status remains `in_progress`.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0.

The Daemon retains ownership of status and declared Verification. No existing
test file, other Task file, Task Graph or production Go code was edited. No
commit, push or pull request was made.
