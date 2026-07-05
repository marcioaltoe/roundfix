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
   **Recurred during the review dogfood (2026-07-05)**: the watch supervisor
   needed two rounds of filter re-tuning because agent-console text (skill
   descriptions) kept matching outcome vocabulary like "Clean".

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

11. **RESOLVED 2026-07-05** — the stale "Knowledge workspace" section in
    `docs/agents/issue-tracker.md` was pruned by Marcio's docs/agents update
    (landed inside a0e572d).
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

## Added while the dogfood Run progressed

15. **Snapshot-diff commit scoping trusts a single-writer worktree.** Task_02's
    commit (eaf2327) swept in concurrent user activity: a skills reinstall
    (`.agents/skills/*`, `skills-lock.json`, a `.claude/skills/knowledge-workspace`
    symlink) and this findings file itself. The before-snapshot rule excludes
    only pre-task dirt; anything that appears in the worktree during the task
    window lands in that task's commit. Proper fix is worktree-per-task
    (work-plan item 4); interim mitigations: document "hands off the worktree
    during an Active Run", and/or journal a warning when a task commit stages
    paths outside the roots the task file references. Size: spec (proper) /
    small (warning heuristic).
    **Escalated at task_03**: commit a0e572d swept 53 files / −2,855 lines of
    concurrent user cleanup (`.scratch/_achieved/**`, `docs/plans/`,
    knowledge-workspace removal, `AGENTS.md`, `docs/agents/*`,
    `skills-lock.json`). Code check ruled out a parsing bug —
    `parsePorcelainPaths` does capture deletions, so pre-task dirt is
    excluded correctly; the sweep is pure timing (the next task's
    before-snapshot lands milliseconds after the previous commit, so any
    user change during a long task window postdates it). The interim
    warning heuristic gains urgency.
    **Product stance (Marcio, 2026-07-05)**: concurrent user work — including
    commits made in parallel — must never become a verification concern or an
    impediment; the Daemon tolerates a multi-writer worktree. This rules out
    turning any warning heuristic into a blocker and confirms
    worktree-per-task (work-plan item 4) as the real fix.

16. **Codex full-access through acpx may lack the danger sandbox preset.**
    task_02 verified `acpx@0.12.0` + pinned `@agentclientprotocol/codex-acp@0.0.44`:
    the adapter's full-access mode carries the sandbox internally but
    advertises no sandbox config option; Roundfix attempts `acpx codex set
    sandbox_mode danger-full-access` and journals
    `codex_sandbox_full_access_unavailable` when unsupported (set-mode failure
    stays fatal). Watch whether `--agent-full-access` under codex/acpx still
    unblocks localhost-network verification; if not, report the gap upstream.
    Size: small (observe + upstream issue).

17. **Review artifacts should live with the Spec, not in a separate root.**
    Today Round/Review Issue artifacts go to the user-scoped Artifact
    Directory (`~/.roundfix/artifacts/reviews/pr-<n>/round-NNN/`, ADR-0003),
    completely apart from the spec tree. Now that implement + review form one
    flow, the proposal (Marcio, 2026-07-05) is: when the Open Pull Request
    belongs to a Spec, store its Rounds/Review Issues inside that Spec's
    folder under `docs/specs/<slug>/`; for spec-less reviews, use a
    PR-referenced folder inside `docs/specs/` following the current layout —
    ending scattered files. Design questions for the PRD grill: how a PR maps
    to a Spec (the `Roundfix-Spec` commit trailers give a natural join; an
    explicit `--spec` flag is the fallback), what supersedes ADR-0003, and the
    big one — in-repo artifacts interact with the clean-worktree Preflight
    Validation and snapshot-diff commits (fetch/resolve would dirty the tree
    mid-Run; artifacts would need committing like `qa/` or explicit
    exclusion). Same problem family as finding 15 — consider solving both in
    one spec. Size: spec (own PRD; supersedes ADR-0003).

18. **Watch should poll first, sleep after.** (Marcio, 2026-07-05, review
    dogfood.) On start, watch should immediately check whether the Open Pull
    Request already has comments/Review Issues — an already-settled review
    with waiting feedback should flow straight into fetch — and apply
    `poll_interval` only between subsequent checks. Same instinct for the
    quiet period: skip or shorten it when the Review Source was already
    settled before the Run began. Today the timers run ahead of the first
    useful check, adding dead wait to the common "review finished long ago"
    case. Size: small (watch loop ordering).

19. **Watch has no stdout contract.** (Review dogfood, 2026-07-05: Run
    `run_20260705T112746Z_d514575a110198e5`, tax-poc PR #9, Clean after 1
    Round.) The entire outcome report — fetch counts, batch settlement,
    verification result, Final Push, `reached Clean after 1 Round(s)` — went
    to stderr and stdout ended empty; exit code is the only machine-readable
    result. Inconsistent with the Implement Command's deterministic stdout
    report and the repo's own stdout-carries-requested-output rule.
    Candidate: give watch (and resolve) a deterministic stdout report shaped
    like implement's per-item lines plus one outcome line. Size: small.

20. **Watch should end on merge-readiness, not only on local Clean.**
    (Marcio, 2026-07-05, review dogfood.) Watch terminates on Max Rounds or
    when the local Review Issue set empties — but the PR's real gate is the
    GitHub status check CodeRabbit reports on the head commit (the merge-box
    "pending check" that blocks squash/merge; GitHub vocabulary: commit
    status / check run, readable via `gh pr checks` or the Checks API).
    Observed live: right after the Final Push, the Run was already Clean
    while the PR showed "CodeRabbit — Waiting for status to be reported —
    Review in progress" on the pushed commit. Proposal: watch's until-clean
    should (optionally) keep watching after the Final Push until the
    CodeRabbit check on the final head SHA reports success with no new
    Review Issues — making Clean mean "ready for squash/merge" — still
    bounded by Max Rounds. Size: small/medium (one more status source in the
    watch loop; the review-status seam already exists). Pairs with finding
    18 (poll-first ordering).

21. **Temp-repo tests are not hermetic against user git config.** (Implement
    dogfood, task_06.) On a machine with global `commit.gpgsign=true`, six
    temporary-repo tests failed with `gpg failed to sign the data`; the
    executing agent had to pass `GIT_CONFIG_*` env overrides to reach a green
    gate (the Daemon's own verbatim re-run then passed, so settlement was
    unaffected — likely gpg-agent caching, which also makes this flaky).
    Root fix: test helpers that create git repos must isolate config
    (`GIT_CONFIG_GLOBAL=/dev/null` + explicit `commit.gpgsign=false`/user
    identity in the temp repo), matching the daemon's own `-c` discipline.
    Size: small.

## Product ideas from the dogfood debrief (Marcio, 2026-07-05)

22. **`roundfix setup` — environment bootstrap.** One command that installs
    the pinned acpx (`npm install -g acpx@<pin>`), validates the environment
    (Node version, acpx probe, runtime adapters), and offers to create the
    Project Config file (`.roundfixrc.yml`) interactively. Overlaps with the
    existing Init Command (User/Project Config creation) — decide whether
    setup extends `roundfix init` or wraps it plus dependency installation.
    Repo-dev extra: a self-test mode that runs the gated real-acpx
    integration suite. Size: small/spec.

23. **`roundfix upgrade` + version freshness check.** A command that updates
    the installed roundfix to the latest released version, plus a passive
    check that compares the running version against the latest release and
    suggests the upgrade when behind (non-blocking, stderr). Needs a release
    channel decision (GitHub Releases + `v*` tags exist as the convention).
    Size: small.

24. **Graceful stop for live Runs, `--force` to kill.** `roundfix stop`
    today releases the lock/settles state — the live process is stopped by
    Ctrl-C. Wanted: stopping an Active Run (spec implement or PR review)
    from another terminal, graceful by default — let the current Batch/Task
    settle (verification + commit) then end the Run — and `--force` for
    immediate cooperative cancel of the Agent. Natural transport: the
    DB-mediated control channel from work-plan item 3 (stop already reads
    that way). Size: small/medium.

25. **Push becomes optional Project Config for spec Runs; PR creation stays
    out of scope.** Direction: a repo-level config key letting a spec Run
    push its branch at a Clean outcome (default off, preserving today's
    behavior); opening pull requests is permanently out of Roundfix's scope.
    Adopting this supersedes the "never push" half of ADR-0013 — needs its
    own ADR when picked up. Size: small (behavior) + ADR.
