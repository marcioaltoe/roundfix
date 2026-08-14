---
spec: 0030-context-driven-agent-instructions
date: 2026-07-16
run: post-spec0031-rerun-06
build: 9645f29
status: closed
verdict: fail
surfaces: [cli, docs, filesystem]
---

# QA report — Context-driven agent instructions after Spec 0031

## Scope and environment

This full rerun executes the saved Spec 0030 public CLI harness against build `9645f29` after
Spec `0031-decision-driven-setup-generation` completed. It uses disposable repository fixtures,
an isolated `HOME`, the bundled setup assets, and the canonical setup snapshot at
`/Users/marcio/dev/skills/setups`. Network access, live Secondbrain reads, skill installation,
and deletion of installed skills remain outside the harness boundary.

## Static gate

Passed through the saved harness: 73 Python tests, 1,272 Go tests, bundled asset validation,
`roundfix skills check`, sync verification, and build.

## Results

| # | Story / criterion / sweep | Actor and surface | Status | Evidence |
| --- | --- | --- | --- | --- |
| QA-01 | Full repository verification gate | Maintainer · CLI | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-01-static-repository-verification-gate.json` |
| QA-02 | Default audit is read-only | Agent · CLI/filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-02-default-audit-is-read-only.json` |
| QA-03 | First apply returns durable questions and complete preview without writes | Agent · CLI/filesystem | **fail** | `evidence/2026-07-16-post-spec0031-rerun-06/qa-03-first-apply-questions-and-preview.json` |
| QA-04 | TypeScript/Bun setup enables complete English Secondbrain guidance | Agent · CLI/docs | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-04-typescript-bun-secondbrain-opt-in.json` |
| QA-05 | Stored decisions support answer-free, byte-idempotent reapply | Agent · CLI/filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-05-stored-decisions-and-idempotency.json` |
| QA-06 | Go CLI/TUI setup omits Secondbrain when disabled | Agent · CLI/docs | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-06-go-cli-tui-secondbrain-opt-out.json` |
| QA-07 | Rust alternate decisions drive generated workflow and runtime guidance | Agent · CLI/docs | **fail** | `evidence/2026-07-16-post-spec0031-rerun-06/qa-07-rust-alternate-decision-behavior.json` |
| QA-08 | Unmarked guide adoption requires confirmation and preserves owner content | Maintainer · CLI/filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-08-unmarked-guide-adoption.json` |
| QA-09 | Secondbrain opt-out removes only managed content | Maintainer · CLI/filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-09-secondbrain-opt-out-cleanup.json` |
| QA-10 | Missing required skill blocks with a stable finding | Agent · CLI | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-10-missing-required-skill.json` |
| QA-11 | Extra skills stay informational and are never removed | Maintainer · CLI/filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-11-optional-extra-skill-report.json` |
| QA-12 | Invalid input is atomic and uses stable codes | Agent · CLI/filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-12-invalid-input-atomicity.json` |
| QA-13 | Bundled setup snapshots match the canonical setup checkout | Maintainer · CLI | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-13-canonical-setup-snapshot-comparison.json` |
| QA-14 | Normal audit remains portable without the canonical checkout | Agent · CLI | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-14-portable-audit-without-canonical-checkout.json` |
| QA-15 | Canonical and embedded skill copies remain synchronized | Maintainer · CLI | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-15-embedded-skill-synchronization.json` |
| QA-16 | CLI help, JSON, streams, and exits remain deterministic | Agent · CLI | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-16-cli-public-contract.json` |
| QA-17 | Nested instructions, installed skills, and Secondbrain data stay untouched | Maintainer · filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-17-safety-non-goals.json` |
| QA-18 | Source skill trees match their pre-run hashes | Maintainer · filesystem | pass | `evidence/2026-07-16-post-spec0031-rerun-06/qa-18-source-repository-restoration.json` |

## Findings

QA-03 and QA-07 failed because the saved Spec 0030 harness still encoded the pre-Spec 0031
question model. QA-03 expected autonomous runtime and verification decisions before autonomous
work was enabled. QA-07 disabled autonomous work while expecting its dependent guidance to be
rendered. The production behavior matches the approved Spec 0031 dependency contract, so this is
QA harness drift rather than a product defect. The harness is corrected before rerun 07.

## Blocked and skipped

None. All 18 rows executed and source restoration passed.

## Coverage

Executed 18 of 18 rows covering the full CLI, docs, filesystem, portability, safety, and source
restoration matrix: 16 passed and 2 failed.

## Final verdict

**Fail.** The implementation passed all product and repository checks, but this historical run
retains two failures caused by the stale saved harness. Rerun 07 validates the corrected harness
against the Spec 0031 decision dependency contract.
