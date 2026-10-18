---
task: task_02
spec: 0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers
status: completed
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


## Result

Implemented the Task 02 slice for ADR-0287 / API Contract 5. Settlement and
Task status remain Daemon-owned; the declared Verification command was not run.

### Implementation

- Added `SC-VERIFY-PATH-UNDECLARED` as a Tasks-stage error with no horizon,
  wired through `Check`, `stagedDetectors` and the characterized corpus codes.
  A missing Task Graph skips with its manifest path.
- The detector reads Task declarations and recorded paths, follows only
  dependency-closure `creates:` declarations, accepts their directory prefixes
  and existing paths via `os.Lstat`, and applies the specified shell-word and
  operand exclusions. It reads files and executes no command.
- Added the guide table row, Roundfix reference bullet and `Verification paths`
  authoring rule. Generated the skill mirrors, Roundfix version 0.1.95,
  write-tasks version 0.0.18, their version records and the Behavior Surface
  Record using the repository-owned commands. No mirror, skill version record
  or Behavior Surface Record was hand-edited.
- Recorded the measured new corpus count of 0 and the required update sentence
  in the golden and its characterization pin, preserving every earlier count.

### Acceptance evidence

1. The 0274 shape raises one error naming the path:
   `TestVerificationPathUndeclaredRaisesTheMissingFile` asserts exactly one
   error, the exact API Contract 5 summary and fix, the Task file and the
   Verification command's line. It also exercises command separators,
   redirections, grep/rg explicit patterns, RTK wrappers and optional `./`,
   checks Tasks-stage activation, and asserts a command-execution marker stays
   absent. Before implementation the focused test did not compile because
   `CodeVerifyPathUndeclared` did not exist; after implementation it passed.
2. Supplied paths and exclusions raise nothing:
   `TestVerificationPathUndeclaredIgnoresWhatThePlanProvides` passed for every
   declaration kind, recorded paths, direct and transitive dependency creates,
   directory prefixes, existing files/directories, grep/rg patterns,
   printf/echo operands, negated test/bracket operands, default and configured
   Spec Roots, completed/in-progress/failed Tasks, QA Tasks and absent graphs.
   It also covers variables, globs, ellipses and unrecognized roots/words.
   `TestVerificationPathUndeclaredDoesNotBorrowOtherDeclarations` passed for
   unrelated creations, dependency declarations other than creates, and the
   direction and segment boundary of directory-prefix matching.
3. The active corpus counts 0 and preserves every earlier count:
   the initial focused corpus sweep measured `SC-VERIFY-PATH-UNDECLARED: 0`.
   Its output was compared with the previous golden before recording; every
   earlier count matched. The updated focused corpus and characterization-pin
   tests passed. A final comparison against HEAD confirmed count preservation.

### Focused checks and regeneration

All commands below exited 0 unless explicitly identified as the initial red
check. These are focused implementation checks and required regeneration,
not Daemon settlement Verification.

- Initial red check: `go test ./internal/speccheck -run '^TestVerificationPathUndeclaredRaisesTheMissingFile$' -count=1`
  exited 1 for the absent detector constant.
- `go test ./internal/speccheck -run '^TestVerificationPathUndeclared' -count=1`
- Final detector and adjacent regression selection:
  `go test ./internal/speccheck -run '^Test(VerificationPathUndeclared|SkillVersionGap|WrapFragileVerification)' -count=1`
- `go test -tags docscontract ./internal/docscontract -run '^TestCheckCorpusGolden$' -count=1`
- `go test ./internal/spec -run '^TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095$' -count=1`
- `make skills-sync`
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
- `make baseline-digests`: `ok: true`, `changed: false`.
- `GOOS=windows GOARCH=amd64 go test -c ./internal/speccheck -o /tmp/roundfix-task02-speccheck-windows.test.exe`
  compiled successfully; Windows execution was not performed.
- `git diff --check` passed. A read-only scope and byte comparison found all
  17 changed paths declared by this Task, canonical/generated skill mirrors
  identical, protected QA settlement text unchanged, and all previous corpus
  counts unchanged. The only pre-existing change was the Daemon's
  `pending` to `in_progress` update in this Task file, which was preserved.

No follow-up implementation was added. The Task Graph, other Task files and
Task status were not edited; no commit, push or pull request was created.
