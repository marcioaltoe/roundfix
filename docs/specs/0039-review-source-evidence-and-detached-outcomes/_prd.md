---
spec: 0039-review-source-evidence-and-detached-outcomes
status: active
created: 2026-07-17
surfaces: [backend, cli, docs]
---

# Review Source evidence and Detached Run outcomes

Roundfix currently treats several different CodeRabbit states as the same signal: a completed check can mean a real review or an explicit skip, an approval can be ignored when a check is missing, and one transient GitHub failure can end a watch before any Round. When failure happens in a Detached Run, zero issue counts and a context-free notification leave users and Supervisors without the evidence needed to recover. Prior dogfood evidence was absorbed into this Spec and remains in Git history; the still-open behavior is documented by the [Vortex detached-watch finding](../../findings/2026-07-16-vortex-pr87-detached-watch-notification.md).

## Project Constraints

- Identifier strategy: not applicable — Review Source Evidence reuses Run IDs,
  Git heads, and Review Source-native identities and creates no project-owned
  Internal Identifier. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — GitHub and CodeRabbit access must
  continue through the repository's existing `gh` and Review Source boundaries;
  this feature adds no authentication provider, credential policy, or HTTP
  route. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0036 preserves the separate
  review-artifact commit, ADR-0042 keeps review Runs in the user checkout,
  ADR-0043 preserves Clean Unverified, ADR-0052 protects terminal completion,
  and ADR-0054 makes head-bound Review Source Evidence authoritative. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-26, the maintainer expressly
  authorizes changes to exactly `.agents/skills/roundfix/SKILL.md` and
  `skills/roundfix/SKILL.md`. On 2026-07-28, the maintainer additionally
  expressly authorizes the deterministic Skill-digest fallout of that edit in
  exactly `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`. No other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- Review Source outcomes reflect what CodeRabbit actually did on the relevant head, including an explicit skipped review.
- Temporary Review Source failures retry within existing Run bounds instead of failing the first poll.
- Waiting users and Supervisors can see the expected head, evidence state, deadline, retry episode, and terminal reason.
- Detached Run notifications and outcome events contain enough context to inspect and recover without database archaeology.
- A Roundfix-owned artifact-only descendant does not force Roundfix to request or await redundant Review Source evidence.

## User Stories

1. As a user whose review was skipped, I want a distinct Review Skipped outcome and actionable reason, so that zero downloaded issues cannot be mistaken for Merge-Ready.
2. As a user on an unstable network, I want transient GitHub failures retried within the existing timeout and Run Budget, so that one brief outage does not end a watch before its first Round.
3. As a user waiting for CodeRabbit, I want the current phase, expected head, deadline, and latest evidence visible, so that an Active Run never looks abandoned.
4. As a Supervisor monitoring a Detached Run, I want a first-party outcome subscription command and a structured terminal record, so that I can react without polling the Run Database or parsing the Console Log.
5. As a user receiving a failure notification, I want the terminal reason, Console Log, Attach command, and next action, so that the notification is operational rather than informational.
6. As an automation consumer, I want unknown Review Issue counts represented as unknown, so that a failed fetch never looks like zero unresolved work.
7. As a user whose current head has an explicit CodeRabbit approval and no unresolved CodeRabbit threads, I want that evidence accepted as Merge-Ready when a check is absent.
8. As a user preserving Roundfix review artifacts in a separate commit, I want Roundfix to recognize its own artifact-only descendant and retain the verified code-head evidence, so that it does not wait for a redundant re-review.

## Core Features

1. CodeRabbit status is represented as typed Review Source Evidence tied to a head. Evidence distinguishes pending, reviewing, reviewed, verified, skipped, and failed states and records the signal type, identity, conclusion, summary, and observed head.
2. An explicit CodeRabbit skip ends the Run as Review Skipped with a non-zero exit code, bounded reason, and actions such as reducing or splitting the pull request. Roundfix does not fetch or persist an empty Round for a skipped review.
3. Merge-Ready accepts either an accepted current-head CodeRabbit check/status or an `APPROVED` current-head CodeRabbit review when no unresolved CodeRabbit threads remain. Every Clean outcome records which evidence proved it.
4. A current head that is exactly a Daemon-created review-artifact-only descendant of an already verified code head inherits that evidence. Roundfix neither requests nor waits for another review of that head. It does not claim to suppress Review Source behavior triggered independently by GitHub.
5. Timeouts, temporary DNS failures, connection resets, and retryable GitHub responses retry using the existing poll interval, Review Source timeout, and Run Budget. One deduplicated retry episode is journaled when retrying begins and another when it recovers or exhausts.
6. Review Source waits expose `WaitingForReview` or `WaitingForReviewCheck`, the expected head, start and deadline, latest evidence, and retry state. Events publish on phase or evidence changes, not every unchanged poll.
7. If status discovery or fetch fails before a Round completes, the terminal report prints `Review Issues: unknown — fetch did not complete` and omits zero-valued issue counts. The outcome event and notification carry the same known/unknown distinction.
8. Every terminal outcome event carries a bounded terminal reason when non-Clean, the next action, and the Console Log and Attach command when available. The stable Run Event Stream projects these fields for Supervisors.
9. Detached startup prints the existing human Attach command plus the exact `roundfix events <run-id> --follow --filter outcome` Supervisor monitor command.
10. A successful, skipped, or failed outcome notification attempt appends a durable Run Event with its route, status, and completion time. Notification delivery remains best-effort and never changes the Run outcome.
11. `notify.command` preserves its existing four environment variables and adds terminal reason, Console Log, Attach command, Review Issue knowledge, and next action.
12. Cleanup-before-Agent behavior follows spec 0037: only registered Agent Sessions are closed, and secondary cleanup warnings follow the primary Review Source failure.

## User Experience

- Review Skipped reports the skipped signal and its Review Source reason, then names the next action. It never prints a zero-issue completion summary.
- During waits, non-TTY output prints one line on phase entry and on evidence or retry changes. The Live Run View shows the same state and derives remaining time from the recorded deadline.
- A failed pre-fetch Run prints the terminal reason followed by `Review Issues: unknown — fetch did not complete`.
- Detached startup includes Run ID, Console Log, human Attach, Supervisor monitor, and Stop commands.
- Native notifications stay short. Failure and unresolved notifications name the next action; command notifications receive the complete structured environment.

## Non-Goals / Out of Scope

- Guaranteed delivery into a specific Codex, Claude, or other host conversation. Roundfix exposes the durable stream and callback context; the host owns consumption.
- Built-in Slack, email, Telegram, or other service integrations.
- Suppressing a Review Source webhook or automatic review that CodeRabbit independently triggers after a push.
- Unbounded retries or a new retry configuration surface.
- Treating a generic pull request approval from another reviewer as CodeRabbit evidence.
- Accepting artifact-only inheritance for user-authored documentation commits or paths outside the resolved Roundfix review-artifact root.
- Changing force-stop, terminal compare-and-set, or Agent Session registry behavior owned by spec 0037.

## Success Metrics

- Explicit skipped-check fixtures produce Review Skipped in every status and merge-readiness path and create zero Round artifacts.
- Each transient outage episode produces at most one retry-start event and one recovery or exhaustion event while remaining bounded by the existing timeout and Run Budget.
- A failed status or fetch fixture emits zero misleading issue-count summaries and reports Review Issues as unknown in text, notification context, and the outcome stream.
- A current-head CodeRabbit approval with zero unresolved CodeRabbit threads produces Clean and records `review_approval` evidence.
- A proven Daemon artifact-only descendant reuses verified evidence without another Roundfix review request or wait; a mixed-path descendant does not.
- Every Detached Run startup report includes a runnable Supervisor monitor command, and every terminal stream record includes a next action for non-Clean outcomes.
- Every notification attempt produces exactly one durable receipt event without changing the terminal outcome on failure.

## Decisions

- Review Skipped is a dedicated terminal outcome, not Clean Unverified or Failed.
- The existing Run Event Stream is the Supervisor subscription contract; no second callback protocol is added.
- The separate review-artifact commit remains, and only a proven Daemon-created artifact-only descendant inherits verified evidence.
- Retry reuses existing polling and budget bounds rather than adding configuration.
- See [ADR-0054](../../adr/0054-review-source-evidence-determines-review-outcomes.md).

## Open Questions

None.
