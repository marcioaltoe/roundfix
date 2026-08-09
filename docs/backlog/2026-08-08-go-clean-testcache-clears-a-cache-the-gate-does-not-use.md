---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-08
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# `go clean -testcache` clears a cache the gate does not use

## Opportunity

`docs/agents/specific-repository.md` tells a reader to run `go clean -testcache`
before trusting `make verify`, and the 2026-08-08 handoff repeats it as the
repair for a gate that reported exit 0 over a failing test. The instruction does
not work.

The Makefile exports its own cache:

```make
GOCACHE ?= $(CURDIR)/.gocache
export GOCACHE
```

A bare `go clean -testcache` runs with the user-level `GOCACHE` — on the
maintainer's machine `~/Library/Caches/go-build` — and never touches
`$(CURDIR)/.gocache`, which was 9.4 GB when this was measured. The two caches
are independent, so cleaning one leaves the other exactly as stale as before.

## Value

Measured on 2026-08-08 during Spec 0088. Task 03 renamed a test the coverage
record still named, which `TestCoverageEquivalence` is built to catch. The Task's
gate ran `go clean -testcache && make verify` and reported:

```text
make verify exit=0
ok  	roundfix/internal/spec	(cached)
```

The regression was real and went unreported. It surfaced one Task later, and
only because that Task re-recorded the coverage record for its own reasons.
Running `GOCACHE="$PWD/.gocache" go clean -testcache` first made the same suite
execute in 20.5 s instead of reporting cached.

This is the same false-green class the repository already has a HARD RULE about
— a gate reporting success it did not earn — and the documented workaround is
part of the defect rather than the repair. Anyone following the current
instruction believes they have a cold cache and does not.

## Shape

Non-binding. The cheapest correct fix is a Makefile target that cleans the cache
the gate actually uses, so no reader has to know which `GOCACHE` is in effect —
something a `verify-cold` or `clean-testcache` target would cover, with the
guidance in `docs/agents/specific-repository.md` and the handoff pointing at it
instead of at bare `go clean -testcache`.

Worth settling in the same work: whether `make verify` should refuse to report
success when its Go suite is entirely cached, since a cached pass proves the
cache and not the tree. Both touch protected tooling and need express
maintainer authorization with bounded files before any Task may run.
