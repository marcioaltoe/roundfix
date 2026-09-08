<!-- setup-context-driven:begin id=guide.secondbrain version=0.0.1 -->

# Secondbrain

- **mandatory**: Consult both the local Secondbrain and Exa MCP before proposing solutions or recommendations, making product, design, or architecture decisions, and authoring an Idea, PRD, or TechSpec. Both sources are required; one does not replace the other.

- **mandatory**: Read `wiki/index.md` first. Then run `qmd query "<question>" --all --files --min-score 0.3`. Inspect `projects/<project>/mirror/` only when the index and query point there, and open only the files required for the task; treat mirrors as references, not workspaces.

- **mandatory**: Use Secondbrain for prior decisions, related project experience, and existing research. Follow the index-first and local query workflow defined in this guide.

- **mandatory**: Use Exa MCP to find and read relevant external sources that support or challenge the proposal. Prefer primary sources, assess their applicability, and distinguish published evidence from inference. This requirement complements the mandatory authoritative documentation workflow for external APIs and libraries.

- **mandatory**: If either source is unavailable or yields no relevant result, record the attempted consultation and its limitation explicitly, then continue with the available evidence. Do not claim validation from an unavailable, unread, or irrelevant source. Another research tool does not satisfy the Exa MCP requirement.

- **mandatory**: Write each Inbox Entry so the triaging session can act on it without the author's context: the observation, its evidence with commands and paths, and the reasoning that makes it actionable. Commit the entry at the moment of capture, because durability is the point of the door. Never commit, edit, or move an entry another session created, even in the same namespace on the same day; two projects reporting one class of defect is signal for Triage, not a merge conflict to resolve.

- **mandatory**: Observing a defect, an improvement worth making, or a feature idea obliges capture: create one pending Inbox Entry under `inbox/<destination>/` for the project that owns the fix, which is frequently not the project the session is running in. Read the destination's existing pending and triaged entries first and extend a strong verified match instead of duplicating it. Capture is an obligation, not a permission; an observation left only in a session transcript is lost when that session ends.

- **mandatory**: The maintainer owns the session-end hook outside this repository; this clause contracts only what that hook writes. Every `capture: auto` draft is always pending triage and never self-triaged.

- **mandatory**: Treat an `empty inbox` as rest and continue the session's current work; do not invent a missing Triage step.

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

- **mandatory**: Triage works pending entries `oldest first`, ordered by `created_at`; take the earliest entry before newer arrivals.

- **mandatory**: `inbox/**` is the only writable Secondbrain namespace; every other Secondbrain path stays read-only. This clause bounds where a session may write, never whether it must — the capture obligation is stated separately.

- **prohibited**: Do not create, edit, rename, move, or delete any Secondbrain file outside `inbox/**`. Do not edit `raw/` or `projects/*/mirror/`, and never copy code or generated artifacts from a mirror without a local source check.

- **mandatory**: Capture substantive external research and relevant findings or articles from Exa MCP as sourced digests in pending research Inbox Entries in the brain's own namespace, for later ingestion and reuse as reference material. Capture relevant articles on subjects important or foundational to the projects even when found outside a dedicated research session. Include each source's title and URL, a concise summary of the content read, the projects or topics it informs, and why it matters; distinguish published evidence from inference. Run the advisory qmd duplicate check first through an authorized access path, verify that returned paths exist, and review substantive overlap; a score alone never decides. A strong verified match routes the digest to extend existing knowledge instead of duplicating it; otherwise create a new pending research Inbox Entry. Follow the Inbox Entry contract. Capture is pending ingestion, which remains the brain's own contract.

- **prohibited**: Do not use Exa or another external research tool to discover or infer local repository code or behavior. Use local code-search tools for that purpose. Never include credentials, private client records, or proprietary source code in external search queries.

- **prohibited**: Never read, copy, or expose `.env` files, tokens, credentials, cookies, private keys, API keys, session material, or unsafe personal and client data. Stop at likely secret-bearing sources and request a safe source.

- **mandatory**: Cite every Secondbrain file used in the final response or handoff by path. Do not claim Secondbrain context when no Secondbrain file was read.

- **mandatory**: Record the Secondbrain file paths and external source URLs used, and explain how they informed, changed, or challenged the decision in the relevant technical artifact. If no artifact is being authored, include this record in the response.

- **mandatory**: Guidance delivered inside `setup-context-driven` markers is owned by the Baseline, not by the repository holding it. Proposing a change to it is an Inbox Entry addressed to the Baseline's owner; editing it locally produces a change the next Baseline update silently overwrites. Verify ownership by looking for the markers before editing any agent guide, because the same file usually carries repository-authored prose outside them that is yours to change.

- **mandatory**: When Secondbrain knowledge must be added or corrected, ask Hermes to ingest or update it instead of writing from this repository.

<!-- setup-context-driven:end id=guide.secondbrain -->
