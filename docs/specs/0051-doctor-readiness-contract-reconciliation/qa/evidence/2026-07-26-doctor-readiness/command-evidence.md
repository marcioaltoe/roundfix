# Command evidence — Doctor readiness reconciliation

Build: `bea42860bc7b1ed0a5d83ed14eb80ebeb58565ac`

## PC-01 — Project Constraints and protected tooling scope

Status: pass.

- Both `_prd.md` and `_techspec.md` contain all four required Project
  Constraint categories. Identifier strategy and authentication/HTTP are
  explicitly not applicable with reasons; active ADR obligations and Tooling
  Authority are applicable. Each row cites its operative source under
  `docs/agents/`.
- The cited ADRs 0049, 0055, 0066, 0072, and 0077 are active and consistent
  with the snapshot.
- `rtk git log --format=... -5` resolved the five Daemon-owned Task commits and
  their `Roundfix-Spec`/`Roundfix-Task` trailers.
- `rtk git -c core.fsmonitor=false diff-tree --no-commit-id --name-only -r
  99af3a78f9cebdfa4f4fd5619aaf267a7f2cbf23` exited 0 and listed only:

  ```text
  docs/specs/0051-doctor-readiness-contract-reconciliation/task_01.md
  go.mod
  go.sum
  ```

- `rtk git -c core.fsmonitor=false diff-tree --no-commit-id --name-only -r
  bea42860bc7b1ed0a5d83ed14eb80ebeb58565ac` exited 0 and listed only:

  ```text
  .agents/skills/roundfix/SKILL.md
  docs/specs/0051-doctor-readiness-contract-reconciliation/task_05.md
  skills/roundfix/SKILL.md
  ```

- `rtk git -c core.fsmonitor=false status --porcelain` exited 0 and showed only
  the new Spec-local `qa/` directory. No current delta overlaps a protected
  tooling path.

The two tooling Tasks match the exact authorization and their Task-file
allowlists. No missing authorization, untraceable scope, or out-of-scope
tooling mutation blocks flow QA.

## SG-01 — Full repository Verification

Status: fail.

`rtk make verify` exited 2 at the repository `test` target:

```text
make: *** [test] Error 1
rtk go test ./...
Go test: 2418 passed, 2 failed, 2 skipped in 23 packages
skills (111 passed, 2 failed)
  [FAIL] TestAuthorialSkillSync/typescript-bun.json
  [FAIL] TestAuthorialSkillSync
```

The exact minimal reproduction was:

```text
$ rtk proxy go test ./skills -run 'TestAuthorialSkillSync/typescript-bun.json' -count=1 -v
=== RUN   TestAuthorialSkillSync
=== RUN   TestAuthorialSkillSync/typescript-bun.json
    baseline_skill_contract_test.go:321: roundfix contentDigest = "d5269bac642a8ef4c1c8439eff308594e36ad786fe14d9f3047a8e731998f59f", want canonical "1e4ea59e0d6e553e42c6c63e16ad95165a78be8bbf31b8e0cd8b56ce13cc9146"
--- FAIL: TestAuthorialSkillSync (0.01s)
    --- FAIL: TestAuthorialSkillSync/typescript-bun.json (0.00s)
FAIL
FAIL roundfix/skills 0.266s
FAIL
```

Root-cause evidence:

- `skills/baseline_skill_contract_test.go:289-324` hashes each canonical
  repository-owned skill and requires every Baseline setup snapshot to carry
  that exact digest.
- `internal/baseline/assets/setups/typescript-bun.json:1030` records the stale
  Roundfix digest
  `d5269bac642a8ef4c1c8439eff308594e36ad786fe14d9f3047a8e731998f59f`.
- The current canonical `.agents/skills/roundfix/` computes
  `1e4ea59e0d6e553e42c6c63e16ad95165a78be8bbf31b8e0cd8b56ce13cc9146`.
- Task commit `bea42860bc7b1ed0a5d83ed14eb80ebeb58565ac`
  changed the canonical and shipped Roundfix Skill pair but did not change
  `internal/baseline/assets/setups/typescript-bun.json` or the related parity
  fixture. `git diff --exit-code` over those Baseline paths exited 0.
- `rtk make skills-sync-check` exited 0 with four passing tests, confirming
  that Task 05's narrower declared gate does not run
  `TestAuthorialSkillSync`.
- `rtk git -c core.fsmonitor=false status --porcelain` still listed only the
  Spec-local `qa/` directory, so QA execution did not cause the mismatch.

Classification: code/product-state caused. The Spec requires the Roundfix
Skill bytes to change while its Non-Goals prohibit Baseline asset changes, but
the repository Verification contract requires those bytes and the Baseline
snapshot digest to stay synchronized. The QA skill requires flow execution to
stop at this point.
