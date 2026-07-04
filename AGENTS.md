# Agent instructions

This repository is Roundfix: a local-first Go CLI and future daemon that picks
up work items (today: PR review issues; next: spec task graphs), resolves them
through the user's selected ACP runtime, and pushes only when nothing
unresolved remains. Stdlib `flag` dispatch and a Bubble Tea v2 TUI.

## High priority

- **MANDATORY**: Use the relevant local skills before changing code, docs,
  tests, workflows, or agent instructions. Skill activation comes BEFORE any
  planning or code generation for that domain.
- **ALWAYS** prefix shell commands with `rtk` when it is available. In command
  chains, prefix each command.
- **MUST** use `rg` / `rg --files` for local code search. Use `context7` for
  external library/API docs and `exa-web-search` for broader web/source
  research. **NEVER** use web research tools to search local code.
- **MUST** run the repo's full verification gate before claiming completion.
  Any format failure, test failure, or build failure is **blocking** — zero
  tolerance. CI validates PR titles only; the local gate is the ONLY gate.
- **NEVER** use workarounds in production code or tests. Fix the root cause.
- **NEVER** hand-edit `go.mod`/`go.sum`. Use `rtk go get` / `rtk go mod tidy`.
- Keep the project KISS: prefer the smallest behavior that satisfies the
  documented product contract.
- **NEVER** copy names, branding, package names, comments, examples, or
  generated artifacts from reference projects into this repository.
- If unexpected user changes exist, read them and work with them. **NEVER**
  revert unrelated work.
- **ABSOLUTELY FORBIDDEN**: `git reset`, `git checkout --`, `git restore`,
  `git clean`, commits, pushes, rebases, or removal of tracked files
  **WITHOUT EXPLICIT USER PERMISSION**. These can permanently lose code.
- Agent-created branches **MUST** use the `ma/` prefix.
- **ALWAYS** use the AskUserQuestion tool for confirmations, clarifying questions, decision points, and any needed user interaction. If this CLI has no such tool, ask as a plain message and stop until the user answers — **NEVER** guess an answer the user can give cheaply.
- **HARD RULE**: before opening any PR, confirm
  `.agents/skills/roundfix/SKILL.md` still matches the shipped CLI behavior
  (commands, flags, output formats, exit codes, Batch contract semantics). If
  the PR changes any of those, the skill update ships in the same PR. The
  project-local copy at `.agents/skills/roundfix/` is canonical and its
  `metadata.version` tracks the released CLI version (the `v*` tag), not an
  independent skill version; the embedded `skills/roundfix/` is generated from
  it with `make skills-sync`, and `make verify` fails on drift.

## Agent docs

Read these only when relevant to the task:

- `CONTEXT.md` — the project glossary (vocabulary contract for code, docs,
  prompts, and TUI copy)
- `docs/adr/` — accepted architectural decisions; flag conflicts before
  overriding them
- `docs/product-brief.md` — the product contract from the grill session;
  supersede its decisions via new ADRs, don't edit history
- Project map: `cmd/roundfix/` is the thin CLI entry point; behavior lives in
  `internal/...` (`internal/cli/` owns parsing, output, and exit behavior;
  `internal/app/` holds app metadata)

**ALWAYS** use canonical terms from `CONTEXT.md` in command names, help text,
issue titles, test names, and user-facing explanations. If the right term is
missing, call out the gap instead of inventing new language.

## Agent skills

### Issue tracker

Tasks live as local markdown under `docs/specs/<feature-slug>/` (the canonical
source — no external tracker). See `docs/agents/issue-tracker.md`.

### Domain docs

This is a single-context repo: root `CONTEXT.md` plus ADRs in `docs/adr/`. See
`docs/agents/domain.md`.

### Spec artifacts

Feature specs live under `docs/specs/<feature-slug>/` (`_idea.md`, `_prd.md`,
`_techspec.md`, `_tasks.md`, `task_NN.md`, `qa/`). Dependencies live only in
`_tasks.md`; task status lives only in each task file's frontmatter. Shipped
specs are archived to `docs/specs/_archived/`.

## Skill dispatch

Before editing, identify the task domain and **activate every matching skill**:

- **Feature discovery or product idea**: Use `brainstorming`; product-level
  ideas go through `write-idea` (scored by `business-analyst`, debated by
  `council`, challenged by `the-fool`)
- **PRD, tech spec, or task breakdown**: Use `write-prd`, `write-techspec`,
  `write-tasks`
- **Executing spec tasks**: Use `implement-task` (one task) or `implement-spec`
  (the whole graph in dependency order)
- **Final QA of a completed spec**: Use `qa-gate`; archive after release with
  `archive-spec`
- **MANDATORY** for CLI behavior, flags, stdout/stderr, exit codes, JSON
  output, dry-run behavior, non-interactive mode, or introspection:
  `agentic-cli-design`
- **ALWAYS USE** `golang-cli` before writing Go command behavior, package
  layout, version output, or command tests. CLI style is stdlib `flag.FlagSet`
  dispatch with a `Run() int` exit-code contract — **no Cobra**.
- **ALWAYS USE** `golang-error-handling` for error paths: `%w` wrapping,
  `errors.Is`/`As`, sentinels
- **ALWAYS USE** `golang-concurrency` before goroutines, channels, worker
  pools, or anything with leak/race exposure
- **ALWAYS USE** `golang-context` for context propagation, cancellation, and
  timeouts
- **Lint config or nolint**: Use `golang-lint` (golangci-lint, vet,
  staticcheck discipline)
- **Tests, fixtures, golden files, integration tests**: Use `golang-testing`
  plus `testing-boss`
- **Bubble Tea or Lip Gloss TUI work**: Use `bubbletea` and `tui-design`
- **Implementation**: Use `coding-guidelines`
- **Bug fix or failing test**: Use `no-workarounds` plus `systematic-debugging`
- **Docs, PRDs, ADRs, issues, PR descriptions**: Use `tech-writer`
- **Commits or PR titles**: Use `conventional-commits`
- **Completion claim**: Use `evidence-gate`
- **Session handoff**: Use `handoff`
- **Roundfix dogfooding or assigned-Batch contract checks**: Use `roundfix`
  when driving Roundfix against an Open Pull Request or validating the Batch
  resolution contract

## CLI behavior

- Design commands for humans **and** agents: deterministic output,
  non-interactive flags, stable exit codes, machine-readable modes.
- **MUST** keep stdout for requested command output only. Diagnostics,
  progress, and warnings go to stderr.
- Command names, flag names, JSON fields, and exit-code contracts are
  **public API** — never change them casually.
- Help text **MUST** be concise, truthful, and backed by implemented behavior.
- Errors **MUST** name the failed operation and the next useful action when
  one is known.

## Go conventions

- **Stdlib first**: no new dependency without a clear job the stdlib cannot
  do. Justify every `go get` in the PR body.
- **Zero test dependencies**: stdlib `testing` only — table tests, hand-rolled
  fakes, buffer-captured CLI runs (`Run(args, &stdout, &stderr) int`). **Do
  NOT introduce** testify, mockery, or TUI test harnesses.
- Errors: wrap with `%w`; **NEVER** `panic` or `log.Fatal` outside
  unrecoverable startup in `main`.
- Context-first signatures for blocking, IO, process, network, database, and
  daemon-boundary operations.
- Every goroutine has an owner, cancellation, and a clear shutdown path. No
  fire-and-forget.
- TUI code uses **Bubble Tea v2 module paths** (`charm.land/bubbletea/v2`,
  `charm.land/lipgloss/v2`) and the v2 API (`tea.KeyPressMsg`, `tea.Key`).
  Drive `model.Update(...)` synchronously in tests — no terminal emulation.
- Keep `cmd/roundfix/main.go` thin. Push behavior into `internal/...`.
- Keep packages cohesive; no generic utility packages unless they remove real
  duplication across multiple packages.
- Prefer dependency injection through small interfaces at the boundary that
  owns the behavior.
- Tests assert **observable behavior**. No production-only hooks for tests.
- Keep exported comments short and useful: comment invariants and protocol
  edge cases, not obvious assignments.

## Verification

The full gate is `make verify` (fmt-check + test + `roundfix skills check` +
build) and it **MUST** pass 100% before any completion claim. For the smallest
relevant gate while iterating:

```bash
rtk gofmt -w <changed-go-files>
rtk go test ./...
rtk go run ./cmd/roundfix --help
```

If concurrency changed, also run:

```bash
rtk go test -race ./...
```

If any required gate fails, report the failing command and do not claim the
task is complete. **Skipping any verification check invalidates the
completion claim.**

## Git and delivery

- **MUST** check `git status --short` before staging; keep unrelated user
  changes out of your diff.
- Use `conventional-commits` for commits and PR titles (check `cog.toml`).
- **NEVER** rewrite unrelated files or format the whole repo unless asked.
- PR bodies summarize changes, call out risk, and list validation commands run.

## Anti-patterns (immediate rejection)

1. Introducing Cobra, testify, or any dependency the stdlib covers
2. Marking a spec task `completed` without fresh verification evidence
3. Tracking progress in `_tasks.md` — status lives only in each `task_NN.md`
4. Asking for confirmation before running spec tasks — invocation is the
   authorization
5. Writing to stdout anything that is not the requested command output
6. Changing exit codes, flags, or JSON fields without treating it as a
   breaking API change
