---
status: accepted
created_at: 2026-10-16T00:00:00Z
updated_at: 2026-10-16T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A test replaces an executable it rebuilds instead of writing into it

F6, the stdout subtest of the item binary fallback test, saw its
`migrate --check` child end by `SIGKILL` four times in full suites, always
after its build copied the fixture over the `bin/roundfix` the previous subtest
had run. On macOS the kernel kills a signed executable that was rewritten in
place while another process still maps it, so every test that rebuilds an
executable now writes a new file and renames it over the old path, and the
diagnosis names the sender from the macOS diagnostic report, which needs no
`sudo`.

## Context

ADR-0266 fixed the class of stray signals: "A Proven Signal is the only signal
Roundfix code or tests send". F6 still recurred after that Spec merged, on
2026-10-14, in the first attempt of a Task of Spec 0268. Reading the code on
2026-10-16 ruled out the other senders a test could reach: the subtest's
context is cancelled only after the subtest returns, the suite guard's mark is
unique to each test binary and its residue check runs only after the package's
tests, and the rule test of ADR-0266 holds every other signal to a proven
target.

Measured on the maintainer's machine (macOS 27.0.1, arm64) on 2026-10-16, with
the Go linker's ad hoc signature on every fixture:

- copying a binary over itself, or over a different binary, and running it, 300
  and 100 times in sequence: no kill;
- copying it over a copy that is still running and then running the new copy:
  killed 40 of 40;
- through the repository's own fixture in a `go test -overlay` probe that wrote
  nothing to the checkout, rebuilding `bin/roundfix` with `cp` while a held
  instance of it runs: the probe ends with the exact F6 notice, "was killed by
  signal SIGKILL (9), no core dump" and no output, 5 of 5; rebuilding by copy
  and rename under the same holder: 5 of 5 pass, with a new inode each time.

Each kill left a report under `~/Library/Logs/DiagnosticReports` with the
termination namespace `CODESIGNING`, the indicator "Taskgated Invalid
Signature" and the signal "SIGKILL (Code Signature Invalid)". A signal another
process sends leaves no such report, so the report tells the kernel from a
sender without root. Apple's "Updating Mac Software"
(https://developer.apple.com/documentation/security/updating-mac-software)
says the same: "macOS caches information about the code’s signature in the
kernel. It doesn’t flush that cache when you modify the file’s contents." and
"write the updated code to a temporary file and replace the existing file with
that temporary one".

Which process still mapped F6's `bin/roundfix` when its stdout subtest copied
over it is not proven: the previous probe had been waited on, and nothing else
in the test runs that file.

## Decision

**A test that rebuilds an executable replaces it.** A test, fixture or declared
build command in a test that puts a compiled executable at a path that may
already hold one writes the new bytes to a new file in the same directory and
renames that file over the path. It never copies or writes into the existing
file. A rename gives the new executable its own in-kernel signature cache while
any process still mapping the old file keeps the old one, so the result does
not depend on which process holds it.

**The diagnosis names the sender without `sudo`.** A recurrence of a child's
Signal Exit is attributed to the kernel's code-signing enforcement when a
diagnostic report for that child carries the `CODESIGNING` namespace, and to
another process when no such report exists. Only the second case needs
`eslogger`, which needs `sudo` and Full Disk Access, and so the maintainer.

**Acceptance stays ten clean full suites.** ADR-0266's acceptance of F6 is
unchanged: ten consecutive full suites at the audited head with no child
killed by a signal.

## Rejected

- **Removing the old file before copying.** It also gives a new file, but it
  leaves a window in which the path is missing; a rename replaces it at once.
- **Waiting for whoever maps the old file.** The holder is not known, and a
  wait would be a timing allowance ADR-0271 forbids.
- **Retrying a probe killed by a signal.** It hides the kill instead of
  removing it.
- **Tracing every signal with `eslogger` first.** It needs `sudo`; the
  diagnostic report answers the kernel case without it.

## Consequences

- The item binary builds of the delivery tests copy to a temporary name in
  `bin/` and rename it over `bin/roundfix`, and a test holds a running copy
  while it rebuilds, so a return to copying in place fails that test on macOS
  every time.
- A user's `delivery.item_binary` build that writes into its output in place
  can meet the same kill when an earlier item binary still runs; the Delivery
  Queue does not change here.
