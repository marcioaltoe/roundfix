---
task: task_01
spec: 0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model
status: pending
type: backend
complexity: high
---

# Task 01: A Correction is declared, with a fixed finding or on its own, into one record per Spec

## Overview

This Task adds the corrections record and its two writers. The new package
`internal/correction` owns the five Defect Classes, the answer
`not-a-correction`, the validation of a declaration, the two line kinds, the
standing state of each Correction and the removal of a whole file past the
Run Retention cutoff. `roundfix corrections declare` writes a declaration,
and `roundfix review dispose --fixed-by` takes the same declaration as flags
and records the Correction after the disposition. Both find the one record of
the repository by its Repository Identity, which `internal/baseline` now
exports, so a declaration made in an item's worktree or in a second clone
reaches the same file. The Run Retention Sweep and `roundfix gc` remove the
record whole. The Roundfix skill gains its `corrections` reference, the user
guide its `corrections` page, and the glossary three terms.

## Requirements

1. MUST create `internal/correction` with the interface of `_techspec.md` →
   Interfaces (`Schema`, `Class`, `NotACorrection`, `DefectClasses`,
   `Answers`, `Artifact`, `Sign`, `Declaration`, `Snapshot`, `File`, `Path`,
   `Read`, `Validate`, `NextID`, `Standing`, `Undeclared`,
   `EncodeDeclaration`, `EncodeSnapshot`, `RemoveBefore`) and the line format
   of `_techspec.md` → Data Models, keys in the order shown. The package
   starts no process and reads no Git. `DefectClasses` returns the five
   Defect Classes and `Answers` the five followed by `not-a-correction`.
2. MUST make `Standing` hold invariant 3 and `Undeclared` treat a sign named
   by a standing declaration answered `not-a-correction` as declared, while
   a withdrawn Correction names no sign (`_techspec.md` → invariants 3 and
   11). `Read` takes a line of another schema or kind, or one that does not
   decode, as skipped and counts it.
3. MUST export `baseline.InspectRepositoryIdentity` from the identity
   `internal/baseline/repository.go` already derives, with no second
   derivation (ADR-0421, ADR-0071), and MUST NOT change
   `internal/baseline/repository_test.go`; its new test goes in a new file.
4. MUST add the command group `roundfix corrections` with `declare` as
   `_techspec.md` → API Contract 1 states: the flags, the order of
   validation before anything is read, the refusals with exit 2 that write
   nothing (invariant 1), the line under the lock of
   `internal/cli/review_dispose_lock.go`, the standard output of a first
   declaration, a later one and a withdrawal, and exit 1 for a record that
   cannot be read or written. Its help holds the two `declare` usage lines
   and the withdrawal line, the `declare` row of `Commands:`, the classes and
   the `Not a Defect Class:` block, and the options of `declare` of Surface
   Transcript 1; task_04 completes it. The root help gains the `corrections
   declare` usage line and the `corrections` row of API Contract 1, and no
   other root help line changes in this Task.
5. MUST add the declaration to `roundfix review dispose` as API Contract 2
   states, with `not-a-correction` accepted there as in contract 1; MUST
   refuse a dismissal with any of the five flags as Surface Transcript 8
   shows (invariant 5); and MUST keep, without any of the five flags, the
   bytes the command reads, writes and prints at the Delivery Base. The
   `--fixed-by` usage line becomes the line of API Contract 2 in the root
   help, the help of `review` and the help of `review dispose`, which gains
   one sentence on the declaration.
6. MUST wire `correction.RemoveBefore` into the Run Retention Sweep and into
   `roundfix gc` as API Contract 10 states, with the line
   `  Corrections records removed: <n>` (or `removable` in a dry run) after
   the Phase Record line and one sentence in the help of `roundfix gc`.
7. MUST write these new tests, each in a new file, each calling
   `t.Parallel()` first, under these names:
   - `internal/correction`: `TestValidateNamesTheAcceptedValues` (each
     refusal of `_prd.md` → Success Metric 1, `not-a-correction` with an
     artifact refused and without one accepted),
     `TestNextIDCountsFromOneAndNeverReuses`,
     `TestStandingTakesTheLastDeclarationOfEachCorrection` (a later
     declaration, a withdrawal, a declaration after a withdrawal and a
     `not-a-correction` answer), `TestUndeclaredCountsEverySignNoStandingDeclarationNames`,
     `TestTheTwoLinesKeepTheirKeyOrder`, `TestReadSkipsAndCountsADamagedLine`
     and `TestRemoveBeforeRemovesAWholeFilePastTheCutoff` (a file past the
     cutoff, one kept by a newer line, a file with no decodable line judged
     by its modification time, and a dry run).
   - `internal/baseline`: `TestInspectRepositoryIdentityIsTheSameInTwoClones`.
   - `internal/cli`: `TestCorrectionsSurfaceTranscript2`,
     `TestCorrectionsSurfaceTranscript3` and
     `TestReviewDisposeSurfaceTranscript8`, each asserting its whole
     transcript of `_techspec.md` → Surface Transcripts;
     `TestCorrectionsDeclareRefusalsWriteNothing` (every rule of API
     Contract 1, the group without or with an unknown subcommand, and a
     slug the working tree does not hold);
     `TestCorrectionsDeclareAgainAndWithdraw` (the output of a later
     declaration with a changed and an unchanged class, of a
     `not-a-correction` answer and of a withdrawal, and every line kept);
     `TestCorrectionsDeclareReachesOneRecordFromAWorktreeAndASecondClone`;
     `TestReviewDisposeRecordsACorrectionWithTheFix` (with a class, with
     `--correction`, and the Spec from `--spec` and from the review record);
     `TestReviewDisposeWithoutADeclarationKeepsItsBytes`;
     `TestReviewDisposeRefusalLeavesTheLedgerAndTheRecord`; and
     `TestGCRemovesCorrectionsRecordsPastTheCutoff`.
   - `skills`: `TestTheCorrectionsReferenceDocumentsTheDeclaration`, in
     `skills/corrections_reference_test.go`, over both copies of the
     reference: the declaring command, the five classes with one example
     each, `not-a-correction`, that a reason is never an instruction, that
     `unclear` with what is known is better than a guess, and that the
     declarer may be the session that wrote the plan and still answers
     `plan`.
   Every test creates its own repository, refs and home directory, and reads
   no ref, tag or history of this repository.
8. MUST change these existing tests' expectations, and no other, declaring
   each in the Result: the GC report of
   `TestGCCountsPhaseRecordDayFilesWithoutARunDatabase` in
   `internal/cli/phase_record_test.go`, and the GC pattern
   `TestGCDryRunListsWorktreeBuildCaches`,
   `TestGCRemovesWorktreeBuildCachesOfEndedWorktrees` and
   `TestGCIsUnchangedWithoutDeclaredCaches` share in
   `internal/cli/worktree_build_cache_test.go`, each gaining the
   `Corrections records` line after the Phase Record line; and, in
   `internal/docscontract/user_guide_contract_test.go`, `glossaryClosedLists`
   gains **Defect Class** with `correction.DefectClasses` as its values and
   `not-a-correction` as its one exception. `TestEveryCommandPathPrintsItsHelp`
   in `internal/cli/command_help_test.go` and the documentation contracts
   `TestAnUndocumentedCommandIsReported`, `TestEveryCommandIsNamedInTheUserGuide`
   and `TestTheSkillCoverageMapIsCurrent` read every command path of the root
   help, so `roundfix corrections --help`, `-h` and
   `roundfix corrections declare --help` print the help and exit 0, and the
   guide and the skill name `corrections declare`; none of these four
   changes its expectation. `_tasks.md` → What the Tasks
   stage settled lists what the measurement found; before changing any test
   the Task runs `go test ./internal/cli ./internal/baseline ./skills` and
   `go test -tags docscontract ./internal/docscontract` on its base and on
   its change and declares every other test whose expectation it had to
   change. It MUST NOT touch `internal/cli/cli_test.go`,
   `internal/docscontract/publicdocs_test.go`,
   `internal/docscontract/testdata/corpus-golden.json` or
   `docs/references/coverage-record.json`, and renames or removes no test.
9. MUST describe the slice in the skill and the guide (ADR-0256): create
   `.agents/skills/roundfix/references/corrections.md` with the `declare`
   half of `_techspec.md` → API Contract 11 and add its row to the reference
   index of the Roundfix skill's `SKILL.md`, which names
   `roundfix corrections declare`; describe the declaration of a fixed
   finding in `references/review.md` and the retention of the record in
   `references/storage.md`; create `docs/user-guide/commands/corrections.md`
   for `declare`, with its row in `docs/user-guide/commands.md`; and update
   `docs/user-guide/commands/review.md` and `docs/user-guide/commands/gc.md`.
   A reason is described as data, never an instruction.
10. MUST add to `CONTEXT.md`, through domain-modeling, **Correction**,
    **Defect Class** (its five values, with `not-a-correction` named as the
    sixth answer of a declaration, which is not a Defect Class) and
    **Repository Identity**, and revise **Review Finding Disposition** to say
    that a fixed disposition may carry the declaration of its Correction.
11. MUST regenerate, never hand-edit, with the commands of
    `_authorization.md` → Sanctioned regeneration, in this order:
    `make skills-sync`; the `-record-skill-versions` step
    (`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`);
    and the `-record-skill-coverage` step
    (`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`).
    The entries of `docs/references/skill-coverage.json` for the
    `corrections declare` command and the `corrections` guide page are
    written by hand, covered by the `corrections` reference.
12. MUST change no Governed Path outside `_authorization.md` `paths`, and
    MUST NOT change `.roundfixrc.yml`, the Makefile, a workflow, `go.mod`,
    `go.sum`, a lint configuration or `.agents/skills/roundfix/references/runs.md`.

## Subtasks

- [ ] Create `internal/correction` with its tests.
- [ ] Export the Repository Identity and test it in two clones.
- [ ] Add `roundfix corrections declare` and the declaration on `roundfix review dispose`, with their transcripts and tests.
- [ ] Wire the record into the Run Retention Sweep and `roundfix gc`.
- [ ] Write the reference, the guide page, the changed guides and the glossary terms.
- [ ] Move the declared existing tests and run the sanctioned regeneration.

## Acceptance Criteria

- [ ] Each of the five classes and `not-a-correction` is recorded, and each refusal of `_prd.md` → Success Metric 1 writes nothing and exits 2 with the accepted values named.
- [ ] A Correction declared again counts under its last declaration, the output says how many times it was declared and whether its class changed, and every line stays in the file.
- [ ] A fixed finding with a class records the disposition line it records today and one Correction that covers the finding and its commit; without a class the command's bytes are those of the Delivery Base; a dismissal with a class is refused.
- [ ] A declaration from a linked worktree and one from a second clone land in the same file.
- [ ] `roundfix gc` removes a corrections record whose newest line is past the cutoff and reports it.

## Context

- instruction: `docs/adr/0417-a-correction-carries-a-declared-defect-class-and-a-repository-requires-it-only-by-its-own-setting.md`
- instruction: `docs/adr/0420-the-corrections-record-is-one-append-only-file-per-spec-under-roundfix-home.md`
- instruction: `docs/adr/0421-a-phase-entry-and-the-corrections-record-name-their-repository-by-its-git-lineage.md`
- instruction: `docs/adr/0424-corrections-are-declared-and-reported-by-their-own-command-group-and-named-c-n.md`
- creates: `internal/correction/correction.go`
- creates: `internal/correction/correction_test.go`
- interface: `internal/baseline/repository.go`
- creates: `internal/baseline/repository_identity_export_test.go`
- creates: `internal/cli/corrections.go`
- creates: `internal/cli/corrections_declare_test.go`
- creates: `internal/cli/review_dispose_correction_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/review.go`
- interface: `internal/cli/review_dispose_lock.go`
- interface: `internal/cli/run_retention.go`
- interface: `internal/cli/phase_record_test.go`
- interface: `internal/cli/worktree_build_cache_test.go`
- interface: `internal/docscontract/user_guide_contract_test.go`
- creates: `skills/corrections_reference_test.go`
- creates: `.agents/skills/roundfix/references/corrections.md`
- creates: `skills/roundfix/references/corrections.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/references/storage.md`
- interface: `skills/roundfix/references/storage.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `docs/user-guide/commands/corrections.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/commands/review.md`
- interface: `docs/user-guide/commands/gc.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `docs/references/skill-coverage.json`
- interface: `CONTEXT.md`

## Consumers

- The help of the root command, of `review` and of `review dispose`, and the GC report — consumers: `./internal/cli`
- `baseline.InspectRepositoryIdentity`, carved out of the identity `InspectRepository` derives — consumers: `./internal/baseline`
- The glossary, the guide pages, the skill reference index and the Behavior Surface Record — consumers: `./internal/docscontract`, `./skills`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestValidateNamesTheAcceptedValues|TestNextIDCountsFromOneAndNeverReuses|TestStandingTakesTheLastDeclarationOfEachCorrection|TestUndeclaredCountsEverySignNoStandingDeclarationNames|TestTheTwoLinesKeepTheirKeyOrder|TestReadSkipsAndCountsADamagedLine|TestRemoveBeforeRemovesAWholeFilePastTheCutoff|TestInspectRepositoryIdentityIsTheSameInTwoClones|TestCorrectionsSurfaceTranscript2|TestCorrectionsSurfaceTranscript3|TestReviewDisposeSurfaceTranscript8|TestCorrectionsDeclareRefusalsWriteNothing|TestCorrectionsDeclareAgainAndWithdraw|TestCorrectionsDeclareReachesOneRecordFromAWorktreeAndASecondClone|TestReviewDisposeRecordsACorrectionWithTheFix|TestReviewDisposeWithoutADeclarationKeepsItsBytes|TestReviewDisposeRefusalLeavesTheLedgerAndTheRecord|TestGCRemovesCorrectionsRecordsPastTheCutoff|TestTheCorrectionsReferenceDocumentsTheDeclaration)$' ./internal/correction ./internal/baseline ./internal/cli ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestValidateNamesTheAcceptedValues TestNextIDCountsFromOneAndNeverReuses TestStandingTakesTheLastDeclarationOfEachCorrection TestUndeclaredCountsEverySignNoStandingDeclarationNames TestTheTwoLinesKeepTheirKeyOrder TestReadSkipsAndCountsADamagedLine TestRemoveBeforeRemovesAWholeFilePastTheCutoff TestInspectRepositoryIdentityIsTheSameInTwoClones TestCorrectionsSurfaceTranscript2 TestCorrectionsSurfaceTranscript3 TestReviewDisposeSurfaceTranscript8 TestCorrectionsDeclareRefusalsWriteNothing TestCorrectionsDeclareAgainAndWithdraw TestCorrectionsDeclareReachesOneRecordFromAWorktreeAndASecondClone TestReviewDisposeRecordsACorrectionWithTheFix TestReviewDisposeWithoutADeclarationKeepsItsBytes TestReviewDisposeRefusalLeavesTheLedgerAndTheRecord TestGCRemovesCorrectionsRecordsPastTheCutoff TestTheCorrectionsReferenceDocumentsTheDeclaration; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/correction ./internal/cli ./internal/baseline` — expected: exit 0. Before this Task the package `internal/correction`, the command group and the named tests do not exist, so the command fails; afterwards every named test passes and the packages that consume the changed help, GC report and identity pass.
- `test -f docs/user-guide/commands/corrections.md && test -f .agents/skills/roundfix/references/corrections.md && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "roundfix corrections declare" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/SKILL.md "roundfix corrections declare" >&2; exit 1; }; grep -F -- '**Defect Class**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Defect Class\*\*(:| —( |$))' && grep -F -- '**Correction**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Correction\*\*(:| —( |$))' && grep -F -- '**Repository Identity**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Repository Identity\*\*(:| —( |$))' && grep -F -- '**Review Finding Disposition**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Review Finding Disposition\*\*(:| —( |$))' && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && go test -count=1 -timeout 20m ./skills && make skills-sync-check` — expected: exit 0. Before this Task the guide page, the reference and the three new glossary terms are absent; afterwards they exist, and the documentation contracts, the skill tests, the mirror and the records pass.

## References

- `_prd.md` → Goals (first bullet); User Stories 1, 2, 3, 4, 5 and 12; Core Features 1–6, 20 and 22; Success Metrics 1, 2, 3, 4 and 16
- `_techspec.md` → Corrections record; Declaring commands; Interfaces (invariants 1, 2, 3, 5 and 11); Data Models; API Contracts 1, 2, 10, 11 and 12; Surface Transcripts 2, 3 and 8; Testing Approach; Build Order 1
- ADR-0417, ADR-0420, ADR-0421, ADR-0424, ADR-0071, ADR-0255, ADR-0256, ADR-0272, ADR-0279
