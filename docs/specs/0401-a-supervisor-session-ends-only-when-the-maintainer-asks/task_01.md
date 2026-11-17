---
task: task_01
spec: 0401-a-supervisor-session-ends-only-when-the-maintainer-asks
status: pending
type: backend
complexity: medium
---

# Task 01: The Phase Record marks no boundary and the phase command prints no notice

## Overview

`roundfix phase start` keeps recording every phase but stops reading the
record back: it writes `"boundary":false` with no `previous_phase` or
`previous_spec`, and prints nothing on standard error, so a session that
goes from phase to phase or from Spec to Spec sees no notice. The reader
the status view uses still reads entries an earlier version marked. The
Roundfix skill's phase reference, the `phase` command page and the Phase
Record section of the usage page describe the record without a boundary,
and the derived records are regenerated.

## Requirements

1. MUST change `internal/phaserecord/record.go` to `_techspec.md` →
   Interfaces: `StartPhase` returns only an `error`; the type `Boundary`,
   the method `Recorder.boundary`, `readDay` and `maxLineBytes` are
   removed; `phaseLine` drops `PreviousPhase` and `PreviousSpec` and keeps
   `Boundary`, always `false`, in its place. The time, the field bounding,
   the directory and file modes and the append stay as they are
   (invariants 1 and 2). `internal/phaserecord/read.go` does not change.
2. MUST change `internal/cli/phase.go` only to call the new `StartPhase`
   and to drop the block that prints the notice, so a successful append
   prints only the result line (invariant 4). The usage text, the
   refusals, their messages and the exit codes stay as they are.
3. MUST follow the new signature in `internal/cli/follow_up_queue_test.go`
   and in the `written_by_start_phase` case of
   `internal/phaserecord/read_test.go`, changing nothing else they assert.
4. MUST replace `TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession`
   in `internal/phaserecord/record_test.go` with
   `TestStartPhaseMarksNoBoundaryWhateverCameBefore`, whose table, with
   every time taken from an injected clock, holds at least these named
   cases, each recording its earlier entries and then a new one for the
   session `s`:
   - `after_follow_up`, `after_idea`, `after_prd`, `after_techspec`,
     `after_tasks`, `after_delivery` and `after_release`: one earlier
     entry of that phase of the same session, of `0412-example` where the
     phase takes a Spec, and a new `techspec` of `0412-example`, or a new
     `follow-up` after `release`;
   - `after_prd_of_another_spec`, `release_of_another_spec_after_release`,
     `release_of_the_same_spec_again` and `release_without_a_spec_again`;
   - `after_a_marked_entry`: a hand-written earlier line with
     `"boundary":true`, `previous_phase` and `previous_spec`, as Roundfix
     0.147.0 wrote it.

   Each case asserts that `StartPhase` returns `nil`, that the new line's
   keys are exactly the nine of `phaseKeys` in order with `"boundary":false`
   and the phase, Spec and session it was given, and that every earlier
   line is byte-for-byte unchanged.
5. MUST replace `TestStartPhaseSkipsUnrelatedAndMalformedHistory` and
   `TestStartPhaseReadsAnOversizedLineWithoutBufferingIt` with
   `TestStartPhaseReadsNoEarlierEntry`, with the subtests
   `an_unreadable_earlier_day`, where a directory stands at the path of the
   day file of the day before the clock's UTC date, and
   `a_malformed_day_file`, where the clock's day file already holds a line
   of 1 MiB that is not JSON, ended by a newline, and a `release` entry of
   the same session. Each asserts that `StartPhase` returns `nil` and
   appends exactly one line with `"boundary":false` and the nine keys.
6. MUST rename the subtest `boundary_keys_only_at_a_boundary` of
   `TestAPhaseEntryIsOneLineOfItsDeclaredKeys` to
   `no_boundary_keys_after_a_release` and make it assert that a
   `follow-up` recorded after a `release` without a Spec of the same
   session has exactly the nine keys and `"boundary":false`. The other
   subtests keep what they assert, with the new signature.
7. MUST add to `TestPhasesReadsBackThePhaseEntriesOfTheWindow` in
   `internal/phaserecord/read_test.go` the case
   `an_entry_an_earlier_version_marked`: the day file also holds, at
   `since`, an `idea` entry of `0001-widget` in `fixture` written with
   `"boundary":true,"previous_phase":"prd","previous_spec":"0001-widget"`,
   and `Phases` returns it, oldest first, beside the other entry, with
   `skipped` `0` (invariant 3).
8. MUST, in `internal/cli/phase_test.go`:
   - change `transcript_3` of `TestPhaseStartSurfaceTranscripts` so that,
     after the earlier `prd` of `0412-example`, the `techspec` command
     prints nothing on standard error and its line holds
     `"boundary":false`; `phaseExpectedLine` writes `"boundary":false` and
     no `previous_` key for every caller;
   - replace `TestPhaseBoundaryNoticeOmitsAnAbsentEarlierSpec` with
     `TestPhaseStartPrintsNoNoticeAfterAnyPhase`, calling `t.Parallel()`
     first, with the subtests `surface_transcript_1`,
     `surface_transcript_2` and `after_a_release_without_a_spec`. The first
     two assert the whole of the Surface Transcript of `_techspec.md` with
     that number, with its earlier entry recorded through the same command
     in the same fixture: standard output, an empty standard error and the
     exit code exactly, the session read from the fixture's
     `CLAUDE_CODE_SESSION_ID`, and the written line equal to
     `phaseExpectedLine`. The third records `release` with no Spec and
     then `follow-up`, and asserts the same.
9. MUST replace, in the canonical `.agents/skills/roundfix/references/phase.md`,
   the paragraph that begins "When the same session last recorded a
   different phase or Spec" with the text of `_techspec.md` → API Contract
   6, word for word, then run `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
10. MUST change `TestThePhaseReferenceDocumentsTheCommandAndTheRecord` in
    `skills/phase_reference_test.go` so that, on both copies, the case
    `the_notice_refuses_nothing` gives way to
    `the_record_is_never_read_back`, which finds "The command never reads
    the record back: whatever the session recorded before, it writes
    `boundary` as `false` and prints nothing on standard error." exactly
    once, and the case `no_boundary_notice` finds "Boundary Notice" zero
    times. The other cases keep their phrases.
11. MUST rewrite, in `docs/user-guide/commands/phase.md`, the paragraphs
    from "When the same non-empty session last recorded" to "stay as they
    are." as: "The command reads no earlier entry. Whatever the session
    recorded before, it writes the entry's `boundary` as `false` and prints
    nothing on standard error. A phase entry ends no session.", and its exit
    `0` row as "One phase entry recorded".
12. MUST rewrite, in the Phase Record section of `docs/user-guide/usage.md`:
    - the `boundary` row as "Always `false`; an entry an earlier Roundfix
      wrote may hold `true`", and the `previous_phase` and `previous_spec`
      rows as "Absent; an entry an earlier Roundfix marked holds the
      earlier phase" and "… the earlier Spec";
    - "only the two boundary fields are conditional." as "a new entry
      never holds `previous_phase` or `previous_spec`.";
    - the paragraphs from "Before appending, the command takes" to "stay
      as they are." as: "The command reads no earlier entry. Whatever the
      session recorded before, it writes `boundary` as `false` and prints
      nothing on standard error. A phase entry ends no session: a
      Supervisor session goes on from phase to phase and ends only when the
      user asks or the work stops."

    Then run
    `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
13. MUST record the Coverage Review of the Behavior Surface
    `guide docs/user-guide/usage.md` in `docs/references/skill-coverage.json`:
    its `review` becomes a dated line naming Spec 0401 task_01 and saying
    that only the Phase Record section changed, which the Roundfix skill's
    phase reference covers, as `_prd.md` → Skills records. Its `skills` and
    `sources` stay as they are.
14. MUST NOT change the record's schema, the key order, the reader
    `Phases`, the usage text, any refusal or exit code, `CONTEXT.md` or any
    ADR.

## Subtasks

- [ ] Remove the writer's lookup and the notice, and follow the new signature.
- [ ] Rewrite the writer's table, history and key tests, and add the reader's marked-entry case.
- [ ] Assert both Surface Transcripts and move `transcript_3`.
- [ ] Rewrite the phase reference, sync it and record the skill version.
- [ ] Rewrite the two user-guide passages, record the surfaces and the Coverage Review.

## Acceptance Criteria

- [ ] After any earlier entry of the same session, of any phase and Spec or marked by an earlier version, `roundfix phase start` writes `"boundary":false` with the nine keys and prints nothing on standard error.
- [ ] The writer appends even where an earlier day file cannot be read.
- [ ] `Phases` still reads an entry an earlier version marked.
- [ ] Both Surface Transcripts hold line for line.
- [ ] The phase reference and both user-guide pages describe the record without a boundary, and none mentions the notice.
- [ ] Restoring the lookup, the mark or the notice makes named subtests of both packages fail.

## Context

- instruction: `docs/adr/0414-a-supervisor-session-ends-only-when-the-maintainer-asks.md`
- interface: `internal/phaserecord/record.go`
- interface: `internal/phaserecord/record_test.go`
- interface: `internal/phaserecord/read_test.go`
- interface: `internal/cli/phase.go`
- interface: `internal/cli/phase_test.go`
- interface: `internal/cli/follow_up_queue_test.go`
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

## Consumers

- The Phase Record's writer, its reader and the phase command — consumers: `./internal/phaserecord`, `./internal/cli`
- The Roundfix skill's phase reference, the Skill Coverage Map and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestStartPhaseMarksNoBoundaryWhateverCameBefore|TestStartPhaseReadsNoEarlierEntry|TestAPhaseEntryIsOneLineOfItsDeclaredKeys|TestPhasesReadsBackThePhaseEntriesOfTheWindow)$' ./internal/phaserecord 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestStartPhaseMarksNoBoundaryWhateverCameBefore/after_prd TestStartPhaseMarksNoBoundaryWhateverCameBefore/after_release TestStartPhaseMarksNoBoundaryWhateverCameBefore/release_of_another_spec_after_release TestStartPhaseMarksNoBoundaryWhateverCameBefore/release_of_the_same_spec_again TestStartPhaseMarksNoBoundaryWhateverCameBefore/after_a_marked_entry TestStartPhaseReadsNoEarlierEntry/an_unreadable_earlier_day TestStartPhaseReadsNoEarlierEntry/a_malformed_day_file TestAPhaseEntryIsOneLineOfItsDeclaredKeys/no_boundary_keys_after_a_release TestPhasesReadsBackThePhaseEntriesOfTheWindow/an_entry_an_earlier_version_marked; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/phaserecord` — expected: exit 0. Before this Task the named cases do not exist and the writer marks a change of phase, so the pass lines are missing; afterwards each case holds and the whole package passes.
- `out="$(go test -count=1 -timeout 30m -v -run '^(TestPhaseStartPrintsNoNoticeAfterAnyPhase|TestPhaseStartSurfaceTranscripts)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPhaseStartPrintsNoNoticeAfterAnyPhase/surface_transcript_1 TestPhaseStartPrintsNoNoticeAfterAnyPhase/surface_transcript_2 TestPhaseStartPrintsNoNoticeAfterAnyPhase/after_a_release_without_a_spec TestPhaseStartSurfaceTranscripts/transcript_3; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; if grep -qF -- "roundfix: notice:" internal/cli/phase.go; then echo 'phase command keeps the notice' >&2; exit 1; fi; go test -count=1 -timeout 30m ./internal/cli` — expected: exit 0. Before this Task the new test does not exist and the command prints the notice; afterwards both Surface Transcripts hold, `transcript_3` prints no notice, and the whole CLI suite passes.
- `out="$(go test -count=1 -timeout 20m -v -run '^TestThePhaseReferenceDocumentsTheCommandAndTheRecord$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in canonical/the_record_is_never_read_back embedded/the_record_is_never_read_back canonical/no_boundary_notice embedded/no_boundary_notice; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestThePhaseReferenceDocumentsTheCommandAndTheRecord/$name " || { echo "missing pass: $name" >&2; exit 1; }; done; for f in docs/user-guide/commands/phase.md docs/user-guide/usage.md; do if tr -s '[:space:]' ' ' < "$f" | grep -qiE -- "boundary notice|a different phase or Spec|continue in a new session"; then echo "old rule in $f" >&2; exit 1; fi; tr -s '[:space:]' ' ' < "$f" | grep -qF -- "The command reads no earlier entry." || { echo "$f lacks the new rule" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && make skills-sync-check` — expected: exit 0. Before this Task the reference and both pages describe the notice; afterwards they describe the record without a boundary, and the skill, mirror, version and surface records pass.

## References

- `_prd.md` → Goals (first and third bullets); Core Features 1 and 4; Success Metrics 1 and 2; Skills
- `_techspec.md` → Interfaces (invariants 1–4); API Contracts 1, 2 and 6; Surface Transcripts 1 and 2; Testing Approach (Task 1); Build Order 1
- [`references/2026-11-16-the-phase-boundary-ends-the-session-after-every-authoring-phase.md`](references/2026-11-16-the-phase-boundary-ends-the-session-after-every-authoring-phase.md), the Backlog Entry of 2026-11-16 this Spec adopts
- ADR-0414, ADR-0407, ADR-0410
