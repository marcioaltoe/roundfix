---
type: fix
status: promoted
created: 2026-10-07
spec: 0245-generated-records-that-hold-across-specs-and-platforms
---

# Generated records break when Specs are authored in parallel or tested on another platform

## Problem

Two checked-in records were wrong for the head CI tested, and an operator had
to rewrite each one by hand.

1. **Baseline module versions.** Specs authored in parallel each write the
   next module version into their TechSpec and Tasks. Spec 0239 took
   `context-workflow` version 22 first, so Spec 0241 shipped asserting a
   version another Spec already owned. The operator bumped it to 23 and
   regenerated the digests (PR #438). Spec 0228 made skill versions pick a
   free number at record time, but nothing does the same for Baseline
   modules.
2. **The coverage record follows the machine that wrote it.**
   `docs/references/coverage-record.json` lists every top-level test the
   suite discovers. On PR #441, re-recording it on macOS added
   `TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded`, a test built only
   on macOS, so Linux CI reported a coverage regression. Linux-only store
   tests show up as additions in the opposite direction. The operator
   removed the entry by hand.

## Direction

- Choose a Baseline module version when the change is recorded (a
  `-record-module-versions` step in the style of 0228), not when the Spec is
  authored. Specs then cite "the next version".
- Make the coverage record platform-neutral. Either list tests with build
  constraints under every supported `GOOS`, or record only tests without
  constraints and keep per-platform sets separate.

## Evidence

- `~/.roundfix-operator/queue-interventions.md`, entries 219 (0241 module
  version) and 224–225 (PR #441 coverage record).
