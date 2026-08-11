# R05 — detector citations and skips

`rtk ... go test ./internal/speccheck -run '^(TestMechanicalAuthPaths|TestMechanicalConsequentOrder|TestMechanicalReportShape|TestMechanicalEvidencePath|TestMechanicalReportsAllFindings|TestMaterializeMechanicalResult|TestMechanicalCorpusNonRegression|TestCheckCorpusBudget)$' -count=1 -v`
passed 18 focused cases.

Observed detector codes at `internal/speccheck/mechanical.go:28-31`:

- `QA-AUTH-PATHS`
- `QA-CONSEQUENT-ORDER`
- `QA-REPORT-SHAPE`
- `QA-EVIDENCE-PATH`

The multi-finding carrier returned all four in one result. Each finding
asserted code, file, positive line, detail, fix, and the known row hint.
Presence-aware carriers returned detector plus missing artifact without a
finding. Materialization preserved the fields in report sections.
