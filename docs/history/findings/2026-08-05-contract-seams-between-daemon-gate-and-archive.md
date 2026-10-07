---
status: done
created_at: 2026-08-05
updated_at: 2026-09-08
absorbed_by: 0127-durable-unattended-spec-workflow
---

# 2026-08-05 — Contract seams between Daemon, gate, and Archive

Source: three consecutive Specs delivered end to end in the `fiscus` repository (`0001-fundacao-de-contextos`, `0002-auth-staff-e-directory`, `0003-design-system-graphite`) — 35 Tasks, 17 Implement Runs, 14 gate executions, 7 merged pull requests. Companion to `2026-08-05-five-frictions-from-a-full-autonomous-spec-night.md` and `2026-08-05-authored-verification-gates-are-untested-code.md`, which cover the gate-authoring and worktree classes. This one records what the **components disagree about**: the Daemon, the QA gate, and Archive each hold a valid contract, and the seams between them cost this session five extra Runs and four supervisor edits to Daemon-owned artifacts.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-05-contract-seams-between-daemon-gate-and-archive.md`.
