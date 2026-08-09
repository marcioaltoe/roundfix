---
spec: 0085-what-an-agent-reads-before-it-decides
status: active
created: 2026-08-09
surfaces: [backend, infra, docs]
---

# What an Agent reads before it decides

Two defects in the same surface: what an Agent is made to read, and what it is
required to consult.

**The archive sits inside the read path.** `docs/adr/` holds 105 decision
records with no structural separation between the 31 carrying an accepted
status, the 20 carrying only a legacy `Status: Accepted` body line, and the 53
carrying no status at all. Archived Specs live under `docs/specs/_archived/` and
archived findings under `docs/findings/_archived/`, so every consumer that wants
to exclude history must know each tree separately. Lifecycle markers do not fix
this: marking ADR-0106 superseded on 2026-08-09 left its bytes exactly where an
Agent loading decision context still reads them.

**The consultation rule has an escape hatch, and it is used.** The Secondbrain
clause is conditional — consult when repository context does not answer, do not
consult when local code, `CONTEXT.md`, ADRs, and repository documentation fully
answer the task. On 2026-08-09 a session formed a design decision about session
warm-up, wrote it into a TechSpec, and only then consulted the Secondbrain. The
consultation changed the design: its production contract holds that a transcript
is not proof of an action, which turned an in-memory check into a published
receipt. The decision reached an artifact before the source that corrected it
was opened.

This is a refactor plus a normative change. It moves no product behavior; it
changes what agents read and what they must read first.

## Project Constraints

- Identifier strategy: not applicable — no new entity or resource identity;
  Spec slugs, ADR numbers, and finding filenames keep the identities they carry.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local filesystem layout and static
  configuration only, with no network surface. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0106 is cited above only as the
  worked example of a marker that changes nothing about where an Agent reads;
  it is superseded by ADR-0108 and this Spec neither revives nor revises it.
  ADR-0081 makes derived pins rewritten by
  `make baseline-digests` sanctioned fallout of an authorized catalog edit;
  ADR-0092 keeps triage intent in a typed backlog, which is where both halves of
  this Spec were filed; ADR-0095 keeps the Secondbrain inbox the single capture
  door and is unaffected, because this Spec changes consultation and not
  capture; ADR-0093 and ADR-0094 govern the citation-based consistency check
  that reads the archive paths this Spec moves; ADR-0091, ADR-0096, and ADR-0097
  govern the authored QA gate; ADR-0104 requires an acceptance row on evidence
  this Spec did not author and holds pull request preparation until it is
  satisfied or carried forward. ADR-0117 places a check with the stage that can produce its defect; it does not change what this Spec delivers, and it moves where this Spec's gate rows run only once Spec 0093 ships. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — Baseline catalog modules, the source-baseline
  corpus and manifest, two setup-owned guides, and the review configuration are
  mutated. Express maintainer authorization:
  `docs/workflow/authorizations/2026-08-09-what-an-agent-reads-before-it-decides.md`.
  Bounded files: `internal/baseline/assets/modules/secondbrain.json`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/secondbrain.md`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/docs-layout.md`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`,
  `docs/agents/secondbrain.md`, `docs/agents/docs-layout.md`, and
  `.coderabbit.yaml`. Source: `docs/agents/agent-instructions.md`.

## Goals

1. One archive root holds every kind of retired artifact, so a consumer learns
   one path instead of one per tree.
2. Retired material leaves the directories an Agent loads by default, and stays
   reachable from the active record that replaced it.
3. Excluding history from review is one path filter, not one per tree.
4. Consulting the Secondbrain is unconditional before a design or architecture
   recommendation and before authoring an Idea, PRD, or TechSpec.
5. Every ADR carries a lifecycle status in one format, so active and retired are
   machine-distinguishable rather than eyeballed.

## Core Features

- **A single archive root.** `_archived/specs/`, `_archived/findings/`,
  `_archived/adr/`, `_archived/backlog/`, with the Archive Command, the
  `archive-spec` skill, and the consistency checker writing and reading there.
- **A forward pointer that survives the move.** An artifact leaving the read
  path names what replaced it, so moving history does not break the trail back
  to it.
- **An unconditional consultation trigger.** The Baseline clause stops offering
  a local-documentation exemption for design formation and for Idea, PRD, and
  TechSpec authoring.
- **One lifecycle format.** The 73 ADRs outside the frontmatter contract are
  normalized, so `status` means the same thing in every record.

## Non-Goals / Out of Scope

- Deleting any ADR, finding, or Spec. The archive is moved, not pruned; the
  graveyard is what answers "did we already try this?".
- Re-litigating the decisions the archived records contain. Moving a record is
  not revisiting it.
- Changing the Secondbrain's capture contract or its inbox. ADR-0095 stands.
- Changing what any Roundfix command does at Run time.

## Decisions

- History moves out of the read path rather than being marked in place, because
  stale documentation is more dangerous for an agent than for a human — it can
  be read as current, and a marker only helps a reader who notices it.
- Nothing is deleted, because 101 of the 105 ADRs are cited by Specs the
  repository requires to stay byte-identical.
- The consultation clause becomes unconditional at two named moments rather than
  always, so the rule stays enforceable instead of becoming background noise.
