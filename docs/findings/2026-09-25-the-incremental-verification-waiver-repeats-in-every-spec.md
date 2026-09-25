---
status: pending
created_at: 2026-09-25
updated_at: 2026-09-25
---

# Triage — The incremental verification waiver repeats in every consumer Spec (2026-09-25)

Triaged from secondbrain `inbox/roundfix/2026-09-17-a-dispensa-do-comando-incremental-virou-paragrafo-em-dez-specs.md` on 2026-09-25 and revalidated against main 7a9b6ec6.
The Profile declares the incremental tier but guides and manifests never publish it, so each adopter Spec hand-writes a waiver (19 copies across 10 Fiscus Specs). Owned by Spec 0121 Core Feature 6, which needs two more governed test files authorized.

## 1. The incremental verification waiver repeats in every consumer Spec

- Symptom / evidence: Fiscus carries 19 copies of the waiver across 10 Specs; `internal/baseline/assets/modules/core.json` keeps the clause mandatory while only `standard-typescript-monorepo.json` declares `verification.incremental`.
- Root cause: The publication half and the single-gate manifest migration stay with Spec 0121 Core Feature 6; Spec 0163 excluded it because two governed test files (`internal/cli/baseline_plan_test.go`, `internal/cli/baseline_release_gate_test.go`) are outside the approved authority.
- Action / suggestion: Route: Spec 0121 Core Feature 6, pending the maintainer's authorization of those two files.
