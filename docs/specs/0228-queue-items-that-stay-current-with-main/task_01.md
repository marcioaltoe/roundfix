---
task: task_01
spec: 0228-queue-items-that-stay-current-with-main
status: pending
type: backend
complexity: medium
---

# Task 01: The owned-skill record command raises a colliding version

## Overview

On 2026-10-04 two queued items both recorded the Roundfix Skill at `0.1.26`,
and a Codex Task agent's hand-written digest made the record command fail
every later attempt with "content changed under version" (Backlog Entry
[queued Specs raise the same skill version](references/2026-10-04-queued-specs-raise-the-same-skill-version.md),
recorded 2026-10-04; operator log entries 150, 154 and 157). This Task makes
the record command choose the version: content not recorded under a version
that is not above every recorded one is raised to the next free patch in both
version fields of the skill and its mirror, then recorded (ADR-0233). The
derived merge of task_02 runs the same command.

## Requirements

1. MUST implement "The record command chooses the version" of the TechSpec in
   the record mode of `TestEveryOwnedSkillVersionIsRecorded`: when a skill's
   content is not recorded under its declared version and that version is not
   higher than every recorded version, the version becomes one patch above the
   highest recorded version, as API Contract 1 states.
2. MUST rewrite every front-matter line matching `^ *version: ` in the
   canonical `.agents/skills/<name>/SKILL.md` and in its mirror
   `skills/<name>/SKILL.md`, compute the digest from the mirror on disk after
   the rewrite, and record that version and digest.
3. MUST refuse, writing nothing, when a skill's mirror differs from its
   canonical copy outside the version lines, with a message naming
   `make skills-sync`.
4. MUST NOT replace or remove a recorded entry, break ascending order, or
   change the check without the flag: unrecorded content still fails with
   today's "content changed under version" and "is not recorded" messages.
5. MUST add the new file `skills/owned_skill_version_raise_test.go` with
   `TestRecordingRaisesAVersionRecordedWithOtherContent`,
   `TestRecordingRaisesAVersionBelowTheHighestRecorded` and
   `TestRecordingRewritesBothVersionFieldsOfASkillAndItsMirror`, the last over
   temporary canonical and mirror roots asserting both fields, the recorded
   digest of the rewritten mirror and the drifted-mirror refusal, as Testing
   Approach 1 describes; no test writes under the repository's real skill
   folders.
6. MUST change, in `TestRecordingNeverReplacesARecordedVersion`, only the
   record-mode cases `changed digest` and `lower unrecorded version`, which
   now expect the raise, and keep every other case and test of
   `skills/owned_skill_versions_test.go` passing unedited.

## Subtasks

- [ ] Let the record decide a raise for colliding or stale content.
- [ ] Rewrite both version fields of the skill and its mirror and record the rewritten digest.
- [ ] Add the raise tests and update the two record-mode cases.

## Acceptance Criteria

- [ ] A skill recorded at `0.1.26` with other content and declared `0.1.26`
      is recorded at `0.1.27`, with both version fields of the skill and its
      mirror rewritten; a skill declared below the highest recorded version is
      raised the same way.
- [ ] A drifted mirror refuses and writes nothing; no recorded digest is ever
      replaced.
- [ ] Without the flag, every existing refusal keeps its message.

## Context

- interface: `skills/owned_skill_versions_test.go`
- creates: `skills/owned_skill_version_raise_test.go`
- instruction: `skills/testdata/owned-skill-versions.json`
- instruction: `docs/adr/0233-a-skill-version-raise-is-regenerated-at-merge-and-a-review-only-correction-returns-to-review.md`
- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRecordingRaisesAVersionRecordedWithOtherContent|TestRecordingRaisesAVersionBelowTheHighestRecorded|TestRecordingRewritesBothVersionFieldsOfASkillAndItsMirror|TestRecordingNeverReplacesARecordedVersion|TestAChangedOwnedSkillUnderARecordedVersionIsRefused|TestAnUnrecordedOwnedSkillVersionIsRefused|TestEveryOwnedSkillVersionIsRecorded)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRecordingRaisesAVersionRecordedWithOtherContent TestRecordingRaisesAVersionBelowTheHighestRecorded TestRecordingRewritesBothVersionFieldsOfASkillAndItsMirror TestRecordingNeverReplacesARecordedVersion TestAChangedOwnedSkillUnderARecordedVersionIsRefused TestAnUnrecordedOwnedSkillVersionIsRefused TestEveryOwnedSkillVersionIsRecorded; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three raise tests do not exist, so the command fails.

## References

- `_prd.md` → Goals 1, 2 and 4; Core Feature 1; Success Metric 1; Success Metric 4
- `_techspec.md` → The record command chooses the version; API Contract 1; Testing Approach 1; Build Order 1
- ADR-0233; ADR-0189
