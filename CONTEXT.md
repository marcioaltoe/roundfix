# Roundfix

Roundfix picks up Work Items — Review Issues from pull request reviews today, Tasks from Spec Task Graphs next — resolves them through the user's selected ACP Runtime, and pushes only when nothing unresolved remains. This glossary defines the product language for that loop.

## Language

**Run**:
A durable attempt to drive one target's Work Items — an Open Pull Request's Review Issues or a Spec's Tasks — to a terminal outcome. One Active Run is allowed per target: (Head Repository, PR Head Branch) for review work, (repository, spec slug) for spec work.
_Avoid_: Session, execution, job

**Fetch Run**:
A short review Run that runs Branch Integrity Preflight, fetches Review Source issues, and persists markdown artifacts from the user's checkout without starting an Agent, creating a Run Worktree, committing, or pushing.
_Avoid_: Standalone fetch, untracked fetch

**Open Pull Request**:
A pull request that GitHub reports as open and still eligible for review-resolution work.
_Avoid_: Closed pull request, merged pull request

**Active Run**:
A Run that has started and has not reached a terminal outcome.
_Avoid_: Open run, live run

**Stop Request**:
A user's explicit request to end an Active Run before it reaches another terminal outcome.
_Avoid_: Pause, retry, failure

**Force Stop**:
The Stop Command mode that proves owner identity, cancels registered Agent Sessions, terminates the recorded owning process, and completes the Run as Stopped only after owner exit is proven. A proven identity mismatch always refuses; only the explicit `--owner-identity-unreadable` last-resort flag permits PID-only termination after the host reports the identity unreadable.
_Avoid_: Lock release, best-effort stop, orphan reclamation

**PR Head Branch**:
The branch on GitHub that supplies the pull request's head commits.
_Avoid_: Local branch, checkout branch

**Head Repository**:
The GitHub repository that owns the PR Head Branch.
_Avoid_: Base repository, local checkout

**Work Item**:
One unit of resolvable work that Roundfix picks up and drives to an outcome: a Review Issue today, a Task next.
_Avoid_: Ticket, job, to-do

**Task Source**:
The origin Roundfix reads Work Items from, such as a Review Source or a Spec's Task Graph.
_Avoid_: Provider, backend, integration

**Review Source**:
The external review system that produces feedback for an Open Pull Request.
_Avoid_: Review Provider, Agent, ACP Runtime

**Review Source Evidence**:
A head-bound Review Source classification for one expected commit: `pending` has no usable expected-head signal, `reviewing` is still in progress, and `reviewed` is complete without proving Merge-Ready. `verified` proves the expected head with no unresolved Review Issues, `skipped` explicitly declines that head, and `failed` records an explicit Review Source failure.
_Avoid_: Check presence, generic approval, inferred clean state

**Spec**:
One Spec's planning artifact set produced by the spec workflow: PRD, Task Graph, Task files, QA evidence, and adopted sources under `references/` with provenance recorded in `references/_index.md`.
_Avoid_: Feature folder, epic, project

**Spec Root**:
The configured directory holding Spec folders, `docs/specs` by default; it may resolve outside the repository working tree, such as a knowledge workspace repository.
_Avoid_: Specs directory, docs folder, knowledge base

**Task**:
One implementable unit of work within a Spec. Its task file is the sole owner of its status.
_Avoid_: Subtask, story, ticket

**Task Type**:
The required classification that routes a Task to its Agent Selection Profile. The valid values are `backend`, `frontend`, `data`, `infra`, `docs`, `test`, and `chore`; use the dominant implementation surface when a Task crosses more than one.
_Avoid_: Task category, work type, inferred type

**Task Graph**:
The Spec's manifest that declares its Tasks and their dependencies as a directed acyclic graph. Dependencies live only in the Task Graph, never in task files.
_Avoid_: Task list, backlog, roadmap

**Spec Context Bundle**:
The bounded Task-start context that gives an Agent the full assigned Task and paths to larger Spec artifacts, relevant interfaces, and files changed by prior Tasks.
_Avoid_: Repository dump, full Spec payload, cold-start exploration

**Project Constraint**:
A confirmed project decision or universal Normative Clause that a Spec records as applicable, not applicable with a reason, or explicitly authorized for a bounded change before its Tasks can execute.
_Avoid_: Optional note, implementation suggestion, inferred permission

**Tooling Authority**:
The universal Normative Clause that forbids changes to linter, formatter, and tool configuration unless the maintainer explicitly authorizes the exact bounded files.
_Avoid_: Tool preference, implicit permission, cleanup authorization

**QA Report**:
The qa-gate evidence report written to a Spec's QA directory, carrying a machine-readable verdict plus `rows_blocked_environment` and `rows_blocked_finding` counts in its frontmatter.
_Avoid_: Test report, QA log

**ACP Runtime**:
A local coding runtime that Roundfix launches through the user's installed tool and authentication setup using Agent Client Protocol stdio. The MVP supports Codex through `codex-acp`, Claude through `claude-agent-acp`, and OpenCode through `opencode acp`; command overrides remain a stdio escape hatch for local testing.
_Avoid_: Review Source, review provider

**Agent Model**:
A runtime-specific model choice that Roundfix explicitly assigns to every Agent Session.
_Avoid_: ACP Runtime, Agent, free-form model override

**Agent Work Category**:
The routing key Roundfix uses to resolve an Agent Selection Profile: `general`, one Task Type, `qa`, or `review`.
_Avoid_: Agent role, benchmark category, automatic route

**Agent Selection**:
One exact ACP Runtime, Agent Model, and reasoning-effort tuple that Roundfix can prove and assign to an Agent Session.
_Avoid_: Model name, runtime default, partial selection

**Adapter Readiness**:
Proof that the effective ACP adapter command has the supported package lineage and version required by Roundfix. Executable presence or a matching command name alone is not proof.
_Avoid_: Adapter found, binary check, PATH readiness

**Exact Agent Selection Proof**:
The token-free disposable Agent Session check that maps one Agent Selection through advertised ACP capabilities, applies its exact model and reasoning assignment, observes matching effective state, and closes the Session successfully.
_Avoid_: Model validity, catalog match, recommendation rank

**Agent Selection Profile**:
The atomic policy for one Agent Work Category, containing one Preferred Selection and a non-empty ordered Fallback Chain. A higher-precedence profile replaces the complete lower-precedence profile rather than merging individual fields.
_Avoid_: Runtime defaults, model preset, partial override

**Agent Selection Profile Readiness**:
The command-scoped result that resolves effective Agent Selection Profiles, deduplicates their Preferred Selections and Fallback Chains, and requires Exact Agent Selection Proof for every distinct tuple before a Run or configuration mutation.
_Avoid_: Single-model check, configured Agent probe, cached readiness

**Preferred Selection**:
The first Agent Selection Roundfix proves and attempts for an Agent Work Category.
_Avoid_: Default model, primary runtime, recommendation winner

**Fallback Chain**:
The non-empty ordered Agent Selection list Roundfix proves with the Preferred Selection before a Run and may activate after notifying the user and Supervisor that the preceding selection failed before Agent work began.
_Avoid_: Dynamic fallback, silent retry, unproven selection

**Default Agent Model**:
The concrete Agent Model Roundfix selects for an ACP Runtime when the user supplies no override. It never inherits the runtime's local model configuration.
_Avoid_: Runtime default, automatic model, local Agent default

**Default Reasoning Effort**:
The runtime-specific reasoning level Roundfix assigns when the value is non-empty. An empty value is valid and means the Agent Model manages reasoning, so Roundfix assigns no reasoning option. It never inherits the runtime's local reasoning configuration.
_Avoid_: Local Agent reasoning, automatic reasoning, reasoning hint

**Model Catalog**:
The ordered set of known Agent Models Roundfix offers for one ACP Runtime during Interactive Input. Its Default label resolves to the Default Agent Model, while non-interactive interfaces may supply a custom value.
_Avoid_: Global model list, model allowlist

**Model Recommendation Ranking**:
The versioned, advisory top-five Agent Selection list Roundfix shows for an Agent Work Category to help configure a profile. It never selects, routes, or changes an Agent Selection automatically.
_Avoid_: Model router, benchmark policy, automatic selection

**Fallback Selection**:
The next configured Agent Selection in a profile's Fallback Chain. Roundfix proves it before the Run, emits a notification before activation, and may switch ACP Runtime automatically only while Agent work has not begun.
_Avoid_: Dynamic fallback, silent model switch, catalog probe winner

**Agent Session**:
The acpx-backed session owned by one Work Item or action. Each Implement Task owns a Task Type-selected Agent Session, requested QA owns a separate `qa` Agent Session, and review work uses a review-selected Agent Session; effective selection and fallback attempts are persisted for that owner.
_Avoid_: ACP session, chat, conversation, thread

**Merge-Ready**:
The state where accepted Review Source Evidence verifies the pushed head commit, or its proven Roundfix artifact-only descendant, with no new Review Issues, letting a watch Run end Clean.
_Avoid_: mergeable, green check, approved

**Wave**:
The set of Tasks whose dependencies are all completed and that may execute concurrently; the scheduler draws from the current Wave up to the configured concurrency.
_Avoid_: Batch, stage, phase

**Task Capacity**:
The maximum number of Task Worktree lifecycles one Implement Run may execute concurrently, configured by `worktree.concurrency`; it limits Agent work and Task Worktree ownership independently from Verification Capacity.
_Avoid_: Verification concurrency, worker count, Agent pool

**Verification Capacity**:
The maximum number of Task Verification attempts one Implement Run may execute concurrently, configured by `verification.concurrency` and defaulting to `1`; an exclusive retry consumes the Run's entire capacity.
_Avoid_: Task Capacity, worktree concurrency, machine-wide test lock

**Task Worktree**:
The ephemeral git worktree one concurrently executing Task runs in — created from the Run Branch tip at Task start and integrated back onto the Run Branch at settlement; kept only when its Task fails.
_Avoid_: Run Worktree, sandbox, scratch dir

**Detached Run**:
A Run started with the detach flag: roundfix re-executes itself as a session leader independent of the caller, reports the Run ID, is followed by humans through Attach or by Supervisors through the Run Event Stream, and is ended through the Stop Command.
_Avoid_: Background job, nohup run, daemon

**Run Worktree**:
The isolated git worktree a spec Run executes in — created on the Run Branch at Run start, recorded on the Run, removed after a Clean integrated spec outcome, and kept as the inspection and settle surface otherwise. Review Runs do not create Run Worktrees.
_Avoid_: Sandbox, scratch dir, user checkout

**Run Branch**:
The named branch (`roundfix/run-<id>`) that carries a spec Run's commits inside its Run Worktree until integration moves them to the user's branch. Review Runs use the user's checkout branch directly; older pending Run Branch work is handled by Branch Integrity Preflight before review work starts.
_Avoid_: Temp branch, detached HEAD, feature branch

**Run Worktree Reconciliation**:
The proof-based classification of a terminal spec Run's retained Git surfaces: `safe` when the Run Branch and recorded target resolve, any present Run Worktree is registered and clean, and the Run Branch tip is an ancestor of the target tip; `superseded` when a QA-report-only Run Branch is older than the target branch's QA Report for the same Spec; `unintegrated` when the same evidence resolves but ancestry is false; `dirty` when a present Run Worktree has tracked or untracked changes; `unknown` when metadata or Git evidence cannot prove another state; and `released` only when both the Run Worktree and Run Branch are absent. Only freshly revalidated `safe` or `superseded` work can be cleaned up.
_Avoid_: GC, force cleanup, manual branch deletion

**Integration Pending**:
A terminal spec Run outcome where the spec Run's work completed but its commits could not be fast-forwarded onto the user's branch (local changes overlap or the branch diverged); the commits stay on the Run Branch and the report names the integration command. Review Runs do not end Integration Pending.
_Avoid_: Failure, conflict, silent divergence

**Max Rounds**:
The configured number of Review Source rounds after which a Run is considered sufficiently reviewed for the developer's final merge, squash, or rebase decision.
_Avoid_: Budget, timeout, token cap

**Max Rounds Reached**:
A terminal Run outcome where the configured review round policy is complete, even if unresolved Review Issues remain for developer judgment.
_Avoid_: Failure, timeout, budget exceeded

**Unresolved Outcome**:
A terminal Run outcome where the cycle completed but unresolved Work Items remain: Unresolved Review Issues blocking Final Push, Tasks not completed, or a failing QA verdict. Distinct from Failed, which means the Run itself broke.
_Avoid_: Failure, crash, partial success

**Clean**:
A terminal Run outcome where the cycle completed with nothing unresolved remaining: no Unresolved Review Issues for review work, or every Task completed — and a passing QA verdict when requested — for spec work.
_Avoid_: Success, done, green

**Clean Unverified**:
A terminal review Run outcome where the cycle completed with nothing unresolved and the Final Push succeeded, but accepted Review Source Evidence for the pushed head never appeared within the grace period, so Merge-Ready was not confirmed. Distinct from Clean by outcome and exit code `3`; an explicit skipped review ends Review Skipped instead.
_Avoid_: Clean, failure, timeout

**Review Skipped**:
A terminal review Run outcome where the Review Source explicitly reports that it did not review the expected head. It is never Merge-Ready, carries the source-reported reason and next action, and uses a non-zero exit code.
_Avoid_: Clean, Clean Unverified, zero Review Issues

**Run Budget**:
A safeguard that stops a Run before it can continue indefinitely and indirectly consume unbounded resources.
_Avoid_: Max rounds, review round limit

**Preflight Validation**:
The early checks Roundfix runs before starting a Run or work that would make the developer wait.
_Avoid_: Best-effort validation, late failure

**Branch Integrity Preflight**:
The deterministic Preflight Validation for review Runs that blocks fetch, resolve, and watch while unintegrated Run Branch commits or another Run remain bound to the PR Head Branch, integrating fast-forwardable work automatically except QA-report-only branches and otherwise naming each pending worktree, run id, and recovery command. Skippable only through an explicit bypass that publishes an audit comment on the pull request.
_Avoid_: Advisory check, best-effort warning, soft gate

**Verification**:
The authoritative command or commands the Daemon runs verbatim in the repository root to decide whether Agent or Settle Command work can be settled and committed. A failure returns only its diagnostics to the Agent Session; for Tasks, a pass is required before status `completed`.
_Avoid_: CI, smoke test, best-effort check

**Waiting for Verification**:
The observable per-Task phase after Agent work is implementation-ready and before the Task acquires Verification Capacity; it is distinct from an Agent that is still working and from a Verification command that has started.
_Avoid_: Queued Agent, pending Task, blocked Run

**Temporary Verification Failure**:
A project-authored Verification command exit with code `75`, eligible for one Daemon-controlled exclusive retry per Task; Roundfix never infers it from logs, timing, or framework-specific error text.
_Avoid_: Flaky test, generic non-zero exit, log-matched infrastructure error

**Verification Feedback**:
The failure diagnostics returned to an Agent Session after the Daemon runs Verification. Passing Verification produces no Agent feedback.
_Avoid_: Full verification output, test log, progress stream

**User Config**:
Configuration that applies to Roundfix runs started by one developer across repositories.
_Avoid_: Global config, machine config

**Roundfix Home**:
The user-scoped directory where Roundfix stores configuration and central state across repositories.
_Avoid_: Workspace, artifact directory

**Project Config**:
Configuration that applies to Roundfix runs inside one repository.
_Avoid_: Local config, repo config

**Artifact Directory**:
The configured directory that overrides review-artifact placement and stores
Artifact Directory-backed Run files such as Detached Run console logs and
opt-in Agent logs. When unset, review artifacts use the Spec tree resolver.
_Avoid_: Workspace, cache, output folder

**Console Log**:
The compact caller-visible text record of a Detached Run's progress. It summarizes Agent file reads and edits while the Run Event Journal retains their lossless payloads.
_Avoid_: Agent log, Run Event Journal, audit log

**Run Outcome Notification**:
A best-effort terminal notice sent through the configured command or native desktop route, carrying the Run outcome and actionable context while its delivery attempt is recorded separately from the outcome.
_Avoid_: Run Event Stream, guaranteed delivery, terminal outcome

**Compatible Artifacts**:
Downloaded markdown artifacts that match the Head Repository, PR Head Branch, and pull request number being resolved.
_Avoid_: Matching by pull request number only, latest artifacts

**Run Database**:
The central Roundfix database that stores Run state and review progress across repositories.
_Avoid_: Artifact directory, state file

**Round**:
One review-resolution cycle within a Run, tied to the pull request state being reviewed.
_Avoid_: Review Round, iteration, pass, cycle

**Review Issue**:
Roundfix's local representation of one unresolved Review Source finding that may need triage or a code change.
_Avoid_: Comment, thread, finding

**Review Issue Fingerprint**:
The stable identity Roundfix uses to recognize the same Review Issue across Rounds.
_Avoid_: File path only, line number only, round-local id

**Source Reference**:
The Review Source-native identity stored as `source_ref` on a Review Issue, such as a CodeRabbit `thread:<id>,comment:<id>` pair.
_Avoid_: Local artifact path, generated issue number

**Duplicated Review Issue**:
A Review Issue that is complete because a newer occurrence with the same Review Issue Fingerprint is being resolved instead.
_Avoid_: Duplicate Review Issue, resolved issue, ignored issue

**Terminal Review Issue**:
A Review Issue whose local outcome is complete for the current Round because it is resolved, invalid, or duplicated.
_Avoid_: Done issue, closed issue

**Settled Review Issue**:
A Review Issue that needs no further Agent work in the current Batch because it is resolved, invalid, duplicated, or failed.
_Avoid_: Terminal issue, done issue

**Failed Review Issue**:
A Review Issue whose latest resolution attempt did not complete: the Agent could not fix it or verification did not pass. It stays Unresolved and is retried when a later Round downloads its still-open Review Source thread.
_Avoid_: Terminal issue, invalid issue, abandoned issue

**Unresolved Review Issue**:
A Review Issue that has been downloaded but has not reached a terminal local outcome.
_Avoid_: Open issue, pending task

**Batch**:
A bounded subset of Work Items assigned to one agent invocation.
_Avoid_: Chunk, group, task

**Outcome Comment**:
The idempotent comment Roundfix publishes on a Review Source thread whose local outcome is not resolved — the triage reason for invalid, the failed step for failed, the canonical thread for duplicated, the revisit plan for unresolved — so the pull request stays auditable without local artifacts.
_Avoid_: Silent resolve, status flip, reply thread

**Final Push**:
The Run-ending push that sends the PR Head Branch after no Unresolved Review Issues remain.
_Avoid_: Batch push, round push, agent push

**Resolve Command**:
The command that runs Agents over downloaded unresolved Review Issues for an Open Pull Request.
_Avoid_: Fix Command, Fetch command, watch command

**Implement Command**:
The command that executes a Spec's Task Graph by running Agents over its Tasks in dependency order.
_Avoid_: Run command, execute command, spec command

**Settle Command**:
The local recovery command that re-runs one failed Task's Verification commands in the current repository. On pass, it settles the Task `completed`, stages all current worktree changes plus the task file, and creates the standard Task commit; it creates no Run and never pushes.
_Avoid_: Retry command, auto-settle, task fix command

**Reconcile Command**:
The support command that inspects terminal spec Run Worktrees and Run Branches and reports their Run Worktree Reconciliation state. It is read-only by default; `--apply` is its only mutation switch, removes only freshly revalidated `safe` or `superseded` work, and has no force bypass.
_Avoid_: GC Command, Settle Command, automatic integration

**Reprocess Command**:
An explicit future command for revisiting selected Terminal Review Issues.
_Avoid_: Include resolved, resolve option

**Init Command**:
The support command that creates User Config or Project Config before operational Runs.
_Avoid_: Bootstrap run, setup run

**Setup Command**:
The support command that verifies and prepares a machine for Roundfix Runs: Node, minimum-supported acpx, Adapter Readiness, generated Agent Selection Profile Readiness, acpx local adapter overrides, User Config, and Project Config. It accepts the minimum tested acpx version and every newer compatible version, never downgrades a newer installation, proves proposed state before requesting authorization, and writes only authorized changes.
_Avoid_: Manual bootstrap checklist, environment wizard

**Upgrade Command**:
The support command that checks or installs the latest released Roundfix binary for the current platform.
_Avoid_: Package manager update, version check only

**Release Plan**:
A read-only classification of committed changes between a base release and a target revision that identifies the required semantic-version increment, proposes the next version, cites its evidence, and states whether explicit human approval or manual impact classification is required.
_Avoid_: Release execution, automatic release, version guess

**Release Plan Command**:
The support command that produces a Release Plan without editing release files, creating or pushing tags, publishing packages, or creating a GitHub Release.
_Avoid_: Release Command, publish command, cut-release command

**Doctor Command**:
The support command that diagnoses a repository and machine's readiness for Roundfix Runs — minimum-supported acpx, Adapter Readiness, Agent Selection Profile Readiness, Repository Skill Set, and codex runtime hygiene. It reports the detected acpx version against the minimum, gives each check a next action, and mutates nothing. It evaluates Repository Skill Set readiness after, and independently from, Agent Selection Profile Readiness; unlike the Doctor Command, the Setup Command prepares the machine.
_Avoid_: Health check run, setup run, environment wizard

**Archive Command**:
The support command that archives a completed Spec: it verifies every Task is completed and QA passed, stamps archive metadata, and moves the Spec folder to the archived spec root. Refuses a Spec with incomplete Tasks or no passing QA verdict.
_Avoid_: Move command, retire run, cleanup command

**GC Command**:
The support command that reclaims Run storage: it prunes the Run Event Journal and artifact directory of terminal Runs older than the Journal Retention window and removes orphaned run artifact directories, reporting what it freed. Never touches Active Runs, `runs` rows, or active-run locks.
_Avoid_: Clean command, vacuum, purge

**Journal Retention**:
The configured age window after which a terminal Run's Run Event Journal and artifact directory become eligible for pruning. Active Runs are never eligible; a retention of zero keeps everything. See ADR-0033.
_Avoid_: Log rotation, TTL, expiry

**Worktree Bootstrap**:
The configured command Roundfix runs once in a newly created Run or Task Worktree, after copying `worktree.copy` files and before Agent work and Verification, to prepare the environment (install dependencies, migrate and seed databases, warm caches). A bootstrap failure ends the Run or settles the Task with a bootstrap-failed outcome. See ADR-0034.
_Avoid_: Setup command, provisioning, install step

**Context-Driven Baseline**:
The portable, versioned set of required project instructions and architectural decisions that makes a repository ready for the CONTEXT-driven workflow.
_Avoid_: Sample docs, optional setup, repository template

**Baseline Profile**:
A versioned composition of modules, Repository Capabilities, decisions, managed artifacts, and verification roles that defines one reproducible Context-Driven Baseline. Roundfix ships built-in profiles; a repository may own additional profiles for its own workflow.
_Avoid_: Agent Selection Profile, setup template, user profile

**Profile Draft**:
A strict repository-owned Baseline Profile document supplied to planning in memory and admitted only after catalog validation and source binding.
_Avoid_: Partial profile, unchecked JSON, Profile override

**Profile Adaptation**:
A maintainer-reviewed Profile Draft that narrows one built-in Baseline Profile by removing repository-inapplicable modules or Repository Capabilities without weakening universal requirements.
_Avoid_: Profile waiver, capability bypass, silent inference

**Baseline Command**:
The public `roundfix baseline` command that authors and validates Baseline Profiles and drives one preflight-to-verification adoption or update flow without requiring an Agent Skill.
_Avoid_: Setup Command, setup skill, configuration wizard

**Source Baseline**:
The immutable Setup Manifest and governed instruction corpus that a setup transition or Baseline Readoption recognizes as its exact origin.
_Avoid_: Legacy Baseline, inferred preimage, current files

**Source Baseline Entry**:
One byte-evidenced structural unit in a Source Baseline that Baseline Readoption classifies and disposes individually as a Normative Clause, recommendation, Operational Contract, or non-governed evidence.
_Avoid_: Inferred rule, category summary, untracked paragraph

**Segmentation Snapshot**:
The immutable, checkout-free semantic-analysis input containing one Source Baseline and the strict byte-range proposal contract.
_Avoid_: Repository prompt, mutable instruction dump, checkout context

**Segmentation Proposal**:
A digest-bound, byte-exhaustive set of ordered ranges that splits Source Baseline Entries without rewriting or omitting source bytes.
_Avoid_: Summary, rewritten rules, partial range list

**Analysis Snapshot**:
The immutable, checkout-free classification input containing segmented Source Baseline Entries and only the semantic destinations active for the confirmed Baseline Profile and project decisions.
_Avoid_: Repository scan, unrestricted prompt, inferred destination list

**Classification Proposal**:
A digest-bound disposition for every entry in one Analysis Snapshot, admitted only after deterministic validation and explicit human review.
_Avoid_: Classification hint, partial proposal, autonomous decision

**Baseline Readoption**:
The confirmation-gated adoption that inventories an incompatible repository state as a Source Baseline and replaces its setup identity without treating existing instructions as disposable.
_Avoid_: Clean install, legacy fallback, automatic overwrite

**Normative Clause**:
An identified repository instruction whose enforcement is mandatory, prohibited, or stop-and-ask.
_Avoid_: Hard rule, guideline, free-form prose

**Normative Clause Manifest**:
The digest-bound inventory that independently accounts for a Source Baseline's Normative Clauses, recommendations, and Operational Contracts.
_Avoid_: Transition mapping, self-declared ledger, category checklist

**Operational Contract**:
An identified structured instruction whose order or shape carries required behavior, such as a template, procedure, or decision matrix.
_Avoid_: Prose summary, optional example, compressed guidance

**Instruction Hierarchy**:
The precedence order from universal instructions through context and documentation, Spec workflow, enabled autonomous work, stack and surface guides, and optional knowledge sources. A narrower guide may add constraints but cannot weaken a universal Normative Clause or confirmed project decision.
_Avoid_: Module list, arbitrary file order, duplicated root policy

**Standard TypeScript Monorepo Profile**:
The opinionated, project-agnostic Context-Driven Baseline profile for the repository's standard TypeScript monorepo stack.
_Avoid_: Project-specific profile, generic TypeScript profile, sample template

**Repository Capability**:
A profile-declared skill, tool, dependency, workspace, or repository contract whose required, recommended, or optional status is evaluated from explicit local evidence.
_Avoid_: Assumed stack, installed package list, inferred readiness

**Skill Activation**:
A declarative mapping from one stable work trigger to one exact required Agent Skill bundle, optionally conditioned on a selected Repository Capability.
_Avoid_: Suggested skill, inferred bundle, partial activation

**Decision Plan**:
The resolved setup proposal produced after every required setup decision has an answer; it is the basis for authorizing setup changes.
_Avoid_: Setup questionnaire, decision draft, configuration prompt

**Change Plan**:
The exact set of repository changes proposed for explicit review and confirmation before a setup or restoration operation mutates files.
_Avoid_: Patch, implicit apply, change preview

**Setup Manifest**:
The setup-owned record of the selected profile, modules, decisions, and managed artifacts that reproduces and audits a repository's Context-Driven Baseline.
_Avoid_: Project Config, lock file, generated guide

**Baseline ADR**:
An Architecture Decision Record whose reserved identity and invariant belong to the Context-Driven Baseline while project-specific notes remain repository-owned.
_Avoid_: Example ADR, project ADR, template copy

**Upgrade Retention Contract**:
The accounting a Context-Driven Baseline transition or Baseline Readoption must satisfy before mutation: every Source Baseline Normative Clause, recommendation, and Operational Contract maps to a current managed target, Repository-Specific Normative Rules, a recognized typed repository document, or an explicit rejection with a recorded reason.
_Avoid_: Best-effort migration, category coverage, silent rule removal

**Repository-Specific Normative Rules**:
Project-authored Normative Clauses that are not portable across the Context-Driven Baseline and cannot be represented by a typed project decision or managed semantic guide. Baseline assigns every representable rule to its semantic owner; only non-empty residual rules live outside setup markers in `docs/agents/specific-repository.md` and remain byte-preserved after confirmed adoption.
_Avoid_: Repository-Owned Extension, duplicate managed rule, baseline rule

**HTTP Contract Decision**:
The repository-owned choice of REST or POST-only application API semantics together with explicit protocol or operational exceptions and their owners.
_Avoid_: HTTP profile, universal REST rule, inferred route style

**Formatter-Stable Output**:
Generated managed Markdown that the target repository's selected formatter leaves unchanged, so apply, formatting, Verification, audit, and reapply compose with no delta.
_Avoid_: Renderer-canonical output, format-after-apply fixup

**Internal Identifier**:
A technical identity for an entity or resource that is generated and controlled by the project.
_Avoid_: External identifier, natural key, business code

**Roundfix Skill**:
A shipped agent skill that teaches an external Agent how to start Roundfix or how to resolve one assigned Batch.
_Avoid_: Runtime, Review Source, plugin

**Repository Skill Set**:
The complete set of Roundfix-owned and externally managed Agent Skills required by a repository workflow. Every required installed artifact matches its authoritative local source; unrelated extras do not join the set.
_Avoid_: Installed skills, plugin set, global skill set

**Interactive Input**:
The TUI flow that collects command parameters before a Run starts.
_Avoid_: Wizard, form, setup screen

**Live Run View**:
The TUI view that shows Review Issues and Run Events for a Run: streaming live while the Run is active, or replayed from the Run Event Journal during Attach.
_Avoid_: Dashboard, report, log file

**Run Browser**:
The read-only, machine-wide TUI list for Run discovery: every repository's Runs with their state, kind, target, Agent, start, duration, branch, and repository — Active Runs by default — from which selecting a Run opens the Live Run View through Attach.
_Avoid_: Run picker, run manager, dashboard

**Cockpit**:
The shared Live Run View composition made of the Phase Row, Work Queue, Session Timeline, footer, and optional Detail Modal.
_Avoid_: Dashboard, reduced pane

**Work Queue**:
The Live Run View pane that lists a Run's Work Items and their current status for both review and spec Runs.
_Avoid_: Issue list, task pane

**Phase Row**:
The Live Run View row that shows a Run's lifecycle position with terminal-text status markers.
_Avoid_: Progress header, pipeline bar

**Session Timeline**:
The Live Run View pane that shows grouped Run Events from the Run Event Journal.
_Avoid_: Console log, agent output pane

**Detail Modal**:
The centered Live Run View overlay that shows one selected Work Item's Review Issue artifact or Task file body read-only.
_Avoid_: Detail pane, inspector panel

**Daemon**:
The Roundfix process that owns the Run lifecycle and Review Source-facing outcomes.
_Avoid_: Orchestrator, controller, manager

**Run Event**:
One ordered product record of something meaningful that happened during a Run, carrying Run identity, Batch when known, event source, event kind, and a structured payload. Producers convert their native models into Run Events; ACP stream updates remain an Agent-internal protocol model.
_Avoid_: Stream update, log line, message

**Run Event Journal**:
The append-only history of Run Events stored in the Run Database, ordered by a per-Run cursor so replay is deterministic and duplicate-free.
_Avoid_: Agent log, log file, event broker

**Run Event Stream**:
A read-only JSONL projection of one Run's Run Event Journal, selected by an explicit Run ID and optionally followed until the Run reaches a terminal outcome. Its stable Supervisor filters are task status, Batch boundary, Verification verdict, and terminal outcome.
_Avoid_: Attach, console log, global event bus

**Attach**:
Viewing a Run by replaying its Run Event Journal and then following new Run Events, without owning, mutating, or stopping the Run.
_Avoid_: Resume, reconnect, takeover

**Agent**:
The local coding assistant invoked by Roundfix to triage and resolve an assigned Batch.
_Avoid_: Review Source, review provider, worker, bot

**Supervisor**:
The external Claude Code role that authors Specs, starts and monitors Runs, and delegates every Work Item to an Agent.
_Avoid_: Fable, Agent, ACP Runtime, Daemon, Orchestrator

**Follow Mode**:
The Live Run View state in which the timeline tail advances automatically as new Run Events arrive; suspended while the user scrolls back, resumed when the viewport returns to the bottom. Scrolling never affects the Run.
_Avoid_: tail mode, auto-scroll, live mode
