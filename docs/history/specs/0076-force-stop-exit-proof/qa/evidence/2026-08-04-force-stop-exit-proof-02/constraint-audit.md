# Project Constraint and Task Graph audit

- Build: `eaebd553ad2b415dbcc48e936b5b8afa980e3a89`
- Spec status: active, non-legacy

## Authored QA node

`_tasks.md` declares `qa: task_03`. Its graph is
`task_01 -> task_02 -> task_03`; `task_03` has `type: qa`, is terminal, and
depends on the only non-QA leaf. Task frontmatter reports `task_01` and
`task_02` as `completed`; Daemon-owned `task_03` remains `pending` during this
gate. The precondition therefore passes.

## Constraint shape and applicability

Both `_prd.md` and `_techspec.md` contain all four required rows:

- Identifier strategy: not applicable because no project-owned Internal
  Identifier is created. Operative source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable because this test-only process
  coordination change opens no transport and touches no credential path.
  Operative source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: ADR-0089 applies to explicit test environment and
  coordination. ADR-0080 governs typed blocked-row verdicts; ADR-0091 governs
  the authored QA Task; ADR-0093 is explicitly accounted for as a
  non-applicable relation candidate because this Spec changes no consistency
  detector. All four ADR files report `status: accepted`.
- Tooling Authority: not applicable. No repository tooling mutation is
  proposed or present. Operative source:
  `docs/agents/agent-instructions.md`.

Every cited operative source exists in the current checkout.

## Git scope and chronology

The planning commit `9dadc83` predates both Daemon-owned Task commits. Exact
`git diff-tree --no-commit-id --name-only -r <commit>` results:

```text
d9be9f5
docs/specs/0076-force-stop-exit-proof/task_01.md
internal/store/process_unix_test.go

c035ebb
docs/specs/0076-force-stop-exit-proof/task_02.md
internal/store/process_unix_test.go
```

Neither Task changes a repository-tooling configuration, script, ignore file,
plugin declaration, or version pin. The post-finding repair commit `c6d49d0`
changes only `internal/cli/implement_test.go`; the gate-reset commit `eaebd55`
changes only `task_03.md`. No tooling authorization choreography applies, and
no Task commit exceeds its assigned Task file plus implementation surface.
