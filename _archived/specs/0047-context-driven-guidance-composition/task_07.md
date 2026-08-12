---
task: task_07
spec: 0047-context-driven-guidance-composition
status: completed
type: infra
complexity: high
---

# Task 07: Synchronize composed Baseline assets

## Overview

Carry the hierarchy, document contracts, semantic destinations, and Profile
adaptation behavior through every maintained embedded and canonical Baseline
asset. Formatter goldens and retention accounting become the distribution
contract for the composed result.

## Requirements

1. MUST update every affected Profile, module, decision effect, template,
   source corpus, coverage declaration, and retention transition coherently.
2. MUST regenerate canonical setup snapshots only from approved immutable
   sources through the Go-owned synchronization command.
3. MUST preserve all upstream-managed skill content and provenance.
4. MUST update formatter fixtures for every affected maintained Profile.
5. MUST make catalog, source accounting, formatter, audit, and empty reapply
   validation fail on any missing asset or clause.
6. MUST keep portable assets free of Fluxus and Oraculum names or policy.

## Subtasks

- [x] Align embedded modules, templates, profiles, and coverage.
- [x] Update source corpora and retention transitions.
- [x] Refresh formatter-compatible golden fixtures.
- [x] Synchronize canonical setup snapshots.
- [x] Add asset completeness and branding guards.

## Acceptance Criteria

- [x] The embedded catalog loads with one deterministic digest.
- [x] Every required clause and semantic destination has complete source
  accounting.
- [x] All maintained Profiles produce formatter-stable generated output.
- [x] Canonical and distributed setup snapshots agree byte-for-byte.
- [x] Upstream skill digests and project-agnostic branding guards pass.

## Context

- instruction: `docs/adr/0059-generated-output-is-formatter-stable-in-the-target-repository.md`
- instruction: `docs/adr/0060-source-baselines-are-exhaustive-and-project-agnostic.md`
- interface: `internal/baseline/assets`
- interface: `internal/baseline/assets_sync.go`
- interface: `internal/baseline/testdata/parity-corpus/v1`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestGuidanceCompositionAssets|TestFormatterComposition|TestSourceBaseline|TestBaselineAssetsSync'` — expected: catalog, accounting, formatter, provenance, sync, and branding cases pass.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline assets sync --help` — expected: the Go-owned synchronization contract remains available.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 1–5; User Stories 1–6; Core Features 1–17; Success Metrics.
- `_techspec.md` → Integration Points; Testing Approach; Build Order 6.
- ADR-0059 → Formatter-Stable Output.
- ADR-0060 → exhaustive project-agnostic assets.

## Result

- The embedded catalog loads deterministically at
  `sha256:267f27fa447014625c70a2efbc53f4cb81f80af0c71eac722cbc1042f782cc6a`;
  missing formatter files, golden drift, and stale accounting targets now
  produce catalog diagnostics.
- The maintained Source Baseline contains 95 independently verified entries.
  Required atomic clauses, direct rules, retained recommendations and
  Operational Contracts are represented, and external-triage, monorepo, and
  skill-dispatch destinations now have source evidence.
- The TypeScript Profile's 14 Markdown goldens were generated through the Go
  planner. The formatter test proves byte agreement, no empty
  repository-specific carrier, successful apply, and verified empty reapply;
  Go and Rust retain their explicit `formatter.kind: none` contracts.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline assets sync
  --source-dir /Users/marcio/dev/skills/setups` refreshed all three setup
  snapshots from clean immutable commit `236847f6956134bf468abb641bac0493a899bca5`.
  The subsequent `--check` returned `setup-context-driven audit: ok`.
- `rtk git diff --name-only -- .agents/skills skills` returned no paths. The
  source corpus and portable generated assets pass the Fluxus/Oraculum,
  generated-marker, and machine-path guards.
- Focused verification passed: `8 passed in 2 packages`. The synchronization
  help command returned exit 0 with the Go-owned source and transaction
  contract. The local pre-feedback `rtk make verify` completed all 2,186
  repository tests, 4 skill runtime contract tests, the Roundfix skill check,
  and the CLI build.

### Verification repair — attempt 1

- The Daemon's first full-gate attempt exposed a deterministic host-global Git
  configuration leak in the asset-sync provenance test: `core.fsmonitor=true`
  could leave a daemon touching the temporary checkout during cleanup.
- Asset-sync source inspection and its temporary Git fixture now pass
  `-c core.fsmonitor=false`, matching the repository's bounded Git inspection
  contract without retries, delays, or suppressed cleanup errors.
- The exact formerly failing subtest passed 20 consecutive runs (`40 passed`).
  The focused asset-sync and Task 07 regression suite then passed
  (`17 passed in 2 packages`). The Daemon owns the authoritative full
  Verification rerun.
