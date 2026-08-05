# QA-05 — Open Pull Request survivor

Status: blocked (environment: no open Pull Request).

The prompt proves that no Pull Request is open. The Run Worktree branch is
never used to infer one.

Equivalent observed evidence on the same build:

- `TestAuditClassifiesPullRequestBranch` passed inside the fresh 12-test
  `internal/specaudit` selection. It records Pull Request #42, requires
  `pull-request` with evidence naming it, and requires no reclaim command.
- `TestAuditReplaysMotivatingSessionResidue` also passed and independently
  covers Pull Request #58 and #68 branches, neither reclaimable.
- Full `rtk make verify` passed all 3,313 Go tests.

Unblocking action: open a Pull Request on `ma/spec-0068-implementation`, then
rerun the gate read-only against its current head.
