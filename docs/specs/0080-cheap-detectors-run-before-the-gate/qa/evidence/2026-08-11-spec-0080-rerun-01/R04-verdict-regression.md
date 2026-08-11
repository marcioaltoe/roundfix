# R04 — verdict and corpus non-regression

Build: `c2372a9f709c9197aa5c5e89fd71da1ab46f07e6`.

- `rtk env GOCACHE=/private/tmp/roundfix-spec0080-rerun-gocache go test
  ./internal/cli -run '^TestRunImplementQAVerdictMatrix$' -count=1 -v`
  exited 0. The public Implement CLI journey passed all five cases: `pass`,
  `partial`, `fail`, `missing_report`, and `unreadable_verdict`. This directly
  rechecks corrective Task 09 and closes prior F-001.
- `rtk env ... go test ./internal/spec -run '^TestQAVerdict' -count=1 -v`
  exited 0. Supported verdicts, newest and same-day rerun selection, missing
  and unreadable reports, and typed environment/finding/declared blocked-count
  semantics all passed.
- `rtk env ... go test ./internal/speccheck -run
  '^(TestMechanicalCorpusNonRegression|TestMechanicalReportShape|TestMaterializeMechanicalResult)$'
  -count=1 -v` exited 0. The corpus stayed non-regressive; green, malformed,
  and absent reports retained their established behavior; materialization
  preserved the verdict-free typed result boundary.

Independent source inspection found no Task 09 change to
`internal/spec/qa.go`. The correction is confined to CLI/Daemon test harness
report completion, so production verdict computation did not move.
