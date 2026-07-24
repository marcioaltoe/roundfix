---
spec: 0046-public-context-driven-baseline-command
status: active
created: 2026-07-23
updated: 2026-07-23
owner: Roundfix maintainers
surfaces: [cli, infra, docs]
---

# Public Context-Driven Baseline Command

The Context-Driven Baseline can currently be configured only through a Python
implementation distributed inside the `setup-context-driven` skill. The live
Fluxus adoption recorded in
[the setup adoption finding](../../findings/2026-07-23-setup-context-driven-adoption-process-improvements.md)
completed safely, but required repository cleanup before source inventory,
repeated large audit responses, a hand-built Decision File, manual HTTP
evidence, an unexplained request for a new PostgreSQL contract document,
file-level plan reconstruction, and correction of persisted verification
commands. Roundfix needs one public Baseline Command that humans and
automations can use without Codex, preserves the safety contracts delivered by
Spec 0045, and makes the observed adoption friction part of the product
contract.

## Goals

- Humans and automations can adopt and update a Context-Driven Baseline through
  one public Roundfix command without requiring Codex, Python, or an installed
  Agent Skill.
- Interactive and non-interactive callers produce the same deterministic
  Decision Plan and Change Plan from equivalent repository state and answers.
- Every correction identified by the Fluxus adoption finding is observable in
  the Baseline Command contract and covered by final QA.
- The `setup-context-driven` skill becomes a thin instructional layer over the
  public CLI, with no independent setup engine or fallback runtime.
- Initial adoption and later updates use one idempotent, confirmation-gated
  workflow without weakening any maintained Python safety contract.
- Roundfix documentation teaches the human, automation, migration, recovery,
  and troubleshooting flows delivered by the Baseline Command.

## User Stories

1. As a repository maintainer, I want `roundfix baseline` to guide adoption or
   update from preflight through verification, so that I do not have to
   coordinate separate setup scripts and Decision Files.
2. As an automation author, I want a non-interactive interface with stable
   machine-readable output, so that I can operate the same baseline workflow
   without Codex or a terminal.
3. As a repository maintainer, I want to preserve selected repository-specific
   instructions or explicitly start greenfield, so that adoption never silently
   discards policy.
4. As a repository maintainer, I want to choose exactly one built-in or
   repository-owned Baseline Profile, so that every contributor and automation
   resolves the same workflow contract.
5. As a maintainer reviewing existing instructions, I want a consolidated,
   editable classification proposal, so that I can approve retained rules
   without reviewing every source entry in a separate prompt.
6. As a maintainer rejecting a plan, I want to select the decision to revisit
   and optionally explain the correction, so that Roundfix can recalculate the
   plan without restarting adoption.
7. As a maintainer confirming a Change Plan, I want a file-level projection
   alongside the complete managed-entry ledger, so that I can understand the
   actual repository mutation without losing ownership evidence.
8. As a maintainer reviewing Repository Capabilities and Verification, I want
   diagnostics to distinguish observed implementation evidence, required
   repository contracts, portable verification roles, and executable
   repository commands, so that the Setup Manifest never claims evidence or
   commands the repository does not have.
9. As an Agent using the setup skill, I want the skill to call the same public
   CLI used by humans and automations, so that its guidance cannot drift from
   executable behavior.
10. As a Roundfix user, I want current command and recovery documentation, so
    that I can complete Baseline adoption without reading implementation
    artifacts.

## Core Features

1. The Baseline Command is the public authority for Baseline Profile
   authoring, validation, audit, decisions, planning, and confirmation-gated
   apply. Its human-readable and machine-readable interfaces expose the same
   behavior without requiring an installed Agent Skill.
2. One `roundfix baseline` invocation starts a deterministic state machine:
   preflight, existing-state detection, instruction-preservation choice,
   profile selection, repository audit, divergence resolution, final plan
   review, confirmation-gated apply, and Baseline verification. Update and
   adoption cannot bypass any applicable state.
3. Preflight detects whether the repository already has a valid Baseline
   Profile. A configured repository enters update with its current profile and
   offers an explicit profile-change action; an unconfigured or incompatible
   repository enters adoption.
4. Roundfix ships built-in Baseline Profiles and supports custom Baseline
   Profiles owned only by the current repository. A maintainer can create,
   configure, inspect, and validate a custom profile before it can participate
   in a Decision Plan. Profile selection presents multiple candidates but
   resolves exactly one profile; profiles are never composed.
5. Adoption offers two explicit instruction-preservation modes. Greenfield
   creates immutable content-addressed backups of root instruction carriers
   and imports none of their rules. Preservation creates the same backups and
   requires human approval of every rule moved into Repository-Specific
   Normative Rules or rejected as conflicting.
6. Audit inventories every bounded instruction and agent-document carrier so
   the Change Plan retains complete source evidence. Automatic backup,
   classification, and mutation of pre-existing instruction carriers are
   limited to root carriers. Nested carriers remain untouched; detected
   conflicts appear as warnings in the plan and final report.
7. Safe root instruction aliases report their in-repository target and content
   identity and preserve the target bytes once. External, escaping, cyclic,
   unreadable, or special-file targets remain opaque and block apply.
8. When root instruction meaning is ambiguous, Roundfix requests a read-only
   ACP analysis from the proven preferred Agent Selection. An invalid,
   incomplete, or unavailable result is discarded before the same immutable
   snapshot is retried with the proven fallback selection. If neither
   selection succeeds, the flow names the backup and Repository-Specific
   Normative Rules destinations and continues through manual classification.
9. ACP analysis produces proposals only. Roundfix presents one consolidated,
   editable review, validates every accepted disposition through the
   deterministic Baseline engine, and writes nothing until the resulting
   Change Plan receives explicit human confirmation.
10. After profile selection, audit evaluates repository evidence and alignment
    with the selected profile. Required divergence blocks planning until the
    maintainer confirms an allowed decision or changes the profile; advisory
    divergence remains visible without being inferred as repository policy.
11. A rejected final plan returns to a structured decision area such as
    profile, repository rules, divergences, or files. The maintainer may add a
    free-form suggestion; ACP analysis can translate it into an in-scope
    proposal, but deterministic validation, plan recalculation, and renewed
    human approval remain mandatory.
12. The Change Plan exposes a derived file-level projection with one record
    per repository path, its file action, before and after identity, and ordered
    managed entries. The complete managed-entry ledger remains canonical and
    contributes to the confirmed plan digest.
13. Apply remains fail-closed, checks the complete repository preimage,
    accepts only the exact current plan digest, verifies every postimage, and
    rolls back incomplete mutation. Empty reapply remains idempotent.
14. Baseline verification checks only Baseline-owned outcomes, including
    carrier relationships, backups, profile and manifest identity, managed
    artifacts, retention accounting, and resolved audit state. Repository
    formatter and Verification commands remain recommendations that the
    Baseline Command never executes.
15. Non-interactive use follows the same state machine without prompts. One
    execution emits a complete machine-readable Change Plan; a second
    execution applies only that exact approved digest. Missing input or stale
    state produces a structured next action and a non-zero exit.
16. Every automatable result provides stable machine-readable output,
    structured failures, schema identity, deterministic ordering, and
    documented exit categories. Requested output goes to stdout; diagnostics
    and progress go to stderr.
17. The deterministic audit engine remains local, read-only, and network-free.
    It reports bounded HTTP route candidates without inferring contract policy,
    distinguishes PostgreSQL implementation evidence from a missing repository
    contract, and validates repository-executable commands before labeling them
    executable.
18. Decision input examples, command help, emitted templates, skill guidance,
    and Roundfix documentation describe the same contract accepted by the
    runtime. Documentation examples must pass the public parser or an
    equivalent documentation contract test.
19. All Spec 0045 safety properties remain binding: byte-exhaustive Baseline
    Readoption, individual dispositions, Upgrade Retention Contract,
    Repository-Specific Normative Rules ownership, formatter-stable managed
    output, declarative setup ownership, exact plan-digest confirmation, and
    atomic rollback.
20. Cutover requires a parity matrix that assigns every maintained Python
    premise, behavior, action, refusal, and supported adoption state to a Go
    destination. User-facing audit and apply become the single Baseline flow;
    maintainer-only asset synchronization or skill-restoration operations may
    remain separate Go operations. No maintained capability can disappear
    silently.
21. The setup skill contains guidance and public CLI recipes only after Go
    parity. It has no independent setup engine, Python fallback, or divergent
    behavior.
22. Roundfix user documentation covers first adoption, update, greenfield,
    instruction preservation, profile changes, plan revision, automation,
    recovery, diagnostics, and the migration from the Python-backed skill.

## User Experience

A maintainer runs `roundfix baseline`. Preflight either identifies the current
Baseline Profile and enters update or begins adoption. Adoption asks whether
root instructions must be preserved. Greenfield archives them without
importing rules; preservation prepares a consolidated classification review.

Roundfix then presents the valid Baseline Profiles and records exactly one
selection. Audit compares repository evidence with that profile, explains
blocking and advisory divergence, and collects only decisions that the
repository owner must make. It never infers HTTP policy, exception ownership,
or other Normative Clauses from implementation evidence.

When instruction analysis needs semantic judgment, ACP produces a proposal
from an immutable snapshot. The maintainer can edit the consolidated proposal
before accepting it. If ACP is unavailable, the same review remains available
for manual classification and names the backup and repository-rules
destinations.

The final review starts with file-level changes and retains the complete
managed-entry and retention ledgers for detailed or machine inspection. A
rejected plan returns to the selected decision area; an optional suggestion can
be analyzed only within Baseline scope. Every accepted revision produces a new
plan digest and requires a new approval.

After approval, Roundfix applies the exact plan and verifies the Baseline
state. It reports repository formatter and Verification commands as
recommendations without running them.

Automation uses the same command family without a TTY. `roundfix baseline
plan` emits the complete JSON plan, and `roundfix baseline apply` receives
that portable plan and its approved digest. Neither subcommand prompts,
substitutes a newer plan, or applies stale input.

## Non-Goals / Out of Scope

- User-scoped custom Baseline Profiles or precedence between repository and
  user profile catalogs.
- Combining multiple Baseline Profiles in one repository state.
- An external profile registry, plugins, or remote profile discovery.
- Keeping the Python setup implementation as a supported fallback after Go
  parity and cutover.
- Changing the existing Setup Command, which continues to prepare a machine
  for Roundfix Runs.
- A separate baseline update path that bypasses preflight, audit, decisions,
  plan confirmation, or Baseline verification.
- Automatically editing, backing up, or removing nested instruction carriers.
- Allowing ACP output to authorize, mutate, or remove repository content.
- Inferring repository policy, Normative Clauses, HTTP mode, exception owner,
  or exception rationale from source code.
- Following or applying through unsafe symlinks.
- Installing application dependencies, provisioning PostgreSQL, connecting to
  live infrastructure, or executing repository scripts during audit, apply, or
  Baseline verification.
- Making repository-authored policy or Repository-Specific Normative Rules
  setup-owned.
- Weakening the Context-Driven Baseline contracts accepted in ADRs 0046, 0047,
  and 0058 through 0064.

## Success Metrics

- 100% of the seven issue categories in the linked Fluxus adoption finding map
  to an observable Core Feature and passing QA evidence.
- 100% of maintained Python premises, behaviors, actions, and supported states
  have an explicit Go destination and parity evidence before Python removal.
- 100% of maintained Baseline Profiles and supported adoption states produce
  equivalent normalized decisions, planned file bytes, plan digests, refusal
  categories, and rollback outcomes before cutover.
- Interactive and non-interactive runs produce identical plan digests for
  100% of equivalent repository snapshots and answers.
- 100% of modified root instruction carriers have an exact immutable
  content-addressed backup included in the confirmed Change Plan.
- 0 nested instruction carriers are mutated by automatic instruction
  preservation; every detected nested conflict appears in the plan and final
  report.
- 0 ACP proposals mutate repository state or bypass consolidated human review,
  deterministic validation, and plan confirmation.
- 0 unsafe instruction targets are followed or read as trusted source.
- 100% of Change Plans expose both one file-level record per affected path and
  the complete managed-entry ledger.
- 0 Setup Manifests label an unvalidated profile expectation as a
  repository-executable command.
- Audit, apply, and Baseline verification perform 0 dependency installations,
  database connections, formatter runs, or repository Verification runs.
- Generated output composes with every maintained profile's formatter,
  repository Verification, final audit, and empty reapply with 0 managed-file
  delta in QA.
- One separately authorized live Fluxus adoption or update completes without
  manual Decision File schema repair, manual route digest calculation,
  ambiguous capability diagnosis, or correction of persisted verification
  commands.
- The shipped setup skill contains 0 executable setup-engine scripts and uses
  only the public Baseline Command.
- 100% of documented Baseline command examples pass their documentation
  contract checks.
- The full repository Verification passes before the Spec can complete.

## Decisions

- The Roundfix CLI is the sole Context-Driven Baseline execution authority;
  the setup skill is a thin instructional layer. See ADR-0066.
- Custom Baseline Profiles are repository-owned and versioned with the
  repository; no user-scoped custom catalog participates. See ADR-0067.
- Human adoption and update use one `roundfix baseline` state machine, while
  automation uses `roundfix baseline plan` and `roundfix baseline apply` over
  the same engine and exact digest. See ADR-0068.
- Ambiguous instruction classification and in-scope plan suggestions use
  read-only ACP proposals with an exact preferred and fallback Agent Selection;
  deterministic validation and human approval remain authoritative. See
  ADR-0069.
- Baseline adoption audits all bounded carriers but automatically preserves and
  changes pre-existing instruction carriers only at the repository root. See
  ADR-0070.
- A greenfield choice creates backups but imports no existing root rule.
- The maintainer selects exactly one Baseline Profile.
- Repository-Specific Normative Rules use the single canonical
  `docs/agents/specific-repository.md` carrier. Baseline creates and links it
  only for non-empty rules and safely migrates either legacy carrier name.
- Baseline verification never executes repository formatter or Verification
  commands; it reports them as recommendations.
- Completion requires full Go parity, a separately authorized live Fluxus
  journey, skill cutover, documentation, and removal of the Python
  implementation.
- Every correction in the linked adoption finding is in scope for this Spec.

## Delivery note

The future Task Graph must contain a dedicated documentation Task. That Task
must update Roundfix's user guide, CLI reference and examples, automation
contract, migration guidance, recovery and troubleshooting guidance, and the
thin setup skill after the public command behavior is stable.

## Open Questions

None.
