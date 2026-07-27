# Live CLI flows

Build: `ef6eb44ad8951112b1c3641bb7fd21793b440f95`

Binary:
`bin/roundfix`, produced by the passing static gate.

Environment: macOS; real user-scoped Run Database; isolated local clone at
`/private/tmp/roundfix-qa-0037-live.qhCESz/repo`; fake `acpx` and GitHub CLI
executables supplied deterministic external boundaries. The scratch clone
started at `95a2c0d25b22db019507bee36cceb2218844a2f4`, stayed clean, and received no
commit or push.

## Force Stop with a registered Agent Session

Run `run_20260727T142906Z_a8d858e652e38cce` entered Task 01 Agent work through
the built detached Implement Command. The fake ACP Runtime blocked at the
public prompt boundary so the owner and registered Task Agent Session remained
active.

The built Force Stop Command exited zero and reported:

```text
Roundfix Run force-stopped
State: Stopped
Force Stop proved the recorded owner process exited, completed the Run as
Stopped, and released its Active Run locks.
Roundfix did not edit user files, commit, push, fetch, or resolve Review Source
threads.
```

The ACP command log records the active Task scope and the Force Stop cleanup:

```text
... sessions ensure --name roundfix-run_20260727T142906Z_a8d858e652e38cce-task_01
... prompt -s roundfix-run_20260727T142906Z_a8d858e652e38cce-task_01 -f -
... cancel -s roundfix-run_20260727T142906Z_a8d858e652e38cce-task_01
... sessions close roundfix-run_20260727T142906Z_a8d858e652e38cce-task_01
```

Independent confirmation:

- `roundfix runs list --state all` reported the Run as `Stopped`.
- `roundfix events ... --filter task-status,verification,outcome` contained
  one Task-start event and exactly one `Stopped` outcome.
- A permitted `ps` search found no process containing the Run ID after the
  report.
- `git worktree list` contained only the scratch user checkout; Force Stop
  reaped the Run Worktree and branch.
- `git status --short` remained empty.

## Lock release and idempotent replay

Run `run_20260727T142631Z_1528df7e84d509ea` was force-stopped through the same
live flow. A second detached Implement Run on the same repository and Spec,
`run_20260727T142726Z_5d351f580c89bf2d`, started successfully immediately
afterward. This independently confirms release of the first Active Run lock.
The second Run was also force-stopped and reaped.

Repeating Force Stop against the first already Stopped Run exited zero and
returned its stored `Stopped` report. A fresh Event Stream read still contained
exactly one outcome event at cursor 7.

Force Stop against the existing terminal
`run_20260727T135811Z_eb4119943fdae042` exited 2:

```text
terminal outcome conflict for Run
"run_20260727T135811Z_eb4119943fdae042": stored "Unresolved", requested
"Stopped"
```

A fresh Run listing and Event Stream read still reported `Unresolved` and the
same outcome events.

## Graceful Stop Request during in-flight work

Run `run_20260727T142819Z_61e165a06756e6e0` was stopped gracefully while its
Task Agent prompt was active. The command exited zero and reported:

```text
Roundfix Stop Request recorded
State: ResolvingWithAgent
Stop Request recorded; the Run stops after the current Work Item settles.
Roundfix recorded the Stop Request in the Run Database only.
```

An immediate `runs list --state active` still showed the Run in
`ResolvingWithAgent`, proving graceful stop did not force-complete the in-flight
Work Item. Force Stop then safely cleaned up the deliberately blocked fixture.

## Graceful Stop Request during Review Source wait

The built Watch Command started detached Run
`run_20260727T143452Z_54843a7fdd2b5ce7` against a deterministic fake CodeRabbit
source. The source reported an in-progress CodeRabbit check for the exact
scratch HEAD, leaving the Run `Active` in `WaitingForReview`.

The built Stop Command recorded a graceful Stop Request. Following the public
Event Stream with:

```text
roundfix events run_20260727T143452Z_54843a7fdd2b5ce7 --follow --filter outcome
```

returned:

```json
{"schema":"roundfix-events/v1","run_id":"run_20260727T143452Z_54843a7fdd2b5ce7","category":"outcome","outcome":"Stopped","summary":"Run reached Stopped."}
```

Independent confirmation:

- A fresh Run listing reported `Stopped` after 33 seconds total runtime.
- The Review Source log contains only the initial PR metadata, check-runs,
  commit-status, and reviews reads. It contains no later call or mutation after
  Stop Request observation.
- The console closed with `Stopped after 0 Round(s)` and states that no later
  verification, commit, push, fetch, or Review Source mutation ran.
- The console reports `Changed paths after Stop Request: none`.
- Scratch `git status --short` stayed empty and HEAD stayed
  `95a2c0d25b22db019507bee36cceb2218844a2f4`.

An earlier fake-source attempt ended `Failed` before Agent work because the QA
fixture returned `{}` instead of a reviews array. The built command reported
that parse failure as the primary reason with no Agent Session close warning.
The corrected fixture rerun above is the credited Review Source journey.

Verdict: pass.
