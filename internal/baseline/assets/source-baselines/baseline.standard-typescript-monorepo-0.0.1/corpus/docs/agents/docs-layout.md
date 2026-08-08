<!-- source-baseline-entry: contract.docs-layout.directory-matrix -->
# Documentation directory matrix

| Directory | Purpose | Lifecycle |
| --- | --- | --- |
| `docs/_inbox/` | Raw, untriaged notes | Promote to the correct durable location, then remove the inbox item. |
| `docs/adr/` | Accepted architecture decisions | Append decisions; never reuse an identity. |
| `docs/agents/` | Agent-facing operating guidance | Regenerate setup-owned guides; preserve repository-owned guides. |
| `docs/design/` | Design and interaction artifacts | Keep while active; archive or prune superseded exploration. |
| `docs/findings/` | Dated investigations and field evidence | Preserve history and append dated evidence or routing updates. |
| `docs/handoffs/` | Resumable session state | Replace with newer handoffs and keep only useful recent history. |
| `docs/references/` | External references with local relevance | Prune references that no longer matter. |
| `docs/specs/` | Local Specs, Task Graphs, Tasks, and QA evidence | Archive only after Tasks complete and QA passes. |
| `docs/user-guide/` | Human-facing documentation and runbooks | Update with the behavior it documents. |
<!-- /source-baseline-entry: contract.docs-layout.directory-matrix -->

<!-- source-baseline-entry: contract.findings.template -->
# Findings template

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
```
<!-- /source-baseline-entry: contract.findings.template -->

<!-- source-baseline-entry: contract.findings.lifecycle -->
# Findings lifecycle

1. Create the dated findings file with `status: pending`, `created_at`, and `updated_at`.
2. Use `pending` while no implementation Spec exists.
3. Use `partial` when a linked Spec intentionally covers only selected observations; record why the remainder is unnecessary.
4. Use `deferred` only when the finding will not be implemented; record the reason.
5. Set `status: done` as soon as an implementation Spec is created and linked. Do not wait for Tasks, QA, archive, merge, or release.
6. Preserve the original observation. Append evidence, root-cause proof, and routing links as dated addenda.
7. Update `updated_at` for every status change or evidence addendum; never change `created_at`.
<!-- /source-baseline-entry: contract.findings.lifecycle -->

<!-- source-baseline-entry: clause.context.adr-01-template -->
When creating a new ADR, prepend this repository-owned lifecycle overlay to the body contract in `.agents/skills/domain-modeling/ADR-FORMAT.md`:

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
<!-- /source-baseline-entry: clause.context.adr-01-template -->

<!-- source-baseline-entry: clause.context.adr-02-active-status -->
Only `accepted` is active. Treat `proposed`, `rejected`, `deprecated`, and `superseded` ADRs as inactive.
<!-- /source-baseline-entry: clause.context.adr-02-active-status -->

<!-- source-baseline-entry: clause.context.adr-03-legacy-compatibility -->
Treat a legacy ADR without lifecycle frontmatter as active unless its body explicitly marks it inactive. Do not rewrite existing ADRs solely to adopt lifecycle metadata.
<!-- /source-baseline-entry: clause.context.adr-03-legacy-compatibility -->

<!-- source-baseline-entry: clause.context.docs-one-job-per-directory -->
Give each documentation directory one job: `docs/_inbox/` for raw notes, `docs/adr/` for decisions, `docs/agents/` for agent guidance, `docs/design/` for design artifacts, `docs/backlog/` for dated, typed intent not yet committed to a Spec, `docs/findings/` for dated investigations, `docs/handoffs/` for session continuity, `docs/references/` for external pointers and durable project reference documents, and `docs/user-guide/` for human documentation. Preserve repository-authored extensions outside setup markers.
<!-- /source-baseline-entry: clause.context.docs-one-job-per-directory -->

<!-- source-baseline-entry: clause.context.docs-upstream-flow -->
Durable knowledge flows upstream only: the project glossary (`CONTEXT.md`) and the agent guides reference accepted ADRs and never reference `docs/specs/` or `docs/findings/` content. Findings are dated reports that become Specs, not reference material; a document meant as durable project reference belongs in `docs/references/`.
<!-- /source-baseline-entry: clause.context.docs-upstream-flow -->

<!-- source-baseline-entry: clause.context.backlog-01-operational-contract -->
Name backlog entries `YYYY-MM-DD-<kebab-slug>.md`. Use this complete copyable Backlog Operational Contract:

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
<!-- /source-baseline-entry: clause.context.backlog-01-operational-contract -->

<!-- source-baseline-entry: clause.context.backlog-02-finding-boundary -->
Keep evidence and intent distinct in both directions: a finding records what happened and is never a commitment; it is evidence-backed, immutable history. A backlog entry records what to do next and is never evidence. A finding may spawn a backlog entry; a backlog entry needs no finding. A `feat` entry is upstream raw material that the spec pipeline may consume, never the `write-idea` artifact itself.
<!-- /source-baseline-entry: clause.context.backlog-02-finding-boundary -->

<!-- source-baseline-entry: clause.context.findings-01-frontmatter -->
Use this complete copyable Findings Operational Contract:

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
<!-- /source-baseline-entry: clause.context.findings-01-frontmatter -->

<!-- source-baseline-entry: clause.context.findings-08-rollup -->
A Rollup is a Finding of `kind: rollup` that consolidates related Findings. It lives beside active Findings under `docs/findings/`, shares their lifecycle contract, and declares a non-empty `members:` list of Finding basenames. Every member must resolve under `docs/findings/` or `docs/findings/_archived/`. Use this extension:

```yaml
kind: rollup
members:
  - YYYY-MM-DD-<finding-slug>.md
```
<!-- /source-baseline-entry: clause.context.findings-08-rollup -->

<!-- source-baseline-entry: clause.context.findings-09-archive -->
Use `docs/findings/_archived/` as the archival home for Findings. Every archived Finding requires an `absorbed_by:` license that resolves to an active Rollup basename or a Spec slug. Use this extension:

```yaml
absorbed_by: <rollup-basename-or-spec-slug>
```
<!-- /source-baseline-entry: clause.context.findings-09-archive -->

<!-- source-baseline-entry: clause.context.findings-10-live-work-health -->
Read a findings directory holding only live work as `health`, not loss: Rollups and `docs/findings/_archived/` hold what was learned. Do not restore absorbed Findings merely to repopulate the active directory.
<!-- /source-baseline-entry: clause.context.findings-10-live-work-health -->

<!-- source-baseline-entry: clause.context.findings-11-rollup-closure -->
When a Rollup has `no open members`, review it as a candidate for its own closure and close it through the existing Finding lifecycle contract when its work is settled.
<!-- /source-baseline-entry: clause.context.findings-11-rollup-closure -->

<!-- source-baseline-entry: clause.context.findings-02-pending -->
Use `pending` when the finding is new and has no implementation Spec.
<!-- /source-baseline-entry: clause.context.findings-02-pending -->

<!-- source-baseline-entry: clause.context.findings-03-partial -->
Use `partial` when a linked Spec covers only the selected implementation scope. Record the reason the remaining observations are unnecessary and link the covering Spec.
<!-- /source-baseline-entry: clause.context.findings-03-partial -->

<!-- source-baseline-entry: clause.context.findings-04-deferred -->
Use `deferred` only when the finding will not be implemented. Record the reason for deferral.
<!-- /source-baseline-entry: clause.context.findings-04-deferred -->

<!-- source-baseline-entry: clause.context.findings-05-done -->
Set `status: done` as soon as the Spec is created and linked; do not wait for its Tasks, QA, archive, or release.
<!-- /source-baseline-entry: clause.context.findings-05-done -->

<!-- source-baseline-entry: clause.context.findings-06-append-evidence -->
Treat findings as immutable history: append evidence and routing links as dated addenda instead of rewriting the original observations.
<!-- /source-baseline-entry: clause.context.findings-06-append-evidence -->

<!-- source-baseline-entry: clause.context.findings-07-update-timestamp -->
Update `updated_at` whenever status changes or an evidence addendum is appended; keep `created_at` as the document creation date.
<!-- /source-baseline-entry: clause.context.findings-07-update-timestamp -->

<!-- source-baseline-entry: clause.context.inbox-01-triage -->
Triage resolves one pending Inbox Entry into exactly one Finding, one Backlog Entry, or one recorded discard. Preserve the ADR-0092 boundary: evidence never becomes intent without a human choice. A minted Finding or Backlog Entry must cite the Inbox Entry's provenance.
<!-- /source-baseline-entry: clause.context.inbox-01-triage -->

<!-- source-baseline-entry: clause.context.inbox-02-fleet-flow -->
- MUST mint each typed Backlog Entry a closed Finding's recorded actions call for, preserving the boundary between evidence and intent. Where a fleet observation is captured before it reaches this repository is the Secondbrain guidance's concern.
<!-- /source-baseline-entry: clause.context.inbox-02-fleet-flow -->

<!-- source-baseline-entry: clause.spec.keep-artifacts-in-spec-folder -->
Keep `_idea.md`, `_prd.md`, `_techspec.md`, `_tasks.md`, Task files, and `qa/` evidence under the Spec folder. Archive only completed Specs with a passing QA verdict under `docs/specs/_archived/`.
<!-- /source-baseline-entry: clause.spec.keep-artifacts-in-spec-folder -->

<!-- source-baseline-entry: clause.spec.specs-are-downstream-artifacts -->
Specs are downstream results of the CONTEXT-driven workflow, never sources it depends on: an archived Spec may be deleted at any time, so durable knowledge a Spec produced must move upstream to its semantic owner — the project glossary, an accepted ADR, an agent guide, or `docs/references/` — before or at archive. The glossary and the agent guides must never reference a Spec.
<!-- /source-baseline-entry: clause.spec.specs-are-downstream-artifacts -->
