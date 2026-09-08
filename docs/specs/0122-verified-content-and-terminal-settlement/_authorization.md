---
status: proposed
granted: null
action: align verified output, QA settlement and user-authorized archive-only QA overrides
consuming: 0122-verified-content-and-terminal-settlement
paths:
  - internal/spec/archive.go
  - internal/spec/archive_test.go
  - internal/cli/archive.go
  - internal/cli/archive_test.go
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/archive-spec/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/qa-gate/SKILL.md
  - skills/archive-spec/SKILL.md
  - internal/baseline/assets/modules/spec-workflow.json
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-tasks/references/task-template.md
  - skills/write-tasks/SKILL.md
  - skills/write-tasks/references/task-template.md
  - .agents/skills/implement-task/SKILL.md
  - skills/implement-task/SKILL.md
  - internal/speccheck/coherence.go
  - docs/agents/spec-routing.md
---

# Proposed authority for Spec 0122

This document is a concrete proposal, **not a grant**. Authoring the plan does
not authorize these protected mutations or turn the PRD's open decisions into
approved policy. Do not execute implementation from this record.

## Proposed bounded mutation

The archive Go files are historically governed. The narrow canonical QA Archive Override policy is separately approved in the dated reference; Go, tests and owned-skill integration remain proposed here. The final implementation must preserve all gates outside the authorized archive-only QA exception. Ordinary daemon and Task model source belongs in the later Implementation Design, not this protected-path list.

The expanded proposal covers the existing archive-layout clause and the
Spec-routing clauses needed for the explicit repair-entry and independent
Verification-group contracts. Task templates and implementation instructions
are bounded above. The named layout/routing guides and Setup Manifest are
managed regeneration outputs, as are the shipped skill copies. None may be edited manually to bypass the approved
canonical source or its derived checks.

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

## Scope refinement from complete triage

The technical candidate now names the additional governed paths above for the
newly exposed residuals. Their presence remains a proposal, not an amendment to
an approved grant. No current tooling, test, template or generator file was
changed by this planning operation. Sanctioned outputs follow source approval.


## Confirmed QA archive override policy — 2026-09-08

The [narrow canonical grant](references/2026-09-08-authorized-qa-archive-override.md)
covers three modules and four generated guides/manifest. The maintainer permits
archival with a QA override when explicitly requested or authorized; that does
not approve this broader source/test/skill proposal or override any current Spec.

The implementation must support approval provenance, unchanged QA/Task evidence,
a QA-only archive exception and intact non-QA/source-integrity gates. It must
not convert an overridden archive into QA pass, Task completion, Run Clean,
review approval or permission for publication. The public archive parser and
its test are explicitly named above for later implementation; final flags and
metadata fields beyond `qa_override` remain to be authored.
