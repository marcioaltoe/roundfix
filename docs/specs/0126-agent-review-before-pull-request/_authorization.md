---
status: proposed
granted: null
created: 2026-09-08
action: implement configurable pre-PR review with codex, claude, coderabbit or none
consuming:
  - 0126-agent-review-before-pull-request
paths:
  - .roundfixrc.yml
  - .coderabbit.yaml
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/autonomous-work.json
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/agents/openai.yaml
  - skills/roundfix/SKILL.md
  - skills/roundfix/agents/openai.yaml
  - docs/agents/agent-instructions.md
  - docs/agents/autonomous-work.md
  - docs/agents/setup-context.json
  - internal/cli/cli_test.go
  - internal/docscontract/publicdocs_test.go
  - skills/baseline_skill_contract_test.go
---

# Proposed implementation authority for configurable pre-PR review

This is a reviewable proposed record, not express maintainer authorization.
`granted: null` means no protected mutation is granted by this artifact. The
Spec-local placement is proposed by Spec 0119; the current checker must not be
assumed to recognize a new authorization schema merely because this file exists.

## Authority already expressed

The maintainer now permits codex, claude, coderabbit or none, preserving the
Codex default and explicit project precedence. This supersedes mandatory review
and complete CodeRabbit removal. Otherwise authorized delivery can omit review
with explicit none, while QA and required checks remain binding.

The [narrow canonical grant](references/2026-09-08-configurable-review-policy-authorization.md)
expressly covers the two canonical modules and generated guides/manifest only.
The implementation paths in this proposed record remain ungranted; the policy
choice does not enable paid calls or implement adapters/configuration parsing.

## Proposed protected scope

- Make CodeRabbit optional through explicit provider selection and compatible configuration migration; preserve unrelated configuration and do not silently change User Config, another repository or the current selected policy.
- Implement and align the owned skill/manifest contract for the already confirmed canonical review policy. Embedded and rendered copies must come from their owners; the narrow source-policy grant is recorded separately.
- Intentionally revise the named affected contract tests to prove the replacement behavior. The scope does not allow suppressing a failure, weakening assertions to obtain green, or changing Verification configuration.
- Keep every path exact. A subsequently discovered governed path requires an amended proposal and express approval before mutation.

## Sanctioned regeneration

After the canonical source edits are expressly approved, the existing
`make skills-sync` and `make baseline-digests` steps may regenerate the declared
copies and deterministic digest fallout according to the repository's sanctioned
rule. Repository guides must be regenerated through the public Baseline command
and its reviewed file plan. Preserve all historical artifacts outside the owned
regeneration outputs; do not hand-edit pins or invent a general fixture grant.

## Pending decisions and execution limits

- Review policy: confirmed — codex, claude, coderabbit or explicit none; default Codex and explicit project precedence. Agent model/effort follows the compatible effective profile.
- Time, subscription quota, API spend, and live-probe allowance: pending; no unlimited default.
- Corrective implementation/re-review cycles and finding dispositions: pending. A rejected false positive needs evidence; accepting a true risk or waiving a failed gate is a separate human decision.
- Public policy schema, legacy-command/configuration migration and mode-specific readiness semantics: pending.
- Archive-first final review and new corrective-Spec handling after a late finding: proposed, pending. This grant proposal neither reopens archived records nor grants a future corrective Spec by inheritance.
- Exact protected-file grant above and any final generated-file plan: pending.
- Branch naming: confirmed on 2026-09-08. Work branches use purpose types and Roundfix Run/Task branches retain their documented namespace. No naming bootstrap exception is needed; Supervisor feature-code implementation remains prohibited.

## Explicit exclusions

No grant is proposed for `.agents/skills/review/SKILL.md` or
`.agents/skills/github-pr-workflow/SKILL.md`; both are upstream-managed. There is
no grant for a release, tag, organization-wide App change, repository protection
change, credential change, paid-API purchase, or deletion of historical Run data.
The supervising session found no direct CodeRabbit reference in those upstream
skills, so retaining them does not leave an identified active integration behind.

## Settlement

The maintainer's eventual decision must replace the proposed status with an
unambiguous grant, record the date and source decision, and settle its limits.
That authorization record must be committed before the protected mutation, in
its own commit. Until then there is no Task Graph, TechSpec implementation
commitment, or runnable authorization conveyed by this file.

## Authoring check

On 2026-09-08, `rtk roundfix spec check
0126-agent-review-before-pull-request --stage prd --run-verification --format json`
was evaluated in the two-Spec command with 0127: exit 0, no findings, fourteen
stage-excluded detectors, and an empty Verification command list. This checks
the PRD stage only. The current authorization detector keys typed validation
from a date in the record filename and does not establish that this proposed
Spec-local record is granted. A green authoring check is not approval or
implementation readiness; Spec 0119 must close that recognition gap.

## Confirmed reviewer decision — 2026-09-08

The maintainer confirmed a Pre-PR Review Policy of codex, claude, coderabbit
or explicit none. ADR-0153 replaces mandatory review and total CodeRabbit
removal; ADR-0151 retains the Codex default and explicit project precedence.
For agent providers, preserve applicable model/effort and profile provenance.
None creates no reviewer/provider call and records configured omission while
QA and required checks remain mandatory. Enabled-provider failure is not none.

This policy correction does not change the current project's selected mode,
install a service, grant paid calls or implement any new configuration field.
The exact provider/profile migration and enabled adapters remain implementation
work. Historical PR-feedback records keep their original meaning.
