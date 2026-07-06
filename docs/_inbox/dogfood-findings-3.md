# Dogfood findings — round 3 (watch on Roundfix's own PR #16)

Third supervised round: PR #16 (the whole 0002–0008 arc, 67 commits) pushed,
reviewed by CodeRabbit, and cleaned by `roundfix watch` running every piece
shipped today — Run Worktree isolation, poll-first, merge-readiness,
`--no-agent-console`, stdout reports — plus one real force-stop recovery.
Run ids: `run_20260705T213559Z_0e380ee839ed44ce` (killed externally,
force-stopped) and `run_20260705T214901Z_770682120bee1cfa` (Clean after 1
Round: 6 resolved, integrated, pushed, threads closed).

Validated in production this round: cold Run Worktree passed `make verify`
with no `worktree.copy` needed here; batch commit integrated into the user
checkout by fast-forward and the Run Worktree was removed after Clean; the
deterministic stdout report matched the documented shape; stderr stayed free
of Agent console noise under `--no-agent-console`.

1. **Force-stop recovered a killed Run in one command — but its debris is
   kept forever.** The externally killed watch left an Active Run, held
   locks, and an orphaned detached agent turn; `roundfix stop --force
   <run-id>` (first production use) cancelled the session, completed the Run
   Stopped, and released everything. Gap: the Stopped Run's kept worktree
   and Run Branch sit at the base commit with zero settled work, yet
   `PruneTerminal` reaps only terminal-Clean Runs — provably valueless
   debris accumulates. Proposal: the preflight sweep (or `stop --force`) may
   also remove kept worktrees whose Run Branch carries no commits beyond its
   base — nothing to lose, cheap ancestry check. Size: small.

2. **CodeRabbit issue titles are summary-table fragments.** The stdout
   report and Work Queue rows rendered titles like `_🩺 Stability &
   Availability_ | _🟡 Minor_ | _⚡ Quick win_` — raw markdown italics,
   emoji, and category/severity fragments instead of a descriptive finding
   title. Fetch-side title derivation should prefer the finding's
   description line and strip markup for the deterministic report. Size:
   small.

3. **Merge-readiness `missing` path observed working as documented.** This
   repo's CodeRabbit posts no commit status, so the confirm phase noted
   `Review Source check missing for the pushed HEAD; treating Run as Clean.`
   and ended — the ADR-0019 degradation exactly as specified. Note for
   docs/expectations: repositories wanting hard merge-readiness need the
   CodeRabbit status check enabled; otherwise Clean keeps the pre-0004
   meaning with the note as the tell. Size: docs note only.

4. **Status-poll stderr repetition.** `Review Source status: reviewing`
   repeats every poll interval on stderr during long review waits; harmless
   for humans, noisy for supervisors even with `--no-agent-console`. A
   changed-only or dedup mode for repeated daemon status lines would finish
   the supervision story. Size: cosmetic.

(The fixed review batch commit template `fix: resolve Roundfix batch 001`
remains tracked by work-plan item 1; not re-recorded here.)

## 0009 cycle (2026-07-05/06)

5. **Removing a config key with a hard-fail broke the user's own config.**
   0009 task_01 (my PRD/techspec decision) made `resolve.concurrent` a
   Preflight hard-fail: `resolve.concurrent has been removed; use
   worktree.concurrency instead`. But the dogfood machine's live
   `~/.roundfix/config.yml` still carried it, so the freshly-built binary
   rejected every Run — implement, resolve, watch, and the 0009 QA — at
   Preflight in ~1s. The QA "ran overnight" only in appearance; it died
   immediately and was never executed. Root cause: hostile migration. A
   removed config key should degrade to a deprecation **warning** (or be
   auto-migrated), never a hard Preflight failure — the very users who set
   the old key are the ones broken. Fix needed in a follow-up: turn the
   `resolve.concurrent` rejection into an ignored-with-warning path.
   Size: small (compat), but a design lesson: never hard-break existing
   config.

6. **External kills orphan acpx adapter processes; they accumulate badly.**
   After a day of externally-killed Runs, `pgrep codex-acp` found **40**
   orphaned adapter processes (PPID=1), each holding memory. acpx spawns
   codex-acp as a detached queue-owner (the resume feature), so killing the
   parent roundfix leaves the adapter running, and `stop --force` cancels
   the acpx *session* cooperatively without reaping the OS process tree.
   Every kill+relaunch leaked more. Candidates: `stop --force` (and the
   preflight sweep) should reap the Run's adapter process tree, and/or a
   `roundfix doctor`/`setup` check that lists and offers to kill orphaned
   adapters. Size: small/medium. Pairs with the acpx detached-owner model
   (round-2 finding 1 territory).

## Supervision lessons (self-inflicted, 2026-07-06)

- A Monitor grep pattern MUST cover failure signatures, not just the happy
  path: the QA monitor watched `^Task|^qa|^Clean|^Unresolved|failed after
  Run start` and missed `Preflight failed`, so a dead process read as
  "still running" for hours. Widen supervision filters to every terminal
  signature.
- `nohup ... & disown` detaches a Run from the harness (surviving external
  task stops — the intended fix) but also removes the harness completion
  notification, so a monitor is the ONLY signal. The monitor must therefore
  be comprehensive.
