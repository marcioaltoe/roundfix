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

- **mandatory**: Give each documentation directory one job: `docs/_inbox/` for raw notes, `docs/adr/` for decisions, `docs/agents/` for agent guidance, `docs/design/` for design artifacts, `docs/findings/` for dated investigations, `docs/handoffs/` for session continuity, `docs/references/` for external pointers and durable project reference documents, and `docs/user-guide/` for human documentation. Preserve repository-authored extensions outside setup markers.

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

<!-- setup-context-driven:end id=guide.docs-layout -->

<!-- setup-context-driven:begin id=guide.spec-docs-layout version=0.0.1 -->

# Spec docs layout

- **mandatory**: Keep `_idea.md`, `_prd.md`, `_techspec.md`, `_tasks.md`, Task files, and `qa/` evidence under the Spec folder. Archive only completed Specs with a passing QA verdict under `docs/specs/_archived/`.

- **mandatory**: Specs are downstream results of the CONTEXT-driven workflow, never sources it depends on: an archived Spec may be deleted at any time, so durable knowledge a Spec produced must move upstream to its semantic owner — the project glossary, an accepted ADR, an agent guide, or `docs/references/` — before or at archive. The glossary and the agent guides must never reference a Spec.

<!-- setup-context-driven:end id=guide.spec-docs-layout -->
