---
task: task_03
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
status: pending
type: backend
complexity: medium
---

# Task 03: Authorization and evidence clauses state today's refusals

## Overview

Five clauses in the `core` and `spec-workflow` modules describe refusals the product no longer makes, or a control it never had. The tooling-scope clause fails a Task for any undeclared path. The choreography clause says an absent record refuses nothing, although `roundfix deliver start` refuses it. The verification clause describes an `execution_approvals` record that no code reads. The authorization clause omits the `paths: []` grant. The outside-evidence clause says the row never blocks. This Task replaces the sentences the TechSpec names and leaves the rest of each clause untouched.

The verification clause keeps its obligation and says who meets it: a Run has committed provenance by construction, and nothing else checks it. The missing check is recorded as an open Backlog Entry and is not implemented here.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST apply, in `internal/baseline/assets/modules/core.json`, the two replacements the TechSpec's "Clause texts" gives for `clause.core.verification-two-tiers` and `clause.core.tooling-commit-choreography`.
2. MUST apply, in `internal/baseline/assets/modules/spec-workflow.json`, the three changes the TechSpec gives for `clause.spec.project-constraints-02-tooling-authorization`, `clause.spec.project-constraints-03-bounded-execution` and `clause.spec.project-constraints-06-outside-evidence`. The second sentence of the bounded-execution clause, about a new or widened grant, MUST stay byte-identical.
3. MUST keep all five clauses `mandatory`, and raise by one the versions of each changed rule, of `guide.agent-instructions`, of `guide.spec-routing`, and of the `core` and `spec-workflow` modules.
4. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a
   golden or a generated guide, and MUST keep the module file's existing
   formatting by replacing strings and version numbers in place.
5. MUST add `internal/baseline/authorization_clauses_test.go` with the two tests the TechSpec's Testing Approach 3 names. The negative test MUST search every module and every formatter golden for the four removed phrases.
6. MUST NOT implement a provenance check, add an approval record, or change `internal/authorization`, `internal/speccheck` or `internal/cli`.
7. MUST NOT edit `internal/docscontract/publicdocs_test.go`. The release-planning sentences that test pins are not part of this Task.

## Subtasks

- [ ] Apply the five clause changes and raise the versions.
- [ ] Regenerate the pins, the two goldens, this repository's two guides and the Setup Manifest.
- [ ] Add the positive and negative clause tests.

## Acceptance Criteria

- [ ] `docs/agents/agent-instructions.md` and `docs/agents/spec-routing.md` carry each new sentence.
- [ ] No module, golden or rendered guide contains `execution_approvals`, `does not by itself refuse implementation`, `may mutate only its bounded` or `never blocks the Spec`.
- [ ] A second managed refresh is a no-op.

## Context

- instruction: `docs/adr/0186-baseline-guidance-states-what-the-product-does-in-adopter-neutral-words.md`
- instruction: `docs/backlog/2026-09-30-authored-verification-runs-without-a-provenance-check.md`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- interface: `docs/agents/agent-instructions.md`
- interface: `docs/agents/spec-routing.md`
- creates: `internal/baseline/authorization_clauses_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAuthorizationClausesStateTodaysRefusals|TestAuthorizationClausesDropTheRemovedPhrases|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAuthorizationClausesStateTodaysRefusals TestAuthorizationClausesDropTheRemovedPhrases TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && for pair in "docs/agents/agent-instructions.md|No approval record makes an untrusted source executable" "docs/agents/agent-instructions.md|grants no governed mutation and no delivery operation" "docs/agents/spec-routing.md|grants its listed operations and bounds no path" "docs/agents/spec-routing.md|the QA gate discloses it" "docs/agents/spec-routing.md|holds Pull Request preparation"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && for pair in "docs/agents/agent-instructions.md|execution_approvals" "docs/agents/agent-instructions.md|does not by itself refuse implementation" "docs/agents/spec-routing.md|may mutate only its bounded" "docs/agents/spec-routing.md|never blocks the Spec"; do file="${pair%%|*}"; phrase="${pair#*|}"; if tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase"; then printf 'stale phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; fi; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task neither new named test exists and the guides still carry the four removed phrases, so the command fails.

## References

- [_techspec.md](_techspec.md) — Clause texts; Testing Approach 3
- `_prd.md` → Goal 2; Core Feature 3; Success Metric 3
- `_techspec.md` → API Contract 2
- [references/2026-09-30-baseline-guides-contradict-the-shipped-product.md](references/2026-09-30-baseline-guides-contradict-the-shipped-product.md)
- ADR-0104, ADR-0130, ADR-0166, ADR-0179, ADR-0186

## Result
