---
spec: 0281-the-codex-full-access-mode-is-the-one-the-adapter-offers
prd: _prd.md
created: 2026-10-19
---

# The Codex full-access mode is the one the adapter offers — Technical Spec

## Executive Summary

One slice changes code and one changes documentation; no package is new.

1. Full-access mode discovery (ADR-0293), in `internal/agent`.
   - `RuntimeSpec` gains an ordered list of full-access candidates.
   - The capability projection keeps the session's `mode` config option.
   - `applyFullAccess` sends the first advertised candidate, sends the Codex
     sandbox preset when the session advertises `sandbox_mode`, and reports a
     degraded effective access policy when no candidate is advertised.
   - The Roundfix `profiles` reference describes it.
2. Documentation.
   - ADR-0011 gets an amendment note.
   - The configuration guide and the glossary describe the selection.
   - The `.roundfixrc.yml` comment loses its codex-acp 1.1.9 claim.
   - The `qa-gate` skill gains its sandboxed browser guidance.

The trade-offs accepted:

- A session that offers no full-access mode keeps running with its default
  access, reported degraded, instead of failing.
- A Codex session that advertises no `sandbox_mode` option receives no preset,
  on the adapter's word that its mode carries the sandbox.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, numbering or
  record schema is minted; one effective access policy text and one Run event
  status are added. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the selection reads the `acpx sessions
  show` output Roundfix already reads, and every test uses the fake acpx
  runner. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0293 governs Build Orders 1 and 2:
  "A full-access request selects the mode the Agent Session advertises".
  ADR-0011, amended by it: "ACP Runtimes can run Batches with explicit full
  access". ADR-0017: "Roundfix drives ACP Runtimes through acpx". ADR-0189:
  "An owned skill's version names its content, and its minimum is the
  bundle". ADR-0285: "Every test file compiles for Windows, and a repository
  contract proves it". ADR-0125: "A spawned fixture is a compiled binary,
  never a written script". ADR-0104: "A Spec accepts on evidence it did not
  author". Set no obligation here, because this Spec changes none of their rules: ADR-0233 (a skill version raise is regenerated at merge, and a review-only correction returns to review); ADR-0240 (one QA partial policy, and rows a Run sandbox cannot reach); ADR-0250 (a module version is chosen when it is recorded, and the coverage record lists every platform); ADR-0256 (a release waits for the skills that describe a changed surface); ADR-0276 (a row the authorization keeps offline qualifies, and a byte invariant compares with the Delivery Base); ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0272 (a review record survives a fix, and a correction inside the archived Spec's declared paths returns to review); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds); ADR-0237 (a Delivery Retry records a merge made outside the queue); ADR-0279 (status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor); ADR-0286 (a fixed round-2 finding gets one Correction Review of its fix). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AJ grant of 2026-10-19 ("Mexer em Governed Paths declarados") and the
  2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0281-the-codex-full-access-mode-is-the-one-the-adapter-offers/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/profiles.md`, `.roundfixrc.yml`,
  `internal/cli/cli_test.go`, `skills/qa-gate/SKILL.md`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-19 at `f382480b` (v0.91.0).

- Runtime capability (`internal/agent/agent.go`).
  - `RuntimeFor` (`:360`) sets `SupportedFullAccessMode` to `full-access` for
    `codex` (`:391`) and `bypassPermissions` for `claude` (`:393`). A command
    override, OpenCode and Cursor get none.
  - With `EnableFullAccess`, it sets `RequestedAccessPolicy` to `full-access`
    and `FullAccessMode` to that name (`:396-399`).
  - `AccessModeFor` (`:45`) returns `SupportedFullAccessMode`, else the legacy
    `FullAccessMode`. `ValidateRequestedAccessPolicy` (`:75`) refuses a
    runtime without one with `access_policy_unsupported`.
- Session setup (`internal/agent/acpx_runner.go`, `selection_assignment.go`).
  - A Run's `ensureSession` (`acpx_runner.go:1261`) ensures the session,
    reads its capabilities (`startSessionSelection`), applies the selection,
    and calls `applyFullAccess` (`:1306`), discarding the policy it returns.
  - The disposable proof `ProveExactSelection` (`selection_assignment.go:99`)
    calls `applyFullAccess` (`:151`) and stores the result as
    `SelectionProof.EffectiveAccessPolicy` (`:157`).
  - `applyFullAccess` (`acpx_runner.go:1364`) always sends `set-mode <name>`.
    A failure returns an `AccessPolicyError` of kind `access_policy_rejected`
    (`:1376-1386`), which fails the session.
  - It sends `set sandbox_mode danger-full-access` only when the runtime is
    `codex` and the mode is literally `full-access` (`:1388`). A failure that
    `isCodexSandboxUnavailable` (`:2253`) recognizes publishes the status
    `codex_sandbox_full_access_unavailable` and returns `full-access
    (degraded: danger-full-access sandbox preset unavailable)` (`:51`,
    `:1395-1401`).
- Capability evidence (`internal/agent/selection_capabilities.go`).
  - `ParseSessionCapabilitySnapshot` (`:313`) reads the `acpx.session.v1`
    record's `config_options` and `ParseSessionConfigOptions` (`:181`)
    projects every select option into `SelectionCapabilities.Options`
    (`SelectCapability`: identifier, current value, values, advertised count).
    The option's category is read (`acpSelectOptionPayload.Category`) but not
    kept.
- Reporting (`internal/cli`).
  - `profiles validate` appends `(effective access policy: <policy>)` to a
    proof whose policy holds `(degraded:` (`profiles_validate.go:421`), and
    `doctor` adds it to its profiles line (`doctor.go:552`).
- Tests that pin today's behavior (`internal/agent/acpx_runner_test.go`):
  `TestDisposableProofAppliesRequestedAccessMode` (`:1165`),
  `TestProofRecordsADegradedFullAccess` (`:1227`),
  `TestACPXRunAppliesFullAccessSessionSetup` (`:2866`),
  `TestACPXRunFailsSetupWhenFullAccessModeFails` (`:2925`) and
  `TestACPXRunWarnsWhenCodexSandboxPresetUnavailable` (`:2951`). The fake's
  `sessions show` (`fakeSelectionStateOutput`, `:5339`) and the
  `sessionCapabilitySnapshotFixture` helper (`:4925`) advertise only `model`
  and the reasoning option. `internal/agent/agent_test.go:107` and `:123`
  pin the Codex name `full-access`, and
  `TestRunResolveAgentFullAccessIsExplicitOptIn`
  (`internal/cli/cli_test.go:6855`) pins it as `wantAccess`.
- The adapter today. A local acpx 0.19.6 record of a codex-acp session
  advertises `mode` (category `mode`) with `read-only`, `workspace-write`,
  `agent`, `agent-full-access`, and no `sandbox_mode`. A claude-agent-acp
  session advertises `mode` with `default`, `acceptEdits`, `plan`, `auto`,
  `bypassPermissions`.
- Documentation.
  - ADR-0011 names `full-access` as the Codex mode.
  - The configuration guide's "Access policy readiness"
    (`docs/user-guide/configuration.md:470-480`) says readiness refuses when
    a runtime has no mode or the adapter refuses it.
  - `.roundfixrc.yml:51-55` says a true value fails every Task, "Measured on
    codex-acp 1.1.9, 2026-08-04".
  - The `qa-gate` skill tells a sandboxed step to mark a permission failure
    `blocked (environment: <error>)` (`.agents/skills/qa-gate/SKILL.md:617`)
    and names no browser driver.

## System Architecture

- `internal/agent/agent.go`: `RuntimeSpec.FullAccessModes`, the candidate
  lists, and the candidate choice.
- `internal/agent/selection_capabilities.go`: `SelectionCapabilities.Modes`,
  the session's mode option.
- `internal/agent/acpx_runner.go`: `applyFullAccess` takes the capabilities,
  and `ensureSession` passes them.
- `internal/agent/selection_assignment.go`: `ProveExactSelection` passes them.
- Documentation as listed in Build Order 2.

## Implementation Design

### Interfaces

```go
// internal/agent
type RuntimeSpec struct {
	// ...existing fields...
	FullAccessModes []string // full-access candidates, preference order
}
type SelectionCapabilities struct {
	// ...existing fields...
	Modes *SelectCapability // the session's mode option; nil when none
}
const AccessModeNotAdvertisedStatus = "full_access_mode_not_advertised"
func chooseFullAccessMode(candidates, advertised []string) (string, bool)
func (runner *ACPXRunner) applyFullAccess(ctx context.Context, req ExecuteRequest,
	capabilities SelectionCapabilities, sink runevent.Sink, codexEnv []string) (AccessPolicy, error)
```

### Data Models

No file format changes. `SelectionProof.EffectiveAccessPolicy` keeps its
type and gains one degraded value.

### API Contracts

1. API Contract: the full-access candidates.
   - `RuntimeFor` sets `FullAccessModes` to `agent-full-access`,
     `full-access` for `codex`, and to `bypassPermissions` for `claude`. No
     other runtime and no command override gets any.
   - `SupportedFullAccessMode` and, with `EnableFullAccess`, `FullAccessMode`
     hold the first candidate, so `AccessModeFor` and
     `ValidateRequestedAccessPolicy` keep their contracts.
   - A `RuntimeSpec` with no `FullAccessModes` uses the name `AccessModeFor`
     returns as its single candidate.
2. API Contract: the advertised modes.
   - `ParseSessionConfigOptions` sets `Modes` to the projected select option
     whose identifier is `mode`, else to the first whose category is `mode`,
     else nil. `Options` and every existing issue rule stay.
   - The advertised modes are `Modes.Values`, in advertised order, or none
     when `Modes` is nil.
3. API Contract: the mode selection.
   - `chooseFullAccessMode` returns the first candidate, in candidate order,
     that the advertised values contain, compared exactly after trimming
     space.
   - `applyFullAccess` sends `set-mode <mode>` with that candidate only. The
     proof passes the capabilities its `startSessionSelection` observed, and
     `ensureSession` passes the ones it observed.
   - A refused `set-mode` keeps today's `AccessPolicyError` of kind
     `access_policy_rejected`, predicate `adapter_accepts_access_mode`.
4. API Contract: the Codex sandbox preset.
   - After a successful `set-mode` on a `codex` runtime, `applyFullAccess`
     sends `set sandbox_mode danger-full-access` when
     `capabilities.Options` holds an option with the identifier
     `sandbox_mode`.
   - A failure keeps today's handling: the status
     `codex_sandbox_full_access_unavailable` and the policy `full-access
     (degraded: danger-full-access sandbox preset unavailable)` for an
     unavailable preset, an error otherwise.
   - Without that option, no preset is sent and the policy is `full-access`.
5. API Contract: the degraded report.
   - When no candidate is advertised, `applyFullAccess` sends no `set-mode`
     and no preset.
   - It returns `full-access (degraded: no advertised mode grants full
     access; advertised modes: <v1>, <v2>, ...)`, or `... advertised modes:
     none)` when `Modes` is nil or empty.
   - It publishes the status `full_access_mode_not_advertised` and warns
     `full-access mode not advertised for Agent Session <name>; advertised
     modes: <list>`.
   - The Run continues; `profiles validate` and `doctor` print the policy
     through their existing degraded rule.
6. API Contract: everything else stays.
   - A default-access request sends no `set-mode` and no preset.
   - A runtime without candidates is refused `access_policy_unsupported`.
   - The sealed Baseline analysis arguments, the model and reasoning
     selection, the pinned adapter versions, and the existing degraded
     sandbox status and text are unchanged.
7. API Contract: the documentation texts.
   - The `qa-gate` skill gains the heading `### Browser driver inside a
     sandbox`. It says that when a daemon-backed browser tool such as
     agent-browser cannot start inside the Run sandbox, the gate drives the
     browser through the chrome-devtools MCP server when the session exposes
     it, because that server runs outside the Agent sandbox. The gate
     records which driver ran each browser row. Without either driver the
     row is `blocked (environment: <error>)`, never retried in a loop.
   - The `.roundfixrc.yml` comment keeps `agent_full_access: false`. It says
     the repository keeps the Agent sandbox on purpose, that a true value
     selects the mode the session advertises (`agent-full-access` on
     codex-acp 2.x), and that the QA gate drives browsers through the
     chrome-devtools MCP server.

### Surface Transcripts

None. No command gains an output line of its own: the new effective access
policy text travels through the existing `profiles validate` and `doctor`
degraded format, and API Contract 5 fixes the text the named test asserts.

### Invariants

```text
1. Roundfix sends set-mode only with a value the Agent Session advertised in its mode option.
2. A full-access request yields the effective policy full-access only after a set-mode of an advertised candidate succeeded and, on a Codex session advertising sandbox_mode, the preset succeeded.
3. A default-access session receives no set-mode and no sandbox preset.
4. A Claude full-access session advertising bypassPermissions is set to bypassPermissions.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, User Story 1, Core Feature 1, Success Metric 1 → the candidates, the advertised modes and the selection (API Contracts 1-3)
- Goal 2, Goal 3, Core Feature 2, Success Metric 2 → the sandbox preset and the Claude path (API Contracts 3-4)
- Goal 4, User Story 2, Core Feature 3, Success Metric 3 → the degraded report (API Contract 5)
- User Story 3, Core Feature 4 → ADR-0011's note, the configuration guide, the `profiles` reference, `.roundfixrc.yml`, `CONTEXT.md` and the `qa-gate` skill (API Contract 7)
- Success Metrics 1-3, API Contract 6 and Invariants 1-4 → the named tests and the scope audit

## Absence Evidence

1. The full-access candidate list — command: `grep -rn 'FullAccessModes' internal`; result: no match at `f382480b`.
2. The session's mode option in the capability projection — command: `grep -rn 'Modes \*SelectCapability' internal/agent`; result: no match at `f382480b`.
3. The degraded report for an unadvertised mode — command: `grep -rn 'full_access_mode_not_advertised' internal docs .agents`; result: no match at `f382480b`.
4. The sandboxed browser guidance — command: `grep -rn 'chrome-devtools' .agents/skills/qa-gate`; result: no match at `f382480b`.

## ADR Alternatives

- ADR-0293: this Spec's decision; it rejects a hard-coded `agent-full-access`, the older `modes.availableModes` field, a name or description heuristic, refusing a session without a candidate, always sending the sandbox preset, and editing `~/.codex/config.toml`.
- ADR-0011: amended by ADR-0293 with a dated note; its rejected alternative of editing the runtime's configuration files stays rejected, and its fail-fast rule keeps holding for an advertised mode the adapter refuses.
- ADR-0017: no rejected alternative is taken; the mode is still set through acpx `set-mode`, and the advertised modes come from acpx's public `sessions show` output, not from a session file.
- ADR-0189: no rejected alternative is taken; the `qa-gate` and Roundfix skill edits raise their recorded versions.
- ADR-0285: no rejected alternative is taken; every new test compiles for Windows.
- ADR-0125: no rejected alternative is taken; the fake acpx runner stays the compiled test binary.
- ADR-0104: no rejected alternative is taken; the outside-evidence rows rest on the adapter's and the protocol's published documentation.
- ADR-0233: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0240: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0250: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0256: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0276: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0167: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0229: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0272: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0197: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0237: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0279: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0286: no rejected alternative is taken; this Spec changes none of its rules.

## Integration Points

acpx 0.19.6 and the ACP adapters behind it, through the `sessions show` and
`set-mode` commands Roundfix already runs. No test starts acpx or an adapter.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, uses the fake acpx runner
of `acpx_runner_test.go` through `newFakeACPXHarness`, and compiles for
Windows. In the new `internal/agent/full_access_mode_test.go`:

- `TestChooseFullAccessModeTakesTheFirstAdvertisedCandidate`: candidate order
  wins over advertised order; no match returns false; spaces are trimmed.
- `TestProofSelectsTheAdvertisedCodexFullAccessMode`: the codex-acp 2.x shape
  (`read-only`, `workspace-write`, `agent`, `agent-full-access`, no
  `sandbox_mode`); one `set-mode agent-full-access`, no `set sandbox_mode`,
  effective policy `full-access`.
- `TestRunAppliesTheSandboxPresetWhenTheSessionAdvertisesIt`: `full-access`
  advertised with a `sandbox_mode` option; a Run sends `set-mode full-access`
  then `set sandbox_mode danger-full-access`.
- `TestRunSelectsBypassPermissionsForClaude`: `bypassPermissions` advertised;
  a Claude Run sends `set-mode bypassPermissions` and no preset.
- `TestFullAccessWithoutAnAdvertisedModeIsDegraded`: `read-only`,
  `workspace-write`, `agent` advertised. The proof sends no `set-mode` and
  records the exact degraded text. A Run publishes
  `full_access_mode_not_advertised`, warns, and reaches its prompt. With no
  mode option, the text ends `advertised modes: none)`.

The fake's default `sessions show` and `sessionCapabilitySnapshotFixture`
gain a `mode` option: for Codex the values `read-only`, `auto`, `full-access`
and a `sandbox_mode` option, for Claude `default`, `bypassPermissions`. That
keeps the five pinned tests' invocation lists, which use the legacy name
`full-access` as their single candidate. Where a pinned test still changes,
the Task declares it and names the contract it moves to.
`agent_test.go` and `cli_test.go` move the Codex name to `agent-full-access`.

## Build Order

1. Full-access mode discovery: the candidates, the mode option, the
   selection, the sandbox preset rule, the degraded report, their tests, the
   pinned tests, the Roundfix `profiles` reference and its skill version (no
   dependency). The reference travels with the code because the change
   alters what `profiles validate`, `doctor` and the Run events report.
2. Documentation: ADR-0011's note, the configuration guide, the
   `.roundfixrc.yml` comment, the glossary, the `qa-gate` skill, its version
   and the Behavior Surface Record (depends on: 1, whose behavior they
   describe and whose skill version record they extend).
3. QA gate (depends on: 2).

No other active Spec is a prerequisite. Specs 0282 and 0283 share no file with
this one.

## Risks & Considerations

- A future codex-acp that advertises a `sandbox_mode` option again receives
  the preset, as 1.x did.
- An adapter that names its full-access mode differently is reported
  degraded until a candidate is added; the report names its modes.
- ACP version 2 removes `session/set_mode`; when acpx follows, the mode moves
  to `set mode <value>` through the same advertised option.

## Glossary

- adds: **Full-Access Mode**

## Decisions

- The advertised modes come from the session's `mode` config option in
  `acpx sessions show`, and the preference is an exact candidate list per
  runtime. See ADR-0293.
- No advertised candidate means a degraded policy, a status and a warning,
  and no `set-mode`. See ADR-0293.
- The Codex sandbox preset is sent when the session advertises
  `sandbox_mode`. See ADR-0293.
- `applyFullAccess` takes the capabilities its callers already observed, so
  no new `sessions show` call is added.
