# Dogfood findings — Implement Command and ACPX Migration

Running log of problems, cosmetic issues, and improvement candidates observed
while shipping spec `0001-implement-command` and dogfooding it against spec
`0002-acpx-migration` (Run `run_20260705T104227Z_b9aabfcb7cb28014`, Codex).
Each entry seeds a future PRD/techspec/issue; nothing here blocks the current
work. Sizes: `cosmetic` (one-line fix), `small` (single task), `spec` (needs
its own PRD/techspec cycle).

## Observed live during the dogfood Run

1. **Task commit titles inherit heading capitalization.** `TaskCommitMessage`
   derives `<type>: <title>` from the `# Task NN: <title>` heading verbatim,
   producing `feat: Build the acpx invocation core` — valid Conventional
   Commits, but the repo's style is lowercase subjects. Fix: lowercase the
   first rune in the derivation, or have `write-tasks` mandate lowercase
   imperative titles. Size: cosmetic.
2. **Run header shows a Run Budget for spec Runs that is never enforced.**
   The implement header prints `Budget: 0s / 2h0m0s` (from `budget.*` config)
   but the Implement Command deliberately does not enforce the Run Budget
   (0001 techspec decision: the graph is finite). Misleading; omit the line
   or label it not-applicable for spec Runs. Size: cosmetic.
3. **Run header shows `Round: -` for spec Runs.** Round is a review-path
   concept; the line is noise on implement Runs. Size: cosmetic.
4. **Non-TTY stderr is a firehose for supervisors.** In non-interactive mode
   stderr carries the boxed Live Run View header plus every Agent console
   event (messages, thoughts, skill text). An automated supervisor has to
   anchor-grep for daemon lines. Candidate: a diagnostics verbosity level
   (e.g. daemon-events-only) for machine supervision. Size: small.

## Carried from 0001 QA report and task Results

5. **Cockpit detail pane is Review-Issue-only.** Enter on a Task should open
   a read-only task-file detail view (task_09 follow-up). Size: small.
6. **`spec.ListActive` silently skips unreadable `_prd.md` files.** Right for
   picker robustness, but there is no diagnostic listing to reveal a broken
   spec folder (task_01 follow-up). Size: small.
7. **Interactive Input has no QA field.** `--qa` is flag-only; the implement
   Interactive Input flow never offers it (task_07 follow-up). Size: small.
8. **Spec-Run agent logs ignore `defaults.artifact_dir`.** Implement Runs
   derive agent logs from `<workdir>/.roundfix/runs/<run-id>/...` because
   `TaskPlan` carries no Artifact Directory, while review Runs log under the
   configured artifact dir. Unify or document the split (0001 task_05 design
   note; store validation also does not require ArtifactDir for implement
   Runs — 0001 task_02 follow-up). Size: small.
9. **QA Report commit message is scoped; this repo's cog rejects scopes.**
   The product contract commits `docs(qa): qa report for <slug> (<verdict>)`,
   but `cog.toml` (`scopes = []`) rejects scoped titles — harmless today
   because CI validates PR titles only, yet the dogfood repo's own lint would
   flag its generated commit. Decide: unscope the generated message or accept
   the exemption explicitly. Size: cosmetic.
10. **Failed-Task resume friction.** A failed Task leaves settled status +
    preserved edits uncommitted; resume preflight demands a clean tree, so
    the user must commit/stash/discard before re-running. Known, documented,
    properly fixed by worktree-per-task (work-plan item 4). Size: spec.

## Repo and docs hygiene

11. **`docs/agents/issue-tracker.md` has a stale "Knowledge workspace"
    section.** It claims `docs` is a symlink into `.knowledge/` and mandates
    `git -C .knowledge` commits; the knowledge-workspace setup was dropped
    (commit 412fad2) and `docs/` is a real directory. Prune the section.
    Size: cosmetic.
12. **Prompt-contract drift remains structural.** `BuildTaskPrompt` mirrors
    the implement-task skill by construction, not mechanism; templating is
    work-plan item 5. Listed for completeness. Size: spec.
13. **`roundfix stop` has no spec-target selector.** Stop accepts run id,
    `--pr`, or head-repo/branch; a `--spec <slug>` selector would give spec
    Runs parity. Run-id works meanwhile (ActiveRunError names it).
    Size: small.
14. **Repo commit-style knowledge is tribal.** `cog.toml` `scopes = []`
    rejecting scoped titles surprised more than one agent session; worth one
    line in the repo agent instructions. Size: cosmetic.
