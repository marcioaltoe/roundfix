# Characterization of the delivery recovery defects (2026-11-23)

Measured by the author of the TechSpec on `fe247d4c` (v0.153.0), Go 1.26.5, macOS 27.2.1. Every measurement ran in temporary homes and temporary repositories with fake runners. Nothing reached the network, an Agent or another repository.

## Method

Each probe is a Go test file kept outside the repository and added to its package with `go test -overlay`, so the repository's tree did not change. The probes assert nothing: they build a state with the packages' own test fixtures, run the production code and log what it did. The two prototypes replace `internal/cli/deliver.go` through the same overlay.

```text
go test -buildvcs=false -count=1 -overlay <overlay.json> ./internal/delivery -run TestProbe0418 -v
go test -buildvcs=false -count=1 -overlay <overlay.json> ./internal/cli -run TestProbe0418 -v
go test -buildvcs=false -count=1 -timeout 20m -overlay <prototype overlay> ./internal/cli
go test -buildvcs=false -count=1 -timeout 20m -tags docscontract -overlay <prototype overlay> ./internal/docscontract
```

A statement below is "measured" when a probe printed it and "inferred" when it follows from reading the code or from the sources and no probe printed it.

## 1. How a Run ends Failed with its QA gate's Task settled failed

Probe: `roundfix implement` through the public dispatch, a Spec with `task_01` and a QA gate's Task, a fake Agent that completes `task_01` and writes a QA Report whose verdict is `fail`, and a fake committer that refuses the commit whose subject begins `docs: qa report for`.

Measured:

```text
implement exit=1 Run state=Failed
stdout:
Failed: 1 completed, 1 failed, 0 skipped, 0 pending.
Tokens: none reported by 2 prompt(s); cost not reported
stderr (one line of it):
roundfix: implement failed after Run start: probe: the commit hook refused the QA Report commit
Run Worktree kept: <path>
daemon.task events of the Run, in order:
{"batch":1,"phase":"started","status":"in_progress","task":"task_01"}
{"phase":"settled","status":"completed","task":"task_01"}
{"phase":"settled","reason":"QA verdict fail","status":"failed","task":"task_qa"}
QA Task file in the Run Worktree: status: failed
```

The path, read in `internal/daemon/task_engine.go` and `internal/cli/implement.go`: the QA step settles the gate's Task `failed` and journals it, and only then prepares and creates the QA Report commit. Every error after the settlement, a failed format command or a refused commit among them, leaves the Task cycle with an error. The Implement Command then completes the Run as Failed and exits `1`. A failing verdict whose QA Report commit succeeds ends the Run Unresolved, which is the case the glossary describes.

Inferred: that the Run quaestor reported, `run_20261122T232957Z_bfb76adf75230a25`, took this path. Its evidence fits it, a QA Task `settled failed` at 00:23:17Z and the outcome "Run reached Failed", and holds no line of the Run's error. The item worktree of that repository installs its dependencies with `bun install`, so a commit hook is plausible, and it is not proved.

Inferred: that on the Run Branch the QA gate's Task keeps the status it had when the Run started. The Task file travels in the QA Report commit (`TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport` pins the commit's paths), and that commit was refused. The probe ran with the fake Run Worktree of the package's fixtures, which has no Run Branch to read.

## 2. What the Queue Owner leaves when its Implement Run fails

Probe: the production Delivery Engine over a temporary Run Database, an item at `running` that records `run-earlier`, a fake runner that returns the executor's error, and a Turn Recorder that keeps every write.

Measured:

```text
item after the failed executor: stage=parked
  blocker="delivery-error: run Implement executor: roundfix implement failed with exit code 1"
  run_id="run-earlier"
turn record write 1: [{probe-spec running advancing}]
turn record write 2: [{probe-spec running turn-wait}]
```

Two facts follow.

- The item parks with the blocker quaestor reported, word for word, and keeps the Run it recorded before. The Run that just failed is not on the item.
- The last write of the Turn Record lists the item at `running` with `turn-wait`. The owner writes it when the step returns and before it parks the item, and writes nothing after the park.

The second fact dates the Turn Record quaestor found. `turn.json` of owner `delivery-85a8e4c6a0c5-63163` holds `"stage":"running","turn":"turn-wait"` and was written at 21:23 BRT of 2026-11-22. The Failed Run ended at 00:23:17Z, which is 21:23:17 BRT. That record is the last write of the owner that ran the Failed Run, not the trace of an owner the retry started.

Inferred: that owner then ended normally. With every item parked, an owner releases its record and exits `0` (`ReleaseIdleDeliveryQueueOwner`), which is why no owner was recorded when the retry ran.

## 3. What a retry leaves when its Queue Owner does not start

Probe: `roundfix deliver retry` through the public dispatch with the production Delivery Engine and the production carry-forward, on a temporary repository whose item is parked with quaestor's blocker and whose Spec has one terminal Implement Run on the Item Branch: `task_01` settled completed with its settlement commit and its passing verdict, the QA gate's Task not completed. The owner start is a fake that prints one line to standard error and returns `2`. The probe ran twice, with the Run completed as Failed and as Unresolved.

Measured, Run Failed:

```text
deliver retry exit=2
stdout:
Retried 0001-widget-flow: delivery-error: run Implement executor: roundfix implement failed with exit code 1 -> running
stderr:
probe: the Queue Owner did not start
queue after the retry: owner_pid=0 item stage=running blocker="" retry_count=1 run_id="<the Failed Run>"
task_01 on the Item Branch: status: pending
deliver status exit=0
stdout:
0001-widget-flow	running	-	<item worktree>
Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none
Usage: no prompts recorded across 1 Run(s); cost not reported
stderr:
```

Measured, Run Unresolved, the control: the same, except that standard output begins with `Carried forward from Run <run-id>: task_01` and `task_01` is `completed` on the Item Branch.

Four facts follow.

- The retry stores the stage and prints `Retried` before it starts the owner. When the start does not succeed the command exits `2`, and the item stays at `running` with no recorded owner and one retry counted.
- That state is the one quaestor measured: stage `running`, `retry_count` 1, the recorded Run still the Failed one, which the retry itself wrote, and a Usage line that counts one Run.
- `deliver status` of that state prints three lines and no `Owner:` line, as quaestor captured.
- A Failed Run is skipped in silence by the carry-forward, and the same Run completed as Unresolved is carried.

Inferred: that the owner start is what did not succeed on 2026-11-22. Three things support it. Quaestor reports that the retry answered `Retried` "sem o relatório de dono destacado", and the four lines of that report are printed only after the new owner's handshake. The sources name one owner directory for that night, the one of section 2. And no code path lets an owner that claimed the queue release its record while an item can advance: it releases when no item can, when `deliver stop` ends it, or when another command proves it gone. The sources hold neither the retry's standard error nor its exit code, so which step of the detached start failed is not established.

### The same through the built binary, with one real cause

Probe: the binary built from the unchanged tree, run in a temporary home against a temporary repository whose Run Database holds one parked item with quaestor's blocker and no Item Branch. The path `<home>/.roundfix/artifacts` is a regular file, so the detached start cannot create its Artifact Directory. Nothing else is faked.

Measured:

```text
$ roundfix deliver retry 0001-widget-flow
exit=2
stdout:
Retried 0001-widget-flow: delivery-error: run Implement executor: roundfix implement failed with exit code 1 -> queued
stderr:
Preflight failed

Reason:
  stat Artifact Directory "<home>/.roundfix/artifacts/<repository id>": stat <home>/.roundfix/artifacts/<repository id>: not a directory

No side effects:
  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

Usage:
  Run 'roundfix deliver --help' for usage.
queue after the retry: owner_pid=0 item stage=queued blocker="" retry_count=1
$ roundfix deliver status
exit=0
stdout:
0001-widget-flow	queued	-	-
Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none
Usage: no Runs recorded
stderr:
```

The item re-enters at `queued` here because it has no Item Branch; with one it re-enters at `running`, as above. This is one cause the QA gate can produce with the built binary and no fake. It is not offered as the cause of 2026-11-22.

## 4. What status shows of a running item

Measured in section 3: the row of a running item names its slug, stage, blocker and worktree. No line names a Run or a Task.

Measured in section 1: the Run Database holds, for each Task of a Run, a `daemon.task` event with the phase `started` and one with the phase `settled`. The QA step journals a settlement for the gate's Task and no start.

## 5. Existing tests the new status lines move

Prototype: `deliveryOwnerLine` prints `Owner: Delivery Queue has no recorded owner and <n> item(s) can still advance; run roundfix deliver resume` when an item can advance and no owner is recorded, and status prints `Run: <slug> no Active Run` for each item at `running`, after the `Turn:` lines.

Measured: `go test ./internal/cli` fails ten tests in nine files, and `go test -tags docscontract ./internal/docscontract` passes.

| Test | File | Line it gains |
| --- | --- | --- |
| `TestDeliverStatusPrintsNoOwnerLineWithoutARecordedOwner` | `internal/cli/deliver_owner_status_test.go` | `Owner:` |
| `TestDeliverStatusPrintsTheItemWorktree` | `internal/cli/deliver_test.go` | `Owner:` |
| `TestDeliverStatusPrintsAnItemWarning` | `internal/cli/deliver_revalidate_test.go` | `Owner:` |
| `TestDeliverStatusPrintsNoWarningLineWithoutAWarning` | `internal/cli/deliver_revalidate_test.go` | `Owner:` |
| `TestDeliverStatusPrintsNoRunsRecorded` | `internal/cli/deliver_usage_test.go` | `Owner:` |
| `TestDeliverStatusPrintsAParkLinePerParkedItem` | `internal/cli/deliver_park_status_test.go` | `Owner:` |
| `TestDeliverStatusShowsUndeclaredCorrectionSigns` | `internal/cli/deliver_corrections_status_test.go` | `Owner:` |
| `TestDeliverStatusNamesTheMachineCapacityHolderAndWaiters` | `internal/cli/machine_verification_status_test.go` | `Owner:` |
| `TestDeliverStatusReproducesSurfaceTranscriptQueueTurns` | `internal/cli/deliver_turn_test.go` | `Run:` |
| `TestDeliverStatusShowsARunSlotWaitTranscript` | `internal/cli/deliver_run_slot_wait_test.go` | `Owner:` and `Run:` |

Each of these builds a queue with an item that can advance and no recorded owner, or with an item at `running`, and compares the whole standard output.

## 6. Existing tests the retry's proof moves

Prototype: the retry buffers what the owner start prints, reads the queue again, and fails when an item can advance and no owner is proven alive.

Measured: `go test ./internal/cli` fails the ten tests of section 5 and one more, `TestDeliverRetryWithAnOpenPullRequestAndNoMergeEvidenceRetriesAsBefore` in `internal/cli/deliver_merged_outside_test.go`:

```text
open retry: exit=2 stdout="" stderr="Retry failed\n\nReason:\n  no Queue Owner holds the Delivery Queue\n" starts=1
```

That test runs the production engine, which stores an advanceable stage, and its fake owner start returns `0` without claiming the queue. Every other retry test keeps passing unchanged: each either hands the item to a recorded owner, or uses a fake engine that leaves the stored item parked, or ends with the item merged, so no item can advance when the command checks.

## 7. The governed set

Measured with `speccheck.GovernedPath` over every file the Build Order names, the ten test files above included. Governed: `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`, the three paths the authorization record grants. Not governed: every Go file under `internal/delivery`, `internal/cli`, `internal/store` and `internal/docscontract`, `CONTEXT.md`, the three ADR files, `.secondbrain-export`, `docs/user-guide/commands/deliver.md`, `skills/roundfix/references/deliver.md`, `skills/testdata/owned-skill-versions.json`, `docs/references/behavior-surfaces.json` and `docs/references/skill-coverage.json`. `docs/agents/setup-context.json` is governed and no step touches it.

## 8. What the documentation part of the update writes

Measured on a scratch clone of this branch: a new accepted ADR numbered 0450 with a true Knowledge Mark, a notice that begins `**Amended by [ADR-0450](...) (2026-11-23).**` in ADR-0279 and in ADR-0395, one commit, then

```text
go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
```

changed three files and no other: `.secondbrain-export` gained the ADR's path in the Standard Export List, and each of the two ADRs gained `amended_by: [ADR-0450]` in its header. `docs/agents/setup-context.json` and every guide stayed byte-identical. The same three lines may be written by hand, which the Secondbrain guide allows for the path of an ADR and the docs-layout guide for the header.

## 9. Consultations

Made on 2026-11-23, before the design was written. Each line says what the source changed or confirmed.

Secondbrain, read only. `wiki/index.md` was consulted first, its first 402 lines read and the rest searched by keyword; the Indexus tool `pesquisar_acervos` then answered with wiki results, and each page below was opened.

- `wiki/concepts/sinalizacao-segura-de-processos-por-pid.md`: a PID is compared with its start identity before it is trusted, and PID reuse is real in this repository's test suite. It confirmed the proof the retry uses, and it is why the tests use a PID that does not exist and the test's own PID, never a PID guessed alive.
- `wiki/concepts/convergencia-observavel-em-sistemas-operacionais.md`: a write that succeeded is not an operation that completed, and a gate that needs a state at its entry confirms it at its exit. It supports proving the owner after starting it.
- `wiki/concepts/erros-acionaveis-em-cli.md`: an error says what happened and names the next step. It shaped the `Item:` and `Next:` lines of the retry's failure report.
- `wiki/concepts/agendas-e-reexecucao-em-motores-de-workflow-durable.md`: a run that started and did not finish is a state worth reporting. It supports the `Run:` line.

Exa, published sources. The quoted passages are the ones the PRD's `## Acceptance evidence` holds; the others are summarized here and not quoted.

- systemd.service(5), https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html, on `Type=simple`, `Type=exec` and `Type=notify`: a start can be reported before the service runs, and a readiness notification is what closes that gap. It changed the retry's rule from "the owner was started" to "the owner is proven alive".
- s6, https://skarnet.org/software/s6/notifywhenup.html: a reliable start notification needs the daemon itself to send it. The owner's handshake is that notification, and the retry keeps waiting for it.
- psutil, https://psutil.io/api/ and https://psutil.io/faq/: a process is identified by its PID and its creation time, and that check is good in practice and not absolute. It is why the TechSpec says the proof is not absolute.
- GitHub Docs, "Re-running workflows and jobs", https://docs.github.com/en/actions/managing-workflow-runs-and-deployments/managing-workflow-runs/re-running-workflows-and-jobs, and Argo Workflows, `argo retry`, https://argo-workflows.readthedocs.io/en/latest/cli/argo_retry/: a retry keeps what succeeded and runs again what failed and what depends on it. It supports carrying the completed Tasks and running the QA gate again.
- Apache Airflow, the DAG run page, https://airflow.apache.org/docs/apache-airflow/stable/core-concepts/dag-run.html: keeping the state of a cleared task is an explicit choice. It matches a carry-forward that only an explicit Delivery Retry performs, under proofs.
- Command Line Interface Guidelines, https://clig.dev/, and git-status(1), https://git-scm.com/docs/git-status: output a script reads is an interface, and changes to it are additive. It supports two new prefixed lines beside lines that keep their bytes.

The ADR sweep and the web and wiki sweep went to the Research Subagent, whose conclusions the author checked against the files it named. The code reads, the probes, the governed measurement and the Indexus search ran in the author's session.
