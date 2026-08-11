<!-- setup-context-driven:begin id=guide.secondbrain version=0.0.1 -->

# Secondbrain

- **mandatory**: Consult the local Secondbrain before making a design or architecture proposal, and before authoring an Idea, PRD, or TechSpec.

- **mandatory**: Read `wiki/index.md` first. Then run `qmd query "<question>" --all --files --min-score 0.3`. Inspect `projects/<project>/mirror/` only when the index and query point there, and open only the files required for the task; treat mirrors as references, not workspaces.

- **mandatory**: Consult the Secondbrain while a decision is being formed — authoring a Spec, choosing an approach, or validating a strategy — for what it holds and this repository does not: the decisions sibling projects in the ecosystem already made and paid for, literature, and general technical knowledge. Report an unreachable Secondbrain or an empty result as a condition in the artifact being written; it is never a reason to stop the work.

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

- **mandatory**: A session that performed substantive external research must capture a digest with its sources for the brain's own namespace. Run the advisory qmd duplicate check first through an authorized access path, verify that returned paths exist, and review substantive overlap; a score alone never decides. A strong verified match routes the digest to extend existing knowledge instead of duplicating it; otherwise create a new pending research Inbox Entry. Ingestion remains the brain's own contract.

- **prohibited**: Never read, copy, or expose `.env` files, tokens, credentials, cookies, private keys, API keys, session material, or unsafe personal and client data. Stop at likely secret-bearing sources and request a safe source.

- **mandatory**: Cite every Secondbrain file used in the final response or handoff by path. Do not claim Secondbrain context when no Secondbrain file was read.

- **mandatory**: Guidance delivered inside `setup-context-driven` markers is owned by the Baseline, not by the repository holding it. Proposing a change to it is an Inbox Entry addressed to the Baseline's owner; editing it locally produces a change the next Baseline update silently overwrites. Verify ownership by looking for the markers before editing any agent guide, because the same file usually carries repository-authored prose outside them that is yours to change.

- **mandatory**: When Secondbrain knowledge must be added or corrected, ask Hermes to ingest or update it instead of writing from this repository.

<!-- setup-context-driven:end id=guide.secondbrain -->
