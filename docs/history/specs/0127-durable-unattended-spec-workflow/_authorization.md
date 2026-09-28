---
status: proposed
granted: null
created: 2026-09-08
action: Add durable unattended delivery for approved Specs and align the owned orchestration instructions
consuming:
  - 0127-durable-unattended-spec-workflow
paths:
  - .agents/skills/implement-spec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/agents/openai.yaml
  - internal/baseline/assets/modules/autonomous-work.json
  - skills/implement-spec/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/roundfix/agents/openai.yaml
  - docs/agents/autonomous-work.md
  - docs/agents/setup-context.json
  - internal/cli/cli_test.go
  - skills/baseline_skill_contract_test.go
  - internal/docscontract/publicdocs_test.go
---

# Proposed authority for durable unattended delivery

This proposed record grants no protected mutation. Its null grant and proposed
status are deliberate: the final architecture, limits, and file boundary await
the maintainer's decision. Spec 0119 proposes this file's future canonical role;
its existence alone does not extend the current checker or execution authority.

## Authority already expressed

In the 2026-09-08 planning conversation, the maintainer confirmed automatic
implementation, the configured review policy, PR creation, and merge after the consuming
Specs and their limits are approved. The configured review policy must be satisfied and required checks must
pass for the current commit; explicit none records intentional review omission. This policy authorizes neither an unapproved queue
nor an unanswered budget; it does not grant a gate waiver. Commit and push of
planning changes are separately authorized in the conversation.

## Proposed protected scope

- Align the owned orchestration skills and canonical autonomous module with Daemon-owned implementation and Verification, terminal QA, the configured review/omission policy, and durable delivery.
- Regenerate only their listed embedded copies and rendered guidance through their existing owners.
- Intentionally revise the named governed contract tests to prove the approved workflow; preserve checks that reject missing authority, stale evidence, and unsafe replay.
- Keep `implement-spec` as the entry point. No kickoff skill, extra window script, package-manager change, build-tool change, or new dependency is included in this proposal.
- Record any additional governed file discovered during technical design as an amended exact proposal before asking for its express grant.

## Sanctioned regeneration

After the source scope is expressly granted, existing `make skills-sync` and
`make baseline-digests` may regenerate the listed skill copies and deterministic
digest fallout under the repository's sanctioned rule. Use the public Baseline
update plan for rendered guides. No hand-edited digest, unrelated fixture edit,
or undocumented generated path is authorized by this proposal.

## Pending decisions and execution limits

- Durable supervision, persistence, action reconciliation, and migration behavior: pending.
- Review policy: confirmed — codex, claude, coderabbit or explicit none; default Codex with explicit project precedence. The implemented policy/adapter prerequisite in Spec 0126 remains pending.
- Queue-wide time, API spend, subscription quota, parallel delivery, and corrective-cycle limits: pending. No missing value means unlimited.
- Cutoff and cancellation handling for delivery already in progress: pending.
- Archive-first final review and authority for a new corrective Spec after a late finding: pending. Existing archive immutability and exact-head evidence remain binding; this proposal grants neither an exception nor automatic approval inheritance.
- The exact protected-file list above and any final generated-file plan: pending.
- Branch naming: confirmed on 2026-09-08. Work branches use purpose types and Run/Task branches retain the Roundfix namespace. No naming bootstrap exception or Supervisor feature-code exception is needed or granted.
- Readiness and settlement prerequisites named by the portfolio: unresolved work remains a prerequisite, never waived by this record.

## Explicit exclusions

No scope is granted for upstream-managed `.agents/skills/review/SKILL.md` or
`.agents/skills/github-pr-workflow/SKILL.md`. No scope is granted for a release,
tag, credential or GitHub protection change, API-credit purchase, destructive
data migration, or deletion of historical Run evidence. A rejected false finding
requires proof; acceptance of a real risk or waiver of a failed gate requires a
separate human decision and is not automatic orchestration behavior.
The supervising session found no direct CodeRabbit reference in those upstream
skills; leaving them unchanged preserves ownership without retaining an
identified active provider dependency.

## Settlement

Record the maintainer's eventual affirmative decision, its date/source, the
bounded paths, and concrete limits before changing this proposal to a grant.
Commit that authorization separately before the governed implementation.
Until then this file supports review of the proposal only: it conveys no Task
Graph, runnable grant, or permission to bypass a blocked prerequisite.

## Authoring check

On 2026-09-08, `rtk roundfix spec check
0127-durable-unattended-spec-workflow --stage prd --run-verification --format json`
was evaluated in the two-Spec command with 0126: exit 0, no findings, fourteen
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


## QA archive override is a separate action

The maintainer permits the capability described in ADR-0154 and Spec 0122.
The queue may consume an already applicable explicit per-Spec archive override
without reconfirmation, preserving approval and QA evidence. This proposed
queue grant supplies no override for a particular Spec and grants no QA success,
Run Clean, publication or merge exception through archival alone.
