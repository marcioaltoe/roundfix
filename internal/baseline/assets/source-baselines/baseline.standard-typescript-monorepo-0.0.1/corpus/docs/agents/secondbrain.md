<!-- source-baseline-entry: contract.secondbrain.protocol -->
# Secondbrain protocol

1. Consult the local knowledge system only when repository code, domain context, ADRs, and local documentation do not answer a business, prior-decision, cross-repository, or shared-architecture question.
2. Read its index first, then use its declared semantic query command. Open only the files required by the result and treat project mirrors as references, never as workspaces.
3. Sessions MAY create files under the knowledge system's `inbox/**`; every other knowledge-system path remains read-only. MUST NOT create, edit, rename, move, or delete raw ingestion data or project mirrors from this repository.
4. Pending Inbox Entries live at their destination namespace root and resolved entries live under `_triaged/`; each entry declares `origin:`, `destination:`, `type-hint:`, `created_at:`, and `capture:`, then gains exactly one of `resolved_to:` or `discarded_reason:` at triage time.
5. The maintainer-owned session-end hook writes `capture: auto` drafts that remain pending Triage and never triage themselves.
6. A session that performs substantive external research captures a sourced digest for the brain's namespace after an advisory duplicate check whose paths and substantive matches are reviewed; a score alone never decides extension versus creation.
7. MUST NOT read, copy, or expose secrets, session material, private keys, credentials, or unsafe personal and client data. Stop at likely secret-bearing sources and request a safe source.
8. Cite every knowledge-system file used in the final response or handoff. Do not claim knowledge-system context when no file was read.
9. Request an authorized ingestion workflow for durable additions or corrections instead of writing from this repository.
<!-- /source-baseline-entry: contract.secondbrain.protocol -->

<!-- source-baseline-entry: clause.secondbrain.01-consult-triggers -->
Consult the local Secondbrain before acting when repository context does not answer business or prior-decision questions, fiscal or tax concepts, cross-project documentation, knowledge about Vortex, Tax, Visio, or Gesttione, or shared architecture patterns. Do not consult it when local code, `CONTEXT.md`, ADRs, and repository documentation fully answer the task.
<!-- /source-baseline-entry: clause.secondbrain.01-consult-triggers -->

<!-- source-baseline-entry: clause.secondbrain.02-query-order -->
Read `wiki/index.md` first. Then run `qmd query "<question>" --all --files --min-score 0.3`. Inspect `projects/<project>/mirror/` only when the index and query point there, and open only the files required for the task; treat mirrors as references, not workspaces.
<!-- /source-baseline-entry: clause.secondbrain.02-query-order -->

<!-- source-baseline-entry: clause.secondbrain.inbox-write-permission -->
- MUST treat `inbox/**` as the only writable Secondbrain namespace and every other Secondbrain path as read-only. This bounds where a session may write, never whether it must.
<!-- /source-baseline-entry: clause.secondbrain.inbox-write-permission -->

<!-- source-baseline-entry: clause.secondbrain.inbox-entry-contract -->
An Inbox Entry uses positional status: pending entries live at the destination namespace root under `inbox/<destination>/`, and resolved entries live under `inbox/<destination>/_triaged/`. Use this complete copyable contract:

```yaml
---
origin: <project-that-observed>
destination: <project-that-triages>
type-hint: <finding-or-intent-hint>
created_at: YYYY-MM-DD
capture: manual # manual | auto
# added at triage time; exactly one:
resolved_to: <repository-relative-artifact-path>
# or discarded_reason: <reason>
---
```
<!-- /source-baseline-entry: clause.secondbrain.inbox-entry-contract -->

<!-- source-baseline-entry: clause.secondbrain.inbox-triage-order -->
Triage works pending entries `oldest first`, ordered by `created_at`; take the earliest entry before newer arrivals.
<!-- /source-baseline-entry: clause.secondbrain.inbox-triage-order -->

<!-- source-baseline-entry: clause.secondbrain.inbox-empty-state -->
Treat an `empty inbox` as rest and continue the session's current work; do not invent a missing Triage step.
<!-- /source-baseline-entry: clause.secondbrain.inbox-empty-state -->

<!-- source-baseline-entry: clause.secondbrain.inbox-auto-capture -->
The maintainer owns the session-end hook outside this repository; this clause contracts only what that hook writes. Every `capture: auto` draft is always pending triage and never self-triaged.
<!-- /source-baseline-entry: clause.secondbrain.inbox-auto-capture -->

<!-- source-baseline-entry: clause.secondbrain.research-capture -->
A session that performed substantive external research must capture a digest with its sources for the brain's own namespace. Run the advisory qmd duplicate check first through an authorized access path, verify that returned paths exist, and review substantive overlap; a score alone never decides. A strong verified match routes the digest to extend existing knowledge instead of duplicating it; otherwise create a new pending research Inbox Entry. Ingestion remains the brain's own contract.
<!-- /source-baseline-entry: clause.secondbrain.research-capture -->
<!-- source-baseline-entry: clause.secondbrain.capture-trigger -->
- MUST capture an observed defect, improvement, or feature idea as one pending Inbox Entry under the namespace of the project that owns the fix, which is frequently not the project the session runs in, after reading that destination's existing entries and extending a strong verified match instead of duplicating it.
<!-- /source-baseline-entry: clause.secondbrain.capture-trigger -->
<!-- source-baseline-entry: clause.secondbrain.capture-self-contained -->
- MUST write each Inbox Entry so a triaging session can act without the author's context, commit it at the moment of capture, and MUST NOT commit, edit, or move an entry another session created.
<!-- /source-baseline-entry: clause.secondbrain.capture-self-contained -->

<!-- source-baseline-entry: clause.secondbrain.prohibit-writes -->
Do not create, edit, rename, move, or delete any Secondbrain file outside `inbox/**`. Do not edit `raw/` or `projects/*/mirror/`, and never copy code or generated artifacts from a mirror without a local source check.
<!-- /source-baseline-entry: clause.secondbrain.prohibit-writes -->

<!-- source-baseline-entry: clause.secondbrain.prohibit-secret-access -->
Never read, copy, or expose `.env` files, tokens, credentials, cookies, private keys, API keys, session material, or unsafe personal and client data. Stop at likely secret-bearing sources and request a safe source.
<!-- /source-baseline-entry: clause.secondbrain.prohibit-secret-access -->

<!-- source-baseline-entry: clause.secondbrain.cite-used-files -->
Cite every Secondbrain file used in the final response or handoff by path. Do not claim Secondbrain context when no Secondbrain file was read.
<!-- /source-baseline-entry: clause.secondbrain.cite-used-files -->

<!-- source-baseline-entry: clause.secondbrain.escalate-durable-updates -->
When Secondbrain knowledge must be added or corrected, ask Hermes to ingest or update it instead of writing from this repository.
<!-- /source-baseline-entry: clause.secondbrain.escalate-durable-updates -->
<!-- source-baseline-entry: clause.secondbrain.baseline-owned-guidance -->
- MUST treat guidance inside `setup-context-driven` markers as Baseline-owned: propose changes to it as an Inbox Entry to the Baseline's owner, never as a local edit the next Baseline update overwrites.
<!-- /source-baseline-entry: clause.secondbrain.baseline-owned-guidance -->
