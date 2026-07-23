---
task: task_16
spec: 0046-public-context-driven-baseline-command
status: pending
type: backend
complexity: high
---

# Task 16: Cut over to the Go Baseline authority

## Overview

Complete the migration only after every maintained behavior has Go parity and
public documentation. The shipped CLI and thin skill then have one runtime
authority, while Python scripts, fallback paths, duplicate assets, and obsolete
gates are removed.

## Requirements

1. MUST prove every compatibility-matrix row against its Go destination before
   deleting any Python implementation or test.
2. MUST remove the executable Python setup engine, Python fallback paths,
   duplicate runtime assets, and obsolete Python-only verification after
   parity passes.
3. MUST retain the standalone compatibility corpus and all Go regression
   coverage required to prevent silent contract loss.
4. MUST make the Go embedded catalog, public Baseline command family, and thin
   setup skill the only shipped execution path.
5. MUST update build, skill-sync, distribution, and documentation checks to
   fail if an executable setup engine or divergent runtime asset reappears.
6. MUST leave unrelated Setup Command behavior unchanged.

## Subtasks

- [ ] Run and settle the complete Go/Python compatibility matrix.
- [ ] Remove Python runtime, fallback, duplicate assets, and obsolete tests.
- [ ] Preserve the standalone corpus under Go-owned verification.
- [ ] Tighten build, distribution, and skill-governance checks.
- [ ] Verify the shipped command and thin skill expose one authority.

## Acceptance Criteria

- [ ] Every maintained matrix row is passing before Python removal.
- [ ] The shipped setup skill contains zero executable setup-engine scripts.
- [ ] No command, document, test, or skill recipe references a Python fallback.
- [ ] The Go test suite proves all exact, semantic, designed-delta, and ancillary destinations.
- [ ] Skill sync rejects reintroduced executable setup-engine content.
- [ ] The existing Setup Command retains its prior public behavior.
- [ ] A source checkout with no Python runtime can build and exercise Baseline.

## Context

- instruction: `docs/adr/0066-context-driven-baseline-execution-belongs-to-the-cli.md`
- instruction: `docs/adr/0072-baseline-go-cutover-preserves-python-contracts.md`
- instruction: `docs/agents/skill-governance.md`
- interface: `Makefile`
- interface: `.agents/skills/setup-context-driven`
- interface: `skills/setup-context-driven`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/baselineacp ./internal/cli ./skills -run 'TestBaselineCompatibilityCorpus|TestNoPythonBaselineRuntime|TestThinSetupSkill|TestSetupCommandCompatibility'` — expected: maintained parity, Python absence, skill authority, and unchanged Setup Command pass.
- `rtk make skills-sync-check` — expected: distributed skill guidance matches the canonical thin skill and contains no executable setup engine.
- `rtk make verify` — expected: the post-cutover gate passes without invoking the removed Python runtime.

## References

- `_prd.md` → Goals 1, 4–5; User Story 9; Core Features 19–21; Non-Goals / Out of Scope; Success Metrics.
- `_techspec.md` → System Architecture; Testing Approach; Build Order 10.
- ADR-0066 → Go CLI runtime authority.
- ADR-0072 → parity-gated Python removal.
