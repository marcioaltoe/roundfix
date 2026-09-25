---
status: approved
granted: 2026-09-09
action: finish the existing cleanup pull request with its required checks
consuming: 0129-spec-authoring-and-gate-recovery
paths:
  - internal/docscontract/corpus_test.go
  - .github/workflows/ci-verify.yml
---

# Authorized cleanup delivery adjustments

The maintainer authorized all adjustments necessary to finish the existing
branch, open/update its Pull Request, squash merge and sync main on 2026-09-09.
The maintainer then rejected starting new implementation Specs and Roundfix Runs.
This records only the two required closure adjustments, applied directly on the
existing branch; it does not start the broader implementation of this Spec.

- Derive the current loop-order test input from the shipped declaration while
  preserving all three negative cases and their error/diagnostic assertions.
- Fetch complete history in the existing CI checkout so its unchanged historical
  authorization test can read the retained main ancestor. Preserve permissions,
  credential handling, action versions, time budget and verification commands.

No check bypass, assertion waiver, new feature, dependency upgrade or new Run is
authorized or required by this bounded delivery record. The initial broader
implementation plan remains separately proposed in `_authorization.md`.

The failure evidence is PR #180 job 102448882825 and the local
`TestCheckLoopOrderDivergent` failure. Secondbrain's historical-authorization
capture supports retaining real Git evidence. Exa read the official
[checkout documentation](https://github.com/actions/checkout), which establishes
the one-commit default and full-history input; it does not validate our code.
