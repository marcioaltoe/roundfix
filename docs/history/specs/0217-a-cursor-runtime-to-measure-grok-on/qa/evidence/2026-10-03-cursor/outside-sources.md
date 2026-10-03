# Outside evidence observed 2026-10-03

Public pages retrieved with the web tool. Shell curl-cffi was refused by the sandbox allowlist; no repeated shell request or escalation was used.

## Cursor ACP documentation

Source: https://cursor.com/docs/cli/acp
Sections: Request flow, Authentication, Cursor extension methods.
The documented ordering is initialize, authenticate (methodId cursor_login), then session/new or session/load. The ask_question and create_plan extensions require a reply before work proceeds. This independently supports the authentication and unattended-session hazards. The supervised record observes a preauthenticated install opening a session without an authenticate message; a typical documentation flow is not a claim that every authenticated install repeats it.

## Cursor model switching thread

Source: https://forum.cursor.com/t/bug-agent-acp-model-switching-updates-session-metadata-but-does-not-change-the-inference-backend/157312
Post 9 publishes configOptions including default[], grok-4-20[thinking=true], and gpt-5.4[reasoning=medium,context=272k,fast=false]. Cursor staff member Colin reported the latest CLI fix on 2026-04-30 (post 22). The published values match the retained fixture semantics; the staff report is evidence of a shipped fix announcement, not proof of current backend identity. The operator measurement independently proved a currently advertised Grok selection.
