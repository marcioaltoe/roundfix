---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-12
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The Secondbrain guide permits capture without ever triggering it

## Symptom

The `guide.secondbrain` module — shipped as `docs/agents/secondbrain.md` in the
consuming repository — does not make capture happen.

An autonomous session in `vortex` on 2026-08-07/08 delivered three Specs,
produced actionable feedback about **Roundfix** with numbers, derived five
authoring rules belonging to the baseline guides, and **wrote zero entries to
`inbox/`** until the maintainer asked why. The mechanism existed and was in use:
`inbox/roundfix/` already held pending and triaged entries. It was not missing
infrastructure.

Reading the guide explains it. Almost every clause contracts **consumption**:
read `wiki/index.md`, run `qmd query`, inspect `projects/<p>/mirror/` only when
index and query point there, cite by path every file used, never claim context
that was not read. The only clause contracting **production** for the general
case is permissive:

> **mandatory**: Sessions MAY create files under the Secondbrain's `inbox/**`;
> this is the only writable Secondbrain namespace.

It grants permission and defines no trigger. "MAY create" never fires on its own.
The one production clause with a real trigger is narrow — it covers substantive
external research — and does not cover the most common case: a session in project
A observes something actionable whose owner is project B. The framing compounds
it: the guide opens with "Consult the local Secondbrain before acting", and the
consumer's `CLAUDE.md` lists the Secondbrain under *optional knowledge sources*,
last in the hierarchy. Everything pushes toward reading.

The cost was not theoretical. In that session the feedback was first written as a
finding in the `vortex` checkout — **with an active Run on the branch** — which is
precisely what the inbox exists to avoid, and the session did it that way because
nothing instructed otherwise.

## Where

The `secondbrain` Baseline module under `internal/baseline/assets/modules/`, which
renders `docs/agents/secondbrain.md`, and the autonomous-loop completion contract
that lists archiving, opening a Pull Request and merging.

## Expected

A production clause with a trigger, in the same register as the others:

> **mandatory**: A session that observes a defect, friction, or improvement whose
> owner is another project MUST create a pending Inbox Entry in that project's
> namespace before the session ends, and MUST NOT record it as a finding in the
> current repository. The current repository records only what it owns.

And the routing consequence, which is written nowhere today:

> **mandatory**: Guides delivered as `setup-context-driven` blocks are
> baseline-owned. Proposed changes to them are Inbox Entries addressed to the
> baseline's project, never local edits: a local edit is overwritten by the next
> baseline update.

The second matters because `autonomous-work.md`, `backend.md`, `domain.md` and
`agent-instructions.md` are all baseline-owned in the consumer, and a session
wanting to improve an authoring rule would naturally try to edit them there.

Worth settling in the same work: naming capture in the autonomous loop's
completion contract, so that "delivered" includes "what this session observed
about other projects was captured". Without that, capture competes with the work
and loses.

## Evidence

Minted from the Inbox Entry
`inbox/roundfix/2026-08-08-guia-secondbrain-sem-gatilho-de-captura.md` in the
Secondbrain, captured from a `vortex` session. Related:
`docs/backlog/2026-08-06-atomic-inbox-capture-helper.md`, which asks for the
mechanism this clause would trigger.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
