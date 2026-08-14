# Public CLI contract

The current built root, Implement, and Attach help all exited 0. No Implement
capacity flag or new top-level exit code appears. Configuration remains the
only capacity surface.

The public CLI runner replayed all three supported Attach forms:

```text
roundfix attach <run-id> --no-input
roundfix attach --no-input <run-id>
roundfix attach <run-id>
```

The built binary independently ran the first two forms against stored Run
`run_20260728T134451Z_ec12a53008910524`; both exited 0 and rendered the same
read-only result. This confirms the historical trailing-flag finding remains
resolved.

The disposable Implement macro suite asserts exact requested stdout and exit
behavior for clean success, deterministic repair, temporary retry, repeated
temporary failure, and Stop. Run identity, progress, diagnostics, waiting,
retry, and event evidence remain on stderr or in the Run Event Journal.
Review Verification and Settle regression checks passed in the full current
suite.
