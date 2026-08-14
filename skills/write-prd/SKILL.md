---
name: write-prd
description: Write a PRD as a durable local artifact at docs/specs/<slug>/_prd.md — research first, clarify with one multiple-choice question at a time, record decisions as ADRs, then write directly without draft-approval rounds.
argument-hint: "<feature description, or nothing after a grilling/brainstorm session>"
metadata:
  category: planning
  tags: [prd, product, requirements, workflow, documentation]
  version: 0.0.2
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
version: 0.0.2
---

# Write PRD

Produce `docs/specs/<feature-slug>/_prd.md` — the product definition the rest of the pipeline (`write-techspec`, `write-tasks`, `implement-task`, `qa-gate`) implements from. The PRD owns _what_ and _why_; the tech spec owns _how_. Keeping that boundary is what lets each artifact stay small and stay true while the code changes underneath it.

## Inputs

`$ARGUMENTS` is a feature description, or empty when the current conversation already explored the feature (a grilling or brainstorm session). Everything already answered in the conversation counts as answered — do not re-interview.

If `docs/specs/<slug>/_idea.md` exists (produced by `write-idea`), it **is** the exploration: read it first, and treat its research, scoring, council insights, chosen direction, and Out of Scope list as answered ground truth. Clarify only what the idea left open.

## Size triage first

Not every change earns a PRD — it pays for itself when there are product decisions to record and multiple tasks to derive. If the work changes no product behavior (a refactor, a bug fix), say so and route to `write-techspec`, which enters the pipeline directly and mints the spec folder with a minimal `_prd.md`. If the change is trivial (one-line fix, typo, config tweak), recommend direct implementation with no spec folder.

## Ground rules

- **Read `CONTEXT.md` and `docs/adr/` before anything else.** Use the glossary's vocabulary throughout — titles, user stories, feature names — and never drift to synonyms the glossary lists under `_Avoid_`. If a concept has no term yet, sharpen one with the user and add it to the glossary right then; a PRD written in fuzzy language produces fuzzy tasks. If either file is missing, proceed silently.
- **Decide, then write.** Once research and clarifications are done, write the file directly and let the user react to the finished artifact. Do not present outline drafts for approval — reviewing a real PRD is faster than reviewing a promise of one.
- **Durability.** No file paths, no code snippets, no line numbers, except the mandatory `docs/agents/` Project Constraint sources and exact files covered by tooling authorization. The PRD may sit in a queue for weeks while the codebase moves; describe behavior and interfaces, which survive refactors.

## Process

### 1. Read the pending inbox

After the ground-rule reads above — `CONTEXT.md` and the active ADRs — and
before exploration, read the repository's pending Inbox Entries from its
brain-side `inbox/<repository>/` namespace through the configured
knowledge-workspace workflow. Treat queued evidence and intent as inputs to the
product decision, and leave Triage to a session of the destination repository.

### 2. Research

Explore before asking anything — questions the codebase can answer are wasted user time:

- Existing behavior and adjacent features the change touches.
- Prior specs under `docs/specs/` and `docs/history/specs/` — overlap with something already built or planned is a finding to surface, not to silently absorb.
- Market/competitor context via web research when the feature is user-facing and positioning matters.

Resolve the effective Project Constraints from the repository's semantic
guides before clarification:

- `docs/agents/domain.md` owns identifier guidance and routes active ADR
  discovery.
- `docs/agents/backend.md` owns authentication and HTTP policy.
- `docs/agents/agent-instructions.md` owns universal Normative Clauses,
  including tooling authority.
- `docs/agents/spec-routing.md` owns the Spec workflow contract.

Read every path that exists and follow its links to the effective decision
values. Never infer a suggested default or treat a missing guide as
authorization. The finished snapshot must classify identifier strategy,
authentication and HTTP, active ADR obligations, and tooling authority as
applicable or not applicable with a reason and cite the operative
`docs/agents/` source path for each row.

### 3. Clarify

Ask **one question per message**, multiple-choice whenever the options are enumerable:

```text
Which failure behavior should an expired import have?

A) Retry automatically up to 3 times  ← suggested: matches the sync retry ADR
B) Fail fast and notify the user
C) Park it for manual review
D) Other — describe
```

Always state a suggested default and the one-line reason. Cover, in order of importance: goals and success criteria, functional scope, non-goals, constraints, risks. Stop asking when the remaining unknowns don't change what gets built.

### 4. Record decisions

A product decision that is hard to reverse, surprising without context, and the result of a real trade-off becomes an ADR at `docs/adr/NNNN-slug.md`, continuing the repository's numbering. Keep it to 1–3 sentences: context, decision, why. Decisions that fail that three-part gate just live in the PRD body.

### 5. Prepare the Spec folder and adopt relied-upon sources

Adoption transfers a source document into the Spec that commits to implementing
it. Resolve the numbered slug with the rule in step 6, then run
`mkdir -p docs/specs/<slug>/references` before the first move. Run these steps
in order after recording decisions and before writing the PRD:

1. **Inventory.** List every inbox note, finding, and backlog entry whose
   content the PRD relies on. A document cited only as background is not
   adopted.
2. **Classify.** A raw inbox note that is source material can move directly. A
   note that records observed behavior is field evidence: promote it to a
   finding first, then adopt the finding. A backlog entry this Spec implements
   is adopted as `backlog`: set its `status: promoted` and `spec` to this
   Spec's slug in the same change that moves it, so the entry's own frontmatter
   and the index agree.
3. **Claim ownership.** Search both `docs/specs/` and
   `docs/history/specs/` for an existing owner. Exactly one Spec owns a shared
   source: the first Spec that commits to implementation. A secondary Spec links
   the owner's post-adoption copy and adopts nothing.
4. **Preflight.** Resolve every adopted source's basename and destination under
   `docs/specs/<slug>/references/` before changing any finding status, writing
   the index, or moving a source. `_index.md` is a reserved basename; reject any
   adopted source with that basename. Reject duplicate source basenames and
   abort if any destination path already exists, including a symbolic link.
   Complete this check for the whole inventory before changing adoption state.
5. **Index.** Write the complete `_index.md` before the first status update or
   source move. Add one row per adopted source using this fixed Markdown table:

   ```markdown
   # Adopted sources

   | source | type | owner | adopted date | path |
   | --- | --- | --- | --- | --- |
   | docs/findings/2026-07-25-example.md | finding | 0060 | 2026-07-30 | 2026-07-25-example.md |
   ```

   `source` is the pre-adoption repository path and is never updated; it is the
   provenance record. `type` identifies the source as `inbox`, `finding`, or
   `backlog`.
   `owner` is the owning four-digit Spec number. `adopted date` is the adoption
   date, and `path` is the current path relative to `_index.md`, so it remains
   valid when the Spec archives. Before continuing, validate the complete index:
   require the fixed header, one complete row per inventoried source, unique
   `source` and `path` values, valid type, owner, and date values, and each
   `path` equal to its `source` basename. Creating this validated index first
   makes an interrupted adoption visible to `archive-spec`, which must reject
   the indexed but incomplete Spec instead of treating it as legacy.
6. **Flip then link.** Before moving a finding, set its frontmatter
   `status: done`, update `updated_at`, and record the owning Spec link in the
   same change. This ordering leaves the status and route visible at the old
   path in Git history.
7. **Move.** Run
   `git mv <source> docs/specs/<slug>/references/<basename>` once for each
   indexed source. Perform one move, never a copy and never a stub. The move
   preserves the basename and every byte. Step 8 may then change only Markdown
   link destinations inside the moved source; never rewrite its observations or
   other source content.
8. **Rewrite and gate.** Search the repository for links to each old path.
   Exclude `docs/history/specs/` from automatic link rewrites; archived Specs
   are immutable historical artifacts. Report links from archived Specs
   separately for explicit policy review. For every other linking file, rewrite
   the destination to the post-adoption path relative to that file and resolve
   every rewritten Markdown link target. This includes link destinations inside
   an adopted source; only Markdown link destinations may change after the
   byte-preserving move. Fail and do not report completion while an adopted
   source still exists at its `source` path or any rewritten link is unresolved.
   Name the offending source or link and repeat the skipped adoption step.

### 6. Write

**HARD RULE — spec folders are numbered `docs/specs/NNNN-<kebab-slug>/`** (zero-padded 4 digits, e.g. `0001-implement-command`). Determine `NNNN` by scanning **both** `docs/specs/` and `docs/history/specs/` for the highest existing prefix and adding 1; use `0001` when no specs exist anywhere. Numbers are never reused and travel with the spec when archived. Never create an unnumbered spec folder. When an `_idea.md` fed this PRD, its folder already carries the number — reuse it, don't mint a new one.

Write `_prd.md` in the Spec folder prepared in step 5, using the template in [references/prd-template.md](references/prd-template.md). If an `_idea.md` fed this PRD, flip its frontmatter `status` to `promoted`. Set the PRD frontmatter carefully — downstream skills parse it:

- `spec` — the folder slug.
- `status: active` — flipped to `archived` by `archive-spec` once the spec completes (every task done, QA passed).
- `surfaces` — every surface the feature touches (`frontend`, `backend`, `cli`, `data`, `infra`, `docs`). `qa-gate` routes browser-based QA from this list, so an omitted `frontend` means the feature ships without browser validation.

`Project Constraints` is body content, never frontmatter. When the PRD
proposes creating, editing, renaming, moving, or deleting protected tooling
configuration, scripts, ignore files, plugin declarations, or version pins,
stop until the maintainer gives express maintainer authorization. Record that
approval and the exact bounded files in the Tooling authority row. A generic
implementation request, setup completion, silence, or authorization without
bounded files does not authorize the mutation.

### 7. Report

Before reporting, re-read the finished artifact. You MUST NOT report completion
or recommend the next pipeline step unless `Project Constraints` is present;
all four rows state applicable or not applicable with a reason; every row cites
an operative `docs/agents/` source; and any protected tooling mutation records
express maintainer authorization plus bounded files. Keep authorization out of
frontmatter.

Then run the checker against the stage that produced the artifact:

```bash
roundfix spec check <slug> --stage prd
```

An error-level finding or a checker execution failure blocks the report. Fix
the PRD and re-run the command; do not report completion or recommend the next
pipeline step while either stands.

A clean PRD-stage result is not full Spec coverage. This stage does not decide
TechSpec coverage mapping or its Vocabulary Contract; Task Graph ADR
accounting, task coverage, context references, Verification independence,
requirement contradictions, or rehearsal declarations; commit-dependent
changed-path scope; rules that are not yet mechanical; or whether the product
goals are correct. Later authoring stages, the full unscoped sweep, and QA
retain those classes. Treat the checker's named skipped detectors as omitted,
not as clean findings.

Reply with the file path, any open questions that survived clarification, and the next step: `write-techspec` for features with architectural decisions to make, `write-tasks` directly when the technical approach is already obvious.

## Anti-patterns

- Technical design in the PRD — schemas, endpoints, package layout belong in `_techspec.md`.
- Re-asking what the conversation or the codebase already answered.
- Several questions in one message, or open-ended questions where options were enumerable.
- Inventing requirements the user never confirmed to make the document look complete.
- Padding user stories with variations nobody asked for — every story must trace to a confirmed need.

## References

- [references/prd-template.md](references/prd-template.md) — the full `_prd.md` template. Read it before writing the file.
