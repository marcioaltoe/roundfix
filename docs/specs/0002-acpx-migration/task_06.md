---
task: task_06
spec: 0002-acpx-migration
status: completed
type: docs
complexity: low
---

# Task 06: Update docs and the Roundfix skill for the acpx dependency

## Overview

Bring the user- and agent-facing docs in line with the new agent layer: the canonical Roundfix skill and the README name the acpx pin and the Node prerequisite, the latency recommendation is recorded, and handoff work-plan item 2 is marked done. Verifiable through the skills drift check inside the full gate.

## Requirements

1. MUST document in the canonical Roundfix skill (and regenerate the embedded copy through the sync target): the agent layer runs through acpx at the pinned version, Node is a prerequisite, the install command, and that Runs drive one Agent Session per Run — without changing any documented command, flag, output, or exit-code text (none changed).
2. MUST update the README (or the equivalent install/usage doc) with the acpx prerequisite and the recommendation to configure direct adapter binaries in acpx config for latency-sensitive setups instead of default npx launches.
3. MUST record handoff work-plan item 2 as done in the handoff document's work plan (a status note, not a rewrite — history stays).
4. MUST verify every term used comes from the glossary (Agent Session, ACP Runtime, Run, Work Item, Stop Request); call out any gap in the Result instead of inventing language.

## Subtasks

- [x] Roundfix skill: acpx pin, Node prerequisite, Agent Session note; embedded copy regenerated
- [x] README/install docs: prerequisite plus latency recommendation
- [x] Handoff work-plan item 2 status note
- [x] Glossary coverage pass

## Acceptance Criteria

- [x] The canonical skill names the exact pinned version shipped by task_03's constant and the install command; the drift check passes inside the full gate.
- [x] The README prerequisite section matches what Preflight Validation actually demands (same version, same command).
- [x] The handoff document shows item 2 closed with a pointer to this spec.
- [x] No new un-glossaried term appears in the updated text.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts validate.
- `make verify` — expected: full gate passes, including the skills drift check.

## References

`_prd.md` → Core Feature 1; Non-Goals (skill semantics unchanged beyond the dependency); User Experience. `_techspec.md` → Integration Points, Build Order 6, Risks (first-run latency). Repo hard rule (canonical skill ships with behavior changes). ADR-0017.

## Result

Updated the canonical Roundfix skill and regenerated `skills/roundfix` with `rtk make skills-sync`. The skill now documents that ACP Runtimes run through acpx `0.12.0`, Node.js 22.13 or newer with npm/npx is required, `npm install -g acpx@0.12.0` installs the pinned dependency, and each Run drives one acpx-backed Agent Session across its Work Items. The public command, flag, stdout, and exit-code text was not changed.

Updated the README requirements to match `internal/agent.PinnedACPXVersion` (`0.12.0`) and the `ACPXProbeError` install command (`npm install -g acpx@0.12.0`). The README also records the TechSpec's latency recommendation: configure direct adapter binaries in acpx config so default adapters do not launch through `npx -y` on first use.

Added a dated status note under handoff work-plan item 2 closing it through Spec `0002-acpx-migration`, with pointers to ADR-0017 and this Spec's Task evidence.

Glossary pass: the updated Roundfix terms are covered by `CONTEXT.md`: Agent Session, ACP Runtime, Run, Work Item, Spec, Task, and Preflight Validation. No new Roundfix domain term was introduced; `adapter binaries` remains lower-case acpx wording required by the TechSpec latency note.

Verification:

- `rtk make skills-sync-check` — passed with no drift output.
- `rtk go run ./cmd/roundfix skills check` — passed: `Roundfix skill check passed: roundfix`.
- `rtk make verify` — initial unisolated run failed in `go test ./...` because this machine's global Git config has `commit.gpgsign=true`; six temporary-repo tests failed before this docs diff was exercised with `gpg failed to sign the data`.
- `rtk env GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=commit.gpgsign GIT_CONFIG_VALUE_0=false make verify` — passed after the final edit: 445 Go tests passed in 16 packages, `roundfix skills check` passed, and `go build -buildvcs=false -o bin/roundfix ./cmd/roundfix` passed.

Acceptance evidence:

- Canonical skill pin/install: `.agents/skills/roundfix/SKILL.md` names acpx `0.12.0` and `npm install -g acpx@0.12.0`; `internal/agent/acpx_runner.go` defines `PinnedACPXVersion = "0.12.0"` and builds the same install command.
- README prerequisite parity: README requirements name Node.js 22.13 or newer, acpx `0.12.0` on `PATH`, and the same Preflight Validation install command.
- Handoff closure: `docs/handoffs/2026-07-04-implementation-daemon-acpx.md` item 2 now has a 2026-07-05 status note pointing to Spec `0002-acpx-migration`.
- Glossary coverage: no new capitalized Roundfix domain term was added beyond terms already defined in `CONTEXT.md`.
