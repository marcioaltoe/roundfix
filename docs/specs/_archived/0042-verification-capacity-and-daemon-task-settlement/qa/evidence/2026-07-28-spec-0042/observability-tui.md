# Live Run View and observability evidence

Live built-CLI PTY Run:
`run_20260728T141606Z_09f53ea96777c352`.

At 100×30 the interactive Live Run View showed:

```text
[verify] Verifying
[queued] Waiting for Verification
```

simultaneously for different Tasks while the aggregate Run phase was
Verifying. It also showed Task Capacity 2 and Verification Capacity 1. After
capacity release, both rows settled `Completed`; pressing `q` restored the
terminal and returned the Clean summary.

Fresh public events and Attach replay retained waiting-before-started order,
both capacities, terminal states, and read-only behavior.

Current-build TUI checks passed for 88/100/120-column stable layout,
too-short-terminal fallback, interleaved Task phases, terminal precedence,
legacy Attach fallback, styled/no-color parity, and review-Run regression.

The native macOS outcome notification warned on stderr in this managed
session. The same canonical `osascript -e 'display notification ...'` command
failed directly with the identical syntax diagnostic before Roundfix was
involved, so QA records it as an environment parity deviation rather than a
Spec-0042 runtime failure.
