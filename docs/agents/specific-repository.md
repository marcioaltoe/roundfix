# Agent instructions

This repository is Roundfix: a local-first Go CLI and future daemon that picks
up work items (today: PR review issues; next: spec task graphs), resolves them
through the user's selected ACP runtime, and pushes only when nothing
unresolved remains. Stdlib `flag` dispatch and a Bubble Tea v2 TUI.

  Any format failure, test failure, or build failure is **blocking** — zero
  tolerance. CI validates PR titles only; the local gate is the ONLY gate.
- Keep the project KISS: prefer the smallest behavior that satisfies the
  documented product contract.
- **NEVER** copy names, branding, package names, comments, examples, or
  generated artifacts from reference projects into this repository.
- Agent-created branches **MUST** use the `ma/` prefix.
- **HARD RULE — roundfix skill sync**: before opening any PR, confirm the
  roundfix skill still matches the shipped CLI behavior; a PR that changes
  CLI behavior ships the skill update too. Contract:
  `docs/agents/skill-dispatch.md`.
- **HARD RULE — skill ownership**: repo-owned authorial workflow skills may
  be adapted locally; every other skill is upstream-managed and **MUST NOT**
  be modified here. Ownership split: `docs/agents/skill-dispatch.md`.
- **HARD RULE — protected-tooling commit choreography**: land prerequisite
  fixes first, the express authorization record with its exact bounded paths
  second, the authorized Task commit third, and everything else afterward.
  Folding a prerequisite fix or the authorization record into the Task commit
  fails the tooling-authority gate.
- **HARD RULE — sanctioned digest regeneration**: after an expressly
  authorized Roundfix-owned Skill or Baseline module edit, run
  `make baseline-digests`. Every derived pin rewritten by that command is
  deterministic fallout of the authorized source edit and needs no separate
  express authorization. A hand-edited pin value remains an unauthorized
  mutation.
- **HARD RULE — release planning**: release work starts with the read-only
  `roundfix release plan` before changelog, version, tag, push, package, asset,
  or GitHub Release mutation. A generic release request authorizes only a
  conclusive patch plan; minor, major, version-zero breaking, and manual
  classification outcomes require the decisions in
  `docs/user-guide/release-runbook.md`.
- **HARD RULE — durable knowledge flows upstream only**: Specs are downstream
  results of the CONTEXT-driven workflow, never sources it depends on — an
  archived Spec may be deleted at any time, so durable knowledge a Spec
  produced moves to its semantic owner (`CONTEXT.md`, an accepted ADR, an
  agent guide, or `docs/references/`) before or at archive. `CONTEXT.md` and
  the agent guides reference accepted ADRs and **NEVER** reference
  `docs/specs/` or `docs/findings/` content; findings are dated reports that
  become Specs, not reference material.
- Project map: `cmd/roundfix/` is the thin CLI entry point; behavior lives in
  `internal/...` (`internal/cli/` owns parsing, output, and exit behavior;
  `internal/app/` holds app metadata)

### Skill dispatch

Skill triggers, ownership, and the Roundfix skill-sync contract live in
`docs/agents/skill-dispatch.md`.

## Verification

The full gate is `make verify` (fmt-check + test + `roundfix skills check` +
build) and it **MUST** pass 100% before any completion claim. For the smallest
relevant gate while iterating:

The Makefile defaults `GOCACHE` to the ignored repository-local `.gocache`
directory only when the environment does not set it; an explicit `GOCACHE`
always wins. After an authorized Roundfix-owned Skill or Baseline module edit,
run `make baseline-digests` before `make verify`. The sanctioned command owns
every derived pin update; never transcribe a digest by hand.

```bash
rtk gofmt -w <changed-go-files>
rtk go test ./...
rtk go run ./cmd/roundfix --help
```

If concurrency changed, also run:

```bash
rtk go test -race ./...
```

- Use `conventional-commits` for commits and PR titles (check `cog.toml`).
- Commit and PR titles are unscoped Conventional Commits subjects here
  (`cog.toml` sets `scopes = []`).
- PR bodies summarize changes, call out risk, and list validation commands run.

## Anti-patterns (immediate rejection)

1. Introducing Cobra, testify, or any dependency the stdlib covers
