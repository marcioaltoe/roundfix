---
type: fix
status: promoted
created: 2026-11-09
spec: 0365-the-capability-projection-reads-grouped-select-options
reason: null
---

# codex-acp 2.6.0 grouped select options break the capability proof

Recorded by the operator on 2026-11-09 (~08:25 BRT) on the MacBook, from the
maintainer's `roundfix doctor` run and two `acpx` session snapshots. The two
captured payloads sit beside this entry:
[codex-acp-2.6.0-config-options.json](codex-acp-2.6.0-config-options.json) and
[codex-acp-2.5.0-config-options.json](codex-acp-2.5.0-config-options.json).

## Symptom

`roundfix doctor` 0.137.1 reports `profiles: failed (runtime="codex",
model="gpt-6.1-sol", reasoning_effort="high"; ... classification:
capability_evidence_invalid; adapter evidence: error="capability evidence
invalid: invalid_current_value, invalid_select_option")`. Every category whose
preferred selection is a Codex tuple fails preflight, so no Run with a Codex
profile can start. Every adopter on a machine that updated the adapter is
blocked the same way.

## Where

- Adapter: `@agentclientprotocol/codex-acp` 2.6.0, published
  2026-11-08T11:40Z and installed by the maintainer's `dev_upgrade` on
  2026-11-09 07:16 BRT on both machines; codex-cli 0.167.0 is unchanged.
- `acpx --format json --json-strict codex sessions show <name>`: the `mode`
  select now advertises grouped options, ACP `SessionConfigSelectGroup`
  entries `{group, name, options: [{value, ...}]}`, with the current value
  `permission-profile:localdev`, a Codex permission profile from the user's
  `~/.codex/config.toml`. The `model`, `reasoning_effort`,
  `collaboration_mode` and `fast-mode` selects stay flat.
- The same probe with codex-acp 2.5.0 (`npx -y
  @agentclientprotocol/codex-acp@2.5.0`) returns only flat selects, with the
  `mode` current value `agent`; with 2.5.0 installed, `roundfix profiles
  validate` passes.
- `internal/agent/selection_capabilities.go`, `projectSelectCapability`,
  reads `options` only as a list of `{value}` entries. A group entry has no
  `value`, so `boundedCapabilityValue("")` fails (`invalid_select_option`)
  and the current value is never advertised (`invalid_current_value`). The
  whole capability evidence is then refused.

## Expected

The capability projection reads a grouped select's values in advertised
order, with the same bounds as a flat select, so the Exact Agent Selection
Proof of a Codex tuple passes on codex-acp 2.6.0 as it does on 2.5.0.

## Evidence

The two captured payloads above. Workaround in place until the fix ships:
codex-acp pinned to 2.5.0 on the MacBook, approved by the maintainer as
temporary.
