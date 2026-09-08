---
status: proposed
granted: null
action: restore verification and regeneration compatibility for delivery of the requested documentation cleanup
consuming: 0120-knowledge-lifecycle-and-durable-capture
paths:
  - internal/baseline/derived_ownership_test.go
  - internal/speccheck/mechanical_test.go
  - internal/speccheck/governed_repocontract_test.go
  - internal/suiteguardcontract/regeneration.go
  - internal/suiteguardcontract/regeneration_test.go
---

# Proposed bounded compatibility work

This five-file delivery proposal supersedes the earlier seven-file question.
It retains only compatibility work necessary to verify and deliver the current
cleanup. Standalone Finding closure and source-adoption detector features
remain with their owning Specs. This is not a grant and does not answer any
pending consumption or execution-limit question.

1. The regeneration fixture currently copies the deleted authorization
   directory. Preserve its required test inputs without recreating that
   documentation tree or weakening its derived-artifact assertions.
2. Historical authorization tests must continue testing real approval evidence;
   missing current files cannot silently turn their checks into successful
   no-op results. Preserve original grants through their observed Git history
   or faithful test evidence, with explicit unavailable-history outcomes.
3. `ReadSanctionedRegenerations` currently reads only the removed directory.
   Support operative Spec-contained grants and retained legacy inputs, without
   treating proposed, null, malformed or unrelated documents as permission.
   The new `regeneration_test.go` is proposed coverage for this reader.

No suppression, weakened assertion, fabricated absorber, restored
`docs/workflow/` directory, new dependency, credential or paid API use is
proposed. The broader Spec 0120 implementation remains separately scoped.
Any granted record must land before its consuming changes; existing red
prerequisites must be repaired before the canonical tooling commit is delivered.
