# General Agent Instructions

Template for Go CLI projects: command-line tools, automation CLIs with stdlib
`flag` dispatch, and Bubble Tea v2 / Lip Gloss TUIs.

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
  Any format failure, vet failure, lint warning, or test failure is
  **blocking** — zero tolerance.
- **HARD RULE — validation configuration is a specification**: NEVER modify lint, OnionCry, formatter, typecheck, test-runner, or verification configuration to eliminate errors or warnings. Fix the reported code, tests, paths, or architecture instead. If the required conforming change is very large and a configuration change may be justified, STOP and ask the user explicitly whether they want to change that configuration before editing it.
- **NEVER** use workarounds in production code or tests. Fix the root cause.
- **NEVER** hand-edit `go.mod`/`go.sum`. Use `rtk go get` / `rtk go mod tidy`.
- **ABSOLUTELY FORBIDDEN**: `git reset`, `git checkout --`, `git restore`,
  `git clean`, or any command that discards working-directory changes
  **WITHOUT EXPLICIT USER PERMISSION**. These can permanently lose code.
- Agent-created branches **MUST** use the `ma/` prefix unless the repo
  documents a different human-owned prefix.
- **ALWAYS** use the AskUserQuestion tool for confirmations, clarifying questions, decision points, and any needed user interaction. If this CLI has no such tool, ask as a plain message and stop until the user answers — **NEVER** guess an answer the user can give cheaply.

## Knowledge workspace

<!-- Shared block: keep this section identical across samples/AGENTS.*.md; edit all copies together. -->

Long-lived documentation (`CONTEXT.md`, `docs/` — including `docs/specs/` and `docs/adr/` — and legacy `.compozy`/`.scratch`) lives in the central knowledge repository, mounted as a sparse checkout at `.knowledge/` (gitignored) and exposed through versioned symlinks. **ALWAYS USE** the `knowledge-workspace` skill before bootstrapping the workspace or committing any documentation change.

- **MUST** commit documentation changes inside the workspace — `git -C .knowledge add projects/<project> && git -C .knowledge commit -m "docs: ..." && git -C .knowledge push origin main` — never in the code repository. This includes spec artifacts and `task_NN.md` status flips: one piece of work often produces two commits (code here, docs in `.knowledge`).
- **NEVER** commit `.knowledge/` and **NEVER** replace the symlinks with real files.
- If `.knowledge/` or the symlinks are missing, run `scripts/knowledge-bootstrap.sh` once (`--adopt-local` when local docs already exist).

## Agent docs

Read these only when relevant to the task:

- `docs/specs/<feature-slug>/` — spec artifacts (`_idea.md`, `_prd.md`,
  `_techspec.md`, `_tasks.md`, `task_NN.md`, `qa/`); completed specs (all
  tasks done, QA passed) move to `docs/specs/_archived/`. Run
  `setup-context-driven` once if the layout is missing.
- `docs/agents/spec-routing.md` — which pipeline stages a change runs
  through (large initiative / feature / refactor-bugfix / trivial) and what
  marks a spec done
- `docs/agents/issue-tracker.md` — the local `docs/specs/` tracker conventions
- `docs/agents/triage-labels.md` — label mapping for issue triage skills
- `docs/agents/domain.md` — how agents consume `CONTEXT.md` and ADRs
- `CONTEXT.md` — project vocabulary, command concepts, domain rules, and
  product decisions
- `docs/adr/` — architectural decisions; flag conflicts before overriding them

**ALWAYS** use canonical terms from `CONTEXT.md` in command names, help text,
issue titles, test names, and user-facing explanations. If the right term is
missing, call out the gap instead of inventing new language.

## Skill dispatch

Before editing, identify the task domain and **activate every matching skill**:

- **Feature discovery or product idea**: Use `brainstorming`; product-level
  ideas go through `write-idea` (scored by `business-analyst`, debated by
  `council`, challenged by `the-fool`)
- **PRD, tech spec, or task breakdown**: Use `write-prd`, `write-techspec`,
  `write-tasks`. Pick the pipeline entry point per
  `docs/agents/spec-routing.md`: large/fuzzy initiatives start at
  `write-idea`, standard features at `write-prd`, refactors and bug fixes at
  `write-techspec` (it mints the spec folder with a minimal `_prd.md`);
  trivial one-line changes skip the pipeline entirely
- **Executing spec tasks**: Use `implement-task` (one task) or `implement-spec`
  (the whole graph in dependency order)
- **Final QA and QA test execution for a completed spec**: Use `qa-gate`
  exclusively. It owns QA planning, execution, evidence, findings, verdict,
  and the dated report. Keep `golang-testing` and `testing-boss` limited to
  writing and maintaining implementation tests; they do not replace the QA
  workflow. On QA pass, the spec is complete and `archive-spec` runs
  automatically (merge/release is a separate user-driven step, never the
  archive gate)
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
- Errors: wrap with `%w`; **NEVER** `panic` or `log.Fatal` outside `main`.
- Context-first signatures for anything that can block or be cancelled.
- Every goroutine has an owner and a shutdown path. No fire-and-forget.
- TUI code uses **Bubble Tea v2 module paths** (`charm.land/bubbletea/v2`,
  `charm.land/lipgloss/v2`) and the v2 API (`tea.KeyPressMsg`, `tea.Key`).
  Drive `model.Update(...)` synchronously in tests — no terminal emulation.
- Keep packages cohesive; no generic utility packages unless they remove real
  duplication across multiple packages.
- Prefer dependency injection through small interfaces at the boundary that
  owns the behavior.
- Tests assert **observable behavior**. No production-only hooks for tests.
- Keep command parsing, execution, and output formatting separable enough to
  test without shelling out for every case.

## Verification

Use the repo's declared verification command when present, commonly
`make verify`. If no aggregate command exists, run at least:

```bash
gofmt -l .
go test ./...
go build ./...
```

Run `go test -race ./...` whenever concurrency changed. Also run any
documented vet, lint, or integration-test command that applies to the touched
files. **Skipping any verification check invalidates the completion claim.**

## Git and delivery

- **MUST** check `git status --short` before staging; keep unrelated user
  changes out of your diff.
- Use `conventional-commits` for commits and PR titles (check `cog.toml` for
  the repo's scope rules).
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
