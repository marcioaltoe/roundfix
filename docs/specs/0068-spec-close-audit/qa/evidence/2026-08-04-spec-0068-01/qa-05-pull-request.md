# QA-05 — Open Pull Request survivor

Status: blocked (environment: no open Pull Request).

The QA prompt proves that no Pull Request is open and states that Pull Request
journeys are environment-blocked. The per-Run branch is not used to infer a
Pull Request.

Fresh equivalent evidence on build `30ec663c`:

- `TestAuditClassifiesPullRequestBranch` passed in the named 17-test audit
  selection. Its real-Git fixture records Pull Request #42, requires kind
  `pull-request`, requires evidence naming that Pull Request, and requires no
  reclaim command.
- `TestAuditReplaysMotivatingSessionResidue` passed independently for Pull
  Request #58 and #68 branches, both without reclaim commands.
- Full `rtk make verify` passed all 3,318 Go tests.

Unblocking action: open a Pull Request on `ma/spec-0068-implementation`, then
rerun this read-only row against the Open Pull Request. The equivalent fixture
evidence satisfies ADR-0080 for this environment block.
