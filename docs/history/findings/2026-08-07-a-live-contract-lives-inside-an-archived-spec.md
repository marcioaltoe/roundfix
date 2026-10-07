---
status: done
absorbed_by: 0083-a-gate-that-can-say-no
created_at: 2026-08-07
updated_at: 2026-09-08
kind: finding
---

# A live contract lives inside an archived Spec (2026-08-07)

`TestCoverageEquivalence` enforces its invariant against `docs/specs/_archived/0071-verification-cost/coverage-record.json` — an 86 KB file inside an archived Spec. Every legitimate test rename must rewrite it, and the repository forbids exactly that.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-07-a-live-contract-lives-inside-an-archived-spec.md`.
