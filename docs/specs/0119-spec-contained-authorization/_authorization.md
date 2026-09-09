---
status: proposed
granted: null
action: establish Spec-contained authorization and command execution authority
consuming: 0119-spec-contained-authorization
paths:
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/speccheck/constraints.go
  - internal/speccheck/constraints_characterization_test.go
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-prd/references/prd-template.md
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/write-techspec/references/techspec-template.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-prd/references/prd-template.md
  - skills/write-techspec/SKILL.md
  - skills/write-techspec/references/techspec-template.md
  - skills/write-tasks/SKILL.md
  - docs/agents/agent-instructions.md
  - docs/agents/spec-routing.md
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
  - internal/speccheck/governed.go
  - internal/speccheck/governed_repocontract_test.go
---

# Proposed authority for Spec 0119

This is a reviewable proposal, **not a grant**. The maintainer requested the
outcome and authorized authoring, commit, and push of the planning documents.
Approval of the exact protected mutations above is still pending.

## Confirmed session authority

On 2026-09-08 the maintainer requested Specs containing their authorizations
and the same requirement in canonical guidance. The confirmed unattended
delivery scope is implementation, the configured review policy, PR, and automatic
merge only after Specs and their limits are approved and required checks pass.
The later confirmed policy permits codex, claude, coderabbit or explicit none;
none records intentional omission and enabled-review failures still block. It
does not approve this proposed implementation schema or waive existing gates.

## Proposed bounded mutation

Update canonical authoring and placement guidance, the Roundfix-owned PRD,
TechSpec and Task-authoring instructions/templates, and the existing
constraint-reader behavior that already lies in the governed set. Preserve
clause identities where the existing obligation can be extended. Ordinary
product-source implementation paths belong to the later TechSpec and do not
become historical tooling grants merely by being discussed here.

Run `make skills-sync` only for the approved source edits, verifying every
other shipped skill remains byte-identical. Run `make baseline-digests` for
the sanctioned derived pins, then the public Baseline update for the named
managed guides and manifest. These regeneration instructions remain proposed
until the source mutations are approved.

## Limits and commit order

- Commit this approval record separately before any consuming tooling commit.
- Do not edit an upstream-managed skill or an archived grant.
- Do not infer approval from silence, preselection, a pending question, or a
  checker that currently fails to validate the new proposed filename.
- No paid API use, release, deployment, destructive cleanup, or branch-policy
  exception is granted by this record.
- Verification remains Daemon-owned. Running commands from third-party Specs
  requires the source-trust decision recorded in the PRD's Open Questions.

## Approval evidence

Pending. Record the maintainer's concrete decision and approved bounded scope
here before changing `status` or `granted`; do not fill either from inference.

## Scope refinement from complete triage

The technical candidate now names the additional governed paths above for the
newly exposed residuals. Their presence remains a proposal, not an amendment to
an approved grant. No current tooling, test, template or generator file was
changed by this planning operation. Sanctioned outputs follow source approval.
