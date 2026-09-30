---
spec: 0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds
prd: _prd.md
created: 2026-09-30
---

# A reviewer that validates its findings and remembers its rounds — Technical Spec

## Executive Summary

The Pre-PR Review Command gains three stages around the reviewer call it
already makes.

1. **Before the call**, it reads the checkout's last record to place the
   review in a Reviewer Lineage: round 1, round 2 or the ceiling. The prompt
   carries the Delivery Conventions and a finding grammar. Round 2's prompt
   carries only the delta diff and round 1's findings with their dispositions.
2. **After the call**, every finding is validated. The anchor is checked
   mechanically against the candidate diff. A sealed, tool-less validator
   prompt may then dismiss a finding, but only inside a region the code
   computes.
3. **Around the call**, a round-1 session that leaves findings standing stays
   open, and round 2 continues it through acpx. At the ceiling no call is made.

The trade-off accepted is that a convention dismissal is an Agent's judgment.
It is bounded by the convention's region and fails closed. In exchange, the
queue stops parking on the delivery's own order without any rule that inspects
the meaning of a finding in Go. A second trade-off: round 2 no longer shows the
reviewer the round-1 part of the diff, and relies on round 1 having reviewed
it.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Conventions carry
  fixed ordinals `C1` to `C4` inside one versioned set,
  `roundfix/delivery-conventions/v1`. Session names keep today's
  `reviewSessionRef` form. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git, local files and the
  configured runtime through acpx; no credential and no network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0196 and ADR-0197 (this Spec) are
  implemented here. ADR-0153, ADR-0169, ADR-0174, ADR-0018, ADR-0051,
  ADR-0165, ADR-0017, ADR-0151, ADR-0189 and the gate ADRs hold as the PRD
  states. ADR-0174: "Each checkout keeps its own review record and answer
  under the Artifact Directory". The lineage reads only that record. ADR-0183
  and ADR-0184 are met by the receipts and Surface Transcripts below. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 recorded in [_authorization.md](_authorization.md); bounded
  files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/review.md`. Sanctioned regeneration:
  `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

Everything lives in the two packages that already own the review.

| Component | File | Change |
| --- | --- | --- |
| Delivery Conventions | `internal/cli/review_conventions.go` (new) | The four conventions: prompt text and region predicate |
| Finding validation | `internal/cli/review_validation.go` (new) | Anchor parser, diff line index, eligibility, record fields |
| Convention validator | `internal/cli/review_convention_validator.go` (new) | Sealed prompt, strict answer parser, fail-closed result |
| Reviewer Lineage | `internal/cli/review_lineage.go` (new) | Lineage decision, round-2 prompt, ceiling, session keep/close |
| Pre-PR Review Command | `internal/cli/review.go` | Wires the stages; record schema; reuse; `dispose` refusal |
| Delivery adapter | `internal/cli/deliver_workflow.go` | `ceiling-closed` maps to Reviewed |
| Agent runner | `internal/agent/agent.go`, `internal/agent/acpx_runner.go`, `internal/agent/acp_stream.go` | `ExecuteResult.ACPSessionID` from the stream |

The Delivery Queue engine (`internal/delivery`) does not change. It already
parks on `Findings` and `Blocked` and advances on `Reviewed`.

```text
runReviewCommand
  └─ decideLineage(prior record, candidate) ─┬─ ceiling → closeAtCeiling (no runner call)
                                             └─ round 1 | round 2
       buildReviewPrompt (+ conventions, grammar) / buildRoundTwoPrompt (+ delta, round-1 findings)
       runConfiguredReviewSession (continued session on round 2)
       classifyReviewCommandResult (unchanged grammar)
       validateReviewFindings (anchor → eligibility → sealed validator)
       settle outcome, keep or close session, persist
```

## Implementation Design

### Interfaces

```go
// internal/cli/review_conventions.go
type deliveryConvention struct {
	ID     string // "C1".."C4"
	Prompt string // one sentence, rendered as "<ID>. <Prompt>"
}
const deliveryConventionsVersion = "roundfix/delivery-conventions/v1"
func deliveryConventions() []deliveryConvention
func conventionRegions(ctx context.Context, repo reviewRepository, anchor reviewFindingAnchor) ([]string, error) // IDs whose region holds the whole anchor

// internal/cli/review_validation.go
func parseReviewFindingAnchor(text string) (reviewFindingAnchor, bool)
func indexReviewDiff(diff string) reviewDiffIndex
func (index reviewDiffIndex) holds(anchor reviewFindingAnchor) bool
func reviewFindingHasFailureClause(text string) bool
func validateReviewFindings(ctx context.Context, in reviewValidationInput) ([]reviewFinding, reviewValidation)

// internal/cli/review_convention_validator.go
type reviewSealedRunner interface {
	RunSealedPrompt(context.Context, agent.SealedPromptRequest) (agent.SealedPromptResult, error)
}
func runConventionValidator(ctx context.Context, runner reviewSealedRunner, runtime agent.RuntimeSpec, workDir string, asked []validatorQuestion) (map[string]validatorVerdict, error)

// internal/cli/review_lineage.go
func decideReviewLineage(ctx context.Context, prior *reviewRecord, candidate reviewRecord, git preflight.GitRunner) (reviewLineagePlan, error)
func buildRoundTwoPrompt(plan reviewLineagePlan, delta string, previous []reviewFinding, dispositions []reviewFindingDisposition) string
func closeAtCeiling(plan reviewLineagePlan, ledger []reviewFindingDisposition, git preflight.GitRunner) (reviewRecord, int)

// internal/agent/agent.go
type ExecuteResult struct { /* existing fields */ ACPSessionID string }
```

### Data Models

The record keeps every existing field, in order. It gains the following:

```go
type reviewFinding struct {
	ID         string                   `json:"id"`
	Text       string                   `json:"text"`
	Anchor     *reviewFindingAnchor     `json:"anchor,omitempty"`
	Validation *reviewFindingValidation `json:"validation,omitempty"`
}
type reviewFindingAnchor struct {
	Path      string `json:"path"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
}
type reviewFindingValidation struct {
	Status string `json:"status"`         // "stands" | "dismissed-by-validation"
	Rule   string `json:"rule,omitempty"` // "unanchored" | "no-failure" | "convention:C1".."convention:C4"
	Reason string `json:"reason"`
}
```

At the record's top level, after `specContextTruncated`, it gains two fields:

- `validation`: `{conventions, validator, reason}`, where `validator` is
  `not-needed`, `ran` or `unavailable`, and `reason` is present only when
  unavailable.
- `lineage`: `{round, previousHead, previousFindings, previousDispositions,
  reviewedHead, session, selection, sessionOpen, acpSessionIds, continued}`.
  `session`, `selection` and `acpSessionIds` are present only when the review
  prepared a session; `sessionOpen` and `continued` are always present.
  task_01 and task_02 write no `lineage`. task_03 adds it without the session
  fields, and task_04 adds those.

Outcomes gain `ceiling-closed`. Validation rules:

- `findings` needs at least one finding whose validation stands and which has
  no evidence-backed `dismissed` disposition.
- `findings-dismissed` needs every finding to be either
  `dismissed-by-validation` or matched by exactly one evidence-backed
  `dismissed` disposition.
- `ceiling-closed` needs `lineage.round` 2, a `lineage.reviewedHead`
  different from `headCommit`, and one disposition for each standing
  finding: `dismissed` with evidence, or `fixed` with a `fixedBy`.

A record without `validation` or `lineage` is read as before, and its findings
count as standing. It then belongs to a lineage at round 1.

### Invariants

1. A finding is `dismissed-by-validation` only by one of three rules:
   - `unanchored`, when its anchor is missing or the candidate diff holds no
     line of it;
   - `no-failure`, when it has no `Failure:` clause and the validator says so;
   - `convention:Cn`, when the whole anchor lies in `Cn`'s region and the
     validator says so.
2. A validator error, timeout, tool use or answer outside its grammar leaves
   every finding it was asked about standing, and the record says
   `unavailable` with the reason.
3. An answer with at least one finding and no readable anchor at all is
   `blocked`.
4. A review at the ceiling never calls the runner, never writes the record or
   the answer, and never removes the answer file.
5. A round-2 prompt carries the diff from the previous head to the current
   head, and never the diff from the merge base.
6. The anchor check reads the full merge-base diff at the current head, in
   every round.
7. A review leaves a session open only when it is round 1 and its outcome is
   `findings`. Every other review ends each session it prepared, as does a
   lineage change that finds an open session.
8. `continued` is true only when both rounds report the same non-empty ACP
   session id.
9. No stage reads another checkout's record, and nothing writes the
   disposition ledger except `roundfix review dispose`.

### The finding grammar and the prompt

`buildReviewPrompt` keeps its signature and every current line. Two things
change:

- After the line that begins `If there are findings, respond with Findings:`,
  it adds: ``Start each finding with its anchor, `path:line` or
  `path:start-end`, naming a line of the candidate diff, and state what breaks
  in a clause that starts with `Failure:`.``
- After the `Head commit` line, it adds a block. The block starts with
  `Delivery Conventions (roundfix/delivery-conventions/v1). A Roundfix
  delivery writes the following by design; do not report them as defects:`,
  followed by one line `Cn. <sentence>` per convention.

The four sentences are:

- **C1.** "A Spec's QA Report records the head it audited and is committed
  after that head, so it never names the commit that records it."
- **C2.** "The Daemon writes a Task file's status and its `## Result`,
  `## Recorded paths` and `## Carry-forward provenance` sections after the
  Task's Verification passes; a Result that calls status Daemon-owned agrees
  with a `completed` status."
- **C3.** "The archive commit moves a completed Spec's directory to the
  archive root and stamps its archive front matter."
- **C4.** "A planning candidate authors a Spec whose Tasks are all pending and
  which has no QA Report; that Spec's own delivery implements it and is
  reviewed then."

An anchor is read from the start of a finding item. The item may open with a
backtick. The anchor is a path of characters other than whitespace, `:` and a
backtick, then `:`, a line, and optionally `-` and an end line. An optional
backtick follows. The next character must be the end of the text, whitespace,
`:`, `,` or `)`. Only the first range counts. A failure clause is the
case-insensitive text `Failure:` followed by at least one non-blank character.

### Regions

The regions are computed from the head's tree with `git show <head>:<path>`,
and from the Spec roots `reviewCandidateSpecRoots` already resolves. Here
`<slug>` is a directory directly under a root.

| Convention | Region |
| --- | --- |
| C1 | Any path under `<root>/<slug>/qa/`, for the Spec Root or the archive root |
| C2 | In `<root>/<slug>/task_*.md`: the front matter (line 1 through its closing `---`), and each `## Result`, `## Recorded paths` or `## Carry-forward provenance` section, from its heading to the line before the next `## ` heading |
| C3 | Any path under the archive root |
| C4 | `<spec-root>/<slug>/_tasks.md`, and the front matter of each `task_*.md` of that slug, when every `task_*.md` there reads `status: pending` and nothing exists under `<slug>/qa/` |

A region holds an anchor only when every line from start to end is inside it.

### Validation

`indexReviewDiff` walks the unified diff that `reviewCandidateDiff` returns,
file by file.

- A hunk `@@ -a,b +c,d @@` holds the lines `c` through `c+d-1`, or the line
  `c` alone when `d` is `0`.
- A deleted file (`+++ /dev/null`) is keyed by its `---` path and holds every
  line.
- A file the diff names with no hunk (binary, mode-only or pure rename) holds
  every line.

For each finding, in order:

1. No readable anchor, or an anchor the index does not hold: the finding is
   dismissed with rule `unanchored`. The reason is `anchor names no line of
   the candidate diff` or `finding has no path:line anchor`.
2. Otherwise the eligible rules are the `convention:Cn` rules whose region
   holds the anchor, plus `no-failure` when the finding has no failure clause.
   A finding with no eligible rule stands, with reason `anchored in the
   candidate diff`.
3. When at least one finding is eligible, one sealed validator prompt asks
   about all eligible findings. Otherwise `validation.validator` is
   `not-needed`.

The validator runs through `reviewSealedRunner`, which the production runner
satisfies, on the selection that ran the review. It uses the existing sealed
limits: 2 MiB of input, 512 KiB of output and 5 minutes. Its input holds:

- the convention list;
- for each eligible finding, its ID, anchor, text and eligible rules;
- up to 60 lines of the anchored file at the head (at the merge base for a
  deleted file), from five lines before the anchor to five lines after it.

The input is labelled untrusted data. The answer must be exactly one JSON
object:
`{"schema":"roundfix/review-validation/v1","findings":[{"id","verdict","rule","reason"}]}`.
`verdict` is `stands` or `dismiss`. Each asked ID must appear exactly once.
`rule` must be one of that finding's eligible rules when `verdict` is
`dismiss`. `reason` must not be blank. Any violation, a runner that does not
implement `reviewSealedRunner`, an error or `ToolUsed` makes the validator
`unavailable`, and every eligible finding stands. The reason is then
`validator unavailable: <cause>`.

The outcome is decided after validation:

| Condition | Outcome | Exit |
| --- | --- | --- |
| A standing finding remains | `findings` | `1` |
| Every finding is dismissed, by validation or by the operator's evidence | `findings-dismissed` | `0` |
| No finding carries a readable anchor | `blocked`, reason `findings name no file and line` | `2` |

For each dismissed finding, standard error prints `roundfix: review finding
<ID> dismissed by validation (<rule>): <reason>`. Reuse at the same head counts
`dismissed-by-validation` as dismissed, and `missing` lists only standing
findings. `roundfix review dispose` refuses a finding dismissed by validation
with `finding "<ID>" was dismissed by validation (<rule>)`.

### The Reviewer Lineage

`decideReviewLineage` compares the checkout's prior record P with the
candidate. Its answer decides the round:

| P | Plan |
| --- | --- |
| absent, another provider or merge base, or a head that is neither the candidate head nor its ancestor | new lineage, round 1; end P's open session first |
| same head | today's reuse path. A blocked P runs its own round again |
| an ancestor head, P `blocked` | run P's `lineage.round` again with P's previous data (round 1 for a legacy record) |
| an ancestor head, P round 1 (or legacy), outcome `reviewed`, `findings` or `findings-dismissed` | round 2, `previousHead` = P's head |
| an ancestor head, P round 2 or `ceiling-closed` | ceiling |

**Round 2.** The delta diff is `git diff --no-ext-diff --no-textconv
--no-color <previousHead> <head> --`. `buildRoundTwoPrompt` renders the
following:

- the round-1 prompt's header, grammar and conventions;
- the line `Previous head: <previousHead>`;
- the sentence `This is round 2 of 2 of this review. The diff below is the
  change since round 1 reviewed <previousHead>.`;
- a `Round 1 findings:` list, with one line per previous finding. Each line
  gives its ID, its validation (`stands`, or `dismissed-by-validation
  (<rule>): <reason>`), its disposition (`fixed by <sha>`, `dismissed:
  <evidence>` or `no disposition`) and its text;
- the instruction `Do not raise again a round-1 finding that was dismissed or
  fixed unless the change below reintroduces it; raise a standing round-1
  finding without a disposition only if the change below leaves it
  unresolved.`;
- the delta between `--- BEGIN ROUND 2 DELTA DIFF ---` and
  `--- END ROUND 2 DELTA DIFF ---`;
- the Spec contexts, as today.

The previous dispositions are the ledger entries for P's repository and head.
The record copies P's findings and those dispositions into
`lineage.previousFindings` and `lineage.previousDispositions`.

**The ceiling.** `closeAtCeiling` takes the lineage's round-2 findings: P's
`findingItems`, whose reviewed head is P's head for round 2, or P's
`lineage.reviewedHead` for `ceiling-closed`. For each standing finding it looks
up a ledger disposition for that head. Either of two dispositions counts:

- `dismissed` with evidence;
- `fixed`, with a `fixedBy` for which `git merge-base --is-ancestor <fixedBy>
  <head>` succeeds.

When every standing finding has one, the review writes a `ceiling-closed`
record at the candidate head and exits `0`. The record carries the round-2
items, the matched dispositions, `lineage.round` 2 and `lineage.reviewedHead`.
Otherwise it prints, and does not persist, a `blocked` record. Its reason
reads `review round ceiling reached: round 2 reviewed <reviewedHead>; dispose
finding <IDs> with roundfix review dispose`. The command exits `2`.
`deliveryReviewResult` maps `ceiling-closed` to `delivery.ReviewOutcomeReviewed`.

### The continued session

- **Round 1** prepares a session named by `reviewSessionRef`, as today.
  - When its outcome is `findings`, it does not call `EndSession`. The record
    stores `session`, `selection` (the index into preferred and fallbacks) and
    `sessionOpen: true`.
  - Every other outcome ends the session.
- **Round 2**, when P's `sessionOpen` is true, first prepares P's `session` on
  selection `selection`.
  - If preparing it fails, it ends that name and runs today's selection loop
    with fresh names, and `continued` is false.
  - It always ends its session.
- **Other closings.** A lineage change ends P's open session. So does a reuse
  at the same head that resolves to `findings-dismissed`. Ending an
  already-closed session is not an error.
- **The session id.** `ExecuteResult.ACPSessionID` is the `sessionId` of the
  first `session/update` notification the acpx runner reads. The ACP schema
  requires that field on every notification. Each round appends its id, or
  an empty string, to `lineage.acpSessionIds`. `continued` compares the two.

### API Contracts

1. API Contract: `roundfix review` prompt — carries the finding-grammar
   sentence and the four Delivery Conventions in every round. Round 2 carries
   the delta diff and the round-1 findings instead of the full diff.
2. API Contract: `roundfix review` record — adds `anchor` and `validation` to
   each finding item, `validation` and `lineage` to the record, and the
   outcome `ceiling-closed`. Exit codes: `0` for `reviewed`,
   `findings-dismissed`, `omitted` and `ceiling-closed`; `1` for `findings`;
   `2` for `blocked` and preflight failures, including the ceiling's blocked
   result.
3. API Contract: `roundfix review` standard error — prints one line per
   finding dismissed by validation. A blocked result keeps
   `roundfix: review blocked: <reason>`.
4. API Contract: `roundfix review dispose <id>` — refuses a finding dismissed
   by validation, with exit `2`. Everything else is unchanged.
5. API Contract: Delivery Queue review step — advances on `ceiling-closed` and
   parks `review-blocked` on the ceiling's blocked result. Park reasons do not
   change.

### Surface Transcripts

Each transcript runs in a checkout of `0001-example` with provider `codex`.
Placeholders in angle brackets stand for commit identifiers, paths and
reasons.

1. Surface Transcript: a Daemon settlement restated as a finding, dismissed as
   `C2`.

   ```transcript
   $ roundfix review
   stdout:
   {"repository":"<repository>","baseCommit":"<merge base>","baseTipCommit":"<base tip>","headCommit":"<head>","provider":"codex","source":"project","outcome":"findings-dismissed","findings":"- `docs/specs/0001-example/task_01.md:4` — The Task is `completed` while its Result says status is Daemon-owned. Failure: an invalid Task state transition.","findingItems":[{"id":"F1","text":"`docs/specs/0001-example/task_01.md:4` — The Task is `completed` while its Result says status is Daemon-owned. Failure: an invalid Task state transition.","anchor":{"path":"docs/specs/0001-example/task_01.md","startLine":4,"endLine":4},"validation":{"status":"dismissed-by-validation","rule":"convention:C2","reason":"<validator reason>"}}],"answerPath":"<answer path>","specs":[],"skippedSpecs":[],"archivedSpecs":[],"specContextTruncated":false,"validation":{"conventions":"roundfix/delivery-conventions/v1","validator":"ran"},"lineage":{"round":1,"session":"<session>","selection":0,"sessionOpen":false,"acpSessionIds":["<acp session id>"],"continued":false}}
   stderr:
   roundfix: review finding F1 dismissed by validation (convention:C2): <validator reason>
   exit: 0
   ```

2. Surface Transcript: one finding outside the diff and one inside it.

   ```transcript
   $ roundfix review
   stdout:
   {"repository":"<repository>","baseCommit":"<merge base>","baseTipCommit":"<base tip>","headCommit":"<head>","provider":"codex","source":"project","outcome":"findings","findings":"<both items>","findingItems":[{"id":"F1","text":"<text of F1>","anchor":{"path":"internal/untouched.go","startLine":10,"endLine":10},"validation":{"status":"dismissed-by-validation","rule":"unanchored","reason":"anchor names no line of the candidate diff"}},{"id":"F2","text":"<text of F2>","anchor":{"path":"review.txt","startLine":2,"endLine":2},"validation":{"status":"stands","reason":"anchored in the candidate diff"}}],"answerPath":"<answer path>","specs":[],"skippedSpecs":[],"archivedSpecs":[],"specContextTruncated":false,"validation":{"conventions":"roundfix/delivery-conventions/v1","validator":"not-needed"},"lineage":{"round":1,"session":"<session>","selection":0,"sessionOpen":true,"acpSessionIds":["<acp session id>"],"continued":false}}
   stderr:
   roundfix: review finding F1 dismissed by validation (unanchored): anchor names no line of the candidate diff
   exit: 1
   ```

3. Surface Transcript: findings that carry no anchor at all.

   ```transcript
   $ roundfix review
   stdout:
   {"repository":"<repository>","baseCommit":"<merge base>","baseTipCommit":"<base tip>","headCommit":"<head>","provider":"codex","source":"project","outcome":"blocked","reason":"findings name no file and line","answerPath":"<answer path>","specs":[],"skippedSpecs":[],"archivedSpecs":[],"specContextTruncated":false,"validation":{"conventions":"roundfix/delivery-conventions/v1","validator":"not-needed"},"lineage":{"round":1,"session":"<session>","selection":0,"sessionOpen":false,"acpSessionIds":["<acp session id>"],"continued":false}}
   stderr:
   roundfix: review blocked: findings name no file and line
   exit: 2
   ```

4. Surface Transcript: round 2 after a fix, on a continued session.

   ```transcript
   $ roundfix review
   stdout:
   {"repository":"<repository>","baseCommit":"<merge base>","baseTipCommit":"<base tip>","headCommit":"<head 2>","provider":"codex","source":"project","outcome":"reviewed","answerPath":"<answer path>","specs":[],"skippedSpecs":[],"archivedSpecs":[],"specContextTruncated":false,"validation":{"conventions":"roundfix/delivery-conventions/v1","validator":"not-needed"},"lineage":{"round":2,"previousHead":"<head 1>","previousFindings":[<round-1 items>],"previousDispositions":[<ledger entries for head 1>],"session":"<session>","selection":0,"sessionOpen":false,"acpSessionIds":["<acp session id>","<acp session id>"],"continued":true}}
   stderr:
   exit: 0
   ```

5. Surface Transcript: a third review with a round-2 finding that has no
   disposition.

   ```transcript
   $ roundfix review
   stdout:
   {"repository":"<repository>","baseCommit":"<merge base>","baseTipCommit":"<base tip>","headCommit":"<head 3>","provider":"codex","source":"project","outcome":"blocked","reason":"review round ceiling reached: round 2 reviewed <head 2>; dispose finding F1 with roundfix review dispose","specs":[],"skippedSpecs":[],"archivedSpecs":[],"specContextTruncated":false,"lineage":{"round":2,"reviewedHead":"<head 2>","sessionOpen":false,"continued":false}}
   stderr:
   roundfix: review blocked: review round ceiling reached: round 2 reviewed <head 2>; dispose finding F1 with roundfix review dispose
   exit: 2
   ```

6. Surface Transcript: the same third review after `roundfix review dispose F1
   --fixed-by <fix>`.

   ```transcript
   $ roundfix review
   stdout:
   {"repository":"<repository>","baseCommit":"<merge base>","baseTipCommit":"<base tip>","headCommit":"<head 3>","provider":"codex","source":"project","outcome":"ceiling-closed","findings":"<round-2 findings>","findingItems":[<round-2 items>],"dispositions":[<fixed disposition of F1>],"specs":[],"skippedSpecs":[],"archivedSpecs":[],"specContextTruncated":false,"lineage":{"round":2,"reviewedHead":"<head 2>","sessionOpen":false,"continued":false}}
   stderr:
   exit: 0
   ```

7. Surface Transcript: disposing a finding that validation dismissed.

   ```transcript
   $ roundfix review dispose F1 --dismiss --evidence "restates the Daemon settlement"
   stdout:
   stderr:
   roundfix: review dispose refused: finding "F1" was dismissed by validation (convention:C2)
   exit: 2
   ```

## Coverage Map

- Goal 1 → Delivery Conventions; Regions; Validation; API Contract 1.
- Goal 2 → Validation (anchor, eligibility); Invariants 1 and 3.
- Goal 3 → Data Models; Validation (standard error); API Contracts 2 and 3.
- Goal 4 → The Reviewer Lineage (round 2); The continued session.
- Goal 5 → The Reviewer Lineage (the ceiling); API Contract 5.
- User Story 1 → Delivery Conventions; Regions.
- User Story 2 → Validation.
- User Story 3 → Data Models; API Contract 4.
- User Story 4 → The Reviewer Lineage (round 2); The continued session.
- User Story 5 → The Reviewer Lineage (the ceiling).
- Core Features 1–2 → The finding grammar and the prompt; Regions.
- Core Feature 3 → Validation (anchor); Invariant 3.
- Core Feature 4 → Validation (validator); Invariant 2.
- Core Feature 5 → Data Models; Validation (outcome table); API Contract 4.
- Core Feature 6 → The Reviewer Lineage.
- Core Feature 7 → The continued session; Invariants 7 and 8.
- Core Feature 8 → The Reviewer Lineage (the ceiling); Invariant 4.
- Core Feature 9 → Build Order 4.
- Success Metrics 1–3 → Testing Approach 1 and 2.
- Success Metrics 4–6 → Testing Approach 3 and 4.
- Success Metric 7 → Testing Approach 6.

## Integration Points

- **acpx and the ACP adapters.** A continued session is an acpx named session
  that was never closed. acpx 0.19.3 reattaches an open record by name and
  working directory. When it must respawn the adapter, it tries
  `session/resume`, then `session/load`, then a new session. The validator
  uses the existing sealed path, which opens and closes its own session and
  denies tools.
- **The disposition ledger.** Read to render round-2 dispositions and to close
  at the ceiling. Written only by `roundfix review dispose`, unchanged.
- **The Delivery Queue.** Reads the record from standard output
  (`deliveryReviewResult`). Only the `ceiling-closed` mapping is new.
- **Spec 0194's reference files.** Documentation goes to the review reference
  and command guide that Spec 0194 creates.

## Testing Approach

Every test uses temporary repositories, temporary Artifact Directories, fake
Agent runners that also implement `SessionPreparer`, `PreparedPromptRunner`,
`EndSession` and `reviewSealedRunner`, and fake ACP streams. None reaches a
reviewer, an adapter, `~/.roundfix`, `~/.acpx` or the network. The existing
fixture `newReviewCommandFixture` has a diff that holds `review.txt` lines 1–2.

1. **Anchor and grammar** (`review_validation_test.go`):
   - the parser reads backticked, bare, ranged and comma-listed anchors, and
     rejects a missing line or a trailing letter;
   - the index covers context lines, a deletion, a deleted file and a
     binary file;
   - Surface Transcripts 2 and 3;
   - the prompt carries the grammar sentence and the four convention lines,
     and today's grammar sentence is unchanged.
2. **Validator** (`review_convention_validator_test.go`):
   - Surface Transcript 1;
   - a `C2` verdict for a finding anchored on Requirements is refused because
     no region holds it, and the finding stands;
   - `no-failure` for a finding with a failure clause is never asked;
   - a validator error, a timeout, `ToolUsed`, a missing ID, a duplicate ID,
     an ineligible rule and a blank reason each leave every asked finding
     standing, with `unavailable`;
   - `C1`, `C3` and `C4` regions each admit and refuse at their edges;
   - Surface Transcript 7.
3. **Lineage** (`review_lineage_test.go`):
   - the lineage table row by row, including a rebase and a blocked P;
   - the round-2 prompt holds the delta's line and not the round-1 diff's
     line, and lists each round-1 finding with its disposition;
   - Surface Transcripts 5 and 6, the runner never called, and the round-2
     record's bytes unchanged after the blocked ceiling;
   - `deliveryReviewResult` maps `ceiling-closed` to Reviewed.
4. **Continued session** (`review_session_test.go`, `agent_session_id_test.go`):
   - no `EndSession` after a round 1 with findings;
   - round 2 prepares the recorded name and ends it, and Surface Transcript 4;
   - a failed preparation falls back to a fresh name with `continued: false`;
   - a lineage change ends the open session;
   - through the fake ACP stream, `ExecuteResult.ACPSessionID` is the first
     notification's `sessionId`.
5. **Characterization.** Existing tests keep their names. Tests whose fake
   answers cite `internal/cli/review.go` or another line outside their fixture
   diff are changed to cite a line the diff holds. The test for the grammar
   sentence stays unchanged and green.
6. **Outside evidence.** The QA gate replays the disposition ledger read-only
   against the repository's history, as the PRD's Acceptance evidence
   describes.
7. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. Delivery Conventions in the prompt, the finding grammar, anchor
   validation, the outcome rules and the `dispose` refusal (depends on: none).
2. The convention validator and regions (depends on: 1 — both edit
   `internal/cli/review.go` and `internal/cli/review_validation.go`).
3. The Reviewer Lineage: round 2, the ceiling and the delivery mapping
   (depends on: 2 — both edit `internal/cli/review.go`).
4. The continued session, the ACP session id, and the skill, mirror, command
   guide and skill version (depends on: 3 — both edit
   `internal/cli/review.go` and `internal/cli/review_lineage.go`).
5. Terminal QA (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **A validator that dismisses a real finding.** It can do so only inside a
  convention region or for a finding without a failure clause, and it gives a
  recorded reason. Replayed over the ledger, no `fixed` finding lies in a
  region.
- **A reviewer that stops anchoring.** The grammar guard blocks instead of
  passing.
- **Round 2 misses a defect in round-1 code.** Accepted. Round 1 reviewed it,
  and a standing round-1 finding is listed for the reviewer to re-raise.
- **Open sessions.** Only after a round 1 with findings. They are closed on
  round 2, on a lineage change and on a dismissed reuse. A removed checkout's
  session stays open, which is a recorded limit.
- **Ceiling dismissal needs the round-2 head.** `dispose --dismiss` still
  requires HEAD to be the reviewed head. After a fix, the operator records
  `--fixed-by`.
- **Older readers.** A Roundfix older than this Spec cannot read a
  `ceiling-closed` record and reviews again. That costs a round and loses
  nothing.

## Decisions

- Validation before the verdict decides; dismissal bounded by a mechanical
  region; fail closed. See ADR-0196.
- A Reviewer Lineage of two rounds, round 2 on the delta with the recorded
  conversation, a continued session measured by ACP session id, and a ceiling
  closed only by dispositions. See ADR-0197.
- Keep `reviewSessionRef` random, and store the name in the record rather
  than deriving it. The lineage is the record, and a derived name would need a
  second identity.
- A blocked ceiling result is printed and not persisted, so `dispose` keeps
  reading the round-2 record.
- The anchor set is the full merge-base diff in every round, so a standing
  round-1 finding can be raised again at its own line.
