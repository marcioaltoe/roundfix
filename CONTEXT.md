# Roundfix

Roundfix picks up Work Items — Review Issues from pull request reviews and Tasks from Spec Task Graphs — resolves them through the user's selected ACP Runtime, and pushes only when nothing unresolved remains. This glossary defines the product language for that loop; the terms under Adopter decisions describe repositories the Context-Driven Baseline configures, not Roundfix itself.

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

**Absence Evidence**:
A TechSpec record of the search that proves each newly declared component absent, or the reason no component is new ([ADR-0265](docs/adr/0265-a-spec-proves-what-it-claims-and-every-active-spec-is-held-to-it.md)).
_Avoid_: Unverified absence, component inventory

**ADR Alternatives**:
A TechSpec line that relates its design to the alternatives rejected by each ADR in its active ADR obligations ([ADR-0265](docs/adr/0265-a-spec-proves-what-it-claims-and-every-active-spec-is-held-to-it.md)).
_Avoid_: Decision summary, replacement decision

**Coverage Cross-Read**:
A TechSpec record of reading each PRD unit against its Coverage Map and naming any unit delivered in another shape ([ADR-0265](docs/adr/0265-a-spec-proves-what-it-claims-and-every-active-spec-is-held-to-it.md)).
_Avoid_: Coverage map, traceability table

**Consumer Declaration**:
A pending non-QA Task's record of the suites that consume each Shared Contract it changes and where those suites run ([ADR-0265](docs/adr/0265-a-spec-proves-what-it-claims-and-every-active-spec-is-held-to-it.md)).
_Avoid_: Test list, regression checklist

**Shared Contract**:
A route, public signature, persisted format or report that code outside a Task's slice consumes ([ADR-0265](docs/adr/0265-a-spec-proves-what-it-claims-and-every-active-spec-is-held-to-it.md)).
_Avoid_: Internal helper, implementation detail

**Integration Checkpoint**:
A non-QA Task between a contract-changing Task and the QA gate whose Verification runs the consuming suites ([ADR-0265](docs/adr/0265-a-spec-proves-what-it-claims-and-every-active-spec-is-held-to-it.md)).
_Avoid_: Final QA, follow-up task

**Run**:
A durable attempt to drive one target's Work Items — an Open Pull Request's Review Issues or a Spec's Tasks — to a terminal outcome. One Active Run is allowed per target: (Head Repository, PR Head Branch) for review work, (repository, spec slug) for spec work.
_Avoid_: Session, execution, job

**Review Run**:
A Run over an Open Pull Request's Review Issues — a Fetch Run, a resolve Run or a watch Run. It executes in the user's clean tracked checkout on the PR Head Branch and creates no Run Worktree ([ADR-0141](docs/adr/0141-review-runs-execute-in-a-clean-tracked-user-checkout.md)).
_Avoid_: PR run, review job, Review Source run

**Spec Run**:
A Run the Implement Command starts over one Spec's Task Graph. It executes in a Run Worktree on a Run Branch, schedules Tasks by Wave and ends at its authored QA gate when the graph declares one; the guides also call it an Implement Run ([ADR-0023](docs/adr/0023-runs-execute-in-per-run-worktrees.md), [ADR-0025](docs/adr/0025-spec-runs-schedule-ready-tasks-concurrently.md)).
_Avoid_: Implement job, task run, graph execution

**Implement Run**:
The Spec Run as the Implement Command and the agent guides name it; the two terms name one Run.
_Avoid_: Spec execution, implement job, implement session

**Fetch Run**:
A short Review Run that runs Branch Integrity Preflight, fetches Review Source issues, and persists markdown artifacts from the user's checkout without starting an Agent, creating a Run Worktree, committing, or pushing.
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

**Proven Signal**:
A signal sent only to the process whose identity was recorded when its PID was learned, or to its process group after proving the leader or the continued presence of live members. A different identity means the recorded process ended and receives no signal ([ADR-0266](docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md)).
_Avoid_: Safe kill, PID kill, best-effort signal

**Signal Exit**:
The end of a child process that a signal caused instead of an exit status, which Roundfix and its test helpers report by the signal's name and number and whether a core was dumped, never as exit code `-1` or `1`.
_Avoid_: exit -1, crash, negative exit code

**Helper Lifeline**:
A test helper's connection to its owning test binary, whose loss ends the helper even when test cleanup cannot run ([ADR-0266](docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md)).
_Avoid_: Owner PID probe, cleanup promise

**Helper Residue**:
A live process carrying a test binary's ownership mark after that package's tests and cleanup grace have ended. It is a test-suite isolation violation, distinct from Process Residue, which concerns a process that outlived its Run ([ADR-0266](docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md)).
_Avoid_: Process Residue, orphan helper, zombie

**Force Stop**:
The Stop Command mode that proves owner identity, cancels registered Agent Sessions, terminates the recorded owning process, and completes the Run as Stopped only after owner exit is proven. A proven identity mismatch always refuses; only the explicit `--owner-identity-unreadable` last-resort flag permits PID-only termination after the host reports the identity unreadable.
_Avoid_: Lock release, best-effort stop, orphan reclamation

**Stop Command**:
The `roundfix stop` command that records a Stop Request for an Active Run, which then stops after the current Work Item settles, or performs a Force Stop with `--force`. It is the only way a Detached Run is ended ([ADR-0022](docs/adr/0022-stop-requests-travel-through-the-run-database.md), [ADR-0028](docs/adr/0028-operational-runs-can-detach-from-their-caller.md)).
_Avoid_: Kill command, cancel command, abort

**Assigned Repair**:
A repair a QA gate's own Task file instructs it to make, as opposed to a finding it observes and reports. The gate performs an Assigned Repair, verifies it, and records what it performed; it may write only the paths its Task names. Anything the Task did not assign stays reported rather than performed. In code and in a report's audit trail the same idea appears lowercase, as an assigned repair.
_Avoid_: Auto-fix, gate edit, remediation

**Governed Path**:
A repository path the tooling-authority rules bind — the configuration, scripts, ignore files, plugin declarations, and version pins of linters, formatters, typecheckers, test runners, architecture checkers, build tools, package managers, and code generators. The changed-path audit judges a Task commit only against its governed paths; an ordinary source, test, or documentation file is not governed and needs no grant. The declared set is held to the record: every path any authorization has ever bounded stays governed.
_Avoid_: Protected file, restricted path, tooling file

**Delivery Base**:
The merge base of the audited head and the Default Branch Head. The pre-PR review diffs the candidate from it. The mechanical stage reads a Spec's Task commits and its grant from it. A QA requirement that holds a file byte-identical compares it with the Delivery Base, never with a named commit ([ADR-0295](docs/adr/0295-reachability-reads-the-fetched-default-branch-and-origin-names-it-when-origin-head-is-unset.md)).
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

**Stale Run Owner Record**:
The owner PID and start identity a terminal Run still records after the live process at that PID has presented a different start identity, which proves the recorded owner ended and the PID now names another process. The Reconcile Command reports it, and `--apply` releases it by clearing the Run's owner PID and identity and journaling them, without sending any signal ([ADR-0313](docs/adr/0313-a-proven-owner-identity-mismatch-releases-a-terminal-runs-owner-record-without-a-signal.md)).
_Avoid_: Orphan PID, dead owner lock, staging owner record

**PR Head Branch**:
The branch on GitHub that supplies the pull request's head commits.
_Avoid_: Local branch, checkout branch

**Head Repository**:
The GitHub repository that owns the PR Head Branch.
_Avoid_: Base repository, local checkout

**Work Item**:
One unit of resolvable work that Roundfix picks up and drives to an outcome: a Review Issue in a Review Run, a Task in a Spec Run.
_Avoid_: Ticket, job, to-do

**Backlog Entry**:
Typed intent for what to do next, with a lifecycle of `open`, `promoted` to a named Spec, or a terminal status (`declined`, `deferred`, `done`, `deprecated`, `superseded`, `closed` or `cancelled`); a terminal entry is deleted and Git keeps its text, while `declined` and `deferred` carry a reason ([ADR-0278](docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md)). Its types express intent through the Conventional Commits vocabulary — `feat`, `fix`, `perf`, `refactor` — and it is never evidence; a `feat` entry is upstream raw material, never the `write-idea` artifact.
_Avoid_: Finding, idea artifact, untyped suggestion

**Cross-project Message**:
A self-contained message a session sends to another project's ai-memory inbox with `memory_message_send` when it observes work that project owns. The receiving project reads its messages oldest first, records each in a Spec, a Backlog Entry or a Finding, or names it when it only informs, and only then consumes it with `memory_message_pop`; without ai-memory the sender gives the text to the user ([ADR-0330](docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md)).
_Avoid_: Inbox Entry, capture, note

**Rollup**:
A Finding of `rollup` kind that consolidates related Findings as declared members and retires by deletion when its work is resolved. Its members resolve as unresolved Findings or through Git, which keeps their retired text ([ADR-0278](docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md)).
_Avoid_: Summary, digest, archive manifest

**Triage**:
The destination-project act that converts one pending Cross-project Message into its contract-true artifact, committed by the destination alone.
_Avoid_: Capture, import, synchronization

**Follow-up Command**:
The read-only `roundfix follow-up` command that lists a repository's Follow-up Items in Follow-up Priority order beside its work in progress, names each source it checked or missed, reports `nothing_to_do` when it finds neither, and with `--status` shows the work in progress, the pending Specs in delivery order, the next release and what waits on the maintainer. It never writes, and the judge never changes its exit code ([ADR-0319](docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md)).
_Avoid_: Triage command, dashboard, inbox sync

**Follow-up Item**:
One pending unit the Follow-up Command lists: an ai-memory message the caller hands in, an open Backlog Entry, an unresolved Finding, a note in `docs/_inbox/` or an active Spec, with its kind, reference, type, date, title and the project that sent it.
_Avoid_: Work Item, ticket, to-do

**Follow-up Priority**:
The order of Follow-up Items in four classes, `unblocking`, `health`, `feature` and `other`, computed in code from four Jev yes/no answers per item, with small items and then the oldest first inside a class; without an answer an item's declared type places it, and never in `unblocking` ([ADR-0319](docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md)).
_Avoid_: Jev ranking, severity, score

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
The recorded fix or evidence-backed dismissal of a finding from a pre-PR review, tied to the head the reviewer examined.
At the same repository, base, head, and provider, the findings verdict is reused without asking the reviewer; only exact evidence-backed dismissals of every finding produce `findings-dismissed`, while a fix requires a fresh review of its changed head. An evidence-backed dismissal may be recorded at a later `HEAD` that descends from the reviewed head when the finding's anchored lines are unchanged; it stays tied to the reviewed head and records `currentHead` ([ADR-0279](docs/adr/0279-status-reports-a-queue-owner-proven-gone-and-a-dismissal-survives-a-fix-that-leaves-its-anchor.md)).
_Avoid_: Review Issue resolution, unrecorded operator decision, finding deletion

**Spec**:
One Spec's artifact set produced by the spec workflow under the Spec Root: the Idea when the route starts there, the PRD, the TechSpec, the Task Graph, the Task files, the authorization record `_authorization.md`, the Glossary Declaration in the PRD or TechSpec, the QA evidence under `qa/`, and the adopted sources under `references/` with provenance recorded in `references/_index.md`.
_Avoid_: Feature folder, epic, project

**Corrective Spec**:
A new Spec, with its own authorization and QA gate, that corrects a finding on an archived Spec; the archived Spec is never edited in place.
_Avoid_: archive edit, inherited authorization, Corrective Task

**Corrective Task**:
A Task appended to an active Spec's Task Graph, before its QA gate, to fix a finding on that Spec before it is archived; an archived Spec gets a Corrective Spec instead.
_Avoid_: Corrective Spec, follow-up Task

**Cause Class**:
The one class from a closed list that `roundfix runs causes` assigns to a failed Verification attempt or a Corrective Task: `scope-or-authorization`, `shared-section-contract`, `repository-convention`, `implementation-defect` or `environment`, or `unclassified` when no signature matches. It is assigned by signature, never by a model judgment ([ADR-0214](docs/adr/0214-why-a-verification-failed-is-classified-by-deterministic-signatures.md)).
_Avoid_: failure category, root cause, model verdict

**Classification Signature**:
An ordered entry of the embedded signature table, with an identifier, a Cause Class, the one evidence source it reads and a regular expression; the first entry that matches classifies the item ([ADR-0214](docs/adr/0214-why-a-verification-failed-is-classified-by-deterministic-signatures.md)).
_Avoid_: heuristic, rule, pattern

**Supersession**:
A durable lifecycle record that names the Spec which delivered another Spec's content. The Archive Command accepts it as proof of completion for the superseded Spec when that Spec has no Task Graph. A cancelled Spec never delivers another Spec's content ([ADR-0262](docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md)).
_Avoid_: Replacement, duplicate, hand-written disposition

**Cancel Command**:
The `roundfix cancel <slug> --reason <text>` command that retires an active Spec which will not be delivered: it writes an Archive Record with the disposition `cancelled`, performs Source Return and removes the Spec folder, which stays in Git at `source_revision`. It refuses while an Active Run or an unfinished Delivery Queue item holds the Spec, an active Spec requires it or Git ignores a file under the Spec folder, never changes a Task status and never commits ([ADR-0262](docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md), [ADR-0269](docs/adr/0269-archive-and-cancel-refuse-ignored-spec-files-and-finish-a-failed-removal-forward.md)).
_Avoid_: Delete command, abandon, archive override

**Ignored Spec File**:
A file under an active Spec folder that Git ignores, so no commit holds it. The Archive Command and the Cancel Command refuse while one exists ([ADR-0269](docs/adr/0269-archive-and-cancel-refuse-ignored-spec-files-and-finish-a-failed-removal-forward.md)).
_Avoid_: Untracked change, scratch file, ignored repository file

**Source Return**:
The Cancel Command's hand-back of a cancelled Spec's adopted sources: a Backlog Entry returns to `docs/backlog/` with `status: open` and no Spec, a Finding returns to `docs/findings/` with `status: pending` and a dated addendum, and every other byte is kept. Other adopted sources stay in Git with the folder ([ADR-0262](docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md)).
_Avoid_: Unadoption, restore, re-triage

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
The required classification that routes a Task to its Agent Selection Profile. The valid values are `backend`, `frontend`, `data`, `infra`, `docs`, `test`, `chore` and `qa`; use the dominant implementation surface when a Task crosses more than one. `qa` names the Spec's authored QA gate, the terminal Task that depends on every leaf, and no other Task carries it ([ADR-0091](docs/adr/0091-the-qa-gate-is-a-task-node-of-its-own-type.md)).
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
The qa-gate evidence report written to a Spec's QA directory, carrying a machine-readable verdict, its Auditing Binary as `auditing_binary`, `auditor_staleness`, and the public-row binary's `--version` line as `user_flow_binary`, plus `rows_blocked_environment`, `rows_blocked_finding`, and `rows_blocked_declared` counts in its frontmatter. The auditor fields are Daemon-owned: the gate keeps the seeded `auditing_binary` and `auditor_staleness` values. A report recording a Precondition Refusal carries `rows_blocked_precondition` beside those counts, plus the `precondition_check` and `precondition_reason` keys that name the refusal; a gate that reached its matrix writes none of those three. A report the Daemon closes without an Agent turn, a Precondition Refusal or a blocking mechanical result, records the audited commit as `build`, as the gate's full report does; the pending report the Daemon seeds for the gate carries none. The pre-PR pull request row, recorded as `blocked (environment: no open Pull Request)` with the pull request row named in its provenance, never decides a qualifying partial and needs no Unreachable Acceptance declaration. A `pending` verdict is never accepted, a report that records no QA row is refused, and a report whose front matter is empty or duplicated is unreadable and refused; only the newest report in the directory is read by a later run's mechanical stage, so a superseded report blocks nothing.
_Avoid_: Test report, QA log

**Auditing Binary**:
The Daemon's Roundfix binary that runs the mechanical stage, distinct from the tree it audits. It carries the version, build commit, and build time (the last two may be empty for a released build), and a QA Report records its formatted identity as `AuditingBinary` / `auditing_binary`; `auditor_staleness` compares its build commit with the Delivery Base, and a stale Auditing Binary is published as a warning while the gate proceeds.
_Avoid_: Audited binary, audited tree, ambiguous build

**Precondition Refusal**:
A QA gate stop at a check that runs before the matrix is built — a strict Spec check being the usual one — recorded as the gate's entire QA Report rather than as a missing one: verdict `fail`, `rows_blocked_precondition` counting the stop, `precondition_check` naming what was checked, `precondition_reason` carrying every refusing code and the sentence beside it, and one terminal row `0 | blocked | precondition` in the Results table. The gate measured no requirement, so row `0` records the refusal itself instead of a result, and the prose justifying it is written as a list rather than a second table. An empty Results table is not the alternative: it refuses every later run on the report instead of on the Spec, and prescribes a repair — materialize every planned row — that a run which never built a matrix cannot perform. It also records `build`, the commit the gate audited, as a full QA Report does, so a reader can tell a refusal an earlier Run left on the branch from the current Run's result; a refusal written before this rule keeps its bytes.
_Avoid_: Empty report, failed gate, skipped QA, precondition error

**External Acceptance Evidence**:
Evidence for at least one named acceptance row that originates outside the Spec's own artifacts — a real repository, a measurement, or published literature — with its origin recorded. A row whose external evidence cannot be obtained is recorded as blocked with its reason. It exists because a rubric and the requirement it measures, written by one author from one premise, confirm rather than test.
_Avoid_: Self-authored rehearsal, internal fixture, passing suite

**Unreachable Acceptance**:
The `## Unreachable Acceptance` declaration a Spec authors with its QA rows, naming each acceptance row whose evidence needs the network, an open pull request or a credential the repository cannot hold. A row blocked for a declared reason can still qualify a declared `partial`; the pre-PR pull request row, a Network-Denied Row and a Network-Prohibited Row need no declaration.
_Avoid_: Skipped row, waived criterion, known limitation

**ACP Runtime**:
A local coding runtime that Roundfix launches through the user's installed tool and authentication setup using Agent Client Protocol stdio. Roundfix supports Codex through `codex-acp`, Claude through `claude-agent-acp`, OpenCode through `opencode acp` and, as an opt-in that enters no built-in profile, Cursor through `cursor-agent acp` ([ADR-0217](docs/adr/0217-a-cursor-selection-names-the-advertised-model-value-and-the-login-stays-the-maintainers.md)); command overrides remain a stdio escape hatch for local testing.
_Avoid_: Review Source, review provider

**Agent Model**:
A runtime-specific model choice that Roundfix explicitly assigns to every Agent Session.
_Avoid_: ACP Runtime, Agent, free-form model override

**Agent Work Category**:
The routing key Roundfix uses to resolve an Agent Selection Profile: `general`, one Task Type (`qa` included), or `review`.
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

**Proof Setup Deadline**:
The wall-clock bound on one attempt of an Exact Agent Selection Proof's
disposable Agent Session setup, set by `defaults.profile_proof_setup_timeout`
(default `30s`). An attempt that exceeds it is retried after 5 seconds and
again after 10 seconds, and a third timeout is classified `temporary` ([ADR-0298](docs/adr/0298-a-timed-out-profile-proof-gets-a-third-attempt-and-its-setup-deadline-is-configurable.md)).
_Avoid_: Run Budget, test deadline, arbitrary timeout

**Runtime Catalogue**:
What an ACP Runtime advertises before Roundfix asks it to apply a specific Agent Selection, established from a disposable Agent Session ensured without the requested model so the answer cannot be contaminated by the question. Membership in it decides refusal where no adapter refuses first; where an adapter refuses, its own message stands ([ADR-0147](docs/adr/0147-the-adapters-refusal-fires-first-and-the-catalogue-is-the-net.md)).
_Avoid_: Model list, capability payload, advertised options

**Selection Encoding**:
How one Agent Selection's reasoning effort is represented in the advertised ACP controls, decided during proof and reported by readiness surfaces. `independent` assigns a separate advertised reasoning option. `model_variant` selects an advertised model identifier that already carries the effort. `model_managed` is an empty reasoning effort on an adapter advertising no reasoning control at all. `runtime_managed` is an empty reasoning effort on an ACP Runtime whose advertised reasoning control Roundfix declines to assign, so the observed value is the Agent Model's own and is not proof-relevant. `runtime_deferred` is a non-empty reasoning effort on an ACP Runtime that advertises the control only once the Agent Session holds a Queue Owner, so proof verifies the value is advertised and the Run applies and observes it after an inert setup turn and before the first work turn.
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
The first Agent Selection that Roundfix proves and attempts for an Agent Work Category.
_Avoid_: Default model, primary runtime, recommendation winner

**Fallback Chain**:
The non-empty ordered Agent Selection list Roundfix proves with the Preferred Selection before a Run and may activate after notifying the user and Supervisor that the preceding selection failed before Agent work began. One exception reaches past Agent Work Started: a Lost Rollout or Provider System Error before the First Handoff, or in a Task of type `qa` before its report, lets the next Fallback Selection take the Task, because the runtime discarded that turn before it could be handed back ([ADR-0245](docs/adr/0245-a-lost-rollout-is-runtime-infrastructure.md), [ADR-0274](docs/adr/0274-a-turn-that-did-not-finish-is-not-accepted.md)).
_Avoid_: Dynamic fallback, silent retry, unproven selection

**Default Agent Model**:
The concrete Agent Model that Roundfix selects for an ACP Runtime when the user supplies no override. It never inherits the runtime's local model configuration.
_Avoid_: Runtime default, automatic model, local Agent default

**Default Reasoning Effort**:
The runtime-specific reasoning level Roundfix assigns when the value is non-empty. An empty value is valid and means the Agent Model manages reasoning, so Roundfix assigns no reasoning option. It never inherits the runtime's local reasoning configuration.
_Avoid_: Local Agent reasoning, automatic reasoning, reasoning hint

**Model Catalog**:
The ordered set of known Agent Models that Roundfix offers for one ACP Runtime during Interactive Input. Its Default label resolves to the Default Agent Model, while non-interactive interfaces may supply a custom value.
_Avoid_: Global model list, model allowlist

**Recommended Profile**:
The dated Agent Selection Profile that Roundfix recommends for one Agent Work Category, from which built-in profiles derive. It never selects, routes, or changes a configuration by itself.
_Avoid_: Model Recommendation Ranking, model router, benchmark policy, automatic selection

**Fallback Selection**:
The next configured Agent Selection in a profile's Fallback Chain. Roundfix proves it before the Run, emits a notification before activation, and may switch ACP Runtime automatically only while Agent work has not begun or under the Fallback Chain's Lost Rollout and Provider System Error exception.
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

**Subsumed Test Selector**:
The Spec Consistency Check error (`SC-VERIFY-SUBSUMED-SELECTOR`) raised when a pending non-QA Task's Verification passes a test selector, a pytest node id such as `tests/test_a.py::test_new`, together with a collection argument of the same invocation that contains it: its file, a directory above it, `.`, or a node id that prefixes it. pytest 9 runs the containing argument whole and drops the selector without an error, so a test that does not exist yet goes unnoticed and the command passes before the work ([ADR-0315](docs/adr/0315-a-test-selector-another-argument-of-its-invocation-contains-is-refused-at-authoring.md)).
_Avoid_: Duplicate test path, redundant node id, shadowed test

**Unobserved Verification**:
The Run Event classification (`verification_unknown`) recorded when a Verification command's verdict could not be observed — a timeout, a partial execution, or a runner error that is not the command's answer. It is distinct from a command that ran and exited non-zero, which has a verdict, and it keeps "the work is wrong" separable from "we did not find out".
_Avoid_: Flaky test, transient failure, retry

**Unsupported Citation**:
The Spec check finding (`SC-CITATION-UNSUPPORTED`) raised when an artifact attributes a subject to a decision record that the record's own text does not carry. It is distinct from an unlisted citation, which asks whether a record was named; this one asks whether the claim about it holds, and reports both texts so a maintainer settles it by reading rather than trusting the checker.
_Avoid_: Broken link, missing ADR, stale reference

**Agent Work Started**:
The status (`agent_work_started`) marking the first Agent output that could have changed something, published once per Agent Session. It is the boundary after which a Fallback Selection may no longer switch ACP Runtime, because a second runtime would inherit state the first one built. Preparing or activating a Session does not reach it. Only the Fallback Chain's Lost Rollout and Provider System Error exception reaches past it.
_Avoid_: Session opened, prompt sent, turn started

**Selection Failure**:
The outcome (`agent_selection_failed`) of a turn that ended before any Agent output, because the selected ACP Runtime would not serve it — exhausted quota, failed authentication, an adapter that will not start. It is a failure of the selection rather than of the work, so the Fallback Chain stays eligible.
_Avoid_: Batch failure, agent crash, runtime error

**Hook Refusal**:
The outcome (`hook_refused`) of a Task commit a repository commit hook rejected after the authoritative Verification had already passed. It is a repository misconfiguration rather than a Task failure — a commit hook may never be stricter than the Verification — so the Task stays `completed`, its verified work stays staged in the surface that produced it, and the Run reports the refusing hook, its exit code, its output, and the Settle Command that recovers the work.
A Delivery Queue commit — the archive commit, the Pull Request Record Commit or the derived merge commit — that a hook refuses is a Hook Refusal too: the item parks `delivery-error` with a reason naming the hook, the refused subject and the hand commit that recovers it, and the queue never bypasses the hook ([ADR-0302](docs/adr/0302-a-queue-commit-names-the-spec-by-number-and-a-hook-refusal-parks-with-the-hand-commit.md)).
_Avoid_: Commit failure, pre-commit error, verification failure

**Mechanical Refusal Code**:
The stable token a pre-QA mechanical check emits when it can refuse from written declarations and repository facts alone, before any Agent Session opens. Four exist: `QA-AUTH-PATHS` (a tooling change outside its authorization's bounded paths), `QA-CONSEQUENT-ORDER` (a consequent fix folded into or ordered before the change that caused it), `QA-REPORT-SHAPE` (a QA Report missing a required structural element), and `QA-EVIDENCE-PATH` (an evidence path a report names but the repository does not carry, or a tracked evidence file that is a gitlink or a non-image binary, as defined by [ADR-0291](docs/adr/0291-qa-evidence-holds-text-and-images-and-the-qa-report-commit-refuses-any-other-binary.md)). The code is the durable name a reader and a later detector both use; the human sentence beside it may be reworded, the token may not. The `QA-PRECONDITION` row a refused mechanical stage writes as its terminal row is a row label beside these codes, not a fifth code ([ADR-0132](docs/adr/0132-a-refused-gate-records-the-refusal-as-its-terminal-row.md)).
_Avoid_: lint code, error code, check name

**Grant Refusal Code**:
The stable token the Spec checker emits when an authorization boundary refuses. `SC-TOOLING-UNAPPROVED` means a Spec claims tooling authorization its cited record does not sustain because the record is proposed, undated, withdrawn, or consumed by another Spec. A record that honestly declares its mutations proposed, claiming no operative grant, is accurate rather than defective and refuses nothing. Like a Mechanical Refusal Code, the token is the durable name a reader and a later detector both use; the sentence beside it may be reworded, the token may not.
_Avoid_: authorization error, permission code, trust error

**Agent Session**:
The acpx-backed session owned by one Work Item or action. Each Task in an Implement Run owns a Task Type-selected Agent Session, the authored QA gate owns a separate `qa` Agent Session, and a Review Run uses a review-selected Agent Session; effective selection and fallback attempts are persisted for that owner.
_Avoid_: ACP session, chat, conversation, thread

**Full-Access Mode**:
The adapter mode an Agent Session is switched to when full access is requested: the first of the runtime's full-access candidates that the session advertises in its `mode` config option ([ADR-0293](docs/adr/0293-a-full-access-request-selects-the-mode-the-agent-session-advertises.md)).
_Avoid_: sandbox mode, YOLO mode

**Lost Rollout**:
The runtime infrastructure failure (`rollout_lost`) in which the ACP Runtime cannot find the persisted record of the Agent Session it is asked to continue. The Daemon recovers it in a new Agent Session without repair or a counted queue retry ([ADR-0245](docs/adr/0245-a-lost-rollout-is-runtime-infrastructure.md)).
_Avoid_: Session loss, missing prompt, retryable Task failure

**Provider System Error**:
The runtime infrastructure failure (`provider_system_error`) of a turn whose stream reports that the model provider failed it — the turn's last `threadStatus` is `systemError` — whatever its stop reason. The Daemon recovers it like a Lost Rollout, without a retry on the same selection ([ADR-0274](docs/adr/0274-a-turn-that-did-not-finish-is-not-accepted.md)).
_Avoid_: Provider failure, capacity reply, Lost Rollout

**Agent Background Task**:
A command an Agent starts through its ACP Runtime's background facility, such as a Claude Code background Bash command or Monitor, which keeps running after its tool call returns. Roundfix refuses them in Run sessions and records each one a turn starts.
_Avoid_: Background job, background command, detached task

**Unfinished Turn**:
The outcome (`turn_unfinished`) of a QA turn that ends with its seeded report at `verdict: pending` after starting an Agent Background Task. The Task of type `qa` settles `failed` with a reason that names each task.
_Avoid_: Pending QA, finished turn, background job

**First Handoff**:
The first Agent turn of a Task that returns to the Daemon ([ADR-0057](docs/adr/0057-daemon-exclusively-owns-implement-task-status.md)'s "Agent handoff"). Before it, the Fallback Chain's Lost Rollout and Provider System Error exception still applies; after it, no Fallback Selection takes over work a turn handed back.
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
The maximum number of Verification attempts for Tasks that one Implement Run may execute concurrently, configured by `verification.concurrency` and defaulting to `1`; an exclusive retry consumes the Run's entire capacity.
_Avoid_: Task Capacity, worktree concurrency, machine-wide test lock

**Task Worktree**:
The ephemeral git worktree one concurrently executing Task runs in — created from the Run Branch tip at Task start and integrated back onto the Run Branch at settlement; kept only when its Task fails.
_Avoid_: Run Worktree, sandbox, scratch dir

**Detached Run**:
A Run started with the detach flag: roundfix re-executes itself as a session leader independent of the caller, reports the Run's identifier, is followed by humans through Attach or by Supervisors through the Run Event Stream, and is ended through the Stop Command.
_Avoid_: Background job, nohup run, daemon

**Run Worktree**:
The isolated git worktree a Spec Run executes in — created on the Run Branch at Run start, recorded on the Run, removed after a Clean integrated spec outcome, and kept as the inspection and settle surface otherwise. Review Runs do not create Run Worktrees.
_Avoid_: Sandbox, scratch dir, user checkout

**Worktree Root**:
The directory two levels above a recorded worktree path, the `worktree.location` it was created under, within which symbolic links are followed; the repository and worktree directories below it must be real directories ([ADR-0268](docs/adr/0268-a-recorded-worktree-path-is-resolved-within-its-worktree-root.md)).
_Avoid_: Worktree directory, Roundfix Home

**Run Branch**:
The named branch (`roundfix/run-<id>`) that carries a Spec Run's commits inside its Run Worktree until integration moves them to the user's branch. Review Runs use the user's checkout branch directly; older pending Run Branch work is handled by Branch Integrity Preflight before review work starts.
_Avoid_: Temp branch, detached HEAD, feature branch

**Branch Disposition**:
A recorded terminal reason a Run Branch no longer needs integration and may be discarded. A superseded Branch Disposition means every branch commit is already reachable from the target, or a later integrated Run covered the same Tasks.
_Avoid_: Branch cleanup, forced deletion, reconcile result

**Run Worktree Reconciliation**:
The proof-based classification of a terminal Spec Run's retained Git surfaces: `safe` when the Run Branch is contained in its target or its content is represented at the merged head; `superseded` when a newer QA Report or the merged-head proof represents its Task or QA Report commits, archived Spec-directory work, or scoped uncommitted leftovers; `unintegrated` when resolved evidence does not represent the Run; `dirty` when tracked or untracked changes lack archived-Spec supersession proof; `unknown` when metadata or Git evidence cannot prove another state; and `released` only when both the Run Worktree and Run Branch are absent. A declared or recorded path covers the dirty paths beneath it, and a kept dirty Run whose Spec is archived at the proof head names the failing proof or first path outside that scope ([ADR-0289](docs/adr/0289-a-delivered-specs-dirty-run-is-superseded-under-its-declared-directories-and-a-kept-one-names-why.md)); only freshly revalidated `safe` or `superseded` work can be cleaned up.
_Avoid_: GC, force cleanup, manual branch deletion

**Integration Pending**:
A terminal Spec Run outcome where the Spec Run's work completed but its commits could not be fast-forwarded onto the user's branch (local changes overlap or the branch diverged); the commits stay on the Run Branch and the report names the integration command. Review Runs do not end Integration Pending.
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

**Failed**:
A terminal Run outcome where the Run itself broke — an adapter, Git, bootstrap or infrastructure error — before it could complete; distinct from an Unresolved Outcome, which finished with work remaining. It exits `1`.
_Avoid_: Unresolved, crash, error outcome

**Stopped**:
The terminal Run outcome a Stop Request or a Force Stop produces: the Run ended before another outcome, its local work preserved and any kept Run Worktree available to Task Carry-Forward. It exits `0`.
_Avoid_: Cancelled, aborted, killed

**Budget Exceeded**:
The terminal Run outcome (`BudgetExceeded`) of a Run whose Run Budget passed: a watch Run settles it from its deadline, and a Spec Run settles it with the configured maximum and the elapsed time, cancels the Agent Sessions it owns and keeps its Run Worktree and Run Branch for Task Carry-Forward ([ADR-0158](docs/adr/0158-a-budget-expired-run-settles-budget-exceeded-and-stays-recoverable.md)).
_Avoid_: Stopped, Timed Out, Max Rounds Reached

**Timed Out**:
The terminal Review Run outcome (`TimedOut`) of a watch Run whose wait for accepted Review Source Evidence reached its deadline before Clean, recorded with the evidence detail it last saw. It exits `1`.
_Avoid_: Budget Exceeded, Clean Unverified, timeout error

**Checkout Moved**:
The terminal Review Run outcome (`CheckoutMoved`) recorded when the user's checkout no longer matches the target recorded on the Run at Preflight Validation, so the Run ends instead of committing to the wrong branch.
_Avoid_: Branch switch, detached HEAD, Stopped

**Fetched**:
The terminal outcome of a Fetch Run: Review Issues and their markdown artifacts persisted, nothing resolved, nothing pushed. It exits `0`.
_Avoid_: Clean, downloaded, synced

**Clean**:
A terminal Run outcome where the cycle completed with nothing unresolved remaining: no Unresolved Review Issues for review work, or every Task completed — and a passing QA verdict when the graph authors a gate — for spec work.
_Avoid_: Success, done, green

**Implement Clean**:
A Spec Run's Clean: every Task completed and, when the graph authors a QA gate, a passing verdict. An archive, a QA Archive Override, a review verdict or a merge never establishes it.
_Avoid_: Archived, QA pass, delivered

**Clean Unverified**:
A terminal Review Run outcome where the cycle completed with nothing unresolved and the Final Push succeeded, but accepted Review Source Evidence for the pushed head never appeared within the grace period, so Merge-Ready was not confirmed. Distinct from Clean by outcome and by exit code `3`, which Review Skipped shares; an explicit skipped review ends Review Skipped instead.
_Avoid_: Clean, failure, timeout

**Review Skipped**:
A terminal Review Run outcome where the Review Source explicitly reports that it did not review the expected head. It is never Merge-Ready, carries the source-reported reason and next action, and exits `3` like Clean Unverified.
_Avoid_: Clean, Clean Unverified, zero Review Issues

**Run Budget**:
A safeguard that stops a Run before it can continue indefinitely and indirectly consume unbounded resources.
_Avoid_: Max rounds, review round limit

**Run Window**:
The durable, repository-scoped bound on when an Implement Run may be created. It governs a Run's start and never its finish: a Run created before the cutoff continues to its own terminal outcome. The implementation spellings `RunWindow` and `run_window` refer to this same term.
_Avoid_: Session, curfew, deadline

**Preflight Validation**:
The early checks Roundfix runs before starting a Run or work that would make the developer wait.
_Avoid_: Best-effort validation, late failure

**Branch Integrity Preflight**:
The deterministic Preflight Validation for Review Runs that blocks fetch, resolve, and watch while unintegrated Run Branch commits or another Run remain bound to the PR Head Branch, integrating fast-forwardable work automatically except QA-report-only branches and otherwise naming each pending worktree, run id, and recovery command. Skippable only through an explicit bypass that publishes an audit comment on the pull request.
_Avoid_: Advisory check, best-effort warning, soft gate

**Verification**:
The authoritative command or commands the Daemon runs verbatim in the repository root to decide whether Agent or Settle Command work can be settled and committed. A failure returns only its diagnostics to the Agent Session; for Tasks, a pass is required before status `completed`.
_Avoid_: CI, smoke test, best-effort check

**Repository Contract Test**:
A Go test under the `docscontract` or `repocontract` build tag that checks the repository itself — its documents, derived artifacts or test wiring — rather than one package's behavior. The Full Contract Run runs every one; `make verify-changed` runs the ones its Contract Relevance selects, and `make verify-docs` runs the `internal/docscontract` tests and the contracts its `repo-test` step lists ([ADR-0252](docs/adr/0252-the-selective-gate-runs-the-repository-contracts-a-change-makes-relevant.md), [ADR-0253](docs/adr/0253-every-repository-contract-test-runs-on-main-and-before-a-release.md)).
_Avoid_: docs test, integration test, CI-only test

**Sequential Test**:
A top-level test that cannot run in parallel because it changes process-wide state and carries a `// Sequential: <reason>` comment naming that state, with at least two words, in its doc comment or before its first statement ([ADR-0259](docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md)).
_Avoid_: Unmarked test, serial test, parallel exception

**Parallel Test Package**:
A package whose top-level tests each call `t.Parallel()` as their first statement or are Sequential Tests, with a repository rule enforcing both the declaration and a ceiling on Sequential Tests that may only fall; files under a `docscontract` or `repocontract` build constraint are excluded ([ADR-0259](docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md)).
_Avoid_: Parallel suite, concurrent package, fast package

**Single-Name Test Binary**:
An executable a test may start that keeps one directory entry while any process may start it; a test that needs it at another path copies it or links it symbolically, because macOS kills a new process whose executable path names an entry that has disappeared ([ADR-0324](docs/adr/0324-a-test-never-gives-a-running-executable-a-second-name.md)).
_Avoid_: hard-linked fixture, shared inode, binary alias

**Real-Time Budget**:
A wall-clock allowance below ten seconds that a test gives work it does not control: a child process, a goroutine or Git; an asserted allowance races work that blocks until cancellation or the test binary's exit ([ADR-0271](docs/adr/0271-a-test-bounds-by-real-time-only-what-it-asserts.md)).
_Avoid_: Run Budget, test deadline, arbitrary timeout

**Contract Relevance**:
The declaration in a Repository Contract Test's source that decides when `make verify-changed` runs it: `always` on every change; by default, when a file in its package directory changes; `relevant`, when its package directory or a declared path changes; or `boundary`, never in that gate. A malformed declaration fails the gate. A change to `go.mod`, `go.sum` or `Makefile`, and a change list that Git cannot produce, each select every contract that is not `boundary`. The selector's summary line names each `relevant` contract it leaves out and each `boundary` contract. The Full Contract Run ignores Contract Relevance ([ADR-0252](docs/adr/0252-the-selective-gate-runs-the-repository-contracts-a-change-makes-relevant.md), [ADR-0253](docs/adr/0253-every-repository-contract-test-runs-on-main-and-before-a-release.md)).
_Avoid_: test tags, impact map, test filter

**Full Contract Run**:
`make verify-contracts`: one run of every Repository Contract Test that `cmd/verify-select` discovers, whatever its Contract Relevance, with no hand-kept list. CI runs it on every push to main and before every release, so a contract the selective gate leaves out of a pull request still runs before the change ships ([ADR-0253](docs/adr/0253-every-repository-contract-test-runs-on-main-and-before-a-release.md)).
_Avoid_: full suite, nightly run, repo-test

**Incremental Verification**:
The selected fast local check recorded by the `verification.incremental` Baseline decision for validating the current change outside a Run while reusing safe local state. In a Daemon-assigned turn the Agent runs focused tests instead, because the Daemon runs the Task's Verification and the repository Verification at settlement ([ADR-0257](docs/adr/0257-inside-a-run-the-daemon-is-the-only-full-gate-and-the-baseline-states-each-rule-once.md)). It is distinct from the complete repository Verification selected by `verification.gate`.
_Avoid_: Complete Verification, CI gate, optional check

**Waiting for Verification**:
The observable per-Task phase after Agent work is implementation-ready and before the Task acquires Verification Capacity; it is distinct from an Agent that is still working and from a Verification command that has started.
_Avoid_: Queued Agent, pending Task, blocked Run

**Temporary Verification Failure**:
A project-authored Verification command exit with code `75`, eligible for one Daemon-controlled exclusive retry per Task; Roundfix never infers it from logs, timing, or framework-specific error text.
_Avoid_: Flaky test, generic non-zero exit, log-matched infrastructure error

**Verification Feedback**:
The failure diagnostics returned to an Agent Session after the Daemon runs Verification. A passing Verification produces no Agent feedback.
_Avoid_: Full verification output, test log, progress stream

**Settlement Check**:
A gate fact the Daemon checks for a non-QA Task before it settles in a Task Graph with an authored QA gate. The checks cover the Spec Consistency Check, the repository Verification at settlement, and authorization of the prospective Task commit.
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

**Artifact Root**:
The directory under which a Run writes `runs/<run-id>`. The recorded Artifact
Directory resolves to this root the way the Run resolved it: a blank value uses
the repository's default root, or, for an older Run, the checkout's default
when only it holds the Run's directory; a relative value resolves under the
Run's Git root; and an absolute value is used as recorded. One resolver serves
the storage report, GC sanitation, Run Retention and the run cause report
([ADR-0301](docs/adr/0301-gc-sanitation-classifies-the-artifact-root-each-run-wrote-to.md),
[ADR-0255](docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md),
[ADR-0214](docs/adr/0214-why-a-verification-failed-is-classified-by-deterministic-signatures.md)).
_Avoid_: Artifact Directory, artifact folder, storage root

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

**Storage Snapshot**:
A consistent read of the Run Database that measures every page group, including committed write-ahead log frames, so grouped bytes equal snapshot bytes exactly ([ADR-0296](docs/adr/0296-the-storage-report-reconciles-every-page-of-one-storage-snapshot.md)).
_Avoid_: Main file size, immutable storage-report measurements

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
The support command that inspects the Run Worktrees and Run Branches of terminal Spec Runs and reports their Run Worktree Reconciliation state. It is read-only by default and carries three named mutation switches, none of which has a force bypass: `--apply` is the only one that releases retained Git surfaces, removing only freshly revalidated `safe` or `superseded` work; `--discard-superseded` records a Branch Disposition before removing a branch it can prove superseded; and `--carry-forward` performs Task Carry-Forward. `--apply` also releases each Stale Run Owner Record without sending a signal ([ADR-0313](docs/adr/0313-a-proven-owner-identity-mismatch-releases-a-terminal-runs-owner-record-without-a-signal.md)).
_Avoid_: GC Command, Settle Command, automatic integration

**Binding Text**:
The planning text a Task binds to: its own Task Graph node and projection row, and the PRD and TechSpec sections its References name. A missing node, row, References bullet or named section binds the whole file; a planning amendment preserves a settled Task only when its Binding Text is unchanged ([ADR-0284](docs/adr/0284-carry-forward-binds-a-task-to-the-planning-text-it-references.md)).
_Avoid_: item-level binding, whole-Spec binding

**Task Carry-Forward**:
The explicit act that hands a settled Task from a terminal Spec Run's Run Branch back to the user's checkout, so work that already ran and passed its Verification is never executed again to reach the same result. It accepts a Run whose outcome is Stopped, Unresolved or Budget Exceeded ([ADR-0158](docs/adr/0158-a-budget-expired-run-settles-budget-exceeded-and-stays-recoverable.md)) and refuses every other terminal outcome. It carries a Task only on proof — a passing Verification verdict, exactly one settlement commit, declared inputs unmoved since settlement except planning amendments that preserve its Binding Text ([ADR-0284](docs/adr/0284-carry-forward-binds-a-task-to-the-planning-text-it-references.md)), a clean checkout, and a repository-local Spec Root — refuses the whole set rather than carrying part of it, and stamps each carried Task with the Run and commit that established it. Reached through the Reconcile Command's `--carry-forward` switch or through a Delivery Retry; it is never automatic. A Task already completed on the target is nothing to carry, so the act does not repeat itself. The implementation spelling `CarryForward` refers to this same term.
_Avoid_: replay, resume, automatic carry, QA row carry-forward

**Delivery Queue**:
The durable, ordered set of Specs that Roundfix advances from an Implement Run through merge with a fixed optional deadline and per-item retry limit. A blocker parks only its item while later items can continue.
_Avoid_: Release queue, Batch, Run list

**Base Test Tools**:
The tools the repository's tests require: `go`, `gofmt`, `git`, `make`, `sh`, `cat`, `kill`, `ps`, `id`, `mkfifo` and `cmd` ([ADR-0294](docs/adr/0294-the-linux-lane-runs-like-ci.md)).
_Avoid_: container tools, optional tools

**External Test Tool**:
A tool a test runs outside the Base Test Tools. A test that needs one skips naming it when it is absent from PATH ([ADR-0294](docs/adr/0294-the-linux-lane-runs-like-ci.md)).
_Avoid_: optional dependency, container tool

**Linux Lane**:
The optional Linux test step of a Delivery Queue candidate's repository gate, covering its changed Go packages before publication. It also runs on demand with `roundfix deliver lane`, independently of the queue opt-in, against a checkout revision's changed Go packages or every package with `--all`. It runs under an init on the image `delivery.linux_lane_image` names, or the default image when empty ([ADR-0294](docs/adr/0294-the-linux-lane-runs-like-ci.md), amending [ADR-0288](docs/adr/0288-the-delivery-queue-gate-runs-the-changed-go-packages-in-a-linux-lane.md)). It reads the repository history through the object store mounted read-only, including the common directory's store for a linked worktree ([ADR-0297](docs/adr/0297-the-linux-lane-reads-the-repository-history-read-only.md)).
_Avoid_: Linux gate, Docker gate, cross-platform Verification

**Lane Environment Skip**:
A skip of an enabled Linux Lane for any reason other than having nothing to test, or of a lane whose fetched configuration cannot be read. It passes the gate and is recorded on the item, shown by `roundfix deliver status`, the queue's pull request body and `roundfix doctor` ([ADR-0303](docs/adr/0303-a-linux-lane-skipped-for-its-environment-is-shown-where-the-delivery-is-read.md)).
_Avoid_: expected skip, lane failure

**Lane Container Failure**:
A Linux Lane failure caused by its container setup, command invocation, a signal or the run cap. It parks the candidate in the environment Park Class ([ADR-0294](docs/adr/0294-the-linux-lane-runs-like-ci.md)).
_Avoid_: change failure, generic gate failure

**Lane Change Failure**:
A Linux Lane failure attributed to the candidate: its Go test run, or the `go list` run that selects its packages exiting below 128. It parks the candidate in the finding Park Class ([ADR-0294](docs/adr/0294-the-linux-lane-runs-like-ci.md), [ADR-0316](docs/adr/0316-the-linux-lane-names-a-failed-go-list-apart-from-go-test.md)).
_Avoid_: container failure, generic gate failure

**Queue Owner**:
The `roundfix deliver` process that holds one Delivery Queue and advances its items in order, waiting on each Prerequisite Spec; a Delivery Retry hands a parked item to the live Queue Owner or starts a new one. When an item can still advance, `roundfix deliver status` reports a recorded Queue Owner proven not running by its PID and start-time identity; status never reclaims it, so the item can still advance after `roundfix deliver resume` ([ADR-0279](docs/adr/0279-status-reports-a-queue-owner-proven-gone-and-a-dismissal-survives-a-fix-that-leaves-its-anchor.md)).
_Avoid_: Scheduler, daemon, deliver loop

**Delivery Plan**:
The read-only account of which Specs have delivery authority, what blocks the others and which production premises ordered Specs share. It reports a prepared queue without recording or running one.
_Avoid_: Delivery Queue, execution plan, approval grant

**Delivery Revalidation**:
The check of a queued Spec against the refreshed default branch immediately before its first Run. A strict finding parks it as `revalidation-failed`; a production-premise overlap records a `premise-changed` warning without parking it. It does not run for an Archived Item Resume.
_Avoid_: Initial validation, retry, approval check

**Archived Item Resume**:
The start of a new Delivery Queue item whose continued Item Branch already archives its Spec in a commit proven to be an exact archive of it: the item enters review with its head as the candidate, without Delivery Revalidation, an Implement Run or a new archive. An item branch whose archive cannot be proven parks as `archived-item-unproven` ([ADR-0312](docs/adr/0312-a-new-delivery-queue-resumes-an-archived-item-branch-at-review.md)).
_Avoid_: Archived retry, requeue

**Delivery Queue Limit**:
A queue's fixed optional deadline, per-item retry allowance or Queue Token Ceiling. It bounds new item starts or explicit retries without changing a Run Budget or an item already past `queued`.
_Avoid_: Run Budget, concurrency

**Pending Question**:
The one operator decision presented for the lowest-position parked Delivery Queue item, with its blocker, answer and number of parked items waiting behind it. Only an explicit Delivery Retry or a new queue answers it.
_Avoid_: Warning, automatic recommendation, blocker list

**Delivery Retry**:
The explicit act that returns one parked Delivery Queue item to the stage supported by its recorded evidence and hands it to a live or newly started Queue Owner. It may perform Task Carry-Forward from every Run of the item's Spec, newest first, and is never automatic. It refuses a Spec whose Archive Record is `cancelled` ([ADR-0262](docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md)).
_Avoid_: Resume Command, automatic retry, replay

**Judge Log**:
The per-month JSONL record of every advisory judgment call: the question, the transport that carried it (OpenRouter or TypeSafe), the answer, probabilities, confidence, latency, input tokens, the model the service reported and the call's cost. Its monthly sum is what the spend ceiling reads ([ADR-0201](docs/adr/0201-the-judge-sends-only-spec-artifacts-and-spends-under-a-monthly-ceiling.md)).
_Avoid_: Run Event Journal, audit log

**Advisory Judgment**:
A typed semantic judgment about a Spec artifact, such as whether a cited ADR supports the attributed claim, or about a source relationship, such as whether a Finding or Backlog Entry shares a Spec's context, that the authoring model must answer but that never gates, never changes another command's exit code, and fails open ([ADR-0200](docs/adr/0200-spec-authoring-gets-an-advisory-judge-that-never-gates.md), [ADR-0209](docs/adr/0209-the-judge-suggests-which-sources-share-a-spec.md)).
_Avoid_: Gate, check, finding

**Grouping Suggestion**:
An Advisory Judgment that an open Finding or Backlog Entry shares the context of a source the Spec already adopts, so the author adopts it into the same Spec or states why it stays apart; it never moves a source by itself ([ADR-0209](docs/adr/0209-the-judge-suggests-which-sources-share-a-spec.md)).
_Avoid_: Merge rule, automatic grouping

**Grouping Bound**:
The size limit on grouping sources into one Spec: four implementation Tasks plus its QA gate. A grouped scope that needs more is split into Specs that each fit, and each source keeps exactly one owning Spec ([ADR-0208](docs/adr/0208-work-items-that-share-a-context-share-a-spec.md)).
_Avoid_: Spec size rule, Task cap

**Setup Snapshot**:
The Baseline's embedded mirror of one upstream skills setup, written only by the asset sync, that lists by name the skills a Baseline Profile takes, with the tree digest of each ([ADR-0191](docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md), [ADR-0241](docs/adr/0241-a-repository-profile-takes-its-skill-contracts-from-the-embedded-setup-snapshots.md)).
_Avoid_: Skill lock, catalog, upstream checkout

**Composed Setup Snapshot**:
A setup snapshot built by name from upstream setups in a fixed order, with duplicate skills dropped, so a built-in profile can span several toolchains ([ADR-0204](docs/adr/0204-a-composed-profile-takes-a-setup-composed-from-upstream-setups-by-name.md)).
_Avoid_: Merged setup, custom setup

**Retired Skill**:
A skill the Baseline no longer requires although the upstream catalog may still list it. No module requires or dispatches it, the asset sync drops it from every Setup Snapshot, and `roundfix baseline update` lists an installed copy for the adopter to delete without deleting it ([ADR-0246](docs/adr/0246-a-retired-skill-leaves-the-baseline-and-the-adopter-removes-its-copy.md)).
_Avoid_: Deprecated skill, obsolete skill

**Missing Required Skill**:
An external skill the adopted modules require and dispatch that is absent from `.agents/skills` or `skills-lock.json`, or no longer matches its lock entry. The `roundfix baseline update` preview lists it under `skills.missing` with its restore command and reports `plan_ready` instead of `current`; `roundfix doctor` and `roundfix skills check` fail on it.
_Avoid_: Uninstalled skill, absent skill, dispatched skill gap

**Park Class**:
The named reason a Delivery Queue item is parked — for example a prerequisite not yet merged, a pull request conflict, an environment-only QA partial or a check that failed twice outside the item's packages — shown with the next command that answers it ([ADR-0192](docs/adr/0192-a-conflict-confined-to-declared-derived-paths-is-resolved-by-regeneration.md), [ADR-0193](docs/adr/0193-a-spec-names-its-prerequisites-and-the-owner-waits-for-them.md)).
_Avoid_: Blocker text, error code

**Prerequisite Spec**:
A Spec another Spec names in its Task Graph as `requires`; the Queue Owner does not start the dependent item until each prerequisite is merged into the default branch ([ADR-0193](docs/adr/0193-a-spec-names-its-prerequisites-and-the-owner-waits-for-them.md)). A cancelled Spec's record never meets a prerequisite ([ADR-0262](docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md)).
_Avoid_: Queue order, dependency Task

**Derived Path Declaration**:
A Project Config entry naming a path pattern and the sanctioned command that regenerates it; a pull request conflict confined to declared derived paths is resolved by regeneration instead of parking ([ADR-0192](docs/adr/0192-a-conflict-confined-to-declared-derived-paths-is-resolved-by-regeneration.md)).
_Avoid_: Generated file list, merge strategy

**Module Version Record**:
A generated record of each Baseline module's versions and content digests, so a module version names one recorded content and the record step can choose the next free version ([ADR-0250](docs/adr/0250-a-module-version-is-chosen-when-recorded-and-the-coverage-record-lists-every-platform.md)).
_Avoid_: Module changelog, manual version note

**Coverage Record**:
A generated record of the tests in the repository's packages, including the release platforms on which platform-limited tests are built, so coverage comparisons are stable on every host ([ADR-0250](docs/adr/0250-a-module-version-is-chosen-when-recorded-and-the-coverage-record-lists-every-platform.md)).
_Avoid_: Host test list, test result report

**Evidence Snapshot**:
The Daemon's record, written when a QA pass closes, of the inputs each carriable passing row observed at the audited head ([ADR-0194](docs/adr/0194-the-daemon-records-what-a-qa-row-observed-and-hands-a-failed-pass-to-the-next.md)).
_Avoid_: QA cache, Agent-written evidence

**Carried Row**:
A QA matrix row a later pass keeps from the prior report because its Evidence Snapshot proves its inputs unmoved; it keeps its original provenance. Rows that read the repository Verification, the pull request or the commit range are never carried ([ADR-0194](docs/adr/0194-the-daemon-records-what-a-qa-row-observed-and-hands-a-failed-pass-to-the-next.md), [ADR-0195](docs/adr/0195-rows-that-read-the-gate-or-the-commits-are-observed-on-every-pass.md)).
_Avoid_: Skipped row, Task Carry-Forward

**Carry Disposition**:
The fixed-list outcome a QA pass records for every row of the prior report: carried, re-run as stale, re-run as always observed, or re-run because it did not pass ([ADR-0195](docs/adr/0195-rows-that-read-the-gate-or-the-commits-are-observed-on-every-pass.md)).
_Avoid_: Row status, verdict

**Delivery Convention**:
A commit shape the delivery itself produces — the QA Report commit, the Daemon's status writes, the archive commit, an all-pending planning candidate — that the pre-PR reviewer is told about and never raises as a defect ([ADR-0196](docs/adr/0196-a-pre-pr-review-finding-parks-only-after-validation.md)). A finding that names a Result Status Note against a `completed` Task is dismissed by validation without the sealed validator ([ADR-0317](docs/adr/0317-a-result-status-note-of-a-completed-task-is-dismissed-without-the-validator.md)).
_Avoid_: Review exemption, ignore list

**Result Status Note**:
A line of a non-QA Task's `## Result` that reports the Task's status as `in_progress` or Daemon-owned, which the Agent writes while the Daemon holds the Task `in_progress`; it does not contradict the `completed` status the Daemon writes after the Task's Verification passes, and the Pre-PR Review Command dismisses a finding that names that pair without the validator ([ADR-0317](docs/adr/0317-a-result-status-note-of-a-completed-task-is-dismissed-without-the-validator.md)).
_Avoid_: Status contradiction, stale status, Result status claim

**Review Segment**:
One of at most four consecutive first-parent commit ranges of a pre-PR review candidate whose diff exceeds the review bound, each within the bound and reviewed in order in one review session, whose findings form one Pre-PR Review Record ([ADR-0282](docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md)).
_Avoid_: Review round, independent review

**Correction Review**:
The one review of the fix after round 2 when every standing round-2 finding has exactly one valid fixed disposition, covering only the delta from the round-2 head. It belongs to round 2 and cannot open another Correction Review ([ADR-0286](docs/adr/0286-a-fixed-round-2-finding-gets-one-correction-review-of-its-fix.md)).
_Avoid_: Third round, new review lineage

**Reviewer Lineage**:
The at-most-two review rounds of one delivery candidate, where round 2 reviews only the delta since round 1 with round 1's findings and dispositions, continuing the same Agent Session when the runtime can resume it. It may include one Correction Review when every standing round-2 finding was fixed ([ADR-0197](docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md), [ADR-0286](docs/adr/0286-a-fixed-round-2-finding-gets-one-correction-review-of-its-fix.md)).
_Avoid_: Review cache, third round

**Turn Usage**:
The token count an ACP adapter reports for one prompt of a Run, recorded per Task and summed per Run; a prompt without a report is recorded as unreported, never as zero ([ADR-0198](docs/adr/0198-a-run-counts-the-tokens-its-adapters-report-by-each-reports-own-scope.md)).
_Avoid_: Cost estimate, spend

**Usage Basis**:
How a Turn Usage was counted: `turn` when the adapter's end-of-prompt count covers the whole turn, `request-sum` when Roundfix sums the adapter's per-request context readings ([ADR-0198](docs/adr/0198-a-run-counts-the-tokens-its-adapters-report-by-each-reports-own-scope.md)).
_Avoid_: Pricing model

**Reported Cost**:
The spend an ACP adapter itself reports for a turn; Roundfix prints only reported cost and never prices tokens from a table ([ADR-0198](docs/adr/0198-a-run-counts-the-tokens-its-adapters-report-by-each-reports-own-scope.md)).
_Avoid_: Estimated cost, price table

**Run Slot Wait**:
The bounded wait of a running Delivery Queue item whose Implement step was refused only by the ceiling on active Implement Runs, shown as its blocker text and polled against the Run Database until a slot frees, `runs.slot_wait` expires or the queue deadline arrives ([ADR-0304](docs/adr/0304-a-queue-item-waits-for-a-run-slot-up-to-runs-slot-wait.md)).
_Avoid_: Waiting stage, admission reservation

**Check Wait**:
The bounded wait of a Delivery Queue item in `checking` for its current-head checks and a mergeable state, and of its publication for an existing pull request to report the candidate head. Each wait lasts `delivery.check_timeout` (User Config, 60 minutes by default) from its own start, restarts once after a check re-run, and ends in a `checks-timeout` park that names the key and the duration ([ADR-0311](docs/adr/0311-the-delivery-queues-check-wait-is-a-user-config-duration-sixty-minutes-by-default.md)).
_Avoid_: CI timeout, check deadline

**Queue Token Ceiling**:
An optional Delivery Queue Limit: once the queue's recorded Turn Usage reaches it, the next item parks and retries are refused, while a running Task is never stopped ([ADR-0199](docs/adr/0199-a-queue-token-ceiling-stops-new-starts-and-never-a-running-task.md)).
_Avoid_: Run Budget, spending limit

**Init Command**:
The support command that creates User Config or Project Config before operational Runs.
_Avoid_: Bootstrap run, setup run

**Setup Command**:
The support command that verifies and prepares a machine for Runs of Roundfix: Node, minimum-supported acpx, Adapter Readiness, generated Agent Selection Profile Readiness, acpx local adapter overrides, User Config, and Project Config. It accepts the minimum tested acpx version and every newer compatible version, never downgrades a newer installation, proves proposed state before requesting authorization, and writes only authorized changes.
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
The support command that produces a Release Plan without editing release files, creating or pushing tags, publishing packages, or creating a GitHub release.
_Avoid_: Release Command, publish command, cut-release command

**Doctor Command**:
The support command that diagnoses a repository and machine's readiness for Runs of Roundfix — minimum-supported acpx, Adapter Readiness, Agent Selection Profile Readiness, Repository Skill Set, codex runtime hygiene, and whether Run storage is reclaimable. It reports the detected acpx version against the minimum, gives each check a next action, and mutates nothing. It evaluates Repository Skill Set readiness after, and independently from, Agent Selection Profile Readiness; unlike the Doctor Command, the Setup Command prepares the machine.
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

**Live Document**:
A tracked Markdown file inside an active Spec folder, an open Backlog Entry, an unresolved Finding, or under `docs/references/`; a citation is a relative link destination or repository-relative path naming a file removed by retirement ([ADR-0260](docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md)).
_Avoid_: Citing file, live source, current document

**Archive Relink**:
The Archive Command's pre-write operation that promotes a cited non-core file to `docs/references/<basename>`, rewrites Live Document citations to the promoted copy or Archive Record, and records rewritten paths under the optional `relinked` list ([ADR-0260](docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md)).
_Avoid_: Archive rewrite, source duplication, link cleanup

**Legacy Link Rebase**:
The rewrite of relative destinations in an Archive Record outcome, its unproven entries and each promoted copy: a folder destination points at its promoted copy, or becomes a Git Locator when no copy survives, never the record itself; every other destination keeps the same repository path. An unproven entry is read as text of the folder's `_prd.md` ([ADR-0278](docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md), [ADR-0308](docs/adr/0308-an-archive-rebases-the-unproven-links-it-records-and-the-record-link-repair-points-a-promoted-target-at-its-copy.md)).
_Avoid_: Archive Relink, broken-link report, arbitrary link rewrite

**Default Branch Head**:
The fetched `origin/<default>` commit, or the local default branch commit only when no remote-tracking branch exists. When origin/HEAD is unset and the current branch does not name the default, exactly one of `origin/main` and `origin/master` names it ([ADR-0295](docs/adr/0295-reachability-reads-the-fetched-default-branch-and-origin-names-it-when-origin-head-is-unset.md)).
_Avoid_: Local-first default, work branch head

**Git Locator**:
The inline code `<path>@<commit>` naming a removed file at a commit that holds it and that the Default Branch Head or a tag reaches, written in place of a link to that file. When no such commit exists it names the path and the record's pull request; outside Markdown it is written as plain text without backticks ([ADR-0299](docs/adr/0299-the-history-conversion-rewrites-every-citation-it-lists-and-the-check-reports-a-removed-history-file-cited-without-a-git-locator.md)).
_Avoid_: Broken link, branch-only revision, Archive Record

**Citation Source**:
A tracked regular text file outside the skill trees and testdata directories, holding no NUL byte in its first 8,000 bytes, in which History Citations are found ([ADR-0299](docs/adr/0299-the-history-conversion-rewrites-every-citation-it-lists-and-the-check-reports-a-removed-history-file-cited-without-a-git-locator.md)).
_Avoid_: Markdown-only source, skill fixture

**Removed History Path**:
A path under the History Root's Findings, Backlog or ADR family that HEAD or its history holds and the working tree lacks ([ADR-0299](docs/adr/0299-the-history-conversion-rewrites-every-citation-it-lists-and-the-check-reports-a-removed-history-file-cited-without-a-git-locator.md)).
_Avoid_: Missing file, retained history file

**History Citation**:
A citation to a file or folder a selected history unit removes, listed in text and JSON and bound into the Plan Digest. As a History Citation, it is found and rewritten in every Citation Source to a Git Locator for removed text, or to the surviving record or promoted copy when one exists ([ADR-0299](docs/adr/0299-the-history-conversion-rewrites-every-citation-it-lists-and-the-check-reports-a-removed-history-file-cited-without-a-git-locator.md)).
_Avoid_: Relocation Citation, broken-link scan, stale citation

**Stale History Citation**:
A citation already unresolved in the working tree whose target remains in the History Full Tag, which `roundfix baseline update` reports with a `git show history-full:<path>` restore command outside the Plan Digest and update state ([ADR-0260](docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md)).
_Avoid_: Broken link, history warning, missing source

**Archive Record**:
The small `<slug>.md` record an archive leaves under the resolved archive root: a file named `NNNN-<slug>.md` whose front matter declares `roundfix/archive-record/v1`; any other Markdown file in that root is not a record. It carries the disposition, QA Report and verdict, promoted references, `source_revision` and, once the delivery pull request is open, `pull_request`; the removed Spec folder stays recoverable in Git at that revision, after a squash merge through Source Revision Recovery. The History Sanitize Command also writes one for a Legacy Archive Folder; a folder without a QA Report receives the disposition `no-qa`, while a folder whose newest QA Report failed without an override receives `failed-qa` and keeps that verdict and report name. `failed-qa` is never a pass or an override ([ADR-0248](docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md), [ADR-0251](docs/adr/0251-a-legacy-archive-folder-is-read-leniently-and-a-failed-qa-keeps-its-verdict.md), [ADR-0260](docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md), [ADR-0261](docs/adr/0261-an-archive-record-names-its-pull-request-and-readers-fetch-its-head.md)). A cancelled Spec's record has the disposition `cancelled` and carries `cancel_reason`, `completed_tasks`, `pending_tasks` and `returned_sources`; only its listed completed Tasks count as completed ([ADR-0262](docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md)). Its summary is the QA Report's Outcome paragraph or the PRD's title paragraph. A summary paragraph that ends with a colon carries the list or paragraph that follows it, on the same line.
_Avoid_: Archive folder, archive log, copied Spec

**Folder Removal Recovery**:
How the Archive Command and the Cancel Command end when the final removal of the Spec folder fails: every write stays in place, and the command prints the command that finishes the removal ([ADR-0269](docs/adr/0269-archive-and-cancel-refuse-ignored-spec-files-and-finish-a-failed-removal-forward.md)).
_Avoid_: Rollback, retry archive, partial cleanup

**Archive Advice**:
Jev's advisory classification of candidate files before the archive cut. It labels reusable knowledge, repository records or transient evidence and never gates the archive.
_Avoid_: Archive verdict, promotion decision, required evidence

**Secret Pattern**:
One reviewed credential shape with a kind name, kept in one table ([ADR-0273](docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md)).
_Avoid_: Secret Match, credential value, ad hoc pattern

**Secret Match**:
One line that matches a Secret Pattern, reported by path, line and kind and never by value ([ADR-0273](docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md)).
_Avoid_: Secret Pattern, secret value, evidence sample

**Secret Rewrite**:
A rewritten head of an archived item parked `secret-found`, bounded by its Run start, whose tree differs from the parked candidate's only at Secret Match lines in paths the blocker names, with the same modes and line counts and no Secret Match left in those files ([ADR-0292](docs/adr/0292-a-task-commit-is-scanned-for-secrets-and-a-retry-accepts-a-secret-rewrite.md)). A Delivery Retry records it as the newest candidate and returns the item to `reviewing`.
_Avoid_: Post-archive correction, reset, arbitrary history rewrite

**Secret Scan**:
A check of the lines a Task commit would add, the files a QA Report commit would stage and the commits a Roundfix push would send; it refuses on a Secret Match ([ADR-0273](docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md), [ADR-0292](docs/adr/0292-a-task-commit-is-scanned-for-secrets-and-a-retry-accepts-a-secret-rewrite.md)).
_Avoid_: Redaction, post-push scan, advisory check

**Archive Command**:
The support command that performs an Archive Relink, then retires a Spec under the archive eligibility contract, writes its Archive Record, and removes the Spec folder while keeping it in Git at `source_revision`. It warns when `source_revision` is not an ancestor of the default branch, and `--pull-request <n>` stores the delivery pull request in the record, also in a record the current branch already added ([ADR-0261](docs/adr/0261-an-archive-record-names-its-pull-request-and-readers-fetch-its-head.md)). The normal QA eligibility and a user-authorized QA Archive Override are distinct dispositions. Neither retirement nor an override fabricates Task completion or a passing QA verdict ([ADR-0260](docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md)).
_Avoid_: Move command, retire run, cleanup command

**Source Revision Recovery**:
A reader's fetch of `refs/pull/<N>/head` when an Archive Record's `source_revision` is missing locally, with the number taken from the record's `pull_request` or from the trailing `(#N)` of the commit that added the record. It writes no ref, never prompts, and leaves the reader's unavailable outcome when it fails ([ADR-0261](docs/adr/0261-an-archive-record-names-its-pull-request-and-readers-fetch-its-head.md)).
_Avoid_: Backfill, archive ref, history restore

**Pull Request Record Commit**:
The Delivery Queue's commit, made after it opens a Spec's pull request, that stores the number in the Archive Record its archive commit added and changes nothing else; the queue accepts it by proof instead of a second repository gate ([ADR-0261](docs/adr/0261-an-archive-record-names-its-pull-request-and-readers-fetch-its-head.md)).
_Avoid_: Archive commit, record backfill, amend

**Spec Number**:
The one to ten ASCII digits that start a Spec slug when a `-` follows them, such as `0088`; a slug without them has no Spec Number. The Delivery Queue names the Spec by it in the Pull Request Record Commit subject, and in the archive commit subject and the pull request title when the slug form would pass 100 characters ([ADR-0302](docs/adr/0302-a-queue-commit-names-the-spec-by-number-and-a-hook-refusal-parks-with-the-hand-commit.md), [ADR-0309](docs/adr/0309-a-queue-archive-subject-and-pull-request-title-keep-the-slug-while-it-fits-100-characters.md)).
_Avoid_: Spec ID, spec prefix

**History Sanitize Command**:
The `roundfix history sanitize` command, a dry run unless `--apply --batch <n>`, that converts Legacy Archive Folders to Archive Records, removes retired families and rewrites their citations to Git Locators after the History Full Tag. `roundfix baseline update` runs the same planning for all eligible units at once ([ADR-0278](docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md)).
_Avoid_: Archive Command, unplanned deletion, cleanup command

**Record Link Repair**: The `record-links` unit of the History Sanitize Command and of the history section of `roundfix baseline update`, which rewrites each broken relative link and self-link of a surviving Archive Record: a target the record promoted points at its `docs/references/` copy, and any other becomes a Git Locator tried against the record's directory, then its `source` folder ([ADR-0308](docs/adr/0308-an-archive-rebases-the-unproven-links-it-records-and-the-record-link-repair-points-a-promoted-target-at-its-copy.md)).
_Avoid_: Live-links repair, record rewrite, Archive Relink

**Pending History**:
The units of a repository's existing history that the History Sanitize Command would plan, which the Managed Refresh includes in its plan and the upgrade notice names. `roundfix baseline update` plans these units with the Baseline Plan ([ADR-0254](docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md)).
_Avoid_: Active history, current history, unplanned archive

**Legacy Archive Folder**:
A Spec folder the Archive Command left under the archive root before archives wrote Archive Records ([ADR-0248](docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md)).
_Avoid_: Archive Record, active Spec, history root

**Refused Unit**:
A pending History Sanitize Command unit that cannot be converted. The plan and the batch name it with its reason on one line, leave it untouched, and do not count it toward `--batch <n>`. The update lists one without blocking the Baseline Plan or changing its state; a unit with uncommitted changes is refused ([ADR-0251](docs/adr/0251-a-legacy-archive-folder-is-read-leniently-and-a-failed-qa-keeps-its-verdict.md), [ADR-0254](docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md)).
_Avoid_: Converted unit, skipped silently, failed batch

**Lenient Legacy Reading**:
The reading of a Legacy Archive Folder that tolerates and names projection rows outside the graph and retired Task types, accepts a legacy list of maps `unproven`, and reads the summary, `legacy: not recorded` for missing legacy override approval or reason, adopted sources and creation-date fallbacks. Legacy Link Rebase is part of this reading for the outcome, the unproven entries and promoted copies. It is never applied to an active Spec ([ADR-0251](docs/adr/0251-a-legacy-archive-folder-is-read-leniently-and-a-failed-qa-keeps-its-verdict.md), [ADR-0254](docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md), [ADR-0277](docs/adr/0277-a-legacy-archive-folder-conversion-keeps-the-summary-links-and-dates-its-folder-held.md), [ADR-0308](docs/adr/0308-an-archive-rebases-the-unproven-links-it-records-and-the-record-link-repair-points-a-promoted-target-at-its-copy.md)).
_Avoid_: Relaxed active reading, manifest rewrite, schema bypass

**Sanitize Batch**:
The next units one `--apply` can convert, delivered as one pull request with the repository gates green and revertible; Refused Units do not count toward the batch. The update converts every convertible unit as one batch in its own change ([ADR-0248](docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md), [ADR-0251](docs/adr/0251-a-legacy-archive-folder-is-read-leniently-and-a-failed-qa-keeps-its-verdict.md), [ADR-0254](docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md)).
_Avoid_: Archive batch, migration wave, unreviewed cleanup

**Reduced History Entry**:
A legacy shape of a retired Finding or Backlog Entry containing its front matter, title, first paragraph and the revision holding its full text. The History Sanitize Command removes it; no retirement writes one ([ADR-0278](docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md)).
_Avoid_: Summary, Archive Record, truncated history

**History Full Tag**:
The annotated `history-full` tag on the last commit before the first batch, which keeps every removed byte reachable. The update creates it at `HEAD` when absent, never moves it, and the operator pushes it with `git push origin history-full` ([ADR-0248](docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md), [ADR-0254](docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md)).
_Avoid_: Lightweight tag, batch tag, archive stamp

**History Root**:
The single repository directory holding only Archive Records; Git keeps the full text of removed material. It is resolved in one place so every consumer uses the same root ([ADR-0278](docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md)).
_Avoid_: Archive folder, `_archived`, attic, trash

**History Relocation**:
One file carried in the Baseline Plan as an ordered ledger entry of source, destination and content identity, never into the Findings, Backlog, ADR or Review families of the History Root. An occupied destination refuses that relocation by name while its siblings still move ([ADR-0278](docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md)).
_Avoid_: Archive move, file migration, relocation payload

**Relocation Citation**:
A citation in a tracked file that resolves before a Baseline Plan's History Relocations and would not resolve after them. It is reported as a warning the Plan Digest binds; planning never rewrites it, and apply writes the same files.
_Avoid_: Broken link, dangling reference, link check

**Worktree Build Cache**:
A repository-declared directory of rebuildable build output retained in a Run Worktree, Task Worktree or Delivery Queue item worktree. It becomes reclaimable when its Run is terminal or its item is merged; active work keeps it ([ADR-0266](docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md)).
_Avoid_: Run artifacts, shared build cache

**Shared Build Cache**:
A build cache several worktrees share outside any one of them, declared in
`worktree.shared_caches`; the GC Command trims it above its declared size by
removing only cache entries no build used for the declared idle time, oldest
first ([ADR-0327](docs/adr/0327-gc-trims-a-declared-shared-build-cache.md)).
_Avoid_: global cache, Go cache cleanup, cache GC

**Item Binary**:
A copy of the Roundfix binary built by a Delivery Queue item for its own steps, retained while that item exists and is not merged. Its copies become reclaimable at merged release and through the GC Command when the item is merged or no longer queued ([ADR-0290](docs/adr/0290-merged-items-give-back-their-item-binaries-and-run-worktrees-leave-only-by-proof.md)).
_Avoid_: owner binary, shared binary

**GC Command**:
The support command that reclaims Run storage through Journal Retention, Run Retention, orphaned run artifacts and declared Worktree Build Caches of ended worktrees, plus the Item Binary sweep for merged or unqueued items ([ADR-0266](docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md), [ADR-0290](docs/adr/0290-merged-items-give-back-their-item-binaries-and-run-worktrees-leave-only-by-proof.md)). It reports only what it reclaimed, keeps Active Runs and active-run locks, and removes a Run row only through Run Retention ([ADR-0255](docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md)).
_Avoid_: Clean command, vacuum, purge

**Journal Retention**:
The configured age window after which a terminal Run's Run Event Journal and artifact directory become eligible for pruning. Active Runs are never eligible; a retention of zero keeps everything ([ADR-0033](docs/adr/0033-the-run-event-journal-is-pruned-by-retention.md)). It is the inner window: Run Retention removes the whole Run past the outer window, so a Journal Retention longer than the Run Retention has no effect past it ([ADR-0255](docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md)).
_Avoid_: Log rotation, TTL, expiry

**Run Retention**:
The User Config window, `store.run_retention_days` of 7, 15 or 30 days and 30 by default, after which a terminal Run leaves the Run Database whole: its row, its Run Event Journal, its Agent Selection records, its token usage and its artifact directory. A Run that is not terminal, that a Delivery Queue references, whose Run Worktree still exists, or whose artifact directory sits under an unproven root is kept ([ADR-0255](docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md)).
_Avoid_: Journal Retention, TTL, purge

**Run Retention Sweep**:
The step that applies Run Retention and then returns freed pages to the filesystem by incremental compaction. It runs at most once a day at the start of a Run or a Delivery Queue under a two-second budget, continuing at the next start when the budget runs out, and without a budget in the GC Command ([ADR-0255](docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md)).
_Avoid_: cleanup job, cron, vacuum

**Worktree Bootstrap**:
The configured command Roundfix runs once in a newly created Run or Task Worktree, after copying `worktree.copy` files and before Agent work and Verification, to prepare the environment (install dependencies, migrate and seed databases, warm caches). A bootstrap failure ends the Run or settles the Task with a bootstrap-failed outcome ([ADR-0034](docs/adr/0034-run-worktrees-can-bootstrap-their-environment.md)).
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

**Baseline Plan**:
The `roundfix/baseline-plan/v1` document a Baseline Command run produces before any mutation: the planned postimages, the canonical managed-entry ledger, the repository and catalog identities, every bounded preimage the plan reads or changes, and the History Relocations and History Citations it carries. Another clone of the same repository can apply it only when all of those match ([ADR-0071](docs/adr/0071-baseline-plans-are-portable-and-preimage-bound.md)).
_Avoid_: Change preview, pending-plan state, patch

**Plan Digest**:
The digest over a Baseline Plan that the operator confirms before apply and that binds the plan's preimages, warnings and History Citations; a plan whose digest no longer matches is not the plan that was confirmed ([ADR-0071](docs/adr/0071-baseline-plans-are-portable-and-preimage-bound.md)).
_Avoid_: File hash, approval token, repository checksum

**Setup Manifest**:
The setup-owned record of the selected profile, modules, decisions, and managed artifacts that reproduces and audits a repository's Context-Driven Baseline. Every applied Baseline change republishes it, so its recorded digests describe the bytes on disk rather than the bytes adoption wrote.
_Avoid_: Project Config, lock file, generated guide

**Managed Region**:
The byte range one setup-owned marker pair delimits inside an instruction carrier, identified by its managed artifact identity. Its content belongs to the Context-Driven Baseline, not to the repository holding it; the bytes outside every marker pair belong to the repository and no Baseline operation changes them.
_Avoid_: Managed file, generated block, setup section

**Managed Refresh**:
The preservation mode that regenerates only Managed Regions, leaves every other byte identical, requires no Source Baseline and no Decision Plan, and takes no root backup because the plan's preimages carry the preservation proof. It also lists History Citations, binds them into the Plan Digest, accepts repeatable `--promote` values, and converts Pending History on approval. Its preservation proof still covers the instruction carriers. It converges: a second run against an unchanged catalog reports the repository current and proposes no change ([ADR-0254](docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md), [ADR-0260](docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md)).
_Avoid_: Baseline Readoption, reapply, overwrite

**Unrecorded Managed Region**:
A Managed Region whose bytes are not the bytes the adopted Setup Manifest recorded. It is classified and reported by path, managed identity, and the lines a refresh removes — never presumed damaged and never silently replaced. Only a managed identity appearing more than once in one carrier blocks, because that refresh has no defensible target.
_Avoid_: Hand-edited marker, corrupted block, drift

**Baseline ADR**:
An architecture decision record whose reserved identity and invariant belong to the Context-Driven Baseline while project-specific notes remain repository-owned.
_Avoid_: Example ADR, project ADR, template copy

**Upgrade Retention Contract**:
The accounting a Context-Driven Baseline transition or Baseline Readoption must satisfy before mutation: every Source Baseline Normative Clause, recommendation, and Operational Contract maps to a current managed target, Repository-Specific Normative Rules, a recognized typed repository document, or an explicit rejection with a recorded reason.
_Avoid_: Best-effort migration, category coverage, silent rule removal

**Repository-Specific Normative Rules**:
Project-authored Normative Clauses that are not portable across the Context-Driven Baseline and cannot be represented by a typed project decision or managed semantic guide. Baseline assigns every representable rule to its semantic owner; only non-empty residual rules live outside setup markers in `docs/agents/specific-repository.md` and remain byte-preserved after confirmed adoption.
_Avoid_: Repository-Owned Extension, duplicate managed rule, baseline rule

**Formatter-Stable Output**:
Generated managed Markdown that the target repository's selected formatter leaves unchanged, so apply, formatting, Verification, audit, and reapply compose with no delta.
_Avoid_: Renderer-canonical output, format-after-apply fixup

**Format Command**:
The Project Config command `verification.format` that the Daemon runs, outside the Agent sandbox, over the files under a Spec's `qa/` directory before the repository Verification precondition (an imported pass) and before the QA Report commit, restoring the original bytes when it fails ([ADR-0249](docs/adr/0249-a-qa-step-formats-its-qa-directory-and-reopen-sees-a-late-dependency.md)).
_Avoid_: Repository Verification, Baseline formatter, whole-tree formatter

**Late Dependency**:
A Task in a completed QA gate's dependency closure that the Task Graph did not hold when the newest QA Report was committed; `roundfix reopen` returns the gate to `pending` over it whatever its status ([ADR-0249](docs/adr/0249-a-qa-step-formats-its-qa-directory-and-reopen-sees-a-late-dependency.md)).
_Avoid_: Stale dependency, missing dependency, unresolved dependency

**Agent Skill**:
A folder holding a `SKILL.md` and its bundled files that an Agent or the Supervisor activates on a declared trigger. Roundfix ships its own as Roundfix Skills and installs upstream ones as Vendored Skills.
_Avoid_: Plugin, prompt template, tool

**Roundfix Skill**:
A shipped agent skill that teaches an external Agent how to start Roundfix or how to resolve one assigned Batch.
_Avoid_: Runtime, Review Source, plugin

**Repository Skill Set**:
The complete set of Roundfix-owned and externally managed Agent Skills required by a repository workflow. Every required installed artifact matches its authoritative local source; unrelated extras do not join the set.
_Avoid_: Installed skills, plugin set, global skill set

**Vendored Skill**:
A skill a repository installs from an upstream source and never edits locally, because the next sync replaces its folder. A divergence is corrected in order: a rule that governs over the skill's default, then a repository-authored skill that supplies only what is missing, then a repository-authored replacement only when the vendored skill is useless or almost useless; an objective defect goes upstream as a pull request ([ADR-0264](docs/adr/0264-roundfix-never-replaces-a-vendored-skill-it-cannot-prove-older.md)).
_Avoid_: Forked skill, local skill copy, edited upstream skill

**Setup Snapshot Difference**:
A required external skill that is installed and matches its `skills-lock.json` entry but whose tree differs from the `treeDigest` its Setup Snapshot pins. Nothing offline tells whether the copy is newer or older, so Doctor names it without a warning and the managed refresh keeps it; only a deliberate per-skill restore replaces it ([ADR-0264](docs/adr/0264-roundfix-never-replaces-a-vendored-skill-it-cannot-prove-older.md)).
_Avoid_: Trailing skill, outdated skill, stale skill

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
The Live Run View pane that lists a Run's Work Items and their current status for both review and Spec Runs.
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
A read-only JSONL projection of one Run's Run Event Journal, selected by the identifier of an explicit Run and optionally followed until the Run reaches a terminal outcome. Its stable Supervisor categories are `task-status`, `batch`, `verification`, `outcome`, `agent-selection` and `usage`; replay includes every one by default, and a filter selects a subset of them. The command skips a record it cannot project with a warning and continues with the next journal entry.
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

**Skill Mode**:
A way of running the Roundfix skill that the user's request selects with the first word after the skill's name, `follow-up`, `status` or `unattended`, with or without a leading `--`; the skill reads it from the request text in every harness, it is not the CLI command `roundfix follow-up`, and a request whose first word names no mode uses the skill as before. [ADR-0320](docs/adr/0320-a-question-call-carries-up-to-three-independent-questions.md)
_Avoid_: Command flag, CLI mode, Follow Mode

**Follow Mode**:
The Live Run View state in which the timeline tail advances automatically as new Run Events arrive; suspended while the user scrolls back, resumed when the viewport returns to the bottom. Scrolling never affects the Run.
_Avoid_: tail mode, auto-scroll, live mode

**Pre-PR Review Policy**:
The repository's choice of Codex, Claude, CodeRabbit, or no reviewer before a
pull request is opened. An enabled review covers the current candidate, while
explicit none records an intentional omission that leaves QA and required
checks binding; provider failure is not an omission.
_Avoid_: Review mode, reviewer profile, automatic review

**Pre-PR Review Provider**:
The selected agent or service that performs review under the Pre-PR Review
Policy. An Agent Selection Profile supplies an agent provider's runtime and
model; explicit none is a policy choice with no provider.
_Avoid_: Review Source, Pre-PR Review Policy, reviewer model

**Pre-PR Review Record**:
The record of one pre-PR review, kept per checkout under its resolved path. Its
verdict is read from the reviewer's final message, as described in
[ADR-0272](docs/adr/0272-a-review-record-survives-a-fix-and-a-correction-inside-the-archived-specs-declared-paths-returns-to-review.md).
_Avoid_: Review Source Evidence, shared review state, reviewer transcript

**Reviewed Head Record**:
The copy of a Pre-PR Review Record kept under the head it reviewed, written for
every verdict that is not a reuse and never replaced by a later round at
another head. A Delivery Retry's correction proof reads it for the parked
candidate, as described in
[ADR-0272](docs/adr/0272-a-review-record-survives-a-fix-and-a-correction-inside-the-archived-specs-declared-paths-returns-to-review.md).
_Avoid_: round record, review history

**QA Archive Override**:
An explicit user request or applicable prior authorization to archive a covered
Spec despite its QA prerequisite being unmet. The recorded exception preserves
actual QA and Task evidence while retaining non-QA completion, source integrity
and separate delivery gates.
_Avoid_: QA pass, QA bypass, Implement Clean

**Clarification Record**:
The `## Decisions` section of a PRD that says who settled its questions: one dated line for each question a person answered, or the autonomous-route line naming the source of the author's decisions ([ADR-0280](docs/adr/0280-a-guide-task-records-its-own-surface-fingerprint-a-prd-records-who-decided-and-the-glossary-is-held-to-its-own-rule.md)).
_Avoid_: Interview transcript, assumed approval, unrecorded clarification

**Glossary Declaration**:
The `## Glossary` section of a Spec's PRD or TechSpec that lists each domain term the Spec adds to or changes in the glossary, and each bolded phrase it declares not a domain term, or `None.`. A Task that declares the glossary file and names the term in its Verification writes each listed term.
_Avoid_: Vocabulary Contract, term list, glossary candidates

**Glossary Gap**:
A domain term a Spec introduces that the glossary does not carry as required: a bolded phrase the glossary lacks and the Glossary Declaration does not cover, a declared term no Task plans to write, or a declared term still missing after its Tasks completed. The Spec Consistency Check reports it, and the QA gate and the Archive Command refuse a Spec that has one.
_Avoid_: Missing term, undocumented token, stale glossary

**Light Tier**:
The dispatch tier that runs a Task on an open model through OpenCode and OpenRouter when its complexity is low, its type is not qa, and it declares no Governed Path; every other Task runs on the standard tier. A light Task whose first Verification fails escalates its one feedback turn to its category's Preferred Selection. A light runtime that offers no mode for the requested access policy is not chosen, and a light session refused before Agent work yields to the category's Preferred Selection ([ADR-0314](docs/adr/0314-a-light-tier-that-cannot-serve-the-run-yields-to-the-category-profile.md)).
_Avoid_: Cheap model, Jev routing, light model profile

**Light Spend Log**:
The per-UTC-month record in Roundfix Home of the cost each Light Tier session reports, summed across every repository and compared with the Light Tier's monthly ceiling. A month that reached the ceiling, or whose log cannot be read, runs light Tasks on their category's profile.
_Avoid_: Judge Log, OpenRouter usage, spend cache

**Jev Ceiling**:
The monthly spending ceiling of the Jev judge, read from User Config only and compared with the Judge Log's sum across every repository on the machine; when it is unset a built-in default applies. The judge's key must carry a monthly credit limit no higher than it.
_Avoid_: Judge budget, key limit, Project Config ceiling

**Stage Key**:
The Roundfix environment variable an OpenRouter stage reads first, the judge's or implementation's own key, before the shared OpenRouter key for Roundfix; the generic OpenRouter variable is never read. Each stage records the variable it used by name, never by value.
_Avoid_: API key value, generic OpenRouter key, shared key

**Run Marker**:
The `ROUNDFIX_RUN_ID` variable naming the Run in each of its Agent Session and Verification processes. While it is non-empty, the judge sends nothing and reports the Run as its skip reason, whatever keys the command holds ([ADR-0283](docs/adr/0283-a-run-marker-keeps-the-judge-silent-inside-a-run.md)).
_Avoid_: Judge Key Variables, Run Worktree location, Roundfix Home

**Judge Key Variables**:
The credential variables withheld from every Agent Session and Verification command of a Run: the judge's declared key variables and the generic TypeSafe key `TYPESAFE_API_KEY`. The Run Marker also keeps the judge silent if a runtime restores a key ([ADR-0283](docs/adr/0283-a-run-marker-keeps-the-judge-silent-inside-a-run.md), amending [ADR-0273](docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md)).
_Avoid_: Stage Key, implementation key, credential value

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
An outside-evidence QA row blocked only because the Run sandbox denied network access, recorded with the host it could not reach. Like the pre-PR pull request row, it never decides a qualifying partial; whoever needs that proof declares it under Unreachable Acceptance.
_Avoid_: Failed outside evidence, skipped row, environment pass

**Network-Prohibited Row**:
An outside-evidence QA row the gate did not look up because the Spec's own authorization forbids network access, recorded as `blocked (environment: network prohibited by Spec authorization)` with the outside-evidence row named in its provenance. Like the Network-Denied Row, it never decides a qualifying partial; any other wording still needs a declaration or the override.
_Avoid_: Network-denied row, authorization override, skipped source

**Pinned Byte Invariant**:
A QA requirement that holds a file byte-identical to its bytes at a named commit instead of at the Delivery Base. It fails once a Spec that merges first changes the file, and `roundfix spec check` reports it as `SC-QA-BYTES-PINNED`.
_Avoid_: Byte invariant, golden pin, frozen file check

**Pre-PR Review Command**:
The command that obtains the Pre-PR Review Provider's review of the current candidate before a pull request opens, reads the verdict from the reviewer's final message and keeps the Pre-PR Review Record. A finding parks delivery only after it is validated.
_Avoid_: Review Run, PR feedback Watch, review command alias

## Adopter decisions

The Context-Driven Baseline records these decisions in the repositories it configures. They describe those repositories, not Roundfix, whose own stack and identifiers follow its ADRs.

**HTTP Contract Decision**:
The repository-owned choice of REST or POST-only application API semantics together with explicit protocol or operational exceptions and their owners.
_Avoid_: HTTP profile, universal REST rule, inferred route style

**Better Auth Base Path**:
The path under which a repository mounts Better Auth's routes, the value of Better Auth's `basePath` option, recorded as the optional `auth.base_path` decision. When it is recorded, the provider's route exception and the matching HTTP Contract Decision exception cover that path followed by `/*`; when it is absent, the recorded exception stands, and Better Auth's own default is `/api/auth` ([ADR-0263](docs/adr/0263-the-better-auth-base-path-is-a-repository-decision-a-baseline-update-records.md)).
_Avoid_: Auth prefix, API prefix, auth scope

**Internal Identifier**:
A technical identity for an entity or resource that is generated and controlled by the project.
_Avoid_: External identifier, natural key, business code

**Frontend Layout Decision**:
The repository's recorded choice of frontend layout, `systems` or `repository-defined`. While none is recorded, the Baseline states the suggested `systems` layout; a recorded value is never overwritten ([ADR-0205](docs/adr/0205-a-repository-records-its-frontend-layout-and-an-unrecorded-layout-follows-the-suggestion.md)).
_Avoid_: Frontend architecture rule, mandated layout

**Standard TypeScript Monorepo Profile**:
The opinionated, project-agnostic Context-Driven Baseline profile for the repository's standard TypeScript monorepo stack.
_Avoid_: Project-specific profile, generic TypeScript profile, sample template
