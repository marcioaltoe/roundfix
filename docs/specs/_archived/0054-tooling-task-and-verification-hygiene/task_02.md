---
task: task_02
spec: 0054-tooling-task-and-verification-hygiene
status: completed
type: infra
complexity: medium
---

# Task 02: Ship the sanctioned regeneration target and a portable cache

## Overview

Expose the update modes as one sanctioned command, make the repository gate
behave identically inside and outside the ACP sandbox by defaulting the Go
build cache to a repository-local path, and stop the bare compile check from
leaving an unignored binary in the tree. This Task owns every authorized
change to the build tooling files.

## Requirements

1. MUST add a regeneration target that runs every update mode from Task 01 in
   dependency order, is idempotent, and is registered as a phony target.
2. MUST default the Go build cache to a repository-local, ignored path when
   the environment does not set one, and MUST leave an explicitly exported
   cache untouched so developer and CI overrides keep winning.
3. MUST ignore the repository-local cache directory and the path a bare
   `go build ./cmd/roundfix` writes, so neither can be reported as untracked
   or swept into a commit.
4. MUST NOT change what the gate verifies beyond cache determinism.
5. MUST change only `Makefile` and `.gitignore`, plus this Task file — the
   exact protected paths this Spec authorizes for build tooling.

## Subtasks

- [ ] Add the regeneration target and register it as phony.
- [ ] Default and export the build cache only when unset.
- [ ] Add the ignore entries for the cache directory and the bare-build
      binary, each with a comment naming its source.

## Acceptance Criteria

- [ ] The regeneration target runs every update mode and, on an unchanged
      repository, leaves the working tree clean; a second consecutive run
      also changes nothing.
- [ ] With no cache exported, the gate runs against the repository-local
      cache; with one exported, that value is used unchanged.
- [ ] `go build ./cmd/roundfix` leaves no path reported as untracked.
- [ ] The changed-path set for this Task is exactly `Makefile`,
      `.gitignore`, and this Task file.

## Context

- interface: `Makefile`
- interface: `.gitignore`

## Verification

- `grep -q '^baseline-digests:' Makefile` — expected: the sanctioned target exists.
- `grep -q 'GOCACHE ?=' Makefile` — expected: the cache defaults only when unset.
- `grep -q '^/roundfix$' .gitignore` — expected: the bare-build binary is ignored.
- `before="$(git status --porcelain --untracked-files=all)" && go build -buildvcs=false ./cmd/roundfix && test "$before" = "$(git status --porcelain --untracked-files=all)"` — expected: the bare build leaves the working-tree status unchanged.
- `git status --porcelain -- Makefile .gitignore` — expected: both files are the only build-tooling paths this Task touched.

## References

`_prd.md` → User Story 3, Core Features 1, 5, 6, Project Constraints:
Tooling authority; `_techspec.md` → Build Order 2 and 3, Interfaces: Make
targets.

## Result

### Implementation

- Added the phony `baseline-digests` target. It runs all five Task 01 update
  modes in dependency order: Skill setup snapshots, normalized catalog,
  parity corpus, maintained Source Baseline, then formatter composition.
- Defaulted and exported `GOCACHE` as `$(CURDIR)/.gocache` with `?=`, so an
  environment-provided value remains authoritative.
- Ignored `/.gocache/` and `/roundfix`, with comments identifying the
  Makefile cache default and bare `go build ./cmd/roundfix` as their sources.
- Left the `verify` target and its dependencies unchanged.

### Focused checks

- Red signal: `rtk make -n baseline-digests` initially failed with
  `No rule to make target 'baseline-digests'`.
- After implementation, `rtk make -n baseline-digests` printed the five
  update commands in dependency order.
- Two consecutive `rtk make baseline-digests` runs passed. Each run reported
  19 passing Skill tests plus one passing test for each of the four Baseline
  selectors.
- The binary diff hash for `internal/baseline/assets` and
  `internal/baseline/testdata` remained
  `8b137891791fe96927ad78e64b0aad7bded08bdc` before and after the
  consecutive regeneration runs.
- With `GOCACHE` removed from the environment, portable Make database output
  resolved `GOCACHE = $(CURDIR)/.gocache`. With
  `GOCACHE=/tmp/roundfix-explicit-cache`, it resolved
  `GOCACHE = /tmp/roundfix-explicit-cache` unchanged.
- `rtk git check-ignore -v roundfix .gocache/probe` matched `/roundfix` and
  `/.gocache/` at their intended `.gitignore` lines.
- `rtk git diff --check` passed.
- An initial GNU Make-only `--eval` cache-introspection probe was unsupported
  by macOS Make. The portable `make -pn` probes above replaced it and passed;
  the failed probe changed no files.

### Acceptance evidence

- **Idempotent sanctioned regeneration:** both consecutive target runs
  passed, the derived-artifact diff hash stayed identical, and no additional
  changed path appeared. The maintained Source Baseline precedes formatter
  regeneration, following Task 01's proven dependency.
- **Portable cache with override precedence:** the Make database showed the
  repository-local default when the environment omitted `GOCACHE` and the
  exact exported override when present. Daemon Verification still owns the
  repository-gate execution.
- **Ignored bare-build output:** `git check-ignore` proved that the root
  `roundfix` path is ignored. The Agent did not run the declared bare-build
  command; the Daemon owns it.
- **Bounded changed paths:** `rtk git diff --name-only` reported exactly
  `.gitignore`, `Makefile`, and this Task file after the final focused checks.

### Verification feedback repair

- Daemon Verification attempt 1 built the bare binary but failed because its
  status filter retained this Task's intentional `M .gitignore` and
  `M Makefile` entries. The diagnostic contained no reported `roundfix` path,
  so the ignore behavior was not the failing component.
- Replaced that self-contradictory check with a before/after porcelain-status
  comparison around the build. The repaired command detects any path the
  build adds while allowing the Task's pre-existing authorized changes.
- The Agent did not rerun the declared command. Focused read-only inspection
  confirmed the Daemon-created `roundfix` path remains matched by
  `.gitignore`, and the changed-path set remains bounded to the authorized
  three files.

### Follow-up

- The TechSpec's illustrative Make snippet lists formatter regeneration
  before maintained Source Baseline regeneration. Task 01's Result proves
  the maintained Source Baseline must run first when corpus bytes change;
  updating the Spec example is outside this Task's authorized path set.
- The Task's declared `## Verification` commands were not rerun by the Agent
  during this repair; the Daemon owns the final rerun and terminal settlement.
