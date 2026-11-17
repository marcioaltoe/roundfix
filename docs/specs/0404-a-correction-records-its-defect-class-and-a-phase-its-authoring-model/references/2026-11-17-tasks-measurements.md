# What the Tasks stage of Spec 0404 measured (2026-11-17)

Recorded by the author of the Tasks on the Mac mini on 2026-11-17, between about 17:15 and 17:40 BRT, in three scratch worktrees of the clone on the branch `docs/plan-0404`, detached at `d65bae06` (the TechSpec amendment on `a8236a87`), and removed afterwards. The Tasks were authored on Opus (`claude-opus-5-5`), as the maintainer's answer "Terminar em Opus agora (Recommended)" decided. Specs 0401, 0402 and 0403 were not on the default branch, so every number below is on a tree without them.

## 1. The existing tests a prototype of the surfaces moves

A prototype changed only what existing tests can observe: the root help gained the two `corrections` usage lines and the `corrections` row, the `--fixed-by` line of `review dispose` and the `phase start` line took their new form, `roundfix corrections` was dispatched to a stub that refuses, `roundfix phase --help` took the text of Surface Transcript 6, an Authoring Phase entry gained the six keys with `unknown` sources, its result line with a session gained the model, effort and runner, the GC report gained `  Corrections records <action>: 0` after the Phase Record line, `ClassifyPark` gained the new blocker in the `finding` class, and Project Config gained `delivery.require_defect_class`. No guide, skill or test changed, and the merge step and the sign reader were not built.

`go test -count=1 -timeout 30m` over `./internal/cli ./internal/phaserecord ./internal/delivery ./internal/config ./internal/baseline ./internal/runcause ./skills` and `go test -tags docscontract ./internal/docscontract` all passed on the unchanged tree. On the prototype:

- `./internal/runcause`, `./skills`, `./internal/config`, `./internal/delivery` and `./internal/baseline` passed. `TestEveryBlockerHasAParkClass` and `TestParkClassesAndNextActions` passed with the new case.
- `./internal/cli` failed in:
  - `TestEveryCommandPathPrintsItsHelp`, for `corrections --help`, `corrections -h`, `corrections declare --help`, `corrections declare -h`, `corrections report --help` and `corrections report -h`: the test reads every command path of the root help, and the stub printed no help. A command that prints its help passes it unchanged.
  - The GC report: `TestGCCountsPhaseRecordDayFilesWithoutARunDatabase` in `phase_record_test.go`, whose whole output is a literal, and `TestGCDryRunListsWorktreeBuildCaches`, `TestGCRemovesWorktreeBuildCachesOfEndedWorktrees` and `TestGCIsUnchangedWithoutDeclaredCaches` in `worktree_build_cache_test.go`, which share one pattern. The substring checks of `gc_run_retention_test.go` and `TestTheRunRetentionSweepRemovesPhaseRecordDayFiles` passed.
  - The phase command: `TestPhaseStartSurfaceTranscripts/transcript_1`, `transcript_2` and `transcript_3`; `TestPhaseStartRecordsEachOfTheSevenPhases/idea`, `prd`, `techspec` and `tasks`; and `TestPhaseStartLoadsNoConfiguration`.
- `./internal/phaserecord` failed in the subtests of `TestAPhaseEntryIsOneLineOfItsDeclaredKeys` and in the boundary tests Spec 0401 replaces. `TestPhasesReadsBackThePhaseEntriesOfTheWindow` passed, because the prototype did not add the six fields to `RecordedPhase`.
- `./internal/docscontract` failed in `TestAnUndocumentedCommandIsReported`, `TestEveryCommandIsNamedInTheUserGuide`, `TestTheSkillCoverageMapIsCurrent` and `TestBehaviorSurfaceBytesComeFromParsedValuesAndGuideContent`, each because the new command paths had neither a help, a guide page nor a skill reference; documenting them is the Task's work, not a changed expectation.
- `TestReviewDisposeHelpPrintsUsage` passed: the new `--fixed-by` line begins with the present one.

## 2. The files the sanctioned regeneration rewrites

A second worktree appended one probe sentence to the canonical `deliver`, `phase`, `review`, `spec-delivery` and `storage` references and the `SKILL.md` of the Roundfix skill, the four authoring skills and the six user-guide pages the Tasks change, and created the `corrections` reference and guide page. Then it ran the three commands of the authorization record's sanctioned regeneration:

- `make skills-sync` exited 0 and wrote the mirrors under `skills/`, the new `corrections` reference among them.
- The `-record-skill-versions` step exited 0 and wrote `skills/testdata/owned-skill-versions.json` and the version fields of the five `SKILL.md` pairs.
- The `-record-skill-coverage` step wrote `docs/references/behavior-surfaces.json` and failed with `guide docs/user-guide/commands/corrections.md: missing Skill Coverage Map entry`: the new guide page needs its entry in `docs/references/skill-coverage.json`, written by hand.

Nothing under `.claude/agents/` changed, and no other file changed.

## 3. The Governed Paths the Tasks declare

A `go test -overlay` probe in `internal/speccheck`, which wrote nothing to the repository, ran `speccheck.GovernedPath` over the 92 distinct paths the five Task files declare. Exactly 16 are governed, the 16 of the authorization record's `paths`. `git status --porcelain` showed only the new Task files afterwards. The record needs no amendment.
