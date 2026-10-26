---
spec: 0325-the-follow-up-mode-consumes-the-handoff-and-every-message
status: active
created: 2026-10-26
surfaces: [docs]
---

# The follow-up mode consumes the handoff and every message

Minimal PRD for a bug fix, minted by `write-techspec`. The Roundfix skill's
`--follow-up` Skill Mode reads the project's ai-memory messages but consumes
them only in its fourth step, so a session that stops earlier leaves every
message it read pending, and it never claims the project's session handoff,
so the previous session's next steps never reach the plan. After this Spec the
mode accepts the project's most recent open handoff during its inventory,
hands the handoff's next steps and open questions to the follow-up command
beside the messages, and leaves no message it handed in pending when the turn
ends, also when it stops early. The change is skill text pinned by a contract
test, and the glossary gains **Session Handoff**.

## Sources

One Backlog Entry, written and adopted at birth by this Spec's planning
change (see `references/_index.md`):

1. [`references/2026-10-26-the-follow-up-mode-leaves-the-handoff-and-handed-in-messages-pending.md`](references/2026-10-26-the-follow-up-mode-leaves-the-handoff-and-handed-in-messages-pending.md),
   `type: fix`, from the maintainer's request in chat of 2026-10-26 at about
   17:30 BRT.

No open Backlog Entry or unresolved Finding shares its context: at the
authoring base `docs/backlog/` and `docs/findings/` hold no entry.

## Prerequisites

None in `requires`. Authored against `feeffa28`, with Specs 0320–0324 active
under `docs/specs/`; this Spec delivers after 0322 and before 0323. Specs
0321, 0322 and 0323 raise the `roundfix` skill version and edit
`skills/testdata/owned-skill-versions.json`; this Spec's Task raises the next
version and records whatever the tree holds then. No active Spec edits
`references/modes.md`, the `## Modes` section of the skill or
`skills/roundfix_modes_test.go`, and none adds a `CONTEXT.md` entry next to
**Cross-project Message** or edits **Follow-up Item**.

## Project Constraints

- Identifier strategy: not applicable — the handoff entries reuse the free
  `id` field of the existing `roundfix/follow-up-messages/v1` schema, and no
  record schema or identifier is minted in code. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added; Roundfix's code still never talks to ai-memory, and
  the skill calls the MCP tools the session exposes. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0330 (ai-memory carries handoffs
  and messages; a project consumes a message only after recording it, or
  after naming it when it only informs; Roundfix's code never talks to
  ai-memory) stays unchanged and is followed. ADR-0320 (the Skill Modes and up
  to three questions per call) stays unchanged: only the follow-up mode's
  inventory, triage and report change. ADR-0338, new, records the decision.
  Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — protected paths limited to the `roundfix`
  skill. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0325-the-follow-up-mode-consumes-the-handoff-and-every-message/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/modes.md`, `skills/roundfix/SKILL.md`.

## Goals

- The `--follow-up` mode takes the project's Session Handoff in its
  inventory: it uses the block the SessionStart hook already claimed, or lists
  the open handoffs, reads a newer cut-off session first, and accepts the
  project's most recent open handoff.
- The handoff's next steps and open questions reach the follow-up command as
  message entries, are ordered with every other item and appear in the plan
  and the report.
- No message the mode handed in stays pending when the turn ends, also when
  the mode stops at step 1, when the user stops it or when no answer comes.
- `--status` and `--unattended` accept no handoff and consume no message.
- A contract test pins this text in the canonical skill and its mirror, and
  `CONTEXT.md` defines **Session Handoff**.

## Core Features

1. Handoff in the inventory. Step 1 of `--follow-up` in
   `.agents/skills/roundfix/references/modes.md` gains the handoff sequence
   and the `handoff:<handoff id>:<n>` message entries; `Optional sources`
   says ai-memory brings the handoff and the messages, how each call is
   scoped, and that `--status` and `--unattended` accept and consume nothing.
2. Every message settled. Step 4 states the three settlements (recorded,
   informational, not yet handled) and that they also run when the mode stops
   early; step 8's report names the accepted handoff, what became of each of
   its items, and each consumed message and its record.
3. The skill summary. The `## Modes` bullets of the skill entry say the same
   in one line each.
4. The contract. `TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage`
   in `skills/roundfix_modes_test.go` pins the phrases in both layouts.
5. The glossary. **Session Handoff** is added, and **Follow-up Item** names
   the handoff's items.

## User Experience

A maintainer types `/roundfix follow-up` in a repository whose previous
session left a handoff and whose inbox holds three messages, one of them
informational. The report lists the handoff it accepted and its two next
steps among the plan's items, names the informational message as consumed,
and names the Backlog Entry that now records each of the other two after the
maintainer stopped the session before step 4. The next session finds no open
handoff and an empty inbox.

## Non-Goals / Out of Scope

- A change to the `roundfix follow-up` command, its flags, its messages schema
  or its output; the handoff entries ride the existing schema.
- A change to the Baseline's ai-memory clauses: the message clause already
  asks to consume only after recording in a Spec, a Backlog Entry or a
  Finding, or after naming an informational message, and the handoff clause
  already routes handoffs through `memory_handoff_*` (see `_techspec.md` →
  Decisions).
- Writing a new handoff at the end of the follow-up turn; the ai-memory
  handoff skill keeps owning when one is written.
- Any ai-memory call from Roundfix's Go code.

## Success Metrics

1. Success Metric: step 1 of `--follow-up` in both layouts names
   `memory_handoff_list`, `newer_unended_session` read with
   `memory_read_session_observations` before `memory_handoff_accept`, the
   `consumed_by_hook` case, and the `handoff:<handoff id>:<n>` entry with its
   `from`, `created_at` and empty `type`.
2. Success Metric: step 4 in both layouts states the early-stop rule and the
   three settlements, each pop naming its `message_id`, and step 8 reports the
   handoff and the consumed messages.
3. Success Metric: `Optional sources` no longer lists messages "without
   consuming them", says `--status` and `--unattended` accept no handoff and
   consume no message, and `--status` says it accepts no handoff.
4. Success Metric: the new contract test passes on the delivered tree and
   fails for each sabotage of `_techspec.md` → Testing Approach; the
   existing `skills` tests and `TestEveryOwnedSkillVersionIsRecorded` pass.
5. Success Metric: `CONTEXT.md` holds **Session Handoff** and the revised
   **Follow-up Item**, and the glossary contracts pass.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- ADR-0330, accepted on 2026-10-24: "a project reads its messages oldest
  first and consumes one only after recording its content in a Spec, a
  Backlog Entry or a Finding, or naming it when it only informs." Each
  settlement this Spec adds follows that order.
- The `ai-memory-handoff` skill vendored in this repository at
  `.agents/skills/ai-memory-handoff/SKILL.md` (upstream content, added by
  Spec 0317 from `marcioaltoe/skills` commit
  `65d5f365cace27c4b468db1723213723b10c870b`, which ADR-0330 records; also
  published at https://www.skills.sh/akitaonrails/ai-memory/ai-memory-handoff,
  read on 2026-10-26): "If the current context already contains a pending
  handoff block, answer from that block directly." and "`memory_handoff_list`
  names it as `newer_unended_session` when it is newer than every open
  handoff." Step 1's order follows it.
- The `ai-memory-messaging` skill vendored at
  `.agents/skills/ai-memory-messaging/SKILL.md`: "`memory_message_pop` claims
  ONE message exactly once — the oldest, or a specific `message_id`." and "A
  pop whose inbox the server would have to guess is refused and claims
  nothing." The ai-memory messaging documentation,
  https://github.com/akitaonrails/ai-memory/blob/main/docs/agent-messaging.md
  (read on 2026-10-26), says the same. The mode pops by `message_id` and
  scopes each call.

## Unreachable Acceptance

- criterion: after release, the maintainer's next `/roundfix follow-up` in a
  repository with an open handoff and pending messages accepts the handoff,
  names its items and leaves the inbox empty when the turn ends.
  reason: it needs a live ai-memory server and an interactive Agent session,
  which no Verification or QA row may reach.
  satisfied-by: the maintainer's or operator's first follow-up session after
  the release that carries this Spec, recorded session evidence.

## Validation

The decisions were checked on 2026-10-26 against these sources.

- Characterization at `feeffa28`: `modes.md` `Optional sources` says "list
  this project's messages without consuming them"; `memory_message_pop`
  appears only in step 4; no file under `.agents/skills/roundfix/` names
  `memory_handoff`. `skills/roundfix_modes_test.go`
  `TestTheFollowUpModeRunsItsStepsInOrder` requires step 4 to name
  `memory_message_pop` and step 8 to name `memory_message_send`, and forbids
  `resolved_to`, `_triaged`, `Inbox Entr` and `secondbrain.root`; the new
  text keeps all four constraints.
- `internal/cli/follow_up_inventory.go` `followUpMessages` accepts any unique
  non-empty `id`, an empty `type`, and uses `from` as the requester and the
  judge's origin, so a `handoff:<handoff id>:<n>` entry with this project as
  sender needs no code change.
- A prototype of the Task in a scratch worktree passed `go test ./skills`,
  `make skills-sync-check skills-check` and the full
  `-tags docscontract ./internal/docscontract`; each sabotage in
  `_techspec.md` → Testing Approach failed the new test.
- Secondbrain, read only: `wiki/index.md`, then `qmd query "ai-memory handoff
  accept consume messages follow-up" --all --files --min-score 0.3`, which
  returned the Roundfix mirror's ADR-0330 and Spec 0317 record, the
  `ai-memory-handoff` and `ai-memory-messaging` skills of the skills mirror,
  and `wiki/entities/ai-memory.md`. That page records the 2026-10-19 adoption
  of ai-memory for session continuity and its single-use proof (an accept
  returned `claimed`, a second `none_pending`), which supports claiming
  exactly one handoff per follow-up. None records a decision to leave
  messages pending.
- Exa: the published ai-memory handoff skill and the messaging document
  quoted under Acceptance evidence; the messaging document changed the
  design by adding the per-call scope and the pop by `message_id`.

## Glossary

- adds: **Session Handoff**
- changes: **Follow-up Item**

## Decisions

- 2026-10-26: Which messages handed to the follow-up may stay pending when the turn ends? — "Toda mensagem até o fim (Recomendado)": none; a recorded one becomes a Spec, Backlog Entry or Finding and is then consumed, an informational one is consumed and cited in the report, one not yet handled becomes a Backlog Entry and is consumed, also when the mode stops early.
- 2026-10-26: What does the follow-up do with the project's handoff? — "Aceitar e usar no plano (Recomendado)": list the open handoffs and accept the project's most recent one; its next steps and open questions become plan items shown in the report; a newer cut-off session is read first.
- 2026-10-26: Where does this Spec deliver? — "Spec 0325 antes da 0323 (Recomendado)": right after 0322, released in v0.125.0 with 0321 and 0322.

## Open Questions

None.
