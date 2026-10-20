---
type: fix
status: promoted
created: 2026-10-20
spec: 0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it
reason: null
---

# Runs causes reads no diagnostic of a Run that recorded no Artifact Root

## Symptom

`roundfix runs causes` shows no diagnostic tail for any Task attempt of an
Implement Run. Every Implement Run records an empty Artifact Directory while
it writes its diagnostics under `~/.roundfix/artifacts/<repository id>/runs/<run-id>/`,
and the report returns an empty tail whenever the recorded directory is empty.

## Where

`internal/runcause/report.go`, `diagnosticTail`, returns `""` when
`run.ArtifactDir == ""` and otherwise builds the Run's directory from the
recorded value literally. It reads the recorded root the way `roundfix gc
sanitize` and `roundfix storage report` read it before they resolved it.

## Expected

The report reads the diagnostic under the Artifact Root the Run wrote to,
through the same resolver the storage report, GC sanitation and Run Retention
use, and keeps its symbolic-link and containment checks.

## Evidence

Code reading while authoring Spec
0289-gc-sanitize-reaches-the-default-artifact-roots, which resolves the root
for the storage report, sanitation and Run Retention only; not reproduced
through the command yet.
