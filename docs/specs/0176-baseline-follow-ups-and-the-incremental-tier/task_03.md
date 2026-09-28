---
task: task_03
spec: 0176-baseline-follow-ups-and-the-incremental-tier
status: pending
type: backend
complexity: medium
---

# Task 03: The docs-layout rule names the override command

## Overview

The clause `clause.spec.keep-artifacts-in-spec-folder` of `rule.spec.docs-layout`
in `internal/baseline/assets/modules/spec-workflow.json` is rendered into every
adopter's mandatory `docs/agents/docs-layout.md`. It tells an agent to record the
approval and "stamp `qa_override: true`" itself. An agent following this
mandatory guide edits `_prd.md` by hand and skips the refusals and provenance of
`roundfix archive <slug> --qa-override --approval <source> --reason <text>`,
contradicting the archive-spec skill.

The module is the canonical source. The guide, this repository's Setup
Manifest and the formatter golden are derived from it, and every agent that
archives a Spec reads the guide.

This is an authorized tooling Task. It may change only the files in its Context,
the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST replace, in the guidance of `clause.spec.keep-artifacts-in-spec-folder`,
   the sentence "Record the approval source, date, covered Spec/revision and
   actual QA outcome or absence; stamp `qa_override: true` and preserve any
   supplied reason." with the two sentences given in the TechSpec's "The
   docs-layout override clause". Those sentences name `roundfix archive <slug>
   --qa-override --approval <source> --reason <text>` and contain `Never
   hand-edit the override stamp`. Every other sentence of the clause MUST stay
   byte-identical, including "A runtime command without override support must
   be reported as unsupported rather than given an invented flag."
2. MUST keep the clause identity and raise by one the versions of
   `rule.spec.docs-layout`, `guide.spec-docs-layout` and the `spec-workflow`
   module.
3. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin or a
   generated guide. The regeneration rewrites:
   - `docs/agents/docs-layout.md` and `docs/agents/setup-context.json`;
   - the formatter golden `docs-layout.md`;
   - the Standard TypeScript Monorepo digest pin;
   - the catalog test data and the plan-characterization goldens listed in the
     Context.

   The parity corpus, the Source Baselines and `internal/baseline/assets/setups/`
   stay byte-identical.
4. MUST add these tests to `internal/baseline/qa_override_clause_test.go`:
   - `TestSpecDocsLayoutClauseRoutesTheOverrideThroughTheCommand`: the embedded
     clause names the command, contains `Never hand-edit the override stamp`,
     and keeps the unsupported-runtime sentence.
   - `TestSpecDocsLayoutClauseNoLongerAsksForAHandStamp`, a separate negative
     test: neither the embedded clause nor the formatter golden
     `docs/agents/docs-layout.md` contains `preserve any supplied reason`.
   - `TestFormatterGoldenDocsLayoutCarriesTheOverrideCommand`: the formatter
     golden `docs/agents/docs-layout.md` contains the command and
     `Never hand-edit the override stamp`.

## Subtasks

- [ ] Change the clause and bump its versions.
- [ ] Regenerate the pins, the golden, the guide and the Setup Manifest.
- [ ] Add the clause and golden tests.

## Acceptance Criteria

- [ ] The module, the formatter golden and `docs/agents/docs-layout.md` name the
      override command and forbid a hand edit of the stamp.
- [ ] None of them asks an agent to stamp the override by hand.
- [ ] A second managed refresh is a no-op.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/unsatisfied-blocking-capabilities.golden.json`
- interface: `docs/agents/docs-layout.md`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/qa_override_clause_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestSpecDocsLayoutClauseRoutesTheOverrideThroughTheCommand|TestSpecDocsLayoutClauseNoLongerAsksForAHandStamp|TestFormatterGoldenDocsLayoutCarriesTheOverrideCommand|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSpecDocsLayoutClauseRoutesTheOverrideThroughTheCommand TestSpecDocsLayoutClauseNoLongerAsksForAHandStamp TestFormatterGoldenDocsLayoutCarriesTheOverrideCommand TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization TestBaselineCompatibilityCorpus; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q "Never hand-edit the override stamp" docs/agents/docs-layout.md && ! grep -q "preserve any supplied reason" docs/agents/docs-layout.md && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task none of the three new named tests exists and `docs/agents/docs-layout.md` still asks for a hand stamp, so the command fails.

## References

- [_techspec.md](_techspec.md) — The docs-layout override clause; Regeneration
  outputs
