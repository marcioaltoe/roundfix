---
status: done
created_at: 2026-10-01
updated_at: 2026-10-01
absorbed_by: 0208-a-baseline-that-follows-the-reshaped-skills-catalog
---

# Baseline: the skills catalog reshaped the Baseline's inputs (marcioaltoe/skills#103)

The maintainer changed the skills repository on 2026-10-01 (PR #103) and asked it to report what the Baseline must adjust before the next Roundfix release. Six Secondbrain inbox entries carry those reports, and this Finding consolidates them.

## 1. The setup presets were renamed, and `go-cli` was retired

- Symptom / evidence:
  - `setups/go-tui.txt` is now `go`, `rust-cli.txt` is now `rust`, and `typescript-bun.txt` is now `typescript`, each with the same skill list.
  - `go-cli` was merged into `go`.
  - The `docs`, `marketing`, `release`, `context-workflow` and `skill-authoring` presets are gone, and `content` is new.
  - The next asset sync fails, because `internal/baseline/assets/setups/{go-tui,rust-cli,typescript-bun}.json` and `go-cli.json` name presets that no longer exist on `main`.
- Action: point the profiles `go-cli-tui`, `rust-cli` and `standard-typescript-monorepo` at the setups `go`, `rust` and `typescript`. Replace the snapshots, retire `go-cli.json`, and move every test, parity fixture and Source Baseline reference that names an old id. Align the active Specs 0206 and 0207, which name the old setups.

## 2. `review` and `triage` were removed while the Baseline still uses them

- Symptom / evidence:
  - `review` (a frozen copy of the two-axis review skill) is listed in every setup snapshot and dispatched by `core` as `trigger.core.review`. Nothing loads it, and `roundfix review` is unrelated.
  - `triage` is required and dispatched by the `typescript` module. It depends on a removed skill, routes work outside the Spec workflow, and sets `disable-model-invocation: true`, so the dispatch never loads it.
- Action:
  - Drop `review` and its dispatch, without renaming it to `code-review`.
  - Drop `triage` from `typescript` and give the `external-triage` module the functions external triage still needs:
    - classify each item and move it through states mapped by the `triage.external` decision;
    - for `needs-info`, ask specific questions and stop;
    - start every forge comment with an AI-generated disclaimer;
    - record `wontfix` decisions so a repeated request gets the same answer;
    - treat a pull request as an issue with code attached;
    - route accepted work into the Spec workflow.
  - Remove this repository's own installed `review` and `triage` copies. Adopters keep their installed trees until they remove them (ADR-0191).

## 3. Skills the catalog removed or now requires

- Symptom / evidence:
  - The 13 removed skills still listed in Baseline snapshots are `accessibility`, `best-practices`, `performance`, `web-quality-audit`, `firecrawl-developer-index`, `firecrawl-monitor`, `firecrawl-research-index`, and `golang-{benchmark,continuous-integration,database,modernize,observability,performance}`. No module requires them.
  - `resolving-merge-conflicts` and `lesson-learned` were removed and are named by no Baseline setup.
  - `typesafe-ai` is in every preset, but no module requires it, so managed repositories never install it.
  - `crafting-effective-readmes` is now in the Go and Rust presets, but only the `typescript` module requires it.
- Action:
  - Let the snapshot refresh drop the removed rows.
  - Require `typesafe-ai` in `core`, with a dispatch for programmable semantic judgment or TypeSafe/Jev work.
  - Require `crafting-effective-readmes` for Go and Rust, or move it to `core`.

Sources: Secondbrain `inbox/roundfix/2026-10-01-skills-*.md` (six entries); https://github.com/marcioaltoe/skills/pull/103.

## Addendum — 2026-10-01 — routed to Spec 0208

Spec `0208-a-baseline-that-follows-the-reshaped-skills-catalog` adopts all three sections. The snapshots take the new setup
names and `go-cli` retires; `review` and `triage` leave the catalog and the
`external-triage` module gains the triage clauses; `core` requires
`typesafe-ai` and `crafting-effective-readmes`. The removed rows of section 3
drop through the snapshot refresh. This repository removes its own `review`
lock entry with `roundfix baseline skills reconcile` and deletes the tree.
Status is `done` because the Spec is created and linked, not because it is
delivered.
