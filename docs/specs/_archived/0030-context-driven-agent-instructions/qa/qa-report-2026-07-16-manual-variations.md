---
spec: 0030-context-driven-agent-instructions
date: 2026-07-16
run: manual-variations-03-with-fix-retest-05
build: 32728c6
status: closed
verdict: fail
surfaces: [cli, docs, filesystem]
---

# QA Report — Context-driven Agent Instructions

## Scope

This run exercises the shipped `setup-context-driven` Python CLI through its
real process boundary against disposable repository fixtures. It answers the
CLI's durable setup decisions with multiple combinations, captures evidence
after every scenario, and verifies that the source repository is restored to
its pre-run production state.

## Summary

The complete manual sweep closed 18 matrix rows: 16 passed and 2 failed. The
static gate passed with 45 Python tests and 1,272 Go tests. The failures expose
one product gap: durable decisions are collected but the current generator
cannot preview their prospective plan or render most selected values into the
managed instructions. This is larger than a safe same-session patch and is
recorded in `docs/findings/2026-07-16-setup-context-driven-manual-qa.md`.

One small independent defect was fixed during QA: implemented audit options
were described as `Reserved` in help. Commit `d19ca85` corrected the canonical
and embedded skill copies and added a regression test. Focused run 05 passed
the full static gate, QA-02 canary, QA-16, and QA-18 restoration.

## QA Matrix

| ID | Scenario | Expected result | Status | Evidence |
| --- | --- | --- | --- | --- |
| QA-01 | Static repository verification gate | `make verify` passes without skipped checks | pass | [run 05](evidence/2026-07-16-help-fix-retest-05/qa-01-static-repository-verification-gate.json) |
| QA-02 | Default audit against an empty repository | Read-only exit `1`; stable `manifest.missing`; no target writes | pass | [run 05](evidence/2026-07-16-help-fix-retest-05/qa-02-default-audit-is-read-only.json) |
| QA-03 | First TypeScript/Bun apply without answers | Exit `3`; only durable decisions requested; complete managed-change preview; no target writes | **fail** | [run 03](evidence/2026-07-16-manual-variations-03/qa-03-first-apply-questions-and-preview.json) |
| QA-04 | TypeScript/Bun monorepo with Secondbrain enabled | Apply and audit succeed; English managed docs, concise pointer, complete read-only guide | pass | [run 03](evidence/2026-07-16-manual-variations-03/qa-04-typescript-bun-secondbrain-opt-in.json) |
| QA-05 | Stored answers and repeat apply | No repeated decisions; second apply is byte-for-byte idempotent | pass | [run 03](evidence/2026-07-16-manual-variations-03/qa-05-stored-decisions-and-idempotency.json) |
| QA-06 | Go CLI/TUI with Secondbrain disabled | Apply and audit succeed; no Secondbrain pointer or guide | pass | [run 03](evidence/2026-07-16-manual-variations-03/qa-06-go-cli-tui-secondbrain-opt-out.json) |
| QA-07 | Rust CLI with alternate valid decisions | Apply and audit succeed; selected workflow, runtime, and verification answers are reflected in generated guidance | **fail** | [run 03](evidence/2026-07-16-manual-variations-03/qa-07-rust-alternate-decision-behavior.json) |
| QA-08 | Existing unmarked guide adoption | Exit `3` until explicit adoption; accepted adoption preserves unrelated root content | pass | [run 03](evidence/2026-07-16-manual-variations-03/qa-08-unmarked-guide-adoption.json) |
| QA-09 | Opt out after prior Secondbrain opt in | Removes only managed Secondbrain content and preserves user-authored bytes | pass | [run 03](evidence/2026-07-16-manual-variations-03/qa-09-secondbrain-opt-out-cleanup.json) |
| QA-10 | Missing required skill | Blocking `skills.required.missing` finding and exit `1` | pass | [run 03](evidence/2026-07-16-manual-variations-03/qa-10-missing-required-skill.json) |
| QA-11 | Extra locked and untracked skills | Hidden by default; opt-in informational findings; never removed | pass | [run 03](evidence/2026-07-16-manual-variations-03/qa-11-optional-extra-skill-report.json) |
| QA-12 | Invalid decisions, profile, manifest, and markers | Stable exit codes/findings; apply performs no partial writes | pass | [run 03](evidence/2026-07-16-manual-variations-03/qa-12-invalid-input-atomicity.json) |
| QA-13 | Canonical setup snapshot comparison | Current `/Users/marcio/dev/skills/setups` matches bundled snapshots | pass | [run 03](evidence/2026-07-16-manual-variations-03/qa-13-canonical-setup-snapshot-comparison.json) |
| QA-14 | Portable execution without canonical checkout | Normal apply/audit succeeds without developer-specific setup source | pass | [run 03](evidence/2026-07-16-manual-variations-03/qa-14-portable-audit-without-canonical-checkout.json) |
| QA-15 | Canonical and embedded skill copies | Repository sync check proves both shipped copies are identical | pass | [run 03](evidence/2026-07-16-manual-variations-03/qa-15-embedded-skill-synchronization.json) |
| QA-16 | CLI help, JSON, stdout/stderr, and unknown invocation | Public help/output/exit contract is deterministic and machine-readable | pass | [run 05](evidence/2026-07-16-help-fix-retest-05/qa-16-cli-public-contract.json) |
| QA-17 | Non-goals and safety boundaries | No nested `AGENTS.md` edits, skill removals, or Secondbrain writes | pass | [run 03](evidence/2026-07-16-manual-variations-03/qa-17-safety-non-goals.json) |
| QA-18 | Source repository restoration | Production files match the post-fix baseline; only declared QA artifacts remain | pass | [run 05](evidence/2026-07-16-help-fix-retest-05/qa-18-source-repository-restoration.json) |

## Environment

- macOS, Python 3.14.6, branch `ma/setup-context-driven-validator`.
- Real CLI process: `.agents/skills/setup-context-driven/scripts/context_setup.py`.
- Disposable repository targets with an isolated `HOME`; no network or live
  Secondbrain access.
- Canonical setup drift source: `/Users/marcio/dev/skills/setups`.
- Final static gate at commit `32728c6`: 45 Python tests, 1,272 Go tests,
  bundled asset validation, `roundfix skills check`, sync check, and build.

## Defects and retests

1. Run 01 aborted after QA-02 because the new evidence writer could not
   serialize set-valued assertion details. Fixed in `4925d69`/`0a422af`; the
   aborted evidence is retained with `abort.json`.
2. Run 02 reached QA-16, then aborted on byte-valued preservation evidence.
   Fixed in `0a422af`; the aborted evidence is retained with `abort.json`.
3. Run 03 completed the entire sweep: 16 pass, 2 fail. This is the canonical
   product-behavior run.
4. Misleading `Reserved` help text was reproduced by a failing regression
   test, fixed in `d19ca85`, and verified by the full gate.
5. Run 04 proved the help fix but exposed harness-created `__pycache__` through
   QA-18. The harness now sets `PYTHONDONTWRITEBYTECODE=1` in `32728c6`.
6. Run 05 passed the affected QA-16 row, QA-02 canary, full gate, and source
   restoration with zero failures.

## Verdict

**Fail.** The implementation is strong on safety, portability, skills, and the
Secondbrain lifecycle, but it does not yet satisfy the promised customizable
instruction workflow. QA-03 and QA-07 require a spec/task implementation that
adds a real prospective preview and makes durable decisions drive generated
rules and text. The corrective work is specified in
`docs/specs/0031-decision-driven-setup-generation/`.
