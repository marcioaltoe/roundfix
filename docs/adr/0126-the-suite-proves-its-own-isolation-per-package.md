---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-14T21:05:30Z
updated_at: 2026-08-14T21:05:30Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# The suite proves its own isolation per package, not once around itself

A no-write invariant checked once around `go test ./...` can say that something
wrote into the repository but never which package did, and the failure surfaces
after every package has finished — far from the code that caused it. The guard
therefore runs per package, from a shared helper each package's `TestMain`
installs, so a violation fails the package that produced it and names the path it
wrote. The cost is that a package without the helper is unguarded, which a
repository-contract test makes visible by enumerating the packages that install
it; that enumeration is a fact a reader can check, where a suite-level pass is a
fact nobody can attribute.
