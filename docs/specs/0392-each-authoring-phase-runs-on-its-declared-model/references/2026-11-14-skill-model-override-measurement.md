# Measurement: does Claude Code apply a skill's `model` frontmatter? (recorded by the operator, 2026-11-14 ~07:35 BRT)

Machine: Mac mini, Claude Code 2.1.325, headless (`claude -p … --model sonnet --output-format stream-json --verbose --setting-sources project --allowedTools Skill --no-session-persistence`), in a throwaway Git repository holding only project skills under `.claude/skills/`. The model of each assistant record was read from `message.model` in the stream.

| Case | Skill frontmatter | How the skill ran | Model of the reply |
| --- | --- | --- | --- |
| A | `model: haiku` | the agent called the Skill tool | `claude-sonnet-5-5` (override not applied) |
| B | `model: haiku` | the user typed `/tier-probe` | `claude-haiku-5-5` (override applied) |
| C | `model: haiku`, `effort: low` | the agent called the Skill tool | `claude-sonnet-5-5` (override not applied) |
| D | `model: haiku`, `context: fork` | the agent called the Skill tool | main-thread replies on `claude-sonnet-5-5`; the run's usage by model lists `claude-haiku-5-5` (US$ 0.0019), which is consistent with the forked skill running on Haiku but was not isolated from the harness's own small-model calls, because case A's usage by model was not captured |

This matches the open Claude Code issue anthropics/claude-code#98898 ("Skill frontmatter model:/effort: ignored when Claude invokes the skill via the Skill tool (applied on /skill-name)", opened 2026-10-02, reproduced by others on 2.1.287 to 2.1.311) and contradicts the Indexus message's item 1 ("inclusive quando o próprio agente chama a skill"). Issue #81618 (`effort` beside `model` drops the model override) was closed on 2026-10-23.

Limits: headless mode only; an interactive session was not measured here. One run per case.
