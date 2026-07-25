---
task: task_16
spec: 0046-public-context-driven-baseline-command
status: completed
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

- [x] Run and settle the complete Go/Python compatibility matrix.
- [x] Remove Python runtime, fallback, duplicate assets, and obsolete tests.
- [x] Preserve the standalone corpus under Go-owned verification.
- [x] Tighten build, distribution, and skill-governance checks.
- [x] Verify the shipped command and thin skill expose one authority.

## Acceptance Criteria

- [x] Every maintained matrix row is passing before Python removal.
- [x] The shipped setup skill contains zero executable setup-engine scripts.
- [x] No command, document, test, or skill recipe references a Python fallback.
- [x] The Go test suite proves all exact, semantic, designed-delta, and ancillary destinations.
- [x] Skill sync rejects reintroduced executable setup-engine content.
- [x] The existing Setup Command retains its prior public behavior.
- [x] A source checkout with no Python runtime can build and exercise Baseline.

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

## Result

### Behavior delivered

- Proved the frozen 243-row compatibility matrix before cutover by running all
  256 maintained Python tests and the complete Go destination packages. The
  matrix retains 240 source-test rows plus three product-delta rows with 163
  exact, 63 semantic, 16 designed-delta, one ancillary, and zero retired
  classifications.
- Removed the canonical and distributed Python setup engines, Python tests,
  obsolete Python-only verification, skill-local references, and duplicate
  runtime assets. Each shipped `setup-context-driven` skill now contains only
  `SKILL.md`.
- Moved the standalone 17-file parity corpus to Go-owned
  `internal/baseline/testdata/parity-corpus/v1` and added an embedded Go gate
  that validates every artifact digest, row, classification, Go destination,
  fixture, adoption state, profile, Plan Digest, and content-addressed blob.
- Replaced the Makefile's Python verification with Go-owned skill-sync,
  distribution, lock-parity, and executable-engine rejection checks. The
  shipped skill validator now rejects any non-guidance file under
  `setup-context-driven`.
- Added a named compatibility test around the unchanged Setup Command's
  healthy-machine, idempotent, exact-profile-proof behavior.

### Acceptance evidence

1. Before removal,
   `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test*.py'`
   passed all 256 maintained tests, and
   `rtk env GOCACHE=/private/tmp/roundfix-task16-go-cache go test -count=1 ./internal/baseline ./internal/baselineacp ./internal/cli ./skills`
   passed every Go destination package.
2. `TestNoPythonBaselineRuntime` and `TestThinSetupSkill` passed and verified
   that the canonical, distributed, and embedded setup skill each contain only
   `SKILL.md`.
3. `TestNoPythonBaselineRuntime` also passed against the Makefile, public user
   documentation, Roundfix skill, and setup skill, proving the shipped
   surfaces contain no executable Python recipe or Python fallback.
4. `TestBaselineCompatibilityCorpus` passed with all exact, semantic,
   designed-delta, and ancillary rows assigned to 26 Go destinations; the full
   Go gate passed 2,013 tests across 22 packages.
5. `TestCheckRejectsExecutableSetupEngineArtifacts` proves a reintroduced
   setup-engine file produces a shipped-skill diagnostic.
   `skills-sync-check` runs this negative guard together with filesystem,
   embedded-skill, and `skills-lock.json` parity checks.
6. `TestSetupCommandCompatibility` passed the pre-existing Setup Command's
   healthy-machine output order, idempotence, exact profile proofs, zero-write
   behavior, and stdout/stderr contract.
7. The source binary built and
   `rtk env -i PATH=/nonexistent /private/tmp/roundfix-task16 baseline profile show rust-cli --format json`
   exited `0` with a valid embedded `roundfix/baseline-result/v1` document,
   proving Baseline can run without resolving Python or another helper.

### Verification

- `rtk env GOCACHE=/private/tmp/roundfix-task16-go-cache go test -count=1 ./internal/baseline ./internal/baselineacp ./internal/cli ./skills -run 'TestBaselineCompatibilityCorpus|TestNoPythonBaselineRuntime|TestThinSetupSkill|TestSetupCommandCompatibility'`
  passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task16-go-cache make skills-sync-check`
  passed four Go-owned cutover guards. The first sandboxed run without the
  explicit cache could not access the host Go build cache; the unchanged
  rerun used the allowed task cache.
- `rtk env GOCACHE=/private/tmp/roundfix-task16-go-cache make verify` passed
  2,013 Go tests, skill sync, shipped skill validation, formatting, and the
  binary build without invoking the removed Python runtime.
- `rtk git diff --check` passed.

No follow-up work was discovered inside this Task's slice.
