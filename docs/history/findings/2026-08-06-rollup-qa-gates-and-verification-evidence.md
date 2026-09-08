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

The QA findings show that a green command can still miss the named test, select
the wrong report, stop before the important journeys, or prove a citation
instead of the behavior it stands for. A gate is trustworthy only when its own
discovery, shell semantics, evidence order, and reachable journeys are tested.

## Consolidated learning

- A Verification must prove that its target exists and ran; negative searches,
  filters, filenames, and report selection need direct characterization.
- Detectors belong near the defect. Deferring every contract to a cold terminal
  QA cycle serializes discovery and spends repair turns on the first visible
  failure.
- A gate must report independent static findings together and preserve public
  journey evidence even when one governance row fails.
- Pull Request journeys and environment-blocked rows need an authored ordering
  and equivalent-evidence contract; a `pass` verdict cannot stand in for an
  unexercised user story.

## Live edge

Specs 0053, 0063, 0065, 0070, 0071, and 0072 absorbed several measured defect
classes. The rollup remains `pending` because gate authoring itself still needs
mechanical discovery and contract checks before Task settlement relies on it.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.

## Addendum — 2026-09-08 — Current triage

Existing deliveries cover reachability/report selection (0053, 0070, 0072),
mechanical and staged checks (0064, 0080, 0093), Verification probes (0090,
0095, 0116), derived QA Verification and combined static/journey findings
(0105), and current-report/precondition-report handling (0113).

Current `internal/daemon/task_engine.go` still completes the QA Task only for
`pass`, while `internal/spec/archive.go` admits the declared-only `partial`
case. The Task verification path checks authored commands before committing,
with no general postcondition over the resulting committed tree. The new
Finding `2026-09-08-verified-executable-source-is-dropped-before-settlement.md`
records a concrete loss at that boundary. Route these to provisional P3,
verified content and terminal settlement.

Gate-cost distribution, load-sensitive checks, and historical test interference
need current measurements in provisional P5. Broader declaration-to-evidence
traceability belongs to provisional P1/P8. No current load experiment or
terminal QA was performed by this triage; this Rollup remains the license for
its 22 archived members.

## Addendum — 2026-09-08 — Complete implementation routing

The maintainer selected the residual work for the queue. Its primary owner is
[0122-verified-content-and-terminal-settlement](../../specs/0122-verified-content-and-terminal-settlement/_prd.md).
- [0124-verification-capacity-and-measured-economics](../../specs/0124-verification-capacity-and-measured-economics/_prd.md): verification economics.
- [0129-spec-authoring-and-gate-recovery](../../specs/0129-spec-authoring-and-gate-recovery/_prd.md): promise and ADR traceability.
- [0128-release-planning-with-bare-stable-tags](../../specs/0128-release-planning-with-bare-stable-tags/_prd.md): bare stable version tags.

`done` records complete routing to implementation Specs, not a passing repair or
QA verdict. This Rollup remains in the active findings directory because its
archived members still name this basename as their absorption license. Those
licenses and original observations are preserved; retirement waits for the
durable replacement contract in 0120. Shipped mechanisms remain regression
obligations rather than duplicate implementation Tasks.

## Addendum — 2026-09-08 — Routing document removed

The maintainer requested removal of `docs/workflow/` and its routing documents.
The earlier citation remains a dated historical observation; its original bytes
can be read at Git revision `6b8ea48725cbca13974eee0b400b3482202874f6`.
The current primary and secondary Spec owners remain those in the complete-triage
addendum above; removing the old plan does not reopen or erase the members.

## Addendum — 2026-09-08 — Archived after complete routing

The maintainer requires terminal Findings and Rollups to leave the active
family directory. Every member now points directly to an existing active or
archived Spec, and this Rollup has its own direct Spec absorber. The earlier
statements retaining this file as an active license root are superseded by
this completed routing migration. Original observations and prior pointers
remain recorded; archival does not claim implementation of pending Specs.
