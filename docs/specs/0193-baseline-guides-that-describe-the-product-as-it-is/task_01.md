---
task: task_01
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
status: pending
type: backend
complexity: medium
---

# Task 01: No clause is rendered twice, and every clause keeps its force

## Overview

The backend module lists one paragraph under two entries of the same rule, so every adopter's backend guide renders it twice. This Task removes the duplicate entry and adds two checks over the embedded catalog: no two clauses share their text, and every clause keeps its enforcement level. The force check is the characterization the later Tasks of this Spec rely on: they rewrite clause guidance and must not change a clause's level.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST remove, from `rule.backend.boundary-contracts` in `internal/baseline/assets/modules/backend.json`, the entry whose `id` is `rule.backend.boundary-contracts`. The entry `clause.backend.boundary-contracts` keeps its guidance and its `mandatory` level.
2. MUST raise by one the versions of `rule.backend.boundary-contracts`, `guide.backend` and the `backend` module.
3. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a
   golden or a generated guide, and MUST keep the module file's existing
   formatting by replacing strings and version numbers in place.
4. MUST add `internal/baseline/clause_characterization_test.go` with the four tests the TechSpec's Testing Approach 1 names:
   - the uniqueness check compares every clause-level and rule-level guidance string after lower-casing and collapsing whitespace;
   - its negative test runs the same function on a copy with one guidance repeated;
   - the force table lists every clause identifier with its enforcement level, written as a literal in the test, and fails for a missing clause, an unlisted clause and a changed level;
   - the rendering test counts the boundary sentence in the backend formatter golden and expects one.
5. MUST keep every other clause, in every module, byte-identical.
6. MUST change no exported function signature and rename or remove no top-level test.

## Subtasks

- [ ] Remove the duplicate entry and raise the three versions.
- [ ] Regenerate the pins, the golden and the Setup Manifest.
- [ ] Add the uniqueness, force and rendering tests, each negative case separate.

## Acceptance Criteria

- [ ] The backend formatter golden holds the boundary paragraph once.
- [ ] A catalog copy with a repeated guidance string is reported by the uniqueness check.
- [ ] The force table matches every clause of the embedded catalog, and the duplicate identifier is absent from both.
- [ ] A second managed refresh is a no-op.

## Context

- instruction: `docs/adr/0186-baseline-guidance-states-what-the-product-does-in-adopter-neutral-words.md`
- interface: `internal/baseline/assets/modules/backend.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/clause_characterization_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestNoTwoBaselineClausesShareText|TestADuplicatedClauseTextIsReported|TestBaselineClauseForceIsCharacterized|TestTheBackendBoundaryParagraphRendersOnce|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestNoTwoBaselineClausesShareText TestADuplicatedClauseTextIsReported TestBaselineClauseForceIsCharacterized TestTheBackendBoundaryParagraphRendersOnce TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization TestBaselineCompatibilityCorpus; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task none of the four new named tests exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Clause texts; Testing Approach 1; Version changes
- `_prd.md` → Goals 3 and 4; Core Feature 1; Success Metrics 1 and 5
- [references/2026-09-30-the-backend-guide-renders-one-clause-twice.md](references/2026-09-30-the-backend-guide-renders-one-clause-twice.md)
- ADR-0186

## Result
