---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-14T21:05:00Z
updated_at: 2026-08-14T21:05:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# A spawned fixture is a compiled binary, never a written script

A test that writes an executable and then executes it races every other goroutine
in the process: between the write's `open` and its `close`, a concurrent
`fork`/`exec` elsewhere in the suite duplicates the still-writable descriptor into
a child, and the kernel then refuses the exec with `ETXTBSY` — `text file busy` —
or the fixture starts before its bytes are flushed and answers nothing. Go carries
this as a known hazard of writing executables from a forking process
(golang/go#22315); this repository has now measured both faces of it, an empty
`--version` probe on 2026-08-10 and a literal `text file busy` on 2026-08-14.
Fixtures are therefore compiled once — a `go test` binary re-executed through
`os.Args[0]`, or a binary built before the suite begins forking — so no test
executes a file that something may still hold open. A retry would hide the window;
compiling removes it.
