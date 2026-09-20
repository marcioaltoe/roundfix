# QA run evidence — Spec 0151 rerun 01

Audited head: `05757c2e1058d85609ade1d68b25943e7aa22588`.
Audited CLI: `roundfix 0.15.0 (05757c2e, built 2026-09-20 00:46:22
-0300)`. The Daemon's auditing binary was
`0.15.0 (30a41b83, built 2026-09-20 00:19:50 -0300)` and was stale because its
build commit predates the audited head.

## Preconditions and static gate

- `roundfix spec check 0151-a-supersession-the-archive-can-see --strict`
  exited `0` with no findings. The installed `roundfix 0.14.1` reported four
  non-equivalent skips: `SC-ROLLUP-MEMBER`, `SC-BACKLOG-UNMOVED`, the
  vocabulary documentation detector, and `SC-REF-UNRESOLVED`.
- `./bin/roundfix spec check
  0151-a-supersession-the-archive-can-see --strict` repeated that result with
  the current-tree binary.
- The Daemon supplied `make verify` as passed with exit `0`; diagnostics were
  removed on success. The QA Agent did not rerun it.
- `_tasks.md` names `task_04` as the unique terminal QA node. Its dependencies
  `task_01`, `task_02`, `task_03`, `task_05`, and `task_06` all declare
  `status: completed`.

## Real repository journey and outside evidence

Outside evidence came from the repository's real active Specs 0128 and 0129
and archived Spec 0147. This Spec did not author those artifacts. Copies were
exercised through the current public binary in disposable Git repository
`/private/tmp/roundfix-qa0151-r2-final.BYuAiL/repo`. This section records the
fresh rerun after the row's input declaration was corrected to cover the full
repository.

1. `roundfix archive 0128-release-planning-with-bare-stable-tags` exited `2`
   and named the absent `_tasks.md` plus the write-tasks remediation.
2. `roundfix supersede --spec 0128-release-planning-with-bare-stable-tags --by
   0147-a-planner-that-reads-both-tag-spellings --reason ...` exited `0`.
3. The repeated archive command exited `0` and moved Spec 0128 to
   `docs/history/specs/0128-release-planning-with-bare-stable-tags`.
4. `diff -qr` between a saved post-amendment copy and the archived destination
   exited `0`.
5. A second public supersede command used the supersession-archived Spec 0128
   as the `--by` deliverer for Spec 0129 and exited `0`.
6. The archived Spec 0128 PRD hash remained
   `45dd42316513ac8a91a4142c7de32d0989087abe6506511cb03233695905d22a`,
   and the full archived directory still matched the saved copy.

## Refusals and amendment preservation

Nine refusal probes ran through the current public binary in disposable Git
repository `/private/tmp/roundfix-qa0151-r3-final.Y9EkUX/repo`. This was a
fresh rerun after correcting the row's input declaration. Each exited `2`
and left `git status --porcelain` empty:

| Probe | Diagnostic anchor |
| --- | --- |
| unknown superseded Spec | `superseded Spec "9998-unknown" is unknown` |
| absent superseding Spec | `neither active nor archived` |
| self-supersession | `cannot supersede itself` |
| draft deliverer | status is `draft`; expected `active` |
| archived deliverer in active root | status is `archived`; expected `active` |
| malformed deliverer | `has malformed _prd.md` and the YAML error |
| statusless deliverer | `frontmatter has no status` |
| unmarked archived deliverer | archived status is `active`; expected `archived` |
| unknown flag | `flag provided but not defined: -unknown` |

An active deliverer then exited `0`, and fresh Git status named only the new
`_supersession.md`. A duplicate supersession exited `2`; the amendment hash
remained `d8c4cdc67ffc15bffd565776b881e749beecce79ec17d3efc2b15484f2d07473`.
A conventionally archived deliverer also exited `0` for a separate target.

The independent focused regression command
`GOCACHE=/private/tmp/roundfix-qa0151-rerun-go-cache go test -count=1
./internal/cli -run
'^(TestSupersede|TestSupersedeRejectsANonActiveDeliverer|TestSupersedeAcceptsAnActiveDeliverer|TestSupersedeAcceptsASupersessionArchivedDeliverer|TestSupersedeRejectsAnUnmarkedArchivedDeliverer|TestArchiveAcceptsARecordedSupersession)$'`
exited `0` and reported `ok roundfix/internal/cli`.

## Task Graph regression

In disposable Git repository
`/private/tmp/roundfix-qa0151-r4-final.LPmoYQ/repo`, a fresh rerun after the
row's input correction, archive of the current
Spec exited `2` before and after the public supersede command. Both runs named
the identical reason: `Task "task_04" is "pending"; archive requires every
Task to be "completed"`. Fresh Git status named only `_supersession.md`.

## Documentation and help

- `cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` passed
  before and after the checks.
- `./bin/roundfix skills check` exited `0` and reported the Roundfix skill set
  passed.
- `./bin/roundfix supersede --help` exited `0` and exposed the invocation,
  options, exit codes, one-file write, and no-Run/no-commit/no-push envelope.
- `./bin/roundfix archive --help` exited `0` and described the supersession
  proof only for a Spec without a Task Graph while preserving the move.
- Focused searches located the invocation, refusal contract, and unchanged
  archive-precondition statement in the canonical skill, mirror, and guide.

The final source-faithful command was rerun after the row's input declaration
was corrected to cover the full repository and exited `0`.

Two preliminary harness assertions failed before this final command: one
incorrectly expected refusal details in CLI help rather than in the authored
documentation, and one used a case-sensitive search for an uppercase sentence.
Direct inspection isolated both as harness mistakes; no product artifact was
changed, and the corrected source-faithful command passed.

## Authorization and commit scope

`git diff-tree --no-commit-id --name-only -r <commit>` enumerated both
authorization commits and every Task commit:

- authorization `2b73ac3f` introduced `_authorization.md`;
- authorization `03cfb8ee` added `internal/spec/archive.go` to the bounded
  paths before Task 02 consumed it;
- Task 01 `7d9e94d8` changed its Task file and ordinary CLI/spec source/tests;
- Task 02 `13358739` changed its Task file, ordinary source/tests, and bounded
  `internal/spec/archive.go`;
- Task 03 `60224664` changed its Task file, ordinary user-guide documentation,
  and the two bounded skill paths;
- Task 05 `0d1d6bec` changed its Task file, `CONTEXT.md`, and ordinary
  CLI/spec source/tests;
- Task 06 `05757c2e` changed its Task file and ordinary CLI source/tests.

The audit was rerun after the row's input declaration was corrected to cover
the full repository. Both authorization ancestry checks and the corrective
chronology check exited `0`; a focused count found exactly the three expected
governed Task paths. No Task changed a governed path outside the exact bounded set, no
authorization or prerequisite fix was folded into a Task commit, and neither
corrective Task changed a governed path. The skill mirror remained identical.
The final Git reads used `-c core.fsmonitor=false` and returned the complete
path lists with exit `0`.

## Promise coverage

After the row's input declaration was corrected to cover the full repository,
the current-tree strict Spec check exited `0` with no findings. Focused direct
searches confirmed that Task 04 cites PRD Goals 1-3, Core Features 1-4,
Success Metrics 1-3, TechSpec Testing Approach 1-6, and API Contracts 1-3.

## Pull Request equivalent evidence

The supervised QA prompt states that no Pull Request is open for
`feat/0151-supersession-the-archive-can-see`, Pull Request journeys are
environment-blocked, and the Run Worktree branch is never pushed and has no
Pull Request of its own. The gate did not query or mutate a Pull Request.

- Approval: none obtained because no Pull Request exists.
- Checks/status: the Daemon-supplied `make verify` result is pass, exit `0`.
- Unresolved review threads: none stands because no Pull Request exists.
- Merge-Ready acceptance: none exists yet.
- Review-artifact ancestry: no review artifact exists yet and none claims a
  head. The audited candidate head is
  `05757c2e1058d85609ade1d68b25943e7aa22588`.

Every unavailable control is explicit and backed by supervised equivalent
evidence rather than inferred as passing.
