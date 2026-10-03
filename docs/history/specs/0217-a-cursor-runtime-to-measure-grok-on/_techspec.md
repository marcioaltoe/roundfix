---
spec: 0217-a-cursor-runtime-to-measure-grok-on
prd: _prd.md
created: 2026-10-01
---

# A Cursor runtime to measure Grok on — Technical Spec

## Executive Summary

`cursor` joins `codex`, `claude` and `opencode` as an ACP Runtime that acpx
reaches through its built-in `cursor` agent, `cursor-agent acp`. Three changes
make it usable. One list of supported runtimes feeds every validation and
message. The capability projection reads a Cursor model value, brackets
included, as one model identity, and a Cursor selection names that value with
an empty reasoning effort (ADR-0217). The adapter check, which every proof,
Run and Doctor line already passes through, asks `cursor-agent status` whether
a login exists and refuses with `cursor_login_required` when it does not.
Documentation and the Roundfix Skill follow. A final live Task measures Grok
on two replayed Tasks against the current default.

The trade-off is fidelity to what Cursor states over uniformity with the
other runtimes. Reading Cursor's bracket as an effort would let a Cursor tuple
look like the others, but it would invent a meaning for `thinking`, `context`
and `fast`. The cost is a model value that is longer and runtime-specific in
configuration, and a projection mode used by one runtime.

## Project Constraints

- Identifier strategy: not applicable — `cursor` is a runtime name in the
  existing Agent Selection vocabulary, and a Cursor model is the advertised
  value. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — no Roundfix HTTP request. Roundfix
  runs `cursor-agent status` with the Run's explicit environment and reads
  only whether it reports a login. It never runs `login` or `logout`, never
  passes `--api-key`, `--auth-token` or `-H`, never reads `CURSOR_API_KEY` or
  `CURSOR_AUTH_TOKEN`, and never prints, logs or stores the account the status
  names. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0217 binds the model-value rule
  and the login rule; ADR-0017, ADR-0037, ADR-0049, ADR-0050, ADR-0105,
  ADR-0107, ADR-0125, ADR-0187, ADR-0189, ADR-0198 and ADR-0211 bind the
  runtime, the proof, the fixtures and the skill as the PRD records; ADR-0108,
  ADR-0180 and ADR-0181 stay unchanged; ADR-0208 explains the split from Spec
  0218. The gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104,
  ADR-0155, ADR-0156 and ADR-0167, and ADR-0093, ADR-0117, ADR-0168, ADR-0176
  and ADR-0183 check consistency. ADR-0166, ADR-0178 and ADR-0182 govern each
  Task commit. ADR-0114 and ADR-0184 apply, and ADR-0069, ADR-0096, ADR-0097,
  ADR-0192, ADR-0194, ADR-0195, ADR-0199, ADR-0209 and ADR-0210 do not, as the PRD records. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 ("considere autorizado a ajustar todas as skills se
  necessário") and the cycle authorization of 2026-10-01 ("Tudo, de A a F"),
  recorded in [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runtime.md`. Sanctioned regeneration:
  `make skills-sync`. The Makefile, CI workflows, lint configuration,
  `.roundfixrc.yml`, `internal/cli/cli_test.go` and `go.mod` stay untouched.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

| Concern | Owner | File |
| --- | --- | --- |
| The supported runtime list | `SupportedRuntimes`, `isSupportedAgent` | `internal/config/config.go` |
| The Cursor selection rule | `normalizeSelection` | `internal/config/profiles.go` |
| The runtime spec and adapter command | `RuntimeFor`, `defaultAdapterCommands`, `resolveAdapterInvocation` | `internal/agent/agent.go`, `internal/agent/acpx_runner.go` |
| Whole-value model parsing | `SelectionRetention.WholeModelValues`, `RetentionFor`, `parseModelCapability` | `internal/agent/selection_capabilities.go` |
| The login check | `checkCursorLogin`, `CursorLoginRequiredError` | new `internal/agent/cursor_login.go`, called from `checkAdapter` in `internal/agent/acpx_runner.go` |
| Help and validation text | `validateAgent`, `displayAgent`, implement and review help | `internal/cli/cli.go`, `internal/cli/implement.go` |
| The Doctor next action | `runtimeHealthChecker.Adapter` | `internal/cli/health.go` |
| Documentation | the guides and the Roundfix Skill | `docs/user-guide/configuration.md`, `docs/user-guide/commands/doctor.md`, `docs/user-guide/commands/profiles.md`, `.agents/skills/roundfix/references/runtime.md` |
| The measurement | the protocol below | `internal/agent/testdata/cursor-session-recorded.json` (new), `internal/agent/cursor_recorded_session_test.go` (new), `docs/specs/0217-a-cursor-runtime-to-measure-grok-on/measurement/grok-through-cursor.md` (new) |

The Daemon, the Fallback Chain and the token record need no change: a Cursor
selection is one more tuple, refused before work like any unproved tuple.

## Implementation Design

### Interfaces

```go
// internal/config/config.go
// SupportedRuntimes is the one ordered list of ACP Runtime names. Every
// validation and every message that lists runtimes reads it.
func SupportedRuntimes() []string // codex, claude, cursor, opencode

// internal/agent/selection_capabilities.go
type SelectionRetention struct {
	Model            string
	ReasoningEffort  string
	WholeModelValues bool // true only for the cursor runtime
}

// internal/agent/cursor_login.go
const CursorLoginRequired = "cursor_login_required"

type CursorLoginRequiredError struct{ Command string }

func (CursorLoginRequiredError) Error() string          // "cursor-agent is not logged in"
func (CursorLoginRequiredError) Classification() string // CursorLoginRequired
func (CursorLoginRequiredError) NextAction() string     // "run cursor-agent login in a terminal yourself"

// checkCursorLogin runs "<executable> status" with environment, bounded by
// ctx and 10 seconds, and returns CursorLoginRequiredError when the command
// fails or its output contains "Not logged in". It returns no account data.
func checkCursorLogin(ctx context.Context, executable string, environment []string) error
```

### The runtime list

`SupportedRuntimes` returns `codex`, `claude`, `cursor`, `opencode`.
`isSupportedAgent` in `internal/config`, `validateAgent` in `internal/cli`, and
the refusals in `RuntimeFor` and `resolveAdapterInvocation` read it, so each
message ends `supported values: codex, claude, cursor, opencode`. The
`--agent` help of `implement` and `review` reads
`Agent runtime. Supported: codex, claude, cursor, opencode`. `displayAgent`
shows `Cursor`. The `runtimes:` section and `defaults.agent`, which are legacy,
keep their three values in one named legacy list, and refuse `cursor` with
`defaults.agent "cursor" is invalid; supported values: codex, claude,
opencode; name cursor in profiles`, composed from that list, so no source
file spells a runtime list by hand.

`RuntimeFor` returns `{ID: "cursor", DisplayName: "Cursor", Protocol: ACP}`
with no full-access mode. `defaultAdapterCommands["cursor"]` is
`cursor-agent acp`, the argv acpx's registry uses, so an acpx `agents.cursor`
override still wins as it does for every runtime. Setup's overrides and the
legacy model catalog are unchanged.

### The Cursor selection rule

`normalizeSelection` refuses a `cursor` selection whose trimmed
`reasoning_effort` is not empty with
`<path>.reasoning_effort must be "" for runtime cursor; Cursor states effort,
thinking, context and speed inside the model value, for example
"grok-4-20[thinking=true]"`. The model is kept verbatim after trimming.
`acpxReasoningEffortConfigKey` returns an error for `cursor`, which is
unreachable after validation.

### Whole-value model parsing

`RetentionFor` sets `WholeModelValues` when the runtime ID, without
`-custom`, is `cursor`. With it set, `parseModelCapability` returns
`ModelCapability{AdapterValue: v, CanonicalModel: v, ModelManaged: true}` for
every bounded value, with or without brackets, `default[]` included, and
dedups on the full value. Without it, parsing is today's, byte for byte.
`modelsForCanonical` then binds a Cursor selection only to the identical
value, the assignment is `model_managed` with `AdapterModel` equal to the
model, and proof sets it through `session/set_config_option` and requires the
complete returned state to report it as `currentValue`, as for every runtime.

The fixture `internal/agent/testdata/cursor-session-new-published.json` holds
the `configOptions` of the published forum response: a `mode` option and a
`model` option of category `model` with the values listed there, including
`default[]`, `grok-4-20[thinking=true]` and
`gpt-5.4[reasoning=medium,context=272k,fast=false]`, current value
`default[]`. It carries no account data.

### The login check

`checkAdapter` resolves the invocation and checks the executable as today.
For runtime `cursor` it then calls `checkCursorLogin` with the executable and
the explicit environment, after the lookup and before any lineage contract,
because Cursor has no npm package to inspect. Its error reaches every caller
unchanged:

- `ProveExactSelection` returns it, so `profiles validate`, `profiles
  configure` and a Run's preflight refuse the tuple before any session opens,
  and a Run's Fallback Chain moves on after notification (ADR-0050).
- The Doctor `adapter:` line reports
  `cursor: cursor-agent is not logged in; classification: cursor_login_required`
  through the existing `Classification()` reading, and the runtime health
  checker in `internal/cli/health.go` takes the next action from the error's
  `NextAction()`, as it takes an install command from `InstallCommand()`.

The status output is read in memory and dropped. No Roundfix string names
`login` as an argument for `cursor-agent`, and no production file contains
`CURSOR_API_KEY`, `CURSOR_AUTH_TOKEN`, `--api-key` or `--auth-token`.

### The measurement protocol

The measurement Task runs only on the maintainer's machine, with
`NODE_OPTIONS` unset, the rebuilt `bin/roundfix`, and the maintainer's
existing Cursor login. It never logs in. Its steps:

1. `cursor-agent status`. Without a login, record `cursor_login_required`
   in the Task's result and stop: the Task fails, and no file below is
   written.
2. Through an ACP client that sends only `initialize`, `authenticate` with
   `cursor_login` when acpx does not, and `session/new` in an empty scratch
   directory, capture the `configOptions` of one real session. Write them,
   with any account field removed, as
   `internal/agent/testdata/cursor-session-recorded.json`, and record whether
   acpx alone could open the session.
3. Pick the first advertised value whose name starts with `grok`, preferring
   one whose brackets state `fast=false` when such a variant exists. Prove it
   with `bin/roundfix profiles validate --category docs` in a scratch clone of
   this repository whose `.roundfixrc.yml` names it as the `docs` preferred
   selection with the current `docs` default as the fallback.
4. Replay two archived Tasks, each on the Grok selection and on the current
   built-in `docs` or `chore` default, in a fresh scratch clone per replay:
   `0210-evidence-snapshots-that-stay-small/task_02` (docs) and
   `0195-owned-skills-and-a-release-step-that-follow-the-bundle/task_06`
   (chore); reserves, in order, `0194-a-skill-and-a-command-guide-read-one-command-at-a-time/task_04`
   and `0202-a-qa-gate-that-reruns-only-stale-rows/task_04`. A replay starts
   from the Spec's squash merge commit with the Task's declared files restored
   to that commit's first parent, and the Spec directory restored under
   `docs/specs/` with that Task `pending`, every other Task removed from the
   graph, and no QA Task. A Task whose Verification passes on that state is
   replaced by the next reserve. The replay runs
   `bin/roundfix implement --spec <slug> --agent <runtime> --model <model>
   --reasoning-effort <effort>` with `--no-input`.
5. Record per replay: the selection, wall time, prompts and Verification
   repairs from `bin/roundfix runs show <run>`, the outcome, whether any
   permission prompt, `cursor/ask_question` or `cursor/create_plan` needed a
   person, and the tokens reported. Record whether `session/load` after
   `session/new` worked, and the billing mode the model value states.
6. Write `docs/specs/0217-a-cursor-runtime-to-measure-grok-on/measurement/grok-through-cursor.md`
   with a `## Measured` table of those rows, the advertised Grok values, and
   a `## Reading` section. When Grok matches the default on outcome with no
   person needed, the Task also mints a dated Backlog Entry under
   `docs/backlog/` proposing a Grok fallback, which the Daemon records as an
   undeclared path; otherwise the reading says why not. It changes no
   profile.

Scratch clones live under the system temporary directory, are clones of this
repository only, and are removed after the record is written. The replay Runs
are recorded in the live Run Database like any Run on this machine.

### Data Models

No Run Database or configuration schema changes. A Cursor tuple is an
`AgentSelection` with `runtime: cursor`, the advertised model value, and
`reasoning_effort: ""`. The recorded fixture is the `configOptions` array of
one ACP `session/new` result.

### API Contracts

1. API Contract: `roundfix implement --agent cursor --model <value>
   --reasoning-effort ""` and `roundfix review` with the same flags accept
   `cursor`; any other runtime name is refused with exit `2` and
   `unsupported Agent "<name>"; supported values: codex, claude, cursor,
   opencode`.
2. API Contract: Project and User Config accept `runtime: cursor` in
   `profiles.<category>.preferred` and `fallbacks[]`; a non-empty
   `reasoning_effort` there is a load error naming the rule.
3. API Contract: `roundfix doctor` reports a Cursor runtime on its
   `adapter:` line; a missing login fails the line with
   `classification: cursor_login_required` and
   `next: run cursor-agent login in a terminal yourself`.

### Surface Transcripts

The fixture for Transcript 1 is a repository whose Project Config names
`cursor grok-4-20[thinking=true] ""` as the preferred selection and
`cursor default[] ""` as the only fallback of every required category, and a
fake `cursor-agent` on `PATH` whose `status` prints `Not logged in`.

1. Surface Transcript: Doctor without a Cursor login.

   ```transcript
   $ roundfix doctor
   stdout:
   adapter: failed (cursor: cursor-agent is not logged in; classification: cursor_login_required; next: run cursor-agent login in a terminal yourself)
   stderr:
   exit: 1
   ```

2. Surface Transcript: an unknown runtime on the one-Run override, in an
   empty Git repository.

   ```transcript
   $ roundfix implement --spec 0300-example --agent gemini --model x --reasoning-effort ""
   stdout:
   stderr:
   Preflight failed

   Reason:
     unsupported Agent "gemini"; supported values: codex, claude, cursor, opencode

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   Usage:
     Run 'roundfix implement --help' for usage.
   exit: 2
   ```

3. Surface Transcript: the `--agent` help line.

   ```transcript
   $ roundfix implement --help
   stdout:
     --agent              Agent runtime. Supported: codex, claude, cursor, opencode
   stderr:
   exit: 0
   ```

Transcript 1 shows only the `adapter:` line; Transcript 3 shows only the
`--agent` line.

## Coverage Map

- Goal 1 → The runtime list; Whole-value model parsing; API Contracts 1 and 2.
- Goal 2 → The Cursor selection rule; Whole-value model parsing.
- Goal 3 → The login check; API Contract 3.
- Goal 4 → The measurement protocol.
- User Story 1 → The runtime list; The Cursor selection rule.
- User Story 2 → Whole-value model parsing; The login check.
- User Story 3 → The login check; Surface Transcript 1.
- User Story 4 → The measurement protocol.
- Core Feature 1 → The runtime list; Surface Transcripts 2 and 3.
- Core Feature 2 → The Cursor selection rule; Whole-value model parsing.
- Core Feature 3 → The login check; Surface Transcript 1.
- Core Feature 4 → Testing Approach 4.
- Core Feature 5 → The measurement protocol.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 1.
- Success Metric 5 → Testing Approach 5.

## Integration Points

- **acpx** `0.12.0` or newer, through its built-in `cursor` agent, which
  spawns `cursor-agent acp`. The installed `0.19.4` lists it.
- **`cursor-agent`**, the maintainer's install, at `~/.local/bin/cursor-agent`
  on the measured machine. Roundfix runs only its `status` subcommand
  directly; acpx runs `acp`.
- **Cursor's service** receives the prompts of a Cursor Agent Session under
  the maintainer's plan; only the measurement Task sends any.

## Testing Approach

Every fake `cursor-agent` is the compiled test binary re-executed through the
`internal/agent` fake adapter provisioning, never a written script
(ADR-0125). No test reaches Cursor or the network.

1. **The runtime list.** `TestSupportedRuntimesIsTheOneList` asserts the
   list and that each refusal message (config preferred and fallback runtime,
   `RuntimeFor`, `resolveAdapterInvocation`, `validateAgent`) is built from it.
   `TestCursorProfileLoadsWithAnEmptyEffort` and
   `TestCursorProfileRefusesAReasoningEffort` load fragments.
   `TestLegacyRuntimesRefuseCursor` keeps the legacy section at three values.
   A built binary's `implement --help` names the four runtimes.
2. **Parsing and assignment.**
   `TestCursorCapabilitiesReadWholeModelValues` parses the published fixture
   under `cursor` with no issue and every value as its own model-managed
   identity. `TestCursorSelectionAssignsTheExactValue` assigns
   `grok-4-20[thinking=true]` as `model_managed` and refuses `grok-4-20` as
   not advertised. `TestOtherRuntimesKeepBracketParsing` parses the same
   fixture under `claude` and gets today's `malformed_model_value`.
3. **The login.** `TestCursorAdapterRefusedWithoutLogin` (fake `status`
   prints `Not logged in`, and a second fake exits `1`) returns
   `CursorLoginRequiredError`. `TestCursorAdapterReadyWithLogin` passes and
   `TestCursorProofRefusesBeforeAnySessionWithoutLogin` records no acpx call.
   `TestDoctorAdapterNamesTheCursorLogin` prints Transcript 1's line. A sweep
   finds no `CURSOR_API_KEY`, `CURSOR_AUTH_TOKEN`, `--api-key`,
   `--auth-token` or a `login` argument in production Go files.
4. **Documentation.** Phrase checks over the configuration guide (task_01),
   the Doctor and profiles guides (task_02) and the skill reference
   (task_03), the skill's recorded version, and `make skills-sync-check`.
5. **The measurement.** `TestCursorRecordedSessionAdvertisesGrok` parses
   `cursor-session-recorded.json` under `cursor`, asserts at least one value
   starting with `grok`, and assigns it. The measurement record exists with
   its `## Measured` table.

## Build Order

1. The runtime list, the selection rule, whole-value parsing and the help
   text, task_01 (depends on: none).
2. The login check in the adapter check and the Doctor line, task_02
   (depends on: 1, because both change `internal/agent/acpx_runner.go`).
3. The Roundfix Skill, task_03 (depends on: 1 and 2, which it describes;
   tasks 01 and 02 write the guides beside their CLI changes).
4. The measurement, task_04 (depends on: 1 and 2, which it runs; not on 3).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **The login blocks the measurement.** On 2026-10-01 `cursor-agent` was not
  logged in. Task 04 then fails with `cursor_login_required` and the Spec
  waits for the maintainer; tasks 01 to 03 are unaffected.
- **acpx and Cursor's authentication.** Cursor's documentation lists an
  `authenticate` call before `session/new`. If acpx does not make it, a
  logged-in session may still be refused; the measurement records which.
- **Blocking extension methods.** A Cursor session may send
  `cursor/ask_question` or `cursor/create_plan`, which acpx may not answer,
  and a Daemon Run could wait on it until its budget. The measurement records
  whether either was sent; a Run is still bounded by its Run Budget.
- **A `status` output change.** The check reads only `Not logged in` and the
  exit status, so a reworded message could read as logged in; the session
  would then be refused by Cursor at `session/new`, which proof also catches.
- **Grok's published standing.** Grok 4.7 scored 25.8 on Terminal-Bench 4.0
  against 59.6 for Opus 5.5, with about three times the output tokens; the
  measurement is a small sample and must not be read as more.

## Decisions

- Whole-value parsing for `cursor` only, chosen by runtime, not by detecting
  bracket shapes, so another runtime's encoding never changes silently.
- The login check lives in `checkAdapter`, the one point proof, Runs,
  sealed sessions and Doctor already share, instead of a new Doctor-only line.
- The measurement replays archived Tasks in scratch clones rather than
  choosing live work, so the comparison is on identical inputs and spends
  nothing on this repository's queue.
- No `grok-build` support: it needs an xAI key Roundfix would hold, and no
  `grok` binary is installed.
