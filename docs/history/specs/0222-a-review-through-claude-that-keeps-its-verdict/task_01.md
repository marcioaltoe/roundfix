---
task: task_01
spec: 0222-a-review-through-claude-that-keeps-its-verdict
status: completed
type: docs
complexity: low
---

# Task 01: The review reference and guide describe a review through Claude that keeps its verdict

## Overview

The Roundfix Skill's `review` reference and the `review` command guide state
what `roundfix review` does after this Spec: a read-only turn that ended after
the session refused a permission is classified by its verdict, a runtime's
`Prompt is too long` answer is named in the reason, and a Claude prompt is
bounded in estimated tokens. The repository's skill-sync rule requires the
skill to ship with the CLI change, and the owned skill's version rises with
its content. This Task is verifiable on its own through the documents and the
skill-version record.

## Requirements

1. MUST describe, in `.agents/skills/roundfix/references/review.md`, API
   Contract 2: the record field `permissionRefused`, classification by the
   verdict after a read-only session refused a permission request, finding ids
   and disposition as for any findings answer, and the blocked-reason suffix,
   containing the exact phrase
   "after the read-only session refused a permission request".
2. MUST describe API Contract 3 with the exact reason prefix
   "review prompt too long:" and say that the line comes from the final answer.
3. MUST describe API Contract 4: for `claude`, the estimate of the whole
   prompt at two bytes per token, the record field `estimatedPromptTokens`,
   the refusal before readiness and any provider call, and the reason
   containing the exact phrase
   "half of its 1000000-token context window"; and that the 917,504-byte diff
   bound now applies to `codex` only, keeping its existing reason text.
4. MUST update the paragraph that lists what blocks a review so that a
   transport anomaly still blocks and the refused-permission exit of a
   read-only `end_turn` turn is named as the exception.
5. MUST make `docs/user-guide/commands/review.md` carry the same text as the
   reference below its first heading line, as it does today.
6. MUST raise the Roundfix Skill's version in both front-matter fields of
   `.agents/skills/roundfix/SKILL.md` from the value it holds when this Task
   starts, run `make skills-sync`, and record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
7. MUST NOT change any other command reference, the skill's other text, or
   the `### QA settlement` section of any skill.

## Subtasks

- [ ] Describe the refused-permission classification and its record field.
- [ ] Describe the prompt-too-long reason.
- [ ] Describe the Claude token bound and limit the byte bound to Codex.
- [ ] Mirror the guide, raise and record the skill version, sync the mirror.

## Acceptance Criteria

- [ ] The reference and the guide carry the three exact phrases and both
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
- instruction: `docs/adr/0227-a-review-through-claude-keeps-its-verdict-and-fits-its-window.md`

## Verification

- `body="$(mktemp)" || exit 1; for phrase in "after the read-only session refused a permission request" "review prompt too long:" "half of its 1000000-token context window" "permissionRefused" "estimatedPromptTokens"; do tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/review.md | grep -qF -- "$phrase" || { printf 'missing phrase in the reference: %s\n' "$phrase" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/review.md | grep -qF -- "$phrase" || { printf 'missing phrase in the guide: %s\n' "$phrase" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md || exit 1; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md || exit 1; tail -n +2 .agents/skills/roundfix/references/review.md > "$body" || exit 1; tail -n +2 docs/user-guide/commands/review.md | cmp - "$body"` — expected: exit 0; before this Task neither document carries the phrases, so the command fails.
- `tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/review.md | grep -qF -- "review prompt too long:" || { printf 'review reference not updated\n' >&2; exit 1; }; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills` — expected: exit 0; before this Task the reference lacks the phrase, so the command fails, and after it the raised version must be recorded for the test to pass.

## References

- `_prd.md` → Core Feature 5; User Stories 1-4; Goals
- `_techspec.md` → Build Order 1; API Contract 2; API Contract 3; API Contract 4
- ADR-0227

## Result

Implemented the review documentation slice. The canonical review reference now
describes refused-permission verdict classification and `permissionRefused`,
the final-answer `review prompt too long:` reason, and Claude's whole-prompt
two-bytes-per-token estimate with `estimatedPromptTokens`, pre-readiness
refusal, and the half-window reason. The 917504-byte bound is explicitly
limited to Codex, and the transport-anomaly blocking paragraph names the
read-only `end_turn` refusal exception. The user guide mirrors the reference
below its heading. The owned Roundfix Skill is version `0.1.24` in both front
matter fields, its mirror is synchronized, and the recorded digest is present.

Focused checks:

- `make skills-sync` — passed.
- `GOCACHE=/private/tmp/roundfix-task-01-gocache go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — passed.
- Whitespace-normalized phrase checks for both review documents — passed for
  all three required phrases and both record field names.
- Canonical/mirror parity checks for the review reference and Roundfix Skill,
  guide/reference body parity, and `git diff --check` — passed.

Acceptance evidence:

- The reference and guide contain the required phrases and fields, and their
  content matches below the first heading line.
- The synchronized skill files are equal, both front matter versions are
  `0.1.24`, and the owned-skill version record contains the generated digest.

## Carry-forward provenance

- Source Run: `run_20261004T145620Z_ee61ac9af5ab0ce2`
- Source commit: `6d068a7943d9a1017c4414b9ea0f8e8550f04f1c`
