---
spec: 0127-durable-unattended-spec-workflow
status: active
created: 2026-09-08
surfaces: [cli, data, docs]
---

# Durable unattended Spec delivery

A Detached Run can survive its initiating session, but that alone does not
advance an approved queue through review, publication, merge, and the next Spec
after the Supervisor disappears. The maintainer wants a fully unattended
workflow with explicit limits and approvals. This PRD is in authoring: it
proposes durable delivery while preserving the existing Daemon ownership and
`implement-spec` entry point. Architecture, limits, and protected-file grants
remain pending; this document does not start or authorize implementation.

## Project Constraints

- Identifier strategy: applicable — preserve Run, Task, Spec, and Git identities and make durable delivery refer to their actual records. New work branches use purpose prefixes; tool-owned Run/Task branches use Roundfix's existing namespace. The maintainer removed the personal-prefix requirement on 2026-09-08, so naming no longer blocks Run creation. No new durable identity format is chosen by this PRD. Source: `docs/agents/domain.md`, `docs/agents/agent-instructions.md`.
- Authentication and HTTP: applicable — delivery uses the repository's existing authenticated runtimes and GitHub boundary, with no new credential storage or authentication policy. No backend HTTP guide exists; absence is not permission to create one. Publication authority and resource limits must be recorded, and a missing or denied credential parks the affected action. Source: `docs/agents/agent-instructions.md`, `docs/agents/cli.md`.
- Active ADR obligations: applicable — preserve execution ownership and evidence semantics while proposing durable delivery across Runs. Source: `docs/agents/domain.md`, `docs/agents/autonomous-work.md`, `docs/agents/spec-routing.md`.
  ADR-0014 applies: the Daemon remains responsible for Task Verification and settlement as the durable Supervisor advances delivery.
  ADR-0019 applies to the inherited Watch meaning: Clean denotes a merge-ready Open Pull Request rather than an empty local queue; the replacement readiness contract is a prerequisite.
  ADR-0020 applies to retained acpx Batch execution: a parsed prompt response remains usable despite a subsequent teardown exit, and that anomaly remains journaled before authoritative Verification.
  ADR-0038 applies: the same Agent Session gets one Verification Feedback repair; queue-level retry proposals cannot enlarge that bound.
  ADR-0056 applies: per-Run Task Capacity and Verification Capacity remain separate and cannot be interpreted as whole-queue concurrency or a machine-wide semaphore.
  ADR-0057 applies: the Daemon is the sole writer of Implement Task status, so recovery and Supervisor narration cannot manufacture a completed Task.
  ADR-0091 applies: each authored QA gate remains a terminal Task in its Spec graph.
  ADR-0096 applies: mechanical facts are checked within the Daemon-owned QA stage before spending an Agent turn, with existing verdict semantics preserved.
  ADR-0117 applies: the workflow checks authoring defects when producing their artifacts and retains commit-dependent and user-surface checks at the gate that can establish them.
  ADR-0127 applies: process residue is a readiness observation, not a synthetic Run or a reason for the inventory to settle work after restart.
  ADR-0139 applies: durable recovery must retain one Active Run per work target and the same-checkout mutation guard rather than starting a competing owner.
  ADR-0142 applies to the existing head-bound evidence contract and historical outcomes; independent pre-PR review must be settled by Spec 0126 before the new delivery path relies on it.
  Preserve the accepted decision to keep `implement-spec` and repair its conflicting instructions rather than introducing a competing implementation loop.
- Tooling authority: applicable — protected changes are proposed but express maintainer authorization for the exact bounded files has not been granted. The proposed record is [_authorization.md](_authorization.md), at `docs/specs/0127-durable-unattended-spec-workflow/_authorization.md`; its proposed status and null grant authorize no mutation. Candidate bounded files: `.agents/skills/implement-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/agents/openai.yaml`, `internal/baseline/assets/modules/autonomous-work.json`, `skills/implement-spec/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/roundfix/agents/openai.yaml`, `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, `internal/cli/cli_test.go`, `skills/baseline_skill_contract_test.go`, `internal/docscontract/publicdocs_test.go`. Deterministic digest fallout follows the sanctioned regeneration rule after source approval. No Task Graph or tooling mutation is authorized by this list. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- An approved Spec queue can reach merge without a live chat session supervising each step.
- Restart resumes from observed durable outcomes without duplicating a Run, Pull Request, or merge.
- Every automatic action stays inside the approved Spec, authority, and resource limits.
- The maintainer can see progress and the exact reason an item stopped, including missing authority or stale evidence.
- `implement-spec`, the canonical autonomous guide, and the Roundfix skill describe one coherent delegation workflow.

## User Stories

1. As a maintainer, I want to approve a queue and its limits once, so that the workflow can deliver each eligible Spec without asking me to approve routine steps again.
2. As a maintainer, I want delivery to recover after the supervising session exits, so that an implemented Spec does not remain stranded between QA, review, and publication.
3. As a maintainer, I want unavailable authority, exhausted limits, or failing evidence to stop the affected action visibly, so that unattended execution cannot manufacture approval.
4. As a Supervisor, I want later Specs rechecked after earlier dependencies land, so that the queue does not execute an obsolete assumption.
5. As a maintainer returning to the repository, I want durable status and receipts, so that I can distinguish an observed merge from an agent's claim.

## Core Features

1. Queue preparation inventories active Specs, backlog intent, Findings, and the destination Inbox under their existing lifecycle rules. It records dependencies and pending decisions and identifies the Specs actually approved for execution. An inventory alone is not implementation authority or source adoption.
2. Durable orchestration owns progress across Spec Runs and delivery actions. The workflow survives loss of the initiating chat and restarts from observed results, including an action that succeeded externally before its local acknowledgement was recorded. It never opens duplicate PRs or merges twice because an acknowledgement was lost.
3. Proposed order is the Daemon-owned Task Graph and terminal QA, archive and commit, independent review of the resulting candidate, publication, current-head checks, and merge. Spec 0126 supplies the independent-review contract. Archive invalidates a pre-archive review, and no failed, skipped, stale, or absent evidence is interpreted as approval.
4. Approval and limits accompany each consuming Spec. A granted routine action need not be reconfirmed; a missing, expired, changed, or insufficient grant stops its dependent action. Only one user question may be pending, and the workflow never treats silence, timeout, or a recommended option as the answer.
5. Run Window and Run Budget retain their existing separate meanings: the window bounds new Run starts, and the budget bounds an individual Run. Queue-wide time, spending, concurrency, and correction behavior must be expressly decided and enforced rather than inferred from those settings.
6. Later Specs are revalidated when their prerequisites change the operative assumptions. A contradiction, new authority requirement, or new irreversible action becomes a durable blocker instead of an improvised implementation decision.
7. Progress records identify actual Runs, candidate commits, review evidence, publication state, and merge outcomes. Recovery reconciles each externally visible action before retrying it. Unresolved findings, exceeded limits, and required maintainer decisions remain visible after restart.
8. The owned `implement-spec` entry point delegates implementation and Verification to Roundfix. It absorbs useful queue preparation from kickoff without introducing another implementation loop, copying a second window script, or authorizing the Supervisor to write code or tests.
9. A blocking review after archive preserves the completed/archived Spec and parks publication. Correction is authored as a new Spec with its own authority and gate; the candidate is reviewed again after that correction is archived. A remaining cycle budget or routine queue-continuation grant does not silently authorize the new Spec.

## User Experience

The maintainer sees the approved queue, its active Spec, current Run and delivery
stage, remaining limits, and any parked decision. After an interruption, the
workflow reports what it recovered and what still requires action. A returning
maintainer can inspect the evidence of publication and merge independently of
the original conversation. When no durable owner is running, status says that
execution stopped instead of promising that work continues.

## Non-Goals / Out of Scope

- Adopting Fluxus kickoff as a second canonical skill or replacing `implement-spec`.
- Letting the Supervisor implement features or tests directly.
- Bypassing QA, review, checks, tooling grants, or the confirmed purpose-based work-branch policy.
- Reopening archived Specs or inheriting review validity across an archive commit.
- Treating a Task Worktree concurrency setting as permission for parallel whole-Spec deliveries.
- Selecting a reviewer or paying for API access without the pending decisions.
- Replacing the existing Run lifecycle with an unrelated orchestration framework.
- Releasing versions, creating tags, or performing destructive migrations under routine delivery authority.

## Success Metrics

| Acceptance observation | Evidence required |
| --- | --- |
| Delivery survives the chat/session boundary | A Spec continues under an identified durable owner after the initiating session exits. |
| The queue survives its own owner's restart | Restart between completed Runs resumes the next eligible delivery step without an additional user prompt for already granted actions. |
| External action replay is safe | Interruptions around push, PR creation, and merge reconcile the observed result without duplicates or lost ownership. |
| Conditional authority remains bounded | Missing approval, exhausted approved limits, stale review, or failed checks prevents the dependent action and remains visible after restart. |
| A late finding cannot mutate completed history | Post-archive review blocks publication, preserves the original Spec, and waits for any missing corrective-Spec authority. |
| Canonical ownership is consistent | Owned instructions consistently assign code/tests and Verification to the runtime and Daemon, with QA retained in the Task Graph. |
| External acceptance evidence | A real repository outside this Spec's fixtures exercises interruption and resumption, with its origin and actual external receipts recorded. Missing evidence is reported under the repository's declared policy. |

## Decisions

- The maintainer confirmed automatic delivery through merge after the Specs and limits are approved, requiring independent review and required checks on the current commit.
- The existing decision keeps `implement-spec` as the Roundfix-owned loop; this proposal repairs the implementation of that decision and evaluates useful kickoff preparation.
- A Detached Run surviving its caller is insufficient proof that a multi-Spec delivery queue survives its owner.
- Proposed archive-first final review preserves the existing archive boundary but can require a new corrective Spec, gate, and grant after a late finding. That trade-off is pending, not an accepted change to QA or archive policy.
- Spec 0119 supplies the proposed authorization convention; Specs 0122, 0125, and 0126 supply settlement, branch identity, and independent-review prerequisites. Runtime readiness and capacity evidence are provided by the corresponding portfolio work.
- This PRD makes no source adoption or executable commitment while the relevant decisions and grants remain pending.

## Open Questions

- Durable owner, persistence/recovery contract, and public command surface — architecture decision pending; no proposed mechanism is represented as implemented.
- Queue-wide time, API/subscription spend, concurrency, and correction limits — maintainer decision pending. Serial whole-Spec delivery and two corrective review cycles are proposals, not defaults already authorized.
- Cutoff behavior for a Run already started, failed delivery after the cutoff, and explicit cancellation — maintainer decision pending; existing Run Window semantics remain unchanged until approved otherwise.
- Independent reviewer policy and its candidate-evidence contract — pending decision and Spec 0126 prerequisite.
- Archive-first final review and the handling/authority of post-archive corrective Specs — proposed sequencing pending confirmation; no artifact-only review exemption or inherited grant is assumed.
- Exact protected-file grants and generated outputs — pending in the proposed authorization record.
- Branch naming is settled and requires no bootstrap exception. Spec 0125 remains relevant to shared repository identity and reconciliation, not permission to start a Run under its existing namespace.

## Research and limitations

Secondbrain consultation started at
`/Users/marcio/dev/secondbrain/wiki/index.md` and used qmd before reading
`/Users/marcio/dev/secondbrain/wiki/concepts/agent-harnesses.md` and
`/Users/marcio/dev/secondbrain/wiki/concepts/agent-workflows-e-loop-engineering.md`.
Their treatment of durable ownership, independent checks, receipts, and stop
conditions informed the proposed outcome and challenged reliance on chat
continuation. These syntheses do not prove a particular recovery architecture.

Exa located and read [OpenAI code review](https://learn.chatgpt.com/docs/code-review)
and [Claude headless execution](https://code.claude.com/docs/en/headless).
They support the native reviewer prerequisite, not automatic end-to-end delivery.
Current local inspection confirmed Run Window and Detached Run primitives and
the instruction conflict. No queue owner, live recovery exercise, reviewer, or
GitHub mutation was executed for this PRD.
