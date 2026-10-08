---
task: task_04
spec: 0251-skills-keep-up-with-the-behavior-they-describe
status: completed
type: backend
complexity: high
---

# Task 04: Spec authoring requires a skills Task for a changed Behavior Surface

## Overview

Adds the authoring rule to the Spec Consistency Check. `SC-SKILLS-UNTASKED`
reports a covered Behavior Surface whose source a non-QA Task declares when no
non-QA Task declares one of its covering skill files and the PRD's Skills
Declaration does not excuse it; `SC-SKILLS-MALFORMED` reports a malformed
Skills Declaration. The write-tasks and write-prd skills, the Roundfix skill's
spec reference, the spec command guide, the repository rules and the glossary
state the rule. The slice is verifiable on its own through fixture Specs.

## Requirements

1. MUST add `internal/speccheck/skills_declaration.go` exporting
   `CodeSkillsUntasked = "SC-SKILLS-UNTASKED"` and
   `CodeSkillsMalformed = "SC-SKILLS-MALFORMED"`, and implement
   `_techspec.md` → API Contract 7 and API Contract 8 exactly, reading the
   worktree map with `internal/skillcoverage` and its horizon from the commit
   that added the PRD and the oldest commit that added the map, in the manner
   of the glossary horizon. Both findings are errors.
2. MUST register `SC-SKILLS-MALFORMED` at the PRD stage and
   `SC-SKILLS-UNTASKED` at the Tasks stage in the staged detector list, run
   both in the full check, run the malformed detector in the PRD and TechSpec
   authoring stages, and list both as skipped wherever the PRD or the Task
   Graph is absent, as the glossary detectors are.
3. MUST add these tests to `internal/speccheck/skills_declaration_test.go`,
   each on a fixture repository in `t.TempDir()` with a committed map, PRD and
   Task Graph:
   - `TestSkillsTaskIsRequiredForAChangedSurface`: a Task declaring a covered
     command's source with no skills Task reports one `SC-SKILLS-UNTASKED`
     naming the surface, the Task and the covering file, and the rendered text
     holds the line prefix `[error] SC-SKILLS-UNTASKED: ` of Surface
     Transcript 4; adding a Task that declares the covering file clears it; a
     source of an uncovered surface never reports.
   - `TestAnUnchangedDeclarationExcusesASurface`: an `- unchanged:` entry
     clears the finding only together with a non-QA Task that declares the
     map.
   - `TestAMalformedSkillsDeclarationIsReported`: a line of another shape, an
     unknown surface and an empty reason each report one
     `SC-SKILLS-MALFORMED` at their line; a PRD without the section reports
     none.
   - `TestTheSkillsRuleStartsAtTheMapHorizon`: a PRD committed before the map
     and a repository without a map list both codes as skipped with their
     reasons; an uncommitted PRD is held.
4. MUST add both codes at 0 to the corpus characterization: the corpus code
   list of `internal/docscontract/corpus_test.go`,
   `internal/docscontract/testdata/corpus-golden.json` and its pin in
   `internal/spec/archive_layout_characterization_test.go`, with the same
   sentence appended to the `update` text in both, changing no existing count.
5. MUST state the rule in the write-tasks skill under a new heading of its
   own, naming `SC-SKILLS-UNTASKED`, the skills Task, and the
   `- unchanged:` entry with its Task that records the Coverage Review; add an
   optional `## Skills` section to the write-prd PRD template; and list both
   codes in the Roundfix skill's spec reference. No text inside any
   `### QA settlement` section changes. Each skill edit MUST be synced with
   `make skills-sync` and recorded with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which chooses the next version. `TestWriteTasksSkillStatesTheDeclaredPathRules`
   MUST require the phrase `SC-SKILLS-UNTASKED`.
6. MUST list both codes in the code table of `docs/user-guide/commands/spec.md`,
   and rewrite the roundfix skill sync rule of
   `docs/agents/specific-repository.md` to name the Skill Coverage Map, the
   release plan's `skill-coverage` check and `SC-SKILLS-UNTASKED`, keeping its
   other rules byte-identical.
7. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   after the guide change.
8. MUST write the glossary term **Skills Declaration** into `CONTEXT.md`
   through the `domain-modeling` skill, citing ADR-0256.

## Subtasks

- [ ] Implement the two detectors and their horizon, and register them.
- [ ] Write the four fixture tests.
- [ ] Characterize the two codes in the corpus golden and its pin.
- [ ] Update the write-tasks, write-prd and Roundfix skills; sync and record.
- [ ] Update the spec guide and the repository rule; re-record the record.
- [ ] Add the glossary term.

## Acceptance Criteria

- [ ] A Spec whose Task changes a covered surface without a skills Task fails `roundfix spec check`, and a skills Task or an excused entry with its map Task clears it.
- [ ] A malformed Skills Declaration is reported at its line.
- [ ] Specs before the map's horizon and repositories without a map are skipped with a reason.
- [ ] The active corpus golden carries both codes at 0 and every other count unchanged.
- [ ] The skills, the guide, the repository rule and the glossary state the rule.

## Context

- creates: `internal/speccheck/skills_declaration.go`
- creates: `internal/speccheck/skills_declaration_test.go`
- interface: `internal/speccheck/coherence.go`
- interface: `internal/speccheck/constraints.go`
- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`
- interface: `internal/spec/archive_layout_characterization_test.go`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-prd/references/prd-template.md`
- interface: `skills/write-prd/references/prd-template.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `skills/baseline_skill_contract_test.go`
- interface: `docs/user-guide/commands/spec.md`
- creates: `docs/references/behavior-surfaces.json`
- interface: `docs/agents/specific-repository.md`
- interface: `CONTEXT.md`
- instruction: `internal/speccheck/glossary.go`
- instruction: `internal/speccheck/surface.go`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestSkillsTaskIsRequiredForAChangedSurface|TestAnUnchangedDeclarationExcusesASurface|TestAMalformedSkillsDeclarationIsReported|TestTheSkillsRuleStartsAtTheMapHorizon)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestSkillsTaskIsRequiredForAChangedSurface TestAnUnchangedDeclarationExcusesASurface TestAMalformedSkillsDeclarationIsReported TestTheSkillsRuleStartsAtTheMapHorizon; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done; docs="$(go test -count=1 -tags docscontract -run '^(TestCheckCorpusGolden|TestTheSkillCoverageMapIsCurrent)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$docs"; exit 1; }; go test -count=1 -run '^TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095$' ./internal/spec || exit 1; go test -count=1 ./skills -run '^(TestEveryOwnedSkillVersionIsRecorded|TestWriteTasksSkillStatesTheDeclaredPathRules)$' || exit 1; grep -qF -- 'SC-SKILLS-UNTASKED' internal/docscontract/testdata/corpus-golden.json || { printf 'the corpus golden does not characterize SC-SKILLS-UNTASKED\n' >&2; exit 1; }; for file in .agents/skills/write-tasks/SKILL.md .agents/skills/roundfix/references/spec.md docs/user-guide/commands/spec.md docs/agents/specific-repository.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'SC-SKILLS-UNTASKED' || { printf 'missing SC-SKILLS-UNTASKED in %s\n' "$file" >&2; exit 1; }; done; grep -q -- '^## Skills$' .agents/skills/write-prd/references/prd-template.md || { printf 'the PRD template has no Skills section\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Skills Declaration**:' || { printf 'missing glossary term Skills Declaration\n' >&2; exit 1; }` — expected: exit 0. Before this Task the four detector tests do not exist, so the first check fails. After it, the detectors pass their fixtures, the corpus golden and its pin carry both codes, the skill record and the write-tasks contract pass, and the skills, guide, rule and glossary state the rule.

## References

- `_prd.md` → Core Feature 5; Core Feature 6; User Story 2; Success Metric 5; Success Metric 6; Glossary
- `_techspec.md` → API Contract 7; API Contract 8; Invariant 7; Surface Transcript 4; Build Order 4
- ADR-0256; ADR-0222; ADR-0187

## Result

Implemented this Task's authoring-rule slice. Task status and authored Verification
remain Daemon-owned; no terminal Task verdict is recorded here.

### Implementation and acceptance evidence

1. **Covered source declarations require a skills Task or a reviewed excuse.**
   Added the error detectors in `internal/speccheck/skills_declaration.go`,
   registered malformed declarations at the PRD stage and untasked surfaces at
   the Tasks stage, and wired the full and authoring checks. Source matching
   uses `skillcoverage.Map.SurfacesForPath`, including prefixes and globs.
   Non-QA `interface:`, `creates:` and `deletes:` source declarations bind a
   surface; only `interface:` or `creates:` skill/map declarations excuse it.
   `TestSkillsTaskIsRequiredForAChangedSurface` exercises the finding, its
   surface/Task/path/covering-file summary, the transcript's error prefix,
   skill updates and uncovered sources. `TestAnUnchangedDeclarationExcusesASurface`
   exercises the excuse's required map Task. Additional QA exclusion fixtures
   prove that a QA Task cannot bind a source, update a skill or record a review.
2. **Malformed declarations identify their PRD line.**
   `TestAMalformedSkillsDeclarationIsReported` exercises another line shape,
   an unknown surface and an empty reason at PRD line 7 through PRD, TechSpec,
   Tasks and full checks. An absent optional section has no malformed finding.
3. **The map horizon preserves historical Specs.**
   `TestTheSkillsRuleStartsAtTheMapHorizon` exercises absent-map skips,
   pre-map PRD skips, held uncommitted PRDs and absent-PRD skips.
   `TestSkillsDeclarationUsesOldestMapAddition` proves deleting and re-adding
   the map does not move its introduction horizon. Missing Task Graphs list
   both codes as skipped with the graph path.
4. **The corpus characterization gains only two zero counts.**
   Added both codes to the corpus code list, golden and archive-layout pin,
   with the same appended update sentence. A focused Python comparison against
   `HEAD` confirmed every existing count is unchanged, the only added keys are
   the two new codes at 0, and the pin agrees on every count and update text.
   Focused current-corpus checks found no errors.
5. **Authoring skills, guide, repository rule and glossary state the rule.**
   Added the write-tasks heading, optional PRD template section, Roundfix spec
   reference and command-guide code rows. Updated only the repository's skill
   sync rule and added **Skills Declaration** to `CONTEXT.md`, citing ADR-0256.
   Synchronized the skills and used the owned-skill recorder to choose the
   next versions. Re-recorded the Behavior Surface Record after the guide
   change; only the spec-guide fingerprint changed. Focused byte comparisons
   confirmed matching edited mirrors, unchanged QA settlement sections and
   unchanged repository rules outside the skill sync clause.

### Focused checks and required generators

- Initial focused check:
  `rtk proxy go test ./internal/speccheck -run TestSkillsTaskIsRequiredForAChangedSurface -count=1`
  exited 1 before implementation because the two new exported codes were absent.
- `rtk proxy go test ./internal/speccheck -run 'Test(Skills|AnUnchanged|AMalformedSkills|TheSkills)' -count=1`
  exited 0 after the detector implementation.
- `rtk proxy go test ./internal/speccheck -count=1` exited 0 after all detector
  and fixture edits, including the added QA and oldest-addition cases.
- `rtk make skills-sync` exited 0.
- `rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  exited 0; the recorder chose the versions and updated both skill copies.
- `rtk proxy go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  exited 0 and regenerated the record. The first sandbox attempt was denied a
  Go build-cache read; the authorized retry succeeded.
- `rtk proxy go test -tags docscontract ./internal/docscontract -run 'TestCheckActiveCorpusHasNoErrors|TestCheckCorpusBudget' -count=1`
  exited 0.
- `rtk proxy go test ./skills -run 'TestAuthorialSkillSync|TestEmbeddedSkillsMatchCanonical|TestOwned' -count=1`
  exited 0.
- Focused Python artifact comparisons exited 0 for corpus/pin equality,
  preserved counts, repository-rule preservation, edited skill mirrors,
  QA settlement preservation, the ADR-backed glossary term and unchanged
  authored Task content/status.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

The declared `## Verification` command and repository-wide delivery gates were
not run in this Agent turn. The Daemon owns Verification and settlement.
No commit, push, Pull Request, Task Graph edit or other Task edit was made.

### Verification Feedback repair — Attempt 1

Inspected the Daemon diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261008T154220Z_960340e21200b0bc/verification/batch-004-attempt-1.log`.
The configured gate reached the owned-skill regeneration contract, whose
repository-copy helper omitted untracked source files. Its copied checker
therefore referenced Task 04's new detector without containing its source.
This was a fixture-copy defect at the pre-commit boundary, rather than a
stale derived record or a detector failure.

Repaired only the already-declared `skills/baseline_skill_contract_test.go`:
`copyTrackedRepository` now copies tracked files and non-ignored new files
using `git ls-files --cached --others --exclude-standard`. Added
`TestRepositoryCopyIncludesNewSourceFiles` to prove new sources and modified
tracked bytes reach the copy while ignored files and Git metadata stay out.
The existing regeneration assertions remain intact.

Focused repair evidence:

- `rtk proxy go test ./skills -run '^TestRepositoryCopyIncludesNewSourceFiles$' -count=1`
  exited 1 before the helper repair because the copied new source was absent,
  then exited 0 after the repair.
- `rtk proxy go test -tags repocontract ./skills -run '^TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical$' -count=1`
  exited 0 after the repair, exercising the exact contract identified by the
  Daemon feedback. The first sandbox invocation was denied a Go cache read;
  the authorized retry succeeded.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

Only the declared test file and this Result addendum changed during the
repair turn. Task status, other Tasks and the Task Graph remain untouched.
The Task's declared Verification command and `make verify-changed` were not
rerun. Full configured Verification and settlement remain with the Daemon;
no commit, push or Pull Request was made.
