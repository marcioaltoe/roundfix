---
status: accepted
knowledge: true
created_at: 2026-11-24T23:35:15Z
updated_at: 2026-11-24T23:35:15Z
deprecated_at: null
superseded_by: null
---

# Roundfix verifies on macOS, and the Delivery Queue has no Linux Lane

Supersedes ADR-0288, ADR-0294, ADR-0297, ADR-0303 and ADR-0316 (2026-11-24); all are retired, and Git keeps their text. Their creation dates are respectively 2026-10-18, 2026-10-19, 2026-10-20, 2026-10-20, 2026-10-22.

Amends ADR-0252 (2026-10-07), ADR-0285 (2026-10-17), ADR-0295 (2026-10-19), ADR-0384 (2026-11-10), ADR-0385 (2026-11-11) and ADR-0423 (2026-11-17).

The maintainer retired Linux verification because Roundfix is used on macOS and the container gate had repeatedly blocked delivery. Continuous integration retains verification before merge.

## Decision

The Delivery Queue repository gate takes one Machine Verification Capacity unit, computes the Verification Core Budget once and supplies it to the repository Verification, releases the unit afterwards, and passes when that Verification passes. A failure parks the item as before. The gate neither reads default-branch configuration for a lane nor resolves Docker, writes a lane log, adds a warning, or adds a pull request paragraph; the engine advances to publishing and preserves existing warnings.

`roundfix deliver lane` is no longer a command: without help it is refused as an unknown command on standard error with exit 2, and with help it prints deliver usage and exits 0. Deliver usage lists plan, start, status, resume, retry and stop.

The keys `delivery.linux_lane` and `delivery.linux_lane_image` are deprecated and ignored, whatever their value. They are removed before validation, each produces its deprecation warning once per load in that order, and a key present in both User and Project Config warns once. The warning is `config: <key> is deprecated and ignored; the Linux Lane was retired, so remove it`. A delivery mapping with only those keys is empty; a configuration read at a commit removes them silently. Other unknown delivery keys remain errors. Doctor no longer runs or reports the lane check or its diagnostic codes; it prints none of `DR-LANE-IMAGE-UNKNOWN`, `DR-LANE-DOCKER-MISSING`, `DR-LANE-DOCKER-UNREACHABLE` and `DR-LANE-IMAGE-MISSING`. Its `verification-cores` line is followed by `codex`, and its other checks retain their order and text.

A stored `linux-lane-container:` blocker keeps its environment Park Class and directs the operator to `roundfix deliver retry <slug>` because the lane was retired. A stored `gate-failed: Linux Lane change failure:` blocker remains a finding with the normal gate-failed action. The retry counts against the item retry limit. Status prints a stored `linux-lane-skipped:` warning unchanged; no new code reads, writes or removes that segment. Low-disk warning updates preserve the other segments and place the low-disk segment last. Existing lane logs are neither read nor deleted.

The CI Verification gate runs its existing steps on `macos-26`. Release verification and the Full Contract Run execute in a `macos-26` job with read-only contents permission, full-history checkout and no persisted credentials; the unchanged publication job needs that job to succeed. The verification job uses the existing Go setup and the existing Verify gate and Verify every contract steps for both tag pushes and manual dispatches. The publication job retains its runner, permissions, environment and step order, including checkout, Go and Node setup, runtime guard, tag validation, preflight, cross-compilation, package comparison, publication and GitHub Release. Tag validation therefore follows verification. The test compile contract runs `go vet -tags "repocontract docscontract" ./...` for Linux and Windows without executing tests on either platform, and detects a Linux-only test compile error. The selective gate and continuous integration run on macOS; failures only Linux or Windows produce remain a known limit. The suite-time budget stays 420 seconds until the suite is measured on the macOS runner.

A test that needs an External Test Tool skips naming it when absent. A tool outside the Base Test Tools is resolved through `testfixture.RequireTool`; the rule in `internal/testfixture` refuses other ways.

## Considered options

A light Linux build-and-vet job was rejected because the maintainer chose to retire Linux verification. No CI tests was rejected because the maintainer wanted a check before merge. Running the whole release on macOS was rejected to preserve the proven publication job and its permissions. Three release jobs validating the tag first were rejected in favor of the bounded two-job change; the trade-off is that a bad tag is refused after verification. The moving `macos-latest` label was rejected so an image change is reviewed; a larger runner was rejected without a measurement. Classifying a stored blocker as `unclassified` was rejected because its generic action would tell the operator to repair a retired gate.

## Consequences

macOS runs the tests; Linux and Windows builds and tests are compiled but not executed. Existing parked items remain actionable through Delivery Retry, and existing warning text remains readable.
