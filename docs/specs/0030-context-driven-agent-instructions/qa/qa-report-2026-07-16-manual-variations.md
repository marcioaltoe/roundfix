---
spec: 0030-context-driven-agent-instructions
date: 2026-07-16
run: manual-variations-01
status: in_progress
verdict: pending
surfaces: [cli, docs, filesystem]
---

# QA Report — Context-driven Agent Instructions

## Scope

This run exercises the shipped `setup-context-driven` Python CLI through its
real process boundary against disposable repository fixtures. It answers the
CLI's durable setup decisions with multiple combinations, captures evidence
after every scenario, and verifies that the source repository is restored to
its pre-run production state.

## QA Matrix

| ID | Scenario | Expected result | Status | Evidence |
| --- | --- | --- | --- | --- |
| QA-01 | Static repository verification gate | `make verify` passes without skipped checks | pending | pending |
| QA-02 | Default audit against an empty repository | Read-only exit `1`; stable `manifest.missing`; no target writes | pending | pending |
| QA-03 | First TypeScript/Bun apply without answers | Exit `3`; only durable decisions requested; complete managed-change preview; no target writes | pending | pending |
| QA-04 | TypeScript/Bun monorepo with Secondbrain enabled | Apply and audit succeed; English managed docs, concise pointer, complete read-only guide | pending | pending |
| QA-05 | Stored answers and repeat apply | No repeated decisions; second apply is byte-for-byte idempotent | pending | pending |
| QA-06 | Go CLI/TUI with Secondbrain disabled | Apply and audit succeed; no Secondbrain pointer or guide | pending | pending |
| QA-07 | Rust CLI with alternate valid decisions | Apply and audit succeed; selected workflow, runtime, and verification answers are reflected in generated guidance | pending | pending |
| QA-08 | Existing unmarked guide adoption | Exit `3` until explicit adoption; accepted adoption preserves unrelated root content | pending | pending |
| QA-09 | Opt out after prior Secondbrain opt in | Removes only managed Secondbrain content and preserves user-authored bytes | pending | pending |
| QA-10 | Missing required skill | Blocking `skills.required.missing` finding and exit `1` | pending | pending |
| QA-11 | Extra locked and untracked skills | Hidden by default; opt-in informational findings; never removed | pending | pending |
| QA-12 | Invalid decisions, profile, manifest, and markers | Stable exit codes/findings; apply performs no partial writes | pending | pending |
| QA-13 | Canonical setup snapshot comparison | Current `/Users/marcio/dev/skills/setups` matches bundled snapshots | pending | pending |
| QA-14 | Portable execution without canonical checkout | Normal apply/audit succeeds without developer-specific setup source | pending | pending |
| QA-15 | Canonical and embedded skill copies | Repository sync check proves both shipped copies are identical | pending | pending |
| QA-16 | CLI help, JSON, stdout/stderr, and unknown invocation | Public help/output/exit contract is deterministic and machine-readable | pending | pending |
| QA-17 | Non-goals and safety boundaries | No nested `AGENTS.md` edits, skill removals, or Secondbrain writes | pending | pending |
| QA-18 | Source repository restoration | Production files match the post-fix baseline; only declared QA artifacts remain | pending | pending |

## Environment

Pending execution evidence.

## Defects and retests

Pending execution evidence.

## Verdict

Pending until every matrix row is closed as `pass`, `fail`, `blocked`, or
`skipped`.
