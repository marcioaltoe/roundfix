---
task: task_05
spec: 0058-npm-trusted-publishing-and-release-preflight
status: completed
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

## Result

Documented the maintainer-side npm Trusted Publishing migration without
changing the tag-driven release procedure or workflow. The runbook now covers
all six package bindings, publish-time-only validation, the publish-free
preflight rehearsal, the evidence-bounded fallback window, and recovery by
failure prefix. The glossary now defines the Publication Preflight and Release
Set using the workflow's canonical vocabulary.

### Focused checks

- `rtk git diff --check` — exited 0 after the final documentation and Result
  edits.
- A Ruby structural assertion over the runbook, `platforms.json`, and
  `CONTEXT.md` — exited 0. It derived the five platform coordinates from the
  package manifest, confirmed the launcher plus all five trusted-publisher
  rows and their shared owner/repository/workflow values, checked the
  publish-time validation, rehearsal, fallback exit and switch, all four
  failure-table rows, and both glossary definitions.
- `rtk git -c core.fsmonitor=false status --short` and
  `rtk git diff --name-only HEAD` — listed only `CONTEXT.md`,
  `docs/user-guide/release-runbook.md`, and this task file. The task file's
  pre-existing change is the Daemon-owned `status: in_progress` transition;
  this implementation did not alter that field.

### Acceptance evidence

1. The ordered setup procedure names `roundfix` and all five
   `@roundfix/cli-*` coordinates. Every row binds GitHub owner `marcioaltoe`,
   repository `roundfix`, and workflow filename `release.yml`; the procedure
   also states the combined repository path `marcioaltoe/roundfix` and full
   workflow path `.github/workflows/release.yml`.
2. The setup section states that npm validates a trusted-publisher binding
   only when `npm publish` attempts OIDC, and that a typo first surfaces as an
   `identity:` authentication failure.
3. The fallback section states what the window protects against, that only
   repository-variable value `1` enables the retry, and that one complete
   tagged release whose fallback record is empty is the only exit evidence. It
   then directs the maintainer to disable the variable and remove the token.
4. The rehearsal section identifies `workflow_dispatch`, its required version
   input, and its publish-free boundary before cross-compilation, npm publish,
   and GitHub Release creation.
5. The failure table defines `registry:`, `undetermined:`, `identity:`, and
   `runtime:` with the emitting workflow stage, meaning, and first recovery
   step for each.
6. `CONTEXT.md` defines Publication Preflight as the read-only whole-set
   eligibility check and Release Set as the indivisible launcher-plus-five
   publication unit.
7. The final changed-path postflight contains only the release runbook, the
   glossary, and this task file; no workflow path changed.

The commands under `## Verification` were not run; the Daemon owns that gate.
