---
task: task_02
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
status: pending
type: backend
complexity: medium
---

# Task 02: The loop clause follows the Delivery Queue

## Overview

`clause.autonomous.loop-01-qa-once` tells every adopter's Agent to archive the candidate and then apply the pre-PR review. The Delivery Queue reviews first and archives second. The clause also never names the Delivery Queue, `roundfix reopen` or Task Carry-Forward, and it cites this repository's Spec 0078 and three of its ADR numbers. This Task replaces the clause with the text the TechSpec fixes, removes the ADR citation from the hook clause, and adds a check that compares the declared order with the order the Delivery Queue runs.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST replace the guidance of `clause.autonomous.loop-01-qa-once` in `internal/baseline/assets/modules/autonomous-work.json` with the text in the TechSpec's "Clause texts", as one JSON string on one line. The clause MUST still begin with `Follow one order per Spec:` and its order sentence MUST end with a full stop, because `SC-LOOP-ORDER-DIVERGENT` reads both.
2. MUST replace, in `clause.autonomous.hook-strictness`, "ADR-0014 makes the Daemon the verification authority" with "The Daemon is the verification authority". Every other sentence of that clause stays byte-identical.
3. MUST keep both clauses `mandatory`, and raise by one the versions of `rule.autonomous.loop`, `rule.autonomous.hook-strictness`, `guide.autonomous-work` and the `autonomous-work` module.
4. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a
   golden or a generated guide, and MUST keep the module file's existing
   formatting by replacing strings and version numbers in place.
5. MUST add `internal/baseline/loop_clause_test.go` and `internal/delivery/loop_clause_order_test.go` with the tests the TechSpec's Testing Approach 2 names. The order test MUST read the module file, not a copy of its text, and MUST take the Delivery Queue's action order from the package's existing fake workflow for a reviewed item, not from a list written in the test.
6. MUST NOT edit `internal/speccheck/citations.go`, any skill, or any existing test.

## Subtasks

- [ ] Replace the loop clause and the hook sentence, and raise the versions.
- [ ] Regenerate the pins, the golden, this repository's guide and the Setup Manifest.
- [ ] Add the wording tests and the order tests, each negative case separate.

## Acceptance Criteria

- [ ] The module, the formatter golden and `docs/agents/autonomous-work.md` declare review before archive and name `roundfix deliver`, `roundfix reopen --spec <slug>` and Task Carry-Forward.
- [ ] None of them contains `Spec 0078` or an ADR number in the loop or hook clause.
- [ ] The declared order equals the Delivery Queue's action order, and the old order is reported as a mismatch.
- [ ] A second managed refresh is a no-op.

## Context

- instruction: `docs/adr/0186-baseline-guidance-states-what-the-product-does-in-adopter-neutral-words.md`
- instruction: `internal/delivery/engine_test.go`
- instruction: `internal/speccheck/citations.go`
- interface: `internal/baseline/assets/modules/autonomous-work.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- interface: `docs/agents/autonomous-work.md`
- creates: `internal/baseline/loop_clause_test.go`
- creates: `internal/delivery/loop_clause_order_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestLoopClauseNamesTheDeliveryQueueAndItsRecoveryActs|TestLoopClauseCitesNoSpecOrDecisionNumber|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestTheLoopClauseOrderMatchesTheDeliveryQueue|TestALoopClauseThatArchivesBeforeReviewIsRefused)$" ./internal/baseline ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestLoopClauseNamesTheDeliveryQueueAndItsRecoveryActs TestLoopClauseCitesNoSpecOrDecisionNumber TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization TestTheLoopClauseOrderMatchesTheDeliveryQueue TestALoopClauseThatArchivesBeforeReviewIsRefused; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && for pair in "docs/agents/autonomous-work.md|apply the configured pre-PR review policy, archive and commit the candidate" "docs/agents/autonomous-work.md|roundfix reopen --spec" "docs/agents/autonomous-work.md|Task Carry-Forward"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && for pair in "docs/agents/autonomous-work.md|Spec 0078" "docs/agents/autonomous-work.md|ADR-0014"; do file="${pair%%|*}"; phrase="${pair#*|}"; if tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase"; then printf 'stale phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; fi; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task none of the four new named tests exists and the guide still orders archive before review, so the command fails.

## References

- [_techspec.md](_techspec.md) — Clause texts; Integration Points; Testing Approach 2
- `_prd.md` → Goal 1; Core Feature 2; Success Metric 2
- `_techspec.md` → API Contract 1
- [references/2026-09-30-baseline-guides-contradict-the-shipped-product.md](references/2026-09-30-baseline-guides-contradict-the-shipped-product.md)
- ADR-0186

## Result
