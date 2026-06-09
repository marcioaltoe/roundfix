---
title: Roundfix MVP creation
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
created_at: 2026-06-09
---

# Roundfix MVP creation

## Problem Statement

Developers who use CodeRabbit and local coding agents still have to run the
pull request review loop by hand. They wait for review results, inspect
unresolved findings, copy work into an agent, verify changes, commit, push,
wait for the next review, and repeat. That manual loop is slow and error-prone:
agents can be given too much authority, pushes can happen before all Unresolved
Review Issues are terminal, and repeated review findings can produce duplicate
work across Rounds.

Roundfix needs an MVP that turns this into a narrow, durable, local-first Go
CLI. The MVP must keep the Daemon in charge of Run state, Review Source access,
verification gates, commit policy, and Final Push, while giving the Agent only a
bounded Batch of Review Issues to triage and resolve.

## Solution

Build the first operational Roundfix MVP around three commands: `fetch`,
`resolve`, and `watch`. The MVP should support CodeRabbit as the first Review
Source and local ACP-capable Agents, starting with Codex as the first end-to-end
runtime path.

The user should be able to run a watch command for an Open Pull Request and let
Roundfix repeat the review-resolution loop until the pull request is clean,
`max_rounds` is reached, a Stop Request occurs, the Run Budget is exceeded, or a
hard failure stops the Run. Roundfix should persist markdown Round and Review
Issue artifacts in the configured Artifact Directory, store durable Run state in
the central Run Database, and push only after no Unresolved Review Issues remain.

## User Stories

1. As a developer, I want Roundfix to detect my current repository, branch, HEAD
   SHA, upstream, and pull request, so that I don't have to provide values that
   can be inferred safely.
2. As a developer, I want Preflight Validation before any long-running work, so
   that setup problems fail quickly with actionable instructions.
3. As a developer, I want Preflight Validation to check required config, git
   state, GitHub access, Review Source availability, Agent availability, and
   Artifact Directory writability, so that a Run doesn't start from a broken
   environment.
4. As a developer, I want Preflight Validation failures to avoid creating Run
   state or artifacts, so that failed setup attempts don't pollute history.
5. As a developer, I want dirty worktree checks outside the Artifact Directory,
   so that Roundfix doesn't mix agent changes with unrelated local work.
6. As a developer, I want dirty worktree errors to show changed paths and next
   steps, so that I know what to commit, stash, or remove before retrying.
7. As a developer, I want User Config and Project Config, so that shared defaults
   and repository-specific policy can be managed outside long command lines.
8. As a developer, I want config precedence to be deterministic, so that CLI
   flags override project settings, project settings override user settings, and
   user settings override built-in defaults.
9. As a developer, I want the Artifact Directory to default locally but be
   configurable, so that each repository can choose where markdown Round
   artifacts live.
10. As a developer, I want Roundfix to validate the Artifact Directory before
    fetching Review Source issues, so that filesystem errors fail before network
    or Agent work starts.
11. As a developer, I want Run state stored centrally, so that future daemon
    modes can track Runs across repositories.
12. As a developer, I want only one Active Run per Head Repository and PR Head
    Branch, so that two Roundfix loops cannot mutate the same pull request head
    at the same time.
13. As a developer, I want simultaneous Active Runs for different PR Head
    Branches, so that unrelated pull requests can be cleaned independently.
14. As a developer, I want `fetch` to create a tracked Fetch Run, so that
    downloaded artifacts have durable history without starting an Agent.
15. As a developer, I want `fetch` to download unresolved CodeRabbit findings,
    so that I can inspect Review Issues locally before deciding to resolve them.
16. As a developer, I want `fetch` to persist one Round of markdown Review Issue
    artifacts, so that Review Source feedback is available to later commands.
17. As a developer, I want `fetch` to stop without committing, pushing, or
    starting an Agent, so that fetching is a safe read-mostly operation.
18. As a developer, I want `resolve` to process already downloaded Compatible
    Artifacts, so that resolution work is separate from Review Source fetching.
19. As a developer, I want `resolve` to fail Preflight Validation when no
    Compatible Artifacts exist, so that it doesn't silently fetch new work.
20. As a developer, I want `resolve` without a Round selector to process all
    downloaded Unresolved Review Issues for the Open Pull Request, so that older
    unresolved Rounds don't get skipped.
21. As a developer, I want `resolve` with a Round selector to limit work to that
    Round, so that I can focus on a bounded artifact set when needed.
22. As a developer, I want repeated Review Issues deduplicated by Review Issue
    Fingerprint, so that the Agent resolves the newest occurrence instead of
    repeating stale work.
23. As a developer, I want duplicate newest-occurrence selection to be
    deterministic, so that Roundfix never picks an issue nondeterministically.
24. As a developer, I want ambiguous duplicate selection to fail during Preflight
    Validation, so that I can fix artifacts instead of trusting a guess.
25. As a developer, I want older duplicate occurrences marked `duplicated` by
    the Daemon, so that local artifact bookkeeping is consistent.
26. As a developer, I want duplicated older occurrences to avoid separate Review
    Source resolution, so that Roundfix doesn't mutate stale source threads.
27. As a developer, I want each Agent invocation to receive only a bounded Batch,
    so that fixes stay reviewable and failures are isolated.
28. As a developer, I want the Agent to read assigned issue files, triage each
    issue, edit code, update tests, update assigned issue statuses, and run
    verification, so that each Batch has a clear completion contract.
29. As a developer, I want the Agent to avoid commits, pushes, and Review
    Source mutations, so that the Daemon owns git and GitHub policy.
30. As a developer, I want Roundfix to verify each successful Batch before a
    commit, so that broken fixes don't enter the local history.
31. As a developer, I want one local commit per successful Batch, so that each
    accepted Batch has an auditable boundary.
32. As a developer, I want successful Batch commits to be allowed while other
    Unresolved Review Issues remain, so that progress is not blocked by later
    Batches.
33. As a developer, I want Final Push to run only after no Unresolved Review
    Issues remain, so that partial cleanup is not pushed as a completed state.
34. As a developer, I want Final Push to include the complete local PR Head
    Branch state, including commits that were already unpushed when the Run
    started, so that the remote branch matches local intent.
35. As a developer, I want Roundfix to reject auto-push without a known upstream
    remote and branch, so that pushes are explicit and predictable.
36. As a developer, I want no force-push in the MVP, so that Roundfix cannot
    rewrite remote branch history.
37. As a developer, I want `watch` to wait for CodeRabbit to finish reviewing
    the current HEAD, so that Roundfix fetches findings for the right commit.
38. As a developer, I want `watch` to use a quiet period after a Review Source
    signal, so that it avoids fetching before review comments settle.
39. As a developer, I want `watch` to repeat fetch and resolve Rounds until the
    pull request is clean or the configured policy stops, so that I don't have to
    babysit the review loop.
40. As a developer, I want `MaxRoundsReached` to be a non-error terminal outcome,
    so that reaching the configured review round policy hands control back
    without implying the Run failed.
41. As a developer, I want any remaining Unresolved Review Issues shown when
    `MaxRoundsReached` occurs, so that I can make the final merge, squash,
    rebase, or manual review decision.
42. As a developer, I want a Run Budget with `max_run_duration`, so that
    unidentified infinite loops or runaway resource usage stop safely.
43. As a developer, I want CodeRabbit review timeouts to produce a clear terminal
    failure, so that I can decide whether to trigger another review manually.
44. As a developer, I want Roundfix to offer a manual `@coderabbitai review`
    action rather than post it automatically, so that review allowance remains
    under my control.
45. As a developer, I want a Stop Request to stop the Active Run cleanly, so that
    I can end work without pretending a Batch succeeded.
46. As a developer, I want Stop Request handling to terminate an active Agent
    gracefully, then kill it after a short grace period if needed, so that
    shutdown is bounded.
47. As a developer, I want Stop Request handling to preserve Agent-created
    worktree changes and report changed paths, so that I can inspect unfinished
    work manually.
48. As a developer, I want a stopped Run to release its Active Run lock, so that
    a later command can start a new Run after Preflight Validation passes.
49. As a developer, I want `resolve` and `watch` to create new Runs rather than
    resume stopped Runs, so that Stop Request semantics stay simple.
50. As a developer, I want CLI exit codes that distinguish success, Run failure,
    Preflight Validation failure, and SIGINT, so that automation can react
    correctly.
51. As a developer, I want missing inputs to open Interactive Input unless
    `--no-input` is used, so that the CLI works both interactively and in
    automation.
52. As a developer, I want Interactive Input to suggest the current or remembered
    pull request and Agent, so that repeated Runs are faster to start.
53. As a developer, I want a Live Run View that shows Run status, Review Issues,
    Agent output, verification output, git state, and keybindings, so that I can
    monitor active work from one terminal screen.
54. As a developer, I want the Review Issue sidebar to group issues by Round,
    severity, status, file, and line, so that I can understand the current Batch
    and remaining work.
55. As a developer, I want the streaming console panel to show Agent and
    verification output, so that I can diagnose failures without opening
    separate logs.
56. As a developer, I want Roundfix to persist Agent output per Run, so that I
    can inspect what happened after the terminal session ends.
57. As a developer, I want CodeRabbit text treated as untrusted input, so that
    reviewer-provided commands or prompt content cannot escape the assigned
    Review Issue contract.
58. As a developer, I want tokens and authorization headers redacted from logs,
    so that local diagnostics do not leak credentials.
59. As a developer, I want a shipped `roundfix-watch` skill, so that agents know
    to start and observe Roundfix instead of scraping review comments manually.
60. As a developer, I want a shipped `roundfix-resolve-round` skill, so that
    child Agents receive durable instructions for one bounded Batch.
61. As a developer, I want the MVP to keep product names and generated artifacts
    specific to Roundfix, so that no reference-project branding leaks into the
    repository.
62. As a maintainer, I want the MVP implemented in small, testable milestones, so
    that each slice can be reviewed without waiting for the entire watch loop.

## Implementation Decisions

- Build a local-first Go CLI with a thin command entry point and internal
  packages for CLI parsing, configuration, preflight, git state, GitHub access,
  Review Source access, Round artifacts, Run storage, Agent runtime execution,
  watch orchestration, TUI views, and shipped skills.
- Keep the initial operational command set to `fetch`, `resolve`, and `watch`.
  Reserve `reprocess` for a later command instead of adding broad
  `include_resolved` behavior to MVP commands.
- Use CodeRabbit as the first Review Source. The Review Source boundary must
  fetch unresolved review threads, map source identifiers into Review Issues,
  expose watch status for the current PR HEAD, and resolve source threads only
  after Daemon-owned verification succeeds.
- Treat reviewer text as untrusted input. Roundfix must never execute
  reviewer-provided commands, shell-interpolate review bodies, or expose tokens
  in logs.
- Use local ACP-capable Agent runtimes through the user's installed tools and
  authentication. Codex is the first end-to-end runtime path. Runtime command
  overrides and probes must produce actionable install or authentication hints.
- Keep the Daemon responsible for Review Source status checks, CodeRabbit
  fetches, Round creation, Agent lifecycle, verification orchestration, commit
  policy, Final Push, Review Source resolution, retry behavior, timeouts, and
  Stop Request handling.
- Keep the Agent responsible only for assigned issue files, triage, code edits,
  test updates, verification commands, and assigned Review Issue status updates.
  The Agent must not commit, push, resolve Review Source threads, or mark issues
  as `duplicated`.
- Implement Preflight Validation as a required gate before any wait, fetch,
  Agent run, commit, push, or Review Source mutation. Preflight failures must
  exit with the setup-error code and must not create Run Database records,
  Active Run locks, markdown artifacts, or diagnostic Run events.
- Support YAML User Config and Project Config. Built-in defaults, User Config,
  Project Config, and CLI flags must apply in that order. Final Push remains
  config-controlled for the MVP instead of an interactive one-off toggle.
- Resolve the Artifact Directory from config or default behavior. Validate it as
  a writable directory before fetching Review Source issues. Roundfix must not
  edit ignore files or decide whether artifacts are tracked by Git.
- Store durable Run state in a central SQLite Run Database under Roundfix Home.
  The Run Database owns Active Run locks, historical Run records, remembered
  interactive defaults, and future multi-repository daemon state.
- Enforce one Active Run per Head Repository and PR Head Branch. Allow
  simultaneous Active Runs for different PR Head Branches.
- Make `fetch` a tracked Fetch Run. It resolves the Open Pull Request, validates
  the Artifact Directory, fetches CodeRabbit issues, persists markdown artifacts,
  records Run history, and exits without starting an Agent, committing, or
  pushing.
- Make `resolve` operate only on downloaded Compatible Artifacts. Compatible
  Artifacts must match Head Repository, PR Head Branch, pull request number, and
  optional Round selector. If none match, fail Preflight Validation and tell the
  user to run `fetch` or `watch`.
- Deduplicate repeated unresolved Review Issues across Compatible Artifact
  Rounds before assigning Agent Batches. Prefer source thread identity when
  present; otherwise use a provider-specific fingerprint. The newest occurrence
  wins by Round, source review submission time, and Round creation time. If that
  order cannot choose one issue, fail Preflight Validation.
- Mark older duplicate occurrences as terminal `duplicated` only after the
  assigned newest occurrence reaches `resolved` or `invalid`. Associate each
  older duplicate occurrence to the newest issue and keep the source mutation
  local-only for older duplicates.
- Treat `pending`, `valid`, `invalid`, `resolved`, `duplicated`, and `failed` as
  the MVP Review Issue statuses. A Review Issue is terminal only when it is
  `resolved`, `invalid`, or `duplicated`.
- Treat a Batch as successful only when every assigned Review Issue is terminal
  and the required verification command succeeds.
- Create one local commit per successful Batch when auto-commit is enabled.
  Successful Batch commits may happen while other Unresolved Review Issues
  remain.
- Run Final Push only when auto-push is enabled, auto-commit is enabled, the
  upstream remote and branch are known, the local branch has commits not present
  on the target branch, and no Unresolved Review Issues remain.
- Use a single explicit push shape from local HEAD to the PR Head Branch. Do not
  force-push in the MVP.
- Implement `watch` as a restartable foreground loop that waits for CodeRabbit
  status on the current HEAD, observes a quiet period, fetches unresolved
  Review Issues, persists a Round, resolves Batches, verifies, commits, performs
  Final Push when allowed, and repeats until a terminal outcome.
- Treat `Clean`, `Fetched`, `MaxRoundsReached`, and clean `Stopped` as
  non-error terminal outcomes. Treat setup failures, Run failures, timeouts,
  budget exhaustion, verification failures, Agent failures, Review Source
  failures, and Final Push failures as distinct error outcomes.
- Implement Stop Request handling as a terminal path that stops active Agent
  work, persists available output, preserves uncommitted worktree changes, shows
  changed paths, releases the Active Run lock, and avoids any further
  verification, commits, pushes, fetches, or Review Source mutations.
- Implement Interactive Input for missing command parameters and explicit
  interactive mode. Run available Preflight Validation before Interactive Input
  and full Preflight Validation afterward.
- Implement the Live Run View as the MVP TUI monitoring surface. It must show
  Run metadata, pipeline state, Review Issues, streaming Agent output,
  verification output, budget status, git state, and keybindings without hiding
  auto-push policy.
- Ship `roundfix-watch` as the user-facing skill and `roundfix-resolve-round` as
  the internal child-agent skill. The shipped skills must use Roundfix language
  only and must preserve the Daemon-owned git and Review Source boundaries.
- Keep the first delivery milestone order narrow: input and preflight, git and
  PR detection, CodeRabbit fetch, Agent Batch resolution, commit and Final Push,
  durable watch loop, then skills installer support.

## Testing Decisions

- Tests should assert external behavior and contracts rather than internal
  implementation details. Prefer command output, exit codes, persisted artifact
  content, Run state transitions, fake boundary calls, and visible TUI state over
  private helper behavior.
- Use the existing CLI runner tests as the highest current seam for command
  parsing, stdout, stderr, and exit code behavior. Extend this seam as command
  behavior becomes operational.
- Add focused unit tests for config loading, config precedence, semantic
  validation, duration parsing, and invalid value diagnostics.
- Add focused unit tests for Artifact Directory resolution, path validation,
  writability checks, and failure messages.
- Add unit tests for markdown Round and Review Issue parsing, writing,
  frontmatter preservation, status transitions, terminal status detection, and
  duplicate bookkeeping.
- Add unit tests for Review Issue Fingerprint selection, newest duplicate
  ordering, older duplicate association, and ambiguous duplicate failure.
- Add unit tests for Run Database migrations, Run creation, Fetch Run terminal
  state, Active Run locking, lock release on terminal outcomes, and remembered
  interactive defaults.
- Add unit tests for git state parsing: git root, branch, HEAD, upstream,
  unpushed commits, dirty worktree, and changed path reporting.
- Add unit tests for GitHub and PR metadata parsing using fake command output or
  fake API clients, not live network calls.
- Add unit tests for CodeRabbit mapping, including REST review comments,
  GraphQL threads, bot filtering, resolved-thread filtering, source thread
  identity, review hashes, and source review submission timestamps.
- Add unit tests for Agent runtime probing, command override selection, startup
  failure diagnostics, graceful termination, forced kill after the grace period,
  and output persistence.
- Add unit tests for child-agent prompt construction to confirm assigned issue
  lists, forbidden actions, verification requirements, and untrusted reviewer
  text handling.
- Add state-machine tests for `fetch`, `resolve`, and `watch` terminal outcomes,
  including `Fetched`, `Clean`, `MaxRoundsReached`, `BudgetExceeded`, `TimedOut`,
  `Failed`, and `Stopped`.
- Add integration tests with a fake Review Source, fake Agent runtime, fake git
  runner, and temporary Artifact Directory so the MVP loop can be verified
  without GitHub or CodeRabbit network access.
- Add integration tests proving Preflight Validation blocks waits, fetches,
  Agent starts, commits, pushes, Review Source mutations, Run records, and
  Active Run locks.
- Add integration tests proving `fetch` creates a tracked Fetch Run and never
  starts an Agent, commits, or pushes.
- Add integration tests proving `resolve` reuses downloaded Compatible Artifacts,
  processes all Rounds by default, honors a Round selector, and fails when no
  Compatible Artifacts exist.
- Add integration tests proving duplicated unresolved Review Issues create one
  Agent assignment, mark older occurrences as `duplicated` after success, and do
  not resolve older source threads separately.
- Add integration tests proving successful Batch commits are allowed while other
  Unresolved Review Issues remain and Final Push is blocked until all remaining
  Review Issues are terminal.
- Add integration tests proving Stop Request handling preserves Agent-created
  worktree changes, reports changed paths, exits cleanly through Roundfix
  controls, releases the Active Run lock, and blocks the next Run if preserved
  changes leave the worktree dirty.
- Add tests for exit code mapping: non-error terminal outcomes exit `0`, Run
  failures exit `1`, Preflight Validation failures exit `2`, and SIGINT exits
  `130`.
- Add TUI model and snapshot tests for Interactive Input defaults, missing input
  errors, Live Run View layout, Review Issue sidebar grouping, streaming console
  output, status strip, git strip, and keybinding labels.
- Keep manual tests for the real GitHub and CodeRabbit path: long review waits,
  review timeout, manual review trigger prompt, dirty worktree rejection, failed
  verification, Final Push, auto-push disabled by config, and Stop Request during
  an active Agent run.
- Use the broad local verification gate before claiming an MVP slice is ready.
  Run race tests when concurrency, Agent process ownership, watch polling, or Run
  Database locking changes.

## Out of Scope

- Building a general workflow engine, PRD/task/spec system, or CI healer.
- Supporting Review Sources other than CodeRabbit in the MVP.
- Adding `reprocess` or any broad `include_resolved` behavior to `fetch`,
  `resolve`, or `watch`.
- Letting Agents commit, push, resolve Review Source threads, mark duplicated
  issues, or mutate GitHub review-thread state directly.
- Pushing after each issue, Batch, or Round.
- Force-pushing or rewriting remote branch history.
- Automatically editing ignore files or deciding whether Artifact Directory
  files are tracked by Git.
- Automatically installing ACP runtimes or falling back to package-manager
  install commands.
- Asking for model-vendor API keys when the selected Agent can use the user's
  local authenticated setup.
- Enforcing token budgets, output-size budgets, or repeated-fingerprint budgets
  as hard MVP gates.
- Running only from webhooks. Polling must remain the fallback and can be the
  MVP default.
- Automatically posting `@coderabbitai review` without explicit user opt-in.
- Building a persistent background service as a separate operating-system daemon.
  The MVP can run the Daemon-owned loop inside a foreground CLI or TUI process.
- Supporting broad multi-user hosted operation. The MVP is local-first.
- Shipping names, branding, comments, examples, or generated artifacts from
  reference projects.

## Further Notes

- This PRD is synthesized from the current Roundfix glossary, product brief, and
  accepted architecture decisions.
- The current implementation is an early scaffold: help and version output work,
  while `fetch`, `resolve`, and `watch` are reserved and not implemented.
- The MVP should stay KISS. Each milestone should deliver one externally
  verifiable behavior before broadening the next boundary.
- Follow-up work should convert this PRD into tracer-bullet implementation
  issues in dependency order.
