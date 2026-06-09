[[CODEX_LOOP name="roundfix-mvp" goal="implement the Roundfix MVP from .scratch/roundfix-mvp/PRD.md and .scratch/roundfix-mvp/issues with all tasks complete, blockers closed, tracking valid, and rtk make verify PASS"]]

Use the codex-loop skill. Track this implementation under
`.codex/loop/roundfix-mvp/`.

You are implementing the Roundfix MVP in `/Users/marcio/dev/roundfix`.
Roundfix is a local-first Go CLI and future daemon for fetching pull request
review issues, resolving Unresolved Review Issues through a selected ACP
Runtime, and pushing only after no Unresolved Review Issues remain.

## Source of truth

Read these before editing code:

- `AGENTS.md`
- `CONTEXT.md`
- `README.md`
- `docs/product-brief.md`
- `docs/adr/`
- `.scratch/roundfix-mvp/PRD.md`
- `.scratch/roundfix-mvp/issues/`
- current implementation under `cmd/roundfix/`, `internal/app/`, and
  `internal/cli/`

Use the glossary terms from `CONTEXT.md`. Respect every accepted ADR. The PRD
and numbered issue files are the implementation contract.

## Current state

The project is an early Go scaffold. The CLI currently supports help and
version output, while `fetch`, `resolve`, and `watch` are reserved command
names. The MVP work starts from the planning tracker committed on branch
`ma/mvp-prd`.

The local Codex Loop tracking state for `roundfix-mvp` may not exist yet. Start
every continuation by running:

```bash
rtk python3 /Users/marcio/.vault/home/.codex/skills/codex-loop/scripts/detect-next.py roundfix-mvp
```

If the action is `bootstrap`, initialize `.codex/loop/roundfix-mvp/` from the
request in this prompt and the 14 issue files listed below. Do not invent a new
task order or collapse the issues unless an issue is already completed by
current code and you record evidence.

## Tracked task order

Use these local issue files as the ordered task list:

1. `.scratch/roundfix-mvp/issues/01-cli-command-contract-and-exit-codes.md`
2. `.scratch/roundfix-mvp/issues/02-preflight-config-and-artifact-directory.md`
3. `.scratch/roundfix-mvp/issues/03-preflight-git-pr-and-push-safety.md`
4. `.scratch/roundfix-mvp/issues/04-run-state-fetch-runs-and-active-locks.md`
5. `.scratch/roundfix-mvp/issues/05-fetch-coderabbit-round-artifacts.md`
6. `.scratch/roundfix-mvp/issues/06-resolve-compatible-artifacts.md`
7. `.scratch/roundfix-mvp/issues/07-resolve-deduplication-and-batches.md`
8. `.scratch/roundfix-mvp/issues/08-agent-runtime-batch-resolution.md`
9. `.scratch/roundfix-mvp/issues/09-daemon-verification-and-batch-commits.md`
10. `.scratch/roundfix-mvp/issues/10-final-push-and-source-resolution.md`
11. `.scratch/roundfix-mvp/issues/11-watch-review-round-loop.md`
12. `.scratch/roundfix-mvp/issues/12-stop-request-and-sigint-shutdown.md`
13. `.scratch/roundfix-mvp/issues/13-tui-input-and-live-run-view.md`
14. `.scratch/roundfix-mvp/issues/14-roundfix-skills-and-installer.md`

Each issue is `ready-for-agent` and AFK. Work in dependency order. Complete
exactly one printed Codex Loop action per continuation.

## Codex Loop discipline

- Run `detect-next.py` first at the start of every continuation.
- Follow the printed action exactly: `bootstrap`, `execute_task`, `verify`,
  `resolve_blocker`, or `done`.
- For `execute_task`, read the printed task file, parent PRD, relevant ADRs,
  current code, and tests. Work only on that task.
- Before marking a task complete, write the per-iteration memory file with the
  objective, files touched, decisions, validation evidence, blockers, and
  next-run notes.
- Update state only through the Codex Loop helper scripts. Do not hand-edit
  `.codex/loop/roundfix-mvp/state.json`.
- Run `validate-tracking.py roundfix-mvp` before each iteration summary.
- If a blocker appears, record it with `update-tracking.py --blocker` and stop
  after the required summary.
- For `done`, require all tasks completed, no blockers, verification `PASS`,
  and `validate-tracking.py roundfix-mvp --expect-done` success.

## Repository rules

- Prefix shell commands with `rtk`.
- Keep `cmd/roundfix/main.go` thin; put behavior under `internal/...`.
- Prefer the Go standard library until a dependency has a clear job.
- Add dependencies with `rtk go get`, not manual `go.mod` edits.
- Pass `context.Context` first for blocking, IO, process, network, database, and
  daemon-boundary operations.
- Return explicit errors with context and wrap underlying errors with `%w`.
- Own goroutines with cancellation and clear shutdown.
- Do not copy names, branding, package names, comments, examples, or generated
  artifacts from reference projects.
- Do not commit, push, rebase, reset, restore, clean, or remove tracked files
  unless the user explicitly asks for that git operation.
- If a new branch is required, it must use the `ma/` prefix.
- If unexpected user changes exist, read them and work with them. Do not revert
  unrelated work.

## Product boundaries

The MVP must preserve these boundaries:

- The Daemon owns Review Source status checks, CodeRabbit fetches, Round
  creation, Agent process lifecycle, verification, Batch commits, Final Push,
  Review Source resolution, retries, timeouts, and Stop Request handling.
- The Agent owns only assigned issue files, triage, code edits, tests,
  verification commands, and assigned Review Issue status updates.
- Agents must not commit, push, resolve Review Source threads, edit unassigned
  issue files, or mark issues as `duplicated`.
- `fetch` creates a tracked Fetch Run and never starts an Agent, commits, or
  pushes.
- `resolve` works over downloaded Compatible Artifacts and never fetches Review
  Source issues.
- `watch` automates the fetch-resolve-review loop until a documented terminal
  outcome.
- Final Push runs only after no Unresolved Review Issues remain.
- `reprocess` and broad `include_resolved` behavior are out of scope.

## Verification gates

For changed Go files:

```bash
rtk gofmt -w <changed-go-files>
```

For every task, run the smallest relevant targeted tests first. Before marking a
task complete, run the relevant broad gate for the slice. Before loop
verification can pass, run:

```bash
rtk make verify
rtk go run ./cmd/roundfix --help
rtk git diff --check
```

If concurrency, Agent process ownership, watch polling, or Run Database locking
changed, also run:

```bash
rtk go test -race ./...
```

Do not mark verification `PASS` from assumptions, partial output, or prior
runs. Record the exact commands and outcomes in Codex Loop memory.

## Completion criteria

The loop is complete only when:

- every issue file in `.scratch/roundfix-mvp/issues/` has been implemented or
  explicitly proven already satisfied;
- all acceptance criteria in those issue files are met;
- `.codex/loop/roundfix-mvp/state.json` is valid and all tracked tasks are
  completed;
- blockers are empty;
- verification is `PASS` with fresh evidence;
- `validate-tracking.py roundfix-mvp --expect-done` succeeds.
