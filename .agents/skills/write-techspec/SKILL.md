---
name: write-techspec
description: Write the technical spec at docs/specs/<slug>/_techspec.md — explore the architecture, settle technical decisions one question at a time, record them as ADRs, and produce a build order that write-tasks can decompose. Entry point for feature work with an approved PRD, and also for refactors and bug fixes, where it mints the spec folder with a minimal _prd.md instead of requiring a product interview.
argument-hint: "<spec slug, path to docs/specs/<slug>/_prd.md, or a refactor/bugfix description>"
metadata:
  category: engineering-design
  tags: [architecture, documentation, workflow]
  version: 0.0.2
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
version: 0.0.2
---

# Write TechSpec

Produce `docs/specs/<slug>/_techspec.md` — the technical answer to the spec's requirements. The PRD said _what_ and _why_; this document decides _how_, _where_, and _with which_ — and its Build Order is what `write-tasks` turns into the task graph, so sequencing quality here becomes execution quality later.

## Preconditions — two entry modes

`$ARGUMENTS` names the spec (slug or path), or describes a refactor/bug fix. Pick the mode by whether product behavior changes:

- **Feature work** (product behavior changes) — `docs/specs/<slug>/_prd.md` must exist; if it doesn't, stop and point the user at `write-prd`. Product decisions need the product conversation first. If the PRD contains low-level technical decisions that belong here, surface that as a finding and propose relocating them rather than silently duplicating.
- **Refactor or bug fix** (no product behavior change) — this skill is the pipeline entry point; no PRD interview happens. When no spec folder exists yet, mint one following `write-prd`'s numbering rule (`<spec-root>/NNNN-<kebab-slug>/`, scanning both the configured Spec Root and its resolved archive directory — `docs/history/specs/` for the built-in `docs/specs` root, or `<spec-root>/_archived/` otherwise — for the highest prefix), and write a **minimal `_prd.md`** carrying only the contract downstream skills parse: the frontmatter (`spec`, `status: active`, `surfaces`), a problem statement, Project Constraints, goals, core features, and non-goals — engineering-framed, a few lines each. It exists so `write-tasks`, `qa-gate`, and `archive-spec` keep a single artifact contract; it is not a product document. If the "refactor" turns out to change product behavior, stop and route to `write-prd`.

## Ground rules

- **Read first**: `CONTEXT.md` (use the glossary vocabulary in every interface and component name), existing `docs/adr/` (respect prior decisions or challenge them explicitly — never contradict one silently), and the full PRD.
- **Decide, then write.** Settle decisions during clarification, record the significant ones as ADRs, then write the finished document. No outline-approval rounds.
- **YAGNI.** Do not propose new packages, layers, or directories when the feature fits as an addition to an existing module. Every new structural element must earn its place against a stated alternative.
- **Reference, don't duplicate.** Point at PRD sections by name ("covers User Stories 3–5"); repeating business context here creates two copies that drift.
- **Use the owned source path.** Link every adopted source at its
  post-adoption path in the owning Spec with a destination relative to `_techspec.md`,
  never at the pre-adoption path kept in the reference index. A secondary Spec
  links the primary owner's copy and adopts nothing.

## Process

### 1. Explore the architecture

Map the code the feature will live in: the modules it extends, the seams it can attach to, existing patterns for the same kind of work (an existing adapter, an existing job runner), and the test infrastructure available. Prefer existing seams to new ones — the ideal number of new seams is zero. Delegate the exploration to a subagent when one is available; the spec needs conclusions, not file dumps.

Resolve the effective Project Constraints from the repository's semantic
guides during exploration:

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

### 2. Clarify technical decisions

Same protocol as the PRD stage — **one question per message, multiple-choice with a suggested default and a one-line rationale** — but scoped to how/where/which: architecture placement, data model, API shape, migration strategy, testing seams, performance and security constraints. Only ask what exploration couldn't settle.

### 3. Record ADRs

Each significant technical decision (hard to reverse + surprising without context + real trade-off) becomes `docs/adr/NNNN-slug.md`, **continuing the same numbering the PRD stage used** — one decision log spans the whole feature. 1–3 sentences each.

### 4. Write

Write `_techspec.md` from the template in [references/techspec-template.md](references/techspec-template.md). Two sections carry the most downstream weight:

- **Every PRD goal and user story maps to a named technical component.** An unmapped story is a design hole; find it now, not during task execution.
- **Build Order** — numbered steps where every step after the first states which previous steps it depends on. `write-tasks` derives the task graph edges from this.

`Project Constraints` is body content, never frontmatter. When the design
proposes creating, editing, renaming, moving, or deleting protected tooling
configuration, scripts, ignore files, plugin declarations, or version pins,
stop until the maintainer gives express maintainer authorization. Record that
approval and the exact bounded files in the Tooling authority row. A generic
implementation request, setup completion, silence, or authorization without
bounded files does not authorize the mutation.

Keep interface sketches under 20 lines each; they document shape, not implementation. Target 1,500–2,500 words — a spec nobody reads protects nobody.

### 5. Report

Before reporting, re-read the finished artifact and any minimal PRD created by
this skill. You MUST NOT report completion or recommend `write-tasks` unless
each new artifact has `Project Constraints`; all four rows state applicable or
not applicable with a reason; every row cites an operative `docs/agents/`
source; and any protected tooling mutation records express maintainer
authorization plus bounded files. Keep authorization out of frontmatter.

Then run the checker against the stage that produced the artifact:

```bash
roundfix spec check <slug> --stage techspec
```

An error-level finding or a checker execution failure blocks the report. Fix
the PRD or TechSpec named by the finding and re-run the command; do not report
completion or recommend `write-tasks` while either stands.

A clean TechSpec-stage result is not full Spec coverage. This stage does not
decide Task Graph ADR accounting, task coverage, context references,
Verification independence, requirement contradictions, or rehearsal
declarations; commit-dependent changed-path scope; rules that are not yet
mechanical; or whether the product goals work through their user-reachable
surfaces. The Task authoring stage, the full unscoped sweep, and QA retain
those classes. Treat the checker's named skipped detectors as omitted, not as
clean findings.

Reply with the file path, the ADRs created, any decisions still open, and the next step: `write-tasks`.

## Anti-patterns

- Restating PRD business context instead of referencing it.
- Speculative generality — abstractions for requirements the PRD doesn't contain.
- Interface sketches that grow into implementations.
- A Build Order without dependency statements — that forces `write-tasks` to guess the graph.
- Contradicting an existing ADR without naming it and proposing to supersede it.

## References

- [references/techspec-template.md](references/techspec-template.md) — the full `_techspec.md` template. Read it before writing the file.
