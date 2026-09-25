---
status: done
created_at: 2026-09-15
updated_at: 2026-09-25
absorbed_by: 0171-a-deterministic-suite-under-load
---

# CLI tests — Operational command tests reach GitHub through the version freshness check (2026-09-15)

This came up while diagnosing why Spec 0138's QA gate could not pass `make verify`
on this machine. The Daemon was not the problem. The QA Agent runs the repository
Verification inside its sandbox, and some tests there need the network. Spec 0139
moves that Verification out of the sandbox. That hides the symptom but leaves
these tests non-hermetic, so the cause is recorded here.

## 1. Tests of operational commands call the real release lookup

- **Symptom / evidence:**
  - `maybeReportVersionFreshness` in `internal/cli/upgrade.go` skips only for a
    development version or an empty home directory. Tests run the checked-in
    `0.12.0` version, so the skip does not apply to them.
  - When the daily cache is stale, the check calls `app.LatestRelease`, which
    resolves the latest release through GitHub.
  - Only `internal/cli/upgrade_test.go` injects fake freshness dependencies,
    through `withVersionFreshnessFakeDeps`. The package-level
    `versionFreshnessDeps` default stays live for every other test that runs
    `fetch`, `resolve`, `watch` or `implement`.
  - Spec 0136's QA Report (`qa-report-2026-09-14-01.md`, Static gate) records
    that `make verify` passed only after "one unchanged authorized retry after
    the sandbox stopped the first process at `api.github.com`". Its CLI-focused
    command "required one unchanged network-enabled retry at its existing GitHub
    boundary".
- **Root cause:** code under test reads a network boundary from the process
  default instead of taking it explicitly. ADR-0089 requires explicit
  environment for exactly this case.
- **Action / suggestion:**
  - Route to a Spec that makes operational command tests inject a
    non-networked release lookup, or disables the freshness check through
    explicit command dependencies.
  - A mutable package-level variable must not be the switch, because parallel
    tests share it.
  - Spec 0139 does not change this. Once the Daemon runs the repository
    Verification outside the sandbox, the lookup can reach the network again,
    but the suite still depends on it.

## Addendum — 2026-09-25 — Revalidated in triage

Revalidated against main 7a9b6ec6: still holds: the default versionFreshnessDeps stays live for fetch, resolve, watch and implement tests; Spec 0165's QA hit api.github.com again. Ranked in the 2026-09-25 triage priority list.

## Addendum — 2026-09-25 — Adopted by Spec 0171

Adopted by Spec 0171-a-deterministic-suite-under-load. Its Task gives every
`internal/cli` test process a release lookup that never leaves the machine,
seeds a fresh version cache in the HOME of every built binary a test runs for an
operational command, and adds guard tests that fail when a live lookup could
happen.
