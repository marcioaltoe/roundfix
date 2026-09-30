---
status: done
created_at: 2026-09-30
updated_at: 2026-09-30
absorbed_by: 0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own
---

# Delivery: the v0.22.0 queue needed twelve manual interventions (2026-09-30)

The v0.22.0 queue (Specs 0192, 0193, 0195 and 0199) started at 11:41 on the v0.21.0 binary with five retries per item. It finished at 18:26. 0199 was the only Spec the queue carried from start to merge, and even 0199 needed two manual steps. This report groups the interventions by class, so Wave 8 ("a queue that classifies parks and recovers") can take them as input.

## 1. A Spec started before the Specs its Verification depends on had merged

- Symptom / evidence: 0195's task_04 failed twice. `TestEveryCheckTheReleaseStepNamesExists` named tests that 0192 and 0193 create. The owner started 0195 while 0192 and 0193 were parked and awaiting retry. After 0195 was rebased, 0192's Roundfix skill text sat under the `0.0.3` that 0195 had recorded, and 0195's own version check refused it. That cost a corrective task_06. Runs: `run_20260930T162501Z_8ec6b68021df9cef` and `run_20260930T203751Z_2b4c20bff8fbb572`.
- Root cause: queue order is not a dependency. A parked item does not hold back its successors, and a Spec cannot say "start only after Spec X merged".
- Action / suggestion: let a Spec declare prerequisite Specs, and have the owner start it only when those Specs are on the default branch.

## 2. Parallel items conflicted on generated Baseline files

- Symptom / evidence: PR #298 (0193) was `CONFLICTING` after 0192 and 0199 merged. The owner waited on checks that never started, then parked with `checks-timeout`. The conflicted files were the digest pin, the catalog snapshots, the four plan goldens and `docs/agents/setup-context.json`. All are derived, and `make baseline-digests` plus the managed refresh regenerate them. The single source conflict was the Standard TypeScript profile. A blanket `--ours` there dropped 0199's REST default, and `TestTheProfileHTTPDefaultMatchesTheDecisionCatalog` caught the loss. 0195's rebase hit the same generated set.
- Root cause: the owner neither detects a conflicting Pull Request nor knows which paths are regenerable.
- Action / suggestion: report `mergeable: CONFLICTING` as its own park reason at once. Resolve a conflict confined to derived paths by rebasing and regenerating through the sanctioned commands. Park only on a source conflict.

## 3. A Spec's design broke an adopter path that only an existing test knew

- Symptom / evidence: 0193 removed a duplicated backend clause entry. The Source Baseline retention classifier then reported that entry `unaccounted`, which would have refused every Standard TypeScript adopter's next `baseline update`. `TestStandardTypeScriptStructuralClauseRetention`, in the governed `plan_test.go`, failed first. Later `TestBaselineUpdateFleetSweep` failed in another package that the first `verify-changed` selection had not run. That took two corrective Tasks and a grant widening (#295).
- Root cause: the TechSpec did not trace the removal through the retention classifier. Each QA run surfaced only the failures in the packages `verify-changed` selected.
- Action / suggestion: authoring rule: a Spec that removes or renames a Baseline clause must name its retention disposition. Operating rule, now in use: simulate a corrective with a full `make verify` before `deliver retry`.

## 4. Environment and timing failures parked items that had passed

- Symptom / evidence:
  - 0192's QA was `partial` because the QA sandbox's proxy refused GitHub (HTTP 403) for the outside-evidence row. The operator fetched the source and archived with `--qa-override`, and the item was delivered by hand (#296). The queue cannot resume an archived item whose head moved.
  - Two time-bound tests failed under load:
    - `TestTaskBudgetReasonNamesTheSettlementThatRenewedIt` (`internal/daemon`, a 200 ms Run Budget that measured 266 ms on the CI runner, PR #297);
    - `TestBatchClosesOnCountLingerAndImmediate` (`internal/store`, its linger subtest, in the owner's repository gate for 0193). It passed 10 of 10 runs in isolation.
  - An orphaned `cli.test implement --spec 0001-widget-flow --detach` ran for five hours. Its working directory was under `TestRunImplementDetachSurvivesCallerProcessGroupKill*`; it was killed by PID.
- Root cause: outside-evidence rows depend on the Agent sandbox's network. Two tests assert wall-clock bounds. One test leaks the detached child it spawns.
- Action / suggestion: classify an environment-only partial whose blocked rows the operator later satisfies, and let `deliver retry` accept an operator-archived head. Replace the two wall-clock assertions with event or fake-clock waits. Make the detach test reap its child.

## 5. False-positive review findings on queue-written artifacts

- Symptom / evidence: 0199's pre-PR review raised two findings. F1 said the QA Report audits the implementation head rather than the commit that records it. F2 said a Task was marked `completed` while its Result says the Daemon owns status. Both describe the product's designed order, and both were dismissed with evidence.
- Root cause: the reviewer is not told that the QA Report commit and the Daemon's settlement writes are expected.
- Action / suggestion: give the review prompt the delivery's own commit conventions: the QA Report commit, Daemon settlement, archive.

## Addendum (2026-09-30): routing

- Sections 1, 2 and 4 (prerequisites, conflicts on derived paths, the environment-only partial and one re-run of an unrelated failed check) → Spec `0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own`.
- Section 5 (false-positive review findings) → Spec `0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds`.
- The full QA re-run after each correction, observed across sections 1 and 3 → Spec `0202-a-qa-gate-that-reruns-only-stale-rows`.
- The two time-bound tests and the leaked detached child → Backlog Entry `docs/backlog/2026-09-30-time-bound-tests-and-a-leaked-detached-child.md`.
- The authoring rule of section 3 (name the retention disposition of a removed Baseline clause) is now an operating rule of the authoring briefing; no Spec is needed.

