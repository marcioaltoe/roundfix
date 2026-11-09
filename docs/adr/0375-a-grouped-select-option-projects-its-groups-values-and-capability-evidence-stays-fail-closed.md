---
status: accepted
created_at: 2026-11-09T11:45:00Z
updated_at: 2026-11-09T11:45:00Z
deprecated_at: null
superseded_by: null
---

# A grouped select option projects its groups' values, and capability evidence stays fail-closed

On 2026-11-09 codex-acp 2.6.0 started advertising its `mode` select as groups
(`sandbox-modes` and `permission-profiles`) with the current value
`permission-profile:localdev`, while `model`, `reasoning_effort` and the
other selects stayed flat. Roundfix read a select's `options` only as a list
of `{value}` entries, so every group entry failed as `invalid_select_option`,
the current value was never found (`invalid_current_value`), and the whole
capability evidence was refused: every Codex profile failed preflight on
every machine that updated the adapter. The Agent Client Protocol defines a
select's `options` as "either a flat array of values or an array of groups.
Values and groups are not mixed in the same array", and each
`SessionConfigSelectGroup` "requires `group` (its unique string identifier),
`name` (its human-readable label), and `options` (an array of
`SessionConfigSelectOption` objects)" (agentclientprotocol.com, protocol v1,
Session Config Options). Roundfix had implemented only the first shape.

## Decision

**A grouped select option is one option whose values are its groups' values,
in advertised order.** The capability projection reads a list whose entries
all carry a nested `options` array as groups and concatenates their values
group by group; a list with no such entry stays a flat list and projects
exactly as before. The flattened values then pass every check a flat list
passes: each value is bounded, a value repeated in two groups is
`duplicate_option_value`, the current value must be one of them or the
option is `invalid_current_value`, the count of values is held to the same
ceiling (`too_many_option_values`), and retention keeps the current value
and the requested Agent Selection's values first (ADR-0105). The projected
option keeps no group identity: the advertised modes that ADR-0293 reads are
these values, so a full-access candidate inside a group is still chosen.

**Every shape the protocol does not define refuses the evidence, whichever
option carries it.** A list that mixes values and groups, a group without a
bounded `group` identifier, or a group nested inside a group is
`invalid_select_option`. A group with no values contributes nothing, as the
protocol's version 2 tells a Client to hide an empty group, but an option
whose groups advertise no value at all is still `invalid_select_option`.
These rules hold for options Roundfix never assigns as much as for `model`
and the reasoning control.

## Considered Options

- Tolerating a malformed option that Roundfix does not assign, dropping it
  and keeping the rest of the evidence, was rejected. The set of options
  Roundfix reads is wider than `model` and the reasoning key: a full-access
  request chooses its mode from the `mode` option (ADR-0293), and the Codex
  sandbox preset depends on whether `sandbox_mode` is advertised. Dropping
  either would turn a malformed payload into a silently degraded or
  unsandboxed session instead of a named refusal. The failure of 2026-11-09
  was a protocol shape Roundfix did not implement, not a malformed option,
  and the refusal is what made it visible within hours.
- Reading only the first group, or only the group that holds the current
  value, was rejected: it hides advertised values, so a full-access mode in
  another group would be reported as not advertised.
- Keeping each group as its own projected option was rejected: groups are
  display headers inside one option with one current value, and Roundfix
  assigns the option, not a group.
- Accepting the protocol version 2 field `groupId` was rejected for this
  change: version 2 also renames the option identifier, which Roundfix does
  not read, so accepting one renamed field would not make a version 2
  payload readable.
