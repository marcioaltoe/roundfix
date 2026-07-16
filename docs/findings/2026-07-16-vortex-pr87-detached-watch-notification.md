---
status: pending
created_at: 2026-07-16
updated_at: 2026-07-16
---

# Detached watch — terminal failure notification lacked actionable context (2026-07-16)

A detached watch for Vortex PR #87 failed after a GitHub API timeout during a brief network
interruption. The user later asked whether the review was still running; only then did the
Supervisor inspect the Run Database and Console Log, discover the terminal failure, and start a
replacement Run. This report extends the shipped [Run Outcome Notifications
spec](../specs/_archived/0019-run-outcome-notifications/_prd.md) with evidence from an unattended
review Run.

Environment:

- Roundfix `0.0.0-dev` at `d3a4096`, built 2026-07-15 14:11:56 -0300.
- Repository `gesttione-solutions/vortex`, PR #87, head `df298591`.
- Failed Run `run_20260716T112753Z_5c13caf149c01780`.
- Replacement Run `run_20260716T120810Z_2baa20044c80de3b`.

## 1. Detached terminal outcome did not reach the supervising conversation

- **Symptom / evidence**: `roundfix watch ... --until-clean --detach` returned the documented
  four-line startup report and exit code `0`. The child later reached `Failed` at
  `2026-07-16T11:34:11.929176Z`, but the supervising Codex conversation received no terminal
  signal. The user had to ask `The review is running?`; the Supervisor then queried
  `~/.roundfix/roundfix.db` and tailed the Console Log to discover the failure.
- **Root cause**: the notification contract is fire-and-forget. The native desktop path targets
  the local desktop, while `notify.command` targets a user-supplied shell integration. Neither
  path acknowledges delivery to the Supervisor that launched the Detached Run. Successful
  delivery is not recorded as a Run Event, so the Run history cannot distinguish a delivered
  notification from one the user or Supervisor never received.
- **Action / suggestion**: add a durable `outcome_notification_sent` Run Event with the selected
  notification route and completion timestamp. Define a Supervisor-facing subscription or
  callback contract for Detached Runs instead of treating a desktop banner as sufficient for
  unattended agent supervision. Keep notification delivery best-effort and separate from the Run
  outcome.

## 2. The outcome notification cannot explain a failure

- **Symptom / evidence**: `notify.Outcome` carries only `RunID`, `State`, `Kind`, and `Target`.
  The native notification body is therefore only `Failed - <target>`. A custom
  `notify.command` receives the same four values through `ROUNDFIX_RUN_ID`,
  `ROUNDFIX_OUTCOME`, `ROUNDFIX_KIND`, and `ROUNDFIX_TARGET`; it receives neither the terminal
  reason nor the Console Log path.
- **Root cause**: the notification payload models Run identity and outcome, but not the evidence
  needed to act on a non-Clean outcome.
- **Action / suggestion**: extend the notification payload with a bounded terminal reason, the
  Console Log path for Detached Runs, and the exact `roundfix attach <run-id>` command. Keep the
  native body short, but make `Failed`, `Unresolved`, `TimedOut`, and `IntegrationPending`
  notifications name the next useful action. Preserve the existing environment variables when
  adding fields to `notify.command`.

## 3. One retryable GitHub timeout ended the watch after zero Rounds

- **Symptom / evidence**: the Console Log ended with:

  ```text
  roundfix: watch failed after Run start: fetch CodeRabbit commit statuses:
  gh api repos/gesttione-solutions/vortex/commits/df2985912e29661adf1e1d1d37573caa1c538a12/status:
  Get "https://api.github.com/repos/gesttione-solutions/vortex/commits/df2985912e29661adf1e1d1d37573caa1c538a12/status":
  dial tcp 4.228.31.149:443: i/o timeout: exit status 1
  ```

  The Run reached `Failed after 0 Round(s)`. Network access was healthy when the Supervisor
  checked later, and an immediate replacement watch reached the settled Review Source and fetched
  26 Review Issues.
- **Root cause**: the GitHub status request timed out. The available evidence does not prove what
  caused the network interruption. Roundfix treated the first transport timeout as terminal
  instead of retrying inside the watch boundary.
- **Action / suggestion**: classify timeouts, temporary DNS failures, connection resets, and
  retryable GitHub responses as transient Review Source failures. Retry them with a bounded policy
  inside the existing review timeout and Run Budget. Emit one deduplicated retry milestone, then
  send a terminal notification only after the retry policy is exhausted.

## 4. Cleanup noise appeared before the actionable failure

- **Symptom / evidence**: before the GitHub timeout, the Console Log printed:

  ```text
  close acpx Agent Session "roundfix-run_20260716T112753Z_5c13caf149c01780":
  No named session "roundfix-run_20260716T112753Z_5c13caf149c01780"
  ```

  No Round or Agent Batch had started. The warning occupied the first failure-shaped line even
  though the GitHub timeout was the reason the Run failed.
- **Root cause**: cleanup attempted to close the derived Agent Session name without evidence that
  the Run had created that session. This reproduces
  [Fluxus PR #53 finding 6](2026-07-15-fluxus-pr-53-watch-stop-reconciliation.md#6-force-stop-tentou-fechar-uma-agent-session-que-nunca-existiu).
- **Action / suggestion**: treat an absent Agent Session as idempotent cleanup when no Agent Batch
  started. Persist session creation and close only registered sessions. Report cleanup warnings
  after the primary failure and label them as secondary.

## 5. Zero issue counts implied knowledge the failed fetch did not have

- **Symptom / evidence**: the failed Run reported:

  ```text
  This Run (Failed after 0 Round(s)): 0 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.
  ```

  The replacement Run then fetched 26 Review Issues from the same PR head. The zero counts were
  true only for local artifacts because the first Run never completed its Review Source fetch;
  they did not mean the PR had no Review Issues.
- **Root cause**: the terminal report uses initialized local counters even when Review Source
  discovery failed before their values became meaningful.
- **Action / suggestion**: when fetch or Review Source status discovery fails, report
  `Review Issues: unknown — fetch did not complete` and omit zero-valued resolution counts. Carry
  the same distinction into notifications and machine-readable Run Events so monitors cannot
  interpret `0 unresolved` as a clean review.

## What worked — keep

- Detached startup named the Run ID, Console Log, Attach command, and Stop command.
- The Run Database and Console Log preserved enough evidence to reconstruct the failure after the
  initiating command had exited.
- Restarting the same watch after network recovery was safe and reached the existing 26 Review
  Issues without manual GitHub mutation.
- The terminal state was `Failed`, not `Clean`; the Run outcome itself did not hide the broken
  Review Source request.
