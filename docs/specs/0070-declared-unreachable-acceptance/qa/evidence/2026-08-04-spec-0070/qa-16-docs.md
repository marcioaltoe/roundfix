# QA-16 — Operator documentation drift

Status: fail

Finding: F-001

The built `roundfix archive --help` accurately states that the newest QA
Report may have `verdict: pass` or a partial verdict whose blocked rows are
covered only by declared Unreachable Acceptance. The real declared-only
journey exited 0 and stamped `unproven`.

Three operator-facing sources still state the old pass-only rule:

- `docs/user-guide/commands.md:600` says the newest report must have
  `verdict: pass`;
- `.agents/skills/roundfix/SKILL.md:1731` says the same;
- the Roundfix Skill's loop guidance at lines 1519 and 1578 routes archive
  only on `verdict: pass`.

Following those sources hides the feature this Spec ships and can direct a
maintainer toward the explicit failed-evidence override even when every unmet
row is validly declared unreachable. The docs surface therefore contradicts
the public CLI and the PRD's primary operator outcome.
