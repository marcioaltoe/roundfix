# Tooling authorization — portable user questions (2026-09-08)

The maintainer requested a researched extension of the canonical question rule
for Codex and OpenCode, preserving the interaction expected from Claude Code:

> Quero que seja adicionado algo para ser usado pelo codex e opencode com o
> objetivo de ter uma pergunta por vez com opções de respostas como temos no
> AskUserQuestion tool, pesquise usando o exa mcp sobre o funcionamento da tool
> do claude code e como pode ser aplicado pelo codex e opencode.

## Authorized scope

- `internal/baseline/assets/modules/core.json`: extend
  `clause.core.ask-user-answerable-decisions` with the native tool mappings,
  one pending question at a time, 2-3 explained choices with a recommendation,
  custom text input, actual tool and mode availability, asynchronous waiting,
  and a textual fallback. Increment the affected module, guide, and rule
  versions. Preserve the clause identity and its decision-authority trigger.
- `docs/agents/agent-instructions.md` and `docs/agents/setup-context.json`:
  regenerate the managed guide and manifest through the public Baseline
  update command using the edited catalog.
- Every derived artifact rewritten by `make baseline-digests` is sanctioned
  fallout of this canonical module edit under
  `docs/agents/specific-repository.md`.

This is a guidance change. It does not install tools, change runtime settings,
enable experimental features, alter ACP transport, or modify upstream skills.
The existing Secondbrain policy changes in this worktree remain a separate
authorized change.

## Research and decision

The existing canonical clause already required a structured interaction tool
or an equivalent textual question, but did not name the equivalents. Local
search found literal `AskUserQuestion` references in upstream skill guidance.
The extension makes those references portable without editing the skills.

| Runtime | Agent-facing tool | Application of the common contract |
| --- | --- | --- |
| Claude Code | `AskUserQuestion` | Send one element in `questions`, with explained choices and single selection for mutually exclusive alternatives. |
| Codex | `request_user_input` when permitted | Send one element in `questions`; the researched schema uses `id`, `header`, `question`, and 2-3 `label`/`description` options. Honor advertised mode availability. |
| Codex runtime with asynchronous input | `request_user_input_async` when exposed | Use the actual schema. This session exposes a question `title` and string `options`; keep one request pending and continue only independent work while waiting. |
| OpenCode | `question` | Send one element in `questions`, with `header`, `question`, and explained `options`; preserve the native custom-answer path. |

The one-question limit is the maintainer's policy: the tools themselves can
accept batches. Choosing 2-3 options fits the researched synchronous Codex
schema and Claude Code's 2-4 option range. The recommendation remains advice,
not a submitted answer. Tool availability is determined by the active session,
not by a client name or a hard-coded assumption that Codex is Plan-only.

### External sources read

- [Claude Code: Handle approvals and user input](https://code.claude.com/docs/en/agent-sdk/user-input), found and fetched through Exa MCP. Establishes the `AskUserQuestion` schema, 1-4 questions, 2-4 choices, single/multiple selection, custom answers, and the distinction between clarifying questions and permission requests. It also documents that the tool is unavailable to SDK subagents; the rule therefore checks actual availability.
- [OpenAI Codex: request_user_input tool source](https://github.com/openai/codex/blob/35aaa5d9/codex-rs/tools/src/request_user_input_tool.rs), found and fetched through Exa MCP. This pinned source defines the tool name, 2-3 mutually exclusive options, recommendation ordering, automatic `Other`, and availability derived from collaboration modes and the Default-mode feature. It supports capability-based wording rather than assuming one universal mode restriction.
- [Codex App Server](https://developers.openai.com/codex/app-server), fetched through both Exa MCP and OpenAI Docs. Documents the client protocol `tool/requestUserInput` and the pending-request lifecycle. The protocol method is not the agent-facing tool name and is not a command for an agent to call directly.
- [OpenCode: Tools — question](https://opencode.ai/docs/tools/#question), fetched through Exa MCP. Establishes the native `question` tool, headers, options, custom answers, and batch support.
- [OpenCode question tool instructions](https://github.com/anomalyco/opencode/blob/dev/packages/opencode/src/tool/question.txt), retrieved through Context7 using `/anomalyco/opencode`. Specifies automatic custom input, label-array answers, single selection by default, and the recommended option first. This branch reference is mutable; the installed runtime's schema remains authoritative.

The asynchronous Codex mapping is grounded in the callable tool schema
advertised in this session. The public sources above do not establish that
every Codex client exposes it. Local Codex and OpenCode executables were found,
but local source checkouts were absent. No live interactive test was performed
in the Claude Code or OpenCode clients.

### Secondbrain context

The session had read `/Users/marcio/dev/secondbrain/wiki/index.md` before
querying. The advisory query was:

```text
rtk proxy qmd query 'perguntas estruturadas AskUserQuestion Codex OpenCode uma pergunta por vez opções decisão usuário' --all --files --min-score 0.3
```

It returned adjacent material, including the read entry
`/Users/marcio/dev/secondbrain/inbox/secondbrain/_triaged/2026-08-17-empurrar-evento-externo-para-dentro-de-sessao-de-agente.md`.
That entry concerns external event ingress, not question widgets. Its boundary
between data and authority informed the rule that a recommendation, preselection,
or timeout is not the user's choice; its dated runtime availability claims were
not reused. No directly matching local question-tool comparison was found.
The query warned that nine documents needed embeddings, limiting coverage.

The capture workflow was checked against
`/Users/marcio/dev/secondbrain/AGENTS.md`,
`/Users/marcio/dev/secondbrain/inbox/README.md`, and
`/Users/marcio/dev/secondbrain/templates/inbox-entry.md`.

The research digest is captured for later ingestion at
`/Users/marcio/dev/secondbrain/inbox/secondbrain/2026-09-08-perguntas-estruturadas-claude-codex-opencode.md`.

## Verification and delivery

Apply this bounded documentation change directly. Regenerate the catalog,
review and apply the public Baseline update plan, then require a recheck with
no pending changes. Run the repository Verification and documentation
contracts. Inspect the generated guide for the approved policy; these checks
validate distribution of instructions, not interactive UI behavior in all
three clients.

When committed, this authorization record must land in its own commit before
the canonical module and generated-artifact changes. The maintainer subsequently
authorized committing and pushing all current Roundfix changes on
`ma/canonical-research-and-user-questions`. The separate Secondbrain digest
follows its existing commit-at-capture contract.
