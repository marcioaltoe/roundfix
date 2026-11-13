---
task: task_01
spec: 0388-a-completion-notice-skips-a-record-the-merge-adopted
status: completed
type: docs
complexity: low
---

# Task 01: A completion notice skips a record the merge adopted, and every guarantee of the section is pinned

## Overview

The closure rule of `## Completion notices` in the roundfix skill's release
reference sends a notice for every record its deletion listing prints, and that
listing prints every Backlog Entry or Finding the merge adopted. This Task
rewrites the closure bullet, adds the paragraph that recognizes an adopted
record with Git alone, restores `memory_message_send` in the send paragraph,
syncs the shipped mirror, raises and records the skill's version, and pins
every guarantee of the section with its own phrase plus an executable check of
the section's Git commands on a temporary repository.

## Requirements

1. MUST answer the Backlog Entry of 2026-11-13,
   [`references/2026-11-13-a-completion-notice-reads-an-adopted-record-as-closed.md`](references/2026-11-13-a-completion-notice-reads-an-adopted-record-as-closed.md)
   (`_prd.md` → Sources 1; reported by argus, ai-memory message `01a1d483`),
   by implementing API Contracts 1 to 4 of `_techspec.md`.
2. MUST, in `.agents/skills/roundfix/references/release.md`, replace the
   closure bullet "- **A record is closed without a Spec.** …" (lines 145 to
   153 at the authoring base) with the bullet and the paragraph of
   `_techspec.md` → Interfaces, byte for byte, one blank line between them and
   one before the paragraph that begins "A Notice Request's work is delivered";
   and replace the paragraph that begins "Send each notice once" with the send
   paragraph of Interfaces, byte for byte. Every other byte of the file stays.
3. MUST run `make skills-sync`, so `skills/roundfix/` equals the canonical
   copy, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   once, after the last edit of the skill, which raises `metadata.version` and
   `version` in both `SKILL.md` copies to the next version and appends it to
   `skills/testdata/owned-skill-versions.json`. Never type a version or a
   digest by hand.
4. MUST, in `skills/roundfix_modes_test.go`:
   - add to the phrase list of
     `TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest`, after its
     20 phrases and in the order of `_techspec.md` → Guarantee inventory, the
     16 phrases the inventory marks "E:", each copied from the inventory;
   - add `TestACompletionNoticeSkipsARecordTheMergeAdopted` and the helpers
     `completionNoticeGitCommands`, `completionNoticeRun`,
     `completionNoticeClassify` and `completionNoticeFixture` of
     `_techspec.md` → Interfaces, behaving as API Contract 3 and Testing
     Approach: for each layout of `roundfixLayoutSources()` a subtest that
     checks the 13 adoption phrases of Testing Approach, then subtests
     `squash` and `merge` that build the fixture with
     `roundfix/internal/gittest` (`InitRepo`, `Run`, `ConfigArgs`,
     `IsolatedEnv`), run the section's own deletion listing, addition listing
     and `git grep` command with the placeholders replaced, under
     `diff.renames` `false`, `true` and `copies`, and compare the three
     classes with the expected records, reporting
     `diff.renames=<v>: <class> records are <got>, want <want>`;
   - each top-level test and subtest calls `t.Parallel()` first; the
     temporary repositories live under `t.TempDir()`; no helper name collides
     with a name the package already declares.
5. MUST NOT change any Go production file, any file under `internal/`,
   `docs/agents/`, `docs/user-guide/` or `docs/adr/`, `CONTEXT.md`, any other
   skill, any roundfix reference other than `references/release.md`, any
   section of it other than `## Completion notices`, the Makefile, the CI
   workflows or the `### QA settlement` section of any skill.

## Subtasks

- [ ] Rewrite the closure bullet, add the adoption paragraph and restore the send channel in the canonical reference.
- [ ] Sync the mirror and record the raised version.
- [ ] Add the guarantee pins, the adoption test and its fixture.

## Acceptance Criteria

- [ ] Both copies of the release reference set aside a record the merge adopted, say how to recognize it with Git alone, ask the user when the evidence disagrees, and name `memory_message_send`.
- [ ] Under every rename setting and both merge shapes, the section's own commands classify the adopted records as adopted, the declined one as closed and the unconfirmed one as one to ask about.
- [ ] Every guarantee of the inventory has its own pinned phrase, and `roundfix` carries its next version, recorded, with the mirror equal to the canonical copy.

## Context

- instruction: `skills/roundfix_layout_test.go`
- instruction: `skills/owned_skill_versions_test.go`
- instruction: `internal/gittest/gittest.go`
- instruction: `docs/agents/docs-layout.md`
- interface: `.agents/skills/roundfix/references/release.md`
- interface: `skills/roundfix/references/release.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `skills/roundfix_modes_test.go`

## Consumers

- the shipped `roundfix` skill folder and the owned-skill version record — consumers: `./skills`

## Verification

- `out="$(go test -timeout 10m -count=1 -v ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestACompletionNoticeSkipsARecordTheMergeAdopted/canonical TestACompletionNoticeSkipsARecordTheMergeAdopted/canonical/squash TestACompletionNoticeSkipsARecordTheMergeAdopted/canonical/merge TestACompletionNoticeSkipsARecordTheMergeAdopted/embedded TestACompletionNoticeSkipsARecordTheMergeAdopted/embedded/squash TestACompletionNoticeSkipsARecordTheMergeAdopted/embedded/merge TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest/canonical TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest/embedded TestEveryOwnedSkillVersionIsRecorded TestAuthorialSkillSync/roundfix; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the adoption test does not exist, so the command fails at its first pass check. Afterwards both layouts pass the adoption phrases and the executable check for both merge shapes, the extended guarantee pins pass, the version record and the `roundfix` mirror pass, and the whole `./skills` package passes.
- `for f in .agents/skills/roundfix/references/release.md skills/roundfix/references/release.md; do t="$(tr -s '[:space:]' ' ' < "$f")"; printf '%s' "$t" | grep -qF -- 'A record the merge adopted is not closed: its notice follows the Spec that adopted it, and the closure event sends nothing for it.' || { printf 'missing adoption rule in %s\n' "$f" >&2; exit 1; }; printf '%s' "$t" | grep -qE -- 'Send each notice once with .memory_message_send., from the session that follows the event' || { printf 'missing send channel in %s\n' "$f" >&2; exit 1; }; done; cmp -s .agents/skills/roundfix/references/release.md skills/roundfix/references/release.md` — expected: exit 0. Before this Task neither copy holds the adoption rule, so the first check fails. Afterwards both copies hold the rule and the send channel, and they are identical.

## References

- `_prd.md` → Sources; Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → Current behavior; Interfaces; API Contracts 1-4; Guarantee inventory; Invariants 1-4; Sabotages; Testing Approach; Build Order 1

## Result

Skills loaded: implement-task, coding-guidelines, golang-testing, testing-boss, tech-writer, evidence-gate, systematic-debugging

Implementation evidence:

- Rewrote the canonical `## Completion notices` closure rule with the adopted-record exception, Git-only adoption checks, the user-question case, and `memory_message_send`; ran `make skills-sync` and confirmed the canonical and shipped references are byte-identical.
- Added the 16 guarantee pins and the executable adoption fixture/test. The fixture uses `roundfix/internal/gittest` and covers canonical/embedded layouts, squash/merge shapes, and `diff.renames=false`, `true`, and `copies`.
- Ran the owned-skill version recorder once after the last skill edit; it raised both `roundfix/SKILL.md` copies and appended the recorded version/digest entry.

Focused checks:

- `go test ./skills -run '^TestACompletionNoticeSkipsARecordTheMergeAdopted$' -count=1 -v`: passed for both layouts, both merge shapes, and all three rename settings.
- `go test ./skills -run '^TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest$' -count=1 -v`: passed for canonical and embedded layouts.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: passed while recording the raised skill version.
- `git diff --check` and `cmp -s .agents/skills/roundfix/references/release.md skills/roundfix/references/release.md`: passed.

Daemon handoff:

- The Task status remains `in_progress`; the Daemon owns settlement.
- The commands in `## Verification`, the repository-wide verification, and any commit/push/PR action were not run.
