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
- Do not edit `raw/`.
- Do not edit `projects/*/mirror/`.
- Do not write to the Secondbrain from this repository.
- If something must become durable knowledge, ask Hermes to ingest or update
  it in the Secondbrain.
- Never read, copy, or expose `.env` files, tokens, credentials, cookies, or
  keys.
