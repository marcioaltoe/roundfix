# Outside evidence assessment

Assessment uses the externally authored passages already quoted in the PRD's Acceptance evidence, as Requirement 7 explicitly requires. No network connection, live documentation fetch or model-provider request was made. Publication contents after the authoring capture are unverified; this row assesses whether those captured passages support the design, not whether the remote sites changed.

| Published source and passage | Design supported | Assessment |
| --- | --- | --- |
| Agent Client Protocol, v2 Prompt Lifecycle RFD, “Report failures after insertion”; https://agentclientprotocol.com/rfds/v2/prompt | ACP v1 prompt failures, including provider outages, use an error response. | Supports failure classification beside an adapter's misleading successful stop reason. The specific Codex metadata detector rests on characterization, not on a normative ACP metadata contract. |
| Agent Client Protocol, Prompt Turn, “Stop Reasons”; https://agentclientprotocol.com/protocol/v1/prompt-turn | end_turn describes the model finishing its response. | Supports rejecting a provider-failed turn even when the adapter returns end_turn. |
| Claude Code, Environment variables, CLAUDE_CODE_DISABLE_BACKGROUND_TASKS entry; https://code.claude.com/docs/en/env-vars | The variable disables Bash/subagent background functionality, automatic backgrounding and the keyboard shortcut. | Supports command refusal. The passage does not explicitly cover Monitor, so task metadata detection remains needed. Environment values are omitted from evidence. |
| Claude Code, Tools reference, Bash backgrounding passage; https://code.claude.com/docs/en/tools-reference | With background functionality disabled, a command reaching its timeout stops. | Supports foreground commands and splitting a longer operation into foreground steps. |

The publisher titles, section names and URLs identify outside authorship. The Spec is only the transport for the recorded quotations; its characterization and its newly written tests are not credited as the outside source.
