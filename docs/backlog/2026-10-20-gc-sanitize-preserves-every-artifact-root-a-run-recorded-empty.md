---
type: fix
status: open
created: 2026-10-20
spec: null
reason: null
---

# GC sanitize preserves every Artifact Root a Run recorded empty

## Symptom

On 2026-10-20 `roundfix gc sanitize` (dry run) on the maintainer's Roundfix
Home listed one Artifact Root, the empty string, classified `unsafe`:
"preserved because recorded Artifact Root "" is not a clean absolute path".
Every Implement Run records an empty root, so sanitation never reaches the
default roots under `~/.roundfix/artifacts/<repository id>/`, which held about
605 MiB that day. The preservation is safe, but the command reports nothing it
could reclaim there.

## Where

`internal/store/journal.go`, `DiscoverArtifactRoots`, passes the recorded
`artifact_dir` as is, and `internal/cli/gc.go`, `classifyGCSanitationRoot`,
refuses any root that is not a clean absolute path. Run Retention
(`internal/cli/run_retention.go`, `runRetentionArtifact`) already resolves an
empty recorded root to the repository key's or the recorded checkout's default
root.

## Expected

Sanitation resolves an empty recorded root the way the Run wrote to it, then
classifies the resolved root with its existing proofs. Removing a directory
stays limited to what those proofs establish.

## Evidence

The dry run above, observed while authoring Spec
0287-reports-that-count-what-they-say, which fixes the same reading in
`roundfix storage report` only and leaves sanitation's deletion behavior to a
Spec of its own; no Finding.
