---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A composed profile takes a setup composed from upstream setups by name

A Baseline Profile names one Setup Snapshot, and each snapshot mirrors one
upstream skill list. No upstream list covers a repository that holds a Go
command-line tool beside a Bun workspace with a Hono backend and a React
frontend, so no built-in profile could select both the Go and the TypeScript
modules without failing the check that a profile's setup lists every skill its
modules name. The planned argus repository has that shape.

A built-in profile may now take a composed Setup Snapshot. A composed snapshot
names its component snapshots by identifier, holds no list of its own choosing,
and carries the union of their skills and activation bundles in component
order. The asset sync writes it after it refreshes the components, and the
catalog refuses a composed snapshot that differs from that union, names an
unknown or composed component, or would merge two different entries under one
skill name or one bundle. The first built-in composed profile is
`go-cli-typescript-monorepo`, on the composed setup `go-cli-typescript-bun`
(`go-cli` then `typescript-bun`). Its guides keep each stack rule scoped to the
language or workspace it governs, as ADR-0190 requires; the Go guide gains the
scope sentence it lacked. The profile declares which toolchain Verification its
root gate runs, and profile alignment reports, without blocking, a root Make
gate that does not reach one of them.

Two alternatives were rejected. A profile naming several setups would change
every reader of a profile's setup, including skill restore, reconcile and
Doctor. Deriving the composition when the catalog loads would leave a snapshot
file that the raw readers of setup files cannot decode.

## Consequences

A later composition is one snapshot declaration and one profile. The composed
snapshot changes only when a component changes, so ADR-0191's rule that a
snapshot follows its upstream by name holds for it through its components. The
gate check reads one Makefile textually: an included file, a shell script or a
differently spelled command is not followed, and the divergence stays advisory.
