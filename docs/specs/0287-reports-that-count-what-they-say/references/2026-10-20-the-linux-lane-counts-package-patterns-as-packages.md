---
type: fix
status: promoted
created: 2026-10-20
spec: 0287-reports-that-count-what-they-say
reason: null
---

# The Linux Lane counts package patterns as packages

## Symptom

`roundfix deliver lane --all` printed `Linux Lane passed 1 package(s)` on
`main` although 43 packages ran in the container: 40 reported `ok` and 3
reported no test files. The number on the pass line is the number of package
arguments the host handed to the container, and `--all` hands it the single
pattern `./...`. The Delivery Queue gate prints the same line, so a gate that
tested a whole module reads as one package.

## Where

`internal/cli/deliver_linux_lane.go`, `runLinuxLane`: the lane log line
(`:239`) and the console line (`:243`) both print `len(result.Packages)`.
`result.Packages` is the pattern list, `./...` in `--all` mode or after a
`go.mod` or `go.sum` change, and the nearest package directories otherwise.
The packages that actually run are chosen inside the container by the lane
script: `go list -e -f '{{if or .GoFiles .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}'`
expands the patterns, and `go test` runs the result. The user guide
(`docs/user-guide/commands/deliver.md`) documents the defect as intended: "`<n>`
counts selected package arguments". The gate and `roundfix deliver lane` share
`runLinuxLane`, so both print the wrong number.

## Expected

The pass line counts the packages `go list` selected in the container, the
packages `go test` ran, in the gate and in `roundfix deliver lane` alike. When
the count is unknown, the line says so instead of printing a number.

## Evidence

The operator's on-demand lane run on `main`, reported in the cycle AK operator
briefing of 2026-10-20 (43 packages ran, the line said 1); no Finding.
