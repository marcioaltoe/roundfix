# Measurement: what Claude Code applies from a subagent file, and what the current Roundfix does with the declaration (recorded by the author, 2026-11-14, 07:33–07:55 BRT)

Machine: Mac mini, Claude Code 2.1.325, headless. Each case ran in a throwaway Git repository created with `mktemp -d` under the authoring clone and deleted afterwards, holding only project subagents under `.claude/agents/` and, for cases 9 to 11, project skills. The session model was `sonnet` (`claude-sonnet-5-5`). Twenty-two runs, two of them repeats of cases 1 and 3, US$ 1.44 in total.

The command, with `<prompt>` and `<tools>` per case:

```bash
claude -p '<prompt>' --model sonnet --output-format stream-json --verbose --setting-sources project --strict-mcp-config --allowedTools <tools> --no-session-persistence
```

`--setting-sources project` keeps the user's hooks out and `--strict-mcp-config` keeps the user's MCP servers out. Every prompt named the subagent and said "passing no model and no effort", except case 8. The model of a subagent was read from `tool_use_result.resolvedModel` of the Agent call's result record, from `message.model` of the subagent's own records when it used a tool, and from `modelUsage` of the final `result` record. Case 2 ran without `--no-session-persistence`, and the effort was read from the `effort` field of the assistant records in the subagent's transcript under `~/.claude/projects/<project>/<session>/subagents/`.

## Results

| Case | Subagent file (or skill) | What ran | Result |
| --- | --- | --- | --- |
| 1 | `model: haiku` | the agent started it with the Agent tool | `resolvedModel` `claude-haiku-5-5`; the session stayed on `claude-sonnet-5-5` |
| 2 | `model: sonnet` with `effort: low`, then with `effort: high`; `model: haiku` without `effort` | the agent started each one | transcript records: `effort: low`, `effort: high`, and `effort: medium` (the session's level) for the one without the field; the session's own records stayed `effort: medium` |
| 3 | `model: opus`, `effort: high` | the agent started it | `resolvedModel` `claude-opus-5-5`: a subagent file raises the model above the session's as well as lowering it |
| 4 | the same `model: sonnet` subagents with `effort: low` and `effort: high`, a counting question, three runs each | the agent started it | `thinking_tokens` 0, 0, 0 at `low` and 106, 95, 96 at `high`, on the same model and prompt |
| 5 | `tools: Read, Grep, Glob`, asked to create a file with a shell or a file-writing tool | the agent started it; the session itself was allowed `Bash`, `Write` and `Edit` | no tool call, no file; it named `Read`, `Grep` and `Glob` as its only tools |
| 6 | no `tools` field, the same request | the agent started it | it called `Bash` and wrote the file: the control of case 5 |
| 7 | `author-probe` (`model: haiku`, `tools: Agent, Read`) told to start `nested-probe` (`model: opus`) | the agent started `author-probe` | `task_started` records with `spawn_depth` 1 and 2; `modelUsage` lists `claude-opus-5-5` for the nested one: a subagent starts a subagent, which runs on its own `model` |
| 8 | `model: haiku` | the agent started it passing `model: "opus"` in the Agent call | `resolvedModel` `claude-opus-5-5`: a model passed in the call overrides the file |
| 9 | skill with `context: fork` and `model: haiku` | the agent called the Skill tool | a forked subagent ran it; its record names `claude-haiku-5-5`, and `modelUsage` lists that model beside the session's |
| 10 | skill with `model: haiku` and `effort: low` in `.agents/skills/<name>/SKILL.md`, reached through the link `.claude/skills/<name>` → `../../.agents/skills/<name>` | the prompt was `/<name>` (typed) | the reply ran on `claude-haiku-5-5`: both keys together apply, through the link |
| 11 | the same skill | the agent called the Skill tool | the reply ran on `claude-sonnet-5-5`: not applied, as in the operator's cases A and C |
| 12 | `model: haiku`, `effort: hihg` | the agent started it | it ran; only the debug log said `has invalid effort 'hihg'. Valid options: low, medium, high, xhigh, max or an integer` |
| 13 | `model: not-a-model-name` | the agent started it | the Agent call returned `Agent terminated early due to an API error: There's an issue with the selected model (not-a-model-name)` (`model_not_found`, HTTP 404) |

In a control without any Haiku subagent or skill (case 3), `modelUsage` listed no Haiku model, so the Haiku usage of cases 1, 2 and 9 belongs to the subagent or the forked skill.

## What it settles

- Claude Code 2.1.325 applies both `model` and `effort` of a subagent file under `.claude/agents/` when the agent starts that subagent (cases 1 to 4). This answers the Idea's first Open Question: yes.
- A `tools` list is a boundary the subagent cannot cross (cases 5 and 6).
- A subagent can start another subagent, which runs on its own model (case 7), so an author subagent can send sweeps to a research subagent.
- The agent can override the declared model by passing one in the Agent call (case 8), so the skill text must forbid it.
- A header with `context: fork` would run an agent-called skill on the header's model (case 9). It is not used: a forked skill starts without the conversation, and Claude Code withholds the question tool from every subagent, so a phase the user starts by typing the skill would lose its interview.
- A mistyped effort is ignored in silence and a mistyped model fails aloud (cases 12 and 13).

## The current Roundfix and a Project Config that holds the declaration

Built from the code of `251ed81e` (v0.143.3) and run with a temporary home in a temporary repository whose `.roundfixrc.yml` holds a `skills:` section:

| Command | Result |
| --- | --- |
| `roundfix skills check` | passes; it reads no Project Config |
| `roundfix skills install --target project` | installs the bundle without a header; it reads no Project Config |
| `roundfix baseline update --repo . --no-skills --no-history --format text` | runs; it reads no Project Config before its history stage |
| `roundfix profiles show`, `doctor`, `update check`, `spec check <slug>`, `runs list` | exit 2 with `parse config "<path>/.roundfixrc.yml": yaml: unmarshal errors:` and `  line 1: field skills not found in type config.configOverlay` |

So a Roundfix older than the release that introduces the section refuses every command that loads configuration, naming the unknown field and no fix, and its installation erases a declared header without a word.

## The other harnesses and the header

The header is written into `.agents/skills/<name>/SKILL.md`, the one copy Claude Code, Codex and OpenCode read. Read on 2026-11-14 through Exa:

- Codex, `codex-rs/skills/src/parser.rs` on the `main` branch of `openai/codex`: `SkillFrontmatter` deserializes `name`, `description` and `metadata` and does not deny unknown fields, so `model` and `effort` are dropped. Every owned skill already carries `argument-hint` and `version`, which that parser does not know either.
- OpenCode, https://opencode.ai/docs/skills/: "All frontmatter is optional at runtime." and, on loading a skill, it "Adds the Markdown body, without frontmatter, to the conversation."

## Limits

Headless mode only; an interactive session was not measured. One run per case, except case 4. Case 4 is indirect evidence of the effort level; case 2 reads it from the transcript. Codex and OpenCode were read, not run.
