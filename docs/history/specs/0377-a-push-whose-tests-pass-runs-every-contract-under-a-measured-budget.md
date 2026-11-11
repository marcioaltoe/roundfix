---
schema: roundfix/archive-record/v1
spec: 0377-a-push-whose-tests-pass-runs-every-contract-under-a-measured-budget
title: A push whose tests pass runs every contract under a measured budget
status: archived
created: "2026-11-10"
archived: "2026-11-11"
disposition: partial
source: docs/specs/0377-a-push-whose-tests-pass-runs-every-contract-under-a-measured-budget
source_revision: 7169c94394afdaa744cebe050a7869bd4f5b5f0a
qa_task: task_02
qa_report: qa-report-2026-11-11-01.md
qa_verdict: partial
unproven:
  - the operator's reading of the first push run of `CI | Verify` on `main` after the merge, with `gh run view <id> --json jobs`, recorded session evidence.
adrs:
  - ADR-0387
  - ADR-0259
sources:
  - 2026-11-10-the-push-gate-on-main-fails-its-suite-time-budget-on-most-runs.md
  - 2026-11-10-ci-verify-push-suite-times.md
regeneration:
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/setup-context.json
promoted: []
type: fix
---

# A push whose tests pass runs every contract under a measured budget

Local QA confirms that the suite-time verdict is recorded for later judgment, that overrun, missing and unreadable verdicts still fail their proper entry points, and that the workflow judges the budget after the contract runs while retaining the measured 420 s default. The required offline archive plan exited 0 with the judge skipped and listed the measurement and adopted operator entry as candidates; both describe the historical premise, and this gate promotes neither into durable reference material. The first post-merge hosted push remains unproven under the PRD declaration and no Pull Request acceptance exists, yielding a qualifying declared partial with no product findings. Offline archive plan (`docs/specs/0377-a-push-whose-tests-pass-runs-every-contract-under-a-measured-budget/qa/evidence/qa-report-2026-11-11-01/archive-plan.txt` (not reachable from any branch or tag)).
