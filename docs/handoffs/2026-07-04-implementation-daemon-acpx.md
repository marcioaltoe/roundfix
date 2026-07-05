# Handoff — roundfix as the generic implementation daemon (ACP → acpx)

- **Date**: 2026-07-04
- **From**: workflow design session in `~/dev/skills` (Claude Fable 5)
- **Repo for this work**: `~/dev/roundfix`, branch `ma/redesign-and-new-functions`

## Mission

Evolve roundfix from a CodeRabbit review-resolution CLI into the generic implementation daemon of the CONTEXT-driven workflow: pick up per-feature task graphs produced by the skills pipeline, execute them via local coding agents over ACP, and modernize the ACP layer to acpx.

## Upstream context (consume, do not rebuild)

The planning side already shipped in marcioaltoe/skills PR #28 (**merged**, <https://github.com/marcioaltoe/skills/pull/28>): `write-prd` → `write-techspec` → `write-tasks` → `implement-task` / `implement-spec` → `qa-gate` → `archive-spec`, plus `setup-workflow` v0.2.0. Read those SKILL.md files for the exact contracts — `implement-task` is the agent-side contract this daemon must mirror. PR #29 (**merged**, <https://github.com/marcioaltoe/skills/pull/29>) then rebuilt the `go-cli` setup around **this repo's profile** (stdlib `flag` dispatch — no Cobra; Bubble Tea v2 on `charm.land`; zero-test-dep stdlib testing; golang error-handling/concurrency/context/lint skills vendored from samber).

The artifact contract the daemon consumes, per feature at `docs/specs/<slug>/` in **target repos**:

- `_prd.md` — frontmatter `spec`, `status: active|archived`, `surfaces: [frontend|backend|cli|data|infra|docs]` (routes QA).
- `_tasks.md` — DAG manifest: frontmatter `schema: spec-tasks/v1`, `graph.nodes[] {id, file, needs[]}`. Dependencies live **only** here.
- `task_NN.md` — frontmatter `task, spec, status: pending|in_progress|completed|failed, type, complexity`. Status lives **only** here. Body: Requirements / Subtasks / Acceptance Criteria / Verification (commands the executor runs verbatim) / References; the executor appends `## Result` with evidence.
- `qa/` — qa-gate evidence reports. Shipped specs move to `docs/specs/_archived/<slug>/`.

Install the setup with: `curl -fsSL https://raw.githubusercontent.com/marcioaltoe/skills/main/install.sh | bash -s -- go-cli` (PRs #28 and #29 are both merged to main).

## Current repo state (at handoff time)

- Branch `ma/redesign-and-new-functions` has **uncommitted deletions** pruning `.agents/skills/` — an installed-skills cleanup in progress. Resolve that first (commit the prune or restore); `skills-lock.json` (~50 entries) still lists the old set.
- Read before anything: `CONTEXT.md` (the glossary — vocabulary contract for code, docs, and prompts), `AGENTS.md`, `docs/adr/` (11 ADRs; continue the numbering), and skim `docs/product-brief.md` (the MVP contract — supersede its decisions via new ADRs, don't edit history).

## Work plan (from the July 2026 architecture analysis)

Priority order. Each item is roughly one spec — dogfood the workflow: run `/setup-workflow` here, then grill → PRD → TechSpec → tasks for each phase.

1. **Task-source abstraction.** The seams exist with exactly one implementation each: `internal/watch` (`Fetcher`/`Resolver`/`StatusSource`) and `internal/reviewsource` (`Source`, CodeRabbit-only). Add a task-file source reading a target repo's `_tasks.md` + task files. Purge PR-centric identity: `PullRequestRef`, the one-Active-Run-per-head-branch lock, the `@coderabbitai review` string in `internal/watch/watch.go`, the `fix: resolve Roundfix batch` commit template, and `gh`-based preflight assumptions.
2. **ACP → acpx.** Today: `internal/agent/acp_runner.go` (~1,150 lines) on `coder/acp-go-sdk` v0.13.5, spawning `codex-acp` / `claude-agent-acp` / `opencode acp` per batch. Evaluate acpx (<https://github.com/openclaw/acpx>, npm `acpx`): persistent named sessions scoped by (agent, cwd), queue-owner processes, `--no-wait`, cooperative cancel, session resume/respawn, `acpx flow run`. Preserve what is good in the current layer: raw-payload journaling of `session/update` (ADR 0008), the daemon/agent ownership split (ADR 0001/0011 — daemon owns verify/commit/push; agent owns edits and status only), path-jailed fs handlers, robust teardown. Decide keep-SDK vs shell-out-to-acpx vs hybrid; record as an ADR.
   Status note (2026-07-05): done by Spec `0002-acpx-migration`; see ADR-0017 and the Spec's Task evidence.
3. **Real daemon + queue.** `roundfix serve`: long-lived process, work-queue table in the existing SQLite store (its WAL single-writer + append-only journal + `attach` replay, ADR 0004/0009, are the right primitives), filesystem/webhook triggers instead of PR polling, `attach` as the primary UX, DB-mediated control channel (pause/resume/reprioritize — `stop` already works this way).
4. **DAG scheduler + worktree-per-task.** Replace `PlanBatches` (dedup + fixed-size chunks) with ready-set computation over `needs`; run parallel agent sessions in isolated git worktrees (`resolve.concurrent` exists in config, unused). This removes the dirty-worktree rejection and the fragile before/after snapshot diffing.
5. **Templated prompt contract.** `BuildPrompt` is Go string-building with repo-specific noise baked in (the `rtk bun run --cwd` lesson ships to every repo). Move to versioned `text/template` prompts with per-repo overrides; align invariants with the `implement-task` skill so agent-side and daemon-side contracts cannot drift.
6. **Structured result artifacts.** Frontmatter `status` plus prose is the only machine-readable agent output today (`SettleAssignedIssues` compensates for forgetful agents). Define a per-task result artifact (files changed, commands run, evidence per acceptance criterion, failure reason) the daemon parses — mirroring the task file's `## Result` section.
7. **Retry/escalation policy.** Per-task retry budgets, agent escalation (codex → claude → needs-human terminal state) surfaced in the journal with non-zero exit. The existing `Unresolved` vs `Failed` distinction is the right foundation.
8. **Permission policy engine.** Stop auto-approving the first `RequestPermission` option and blanket `danger-full-access`/`bypassPermissions`; per-task-class allowlists of tools/paths/commands, deny-with-journal.
9. **Slim `internal/cli`** (2,071 lines mixing flags, TUI, assembly, reporting) into a composition-root package so the daemon binary, CLI, and tests share one wiring path. The liftable core is `daemon`, `watch`, `runevent`, `store`, `agent`.

Small cleanups also pending: empty scaffolds `skills/roundfix-resolve-round/` and `skills/roundfix-watch/`; the shipped `skills/roundfix/SKILL.md` duplicates BuildPrompt's contract and drifts independently.

## Vocabulary discipline

`CONTEXT.md` stays glossary-only (`**Term**:` / one-sentence definition / `_Avoid_:`). New concepts from this work (Task Source, Work Queue, Wave, Worktree Session, Result Artifact, Escalation, …) get terms there the moment they are decided, and code identifiers follow the glossary. Decisions go to `docs/adr/00NN-*.md`, 1–3 sentences each.

## Suggested skills

Install the **`go-cli`** setup — it now carries the whole spec workflow plus golang-cli, golang-testing, golang-error-handling, golang-concurrency, golang-context, golang-lint, bubbletea v2, tui-design, agentic-cli-design, qa-gate, evidence-gate, tech-writer, and handoff (no cobra; this repo uses stdlib flag dispatch). Optional per-project extras for the redesign phase, installed individually (`bunx skills add marcioaltoe/skills/skills/<collection>/<name>`): `architectural-analysis`, `refactoring-analysis`, `tactical-ddd`, `git-rebase`, `lesson-learned`. When roundfix starts cutting versioned releases, `cut-release` (08-release) exists in the catalog. Session rhythm: `grill-with-docs` before deciding, `write-prd`/`write-techspec`/`write-tasks` to spec each phase, `implement-task` discipline while executing, `evidence-gate` before any completion claim.

## Verification gates

`make verify` must pass 100% (fmt + lint + test + build + `roundfix skills check`). Agent branches use the `ma/` prefix; commits follow Conventional Commits (check this repo's `cog.toml` for scope rules).

## First steps for the fresh session

1. Read `CONTEXT.md`, `AGENTS.md`, `docs/adr/`; skim `docs/product-brief.md`. Note: `AGENTS.md` was already replaced (2026-07-04, uncommitted) with the marcioaltoe/skills `AGENTS.go-cli.md` template adapted to this repo — the old skill-router drift is gone; commit it together with the `.agents/skills` prune resolution.
2. Resolve the `.agents/skills` prune (commit or restore), then install the suggested skills.
3. Run `/setup-workflow` to scaffold `docs/specs/` in this repo.
4. `/grill-with-docs` on work-plan item 1 (task-source abstraction) → `/write-prd` → `/write-techspec` → `/write-tasks` → implement.
