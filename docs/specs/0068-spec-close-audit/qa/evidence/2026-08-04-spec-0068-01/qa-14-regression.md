# QA-14 — Existing command regressions

Status: pass.

- Built `bin/roundfix spec check 0068-spec-close-audit` exited 0 with no
  finding; its two skips identify optional absent sections.
- Built `bin/roundfix reconcile` and `--format json` exited 0 in dry-run mode
  without applying anything.
- The fresh 15-test CLI selection passed established Spec Check output,
  diagnostics, strict mode, JSON, discovery, and Spec Audit behavior.
- The fresh six-test worktree selection passed deleted-target safety and
  preservation cases.
- Full `rtk make verify` passed on the same build.

Pre/post snapshots confirm the public regression commands did not change refs,
worktrees, status, or the Run Database.
