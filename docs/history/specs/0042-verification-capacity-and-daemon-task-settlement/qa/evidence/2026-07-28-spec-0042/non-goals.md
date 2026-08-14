# Non-Goal and regression evidence

Current-build public, focused, full, and race checks preserve:

- the complete Task Verification command sequence and existing timeouts;
- exit `75` as the sole explicit temporary signal;
- no log-text, framework, port, package, database, listener, or timeout
  heuristic;
- one Task-scoped exclusive retry and one independent deterministic Agent
  repair;
- no configurable retry count, backoff, or new terminal Run outcome;
- per-Implement-Run capacity only, with no machine-wide or cross-Run claim;
- unchanged review Batch Verification;
- failed Task Worktree retention, dependency blocking, Settle recovery, and
  Run integration.

Focused `Test(Settle|Implement).*TaskStatus`, Agent-failure/Stop/status
matrices, public help/config inspection, the full 2,727-test gate, and the full
race suite all passed.
