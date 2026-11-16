# Where the adopters on the authoring host keep their own rules

The author of Spec 0399 recorded this inventory on 2026-11-16, about 19:30 BRT, on the MacBook. It is read-only. Nothing was committed in any adopter, no message was sent to another project, and each adopter's `git status --porcelain` was the same before and after.

## Method

1. The author searched the host for adopted repositories. These are the repositories that hold `docs/agents/setup-context.json` within five levels of the home folder, excluding `Library`, `node_modules` and Roundfix's own stores and worktrees.
2. For each one it read the Setup Manifest's `managedArtifacts` and took the Markdown carriers it lists, plus a root `CLAUDE.md` when that is a regular file and not a symbolic link. In every repository here, `CLAUDE.md` is a symbolic link to `AGENTS.md`.
3. In each carrier it listed every passage outside the `setup-context-driven` marker pairs. A passage is a maximal run of lines between two Managed Regions, or between a Managed Region and the start or end of the file. Leading and trailing blank lines are trimmed.
4. A block that another tool writes and delimits with its own `<!-- BEGIN:<name> -->` and `<!-- END:<name> -->` pair was set apart and is not counted as a passage.
5. It ran the installed Roundfix 0.147.0 plan without confirmation, which writes nothing: `roundfix baseline update --repo . --no-skills --no-history --format text`. This shows what the tool reports today.

## Repositories

| Checkout | Commit | Roundfix that wrote the Baseline | Baseline Profile |
| --- | --- | --- | --- |
| `~/dev/argus` | `fc53eb7c` | 0.143.0 | `argus-electron-app` |
| `~/dev/quaestor` | `fdd3d184` | 0.145.1 | `python` |
| `~/srv/quaestor` | `fdd3d184` | 0.145.1 | `python` (deployment checkout of the same repository and commit) |
| `~/srv/fiscus` | `55bac1b1` | 0.142.0 | `standard-typescript-monorepo` |
| `~/srv/fiscus-mcp` | `55bac1b1` | 0.142.0 | `standard-typescript-monorepo` (second checkout of the same repository and commit) |
| `~/srv/indexus` | `3e48f0d0` | 0.145.1 | `indexus-backend` |
| `~/dev/roundfix` | `bcf0e675` | 0.145.1 | `go-cli-tui` (this repository) |

These are five distinct repositories. Vortex, Fluxus, Oraculum, Conexus and Pantheon have no checkout on this host, so they are not in this inventory.

## Passages outside every Managed Region

| Repository | Passage | First line |
| --- | --- | --- |
| quaestor | `docs/agents/skill-dispatch.md:109-115` | `## Repository-specific rules` |
| indexus | `docs/agents/skill-dispatch.md:183-189` | `## Repository-specific rules` |
| fiscus | `AGENTS.md:1-3` | `### Repository review policy` |
| fiscus | `AGENTS.md:89-91` | `### Product tax review (client registers)` |
| fiscus | `docs/agents/agent-instructions.md:1-5` | `## Maintainer review override` |
| fiscus | `docs/agents/autonomous-work.md:1-6` | `## Repository delivery override` |
| fiscus | `docs/agents/backend.md:37-53` | `<!-- repository-authored: fora dos marcadores de setup, preservado por docs/agents/docs-layout.md -->` |
| fiscus | `docs/agents/secondbrain.md:31-62` | `<!-- repository-authored: maintainer instruction, 2026-09-08 -->` |
| argus | none | |
| roundfix | none | |

The shapes the Baseline has to recognise are these:

- **A section appended after the guide's Managed Region.** This is the quaestor and indexus `skill-dispatch.md` case. Both added the same three rules on the authoring models and the research subagent, under a heading that says they are "authored by this repository, outside the Baseline markers".
- **A section before the first Managed Region.** These are the fiscus overrides at the top of `agent-instructions.md`, `autonomous-work.md` and the root `AGENTS.md`.
- **A section between two Managed Regions of the root instruction file.** This is fiscus `AGENTS.md:89-91`. It is followed by a Turborepo block (`<!-- BEGIN:turborepo-agent-rules -->` to `<!-- END:turborepo-agent-rules -->`, lines 93-102), which the `turbo` tool writes and re-adds by itself. That block is not a repository rule and must not be reported.
- **Several sections after the Managed Region, each with an HTML comment naming its origin.** This is fiscus `secondbrain.md:31-62`. It is one passage, because no Managed Region separates the sections.
- **A passage that contradicts the Managed Region above it.** This is fiscus `backend.md:37-53`. It states that it contradicts the setup-owned text, and it cites `docs/agents/docs-layout.md` as the reason it may stay outside the markers.

## Inside the Managed Regions

The 0.147.0 plan reported no Unrecorded Managed Region in any of the five repositories. No adopter on this host holds a hand edit inside a marker pair today.

## What Roundfix 0.147.0 reports for these passages

Nothing. The plans of quaestor, fiscus and indexus listed only the managed files they would refresh, for example `update docs/agents/autonomous-work.md (1 managed entries)`. They printed no line about a passage outside a Managed Region. A managed refresh preserves those bytes and says nothing about them.

## What the Baseline text said when these passages were written

The adopters followed the text they received. The guide template of `docs/agents/agent-instructions.md` says: "Repository authors own project-specific extensions outside setup markers and may add stricter rules". The guide template of `docs/agents/skill-dispatch.md` says: "repository-authored dispatch extensions may add stricter triggers". The clause on documentation directories says: "Preserve repository-authored extensions outside setup markers". The glossary entry **Repository-Specific Normative Rules** in this repository's `CONTEXT.md` places them only in `docs/agents/specific-repository.md`, but no text an adopter receives says so.

## Use in Spec 0399

The Tasks build their test fixtures from these shapes, written in neutral text. Each passage above becomes one expected entry, and the Turborepo block, the clean argus and roundfix carriers, and the `CLAUDE.md` symbolic link become cases that must report nothing. The QA gate re-reads the five repositories read-only with the built binary. It accepts the table above as the expected answer for each adopter whose commit is unchanged, and for a changed commit it compares the report with a fresh read made by the same method.
