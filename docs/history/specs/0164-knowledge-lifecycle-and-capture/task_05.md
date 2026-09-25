---
task: task_05
spec: 0164-knowledge-lifecycle-and-capture
status: completed
type: docs
complexity: low
---

# Task 05: Capture names its publication owner

## Overview

The Secondbrain guide says "Commit the entry at the moment of capture", while the installed autosync job commits and pushes the shared checkout, and nothing distinguishes a local file from a published one.

This is an authorized tooling Task. It may change only `internal/baseline/assets/modules/secondbrain.json`, `docs/agents/secondbrain.md`, `docs/agents/setup-context.json`, the derived pins the sanctioned regeneration rewrites, and this Task file. Stop before any other governed mutation. The bounded set comes from [_authorization.md](_authorization.md).

## Requirements

1. MUST replace the sentence "Commit the entry at the moment of capture, because durability is the point of the door." in `clause.secondbrain.capture-self-contained` of `internal/baseline/assets/modules/secondbrain.json` with a contract that uses the literal phrase `publication owner`: an installed scheduler that commits and pushes the checkout, or the capturing session when none is installed.
2. MUST state that under a scheduler the session writes the entry and does not commit, and that without one the session commits and pushes only its own entry's path, never stages another session's change, and never installs a second schedule.
3. MUST name three evidence levels with the literal phrases `captured locally`, `committed locally` and `confirmed on the remote`, and state that a session reports only the level it observed and that a scheduler's presence or last exit status is never confirmation.
4. MUST keep the rest of that clause and every other clause, bumping versions by the catalog convention.
5. MUST render `docs/agents/secondbrain.md` and `docs/agents/setup-context.json` through `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`, then run `make baseline-digests`, and MUST NOT hand-edit a pin.
6. MUST NOT change any file of the companion Secondbrain repository.

## Subtasks

- [ ] Rewrite the capture sentence in the module.
- [ ] Regenerate the guide, manifest and pins.

## Acceptance Criteria

- [ ] The module and the rendered guide name the publication owner and the three evidence levels.
- [ ] Neither commands a commit at capture.
- [ ] The catalog stays compatible after regeneration.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/baseline/assets/modules/secondbrain.json`
- interface: `docs/agents/secondbrain.md`

## Verification

- `for f in internal/baseline/assets/modules/secondbrain.json docs/agents/secondbrain.md; do grep -q "publication owner" "$f" || exit 1; grep -q "captured locally" "$f" || exit 1; grep -q "committed locally" "$f" || exit 1; grep -q "confirmed on the remote" "$f" || exit 1; if grep -q "Commit the entry at the moment of capture" "$f"; then exit 1; fi; done; go test -count=1 -run "^TestCatalogCompatibility$" ./internal/baseline` — expected: exit 0; before this Task neither file names the publication owner and both still command the commit, so the command fails.

## References

- [_techspec.md](_techspec.md) — Capture ownership

## Result

Implementation:

- Replaced the unconditional capture-time commit instruction with a publication-owner contract: an installed scheduler owns commit/push when present; otherwise the capturing session owns only its entry path and does not install another schedule.
- Added the three observed evidence levels — `captured locally`, `committed locally`, and `confirmed on the remote` — and the rule that scheduler presence or exit status is not remote confirmation.
- Bumped the Secondbrain module, supporting-guide, and owning-rule versions from 9 to 10, then regenerated the managed guide and setup manifest through the sanctioned Baseline renderer. The sanctioned digest regeneration updated only derived catalog pins and fixture/golden snapshots.

Focused checks:

- `GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 ./internal/baseline`: passed.
- `jq` module-version assertion: passed; module, supporting guide, and `rule.secondbrain.read-only` are all version 10.
- Required publication-owner and evidence-level literals were present in the module and rendered guide; the old capture-time commit sentence was absent.
- `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`: applied and verified after rerunning with task-scoped cache and required access to the shared transaction lock.
- `make baseline-digests`: passed with no further changes after regeneration.
- `git diff --check`: passed; changed paths are limited to the authorized module, rendered outputs, sanctioned derived pins, and this Task file. No companion Secondbrain repository path changed.

Acceptance evidence:

- The module and rendered guide carry the publication owner and all three evidence levels.
- The unconditional `Commit the entry at the moment of capture` instruction was removed; commit/push behavior is conditional on scheduler ownership and scoped to the capturing session's path.
- Catalog version and digest regeneration completed through the sanctioned commands; the focused baseline package tests passed.

Not run:

- The declared `## Verification` command remains daemon-owned and was not run in this turn.

## Carry-forward provenance

- Source Run: `run_20260924T223414Z_54e6f85ae98b6550`
- Source commit: `c9718dbc88bf866e7493a1ae053ccd01c70f8f31`
