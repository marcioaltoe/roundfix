# Roundfix

Roundfix picks up Work Items — Review Issues from pull request reviews today, Tasks from Spec Task Graphs next — resolves them through the user's selected ACP Runtime, and pushes only when nothing unresolved remains. This glossary defines the product language for that loop.

## Language

**Behavior Surface**:
One user-visible unit of Roundfix behavior: a command's help, the set of configuration keys, the table of exit codes, or one user-guide page ([ADR-0256](docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md)).
_Avoid_: Source file, command family, feature

**Skill Coverage Map**:
The authored inventory of Behavior Surfaces and the owned skill files that describe each, or a reason no skill does, together with their sources and any Coverage Review ([ADR-0256](docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md)).
_Avoid_: Skill index, generated coverage, fingerprint map

**Behavior Surface Record**:
The generated inventory of one fingerprint per Behavior Surface, kept separate from the authored Skill Coverage Map ([ADR-0256](docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md)).
_Avoid_: Coverage map, authored snapshot, skill digest

**Coverage Review**:
A dated explanation in a Skill Coverage Map entry of why a behavior change needs no skill text; it applies only to the release range in which the explanation changed ([ADR-0256](docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md)).
_Avoid_: Permanent waiver, review approval, uncovered reason

**Skills Declaration**:
A Spec's optional declaration of the covered Behavior Surfaces whose changes need no skill text, each with a reason and a planned Coverage Review ([ADR-0256](docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md)).
_Avoid_: Skills inventory, permanent waiver, Coverage Review

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

**Assigned Repair**:
A repair a QA gate's own Task file instructs it to make, as opposed to a finding it observes and reports. The gate performs an Assigned Repair, verifies it, and records what it performed; it may write only the paths its Task names. Anything the Task did not assign stays reported rather than performed. In code and in a report's audit trail the same idea appears lowercase, as an assigned repair.
_Avoid_: Auto-fix, gate edit, remediation

**Governed Path**:
A repository path the tooling-authority rules bind — the configuration, scripts, ignore files, plugin declarations, and version pins of linters, formatters, typecheckers, test runners, architecture checkers, build tools, package managers, and code generators. The changed-path audit judges a Task commit only against its governed paths; an ordinary source, test, or documentation file is not governed and needs no grant. The declared set is held to the record: every path any authorization has ever bounded stays governed.
_Avoid_: Protected file, restricted path, tooling file

**Delivery Base**:
The merge base of the audited head and the repository default branch. The pre-PR review diffs the candidate from it. The mechanical stage reads a Spec's Task commits and its grant from it.
_Avoid_: Run start head, audited head, delivery target

**Undeclared Governed Path**:
A Governed Path named by a pending non-QA Task but omitted from its authorization record or a present Tooling authority row.
_Avoid_: Unauthorized edit, inferred path, ordinary path

**Undocumented CLI Surface**:
A CLI source or contract-test path named by a pending non-QA Task that names no guide in itself or its dependency ancestry.
_Avoid_: Missing documentation, undocumented command, downstream guide

**Process Residue**:
A process Roundfix started that outlived the Run that started it, so no live Run record owns it. The readiness diagnostic reports residue with each process's age, its consumed CPU time, and its originating Run when the Run Database still knows it; reporting settles nothing and creates no Run record.
_Avoid_: Orphan, zombie, leaked run, stale process

**PR Head Branch**:
The branch on GitHub that supplies the pull request's head commits.
_Avoid_: Local branch, checkout branch

**Head Repository**:
The GitHub repository that owns the PR Head Branch.
_Avoid_: Base repository, local checkout

**Work Item**:
One unit of resolvable work that Roundfix picks up and drives to an outcome: a Review Issue today, a Task next.
_Avoid_: Ticket, job, to-do

**Backlog Entry**:
Typed intent for what to do next, with a lifecycle of `open`, `promoted` to a named Spec, or `declined` with a reason; it is never evidence, just as a finding is never a commitment. Its types are the Conventional Commits intent vocabulary — `feat`, `fix`, `perf`, `refactor` — so one word carries intent from entry to Spec to commit; a `feat` entry is upstream raw material, never the `write-idea` artifact.
_Avoid_: Finding, idea artifact, untyped suggestion

**Inbox Entry**:
A fleeting capture born committed in the Secondbrain's inbox under its destination's namespace. It carries its origin, destination, advisory type hint, and capture mode, and remains pending until Triage resolves it into exactly one Finding, one Backlog Entry, or one recorded discard.
_Avoid_: Project inbox item, uncommitted note, draft artifact

**Rollup**:
A Finding of `rollup` kind that consolidates related Findings as declared members, supersedes those members, and licenses their archival.
_Avoid_: Summary, digest, archive manifest

**Triage**:
The destination-project act that converts one pending Inbox Entry into its contract-true artifact, committed by the destination alone.
_Avoid_: Capture, import, synchronization

**Task Source**:
The origin Roundfix reads Work Items from, such as a Review Source or a Spec's Task Graph.
_Avoid_: Provider, backend, integration

**Review Source**:
The external review system that produces feedback for an Open Pull Request.
This legacy PR-feedback concept is distinct from a Pre-PR Review Provider.

_Avoid_: Review Provider, Agent, ACP Runtime

**Review Source Evidence**:
A head-bound Review Source classification for one expected commit: `pending` has no usable expected-head signal, `reviewing` is still in progress, and `reviewed` is complete without proving Merge-Ready. `verified` proves the expected head with no unresolved Review Issues, `skipped` explicitly declines that head, and `failed` records an explicit Review Source failure.
_Avoid_: Check presence, generic approval, inferred clean state

**Review Finding Disposition**:
A Pre-PR Review finding's recorded fix or evidence-backed dismissal, tied to the head the reviewer examined.
At the same repository, base, head, and provider, the findings verdict is reused without asking the reviewer; only exact evidence-backed dismissals of every finding produce `findings-dismissed`, while a fix requires a fresh review of its changed head.
_Avoid_: Review Issue resolution, unrecorded operator decision, finding deletion

**Spec**:
One Spec's planning artifact set produced by the spec workflow: PRD, Task Graph, Task files, QA evidence, and adopted sources under `references/` with provenance recorded in `references/_index.md`.
_Avoid_: Feature folder, epic, project

**Corrective Spec**:
A new Spec, with its own authorization and QA gate, that corrects a finding on an archived Spec; the archived Spec is never edited in place.
_Avoid_: archive edit, inherited authorization; a Corrective Task fixes an active Spec, not an archived one

**Corrective Task**:
A Task appended to an active Spec's Task Graph, before its QA gate, to fix a finding on that Spec before it is archived; an archived Spec gets a Corrective Spec instead.
_Avoid_: Corrective Spec, follow-up Task

**Cause Class**:
The one class from a closed list that `roundfix runs causes` assigns to a failed Verification attempt or a Corrective Task: `scope-or-authorization`, `shared-section-contract`, `repository-convention`, `implementation-defect` or `environment`, or `unclassified` when no signature matches. It is assigned by signature, never by a model judgment (ADR-0214).
_Avoid_: failure category, root cause, model verdict

**Classification Signature**:
An ordered entry of the embedded signature table, with an identifier, a Cause Class, the one evidence source it reads and a regular expression; the first entry that matches classifies the item (ADR-0214).
_Avoid_: heuristic, rule, pattern

**Supersession**:
A durable lifecycle record that names the Spec which delivered another Spec's content. The Archive Command accepts it as proof of completion for the superseded Spec when that Spec has no Task Graph.
_Avoid_: Replacement, duplicate, hand-written disposition

**Spec Root**:
The configured directory holding Spec folders, `docs/specs` by default; it may resolve outside the repository working tree, such as a knowledge workspace repository.
_Avoid_: Specs directory, docs folder, knowledge base

**Task**:
One implementable unit of work within a Spec. Its task file is the sole owner of its status.
_Avoid_: Subtask, story, ticket

**Recorded Path**:
A path the Daemon records under a Task file's `## Recorded paths` section when
the Task changed it without declaring it. The QA scope audit counts it as
declared and names it in the scope row; recording discloses a change and
reserves nothing, while a Governed Path still needs its authorization.
_Avoid_: Declared path, reserved path, authorized path

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
The qa-gate evidence report written to a Spec's QA directory, carrying a machine-readable verdict, its Auditing Binary as `auditing_binary`, `auditor_staleness`, and the public-row binary's `--version` line as `user_flow_binary`, plus `rows_blocked_environment`, `rows_blocked_finding`, and `rows_blocked_declared` counts in its frontmatter. The auditor fields are Daemon-owned: the gate keeps the seeded `auditing_binary` and `auditor_staleness` values. A report recording a Precondition Refusal carries `rows_blocked_precondition` beside those counts, plus the `precondition_check` and `precondition_reason` keys that name the refusal; a gate that reached its matrix writes none of those three. The pre-PR Pull Request row, recorded as `blocked (environment: no open Pull Request)` with the Pull Request row named in its provenance, never decides a qualifying partial and needs no Unreachable Acceptance declaration. A `pending` verdict is never accepted, a report that records no QA row is refused, and a report whose front matter is empty or duplicated is unreadable and refused; only the newest report in the directory is read by a later run's mechanical stage, so a superseded report blocks nothing.
_Avoid_: Test report, QA log

**Auditing Binary**:
The Daemon's Roundfix binary that runs the mechanical stage, distinct from the tree it audits. It carries the version, build commit, and build time (the last two may be empty for a released build), and a QA Report records its formatted identity as `AuditingBinary` / `auditing_binary`; `auditor_staleness` compares its build commit with the Delivery Base, and a stale Auditing Binary is published as a warning while the gate proceeds.
_Avoid_: Audited binary, audited tree, ambiguous build

**Precondition Refusal**:
A QA gate stop at a check that runs before the matrix is built — a strict Spec check being the usual one — recorded as the gate's entire QA Report rather than as a missing one: verdict `fail`, `rows_blocked_precondition` counting the stop, `precondition_check` naming what was checked, `precondition_reason` carrying every refusing code and the sentence beside it, and one terminal row `0 | blocked | precondition` in the Results table. The gate measured no requirement, so row `0` records the refusal itself instead of a result, and the prose justifying it is written as a list rather than a second table. An empty Results table is not the alternative: it refuses every later run on the report instead of on the Spec, and prescribes a repair — materialize every planned row — that a run which never built a matrix cannot perform.
_Avoid_: Empty report, failed gate, skipped QA, precondition error

**External Acceptance Evidence**:
Evidence for at least one named acceptance row that originates outside the Spec's own artifacts — a real repository, a measurement, or published literature — with its origin recorded. A row whose external evidence cannot be obtained is recorded as blocked with its reason. It exists because a rubric and the requirement it measures, written by one author from one premise, confirm rather than test.
_Avoid_: Self-authored rehearsal, internal fixture, passing suite

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

**Runtime Catalogue**:
What an ACP Runtime advertises before Roundfix asks it to apply a specific Agent Selection, established from a disposable Agent Session ensured without the requested model so the answer cannot be contaminated by the question. Membership in it decides refusal where no adapter refuses first; where an adapter refuses, its own message stands (ADR-0147).
_Avoid_: Model list, capability payload, advertised options

**Selection Encoding**:
How one Agent Selection's reasoning effort is represented in the advertised ACP controls, decided during proof and reported by readiness surfaces. `independent` assigns a separate advertised reasoning option. `model_variant` selects an advertised model identifier that already carries the effort. `model_managed` is an empty reasoning effort on an adapter advertising no reasoning control at all. `runtime_managed` is an empty reasoning effort on an ACP Runtime whose advertised reasoning control Roundfix declines to assign, so the observed value is the Agent Model's own and is not proof-relevant. `runtime_deferred` is a non-empty reasoning effort on an ACP Runtime that advertises the control only once the Agent Session holds a queue owner, so proof verifies the value is advertised and the Run applies and observes it after an inert setup turn and before the first work turn.
_Avoid_: Effort mode, reasoning strategy, assignment style

**Agent Selection Profile**:
The atomic policy for one Agent Work Category, containing one Preferred Selection and a non-empty ordered Fallback Chain. A higher-precedence profile replaces the complete lower-precedence profile rather than merging individual fields.
_Avoid_: Runtime defaults, model preset, partial override

**Profile Deviation**:
The dated, reasoned record on a configured Agent Selection Profile that its difference from the Recommended Profile is deliberate, which holds for the snapshot it names.
_Avoid_: Permanent pin, recommendation waiver

**Agent Selection Profile Readiness**:
The command-scoped result that resolves effective Agent Selection Profiles, deduplicates their Preferred Selections and Fallback Chains, and requires Exact Agent Selection Proof for every distinct tuple before a Run or configuration mutation.
_Avoid_: Single-model check, configured Agent probe, cached readiness

**Preferred Selection**:
The first Agent Selection Roundfix proves and attempts for an Agent Work Category.
_Avoid_: Default model, primary runtime, recommendation winner

**Fallback Chain**:
The non-empty ordered Agent Selection list Roundfix proves with the Preferred Selection before a Run and may activate after notifying the user and Supervisor that the preceding selection failed before Agent work began. A Lost Rollout before the First Handoff, or in a QA Task before its report, lets the next Fallback Selection take the Task.
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

**Recommended Profile**:
The dated Agent Selection Profile Roundfix recommends for one Agent Work Category, from which built-in profiles derive. It never selects, routes, or changes a configuration by itself.
_Avoid_: Model Recommendation Ranking, model router, benchmark policy, automatic selection

**Fallback Selection**:
The next configured Agent Selection in a profile's Fallback Chain. Roundfix proves it before the Run, emits a notification before activation, and may switch ACP Runtime automatically only while Agent work has not begun. A Lost Rollout before the First Handoff, or in a QA Task before its report, is the exception that lets this selection take the Task.
_Avoid_: Dynamic fallback, silent model switch, catalog probe winner

**Negative Control**:
The defect a Task declares its Verification must catch, named in the Task's own `## Negative Control` section and recorded beside its outcome. It is authored rather than synthesised, because manufacturing one means mutating the repository under test. A Task declaring none is recorded as having none, which is a weaker gate stated honestly.
_Avoid_: Failing test, mutation, sanity check

**Vacuous Verification**:
The Run Event classification (`verification_vacuous`) recorded when a Task's Verification command exits zero against the tree as it stands before the Agent runs. A command that passes before the work happened cannot be evidence that the work happened, so the Task is refused at dispatch with the offending command named and no Agent turn spent.
_Avoid_: Weak test, trivial gate, false green

**Inverted Verification Exit**:
The Spec Consistency Check error (`SC-VERIFY-INVERTED-EXIT`) raised when an authored Verification command uses a measured shell form whose exit status reverses or ignores the condition its output appears to assert. The finding names the matched form and a replacement that exits zero when the asserted condition holds.
_Avoid_: Verification failure, shell lint, non-zero result

**Wrap-Fragile Phrase Check**:
The Spec Consistency Check error (`SC-VERIFY-WRAP-FRAGILE`) raised when a pending Task uses a line-bound multi-word phrase check against Markdown that can miss text wrapped across lines. The finding names the phrase and file and gives a wrap-tolerant replacement.
_Avoid_: Line-bound grep, phrase grep, wrapped-text failure

**Non-Hermetic Verification**:
The Spec Consistency Check error (`SC-VERIFY-NON-HERMETIC`) raised when an authored Verification command depends on an undeclared environment variable or a pre-existing path outside the repository. A command-local variable or a path the Task creates before reading is not an external dependency.
_Avoid_: External-state Verification, environment guard, temporary-path check

**Unobserved Verification**:
The Run Event classification (`verification_unknown`) recorded when a Verification command's verdict could not be observed — a timeout, a partial execution, or a runner error that is not the command's answer. It is distinct from a command that ran and exited non-zero, which has a verdict, and it keeps "the work is wrong" separable from "we did not find out".
_Avoid_: Flaky test, transient failure, retry

**Unsupported Citation**:
The Spec check finding (`SC-CITATION-UNSUPPORTED`) raised when an artifact attributes a subject to a decision record that the record's own text does not carry. It is distinct from an unlisted citation, which asks whether a record was named; this one asks whether the claim about it holds, and reports both texts so a maintainer settles it by reading rather than trusting the checker.
_Avoid_: Broken link, missing ADR, stale reference

**Agent Work Started**:
The status (`agent_work_started`) marking the first Agent output that could have changed something, published once per Agent Session. It is the boundary after which a Fallback Selection may no longer switch ACP Runtime, because a second runtime would inherit state the first one built. Preparing or activating a Session does not reach it. The exception is a Lost Rollout before the First Handoff, or in a QA Task before its report, which lets the next Fallback Selection take the Task.
_Avoid_: Session opened, prompt sent, turn started

**Selection Failure**:
The outcome (`agent_selection_failed`) of a turn that ended before any Agent output, because the selected ACP Runtime would not serve it — exhausted quota, failed authentication, an adapter that will not start. It is a failure of the selection rather than of the work, so the Fallback Chain stays eligible.
_Avoid_: Batch failure, agent crash, runtime error

**Hook Refusal**:
The outcome (`hook_refused`) of a Task commit a repository commit hook rejected after the authoritative Verification had already passed. It is a repository misconfiguration rather than a Task failure — a commit hook may never be stricter than the Verification — so the Task stays `completed`, its verified work stays staged in the surface that produced it, and the Run reports the refusing hook, its exit code, its output, and the Settle Command that recovers the work.
_Avoid_: Commit failure, pre-commit error, verification failure

**Mechanical Refusal Code**:
The stable token a pre-QA mechanical check emits when it can refuse from written declarations and repository facts alone, before any Agent Session opens. Four exist: `QA-AUTH-PATHS` (a tooling change outside its authorization's bounded paths), `QA-CONSEQUENT-ORDER` (a consequent fix folded into or ordered before the change that caused it), `QA-REPORT-SHAPE` (a QA Report missing a required structural element), and `QA-EVIDENCE-PATH` (an evidence path a report names but the repository does not carry). The code is the durable name a reader and a later detector both use; the human sentence beside it may be reworded, the token may not.
_Avoid_: lint code, error code, check name

**Grant Refusal Code**:
The stable token the Spec checker emits when an authorization boundary refuses. `SC-TOOLING-UNAPPROVED` means a Spec claims tooling authorization its cited record does not sustain because the record is proposed, undated, withdrawn, or consumed by another Spec. A record that honestly declares its mutations proposed, claiming no operative grant, is accurate rather than defective and refuses nothing. Like a Mechanical Refusal Code, the token is the durable name a reader and a later detector both use; the sentence beside it may be reworded, the token may not.
_Avoid_: authorization error, permission code, trust error


**Agent Session**:
The acpx-backed session owned by one Work Item or action. Each Implement Task owns a Task Type-selected Agent Session, requested QA owns a separate `qa` Agent Session, and review work uses a review-selected Agent Session; effective selection and fallback attempts are persisted for that owner.
_Avoid_: ACP session, chat, conversation, thread

**Lost Rollout**:
The runtime infrastructure failure (`rollout_lost`) in which the ACP Runtime cannot find the persisted record of the Agent Session it is asked to continue. The Daemon recovers it in a new Agent Session without repair or a counted queue retry (ADR-0245).
_Avoid_: Session loss, missing prompt, retryable Task failure

**First Handoff**:
The first Agent turn of a Task that returns to the Daemon (ADR-0057's "Agent handoff"). Before it, a Lost Rollout lets the Fallback Chain take the Task.
_Avoid_: Session opening, first prompt, Agent Work Started

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

**Branch Disposition**:
A recorded terminal reason a Run Branch no longer needs integration and may be discarded. A superseded Branch Disposition means every branch commit is already reachable from the target, or a later integrated Run covered the same Tasks.
_Avoid_: Branch cleanup, forced deletion, reconcile result

**Run Worktree Reconciliation**:
The proof-based classification of a terminal spec Run's retained Git surfaces: `safe` when the Run Branch is contained in its target or its content is represented at the merged head; `superseded` when a newer QA Report or the merged-head proof represents its Task or QA Report commits; `unintegrated` when resolved evidence does not represent the Run; `dirty` when a present Run Worktree has tracked or untracked changes; `unknown` when metadata or Git evidence cannot prove another state; and `released` only when both the Run Worktree and Run Branch are absent. Only freshly revalidated `safe` or `superseded` work can be cleaned up.
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

**Run Window**:
The durable, repository-scoped bound on when an Implement Run may be created. It governs a Run's start and never its finish: a Run created before the cutoff continues to its own terminal outcome. The implementation spellings `RunWindow` and `run_window` refer to this same term.
_Avoid_: Session — Agent Session names a different concept; curfew and deadline suggest the wrong boundary

**Preflight Validation**:
The early checks Roundfix runs before starting a Run or work that would make the developer wait.
_Avoid_: Best-effort validation, late failure

**Branch Integrity Preflight**:
The deterministic Preflight Validation for review Runs that blocks fetch, resolve, and watch while unintegrated Run Branch commits or another Run remain bound to the PR Head Branch, integrating fast-forwardable work automatically except QA-report-only branches and otherwise naming each pending worktree, run id, and recovery command. Skippable only through an explicit bypass that publishes an audit comment on the pull request.
_Avoid_: Advisory check, best-effort warning, soft gate

**Verification**:
The authoritative command or commands the Daemon runs verbatim in the repository root to decide whether Agent or Settle Command work can be settled and committed. A failure returns only its diagnostics to the Agent Session; for Tasks, a pass is required before status `completed`.
_Avoid_: CI, smoke test, best-effort check

**Repository Contract Test**:
A Go test under the `docscontract` or `repocontract` build tag that checks the repository itself — its documents, derived artifacts or test wiring — rather than one package's behavior. The Full Contract Run runs every one; `make verify-changed` runs the ones its Contract Relevance selects, and `make verify-docs` runs the `internal/docscontract` tests and the contracts its `repo-test` step lists (ADR-0252, ADR-0253).
_Avoid_: docs test, integration test, CI-only test

**Contract Relevance**:
The declaration in a Repository Contract Test's source that decides when `make verify-changed` runs it: `always` on every change; by default, when a file in its package directory changes; `relevant`, when its package directory or a declared path changes; or `boundary`, never in that gate. A malformed declaration fails the gate. A change to `go.mod`, `go.sum` or `Makefile`, and a change list that Git cannot produce, each select every contract that is not `boundary`. The selector's summary line names each `relevant` contract it leaves out and each `boundary` contract. The Full Contract Run ignores Contract Relevance (ADR-0252, ADR-0253).
_Avoid_: test tags, impact map, test filter

**Full Contract Run**:
`make verify-contracts`: one run of every Repository Contract Test that `cmd/verify-select` discovers, whatever its Contract Relevance, with no hand-kept list. CI runs it on every push to main and before every release, so a contract the selective gate leaves out of a pull request still runs before the change ships (ADR-0253).
_Avoid_: full suite, nightly run, repo-test

**Incremental Verification**:
The selected fast local check recorded by the `verification.incremental` Baseline decision for validating the current change outside a Run while reusing safe local state. In a Daemon-assigned turn the Agent runs focused tests instead, because the Daemon runs the Task's Verification and the repository Verification at settlement (ADR-0257). It is distinct from the complete repository Verification selected by `verification.gate`.
_Avoid_: Complete Verification, CI gate, optional check

**Waiting for Verification**:
The observable per-Task phase after Agent work is implementation-ready and before the Task acquires Verification Capacity; it is distinct from an Agent that is still working and from a Verification command that has started.
_Avoid_: Queued Agent, pending Task, blocked Run

**Temporary Verification Failure**:
A project-authored Verification command exit with code `75`, eligible for one Daemon-controlled exclusive retry per Task; Roundfix never infers it from logs, timing, or framework-specific error text.
_Avoid_: Flaky test, generic non-zero exit, log-matched infrastructure error

**Verification Feedback**:
The failure diagnostics returned to an Agent Session after the Daemon runs Verification. Passing Verification produces no Agent feedback.
_Avoid_: Full verification output, test log, progress stream

**Settlement Check**:
A gate fact the Daemon checks for a non-QA Task before it settles in a Task Graph with an authored QA gate. The checks cover Spec Consistency, the repository Verification at settlement, and authorization of the prospective Task commit.
_Avoid_: QA settlement, archive check, post-Run check

**Repeated Failure**:
A Verification failure whose normalised diagnostic signature matches an earlier failure of the same Work Item. The signal names the earlier Run and attempt in the Run record, Run Event Stream, and Verification Feedback.
_Avoid_: New failure, duplicate error, retried failure

**User Config**:
Configuration that applies to Roundfix runs started by one developer across repositories.
_Avoid_: Global config, machine config

**Roundfix Home**:
The user-scoped directory where Roundfix stores configuration and central state across repositories.
_Avoid_: Workspace, artifact directory

**Project Config**:
Configuration that applies to Roundfix runs inside one repository.
A Task commit includes Project Config only when the frozen Spec authorization bounds `.roundfixrc.yml`; otherwise the Task fails with `Project Config outside the Spec's authorization`, while Batch and QA Report commits never stage Project Config and report the exclusion.
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
The local recovery command that re-runs one Task's Verification commands in the surface that still holds its work — a kept Task Worktree, a kept Run Worktree, or the current repository. It recovers a Task that is failed, or completed with uncommitted work; the second is the state a Hook Refusal leaves behind, and a completed Task whose surfaces are all clean was committed already and settles nothing. On pass, it settles the Task `completed`, stages all current worktree changes plus the task file, and creates the standard Task commit; it creates no Run and never pushes.
_Avoid_: Retry command, auto-settle, task fix command

**Reconcile Command**:
The support command that inspects terminal spec Run Worktrees and Run Branches and reports their Run Worktree Reconciliation state. It is read-only by default and carries three named mutation switches, none of which has a force bypass: `--apply` is the only one that releases retained Git surfaces, removing only freshly revalidated `safe` or `superseded` work; `--discard-superseded` records a Branch Disposition before removing a branch it can prove superseded; and `--carry-forward` performs Task Carry-Forward.
_Avoid_: GC Command, Settle Command, automatic integration

**Task Carry-Forward**:
The explicit act that hands a settled Task from a terminal spec Run's Run Branch back to the user's checkout, so work that already ran and passed its Verification is never executed again to reach the same result. It accepts a Run whose outcome is Stopped or Unresolved and refuses every other terminal outcome. It carries a Task only on proof — a passing Verification verdict, exactly one settlement commit, declared inputs unmoved since settlement, a clean checkout, and a repository-local Specs Root — refuses the whole set rather than carrying part of it, and stamps each carried Task with the Run and commit that established it. Reached through the Reconcile Command's `--carry-forward` switch or through a Delivery Retry; it is never automatic. A Task already completed on the target is nothing to carry, so the act does not repeat itself. The implementation spelling `CarryForward` refers to this same term.
_Avoid_: replay, resume, automatic carry, QA row carry-forward

**Delivery Queue**:
The durable, ordered set of Specs that Roundfix advances from an Implement Run through merge with a fixed optional deadline and per-item retry limit. A blocker parks only its item while later items can continue.
_Avoid_: Release queue, Batch, Run list

**Delivery Plan**:
The read-only account of which Specs have delivery authority, what blocks the others and which production premises ordered Specs share. It reports a prepared queue without recording or running one.
_Avoid_: Delivery Queue, execution plan, approval grant

**Delivery Revalidation**:
The check of a queued Spec against the refreshed default branch immediately before its first Run. A strict finding parks it as `revalidation-failed`; a production-premise overlap records a `premise-changed` warning without parking it.
_Avoid_: Initial validation, retry, approval check

**Delivery Queue Limit**:
A queue's fixed optional deadline, per-item retry allowance or Queue Token Ceiling. It bounds new item starts or explicit retries without changing a Run Budget or an item already past `queued`.
_Avoid_: Run Budget, concurrency

**Pending Question**:
The one operator decision presented for the lowest-position parked Delivery Queue item, with its blocker, answer and number of parked items waiting behind it. Only an explicit Delivery Retry or a new queue answers it.
_Avoid_: Warning, automatic recommendation, blocker list

**Delivery Retry**:
The explicit act that returns one parked Delivery Queue item to the stage supported by its recorded evidence and hands it to a live or newly started queue owner. It may perform Task Carry-Forward from every Run of the item's Spec, newest first, and is never automatic.
_Avoid_: Resume Command, automatic retry, replay

**Judge Log**:
The per-month JSONL record of every advisory judgment call: the question, the transport that carried it (OpenRouter or TypeSafe), the answer, probabilities, confidence, latency, input tokens, the model the service reported and the call's cost. Its monthly sum is what the spend ceiling reads (ADR-0201).
_Avoid_: Run Event Journal, audit log

**Advisory Judgment**:
A typed semantic judgment about a Spec artifact, such as whether a cited ADR supports the attributed claim, or about a source relationship, such as whether a Finding or Backlog Entry shares a Spec's context, that the authoring model must answer but that never gates, never changes another command's exit code, and fails open (ADR-0200, ADR-0209).
_Avoid_: Gate, check, finding

**Grouping Suggestion**:
An Advisory Judgment that an open Finding or Backlog Entry shares the context of a source the Spec already adopts, so the author adopts it into the same Spec or states why it stays apart; it never moves a source by itself (ADR-0209).
_Avoid_: Merge rule, automatic grouping

**Grouping Bound**:
The size limit on grouping sources into one Spec: four implementation Tasks plus its QA gate. A grouped scope that needs more is split into Specs that each fit, and each source keeps exactly one owning Spec (ADR-0208).
_Avoid_: Spec size rule, Task cap

**Frontend Layout Decision**:
The repository's recorded choice of frontend layout, `systems` or `repository-defined`. While none is recorded, the Baseline states the suggested `systems` layout; a recorded value is never overwritten (ADR-0205).
_Avoid_: Frontend architecture rule, mandated layout

**Composed Setup Snapshot**:
A setup snapshot built by name from upstream setups in a fixed order, with duplicate skills dropped, so a built-in profile can span several toolchains (ADR-0204).
_Avoid_: Merged setup, custom setup

**Retired Skill**:
A skill the Baseline no longer requires although the upstream catalog may still list it. No module requires or dispatches it, the asset sync drops it from every Setup Snapshot, and `roundfix baseline update` lists an installed copy for the adopter to delete without deleting it (ADR-0246).
_Avoid_: Deprecated skill, obsolete skill

**Park Class**:
The named reason a Delivery Queue item is parked — for example a prerequisite not yet merged, a Pull Request conflict, an environment-only QA partial or a check that failed twice outside the item's packages — shown with the next command that answers it (ADR-0192, ADR-0193).
_Avoid_: Blocker text, error code

**Prerequisite Spec**:
A Spec another Spec names in its Task Graph as `requires`; the queue owner does not start the dependent item until each prerequisite is merged into the default branch (ADR-0193).
_Avoid_: Queue order, dependency Task

**Derived Path Declaration**:
A Project Config entry naming a path pattern and the sanctioned command that regenerates it; a Pull Request conflict confined to declared derived paths is resolved by regeneration instead of parking (ADR-0192).
_Avoid_: Generated file list, merge strategy

**Module Version Record**:
A generated record of each Baseline module's versions and content digests, so a module version names one recorded content and the record step can choose the next free version (ADR-0250).
_Avoid_: Module changelog, manual version note

**Coverage Record**:
A generated record of the tests in the repository's packages, including the release platforms on which platform-limited tests are built, so coverage comparisons are stable on every host (ADR-0250).
_Avoid_: Host test list, test result report

**Evidence Snapshot**:
The Daemon's record, written when a QA pass closes, of the inputs each carriable passing row observed at the audited head (ADR-0194).
_Avoid_: QA cache, Agent-written evidence

**Carried Row**:
A QA matrix row a later pass keeps from the prior report because its Evidence Snapshot proves its inputs unmoved; it keeps its original provenance. Rows that read the repository Verification, the Pull Request or the commit range are never carried (ADR-0194, ADR-0195).
_Avoid_: Skipped row, Task Carry-Forward

**Carry Disposition**:
The fixed-list outcome a QA pass records for every row of the prior report: carried, re-run as stale, re-run as always observed, or re-run because it did not pass (ADR-0195).
_Avoid_: Row status, verdict

**Delivery Convention**:
A commit shape the delivery itself produces — the QA Report commit, the Daemon's status writes, the archive commit, an all-pending planning candidate — that the pre-PR reviewer is told about and never raises as a defect (ADR-0196).
_Avoid_: Review exemption, ignore list

**Reviewer Lineage**:
The at-most-two review rounds of one delivery candidate, where round 2 reviews only the delta since round 1 with round 1's findings and dispositions, continuing the same Agent Session when the runtime can resume it (ADR-0197).
_Avoid_: Review cache, third round

**Turn Usage**:
The token count an ACP adapter reports for one prompt of a Run, recorded per Task and summed per Run; a prompt without a report is recorded as unreported, never as zero (ADR-0198).
_Avoid_: Cost estimate, spend

**Usage Basis**:
How a Turn Usage was counted: `turn` when the adapter's end-of-prompt count covers the whole turn, `request-sum` when Roundfix sums the adapter's per-request context readings (ADR-0198).
_Avoid_: Pricing model

**Reported Cost**:
The spend an ACP adapter itself reports for a turn; Roundfix prints only reported cost and never prices tokens from a table (ADR-0198).
_Avoid_: Estimated cost, price table

**Queue Token Ceiling**:
An optional Delivery Queue Limit: once the queue's recorded Turn Usage reaches it, the next item parks and retries are refused, while a running Task is never stopped (ADR-0199).
_Avoid_: Run Budget, spending limit

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

**Release Set**:
The launcher package and five platform package coordinates that form one indivisible Roundfix publication unit. Every coordinate must be eligible before any one is published.
_Avoid_: Package list, publish targets, artifact set

**Publication Preflight**:
A read-only eligibility check of the entire Release Set against registry truth for an exact target version. It stops publication unless every coordinate is proven eligible.
_Avoid_: Publish dry run, identity check, release plan

**Lagging Surface**:
A Behavior Surface that changed, appeared or disappeared since the base release without a covering skill change, an applicable changed Coverage Review or an uncovered declaration. It blocks a Release Plan that proposes a version ([ADR-0256](docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md)).
_Avoid_: Stale skill, undocumented commit

**Release Plan**:
A read-only classification of committed changes between a base release and a target revision that identifies the required semantic-version increment, proposes the next version, cites its evidence, and states whether explicit human approval or manual impact classification is required. A Lagging Surface blocks a plan that proposes a version without changing its decision state or proposed version ([ADR-0256](docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md)).
_Avoid_: Release execution, automatic release, version guess

**Release Plan Command**:
The support command that produces a Release Plan without editing release files, creating or pushing tags, publishing packages, or creating a GitHub Release.
_Avoid_: Release Command, publish command, cut-release command

**Doctor Command**:
The support command that diagnoses a repository and machine's readiness for Roundfix Runs — minimum-supported acpx, Adapter Readiness, Agent Selection Profile Readiness, Repository Skill Set, codex runtime hygiene, and whether Run storage is reclaimable. It reports the detected acpx version against the minimum, gives each check a next action, and mutates nothing. It evaluates Repository Skill Set readiness after, and independently from, Agent Selection Profile Readiness; unlike the Doctor Command, the Setup Command prepares the machine.
_Avoid_: Health check run, setup run, environment wizard

**Migrate Command**:
The support command that upgrades an older Run Database to the binary's schema version under the machine-wide write lock. It writes nothing when the Run Database is current, absent, or newer than the binary.
_Avoid_: Automatic migration, database downgrade, schema compatibility mode

**Spec Consistency Check**:
The read-only, pre-Run support command that compares a Spec's written citations, declarations, and cross-references. It uses an ADR horizon: `SC-ADR-RELATED` reports an ADR for a committed Spec only when the commit that added the ADR is an ancestor of the commit that added the Spec's `_prd.md`. Its Claim Receipt and Surface Transcript declaration gaps begin at the contract horizon, the commit that added the concrete-contract guide to the write-techspec skill. It reports consistency findings and never edits artifacts or emits a QA verdict. Its Glossary Declaration gap begins at the glossary horizon, the commit that added the glossary guide to the write-prd skill, and a Spec that carries a Glossary Declaration is checked at any age.
_Avoid_: QA gate, Spec validator, inference engine

**Claim Receipt**:
A source-and-quote pair in the same paragraph as an attribution to a decision record. The Spec Consistency Check proves that the verbatim quote occurs in the named source after whitespace normalization; it does not judge whether the quote supports the claim.
_Avoid_: Citation proof, supported claim

**Surface Transcript**:
A numbered TechSpec declaration of one command surface containing its command, standard output, standard error, and exit code. The QA gate reproduces it through the built product and compares the streams using the transcript's exact-match and controlled-variation conventions.
_Avoid_: CLI example, expected output prose

**Consistency Finding Severity**:
Each Spec Consistency Check finding is an `error` when the check locates both sides of a contradiction, or a `gap` when it surfaces a candidate it cannot settle. The `SC-*` diagnostic codes are stable and never renumbered once shipped.
_Avoid_: QA verdict, priority, confidence

**Metric Declaration Gap (`SC-METRIC-UNDECLARED`)**:
The Spec Consistency Check gap raised when a Spec's `_prd.md` does not declare numbered Success Metrics or an explicit `None.` with a reason.
_Avoid_: Missing metric, untraced metric

**Contract Declaration Gap (`SC-CONTRACT-UNDECLARED`)**:
The Spec Consistency Check gap raised when a Spec's `_techspec.md` does not declare numbered API Contracts or an explicit `None.` with a reason.
_Avoid_: Missing contract, untraced contract

**Archive Record**:
The small `<slug>.md` record an archive leaves under the resolved archive root. It carries the disposition, QA Report and verdict, promoted references, and `source_revision`; the removed Spec folder stays recoverable in Git at that revision. The History Sanitize Command also writes one for a Legacy Archive Folder; a folder without a QA Report receives the disposition `no-qa`, while a folder whose newest QA Report failed without an override receives `failed-qa` and keeps that verdict and report name. `failed-qa` is never a pass or an override (ADR-0248, ADR-0251).
_Avoid_: Archive folder, archive log, copied Spec

**Archive Advice**:
Jev's advisory classification of candidate files before the archive cut. It labels reusable knowledge, repository records or transient evidence and never gates the archive.
_Avoid_: Archive verdict, promotion decision, required evidence

**Archive Command**:
The support command that retires a Spec under the archive eligibility contract, writes its Archive Record, and removes the Spec folder while keeping it in Git at `source_revision`. Normal QA eligibility and a user-authorized QA Archive Override are distinct dispositions. Neither retirement nor an override fabricates Task completion or a passing QA verdict.
_Avoid_: Move command, retire run, cleanup command

**History Sanitize Command**:
The `roundfix history sanitize` command, a dry run unless `--apply --batch <n>`, that converts the existing history after the History Full Tag. `roundfix baseline update` runs the same planning for every convertible unit at once, while this command remains the path for reviewed batches (ADR-0254).
_Avoid_: Archive Command, history deletion, cleanup command

**Pending History**:
The units of a repository's existing history that the History Sanitize Command would plan, which the Managed Refresh includes in its plan and the upgrade notice names. `roundfix baseline update` plans these units with the Baseline Plan (ADR-0254).
_Avoid_: Active history, current history, unplanned archive

**Legacy Archive Folder**:
A Spec folder the Archive Command left under the archive root before archives wrote Archive Records (ADR-0248).
_Avoid_: Archive Record, active Spec, history root

**Refused Unit**:
A pending History Sanitize Command unit that cannot be converted. The plan and the batch name it with its reason on one line, leave it untouched, and do not count it toward `--batch <n>`. The update lists one without blocking the Baseline Plan or changing its state; a unit with uncommitted changes is refused (ADR-0251, ADR-0254).
_Avoid_: Converted unit, skipped silently, failed batch

**Lenient Legacy Reading**:
The reading of a Legacy Archive Folder's Task Graph that tolerates and names projection rows outside the graph and retired Task types. It also accepts a legacy list of maps `unproven`, one text line per map. It is never applied to an active Spec (ADR-0251, ADR-0254).
_Avoid_: Relaxed active reading, manifest rewrite, schema bypass

**Sanitize Batch**:
The next units one `--apply` can convert, delivered as one Pull Request with the repository gates green and revertible; Refused Units do not count toward the batch. The update converts every convertible unit as one batch in its own change (ADR-0248, ADR-0251, ADR-0254).
_Avoid_: Archive batch, migration wave, unreviewed cleanup

**Reduced History Entry**:
A retired Finding or Backlog Entry cut to its front matter, title, first paragraph and the revision holding its full text (ADR-0248).
_Avoid_: Summary, Archive Record, truncated history

**History Full Tag**:
The annotated `history-full` tag on the last commit before the first batch, which keeps every removed byte reachable. The update creates it at `HEAD` when absent, never moves it, and the operator pushes it with `git push origin history-full` (ADR-0248, ADR-0254).
_Avoid_: Lightweight tag, batch tag, archive stamp

**History Root**:
The single repository directory holding Archive Records, Reduced History Entries and retired ADRs, with the full bytes in Git. Retired material leaves the directory an Agent loads by default and stays greppable as the record of what was built and why; one path filter excludes the whole tree from review. Resolved in exactly one place, so no consumer carries a path of its own (ADR-0248).
_Avoid_: Archive folder, `_archived`, attic, trash

**History Relocation**:
One retired file that belongs under the History Root and is not there, carried in the Baseline Plan as an ordered ledger entry of source, destination, and content identity rather than file bytes. A relocation whose destination is occupied is a collision: it is refused by name, its siblings still move, and the repository is not current until it is resolved.
_Avoid_: Archive move, file migration, relocation payload

**Relocation Citation**:
A citation in a tracked file that resolves before a Baseline Plan's History Relocations and would not resolve after them. It is reported as a warning the Plan Digest binds; planning never rewrites it, and apply writes the same files.
_Avoid_: Broken link, dangling reference, link check

**GC Command**:
The support command that reclaims Run storage: it prunes the Run Event Journal and artifact directory of terminal Runs older than the Journal Retention window, removes orphaned run artifact directories, runs the Run Retention Sweep without a budget, and reports what it freed. A Run it already emptied is not counted again. It never touches Active Runs or active-run locks, and it removes a `runs` row only through Run Retention (ADR-0255).
_Avoid_: Clean command, vacuum, purge

**Journal Retention**:
The configured age window after which a terminal Run's Run Event Journal and artifact directory become eligible for pruning. Active Runs are never eligible; a retention of zero keeps everything. See ADR-0033. It is the inner window: Run Retention removes the whole Run past the outer window, so a Journal Retention longer than the Run Retention has no effect past it (ADR-0255).
_Avoid_: Log rotation, TTL, expiry

**Run Retention**:
The User Config window, `store.run_retention_days` of 7, 15 or 30 days and 30 by default, after which a terminal Run leaves the Run Database whole: its row, its Run Event Journal, its Agent Selection records, its token usage and its artifact directory. A Run that is not terminal, that a Delivery Queue references, whose Run Worktree still exists, or whose artifact directory sits under an unproven root is kept (ADR-0255).
_Avoid_: Journal Retention, TTL, purge

**Run Retention Sweep**:
The step that applies Run Retention and then returns freed pages to the filesystem by incremental compaction. It runs at most once a day at the start of a Run or a Delivery Queue under a two-second budget, continuing at the next start when the budget runs out, and without a budget in the GC Command (ADR-0255).
_Avoid_: cleanup job, cron, vacuum

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
The setup-owned record of the selected profile, modules, decisions, and managed artifacts that reproduces and audits a repository's Context-Driven Baseline. Every applied Baseline change republishes it, so its recorded digests describe the bytes on disk rather than the bytes adoption wrote.
_Avoid_: Project Config, lock file, generated guide

**Managed Region**:
The byte range one setup-owned marker pair delimits inside an instruction carrier, identified by its managed artifact identity. Its content belongs to the Context-Driven Baseline, not to the repository holding it; the bytes outside every marker pair belong to the repository and no Baseline operation changes them.
_Avoid_: Managed file, generated block, setup section

**Managed Refresh**:
The preservation mode that regenerates only Managed Regions, leaves every other byte identical, requires no Source Baseline and no Decision Plan, and takes no root backup because the plan's preimages carry the preservation proof. It also plans the Pending History, binds it into the Plan Digest and converts it on approval. Its preservation proof still covers the instruction carriers. It converges: a second run against an unchanged catalog reports the repository current and proposes no change (ADR-0254).
_Avoid_: Baseline Readoption, reapply, overwrite

**Unrecorded Managed Region**:
A Managed Region whose bytes are not the bytes the adopted Setup Manifest recorded. It is classified and reported by path, managed identity, and the lines a refresh removes — never presumed damaged and never silently replaced. Only a managed identity appearing more than once in one carrier blocks, because that refresh has no defensible target.
_Avoid_: Hand-edited marker, corrupted block, drift

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

**Format Command**:
The Project Config command `verification.format` that the Daemon runs, outside the Agent sandbox, over the files under a Spec's `qa/` directory before the repository Verification precondition (an imported pass) and before the QA Report commit, restoring the original bytes when it fails (ADR-0249).
_Avoid_: Repository Verification, Baseline formatter, whole-tree formatter

**Late Dependency**:
A Task in a completed QA gate's dependency closure that the Task Graph did not hold when the newest QA Report was committed; `roundfix reopen` returns the gate to `pending` over it whatever its status (ADR-0249).
_Avoid_: Stale dependency, missing dependency, unresolved dependency

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
A read-only JSONL projection of one Run's Run Event Journal, selected by an explicit Run ID and optionally followed until the Run reaches a terminal outcome. Its stable Supervisor filters are task status, Batch boundary, Verification verdict, and terminal outcome. The command skips a record it cannot project with a warning and continues with the next journal entry.
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


**Pre-PR Review Policy**:
The repository's choice of Codex, Claude, CodeRabbit, or no reviewer before a
Pull Request is opened. An enabled review covers the current candidate, while
explicit none records an intentional omission that leaves QA and required
checks binding; provider failure is not an omission.
_Avoid_: Review mode, reviewer profile, automatic review

**Pre-PR Review Provider**:
The selected agent or service that performs review under the Pre-PR Review
Policy. An Agent Selection Profile supplies an agent provider's runtime and
model; explicit none is a policy choice with no provider.
_Avoid_: Review Source, Pre-PR Review Policy, reviewer model

**Pre-PR Review Record**:
The record of one Pre-PR Review, kept per checkout. Its verdict is read from
the reviewer's final message.
_Avoid_: Review Source Evidence, shared review state, reviewer transcript


**QA Archive Override**:
An explicit user request or applicable prior authorization to archive a covered
Spec despite its QA prerequisite being unmet. The recorded exception preserves
actual QA and Task evidence while retaining non-QA completion, source integrity
and separate delivery gates.
_Avoid_: QA pass, QA bypass, Implement Clean

**Glossary Declaration**:
The `## Glossary` section of a Spec's PRD or TechSpec that lists each domain term the Spec adds to or changes in the glossary, and each bolded phrase it declares not a domain term, or `None.`. A Task that declares the glossary file and names the term in its Verification writes each listed term.
_Avoid_: Vocabulary Contract, term list, glossary candidates

**Glossary Gap**:
A domain term a Spec introduces that the glossary does not carry as required: a bolded phrase the glossary lacks and the Glossary Declaration does not cover, a declared term no Task plans to write, or a declared term still missing after its Tasks completed. The Spec Consistency Check reports it, and the QA gate and the Archive Command refuse a Spec that has one.
_Avoid_: Missing term, undocumented token, stale glossary

**Light Tier**:
The dispatch tier that runs a Task on an open model through OpenCode and OpenRouter when its complexity is low, its type is not qa, and it declares no Governed Path; every other Task runs on the standard tier. A light Task whose first Verification fails escalates its one feedback turn to its category's Preferred Selection.
_Avoid_: Cheap model, Jev routing, light model profile

**Light Spend Log**:
The per-UTC-month record in Roundfix Home of the cost each Light Tier session reports, summed across every repository and compared with the Light Tier's monthly ceiling. A month that reached the ceiling, or whose log cannot be read, runs light Tasks on their category's profile.
_Avoid_: Judge Log, OpenRouter usage, spend cache

**Jev Ceiling**:
The monthly spending ceiling of the Jev judge, read from User Config only and compared with the Judge Log's sum across every repository on the machine; when it is unset a built-in default applies. The judge's key must carry a monthly credit limit no higher than it.
_Avoid_: Judge budget, key limit, Project Config ceiling

**Stage Key**:
The Roundfix environment variable an OpenRouter stage reads first, the judge's or implementation's own key, before the shared Roundfix OpenRouter key; the generic OpenRouter variable is never read. Each stage records the variable it used by name, never by value.
_Avoid_: API key value, generic OpenRouter key, shared key

**Tested Base**:
The default-branch commit a failed required check actually tested, read from the check job's annotation. A failure whose Tested Base does not contain the current default-branch tip is stale, and the Delivery Queue re-runs it instead of parking the item.
_Avoid_: PR base, merge base, Delivery Base

**Merge Evidence**:
Proof that a Spec was delivered: the default-branch head holds its archived PRD, and the Delivery Commit that added that file is not reachable from the Run Branch or Item Branch being proven. With it, reconcile releases that Spec's terminal Runs and Item Branches without comparing content.
_Avoid_: Content proof, merged flag, merge commit

**Delivery Commit**:
The default-branch commit that added a Spec's archived PRD. Merge Evidence and a Delivery Retry that records a merge made outside the queue both read it.
_Avoid_: Merge commit, archive commit, release commit

**Item Branch**:
The branch a Delivery Queue item builds its candidate on, named `roundfix/deliver-<slug>-<16 hex>`. Once no live item names it, Merge Evidence releases it, with its worktree only when that worktree is clean.
_Avoid_: Run Branch, PR Head Branch, deliver branch

**Network-Denied Row**:
An outside-evidence QA row blocked only because the Run sandbox denied network access, recorded with the host it could not reach. Like the pre-PR Pull Request row, it never decides a qualifying partial; whoever needs that proof declares it under Unreachable Acceptance.
_Avoid_: Failed outside evidence, skipped row, environment pass

**Pre-PR Review Command**:
The command that obtains the Pre-PR Review Provider's review of the current candidate before a Pull Request opens, reads the verdict from the reviewer's final message and keeps the Pre-PR Review Record. A finding parks delivery only after it is validated.
_Avoid_: Review Run, PR feedback Watch, review command alias
