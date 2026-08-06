---
task: task_05
spec: 0079-one-door-for-fleet-knowledge
status: completed
type: docs
complexity: medium
---

# Task 05: Pilot the door with this Spec's own material

## Overview

The proof the PRD requires before any clause binds: the save–retrieve–triage
cycle exercised end to end with real material, witnessed by a committed
pilot report in this repository. Two captures go through the door — this
Spec's research digest destined to the brain itself, and at least one
cross-project observation destined to `roundfix` — and the roundfix-destined
entry is retrieved and triaged into this repository. Verification asserts
only this repository's own outputs; brain-side facts are witnessed as
recorded evidence in the report.

**External gate**: the brain-side `inbox/` namespace, its `AGENTS.md`
carve-out, and the qmd collection must exist under the secondbrain
repository's own contract. When any of these is absent, or the Agent Session
cannot write outside its worktree, record the exact blocker in the report
draft and settle failed naming the dependency — never fake the proof.

## Requirements

1. MUST verify the external gate first: the brain's inbox namespace exists
   and its local guide permits inbox creation; on absence, stop with the
   named blocker.
2. MUST capture this Spec's research digest (sources cited) as an inbox
   entry under the brain's own namespace with `destination: secondbrain`,
   running the advisory qmd duplicate check first and recording its result;
   a strong match extends existing knowledge instead of duplicating it.
3. MUST capture at least one cross-project observation with
   `destination: roundfix`, `capture: manual`, and a truthful `origin`,
   following the entry contract task_02 shipped.
4. MUST commit each capture in the brain at capture time and record, in the
   pilot report, each entry's brain-relative path, its brain commit hash,
   and the measured capture-to-durability seconds.
5. MUST retrieve and triage the roundfix-destined entry in this repository:
   the entry carries intent, so triage mints one typed Backlog Entry under
   `docs/backlog/` per the ADR-0092 boundary, citing the inbox entry's
   brain-relative path as provenance; the brain-side entry moves to
   `_triaged/` gaining `resolved_to`, witnessed in the report. A new
   Finding is out of this task's scope — the sweep owns the findings
   directory in this wave.
6. MUST audit non-interference: no capture touched any project checkout;
   the report records the destination checkout's clean `git status` before
   triage.
7. MUST write the pilot report at
   `docs/specs/0079-one-door-for-fleet-knowledge/pilot-report.md` carrying,
   as literal tokens: `destination: secondbrain`, `destination: roundfix`,
   `brain commit: <hash>` per capture, `capture_to_durability_seconds:` per
   capture, the qmd check result, `resolved_to`, and the friction
   observations the maintainer needs before obligations bind (task_06).

## Subtasks

- [ ] Gate check, then the two captures with brain commits.
- [ ] Retrieve, triage, and mint the Backlog Entry with provenance.
- [ ] Write the pilot report with recorded evidence and friction notes.

## Rehearsal Cases

- Case: research digest captured for the brain itself after the advisory qmd
  duplicate check; Observation: inbox entry committed under the brain's own
  namespace, check result and durability seconds recorded in the report.
- Case: cross-project observation captured for roundfix and triaged in the
  destination; Observation: typed Backlog Entry minted in `docs/backlog/`
  citing the entry's brain path; entry resolved under `_triaged/` with
  `resolved_to`, per the report.
- Case: capture attempted while the destination has an Active Run;
  Observation: capture lands only in the brain, destination checkout
  untouched — clean `git status` recorded before triage.

## Acceptance Criteria

- [ ] The pilot report exists with every literal token from Requirement 7.
- [ ] A typed Backlog Entry exists citing its inbox provenance.
- [ ] Both captures' durability measurements are under the one-minute
      target, or the miss is recorded with its cause.
- [ ] On a blocked gate: no fabricated evidence — the report draft names the
      missing dependency and the task settles failed.

## Context

- instruction: docs/agents/secondbrain.md
- instruction: docs/agents/docs-layout.md

## Verification

- `test -f docs/specs/0079-one-door-for-fleet-knowledge/pilot-report.md`
  — expected: exit 0; the committed pilot report exists.
- `grep -q 'destination: secondbrain' docs/specs/0079-one-door-for-fleet-knowledge/pilot-report.md && grep -q 'destination: roundfix' docs/specs/0079-one-door-for-fleet-knowledge/pilot-report.md`
  — expected: exit 0; both destinations were exercised.
- `grep -qE 'brain commit: [0-9a-f]{7,40}' docs/specs/0079-one-door-for-fleet-knowledge/pilot-report.md`
  — expected: exit 0; durability is witnessed by recorded brain commits.
- `grep -q 'capture_to_durability_seconds:' docs/specs/0079-one-door-for-fleet-knowledge/pilot-report.md && grep -qi 'qmd' docs/specs/0079-one-door-for-fleet-knowledge/pilot-report.md`
  — expected: exit 0; the durability audit and the advisory duplicate check
  are recorded.
- `grep -q 'resolved_to' docs/specs/0079-one-door-for-fleet-knowledge/pilot-report.md`
  — expected: exit 0; the triage resolution is witnessed.
- `grep -rq 'inbox/roundfix/' docs/backlog`
  — expected: exit 0; the minted Backlog Entry cites its inbox provenance —
  the save–retrieve–triage relation closed in this repository.

## References

- `_prd.md` → User Stories 1, 2, 5; Core Features 1–3, 9; Success Metrics
  (durability; interference; research capture); Decisions (the pilot
  precedes the rules).
- `_techspec.md` → Integration Points; Testing Approach (pilot hermeticity);
  Build Order 4; Risks.
- ADR-0095, ADR-0092.

## Result

### Implementation

- Verified the Secondbrain gate before capture: the two destination
  namespaces existed, the brain's own `AGENTS.md` permitted create-only inbox
  writes, the entry contract existed, and the advisory qmd collection returned
  results after the Agent Session received filesystem permission.
- Captured the Spec research digest at
  `inbox/secondbrain/2026-08-06-fleet-knowledge-door-research-digest.md` with
  cited Roundfix and external sources. Brain commit
  `93bd7b35b77a356a3d2d50989de93565b1e85dbb` made it durable in 50 seconds.
- Captured the cross-project intent at
  `inbox/roundfix/2026-08-06-automate-inbox-capture-durability.md` with
  `origin: secondbrain`, `destination: roundfix`, and `capture: manual`.
  Brain commit `e982c0137ad61b5cf82de0fa8f63b9c5c325412e` made it durable in
  24 seconds.
- Recorded the clean pre-triage destination checkout, then minted exactly one
  `feat` Backlog Entry with inbox provenance at
  `docs/backlog/2026-08-06-atomic-inbox-capture-helper.md`. Moved the consumed
  brain entry under `_triaged/`, added `resolved_to`, and made that brain state
  durable in commit `4aa955c3ec2c03ef5799ea0745bc8a8938ce6530`.
- Added `pilot-report.md` with the gate, qmd result, both paths and capture
  commits, durability measurements, non-interference evidence, triage
  resolution, friction, and explicitly unexercised scope.

### Focused checks

- The initial advisory qmd query reached the external tool but failed with
  `SQLITE_CANTOPEN` under the restricted Agent filesystem. The authorized retry
  returned ranked collection results and the 248-document embedding warning;
  inspection found no substantive existing digest to extend.
- Fresh `rtk git -c core.fsmonitor=false status --short
  --untracked-files=all` in `/Users/marcio/dev/roundfix` produced no paths
  immediately before triage. The same check in the Secondbrain was clean after
  capture and triage commits.
- A focused `rtk rg` report scan found both literal destinations, both 40-byte
  `brain commit:` values, both durability fields, `qmd check result:`, and
  `resolved_to:`. A separate Backlog scan found the `feat`/`open` frontmatter
  and exact pending-entry provenance.
- `rtk git show` at each capture commit reproduced the corresponding pending
  entry, and `rtk git show HEAD:` reproduced the triaged entry with the exact
  `resolved_to` artifact path.
- `rtk git diff --check` exited 0 for the tracked Task-file diff. A focused
  trailing-whitespace scan found no matches in the Task file, pilot report, or
  Backlog Entry.
- The Daemon-owned commands under `## Verification` were not run in this Agent
  turn. Task status remains Daemon-owned.

### Acceptance-criterion evidence

1. `pilot-report.md` records both literal destinations, both full `brain
   commit:` tokens, both `capture_to_durability_seconds:` values, the qmd check
   result, and `resolved_to`, plus the maintainer friction audit.
2. `docs/backlog/2026-08-06-atomic-inbox-capture-helper.md` follows the complete
   `feat` operational contract and cites
   `inbox/roundfix/2026-08-06-automate-inbox-capture-durability.md` as its
   provenance.
3. The two capture measurements are 50 seconds and 24 seconds, both below the
   one-minute target. Their brain commits are recorded in the pilot report and
   resolve to the captured entry contents.
4. The external gate was not blocked. The report preserves the initial qmd
   filesystem denial as friction and distinguishes it from the successful
   authorized duplicate check; no brain-side fact is inferred from a missing
   artifact.

### Follow-up

- The open Backlog Entry carries the future atomic-helper intent. Implementing
  that helper, binding inbox obligations, ingesting the research digest, or
  creating a Finding remains outside this Task's slice.
