---
spec: 0054-tooling-task-and-verification-hygiene
status: active
created: 2026-07-28
surfaces: [backend, infra, docs]
---

# Tooling task and verification environment hygiene

Four recurring, process-only failure modes cost the dogfood loop roughly three
blocked Runs per day on 2026-07-27/28, and none of them was a defect in the
product being built. Editing a Roundfix-owned Skill strands its own Task
because five derived digest pins have no regeneration path; landing any
protected-tooling change requires a commit choreography no document states and
a repository that is already green everywhere; `make verify` fails before
compilation when the ACP sandbox denies the host Go cache; and a bare
`go build ./cmd/roundfix` drops an unignored 20 MiB binary that Daemon commits
sweep in. Evidence:
[derived digest pins have no regeneration path](../../findings/2026-07-27-derived-skill-digest-pins-have-no-regeneration-path.md),
[tooling Tasks need a green repo and an undocumented commit choreography](../../findings/2026-07-28-tooling-tasks-need-a-green-repo-and-an-undocumented-commit-choreography.md),
[sandboxed Agents cannot reach the default Go cache](../../findings/2026-07-27-sandboxed-agents-cannot-reach-the-default-go-cache.md), and
[bare go build writes an untracked root binary](../../findings/2026-07-27-bare-go-build-writes-an-untracked-root-binary.md).

## Project Constraints

- Identifier strategy: not applicable — this Spec creates no project-owned
  Internal Identifier; Skill names, digest values, and Make target names keep
  their existing identities. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — every change is local build
  tooling, digest computation, Daemon staging, documentation, or policy; no
  authentication or HTTP surface is touched. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0045 keeps review Runs requiring a
  clean tracked working tree, which the stray-binary fix protects; ADR-0038
  keeps the single Verification repair contract; ADR-0057 keeps Task status
  Daemon-owned; ADR-0081 makes sanctioned digest regeneration fallout of the
  authorized Skill edit. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-28, the maintainer expressly
  authorizes changes to exactly `Makefile` and `.gitignore`, to exactly
  `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`, and to
  the deterministic Skill-digest fallout in exactly
  `internal/baseline/assets/setups/go-cli.json`,
  `internal/baseline/assets/setups/rust-cli.json`,
  `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`. No other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- Editing a Roundfix-owned Skill plus one sanctioned regeneration command
  leaves the repository gate green with no hand-transcribed hashes and no
  per-Spec authorization amendment.
- The repository gate behaves identically inside and outside the ACP
  sandbox, with no Agent-discovered environment workarounds.
- A bare compile check can no longer litter the working tree or reach a
  Daemon commit.
- A protected-tooling Task that cannot succeed fails before the Agent works,
  naming the real precondition, and the commit choreography it needs is a
  documented rule rather than folklore.

## User Stories

1. As an Agent assigned a Skill-alignment Task, I want a sanctioned command
   that regenerates every derived digest pin from the canonical sources, so
   that my Task's Verification is satisfiable without exceeding my
   authorization.
2. As a maintainer, I want stale digest pins to fail the gate with a message
   naming the regeneration command, so that the failure reads as a stale
   snapshot instead of a broken catalog.
3. As an Agent running the repository gate inside the ACP sandbox, I want
   the Go build cache to default to a repository-local path, so that the
   gate compiles without my discovering an environment workaround.
4. As a supervisor landing a protected-tooling change, I want the commit
   choreography — the authorization record and any prerequisite fix each in
   their own commit before the Task commit, then the rest — stated in the agent guidance, so that landing it takes one attempt
   instead of three QA cycles.
5. As a maintainer, I want a Task whose Verification demands a green
   repository to be refused at assignment when the repository is red for
   unrelated reasons, so that a correct Task is not settled failed on
   someone else's breakage.
6. As a user whose Agent ran a bare compile check, I want the resulting
   binary ignored and unstageable, so that no Task commit ever ships an
   executable again.

## Core Features

1. A sanctioned regeneration command — a Make target — recomputes every
   derived digest pin from the canonical Skill sources: the per-skill
   content digests and top-level digests of every catalog setup snapshot,
   the normalized catalog snapshot and its digest, and the parity-corpus
   fixture and manifest rows (up to seven files; a non-roundfix Skill edit
   cascades into all three setup snapshots). Running it after any
   Roundfix-owned Skill edit leaves the gate green, and it is idempotent.
   The same command covers the second derived chain, exercised manually on
   2026-07-28: editing a Baseline module's clauses regenerates the
   source-baseline manifest's byte-range entries with its identity and
   index digests, the formatter golden fixtures with the profile's pinned
   golden digest, and the catalog compatibility fixtures.
2. A stale pin fails verification with a diagnostic naming the regeneration
   command instead of raw digest mismatches fanning out across the Baseline
   suites.
3. Standing authorization policy per ADR-0081: derived pins regenerated by
   the sanctioned command are deterministic fallout of an expressly
   authorized Skill edit and need no separate per-Spec express
   authorization; the policy is recorded in the repository's agent guidance
   and the Skill-alignment Task template stops enumerating pin paths.
4. The commit choreography for protected-tooling changes is documented as a
   rule in the agent guidance: the express authorization record and any
   prerequisite fix are each their own commit and each lands before the
   authorized Task commit, in either relative order, with everything else
   after — and the QA
   gate's tooling audit reports every authorization-shape problem it finds
   in one pass rather than one per rerun.
5. The repository gate's Go build cache defaults to a repository-local,
   ignored path when the environment sets none, making the gate
   deterministic under the ACP sandbox; an explicitly exported cache still
   wins.
6. The bare-build binary path is ignored, and Daemon Task staging refuses to
   stage executable files as Task output, reporting the refused path instead
   of committing it.
7. A Task whose declared Verification requires the repository-wide gate is
   preceded by a repository-green precondition check; a red repository
   fails the Task assignment early with the failing evidence, before any
   Agent session is created.

## User Experience

- `make baseline-digests` (name settled in the TechSpec) prints each
  regenerated artifact and is idempotent; a second run reports no changes.
- A stale-pin gate failure prints one line naming the command to run.
- `make verify` output is byte-comparable between a developer shell and a
  sandboxed Agent session.
- A refused executable staging names the file, its mode, and the rule; a
  refused red-repo tooling assignment names the failing check it saw.

## Non-Goals / Out of Scope

- Weakening the QA gate's tooling-authorization audit or the Agents' refusal
  to exceed authorization — the fix is tooling and policy, not looser
  boundaries.
- Editing the setup-owned managed blocks of the agent guidance; policy text
  lands in repository-authored carriers.
- Changing production Run Worktree git-maintenance behavior (the disposable
  test-repo hardening shipped with Spec 0042).
- Shrinking the digest chain itself (dropping the parity-corpus pin) —
  considered, deferred until the regeneration command proves insufficient.
- Any change to what `make verify` verifies beyond cache determinism and the
  stale-pin diagnostic.

## Success Metrics

- Editing a Roundfix-owned Skill and running the sanctioned command leaves
  `make verify` green with no hand-edited hashes; a Skill-alignment Task
  lands in one attempt with no QA authorization finding.
- A stale pin fails with the regeneration command named in the output.
- `make verify` succeeds from a Task Worktree under the ACP sandbox with no
  Agent-set environment variables.
- `go build ./cmd/roundfix` leaves nothing that `git status` reports; a Task
  commit that would include an executable file is refused and reported.
- A tooling Task assigned against a red repository reports the precondition
  instead of failing on an unrelated test after the Agent's work.
- The commit choreography is discoverable from the agent guidance without
  reading a QA report.

## Decisions

- Derived digest pins regenerated by the sanctioned command are fallout of
  the authorized Skill edit and need no separate express authorization. See
  [ADR-0081](../../adr/0081-sanctioned-digest-regeneration-is-fallout-of-the-authorized-edit.md).
- The Go cache default changes only when the environment sets none —
  explicit `GOCACHE` always wins, preserving developer and CI overrides.
- The executable-staging guard lives in the Daemon so it protects every
  repository, with the `.gitignore` entry as this repository's second layer.

## Open Questions

None.
