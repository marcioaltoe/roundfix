# Skill Modes

The Roundfix skill reads an optional Skill Mode from the user's request and follows the selected mode while keeping command-specific details in the relevant reference.

## How a mode is chosen

Claude Code appends the text typed after the skill's name to the skill as an `ARGUMENTS:` line, while Codex and OpenCode keep it in the message, so the skill reads the mode from the request text in all three. The mode is the first word after the skill's name when that word is `follow-up`, `status` or `unattended`, with or without a leading `--`, so `/roundfix follow-up`, `/roundfix --follow-up`, `$roundfix status` and `@roundfix --unattended` each select a mode. This is the Skill Mode, not the CLI command `roundfix follow-up`, which the follow-up mode runs as one of its steps. A request whose first word is anything else, or that has none, uses the skill as before (for example `/roundfix watch 123`), and a mode word later in such a request selects nothing. When the words after the mode word name another mode, with or without `--`, the skill asks one question that names them, recommending `follow-up` when it is named and `status` otherwise. Each mode runs Roundfix's read-only follow-up command; when the installed Roundfix lacks it, the mode says so, names `roundfix upgrade` and stops.

## Optional sources

The Secondbrain is a read-only knowledge source when the repository has `docs/agents/secondbrain.md`; consult it while authoring as that guide says. It holds no inbox, and the follow-up command reads no Secondbrain. ai-memory is a source when the session exposes its MCP tools: list this project's messages without consuming them, write them to a JSON file in a temporary directory using the `roundfix/follow-up-messages/v1` schema documented by the command, and hand the file to `roundfix follow-up --messages <path> --format json`; consume a message only after recording it as a Spec, Backlog Entry, Finding or answer. Jev is used by the follow-up command when a judge key is configured; relay the command's own report when Jev was unavailable. A release-announcement step exists only when the user's own instructions define one; this skill never names one. A missing source never stops a mode; its report names each missing source and why.

## --status

Run `roundfix follow-up --status --format json` and report, in a few lines, the work in progress, the pending Specs in delivery order, the next release and what it brings, the blockers, and what waits on the user. The command's status view reads `runs`, `delivery-queue`, `pull-requests`, `git`, `specs` and the Release Plan; it writes nothing, asks nothing, consumes no ai-memory message and arms nothing. Read the command's `state`, `sources`, `judge` and `items` fields as documented in [Follow-up Command](follow-up.md); the message-file option is `--messages <path>`, its schema is `roundfix/follow-up-messages/v1`, and output formats are `--format <text|json>`.

## --follow-up

1. Inventory: run `roundfix follow-up --format json`, handing it the ai-memory message file when that source exists; when it reports nothing to do and nothing in progress, report what it checked and stop.
2. Plan: show a short plan in stages, grouping sources into Specs as `docs/agents/spec-routing.md` allows and following the command's priority order.
3. Run `--unattended` to gather every decision and authorization the plan needs.
4. Triage, after the answers: only ai-memory messages are triaged, as the cross-project continuity rule of `docs/agents/agent-instructions.md` and ADR-0330 describe, because the Secondbrain holds no inbox to triage. When the ai-memory source exists, consume each message the mode handed in (`memory_message_pop`) only after its content is recorded in a Spec, a Backlog Entry or a Finding, or when it is informational and the report names it, answering the sender with `memory_message_send` when a reply is due, so no handled message is listed again.
5. Author each Spec through the repository's Spec routing, consulting the Secondbrain and Exa before the PRD and TechSpec, with the answers in its `_authorization.md` and PRD.
6. Run `roundfix spec check <slug> --strict --run-verification` on each committed plan and fix what it reports.
7. Queue the Specs with `roundfix deliver` in delivery order and follow them to merge, and to release only when the user granted release.
8. Before the turn ends: arm a detached Run, the Delivery Queue or a scheduled wake-up, or say plainly that work stops and what remains; send each finding that belongs to another project as one self-contained ai-memory message to that project (`memory_message_send`), or, when the session exposes no ai-memory tools, put its text in the report for the user to deliver; and report what was done, what is queued or armed, which sources were missing and what still waits on the user.

## --unattended

Gather missing information, decisions, and explicit authorizations for push, merge and release, and for a release-announcement step only when the user's own instructions define one.

Ask at most three questions per call, never more than the tool accepts, as the stop-and-ask clause of `docs/agents/agent-instructions.md` allows, and put only independent questions in one call.

Answers and authorizations about an item that becomes a Spec go into that Spec's `_authorization.md`, and its decisions also go into the PRD's `## Decisions`.

An item that does not become a Spec now becomes a Backlog Entry or a Finding that records the information and decisions and never an authorization.

Each question offers 2-3 options with the recommended one first and `(Recommended)` in its label. Without a structured tool the same questions come numbered in the chat and the mode waits. A recommendation, silence, a dismissal or a timeout is never an answer. An authorization is granted only to a consuming Spec, so a later session asks again when it writes that Spec.

When run alone, arm the continuation as `--follow-up` step 8 does: arm a detached Run, the Delivery Queue or a scheduled wake-up, or say plainly that work stops and what remains, send another project's finding through ai-memory when available (or put it in the report when not), and report what was done, what is queued or armed, missing sources and what still waits on the user.
