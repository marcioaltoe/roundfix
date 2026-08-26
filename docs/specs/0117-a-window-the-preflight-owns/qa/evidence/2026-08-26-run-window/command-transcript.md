# Command transcript — Run Window QA

Build under test: `79f3a2eb957d14874c1636174ea664ab703c2661`.

All live CLI flows used the rebuilt
`/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260826T194844Z_4300b88f97ac6163/bin/roundfix`
binary, isolated homes under `/tmp/roundfix-qa-0117.OSwOS6`, and `TZ=Etc/GMT-3`
where the night-start scenario required a local time after 23:00.

## Authoring and static gates

```text
$ rtk ./bin/roundfix spec check 0117-a-window-the-preflight-owns --strict
Spec 0117-a-window-the-preflight-owns
No findings.
Skipped:
  SC-REF-UNRESOLVED: missing docs/specs/0117-a-window-the-preflight-owns/references/_index.md
Verification: not run (use --run-verification).
exit 0
```

```text
$ rtk make verify
all reported packages passed except internal/cli
TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion: operation not permitted while enumerating the process table
TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner: operation not permitted while enumerating the process table
exit 2 (make test target failure)
```

The same focused tests failed for the same process-table denial on the clean,
unchanged Spec target checkout at `730b26ab7881907e8dc3cb0e90114b8c0f781519`:

```text
$ rtk env GOCACHE=/tmp/roundfix-qa-base-cache go test -count=1 ./internal/cli -run 'TestRunForceStop(OwnerProcessIntegrationProvesExitBeforeStoreCompletion|LegacyRunWithoutOwnerIdentityStillStopsOwner)$'
read process table for non-session owner: operation not permitted
FAIL
```

The failure is environment-caused and unrelated to every Run Window row.

After the report and this evidence file were written:

```text
$ rtk make verify-docs
ok roundfix/internal/docscontract
ok roundfix/internal/baseline
ok roundfix/skills
Spec 0117-a-window-the-preflight-owns: No findings.
exit 0
```

## Night-start, restart, repository scope, and set semantics

```text
$ TZ=Etc/GMT-3 roundfix window set 07:00
Run Window set for <run-worktree>.
Cutoff: 2026-08-27 07:00 +03
exit 0

$ (from <run-worktree>/internal, new process) TZ=Etc/GMT-3 roundfix window show
Run Window for <run-worktree>
Cutoff: 2026-08-27 07:00 +03
Current time: 2026-08-26 23:45 +03
Remaining: 7h14m14s
The Run Window bounds when a Run may start; budget.max_run_duration bounds how long one may run.
exit 0
```

The first persisted database row was:

```text
<run-worktree>|1787803200|1787777135
```

A repeat set requested `08:00` without force and preserved those exact values:

```text
Run Window already set for <run-worktree>; unchanged without --force.
Cutoff: 2026-08-27 07:00 +03
1787803200|1787777135
```

The forced set replaced the row, and a past forced set exited `2` without
changing the replacement:

```text
Run Window replaced for <run-worktree>.
Cutoff: 2026-08-27 08:00 +03
1787806800|1787777168

Preflight failed
Reason:
  cutoff 2026-08-26 22:00 +03 must be in the future; the current time is 2026-08-26 23:46 +03
exit 2

1787806800|1787777168
```

`window clear` followed by a new-process `window show` reported no window.
Setting a separate window from `/Users/marcio/dev/roundfix` with the same Run
Database did not affect the Run Worktree's absent state, proving Git-root
scope through the public command.

## Closed-window refusal

The public command set `2026-08-26T23:52` while it was future. After the cutoff:

```text
$ roundfix implement --spec 0117-a-window-the-preflight-owns --agent-command "codex-acp --stdio" --no-input --no-agent-console
Preflight failed

Reason:
  the Run Window for <isolated-repository> closed at 2026-08-26 23:52 +03; the time is 2026-08-26 23:52 +03

No side effects:
  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

Next action:
  move the window with `roundfix window set <HH:MM> --force`, or remove it with `roundfix window clear`
exit 2
```

Independent public and database reads after the refusal:

```text
$ roundfix runs list --state all --limit 0
No Runs found.

$ sqlite3 <isolated-home>/.roundfix/roundfix.db 'SELECT COUNT(*) FROM runs;'
0

$ find <isolated-home>/.roundfix/worktrees -maxdepth 3 -type d
<no output>
```

## Crossing and start-only control

The isolated scripted ACP adapter is a parity deviation used only to hold the
Run across the wall-clock boundary and then produce a named runtime failure.
It does not implement or inspect Run Window behavior.

```text
Run Window: closes 2026-08-26 23:55 +03, in 102s; max_run_duration is 2h, so this Run may run past it.
Implement Run: run_20260826T205318Z_8513a3442369b65c
```

The scripted adapter entered its prompt before the cutoff. It was released
after 23:55 and failed intentionally. The public reads then showed:

```text
$ roundfix runs list --state all --limit 0
run_20260826T205318Z_8513a3442369b65c  Failed  implement  spec:0117-a-window-the-preflight-owns  codex-custom  2026-08-26T20:53:18Z  1m  ma/qa-live

$ roundfix events run_20260826T205318Z_8513a3442369b65c --filter outcome
{"schema":"roundfix-events/v1","run_id":"run_20260826T205318Z_8513a3442369b65c","category":"outcome","time":"2026-08-26T20:55:11.586685Z","cursor":13,"outcome":"Failed","summary":"Run reached Failed.","reason":"The Run failed before it could complete.","next_action":"Inspect the diagnostics, correct the failure, and start another Run."}

$ TZ=Etc/GMT-3 date '+%Y-%m-%dT%H:%M:%S %Z'
2026-08-26T23:55:27 +03

$ git rev-parse HEAD
79f3a2eb957d14874c1636174ea664ab703c2661
```

The window admitted the Run before 23:55 and did not stop it. The Run reached
its own terminal outcome after 23:55 from the injected runtime cause. No commit
was created in the isolated checkout.

## Schema migration

Before the public `window show`, a copied historical database reported:

```text
PRAGMA user_version = 12
runs = 1
run_windows table count = 0
```

The first `window show` returned the normal absent-window state. Afterwards:

```text
PRAGMA user_version = 13
runs = 1
run_windows table count = 1
run_windows rows = 0
```

A second `window show` returned the same absent state and retained version 13,
one pre-existing Run, and zero windows.

## Focused checks and docs

```text
$ rtk env GOCACHE=/tmp/roundfix-qa-focused-cache go test -count=1 ./internal/store -run 'Test(OpenMigratesV12RunDatabaseAddingRunWindows|RunWindowPersistsAndClearsByGitRoot|SetRunWindowPreservesExistingWindowWithoutReplace)$'
ok roundfix/internal/store 1.227s

$ rtk env GOCACHE=/tmp/roundfix-qa-focused-cache go test -count=1 ./internal/cli -run '^Test(Window|ImplementRunWindow)'
ok roundfix/internal/cli 1.976s

$ rtk env GOCACHE=/tmp/roundfix-qa-focused-cache go test -count=1 ./internal/watch
ok roundfix/internal/watch 0.374s

$ rtk env GOCACHE=/tmp/roundfix-qa-docs-cache go test -count=1 -tags docscontract ./internal/docscontract -run '^TestCheckActiveCorpusHasNoErrors$'
ok roundfix/internal/docscontract 0.512s

$ rtk env GOCACHE=/tmp/roundfix-qa-docs-cache go test -count=1 ./internal/config -run 'Test(DefaultConfigYAML|RenderedConfig)'
ok roundfix/internal/config 0.356s
```

`roundfix window --help`, the public command guide, the configuration guide,
the rendered configuration, and `CONTEXT.md` all state that the Run Window
bounds when an Implement Run may start while `budget.max_run_duration` bounds
how long one may run.

## Outside evidence and scope

The outside-evidence source predates the Spec implementation and is present at
`origin/main` commit `2aa4511e80b7e90e051abb98efa7ca1ee90bc36c`:
`docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md`. It records the
pre-Spec session-cutoff script, the failure mode in which a skipped caller-side
check lets Runs open all night, the `window.sh` header's naive-comparison bug,
and the requirement that Preflight own the bound.

`git diff 730b26ab..HEAD -- internal/watch skills` was empty. The five Task
commits changed only source, tests, fixtures, Task files, glossary, and user
guides; none changed a governed tooling path. Run Window reads occur only in
the `window` and `implement` command paths. The delivered feature adds no
timer, scheduler, supervisor-loop skill, or window-driven Stop Request.
