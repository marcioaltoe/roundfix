<!-- setup-context-driven:begin id=guide.secondbrain version=0.0.1 -->

# Secondbrain

- **mandatory**: Consult the local Secondbrain before acting when repository context does not answer business or prior-decision questions, fiscal or tax concepts, cross-project documentation, knowledge about Vortex, Tax, Visio, or Gesttione, or shared architecture patterns. Do not consult it when local code, `CONTEXT.md`, ADRs, and repository documentation fully answer the task.

- **mandatory**: Read `wiki/index.md` first. Then run `qmd query "<question>" --all --files --min-score 0.3`. Inspect `projects/<project>/mirror/` only when the index and query point there, and open only the files required for the task; treat mirrors as references, not workspaces.

- **mandatory**: An Inbox Entry uses positional status: pending entries live at the destination namespace root under `inbox/<destination>/`, and resolved entries live under `inbox/<destination>/_triaged/`. Use this complete copyable contract:

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

- **mandatory**: Sessions MAY create files under the Secondbrain's `inbox/**`; this is the only writable Secondbrain namespace. Every other Secondbrain path stays read-only.

- **prohibited**: Do not create, edit, rename, move, or delete any Secondbrain file outside `inbox/**`. Do not edit `raw/` or `projects/*/mirror/`, and never copy code or generated artifacts from a mirror without a local source check.

- **prohibited**: Never read, copy, or expose `.env` files, tokens, credentials, cookies, private keys, API keys, session material, or unsafe personal and client data. Stop at likely secret-bearing sources and request a safe source.

- **mandatory**: Cite every Secondbrain file used in the final response or handoff by path. Do not claim Secondbrain context when no Secondbrain file was read.

- **mandatory**: When Secondbrain knowledge must be added or corrected, ask Hermes to ingest or update it instead of writing from this repository.

<!-- setup-context-driven:end id=guide.secondbrain -->
