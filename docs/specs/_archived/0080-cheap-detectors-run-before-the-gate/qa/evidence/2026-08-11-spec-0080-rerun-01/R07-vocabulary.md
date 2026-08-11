# R07 — Mechanical Refusal Code vocabulary

Build: `c2372a9f709c9197aa5c5e89fd71da1ab46f07e6`.

- The built-from-current-source public command
  `rtk env GOCACHE=/private/tmp/roundfix-spec0080-rerun-gocache go run
  -buildvcs=false ./cmd/roundfix spec check
  0080-cheap-detectors-run-before-the-gate --strict` exited 0 with
  `No findings.` No detector was listed as skipped.
- The TechSpec Vocabulary Contract maps emission from
  `internal/speccheck/mechanical.go` through pattern
  `QA-(AUTH-PATHS|CONSEQUENT-ORDER|REPORT-SHAPE|EVIDENCE-PATH)` to
  `CONTEXT.md`.
- `internal/speccheck/mechanical.go:28-31` emits the four stable codes.
  `CONTEXT.md:202-205` owns them under `Mechanical Refusal Code`, explains
  each meaning, and marks the token stable.
- `rtk env ... go test ./internal/speccheck -run 'Vocabulary' -count=1 -v`
  exited 0. It passed documented, undocumented, repeated-token, invalid-pattern,
  unreadable-path, absent-contract, and fix-shape cases.

The prior report's F-002 is closed on the current build: durable ownership now
exists and strict checking verifies it rather than skipping the detector.
