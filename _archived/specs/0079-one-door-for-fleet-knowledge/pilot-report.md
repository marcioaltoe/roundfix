# Pilot report — fleet knowledge capture door

The 2026-08-06 pilot exercised the save–retrieve–triage cycle with the
Secondbrain repository and this Roundfix checkout. Both captures became
durable inside the one-minute target, capture left the destination checkout
untouched, and the intent entry resolved to one typed Roundfix Backlog Entry.

## External gate

The gate was open before capture:

- `/Users/marcio/dev/secondbrain/inbox/` existed with `roundfix/` and
  `secondbrain/` destination namespaces.
- `/Users/marcio/dev/secondbrain/AGENTS.md` named `inbox/<dest>/` as the
  create-only territory for any fleet session and pointed to
  `inbox/README.md` for the entry contract.
- The advisory qmd collection answered the required query after filesystem
  permission was granted to the Agent Session.
- The Secondbrain worktree had no changed paths before the first capture.

qmd check result: the query `fleet knowledge inbox save retrieve triage
research capture roundfix` returned results. Its 0.88 top result named the
stale path
`qmd://projects/roundfix/mirror/docs/specs/archived/0030-context-driven-agent-instructions/prd.md`,
which did not exist in the current mirror and concerned older general context
guidance rather than this pilot's research digest. The remaining results were
generic inbox indexes or unrelated project material. No substantive existing
page could be extended, so the pilot created a new research entry.

## Capture 1 — research digest

- origin: roundfix
- destination: secondbrain
- capture: manual
- brain-relative path:
  `inbox/secondbrain/2026-08-06-fleet-knowledge-door-research-digest.md`
- brain commit: 93bd7b35b77a356a3d2d50989de93565b1e85dbb
- capture_to_durability_seconds: 50
- target: met; 10 seconds below one minute

The entry digests the adopted Spec investigation, the ADR-0092 evidence versus
intent boundary, and ADR-0095's single-door decision. It cites the three
Roundfix source paths plus the SRE postmortem-culture and operational-knowledge
sources used by the adopted investigation. The entry remains pending for the
Secondbrain's own ingestion process.

## Capture 2 — cross-project intent

- origin: secondbrain
- destination: roundfix
- capture: manual
- brain-relative path:
  `inbox/roundfix/2026-08-06-automate-inbox-capture-durability.md`
- brain commit: e982c0137ad61b5cf82de0fa8f63b9c5c325412e
- capture_to_durability_seconds: 24
- target: met; 36 seconds below one minute

The observation arose while operating the Secondbrain door: the first capture
needed manual file creation, exact-path staging, staged-diff review, commit-hash
retrieval, and elapsed-time calculation. It carries `type-hint: feat` because
the destination must decide whether an atomic helper becomes a future product,
skill, or workflow capability.

## Non-interference audit

Immediately before triage, the destination checkout
`/Users/marcio/dev/roundfix` was on branch
`ma/0079-one-door-for-fleet-knowledge`. The command
`rtk git -c core.fsmonitor=false status --short --untracked-files=all`
produced no paths. Both captures existed only in the Secondbrain at that
point; neither capture touched a Roundfix checkout.

The active Task worktree already carried the Daemon-owned `pending` to
`in_progress` status change in `task_05.md`. Capture added no path there; the
Backlog Entry and this report appeared only after destination triage began.

## Retrieval and triage

The Roundfix session retrieved the intent entry from brain commit
`e982c0137ad61b5cf82de0fa8f63b9c5c325412e` and made the explicit ADR-0092
choice to preserve it as intent, not evidence. Triage minted exactly one
artifact:

- typed artifact: `docs/backlog/2026-08-06-atomic-inbox-capture-helper.md`
- type: feat
- status: open
- provenance:
  `inbox/roundfix/2026-08-06-automate-inbox-capture-durability.md`

The brain-side entry moved to
`inbox/roundfix/_triaged/2026-08-06-automate-inbox-capture-durability.md` and
gained this resolution field:

`resolved_to: docs/backlog/2026-08-06-atomic-inbox-capture-helper.md`

The triage move is durable in brain commit
`4aa955c3ec2c03ef5799ea0745bc8a8938ce6530`. No Finding was created; findings
are outside this Task's slice.

## Friction before obligations bind

- The first sandboxed qmd invocation failed with `SQLITE_CANTOPEN` because the
  Agent Session could not open qmd's database. The same query succeeded after
  explicit filesystem permission. A binding workflow must name this access
  requirement or provide an authorized execution path.
- qmd warned that 248 documents still needed embeddings and returned a stale
  path as its highest-scoring result. Duplicate checks remain advisory and
  require path existence plus substantive review; score alone cannot decide
  extension versus creation.
- The first manual capture consumed 50 of the 60 available seconds. The second
  took 24 seconds after the access path was established. The open Backlog Entry
  records the non-binding intent for an atomic helper that reports path, commit,
  and elapsed time.
- Capture and triage required three narrowly scoped brain commits: one per
  capture and one for the positional move plus `resolved_to`. Any future helper
  must retain exact-path staging so unrelated brain work cannot ride along.

## Not exercised

The pilot did not ingest the research digest into the curated wiki, automate
triaged-entry retention, create a Finding, or add a capture CLI. Those actions
belong to the Secondbrain ingestion contract or later work.
