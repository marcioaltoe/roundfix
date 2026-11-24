---
spec: 0419-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed
prd: _prd.md
created: 2026-11-23
---

# Round 2 of a pre-PR review opens a new session when the configured model changed — Technical Spec

## Executive Summary

Round 2 of a Reviewer Lineage continues the Agent Session round 1 left open by preparing its recorded name on the recorded index of the review profile. Nothing compares the model, because the lineage does not record one. When the profile's model changed in between, Codex answers on the continued session with a notice before the verdict, and the Pre-PR Review Command blocks a no-findings verdict beside other content (`_prd.md`, reported by fiscus). The fix is three changed lines in code that exists: the lineage gains the model of the session each round used, every round that sends a prompt records it, and round 2 prepares the recorded name only when the selection at the recorded index still names that model. In every other case round 2 takes the path it already has for a session that cannot be prepared: it ends the recorded session and opens a new one with the round-2 prompt.

The trade-off accepted is the reviewer's continued conversation. A round 2 whose model cannot be proven the same starts in a new session, and that includes every lineage that is open when Roundfix is updated, whose record holds no model. The design accepts that loss because the round-2 prompt already carries what the reviewer needs (ADR-0197: "so it never depends on the reviewer's memory"), and because the alternative, taking a runtime's notice out of the answer before classifying it, would make Roundfix track the wording of each runtime's notices and loosen a rule an accepted decision keeps strict (ADR-0174: "Both verdicts together, or a no-findings verdict beside other content, still block."). The measurements behind the design are in [`references/sources.md`](references/sources.md), sections 5 and 6.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted. The design records an Agent Model name the review profile already holds, and the decision record takes the number 0451. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added. The Pre-PR Review Command reaches the reviewer through the Agent runner as it does today. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the PRD's row stands. ADR-0197 is the decision this design amends in part, by the maintainer's decision the PRD quotes, and ADR-0451 is the record a Task writes for it. ADR-0174, ADR-0256, ADR-0349 and ADR-0233 are the ones it obeys or uses. ADR-0196, ADR-0250, ADR-0272, ADR-0286, ADR-0366, ADR-0391 and ADR-0412 are listed only to record why they do not apply. This stage adds no ADR to the PRD's list, and `## ADR Alternatives` relates each one. Source: `docs/agents/domain.md`, `docs/agents/docs-layout.md`.
- Tooling authority: applicable — express maintainer authorization: the maintainer's three decisions of 2026-11-23, quoted with their provenance in the authorization record, which grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release` at a `patch` Release Ceiling. The design changes no Makefile, workflow, Go module, lint or formatter file. The governed set was measured with `speccheck.GovernedPath` over every file the Build Order names, in a `go test -overlay` probe that writes nothing to the repository: the three granted paths are the only governed ones. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`; Spec-contained authorization record: `docs/specs/0419-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/review.md`, `skills/roundfix/SKILL.md`.

## System Architecture

Every component extends something that exists. The new files are two test files and one ADR.

### The lineage of the Pre-PR Review Record

`reviewLineage` in `internal/cli/review_lineage.go` holds the session name, the selection index, `sessionOpen`, the ACP session identifiers and `continued`. It gains the model, after the selection index.

### Where a round records its session

In `internal/cli/review.go`, the block that runs when a prompt was sent (lines 933 to 943 at `50143e74`) writes the session name, the selection index and the ACP session identifier into the lineage. It also writes the model of the runtime the prompt went to. A record written without a prompt, a `ceiling-closed` record for one, gains none, and a reuse at the same head keeps the record it reuses.

### The round-2 branch of the selection loop

The same file prepares the recorded session name when the plan is round 2, is not a Correction Review, and the prior lineage has `sessionOpen` (lines 1757 to 1787). Its inner condition accepts any recorded index inside the profile. It now also requires the recorded model to be present and equal to the model of the selection at that index. When the condition fails, control reaches the two steps that follow the branch today and do not change: `endOpenReviewSession` ends the recorded session, and the selection loop starts at the preferred selection with a new random name. The prompt is the one the plan already built for round 2.

### The decision record, the glossary, the guide and the skill

The Task writes ADR-0451 with the notice in ADR-0197 and its line of the Standard Export List, revises the glossary entry Reviewer Lineage, rewrites two paragraphs of the review page of the user guide and of the canonical review reference of the Roundfix skill, and regenerates the mirror, the version fields and the Behavior Surface Record.

## Implementation Design

### Interfaces

```go
// internal/cli/review_lineage.go
type reviewLineage struct {
	// ... the fields before Selection, unchanged and in their order ...
	Selection *int `json:"selection,omitempty"`
	// Model is the Agent Model of the Agent Session this round's prompt was
	// sent to. It is empty when the round sent no prompt and in a record an
	// earlier Roundfix wrote.
	Model       string `json:"model,omitempty"`
	SessionOpen bool   `json:"sessionOpen"`
	// ... the fields after SessionOpen, unchanged and in their order ...
}
```

```go
// internal/cli/review.go, the round-2 branch. Today:
//   if selection >= 0 && selection < len(selections) { prepare prior.Session ... }
// After this Spec the same block runs only when, in addition:
//   prior.Model != "" && prior.Model == strings.TrimSpace(selections[selection].Model)

// internal/cli/review.go, the block guarded by promptSent, beside Session and Selection:
//   record.Lineage.Model = specContext.runtime.Model
```

### Invariants

```text
1. Round 2 prepares the recorded session name only when the prior lineage holds
   a model and the Agent Selection at the recorded index names the same model,
   compared as exact strings after trimming.
2. When round 2 does not continue, the only call it makes on the recorded
   session is its end, and that end precedes the preparation of any new session.
3. A round 2 on a new session sends the round-2 delta prompt, stays round 2 and
   records continued false; the two-round ceiling and the Correction Review
   read its record as they read any round-2 record.
4. A round whose prompt was sent records the model of the session that received
   it; a record written without a prompt gains none.
5. The classifier of the reviewer's answer keeps its code and its results.
6. The help of the review commands keeps its bytes, and no output line, flag,
   exit code or key of Project Config is added.
7. A test finds ADR-0451 by its title, never by its number.
8. A record that carries the model stays readable by a Roundfix released
   before this Spec, and a record without it stays readable after.
```

### Data Models

The Pre-PR Review Record gains one optional string, `lineage.model`, written after `lineage.selection`. No other field moves, and a lineage without a model encodes the bytes it encodes today. There is no migration: a record without the field is read as a record whose model is unknown.

```json
{"round":2,"previousHead":"<commit>","session":"roundfix-pre-pr-review-<head>-<random>","selection":0,"model":"gpt-6.1-sol","sessionOpen":false,"acpSessionIds":["<id>","<id>"],"continued":false}
```

### API Contracts

1. API Contract: the session of round 2. Input: a plan at round 2 that is not a Correction Review, whose prior lineage has `sessionOpen: true`. Round 2 prepares the recorded session name on the selection at the recorded index when that index is inside the review profile and invariant 1 holds; everything after that preparation is as today, `continued` included. Otherwise it prepares nothing on the recorded name, ends it, and runs the selection loop from the preferred selection with new random names and the round-2 prompt. Three states take the second path: the model at the recorded index changed, the record holds no model, and the recorded index is outside the profile, which already does so today. Failure modes are today's: a failed end of the recorded session returns its error before any prompt; a failed preparation of a new session retries once and falls back as the loop does. Standard output and standard error gain no line.
2. API Contract: the record field. `lineage.model` holds the model of the Agent Session the round's prompt was sent to: the preferred selection's model, or the model of the fallback the round ended on. A round 2 that continued the round-1 session records the same model as round 1; a round 2 on a new session records the new one; a Correction Review records the model of its own new session. The field is absent from a record that sent no prompt. It appears wherever the record does: on the standard output of `roundfix review`, in the checkout's record file and in the Reviewed Head Record.
3. API Contract: the decision record. `docs/adr/0451-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed.md` is accepted, carries a true Knowledge Mark and has the title `Round 2 of a pre-PR review opens a new session when the configured model changed`. It states API Contracts 1 and 2 as rules, names the report of the Fiscus of 2026-11-23 as its context, records the five options this design gives up (`## Decisions`, the first five lines), and says that it amends ADR-0197 in part. ADR-0197 gains `ADR-0451` at the end of its `amended_by` list, takes ADR-0451's `created_at` as its `updated_at`, and gains one paragraph after the notice of ADR-0391 that begins `**Amended by [ADR-0451](<file>) (<date>).**`, with the date of the new record's `created_at`; every other line of that file keeps its bytes. `.secondbrain-export` gains the ADR's path once, in order, inside the Standard Export List. The list entry, the notice and the export line are written by hand.
4. API Contract: the glossary. In `CONTEXT.md`, the first sentence of the entry Reviewer Lineage ends, word for word, "continuing the same Agent Session when the runtime can resume it and the review profile still names the Agent Model that session was opened with, and otherwise opening a new one with the same prompt." The entry's list of decision records gains a link to ADR-0451 after the link to ADR-0391, written as its neighbours are. No other entry changes.
5. API Contract: the guide and the skill. In `docs/user-guide/commands/review.md` and in `.agents/skills/roundfix/references/review.md`, the paragraph that begins "When round 1 leaves standing findings" becomes the first paragraph below, and the first sentence of the paragraph after it becomes the second line below; both files wrap them as they wrap their neighbours, and the rest of each file keeps its bytes.

   ```text
   When round 1 leaves standing findings, its reviewer session stays open. The checkout record stores the random session name, the selection index into the review profile's preferred selection and fallbacks, the model of that selection as `lineage.model`, and `sessionOpen: true`. Round 2 prepares that name on that selection through acpx only when the selection at that index still names the recorded model. When the model changed, or when the record holds no model because an earlier Roundfix wrote it, Roundfix does not prepare that name: it ends it and tries the normal selection loop with fresh random names. It does the same when preparation fails. The round-2 prompt carries the recorded findings even when resume is unavailable.

   Every round records the ACP session ID reported by its prompt stream in `lineage.acpSessionIds`, or an empty string when absent, and the model of its session in `lineage.model`.
   ```

   `make skills-sync` then writes the mirror `skills/roundfix/references/review.md`, the skill record command raises both version fields to the next version and records it, and the Behavior Surface record command re-records the fingerprint of the guide page. The commands and their outputs are the ones `_authorization.md` declares under Sanctioned regeneration.

### Surface Transcripts

None. No line of standard output or standard error and no exit code is added or reworded. What changes is which outcome a round 2 reaches after a changed model, and one optional field of the record, which API Contract 2 fixes. A transcript is compared through the built product, and the built Pre-PR Review Command reaches a reviewer only through a real Agent session, which this Spec's authorization keeps out of its tests and its QA.

## Coverage Map

- Goal 1 → The round-2 branch of the selection loop (API Contract 1, invariants 1 to 3)
- Goal 2 → The lineage of the Pre-PR Review Record, Where a round records its session (API Contract 2, invariant 4)
- Goal 3 → The round-2 branch of the selection loop, whose first path is today's (API Contract 1, invariant 1)
- Goal 4 → The decision record, the glossary, the guide and the skill (API Contracts 3 to 5)
- Core Feature 1 → `reviewLineage.Model` and the block guarded by `promptSent` (API Contract 2)
- Core Feature 2 → the condition of the round-2 branch (API Contract 1, invariant 1)
- Core Feature 3 → the condition's second path, `endOpenReviewSession` and the selection loop, both unchanged (API Contract 1, invariants 2 and 3)
- Core Feature 4 → invariants 5 and 6
- Core Feature 5 → API Contracts 3, 4 and 5, invariant 7
- Success Metric 1 → `TestReviewRoundTwoOpensANewSessionWhenTheConfiguredModelChanged`
- Success Metric 2 → `TestReviewRoundTwoOpensANewSessionWhenTheRecordHoldsNoModel`
- Success Metric 3 → `TestReviewRoundTwoContinuesTheRecordedSession` and `TestReviewRoundTwoContinuesTheRecordedFallbackSelection`, both existing and unchanged
- Success Metric 4 → `TestReviewRecordsTheModelOfTheSessionOfEachRound`, with the existing `ceiling-closed` test of `internal/cli/review_lineage_selection_test.go` unchanged
- Success Metric 5 → `TestTheReviewModelDecisionAmendsTheLineageRule`, `TestTheReviewerLineageEntryStatesTheModelRule`, `TestTheReviewGuideAndReferenceDescribeTheModelRule`, with the existing documentation, skills and glossary contracts
- Success Metric 6 → `TestReviewRoundTwoOpensANewSessionWhenTheConfiguredModelChanged`, whose fake runner returns the answer of `references/sources.md`, section 1, and its subtest `the_classifier_still_refuses_the_reported_answer`

Deviations from the PRD: None.

## Absence Evidence

1. A model in the lineage or anywhere in the two changed source files — command: `rg -n 'Model' internal/cli/review.go internal/cli/review_lineage.go`; result: no match.
2. A comparison before round 2 prepares the recorded session — command: `rg -n 'prior\.Selection|selection < len\(selections\)' internal/cli/review.go`; result: three matches, at lines 1760, 1761 and 1763, which read the recorded index and test only that it is inside the profile.
3. The field in the guide, the skill or the glossary — command: `rg -n 'lineage\.model' internal docs/user-guide .agents/skills/roundfix skills/roundfix CONTEXT.md`; result: no match, and `rg -n 'Round 2 prepares that name on that selection' docs/user-guide .agents/skills/roundfix skills/roundfix` finds the sentence to replace once in each of the three files.
4. The decision record — command: `ls docs/adr | rg '^0451'`; result: no match, and `rg -l -i 'model changed' docs/adr CONTEXT.md docs/user-guide .agents/skills/roundfix` finds no file.
5. The six tests and the fake runner — command: `rg -n 'TestReviewRoundTwoOpensANewSession|TestReviewRecordsTheModel|TestTheReviewModelDecision|TestTheReviewerLineageEntry|TestTheReviewGuideAndReference|modelNoticeReviewRunner' internal skills`; result: no match, and `ls internal/cli/review_model_change_test.go internal/docscontract/review_model_decision_test.go` finds neither file.

## ADR Alternatives

- ADR-0451: the record a Task of this Spec writes, absent until then; it rejects nothing yet. API Contract 3 fixes what it states and the five options it records as given up.
- ADR-0197: amends it in part, by the maintainer's decision of 2026-11-23 that the PRD's Core Feature 5 names. Its rule for the session gains a condition (ADR-0197: "round 2 continues it through acpx, on the same selection"), and the path the design takes for a changed model is the one it already gives a session that cannot continue (ADR-0197: "A runtime that cannot resume gets a new session, and round 2 still has the recorded conversation."). The lineage is not started again, so the ceiling it decided stands, and the behavior it replaced, a new session that rereads the whole candidate at every review, does not return.
- ADR-0174: obeys it. The verdict is still read from the last message and the classifier is not touched (ADR-0174: "The classifier's strictness is unchanged inside the final message."). The design removes the cause of the extra content instead of the content.
- ADR-0256: obeys it. The review page of the user guide is a Behavior Surface whose fingerprint moves, and the covering skill reference changes in the same range (ADR-0256: "is a Lagging Surface unless a covering skill file changed in the same range"). The record is regenerated by its command, never judged by a model, which it rejected.
- ADR-0349: uses it. The authorization record carries a `patch` Release Ceiling and the Spec is a `fix`. No Task performs the release.
- ADR-0233: uses its first two rules for the skill's version, which Spec 0418 raises too: the record command chooses the next version on the tree it runs on, and the two version fields are line-scoped, so a conflict at the Pull Request takes the default branch's lines and regenerates. The option it rejected, refreshing each item from the default branch before its Run, is not relied on. Its third rule is not touched.
- ADR-0196: not applicable — the validation of a finding and the moment a finding parks are not changed, and the reported verdict holds no finding to validate.
- ADR-0250: not applicable — no Baseline module and no coverage record is changed; the version this Spec raises is an owned skill's.
- ADR-0272: not applicable — the directory of a checkout's record and the Reviewed Head Record are not changed; the head record carries the new field only because it holds the same content as the current record.
- ADR-0286: not applicable — a Correction Review never continues a session (ADR-0286: "it opens a fresh Agent Session"), so the new condition does not reach it; it records the model of its own session as any round that sends a prompt does.
- ADR-0366: not applicable — the delta of round 2 after a merge of the default branch, and the measurement of a Review Segment, are not changed.
- ADR-0391: not applicable — the Conflicted Carry and the delta that starts from it are not changed.
- ADR-0412: not applicable — the release plan's reading of a commit is not changed.

## Integration Points

- The Agent runner: `PrepareSession`, `RunPrepared` and `EndSession`, called as today. Round 2 skips one `PrepareSession` it used to make, on the recorded name, when the model is not proven the same.
- The Pre-PR Review Record: one optional field more, read and written by the same functions. The Delivery Queue reads the record from the command's output and takes its stages from the outcome; its code reads no field of the lineage.
- The Delivery Queue: unchanged. An item whose review used to park as `review-blocked` after a changed model now reaches the stage its verdict supports. The `Review:` line of `roundfix deliver status` and `roundfix deliver retry` is printed only for a lineage that starts again, which a changed model does not cause.
- The Secondbrain: reads `.secondbrain-export` after the merge and receives ADR-0451.
- The Fiscus: the operator tells it after the release, as the maintainer's first decision asks; no Task sends a message.

## Testing Approach

No new seam. The tests use what the review session tests already inject: `newReviewCommandFixture`, `runLineageReview`, `commitLineageFile`, `disposeLineageFinding` and the fake `sessionReviewRunner`, over a temporary repository, a temporary home and a temporary Artifact Directory.

New, in `internal/cli/review_model_change_test.go`:

- A fake runner, `modelNoticeReviewRunner`, that embeds `sessionReviewRunner`, remembers the model each session name was first prepared with, and puts the notice of `references/sources.md`, section 1, with the two model names, before the answer of a session name prepared again with another model. The source of the probe that proved it, and that proved the three tests below, is in that file, section 6.
- `TestReviewRoundTwoOpensANewSessionWhenTheConfiguredModelChanged`. The review profile's preferred selection names `gpt-6-luna`; round 1 answers with one finding; a commit fixes it and the disposition is recorded; the profile changes to `gpt-6.1-sol`; round 2 runs. It asserts exit `0`, the outcome `reviewed`, `round: 2`, `continued: false`, `sessionOpen: false`, empty standard error, `lineage.model` of `gpt-6-luna` in round 1 and `gpt-6.1-sol` in round 2, exactly two preparations with different session names and those two models, the order of the runner's events (prepare the round-1 name, end it, prepare the new name, end it), and a round-2 prompt that holds the round-2 delta marker and the fixing commit and not the full candidate marker. On the unchanged tree it fails: the record of round 1 holds no model, and round 2 exits `2` with the blocked reason of the PRD's Success Metric 1. Its subtest `the_classifier_still_refuses_the_reported_answer` gives the reported answer to a round 1 as the reviewer's only message and asserts exit `2` and the outcome `blocked` with that reason, which keeps invariant 5 honest.
- `TestReviewRoundTwoOpensANewSessionWhenTheRecordHoldsNoModel`. After a round 1 with a finding, the test removes `lineage.model` from the checkout's record file, as a record of an earlier Roundfix reads, and runs round 2 with the profile unchanged. It asserts a new session name, the end of the round-1 name before it, `continued: false` and the model in the round-2 record. On the unchanged tree it fails because round 2 prepares the round-1 name.
- `TestReviewRecordsTheModelOfTheSessionOfEachRound`. A round 1 on the preferred selection and a round 1 that ended on a fallback each record that selection's model, in the command's output and in the record file; a round 2 that continued records the same model as round 1. On the unchanged tree it fails because the record holds no model.

New, in `internal/docscontract/review_model_decision_test.go`, with the package's build tag and marker:

- `TestTheReviewModelDecisionAmendsTheLineageRule`: API Contract 3 by invariant 7, as `TestTheIndexusDecisionAmendsTheConsultationRule` does for its ADR, including the export line and the bytes of ADR-0197 before the new notice.
- `TestTheReviewerLineageEntryStatesTheModelRule`: the clause of API Contract 4 in the entry Reviewer Lineage of `CONTEXT.md`, and its link to the record the first test found.
- `TestTheReviewGuideAndReferenceDescribeTheModelRule`: in the guide page, the canonical reference and its mirror, compared with whitespace collapsed, the four phrases "the model of that selection as `lineage.model`", "only when the selection at that index still names the recorded model", "or when the record holds no model because an earlier Roundfix wrote it" and "and the model of its session in `lineage.model`", and the absence of "Round 2 prepares that name on that selection through acpx. If preparation".

Existing tests that keep passing unchanged, and are the controls: `TestReviewRoundTwoContinuesTheRecordedSession`, `TestReviewRoundTwoContinuesTheRecordedFallbackSelection`, `TestReviewRoundTwoFallsBackToAFreshSessionWhenPreparationFails`, `TestReviewLineageChangeEndsTheOpenSession`, and the test of `internal/cli/review_lineage_selection_test.go` that pins the bytes of a `ceiling-closed` lineage. With the prototype laid over the tree, every test of `internal/cli` and of `internal/delivery` passed, and a probe of the three new tests of `internal/cli` passed over the prototype and failed on the unchanged tree (`references/sources.md`, section 6). No existing test, golden or fixture of those two packages moves.

The QA gate runs the six new tests and the controls by name, with three sabotages: remove the comparison of the model (the first new test must fail); stop recording the model (`TestReviewRoundTwoContinuesTheRecordedSession` must fail, because round 2 can no longer prove the model); make the condition always false (the same control must fail). No test reads the network, the real home, another repository or the clock, and none starts an Agent session.

## Build Order

1. The session rule, its record and its documents, in one change that answers the finding reported by fiscus on 2026-11-23 (ai-memory message `01a20a55-66b6-7b31-bd03-ed4f3d6a1fbc`). In this order:
   - the lineage field, the block that records it and the condition of the round-2 branch (API Contracts 1 and 2), with the three tests of `internal/cli/review_model_change_test.go`;
   - ADR-0451, the notice and the list entry in ADR-0197, and the export line (API Contract 3);
   - the glossary entry Reviewer Lineage (API Contract 4);
   - the guide page and the canonical review reference (API Contract 5), then `make skills-sync`, the skill record command and the Behavior Surface record command, in that order;
   - the three tests of `internal/docscontract/review_model_decision_test.go`.

   Files: `internal/cli/review.go`, `internal/cli/review_lineage.go`, `internal/cli/review_model_change_test.go` (new), `internal/docscontract/review_model_decision_test.go` (new), `docs/adr/0451-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed.md` (new), `docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md`, `.secondbrain-export`, `CONTEXT.md`, `docs/user-guide/commands/review.md`, `.agents/skills/roundfix/references/review.md`, `skills/roundfix/references/review.md`, `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/testdata/owned-skill-versions.json` and `docs/references/behavior-surfaces.json`.
2. The QA gate (depends on: 1).

Step 1 is one context and one Task. The code is three lines, and the documents describe only those lines. The guide page, the skill reference and the Behavior Surface Record cannot be split from one another, because a changed guide page fails the surface contract until its record is re-recorded and fails the release check until its skill changes. A second Task for the documents alone would pay an Agent turn and a repository Verification for about ten changed lines of prose.

`docs/references/skill-coverage.json` does not change: the map names no source file this design adds, and the three surfaces whose sources change are covered by the skill reference this step edits, so the PRD needs no `## Skills` entry.

Files shared with Spec 0418, which is delivered first: the two `SKILL.md` version fields, `skills/testdata/owned-skill-versions.json`, `docs/references/behavior-surfaces.json`, `CONTEXT.md` and `.secondbrain-export`. Step 1 regenerates the version fields and the two records on the tree it runs on and writes the next version, never a number, and no test or Verification of this Spec pins a version or a digest of a shared file. No Verification of this Spec depends on Spec 0418. The Task Graph still names Spec 0418 in `requires`, because both Specs change these files and the maintainer's third decision delivers this Spec after it: the Delivery Queue then starts this Spec only once Spec 0418 has merged, and step 1 runs on a tree that holds it.

## Risks & Considerations

- The tests use a fake runner; the real Codex adapter is not run. What supports the fix is that Codex compares the model a session was recorded with, and a new session has no recorded model (`references/sources.md`, section 7). What confirms it is the Fiscus's Delivery Retry of its Spec 0108 on the released version, which the operator's notice asks for.
- The report names the two models, not what changed the profile. Either cause, an edit of the adopter's review profile or an update of Roundfix that moves the built-in one, changes the model at the recorded index, which is what the condition reads. In v0.153.0 no path prepares a recorded session name on another index than the one it was opened on.
- Any other notice an adapter sends as agent text still blocks a no-findings verdict. That is outside this Spec (`_prd.md`, Non-Goals), and the classifier's rule is the reason it stays visible.
- Every lineage that is open when Roundfix is updated loses its continuation once, because its record holds no model. The round-2 prompt carries the findings and their dispositions either way.
- The model is compared as an exact string. Two names of one model, an alias and its full identifier, count as a change and cost a continuation, never a wrong one.
- A changed reasoning effort still continues the session. No report says a runtime answers differently for it; if one does, the same condition can take it.
- `lineage.model` is a new optional field of a machine-readable record, released under a `fix` at `patch`. It is additive, a reader that ignores unknown fields is unaffected, and the release notes name it.
- A Roundfix released before this Spec that runs round 2 over a record this Spec wrote ignores the model and continues as it always did.
- ADR-0451's number can be taken while the Pull Request waits. The operator then renames the file, its link in the notice of ADR-0197, its link in the glossary and its export line; invariant 7 keeps the test independent of the number.

## Glossary

- changes: **Reviewer Lineage**

## Decisions

- Not taken: removing or ignoring the runtime's notice before the answer is classified. The new session covers the reported case, and the classifier stays as strict as it is (ADR-0174: "The classifier's strictness is unchanged inside the final message."). See ADR-0451.
- Not taken: continuing the session after setting its recorded model back on it. The review would run on a model the operator no longer configures. See ADR-0451.
- Not taken: comparing the reasoning effort or the whole Agent Selection. See ADR-0451 and the first Delegated Decision below.
- Not taken: starting a new Reviewer Lineage at round 1 when the model changed. It would reread the whole candidate and reset the two-round ceiling. See ADR-0451.
- Not taken: continuing a session whose record holds no model. See ADR-0451 and the second Delegated Decision below.
- The model is recorded as the runtime holds it and compared with the selection's model after trimming; both are the same string, because the runtime copies the selection's model unchanged (`references/sources.md`, section 6).
- The condition is written where the recorded index is already checked, with no new function: a round 2 that does not continue needs no code of its own, since the two steps after the branch already end the recorded session and open a new one.
- ADR-0451 is written by the Task, not by the planning change, so the notice in ADR-0197 lands in the change that adds it (`_prd.md` → Decisions).
- Delegated Decision of 2026-11-23: What does round 2 compare before it continues the round-1 session?
- Delegated Decision of 2026-11-23: A round-1 record written before this fix holds no model. Does round 2 continue its session or open a new one?
- Delegated Decision of 2026-11-23: Which prompt does the new session of round 2 receive?
- Delegated Decision of 2026-11-23: Does the review print why it did not continue the round-1 session?
