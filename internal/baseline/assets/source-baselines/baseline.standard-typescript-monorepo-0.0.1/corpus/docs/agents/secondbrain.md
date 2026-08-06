<!-- source-baseline-entry: contract.secondbrain.protocol -->
# Secondbrain protocol

1. Consult the local knowledge system only when repository code, domain context, ADRs, and local documentation do not answer a business, prior-decision, cross-repository, or shared-architecture question.
2. Read its index first, then use its declared semantic query command. Open only the files required by the result and treat project mirrors as references, never as workspaces.
3. Sessions MAY create files under the knowledge system's `inbox/**`; every other knowledge-system path remains read-only. MUST NOT create, edit, rename, move, or delete raw ingestion data or project mirrors from this repository.
4. Pending Inbox Entries live at their destination namespace root and resolved entries live under `_triaged/`; each entry declares `origin:`, `destination:`, `type-hint:`, `created_at:`, and `capture:`, then gains exactly one of `resolved_to:` or `discarded_reason:` at triage time.
5. MUST NOT read, copy, or expose secrets, session material, private keys, credentials, or unsafe personal and client data. Stop at likely secret-bearing sources and request a safe source.
6. Cite every knowledge-system file used in the final response or handoff. Do not claim knowledge-system context when no file was read.
7. Request an authorized ingestion workflow for durable additions or corrections instead of writing from this repository.
<!-- /source-baseline-entry: contract.secondbrain.protocol -->

<!-- source-baseline-entry: clause.secondbrain.01-consult-triggers -->
Consult the local Secondbrain before acting when repository context does not answer business or prior-decision questions, fiscal or tax concepts, cross-project documentation, knowledge about Vortex, Tax, Visio, or Gesttione, or shared architecture patterns. Do not consult it when local code, `CONTEXT.md`, ADRs, and repository documentation fully answer the task.
<!-- /source-baseline-entry: clause.secondbrain.01-consult-triggers -->

<!-- source-baseline-entry: clause.secondbrain.02-query-order -->
Read `wiki/index.md` first. Then run `qmd query "<question>" --all --files --min-score 0.3`. Inspect `projects/<project>/mirror/` only when the index and query point there, and open only the files required for the task; treat mirrors as references, not workspaces.
<!-- /source-baseline-entry: clause.secondbrain.02-query-order -->

<!-- source-baseline-entry: clause.secondbrain.inbox-write-permission -->
Sessions MAY create files under the Secondbrain's `inbox/**`; this is the only writable Secondbrain namespace. Every other Secondbrain path stays read-only.
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
