# Live Run View, documentation, and Skills

The current built CLI attached read-only to the active Spec Run at 100×30 with
`ROUNDFIX_COLOR=never`. The Live Run View rendered all eight Tasks as
Completed, stable two-pane geometry, text status markers, the Phase Row, Work
Queue, Session Timeline, and keyboard help. Keyboard-only Tab, Down, Enter,
Esc, End, and `q` opened task_02 detail, closed it, returned to follow, and
detached without stopping or mutating the Run.

Fresh synchronous model tests passed for 88, 100, and 120-column supported
widths, degenerate truncation, interleaved Agent working / Waiting for
Verification / Verifying phases, terminal precedence, styled/no-color text
parity, focus, detail, viewport, and review-Run compatibility. The current
built non-interactive Attach replay independently rendered both capacity
labels from stored Run evidence.

`TestCommandUsage`, `TestDocumentationContract`, and current CLI help passed.
The documented 2:1 Config, per-Run scope, event filters, Attach flow, exit-75
protocol, one repair, Daemon-owned status, and Task Type-selected Agent
Session wording match the runtime.

Both canonical/generated Skill pairs compare byte-identically. With a writable
Go cache, `make skills-sync-check` passed four policy tests and
`roundfix skills check` passed all fourteen shipped Skills.
