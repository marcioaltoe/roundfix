# R05 — detector citations and skips

Build: `c2372a9f709c9197aa5c5e89fd71da1ab46f07e6`.

`rtk env GOCACHE=/private/tmp/roundfix-spec0080-rerun-gocache go test
./internal/speccheck -run
'^(TestMechanicalAuthPaths|TestMechanicalConsequentOrder|TestMechanicalReportShape|TestMechanicalEvidencePath|TestMechanicalReportsAllFindings|TestMaterializeMechanicalResult)$'
-count=1 -v` exited 0.

The selection passed exact bounded-path acceptance, every outside path in one
result, consequent fix after/before/folded cases, terminal and malformed report
shape, evidence-path resolution, multi-finding aggregation, and typed
materialization. Presence-aware cases recorded detector plus missing
authorization, consequent declaration, or report instead of inventing a
failure.

`TestMechanicalReportsAllFindings` verifies every emitted finding has a
non-empty stable code, file, positive line, detail, fix, and row hint where
known. `rtk env ... go test ./internal/speccheck -run
'^TestRenderResultTextAndJSON$' -count=1 -v` also exited 0, proving both public
result renderings preserve the structured facts.

The same detector selection was re-run in a fresh process after the earlier
R04 corpus run, satisfying the persistence check.
