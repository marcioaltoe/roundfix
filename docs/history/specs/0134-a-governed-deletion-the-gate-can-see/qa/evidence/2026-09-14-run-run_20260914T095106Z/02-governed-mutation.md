# Governed-mutation evidence

## Historical characterization

A local clone checked out commit `c03a9f305c7552e2ee9403da921b1cdbccdd495f`. On Go 1.26.7,

`rtk env GOTOOLCHAIN=go1.26.7 GOCACHE=/tmp/roundfix-go-build-0134-history go test -count=1 -v ./internal/daemon -run '^TestGovernedMutationClassificationCharacterization$'`

exited 0. Its three passing subtests recorded addition as a governed mutation, removal as not yet a governed mutation, and rename from a governed path as not yet a governed mutation. The current-tree run of the same test also exited 0 and reported addition, removal, and governed-source rename as governed mutations. `git show b8869c7a -- internal/daemon/task_engine_test.go` shows that only the latter two outcomes moved.

## Rebuilt CLI: governed deletion refusal

Scratch repository: `/private/tmp/roundfix-qa-0134-settle.pDRpoY/repo`. The committed authorization granted `implement` but not `commit`; the only dirty path was a deleted tracked `Makefile`.

`rtk env HOME=/private/tmp/roundfix-qa-0134-settle.pDRpoY/home GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null GIT_OPTIONAL_LOCKS=0 <worktree>/bin/roundfix settle --spec qa-smoke --task task_01`

exited 2 and reported:

```text
Preflight failed

Reason:
  refuse Settle commit: authorization operation "commit" is not permitted by record "docs/specs/qa-smoke/_authorization.md"

No side effects:
  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.
```

A fresh reread kept HEAD at `25cc550e9c806171927ebc6d33934f4e7d364bab`, kept `git status --porcelain=v1` at ` D Makefile`, kept the Task SHA-256 at `c2d011a5359377ee17033f6949e2672d682daacd6dd245ab2602f3ce77118c0d`, and found no `should-not-run` file. This proves the refusal happened before Verification and mutation.

## Rebuilt CLI: ordinary deletion control

Scratch repository: `/private/tmp/roundfix-qa-0134-ordinary.ui4Ns2/repo`. It had no authorization record and one deleted tracked ordinary source path, `internal/ordinary.go`.

The same rebuilt binary ran `settle --spec qa-smoke --task task_01`, exited 0, printed `verify true — ok`, printed `commit internal/ordinary.go — deleted`, and settled the Task at commit `d180ed2`. A fresh reread found a clean worktree, `status: completed`, and commit paths limited to the Task file plus the deleted ordinary source.

## Current reader checks

The fresh Go 1.26.7 focused run across `internal/daemon`, `internal/cli`, and `internal/speccheck` exited 0 for:

- `TestGovernedRemovalRequiresOperation`
- `TestGovernedRenameClassifiesFromSource`
- `TestOrdinaryRemovalDoesNotRequireOperation`
- `TestSettleClassifiesGovernedRemoval`
- `TestSettleRefusesMissingCommitAuthority`
- `TestSettleCommitsOrdinaryWorkWithoutRecord`
- `TestGovernedMutationRefusesMissingImplementAuthority`
- `TestGovernedMutationDetectionUsesTheUnfilteredSnapshot`
- `TestGovernedChangeStillRefusesWithoutRecord`
- `TestAuditRefusesOutOfGrantUnderBothCitationForms`
- `TestAuditRefusesSelfApprovalAndRetroactiveGrants`

These readers preserved the existing `commit`, authorization-record path, `QA-AUTH-PATHS`, and missing-`implement` refusal conditions while adding deletion and governed-source rename.
