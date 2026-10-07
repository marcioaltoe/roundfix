---
schema: roundfix/archive-record/v1
spec: 0171-a-deterministic-suite-under-load
title: A deterministic suite under load
status: archived
created: "2026-09-25"
archived: "2026-09-25"
disposition: pass
source: docs/history/specs/0171-a-deterministic-suite-under-load
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_06
qa_report: qa-report-2026-09-25-03.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0089
sources:
  - 2026-09-25-implement-tests-fail-the-qa-gate-under-load.md
  - 2026-09-15-operational-command-tests-reach-github.md
  - 2026-09-16-the-adapter-fixture-symlinks-a-relative-test-binary-path.md
regeneration: []
promoted: []
pull_request: "258"
delivery_commit: fc296df02b07b3e56059878caf8cea58b0aa0de9
---

# A deterministic suite under load

QA gates run the repository Verification on a machine that is also running other Runs, and in the week to 2026-09-25 four Runs were spent on tests that pass in isolation. Spec 0163's gate failed on `TestRunImplementQueuedCancellationStartsNoChildAndKeepsResumableTasks` ("timed out waiting for 2 Agent starts; got 1" after 91 s) and then on `TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork` ("expected Clean exit, got 1"). Spec 0164's gate failed on `TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement`. A documentation-only `make verify` failed on `TestBootstrapSerializesAcrossSiblings`. All four pass three times out of three in isolation, on the Spec branches and on `main`.
