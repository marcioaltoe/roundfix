---
task: task_02
spec: 0335-follow-up-takes-every-documented-type-and-an-absent-spec-root-is-empty
status: completed
type: backend
complexity: medium
---

# Task 02: The follow-up command takes every documented message type, and its guides name the JSON it emits

## Overview

The `--messages` parser and the type fallback of the Follow-up Priority read
one table of accepted message types, a refused type's error lists that
table, and both follow-up guides name the table and the JSON fields the
command emits, held there by a documentation contract test. The modes
reference tells the agent that writes the messages file to use only those
types. This answers the first and third frictions of the adopted Backlog
Entry of 2026-10-30 (the follow-up command rejects documented message types
and an empty Spec Root) and documents the second for the `specs` source.

## Requirements

1. MUST characterize first (lesson 9): write
   `TestFollowUpAcceptsEveryDocumentedMessageType` in
   `internal/cli/follow_up_test.go` and observe on the current tree that the
   message with type `other` makes the command exit 2 with `invalid type
   "other"`; then implement.
2. MUST implement `_techspec.md` → API Contract 4 in
   `internal/cli/follow_up_inventory.go`: the `followUpMessageTypes` table
   (`feat`, `feature`, `finding`, `fix`, `information`, `other`, `perf`,
   `refactor`, each with its class), the parser reading it, the item keeping
   the given name as its `type`, an empty `type` still stored as `message`,
   and the error `invalid type %q; use one of <names joined by ", ">, or omit
   type` built from the table inside today's `messages file %q: ` prefix.
3. MUST implement API Contract 5 in `internal/cli/follow_up.go`:
   `followUpOrder` takes an unanswered item's class from `followUpTypeClass`,
   which reads the same table, gives `health` to `Kind == "finding"` and
   `other` to any unlisted type, so every item the current code classifies
   keeps its class (Invariant 5).
4. MUST add the tests `_techspec.md` → Testing Approach names for task_02,
   each failing when its behavior is removed:
   `TestFollowUpAcceptsEveryDocumentedMessageType` (one message per table
   name and one without `type`, no judge key: exit 0, each item's `type` as
   given and its `priority_name` as `_prd.md` → Success Metric 1 states,
   written as literal expectations) and
   `TestFollowUpRefusesAnUnknownMessageTypeNamingTheAcceptedOnes` (type
   `chore`: stdout empty, exit 2, stderr exactly the error line built from
   `followUpMessageTypes` and then `Run 'roundfix follow-up --help' for
   usage.`) in `internal/cli/follow_up_test.go`; and
   `TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes` in the existing
   `internal/cli/live_documents_docscontract_test.go`, which reads
   `docs/user-guide/commands/follow-up.md`,
   `.agents/skills/roundfix/references/follow-up.md` and
   `skills/roundfix/references/follow-up.md`, takes the backticked names
   after each sentence of API Contract 6 up to its period, and requires them
   equal, in order, to the JSON tags of the named type read by reflection,
   and to the table's names for the message type sentence. A guide without a
   sentence fails the test. Each top-level test calls `t.Parallel()` first.
5. MUST rewrite the field and type prose of
   `docs/user-guide/commands/follow-up.md` and the canonical
   `.agents/skills/roundfix/references/follow-up.md` with the sentences of
   API Contract 6: "The list report carries these JSON fields:", "The status
   report carries these JSON fields:", "Each source carries these JSON
   fields:", "The judge carries these JSON fields:", "Each item carries these
   JSON fields:" and "A message type is empty or one of", each followed by its
   names. Both say which class each type takes without a judge answer, that a
   refused type's error lists the accepted ones, and that the `specs` source
   of a repository whose Spec Root is absent is `checked` with count `0`
   (ADR-0348). The reference MUST no longer contain `class, kind,
   reference`, `requester` or `blocking project` as item field names, and its
   judge sentence names `cost_usd`.
6. MUST add to the canonical `.agents/skills/roundfix/references/modes.md`,
   in step 1 of the `--follow-up` mode where the messages file is written,
   that a message entry gets a `type` only when the message's type hint is
   one of the types the Follow-up Command reference lists, "and omit the type
   field otherwise". New text goes under existing headings, never inside a
   `### QA settlement` section.
7. MUST copy the canonical skill to `skills/` with `make skills-sync`, then
   raise both version fields of `.agents/skills/roundfix/SKILL.md` and its
   mirror with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
   -record-skill-versions -count=1`, recording
   `skills/testdata/owned-skill-versions.json`.
8. MUST re-record the Behavior Surface fingerprint as the last step, after
   the last guide edit: `go test -tags docscontract ./internal/docscontract
   -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
   -count=1`, recording `docs/references/behavior-surfaces.json`.

## Subtasks

- [ ] Characterize the refusal of type `other`.
- [ ] Add the type table and read it in the parser, its error and the ordering fallback.
- [ ] Add the two command tests and the guide contract test.
- [ ] Rewrite the follow-up guide and reference prose; add the modes rule.
- [ ] Sync, record the skill version, re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] Every documented type is accepted with its class, and `chore` is refused with the accepted list.
- [ ] Both follow-up guides name the JSON fields and the types the command uses, and the contract test fails when either drifts.
- [ ] The modes reference tells the agent to omit an undocumented type.

## Context

- instruction: `docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md`
- instruction: `docs/adr/0348-an-absent-spec-root-inside-the-repository-holds-no-spec.md`
- interface: `internal/cli/follow_up_inventory.go`
- interface: `internal/cli/follow_up.go`
- interface: `internal/cli/follow_up_test.go`
- interface: `internal/cli/live_documents_docscontract_test.go`
- interface: `docs/user-guide/commands/follow-up.md`
- interface: `.agents/skills/roundfix/references/follow-up.md`
- interface: `skills/roundfix/references/follow-up.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `roundfix/follow-up-messages/v1` accepted types and the Follow-up Priority fallback — consumers: `./internal/cli`
- The owned-skill version record and the Behavior Surface record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestFollowUp.*)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestFollowUpAcceptsEveryDocumentedMessageType TestFollowUpRefusesAnUnknownMessageTypeNamingTheAcceptedOnes TestFollowUpOrdersByTypeWithoutTheJudge TestFollowUpRefusesAnInvalidMessagesFile TestFollowUpOrdersItemsByTheFourAnswers; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; doc="$(go test -count=1 -timeout 20m -v -tags docscontract -run '^TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes$' ./internal/cli 2>&1)" || { printf '%s\n' "$doc"; exit 1; }; printf '%s\n' "$doc" | grep -q -- '--- PASS: TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes ' || { echo 'missing pass: TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes' >&2; exit 1; }; GOOS=windows go build -buildvcs=false ./internal/cli` — expected: exit 0. Before this Task the type tests and the guide contract test do not exist, so the first pass check fails. Afterwards they pass with the existing follow-up tests that consume the parser and the ordering, and the package builds for Windows.
- `for phrase in 'Each item carries these JSON fields:' 'The judge carries these JSON fields:' 'A message type is empty or one of'; do tr -s '[:space:]' ' ' < docs/user-guide/commands/follow-up.md | grep -qF -- "$phrase" || { printf 'missing phrase in the follow-up guide: %s\n' "$phrase" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/follow-up.md | grep -qF -- "$phrase" || { printf 'missing phrase in the follow-up reference: %s\n' "$phrase" >&2; exit 1; }; done; if grep -qE 'requester|blocking project' .agents/skills/roundfix/references/follow-up.md; then echo 'the follow-up reference still names prose item fields' >&2; exit 1; fi; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/modes.md | grep -qF -- 'and omit the type field otherwise' || { echo 'missing phrase in the modes reference' >&2; exit 1; }; for ref in follow-up modes; do cmp ".agents/skills/roundfix/references/$ref.md" "skills/roundfix/references/$ref.md" || exit 1; done; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither guide carries the field sentences, so the first phrase check fails. Afterwards both carry them, the reference no longer uses the prose field names, the modes reference carries its rule, the mirrors equal their canonical copies, and the owned-skill version and documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Features 1, 4, 5; User Experience; Non-Goals / Out of Scope; Success Metric 1; Success Metric 2; Success Metric 5
- `_techspec.md` → Current behavior 3-5; Interfaces; API Contracts 4-6; Invariants 5-6; Testing Approach; Build Order 2
- ADR-0319; ADR-0348
- [Adopted Backlog Entry](references/2026-10-30-the-follow-up-command-rejects-documented-message-types-and-an-empty-spec-root.md)


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-testing, testing-boss, systematic-debugging, no-workarounds, tech-writer, evidence-gate.

Implemented this Task's message-type and guide-contract slice. Task status and
checkboxes remain Daemon-owned. The initial worktree change was this Task's
`pending` → `in_progress` status; this turn preserved it.

### Acceptance evidence

- Every documented type is accepted with its class, and `chore` is refused
  with the accepted list: `followUpMessageTypes` now supplies the parser,
  refusal list and `followUpTypeClass`. The fallback preserves the Finding
  kind override and assigns unlisted types to `other`; accepted names are
  retained and omitted/empty types remain `message`.
  `TestFollowUpAcceptsEveryDocumentedMessageType` asserts literal expectations
  for all eight names and an omitted type through the CLI without a judge
  key. `TestFollowUpRefusesAnUnknownMessageTypeNamingTheAcceptedOnes` asserts
  empty stdout, exit 2 and the entire table-derived error and usage line.
- Both follow-up guides name the emitted JSON fields and accepted types:
  the new docscontract test reads the user guide, canonical reference and
  mirror, requires all six sentences, and compares their ordered backticked
  names against reflected JSON tags or the shared type table. Both guides
  describe each type's fallback class and absent Spec Root `specs: checked 0`.
  The reference names `cost_usd` and no longer names the old prose item fields.
- The modes reference tells the agent to omit an undocumented type: step 1
  now permits only type hints listed by the Follow-up Command reference and
  says `and omit the type field otherwise`. A focused Python inspection
  confirmed this rule within step 1 in both copies. Existing mode order and
  handoff/message-settlement tests still pass. No QA settlement text changed.

### Characterization and focused checks

- Before production edits:
  `go test ./internal/cli -run '^TestFollowUpAcceptsEveryDocumentedMessageType$/other$' -count=1 -v`
  failed at the new test with command exit 2, empty stdout and
  `invalid type "other"`, followed by the usage line.
- Before guide edits:
  `go test -tags docscontract ./internal/cli -run '^TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes$' -count=1 -v`
  failed for all three copies, identifying the six missing sentences.
- After implementation:
  `go test ./internal/cli -run '^TestFollowUp(AcceptsEveryDocumentedMessageType|RefusesAnUnknownMessageTypeNamingTheAcceptedOnes|OrdersByTypeWithoutTheJudge|RefusesAnInvalidMessagesFile|OrdersItemsByTheFourAnswers)$' -count=1 -v`
  exited 0; the two new command tests and all three selected existing tests
  passed, including the nine accepted/omitted-type cases.
- After guide edits and synchronization, the focused docscontract command
  above exited 0 for all three copies.
- Negative probes used a temporary compiled docscontract test binary,
  isolated Git fixtures, and `go test -overlay` source mutations; repository
  inputs were not mutated. All 18 probes failed at the intended assertions:
  five per guide copy (missing/added item field, missing/added message type,
  missing judge sentence), plus parser rejection of `other`, loss of the
  `finding` type's `health` fallback, and loss of the accepted list in the
  refusal. Isolated guide baselines passed before and after the mutations.
- `go test ./skills -run '^TestTheFollowUpMode(RunsItsStepsInOrder|TakesTheHandoffAndSettlesEveryMessage)$' -count=1 -v`
  exited 0 for canonical and embedded layouts. Focused Python inspection also
  confirmed byte equality of both SKILL.md files and both changed references,
  both version fields at `0.1.163`, and absence of the old prose field names.

### Sanctioned regeneration

- `make skills-sync` copied the canonical references to `skills/`.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`
  exited 0, raised both version fields in both SKILL.md files from `0.1.162`
  to `0.1.163`, and appended the skill digest to the owned-skill version record.
- The repository-required `make baseline-digests` exited 0 with
  `changed: false`; no Baseline artifact changed.
- As the final implementation/regeneration step after the last guide edit,
  `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  exited 0 and recorded `docs/references/behavior-surfaces.json`.

Declared Verification, Windows build and repository-wide gates were not run;
they remain with the Daemon. No commit, push or pull request was performed.
No follow-up implementation was added outside this slice.
