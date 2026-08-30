# QA evidence — Spec 0116

Audited HEAD: `9b4292352bbefff60eea987734f1bdc976e94b3f`.

## Auditing Binary

- Installed comparison binary: `roundfix 0.9.0 (0d2619c7, built 2026-08-30
  10:20:10 -0300)`. Commit `0d2619c7` predates the audited HEAD.
- Application under test after the static gate rebuilt it: `roundfix 0.9.0
  (9b429235-dirty, built 2026-08-30 15:19:11 -0300)`.
- The `-dirty` suffix is caused by this Daemon-owned untracked QA report, so
  the report records staleness as unknown rather than silently calling an
  unresolvable stamp current.

## Strict precondition and clean non-probing verdict

Command:

```text
rtk ./bin/roundfix spec check 0116-a-verdict-that-states-its-own-scope --strict
```

Exit: `0`.

```text
Spec 0116-a-verdict-that-states-its-own-scope
No findings. Authored Verification commands were not executed.
```

The replaced trailing text `Verification: not run (use --run-verification).`
did not appear.

## Authoring-skill probing journeys

The canonical guidance names these exact commands:

```text
roundfix spec check <slug> --stage prd --run-verification
roundfix spec check <slug> --stage techspec --run-verification
roundfix spec check <slug> --run-verification
```

Each command was run through `bin/roundfix` for Spec 0116. Each exited `1`
because it reached the prober on an intentionally out-of-order, completed
graph, and each began:

```text
Spec 0116-a-verdict-that-states-its-own-scope
No findings. Authored Verification commands executed: 8.
Verification tree: HEAD
```

All three runs classified the completed implementation Tasks `task_01` through
`task_05`, `task_07`, and `task_08` as `vacuous — ... (exited zero before
work)`. They classified pending `task_06` as `honest — ... (exited non-zero
before work)`. The first sandboxed full-sweep attempt was environment-blocked
because the disposable checkout could not write Git worktree metadata outside
the sandbox; the permitted run completed and produced the observations above.

The QA-gate guidance instead names `roundfix spec check <slug> --strict` in its
command and both refusal examples. Beside it, the guide states that a probe
asks whether a command already passes before work exists and would therefore
call every completed Task vacuous at the terminal gate.

Canonical and generated copies of `write-prd`, `write-techspec`, `write-tasks`,
and `qa-gate` compared byte-identical with `cmp`.

## Clean probing persistence probe

A disposable Git fixture used the public binary with two authored commands,
`test -f task-output-one.txt` and `test -f task-output-two.txt`. Neither file
existed at HEAD. This command ran twice:

```text
bin/roundfix spec check clean --strict --run-verification
```

Both runs exited `0` and reported:

```text
Spec clean
No findings. Authored Verification commands executed: 2.
Verification tree: HEAD
- task_01: honest — "test -f task-output-one.txt" (exited non-zero before work)
- task_01: honest — "test -f task-output-two.txt" (exited non-zero before work)
```

The fixture also printed expected skips for repository artifacts it did not
copy; no finding was present and the two-command observable persisted on the
fresh second run.

## Finding-shape probe

The public binary ran against the repository's `tooling-unauthorized` fixture.
It exited `1` and retained all public diagnostic fields:

```text
[error] SC-TOOLING-UNAUTHORIZED: docs/specs/tooling-unauthorized/_prd.md cites docs/workflow/authorizations/other-spec.md, but that record does not name Spec tooling-unauthorized
  at docs/specs/tooling-unauthorized/_prd.md:12
  at docs/workflow/authorizations/other-spec.md:1
  fix: Add Spec tooling-unauthorized to docs/workflow/authorizations/other-spec.md or cite the authorization record that already names it.
```

## Static and focused gates

- `rtk make verify` initially reached every package but two unrelated
  force-stop integration tests failed because the sandbox denied process-table
  reads with `operation not permitted`.
- The permitted full-access `rtk make verify` rerun exited `0`. It passed every
  package, the focused skill tests, `roundfix skills check`, and rebuilt
  `bin/roundfix`.
- The focused auditor, staleness, verdict-line, CLI, report-reader, refusal
  writer, and ordinary mechanical-report set passed `63` tests in `5` packages.
- The unchanged Daemon pre-work classifications passed `4` focused tests:
  vacuous refusal, ordinary failing path, unknown verdict, and repository-gate
  precondition followed by vacuous refusal.
- `rtk go test -count=1 -tags docscontract ./internal/docscontract` passed `47`
  tests.
- `rtk make skills-sync-check` exited `0`.

## Tooling authorization and chronology

`git diff-tree --no-commit-id --name-only -r` established:

- `4ff24609`: only the authorization record.
- `5b965b6e`: the four canonical authoring skills, their four generated copies,
  and `task_05.md`.
- `66e613e4`: canonical/generated `qa-gate` and `task_07.md`.
- `9b429235`: canonical/generated `qa-gate`, `task_08.md`, and the ordinary
  production/report expectation pair `internal/daemon/task_engine.go` and
  `internal/daemon/task_engine_test.go`.

`git merge-base --is-ancestor` exited `0` for authorization → `task_05` →
corrective `task_07` → corrective `task_08`. No authorization, prerequisite,
or corrective change is folded into the commit it precedes or repairs. Every
governed path is in the grant's exact eight-path list; each Task commit also
touches only its own Task file and ordinary production/test files where
applicable.

## Glossary and scope

The docs contract passed with the emitted spellings `AuditingBinary`,
`auditing_binary`, and `auditor_staleness`. `CONTEXT.md` defines `Auditing
Binary`, updates `QA Report`, and records the avoided ambiguous terms. The full
code diff changes no worktree/collision package or Verification prober file.
The only `internal/daemon/task_engine.go` production hunk is in ordinary QA
Report materialization; the Daemon's pre-work probe functions are unchanged.

## Outside evidence

The adopted reference traces to the Roundfix Finding added by commit
`6a27b556` on 2026-08-14, twelve days before this Spec. Git history records its
rename into this Spec's `references/` directory. The reference identifies the
origin as three freshly decomposed Specs in the separate `fluxus` repository
and records eight authored commands: three already-existing suite/file checks
and five repository verification gates. The same document's eleven-command
Roundfix measurement is corroboration only and was not used as the external
origin.
