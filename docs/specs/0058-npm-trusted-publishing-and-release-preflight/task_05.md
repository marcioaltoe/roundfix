---
task: task_05
spec: 0058-npm-trusted-publishing-and-release-preflight
status: pending
type: docs
complexity: medium
---

# Task 05: Document the migration, window, and vocabulary

## Overview

The OIDC migration depends on six per-package trusted-publisher configurations
created by hand on the registry's website, and the token fallback is a state
with an exit condition that only evidence can satisfy. Neither is discoverable
from the workflow. This task records the setup, the rollback window, the
rehearsal trigger, and the failure vocabulary in the release runbook, and adds
the Spec's new vocabulary to the glossary.

## Requirements

1. MUST document the per-package trusted-publisher setup as an ordered
   maintainer procedure covering the launcher and all five platform packages,
   naming the repository and workflow filename each configuration must match.
2. MUST document that the registry validates a trusted-publisher configuration
   only at publish time, so a typo surfaces as an authentication failure rather
   than at setup.
3. MUST document the fallback window: what it protects against, how it is
   switched off, and the exact evidence that closes it — one release whose
   fallback record is empty.
4. MUST document the rehearsal trigger as the way to exercise the preflight
   without cutting a tag.
5. MUST document the failure vocabulary, giving each prefix its meaning, the
   stage that emits it, and the maintainer's first recovery step.
6. MUST add the Spec's new canonical terms to the glossary, at minimum
   Publication Preflight and Release Set, in the vocabulary the workflow and
   Spec already use.
7. MUST NOT alter the release procedure the runbook already documents beyond
   what these additions require, and MUST NOT change any workflow file.

## Subtasks

- [ ] Write the trusted-publisher setup procedure for all six coordinates.
- [ ] Document the fallback window and its evidence-based exit condition.
- [ ] Document the rehearsal trigger.
- [ ] Document the failure vocabulary with recovery steps.
- [ ] Add the new canonical terms to the glossary.

## Acceptance Criteria

- [ ] The runbook names all six coordinates requiring trusted-publisher setup
      and states the repository and workflow filename each must match.
- [ ] The runbook states that configuration is validated only at publish time.
- [ ] The runbook states the fallback window's exit condition as a release
      whose fallback record is empty.
- [ ] The runbook documents the rehearsal trigger as publish-free.
- [ ] The runbook documents every failure prefix the workflow emits with a
      recovery step for each.
- [ ] The glossary defines Publication Preflight and Release Set.
- [ ] `git status --porcelain` shows no path outside the release runbook, the
      glossary, and this task file.

## Context

- instruction: `docs/agents/docs-layout.md`
- interface: `docs/user-guide/release-runbook.md`
- interface: `CONTEXT.md`

## Verification

- `grep -qi 'trusted publish' docs/user-guide/release-runbook.md` — expected:
  exit 0; the migration is documented.
- `grep -q 'roundfix' docs/user-guide/release-runbook.md && grep -q 'cli-darwin-arm64' docs/user-guide/release-runbook.md && grep -q 'cli-win32-x64' docs/user-guide/release-runbook.md`
  — expected: exit 0; the coordinate list is complete at both ends.
- `grep -q 'release.yml' docs/user-guide/release-runbook.md` — expected: exit 0;
  the workflow filename each configuration must match is named.
- `grep -qi 'rollback window\|fallback window' docs/user-guide/release-runbook.md`
  — expected: exit 0; the window is documented.
- `grep -q 'registry:' docs/user-guide/release-runbook.md && grep -q 'identity:' docs/user-guide/release-runbook.md && grep -q 'undetermined:' docs/user-guide/release-runbook.md && grep -q 'runtime:' docs/user-guide/release-runbook.md`
  — expected: exit 0; every failure prefix is documented.
- `grep -q 'Publication Preflight' CONTEXT.md` — expected: exit 0.
- `grep -q 'Release Set' CONTEXT.md` — expected: exit 0.
- `git diff --name-only HEAD -- .github/ | grep -q . && exit 1 || exit 0` —
  expected: exit 0; this task changed no workflow file.

## References

- `_prd.md` → User Experience (runbook documents OIDC setup, rollback window,
  preflight vocabulary); Core Features 4.
- `_techspec.md` → Build Order 6; Risks & Considerations.
- ADR-0084.
