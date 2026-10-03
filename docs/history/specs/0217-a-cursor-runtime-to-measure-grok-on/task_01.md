---
task: task_01
spec: 0217-a-cursor-runtime-to-measure-grok-on
status: completed
type: backend
complexity: high
---

# Task 01: `cursor` is an ACP Runtime whose selection names the advertised value

## Overview

Roundfix accepts three ACP Runtimes, and each message that lists them spells
the list out by hand. Cursor advertises each model as one value whose
brackets carry its parameters, and the capability projection reads that
bracket as a reasoning effort and refuses `default[]`. This Task adds `cursor`
through one list of runtimes, reaches it through acpx's built-in `cursor`
agent, requires a Cursor selection to name the advertised value with an empty
effort (ADR-0217), and parses Cursor values as whole model identities.

## Requirements

1. MUST characterize, before changing them, today's refusal of
   `runtime: cursor` in a profile and today's `malformed_model_value` for the
   published Cursor fixture, in the new tests below, and then change them as
   the TechSpec "The runtime list" and "Whole-value model parsing" state.
2. MUST add `SupportedRuntimes()` in `internal/config` returning `codex`,
   `claude`, `cursor`, `opencode`, and MUST make `isSupportedAgent`, the
   profile runtime refusal, `validateAgent` in `internal/cli`, and the
   refusals of `RuntimeFor` and `resolveAdapterInvocation` read it. No
   production Go file may keep a runtime list spelled
   `values: codex, claude, opencode` or `Supported: codex, claude, opencode`.
   The skill install targets of `roundfix skills`, which also name `codex`,
   `claude` and `opencode`, are not runtimes and stay unchanged.
3. MUST keep the legacy `runtimes:` section and `defaults.agent` at their
   three values, read from one named legacy list rather than a literal, and
   refuse `cursor` there with `defaults.agent "cursor" is invalid; supported
   values: codex, claude, opencode; name cursor in profiles`, composed from
   that list.
4. MUST make `RuntimeFor` return `cursor` with display name `Cursor`, ACP
   protocol and no full-access mode, and MUST map `cursor` to
   `cursor-agent acp` in `defaultAdapterCommands`.
5. MUST refuse a `cursor` selection whose trimmed `reasoning_effort` is not
   empty, with the message of the TechSpec "The Cursor selection rule", and
   keep its model verbatim after trimming.
6. MUST add `WholeModelValues` to `SelectionRetention`, set by `RetentionFor`
   only for `cursor` (with or without `-custom`), and make
   `parseModelCapability` read every bounded value as one model-managed
   identity under it, `default[]` included. Parsing for every other runtime
   MUST stay byte-for-byte today's.
7. MUST add the fixture `internal/agent/testdata/cursor-session-new-published.json`
   holding the `configOptions` the TechSpec describes, with no account data.
8. MUST change the `--agent` help line of `implement`, `resolve` and `watch`
   to `Agent runtime. Supported: codex, claude, cursor, opencode`, and
   `displayAgent` to show `Cursor`.
9. MUST update the pinned refusal in `internal/config/config_test.go` to the
   four-runtime message, the declared break of the PRD.
10. MUST add to `docs/user-guide/configuration.md`, under
    `## Agent selection profiles`, a `### Cursor runtime` subsection that
    states: `cursor` is opt-in and reached as `cursor-agent acp` through
    acpx; a Cursor selection names the advertised model value verbatim, with
    the example
    `{runtime: cursor, model: "grok-4-20[thinking=true]", reasoning_effort: ""}`;
    a non-empty effort is refused; the value states speed, so `fast=false` is
    how a maintainer avoids Fast billing when Cursor offers that variant; and
    `cursor` is not in the legacy `runtimes:` section.
11. MUST NOT change the Recommended Profile, the built-in profiles, the model
    picker catalog, Setup's acpx overrides or `internal/cli/cli_test.go`.

## Subtasks

- [ ] Characterize today's refusals, then add the one runtime list.
- [ ] Add the runtime spec, the adapter command and the selection rule.
- [ ] Parse Cursor model values whole, with the published fixture.
- [ ] Update help, display, the pinned message and the configuration guide.

## Acceptance Criteria

- [ ] A profile naming `cursor` and `grok-4-20[thinking=true]` with an empty
      effort loads; one with `reasoning_effort: high` is refused with the rule.
- [ ] The published fixture projects under `cursor` with no issue, and the
      same fixture under `claude` still yields `malformed_model_value`.
- [ ] A built binary's `implement --help` lists four runtimes.

## Context

- creates: `internal/config/cursor_runtime_test.go`
- creates: `internal/agent/cursor_capabilities_test.go`
- creates: `internal/agent/testdata/cursor-session-new-published.json`
- interface: `internal/config/config.go`
- interface: `internal/config/profiles.go`
- interface: `internal/config/config_test.go`
- interface: `internal/agent/agent.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/selection_capabilities.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/implement.go`
- instruction: `docs/adr/0217-a-cursor-selection-names-the-advertised-model-value-and-the-login-stays-the-maintainers.md`
- instruction: `internal/agent/selection_assignment.go`
- interface: `docs/user-guide/configuration.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestSupportedRuntimesIsTheOneList|TestCursorProfileLoadsWithAnEmptyEffort|TestCursorProfileRefusesAReasoningEffort|TestLegacyRuntimesRefuseCursor|TestCursorCapabilitiesReadWholeModelValues|TestCursorSelectionAssignsTheExactValue|TestOtherRuntimesKeepBracketParsing|TestCursorRuntimeSpecUsesCursorAgentACP)$' ./internal/config ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSupportedRuntimesIsTheOneList TestCursorProfileLoadsWithAnEmptyEffort TestCursorProfileRefusesAReasoningEffort TestLegacyRuntimesRefuseCursor TestCursorCapabilitiesReadWholeModelValues TestCursorSelectionAssignsTheExactValue TestOtherRuntimesKeepBracketParsing TestCursorRuntimeSpecUsesCursorAgentACP; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 ./internal/config || exit 1; for phrase in '### Cursor runtime' 'cursor-agent acp' 'grok-4-20[thinking=true]' 'fast=false'; do tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "$phrase" || { printf 'missing phrase in configuration guide: %s\n' "$phrase" >&2; exit 1; }; done; if grep -rnE --include='*.go' '(values|Supported): codex, claude, opencode' internal cmd | grep -v '_test.go:'; then printf 'a production file still lists three runtimes\n' >&2; exit 1; fi` — expected: exit 0; before this Task the named tests do not exist, so the command fails.
- `tmp="$(mktemp -d)" || exit 1; go build -buildvcs=false -o "$tmp/roundfix" ./cmd/roundfix || exit 1; for command in implement resolve watch; do "$tmp/roundfix" "$command" --help > "$tmp/$command.txt" 2>&1; grep -qF 'Agent runtime. Supported: codex, claude, cursor, opencode' "$tmp/$command.txt" || { printf 'help of %s does not list cursor\n' "$command" >&2; exit 1; }; done` — expected: exit 0; before this Task the help lists three runtimes, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1-2; User Stories 1-2; Core Features 1-2; Success Metrics 1, 2 and 4; Declared breaks
- [_techspec.md](_techspec.md) — Interfaces; The runtime list; The Cursor selection rule; Whole-value model parsing; Surface Transcripts 2 and 3; API Contracts 1-2; Testing Approach 1-2; Build Order 1
- ADR-0217; ADR-0037; ADR-0049; ADR-0105

## Result

Implemented this Task's runtime, selection, projection, and help slice for
Daemon Verification. Task status remains Daemon-owned. No declared Verification
command was run, and no commit, push, or Pull Request was created.

### Implementation

- `config.SupportedRuntimes()` returns a caller-owned ordered list of `codex`,
  `claude`, `cursor`, and `opencode`. Profile validation, CLI validation,
  Runtime and adapter refusals, and the three command help lines read it.
  A separate named legacy list governs `runtimes:` and `defaults.agent`;
  Cursor's legacy default refusal directs the maintainer to profiles.
- Cursor uses ACP, display name `Cursor`, and the default adapter argv
  `cursor-agent acp`, with no full-access mode. Existing configured adapter
  precedence and reasoning-control refusal remain in place.
- Cursor selections preserve the trimmed advertised model value and refuse
  non-empty trimmed effort with the authored rule, in preferred and fallback
  positions. Cursor and `cursor-custom` retention reads bounded values as
  whole model-managed identities and deduplicates by the full value. The
  existing parser branch for other runtimes is unchanged.
- Added the published model option fixture, including all 25 model values and
  the mode option, with model current value `default[]` as specified. Source:
  [Cursor forum response, post 9](https://forum.cursor.com/t/157312/9).
  Only public option fields are retained; there is no account or session data.
  The first HTTP fetch was sandbox-blocked; the web reader supplied the public
  response. No Cursor service or local account was contacted.
- Documented Cursor under Agent selection profiles, including opt-in use,
  exact value selection, empty effort, `fast=false`, and the legacy exclusion.
  Updated the intentionally changed pinned profile-runtime refusal.
- Moved `TestModelCatalogRetainsOfficialCodexIdentifiers` from config's internal
  test file into the new external `config_test` file, preserving every catalog
  assertion and its `internal/config` package ownership. This avoids a test
  import cycle now that agent reads config's runtime list. The model catalog
  and recommended/built-in profiles are untouched.

### Characterization before production changes

New tests first asserted today's Cursor profile refusal and today's published
fixture `malformed_model_value` refusal under Cursor. With an isolated writable
Go cache, `go test ./internal/config ./internal/agent -run Characterization
-count=1` exited 0 before production edits. Those tests then became the authored
Cursor regression tests; the other-runtime test retains the old parser refusal.

The initial attempt using the host Go cache was refused by the sandbox. All
subsequent Go checks used `GOCACHE=/private/tmp/roundfix-0217-task01-cache`.

### Focused checks and acceptance evidence

- Profile acceptance and effort refusal: the final focused command
  `GOCACHE=/private/tmp/roundfix-0217-task01-cache rtk proxy go test
  ./internal/config ./internal/agent -run
  'Cursor|SupportedRuntimes|LegacyRuntimes|OtherRuntimes|SelectionCapabilities|SelectionAssignment|ModelCatalogRetains|AgentSelectionProfileRejects'
  -count=1` exited 0 for both packages. `TestCursorProfileLoadsWithAnEmptyEffort`
  checks empty and whitespace-only effort, model trimming, and Cursor fallback
  loading. `TestCursorProfileRefusesAReasoningEffort` checks the exact rule
  at preferred and fallback paths.
- Published fixture acceptance and preserved Claude refusal: that same focused
  check includes `TestCursorCapabilitiesReadWholeModelValues` (Cursor and
  custom Cursor, all 25 whole identities), `TestCursorSelectionAssignsTheExactValue`
  (exact Grok value assigns as `model_managed`; the bare prefix is refused),
  and `TestOtherRuntimesKeepBracketParsing` (Claude, Codex, OpenCode and custom
  Claude still return only `malformed_model_value`).
- Built help acceptance: `GOCACHE=/private/tmp/roundfix-0217-task01-cache rtk
  proxy go build -o /private/tmp/roundfix-0217-task01 ./cmd/roundfix` exited 0.
  A Python subprocess smoke check ran that binary's help for `implement`,
  `resolve`, and `watch`: each exited 0 and its `--agent` line read
  `Agent runtime. Supported: codex, claude, cursor, opencode`. The same smoke
  check in an isolated empty Git repository exercised the implement override
  with an unknown runtime; it exited 2 and reported the four-runtime refusal.
- Additional config load regressions: `GOCACHE=/private/tmp/roundfix-0217-task01-cache
  rtk proxy go test ./internal/config -run
  'Load|AgentSelectionProfile|Legacy|SupportedRuntimes|Cursor' -count=1` exited 0.
- `GOCACHE=/private/tmp/roundfix-0217-task01-cache rtk proxy go vet
  ./internal/config ./internal/agent ./internal/cli` exited 0.
- Production source inspection found no old three-runtime diagnostic/help
  phrase. `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

### Handoff boundary

The only pre-existing changed path was this Task file's Daemon-owned transition
from `pending` to `in_progress`; it was preserved. No other Task or Task Graph
file was edited. Setup overrides, skill install targets, model picker catalog,
`internal/cli/cli_test.go`, and the Recommended Profile remain untouched.
Login checks, the Roundfix Skill update, and live measurement remain assigned
to their subsequent Tasks. Declared Verification and terminal settlement remain
with the Daemon; the focused evidence above is not their verdict.

### Verification Feedback repair — attempt 1

Inspected the Daemon diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261002T190644Z_6dc4ccb1092ddf45/verification/batch-001-attempt-1.log`.
The configured `make verify-changed` failed because the coverage guard keys
existing test names by package: moving the catalog test into `internal/agent`
removed its recorded identity from `internal/config`. The diagnostic's coverage
additions were informational; the missing config test was its sole regression.

Moved that test into `internal/config/cursor_runtime_test.go`, now an external
`config_test` package. The catalog assertions remain unchanged. Cursor tests
in that file now exercise exported config loading and resolution; supported
runtime acceptance is checked through actual profile loading. This preserves
the original coverage ownership while avoiding the config/agent test import
cycle. No coverage record or Verification configuration was changed.

Fresh focused evidence (each Go command used
`GOCACHE=/private/tmp/roundfix-0217-task01-cache` and `rtk proxy`):

- Before the repair, `go test ./internal/config -list
  '^TestModelCatalogRetainsOfficialCodexIdentifiers$'` did not list the test.
  After the repair, the same discovery check listed it under `internal/config`.
- `go test ./internal/config ./internal/agent -run
  'Cursor|SupportedRuntimes|LegacyRuntimes|OtherRuntimes|ModelCatalogRetains'
  -count=1` exited 0 for both packages.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1` exited 0,
  directly checking the reported coverage regression.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

Task status and declared Verification remain untouched. The configured
Verification sequence was not rerun; its next run and terminal settlement
remain with the Daemon.

## Carry-forward provenance

- Source Run: `run_20261002T190644Z_6dc4ccb1092ddf45`
- Source commit: `4612c792ab521f642977c05e69007c3014634d9b`
