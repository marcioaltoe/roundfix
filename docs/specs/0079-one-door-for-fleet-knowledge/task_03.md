---
task: task_03
spec: 0079-one-door-for-fleet-knowledge
status: completed
type: docs
complexity: medium
---

# Task 03: Sweep the legacy findings behind their first rollups

## Overview

One sweep, rollup-first, reviewable in a single diff. The first rollups are
written as findings of `kind: rollup`, every legacy finding under
`docs/findings/` receives its lifecycle status, absorbed findings move to
`docs/findings/_archived/` carrying their `absorbed_by` license, and anything
no rollup references defaults to `deferred`. The sweep lands before the
checks (task_04) so the repository is conformant on the day the rules start
demanding it.

## Requirements

1. MUST group the existing findings by theme and write rollups consolidating
   the groups that have genuinely related members — quality of grouping over
   quantity of rollups; a rollup is the readable front page of what was
   learned in its area, following the contract task_02 adopted into the
   layout guide.
2. MUST give every rollup `kind: rollup` and a non-empty `members:` list of
   resolvable finding basenames, and MUST give the rollup itself a lifecycle
   status like any finding.
3. MUST stamp every finding under `docs/findings/` with a lifecycle
   `status:` from exactly `pending`, `partial`, `deferred`, or `done`,
   choosing from the finding's own content: findings whose observations a
   shipped Spec provably absorbed are `done`; findings a rollup absorbs and
   supersedes are `done` and archived; unreferenced findings default to
   `deferred`.
4. MUST move each absorbed finding to `docs/findings/_archived/` via
   `git mv`, adding `absorbed_by:` pointing at the absorbing rollup basename
   or Spec folder slug; the move preserves the basename.
5. MUST stamp by editing frontmatter only — every finding's observation body
   stays byte-identical; a rollup summarizes, it never rewrites its members.
6. MUST rewrite repository links that point at moved findings, excluding
   `docs/specs/_archived/` (immutable history); links from archived Specs
   are reported in the Result instead of rewritten.
7. MUST leave the active set as live work only; the fifteen-or-fewer target
   is the QA gate's observation, not this task's assertion — do not force
   archival to hit a number.

## Subtasks

- [ ] Read and group the findings; draft the rollup set.
- [ ] Write rollups; stamp every finding's lifecycle.
- [ ] Archive absorbed findings with licenses; rewrite links.

## Acceptance Criteria

- [ ] Every file under `docs/findings/` carries a valid lifecycle status.
- [ ] At least one rollup exists; every rollup's members resolve.
- [ ] Every archived finding's license resolves to its rollup or Spec.
- [ ] No observation body changed; the diff shows frontmatter additions and
      moves only.

## Context

- instruction: docs/agents/docs-layout.md

## Verification

- `output="$(grep -L '^status:' docs/findings/*.md 2>/dev/null)"; [ -z "$output" ]`
  — expected: exit 0; every active finding is stamped (fails red today on
  the unstamped legacy set).
- `output="$(grep -h '^status:' docs/findings/*.md docs/findings/_archived/*.md 2>/dev/null | sort -u | grep -vE '^status: (pending|partial|deferred|done)$')"; [ -z "$output" ]`
  — expected: exit 0; only the four contract lifecycle values occur.
- `grep -q '^kind: rollup' docs/findings/*.md`
  — expected: exit 0; at least one rollup exists among the active findings.
- `status=0; for r in docs/findings/*.md; do grep -q '^kind: rollup' "$r" || continue; members="$(awk '/^members:/{f=1;next} f&&/^  - /{print $2} f&&!/^  - /{f=0}' "$r")"; [ -n "$members" ] || { echo "empty rollup: $r"; status=1; }; for m in $members; do [ -e "docs/findings/$m" ] || [ -e "docs/findings/_archived/$m" ] || { echo "unresolved member: $r -> $m"; status=1; }; done; done; exit $status`
  — expected: exit 0; every rollup declares members and every member
  resolves, active or archived.
- `status=0; found=0; for f in docs/findings/_archived/*.md; do [ -e "$f" ] || continue; found=1; t="$(sed -n 's/^absorbed_by:[[:space:]]*//p' "$f" | head -1)"; [ -n "$t" ] || { echo "unlicensed: $f"; status=1; continue; }; [ -e "docs/findings/$t" ] || [ -d "docs/specs/$t" ] || [ -d "docs/specs/_archived/$t" ] || { echo "unresolved license: $f -> $t"; status=1; }; done; [ "$found" = "1" ] || { echo "no archived findings"; status=1; }; exit $status`
  — expected: exit 0; the archive exists and every archived finding's
  license resolves.
- `status=0; for f in docs/findings/_archived/*.md; do [ -e "$f" ] || continue; t="$(sed -n 's/^absorbed_by:[[:space:]]*//p' "$f" | head -1)"; case "$t" in *.md) grep -q "$(basename "$f")" "docs/findings/$t" 2>/dev/null || { echo "rollup does not name back: $f"; status=1; };; esac; done; exit $status`
  — expected: exit 0; every rollup-absorbed finding is named back by its
  rollup — the relation is bidirectional.

## References

- `_prd.md` → Core Features 4–6; User Story 3; Success Metrics (active
  findings; lifecycle coverage).
- `_techspec.md` → Testing Approach (relations, never counts); Build
  Order 2.
- `references/2026-08-06-findings-accumulate-faster-than-they-become-specs.md`
  → the measured accumulation and the rollup prior art.

## Result

### Implementation

- Consolidated the legacy corpus into six active, `pending` rollups covering
  Baseline and derived tooling, Agent selection and execution environments,
  Run lifecycle and branch integrity, QA gates and Verification evidence, Spec
  authoring and contract enforcement, and review and delivery convergence.
  Their frontmatter names 63 unique legacy members.
- Stamped all 65 legacy findings `done`, updated their lifecycle date, and
  added an `absorbed_by:` license in frontmatter. Sixty-three licenses point
  to one of the six rollups. The npm publishing and Spec-owned reference
  lifecycle findings retain their more precise shipped-Spec ownership through
  `0058-npm-trusted-publishing-and-release-preflight` and
  `0060-spec-owned-reference-lifecycle`.
- Moved all 65 absorbed findings with `git mv` into
  `docs/findings/_archived/`, preserving every basename. No unreferenced
  finding remained to receive the `deferred` default; the active set consists
  only of the six live rollup fronts.
- Rewrote 16 current references across 12 handoff, workflow, user-guide,
  review-evidence, and speccheck fixture files to the `_archived/` paths.
  Provenance in the adopted-source index remains the original pre-adoption
  path by contract.
- Preserved immutable archived Spec history. The retained inventory contains
  56 references across these 37 archived Spec artifacts:
  - `0037-terminal-outcome-integrity/{_prd.md,task_06.md}`
  - `0038-terminal-run-worktree-reconciliation/task_06.md`
  - `0039-review-source-evidence-and-detached-outcomes/{_prd.md,task_08.md,qa/qa-report-2026-07-28-04.md}`
  - `0046-public-context-driven-baseline-command/{_prd.md,task_06.md,task_17.md}`
  - `0047-context-driven-guidance-composition/_prd.md`
  - `0048-context-driven-project-decisions-and-spec-constraints/_prd.md`
  - `0052-claude-adapter-standardization/_prd.md`
  - `0053-qa-gate-reachability-and-verdict-semantics/_prd.md`
  - `0054-tooling-task-and-verification-hygiene/{_prd.md,qa/evidence/2026-07-29-supervisor-blocked-rows.md}`
  - `0055-owner-identity-without-fork/_prd.md`
  - `0056-profiles-configure-merge-semantics/_prd.md`
  - `0057-baseline-capability-evidence-and-retention/_prd.md`
  - `0058-npm-trusted-publishing-and-release-preflight/{_prd.md,task_07.md,qa/qa-report-2026-07-31.md}`
  - `0059-run-storage-compaction-and-global-sanitation/_prd.md`
  - `0060-spec-owned-reference-lifecycle/{_prd.md,task_03.md,qa/evidence/2026-07-31-spec-0060/lifecycle-rehearsal.md}`
  - `0061-repository-derived-skill-requirements/{_prd.md,qa/evidence/2026-07-29-supervisor-full-access-gate.md}`
  - `0062-baseline-digest-regeneration-bootstrap/_prd.md`
  - `0063-qa-cycle-economics/_prd.md`
  - `0065-loop-order-and-verification-honesty/_prd.md`
  - `0066-run-teardown-reclaims-what-it-created/_prd.md`
  - `0067-derived-artifact-regeneration-boundary/{_prd.md,task_01.md}`
  - `0068-spec-close-audit/{_prd.md,task_05.md}`
  - `0070-declared-unreachable-acceptance/task_04.md`
  - `0077-a-green-check-is-not-a-review/_prd.md`

### Focused checks

- A custom frontmatter and relation audit reported
  `active=6 archived=65 rollup_members=63 errors=0`. It checked every active
  lifecycle, every archived `done` status and license, member resolution,
  bidirectional rollup membership, unique rollup ownership, and both
  Spec-folder licenses.
- A `HEAD`-based exact body audit compared each legacy file from its first H1
  through EOF after the move and reported
  `focused body audit: checked=65 mismatches=0`.
- A current-link audit, excluding immutable archived Specs, archived finding
  observation bodies, and the adopted-source provenance index, reported
  `focused current-link audit: errors=0` after the 16 rewrites.
- A separate archived-Spec inventory reported
  `archived Spec references retained: refs=56 files=37`, producing the file
  list recorded above.
- The rollup prose scan found none of the technical-writing skill's banned
  filler terms. `rtk git diff --check` exited 0 after the implementation edits.
- The commands under this Task's `## Verification` and the repository-wide
  `make verify` gate were not run; the Daemon owns those checks and settlement.

### Acceptance-criterion evidence

1. The focused relation audit parsed every active and archived frontmatter and
   found only valid lifecycle values, with zero errors across 71 files.
2. Six active files declare `kind: rollup`; their 63 unique members all resolve
   under `_archived/`, and no member appears in more than one rollup.
3. Every archived finding has a resolvable license: 63 resolve to active
   rollups that name the basename back, and two resolve to archived Spec
   folders 0058 and 0060.
4. The exact `HEAD` comparison found zero body mismatches across all 65 legacy
   findings. Their diffs contain only frontmatter lifecycle/license edits and
   basename-preserving moves; rollup summaries are new files and do not alter
   member observations.
