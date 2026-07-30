---
status: open
created_at: 2026-07-30
updated_at: 2026-07-30
---

# Run termination does not reach the acpx child (2026-07-30)

Checking whether any Run was still active turned up four `acpx` processes that
had been spinning for **three days and six hours**, since 2026-07-27 11:26.
They belonged to Spec 0037's live QA fixture. Their parent was gone, the
worktrees they pointed at were gone, and nothing had ever told them to stop.

## 1. A terminated Run leaves its acpx child running

- **Symptom / evidence**: four processes matching
  `/private/tmp/roundfix-qa-0037-live.qhCESz/bin/acpx … codex prompt -s roundfix-run_<id>-task_01`,
  each with `PPID 1` — reparented to init because whatever started them exited
  without terminating them.

  ```text
  PID    STARTED                       ELAPSED
  32336  Mon Jul 27 11:26:32 2026      03-06:01:15
  35354  Mon Jul 27 11:27:27 2026      03-06:00:20
  40500  Mon Jul 27 11:28:20 2026      03-05:59:27
  45096  Mon Jul 27 11:29:06 2026      03-05:58:41
  ```

  They were not hung. Each was in the fixture's prompt wait loop:

  ```sh
  *" prompt "*)
    cat >/dev/null
    : > "${ROUNDFIX_QA_PROMPT_STARTED:?}"
    while [ ! -f "${ROUNDFIX_QA_PROMPT_RELEASE:?}" ]; do
      sleep 0.05
    done
  ```

  The release file was never created — the fixture directory held
  `prompt-started`, `prompt-started-2`, and `prompt-started-3` and no release
  marker. So each process forked a `sleep` twenty times a second for
  280,000 seconds: roughly 5.6 million process spawns per orphan.

- **Root cause**: Run teardown does not propagate to the ACP Runtime child. The
  `--cwd` paths those processes carried,
  `/Users/marcio/.roundfix/worktrees/repo-951540fa/run_<id>`, no longer existed
  — so worktree cleanup had run and completed while the processes it was
  cleaning up around kept running. Cleanup removes the filesystem and forgets
  the process tree.

  Spec 0037 gave Force Stop a real ownership proof for the *owner* process.
  Nothing plays the same role for the grandchild: `acpx`, and in production the
  `codex` it spawns.

- **Action / suggestion**: terminate the ACP Runtime child as part of Run
  teardown, on every terminal path — success, failure, Force Stop, and harness
  exit — and prove it with an integration test that asserts no descendant
  survives the Run. A process group or a recorded child PID plus the same
  liveness proof the owner already uses would both work. Until then a
  crash-terminated Run leaves an agent session alive with credentials and a
  writable checkout, which is a larger problem than the wasted cycles.

  Whether production Runs leak the same way is untested. The direct evidence
  here is a QA fixture, but the fixture is a stand-in for `acpx` at exactly the
  seam where teardown should have reached it, and the missing teardown is on
  Roundfix's side of that seam.

## 2. QA fixture directories are never reclaimed

- **Symptom / evidence**: `/private/tmp/` held 54 `roundfix-qa-*` directories
  totalling **7.1 GB**, dated 2026-07-24 through 2026-07-30, from Specs 0036,
  0037, 0038, 0042, 0047, 0050, 0052, 0053, 0054, and 0061 — every one of them
  archived. Most are Go build caches created by QA Agents running tests in
  sandboxed environments (`roundfix-qa-0042-rerun03-gocache` alone was 582 MB).
- **Root cause**: QA Agents create scratch directories under `TMPDIR` and
  nothing removes them when the Spec closes. `roundfix reconcile` reclaims Run
  Worktrees and the Run Database; it has no view of agent-created scratch.
- **Action / suggestion**: either place agent scratch under a Roundfix-owned
  path the GC Command already knows about, or teach the GC Command the
  `roundfix-qa-*` convention so archived Specs release their fixtures. 7 GB
  accumulated in six days of ordinary use.

## What worked — keep

- `roundfix reconcile` was accurate throughout: 136 Runs, all `released`,
  nothing preserved. The leak is entirely outside what it tracks, which is why
  a clean reconcile report did not contradict the evidence.
- The fixture's wait loop is the right shape for a test — it blocks until
  released rather than sleeping a fixed interval. The defect is that nothing
  released it and nothing killed it, not the loop itself.

## Resolution taken

The four processes were terminated with `SIGTERM` — all four exited, so no
`SIGKILL` was needed — and the 54 fixture directories were removed, reclaiming
7.1 GB.

Two Run Worktrees from other repositories were also removed, reclaiming a
further 836 MB. Both turned out to be **orphaned**: their parent repositories
no longer exist, so neither worktree was usable.

| worktree | parent gitdir | state |
| --- | --- | --- |
| `gss-ffb67011` (836 MB) | `~/dev/archive/gss/.git/worktrees/…` | parent gone |
| `rf-red-gate-72ea9b8f` (20 KB) | `/private/tmp/rf-red-gate/.git/worktrees/…` | parent gone |

Neither carried uncommitted work. This is a third instance of the same shape as
finding 1: a Run Worktree outlives the thing that owned it. Removing a
repository, or letting a `/tmp` test repository be reclaimed, leaves its Run
Worktrees behind under `~/.roundfix/worktrees/` with nothing able to classify
them — `roundfix reconcile` reads the Run Database, and a Run whose repository
is gone has no repository to reconcile against. The GC Command could detect an
unresolvable parent gitdir and offer the removal.
