---
task: task_04
spec: 0252-one-gate-per-run-and-a-smaller-baseline
status: completed
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

## Result

Implemented the B07 root reading-set change from ADR-0257. The two source
templates contain the exact whole TechSpec texts, each ending in one newline.
Their template versions and root-block versions each rose by one, preserving
the token lists. The Module Version Record step chose and recorded the changed
`context-workflow` and `secondbrain` module versions; no clause, guide template,
or other root template changed.

### Acceptance evidence

- Domain and docs-layout: `AGENTS.md` keeps the domain guide mandatory and
  requires docs-layout before creating, changing, moving or retiring
  `CONTEXT.md` or a document under `docs/`, and before a test or build reads
  one. `TestTheRootNamesWhenToReadTheDocsLayoutGuide` reads the embedded
  golden through `catalog.Asset`, requires the complete rendered sentence
  exactly once, and rejects both old unconditional sentences.
- Secondbrain: `AGENTS.md` names consultation, writing, and Idea/PRD/TechSpec
  authoring as the triggers and explicitly says a Run session skips it.
  `TestTheRootNamesWhenToReadTheSecondbrainGuide` applies the same exact-once
  and old-sentence checks to the embedded golden.
- Golden and no-op refresh: both roots carry the same sentences. The second
  required text refresh exited 0 with `File changes: 0` and
  `Idempotence: verified`. After both sabotages were restored and regenerated,
  a final refresh again exited 0 with those same results. A comparison against
  snapshots taken before edits confirmed every byte outside the two targeted
  root blocks remains identical in both `AGENTS.md` files.

### Focused checks and required regeneration

- Initial `go test ./internal/baseline -run '^TestTheRootNamesWhenToRead'
  -count=1`: exit 1; both newly added tests rejected the original golden's
  missing conditional sentences and the two old unconditional sentences.
- Required `go test ./internal/baseline -run
  '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions
  -count=1`: exit 0 using `GOCACHE=/tmp/roundfix-task04-gocache`. The first
  attempt was blocked by shared Go-cache filesystem access; no assertion
  failed on that attempt.
- Required `make baseline-digests`: exit 0. It also regenerated each sabotage
  and each restored source successfully. All generators ran with the same
  task-local Go cache.
- Required `go run -buildvcs=false ./cmd/roundfix baseline update --repo .
  --no-skills --yes --format text`: the sandboxed attempt could not open
  the Git worktree transaction lock. The authorized rerun with filesystem
  escalation exited 0, applied two files and verified the approved postimages.
  The next invocation exited 0 with zero file changes; a final invocation
  after source restoration also exited 0 with zero file changes.
- Final `go test ./internal/baseline -run '^TestTheRootNamesWhenToRead'
  -count=1 -v`: exit 0; explicit PASS for both new tests.
- `gofmt -d internal/baseline/root_reading_set_test.go`: no output, exit 0.
- Changed-path postflight: all 19 changed/new paths belong to this Task's
  declared Context, including this Task file and the new test. Only the Task
  file was already modified on entry; its existing content was preserved.

### Sabotage evidence

1. In `templates/root/context-workflow.md`, changed `Domain rules are
   mandatory` to `Domain rules are optional`, preserving all template tokens.
   Ran `make baseline-digests` (exit 0), then
   `go test ./internal/baseline -run
   '^TestTheRootNamesWhenToReadTheDocsLayoutGuide$' -count=1` (exit 1):
   `TestTheRootNamesWhenToReadTheDocsLayoutGuide` reported the required
   sentence occurs 0 times. Restored the original source bytes and ran
   `make baseline-digests` again (exit 0).
2. In `templates/root/secondbrain.md`, replaced `a Run session, which cannot
   reach it, skips it` with `a Run session reads it`, preserving the token.
   Ran `make baseline-digests` (exit 0), then
   `go test ./internal/baseline -run
   '^TestTheRootNamesWhenToReadTheSecondbrainGuide$' -count=1` (exit 1):
   `TestTheRootNamesWhenToReadTheSecondbrainGuide` reported the required
   sentence occurs 0 times. Restored the original source bytes and ran
   `make baseline-digests` again (exit 0). Both final focused tests passed.

### Files rewritten by the required commands

Module Version Record:

- `internal/baseline/assets/modules/context-workflow.json`
- `internal/baseline/assets/modules/secondbrain.json`
- `internal/baseline/module-versions.json`

`make baseline-digests`:

- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/AGENTS.md`
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- `internal/baseline/testdata/catalog.diagnostics.golden.json`
- `internal/baseline/testdata/catalog.digest`
- `internal/baseline/testdata/catalog.normalized.json`
- `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`

Managed Refresh:

- `AGENTS.md`
- `docs/agents/setup-context.json`

The other edited sources are the two root templates and `templates/index.json`;
the new test is `internal/baseline/root_reading_set_test.go`. This Result is the
only agent-authored change to this Task file. Declared Task Verification,
selected repository Verification and incremental Verification remain for the
Daemon. Task status and graph were not edited; no commit, push or Pull Request
was made. No follow-up implementation was added.

## Carry-forward provenance

- Source Run: `run_20261008T180729Z_1c5cbdf685b778d0`
- Source commit: `350e2993ae48a816927a2ec6d50d7b1e636ac169`
