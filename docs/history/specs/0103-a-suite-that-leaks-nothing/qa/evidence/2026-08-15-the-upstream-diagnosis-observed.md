---
spec: 0103-a-suite-that-leaks-nothing
date: 2026-08-15
kind: outside-evidence
row: R11
source: https://github.com/golang/go/issues/22315
---

# The upstream diagnosis, observed

Fetched from the operator's machine on 2026-08-15, outside the Agent Session
whose sandbox has no network. This is the outside source Spec 0103 rests on: it
was written by the Go project, not by this repository, and it describes the
failure mode before this repository ever measured it.

**golang/go#22315 — "os: StartProcess ETXTBSY race on Unix systems."** Open,
assigned to Backlog.

The mechanism it describes:

- One thread writes an executable and opens it with `O_CLOEXEC`. Another thread
  forks and execs something else. The first thread's descriptor leaks into the
  second's child before that child reaches `exec`.
- If the writer then closes its descriptor and execs before the inherited copy is
  closed, the exec fails with `ETXTBSY`, because Unix refuses to execute a file
  that any process anywhere still holds open for writing.
- The issue states the general case plainly: this race is faced "every time anyone
  writes a program that both writes and executes a program."

The workarounds it discusses are a sleep-and-retry loop of up to one second in
`os.StartProcess`, and locking that synchronises close against fork+exec — for
which the issue records no portable implementation. `vfork(2)` does not help.

Two things follow for this Spec, and both are why ADR-0125 chose what it chose.

The diagnosis is not this repository's invention. The `text file busy` failure
observed on 2026-08-14 in `TestProjectDecisionJourney`, and the empty `--version`
probes observed on 2026-08-10 and again on 2026-08-14, are the same upstream race
caught at two instants — the exec refused, and the exec succeeding against a file
whose bytes were not yet readable.

And the remedy the issue offers is the one this Spec refused. Retrying is what
upstream can do from inside `os.StartProcess`, where the write is somebody else's
code. A test suite owns both halves, so it can stop writing the executable
altogether: hard-link a binary that was compiled before any forking began. There
is then no descriptor to leak and no window to retry into.
