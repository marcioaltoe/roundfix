---
status: accepted
created_at: 2026-10-22T22:00:00Z
updated_at: 2026-10-22T22:00:00Z
deprecated_at: null
superseded_by: null
---

# A question call carries up to three independent questions

Since 2026-09-08 the Baseline's stop-and-ask clause
`clause.core.ask-user-answerable-decisions` required exactly one question per
call or message, by the maintainer's choice, although every structured
question tool accepts several: Claude Code's `AskUserQuestion` one to four,
Codex's `request_user_input` one to three, OpenCode's `question` a list. On
2026-10-22 the maintainer asked for a follow-up session that gathers every
decision a repository's pending work needs before it continues alone, and
answered "Até 3 no geral" for the number of questions per call, for every
adopter.

## Decision

**Up to three, independent, in one pending call.** An Agent asks at most three
questions per call or message, never more than the tool in use accepts, and
keeps only that call pending. Questions share a call only when they are
independent of each other; a question whose options depend on another answer
waits for that answer. Every other part of the clause stays: the tool the
session really exposes, options with consequences, the recommended option
first with `(Recommended)`, custom answers, a numbered chat fallback when no
structured tool exists, and never treating a recommendation, silence, a
dismissal or a timeout as the user's choice. The clause keeps its identity and
its stop-and-ask enforcement, so a managed refresh records it `retained` and
rewrites the adopter's `docs/agents/agent-instructions.md`.

**The Roundfix skill gathers decisions that way and records each answer where
the work reads it.** The skill's `--unattended` Skill Mode, which its
`--follow-up` mode calls, lists the information, decisions and authorizations
the pending work needs and asks for them in such calls. An answer about an item
that becomes a Spec goes into that Spec's `_authorization.md`, with its
decisions also in the PRD's Clarification Record. An item that does not become
a Spec now becomes a Backlog Entry or a Finding that keeps the information and
decisions and never an authorization: an authorization is granted to a
consuming Spec, so a later session asks again when it writes that Spec.

## Considered Options

- Keeping one question per call was rejected by the maintainer: a session that
  must gather a dozen decisions before working unattended interrupts him a
  dozen times.
- Four per call, Claude Code's limit, was rejected: Codex accepts three, and
  one portable number serves every harness.
- Each tool's own limit was rejected for the same reason: the rule would read
  differently in each harness.
- Recording authorizations in a Backlog Entry was rejected by the maintainer:
  an authorization without a consuming Spec bounds no change, and a stale one
  would be read as current when the entry becomes a Spec.

## Consequences

- Adopters receive the reworded clause at their next managed refresh; no
  repository decision is asked and nothing else in their guides changes.
- A skill that interviews one question per message, because each answer
  shapes the next question, still complies: dependent questions wait under the
  new clause too.
- The `--follow-up` and `--status` modes read the inventory that Roundfix's
  read-only follow-up command produces; the skill adds what only the Agent can
  see, such as ai-memory messages, and never invents a command or flag the
  installed Roundfix lacks.
