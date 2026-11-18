### Local Spec workflow

- Route and execute local Specs through {{reference.spec-routing}} and {{reference.issue-tracker}}.
- In Claude Code, a skill's `model` and `effort` header applies only when a
  person types the skill (`/skill`); a skill an Agent activates runs on the
  session's model (anthropics/claude-code#98898, measured 2026-11-14). To
  hold an Authoring Phase to its declared model there, hand the phase to its
  Author Subagent, `.claude/agents/roundfix-<phase>-author.md`, which
  `roundfix skills install --target project` generates from
  `skills.authoring.<harness>`; when that model's usage limit ends the
  subagent, run `roundfix phase fallback` and restart it on the model it prints.
  This bullet does not apply to Codex or OpenCode, whose models come from
  the Agent Selection Profiles.
