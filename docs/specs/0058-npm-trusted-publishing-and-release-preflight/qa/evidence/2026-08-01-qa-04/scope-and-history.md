# Scope and protected-tooling history

Build: `e45dd37d2f2ced6dcaa3533fcea939a867b3ea6c`

`rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-08-01-qa-04/current_scope_probe.rb` exited 0.

- All eight Task files are `completed`; each has a `## Result` and acceptance evidence.
- The PRD and TechSpec each account for identifier strategy, authentication and HTTP, active ADR obligations, and tooling authority. Both cite the operative `docs/agents/` sources.
- ADR-0031 and ADR-0048 are active legacy ADRs; ADR-0082 and ADR-0084 carry `status: accepted`. Their npm package layout, read-only Release Plan, all-or-nothing Release Set, and bounded fallback obligations match the Spec artifacts.
- Commit `397227ff` contains the express authorization for exactly `.github/workflows/release.yml` and precedes every tooling Task.
- Exact `git diff-tree --no-commit-id --name-only -r <commit>` assertions passed for tooling Tasks `21bc4bf`, `8d14a67`, `b0052e9`, `47de307`, `ab34e03`, and `e45dd37`. Each contains only `.github/workflows/release.yml`, its own Task file, and Task 02's Spec-owned fixtures where applicable.
- Planning/remediation and consequent-fix chronology remains separated. The new Task 08 graph commit `a8276a4` precedes Task commit `e45dd37`; the latter changes only the authorized workflow and `task_08.md`.
- Review artifacts settle every fetched round-002 issue as two `resolved` and two `invalid` records.
- The worktree delta contains only this QA report and evidence.

Fresh direct `git merge-base --is-ancestor` checks also exited 0 for authorization → first tooling Task, authorization → Task 08, Task 08 graph → Task 08 implementation, and Task 08 implementation → current `HEAD`.
