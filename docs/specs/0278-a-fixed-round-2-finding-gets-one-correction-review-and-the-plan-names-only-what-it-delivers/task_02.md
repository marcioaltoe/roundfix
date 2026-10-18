---
task: task_02
spec: 0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers
status: pending
type: backend
complexity: high
---

# Task 02: The Spec Consistency Check refuses a Verification path the plan does not provide

## Overview

On 2026-10-17 a Run of Spec 0274 failed twice. task_03's Verification grepped
a file its Task no longer declared, that no Task created and that did not
exist (`_prd.md` → Sources 2). This Task adds `SC-VERIFY-PATH-UNDECLARED`, an
error from the Tasks stage, that reports such a path before any Run starts.
It is wired the way `SC-SKILL-VERSION-UNTASKED` was. It is verifiable on its
own through check tests on temporary repositories and the corpus golden.

## Requirements

1. MUST answer the operator observation of 2026-10-17 (Spec 0274, pull
   requests #559 and #560) cited in `_prd.md` → Sources 2 by implementing API
   Contract 5 of `_techspec.md` (ADR-0287).
2. MUST define `CodeVerifyPathUndeclared` as `SC-VERIFY-PATH-UNDECLARED` in
   the new `internal/speccheck/verification_path.go`. It is reported with
   `SeverityError`, with no horizon, for each pending Task that is not QA, and
   repositories without a Task Graph skip it with the manifest path.
   - It MUST read only the Spec's Task files and the checked tree's file
     system.
   - It MUST execute no command.
   - The word rules, the exclusions, the declaration sources (`creates:`,
     `interface:`, `deletes:`, `instruction:`, `## Recorded paths` through
     `spec.RecordedTaskPaths`), the dependency closure and the existence test
     are API Contract 5's.
   - The summary, location and fix texts are API Contract 5's.
3. MUST wire the code like `SC-SKILL-VERSION-UNTASKED`, and MUST keep every
   earlier count unchanged:
   - a `stagedDetectors` entry with `StageTasks` in
     `internal/speccheck/coherence.go`;
   - a call from `Check` in `internal/speccheck/constraints.go`;
   - a `corpusFindingCodes` entry in `internal/docscontract/corpus_test.go`;
   - a count of 0 in `internal/docscontract/testdata/corpus-golden.json`, whose
     `update` text gains the sentence "SC-VERIFY-PATH-UNDECLARED joins the
     characterized code set, with every earlier count unchanged.";
   - the same count and sentence in the pin of
     `internal/spec/archive_layout_characterization_test.go`.
4. MUST add `TestVerificationPathUndeclaredRaisesTheMissingFile` and
   `TestVerificationPathUndeclaredIgnoresWhatThePlanProvides` to the new
   `internal/speccheck/verification_path_test.go`, each calling `t.Parallel()`
   first and building a temporary repository with the fixtures already used
   beside them.
   - The first reproduces the 0274 shape: a pending Task whose Verification
     greps a missing, undeclared file under an existing top-level directory.
     It asserts one error naming the Task, the line and the path.
   - The second asserts no finding for each of these:
     - a declared path, a recorded path, and a path a needed Task creates;
     - a directory prefix of a declared path, and an existing path;
     - a grep pattern, a `printf` operand, and a negated `test` operand;
     - a path under the Spec Root, a completed Task, and the QA Task.
5. MUST document the code:
   - a row in the check table of `docs/user-guide/commands/spec.md`;
   - a bullet in `.agents/skills/roundfix/references/spec.md`, each naming
     `SC-VERIFY-PATH-UNDECLARED`;
   - a rule headed `Verification paths` in `.agents/skills/write-tasks/SKILL.md`,
     which says a Task declares every repository path its Verification names
     unless the path exists or a Task it needs creates it.

   No text inside any `### QA settlement` section may change.
6. MUST run these commands, because this Task edits a user guide page, the
   Roundfix reference and the `write-tasks` skill. No record or mirror is
   hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to raise and record the Roundfix and `write-tasks` skill versions;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Add the detector and its wiring.
- [ ] Re-record the corpus golden and its pin.
- [ ] Add the two tests.
- [ ] Document the code in the guide, the reference and the `write-tasks` rule.
- [ ] Sync the mirrors, record the versions and the Behavior Surface Record.

## Acceptance Criteria

- [ ] The 0274 shape raises one `SC-VERIFY-PATH-UNDECLARED` error naming the path.
- [ ] Every path the plan provides, every excluded word, completed Tasks and the QA Task raise nothing.
- [ ] The active corpus counts 0 for the new code, and every earlier count is unchanged.

## Consumers

- The Spec Consistency Check's codes and corpus golden — consumers: `./internal/speccheck`, `./internal/docscontract`, `./internal/spec`

## Context

- instruction: `docs/adr/0287-a-plan-is-checked-for-undeclared-verification-paths-and-every-settlement-checks-the-setup-manifest.md`
- instruction: `internal/speccheck/skill_version.go`
- instruction: `internal/speccheck/verification.go`
- creates: `internal/speccheck/verification_path.go`
- creates: `internal/speccheck/verification_path_test.go`
- interface: `internal/speccheck/coherence.go`
- interface: `internal/speccheck/constraints.go`
- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`
- interface: `internal/spec/archive_layout_characterization_test.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestVerificationPathUndeclaredRaisesTheMissingFile|TestVerificationPathUndeclaredIgnoresWhatThePlanProvides)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestVerificationPathUndeclaredRaisesTheMissingFile TestVerificationPathUndeclaredIgnoresWhatThePlanProvides; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md internal/docscontract/testdata/corpus-golden.json internal/spec/archive_layout_characterization_test.go internal/docscontract/corpus_test.go; do grep -qF -- 'SC-VERIFY-PATH-UNDECLARED' "$f" || grep -qF -- 'CodeVerifyPathUndeclared' "$f" || { printf 'missing SC-VERIFY-PATH-UNDECLARED in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < .agents/skills/write-tasks/SKILL.md | grep -qF -- 'Verification paths' || { printf 'missing write-tasks rule Verification paths\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 ./internal/speccheck || exit 1; go test -count=1 -run '^TestArchiveLayoutCharacterization' ./internal/spec || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the two tests do not exist, so the command fails at its first pass check. Afterwards these checks pass: the check tests, the documentation, the corpus golden and its pin, the skill versions and the Behavior Surface Record.

## References

- `_prd.md` → Sources; Goals; User Story 2; Core Feature 3; Core Feature 5; Success Metric 3
- `_techspec.md` → API Contract 5; API Contract 7; Invariant 3; Current behavior; Testing Approach; Build Order 2
- ADR-0287; ADR-0117; ADR-0280; ADR-0189; ADR-0256; ADR-0259
