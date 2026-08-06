---
status: done
created_at: 2026-08-06
updated_at: 2026-08-06
---

# Knowledge flow — findings accumulate faster than they become Specs (2026-08-06)

Investigation prompted by the maintainer at the close of the v0.4.0 session:
findings pile up, cross-project feedback travels by hand, and the secondbrain
receives raw files instead of curated knowledge. Measured across the local
fleet and checked against the secondbrain corpus and external practice. The
proposal this evidence feeds lives with the knowledge-pipeline idea under
`docs/specs/`; this document records what is true today.

## 1. The corpus grows superlinearly and nothing closes it

- Symptom / evidence: `docs/findings/` holds 63 files, 500K — 28 created in all
  of July, 35 in the first six days of August. Fleet-wide the mirrors show 111
  findings across 8 projects (roundfix 60, vortex 26, oraculum 10, pantheon 8,
  tax-poc 3, gss 2, fiscus 1, fluxus 1). None of the 63 local files carries the
  lifecycle frontmatter, because the contract that defines it
  (`status: pending | partial | deferred | done`, in
  `docs/agents/docs-layout.md`) shipped **tonight**, with Spec 0075.
- Root cause: until 0075 there was no lifecycle, no rollup instrument, no
  archival path, and there is still no enforcement — a finding is written once
  and never revisited. Google's SRE workbook names the failure mode exactly:
  rewarding postmortem writing without closing action items yields "an
  unvirtuous cycle of unclosed postmortems"
  (https://sre.google/workbook/postmortem-culture/). Organizing knowledge by
  dated report is the "project mindset" the Communication Patterns book warns
  about — knowledge filed by transitory container gets lost when the container
  ends (`wiki/concepts/communication-patterns-capitulo-14-10-knowledge-management-principles.md`).
- Action / suggestion: migrate the 63 to the 0075 contract, add a rollup
  instrument, an archival path, and `SC-*`-style enforcement — scoped in the
  knowledge-pipeline idea.

## 2. Cross-project findings already flow, by hand, and they collide with Active Runs

- Symptom / evidence: during this session the maintainer hand-pasted vortex
  feedback (the `StaleGateError` reproduction on task_08) and fiscus feedback
  (the `busy_timeout` confirmation) into the conversation; a spec-0002 findings
  addendum authored by another session rode along inside this session's commit
  because ride-along is the convention. The interference is measured: review
  Runs execute in the user checkout and Batch commits stage every path changed
  since the snapshot, so twice tonight findings edits had to wait for a Run to
  free the checkout, and once a dirty findings file blocked `resolve`
  Preflight outright.
- Root cause: no origin marking, no transport other than a human, and a
  destination convention (ride-along) that conflicts by construction with the
  destination's own Active Runs.
- Action / suggestion: the pipeline needs an explicit cross-project contract —
  origin recorded in frontmatter, destination always owns the commit, and a
  transport that cannot touch a checkout an Active Run owns. Candidates to
  weigh: the destination's `docs/_inbox/` with a triage step, a typed GitHub
  issue form (no repo file until triage), or a brain-mediated inbox — each has
  a different Run-interference profile.

## 3. The secondbrain receives the raw pile, not the knowledge

- Symptom / evidence: `.secondbrain-export` publishes `docs/` wholesale, so all
  60 roundfix findings mirror into `projects/roundfix/mirror/` uncurated. The
  brain's own architecture has a curated layer (`wiki/`) fed by deliberate
  ingest, and the mirror layer explicitly is "snapshot, replaced on merge" —
  findings bypass curation entirely.
- Root cause: export is path-based with no lifecycle awareness; an archived or
  superseded finding mirrors identically to an open one.
- Action / suggestion: archival must become a first-class state the export can
  see — either an `_archived/` move the mirror preserves distinctly, or
  status-aware export. The operational-knowledge literature converges on the
  same three-stage shape: capture → curation → retention, "a good system
  rejects most inputs", every published item carries owner, review date, and
  archive path
  (https://us.fitgap.com/stack-guides/automate-knowledge-capture-from-operational-systems-without-creating-a-data-swamp,
  https://prepared.cloud/knowledge-base-governance-checklist-owners-review-dates-and-archive-rules).

## 4. The consolidation instrument was already invented, ad hoc, by vortex

- Symptom / evidence:
  `projects/vortex/mirror/docs/findings/2026-07-27-roteiro-implementacao-findings-abertos.md`
  sequences 30+ open findings into phased implementation order, carries
  `status:` frontmatter, and states what blocks what — a rollup, written as
  just another finding because no home existed for it. Google's practice is the
  same instrument institutionalized: a central searchable repository (Requiem)
  plus periodic rollups ("weekly outage report", "greatest hits") and
  action items filed into the tracker rather than left in documents.
- Root cause: the fleet needs rollups and has no typed place for them.
- Action / suggestion: the pipeline should define the rollup as a first-class
  artifact — the map that supersedes its members and licenses their archival.

## 5. The intent door exists since tonight, and nothing routes through it

- Symptom / evidence: Spec 0075 shipped `docs/backlog/` — typed intent
  (`feat | fix | perf | refactor`), lifecycle (`open | promoted | declined`),
  promotion moves the entry into the adopting Spec's `references/`. But zero
  authorial skills reference it (`grep -rl backlog .agents/skills/` matches
  only an unrelated the-fool reference), so `write-idea`/`write-prd`/
  `brainstorming` neither consume `feat` entries nor spawn backlog entries from
  findings. The boundary is already right and matches the external evidence —
  "a ticket is not knowledge; it is evidence": findings are evidence and never
  commitments, backlog entries are intent and never evidence.
- Root cause: 0075 built the door this same night; routing was out of its
  scope, and its Non-Goals explicitly deferred finding migration.
- Action / suggestion: reinforcement lands in the authorial skills and agent
  guides (protected tooling — needs its named grant), and the finding→backlog
  spawn becomes the pipeline's standard closing move, mirroring SRE action
  items filed as bugs.

## What worked — keep

- The lifecycle states in the note-lifecycle literature map one-to-one onto
  contracts this repository already has: fleeting→permanent→archived with
  kind-determined defaults and advisory review
  (https://github.com/mrosnerr/open-zk-kb/blob/main/docs/note-lifecycle.md) is
  the findings `pending→done` contract plus the backlog `open→promoted`
  contract plus the missing archival state — the design is confirmation, not
  invention.
- The two-directory trust rule — one directory strictly "not yet real", one
  strictly current-state, incremental graduation, delete when empty
  (https://holdex.io/insights/spec-lifecycle) — is exactly how this
  repository's specs already archive; extending it to findings is symmetric.

## Addendum — 2026-08-06 — routed to Spec 0079

Adopted by [Spec 0079 — one door for fleet knowledge](../_idea.md),
which carries every action this investigation proposed. Status flipped to
`done` on adoption, per the findings contract.
