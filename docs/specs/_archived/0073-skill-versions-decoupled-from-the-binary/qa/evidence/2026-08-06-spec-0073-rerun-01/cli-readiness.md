# CLI readiness journeys

The current built binary's `skills check` exited 0. Doctor independently
reported `skills: ok (38 required: 14 Roundfix-owned, 24 external)`; unrelated
adapter and profile checks failed because the sandbox cannot initialize Codex
state under `/Users/marcio/.codex`, but the independent Skill check completed.

A scratch clone at build `5553c3aa` supplied rebuilt binaries and matching
installed Skill fixtures:

- Declared `0.1.0` against minimum `0.0.2`: `skills check` exited 0 and Doctor
  reported `skills: ok`.
- Declared `0.0.1`: `skills check` exited 1 and Doctor reported `skills:
  failed`; both named `roundfix`, required `0.0.2`, found `0.0.1`, and the
  project-install upgrade action.
- No top-level declaration: `skills check` exited 0 with `Roundfix skill check
  unversioned: roundfix`; Doctor reported `skills: unversioned`.
- An unreadable installed `roundfix/SKILL.md` produced `skills: unversioned`.
  Moving the complete installed Skill directory away produced `skills: failed
  (missing: roundfix; next: roundfix skills install --target project)`.

The focused call-ledger and shared-surface tests also passed for satisfies,
below, unversioned, missing, and third-party fixtures.

