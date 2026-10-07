---
status: done
created_at: 2026-08-06
updated_at: 2026-09-08
kind: rollup
absorbed_by: 0122-verified-content-and-terminal-settlement
members:
  - 2026-08-06-a-test-mutates-the-repository-another-test-is-reading.md
  - 2026-08-26-a-qa-tasks-authored-verification-is-never-executed.md
  - 2026-08-26-the-daemon-has-no-postcondition-so-a-task-settles-clean-breaking-the-tree.md
  - 2026-07-28-qa-gate-cannot-reach-pull-request-journeys.md
  - 2026-07-28-same-day-qa-reruns-are-ignored-by-the-verdict-selector.md
  - 2026-07-29-qa-cycle-cost-is-cold-environments-and-agent-turns.md
  - 2026-07-29-qa-cycle-latency-and-detector-placement.md
  - 2026-07-29-qa-gate-round-economics.md
  - 2026-07-30-the-autonomous-loop-orders-qa-before-its-own-preconditions.md
  - 2026-07-31-a-rehearsal-task-can-settle-completed-without-rehearsing.md
  - 2026-08-02-a-verification-naming-a-missing-test-passes-vacuously.md
  - 2026-08-03-a-200ms-attach-budget-fails-under-ci-load.md
  - 2026-08-03-a-recording-order-defect-blocked-every-public-qa-row.md
  - 2026-08-04-a-spec-archives-with-pass-while-a-user-story-was-never-exercised.md
  - 2026-08-04-a-static-gate-row-reported-one-instance-per-cycle.md
  - 2026-08-04-fail-fast-verification-spends-the-single-repair-turn-on-the-first-of-n-defects.md
  - 2026-08-04-pre-contract-spec-graphs-run-with-no-qa-gate-and-say-nothing.md
  - 2026-08-05-an-absence-grep-rejects-the-work-it-was-written-to-protect.md
  - 2026-08-05-archive-refuses-a-spec-whose-graph-declined-the-qa-gate.md
  - 2026-08-05-authored-verification-gates-are-untested-code.md
  - 2026-08-06-the-gate-checks-that-adrs-were-cited-not-that-they-were-obeyed.md
  - 2026-08-06-the-qa-gate-and-the-pull-request-cannot-both-be-current.md
---

# QA gates and Verification evidence — authored checks are production code (2026-08-06)

The QA findings show that a green command can still miss the named test, select the wrong report, stop before the important journeys, or prove a citation instead of the behavior it stands for. A gate is trustworthy only when its own discovery, shell semantics, evidence order, and reachable journeys are tested.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-06-rollup-qa-gates-and-verification-evidence.md`.
