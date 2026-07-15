# Secondbrain

Use the local Secondbrain only when this repository does not already answer the question through `CONTEXT.md`, ADRs, local documentation, or source code.

## When to consult it

Consult the Secondbrain before acting when work depends on:

- business context, prior decisions, or cross-project documentation;
- fiscal or tax concepts;
- knowledge about Vortex, Tax, Visio, or Gesttione;
- shared architecture patterns that might have been decided outside this repository.

Do not consult it when local repository context fully answers the task.

## Query order

1. Read `~/dev/secondbrain/wiki/index.md` first to identify the relevant area.
2. Run `qmd query "<question>" --all --files --min-score 0.3` to find candidate files.
3. For mirrored projects, inspect `~/dev/secondbrain/projects/<project>/mirror/` only after the index and query point there.
4. Open only the files needed for the task.
5. Cite every Secondbrain file used in the final response or handoff.

## Read-only boundaries

- You may read `wiki/`, `shared/`, `raw/`, and `projects/*/mirror/` when the task requires it.
- Do not write to the Secondbrain from this repository.
- Do not edit raw/ files.
- Do not edit projects/*/mirror/ files.
- Do not create, rename, move, or delete files inside the Secondbrain.
- If knowledge must be added or corrected, ask Hermes to ingest or update it.

## Secret safety

- Never read, copy, or expose `.env` files, tokens, credentials, cookies, private keys, API keys, or session material.
- If search results point to a likely secret-bearing file, stop reading that file and report that the task needs a safe source.
- Do not paste Secondbrain excerpts that contain personal, client, fiscal, or credential-like data unless the user explicitly asks and the content is safe to share.

## Project mirrors

Project mirrors are references, not workspaces. Use a mirror to understand another project's vocabulary, decisions, or architecture. Do not treat it as the source of truth for this repository, and do not copy code or generated artifacts from a mirror into this repository without a separate local source check.

## Citations

Every answer that uses Secondbrain context must cite each Secondbrain file read. Cite paths, not vague source names. If no Secondbrain file was needed, do not mention it.
