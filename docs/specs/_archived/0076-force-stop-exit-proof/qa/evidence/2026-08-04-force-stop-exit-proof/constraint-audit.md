# Constraint and gate audit

Build: `c035ebb19dcb6eb81844f5195a0b89abbf99e4e1`

## Authored gate

- `_tasks.md` declares `qa: task_03`.
- `task_03` has `type: qa`, is the only terminal node, and depends on
  `task_02`, the graph's only non-QA leaf.
- `task_01.md` and `task_02.md` both carry `status: completed` before this
  gate; `task_03.md` remains Daemon-owned `status: pending`.

Result: the authored gate is runnable and closes every non-QA leaf.

## Project Constraints

- Identifier strategy: not applicable. The PRD and TechSpec both cite
  `docs/agents/domain.md`; neither Task commit changes an identifier-bearing
  product surface.
- Authentication and HTTP: not applicable. The PRD and TechSpec both cite
  `docs/agents/agent-instructions.md`; the Task commits change no transport,
  credential, authentication, or authorization path.
- Active ADR obligations: ADR-0080, ADR-0089, ADR-0091, and ADR-0093 all carry
  `status: accepted`. ADR-0089 applies to explicit test-process coordination;
  ADR-0080 and ADR-0091 govern this report and graph; the PRD explicitly
  records why ADR-0093 does not apply to the implementation surface.
- Tooling Authority: not applicable. No protected tooling path appears in
  either Task commit, so there is no tooling authorization, prerequisite fix,
  consequent fix, or derived-pin choreography to audit.

All cited operative source paths exist.

## Git evidence

The planning commit predates both Task commits, and the Task commits are in
dependency order. Both ancestry checks exited 0:

```text
git merge-base --is-ancestor 9dadc83 d9be9f5
git merge-base --is-ancestor d9be9f5 c035ebb
```

Task 01 commit `d9be9f58913a91a6f32795c0d80f3f3d1254e274` changed exactly:

```text
docs/specs/0076-force-stop-exit-proof/task_01.md
internal/store/process_unix_test.go
```

Task 02 commit `c035ebb19dcb6eb81844f5195a0b89abbf99e4e1` changed exactly:

```text
docs/specs/0076-force-stop-exit-proof/task_02.md
internal/store/process_unix_test.go
```

The changed paths were resolved with
`git diff-tree --no-commit-id --name-only -r <commit>`, not inferred from Task
Results.
