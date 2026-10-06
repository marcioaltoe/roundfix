---
task: task_01
spec: 0237-review-selection-failures-that-say-why
status: completed
type: docs
complexity: low
---

# Task 01: The review reference and guide describe the failing step, the adapter's message and the retry

## Overview

Describes, in the Roundfix Skill's `review` reference, its mirror and the
`review` command guide, what a blocked review now says about a runtime
failure and when it retries before the prompt, so the skill ships with the
CLI behavior task_03 changes. It answers the Backlog Entry "A pre-PR review
agent selection failure says nothing actionable" of 2026-10-06. Verifiable on
its own by the phrases, the mirrors and the recorded skill version.

## Requirements

1. MUST describe, in `.agents/skills/roundfix/references/review.md`, API
   Contract 3: the record fields `failedStep` and `adapterMessage`, that the
   step is the protocol request that failed (`initialize`, `session/load` or
   `session/new`, `session/set_model`, `session/prompt`) or the session
   preparation command (`sessions ensure`, `set <key>`, `set-mode`), and that
   the message is only the adapter's error message, on one line of at most 512
   bytes.
2. MUST describe API Contract 4: a failure during Agent Selection before the
   review prompt was sent is retried once on the same selection; the stderr
   notice containing the exact phrase "retrying selection"; the record field
   `selectionRetries` with `selection`, `step` and `message`; and the reason
   suffix containing the exact phrase
   "after one automatic retry before the prompt".
3. MUST update the paragraph that says when a configured selection fallback is
   eligible so that it states the fallback activates only after the retry also
   failed before the prompt, that a failure after the prompt is never retried,
   and that no failure selects `none`.
4. MUST make `docs/user-guide/commands/review.md` carry the same text as the
   reference below its first heading line, as it does today.
5. MUST raise the Roundfix Skill's version in both front-matter fields of
   `.agents/skills/roundfix/SKILL.md` from the value it holds when this Task
   starts, run `make skills-sync`, and record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   Running those commands is an implementation step, never part of
   Verification, and no version is written by hand.
6. MUST NOT change any other command reference, the skill's other text, or
   the `### QA settlement` section of any skill.

## Subtasks

- [ ] Describe the failing step, the adapter's message and their record fields.
- [ ] Describe the single retry before the prompt, its notice and record field.
- [ ] Update the fallback-eligibility paragraph.
- [ ] Mirror the guide, raise and record the skill version, sync the mirror.

## Acceptance Criteria

- [ ] The reference and the guide carry the two exact phrases and the three
      record field names, and differ only in their first heading line.
- [ ] `skills/roundfix/references/review.md` and `skills/roundfix/SKILL.md`
      equal their canonical files, and the raised version is recorded.

## Context

- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `docs/user-guide/commands/review.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `docs/adr/0242-a-blocked-review-names-the-protocol-step-and-retries-once-before-the-prompt.md`

## Verification

- `body="$(mktemp)" || exit 1; for phrase in "retrying selection" "after one automatic retry before the prompt" "failedStep" "adapterMessage" "selectionRetries"; do tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/review.md | grep -qF -- "$phrase" || { printf 'missing phrase in the reference: %s\n' "$phrase" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/review.md | grep -qF -- "$phrase" || { printf 'missing phrase in the guide: %s\n' "$phrase" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md || exit 1; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md || exit 1; tail -n +2 .agents/skills/roundfix/references/review.md > "$body" || exit 1; tail -n +2 docs/user-guide/commands/review.md | cmp - "$body"` — expected: exit 0; before this Task neither document carries the phrases, so the command fails.
- `tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/review.md | grep -qF -- "selectionRetries" || { printf 'review reference not updated\n' >&2; exit 1; }; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills` — expected: exit 0; before this Task the reference lacks the field name, so the command fails, and after it the raised version must be recorded for the test to pass.

## References

- `_prd.md` → Core Feature 4; User Stories 1-3; Goals
- `_techspec.md` → Build Order 1; API Contract 3; API Contract 4; Vocabulary Contract
- ADR-0242
- ADR-0189

## Result

Implementation is ready for Daemon Verification. The review reference now
documents `failedStep`, `adapterMessage`, the bounded protocol or session
preparation step, and `selectionRetries` with the single pre-prompt retry,
stderr notice, retry reason suffix, and fallback eligibility rules. The command
guide carries the same body below its heading. The Roundfix Skill version was
raised from `0.1.40` to `0.1.41`, its shipped mirrors were regenerated, and the
owned version record was generated.

Focused evidence:

- `make skills-sync` — passed.
- `GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-deliver-0237-review-selection-failures-that-say-why-bb5befdc69932e64-0060f39f-5068b8b4/run_20261006T151517Z_63f7ac2ee342847c.task_01/.gocache go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — passed, 1 test; recorded Roundfix Skill version `0.1.41`.
- `git diff --check` — passed.
- Phrase and field search over the reference and guide found `retrying selection`, `after one automatic retry before the prompt`, `failedStep`, `adapterMessage`, and `selectionRetries` in both documents.
- The guide and reference bodies matched after their first heading lines; canonical and shipped Roundfix skill files were regenerated together by `make skills-sync`.

The first version-recording attempt against the shared Go cache was blocked by
cache permissions; the same required command passed with the task-scoped
absolute `GOCACHE` above. The declared Verification commands remain for the
Daemon and were not run as the terminal gate.
