---
spec: 0179-review-findings-with-evidence-and-no-unselected-providers
status: active
created: 2026-09-28
surfaces: [backend, cli, docs]
---

# Review findings with evidence, and no unselected providers

## Executive Summary

Scope every CodeRabbit request to the repository that selected it or to the
legacy PR-feedback commands. Give each pre-PR review finding an ordinal
identity and one evidence-backed disposition in an append-only ledger. Let a
findings verdict stand for its head until every finding is dismissed there,
with a fresh review for a changed candidate. Park a blocking review of a
candidate that archives a Spec as `corrective-spec-required`.

## Project Constraints

- Identifier strategy: applicable — a finding's identity is `F<n>` within the
  review record of one head; a disposition is keyed by repository, head
  commit, finding identity and finding text. No other identifier changes.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and the Run Database
  only; no credential, no network call and no provider request is added.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0014, ADR-0057, ADR-0080, ADR-0091,
  ADR-0093, ADR-0096, ADR-0097, ADR-0104, ADR-0110, ADR-0117, ADR-0120,
  ADR-0130, ADR-0142, ADR-0151, ADR-0153, ADR-0154, ADR-0155 and ADR-0156
  hold; ADR-0020, ADR-0038, ADR-0056, ADR-0127, ADR-0159 and ADR-0160 do not
  apply. This Spec adds ADR-0165. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Unselected providers receive no request

The one Roundfix code path that publishes a CodeRabbit request is the legacy
review-request publication of `watch` and `resolve` (`requester.RequestReview`
in `internal/cli/cli.go` and `internal/watch/watch.go`), gated by an explicit
`--source coderabbit` run and `review_source.request_review`. Preflight's
`validateReviewRequestCoherence` in `internal/preflight/preflight.go` reads
`.coderabbit.yaml` for `watch` and `resolve` only. Neither changes.
`runReviewCommand` in `internal/cli/review.go` refuses `coderabbit` before any
call, and `commandDeliveryWorkflow.Publication` in
`internal/cli/deliver_workflow.go` writes a title and body with no request.
Tests pin both, so a later change that adds a request or a `.coderabbit.yaml`
read to either path fails.

The universal assumption lives in guidance. The Roundfix skill's section "A
review only happens when it is asked for" tells every agent to request
CodeRabbit review on every pull request "across these repositories". That
section is scoped to a repository whose Pre-PR Review Policy selects
`coderabbit` and to the legacy commands `fetch`, `watch` and `resolve`. It says
that a repository whose policy selects `codex`, `claude` or `none` is never
asked for a CodeRabbit review. The sentence "The supported Review Source is
`coderabbit`" in `docs/user-guide/commands.md` and `docs/user-guide/usage.md`,
and the `review_source.name` row of `docs/user-guide/configuration.md`, add
that the Review Source is the legacy PR-feedback source, read only by `fetch`,
`watch` and `resolve`, and never selects or requests a pre-PR reviewer.
`defaultConfigYAML` in `internal/config/config.go` writes the same statement as
a comment above `review_source:`, so a generated config says it too. The
generated config still writes no `pre_pr_review` key, so it opts the
repository into neither a service nor `none`.

## Finding identity and dispositions

`classifyReviewCommandResult` sets `findingItems` on a findings record from
`splitReviewFindings(findings)`. A line with no leading whitespace that starts
with `- `, `* `, or ASCII digits followed by `.` or `)` and a space opens a new
finding; any other line continues the current one. Text before the first
marker is one finding when it is not blank, and text with no marker at all is
one finding. Markers and surrounding whitespace are trimmed, and blank findings
are dropped. Identities are `F1`, `F2`, … in order. A record read from disk
without `findingItems` derives them from `findings` the same way.
`validateReviewRecord` keeps accepting a findings record without items.
`buildReviewPrompt` asks for each finding as a list item that starts with `- `
and names its file and line.

`roundfix review dispose <finding-id>` is parsed when the first argument of
`review` is `dispose`. It reads `pre-pr-review.json` from the Artifact
Directory and requires a record whose repository is the current Git root and
whose outcome is `findings` or `findings-dismissed`. Exactly one form is
accepted:

- `--dismiss --evidence <text>` requires non-blank evidence and a current
  `HEAD` equal to the record's `headCommit`.
- `--fixed-by <commit>` resolves the commit, and requires it to differ from the
  record's head, to have that head as an ancestor, and to be reachable from the
  current `HEAD`.

On success it appends one JSON line to `pre-pr-review-dispositions.jsonl`
beside the record and prints the same line on stdout. The line carries
`repository`, `headCommit`, `finding`, `text` (copied from the record),
`disposition` (`dismissed` or `fixed`), `evidence` or `fixedBy`, and
`recordedAt` in RFC 3339 UTC. Each refusal exits `2` with
`roundfix: review dispose refused: <reason>` on stderr and appends nothing:

- no record, a record of another repository, or an outcome other than
  `findings` or `findings-dismissed`;
- an unknown finding identity;
- both forms, or neither;
- blank evidence;
- a dismissal when `HEAD` moved;
- a `--fixed-by` commit that fails any of its three checks;
- a finding that already has a disposition for the same repository, head,
  identity and text.

The top-level usage in `internal/cli/cli.go` and `commandUsage("review")` list
both forms. The review sections of `docs/user-guide/commands.md` and the
Roundfix skill describe items, the two forms, the ledger and the refusals.
`CONTEXT.md` gains a Review Finding Disposition entry.

## The verdict stands for its head

`runReviewCommand` decides reuse after it resolves the base commit and the
Artifact Directory and confirms the policy is `codex` or `claude`. It decides
before `removeReviewAnswer`, profile resolution and readiness. When
`pre-pr-review.json` names the same repository, base commit, head commit and
provider with outcome `findings` or `findings-dismissed`, the reviewer is not
asked again. The command joins the ledger entries with the same repository,
head, identity and text into `dispositions`, and sets `reused: true`:

- When every finding has a `dismissed` entry, the outcome is
  `findings-dismissed` and the exit is `0`.
- Otherwise the outcome is `findings` and the exit is `1`, with one stderr line
  naming each finding without a disposition.

A `fixed` entry never counts toward clearing. A record whose head, base or
provider differs, or whose outcome is `reviewed`, `blocked` or `omitted`, is
never reused. `validateReviewRecord` accepts `findings-dismissed` only with
findings text, no reason, at least one item and one `dismissed` disposition per
item.

`commandDeliveryWorkflow.Review` maps records through a new
`deliveryReviewResult(record, head)`, which turns `findings-dismissed` into
`delivery.ReviewOutcomeReviewed` for its head. A retried `review-findings`
item at an unchanged head therefore advances once its findings are dismissed
and parks again while any stands. Nothing re-asks the reviewer until the head
moves. The review and delivery sections of `docs/user-guide/commands.md` and
the Roundfix skill describe the reuse, the outcome and its exit code.

## A blocking review after archive

`reviewCandidateSpecContexts` also returns the slugs of changed Spec folders
under any root other than the active Spec Root, which
`reviewCandidateSpecRoots` returns first. That is the resolved archive root,
and each slug counts whether its context is carried or skipped. The record
carries them sorted in `archivedSpecs`, which is always present and is `[]` on
paths that compute no Spec context. When a findings record names archived
Specs, `finishReviewCommand` writes one stderr line naming them. The line says
that an archived Spec is never corrected in place and tells the operator to
author a corrective Spec with its own authorization and QA gate.

`delivery.ReviewResult` gains `ArchivedSpecs`, filled by
`deliveryReviewResult`. `reviewCandidate` in `internal/delivery/engine.go`
parks a findings result with archived Specs as
`corrective-spec-required: <slug>, <slug>` (new constant
`BlockerCorrectiveSpecRequired`) and keeps `review-findings` for findings with
none. `Engine.Retry` handles a blocker with that prefix before its other
branches. It resolves the item branch and head, then refuses with exit `2`
when the head differs from the parked candidate head, naming the archived
Specs and the corrective Spec requirement and leaving the item unchanged.
When the head is unchanged it returns the item to `reviewing` without
carry-forward, so the recorded dispositions decide. The delivery sections of
`docs/user-guide/commands.md` and the Roundfix skill add the blocker and the
retry rule, and `CONTEXT.md` gains a Corrective Spec entry.

The policy these rules enforce, confirmed by the maintainer on 2026-09-28, is
recorded as
`docs/adr/0165-a-blocking-review-after-archive-parks-publication-for-a-corrective-spec.md`
with `status: accepted`:

- An archived Spec is never edited to absorb a finding.
- Publication parks as `corrective-spec-required: <slug>`.
- The correction is a new corrective Spec with its own `_authorization.md` and
  QA gate.
- No Run budget, corrective-Task ceiling or queue grant authorizes it, and
  Roundfix never authors or starts it.
- How the corrective work reaches the parked candidate is a recorded limit for
  a future Spec.

## API Contracts

1. `pre-pr-review.json` gains `findingItems` (`[{id, text}]`) on findings
   records, `dispositions`, `reused`, and `archivedSpecs` (always present), plus
   the outcome `findings-dismissed`. Every existing field keeps its name and
   meaning.
2. `pre-pr-review-dispositions.jsonl` holds one JSON object per line with
   `repository`, `headCommit`, `finding`, `text`, `disposition`, `evidence` or
   `fixedBy`, and `recordedAt`, and is only ever appended.
3. `roundfix review dispose <finding-id> (--dismiss --evidence <text> |
   --fixed-by <commit>)` exits `0` and prints the appended entry, or exits `2`
   and appends nothing.
4. `roundfix review` exits `0` for `findings-dismissed` and `1` for standing
   findings, and asks no reviewer when it reuses a record.
5. The Delivery Queue blocker `corrective-spec-required: <slug>[, <slug>]`
   parks a findings verdict on a candidate that archives a Spec; `roundfix
   deliver retry` refuses such an item with exit `2` when its head moved and
   returns it to `reviewing` when unchanged.
6. The generated User Config and Project Config label `review_source` as the
   legacy PR-feedback Review Source read only by `fetch`, `watch` and
   `resolve`, and write no `pre_pr_review` key.

## Coverage Map

- Goal 1 → Unselected providers receive no request; API Contract 6.
- Goal 2 → Finding identity and dispositions; API Contracts 1-3.
- Goal 3 → The verdict stands for its head; API Contracts 1 and 4.
- Goal 4 → A blocking review after archive; API Contracts 1 and 5.
- Core Feature 1 → Unselected providers receive no request.
- Core Feature 2 → Finding identity and dispositions.
- Core Feature 3 → The verdict stands for its head.
- Core Feature 4 → A blocking review after archive.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 4.
- Success Metric 5 → Testing Approach 5.
- API Contracts 1-3 → Finding identity and dispositions.
- API Contract 4 → The verdict stands for its head.
- API Contract 5 → A blocking review after archive.
- API Contract 6 → Unselected providers receive no request.

## Integration Points

- **Spec 0175.** Active and ahead in the queue; it edits
  `internal/delivery/engine.go`, `internal/cli/deliver_workflow.go`,
  `docs/user-guide/commands.md`, the Roundfix skill and `CONTEXT.md`, so each
  Task here revalidates those files after it merges.
- **Spec 0173.** Added `roundfix deliver retry`, whose re-entry table this Spec
  extends with one row.
- **Specs 0153 and 0160.** Shipped the review record, the Codex and Claude
  review path and the Delivery Queue's review stage that this Spec builds on.

## Testing Approach

1. **Unselected providers.** `roundfix review` under `codex` and `none` in a
   repository whose `.coderabbit.yaml` and `review_source.request_review` form
   the pair Preflight refuses for `watch`; the delivery publication under every
   policy; the generated User Config and Project Config; and
   `TestRunEnforcesReviewRequestCoherence` unchanged as the mirror.
2. **Dispositions.** Items from marked and unmarked answers; one test per
   dismissal and fix; one test per refusal asserting exit `2` and
   byte-identical ledger.
3. **Head-bound verdict.** A reused record with every finding dismissed, with
   one standing, and with a dismissal of different text, each asserting zero
   reviewer calls. A moved head, a different base and a different provider each
   asserting one call. A `reviewed` and a `blocked` record each asserting one
   call. The delivery mapping of both findings outcomes.
4. **After archive.** `archivedSpecs` for a candidate that archives a Spec and
   `[]` for an active one, with the stderr line only on findings. The engine
   parks `corrective-spec-required` with archived Specs and `review-findings`
   without. The retry refuses a moved head and returns an unchanged head to
   `reviewing`.
5. **Outside evidence.** The terminal QA Task replays the dispositions of
   #257, #258 and #259 through the built binary.
6. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. Unselected providers receive no request (depends on: none).
2. Finding identity and dispositions (depends on: 1).
3. The verdict stands for its head (depends on: 2).
4. A blocking review after archive (depends on: 3).
5. Terminal QA (depends on: 1, 2, 3, 4, 6, 7).
7. The revision option guard (depends on: 6).
6. The atomic disposition reservation (depends on: 4).

## Risks & Considerations

- **One file, four Tasks.** Every Task edits `docs/user-guide/commands.md` and
  the Roundfix skill, and Tasks 2, 3 and 4 edit `internal/cli/review.go`, so
  the graph is a chain; each Task puts its new tests in files of its own.
- **Review shopping.** Without the reuse rule, a retry re-asks the reviewer at
  an unchanged head until the answer changes, the permissive direction ADR-0153
  forbids. The rule is covered by tests that count reviewer calls, not by the
  record alone.
- **A stale skill.** Each guidance change ships in the Task that changes the
  behavior it describes.
- **Parallel ADR ordinals.** Specs authored in parallel claim the ordinals
  after 0163 below this one; this Spec claims ADR-0165, and the ordinal check
  reports a collision once the Specs meet on `main`.
