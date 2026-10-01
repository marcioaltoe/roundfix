## A review only happens when it is asked for

These CodeRabbit request instructions apply only when the repository's Pre-PR
Review Policy selects `coderabbit`, or when the legacy PR-feedback workflow is
using `fetch`, `watch`, or `resolve` with CodeRabbit as its Review Source. A
repository whose Pre-PR Review Policy selects `codex`, `claude`, or `none` is
never asked for a CodeRabbit review by the pre-PR workflow.

For that scoped case, automatic CodeRabbit review is **off** by deliberate
configuration: automatic incremental review fires on every push and burns the
hourly allowance while a review is still being worked. The consequence is the
rule that is easiest to forget — **a pull request gets no review unless someone
requests one.**

Requesting it, either way works:

- add the `coderabbit:review` tag to the pull request description, or
- comment `@coderabbitai review` on the pull request.

`@coderabbitai review` is **incremental**: it covers only what changed since the
last review. Use `@coderabbitai full review` for a pass over the whole pull
request — after many incremental rounds, or when the earlier reviews may have
missed something.

Request again after any of these, because none of them triggers a review on its
own: commits pushed after the first review, a batch of fixes landing, or a
rebase that changes the head.

**A green check is not evidence of a review.** When the allowance is exhausted,
CodeRabbit posts a rate-limit comment and a check named `Review rate limited`
that **passes by design**, so it never blocks a merge on a protected branch. The
comment is the authoritative signal that no review ran. Reading that green check
as "reviewed" is how a pull request reaches `main` unreviewed.

The allowance is roughly ten pull request reviews per hour, shared across the
`marcioaltoe` and `gesttione-solutions` organizations, as a rolling window
rather than a quota that resets on the hour. Comment `@coderabbitai rate limit`
to see what remains without spending a review.

`review_source.include_nitpicks` defaults to `false`, so CodeRabbit findings
whose severity is `nitpick` do not become Review Issues unless User Config or
Project Config explicitly sets the key to `true`.

`review_source.request_review` defaults to `false`, and
`review_source.request_command` defaults to `@coderabbitai review`. Project
Config overrides User Config, which overrides these built-in defaults. When
`request_review` is enabled, `watch` and `resolve` publish one idempotent
review request for the pushed head after each Round's Final Push. `fetch`
never publishes a review request.

Publishing the request is not Review Source Evidence and does not advance the
Round. `watch` and `resolve` still wait for Evidence bound to the pushed head.
If the Review Source explicitly refuses the requested review, that refusal
ends the Run; Roundfix does not retry the request, back off, or wait for review
capacity.

Preflight Validation for `watch` and `resolve` reads `.coderabbit.yaml` and
defines `pushTriggersReview` as `auto_review.enabled` and
`auto_review.auto_incremental_review` being enabled with
`auto_review.auto_pause_after_reviewed_commits: 0`. A finite pause cannot
guarantee a review for every pushed head. Absent, unreadable, or omitted values
use the Review Source defaults, including the finite pause default of `5`.
Preflight exits `2` for either incoherent pair:

- `pushTriggersReview=false` with `review_source.request_review=false` would
  strand the Run waiting for a review nobody requests; set
  `review_source.request_review` to `true` in Project Config.
- `pushTriggersReview=true` with `review_source.request_review=true` would
  request a duplicate review after every push; set
  `review_source.request_review` to `false` in Project Config.

The refusal names all three `.coderabbit.yaml` values and
`review_source.request_review`, plus the Project Config change that repairs
the pair. `fetch` is exempt because it neither pushes nor requests a review.

## User-Facing Review Runs

1. Prefer `roundfix` commands over manual GitHub scraping.
2. Inspect the current repository and Open Pull Request only when Roundfix needs
   missing command input.
3. Start the watched loop with:

   ```bash
   roundfix watch --source coderabbit --pr <number> [--spec <slug>] --until-clean
   ```

4. Let Roundfix own Branch Integrity Preflight, Review Source waits,
   CodeRabbit fetches, Round creation, Agent lifecycle, verification, Batch
   commits, Final Push, Review Source resolution, Outcome Comments, retries,
   timeouts, and Stop Request handling.
5. Use the bounded `roundfix runs list` (Active Runs by default; widen with
   `--state all` or `--limit 0`) or the Run Browser (`roundfix attach` with
   no argument at an interactive terminal) when the Run ID was not captured.
6. Report the Run ID, Open Pull Request, Review Source, Agent, and current Run
   state whenever you summarize progress. Include Agent Model and Default
   Reasoning Effort when the Run starts Agent work.
7. For unattended waits, follow `roundfix events <run-id> --follow` and parse
   JSONL from stdout. Use the Live Run View for human inspection.

Useful commands:

```bash
roundfix fetch --source coderabbit --pr <number> [--spec <slug>]
roundfix resolve --pr <number> [--spec <slug>]
roundfix watch --source coderabbit --pr <number> [--spec <slug>] --until-clean
roundfix resolve --pr <number> [--spec <slug>] --detach
roundfix watch --source coderabbit --pr <number> [--spec <slug>] --until-clean --detach
roundfix implement --spec <slug>
roundfix implement --spec <slug> --detach
roundfix spec check
roundfix spec check <slug> --format json --strict
roundfix runs list
roundfix runs list --state all --limit 0
roundfix runs
roundfix attach
roundfix attach <run-id>
roundfix events <run-id>
roundfix events <run-id> --follow
roundfix events <run-id> --filter verification,outcome
roundfix settle --spec <slug> --task <task_id>
roundfix archive <slug>
roundfix release plan
roundfix release plan --impact <none|patch|minor|major> --reason "<classification reason>"
roundfix gc --dry-run
roundfix gc
roundfix stop --spec <slug>
roundfix stop --force --spec <slug>
roundfix setup --yes
roundfix setup --no-input
roundfix doctor
roundfix upgrade --check
roundfix skills list
roundfix skills check
```

## Review checkout and spec worktree isolation

Review Runs (`fetch`, `resolve`, and `watch`) execute in the user's checkout on
the checked-out PR Head Branch and create no Run Worktree. `fetch` starts no
Agent. `resolve` and `watch` start the Agent from the same checkout, so a
review fix is always a delta over the pull request branch that Final Push
updates.

Roundfix never checks out a branch or moves the working tree. Before any review
Run starts, Preflight Validation resolves the Open Pull Request's PR Head
Branch and validates the checkout against it. When the branches differ,
`fetch`, `resolve`, and `watch` refuse with exit `2`; the diagnostic names the
PR Head Branch and its revision and the checkout branch and its revision. The
refusal creates no Run and has no side effects: Roundfix does not fetch Review
Source issues, start an Agent, commit, push, or move the working tree.
Recover by placing the checkout on the named PR Head Branch with the printed
command shape, then rerun the review command:

```bash
git switch -- '<PR Head Branch>'
```

After a Run starts, Roundfix revalidates the recorded PR Head Branch and
expected revision before each Batch and write boundary. If the checkout moves,
the Run stops before that write and reaches the terminal `CheckoutMoved`
outcome. This is an environmental interruption, not a Review Issue failure or
a `Failed` Run: affected Review Issues stay unsettled so the unchanged work can
be retried after the operator restores the PR Head Branch checkout. Roundfix
does not restore or otherwise move the working tree itself.

Branch Integrity Preflight runs before any fetch, Agent Session, Review Source
comment, code change, commit, or push for `fetch`, `resolve`, and `watch`.

- The preflight enumerates pending `roundfix/run-*` Run Branch work and kept
  worktrees bound to the PR Head Branch. Fast-forwardable work is integrated
  automatically and journaled before the review Run continues. For a
  failed-cycle set proven from its QA Reports, preflight disregards only the
  branches proven superseded by the current evidence branch and leaves those
  superseded Git refs unchanged. The current evidence branch remains subject
  to normal automatic integration or refusal. Reclaim superseded branches
  separately with `roundfix reconcile --apply`.
- Non-fast-forward pending work refuses the command with exit `2`, names each
  pending Run Branch and worktree, and prints the recovery command
  `git merge --ff-only <branch>`.
- Another Active Run bound to the Head Repository and PR Head Branch refuses
  the command with exit `2` and names both `roundfix stop --run-id <id>` and
  `roundfix stop --force --run-id <id>`.
- `--skip-branch-integrity` is the only bypass. It skips pending Run Branch
  and Active Run guardrails only after Roundfix publishes a pull request audit
  comment naming the run id, actor, time, skipped guardrails, ignored pending
  work, and ignored Active Runs. If that comment cannot be published, the
  command fails preflight with exit `2`.
- `resolve` and `watch` also require a clean tracked working tree before Agent
  work starts. Dirty tracked paths refuse with exit `2`; untracked files are
  allowed because Batch commits stage only paths changed since the Batch
  snapshot. After a failed Batch, dirty tracked files in the checkout are
  Agent work by construction.

This is the only Branch Integrity Preflight relaxation: proven-superseded QA
Run Branch work no longer blocks a review Run. Non-fast-forward or ambiguous
pending work, preserved branch-set evidence, another Active Run, a dirty
tracked review checkout, and every other existing refusal still block.

Review Runs have no Integration Pending outcome. They either mutate the user's
checkout directly, stop before side effects through Preflight Validation, or
end with a review outcome such as Clean, CleanUnverified, MaxRoundsReached,
TimedOut, Failed, Stopped, ReviewSkipped, or Unresolved.

Review Source Evidence is bound to the expected head. `pending` means no usable
expected-head signal exists; `reviewing` means a current-head CodeRabbit check
or status is still in progress; `reviewed` means CodeRabbit produced a
current-head result that does not prove Merge-Ready; and `verified` requires a
recognised review-completed current-head CodeRabbit check or commit status, or
a current-head CodeRabbit `APPROVED` review, with zero unresolved CodeRabbit
threads. A stale signal never verifies the expected head. An explicit Review
Source refusal resolves to `skipped` evidence and never verifies a head;
`watch --until-clean` will not merge that head or clear it for merge. An
unrecognised signal resolves to `pending`, even when its check conclusion is
success: a green check is not evidence that a review ran. `failed` records an
explicit current-head Review Source failure.

A refusal is not a transient Review Source failure. Roundfix does not
automatically retrigger the review or retry a refused head; the follow-on work
owns that policy.

Roundfix retries only typed transient Review Source failures: a context
deadline not caused by Run cancellation, temporary DNS failure, connection
reset, HTTP `429`, or GitHub `5xx` response. One retry episode records
`started`, then `recovered` or `exhausted`. Retry sleeps reuse the existing poll
interval and remain bounded by the Review Source timeout and Run Budget; there
is no retry configuration or log-text inference.

Spec Runs (`implement`) keep worktree isolation because Task concurrency needs
it:

- `worktree.location` sets the parent directory with Project Config > User
  Config > built-in default precedence. The built-in default is
  `~/.roundfix/worktrees`.
- Each spec Run Worktree is created at
  `<worktree.location>/<repo-slug>/<run-id>` on a Run Branch named
  `roundfix/run-<id>`. The Run row records the path as `work_dir`.
- Each concurrent Task runs in a sibling Task Worktree at
  `<worktree.location>/<repo-slug>/<run-id>.<task_id>` on a Task Branch named
  `roundfix/run-<id>-<task_id>`. Roundfix always appends the repo slug and Run
  ID segments plus the Task suffix; those final path segments are not
  configurable.
- Spec Run startup reports the execution workspace on stderr with
  `Run Worktree: <path>`. Terminal outcomes that keep the workspace report
  `Run Worktree kept: <path>`.
- Integrated Clean spec outcomes remove the Run Worktree with
  `git worktree remove --force` and delete the Run Branch. If cleanup fails
  after integration, the Run stays Clean: stderr prints exactly one warning
  shaped as
  `roundfix: Run Worktree cleanup failed; kept <path>: <reason>`, the Daemon
  journals one Run Event, and the exit code and stdout report stay unchanged.
  The kept path remains available for manual inspection and later terminal
  Worktree reaping. Integration Pending, Unresolved, Failed, Stopped, and any
  other non-integrated spec outcome keep the Run Worktree and Run Branch.
- Worktree Bootstrap runs `worktree.bootstrap` once in each newly created spec
  Run or Task Worktree after `worktree.copy` and before Agent work and
  Verification. Empty `worktree.bootstrap` skips the step. The command runs in
  the worktree root and is bounded by `worktree.bootstrap_timeout`, which
  defaults to `10m`.
- A Worktree Bootstrap start failure, non-zero exit, or timeout fails the
  owning spec Run for a Run Worktree or settles only the owning Task failed for
  a Task Worktree. The failure reason is shaped as
  `worktree bootstrap failed: <command>: <reason>`, and bootstrap output
  streams to stderr and the Run Event Journal.
- Roundfix owns invoking and timing the Worktree Bootstrap command. Dependency
  installation, database provisioning, migrations, seeding, and cache strategy
  belong in the configured command.
- The built-in Artifact Directory default is Roundfix Home
  `artifacts/<repo-id>`. Explicit `defaults.artifact_dir` values, including
  repository-relative values, continue to override the built-in default and
  the review-artifact Spec tree resolver.

Spec Run integration uses porcelain git only. When spec Run integration cannot
fast-forward the user's branch, the Run ends Integration Pending, exits `1`,
keeps the Run Worktree and Run Branch, and prints:

```text
IntegrationPending: X completed, Y failed, Z skipped, W pending; integrate with git merge --ff-only roundfix/run-<id>
```

Completed Task Worktree commits integrate onto the Run Branch through a
serialized queue. The first compatible Task can fast-forward; later compatible
Tasks cherry-pick onto the Run Branch. A conflict settles that Task `failed`,
keeps its Task Worktree and Task Branch, and records a reason shaped like
`integration conflict: <path>`.

Review Run output and completion contract:

- With `--until-clean`, a Watch Run ends Clean only after there are no
  Unresolved Review Issues and the Review Source check on the final pushed
  commit reports success. If the Review Source check never appears within the
  grace period after Final Push, watch ends CleanUnverified, exits `3`, and
  reports the next action: confirm the pull request's Review Source check
  before merging. Pending or failing checks keep the Run inside the existing
  review timeout and Max Rounds bounds.
- `watch` and `resolve` write diagnostics, progress, the Run ID, and Agent
  output to stderr. stdout is reserved for the deterministic report at Run
  end.
- The report has one line per Review Issue in Round/fetch order, followed by
  this-Run counts and pull request cumulative counts. The CLI fixtures assert
  this shape:

  ```text
  issue 001 resolved — major: handle test issue
  This Run (Clean after 1 Round(s)): 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.
  Pull Request cumulative: 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.
  ```

  Review Issue statuses in the first line are `resolved`, `invalid`,
  `failed`, `duplicated`, or `unresolved`. Failed, invalid, and unresolved
  lines include — `reason: <terminal_reason>` when the issue artifact carries
  one. `resolve` uses the same report shape with `1 Round(s)`.
- Before a Review Source fetch completes, Review Issue knowledge is unknown.
  The report omits every zero-valued status summary and prints only:

  ```text
  Review Issues: unknown — fetch did not complete.
  ```

  A completed fetch that returns zero Review Issues is known-zero and keeps the
  two count lines. An explicit structured skip instead ends Review Skipped with
  exit `3`, fetches no Review Issues, and prints:

  ```text
  Review Source: skipped — reason: <Review Source reason>
  Next action: Reduce or split the pull request, then request another Review Source review.
  ```

  Review Skipped is not Clean, Clean Unverified, or a successful zero-issue
  Round.
- Roundfix publishes Outcome Comments on Review Source threads for
  non-resolved outcomes. Invalid and duplicated issues get the comment before
  the thread resolves. Failed issues stay open with the failed-step comment.
  Run-end unresolved issues stay open with the revisit-plan comment. Each
  comment carries an idempotency marker and each propagation is journaled with
  the Review Issue reference.
- `--no-agent-console` is available on `resolve`, `watch`, and `implement`.
  In non-TTY mode it hides Agent-source console events from stderr while
  keeping Daemon/progress lines. The Run Event Journal still records both
  Agent-source and Daemon-source events. The flag is rejected before Run
  creation when it conflicts with Interactive Input or the Live Run View.

## Review Artifact Storage

For `fetch`, `resolve`, and `watch`, Roundfix resolves the directory under
which `round-*` is written with this ADR-0029 hierarchy:

- Explicit `--artifact-dir` or `defaults.artifact_dir` preserves the legacy
  layout: `<artifact_dir>/reviews/pr-<number>/round-*`.
- Otherwise, explicit `--spec <slug>` wins. If `<specs.root>/<slug>/` exists,
  artifacts go to `<specs.root>/<slug>/reviews/round-*`.
- Otherwise, Roundfix uses the newest `Roundfix-Spec: <slug>` trailer on the
  PR head commit when that Spec folder exists.
- Without a valid Spec association, artifacts go to
  `<specs.root>/_reviews/pr-<number>/round-*`.

Unknown or invalid trailer slugs are treated as no association.

After a clean integration, `resolve` and `watch` commit the Run's review
artifacts in one separate Daemon-owned docs commit and run Final Push from the
user checkout so the commit rides it (ADR-0036). The commit subject shape is
`docs: review round NNN for pr <n>` for a single Round scope and
`docs: review rounds for pr <n>` for an all-Rounds scope, and the progress
line is `Review artifacts commit created: <subject>`. `fetch` still never
commits; `auto_commit: false` disables the review artifact commit along with
every other Daemon commit. Review artifact roots outside the repository — an
explicit external Artifact Directory, an external Spec Root, or a path
crossing a symbolic link — are never staged; the Run reports them kept
outside the repository and proceeds. Agents never create this commit by hand;
the Daemon owns it.

After Final Push, an exact Daemon-created artifact-only descendant can inherit
its verified parent Evidence without another Roundfix review request or wait.
Roundfix requires the recorded artifact commit to be the current head with
exactly the recorded verified parent and the Daemon-generated subject. Its
non-empty diff must stay entirely below the resolved in-repository review root
without a symbolic-link crossing, and refreshed parent Evidence must still be
verified with zero unresolved CodeRabbit threads. A user-authored, empty,
mixed-path, out-of-root, wrong-parent, stale, or unresolved descendant returns
to normal current-head Evidence polling.

