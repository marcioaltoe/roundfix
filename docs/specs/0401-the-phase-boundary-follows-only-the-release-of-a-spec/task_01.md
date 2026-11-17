---
task: task_01
spec: 0401-the-phase-boundary-follows-only-the-release-of-a-spec
status: pending
type: backend
complexity: medium
---

# Task 01: The Boundary Notice follows only a release

## Overview

`roundfix phase start` marks a Phase Record entry and prints the Boundary
Notice only when the same session's latest phase entry is `release`, so a
session that goes from one Authoring Phase to the next sees no notice. The
Roundfix skill's phase reference, the `phase` command page, the Phase Record
section of the usage page and the glossary's **Boundary Notice** describe
the new rule, and the derived records are regenerated.

## Requirements

1. MUST change only the crossing condition of `Recorder.boundary` in
   `internal/phaserecord/record.go` to `_techspec.md` → Interfaces,
   invariant 1: a crossing needs the same session's latest phase entry to
   be `release`, and a `release` whose Spec equals that entry's Spec, both
   empty included, never crosses. The lookup of the latest entry, the
   seven dates it reads, the skipped lines, the written fields and the
   `Boundary` value it returns stay as they are (invariants 2–4).
2. MUST rewrite `TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession`
   in `internal/phaserecord/record_test.go` so that, with every time taken
   from the injected clock, its table holds at least these named cases and
   asserts the returned `Boundary`, the written `boundary` value and, when
   it crosses, `previous_phase` `release` with its `previous_spec`:
   - `after_<phase>_of_the_same_spec` and `after_<phase>_of_another_spec`
     for each of `idea`, `prd`, `techspec`, `tasks` and `delivery`, and
     `after_follow_up`: no crossing;
   - `delivery_then_release_of_another_spec`: no crossing;
   - `after_release_<phase>` for each of the seven phases but `release`:
     crossing;
   - `release_of_another_spec_after_release`: crossing;
   - `release_of_the_same_spec_again` and `release_without_a_spec_again`:
     no crossing;
   - `another_session`, `no_session`, `six_days_old` (crossing after a
     release) and `eight_days_old` (none), and `latest_entry_decides`,
     where a later `follow-up` hides an earlier `release`.
3. MUST change `TestStartPhaseSkipsUnrelatedAndMalformedHistory` and
   `TestStartPhaseReadsAnOversizedLineWithoutBufferingIt` only so that the
   valid earlier entry they must find is a `release` of `0412-example`,
   and their wanted `Boundary` follows it; each keeps what it proves.
4. MUST add `TestPhaseBoundaryNoticeFollowsOnlyARelease` to
   `internal/cli/phase_test.go`, calling `t.Parallel()` first, with the
   subtests `surface_transcript_1` and `surface_transcript_2`. Each asserts
   the whole of the Surface Transcript of `_techspec.md` with that number,
   with its earlier entry recorded through the same command in the same
   fixture: standard output, standard error and the exit code exactly, the
   session read from the fixture's `CLAUDE_CODE_SESSION_ID`, and the
   written line's `boundary` value and, for the second, its
   `previous_phase` and `previous_spec`.
5. MUST change `transcript_3` of `TestPhaseStartSurfaceTranscripts`, whose
   contract moves, so its earlier entry is `release --spec 0412-example`
   and its notice names `release of 0412-example`; `phaseExpectedLine`
   follows it. `TestPhaseBoundaryNoticeOmitsAnAbsentEarlierSpec` keeps its
   assertions.
6. MUST replace, in the canonical `.agents/skills/roundfix/references/phase.md`,
   the paragraph that begins "When the same session last recorded a
   different phase or Spec" with the text of `_techspec.md` → API Contract
   6, word for word, then run `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
7. MUST add to `TestThePhaseReferenceDocumentsTheCommandAndTheRecord` in
   `skills/phase_reference_test.go` two cases on both copies: the new
   paragraph's sentence "When the same session's latest phase entry of the
   last seven days is `release` and the new entry is not a `release` of the
   same Spec" exactly once, and "last recorded a different phase or Spec"
   absent.
8. MUST rewrite the boundary paragraph of `docs/user-guide/commands/phase.md`
   and the matching paragraph of the Phase Record section of
   `docs/user-guide/usage.md` to invariant 1: the notice block shows
   `this session already recorded release of <spec>`, ` of <spec>` is
   omitted when the earlier release names no Spec, a `release` of the same
   Spec again prints none, and no other earlier phase prints one, because a
   Phase Boundary follows a Spec's delivery and release, never an
   Authoring Phase. Both pages drop "A different phase or Spec" and "Three
   exceptions" / "three exceptions". Then run
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
9. MUST record the Coverage Review of the Behavior Surface
   `guide docs/user-guide/usage.md` in `docs/references/skill-coverage.json`:
   its `review` becomes a dated line naming Spec 0401 task_01 and saying
   that only the Phase Record section changed, which the Roundfix skill's
   phase reference covers, as `_prd.md` → Skills records. Its `skills` and
   `sources` stay as they are.
10. MUST revise **Boundary Notice** in `CONTEXT.md` through
   `domain-modeling` to this definition, keeping its `_Avoid_` line: "The
   one line `roundfix phase start` prints on standard error when the same
   session's latest phase entry, in the last seven days of the Phase
   Record, is `release`. It names that release and its Spec and marks the
   new entry. A `release` of the same Spec again prints none, and no other
   earlier phase prints one. It refuses nothing: standard output and the
   exit code are those of the command without it", followed by the links
   to ADR-0407 and ADR-0413.
11. MUST NOT change `internal/cli/phase.go`, the notice text, the record's
    schema, fields or key order, the usage text or any exit code.

## Subtasks

- [ ] Change the crossing condition and rewrite the boundary table test.
- [ ] Move the two history tests to a `release` earlier entry.
- [ ] Assert both Surface Transcripts and move `transcript_3`.
- [ ] Rewrite the phase reference, sync it and record the skill version.
- [ ] Rewrite the two user-guide paragraphs, record the surfaces and the Coverage Review.
- [ ] Revise **Boundary Notice** in the glossary.

## Acceptance Criteria

- [ ] After any earlier phase but `release`, of any Spec, `roundfix phase start` writes `"boundary":false` and prints nothing on standard error.
- [ ] After an earlier `release`, any other phase or a `release` of another Spec is marked and prints the unchanged notice naming that release; the same release again is not.
- [ ] Both Surface Transcripts hold line for line.
- [ ] The phase reference, both user-guide pages and the glossary describe the rule, and none says that a different phase or Spec crosses a boundary.
- [ ] Restoring the old condition makes named subtests of both packages fail.

## Context

- instruction: `docs/adr/0413-a-phase-boundary-follows-only-the-release-of-a-spec.md`
- interface: `internal/phaserecord/record.go`
- interface: `internal/phaserecord/record_test.go`
- interface: `internal/cli/phase_test.go`
- interface: `.agents/skills/roundfix/references/phase.md`
- interface: `skills/roundfix/references/phase.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `skills/phase_reference_test.go`
- interface: `docs/user-guide/commands/phase.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `docs/references/skill-coverage.json`
- interface: `CONTEXT.md`

## Consumers

- The Phase Record's boundary mark and the Boundary Notice — consumers: `./internal/phaserecord`, `./internal/cli`
- The Roundfix skill's phase reference, the Skill Coverage Map and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession$' ./internal/phaserecord 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in after_prd_of_the_same_spec after_tasks_of_another_spec after_delivery_of_the_same_spec delivery_then_release_of_another_spec after_release_follow_up after_release_techspec release_of_another_spec_after_release release_of_the_same_spec_again release_without_a_spec_again latest_entry_decides; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession/$name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/phaserecord` — expected: exit 0. Before this Task the named cases do not exist and the old rule marks a change of phase, so the pass lines are missing; afterwards each case holds and the whole package passes.
- `out="$(go test -count=1 -timeout 30m -v -run '^(TestPhaseBoundaryNoticeFollowsOnlyARelease|TestPhaseStartSurfaceTranscripts|TestPhaseBoundaryNoticeOmitsAnAbsentEarlierSpec)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPhaseBoundaryNoticeFollowsOnlyARelease/surface_transcript_1 TestPhaseBoundaryNoticeFollowsOnlyARelease/surface_transcript_2 TestPhaseStartSurfaceTranscripts/transcript_3 TestPhaseBoundaryNoticeOmitsAnAbsentEarlierSpec; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/cli` — expected: exit 0. Before this Task the transcript test does not exist; afterwards both Surface Transcripts hold, `transcript_3` names the earlier release, and the whole CLI suite passes.
- `out="$(go test -count=1 -timeout 20m -v -run '^TestThePhaseReferenceDocumentsTheCommandAndTheRecord$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/phase.md | grep -qE -- "latest phase entry of the last seven days is .release. and the new entry is not a .release. of the same Spec" || { echo 'phase reference lacks the rule' >&2; exit 1; }; if tr -s '[:space:]' ' ' < docs/user-guide/commands/phase.md | grep -qF -- "a different phase or Spec"; then echo 'phase page keeps the old rule' >&2; exit 1; fi; if tr -s '[:space:]' ' ' < docs/user-guide/usage.md | grep -qF -- "A different phase or Spec crosses a Phase Boundary"; then echo 'usage page keeps the old rule' >&2; exit 1; fi; tr -s '[:space:]' ' ' < docs/user-guide/commands/phase.md | grep -qF -- "this session already recorded release of <spec>" || { echo 'phase page lacks the release notice' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && make skills-sync-check` — expected: exit 0. Before this Task the reference and both pages state the old rule; afterwards they state the new one, and the skill, mirror, version and surface records pass.
- `grep -F -- '**Boundary Notice**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Boundary Notice\*\*(:| —( |$))' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qE -- "latest phase entry, in the last seven days of the Phase Record, is .release." || { echo 'glossary lacks the revised Boundary Notice' >&2; exit 1; }` — expected: exit 0. Before this Task the entry says the notice follows a different phase or Spec; afterwards it follows only a release.

## References

- `_prd.md` → Goals (first and third bullets); Core Features 1 and 4; Success Metrics 1 and 2; Skills; Glossary
- `_techspec.md` → Interfaces (invariants 1–4); API Contracts 1, 2 and 6; Surface Transcripts 1 and 2; Testing Approach (Task 1); Build Order 1; Glossary
- [`references/2026-11-16-the-phase-boundary-ends-the-session-after-every-authoring-phase.md`](references/2026-11-16-the-phase-boundary-ends-the-session-after-every-authoring-phase.md), the Backlog Entry of 2026-11-16 this Spec adopts
- ADR-0413, ADR-0407, ADR-0410
