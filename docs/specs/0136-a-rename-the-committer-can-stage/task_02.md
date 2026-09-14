---
task: task_02
spec: 0136-a-rename-the-committer-can-stage
status: pending
type: backend
complexity: low
---

# Task 02: Derive the final push's authority from the changed paths

## Overview

The Run's final push passes governed mutation as "an authorization record was
granted", which describes the record rather than the work. A Spec whose record
grants implement and commit but not push, and whose Run changed only ordinary
files, is then refused for authority its change never needed. The Run already
knows which paths it changed; the decision must come from those.

## Requirements

1. MUST derive whether the final push requires push authority from the paths
   the Run actually changed, not from whether an authorization record resolved
   to granted.
2. MUST keep requiring push authority when the Run did change a Governed Path.
3. MUST keep the final push a separate explicit operation, unchanged in every
   other respect.
4. MUST NOT change what any operation permits, or any refusal token.

## Subtasks

- [ ] Derive the flag from the Run's changed paths.
- [ ] Prove an ordinary Run is no longer refused.
- [ ] Prove a governed Run is still refused.

## Acceptance Criteria

- [ ] An ordinary Run under a record that grants implement and commit but not
      push completes its final push; today it is refused.
- [ ] A Run that changed a Governed Path under the same record is still refused
      at the final push.
- [ ] The final push remains a separate explicit operation.

## Context

- interface: `internal/cli/implement.go`

## Verification

- `grep -q 'func TestFinalPushAuthorityFollowsTheChangedPaths' internal/cli/implement_test.go && go test -count=1 ./internal/cli -run '^TestFinalPushAuthorityFollowsTheChangedPaths$'` — the flag follows the change, for an ordinary and a governed Run; this fails today.
- `grep -q 'func TestFinalPushAuthorityFollowsTheChangedPaths' internal/cli/implement_test.go || exit 1; go test -count=1 ./internal/daemon -run '^TestFinalPushIsASeparateExplicitOperation$'` — the final push is still a separate explicit operation.

## References

- `_prd.md` → Goal 2; Core Feature 4; Declared intentional breaks 2.
- `_techspec.md` → Implementation Design: Ask the change, not the record;
  Testing Approach observation 5; Build Order 2.
