---
status: proposed
granted: null
created: 2026-09-08
action: Replace operational CodeRabbit review with independent local review before publication
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

# Proposed authority for independent local review

This is a reviewable proposed record, not express maintainer authorization.
`granted: null` means no protected mutation is granted by this artifact. The
Spec-local placement is proposed by Spec 0119; the current checker must not be
assumed to recognize a new authorization schema merely because this file exists.

## Authority already expressed

In the 2026-09-08 planning conversation, the maintainer requested removal of the
CodeRabbit dependency and review before opening a Pull Request, and confirmed
automatic delivery through merge after the Spec and limits are approved.
That policy requires independent review and required checks on the current
commit. It does not approve a reviewer choice, unlimited spending, a failed-gate
waiver, or the exact proposed files above. Commit and push of planning changes
are separately authorized in the conversation.

## Proposed protected scope

- Retire the repository's obsolete CodeRabbit configuration in the same bounded change that makes its removal valid; preserve unrelated configuration. Do not silently modify User Config or another repository.
- Replace the two canonical modules' operational review obligations and the Roundfix-owned skill/manifest contract. Embedded and rendered copies must come from their owners.
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

- Reviewer policy, effective models, and permitted fallback behavior: pending.
- Time, subscription quota, API spend, and live-probe allowance: pending; no unlimited default.
- Corrective implementation/re-review cycles and finding dispositions: pending. A rejected false positive needs evidence; accepting a true risk or waiving a failed gate is a separate human decision.
- Public legacy-command/configuration migration and readiness semantics: pending.
- Archive-first final review and new corrective-Spec handling after a late finding: proposed, pending. This grant proposal neither reopens archived records nor grants a future corrective Spec by inheritance.
- Exact protected-file grant above and any final generated-file plan: pending.
- The Spec 0125 branch-prefix conflict or an expressly bounded bootstrap method: pending. No non-`ma/` agent-created branch and no Supervisor feature-code implementation is approved.

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
