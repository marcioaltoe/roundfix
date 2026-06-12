# Product Brief

## Purpose

Build a focused Go tool that watches GitHub pull requests, imports CodeRabbit review feedback, uses local coding runtimes to resolve valid issues, commits successful batches, pushes the complete branch state only after no Unresolved Review Issues remain, and repeats until CodeRabbit reports the PR is clean, a configured review round limit is reached, or a budget safeguard stops the run.

This should not be a full workflow/orchestration system. It should be a narrow, durable review-resolution loop.

## Product Goal

Given an open pull request with CodeRabbit review comments, the tool should:

1. Run Preflight Validation and fail early with actionable instructions when the environment is not ready.
2. Detect the current repository, branch, PR, upstream, and HEAD SHA.
3. Wait for CodeRabbit to finish reviewing the current HEAD.
4. Fetch unresolved CodeRabbit review threads.
5. Persist those findings as local round artifacts.
6. Spawn a selected Agent runtime, such as Codex, Claude Code, or OpenCode, to resolve one bounded batch.
7. Verify the batch.
8. Create one local commit per successful Batch.
9. Push the complete branch state only after no Unresolved Review Issues remain.
10. Wait for CodeRabbit's next review.
11. Repeat until clean, timed out, stopped, budget-limited, or `max_rounds` is reached.

## Recommended Language

Use Go.

The hard parts are orchestration, subprocess management, GitHub API calls, git state, config, SQLite state, polling, retries, and a terminal UI. Go is a better fit for shipping the first product quickly and keeping the code operationally simple.

Rust is a good fit only if the project prioritizes strict state-machine modeling over iteration speed. For the MVP, choose Go.

## Non-Goals

- Do not build a full PRD/task/spec workflow system.
- Do not build a general CI healer in v1.
- Do not support every Review Source in v1.
- Do not let child agents directly mutate GitHub review-thread state.
- Do not push per issue, Batch, or Round.
- Do not hide git pushes or GitHub mutations behind ambiguous agent behavior.
- Do not mention internal reference project names in repository docs, code, comments, tests, or generated artifacts.
- Do not treat `max_rounds` as a resource-consumption safeguard.
- Do not enforce token, output-size, or repeated-issue budgets in the MVP.
- Do not provide an `include_resolved` option on `resolve` or `watch`.

## Core Design Principle

The daemon owns the loop. The agent owns only the code change.

The daemon should control:

- Review Source status checks
- CodeRabbit source fetches
- Round creation
- Agent process lifecycle
- Verification gate orchestration
- Commit and push policy
- CodeRabbit thread resolution
- Retry, timeout, and stop behavior

The child agent should control:

- Reading assigned issue files
- Triage of valid versus invalid findings
- Code edits
- Test updates
- Verification commands

## MVP Command

```bash
roundfix watch --source coderabbit --pr 123 --agent codex \
  --until-clean --max-rounds 6
```

## CLI Commands

```bash
roundfix init [--scope project|user]
roundfix --version
roundfix -v
roundfix fetch --source coderabbit --pr 123
roundfix resolve --pr 123 --agent codex
roundfix watch --source coderabbit --pr 123 --agent codex --until-clean
roundfix skills check
roundfix skills install
```

The MVP has three operational commands and two support command groups:

- `fetch`: download unresolved Review Source issues for an Open Pull Request and persist markdown artifacts.
- `resolve`: iterate over downloaded unresolved Review Issues for an Open Pull Request and resolve them with the configured Agent.
- `watch`: automate `fetch` and `resolve` while monitoring Review Source publication across configured review rounds.
- `init`: create User Config or Project Config. If `--scope` is omitted, prompt for the scope and default to Project Config.
- `skills`: validate or install the shipped `roundfix` skill into the current project or Codex, Claude Code, and OpenCode-compatible skill directories.

`roundfix fetch` creates a tracked Fetch Run in the Run Database. It resolves the Open Pull Request, validates the Artifact Directory, fetches Review Source issues, persists markdown artifacts, and then stops without starting an Agent, committing, or pushing.

`roundfix resolve --pr <n>` works only when the pull request is open. It does not fetch new Review Source issues. By default, it reuses existing downloaded markdown artifacts for that pull request, reads all downloaded Unresolved Review Issues from all Compatible Artifact Rounds, processes them in batches as needed, and updates Run history in the Run Database. Passing `--round <n>` limits resolution to that Round.

Compatible Artifacts match the requested Head Repository, PR Head Branch, and pull request number. If `--round <n>` is passed, they must also match that Round. If no Round is passed, all Compatible Artifact Rounds are included.

If no Compatible Artifacts exist for the requested pull request and Round scope, `resolve` must fail during Preflight Validation before creating a new Run. The error should tell the user to run `roundfix fetch --source coderabbit --pr <n>` or use `roundfix watch`.

When `resolve` includes multiple Compatible Artifact Rounds, Roundfix must deduplicate Review Issues before assigning batches. Deduplication uses the Review Issue Fingerprint: prefer `source_ref` when present, otherwise use a provider-specific fingerprint such as `review_hash`. For duplicated unresolved Review Issues, process only the newest occurrence and associate older occurrences to that newest issue.

KISS resolve loop:

- If a Review Issue belongs to the current pull request, is unresolved, and is not duplicated, resolve it.
- If a Review Issue is duplicated, resolve the newest occurrence and mark older occurrences locally as `duplicated`.
- When marking an older occurrence as `duplicated`, associate it to the newest occurrence with `duplicate_of`.
- A successful Batch may be committed even if other Unresolved Review Issues remain.
- Final Push must not run while any Unresolved Review Issues remain.

Newest occurrence resolution:

1. Higher `round` wins.
2. If the duplicate occurrences are in the same Round, later `source_review_submitted_at` wins.
3. If `source_review_submitted_at` is missing or tied, later `round_created_at` wins.
4. If the newest occurrence is still ambiguous, fail Preflight Validation with an ambiguous artifacts error instead of choosing nondeterministically.

Future command:

- `reprocess`: explicitly revisit selected Terminal Review Issues.

`roundfix reprocess` is separate from `resolve` and `watch`, and should be implemented after the MVP. It never runs implicitly and must target specific Terminal Review Issues by issue id, issue range, pull request selector, source ref, or another explicit selector. Roundfix must not provide a broad `include_resolved` option that silently mixes Terminal Review Issues into normal resolution.

## Current Go Project Layout

```text
cmd/roundfix/main.go

internal/cli/                 # CLI parsing, output, and exit behavior
internal/tui/                 # Interactive Input and Live Run View rendering
internal/watch/               # Durable review loop state machine
internal/reviewsource/        # Review Source interface
internal/reviewsource/coderabbit/ # CodeRabbit implementation
internal/agent/               # Agent runtimes: codex, claude, opencode
internal/runevent/            # Run Event type, sink seam, and fanout
internal/rounds/              # Markdown issue artifacts and frontmatter
internal/store/               # Global SQLite Run Database
internal/config/              # User and project config loading and validation
internal/preflight/           # Git, PR, worktree, config, and push safety checks
internal/daemon/              # Verification, commits, source resolution, and push
skills/                       # Agent skills shipped with the tool
docs/                         # User docs after MVP stabilizes
```

## Current Implementation Choices

- CLI: Go standard library `flag` and explicit command dispatch.
- Interactive Input and Live Run View: prompt collection plus ACP cockpit-style terminal rendering owned by `internal/tui`.
- State: SQLite Run Database
- SQLite driver: `modernc.org/sqlite`
- Config: YAML user and project config files through `gopkg.in/yaml.v3`
- GitHub: `gh` subprocess for REST, GraphQL, PR metadata, review status, and thread resolution boundaries.

Cobra remains a future option if the CLI grows enough to justify a framework. The Live Run View now uses Bubble Tea and Lipgloss because ACP streaming needs a real terminal render loop.

## Review Source Interface

Start with a small Review Source interface.

```go
type FetchRequest struct {
    Source          string
    PRNumber        string
    BaseRepository  string
    HeadRepository  string
    HeadBranch      string
    HeadSHA         string
    IncludeNitpicks bool
}

type ReviewItem struct {
    Title       string
    File        string
    Line        int
    Severity    string
    Author      string
    Body        string
    SourceRef  string

    ReviewHash              string
    SourceReviewID          string
    SourceReviewSubmittedAt time.Time
}

type ResolvedIssue struct {
    FilePath  string
    Status    string
    SourceRef string
}

type Source interface {
    FetchReviews(ctx context.Context, req FetchRequest) ([]ReviewItem, error)
}

type WatchStatusRequest struct {
    Source         string
    PRNumber       string
    BaseRepository string
    HeadRepository string
    HeadBranch     string
    HeadSHA        string
}
```

## CodeRabbit Requirements

The CodeRabbit Review Source should:

- Use GitHub repository metadata from the current checkout.
- Fetch PR metadata to get the current head SHA.
- Fetch pull request review comments.
- Fetch GraphQL review threads.
- Map REST review comments to GraphQL threads via comment database IDs.
- Accept both `coderabbitai[bot]` and `coderabbitai`, because GitHub REST and GraphQL can expose the CodeRabbit actor differently.
- Skip resolved review threads.
- Preserve `thread:<id>,comment:<id>` in `source_ref`.
- Resolve source threads only after the child agent completed and local issue files are terminal.
- Duplicated older occurrences are local artifact bookkeeping; do not resolve their Review Source threads separately.
- The current MVP fetches unresolved inline CodeRabbit review threads. Review-body summaries and outside-diff comments remain future parser work.
- Future review-body parsing should hash review-body comments so resolved nitpicks or outside-diff findings do not reappear in later rounds unless the source review is newer.
- Use `source_ref` or `review_hash` as the Review Issue Fingerprint for deduplicating repeated findings across Rounds.

## Watch Loop

The watch loop should be durable and restartable.

State machine:

```text
Idle
WaitingForReview
FetchingIssues
PersistingRound
ResolvingWithAgent
VerifyingRound
Committing
Pushing
WaitingForNextReview
Clean
Fetched
MaxRoundsReached
BudgetExceeded
TimedOut
Failed
Stopped
```

`MaxRoundsReached` is a non-error terminal state. It means Roundfix completed the configured review round policy and is handing control back to the developer for the final merge, squash, rebase, or manual review decision. If unresolved Review Issues remain, Roundfix must show them clearly, but CLI commands should exit successfully unless a separate hard error occurred.

`Fetched` is a terminal state for Fetch Runs. It means Roundfix fetched Review Source issues and persisted markdown artifacts without entering the review-resolution loop.

## CLI Exit Codes

Exit codes should make automation distinguish setup problems from Run failures.

- `0`: non-error terminal outcomes such as `Clean`, `Fetched`, `MaxRoundsReached`, and intentional `Stopped` when shutdown completes cleanly through the TUI or another Roundfix control.
- `1`: failures after a Run has started, including `Failed`, `TimedOut`, `BudgetExceeded`, verification failure, Agent failure, Review Source timeout, and Final Push failure.
- `2`: failures before a Run is created, including Preflight Validation failures, config syntax errors, invalid config values, missing required input, invalid command input, missing authentication, unavailable selected Agent command, invalid Artifact Directory, and dirty worktree rejection.
- `130`: SIGINT/Ctrl-C interruption.

## Stop Behavior

A Stop Request is user intent, not a successful batch outcome.

Rules:

- A Stop Request moves the Active Run toward `Stopped`.
- If no child process is active, stop immediately and mark the Run `Stopped`.
- If an Agent is active, send the runtime's graceful cancel or terminate signal.
- Wait up to `10s` for the Agent to exit.
- If the Agent is still running after `10s`, kill the process.
- Persist available Agent output before finishing shutdown.
- Preserve any uncommitted worktree changes left by the Agent.
- Show changed paths and statuses after shutdown so the user can inspect them.
- After a Stop Request, do not start another Agent, run verification, create commits, push, fetch more Review Source issues, or resolve Review Source threads.
- Do not revert, stash, stage, or commit Agent changes automatically after a Stop Request.
- `Stopped` is terminal and releases the Active Run lock.
- Do not resume a `Stopped` Run.
- A later `resolve` or `watch` command creates a new Run after Preflight Validation passes.
- If preserved Agent changes left the worktree dirty, Preflight Validation should reject the new Run before it is created.
- A clean stop requested through the TUI or another Roundfix control exits with `0`.
- SIGINT/Ctrl-C should use the same shutdown path where possible, but the CLI exits with `130`.

## Preflight Validation

Roundfix should fail first. Validate as much as possible before the user waits, before starting a long TUI flow, before fetching Review Source issues, and before launching an Agent.

Preflight order:

1. Parse User Config and Project Config.
2. Resolve command flags and inferable values.
3. Run all validations that do not depend on missing required parameters.
4. If required parameters are still missing, open Interactive Input.
5. Run full Preflight Validation again with the collected parameters.
6. Only then start Review Source waits, issue fetches, Agent runs, commits, or Final Push.

Preflight must validate:

- Config file syntax and semantic ranges.
- Git root detection.
- Current branch and HEAD.
- Open Pull Request existence and state.
- PR Head Branch and Head Repository.
- Active Run uniqueness for the same Head Repository and PR Head Branch.
- Artifact Directory resolution and writability.
- Compatible Artifacts for `resolve`.
- Run Database creation/opening/migration.
- GitHub authentication and minimum API access.
- Review Source availability for `fetch` and `watch`.
- Selected Agent command availability for `resolve` and `watch`.
- Upstream remote and branch when Final Push is enabled.
- Dirty worktree status outside the Artifact Directory.

Preflight failure output must include:

- The failed check.
- The detected value when useful.
- The exact file, config key, command, branch, or path involved.
- What the user should do before running Roundfix again.
- Confirmation that Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

Preflight failures:

- Do not create Run Database records.
- Do not create an Active Run.
- Do not create or reserve an Active Run lock.
- Do not persist diagnostic failure events in `roundfix.db`.
- Do not write markdown artifacts.

Dirty worktree behavior:

- Fail before waiting or fetching when unrelated uncommitted changes exist outside the Artifact Directory.
- Show the changed paths and statuses.
- Tell the user to commit, stash, or remove those changes before continuing.
- Do not stage, commit, stash, or revert user changes automatically.
- Do not treat Roundfix-owned markdown changes inside the Artifact Directory as a dirty-worktree blocker.
- Do not treat Project Config changes at `<repo>/.roundfixrc.yml` as a dirty-worktree blocker.
- Exclude Project Config changes from Batch commits so local setup changes do not mix with review fixes.

Loop logic:

1. Run Preflight Validation.
2. Resolve repo, branch, PR, upstream, and HEAD.
3. Resolve and validate the Artifact Directory.
4. Wait until CodeRabbit has reviewed or settled the current HEAD.
5. Wait a quiet period after the Review Source signal.
6. Fetch unresolved review issues.
7. If no issues remain, mark the watch clean.
8. Persist issues into the next round directory.
9. Spawn a child agent run for the round.
10. Wait for the child run to complete.
11. Verify all issue files are terminal.
12. Commit each successful Batch locally when `auto_commit` is enabled.
13. If no Unresolved Review Issues remain and `auto_push` is enabled, push the complete local branch state, including commits that were already unpushed when the Run started.
14. Repeat if `until_clean` is enabled.

## Handling Slow CodeRabbit Reviews

Waiting is a normal state, not a failure.

Defaults:

```yaml
watch:
  poll_interval: 30s
  review_timeout: 30m
  quiet_period: 30s
  max_rounds: 6
  until_clean: true
  auto_commit: true
  auto_push: true
```

`max_rounds` is the review-completion policy: after this many Review Source rounds, Roundfix stops waiting for additional review-resolution cycles and leaves the pull request ready for the developer's final merge, squash, or rebase decision. It is distinct from the Run Budget, which exists only to stop unidentified infinite loops or runaway resource consumption.

## Run Budget

The MVP Run Budget has one hard gate: `max_run_duration`.

`max_run_duration` limits total wall-clock time for one Run. If the limit is reached, Roundfix stops the Run with `BudgetExceeded`. This safeguard exists to prevent unidentified infinite loops or runaway agent usage. It is separate from `max_rounds`, which describes how many Review Source rounds are enough before the developer makes the final merge, squash, or rebase decision.

Token usage, child-agent output size, and repeated Review Issue fingerprints may be recorded as diagnostics when available, but they are not hard budget gates in the MVP.

Rules:

- Use webhooks if available, but always keep polling as a fallback.
- Poll GitHub status/check data for CodeRabbit status on the current HEAD.
- Compare CodeRabbit review commit SHA with the current PR head SHA.
- If CodeRabbit is pending, show `WaitingForReview`.
- If the review looks complete, wait `quiet_period` before fetching comments.
- If `review_timeout` is reached, stop the round with `TimedOut`.
- Offer a manual action to post `@coderabbitai review`.
- Do not automatically post `@coderabbitai review` unless the user opted in, because it may consume review allowance.
- If CI is slow, allow users to delay review triggering until CI checks complete.

## Git Safety

Pushes must be boring and explicit.

Requirements:

- Detect current branch.
- Detect current HEAD SHA.
- Detect dirty worktree.
- Detect upstream remote and branch.
- Detect unpushed commit count.
- Reject auto-push unless a remote and branch are known.
- Do not push at Run start.
- Final Push is enabled by default.
- The Final Push sends the complete local branch state, including commits that were already unpushed when the Run started.
- Disable Final Push explicitly through User Config or Project Config.
- Use only this push shape:

```bash
git push <remote> HEAD:<branch>
```

- If `auto_push` is enabled, require `auto_commit` to be true.
- Create one local commit per successful Batch when `auto_commit` is enabled.
- Commit creation is allowed even when other Unresolved Review Issues remain.
- Do not push after each Batch or Round.
- Push only when no Unresolved Review Issues remain for the current pull request.
- Before pushing the final result, verify the local branch has commits not present on the target remote branch.
- Do not force-push in v1.
- Reject unrelated uncommitted user changes before waiting for Review Source work.
- Do not push if unrelated uncommitted user changes exist.
- The dirty-worktree error should tell the user to commit, stash, or remove the listed changes before running Roundfix again.

## Round Artifacts

Use local markdown files for review rounds.

By default, Roundfix stores markdown artifacts in `<git-root>/.roundfix/`. The Artifact Directory must be configurable from User Config or Project Config so a repository can choose a different controlled folder, such as `~/.roundfix/artefacts` or `<repo>/docs/reviews`. The configured path always identifies a directory that Roundfix owns for markdown Round and Review Issue artifacts.

The Artifact Directory does not store the Run Database. If the user chooses a repository path for the Artifact Directory, only markdown artifacts belong there.

Artifact Directory resolution:

- Empty `artifact_dir`: use `<git-root>/.roundfix/`.
- Absolute `artifact_dir`: expand `~` and use the path as configured.
- Relative `artifact_dir`: resolve against `<git-root>`.
- No Git root: require an absolute `artifact_dir` or fail commands that depend on a pull request.

Artifact Directory validation:

- Resolve the final Artifact Directory before fetching Review Source issues.
- Create the directory if it does not exist.
- Reject the path if it exists and is not a directory.
- Test that Roundfix can write to the directory, such as by creating and deleting a temporary validation file.
- If validation fails, show the concrete filesystem error and do not fetch pull request issues.

Artifact Directory git tracking:

- Roundfix does not decide whether Artifact Directory files are tracked or ignored by Git.
- Roundfix must not edit `.gitignore` automatically.
- Roundfix does not warn about Artifact Directory files being tracked, ignored, or untracked.
- The user or repository owner decides whether to add the Artifact Directory to `.gitignore`, version it, or manage it another way.
- Each successful `fetch` with automatic Round selection reuses an existing matching Round when the same HEAD already has the same Review Issue fingerprints.
- If the fetched payload is new, automatic Round selection writes the next Round directory. Roundfix does not upsert Round artifacts in place.
- An explicit existing `--round <n>` is rejected instead of overwritten.
- Repeated findings across different Rounds are handled during `resolve` by Review Issue Fingerprint deduplication, not by deleting old artifacts at fetch time.

```text
<artifact-dir>/
  reviews/
    pr-123/
      round-001/
        issue_001.md
        issue_002.md
      round-002/
        issue_001.md
```

## Run Database

The global Roundfix directory is `~/.roundfix/`. The Run Database is named `roundfix.db` and lives at `~/.roundfix/roundfix.db`. It stores Run state and review progress centrally so future daemon modes can track multiple Runs across different repositories at the same time.

Concurrency rules:

- The MVP may run multiple Active Runs at the same time when they target different PR Head Branches.
- Roundfix must reject a new Run when another Active Run already targets the same Head Repository and PR Head Branch.
- The Run Database should enforce a logical lock keyed by Head Repository owner, Head Repository name, and PR Head Branch.
- The logical lock exists only while the Run is active.
- Terminal outcomes, including `Stopped`, release the logical lock.
- Terminal Runs remain historical records and are not resumed by `resolve` or `watch`.
- For fork pull requests, the lock uses the Head Repository, not the base repository.
- If a duplicate Run is rejected, show the existing `run_id`, Head Repository, PR Head Branch, base repository, pull request, and current state.

Issue file example:

```markdown
---
source: coderabbit
pr: "123"
round: 1
round_created_at: 2026-06-08T12:00:00Z
status: pending
file: apps/api/src/auth.go
line: 88
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kw...,comment:PRRC_kw...
review_hash: ""
duplicate_of: ""
source_review_id: "987654"
source_review_submitted_at: 2026-06-08T12:00:00Z
---

# Issue 001: Missing error handling

## Review Comment

<original CodeRabbit comment>

## Triage

- Decision: `UNREVIEWED`
- Notes:
```

Allowed statuses:

```text
pending
valid
invalid
resolved
duplicated
failed
```

A Review Issue is terminal for the current Round only when its status is `resolved`, `invalid`, or `duplicated`. A Batch succeeds only when every assigned Review Issue is terminal. The daemon may resolve Review Source threads for `resolved` and `invalid` Review Issues after verification passes, because valid fixes and invalid false positives should not reappear indefinitely. Duplicated older occurrences are marked locally and do not resolve Review Source threads separately.

When multiple unresolved Review Issues share the same Review Issue Fingerprint across Compatible Artifact Rounds, the daemon assigns only the newest occurrence to an Agent. Newest means highest Round, then latest `source_review_submitted_at`, then latest `round_created_at`. If that order still cannot choose one occurrence, Preflight Validation fails with an ambiguous artifacts error. Older duplicate occurrences remain as historical context and should not create separate Agent work.

After the assigned newest occurrence reaches `resolved` or `invalid`, the daemon marks older duplicate occurrences as `duplicated` and sets `duplicate_of` to the newest occurrence. The Agent must not mark issues as `duplicated`.

## Run Events

Roundfix records Run activity as Run Events so the Run Database can answer what happened, in what order, for which Run, Batch, Review Issue, and tool call.

- Every meaningful Run occurrence becomes a Run Event: Run state transitions, daemon decisions, Agent stream output, verification milestones, commit and Final Push decisions, Review Source resolution, Stop Requests, and terminal outcomes.
- Run Events live in the Run Event Journal inside the Run Database, ordered by a per-Run monotonic cursor that readers treat as an opaque replay position.
- Producers publish Run Events through one sink seam owned by `internal/runevent`. The Agent log, the Run Event Journal, and the Live Run View are sink adapters; producers never know which adapters are attached.
- The journal sink is critical: once a Run has started, a journal append failure fails the Run. Live Run View and log sinks are best-effort and must never block or fail producers.
- Agent Run Event payloads store the raw ACP session update JSON (ADR 0008). Readers skip unknown event kinds instead of failing, so journals stay readable across ACP Runtime versions.
- The Live Run View renders live Run Events and replayed Run Events through the same timeline renderer, so Attach shows the same cockpit as a live Run.
- Attach is non-mutating: it replays the Run Event Journal and follows new Run Events. Detach leaves the Run active. Stop Request remains the only way to end an Active Run.

## Child Agent Contract

The child agent receives a strict prompt and a bounded list of issue files.

The child agent must:

1. Read every assigned issue file completely.
2. Triage each issue.
3. Set `status: valid` or `status: invalid` and record triage notes.
4. For valid issues, make production-quality code changes.
5. Add or update tests when behavior changes.
6. Set issue files to `status: resolved` only after changes and verification.
7. Run real repository verification.
8. Leave commit creation to the daemon.
9. Never push.
10. Never call `gh` to resolve CodeRabbit threads.
11. Never edit issue files outside the assigned batch.
12. Never set `status: duplicated`; duplicated status is daemon-owned.

The daemon resolves Review Source threads after the batch succeeds for assigned `resolved` and `invalid` issues.

## Agent Runtime Compatibility

Public CLI and TUI surfaces should label the selected coding runtime as `Agent`. Use `ACP Runtime` in implementation and technical docs only when the protocol boundary matters.

MVP runtimes:

```text
codex       -> codex-acp, fallback npx --yes @zed-industries/codex-acp
claude      -> claude-agent-acp, fallback npx --yes @agentclientprotocol/claude-agent-acp
opencode    -> opencode acp
```

Future runtimes:

```text
cursor-agent
gemini
```

Agent runtime interface:

```go
type RuntimeSpec struct {
    ID              string
    DisplayName     string
    Protocol        string
    Command         string
    Args            []string
    ProbeArgs       []string
    Fallbacks       []RuntimeLauncher
    DefaultModel    string
    Model           string
    SupportsAddDirs bool
    BootstrapModel  bool
    FullAccessMode  string
}
```

Requirements:

- Probe runtime availability before starting a run.
- Show actionable install hints.
- Keep model selection runtime-specific.
- Treat Agent runtime startup as using the user's local installed tool, login, subscription, and model-vendor credentials.
- Do not ask for model-vendor API keys when the selected Agent runtime can use the user's local authenticated setup.
- Support command overrides because ACP adapter command names can differ by installation.
- Use explicit adapter fallback commands only for known ACP runtimes, matching the selected Agent and never silently switching to another Agent.
- If the selected Agent command fails to probe or start, show the concrete command, error, and install/authentication hint, then stop.
- Let the user choose another Agent explicitly if they want to retry with a different local setup.
- Drive Codex, Claude, and OpenCode through their ACP stdio protocol, not by streaming Markdown prompts directly into a JSON-RPC server.
- Publish ACP session updates as Run Events through one event sink boundary that feeds the Agent log, the Run Event Journal, and the Live Run View.
- Keep the raw ACP session update JSON as the durable Run Event payload; bounded text summaries serve list rendering.
- Support headless streaming logs through a writer sink adapter, not a separate output path inside the runner.
- Persist agent output per run.
- Treat reviewer text as untrusted input.

## TUI Requirements

The MVP terminal UI has two responsibilities:

1. Collect optional parameters for `fetch`, `resolve`, and `watch`.
2. Show a readable Live Run View for Fetch, Resolve, and Watch Runs.

Interactive input requirements:

- Run all possible Preflight Validation before opening Interactive Input.
- Open Interactive Input when a required command parameter is missing.
- Open Interactive Input when the user passes `--interactive`.
- If all required parameters are available through flags, config, or inference, run the command directly.
- `--no-input` disables Interactive Input and fails with a clear error when required parameters are missing.
- Run full Preflight Validation after Interactive Input and before waiting for Review Source work.
- `fetch` should ask for pull request number, Review Source, and Round.
- `resolve` should ask for pull request number, Round, Artifact Directory override, concurrent jobs, batch size, Agent, model override, additional directories, reasoning effort, dry run, and auto commit.
- `watch` should ask for pull request number, Review Source, Agent, max rounds, and max run duration in the MVP. More advanced knobs can be added when the corresponding runtime behavior is implemented.
- Auto push is resolved from User Config or Project Config and shown as read-only Run state in the TUI.
- When possible, infer the current pull request number and suggest it.
- Remember the last pull request number and Agent in the Run Database and suggest them on the next run.
- Prefer configured Agent defaults from User Config or Project Config before falling back to remembered values.
- Use concise UI labels that match the glossary. For example, use `Round` instead of `Review Round`.
- In `fetch`, an empty or `auto` Round creates the next available Round.
- In `resolve`, Round defaults to `all`, which processes every downloaded Unresolved Review Issue across all Compatible Artifact Rounds for the pull request.
- In `resolve`, setting Round to a number limits processing to that Round.
- In `watch`, Round selection is automatic and the user configures `max_rounds`.

Live Run View requirements:

- Header: command, repo, PR, PR Head Branch, Review Source, Agent, HEAD.
- Target block: PR, repository, branch, source, Agent, and HEAD.
- Run block: run ID, state, Round progress, budget, git state, auto-commit, auto-push, and last push.
- Split pane: Review Issues on the left and Agent Console on the right.
- Review Issues pane: downloaded or assigned Review Issues grouped by Round, severity, status, file, line, and issue title.
- Agent Console pane: ACP session timeline including assistant text, tool calls, plan/status updates, daemon messages, and verification output where available.
- Use Roundfix branding only.
- In a TTY, render the Agent stream through an ACP cockpit-style Bubble Tea view; outside a TTY, print a readable text stream.

Current MVP Live Run View example:

```text
Roundfix fetch

Target:
  PR: #123 owner/project
  Branch: feature/auth
  Source: CodeRabbit
  Agent: Codex
  HEAD: abc123

Run:
  ID: run_20260610T001344Z_639014c1cb9ac528
  State: FetchingIssues
  Round: 0 / 6
  Budget: 0s / 2h0m0s
  Git: clean, 0 unpushed commit(s), upstream origin/feature/auth
  Auto-commit: on
  Auto-push: off
  Last push: disabled

+-----------------------------------------+-----------------------------------------------+
| Review Issues                           | Agent Console                                 |
+-----------------------------------------+-----------------------------------------------+
| none                                    | Fetching Review Source issues...              |
+-----------------------------------------+-----------------------------------------------+
Keys: Ctrl-C stop
```

A future full-screen TUI may add keyboard focus, detach, manual fetch, manual resolve, push, and Review Source trigger controls.

## Skill To Ship

Ship one `roundfix` agent skill with the Go tool.

### User-Facing Mode

Purpose: starting and observing the tool.

Trigger phrases:

- resolve CodeRabbit comments
- watch this PR
- run roundfix until clean
- clean up review bot feedback

Instructions:

- Prefer `roundfix` commands over manual GitHub scraping.
- Inspect current repo and PR if needed.
- Start `roundfix watch`.
- Prefer daemon or TUI mode for long waits.
- Report the run ID, PR, Review Source, Agent, and current state.
- Do not manually resolve CodeRabbit threads unless the tool is unavailable.

### Assigned Batch Mode

Purpose: internal child-agent guidance for one bounded batch.

Instructions:

- Read the assigned issue files.
- Triage each as valid or invalid.
- Resolve valid issues.
- Update issue frontmatter.
- Run verification.
- Never create commits.
- Do not push.
- Do not call Review Source-specific GitHub mutations.

Suggested skill layout:

```text
skills/
  roundfix/
    SKILL.md
    agents/openai.yaml
```

## Config Example

Roundfix supports User Config and Project Config files. User Config applies across repositories. Project Config applies inside one repository. Built-in defaults set `auto_commit = true` and `auto_push = true`. Disabling Final Push is an explicit opt-out in User Config or Project Config, not an interactive prompt.

Paths:

- User Config: `~/.roundfix/config.yml`
- Run Database: `~/.roundfix/roundfix.db`
- Project Config: `<repo>/.roundfixrc.yml`

`roundfix init` creates these files. If `--scope` is omitted, Roundfix asks where to write the file and defaults to Project Config when the user presses Enter.

Project Config discovery:

1. Detect the Git root with `git rev-parse --show-toplevel`.
2. Load `<git-root>/.roundfixrc.yml` if it exists.
3. If no Git root exists, do not load Project Config automatically.

Precedence:

1. Built-in defaults.
2. User Config.
3. Project Config.
4. CLI flags for exposed command parameters.

`watch.auto_push` is config-only in the MVP. Do not expose a one-off CLI flag or Interactive Input toggle for disabling Final Push.

```yaml
defaults:
  agent: codex
  model: ""
  auto_commit: true
  verification: make verify
  artifact_dir: ""

review_source:
  name: coderabbit
  include_nitpicks: true

watch:
  until_clean: true
  max_rounds: 6
  poll_interval: 30s
  review_timeout: 30m
  quiet_period: 30s
  auto_push: true
  push_remote: ""
  push_branch: ""

budget:
  enabled: true
  max_run_duration: 2h

resolve:
  batch_size: 3
  concurrent: 1
```

## Implemented MVP Capabilities

### Interactive Input

- Add Preflight Validation before and after Interactive Input.
- Collect optional parameters for `fetch`, `resolve`, and `watch`.
- Suggest current or remembered pull request number when available.
- Suggest configured or remembered Agent when available.
- Start the selected command with the collected parameters.

### Git and PR Detection

- Detect repo, branch, HEAD, upstream, dirty status, and unpushed commits.
- Detect current PR with `gh pr view`.

### CodeRabbit Fetch

- Fetch unresolved CodeRabbit review threads.
- Persist markdown issue files.
- Show them in the Live Run View and fetch summary.

### Agent Batch Resolve

- Spawn the selected Agent runtime.
- Pass assigned issue files.
- Persist logs.
- Verify issue statuses after child completion.
- Render live Agent and verification output in the Live Run View.
- Render Review Issues grouped by Round.

### Commit and Push

- Add one commit per successful Batch.
- Add explicit final auto-push.
- Enforce push safety.

### Durable Watch Loop

- Poll CodeRabbit status.
- Handle quiet period and timeout.
- Repeat until clean or max rounds.

### Skills Installer

- Install `roundfix`.
- Install to `<repo>/.agents/skills/roundfix` by default.
- Support explicit Codex, Claude Code, and OpenCode Agent skill directories.
- If `<repo>/.claude/skills` exists, ask whether to create `.claude/skills/roundfix` as a symlink to the project-local skill.

## Testing Strategy

Unit tests:

- Review Source parsing.
- GraphQL thread mapping.
- CodeRabbit author normalization across REST and GraphQL actor names.
- Round artifact parsing and writing.
- Artifact Directory resolution and validation.
- Active Run uniqueness by Head Repository and PR Head Branch.
- Fetch Run state transitions.
- Resolve command iterates over all downloaded Unresolved Review Issues across all Compatible Artifact Rounds for an Open Pull Request by default.
- Resolve command reuses existing downloaded artifacts and fails Preflight Validation when none match the requested pull request and Round scope.
- Resolve command deduplicates repeated Review Issues across Compatible Artifact Rounds by Review Issue Fingerprint before batching.
- Resolve command chooses the newest duplicate by Round, `source_review_submitted_at`, then `round_created_at`.
- Resolve command fails Preflight Validation when duplicate newest occurrence selection is ambiguous.
- Duplicated Review Issues use terminal `status: duplicated`, assigned only by the daemon.
- Duplicated older occurrences are local-only and are associated to the newest occurrence with `duplicate_of`.
- `postponed` is not an MVP Review Issue status.
- Git state parsing.
- Preflight Validation ordering and failure messages.
- Preflight Validation does not create Run Database records on failure.
- CLI exit code mapping.
- Stop Request state transitions and Agent termination grace period.
- Stop Request leaves Agent-created worktree changes intact.
- `Stopped` releases the Active Run lock.
- `resolve` and `watch` create a new Run instead of resuming a `Stopped` Run.
- Watch state transitions.
- Config validation.
- TUI input defaults for current PR, remembered PR, configured Agent, and remembered Agent.
- Split-pane Live Run View rendering with Review Issues on the left and Agent Console on the right.

Integration tests:

- Fake Review Source.
- Fake agent runtime.
- Fake git runner.
- Invalid Artifact Directory blocks Review Source fetch.
- Preflight failure blocks Review Source waits, issue fetches, Agent runs, commits, and pushes.
- Preflight failure does not create a Run Database record or Active Run lock.
- Preflight failure exits with `2`; Run failure exits with `1`; non-error terminal outcomes exit with `0`.
- Intentional stop through Roundfix controls exits with `0`; SIGINT exits with `130`.
- Stopping during an Agent run waits up to `10s`, then kills the child process if needed, without commit, push, verification, fetch, or Review Source thread resolution.
- Stopping during an Agent run leaves uncommitted Agent changes intact and reports changed paths.
- A Run stopped with leftover Agent changes releases its lock, and the next `resolve` or `watch` attempt is blocked by dirty-worktree Preflight Validation.
- After a clean `Stopped` Run with a clean worktree, the next `resolve` or `watch` creates a new Run instead of resuming the stopped one.
- Duplicate Active Run for the same Head Repository and PR Head Branch is rejected.
- `roundfix fetch` creates a tracked Fetch Run and never starts an Agent.
- `roundfix resolve --pr <n>` updates all downloaded Unresolved Review Issues across all Compatible Artifact Rounds by default and does not fetch Review Source issues.
- `roundfix resolve --pr <n>` with no Compatible Artifacts exits with `2` before creating a Run.
- `roundfix resolve --pr <n>` deduplicates repeated unresolved Review Issues and assigns only the newest occurrence while preserving older Round references.
- Duplicated Review Issues across Compatible Artifact Rounds produce one Agent assignment, not one assignment per Round occurrence.
- Older duplicate Review Issue occurrences are marked `duplicated` after the assigned newest occurrence reaches `resolved` or `invalid`.
- Older duplicate Review Issue occurrences are associated to the newest occurrence with `duplicate_of`.
- Older duplicate Review Issue occurrences are local-only and do not resolve Review Source threads separately.
- `duplicated` is terminal and does not count as an Unresolved Review Issue for commit or Final Push gating.
- Final Push is blocked while any Unresolved Review Issues remain, but successful Batch commits are still allowed.
- Ambiguous duplicate newest occurrence selection exits with `2` before creating a Run.
- Full watch loop without network.
- TUI output tests for interactive inputs and Live Run View state.

Manual tests:

- Real GitHub PR with CodeRabbit.
- Long CodeRabbit wait.
- Timeout and manual trigger.
- Final auto-push default includes commits that were already unpushed when the Run started.
- Auto-push disabled through config.
- Auto-push runs only after all downloaded Unresolved Review Issues are resolved.
- Failed verification.
- Dirty worktree rejection.
- Stop during Agent run with leftover worktree changes.

Future tests:

- Review-body, outside-diff, and nitpick parsing after those CodeRabbit sources become first-class Review Issue artifacts.
- Full-screen TUI keyboard navigation after a real interactive TUI is implemented.

## Security Rules

- Treat all CodeRabbit text as untrusted input.
- Never execute reviewer-provided commands.
- Never shell-interpolate review bodies.
- Read only cited files unless broader context is needed for the change.
- Do not expose tokens in logs.
- Redact GitHub tokens and authorization headers.
- Keep GitHub mutations in daemon code, not agent instructions.

## Initial README Pitch

Roundfix is a local-first PR review cleanup tool. It watches CodeRabbit feedback, dispatches your local coding agent to resolve valid issues, commits successful batches, pushes after all downloaded Unresolved Review Issues are resolved, then waits for the next review round until the PR is clean.

It is designed for developers who already use Codex, Claude Code, or OpenCode and want the review-resolve-push-repeat loop to run safely without babysitting GitHub.

## Maintenance Recommendation

Keep the first product narrow and reliable:

1. Preserve the command boundaries: `fetch` downloads, `resolve` works over downloaded artifacts, and `watch` automates both.
2. Prefer the Go standard library until an added dependency removes real complexity.
3. Keep GitHub mutations in daemon-owned code.
4. Keep Agents limited to assigned issue files, code edits, tests, and verification.
5. Treat CodeRabbit text as untrusted input.
6. Add more Review Sources and runtime protocol adapters only after CodeRabbit plus Codex remains stable end to end.
