---
spec: 0360-the-build-cache-stays-within-a-bound-and-a-full-disk-is-named
status: active
created: 2026-11-06
surfaces: [backend, infra, docs]
type: fix
---

# The build cache stays within a bound, and a full disk is named

Minimal PRD for a bug fix, minted by `write-techspec`. Reported by oraculum on
2026-11-05 and by the maintainer's machine the same evening: the Mac mini's
system volume filled twice in one day, first stopping oraculum's Postgres in
the middle of a Run and then parking three Delivery Queue items as a generic
`delivery-error`, and both times the largest temporary directory was this
repository's Makefile build cache, `$TMPDIR/roundfix-go-build-<uid>`, at 25
GB. After this Spec the Makefile holds that cache to a Build Cache Bound
before it builds or tests, `roundfix doctor` and the Delivery Queue warn when
the cache's, the worktrees' or the temporary directory's filesystem drops
below a Free-Space Floor, and a park caused by a full filesystem says so,
with what to free, instead of `delivery-error`.

## Sources

Two reports of one incident, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-05-the-roundfix-build-cache-fills-the-mac-mini-internal-disk.md`](references/2026-11-05-the-roundfix-build-cache-fills-the-mac-mini-internal-disk.md),
   `type: fix`, reported by oraculum (ai-memory message
   `01a1ad41-0755-7203-afc7-2a9d34575cc6`, 2026-11-05 16:48 BRT), with the
   operator's note of 2026-11-06 01:58 BRT.
2. [`references/2026-11-06-the-repository-build-cache-stays-within-a-size-bound.md`](references/2026-11-06-the-repository-build-cache-stays-within-a-size-bound.md),
   `type: fix`, the maintainer's machine incident of 2026-11-05 19:07 BRT,
   open on `docs/backlog/` at the authoring base.

At the authoring base `docs/backlog/` holds one other open entry, the release
workflow and a transparency-log 409, which shares no context with this one,
and `docs/findings/` holds none.

## Prerequisites

None in `requires`. Authored against `6ceb7ff5` (v0.136.0 plus the backlog
commit). Spec 0359, authored in parallel, delivers first: it owns the
Delivery Queue's check wait in `internal/delivery` (`checkCandidate`) and the
`gh` checks reader in `internal/cli`. This Spec owns the Makefile, the
`internal/sharedcache` bound, the Doctor `disk-space:` line and the
Disk-Full Park. Files it may share with 0359, kept to narrow edits named by
function: `internal/delivery/engine.go` (`runCandidate`, `gateCandidate`, the
`delivery-error` park in `Run`, `Retry`'s retry-limit check and
`EngineDependencies`), `internal/delivery/park_class.go` (`ClassifyPark`),
`internal/cli/deliver_workflow.go` (`stepExecutable`, `Gate`,
`commandDeliveryEngineDependencies`), the `deliver` guide and both copies of
the Roundfix `deliver` reference, the owned Roundfix skill version and its
record, `CONTEXT.md`, and the Behavior Surface Record.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the new blocker and warning names are strings of the existing
  Delivery Queue item record. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; every change reads local files and filesystem statistics.
  Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0370 (new) records the decision:
  "The Makefile bounds its build cache, and the queue warns of low disk and
  names a full one". It amends ADR-0327 (the Shared Build Cache trim, whose
  Considered Options rejected a Makefile trim) and ADR-0334 (the automatic
  trims, which deferred a free-space floor for Runs), whose amendment notices
  task_01 and task_03 write. Their trim rules, the thirty-second budget and the
  `shared_cache_trim` Run Event keep their records unchanged. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AZ batch 3 answers of 2026-11-06 quoted in the authorization record,
  which grant the Makefile when the design needs it; the Governed Paths are
  Makefile, `.roundfixrc.yml`, the setup manifest task_01 regenerates after
  the Makefile and the repository rules guide task_01 changes, and the
  Roundfix skill files tasks 02 and 03 change. No workflow, Go module, lint or
  formatter change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0360-the-build-cache-stays-within-a-bound-and-a-full-disk-is-named/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/setup.md`, `.roundfixrc.yml`, `Makefile`, `docs/agents/setup-context.json`, `docs/agents/specific-repository.md`, `skills/roundfix/SKILL.md`.

## Goals

- A day of Runs on this repository cannot grow the Makefile's build cache
  without bound: every top-level make that builds or tests first trims it
  above 10,240 MiB, and no trim ever removes an entry a running build may
  still read.
- The bound never fails a build, never repeats inside a sub-make, and never
  touches a `GOCACHE` the caller set unless asked.
- `roundfix doctor` warns, beside the worktree location and the declared
  Shared Build Caches, when the temporary directory's filesystem has less
  than 10,240 MiB free.
- The Delivery Queue warns, before an item's Implement Run and before its
  repository gate, when one of those filesystems is below the same floor,
  in `roundfix deliver status` and the Queue Owner's log, and still advances.
- An item whose item binary build or repository gate failed because a
  filesystem was full parks as `disk-full: <path>: <N> MiB free` in the
  environment Park Class, with a next action naming what to free, and its
  retry is not counted.

## Core Features

1. The Build Cache Bound: a `gocache-bound` make target, a prerequisite of
   every target that builds or tests, runs `cmd/gocache-bound`, which
   applies the Shared Build Cache trim rules (idle entries only, oldest
   first, at least two idle hours) to the cache the Makefile chose; the
   Project Config declaration of the same cache takes the same values.
2. The Free-Space Floor in Doctor: the `disk-space:` line also inspects the
   temporary directory.
3. The Free-Space Floor in the Delivery Queue: a `low-disk` warning segment
   before the Implement Run and the repository gate.
4. The Disk-Full Park: a `disk-full` blocker built from an `ENOSPC` error or
   a retained log that says `no space left on device`, an environment Park
   Class with its next action, and an uncounted retry.
5. The rule in writing: the `deliver` and `doctor` guides, the README's
   Development section with the machine options (another volume for
   `TMPDIR` or `GOCACHE`, a User Config declaration of Go's user cache), the
   Roundfix `deliver` and `setup` references, the repository rules guide,
   three glossary entries, ADR-0370 and the amendment notices of ADR-0327 and
   ADR-0334.

## Non-Goals / Out of Scope

- Clearing the cache (`go clean -cache`) or removing an entry used in the
  last two hours, by any command.
- Go's default user cache (`~/Library/Caches/go-build`), which Task
  Verification commands that call `go` directly fill: a machine can declare
  it in User Config `worktree.shared_caches`, and the guide says so.
- Moving `GOCACHE` or `TMPDIR` to another volume: a machine choice the guide
  names as an option.
- A bound for adopters' build caches: the Makefile is this repository's own.
  The Doctor and Delivery Queue floor and the Disk-Full Park reach adopters.
- A `run-unresolved` park whose Task Verification filled the disk: it keeps
  its blocker, and `roundfix runs causes` already names its `disk-full`
  signature. The Linux Lane's container failures keep theirs.
- Refusing to start a Run or a gate below the floor.
- The CI workflows, `go.mod`, `go.sum`, lint and formatter configuration.

## Success Metrics

1. Success Metric: with a Go stub, `make test` in a temporary directory with
   `TMPDIR` set and no `GOCACHE` invokes, before the test, exactly
   `run -buildvcs=false ./cmd/gocache-bound -dir <TMPDIR>/roundfix-go-build-<uid> -max-mb 10240 -idle 2h`;
   `make vet test` invokes it once, a sub-make (`MAKELEVEL=1`) and an
   explicit `GOCACHE` invoke it never, and a failing bound prints
   `gocache-bound: warning: the Build Cache Bound did not run (exit 1); the build continues`
   while `make` exits `0`. At `6ceb7ff5` no bound is invoked.
2. Success Metric: over a cache whose entries total 1,843,200 bytes in three
   600 KiB entries idle for ten hours, five hours and thirty minutes,
   `gocache-bound -max-mb 1 -idle 2h` removes the two idle entries oldest
   first and keeps the third, and with every entry younger than two hours
   it removes nothing and prints the over-bound warning, exiting `0` both
   times.
3. Success Metric: `roundfix doctor`'s `disk-space:` line names
   `temporary directory <path>` with `DR-DISK-LOW` when its filesystem has
   10,239 MiB free and adds nothing at 10,240 MiB; Doctor's exit code is the
   same either way.
4. Success Metric: a Delivery Queue item whose locations have 10,239 MiB free
   before its Implement Run is shown by `roundfix deliver status` with
   `Warning: <slug> low-disk: ...` and still reaches its Run; with the space
   back above the floor at the gate, the segment is gone.
5. Success Metric: an item whose `make build` log ends with
   `write /tmp/x/_pkg_.a: no space left on device` parks as
   `disk-full: /tmp/x/_pkg_.a: <N> MiB free` in the `environment` Park
   Class, where at `6ceb7ff5` it parks as `delivery-error: ... build item
   binary: verification command "make build" failed ...`; a repository
   gate that fails the same way parks `disk-full` instead of `gate-failed`;
   and a Delivery Retry from it leaves the retry count unchanged.
6. Success Metric: the guides, both copies of the Roundfix references, the
   repository rules guide, `CONTEXT.md` and the amended ADRs state the rule,
   and the documentation contract holds the `disk-full` next action in the
   `deliver` guide and both reference copies.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The oraculum report quoted in the adopted reference, written before this
  Spec existed: "o maior item temporário do disco interno é
  `$TMPDIR/roundfix-go-build-501` ... com 25 GiB" and "O volume está em 99%
  (3,3 GiB livres) e caiu de 7 GiB para 3,3 GiB em ~30 minutos", measured on
  the Mac mini on 2026-11-05.
- golang/go issue 31948, "cmd/go: concurrent build and cache clean is unsafe"
  (https://github.com/golang/go/issues/31948, opened 2019-05-09, read through
  Exa on 2026-11-06): "If you start a build with `go build`, and run `go
  cache -clean` while it is building, the build can fail, because
  intermediate artifacts get deleted." This is why the bound removes only
  idle entries and never clears the cache.
- golang/go issue 29561, "cmd/go: clean GOCACHE based on disk usage"
  (https://github.com/golang/go/issues/29561, read through Exa on
  2026-11-06), where the go command's author describes its trim: "we just
  scan the tree for files that haven't been used in 5 days and delete them";
  the go command has no size limit of its own.
- `docs/references/2026-10-24-test-suite-audit.md`, committed by Spec 0311
  before this Spec existed: the same cache "held 31 GB" on 2026-10-24.

## Unreachable Acceptance

- criterion: after release, on the operator's delivery machine, the
  Makefile cache stays near its 10,240 MiB bound across a day of queued
  Runs, and no Delivery Queue item parks `delivery-error` for a full disk.
  reason: it needs the released binary, the real Roundfix Home and the
  machine's real disk, which no Verification or QA row may read or fill.
  satisfied-by: the operator's queue records and `du` of the cache after the
  release that carries this Spec, recorded session evidence.

## Validation

The decisions were checked on 2026-11-06 against these sources.

- Characterization at `6ceb7ff5`:
  - the Makefile exports `GOCACHE` and runs no trim; `.roundfixrc.yml`
    declares the cache with `max_size_mb: 20480` and `idle: 4h`, trimmed only
    at Run starts, at gates and by `roundfix gc` (ADR-0327, ADR-0334);
  - `defaultDoctorDiskSpaceResults` (`internal/cli/doctor_disk_space.go`)
    inspects the worktree location and each declared Shared Build Cache
    against 10,240 MiB, and nothing else;
  - `stepExecutable` (`internal/cli/deliver_workflow.go`) returns `build item
    binary: verification command "make build" failed; diagnostics: <log>:
    exit status <n>`, which `Run` (`internal/delivery/engine.go`) parks as
    `delivery-error: ...` in the environment Park Class; the build output,
    `no space left on device` included, stays only in the log;
  - `Gate` returns a failed result without a blocker, parked `gate-failed`
    in the finding Park Class;
  - only `runtime-infrastructure` is exempt from the retry limit (`Retry` in
    `internal/delivery/engine.go`, `RetryDeliveryQueueItem` in
    `internal/store/delivery.go`).
- Measurement on this machine with a temporary `GOCACHE` in this Spec's
  scratch directory: one fresh clone's `go vet ./...` and `go test -run '^$'
  ./...` wrote 507,664 KiB in 4,293 files; a second clone at another path
  added 301,572 KiB more, the per-worktree growth ADR-0327 names.
- A prototype of task_01 on a scratch copy at `6ceb7ff5`: `make build` over
  an 800 MB cache with `GOCACHE_BOUND_MB=700` removed 171 idle entries and
  warned for the rest; with the bound under `go test`, the `verifyselect`,
  `skills`, `suiteguardcontract`, `baseline` and setup-manifest suites
  passed except `TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped`
  (`_techspec.md` → Declared breaks), and each Makefile sabotage of Testing
  Approach failed its named row.
- Secondbrain, read only: `wiki/index.md`, then `qmd query "go build cache
  GOCACHE size bound disk full no space left on device trim" --all --files
  --min-score 0.3`. Its best matches were this repository's mirrored Backlog
  Entry, ADR-0327, ADR-0334, the test-suite audit and the Spec 0321 record,
  and `wiki/concepts/nobreak-e-armazenamento-do-brain-macmini.md`, which
  records the 256 GB internal disk, the external `Brain-SSD` and the plan to
  put "cache Go via `GOCACHE`, worktrees do roundfix" there: the machine
  option the guide names, not this Spec's fix.
- Exa: golang/go issues 31948 and 29561 quoted above, issue 69879 (a
  declined proposal to make the five-day cutoff configurable, whose last
  comment reports concurrent agents filling the cache "several times a
  day"), and issue 81830, which lists "concurrent `clean` and `build` is
  unsafe" and "cache can be trimmed while still in use" among the cache's
  known problems.

## Glossary

- adds: **Build Cache Bound**
- adds: **Free-Space Floor**
- adds: **Disk-Full Park**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AZ batch 3 brief for Spec 0360 of 2026-11-06, the maintainer's answers below and the adopted reports.

The maintainer's answers (AskUserQuestion, 2026-11-06), verbatim:

- Question "O cache de build das Runs do Roundfix
  ($TMPDIR/roundfix-go-build-501) chegou a 26 GB e encheu o disco interno do
  mini duas vezes ontem. Às 15:58 BRT derrubou o Postgres do Oraculum, e às
  19:07 BRT parou a fila. Acabei de limpá-lo. Como evitar que volte?" →
  "Spec: limite + aviso (Recommended)" — option text: "O Makefile limpa ou
  poda o cache acima de um tamanho antes de compilar. O doctor e a fila
  avisam quando o disco do cache fica com pouco espaço, e um park por disco
  cheio diz isso em vez de 'delivery-error'."
- Rules: "Mesmas do AZ (Recommended)" — option text: "Sem supervisão, merge
  pela fila, conserto e retry quando a entrega parar. Release assim que cada
  Spec ficar pronta, com notas aos adotantes e resposta ao Fluxus e ao
  Oraculum."
- The maintainer's standing instruction of 2026-11-05: "Quero que faça os
  releases hoje assim que estiverem prontos. E se existe no baseline regra
  para um release por dia por algum motivo deve ser removida."

Taken from those answers, the reports and the characterization, with no
further interview:

- The Makefile prunes and never clears: "limpa ou poda" allows either, and a
  clear during a concurrent build can fail that build (golang/go issue
  31948), so the bound removes only entries idle for at least two hours.
- The bound is 10,240 MiB and the idle time two hours, the minimum the
  Shared Build Cache rules accept; the Project Config declaration of the same
  cache moves from 20,480 MiB and four hours to the same values.
- The bound runs once per top-level make, from `go run ./cmd/gocache-bound`,
  reusing `internal/sharedcache`; it fails open with a warning.
- The Free-Space Floor stays Doctor's 10,240 MiB and becomes one rule shared
  with the Delivery Queue; it warns and never fails or parks.
- The Disk-Full Park covers the item binary build and the repository gate,
  the queue's own builds; a Run's Task Verification keeps `run-unresolved`.
- A retry from `disk-full` is not counted, as one from
  `runtime-infrastructure` is not.
- ADR-0370 records the decision and amends ADR-0327 and ADR-0334.
  Three implementation Tasks in series and a QA gate.

## Open Questions

None.
