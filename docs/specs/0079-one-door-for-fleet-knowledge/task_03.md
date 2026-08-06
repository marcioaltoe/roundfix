---
task: task_03
spec: 0079-one-door-for-fleet-knowledge
status: pending
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
