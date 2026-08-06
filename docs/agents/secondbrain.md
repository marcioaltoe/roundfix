# Secondbrain

Consult the local Secondbrain when work requires business context, prior
decisions, cross-project documentation, fiscal or tax concepts, knowledge of
Vortex, Tax, Visio, or Gesttione, or shared architecture patterns.

## Location

The Secondbrain is available at:

`~/dev/secondbrain`

## When to use it

Consult the Secondbrain before acting when a task depends on:

- business context or prior decisions not documented in this repository;
- documentation or patterns shared across projects;
- fiscal or tax concepts;
- knowledge about Vortex, Tax, Visio, or Gesttione;
- architecture decisions that might have been made in another project.

Do not consult the Secondbrain when the code, `CONTEXT.md`, ADRs, and local
documentation fully answer the task.

## How to use it

1. Read `~/dev/secondbrain/wiki/index.md`.
2. Run `qmd query "<question>" --all --files --min-score 0.3`.
3. For mirrored projects, inspect
   `~/dev/secondbrain/projects/<project>/mirror/`.
4. Open only the files required for the task.
5. Cite every Secondbrain file used in the response.

## Access rules

- You may read `wiki/`, `shared/`, `projects/*/mirror/`, and `raw/`.
- Sessions may create files under `inbox/**`.
- Every other Secondbrain path is read-only.
- Do not edit `raw/`.
- Do not edit `projects/*/mirror/`.
- If something must become durable knowledge, ask Hermes to ingest or update
  it in the Secondbrain.
- Never read, copy, or expose `.env` files, tokens, credentials, cookies, or
  keys.

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

<!-- roundfix:repository-rule:begin id=rule.8d99a6ddc2032d2e4df584b4e9402d7e733345ebb279cfe371423d379054c639 -->
### Secondbrain

When work depends on business context, prior decisions, cross-project
documentation, fiscal or tax concepts, Vortex, Tax, Visio, Gesttione, or shared
architecture patterns, read `docs/agents/secondbrain.md` before acting. Skip it
for self-contained repository work that the local code and docs fully answer.

<!-- roundfix:repository-rule:end id=rule.8d99a6ddc2032d2e4df584b4e9402d7e733345ebb279cfe371423d379054c639 -->

<!-- roundfix:repository-rule:begin id=rule.cf2bb6e8555b28a55a64f750621e98ecc8adec17060fed22824dacd2bb227bd9 -->
The Secondbrain's `inbox/**` is the only writable namespace from this repo;
every other path is read-only, and responses must cite every Secondbrain file
used.


<!-- roundfix:repository-rule:end id=rule.cf2bb6e8555b28a55a64f750621e98ecc8adec17060fed22824dacd2bb227bd9 -->
