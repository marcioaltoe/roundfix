---
task: task_02
spec: 0239-a-glossary-every-spec-keeps-current
status: completed
type: backend
complexity: medium
---

# Task 02: The Baseline's domain clauses say the glossary and the ADRs are the base

## Overview

The CONTEXT workflow's domain clauses ask for a glossary check at the close of
a Spec, but they do not say that the domain context and the ADRs are the base,
that a Spec writes the terms it introduces in its own Task, or where the
upstream `domain-modeling` skill must write now that it names its file
differently ([the adopted Backlog Entry of 2026-10-06](references/2026-10-06-the-context-driven-loop-does-not-keep-context-md-current.md)).
This Task rewords `clause.context.read-domain-contract` and
`clause.domain.glossary-currency` in place with `_techspec.md` → Exact clause
texts, raises the versions, regenerates the derived files and refreshes this
repository's guides (ADR-0244, ADR-0186). It is verifiable on its own: the
domain guide states both clauses and an adopter's plan retains them.

This is an authorized tooling Task. It may change only the files in its
Context, the derived files the sanctioned regeneration rewrites, and this Task
file.

## Requirements

1. MUST replace the `guidance` of `clause.context.read-domain-contract` and of
   `clause.domain.glossary-currency` in
   `internal/baseline/assets/modules/context-workflow.json` with the exact texts
   of `_techspec.md` → Exact clause texts, editing the strings in place so every
   other byte of the file stays as it is. Each clause MUST keep its `id` and its
   `mandatory` enforcement and MUST NOT gain a `replaces` list.
2. MUST raise by one, from the value on this Task's starting commit, the
   versions `_techspec.md` → Version changes lists.
3. MUST create `internal/baseline/glossary_clauses_test.go` with
   `TestTheGlossaryClausesCarryTheirForceAndText`,
   `TestTheGlossaryClausesRenderInTheDomainGuide` and
   `TestAnAdopterRetainsTheGlossaryClauses`, following
   `internal/baseline/grouped_sources_clauses_test.go`. The texts are literals
   in the test, so a reworded clause fails it; the rendering test MUST require
   each clause once as a `mandatory` bullet in the Standard TypeScript Monorepo
   golden `docs/agents/domain.md`; the retention test MUST build a Source
   Baseline adopter in a temporary Git repository and require a ready Managed
   Refresh plan that records both clauses `retained` and no `unaccounted`
   clause.
4. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`;
   a second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a
   golden, a snapshot or a generated guide. The authoring probe measured the
   rewritten set as exactly the derived files in this Task's Context; the
   Result MUST name every file the commands rewrote.
5. MUST NOT change the Source Baseline corpus, manifest or index, the force
   record, any other clause, or any skill.
6. MUST prove each new gate can fail. The Result MUST record one sabotage of a
   clause text (for example dropping its last sentence) and one of the
   retention (for example changing a clause's enforcement), each with the test
   that failed, and that the source was restored and regenerated.

## Subtasks

- [ ] Reword the two clauses and raise the versions.
- [ ] Add the clause test file.
- [ ] Regenerate, refresh twice, and record each sabotage.

## Acceptance Criteria

- [ ] `docs/agents/domain.md` states both reworded clauses, each labelled
      `mandatory`.
- [ ] The Standard TypeScript Monorepo golden states the same clauses.
- [ ] A Source Baseline adopter's Managed Refresh plan is ready and records
      both clauses `retained`.
- [ ] The force record, the duplicate-text check, the record-citation check
      and the catalog validation pass, and a second Managed Refresh is a no-op.

## Context

- instruction: `docs/adr/0244-a-spec-declares-the-domain-terms-it-introduces-and-the-check-holds-them-to-the-glossary.md`
- instruction: `docs/adr/0186-baseline-guidance-states-what-the-product-does-in-adopter-neutral-words.md`
- instruction: `internal/baseline/grouped_sources_clauses_test.go`
- interface: `internal/baseline/assets/modules/context-workflow.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/domain.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/domain.md`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/glossary_clauses_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheGlossaryClausesCarryTheirForceAndText|TestTheGlossaryClausesRenderInTheDomainGuide|TestAnAdopterRetainsTheGlossaryClauses|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestShippedGuidanceCitesNoRepositoryRecord|TestReadoptionCompatibilityMaintainedFixture|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus)$" ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheGlossaryClausesCarryTheirForceAndText TestTheGlossaryClausesRenderInTheDomainGuide TestAnAdopterRetainsTheGlossaryClauses TestBaselineClauseForceIsCharacterized TestNoTwoBaselineClausesShareText TestShippedGuidanceCitesNoRepositoryRecord TestReadoptionCompatibilityMaintainedFixture TestCatalogCompatibility; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three glossary clause tests do not exist, so the command fails.
- `for phrase in "- **mandatory**: The repository's selected domain context and its accepted ADRs are the base of the CONTEXT-driven workflow." "in one of its own Tasks, whose Verification names the term." "it writes to the selected domain context." "Work outside a Spec checks at its close whether it introduced, changed, or retired a term"; do tr -s '[:space:]' ' ' < docs/agents/domain.md | grep -qF -- "$phrase" || { printf 'missing phrase in docs/agents/domain.md: %s\n' "$phrase" >&2; exit 1; }; done; plan="$(go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json)" || exit 1; printf '%s\n' "$plan" | grep -qF -- '"state":"current"' || { printf 'the guides are not refreshed: %s\n' "$plan" >&2; exit 1; }` — expected: exit 0; before this Task the domain guide lacks the reworded clauses, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 3; User Story 4; Core Feature 5; Success Metric 4
- [_techspec.md](_techspec.md) — Exact clause texts; Version changes; Retention; API Contract 6; Testing Approach 4; Build Order 2
- ADR-0244; ADR-0186; ADR-0189

## Result

Implemented the two exact TechSpec clause texts in place, preserving their ids,
`mandatory` enforcement and absence of `replaces`. Starting-commit versions
were module 21, `rule.context.domain-docs` 5 and `guide.domain` 6; they are now
22, 6 and 7. The module diff contains only those three version lines and the
two guidance lines. Added `internal/baseline/glossary_clauses_test.go` following
the grouping-clause tests, with independent text literals, exactly-once forced
bullet assertions, and a temporary Git adopter built by
`newClauseReplacementAdopter` in Managed Refresh mode.

### Focused evidence by acceptance criterion

All successful Go and Make commands used
`GOCACHE=/private/tmp/roundfix-0239-task02-gocache`; the default cache was
sandbox-denied on the first test attempt.

- Repository domain guide: an exact-text inspection required each authored
  clause once as a `mandatory` bullet in `docs/agents/domain.md`; passed.
  The public refresh applied and verified its two postimages.
- Standard TypeScript Monorepo golden: the same inspection and
  `TestTheGlossaryClausesRenderInTheDomainGuide` passed, requiring each whole
  clause once with the `mandatory` label.
- Source Baseline adopter: `TestAnAdopterRetainsTheGlossaryClauses` passed.
  It required a ready Managed Refresh plan, both clause dispositions and
  retention rows `retained`, and no `unaccounted` clause anywhere.
- Preserved force and catalog contracts:
  `go test -count=1 -v ./internal/baseline -run 'TestTheGlossaryClauses|TestAnAdopterRetainsTheGlossaryClauses|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestShippedGuidanceCitesNoRepositoryRecord'`
  exited 0 after restoration, with all six gates passing. The final
  `make baseline-digests` exited 0, including strict catalog validation.
  `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
  applied the first refresh and verified the postimages; the second invocation
  exited 0 with `File changes: 0` and `Idempotence: verified`.
  The first refresh initially encountered a sandbox denial creating its
  Git-private journal; the authorized retry with filesystem escalation applied
  successfully. Both invocations reported the existing nested-carrier warnings
  for the golden and Source Baseline AGENTS.md; neither nested carrier changed.

### Mutation evidence

Before changing the module, the two new text/rendering gates both failed on
the original texts. This established the red starting signal without running
the Task's declared Verification.

1. Text sabotage: removed the last sentence of
   `clause.domain.glossary-currency` in the module, then ran
   `make baseline-digests` so the catalog and rendered golden consistently
   carried that incorrect text. The focused command
   `go test -count=1 -v ./internal/baseline -run '^TestTheGlossaryClauses(CarryTheirForceAndText|RenderInTheDomainGuide)$'`
   exited 1. Both `TestTheGlossaryClausesCarryTheirForceAndText` and
   `TestTheGlossaryClausesRenderInTheDomainGuide` failed for
   `clause.domain.glossary-currency`; the other clause still passed.
   Restored the exact module bytes saved before sabotage and ran
   `make baseline-digests` again, exit 0. No derived file was hand-edited.
2. Retention sabotage: a temporary Go overlay of the new test changed
   `clause.domain.glossary-currency` from `mandatory` to `prohibited` only in
   the adopter helper's cloned in-memory catalog, after creating the adopter.
   `go test -overlay=/private/tmp/roundfix-0239-task02-retention-overlay.json -count=1 -v ./internal/baseline -run '^TestAnAdopterRetainsTheGlossaryClauses$'`
   exited 1 in `TestAnAdopterRetainsTheGlossaryClauses`: the refresh was
   `action_required`, with one `unaccounted` clause,
   `clause.domain.glossary-currency`. The shipped source enforcement and
   Source Baseline never changed. The final focused run omitted the overlay
   and passed; the source had already been restored and regenerated.

### Files rewritten by the sanctioned commands

`make baseline-digests` rewrote exactly these nine derived files:

- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/domain.md`
- `internal/baseline/testdata/catalog.diagnostics.golden.json`
- `internal/baseline/testdata/catalog.digest`
- `internal/baseline/testdata/catalog.normalized.json`
- `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`

The first public Managed Refresh rewrote exactly:

- `docs/agents/domain.md`
- `docs/agents/setup-context.json`

The second refresh rewrote no file. The source module, new Go test and this
Result are the other Task-owned changes. The only pre-existing changed path
was this Task file's Daemon-owned `status: in_progress` change, preserved.
No Source Baseline corpus, manifest or index, force record, other clause,
skill, Task Graph or other Task file changed. Declared Verification and Task
settlement remain Daemon-owned; no commit, push or Pull Request was made.

### Incremental check and final postflight

`GOCACHE=/private/tmp/roundfix-0239-task02-gocache rtk make verify-incremental`
exited 0 with filesystem/process-table escalation and no concurrent edits:
formatting, vet, package tests, skill sync/check and build passed. Its first
sandboxed attempt exited 2: two existing force-stop integration tests could
not enumerate the process table, and writing this Result while tests were
running triggered the repository mutation guard. The retry corrected both
execution conditions, with no production or test changes and no bypass.

`git -c core.fsmonitor=false diff --check` exited 0. A postflight compared
tracked and untracked changes to this Task's Context and found all 14 paths
within the declared set. The Task's authored bytes before Result match the
starting commit except for the pre-existing Daemon-owned `in_progress` status.
