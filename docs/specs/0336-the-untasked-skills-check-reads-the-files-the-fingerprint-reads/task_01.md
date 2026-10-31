---
task: task_01
spec: 0336-the-untasked-skills-check-reads-the-files-the-fingerprint-reads
status: pending
type: backend
complexity: medium
---

# Task 01: The untasked-skills check and the Behavior Surface fingerprint share one file predicate

## Overview

The Skill Coverage Map package gains one exported predicate that says whether
a Behavior Surface fingerprint can read a file. The map's path lookup and the
`config keys` fingerprint both call it, so a Go test file never asks for a
skills Task while a source file still does. The spec command guide and the
Roundfix skill's spec reference state the rule. This answers the adopted
Backlog Entry of 2026-10-30, "The untasked-skills check counts test files the
fingerprint ignores"
([references](references/2026-10-30-the-untasked-skills-check-counts-test-files-the-fingerprint-ignores.md)).

## Requirements

1. MUST add `FingerprintReads(filePath string) bool` to
   `internal/skillcoverage/skillcoverage.go` as `_techspec.md` → API
   Contract 1 states: `false` exactly when `path.Base(filePath)` ends in
   `_test.go`, `true` otherwise, with a doc comment that says no fingerprint
   reads a Go test file.
2. MUST make `Map.SurfacesForPath` return nil when `FingerprintReads`
   refuses the path, before any source is matched (API Contract 2). MUST NOT
   change `internal/speccheck/skills_declaration.go`; the finding's code,
   summary, location and fix keep their bytes.
3. MUST replace the `strings.HasSuffix(entry.Name(), "_test.go")` filter of
   the `config keys` read in `computeBehaviorSurfaces`
   (`internal/docscontract/skill_coverage_test.go`) with
   `!skillcoverage.FingerprintReads(entry.Name())` (API Contract 3), keeping
   the directory and `.go` filters. Change no other line of that file.
4. MUST add `TestAGoTestFileBelongsToNoBehaviorSurface` to
   `internal/skillcoverage/skillcoverage_test.go` with the map and the eight
   paths `_techspec.md` → Testing Approach lists, asserting both
   `FingerprintReads` and the surface ids `SurfacesForPath` returns. The
   cases `internal/config/config_testing.go` and
   `internal/config/test_helpers.go` must stay read, so a predicate that
   matches any name containing `_test` fails.
5. MUST add `TestAGoTestFileUnderACoveredSourceNeedsNoSkillsTask` to
   `internal/speccheck/skills_declaration_test.go`, calling `t.Parallel()`
   first and using `skillsRepo`, `skillsTasks`, `skillsCheck` and
   `skillsCount` with the cases `_techspec.md` → Testing Approach lists:
   for each of `interface`, `creates` and `deletes`, zero
   `SC-SKILLS-UNTASKED` for `internal/cli/example_test.go` and
   `internal/cli/sub/other_test.go`, and exactly one, whose summary names the
   path, for `internal/cli/example_testing.go` and
   `internal/cli/sub/other.go`; then a test-only `task_01` beside a
   `task_02` declaring `internal/cli/example.go` gives exactly one, naming
   `task_02`. Keep every existing test of the file unchanged.
6. MUST add one sentence containing the phrase `is never a Behavior Surface
   source` (API Contract 4): in `docs/user-guide/commands/spec.md`, to the
   paragraph that begins "The skills rule reads", and in the canonical
   `.agents/skills/roundfix/references/spec.md`, to the paragraph that begins
   "The optional Skills Declaration lives". The sentence says that a Go test
   file, whose name ends in `_test.go`, is never a Behavior Surface source,
   because no fingerprint reads it, so declaring one needs no skills Task
   even when a source pattern matches it. Change no table row and no other
   paragraph; Specs 0334 and 0335 edit other passages of these files first.
7. MUST copy the canonical skill to `skills/` with `make skills-sync`, then
   raise both version fields of `.agents/skills/roundfix/SKILL.md` and its
   mirror with `go test ./skills -run
   '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
   -count=1`, recording `skills/testdata/owned-skill-versions.json`. Never
   hand-edit a mirror or a version field. Specs 0333, 0334 and 0335 raise
   this version first; the record command raises above whatever is recorded.
8. MUST re-record the Behavior Surface Record as the LAST step, after the
   last edit of the guide, with `go test -count=1 -tags docscontract
   ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'
   -record-skill-coverage`, and commit
   `docs/references/behavior-surfaces.json`. Only the
   `guide docs/user-guide/commands/spec.md` fingerprint moves; command help
   does not change.
9. MUST NOT touch `docs/references/skill-coverage.json`, `CONTEXT.md`, any
   `### QA settlement` section, any Baseline module, `docs/agents/`, the
   Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize: a Task declaring only `internal/cli/example_test.go` in the `skillsRepo` fixture reports `SC-SKILLS-UNTASKED` today.
- [ ] Add the predicate and the lookup guard, and point the fingerprint at the predicate.
- [ ] Add both regression tests.
- [ ] Add the sentence to the guide and the reference, sync the mirror, record the skill version, re-record the Behavior Surface fingerprints last.

## Acceptance Criteria

- [ ] A Go test file matched by a covered surface's source gives no `SC-SKILLS-UNTASKED`; a non-test source still gives exactly one.
- [ ] The `config keys` fingerprint and the lookup call `skillcoverage.FingerprintReads`, and no recorded fingerprint moves for an unchanged source.
- [ ] The guide and both reference copies carry the rule, the mirror equals its canonical copy, and the owned-skill version and Behavior Surface records pass.

## Context

- instruction: `docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md`
- instruction: `internal/speccheck/skills_declaration.go`
- interface: `internal/skillcoverage/skillcoverage.go`
- interface: `internal/skillcoverage/skillcoverage_test.go`
- interface: `internal/speccheck/skills_declaration_test.go`
- interface: `internal/docscontract/skill_coverage_test.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Skill Coverage Map lookup `SurfacesForPath` and the predicate `FingerprintReads` — consumers: `./internal/skillcoverage`, `./internal/speccheck`, `./internal/docscontract`
- The owned-skill version record and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestAGoTestFileUnderACoveredSourceNeedsNoSkillsTask$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestAGoTestFileUnderACoveredSourceNeedsNoSkillsTask ' || { echo 'missing pass: TestAGoTestFileUnderACoveredSourceNeedsNoSkillsTask' >&2; exit 1; }; go test -count=1 -timeout 20m -run 'Skill' ./internal/speccheck` — expected: exit 0. Before this Task the test does not exist, so the pass line is missing; afterwards it passes with every skills test of the Spec Check.
- `out="$(go test -count=1 -timeout 20m -v -run '^TestAGoTestFileBelongsToNoBehaviorSurface$' ./internal/skillcoverage 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestAGoTestFileBelongsToNoBehaviorSurface ' || { echo 'missing pass: TestAGoTestFileBelongsToNoBehaviorSurface' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/skillcoverage` — expected: exit 0. Before this Task the test does not exist; afterwards it passes with the whole package.
- `grep -qF 'skillcoverage.FingerprintReads(entry.Name())' internal/docscontract/skill_coverage_test.go || { echo 'the config keys read does not call the shared predicate' >&2; exit 1; }; if grep -qF 'HasSuffix(entry.Name(), "_test.go")' internal/docscontract/skill_coverage_test.go; then echo 'the config keys read keeps its own test-file filter' >&2; exit 1; fi; for f in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'is never a Behavior Surface source' || { printf 'missing rule in %s\n' "$f" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the fingerprint keeps its own filter, so the first check fails. Afterwards it calls the shared predicate, the guide and both reference copies carry the rule, the mirrors equal their canonical copies, and the owned-skill version and documentation contracts, the Behavior Surface fingerprint and `TestBehaviorSurfaceBytesComeFromParsedValuesAndGuideContent` among them, pass.

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 4; Success Metric 5
- `_techspec.md` → Current behavior 1-5; API Contracts 1-4; Invariants 1-4; Declared breaks; Testing Approach; Build Order 1
