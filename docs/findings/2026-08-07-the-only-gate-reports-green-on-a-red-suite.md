---
status: pending
created_at: 2026-08-07
updated_at: 2026-08-07
kind: finding
---

# The only gate reports green on a red suite (2026-08-07)

`make verify` exits `0` on a working tree whose Go test suite exits `1`. The
repository's own rule says the local gate is the only gate; that gate is
currently capable of certifying a red build as verified.

## Reproduction

Same commit, same `GOCACHE`, same working directory:

```
GOCACHE=$PWD/.gocache go test -parallel 16 ./...        → exit 1
                                                          FAIL roundfix/internal/spec
GOCACHE=$PWD/.gocache rtk go test -parallel 16 ./...    → exit 0
                                                          "Go test: 3647 passed in 26 packages"
make verify                                              → exit 0
```

The wrapper does not merely lose the status — it omits the failing package from
its summary, so the output reads as a complete pass. `make verify` inherits this
because the Makefile defines `GO := $(RTK) go`, so every gate invocation goes
through the wrapper.

## Bound of the defect

Not universal. `rtk go test ./internal/spec` — a single package — propagates
exit `1` correctly, and earlier in the same session the wrapper surfaced real
failures through `make verify` (`TestCheckCorpusGolden`,
`TestHumanBaselineDecisionDefaults`) with a correct non-zero status.

The masking was observed on the full `./...` set with a test that emits a large
volume of `t.Log` output alongside its `t.Error` calls: `TestCoverageEquivalence`
logged 278 additions and 2 regressions in the observed run. Whether the trigger
is package count, output volume, or log-to-error ratio is **not established**,
and nobody should act on one of those explanations without reproducing it.

## Why it matters

Every "gate green" claim made through `make verify` is unverified until this is
understood. In this session that includes claims used to authorize commits. The
repository already carries a HARD RULE for the adjacent hazard — never let a
pipe hide a gate's exit status, written after a commit landed on a failing gate
on 2026-07-29 — and the same failure mode now arrives through the Makefile's own
tool wrapper rather than through a caller's pipe.

The rule's remedy ("run the gate on its own, capture `$?`") does not help here:
the gate was run on its own and its `$?` was `0`.

## Interim practice

Until this is resolved, a truthful gate reading requires the unwrapped command:

```
GOCACHE=$PWD/.gocache go test -parallel 16 ./...
```

## Route

Not fixed here. Two candidate owners and the choice matters:

- the wrapper, if it can lose a package-level failure in a multi-package run;
- the Makefile, if the gate should not route its authoritative test invocation
  through an output-filtering tool at all.

The second is worth weighing on its own merits regardless of the first: a gate
whose exit status depends on a summarizing wrapper has a failure mode that a
direct invocation does not.
