---
spec: 0304-the-review-drops-a-result-status-note-and-a-stale-fingerprint-names-its-command
prd: _prd.md
created: 2026-10-22
---

# The review drops a Result Status Note, and a stale fingerprint names its command — Technical Spec

## Executive Summary

`validateReviewConventions` in `internal/cli/review_convention_validator.go`
sends every anchored finding inside a convention region to the sealed
validator, and only the validator's verdict can dismiss it. For the pair "the
Result says `in_progress`, the front matter says `completed`" that verdict
rests on C2's last sentence, "A Result that calls status Daemon-owned agrees
with a `completed` status.", which the validator read twice as a requirement.
C2 also says the Daemon writes the Result, which the Daemon's status ownership contradicts: the
Agent writes it during its turn. This Spec rewords C2 (version 4) and adds a
mechanical step before the validator, `reviewResultStatusNote`, that dismisses
the pair with rule `convention:C2` and a fixed reason when the head tree and
the finding's text both prove it. The stale Behavior Surface fingerprint
problem in `internal/docscontract/skill_coverage_test.go` gains the re-record
command, held equal to the `.roundfixrc.yml` declaration by a contract test.

The trade-off accepted: a finding that names the pair and also another failure
in the same list item is dismissed with the pair; ADR-0317 records why.

## Project Constraints

- Identifier strategy: not applicable — no identifier or record schema is
  minted; the dismissal uses the existing rule value `convention:C2`. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, credential or
  recipient is added; tests use temporary repositories, a fake reviewer and a
  fake sealed validator. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Order 1 follows ADR-0317 and
  keeps the fail-closed validator for every other finding and the Daemon's
  ownership of Task status, both cited in ADR-0317. It is named with its
  title in `_prd.md` → Project Constraints. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 ("same as cycle AO"), carrying the cycle AJ grant
  "Mexer em Governed Paths declarados". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0304-the-review-drops-a-result-status-note-and-a-stale-fingerprint-names-its-command/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/review.md`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Read and characterized on 2026-10-22 at `42f523c2`. Line numbers move when
other Specs merge; the function names do not.

1. `deliveryConventions` (`internal/cli/review_conventions.go:23`) returns
   C1 to C5 under `deliveryConventionsVersion`
   `roundfix/delivery-conventions/v3`. C2 says "The Daemon writes a Task
   file's status after the Task's Verification passes, its `## Result` and
   `## Carry-forward provenance` sections for non-QA Tasks only, …" and ends
   "A Result that calls status Daemon-owned agrees with a `completed`
   status." `buildReviewPrompt` (`internal/cli/review.go`) and
   `runConventionValidator` carry the list verbatim.
2. The Daemon's status ownership decision and `taskExecutionInvariants` (`internal/agent/spec_prompt.go`)
   say otherwise about the Result: the Daemon writes `in_progress` before the
   Agent turn, the Agent appends `## Result`, and the Daemon writes
   `completed` after Verification. The prompt line "The Daemon owns Task
   status during Implement; do not edit the task file's status field." is
   what the Agents echo as "status remains `in_progress`".
3. `conventionRegions` gives `C2` when the whole anchor lies in a Task file's
   front matter or its `## Result`, `## Recorded paths` or `## Carry-forward
   provenance` sections (`reviewTaskSettlementRegion`).
   `validateReviewConventions` asks the validator about every anchored finding
   with an eligible rule, marks it `dismissed-by-validation` only on a
   `dismiss` verdict, and returns exit 1 (`exitRunFailed`) when nothing is
   asked or any finding stands.
4. Characterization, a `go test -overlay` probe in `internal/cli` that wrote
   nothing to the repository: through `roundfix review` with a fake reviewer
   answering the 0297 finding text and a fake validator answering the 0297
   `stands` verdict, over a committed Task of the 0297 shape (`status:
   completed`, Result "they remain Daemon-owned. Task status remains
   `in_progress`"), the validator was asked once with `[convention:C2]`, its
   input held the v3 C2 sentence, and the review exited 1 with outcome
   `findings` and the finding standing.
5. The real records (`pre-pr-review.json` of review heads `61dec34e` and
   `1aa6ae3a`) hold validator reasons citing "C2's required agreement with a
   completed status" and "the required completed-status agreement".
6. `TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem`
   (`internal/docscontract/command_documentation_test.go`) requires the
   version and each convention text in `docs/user-guide/commands/review.md`
   and both copies of the Roundfix review reference;
   `TestDeliveryConventionC2ExemptsTheQATask`
   (`internal/cli/review_convention_qa_task_test.go`) pins the v3 version and
   C2 text; `TestReviewValidatorDismissesADaemonSettlementAsC2` sends the pair
   "The Task is `completed` while its Result says status is Daemon-owned." to
   the validator over a Result "Status is Daemon-owned.".
7. `skillCoverageProblems` (`internal/docscontract/skill_coverage_test.go`)
   reports `<id>: stale fingerprint` and nothing else;
   `TestAStaleBehaviorSurfaceIsReported` pins that text and the `missing` and
   `extra` texts. `.roundfixrc.yml` declares, under `delivery.derived_paths`
   for `docs/references/behavior-surfaces.json`, the regeneration `go test
   -count=1 -tags docscontract ./internal/docscontract -run
   '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

Sweep of the wording and the readers (`grep -rn
"delivery-conventions/v\|calls status Daemon-owned\|stale fingerprint"` over
`internal`, `docs`, `.agents`, `skills` and `CONTEXT.md`, excluding
`docs/history`): the version and C2 text appear only in the two Go files and
the three documents above; "stale fingerprint" only in
`skill_coverage_test.go`. The rule `convention:C2` is read by
`reviewFindingDismissedByValidation`, `review dispose` (which refuses a
validation dismissal) and the stderr diagnostic; none changes.

## System Architecture

- `internal/cli/review_conventions.go`: C2 version 4.
- `internal/cli/review_convention_validator.go`: `reviewResultStatusNote`,
  its fixed reason, and one settlement helper for "no question left" and
  "after the verdicts".
- `internal/docscontract/skill_coverage_test.go`: the re-record hint.
- `docs/user-guide/commands/review.md`, the Roundfix review reference and its
  mirror, the Roundfix skill's versions, `CONTEXT.md`.

No package, import edge or schema changes.

## Implementation Design

### Interfaces

```go
// reviewResultStatusNoteReason is the fixed reason of the dismissal.
const reviewResultStatusNoteReason = "Result Status Note: the Agent wrote this Result while the Daemon held the Task in_progress, and the Daemon wrote completed after the Task's Verification passed (C2)"

// reviewResultStatusNote reports whether a finding names the C2 pair of a
// completed non-QA Task whose Result holds a Result Status Note.
func reviewResultStatusNote(ctx context.Context, repo reviewRepository, finding reviewFinding, conventions []string) (bool, error)

// In internal/docscontract (test-only code).
const behaviorSurfaceRecordCommand = "go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage"
```

### Data Models

None changes shape. A dismissed finding's `validation` is `{"status":
"dismissed-by-validation", "rule": "convention:C2", "reason":
reviewResultStatusNoteReason}`; the record's `validation.conventions` reads
`roundfix/delivery-conventions/v4`.

### API Contracts

1. API Contract: C2 version 4. `deliveryConventionsVersion` is
   `roundfix/delivery-conventions/v4` and C2's text is exactly: "The Daemon
   writes a Task file's status after the Task's Verification passes, its `##
   Carry-forward provenance` section for non-QA Tasks only, and `## Recorded
   paths` only when a Task changed a path it did not declare; a QA Task
   settles on its committed QA Report and carries none of these sections. The
   Agent writes a non-QA Task's `## Result` during its turn, while the Daemon
   holds the Task `in_progress`, so a Result that reports the status as
   `in_progress`, unchanged or Daemon-owned states what was true when it was
   written and does not contradict a `completed` status." C1, C3, C4 and C5
   keep their bytes.
2. API Contract: the Result Status Note test. `reviewResultStatusNote`
   returns true only when all hold: the finding has an anchor and
   `conventions` contains `C2`; the lower-cased finding text contains
   `completed` and one of `in_progress`, `in progress`, `daemon-owned`; the
   anchored file at the review head has front matter with exactly one line
   starting `status:` whose trimmed value is `completed` and no line starting
   `type:` whose trimmed value is `qa`; and a line after the front matter,
   inside the `## Result` section (from that heading to the next `## `
   heading), contains `in_progress` or `daemon-owned` once lower-cased. A Git
   read failure returns an error.
3. API Contract: the dismissal. In `validateReviewConventions`, after
   `conventionRegions` and before the rule list is built, a finding for which
   `reviewResultStatusNote` is true gets the validation of Data Models and is
   not asked; an error from it makes the validator `unavailable` as a
   `conventionRegions` error does. When no finding is left to ask, and after
   the verdicts are applied, the same helper settles the record: outcome
   `findings-dismissed` and exit 0 when every finding is dismissed by
   validation or by the operator's evidence, exit 1 otherwise. A record
   whose validator was never asked keeps `validator` `not-needed`.
4. API Contract: unchanged paths. Anchor validation, the validator prompt,
   grammar and fail-closed rules, `review dispose`, the record schema and the
   stderr diagnostic line format keep their bytes; the diagnostic prints the
   new reason through the existing format.
5. API Contract: the stale fingerprint problem.
   `skillCoverageProblems` reports `<id>: stale fingerprint; re-record
   docs/references/behavior-surfaces.json after the last edit of its sources
   with: ` followed by `behaviorSurfaceRecordCommand`, the path written
   through `skillcoverage.RecordPath`. Every other problem keeps its bytes.
6. API Contract: one command. `TestTheStaleFingerprintProblemNamesTheProjectRegeneration`
   parses `.roundfixrc.yml` with `gopkg.in/yaml.v3`, an existing module
   dependency, collects the `regenerate` of every `delivery.derived_paths`
   entry whose `paths` holds `skillcoverage.RecordPath`, and requires exactly
   `[behaviorSurfaceRecordCommand]`; it also requires the flag
   `record-skill-coverage` to be registered in the package and the command to
   name ` -record-skill-coverage` and `-run '^TestTheSkillCoverageMapIsCurrent$'`.
7. API Contract: documentation. The review guide and both reference copies
   carry version 4 and the new C2 text, a paragraph describing the Result
   Status Note dismissal, and a validator sentence that excludes it; the
   glossary adds **Result Status Note** and revises **Delivery Convention**.

### Surface Transcripts

None. No command, flag, exit code or stream placement changes; the stderr
dismissal line keeps its existing format with the new reason, asserted whole
by the tests of API Contract 3.

### Invariants

```text
1. A finding is dismissed without the validator only when the head tree and its text both prove the pair of API Contract 2.
2. Every other anchored finding with an eligible rule reaches the validator, which still fails closed.
3. A dismissed finding keeps the rule convention:C2 and records the fixed reason; it is never dropped from the record.
4. A review whose findings are all dismissed exits 0 with findings-dismissed whether or not the validator ran.
5. C1, C3, C4, C5 and the documents' other convention texts keep their bytes.
6. The stale fingerprint problem names exactly the command .roundfixrc.yml declares for the record; other problems keep their bytes.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 2, Success Metrics 1 and 3, Invariants 1, 3 and 4 → the Result Status Note test and the dismissal (API Contracts 2 and 3)
- Goal 2, Success Metrics 2 and 3, Invariant 2 → the dismissal's negative cases and unchanged paths (API Contracts 2-4)
- Goal 3, Core Features 1 and 3, Success Metric 4, Invariant 5 → C2 version 4 and documentation (API Contracts 1 and 7)
- Goal 4, Core Feature 4, Success Metric 5, Invariant 6 → the stale fingerprint problem and one command (API Contracts 5 and 6)

## Absence Evidence

1. `reviewResultStatusNote` — command: `grep -rn -e ResultStatusNote -e 'Result Status Note' internal docs .agents skills CONTEXT.md`; result: no match at `42f523c2`; `validateReviewConventions` dismisses only on a validator verdict.
2. The re-record hint — command: `grep -rn -e 'stale fingerprint;' -e behaviorSurfaceRecordCommand internal`; result: no match at `42f523c2`; the problem is the bare `<id>: stale fingerprint`.
3. A contract holding the problem's command to `.roundfixrc.yml` — command: `grep -rn -e derived_paths internal/docscontract`; result: no match at `42f523c2`.

## ADR Alternatives

- ADR-0317: no rejected alternative is taken; C2 is not only reworded, the Agent prompt is left alone, and neither the anchor alone nor every Result anchor dismisses a finding. The validator decision it refines keeps its own rejected alternative, a judgment model as validator: the new step is code, not a model, and the validator keeps failing closed for every other finding.

## Integration Points

Git, read at the review head through the existing `reviewRepository`; the
fake reviewer and sealed validator in tests. No test, Verification command or
QA row calls a provider, reads the real `~/.roundfix` or reaches the network.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows and
passes no real-time duration.

- `internal/cli/review_result_status_note_test.go` (new), with the existing
  `validatorFixture`, `conventionValidatorRunner` and `validatorQuestions`:
  - `TestReviewDismissesAResultStatusNoteOfACompletedTaskWithoutTheValidator`,
    subtests `0297 shape` (a completed backend Task whose Result ends with the
    0297 paragraph, anchored on its first line, which holds no note itself),
    `0290 shape` (a completed docs Task whose Result opens with the 0290
    sentence, anchored on the blank line after `## Result`) and `front matter
    anchor` (anchored on the `status:` line): exit 0, `findings-dismissed`,
    `validator` `not-needed`, no sealed request, rule `convention:C2`, the
    reason constant, and the whole stderr line.
  - `TestReviewKeepsAContradictionThatIsNotAResultStatusNote`, subtests
    `failed Task`, `QA Task`, `Result without a note`, `note outside the
    Result` (the note only under `## Requirements`) and `finding names another
    failure` (text with `completed` but neither status token): one question
    with first rule `convention:C2`, the fake verdict `stands` kept, exit 1;
    and `anchor outside the settlement region` (an anchor under `##
    Requirements` with the pair's text): no question, the finding stands,
    exit 1.
- Declared breaks: `TestDeliveryConventionC2ExemptsTheQATask` moves to the v4
  version and C2 text; `TestReviewValidatorDismissesADaemonSettlementAsC2`
  states the pair, which is now dismissed without a question, so its finding
  text becomes "The Task is `completed` while its Result gives no evidence.
  Failure: an invalid Task state transition." and every other assertion of it
  stays.
- `internal/docscontract/skill_coverage_test.go`:
  `TestAStaleBehaviorSurfaceIsReported` expects the new stale text (its
  `missing` and `extra` cases keep theirs), and
  `TestTheStaleFingerprintProblemNamesTheProjectRegeneration` (new) holds API
  Contract 6.
- Prototype record: the design was prototyped on 2026-10-22 in overlay files
  outside the tree; these sabotages each failed the named case: completed
  status check → `failed Task`; QA type check → `QA Task`; Result section
  restriction → `note outside the Result`; Result Status Note requirement →
  `Result without a note` and `note outside the Result`; finding text check →
  `finding names another failure`; C2 region requirement → `anchor outside the
  settlement region`; settling without a question returns exit 1 → every
  subtest of the dismissal test; stale problem without the hint →
  `TestAStaleBehaviorSurfaceIsReported/stale`; hint added to the missing
  problem → `/missing`; command constant changed →
  `TestTheStaleFingerprintProblemNamesTheProjectRegeneration`.

## Build Order

1. C2 version 4, the Result Status Note dismissal and its tests, the two
   declared breaks, the review guide, the Roundfix review reference, the
   skill's versions, the Behavior Surface Record, **Result Status Note** and
   **Delivery Convention** (no dependency; the convention text and its
   documents must change together, because the documentation contract holds
   them equal).
2. The stale fingerprint problem and its contract test (no dependency: it
   touches only `skill_coverage_test.go`).
3. QA gate (depends on: 1, 2).

No Spec is listed in `requires`. Specs 0299, 0300 and 0301 to 0303 deliver
first and may raise the Roundfix skill's version; this Spec rebases on them
and touches none of their files apart from the skill's version fields.

## Risks & Considerations

- A reviewer that puts two failures in one list item, one of them the pair,
  loses the other with the dismissal; the prompt already asks for one finding
  per item.
- A Result Status Note anywhere in the Result enables the dismissal for an
  anchor elsewhere in the C2 region, because the 0290 and 0297 anchors fell on
  lines without the note; the finding's text must still name the pair.
- A completed Task whose Agent hand-edited the status is indistinguishable in
  the diff from the Daemon's write; the Daemon's normalization already makes an
  Agent-written status no verdict.

## Glossary

- adds: **Result Status Note**
- changes: **Delivery Convention**

## Decisions

- C2 is reworded and versioned, and the pair is dismissed by code before the
  validator; see ADR-0317.
- The dismissal reuses `convention:C2`, so `review dispose` and every rule
  reader keep their behavior.
- The re-record command is a test constant held equal to `.roundfixrc.yml`.
