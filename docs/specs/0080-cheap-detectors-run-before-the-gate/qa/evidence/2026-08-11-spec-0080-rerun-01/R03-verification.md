# R03 — verification tiers and F-003 reproduction

Build: `c2372a9f709c9197aa5c5e89fd71da1ab46f07e6`.

## Required gates

- `rtk make verify` first reached the managed sandbox's denied
  `api.github.com` boundary. The identical authorized rerun reached the Go
  suite and exited 2. `internal/daemon` failed
  `TestQAGatePromptUsesTaskContextBuilderAndPreviousReportIdentity` because the
  mechanical stage could not resolve a Git head (`git` exit 128).
- `rtk make verify-incremental` exited 2 after 1.29 seconds with the same
  failing test. It met the 10-second time budget but did not produce a valid
  local verdict.
- The first authorized full run took about 87 seconds and exited 2. After an
  authorized `rtk go clean -testcache`, the repeated `rtk make verify` exited 2
  at the same test. The second run reused compile artifacts but reran the
  failing case under a new Run id.

## Minimal reproduction

`rtk sh qa/evidence/2026-08-11-spec-0080-rerun-01/R03-daemon-head-reproduction.sh`
exited 1 twice. The focused process took 0.38 seconds on the recorded run and
reported:

```text
task_context_test.go:249: TaskCycle: run QA mechanical stage ...:
resolve mechanical-stage Git head: exit status 128
--- FAIL: TestQAGatePromptUsesTaskContextBuilderAndPreviousReportIdentity
```

## Root-cause trace

`newTaskCycleFixture` in `internal/daemon/task_engine_test.go` assigns
`gitRoot := t.TempDir()` and writes the Spec, but does not initialize a Git
repository. The prompt integration case in
`internal/daemon/task_context_test.go:216-265` uses that fixture and does not
inject a `QAMechanicalStage` fake. Task 03 made `runQAGate` call the real stage
before the prompt Agent, and `mechanicalHead` in
`internal/speccheck/mechanical.go:830-839` unconditionally executes
`git -C <root> rev-parse --verify HEAD^{commit}`. Git therefore exits 128
before the test can observe the prompt.

History confirms the prompt test and its non-Git fixture are byte-present on
both sides of Task 03 commit `1432dd1e`; that commit added the unconditional
mechanical stage but did not adapt this case. Corrective Task 09 commit
`c2372a9f` changed `internal/cli/implement_test.go` and
`internal/daemon/task_engine_test.go`, not `task_context_test.go`.

## Contract inspection

`Makefile` keeps the targets distinct:

```text
verify: fmt-check $(VERIFY_TEST_TARGET) skills-sync-check skills-check build
verify-incremental: fmt-check test skills-sync-check skills-check build
```

The adopted clauses in `docs/agents/agent-instructions.md:26` and
`docs/agents/spec-routing.md:15` assign incremental validity to the local tier
and complete fresh judgement to CI. Both current targets are red, so R03 fails
despite the incremental command's short elapsed time.
