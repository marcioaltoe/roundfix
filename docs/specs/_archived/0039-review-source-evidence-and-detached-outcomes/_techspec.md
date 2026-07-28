---
spec: 0039-review-source-evidence-and-detached-outcomes
prd: _prd.md
created: 2026-07-17
---

# Review Source evidence and Detached Run outcomes — Technical Spec

## Executive Summary

The CodeRabbit adapter gains one shared evidence classifier used by both the pre-fetch wait and Merge-Ready confirmation, eliminating the current disagreement between `WatchStatus` and `HeadCheck`. The watch state machine gains Review Skipped, bounded transient retry, issue-count knowledge, and phase/deadline projection. Terminal context flows once through the outcome event and notification boundary, while the existing Run Event Stream remains the Supervisor subscription. A narrowly proven Daemon-created artifact-only commit may inherit evidence from its parent. The primary trade-off is CodeRabbit-specific precision: CodeRabbit signal parsing becomes richer, but the generic watch loop consumes Review Source-neutral evidence and does not embed GitHub response shapes.

## Project Constraints

- Identifier strategy: not applicable — the design carries existing Run, head,
  and Review Source-native identities without creating a new Internal
  Identifier. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the CodeRabbit adapter must reuse the
  existing authenticated `gh` boundary and preserve repository-owned HTTP and
  credential policy. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0036, ADR-0042, ADR-0043, ADR-0052,
  and ADR-0054 govern artifact commits, review checkout ownership, Clean
  Unverified, terminal completion, and accepted Review Source Evidence.
  Source: `docs/agents/domain.md`.
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

## System Architecture

- **`internal/reviewsource`** defines Review Source-neutral Evidence and typed transient failures.
- **`internal/reviewsource/coderabbit`** gathers check runs, commit statuses, CodeRabbit reviews, and unresolved threads, then classifies them through one evidence hierarchy. Both initial wait and Merge-Ready calls use this classifier.
- **`internal/watch`** consumes Evidence, manages retry episodes inside existing bounds, tracks whether a fetch completed, and returns terminal context with its Result.
- **`internal/cli`** renders Review Skipped and unknown counts, projects wait phases, proves artifact-only lineage through existing Git/commit seams, enriches Detached Run startup, and builds terminal context.
- **`internal/notify`** returns a delivery receipt and accepts additive outcome context while preserving current command environment variables.
- **`internal/runevent`** projects terminal reason, next action, issue knowledge, Console Log, Attach command, and evidence into the existing `roundfix-events/v1` outcome record. Additive optional fields preserve existing consumers.

## Implementation Design

### Interfaces

```go
type EvidenceState string
const (
    EvidencePending EvidenceState = "pending"
    EvidenceReviewing EvidenceState = "reviewing"
    EvidenceReviewed EvidenceState = "reviewed"
    EvidenceVerified EvidenceState = "verified"
    EvidenceSkipped EvidenceState = "skipped"
    EvidenceFailed EvidenceState = "failed"
)

type Evidence struct {
    State, Kind, Identity, HeadSHA, Detail, Reason string
}
```

The CodeRabbit client exposes one evidence operation with the Open Pull Request and expected head. `Kind` is a stable Review Source-neutral value such as `check_run`, `commit_status`, `review_approval`, or `artifact_only_descendant`.

```go
type TransientError struct { Operation string; Err error }
func IsTransient(error) bool

type NotificationReceipt struct {
    Route, Status string
    CompletedAt time.Time
}
```

`Notifier.Notify` returns a receipt plus error. `sent`, `skipped`, and `failed` are receipt statuses; disabled or unavailable native routes return `skipped`, not false success.

### Data Models

`watch.Result` gains:

- `ReviewIssuesKnown bool` — true only after a fetch completes;
- `TerminalReason` and `NextAction` — bounded single-line values;
- `Evidence` — the accepted or terminal Review Source Evidence;
- `VerifiedHeadSHA` — the head proven Merge-Ready before optional artifact inheritance.

The internal terminal-context value passed to event and notification boundaries includes Run identity, outcome, reason, next action, issue knowledge, Console Log, Attach command, and Evidence. It is not a new database row: the durable representation is the terminal outcome Run Event, and notification attempts are separate receipt events.

The Run Database schema version advances even though no column is required. The version fence records that `runs.state` may contain ReviewSkipped; older binaries must refuse the newer database instead of treating that unknown terminal state as Active.

CodeRabbit check-run parsing adds the output title and summary fields needed to distinguish successful review from explicit skip. Review parsing already carries state and commit SHA; classification now requires `APPROVED` instead of accepting any published review on the current head.

### API Contracts

#### Review Source evidence hierarchy

For the expected head, CodeRabbit classification applies this order:

1. An explicit skip signal in a CodeRabbit check output produces `skipped`, preserving the Review Source reason.
2. A pending or in-progress CodeRabbit check/status produces `reviewing`.
3. A successful CodeRabbit check or commit status produces `verified` when the current fetch confirms zero unresolved CodeRabbit threads; otherwise it produces `reviewed` so fetch proceeds.
4. A CodeRabbit review with state `APPROVED` on the expected head plus zero unresolved CodeRabbit threads produces `verified` with kind `review_approval`.
5. A CodeRabbit review on the expected head with another state proves `reviewed`, never verified approval.
6. Signals tied to another head are retained in detail but do not verify the expected head.
7. No usable signal produces `pending`; an explicit Review Source failure produces `failed`.

Each changed observation publishes a `daemon.review_status` event carrying state, kind, identity, expected and observed heads, conclusion, and bounded detail. Unchanged polls do not append duplicate events.

#### Review Skipped

`ReviewSkipped` is added to terminal Run states, Run Browser rendering, outcome stream projection, and Watch exit mapping. It uses the existing non-clean/unverified exit code `3`, preserving `1` for broken Runs and `2` for Preflight Validation. The final report names the Review Source reason and action and does not print resolution counts.

#### Transient retry

The CodeRabbit adapter wraps context timeouts not caused by Run cancellation, temporary DNS errors, connection resets, HTTP 429, and GitHub 5xx responses as `TransientError`. The watch loop retries after the configured poll interval while Review Source timeout and Run Budget remain. The first failure in an episode publishes `daemon.retry` with operation and reason; success publishes recovered, and bound exhaustion publishes exhausted. A graceful Stop Request from spec 0037 interrupts the retry sleep and wins over retry.

#### Wait projection

Wait events include `phase`, `expected_head`, `started_at`, `deadline`, `evidence_state`, `evidence_kind`, and retry status. `WaitingForReview` covers the pre-fetch wait; `WaitingForReviewCheck` covers Merge-Ready confirmation. The Live Run View computes remaining time locally from the deadline. Non-TTY progress prints phase entry and evidence/retry changes only.

#### Artifact-only descendant

`maybeCommitReviewArtifacts` returns the created commit SHA, parent SHA, and resolved review root. Evidence inheritance is allowed only when:

- the parent SHA already has accepted `verified` evidence;
- the current head is exactly the Daemon-created review-artifact commit;
- the diff from parent to current contains at least one path and every path is under the resolved review root;
- no unresolved CodeRabbit threads remain.

The resulting Evidence kind is `artifact_only_descendant`, with both SHAs in the event payload. Any mismatch falls back to normal current-head evidence polling. Roundfix does not publish a new review request for an inherited head.

#### Terminal report, stream, and notification

When `ReviewIssuesKnown` is false, stdout contains only:

```text
Review Issues: unknown — fetch did not complete.
```

The five zero-valued status counts are omitted. The `roundfix-events/v1` outcome record adds optional `reason`, `next_action`, `review_issues_known`, `console_log`, `attach_command`, `evidence_kind`, and `evidence_head_sha` fields.

Detached startup becomes five lines: Run ID, Console Log, Attach, Supervisor monitor, and Stop. The monitor command is `roundfix events <run-id> --follow --filter outcome`.

The existing notification variables remain unchanged. Additive variables are `ROUNDFIX_REASON`, `ROUNDFIX_CONSOLE_LOG`, `ROUNDFIX_ATTACH_COMMAND`, `ROUNDFIX_REVIEW_ISSUES_KNOWN`, and `ROUNDFIX_NEXT_ACTION`. After each attempt, a Daemon status event records `outcome_notification_sent`, `outcome_notification_skipped`, or `outcome_notification_failed`, plus route and completion time. Notification errors remain warnings and never alter the Run state or exit code.

#### Artifact persistence

A skipped review never calls the Round fetcher and creates no Round directory. A failed fetch promotes no partial Round: writes stay temporary until Review Source discovery completes, then publish atomically. A successfully fetched zero-issue Round remains valid evidence and follows the existing artifact commit policy.

## Coverage Map

- Goals 1 and 5, Stories 1, 7, and 8 → shared Evidence classifier, Review Skipped, approval evidence, artifact-only inheritance.
- Goal 2, Story 2 → typed transient failure and bounded retry episode.
- Goal 3, Stories 3, 4, and 6 → wait projection, issue knowledge, enriched outcome stream.
- Goal 4, Stories 4 and 5 → Detached monitor line, terminal context, notification receipt.
- Core Feature 12 → spec 0037 registered-session cleanup dependency.

## Integration Points

- **GitHub through `gh`**: check runs with output details, commit statuses, pull request reviews, and CodeRabbit review threads. Existing authentication and repository metadata remain unchanged.
- **Git**: local proof for the exact Daemon-created review-artifact commit and its diff root.
- **Native notification tools and `notify.command`**: existing best-effort routes with receipt metadata.
- **Supervisor**: existing `roundfix events` JSONL stream; no host-specific API is introduced.

## Testing Approach

- CodeRabbit table tests cover skipped summaries, pending checks, success/failure, approval/commented reviews on current and stale heads, unresolved threads, and evidence precedence.
- Transient classification tests cover timeout, temporary DNS, reset, 429, 5xx, permanent authentication/validation errors, and parent Run cancellation.
- Watch tests use fake clock/sleeper to prove retry bounds, deduplicated episode events, Stop Request interruption, wait deadlines, Review Skipped, and `ReviewIssuesKnown` transitions.
- Real temporary Git repositories prove exact artifact-only inheritance, mixed-path refusal, wrong-parent refusal, empty-diff refusal, and user-authored docs refusal.
- CLI tests assert Review Skipped exit `3`, unknown-count output, five-line Detached report, enriched outcome JSONL, and no skipped/failed partial Round artifact.
- Notification tests assert old and new environment variables, bounded native text, sent/skipped/failed receipts, one receipt event, and unchanged Run outcome on failure.
- Regression scenarios reproduce Fluxus's 306-file skipped review and Vortex's one-time GitHub timeout before Round 001.

## Build Order

1. Review Source-neutral Evidence and transient error contracts with CodeRabbit signal parsing and tests.
2. Shared CodeRabbit evidence hierarchy for pre-fetch and Merge-Ready paths (depends on: 1).
3. Review Skipped state, exit/report behavior, and no-artifact skip path (depends on: 2).
4. Bounded retry episodes and wait-phase projection, using spec 0037 stop-aware waits (depends on: 1, 2, spec 0037).
5. Review Issue knowledge and enriched terminal outcome stream/report (depends on: 3, 4).
6. Notification context, receipts, durable receipt events, and Detached monitor command (depends on: 5).
7. Daemon artifact-commit identity and verified evidence inheritance (depends on: 2, 5).
8. User guide, command help, CONTEXT vocabulary, ADR references, and finding
   traceability (depends on: 3, 4, 6, 7).
9. Dedicated tooling-only update of
   `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`, with
   direct byte-identical edits and read-only sync verification (depends on: 8).

## Risks & Considerations

- CodeRabbit wording can change. Skip classification must use structured check fields where available and keep unknown text as reviewed or pending rather than guessing skipped.
- HTTP status extraction from `gh` errors is imperfect. Only positively classified transient failures retry; authentication and malformed-request failures remain terminal.
- `roundfix-events/v1` gains optional fields but no removed or renamed fields. Consumers that ignore unknown fields remain compatible.
- Native notification success means the local tool accepted the request, not that a person saw it. The receipt records route completion, not human acknowledgement.
- Artifact-only inheritance prevents Roundfix from requesting or awaiting redundant evidence; it cannot stop CodeRabbit from independently reacting to a GitHub push.
- Spec 0037 must land first for stop-aware retries, registered-session cleanup, and winner-only outcome publication.

Rollout ships the state reader, writer, reporting, and schema-version fence in one binary. The migration only advances the database version before any ReviewSkipped value can be written. Because older binaries do not understand that state, binary downgrade after migration is intentionally unsupported; recovery uses a forward fix or a pre-migration Run Database backup rather than silently reopening terminal Runs.

## Decisions

- One evidence classifier feeds both watch phases.
- Review Skipped is terminal with exit code `3`.
- Existing timeout, polling, and Run Budget settings bound retry; no new config is added.
- Existing Run Event Stream is the Supervisor subscription; notification receipts remain Run Events.
- Exact Daemon artifact-only descendants inherit verified parent evidence while ADR-0036 remains in force.
- Protected Roundfix Skill publication is isolated in one Task whose changed
  files are the two authorized `SKILL.md` paths and its own Task file; it does
  not run the broad `make skills-sync` mutation target.
- See [ADR-0054](../../adr/0054-review-source-evidence-determines-review-outcomes.md).
