---
spec: 0029-launch-and-recovery-fixes
status: archived
created: 2026-07-15
surfaces: [cli, docs]
archived: "2026-07-15"
source_slug: 0029-launch-and-recovery-fixes
---


# Launch and Recovery Fixes

Engineering-framed minimal PRD (bug-fix route): three launch/recovery defects reproduced in the 2026-07-14/15 dogfood sessions, each breaking a contract an earlier spec already defined. Field evidence and diagnosis live in `docs/findings/2026-07-14-implement-0027-launch-findings.md` and `docs/findings/2026-07-15-retrospectiva-workflow-skills-e-versao-dev.md` §6.

## Problem

1. **`--detach` kills its child silently.** The parent's handshake deadline (10s) is shorter than a real Preflight Validation, whose model-availability probe measured 11.4s on the dogfood machine. The parent times out, SIGKILLs the child mid-preflight, relays the empty console temp, and exits 1 with no message — the Detached Run contract promises the child's stderr and a report, and delivers neither.
2. **A model rejection at Batch time is an opaque `agent/protocol error`.** codex-acp assembles its advertised-model list per session (dynamic fetch with a fallback), so the Preflight probe and a later Batch session can see different lists — on 2026-07-14 the probe accepted `gpt-5.6-sol` and every Batch died. Upstream nondeterminism cannot be fully preflighted; the Batch-time failure must become first-class and actionable instead.
3. **The Settle Command resolves a stale kept worktree as its surface.** Surface resolution prefers the latest kept Run's worktree without checking the Task's status there, refusing with "status is pending" while the authoritative checkout has the Task `failed` (reproduced 2026-07-14 during 0027 recovery).

## Goals

- A Detached Run start either produces the four-line report or a diagnostic that names what happened (phase, child exit code or signal, and the child's output or its absence) — never a silent exit 1.
- A slow but healthy Preflight Validation never causes a Detached Run start to fail.
- A Batch that dies because the Agent Session rejects the Agent Model reports the model, the runtime's advertised models, and the recovery paths — same quality as the Preflight rejection — and the Doctor Command shows the effective model probe result.
- Settle selects a surface where the target Task is actually `failed`, and always names the surface it chose.

## Core Features

1. Two-phase Detached Run handshake: the child signals liveness immediately, then reports the run id when the Run exists; the parent bounds each phase separately, and every failure branch prints an explicit diagnostic.
2. Batch-time Agent Model rejections are classified as selection failures with the advertised-model list and recovery actions, not generic protocol errors.
3. The Doctor Command reports the configured Agent Model's probe outcome, including the advertised list on failure.
4. Settle surface resolution prefers surfaces where the Task is `failed`, reports the chosen surface, and names each candidate's status when none qualifies.
5. The Roundfix Skill's Detached Run, Doctor, and Settle sections ship with the behavior (skill-sync rule).

## Non-Goals / Out of Scope

- Fixing acpx/codex-acp's advertised-list nondeterminism upstream.
- Settle pathspec scoping (deliberately deferred by 0028 pending field evidence).
- Any change to the four-line detach report, exit-code contract, or other public output on the success paths.
