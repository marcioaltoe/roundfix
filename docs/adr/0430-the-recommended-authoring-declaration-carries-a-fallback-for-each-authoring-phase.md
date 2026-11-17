---
status: accepted
created_at: 2026-11-17T22:45:00Z
updated_at: 2026-11-17T22:45:00Z
deprecated_at: null
superseded_by: null
---

# The Recommended Authoring Declaration carries a fallback for each Authoring Phase

The Recommended Authoring Declaration that the Doctor Command and the Update
Check print gains a fallback for each Authoring Phase:

- `idea`, `prd` and `techspec` fall back to `opus`/`xhigh`;
- `tasks`, declared `opus`/`high`, falls back to `sonnet`/`high`;
- `research` declares none.

The repository contract test compares a fallback in this repository's Project
Config with the recommendation only where one is declared, because the key
reaches that file only after the release. This amends ADR-0415.

## Context

ADR-0415 makes the recommendation a table in the binary that is held equal to
this repository's declaration and to the configuration guide by a repository
contract test. ADR-0427 adds the fallback key. Asked on 2026-11-17 whether the
recommended block carries the fallbacks, the maintainer answered "Incluir os
reservas (Recommended)", with the values of this repository.

A Roundfix older than the release refuses a Project Config that holds the
new key, and the binary of the Delivery Queue reads this repository's Project
Config from the default branch (ADR-0408). This repository therefore gains its
fallbacks in an operator's step after the release, while the table and the
guide gain them in the change that teaches the binary the key.

## Decision

**The table.** The Recommended Authoring Declaration holds the fallbacks
above, and its rendered block prints each as a `fallback` mapping under its
phase, after `reasoning_effort`.

**The contract.**
- The configuration guide's block equals the rendered block, fallbacks
  included.
- For each of the five keys, this repository's Project Config declares the
  recommendation's model and reasoning effort, as before.
- A fallback that the Project Config declares equals the recommendation's,
  and a fallback it does not declare is not compared.

**This repository.** The operator adds the fallbacks after the release; from
then on the contract holds them too.

## Considered Options

- **Leave the recommendation without fallbacks.** Rejected by the maintainer.
- **Compare every fallback with this repository from the start.** Rejected:
  the Task could not pass, because the key cannot reach this repository before
  the release.
- **Skip the fallbacks in the contract.** Rejected: once the operator adds
  them, a drift between the file and the table would go unseen.
- **Add the fallbacks to this repository in the same change.** Rejected by
  ADR-0408.

## Consequences

- Until the operator's step, this repository's Project Config is behind the
  recommendation, and the contract says nothing about it.
- An adopter that pastes the recommended block needs a Roundfix that knows the
  key, which every binary that prints the block does.
