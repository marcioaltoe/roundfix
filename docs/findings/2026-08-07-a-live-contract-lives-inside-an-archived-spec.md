---
status: deferred
created_at: 2026-08-07
updated_at: 2026-08-26
kind: finding
---

# A live contract lives inside an archived Spec (2026-08-07)

`TestCoverageEquivalence` enforces its invariant against
`docs/specs/_archived/0071-verification-cost/coverage-record.json` — an 86 KB
file inside an archived Spec. Every legitimate test rename must rewrite it, and
the repository forbids exactly that.

## The collision

Two repository rules apply to the same file and point opposite ways:

- *Keep completed or archived legacy Specs byte-identical.* Updating the record
  edits an archived Spec.
- *Durable knowledge flows upstream only — an archived Spec may be deleted at
  any time, so durable knowledge a Spec produced moves to its semantic owner
  before or at archive.* Leaving the record there makes the whole coverage
  contract deletable by an ordinary archive cleanup.

The test itself expects the file to be rewritten: it ships an
`-update-coverage-record` flag and its read error tells the maintainer to
"regenerate deliberately" with it. So the sanctioned repair path is the one the
archival rule forbids.

The test already anticipates a non-archived home — it falls back to
`docs/specs/0071-verification-cost/coverage-record.json` when the archived path
is absent — but no active copy exists, so the archived file is the live
contract.

## How it surfaced

Spec 0082's task_06 replaced two CLI test functions with new ones. The record
still names the old identities, so the invariant is red:

```
coverage regression: package "roundfix/internal/cli" no longer executes
  "TestHumanBaselineIncompatibleManifestKeepsValidDefaults"
  "TestHumanBaselineUpdate"
```

The failure is correct — the record and the suite genuinely disagree. What has
no sanctioned answer is where to record the agreement.

## Why it matters

Any Spec that renames or replaces a test hits this, which is most Specs. The
maintainer is left choosing between a red gate and an edit the rules forbid,
and neither choice is recorded anywhere as the intended one. That is a
governance gap, not a test defect.

## Route

Not fixed here, because the fix is a decision rather than an edit: move the
coverage record to a semantic owner outside `docs/specs/` and point both paths
at it, or declare this file an explicit, documented exception to the
archived-Spec rule. Rewriting it in place without settling that would silently
establish the exception by precedent.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
