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

- **mandatory**: For an ADR that carries lifecycle frontmatter, only `accepted` is active. Treat `proposed`, `rejected`, `deprecated`, and `superseded` ADRs as inactive.

- **mandatory**: Treat a legacy ADR without lifecycle frontmatter as active unless its body explicitly marks it inactive. Do not rewrite existing ADRs solely to adopt lifecycle metadata.

- **mandatory**: Name backlog entries `YYYY-MM-DD-<kebab-slug>.md`. Use this complete copyable Backlog Operational Contract:

```markdown
---
type: feat # feat | fix | perf | refactor
status: open # open | promoted | declined | deferred | done | deprecated | superseded | closed | cancelled
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

Keep `open` entries in `docs/backlog/`. When a Spec adopts an entry, set `status: promoted` and `spec` to that Spec's slug, then move the entry to `docs/specs/<slug>/references/`; git history remains the discovery trail. Set `status: declined` or `status: deferred` only with a non-null `reason`. The current type set is open: a new type must be a Conventional Commits type that expresses intent. Adding a type is a contract change that requires a corpus re-record, never an informal addition. Use `refactor` as the canonical token, never an abbreviation.

`docs/backlog/` holds unresolved `open` intent only. A `promoted` entry is adopted source material and moves to its consuming Spec's `references/` under the adoption contract; it is not a terminal implementation verdict. Every terminal Backlog Entry is deleted in the change that records its disposition under the retirement clause. An unknown status is a declaration error, never an inferred terminal disposition.

- **mandatory**: Keep evidence and intent distinct in both directions: a finding records what happened and is never a commitment; it is evidence-backed, immutable history. A backlog entry records what to do next and is never evidence. A finding may spawn a backlog entry; a backlog entry needs no finding. A `feat` entry is upstream raw material that the spec pipeline may consume, never the `write-idea` artifact itself.

- **mandatory**: Give each documentation directory one job: `docs/_inbox/` for raw notes only when a repository has no Secondbrain, `docs/adr/` for decisions, `docs/agents/` for agent guidance, `docs/design/` for design artifacts, `docs/backlog/` for dated, typed intent not yet committed to a Spec, `docs/findings/` for dated investigations, `docs/handoffs/` for session continuity, `docs/references/` for external pointers and durable project reference documents, and `docs/user-guide/` for human documentation. Preserve repository-authored extensions outside setup markers. Each directory holds live work only. The history root `docs/history/` holds Archive Records, and Git keeps their full text. Spec-contained authorization records live at `<spec-root>/<slug>/_authorization.md`; the preserved legacy location `docs/workflow/authorizations/` remains readable for historical grants. When the runtime's archive leaves an Archive Record, the record names each Spec-owned adopted reference, and the reference leaves the tree with the Spec folder.

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

- **mandatory**: Set `status: done` as soon as the Spec is created and linked; do not wait for its Tasks, QA, archive, or release. Complete source adoption under the Spec's adoption contract. `done` records routing completion, not verified implementation; deletion follows the retirement clause.

- **mandatory**: Treat findings as immutable history: append evidence and routing links as dated addenda instead of rewriting the original observations.

- **mandatory**: Update `updated_at` whenever status changes or an evidence addendum is appended; keep `created_at` as the document creation date.

- **mandatory**: A Rollup is a Finding of `kind: rollup` that consolidates related Findings and shares their lifecycle contract. Only a Rollup with unresolved work belongs under `docs/findings/`; a terminal Rollup is deleted under the retirement clause. It declares a non-empty `members:` list of Finding basenames. Every member must resolve under `docs/findings/` or in Git. Use this extension:

```yaml
kind: rollup
members:
  - YYYY-MM-DD-<finding-slug>.md
```

- **mandatory**: A retirement deletes the file. A Finding, Rollup or Backlog Entry with a terminal status (`deferred` included), a handoff on the user's explicit confirmation, a finished Review Artifact, a rejected or deprecated ADR, and a superseded ADR once an accepted successor names it with its date are deleted in the change that records the disposition. Git keeps the text. A Spec that absorbs a source adopts it and its Archive Record names it; any other retirement names its reason in the commit. A citation of a deleted file becomes the Git Locator `roundfix spec check` prints. The runtime's `roundfix history sanitize` removes what an earlier layout left.

- **mandatory**: Triage resolves one pending cross-project message into exactly one Finding, one Backlog Entry, or one recorded discard. Preserve the boundary between evidence and intent: evidence never becomes intent without a human choice. A minted Finding or Backlog Entry must cite the message's provenance: the sending project, the message id and its date.

- **mandatory**: When a Finding's lifecycle closes, mint each typed Backlog Entry that its recorded actions call for, while preserving the boundary between evidence and intent. Where a fleet observation is captured before it reaches this repository is the cross-project continuity rule's concern, not this one's.

- **mandatory**: Before minting a Finding or Backlog Entry, including in Triage, look for an existing one that is not yet implemented and that the new observation or intent fits: an `open` Backlog Entry in `docs/backlog/` or an unresolved Finding in `docs/findings/`. Revise a fitting Backlog Entry in place under its existing name; extend a fitting Finding with a dated addendum, because a Finding is immutable history. Mint a new file only when none fits. A Triage that extends an existing entry resolves its message into that entry, which cites the message's provenance.

<!-- setup-context-driven:end id=guide.docs-layout -->
<!-- setup-context-driven:begin id=guide.spec-docs-layout version=0.0.1 -->

# Spec docs layout

- **mandatory**: Keep `_idea.md`, `_prd.md`, `_techspec.md`, `_tasks.md`, Task files, `_authorization.md`, and `qa/` evidence under the Spec folder. Archive Specs under the resolved archive root (`docs/history/specs/` for the built-in root) after their non-QA Tasks are completed and their QA meets the normal archive eligibility contract, or when an explicit user request or applicable prior user authorization permits a QA Archive Override. A generic archive request alone does not silently waive QA. Perform the override only through the runtime's archive override command, such as `roundfix archive <slug> --qa-override --approval <source> --reason <text>`; that command records the approval source, date, covered Spec revision and actual QA outcome or absence, stamps `qa_override: true` and preserves the reason. Never hand-edit the override stamp or any other archive front matter, because a hand edit skips the command's refusals and provenance. Do not ask again for an already applicable approval. The override may waive the named terminal QA gate's completion/evidence prerequisite for archive, but never changes its Task status or report verdict to completed/pass. Missing, failed, incomplete, stale or malformed QA remains recorded as observed. Non-QA Task completion, self-contained references, safe destination and execution ownership remain required. An archive override is not Implement Clean, QA success, review approval, permission to publish/merge/release or a waiver of required checks. Declared QA decline and qualifying declared-only partial evidence are evaluated under their own contracts, not automatically called overrides. A runtime command without override support must be reported as unsupported rather than given an invented flag. A runtime whose archive leaves an Archive Record, such as `roundfix archive <slug>`, writes `<archive-root>/<slug>.md` and removes the Spec folder in the same change; the folder stays in Git at the record's `source_revision`, so copy every file a later reader needs to `docs/references/`, an accepted ADR, the glossary, an agent guide or the Secondbrain inbox before or at archive. A Spec dropped before delivery is cancelled, never deleted by hand: a runtime cancel command, such as `roundfix cancel <slug> --reason <text>`, writes an Archive Record with the disposition `cancelled` that keeps the reason, the date and the Tasks completed and pending at that moment, returns each adopted Backlog Entry to `docs/backlog/` with `status: open` and each adopted Finding to `docs/findings/` with `status: pending`, and removes the Spec folder in the same change. It refuses while a Run or a delivery queue item is active for the Spec. A cancelled Spec never meets a prerequisite and never supersedes another Spec. A runtime without a cancel command must report cancellation as unsupported rather than deleting the folder.

- **prohibited**: Do not make a test, fixture, or build step read a file under the Spec root or assume that a Spec is still active; Specs archive and may be deleted. A check whose subject is the Spec artifacts themselves is exempt.

- **mandatory**: Specs are downstream results of the CONTEXT-driven workflow, never sources it depends on: an archived Spec may be deleted at any time, so durable knowledge a Spec produced must move upstream to its semantic owner — the project glossary, an accepted ADR, an agent guide, or `docs/references/` — before or at archive. The glossary and the agent guides must never reference a Spec.

<!-- setup-context-driven:end id=guide.spec-docs-layout -->
