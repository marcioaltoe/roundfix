<!-- setup-context-driven:begin id=guide.secondbrain version=0.0.1 -->

# Secondbrain

- **mandatory**: Consult both the local Secondbrain and Exa MCP before proposing solutions or recommendations, making product, design, or architecture decisions, and authoring an Idea, PRD, or TechSpec. Both sources are required; one does not replace the other.

- **mandatory**: Read `wiki/index.md` first. Then run `qmd query "<question>" --all --files --min-score 0.3`. Inspect `projects/<project>/mirror/` only when the index and query point there, and open only the files required for the task; treat mirrors as references, not workspaces.

- **mandatory**: Use Secondbrain for prior decisions, related project experience, and existing research.

- **mandatory**: Use Exa MCP to find and read relevant external sources that support or challenge the proposal. Prefer primary sources, assess their applicability, and distinguish published evidence from inference. This requirement complements the mandatory authoritative documentation workflow for external APIs and libraries.

- **mandatory**: If either source is unavailable or yields no relevant result, record the attempted consultation and its limitation explicitly, then continue with the available evidence. Do not claim validation from an unavailable, unread, or irrelevant source. Another research tool does not satisfy the Exa MCP requirement.

- **prohibited**: Do not create, edit, rename, move, or delete any Secondbrain file, `inbox/` included. Do not edit `raw/` or `projects/*/mirror/`, and never copy code or generated artifacts from a mirror without a local source check.

- **mandatory**: Send substantive external research and relevant findings or articles from Exa MCP as a sourced digest in one message to the project that owns the Secondbrain, for its own ingestion and later reuse as reference material, including relevant articles found outside a dedicated research session. Include each source's title and URL, a concise summary of the content read, the projects or topics it informs, and why it matters; distinguish published evidence from inference. Run the advisory qmd duplicate check first through an authorized access path, verify that returned paths exist, and name a strong verified match in the message so the owner extends existing knowledge instead of duplicating it; a score alone never decides.

- **prohibited**: Never read, copy, or expose `.env` files, tokens, credentials, cookies, private keys, API keys, session material, or unsafe personal and client data. Stop at likely secret-bearing sources and request a safe source.

- **mandatory**: Cite every Secondbrain file used in the final response or handoff by path. Do not claim Secondbrain context when no Secondbrain file was read.

- **mandatory**: Record the Secondbrain file paths and external source URLs used, and explain how they informed, changed, or challenged the decision in the relevant technical artifact. If no artifact is being authored, include this record in the response.

- **mandatory**: Guidance delivered inside `setup-context-driven` markers is owned by the Baseline, not by the repository holding it. Proposing a change to it is a message to the Baseline's owner; editing it locally produces a change the next Baseline update silently overwrites. Verify ownership by looking for the markers before editing any agent guide, because the same file usually carries repository-authored prose outside them that is yours to change.

- **mandatory**: When Secondbrain knowledge must be added or corrected, ask Hermes to ingest or update it instead of writing from this repository.

<!-- setup-context-driven:end id=guide.secondbrain -->
