# Roundfix

Roundfix picks up Work Items — Review Issues from pull request reviews today, Tasks from Spec Task Graphs next — resolves them through the user's selected ACP Runtime, and pushes only when nothing unresolved remains. This glossary defines the product language for that loop.

## Language

**Run**:
A durable attempt to drive one target's Work Items — an Open Pull Request's Review Issues or a Spec's Tasks — to a terminal outcome. One Active Run is allowed per target: (Head Repository, PR Head Branch) for review work, (repository, spec slug) for spec work.
_Avoid_: Session, execution, job

**Fetch Run**:
A short Run that fetches Review Source issues and persists markdown artifacts without starting an Agent.
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

**Spec**:
One feature's planning artifact set produced by the spec workflow: PRD, Task Graph, Task files, and QA evidence.
_Avoid_: Feature folder, epic, project

**Task**:
One implementable unit of work within a Spec. Its task file is the sole owner of its status.
_Avoid_: Subtask, story, ticket

**Task Graph**:
The Spec's manifest that declares its Tasks and their dependencies as a directed acyclic graph. Dependencies live only in the Task Graph, never in task files.
_Avoid_: Task list, backlog, roadmap

**QA Report**:
The qa-gate evidence report written to a Spec's QA directory, carrying a machine-readable verdict in its frontmatter.
_Avoid_: Test report, QA log

**ACP Runtime**:
A local coding runtime that Roundfix launches through the user's installed tool and authentication setup using Agent Client Protocol stdio. The MVP supports Codex through `codex-acp`, Claude through `claude-agent-acp`, and OpenCode through `opencode acp`; command overrides remain a stdio escape hatch for local testing.
_Avoid_: Review Source, review provider

**Agent Session**:
The persistent acpx-backed session through which one Run drives its Agent across Work Items — created when the Run starts Agent work, named by the Run, and closed at the Run's terminal outcome.
_Avoid_: ACP session, chat, conversation, thread

**Merge-Ready**:
The state where the pull request's Review Source status check on the pushed head commit reports success with no new Review Issues, letting a watch Run end Clean.
_Avoid_: mergeable, green check, approved

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

**Run Budget**:
A safeguard that stops a Run before it can continue indefinitely and indirectly consume unbounded resources.
_Avoid_: Max rounds, review round limit

**Preflight Validation**:
The early checks Roundfix runs before starting a Run or work that would make the developer wait.
_Avoid_: Best-effort validation, late failure

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
The directory where Roundfix stores markdown Round and Review Issue artifacts.
_Avoid_: Workspace, cache, output folder

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

**Final Push**:
The Run-ending push that sends the PR Head Branch after no Unresolved Review Issues remain.
_Avoid_: Batch push, round push, agent push

**Resolve Command**:
The command that runs Agents over downloaded unresolved Review Issues for an Open Pull Request.
_Avoid_: Fix Command, Fetch command, watch command

**Implement Command**:
The command that executes a Spec's Task Graph by running Agents over its Tasks in dependency order.
_Avoid_: Run command, execute command, spec command

**Reprocess Command**:
An explicit future command for revisiting selected Terminal Review Issues.
_Avoid_: Include resolved, resolve option

**Init Command**:
The support command that creates User Config or Project Config before operational Runs.
_Avoid_: Bootstrap run, setup run

**Roundfix Skill**:
A shipped agent skill that teaches an external Agent how to start Roundfix or how to resolve one assigned Batch.
_Avoid_: Runtime, Review Source, plugin

**Interactive Input**:
The TUI flow that collects command parameters before a Run starts.
_Avoid_: Wizard, form, setup screen

**Live Run View**:
The TUI view that shows Review Issues and Run Events for a Run: streaming live while the Run is active, or replayed from the Run Event Journal during Attach.
_Avoid_: Dashboard, report, log file

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

**Attach**:
Viewing a Run by replaying its Run Event Journal and then following new Run Events, without owning, mutating, or stopping the Run.
_Avoid_: Resume, reconnect, takeover

**Agent**:
The local coding assistant invoked by Roundfix to triage and resolve an assigned Batch.
_Avoid_: Review Source, review provider, worker, bot

**Follow Mode**:
The Live Run View state in which the timeline tail advances automatically as new Run Events arrive; suspended while the user scrolls back, resumed when the viewport returns to the bottom. Scrolling never affects the Run.
_Avoid_: tail mode, auto-scroll, live mode
