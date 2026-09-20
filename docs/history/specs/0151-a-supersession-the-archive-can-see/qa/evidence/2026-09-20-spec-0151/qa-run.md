# QA run evidence — Spec 0151

Audited head: `0d1d6becd2bf6675ed0d6716bbc22da4920bd501`.
Auditing binary: `roundfix 0.15.0 (0d1d6bec, built 2026-09-20 00:05:48 -0300)`.
Disposable repositories lived under
`/private/tmp/roundfix-qa0151-20260920.d5PosW/`.

## Preconditions and static gate

- `./bin/roundfix spec check 0151-a-supersession-the-archive-can-see --strict`
  exited `0` with `No findings. Authored Verification commands were not
  executed.` It listed the same four non-equivalent skips retained in the QA
  report.
- The Daemon supplied `make verify` as passed with exit `0`; diagnostics were
  removed on success. The QA Agent did not rerun it.
- `_tasks.md` names `task_04` as the unique terminal QA node. Its dependencies
  `task_01`, `task_02`, `task_03`, and `task_05` are all `completed`.

## Real repository journey and outside evidence

Outside evidence came from the repository's real active Spec
`docs/specs/0128-release-planning-with-bare-stable-tags/` and archived Spec
`docs/history/specs/0147-a-planner-that-reads-both-tag-spellings/`, neither of
which this Spec authored. Copies were exercised through the public binary in a
disposable Git repository.

1. `roundfix archive 0128-release-planning-with-bare-stable-tags` exited `2`
   and named the absent `_tasks.md` plus the write-tasks remediation.
2. `roundfix supersede --spec 0128-release-planning-with-bare-stable-tags --by
   0147-a-planner-that-reads-both-tag-spellings --reason ...` exited `0`.
3. A fresh copy of the amended tree was taken after the first process exited.
4. The repeated archive command exited `0` and moved the Spec to
   `docs/history/specs/0128-release-planning-with-bare-stable-tags`.
5. `diff -qr` between the saved amended tree and the archived destination
   exited `0`, proving the move preserved every byte.

## Refusals and amendment preservation

Each refusal ran through the public binary in a committed disposable Git
repository. After the first seven refusals, `git status --porcelain` was empty.

| Probe | Exit | Diagnostic |
| --- | ---: | --- |
| unknown superseded Spec | 2 | `superseded Spec "9998-unknown" is unknown` |
| absent superseding Spec | 2 | `neither active nor archived` |
| self-supersession | 2 | `cannot supersede itself` |
| draft deliverer | 2 | status is `draft`; expected `active` |
| archived deliverer in active root | 2 | status is `archived`; expected `active` |
| malformed deliverer frontmatter | 2 | `has malformed _prd.md` plus the YAML error |
| deliverer with no status | 2 | `frontmatter has no status` |
| unknown flag | 2 | `flag provided but not defined: -unknown` |

An active deliverer then succeeded with exit `0`. A fresh Git read named only
the new `_supersession.md`; its frontmatter and body carried the deliverer,
date, and reason. A second supersession exited `2` with `already carries a
supersession`; the amendment SHA-256 remained
`028d6222ffddad7d6c5eaaa46a6068777c98d83cbc4619d7b16436d2d1be349d`.
A separate untouched target accepted the real archived Spec 0147 with exit
`0`, confirming both accepted deliverer states.

The focused independent regression command
`GOCACHE=/private/tmp/roundfix-qa0151-go-cache go test -count=1 ./internal/cli
-run '^(TestSupersede|TestSupersedeRejectsANonActiveDeliverer|TestSupersedeAcceptsAnActiveDeliverer|TestArchiveAcceptsARecordedSupersession)$'`
also exited `0`.

## Task Graph regression

In a separate committed disposable repository, archive of the copied current
Spec exited `2` with `Task "task_04" is "pending"` before supersession. The
public supersede command then exited `0`; archive again exited `2` with the
identical Task diagnostic. A fresh Git read named only `_supersession.md`, and
`git diff --exit-code` excluding that file exited `0`.

## Documentation and help

- `cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` exited `0`.
- `./bin/roundfix skills check` exited `0` and reported the Roundfix skill set
  passed.
- `./bin/roundfix supersede --help` exited `0` and exposed the invocation,
  options, exit codes, one-file write, and no-Run/no-commit/no-push envelope.
- `./bin/roundfix archive --help` exited `0` and described the supersession
  proof only for a Spec without a Task Graph while preserving the move.
- Focused `rg` reads located the command and archive/refusal contract in the
  canonical skill, mirror, and user guide.

## Authorization and commit scope

`git diff-tree --no-commit-id --name-only -r <commit>` was run for both
authorization commits and every Task-related commit:

- authorization `2b73ac3f` introduced `_authorization.md`;
- authorization widening `03cfb8ee` added the historically governed
  `internal/spec/archive.go` path before consuming implementation;
- Task 01 `7d9e94d8` changed its Task file and ordinary CLI/spec source/tests;
- Task 02 `13358739` changed its Task file, ordinary source/tests, and the
  expressly bounded `internal/spec/archive.go`;
- Task 03 `60224664` changed its Task file, ordinary user-guide documentation,
  and the two expressly bounded skill paths;
- corrective Task authoring `86984d4f` changed only Spec artifacts;
- Task 05 `0d1d6bec` changed its Task file, `CONTEXT.md`, and ordinary
  CLI/spec source/tests.

Both authorization commits are ancestors of the current head. No Task changed
a governed path outside the exact bounded set; the authorization and widening
were separate prior commits; no consequent fix or derived pin was folded into
a Task commit; and the canonical and mirror skill files remain byte-identical.

## Pull Request equivalent evidence

The supervised QA prompt states that no Pull Request is open for
`feat/0151-supersession-the-archive-can-see`, that Pull Request journeys are
environment-blocked, and that the Run Worktree branch is never pushed and has
no Pull Request of its own. The gate did not query or mutate a Pull Request.

- Approval: none obtained because no Pull Request exists.
- Checks/status: the Daemon-supplied `make verify` result is pass, exit `0`.
- Unresolved review threads: none stands because no Pull Request exists.
- Merge-Ready acceptance: none exists yet; it follows this terminal gate.
- Review-artifact ancestry: no review artifact exists yet and none claims a
  head. The audited candidate head is `0d1d6becd2bf6675ed0d6716bbc22da4920bd501`.

Every unavailable control is explicit rather than inferred as passing.
