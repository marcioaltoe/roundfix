<!-- setup-context-driven:begin id=guide.docs-layout version=0.0.1 -->

# Docs layout

- **mandatory**: When creating a new ADR, prepend this repository-owned lifecycle overlay to the body contract in `.agents/skills/domain-modeling/ADR-FORMAT.md`:

```markdown
---
status: proposed # proposed | accepted | rejected | deprecated | superseded
created_at: YYYY-MM-DDTHH:MM:SSZ
updated_at: YYYY-MM-DDTHH:MM:SSZ
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# <Short title of the decision>

<One to three sentences describing the context, decision, and reason.>
```

- **mandatory**: Only `accepted` is active. Treat `proposed`, `rejected`, `deprecated`, and `superseded` ADRs as inactive.

- **mandatory**: Treat a legacy ADR without lifecycle frontmatter as active unless its body explicitly marks it inactive. Do not rewrite existing ADRs solely to adopt lifecycle metadata.

- **mandatory**: Name backlog entries `YYYY-MM-DD-<kebab-slug>.md`. Use this complete copyable Backlog Operational Contract:

```markdown
---
type: feat # feat | fix | perf | refactor
status: open # open | promoted | declined | done | deprecated | superseded | closed | cancelled
created: YYYY-MM-DD
spec: null # Spec slug when status: promoted
reason: null # required for terminal closure without a consuming Spec
---
```

For `type: feat`, use:

```markdown
# <Title — the intent in one line>

## Opportunity

<What could exist and for whom.>

## Value

<Why it would matter; the hypothesis.>

## Shape

<The rough form of a solution, explicitly non-binding.>
```

For `type: fix`, use:

```markdown
# <Title — the defect in one line>

## Symptom

<What misbehaves, as a user or operator sees it.>

## Where

<The surface, command, or package, as known.>

## Expected

<The behavior that should replace it.>

## Evidence

<A finding link when one exists; `none yet` is honest.>
```

For `type: perf`, use:

```markdown
# <Title — the cost in one line>

## Slow

<What is slow, for whom, and in which operation.>

## Measured

<The number that says so and how it was measured.>

## Target

<The number that would settle it.>
```

For `type: refactor`, use:

```markdown
# <Title — the tangle in one line>

## Tangled

<What resists change, and where it is duplicated or coupled.>

## Cost

<What it makes slow, risky, or wrong to touch.>

## Shape

<The structure that would replace it, explicitly non-binding.>
```

Keep `open` entries in `docs/backlog/`. When a Spec adopts an entry, set `status: promoted` and `spec` to that Spec's slug, then move the entry to `docs/specs/<slug>/references/`; git history remains the discovery trail. Set `status: declined` only with a non-null `reason`. The current type set is open: a new type must be a Conventional Commits type that expresses intent. Adding a type is a contract change that requires a corpus re-record, never an informal addition. Use `refactor` as the canonical token, never an abbreviation.

`docs/backlog/` holds unresolved `open` intent only. A `promoted` entry is adopted source material and moves to its consuming Spec's `references/` under the adoption contract; it is not a terminal implementation verdict. Every terminal Backlog Entry (`declined`, `done`, `deprecated`, `superseded`, `closed`, or `cancelled`) moves to `docs/history/backlog/` in the same operation that records its terminal status. Preserve its original intent, the true disposition, any consuming Spec and the reason for closure. No absorption license is required for a Backlog Entry. An unknown status is a declaration error, never an inferred terminal disposition.

- **mandatory**: Keep evidence and intent distinct in both directions: a finding records what happened and is never a commitment; it is evidence-backed, immutable history. A backlog entry records what to do next and is never evidence. A finding may spawn a backlog entry; a backlog entry needs no finding. A `feat` entry is upstream raw material that the spec pipeline may consume, never the `write-idea` artifact itself.

- **mandatory**: Give each documentation directory one job: `docs/_inbox/` for raw notes, `docs/adr/` for decisions, `docs/agents/` for agent guidance, `docs/design/` for design artifacts, `docs/backlog/` for dated, typed intent not yet committed to a Spec, `docs/findings/` for dated investigations, `docs/handoffs/` for session continuity, `docs/references/` for external pointers and durable project reference documents, and `docs/user-guide/` for human documentation. Preserve repository-authored extensions outside setup markers. Each of these directories holds live work only: retired content moves under the single history root `docs/history/`, family by family (`docs/history/backlog/`, `docs/history/findings/`, `docs/history/handoffs/`). Handoffs archive only on the user's explicit confirmation, and when confirmed, every handoff archives together — the active `docs/handoffs/` directory is the capture door, never a shelf. Findings, Backlog Entries and Rollups with a terminal lifecycle status must leave their active directory for the matching history family; preserve valid provenance and update dependent references as part of that transition. Spec-owned adopted references retain their Spec ownership and archive with the Spec; this rule clears active family directories without dismantling self-contained Specs.

- **mandatory**: Durable knowledge flows upstream only: the project glossary (`CONTEXT.md`) and the agent guides reference accepted ADRs and never reference `docs/specs/` or `docs/findings/` content. Findings are dated reports that become Specs, not reference material; a document meant as durable project reference belongs in `docs/references/`.

- **mandatory**: Use this complete copyable Findings Operational Contract:

```markdown
---
status: pending # pending | partial | deferred | done | deprecated | superseded | closed | cancelled
created_at: YYYY-MM-DD
updated_at: YYYY-MM-DD
---

# <Area> — <short title> (YYYY-MM-DD)

<Two to four sentences describing the session or investigation, the attempted outcome, and links to adjacent evidence.>

## 1. <Finding title — symptom, not hypothesis>

- Symptom / evidence: <observed behavior, command output, identifiers, and paths needed to reproduce>
- Root cause: <proven cause, or `unknown` with what was ruled out>
- Action / suggestion: <fix, mitigation, or route to a Spec, direct change, or upstream report>

## 2. <Next finding>

## What worked — keep

<Optional evidence about behavior worth preserving.>

## Addendum — YYYY-MM-DD — <short title>

<Append new evidence, root-cause proof, status context, or routing links without rewriting the original observation.>
```

- **mandatory**: Use `pending` when the finding is new and has no implementation Spec.

- **mandatory**: Use `partial` when a linked Spec covers only the selected implementation scope. Record the reason the remaining observations are unnecessary and link the covering Spec.

- **mandatory**: Use `deferred` only when the finding will not be implemented. Record the reason for deferral.

- **mandatory**: Set `status: done` as soon as the Spec is created and linked; do not wait for its Tasks, QA, archive, or release. Complete source adoption or move the terminal Finding to `docs/history/findings/` in that same lifecycle operation, preserving its valid absorption pointer. `done` records routing completion, not verified implementation.

- **mandatory**: Treat findings as immutable history: append evidence and routing links as dated addenda instead of rewriting the original observations.

- **mandatory**: Update `updated_at` whenever status changes or an evidence addendum is appended; keep `created_at` as the document creation date.

- **mandatory**: A Rollup is a Finding of `kind: rollup` that consolidates related Findings and shares their lifecycle contract. Only a Rollup with unresolved work belongs under `docs/findings/`; a terminal Rollup belongs under `docs/history/findings/`. It declares a non-empty `members:` list of Finding basenames. Every member must resolve under `docs/findings/` or `docs/history/findings/`. Use this extension:

```yaml
kind: rollup
members:
  - YYYY-MM-DD-<finding-slug>.md
```

- **mandatory**: The active `docs/findings/` directory holds unresolved Findings and unresolved Rollups only. Every terminal Finding or Rollup still in an active family directory (`done`, `deferred`, `deprecated`, `superseded`, `closed`, or `cancelled`) moves to `docs/history/findings/` in the same operation that records its terminal status; do not wait for the consuming Spec's implementation, QA, merge, or release. Preserve the original observation and record the true disposition and its reason in a dated addendum. An unknown status is a declaration error, never an inferred terminal disposition. A Spec-owned adopted reference retains its Spec ownership and travels with that Spec at archive; do not extract it from `references/` merely because its preserved source status is terminal. When absorption actually occurred, the archived Finding or Rollup retains a valid `absorbed_by:` pointer to the active or archived Spec that absorbs its content, or to an unresolved active Rollup:

```yaml
absorbed_by: <active-rollup-basename-or-active-or-archived-spec-slug>
```

When any terminal Finding has no consuming Spec or Rollup, record a non-empty `closure_reason` and `closure_evidence` source locator instead of inventing an absorber. Preserve its original observation and dated disposition; absence of a Spec never prevents a legitimate terminal record from entering history. An invalid existing `absorbed_by` must be repaired, not hidden by these closure fields.

Before retiring a Rollup, transfer each member's absorption pointer to its true remaining owner and record the retiring Rollup's own absorption or evidenced terminal closure. Preserve the prior routing in a dated addendum and validate all member and absorption references before completing the move. Referenced history is a migration obligation, not an exemption that keeps a terminal Rollup active indefinitely. Do not invent an owner or rewrite original evidence merely to satisfy the archive check.

- **mandatory**: Read a findings directory holding only live work as `health`, not loss: Rollups and `docs/history/findings/` hold what was learned. Do not restore absorbed Findings merely to repopulate the active directory.

- **mandatory**: When a Rollup has no unresolved work of its own because its residuals were absorbed by Specs or received an evidenced terminal disposition, record its terminal status, migrate its member links and archive it under `docs/history/findings/`. A linked implementation Spec may still be pending: closure of the routing record does not assert delivery of that Spec. Do not keep a terminal Rollup under `docs/findings/` solely because archived members previously named it.

- **mandatory**: Triage resolves one pending Inbox Entry into exactly one Finding, one Backlog Entry, or one recorded discard. Preserve the ADR-0092 boundary: evidence never becomes intent without a human choice. A minted Finding or Backlog Entry must cite the Inbox Entry's provenance.

- **mandatory**: When a Finding's lifecycle closes, mint each typed Backlog Entry that its recorded actions call for, while preserving the boundary between evidence and intent. Where a fleet observation is captured before it reaches this repository is the Secondbrain guidance's concern, not this one's.

<!-- setup-context-driven:end id=guide.docs-layout -->
<!-- setup-context-driven:begin id=guide.spec-docs-layout version=0.0.1 -->

# Spec docs layout

- **mandatory**: Keep `_idea.md`, `_prd.md`, `_techspec.md`, `_tasks.md`, Task files, and `qa/` evidence under the Spec folder. Archive Specs under the resolved archive root (`docs/history/specs/` for the built-in root) after their non-QA Tasks are completed and their QA meets the normal archive eligibility contract, or when an explicit user request or applicable prior user authorization permits a QA Archive Override. A generic archive request alone does not silently waive QA. Record the approval source, date, covered Spec/revision and actual QA outcome or absence; stamp `qa_override: true` and preserve any supplied reason. Do not ask again for an already applicable approval. The override may waive the named terminal QA gate's completion/evidence prerequisite for archive, but never changes its Task status or report verdict to completed/pass. Missing, failed, incomplete, stale or malformed QA remains recorded as observed. Non-QA Task completion, self-contained references, safe destination and execution ownership remain required. An archive override is not Implement Clean, QA success, review approval, permission to publish/merge/release or a waiver of required checks. Declared QA decline and qualifying declared-only partial evidence are evaluated under their own contracts, not automatically called overrides. A runtime command without override support must be reported as unsupported rather than given an invented flag.

- **mandatory**: Specs are downstream results of the CONTEXT-driven workflow, never sources it depends on: an archived Spec may be deleted at any time, so durable knowledge a Spec produced must move upstream to its semantic owner — the project glossary, an accepted ADR, an agent guide, or `docs/references/` — before or at archive. The glossary and the agent guides must never reference a Spec.

<!-- setup-context-driven:end id=guide.spec-docs-layout -->
