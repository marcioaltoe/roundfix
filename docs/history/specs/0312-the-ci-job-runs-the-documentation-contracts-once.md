---
schema: roundfix/archive-record/v1
spec: 0312-the-ci-job-runs-the-documentation-contracts-once
title: The CI job runs the documentation contracts once
status: archived
created: "2026-10-24"
archived: "2026-10-24"
disposition: partial
source: docs/specs/0312-the-ci-job-runs-the-documentation-contracts-once
source_revision: 6e26879302527af365dc1ae796a19ae49f50a3d9
qa_task: task_02
qa_report: qa-report-2026-10-24.md
qa_verdict: partial
unproven:
  - the operator's reading of that pull request's `ci-verify` run after this Spec merges.
adrs: []
sources:
  - 2026-10-22-the-ci-verify-job-runs-the-documentation-contract-package-twice.md
regeneration: []
promoted: []
---

# The CI job runs the documentation contracts once

The pull-request Verify docs step skips the documentation-contract package already run by Verify changed, while standalone verify-docs and pushes retain the package run and the always-selected guard catches marker regressions. Seven rows passed; the GitHub expression source and Pull Request controls retain permitted environment blocks, and Success Metric 3 remains unproven until the operator reads the first post-merge CI timing. With a temporary HOME and all Judge Key Variables plus TYPESAFE_API_KEY unset, `./bin/roundfix archive 0312-the-ci-job-runs-the-documentation-contracts-once --plan --offline` exited 0 and reported `judge: skipped (offline)`; its transcript (`docs/specs/0312-the-ci-job-runs-the-documentation-contracts-once/qa/evidence/2026-10-24-offline/archive-plan.txt` (not reachable from any branch or tag)) records the proposed Archive Record and retained reference candidate. No archive mutation, promotion, commit, push or Pull Request occurred.
