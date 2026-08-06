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
status: open # open | promoted | declined
created: YYYY-MM-DD
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
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

- **mandatory**: Keep evidence and intent distinct in both directions: a finding records what happened and is never a commitment; it is evidence-backed, immutable history. A backlog entry records what to do next and is never evidence. A finding may spawn a backlog entry; a backlog entry needs no finding. A `feat` entry is upstream raw material that the spec pipeline may consume, never the `write-idea` artifact itself.

- **mandatory**: Give each documentation directory one job: `docs/_inbox/` for raw notes, `docs/adr/` for decisions, `docs/agents/` for agent guidance, `docs/design/` for design artifacts, `docs/backlog/` for dated, typed intent not yet committed to a Spec, `docs/findings/` for dated investigations, `docs/handoffs/` for session continuity, `docs/references/` for external pointers and durable project reference documents, and `docs/user-guide/` for human documentation. Preserve repository-authored extensions outside setup markers.

- **mandatory**: Durable knowledge flows upstream only: the project glossary (`CONTEXT.md`) and the agent guides reference accepted ADRs and never reference `docs/specs/` or `docs/findings/` content. Findings are dated reports that become Specs, not reference material; a document meant as durable project reference belongs in `docs/references/`.

- **mandatory**: Use this complete copyable Findings Operational Contract:

```markdown
---
status: pending # pending | partial | deferred | done
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

- **mandatory**: Set `status: done` as soon as the Spec is created and linked; do not wait for its Tasks, QA, archive, or release.

- **mandatory**: Treat findings as immutable history: append evidence and routing links as dated addenda instead of rewriting the original observations.

- **mandatory**: Update `updated_at` whenever status changes or an evidence addendum is appended; keep `created_at` as the document creation date.

- **mandatory**: A Rollup is a Finding of `kind: rollup` that consolidates related Findings. It lives beside active Findings under `docs/findings/`, shares their lifecycle contract, and declares a non-empty `members:` list of Finding basenames. Every member must resolve under `docs/findings/` or `docs/findings/_archived/`. Use this extension:

```yaml
kind: rollup
members:
  - YYYY-MM-DD-<finding-slug>.md
```

- **mandatory**: Use `docs/findings/_archived/` as the archival home for Findings. Every archived Finding requires an `absorbed_by:` license that resolves to an active Rollup basename or a Spec slug. Use this extension:

```yaml
absorbed_by: <rollup-basename-or-spec-slug>
```

- **mandatory**: Triage resolves one pending Inbox Entry into exactly one Finding, one Backlog Entry, or one recorded discard. Preserve the ADR-0092 boundary: evidence never becomes intent without a human choice. A minted Finding or Backlog Entry must cite the Inbox Entry's provenance.

<!-- setup-context-driven:end id=guide.docs-layout -->

<!-- setup-context-driven:begin id=guide.spec-docs-layout version=0.0.1 -->

# Spec docs layout

- **mandatory**: Keep `_idea.md`, `_prd.md`, `_techspec.md`, `_tasks.md`, Task files, and `qa/` evidence under the Spec folder. Archive only completed Specs with a passing QA verdict under `docs/specs/_archived/`.

- **mandatory**: Specs are downstream results of the CONTEXT-driven workflow, never sources it depends on: an archived Spec may be deleted at any time, so durable knowledge a Spec produced must move upstream to its semantic owner — the project glossary, an accepted ADR, an agent guide, or `docs/references/` — before or at archive. The glossary and the agent guides must never reference a Spec.

<!-- setup-context-driven:end id=guide.spec-docs-layout -->
