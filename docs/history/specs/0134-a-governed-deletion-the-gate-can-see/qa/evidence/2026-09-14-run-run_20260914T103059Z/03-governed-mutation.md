# Governed-mutation evidence

## Historical and assembled characterization

A fresh local clone checked out
`c03a9f305c7552e2ee9403da921b1cdbccdd495f`. On Go 1.26.7,
`TestGovernedMutationClassificationCharacterization` exited 0 and passed:

- addition is a governed mutation;
- removal is not yet a governed mutation;
- rename to an ungoverned path is not yet a governed mutation.

The current-tree run exited 0 and passed addition, removal, and governed-source
rename as governed mutations. `git diff e9746a45..HEAD --
internal/daemon/task_engine_test.go` shows that the addition assertion stayed
unchanged; Task 02 moved removal and rename and added their Task-cycle controls.

## Current reader checks

The fresh Go 1.26.7 focused run across `internal/daemon`, `internal/cli`,
`internal/speccheck`, `internal/suiteguardcontract`, and `internal/baseline`
exited 0. Its governed-mutation cases passed:

- `TestGovernedRemovalRequiresOperation`;
- `TestGovernedRenameClassifiesFromSource`;
- `TestOrdinaryRemovalDoesNotRequireOperation`;
- `TestSettleClassifiesGovernedRemoval`;
- `TestSettleRefusesMissingCommitAuthority`;
- `TestSettleCommitsOrdinaryWorkWithoutRecord`;
- `TestGovernedMutationRefusesMissingImplementAuthority`;
- `TestGovernedMutationDetectionUsesTheUnfilteredSnapshot`;
- `TestGovernedChangeStillRefusesWithoutRecord`;
- `TestAuditRefusesOutOfGrantUnderBothCitationForms`;
- `TestAuditRefusesSelfApprovalAndRetroactiveGrants`.

## Rebuilt CLI: governed deletion refusal

Scratch repository:
`/private/tmp/roundfix-qa-0134-settle.xAryk8/repo`. Its committed authorization
granted `implement` but omitted `commit`; the only dirty path was a deleted
tracked `Makefile`.

The rebuilt binary ran `settle --spec qa-smoke --task task_01`, exited 2, and
reported:

```text
Preflight failed

Reason:
  refuse Settle commit: authorization operation "commit" is not permitted by record "docs/specs/qa-smoke/_authorization.md"

No side effects:
  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.
```

A fresh reread kept HEAD at
`4e585f014fd423e48af1e74dbe9c7c43babc3a71`, kept ` D Makefile`, and kept the
Task SHA-256 at
`65a8aae642d27c7e8b872346c8d331029935a478d93b85988a0ba9c8266be6ef`.
`should-not-run` remained absent. The isolated home contained no Run Database.

## Rebuilt CLI: ordinary deletion control

Scratch repository:
`/private/tmp/roundfix-qa-0134-ordinary.vxiwYq/repo`. It had no authorization
record and one deleted tracked ordinary source path, `internal/ordinary.go`.

The rebuilt binary ran `settle --spec qa-smoke --task task_01`, exited 0,
printed `verify test ! -e internal/ordinary.go — ok`, and settled at
`0a208c3d718bf2be6dfb8fc864a8cf5ac414d03b`. A fresh reread found a clean
worktree, `status: completed`, the source still absent, and commit paths limited
to the Task file plus the deleted ordinary source.
