---
schema: roundfix/archive-record/v1
spec: 0213-a-test-suite-that-does-not-flake
title: A test suite that does not flake
status: archived
created: "2026-10-01"
archived: "2026-10-02"
disposition: qa-override
source: docs/history/specs/0213-a-test-suite-that-does-not-flake
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-02.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: QA partial with zero findings; the only blocked rows (10, 13) need an open Pull Request, and row 10 proves the ETXTBSY fixture change on Linux, which the PR's CI Verification gate runs
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 4bd6116ffa23b780f7faa5b93aaa81f4d24f5e62
adrs:
  - ADR-0125
  - ADR-0213
sources:
  - 2026-09-30-time-bound-tests-and-a-leaked-detached-child.md
regeneration: []
promoted: []
pull_request: "341"
delivery_commit: 6364d3c9192a03ab8e6c505c71fea4baf6fa8810
---

# A test suite that does not flake

Six tests fail now and then in `make verify` or in the CI Verification gate, then pass on rerun. Each failure parks a delivery that was otherwise ready, and the operator retries it by hand (interventions 59 and 60 on 2026-10-01). The adopted Backlog Entry ([references/2026-09-30-time-bound-tests-and-a-leaked-detached-child.md](references/2026-09-30-time-bound-tests-and-a-leaked-detached-child.md)) names three of the tests, and its 2026-10-01 addendum names three more. The authoring session reproduced five of the six, and reading the code gives the cause of each:
