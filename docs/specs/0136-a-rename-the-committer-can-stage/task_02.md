---
task: task_02
spec: 0136-a-rename-the-committer-can-stage
status: completed
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

## Result

The final-push caller now reads the paths changed since the Run's initial HEAD
and classifies those paths with the existing Governed Path source of truth. It
passes that classification to the unchanged explicit `FinalPush` operation;
authorization-record outcome no longer stands in for mutation state.

Focused-check evidence:

- Before the production change,
  `rtk proxy env GOCACHE=/private/tmp/roundfix-go-build-0136-task-02 rtk proxy go test -v ./internal/cli -run 'FinalPushAuthorityFollowsTheChangedPaths/ordinary'`
  reached the full Implement flow and failed after Task settlement with the
  existing missing-`push` refusal for the ordinary path
  `internal/cli/implement.go`.
- After the change,
  `rtk proxy env GOCACHE=/private/tmp/roundfix-go-build-0136-task-02 rtk go test ./internal/cli -run 'FinalPushAuthorityFollowsTheChangedPaths/(ordinary|governed)'`
  passed all three reported tests: the parent and both named subtests.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-go-build-0136-task-02 rtk go test ./internal/cli -run 'TestRunImplementAutoPush(OutcomeMatrix|MissingUpstreamWarnsAndStaysClean|FailureEndsFailedAndJournalsPush)$'`
  passed all eight reported tests, covering the adjacent Clean-only, disabled,
  missing-upstream, unresolved, stopped, QA-failed, and pusher-failure paths.
- The first unprivileged focused attempt did not reach the test because the
  sandbox denied the default Go build cache; the unchanged retry used the
  writable task cache above. `rtk git diff --check` then passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-go-build-0136-task-02 rtk make verify-incremental`
  stopped in `fmt-check` before the test stage because the unchanged files
  `internal/cli/baseline_skills_restore_test.go` and
  `internal/cli/baseline_assets_sync_test.go` need formatting. Neither file is
  in this Task's diff, so this slice leaves them unchanged.

Acceptance evidence:

- The ordinary-Run subtest uses a record granting `implement` and `commit` but
  not `push`; it observes exit zero, one push call, and the pushed target
  `origin/ma/widget-flow`.
- The governed-Run subtest uses the same record with changed path `Makefile`;
  it observes no push call and the unchanged diagnostic
  `authorization operation "push" is not permitted`.
- The caller still invokes `engine.FinalPush` after a Clean Task cycle. The
  engine operation and its separate-operation test were not edited, and the
  focused auto-push matrix still observes exactly one final push only on the
  eligible Clean path.

## Carry-forward provenance

- Source Run: `run_20260914T134813Z_e15c7ee834fbe8bf`
- Source commit: `f9feb1b3efa752df674b444598a08bb123735835`
