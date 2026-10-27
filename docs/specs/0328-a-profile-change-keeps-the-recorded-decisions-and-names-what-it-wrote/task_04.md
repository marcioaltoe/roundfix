---
task: task_04
spec: 0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote
status: pending
type: backend
complexity: high
---

# Task 04: A restore keeps the repository's Claude Skill Links

## Overview

In a repository whose `.claude/skills` keeps per-skill symbolic links to
`.agents/skills`, `roundfix baseline skills restore` plans a Claude Skill
Link for each selected skill that lacks a `.claude/skills/<name>` entry,
binds it into the Plan Digest, and writes it with the skill tree and the lock
inside the one recoverable transaction, which learns a symbolic-link
postimage. Repositories without such links get nothing new. This answers
friction (e) of the adopted Backlog Entry (Quaestor, 2026-10-26).

## Requirements

1. MUST implement API Contract 6 in `buildSkillsRestorePlan`: detect the
   link-keeping `.claude/skills` exactly as the contract defines (a real
   directory, `.claude` not a symbolic link, at least one entry that is a
   symbolic link whose target is `../../.agents/skills/<its own name>`),
   plan `create-link` with target `../../.agents/skills/<name>` for each
   selected skill without an entry, include each in the Plan Digest, and
   plan it also when the skill tree itself needs no change, so a restore of
   an installed skill that lacks its link previews exit `3`.
2. MUST marshal a `create-link` planned change as
   `{"action","path","skill","target"}` and print it in the text preview as
   `- create-link .claude/skills/<name> [<name>]`.
3. MUST carry each link through `buildRestoreTransactionDocument` as a
   postimage of kind `symlink` whose content is the target, with a missing
   preimage, and teach `internal/baseline/transaction.go` to validate,
   journal, stage, replace after the skill tree, verify with `Lstat` and
   `Readlink`, and roll back that kind, keeping every existing refusal of a
   symbolic-link destination or parent for regular-file postimages and every
   existing transaction behavior unchanged (Invariant 6).
4. MUST leave an existing `.claude/skills/<name>` entry of any kind
   untouched and unfollowed, and plan no link where `.claude/skills` is
   absent, is a symbolic link, or holds no Claude Skill Link (Invariant 5).
5. MUST keep the update's skills stage working through `baseline.RestoreSkills`
   unchanged in code; a skill it restores gains its link by the same plan.
6. MUST add the tests `_techspec.md` → Testing Approach names for task_04:
   `TestSkillsRestoreLinksARestoredSkillWhereClaudeLinksAreKept`,
   `TestSkillsRestoreAddsNoLinkWithoutClaudeLinks` (absent, symbolic-link
   and copies-only `.claude/skills` as subtests),
   `TestSkillsRestoreLeavesAnExistingClaudeEntry`,
   `TestSkillsRestoreRollsBackTheClaudeLinkWithTheTree` (a fault after the
   link is replaced leaves neither the link nor the new tree),
   `TestSkillsRestoreLinksAnInstalledSkillThatLacksItsLink`,
   `TestTransactionWritesVerifiesAndRollsBackASymbolicLinkPostimage` and
   `TestBaselineSkillsRestoreTextListsTheClaudeLink`, using the offline
   restore source the existing restore tests build. Each fails when its
   behavior is removed. Test files compile for Windows.

## Subtasks

- [ ] Characterize today's restore in a repository with a Claude Skill Link: no link is planned.
- [ ] Add the symbolic-link postimage to the transaction with its rollback.
- [ ] Plan the links in the restore and carry them into the transaction document.
- [ ] Print and marshal the new planned change.
- [ ] Add the named tests.

## Acceptance Criteria

- [ ] A restore of a new skill in a repository that keeps Claude Skill Links writes `.claude/skills/<name>` pointing at `../../.agents/skills/<name>`.
- [ ] A restore of an installed skill that lacks its link previews and writes only the link.
- [ ] No link is written where the repository keeps none, and an existing entry is never changed.
- [ ] A failed apply leaves no new link.

## Context

- instruction: `docs/user-guide/commands/baseline.md`
- instruction: `.agents/skills/roundfix/references/baseline.md`
- instruction: `docs/adr/0341-a-profile-change-carries-the-recorded-decisions-and-a-restore-keeps-the-claude-skill-links.md`
- instruction: `internal/baseline/skills_restore_test.go`
- instruction: `internal/cli/cli.go`
- instruction: `internal/cli/baseline_update.go`
- interface: `internal/baseline/skills_restore.go`
- interface: `internal/baseline/transaction.go`
- interface: `internal/cli/baseline_skills_restore.go`
- creates: `internal/baseline/skills_restore_claude_link_test.go`
- interface: `internal/baseline/transaction_test.go`
- interface: `internal/cli/baseline_skills_restore_test.go`

## Consumers

- The file transaction and the restore payload — consumers: `./internal/baseline`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestSkillsRestore.*|TestTransaction.*|TestHistoryMove.*|TestBaselineSkillsRestore.*|TestBaselineUpdatePreviewNamesAMissingRequiredSkill)$' ./internal/baseline ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSkillsRestoreLinksARestoredSkillWhereClaudeLinksAreKept TestSkillsRestoreAddsNoLinkWithoutClaudeLinks TestSkillsRestoreLeavesAnExistingClaudeEntry TestSkillsRestoreRollsBackTheClaudeLinkWithTheTree TestSkillsRestoreLinksAnInstalledSkillThatLacksItsLink TestTransactionWritesVerifiesAndRollsBackASymbolicLinkPostimage TestBaselineSkillsRestoreTextListsTheClaudeLink TestSkillsRestoreOfflinePreviewApplyAndIdempotence TestTransactionRollback; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go build -buildvcs=false ./internal/baseline ./internal/cli` — expected: exit 0. Before this Task the seven new tests do not exist, so the first pass check fails. Afterwards they pass with the existing restore, transaction and history-move tests, and both packages build for Windows.

## References

- `_prd.md` → Goals; Core Feature 6; User Experience; Success Metric 5; Acceptance evidence
- `_techspec.md` → Current behavior 6; Interfaces; Data Models; API Contract 6; Invariants 5-6; Testing Approach; Build Order 3; Risks & Considerations
- ADR-0341
