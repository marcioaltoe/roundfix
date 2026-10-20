---
type: fix
status: promoted
created: 2026-10-20
spec: 0287-reports-that-count-what-they-say
reason: null
---

# The storage report misses the default Artifact Roots

## Symptom

On 2026-10-20 `roundfix storage report` on the maintainer's machine listed one
Artifact Root, `.`, as `missing` with 508 Run rows and 0 bytes, `Measured
bytes: 0`, and 0 artifact bytes for every repository and every Run state. This
repository's `.roundfixrc.yml` sets `defaults.artifact_dir: ""`, so every Run
wrote its artifacts under the default root of its repository,
`~/.roundfix/artifacts/<repository id>/`. That tree held 37 repository
directories and about 605 MiB (`du -sk` reported 619,428 KiB), with 336
`runs/<run-id>` directories.

## Where

- `internal/store/journal.go`, `storageRunRecords` and
  `storageFilesystemGroups`: each Run's recorded `artifact_dir` is cleaned with
  `filepath.Clean`, so an empty value becomes `.`, and `storageDirectoryBytes`
  and `storageRunArtifactBytes` treat `.` as absent. A relative value is
  measured against the process's working directory.
- `internal/cli/gc.go`, `printStorageReport` (`:349` onward), prints the group.
- The recorded value is empty for every Implement Run: `runImplement` resolves
  the directory with `ValidateArtifactDirectory` and writes there, but
  `store.CreateRunRequest` carries no `ArtifactDir`. Review Runs record the
  configured value.
- `internal/store/journal.go`, `DiscoverArtifactRoots` (`:1583`), gives `roundfix
  gc sanitize` the same recorded value: its dry run on 2026-10-20 printed
  `Artifact Root: ` classified `unsafe`, "preserved because recorded Artifact
  Root "" is not a clean absolute path".
- Run Retention already resolves an empty recorded root
  (`internal/cli/run_retention.go`, `runRetentionArtifact`): the repository
  key's default root, or the recorded checkout's default root.

## Expected

The storage report measures each Run's artifacts where the Run wrote them: an
empty recorded root resolves to its repository's default root, a relative one
resolves against the Run's Git root, and an absolute one stays as recorded. The
report stays read-only and creates no directory.

## Evidence

`roundfix storage report` and `roundfix gc sanitize` (dry run) on the
maintainer's Roundfix Home, 2026-10-20, read-only; no Finding.
