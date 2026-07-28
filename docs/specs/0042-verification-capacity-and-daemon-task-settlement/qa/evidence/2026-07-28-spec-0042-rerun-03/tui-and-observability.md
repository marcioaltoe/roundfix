# Live Run View and observability

Fresh synchronous TUI and Attach tests passed for both capacity labels, legacy
fallback, interleaved Agent working / Waiting for Verification / Verifying
phases, Verification Feedback, terminal precedence, styled/no-color parity,
review compatibility, keyboard focus, detail, and viewport behavior.

The built current CLI independently replayed stored Spec Run
`run_20260728T134451Z_ec12a53008910524`:

```text
roundfix attach run_20260728T134451Z_ec12a53008910524 --no-input
roundfix attach --no-input run_20260728T134451Z_ec12a53008910524
```

Both forms exited 0 and rendered the same read-only Unresolved Run with Task
Capacity 1 and Verification Capacity 1. A public `roundfix events` replay with
the documented `task-status,verification,outcome` filter exited 0 and emitted
`roundfix-events/v1` JSON.

A live `ROUNDFIX_COLOR=never` TUI pass exercised both an 80-column degenerate
view and a 100×30 two-pane view. At 100×30 it rendered all eight Tasks with
textual `Completed` markers, stable Work Queue and Session Timeline panes,
Phase Row, summary, and key guidance. Keyboard-only `d` opened task_01 detail,
`j` scrolled its source, Escape closed it, and `q` detached. The terminal
restored and the Run remained Unresolved; Attach did not mutate or stop it.

Current model tests cover supported 88, 100, and 120-column layouts, bounded
degenerate truncation, and the interleaved nonterminal phases that this
historical terminal Run cannot display live.
