---
task: task_07
spec: 0058-npm-trusted-publishing-and-release-preflight
status: pending
type: docs
complexity: low
---

# Task 07: Document the registry-side token shutdown

## Overview

Closing the fallback window currently reads as a repository-side operation:
remove the variable, remove the secret, remove the fallback branch. The Spec's
whole purpose is that a long-lived credential can no longer publish these
packages, and that only happens when token publication is disallowed on the
registry itself, per package. This task completes the closing procedure so the
window ends with the credential actually powerless.

## Requirements

1. MUST document, as part of closing the fallback window, that token
   publication is disallowed on the registry for the launcher and all five
   platform packages, naming it as a per-package setting.
2. MUST place the registry-side step in the closing procedure alongside the
   repository-side removals, so the two are performed together rather than
   read as alternatives.
3. MUST state the ordering constraint: the registry-side shutdown happens only
   after a release whose fallback record is empty, because disallowing tokens
   while the fallback is still relied upon would break the release path.
4. MUST state how the maintainer confirms the shutdown took effect.
5. MUST NOT change the preflight, publication, or rollback-window behavior
   already documented, and MUST NOT modify any workflow file.

## Subtasks

- [ ] Add the per-package registry-side token shutdown to the closing
      procedure.
- [ ] State the ordering constraint against the empty fallback record.
- [ ] State the confirmation step.

## Acceptance Criteria

- [ ] The runbook's window-closing procedure includes disallowing token
      publication on the registry, described as a per-package setting covering
      all six coordinates.
- [ ] The procedure states that the registry-side shutdown follows a release
      with an empty fallback record.
- [ ] The procedure states how the maintainer confirms the shutdown.
- [ ] The repository-side removals already documented remain present and
      unchanged in meaning.
- [ ] `git status --porcelain` shows no path outside the release runbook and
      this task file.

## Context

- interface: `docs/user-guide/release-runbook.md`

## Verification

- `grep -qiE 'disallow|disable' docs/user-guide/release-runbook.md` — expected:
  exit 0; the shutdown is described.
- `grep -qi 'token' docs/user-guide/release-runbook.md` — expected: exit 0.
- `grep -qi 'fallback' docs/user-guide/release-runbook.md` — expected: exit 0;
  the ordering constraint references the fallback record.
- `grep -q 'cli-darwin-arm64' docs/user-guide/release-runbook.md && grep -q 'cli-win32-x64' docs/user-guide/release-runbook.md`
  — expected: exit 0; the coordinate list the setting applies to is complete.
- `rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-07-31-qa-01/docs_scope_probe.rb`
  — expected: the probe that proved the omission now finds the registry-side
  shutdown documented.
- `git diff --name-only HEAD -- .github/ | grep -q . && exit 1 || exit 0` —
  expected: exit 0; this task changed no workflow file.

## References

- `_prd.md` → Core Features 4.
- `qa/qa-report-2026-07-31.md` → QA-003.
- `docs/findings/2026-07-25-npm-trusted-publishing-and-release-preflight.md` →
  the Trusted Publishing recommendation this Spec implements.
