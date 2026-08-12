---
spec: 0030-context-driven-agent-instructions
date: 2026-07-16
run: post-spec0031-rerun-07
build: 9645f29
status: closed
verdict: pass
surfaces: [cli, docs, filesystem]
---

# QA report — Context-driven agent instructions after Spec 0031

## Scope and environment

This full rerun executes the saved Spec 0030 public CLI harness against the worktree based on
build `9645f29`, after aligning the harness with Spec 0031's approved dependent-decision contract.
It uses disposable repository fixtures, an isolated `HOME`, the bundled setup assets, and the
canonical setup snapshot at `/Users/marcio/dev/skills/setups`. Network access, live Secondbrain
reads, skill installation, and deletion of installed skills remain outside the harness boundary.

## Static gate

Passed through the saved harness: 73 Python tests, 1,272 Go tests, bundled asset validation,
`roundfix skills check`, sync verification, and build.

## Results

| # | Story / criterion / sweep | Actor and surface | Status | Evidence |
| --- | --- | --- | --- | --- |
| QA-01 | Full repository verification gate | Maintainer · CLI | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-01-static-repository-verification-gate.json` |
| QA-02 | Default audit is read-only | Agent · CLI/filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-02-default-audit-is-read-only.json` |
| QA-03 | First apply asks entry decisions, defers dependent decisions, and previews without writes | Agent · CLI/filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-03-first-apply-questions-and-preview.json` |
| QA-04 | TypeScript/Bun setup enables complete English Secondbrain guidance | Agent · CLI/docs | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-04-typescript-bun-secondbrain-opt-in.json` |
| QA-05 | Stored decisions support answer-free, byte-idempotent reapply | Agent · CLI/filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-05-stored-decisions-and-idempotency.json` |
| QA-06 | Go CLI/TUI setup omits Secondbrain when disabled | Agent · CLI/docs | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-06-go-cli-tui-secondbrain-opt-out.json` |
| QA-07 | Rust alternate autonomous decisions drive generated runtime and verification guidance | Agent · CLI/docs | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-07-rust-alternate-decision-behavior.json` |
| QA-08 | Unmarked guide adoption requires confirmation and preserves owner content | Maintainer · CLI/filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-08-unmarked-guide-adoption.json` |
| QA-09 | Secondbrain opt-out removes only managed content | Maintainer · CLI/filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-09-secondbrain-opt-out-cleanup.json` |
| QA-10 | Missing required skill blocks with a stable finding | Agent · CLI | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-10-missing-required-skill.json` |
| QA-11 | Extra skills stay informational and are never removed | Maintainer · CLI/filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-11-optional-extra-skill-report.json` |
| QA-12 | Invalid input is atomic and uses stable codes | Agent · CLI/filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-12-invalid-input-atomicity.json` |
| QA-13 | Bundled setup snapshots match the canonical setup checkout | Maintainer · CLI | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-13-canonical-setup-snapshot-comparison.json` |
| QA-14 | Normal audit remains portable without the canonical checkout | Agent · CLI | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-14-portable-audit-without-canonical-checkout.json` |
| QA-15 | Canonical and embedded skill copies remain synchronized | Maintainer · CLI | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-15-embedded-skill-synchronization.json` |
| QA-16 | CLI help, JSON, streams, and exits remain deterministic | Agent · CLI | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-16-cli-public-contract.json` |
| QA-17 | Nested instructions, installed skills, and Secondbrain data stay untouched | Maintainer · filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-17-safety-non-goals.json` |
| QA-18 | Source skill trees match their pre-run hashes | Maintainer · filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-07/qa-18-source-repository-restoration.json` |

## Findings

None. The earlier QA-03 and QA-07 failures were harness contract drift, not production defects.
After aligning those assertions with the approved dependent-decision behavior, every row passed.

## Blocked and skipped

None. All 18 rows executed and source restoration passed.

## Coverage

Executed 18 of 18 rows covering the full CLI, docs, filesystem, portability, safety, and source
restoration matrix. All 18 passed.

## Final verdict

**Pass.** The complete saved Spec 0030 matrix passes after Spec 0031. Entry decisions are asked
first, dependent runtime and verification decisions are deferred until autonomous work is enabled,
selected values drive generated guidance, and all safety, portability, idempotency, skills, and
restoration checks remain green.
