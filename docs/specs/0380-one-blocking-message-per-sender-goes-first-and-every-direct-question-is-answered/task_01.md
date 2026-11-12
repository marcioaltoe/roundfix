---
task: task_01
spec: 0380-one-blocking-message-per-sender-goes-first-and-every-direct-question-is-answered
status: completed
type: backend
complexity: medium
---

# Task 01: The Follow-up Command gives the declared precedence to one message per sender

## Overview

Of the message entries with `blocked_work` that share a sender, only the
oldest keeps the `declared` basis that ranks it first without the judge.
Every later declared entry of that sender goes to the judge with the
undeclared items, keeps its `blocked_work`, and names the holding entry in a
new optional key `precedence_held_by` and on its text line. The command guide
and both `follow-up` reference copies document the bound.

## Requirements

1. MUST add `PrecedenceHeldBy string` with JSON tag
   `precedence_held_by,omitempty` to `followUpItem` in
   `internal/cli/follow_up.go`, declared right after `BlockedWork`.
2. MUST add `followUpDeclaredPrecedence(items []followUpItem) (declared, ranked []followUpItem)`
   to `internal/cli/follow_up.go` as `_techspec.md` → API Contract 1 says,
   with the helpers `followUpSender` and `followUpDeclaredBefore` that
   `_techspec.md` → Interfaces names: sender `strings.TrimSpace(*RequestedBy)`
   compared byte for byte, computed only by `followUpSender`; holding item
   the earliest `created_at` by `followUpDate`, then the lowest `Ref`;
   both results in input order and never nil; every other declared item of
   that sender in `ranked` with `PrecedenceHeldBy` set to the holding `Ref`
   and its `BlockedWork` kept.
3. MUST change `runFollowUpCommand` as API Contract 2 says: Waiting Items
   split as today, the rest through `followUpDeclaredPrecedence`, `declared`
   to `followUpDeclaredOrder` and `ranked` to the judge and `followUpOrder`,
   both unchanged. MUST append to a text line, after the existing
   ` — blocks <blocks>` part, ` — declares blocked work: <blocked_work> — precedence held by <precedence_held_by>`
   when `PrecedenceHeldBy` is non-empty (API Contract 3). MUST NOT read a
   message body to decide a declaration or a sender (Invariant 1).
4. MUST create `internal/cli/follow_up_declared_precedence_test.go` with the
   three tests `_techspec.md` → Testing Approach names for `internal/cli`,
   each calling `t.Parallel()` first and built on the inputs it lists:
   `TestFollowUpDeclaredPrecedenceHoldsOneEntryPerSender` (calls
   `followUpDeclaredPrecedence` directly),
   `TestFollowUpSendsALaterDeclaredMessageToTheJudge` (fake judge; also calls
   `renderFollowUpText` on the returned document) and
   `TestFollowUpLaterDeclaredMessageTranscript` (Surface Transcript 1 whole:
   every line of standard output, empty standard error, exit `0`, with only
   the transcript's placeholders replaced). No existing test changes.
5. MUST change `docs/user-guide/commands/follow-up.md` and the canonical
   `.agents/skills/roundfix/references/follow-up.md` as API Contract 5 says,
   word for word, then run `make skills-sync` to copy the reference to
   `skills/roundfix/references/follow-up.md`, and raise the Roundfix skill
   version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   never hand-edit a mirror, a version field or the record.
6. MUST, as the Task's last step, after its last guide edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   Of the Governed Paths the Task may change only
   `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/follow-up.md` and
   `skills/roundfix/SKILL.md`.

## Subtasks

- [ ] Add the key and the per-sender selection.
- [ ] Route held-back declared entries to the judge and render their line.
- [ ] Add the three tests, the transcript included.
- [ ] Document the bound, sync and record the skill version, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Only the oldest declared entry of each sender ranks with basis `declared`; a later one is judged or typed, keeps `blocked_work` and names the holding entry.
- [ ] A file with no two declared entries from one sender keeps its report bytes.
- [ ] The guide, both reference copies, the skill version and the Behavior Surface Record agree with the command.

## Context

- instruction: `docs/adr/0390-one-blocking-message-per-sender-takes-precedence-and-every-direct-question-is-answered.md`
- instruction: `docs/adr/0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md`
- instruction: `internal/cli/follow_up_test.go`
- instruction: `internal/cli/follow_up_blocked_work_test.go`
- interface: `internal/cli/follow_up.go`
- creates: `internal/cli/follow_up_declared_precedence_test.go`
- interface: `internal/cli/live_documents_docscontract_test.go`
- interface: `docs/user-guide/commands/follow-up.md`
- interface: `.agents/skills/roundfix/references/follow-up.md`
- interface: `skills/roundfix/references/follow-up.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The item's `precedence_held_by` and the per-sender `declared` basis — consumers: `./internal/cli`, `./skills`
- The command guide and the `follow-up` reference — consumers: `./internal/docscontract`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestFollowUpDeclaredPrecedenceHoldsOneEntryPerSender|TestFollowUpSendsALaterDeclaredMessageToTheJudge|TestFollowUpLaterDeclaredMessageTranscript)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestFollowUpDeclaredPrecedenceHoldsOneEntryPerSender TestFollowUpSendsALaterDeclaredMessageToTheJudge TestFollowUpLaterDeclaredMessageTranscript; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m -run 'FollowUp' ./internal/cli && go test -count=1 -timeout 30m -tags docscontract -run 'TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes' ./internal/cli` — expected: exit 0. Before this Task the three tests do not exist, so their pass lines are missing; afterwards they pass, every existing Follow-up Command test passes unchanged, and the guides list the item's fields in struct order.

- `for f in docs/user-guide/commands/follow-up.md .agents/skills/roundfix/references/follow-up.md skills/roundfix/references/follow-up.md; do tr -s '[:space:]' ' ' < "$f" | grep -qE -- 'holds the declared precedence' || { echo "missing bound: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qE -- 'never replaces the maintainer.s judgment' || { echo "missing judgment: $f" >&2; exit 1; }; grep -q 'precedence_held_by' "$f" || { echo "missing key: $f" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no guide names the bound or the key; afterwards the guide and both reference copies document them, and the skills suite with its recorded version and the documentation contracts with the Behavior Surface Record pass.

## References

- `_prd.md` → Goal 2; Goal 4; User Story 2; User Story 3; Core Feature 3; Core Feature 5; Success Metric 2; Success Metric 3; Success Metric 5
- `_techspec.md` → Current behavior 1-4; API Contracts 1-5; Surface Transcript 1; Invariants 1-4; Declared breaks; Testing Approach; Build Order 1
- ADR-0390


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate, typesafe-ai, golang-safety, systematic-debugging.

The command selects one holding declaration per trimmed, byte-compared sender,
using the earliest parsed date and then the lowest reference. Selection keeps
both slices non-nil and in input order without mutating the inventory. Waiting
Items retain their existing split. Later declarations reach the existing judge
and ordering code, retain `blocked_work`, and name their holder through the
optional `precedence_held_by` field and the prescribed text suffix. Neither
sender selection nor declaration detection reads a message body.

Acceptance evidence:

- One declaration per sender: the new selection test covers the six prescribed
  inputs, including the reference tie-break and trailing sender space, holder
  references, retained blocked work, input order, unchanged input and empty
  non-nil results. The fake-judge test requires exactly the two prescribed
  requests and proves that the later declaration ranks as health with basis
  `judge`, null `blocks`, and its holding reference. The third test matches the
  entire Surface Transcript 1 stdout, empty stderr and exit 0 without a judge
  request. Existing tests were not edited.
- Byte compatibility: a temporary `go test -overlay` probe,
  `TestFollowUpDistinctSenderCompatibilityProbe`, compared the distinct-sender
  fixture against `013af3b4:internal/cli/follow_up.go`. Text and JSON were
  identical after normalizing only the independently created fixture roots and
  Git heads. Evidence: `/private/tmp/roundfix-0380-compat-3ww_0zzh/`.
  The initial probe identified differing temporary commit hashes; correcting
  that fixture normalization required no product change.
- Documentation and records: all three guides have the exact prescribed
  Declared blocked work section and JSON fields in struct order, checked by
  a focused Python comparison. The canonical reference and embedded reference
  match byte for byte. The version recorder raised both Roundfix SKILL.md
  copies from `0.1.218` to `0.1.219` and recorded the new digest. The final
  Behavior Surface Record regeneration changed only the follow-up guide's
  fingerprint. No mirror, version field or generated record was hand-edited.

Focused checks and regeneration (all exit 0 unless noted):

- Initial red: `go test ./internal/cli -run '^TestFollowUpDeclaredPrecedenceHoldsOneEntryPerSender$' -count=1 -timeout=5m`
  exited 1 because the selection function and field did not exist yet.
- `go test ./internal/cli -run '^TestFollowUp(DeclaredPrecedence|SendsALater|LaterDeclared)' -count=1 -timeout=5m`
- `go test ./internal/cli -run '^TestFollowUp(DeclaredBlockedWorkTranscripts|OnlyDeclaredItemsAskNoJudge|OrdersByTypeWithoutTheJudge|ListHelpMatchesTranscript|SaysNothingToDo)$' -count=1 -timeout=5m`
- `make skills-sync`
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
- `make baseline-digests`: no derived changes.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  ran after the last guide edit as the final implementation step.
- `git diff --check` and the bounded changed-path inspection: clean; only this
  Task's declared paths changed, including the Daemon's pre-existing status
  update in this Task file.

The executor's symlinked writable-root rejection required running commands
with escalation at the physical checkout. Judge requests used the existing
per-test fake transport; no live API request or external documentation lookup
was made under the Spec's network restriction. No Task status, Task Graph,
other Task, commit, push or pull request was changed by this Agent. Declared
Verification and settlement remain with the Daemon.
