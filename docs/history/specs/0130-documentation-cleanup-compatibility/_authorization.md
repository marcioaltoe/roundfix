---
status: approved
granted: 2026-09-09
action: restore verification and regeneration compatibility for the requested documentation cleanup
consuming: 0130-documentation-cleanup-compatibility
paths:
  - internal/baseline/derived_ownership_test.go
  - internal/speccheck/mechanical_test.go
  - internal/speccheck/governed_repocontract_test.go
  - internal/suiteguardcontract/regeneration.go
  - internal/suiteguardcontract/regeneration_test.go
---

# Approved bounded cleanup compatibility repairs

On 2026-09-09 the maintainer answered "Autorizo os reparos" to the explicit
request approving these five paths to unblock PR #180. This grants the bounded
proposal previously recorded here; the earlier seven-file proposal is superseded.
Spec 0130 is the executable repair slice of the broader Spec 0120 plan.

1. Preserve the regeneration fixture's required inputs without recreating the
   removed documentation tree in the delivered checkout or weakening assertions.
2. Preserve real historical authorization evidence through Git objects or faithful
   isolated test evidence, with explicit outcomes for unavailable history and
   non-vacuous controlled coverage.
3. Discover operative Spec-contained regeneration grants and retained legacy
   inputs, rejecting proposed, null, malformed or unrelated records as permission.
   Resolve command-only grants through the existing output ownership declarations;
   compare the result with the Baseline resolver, without a second exemption list.

The allowed implementation paths are exactly those above, plus each assigned
Task's own evidence file under the standard Task contract. Supervisor-authored
Spec and QA evidence are ordinary documentation. This grants no new dependency,
credential access, paid API use, unrelated feature, verification bypass or change
to test-runner configuration. The broader Spec 0120 remains separately scoped.

This authorization must land in main before the consuming squash merge. The
already approved canonical policies and their sanctioned digest regeneration
remain covered by the grants merged in PR #179.
