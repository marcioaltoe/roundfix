---
task: task_04
spec: 0325-the-follow-up-mode-consumes-the-handoff-and-every-message
status: completed
type: test
complexity: low
---

# Task 04: The modes test catches the restored "without consuming" sentence and names the status refusal

## Overview

Corrective Task 1 of 2. It answers findings F-01 and F-02 of the QA Report of
2026-10-26 (Run `run_20261026T232417Z_2824de576977e0a5`, rows Q05f and Q05c,
task_03 Requirement 5). The delivered skill text is correct; only
`TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage` in
`skills/roundfix_modes_test.go` is too weak or names the wrong diagnostic:

- F-01: replacing "List this project's messages with `memory_message_list`"
  with "List this project's messages without consuming them" in
  `skills/roundfix/references/modes.md` keeps the test green. The negative
  check searches the lowercase phrase "list this project's messages without
  consuming them" case-sensitively, the mutation starts with an uppercase
  "List", and no positive assertion requires the removed
  `memory_message_list` sentence.
- F-02: removing "accepts no handoff, " from `## --status` fails the test, but
  with "--status omits its handoff refusal" instead of the diagnostic task_03
  Requirement 5 declares, "--status section does not say it accepts no
  handoff".

## Requirements

1. MUST make the `Optional sources` check in
   `TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage` require the
   sentence "List this project's messages with `memory_message_list`"
   (whitespace-normalized, as the other phrases are) and fail with the exact
   message `Optional sources still lists messages only to leave them pending`
   when that sentence is missing or when the section contains "without
   consuming them" in any letter case. Keep every other `Optional sources`
   phrase check unchanged.
2. MUST change the `--status` diagnostic to the exact message
   `--status section does not say it accepts no handoff`, keeping the phrase
   it checks unchanged.
3. MUST keep both layouts (`canonical` and `embedded`) running the same
   checks, and every other assertion and diagnostic of the test unchanged.
4. MUST prove both negative controls before handing back, each in its own
   temporary `git archive` copy of the tree outside the checkout, running
   `go test -timeout 10m -count=1 -run '^TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage$' ./skills`
   there:
   - replace "List this project's messages with `memory_message_list`" with
     "List this project's messages without consuming them" in
     `skills/roundfix/references/modes.md`: the `embedded` subtest fails with
     `Optional sources still lists messages only to leave them pending`;
   - remove "accepts no handoff, " from `## --status` in
     `skills/roundfix/references/modes.md`: the `embedded` subtest fails with
     `--status section does not say it accepts no handoff`.
   Write nothing under the checkout for these proofs, and leave no copy behind.
5. MUST NOT change the skill text, its mirrors, its version record, any other
   test, guide, record or planning file. The findings are test gaps, not
   product defects.

## Subtasks

- [ ] Require the `memory_message_list` sentence and reject "without consuming them" in any case, with the declared diagnostic.
- [ ] Rename the `--status` diagnostic to the declared one.
- [ ] Prove both sabotages fail with the declared messages in temporary copies.

## Acceptance Criteria

- [ ] The test passes on the delivered skill text in both layouts.
- [ ] The F-01 mutation fails `embedded` with `Optional sources still lists messages only to leave them pending`.
- [ ] The F-02 mutation fails `embedded` with `--status section does not say it accepts no handoff`.

## Consumers

None. The Task changes one test file and no contract another package reads.

## Context

- interface: `skills/roundfix_modes_test.go`
- instruction: `skills/roundfix/references/modes.md`
- instruction: `docs/specs/0325-the-follow-up-mode-consumes-the-handoff-and-every-message/task_03.md`

## Verification

- `grep -qF 'Optional sources still lists messages only to leave them pending' skills/roundfix_modes_test.go && grep -qF -- '--status section does not say it accepts no handoff' skills/roundfix_modes_test.go && go test -timeout 10m -count=1 -run '^TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage$' ./skills && go test -timeout 15m -count=1 ./skills` — expected: exit 0. Before this Task neither declared diagnostic exists in the test file, so the first check fails; afterwards both diagnostics exist, the named test passes in both layouts and the whole `skills` package passes.

## References

- `_techspec.md` → Testing Approach
- `task_03.md` → Requirement 5

## Result

Corrected the `Optional sources` and `--status` checks in
`TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage`
(`skills/roundfix_modes_test.go`). No skill text, mirror, version record,
guide, record or planning file was touched; only the one test file changed.

### Requirement 1 — `Optional sources` requires the `memory_message_list` sentence and rejects "without consuming them" in any case

The single old negative-only check

```go
if strings.Contains(optionalSources, roundfixWhitespace("list this project's messages without consuming them")) {
    t.Error("Optional sources still leaves messages unconsumed")
}
```

became

```go
if !strings.Contains(optionalSources, roundfixWhitespace("List this project's messages with `memory_message_list`")) ||
    strings.Contains(strings.ToLower(optionalSources), "without consuming them") {
    t.Error("Optional sources still lists messages only to leave them pending")
}
```

Evidence: the positive phrase is whitespace-normalized like the other
`Optional sources` phrases, the forbidden phrase is matched with
`strings.ToLower` so any letter case (including the mutation's leading
"List this project's messages without consuming them") is caught, and the
diagnostic is exactly `Optional sources still lists messages only to leave
them pending`. Every other `Optional sources` phrase check is unchanged, and
no other assertion or diagnostic in the test changed.

### Requirement 2 — `--status` diagnostic renamed

The phrase checked is unchanged (`accepts no handoff, consumes no ai-memory
message and arms nothing`); only the failure message changed from
`--status omits its handoff refusal` to the exact
`--status section does not say it accepts no handoff`.

### Requirement 3 — both layouts, same checks

The change sits inside the existing
`for name, source := range roundfixLayoutSources()` loop, so `canonical`
(`../.agents/skills`) and `embedded` (the `go:embed`ed `skills/roundfix`)
run the identical corrected checks. No other loop, assertion or diagnostic
was modified.

### Requirement 4 — both negative controls proven in temporary `git archive` copies (outside the checkout)

Both controls used a fresh `mktemp -d` directory under
`/var/folders/.../T/opencode/`, populated with `git archive HEAD | tar -x`,
with this Task's working-tree `skills/roundfix_modes_test.go` overlaid (the
audited head carries the corrected test), and the mutation applied to
`skills/roundfix/references/modes.md`. Both copies were removed afterwards;
nothing was written under the checkout for these proofs.

- F-01: replaced the `Optional sources` sentence
  `List this project's messages with (memory_message_list)` with
  `List this project's messages without consuming them`. Running
  `go test -timeout 10m -count=1 -run '^TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage$' ./skills`
  in the copy produced:

  ```text
  --- FAIL: TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage (0.00s)
      --- FAIL: TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage/embedded (0.00s)
          roundfix_modes_test.go:155: Optional sources still lists messages only to leave them pending
  ```

- F-02: removed `accepts no handoff, ` from `## --status`. The same command
  in the copy produced:

  ```text
  --- FAIL: TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage (0.00s)
      --- FAIL: TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage/embedded (0.00s)
          roundfix_modes_test.go:159: --status section does not say it accepts no handoff
  ```

Both fail on the `embedded` layout with the exact declared message.

### Acceptance Criteria

- Delivered skill text passes in both layouts:
  `go test -timeout 10m -count=1 -run '^TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage$' ./skills`
  → `ok  roundfix/skills`.
- F-01 mutation fails `embedded` with `Optional sources still lists messages
  only to leave them pending` (evidence above).
- F-02 mutation fails `embedded` with `--status section does not say it
  accepts no handoff` (evidence above).

### Focused checks run in the checkout

- `go test -timeout 10m -count=1 -run '^TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage$' ./skills`
  → `ok  roundfix/skills  0.292s`.
- `go test -timeout 15m -count=1 ./skills` → `ok  roundfix/skills  1.139s`.

Not run here (left to the Daemon): the task file's `## Verification`
command, which also runs the whole `skills` package and the two `grep -qF`
checks for the declared diagnostics; both declared diagnostics are now
present in `skills/roundfix_modes_test.go`.

### Changed files

- `skills/roundfix_modes_test.go` — test only.
- `docs/specs/0325-the-follow-up-mode-consumes-the-handoff-and-every-message/task_04.md` — this Result section.

## Carry-forward provenance

- Source Run: `run_20261027T025147Z_96b54c79ec0f7eed`
- Source commit: `f003fa82ce4b54ee729b9179f03a6959cd4c9971`
