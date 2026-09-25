---
status: proposed
granted: null
action: measure and then deliberately align repository Verification coverage, fresh-run policy, and justified capacity
consuming: 0124-verification-capacity-and-measured-economics
paths:
  - Makefile
  - .github/workflows/ci-verify.yml
---

# Proposed authority for Spec 0124

This document is a concrete proposal, **not a grant**. Authoring the plan does
not authorize these protected mutations or turn the PRD's open decisions into
approved policy. Do not execute implementation from this record.

## Proposed bounded mutation

These are the only proposed test/build/CI configuration paths. Ordinary fixture and regression-test source is not added to the historical tooling set. No linter suppression, new package, system-wide runtime setting, or unrelated workflow edit is permitted. A timing or capacity number must be justified by captured evidence before the consuming Task is approved.

The PRD states the proposed behavior. Its ordinary implementation-source
inventory belongs in the later TechSpec, outside this historical protected
path list. Every listed path is an existing file except an explicitly named
new ownership declaration; no wildcard or directory-wide grant is proposed.

## Proposed regeneration

None proposed.

Generation is limited to approved source mutations and their declared owned
outputs. Read and audit the ownership declarations, preserve unrelated outputs,
and refuse manual digest edits. Where named managed guidance changes, use the
public Baseline workflow after its exact plan is reviewable. These commands
are proposed here; they are not authorization to run a mutation now.

## Proposed limits

- Approval applies only to this consuming Spec and its selected actions and
  exact protected paths. An executor cannot extend its own grant.
- No dependency/tool version upgrade, secret access, production mutation,
  release, deployment, or destructive cleanup is included.
- Live adapter or load experiments require their scope and resource/spending
  limit to be recorded before execution. Silence grants nothing.
- Preserve required checks and honest failure outcomes. Unknown causes remain
  unknown; no suppression, blanket skip, retry-until-green, or invented proof.
- Record approval before any consuming tooling commit. A prerequisite repair
  has its own earlier commit; consequent repairs follow their cause and stay
  within their approved boundary.

## Approval checkpoint

Pending. The maintainer must confirm the proposed action, exact protected
paths, unsettled product choices, and applicable experiment limits. Record the
decision and date before changing `status` or `granted`. No value is inferred
from the request to author Specs, a default answer, or a checker result.

Research and provisional source provenance are recorded in [_prd.md](_prd.md).
The planning documents may be reviewed and committed without activating this
proposed grant; implementation remains blocked on the checkpoint above.
