---
spec: 0365-the-capability-projection-reads-grouped-select-options
prd: _prd.md
created: 2026-11-09
---

# The capability projection reads grouped select options — Technical Spec

## Executive Summary

`projectSelectCapability` (`internal/agent/selection_capabilities.go`)
learns the protocol's second shape: when every entry of a select's `options`
carries a nested `options` array, the entries are groups and their values
are concatenated, group by group, into the one list the existing checks
already read. A list with no such entry is today's flat list and takes
today's path unchanged. The trade-off this design accepts is that a
malformed option Roundfix never assigns still refuses the whole evidence
(ADR-0375), in exchange for never turning a malformed payload into a
silently degraded or unsandboxed session. Two recorded adapter payloads
become fixtures, and the proof is exercised end to end through the existing
fake `acpx` harness.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; a local `acpx` JSON payload is parsed. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0375 (new) records the rule.
  ADR-0105: "A capability projection bounds what it keeps, not what a
  runtime advertises", applied to the flattened values. ADR-0293: "A
  full-access request selects the mode the Agent Session advertises", whose
  advertised modes are the flattened values. ADR-0011: "ACP Runtimes can run
  Batches with explicit full access", cited by ADR-0293 and unchanged. ADR-0314: "A Light Tier that cannot serve the Run yields to the category
  profile", cites ADR-0011 and does not apply: no Light Tier planning
  changes. They are
  reached only through citations and do not apply either: ADR-0335 ("A stop
  that is not the work's fault does not cost the Run"), ADR-0274 ("A turn
  that did not finish is not accepted"), ADR-0276 ("A row the authorization
  keeps offline qualifies, and a byte invariant compares with the Delivery
  Base") and ADR-0295 ("Reachability reads the fetched default branch, and
  origin names it when origin/HEAD is unset"); no Run outcome, turn
  acceptance, QA row rule or reachability rule changes.
  Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BA answers of 2026-11-09 quoted in the authorization record; no
  Makefile, workflow, Go module, lint or formatter change and no Governed
  Path. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0365-the-capability-projection-reads-grouped-select-options/_authorization.md`;
  bounded files: none (`paths: []`).

## System Architecture

- `internal/agent/selection_capabilities.go`: `acpSelectValuePayload` gains
  the group fields; `projectSelectCapability` obtains its value entries from
  one new helper, `flattenSelectValues`, before the existing loop that
  bounds values, finds duplicates and finds the current value. Retention
  (`retainAdvertisedValues`), `ParseSessionConfigOptions`,
  `ParseSessionCapabilitySnapshot` and the `Modes` selection are unchanged
  and read the flattened values.
- `internal/agent/acpx_runner.go`, `applyFullAccess`: unchanged; it reads
  `capabilities.Modes.Values`, which now include the values of every group.
- `internal/agent/testdata/`: the two recorded payloads, byte copies of
  `references/codex-acp-2.6.0-config-options.json` and
  `references/codex-acp-2.5.0-config-options.json`.
- Unchanged by design: the issue vocabulary and `MaxCapabilityDiagnosticIssues`,
  the Doctor and `profiles validate` output, every guide and skill.

## Implementation Design

### Interfaces

```go
// internal/agent/selection_capabilities.go
type acpSelectValuePayload struct {
	Value   string                   `json:"value"`
	Group   *string                  `json:"group"`
	Options *[]acpSelectValuePayload `json:"options"`
}

// flattenSelectValues returns a select's value entries: the flat list as
// given, or every group's values in advertised order. It records an issue
// and returns false for a shape the protocol does not define.
func flattenSelectValues(entries []acpSelectValuePayload, issues *capabilityIssueSet) ([]acpSelectValuePayload, bool)
```

### Data Models

No changes. `SelectCapability` keeps `ID`, `CurrentValue`, `Values` and
`AdvertisedCount`; a grouped option projects without group identity.

### API Contracts

1. API Contract: flat select. A list in which no entry carries a non-null
   `options` array projects exactly as at `2e59aebe`: same values, order,
   `AdvertisedCount`, retention and issues.
2. API Contract: grouped select. A list in which every entry carries a
   non-null `options` array is a list of groups. Each group's `group` must
   be present and bounded (`boundedCapabilityID`). Its values are appended in
   advertised order, group after group; a group with an empty `options`
   array contributes nothing. The flattened entries then pass the flat
   checks: each value bounded or `invalid_select_option`, a repeated value
   anywhere in the option `duplicate_option_value`, the current value among
   them or `invalid_current_value`, retention over the flattened list and
   `AdvertisedCount` equal to the number of distinct valid values.
3. API Contract: refusals. `invalid_select_option` for a list mixing value
   entries and groups, a group without `group` or with an unbounded one, an
   entry inside a group that carries its own `group` or `options`, and an
   option whose flattened list is empty. `too_many_option_values` when the
   outer entry count or the flattened value count exceeds
   `maxAdvertisedCapabilityValues`. Every refusal refuses the whole
   evidence, whichever option carries it (ADR-0375).
4. API Contract: proof. Over the recorded 2.6.0 payload,
   `ProveExactSelection` for `codex`, `gpt-6.1-sol`, `high` returns a proof
   with encoding `independent`; with a full-access request its effective
   policy is `full-access` and the only access call is `set-mode
   agent-full-access`; without one it sends no `set-mode`. Over the 2.5.0
   payload it returns the same assignment and calls.

### Surface Transcripts

None. No command output changes for any payload that passes today, and the
command that fails today needs the real adapter; the proof is asserted
through the fake `acpx` harness and the live run is the PRD's Unreachable
Acceptance entry.

## Coverage Map

- Goal 1, Story 1, Core Feature 2, Success Metrics 1-3 → API Contracts 2
  and 4, Build Orders 1-2.
- Goal 2, Story 2, Success Metric 3 → API Contract 4, Build Order 2.
- Goal 3, Success Metric 2 → API Contract 1, Build Order 1.
- Goal 4, Story 3, Core Feature 3, Success Metric 4 → API Contracts 2-3,
  Build Order 1.
- Core Feature 1 → API Contract 2, Build Order 1.

Deviations from the PRD: None.

## Absence Evidence

1. `flattenSelectValues` and the group fields — command: `grep -n 'json:"group"\|flattenSelectValues' internal/agent/selection_capabilities.go`; result: no match; `acpSelectValuePayload` holds only `Value`, and `projectSelectCapability` iterates `*raw.Options` directly.
2. The recorded payload fixtures — command: `ls internal/agent/testdata`; result: `cursor-session-new-published.json`, `cursor-session-recorded.json` and `usage/`; no codex-acp payload.
3. The grouped and proof tests — command: `grep -rln '"group"' internal/agent`; result: no match; no test or source names the protocol's group field, so no test feeds a grouped select.

## ADR Alternatives

- ADR-0375: the design takes none of its rejected options: no tolerance for
  a malformed unassigned option, no single-group read, no per-group
  options, no version 2 field names.
- ADR-0105: rejected raising the value cap and refusing by size; the
  flattened count is held to the same absolute ceiling and retention stays
  by reference, never by position within a group.
- ADR-0293: rejected hard-coding a mode, reading `modes.availableModes` and
  choosing by name; the full-access mode is still the first candidate found
  among the advertised values, which now include grouped ones.
- ADR-0314: does not apply; its rejected options concern Light Tier
  planning, which this design does not touch.
- ADR-0335: does not apply; its rejected options concern how a Run stop is
  charged, which this design does not touch.
- ADR-0274: does not apply; its rejected options concern accepting an
  unfinished turn, which this design does not touch.
- ADR-0276: does not apply; its rejected options concern QA rows and byte
  invariants, which this design does not touch.
- ADR-0295: does not apply; its rejected options concern default-branch
  reachability, which this design does not touch.
- ADR-0011: rejected editing `~/.codex/config.toml`; nothing here reads or
  writes it, and the `permission-profile:localdev` value is only projected.

## Integration Points

The `acpx` strict-JSON output of `sessions show` and `set`, already read
through `ParseSessionCapabilitySnapshot` and `ParseSessionConfigOptions`;
nothing new is spawned or read.

## Testing Approach

- `internal/agent/selection_capabilities_grouped_test.go` (new): the two
  recorded payloads through `ParseSessionCapabilitySnapshot`, and a table of
  constructed payloads for API Contracts 2-3, including a grouped model
  option above the retention bound whose requested model sits in the last
  group.
- `internal/agent/selection_proof_grouped_mode_test.go` (new): API Contract
  4 through `ProveExactSelection` and the fake `acpx` harness, reusing
  `newFakeACPXHarness`, `fullAccessRuntimeForTest`, `selectedRuntime`,
  `assertFullAccessCalls` and `assertLastInvocationClosesDisposable`.
- Existing seams that must keep passing unchanged:
  `TestSelectionCapabilitiesOfficialAndLegacyFixtures`,
  `TestParseSessionConfigOptionsRejectsInvalidEvidence`,
  `TestSelectionCapabilitiesProjectsTheSessionModeOption`, the
  `TestCharacterizationInvariant*` and `TestRetention*` tests, and the
  `full_access_mode_test.go` tests.
- No test starts an adapter, reaches the network or reads `~/.codex` or
  `~/.roundfix`.

## Build Order

1. The grouped projection, the two fixtures and their tests.
2. The Exact Agent Selection Proof over both recorded payloads, with and
   without full access (depends on: 1).
3. Final QA gate (depends on: 2).

## Risks & Considerations

- A later adapter could send a shape outside the protocol; it is refused by
  name, as today, and the operator pins the adapter until a fix ships.
- A Codex user whose permission profile is the current mode keeps it until
  a full-access request sets `agent-full-access`; Roundfix never chooses a
  permission profile.
- Rebase order: this Spec delivers first; 0366 and 0367 touch no file it
  changes.

## Glossary

None.

## Decisions

- Flatten groups into one option's values and keep every flat bound. See
  ADR-0375.
- A malformed option refuses the evidence even when Roundfix never assigns
  it. See ADR-0375.
- Two implementation Tasks and the QA gate, in series.
