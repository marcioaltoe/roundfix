---
task: task_05
spec: 0164-knowledge-lifecycle-and-capture
status: pending
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
