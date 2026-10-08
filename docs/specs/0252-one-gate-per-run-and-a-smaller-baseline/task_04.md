---
task: task_04
spec: 0252-one-gate-per-run-and-a-smaller-baseline
status: pending
type: backend
complexity: medium
---

# Task 04: The root block says when a session needs the docs-layout and Secondbrain guides

## Overview

The root block calls the docs-layout guide mandatory for every session and
points every session at the Secondbrain guide. Together they are about 23 KB,
mostly record templates and Secondbrain rules that a Run's code Task never
uses and whose Secondbrain it cannot reach (finding B07 of the Baseline audit
of 2026-10-08, `docs/references/2026-10-08-baseline-audit.md`). This Task
rewrites the two root templates with `_techspec.md` → Template texts
(ADR-0257). It is verifiable on its own: this repository's `AGENTS.md` keeps
the domain guide mandatory and names when each of the two guides applies.

This is an authorized tooling Task. It may change only the files in its
Context, the derived files the sanctioned regeneration rewrites, and this Task
file.

## Requirements

1. MUST answer finding B07 of the Baseline audit of 2026-10-08 by replacing
   `internal/baseline/assets/templates/root/context-workflow.md` and
   `internal/baseline/assets/templates/root/secondbrain.md` with the whole
   texts of `_techspec.md` → Template texts, keeping each template's token
   list in `templates/index.json`.
2. MUST raise by one the task_04 versions of `_techspec.md` → Version changes,
   then run
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
   `make baseline-digests` and
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. `AGENTS.md` and the golden
   `AGENTS.md` change only inside the two root blocks, and the Result MUST
   name every file the commands rewrote.
3. MUST create `internal/baseline/root_reading_set_test.go` with
   `TestTheRootNamesWhenToReadTheDocsLayoutGuide` and
   `TestTheRootNamesWhenToReadTheSecondbrainGuide`. They read the golden
   `formatter-fixtures/standard-typescript-monorepo/golden/AGENTS.md` through
   `catalog.Asset`. Each requires its root sentence exactly once with the
   guide references rendered, and requires that "Domain and documentation
   rules are mandatory" and "Optional cross-project knowledge follows" no
   longer appear.
4. MUST prove each new gate can fail. The Result MUST record one sabotage of
   each template, with the test that failed, and that the source was restored
   and regenerated.
5. MUST NOT change any other root template, any guide template, any clause,
   or the bytes of `AGENTS.md` outside its managed markers.

## Subtasks

- [ ] Rewrite the two root templates and raise the versions.
- [ ] Record, regenerate and refresh twice.
- [ ] Add the root tests and record each sabotage.

## Acceptance Criteria

- [ ] `AGENTS.md` keeps the domain guide mandatory and names the work that
      needs the docs-layout guide.
- [ ] `AGENTS.md` names the work that needs the Secondbrain guide and says a
      Run session skips it.
- [ ] The golden `AGENTS.md` carries the same sentences, and a second refresh
      of this repository is a no-op.

## Context

- instruction: `docs/adr/0257-inside-a-run-the-daemon-is-the-only-full-gate-and-the-baseline-states-each-rule-once.md`
- instruction: `docs/references/2026-10-08-baseline-audit.md`
- interface: `internal/baseline/assets/templates/root/context-workflow.md`
- interface: `internal/baseline/assets/templates/root/secondbrain.md`
- interface: `internal/baseline/assets/templates/index.json`
- interface: `internal/baseline/assets/modules/context-workflow.json`
- interface: `internal/baseline/assets/modules/secondbrain.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/AGENTS.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `AGENTS.md`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/root_reading_set_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheRootNamesWhenToReadTheDocsLayoutGuide|TestTheRootNamesWhenToReadTheSecondbrainGuide|TestEveryBaselineModuleVersionIsRecorded|TestCatalogCompatibility|TestFormatterComposition|TestBaselinePlanCharacterization)$" ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheRootNamesWhenToReadTheDocsLayoutGuide TestTheRootNamesWhenToReadTheSecondbrainGuide TestEveryBaselineModuleVersionIsRecorded TestCatalogCompatibility; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two root tests do not exist, so the command fails.
- `root="$(tr -s '[:space:]' ' ' < AGENTS.md)"; for phrase in "Domain rules are mandatory" "whose rules then apply, before creating, changing, moving, or retiring" "and before a test or build step reads one" "before consulting or writing the Secondbrain or authoring an Idea, PRD, or TechSpec" "a Run session, which cannot reach it, skips it"; do printf '%s\n' "$root" | grep -qF -- "$phrase" || { printf 'missing phrase in AGENTS.md: %s\n' "$phrase" >&2; exit 1; }; done; for phrase in "Domain and documentation rules are mandatory" "Optional cross-project knowledge follows"; do if printf '%s\n' "$root" | grep -qF -- "$phrase"; then printf 'old root sentence remains: %s\n' "$phrase" >&2; exit 1; fi; done; plan="$(go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json)" || { printf '%s\n' "$plan"; exit 1; }; printf '%s\n' "$plan" | grep -qF -- '"state":"current"' || { printf 'the guides are not refreshed: %s\n' "$plan" >&2; exit 1; }` — expected: exit 0; before this Task the root block lacks the conditional sentences, so the command fails at its first phrase.

## References

- [_prd.md](_prd.md) — Goal 5; User Story 5; Core Feature 5; Success Metric 5
- [_techspec.md](_techspec.md) — API Contract 3; Template texts; Version changes; Derived files; Testing Approach; Build Order 4
- ADR-0257; ADR-0250
