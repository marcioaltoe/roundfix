---
type: fix
status: promoted
created: 2026-09-30
spec: 0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need
reason: null
---

# An installed skill that matches its lock can trail the snapshot

## Symptom

`roundfix baseline update` restores a required upstream skill only when it is missing, when its lock entry is missing, or when its installed bytes differ from the lock entry's `computedHash`. It never compares the installed bytes with the `treeDigest` the Setup Snapshot pins. A lock written by another skills tool, which records `ref: main`, therefore keeps an installed skill at whatever `main` held on the day it was installed, and Doctor reports `skills: ok`.

## Where

`skills/repository.go` (`checkRepositoryWithReadiness`, external skills) and `internal/cli/baseline_update.go` (`baselineUpdateExternalDrift`).

## Expected

A required upstream skill whose installed tree differs from the snapshot's `treeDigest` is reported, and the managed refresh restores it to the snapshot's commit, or the difference is shown to be intended.

## Evidence

Measured on 2026-09-30 in this repository against `marcioaltoe/skills` at `a4e18e4`: ten required upstream skills (`bubbletea`, `golang-concurrency`, `golang-context`, `golang-error-handling`, `golang-lint`, `golang-testing`, `no-workarounds`, `systematic-debugging`, `testing-boss`, `tui-design`) match their lock entries and differ from `a4e18e4`, across 40 files. Spec 0200 refreshes the snapshot and restores only what the public update restores, so these ten stay as they are.
