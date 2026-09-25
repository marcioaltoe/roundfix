---
status: proposed
granted: null
action: document proved model capability and readiness semantics without changing configured profiles or tool versions
consuming: 0123-runtime-readiness-and-model-capabilities
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
---

# Proposed authority for Spec 0123

This document is a concrete proposal, **not a grant**. Authoring the plan does
not authorize these protected mutations or turn the PRD's open decisions into
approved policy. Do not execute implementation from this record.

## Proposed bounded mutation

Only the shipped command contract is protected in this bounded proposal. Product catalog entries and runner/source tests are ordinary code and are intentionally absent from paths. No adapter dependency, package lock, runtime version pin, Makefile, or repository profile is authorized by this record. Any later need for one requires a new exact proposal.

The PRD states the proposed behavior. Its ordinary implementation-source
inventory belongs in the later TechSpec, outside this historical protected
path list. Every listed path is an existing file except an explicitly named
new ownership declaration; no wildcard or directory-wide grant is proposed.

## Proposed regeneration

- `make skills-sync`
- `make baseline-digests`

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
