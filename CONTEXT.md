# Roundfix

Roundfix coordinates repeated review-resolution cycles for pull requests. This glossary defines the product language used to describe that loop.

## Language

**Run**:
A durable attempt to clean one pull request by coordinating review rounds until it reaches a terminal outcome.
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

**Review Source**:
The external review system that produces feedback for an Open Pull Request.
_Avoid_: Review Provider, Agent, ACP Runtime

**ACP Runtime**:
An ACP-compatible local coding runtime that Roundfix launches through the user's installed tool and authentication setup.
_Avoid_: Review Source, review provider

**Max Rounds**:
The configured number of Review Source rounds after which a Run is considered sufficiently reviewed for the developer's final merge, squash, or rebase decision.
_Avoid_: Budget, timeout, token cap

**Max Rounds Reached**:
A terminal Run outcome where the configured review round policy is complete, even if unresolved Review Issues remain for developer judgment.
_Avoid_: Failure, timeout, budget exceeded

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

**Unresolved Review Issue**:
A Review Issue that has been downloaded but has not reached a terminal local outcome.
_Avoid_: Open issue, pending task

**Batch**:
A bounded subset of Review Issues assigned to one agent invocation.
_Avoid_: Chunk, group, task

**Final Push**:
The Run-ending push that sends the PR Head Branch after no Unresolved Review Issues remain.
_Avoid_: Batch push, round push, agent push

**Resolve Command**:
The command that runs Agents over downloaded unresolved Review Issues for an Open Pull Request.
_Avoid_: Fix Command, Fetch command, watch command

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
The TUI view that shows Review Issues and streaming Agent output while a Run is active.
_Avoid_: Dashboard, report, log file

**Daemon**:
The Roundfix process that owns the Run lifecycle and Review Source-facing outcomes.
_Avoid_: Orchestrator, controller, manager

**Agent**:
The local coding assistant invoked by Roundfix to triage and resolve an assigned Batch.
_Avoid_: Review Source, review provider, worker, bot
