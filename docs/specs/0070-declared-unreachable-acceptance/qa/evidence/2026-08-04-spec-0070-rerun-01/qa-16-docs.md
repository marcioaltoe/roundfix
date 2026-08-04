# QA-16 — Operator documentation correction

Status: pass

The previously failing Archive Command documentation now matches the current
public CLI and the observed archive boundary.

- `rtk grep -n -C 3 "verdict: pass\\|unproven\\|finding-blocked\\|environment-blocked\\|qa_override" docs/user-guide/commands.md .agents/skills/roundfix/SKILL.md` found the two archive-eligible shapes, full declaration coverage, the `unproven` record, and every unchanged refusal in both operator-facing sources.
- `rtk cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` exited 0.
- The built `/private/tmp/roundfix-qa70-a10638b archive --help` exited 0 and described `verdict: pass` plus declared-only partial acceptance.
- Following the written contract against a fresh declared-only Spec Root exited 0, moved the active Spec, and a fresh read found `unproven: [a maintainer performs the live publication and records it]`. Repeating the command exited 2 because the move persisted.
- Following the documented finding-blocked refusal exited 2 with `rows_blocked_finding is 1; expected 0`; a fresh read retained the active Spec and found no archive destination.

F-001 from `qa-report-2026-08-04.md` is resolved on build
`a10638b9245880c4681fcf622048e13c6e72ce6b`.
