---
status: accepted
created_at: 2026-08-02T00:00:00Z
updated_at: 2026-08-02T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Executable capability discovery resolves links without executing

Baseline capability discovery deliberately never runs a candidate binary —
planning stays offline and side-effect free, and executing an arbitrary PATH
entry to prove it exists would be both slow and unsafe. The current
implementation enforces that by `Lstat`-ing each PATH candidate and requiring a
regular executable file, which silently rejects every symlinked install: on a
Homebrew or Docker Desktop machine the tool is present and working while
Roundfix reports it missing. Discovery therefore resolves a bounded symlink
chain to its target and judges the target, while still never executing it, and
reports the inspected candidate and the reason it was rejected instead of an
empty result. The chain is bounded and a cycle, a broken link, or a
non-executable target each produce a distinct diagnostic rather than a generic
absence. Executing the candidate to prove it works was rejected because it
breaks the offline, side-effect-free guarantee that makes planning safe to
re-run; trusting the link without inspecting its target was rejected because a
dangling link would then read as a satisfied capability.
