# QA Report contract

The closed report is
`docs/specs/0076-force-stop-exit-proof/qa/qa-report-2026-08-04.md`, the first
report for 2026-08-04.

Expected contract at close:

- `status: closed`;
- `verdict: fail`, mechanically required by one failed row;
- `rows_blocked_environment: 0`;
- `rows_blocked_finding: 0`;
- 13 terminal rows: 12 pass, 1 fail, 0 blocked, 0 skipped, 0 pending;
- direct helper liveness evidence at `helper-liveness.md`;
- independently observed premature-exit failure at
  `premature-exit-mutation.md`;
- every matrix evidence path resolves under this directory or to the report.

The Daemon retains ownership of `task_03` status and the report commit.
