---
task: task_02
spec: 0096-a-failure-the-agent-can-read
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 02: A signature that survives a timestamp

## Overview

Two runs of the same failure differ by timestamps, temporary paths, durations and
process identifiers, so comparing captured diagnostics byte for byte calls every
repetition new. This slice builds the normalised form the repetition check
compares, and nothing else consumes it yet.

## Requirements

1. MUST produce a stable signature from a failing command and its captured
   diagnostic.
2. MUST compare equal for two diagnostics differing only in timestamps,
   durations, temporary directory paths, process identifiers, or run identifiers.
3. MUST compare unequal for two diagnostics whose assertions differ.
4. MUST include the failing command in the signature, so the same output from
   different commands is not one failure.
5. MUST NOT read or write any store; it is a pure function of its inputs.

## Subtasks

- [ ] Normalise the volatile spans.
- [ ] Hash command and normalised diagnostic together.
- [ ] Cover each volatile class and the genuinely-different case.

## Acceptance Criteria

- [ ] Diagnostics differing only by timestamp, duration, temporary path, process
      id, or run id produce one signature.
- [ ] Diagnostics with different assertions produce different signatures.
- [ ] The same diagnostic from two different commands produces different
      signatures.
- [ ] The function touches no store.

## Verification

- `go test -count=1 ./internal/daemon -run 'TestDiagnosticSignature' -v > /tmp/0096-t02.log 2>&1; s=$?; grep -q '^--- PASS: TestDiagnosticSignature' /tmp/0096-t02.log || { cat /tmp/0096-t02.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0096-t02.log || { echo 'the suite produced no output'; exit 1; }; grep -qi 'no tests to run' /tmp/0096-t02.log && { echo 'the suite selected no cases'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0096-t02.log > /tmp/0096-t02-n.txt; test "$(cat /tmp/0096-t02-n.txt)" -ge 6 || { echo "expected a case per volatile class plus the different-assertion and different-command cases, got $(cat /tmp/0096-t02-n.txt)"; cat /tmp/0096-t02.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving each normalisation class is exercised on its own.
- `f=$(grep -rl 'func DiagnosticSignature' internal/daemon --include='*.go' | grep -v '_test.go' | head -1); test -n "$f" || { echo 'no signature function exists'; exit 1; }; grep -qE 'store\.|Store\b' "$f" && { echo "FAIL: $f reads a store; the signature must be a pure function"; exit 1; }; exit 0` — expected: exits 0, proving the function exists and stays pure. Fails today.

## Context

- interface: `internal/daemon/task_engine.go`

## References

`_techspec.md` → Build Order 2; Implementation Design, Interfaces; Risks &
Considerations, the collision. `_prd.md` → Core Feature 2; Open Questions.
ADR-0136.

## Result

Implemented a pure `DiagnosticSignature` function that replaces timestamps,
durations, Unix and Windows temporary paths, process identifiers, and Run IDs
before hashing a length-delimited failing command with the normalized diagnostic.

Focused-check evidence:

- Before implementation, `rtk go test ./internal/daemon -run
  '^TestDiagnosticSignature$'` failed to compile because `DiagnosticSignature`
  did not exist.
- After implementation, the same focused command reported 8 passing tests: the
  parent plus separate timestamp, duration, temporary-path, process-id, run-id,
  different-assertion, and different-command subtests.
- `rtk go test ./internal/daemon` reported 233 passing tests.
- A source scan for `store\.` or `Store` in
  `internal/daemon/diagnostic_signature.go` returned no matches.

Acceptance evidence:

- `TestDiagnosticSignature/normalizes_timestamp`,
  `/normalizes_duration`, `/normalizes_temporary_directory_path`,
  `/normalizes_process_identifier`, and `/normalizes_run_identifier` compare
  each volatile class independently and produce equal signatures.
- `TestDiagnosticSignature/preserves_different_assertion` produces unequal
  signatures when the assertion changes.
- `TestDiagnosticSignature/includes_failing_command` produces unequal signatures
  for identical diagnostics from different commands.
- `DiagnosticSignature` depends only on its command and diagnostic parameters;
  its source imports hashing, encoding, regular-expression, and numeric-formatting
  standard-library packages and has no persistence dependency.

The Daemon-owned Verification commands were not run in this Agent turn.
