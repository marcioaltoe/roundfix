---
task: task_04
spec: 0069-review-run-targets-its-pull-request
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

The authored terminal gate. The Daemon executes `qa-gate` as this node once
task_03 settles `completed`.

This Spec exists because two legitimate security findings failed for an
environmental reason, so this gate's most valuable rows are the ones proving
that an environmental stop no longer reads as a defect.

## Requirements

1. MUST run only after task_03 settles `completed`.
2. MUST observe that a Run started on a branch other than the Pull Request's
   head refuses at Preflight with exit `2`, creating no Run and writing
   nothing.
3. MUST observe that a checkout moved mid-Run reaches the interruption outcome
   with its Review Issues left unsettled, rather than accepting the Task's
   claim.
4. MUST confirm every review artifact commit in the exercised Runs landed on
   the Pull Request's head branch, read from Git rather than from a log.
5. MUST confirm a Run whose checkout matches behaves as before.
6. MUST confirm Roundfix moved no working tree during any observed Run.
7. MUST classify any finding by user impact and record typed blocked-row
   counts.

## Acceptance Criteria

- [ ] The gate runs only after task_03 settles `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] The report records the Preflight refusal observed end to end.
- [ ] The report records the mid-Run interruption observed independently.

## Verification

- `sh -c 'report=$(ls -1t docs/specs/0069-review-run-targets-its-pull-request/qa/qa-report-*.md 2>/dev/null | head -n 1); test -n "$report" && rg -q "^status: closed$" "$report" && rg -q "^verdict: pass$" "$report" && rg -q "^rows_blocked_environment: [0-9]+$" "$report" && rg -q "^rows_blocked_finding: [0-9]+$" "$report" && rg -q "^rows_blocked_declared: [0-9]+$" "$report" && ! rg -q "^\\| QA-[0-9]+ \\|.*\\| (pending|planned) \\|" "$report"'`
  — expected: exit 0; the latest report is closed and passing, carries typed
  blocked-row counts, and leaves no planned row pending.
- `sh -c 'report=$(ls -1t docs/specs/0069-review-run-targets-its-pull-request/qa/qa-report-*.md 2>/dev/null | head -n 1); row=$(rg -i "^\\| QA-[0-9]+ \\| .*Preflight.*refus" "$report"); test -n "$row" && printf "%s\\n" "$row" | rg -q "\\| pass \\|" && printf "%s\\n" "$row" | rg -q "\\[[^]]+\\]\\([^)]*\\.md"'`
  — expected: exit 0; a passing result row links evidence for the Preflight
  refusal.
- `sh -c 'report=$(ls -1t docs/specs/0069-review-run-targets-its-pull-request/qa/qa-report-*.md 2>/dev/null | head -n 1); row=$(rg -i "^\\| QA-[0-9]+ \\| .*(mid-Run.*interrupt|checkout.*moved.*interrupt)" "$report"); test -n "$row" && printf "%s\\n" "$row" | rg -q "\\| pass \\|" && printf "%s\\n" "$row" | rg -q "\\[[^]]+\\]\\([^)]*\\.md"'`
  — expected: exit 0; a passing result row links evidence for the mid-Run
  interruption.
- `sh -c 'report=$(ls -1t docs/specs/0069-review-run-targets-its-pull-request/qa/qa-report-*.md 2>/dev/null | head -n 1); row=$(rg -i "^\\| QA-[0-9]+ \\| .*artifact commit.*head branch" "$report"); test -n "$row" && printf "%s\\n" "$row" | rg -q "\\| pass \\|" && printf "%s\\n" "$row" | rg -q "\\[[^]]+\\]\\([^)]*\\.md"'`
  — expected: exit 0; a passing result row links Git evidence that artifact
  commits landed on the head branch.
- `sh -c 'report=$(ls -1t docs/specs/0069-review-run-targets-its-pull-request/qa/qa-report-*.md 2>/dev/null | head -n 1); row=$(rg -i "^\\| QA-[0-9]+ \\| .*(checkout.*match.*unchanged|matching checkout.*unchanged)" "$report"); test -n "$row" && printf "%s\\n" "$row" | rg -q "\\| pass \\|" && printf "%s\\n" "$row" | rg -q "\\[[^]]+\\]\\([^)]*\\.md"'`
  — expected: exit 0; a passing result row links non-regression evidence for a
  matching checkout.
- `sh -c 'report=$(ls -1t docs/specs/0069-review-run-targets-its-pull-request/qa/qa-report-*.md 2>/dev/null | head -n 1); row=$(rg -i "^\\| QA-[0-9]+ \\| .*(working tree.*not moved|moved no working tree)" "$report"); test -n "$row" && printf "%s\\n" "$row" | rg -q "\\| pass \\|" && printf "%s\\n" "$row" | rg -q "\\[[^]]+\\]\\([^)]*\\.md"'`
  — expected: exit 0; a passing result row links evidence that Roundfix moved
  no working tree.

## References

- `_prd.md` → Goals; Success Metrics.
- `_techspec.md` → Testing Approach.
- ADR-0052, ADR-0080, ADR-0091.
