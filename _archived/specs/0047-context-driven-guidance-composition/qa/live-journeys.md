---
spec: 0047-context-driven-guidance-composition
status: pending
scope: live-qa
---

# Live Baseline guidance-composition journeys

This plan reserves the repository-backed acceptance evidence for a separately
authorized `qa-gate` run. Task 09 supplies hermetic release-gate coverage only;
it does not execute or claim a verdict for Fluxus or Oraculum.

## Entry criteria

- Build and record the Roundfix binary identity from the revision under QA.
- Record the source revision and starting tree state for each live repository.
- Obtain explicit authorization for each Fluxus and Oraculum checkout.
- Run the hermetic Task 09 Verification commands successfully before live QA.

## Evidence matrix

| ID | Repository journey | Required observations | Required retained evidence |
| --- | --- | --- | --- |
| LIVE-01 | Fluxus greenfield | The real binary plans and applies the selected maintained Profile; every accepted rule has one semantic owner; no manual redistribution is needed; no generic or empty repository-specific carrier or root pointer exists. Formatter and repository Verification run outside Baseline, the post-apply audit is ready, and a fresh Plan has zero file changes. | Repository revision and clean starting tree; binary identity; exact plan and apply commands with exits; Plan Digest; complete managed-entry and Upgrade Retention Contract ledgers; before/after paths and identities; formatter and repository Verification output; audit output; fresh empty Plan. |
| LIVE-02 | Fluxus update | The real binary inventories the existing Baseline and recognized carriers, preserves exact rule bytes, shows every semantic or residual disposition before confirmation, removes zero-residual legacy carriers and pointers, and applies without manual rule movement. Formatter and repository Verification run outside Baseline, the audit is ready, and a fresh Plan has zero file changes. | Repository revision and starting Setup Manifest; binary identity; exact plan and apply commands with exits; confirmed Decision Document; Plan Digest; complete managed-entry and Upgrade Retention Contract ledgers; exact redistributed byte evidence; removed or retained carrier evidence; formatter and repository Verification output; audit output; fresh empty Plan. |
| LIVE-03 | Oraculum backend-only TypeScript divergence | The built-in TypeScript Profile blocks on profile-specific divergence. A maintainer reviews the removed modules and capabilities, the repository-owned Profile adaptation remains catalog-valid, universal required capabilities remain required and satisfied, and the Profile file appears only in the confirmed Plan. Apply verifies that Profile postimage; repository Verification runs outside Baseline; the repeated audit is ready; a fresh Plan has zero file changes. | Repository revision and clean starting tree; binary identity; blocked built-in Profile audit; reviewed Profile draft and source Profile identity; exact plan and apply commands with exits; Plan Digest; Profile managed-entry and file-change ledger entries; universal-capability audit evidence; repository Verification output; fresh empty Plan. |

## Exit criteria

- All three journeys have fresh evidence from the same Roundfix revision.
- Every plan/apply command uses the real binary and records stdout, stderr, and
  exit status separately.
- Every P0 journey applies only its confirmed Plan Digest and reaches a ready
  post-apply audit.
- Fluxus greenfield and update both finish with no empty carrier and no manual
  rule redistribution.
- Oraculum finishes with a verified repository-owned Profile postimage and no
  waiver of a universal required capability.
- Formatter and repository Verification evidence postdates apply, and the
  final fresh Plan for every journey has zero file changes.

## Not executed by Task 09

No live repository was read or changed for this task. The `qa-gate` report must
record PASS, FAIL, or blocked results after the separately authorized runs and
must list any untested row with its reason.
