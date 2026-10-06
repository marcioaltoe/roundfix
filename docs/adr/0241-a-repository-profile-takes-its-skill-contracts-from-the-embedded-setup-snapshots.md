---
status: accepted
created_at: 2026-10-06T00:00:00Z
updated_at: 2026-10-06T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A repository profile takes its skill contracts from the embedded Setup Snapshots

A built-in Baseline Profile names one Setup Snapshot, and skill restore,
lock reconciliation and the snapshot comparison of ADR-0221 read their skill
contracts from it. A repository-owned profile names none: its file holds
modules, decisions, capabilities, templates and values only (ADR-0067). The
three readers resolved only built-in profiles, so on Roundfix 0.43.0
`roundfix baseline update` exited 1 with "Unknown built-in Baseline Profile"
in every repository adopted with its own profile, Doctor reported the
comparison unavailable, and the restore command that profile alignment
recommends refused the profile it was given.

The skill readers now resolve a profile exactly as the rest of `baseline
update` and `baseline profile validate` do: a built-in ID, or the one file
`.roundfix/baseline/profiles/<id>.json`. A built-in profile keeps its Setup
Snapshot. A repository profile requires the skills its own modules require,
and the contract of each of those skills that an embedded Setup Snapshot
pins from an upstream source is the one entry every embedded snapshot
listing that skill agrees on. Snapshots that pin different trees for one
such skill are refused by name, never resolved by catalog order. A profile
that resolves neither way is reported with the repository path that was
searched, not as an unknown built-in profile.

Two alternatives were rejected. Binding the profile to the closest built-in
profile it adapts (ADR-0219) failed on the first measured case: a profile
written by `roundfix baseline profile init --from go-cli-tui` on 2026-10-06
is not a valid adaptation of any built-in profile, because the
`autonomous-work` module requires a decision the copy does not list, and a
profile that keeps only shared modules ties between profiles. Adding a
`setup` field to the repository profile would change the strict schema every
adopted profile already passes.

## Consequences

A repository profile's restorable skill set is the external skills its
modules require, not every skill of some upstream list, so a restore without
`--skill` never installs another stack's skills. A restore or reconcile plan
for a repository profile carries a null `setup`. Every embedded snapshot pins
the same upstream commit today, and a test holds them in agreement, so the
refusal is reachable only through a catalog change that a test also
reports. Roundfix-owned skills refresh through `baseline update` in these
repositories again: the owned install ran before the comparison, but the
command then failed, so its preview never reached `current` and its
`--yes` run never reported the refresh as verified.
