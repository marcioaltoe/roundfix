---
task: task_01
spec: 0361-a-session-starts-with-an-update-check-and-the-update-mode-applies-it
status: completed
type: backend
complexity: medium
---

# Task 01: Project Config chooses notify or auto in update.mode

## Overview

User Config and Project Config gain a top-level `update:` mapping whose one
key, `mode`, selects `notify` (the built-in default) or `auto`, Project Config
winning, and any other value is refused at load naming the key. The
configuration guide documents the key, and the Behavior Surface Record is
re-recorded for the changed key set and guide.

## Requirements

1. MUST add to `internal/config/config.go` the `Update` section of
   `_techspec.md` → Interfaces (`Config.Update`, `updateOverlay` under the
   `update` YAML key in `configOverlay`), the built-in default `notify` in
   `Builtin()`, the overlay merge for User Config and Project Config with
   Project Config applied last, and the `Validate` refusal of API Contract 1,
   word for word: `update.mode must be notify or auto, got "<value>"`.
2. MUST NOT change `DefaultConfigYAML`, `DefaultProjectConfigYAML` or any other
   key's default, validation or message (`_techspec.md` → Invariant 5).
3. MUST add `TestUpdateModeIsNotifyOrAuto` in
   `internal/config/update_mode_test.go`, following
   `internal/config/delivery_review_concurrency_test.go` (temporary home and
   work directory, `gittest.MarkRepository`, `Load(LoadOptions{...})`), with
   the rows `_techspec.md` → Testing Approach names: `absent`, `project auto`,
   `user auto`, `project over user`, `invalid project`, `invalid user` and
   `templates unchanged` (both init templates contain no `update:` line).
4. MUST add, word for word, the `update.mode` row of API Contract 1 to the
   General settings table of `docs/user-guide/configuration.md`, after the
   `specs.root` row.
5. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last guide edit; never hand-edit it.

## Subtasks

- [ ] Add the `Update` section, its default, merge and validation.
- [ ] Add the configuration test.
- [ ] Document the key, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] An absent key reads `notify`; Project Config `auto` over User Config `notify` reads `auto`; User Config alone sets it.
- [ ] Any other value is refused naming `update.mode`, in either file.
- [ ] The init templates keep their bytes, and the configuration guide documents the key.

## Context

- instruction: `docs/adr/0371-a-session-starts-with-an-update-check-and-the-skill-applies-the-update.md`
- instruction: `internal/config/delivery_review_concurrency_test.go`
- interface: `internal/config/config.go`
- creates: `internal/config/update_mode_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `Config.Update` field and the `update.mode` refusal — consumers: `./internal/config`, `./internal/cli`
- The configuration guide and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestUpdateModeIsNotifyOrAuto$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestUpdateModeIsNotifyOrAuto ' || { echo 'missing pass: TestUpdateModeIsNotifyOrAuto' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/config && go build -buildvcs=false ./internal/cli && GOOS=windows go build -buildvcs=false ./internal/config` — expected: exit 0. Before this Task the test does not exist, so its pass line is missing; afterwards the key is read and refused as specified, the config suite passes, and the CLI still compiles against it.

- `grep -q '^| .update\.mode. | .notify. | Selects what a session does with the report of .roundfix update check.: .notify. reports a newer release' docs/user-guide/configuration.md || { echo 'missing update.mode row' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the configuration guide has no `update.mode` row; afterwards the row is present and the documentation contracts, the Behavior Surface fingerprint included, pass.

## References

- `_prd.md` → Goal 4; Core Feature 1; Success Metric 4
- `_techspec.md` → Current behavior 3; Interfaces; API Contract 1; Invariant 5; Declared breaks; Testing Approach; Build Order 1
- ADR-0371

## Result

Skills loaded: implement-task, coding-guidelines, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate.

Implemented only Task 01: `Config.Update`, its `notify` default, the `update`
YAML overlay and the shared validation used by `Validate` and each input file.
The existing load order applies User Config first and Project Config last.
Per-file validation retains the existing `parse config "<path>"` wrapper and
refuses an invalid User Config even when Project Config supplies a valid mode.
No other key's defaults, validators or messages changed.

Focused evidence:

- Before implementation, `go test -count=1 -timeout 2m ./internal/config -run
  '^TestUpdateModeIsNotifyOrAuto/'` exited 1: `Config` had no `Update` field.
- After implementation, `go test -count=1 -timeout 2m -v ./internal/config -run
  '^TestUpdateModeIsNotifyOrAuto/'` exited 0. The required rows `absent`,
  `project auto`, `user auto` and `project over user` passed, proving `notify`
  inheritance and `auto` from either file. The additional `project notify over
  user` row proves Project Config can also override User Config `auto`.
- The `invalid project` and `invalid user` rows passed for `always`, an empty
  string, `AUTO` and ` auto `. They assert the exact refusal from both
  `Validate` and `Load`, including the offending file's path. The invalid-user
  cases include a valid project override to prove invalid values are refused
  before another file can hide them.
- The `templates unchanged` row passed: neither init template has an `update:`
  line. A Python comparison against `git show HEAD:internal/config/config.go`
  confirmed byte-identical source for `DefaultConfigYAML`,
  `DefaultProjectConfigYAML`, `defaultConfigYAML` and `renderBuiltinProfilesYAML`.
  The new built-in field is not consumed by template rendering.
- A Python comparison with API Contract 1 confirmed the new General settings
  row is word-for-word identical and immediately follows `specs.root`.

The last implementation step after the final guide edit was the sanctioned
recorder: `go test -count=1 -tags docscontract ./internal/docscontract -run
'^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` exited 0
(`ok roundfix/internal/docscontract`). It regenerated
`docs/references/behavior-surfaces.json`; the record was never hand-edited.

Declared Task Verification and repository Verification remain for the Daemon;
no Task status, Task Graph or other Task file was edited. No follow-up work
was identified in this slice.
