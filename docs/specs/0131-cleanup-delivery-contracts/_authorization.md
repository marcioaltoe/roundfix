---
status: approved
granted: 2026-09-09
action: finish the existing cleanup branch and its required delivery checks
consuming: 0131-cleanup-delivery-contracts
paths:
  - internal/docscontract/corpus_test.go
  - .github/workflows/ci-verify.yml
---

# Approved final cleanup delivery repairs

The maintainer authorized on 2026-09-09:

> Considere autorizada as mudanças e ajustes necessários para finalizar o trabalho nessa branch e fazemos o pull request, squash merge e sync main.

This applies the approval to the two observed remaining failures: derive the
loop-order mutation fixture from the current canonical declaration, preserving
its strict negative assertions; and make the existing Verification CI checkout
include the history required by the historical authorization tests.

The Task owns exactly the two implementation files above plus its Task evidence.
Record this grant before the consuming commit on the existing cleanup branch.
The maintainer requested delivery and squash merge on that branch; do not create
another long-lived repair or authorization branch. A Roundfix Run branch is
transient and must be reclaimed after integration. This record adds to the
already merged grants from PRs #179 and #181; it does not rewrite them.

The same maintainer instruction covers other repairs proven necessary to finish
this branch: record their exact scope before consumption if such a failure is
observed. No weakening of checks, new feature, dependency upgrade or release is
part of this bounded Task.
