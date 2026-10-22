---
spec: 0300-reconcile-releases-a-stale-run-owner-record
status: active
created: 2026-10-22
surfaces: [backend, cli, docs]
---

# Reconcile releases a stale Run owner record

Minimal PRD for a bug fix, minted by `write-techspec`. A terminal Run keeps the
PID and start identity of the process that owned it. When that process has
exited and the system has given its PID to another process, `roundfix
reconcile` reads a different start identity, refuses to prove the owner, and
lists the Run as a preserved process candidate on every invocation, with no
command that ever ends it. After this Spec the Reconcile Command reports such
a record as a Stale Run Owner Record and `--apply` releases it without sending any
signal, while an identity it cannot read stays preserved as today. The
TechSpec holds the design.

## Sources

One open Backlog Entry, adopted by this Spec's planning change into
`references/` (see `references/_index.md`):

1. [`references/2026-10-21-reconcile-preserves-forever-a-process-candidate-whose-owner-pid-now-belongs-to-another-process.md`](references/2026-10-21-reconcile-preserves-forever-a-process-candidate-whose-owner-pid-now-belongs-to-another-process.md),
   type `fix`, written on 2026-10-21 from the maintainer's `roundfix
   reconcile` output (roundfix 0.101.0): three `kind=process` preserved
   candidates refused with `owner process identity cannot be proven: live
   process start identity does not match the recorded owner identity`. On
   2026-10-22 at about 00:30 BRT, with 0.106.0, the same command reported
   `preserved=4`: Runs `run_20261014T174911Z_8d771e92d24f52c5` (PID 13183),
   `run_20261012T110224Z_380dd005763aa861` (57765),
   `run_20261001T000219Z_0e911d72b0d72d2f` (74439) and
   `run_20260929T023407Z_5ea4db7269de3579` (25814).

No other open Backlog Entry or unresolved Finding shares this context: at the
authoring base this entry was the only file in `docs/backlog/`, and
`docs/findings/` holds no entry.

## Prerequisites

None in `requires`. This Spec is authored against `34b4e246` (v0.107.0), with
no Spec active, and delivers after the cycle AO Specs 0297, 0298 and 0299.
None of them changes `internal/store/process.go`, the reconcile report or the
reconcile documents; they raise the Roundfix skill's version, so this Spec's
documentation Task raises it over whatever the tree holds when it starts and
rebases over their bumps.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the release clears two existing columns of an existing Run row
  and journals an existing event kind. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the change reads the local process table
  and writes the local Run Database, and every test uses temporary homes,
  repositories and fake process tables. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0313: "A proven owner identity
  mismatch releases a terminal Run's owner record without a signal", which
  builds on the Proven Signal rule and on proof-based terminal reconciliation
  that the glossary entries Proven Signal and Reconcile Command
  carry. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AO approval of 2026-10-22 through AskUserQuestion, authorization "the same
  rules as cycle AN", which carries the cycle AJ grant "Mexer em Governed
  Paths declarados" ("'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica)."). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0300-reconcile-releases-a-stale-run-owner-record/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/reconcile.md`,
  `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint,
  formatter, Baseline module, Project Config or Run Database schema change.

## Goals

- A terminal Run whose recorded owner PID now presents a different start
  identity is reported once as a Stale Run Owner Record, not as a preserved
  candidate.
- `roundfix reconcile --apply` releases that record and sends no signal to
  any process; a later reconcile no longer reports it.
- An owner identity that cannot be read or compared, a PID that cannot be
  proven for another reason, and a Run recorded without an owner identity
  stay preserved exactly as today.
- Force Stop keeps refusing to signal a PID whose identity does not match.
- The `reconcile` user guide, the Roundfix `reconcile` reference and the
  glossary describe the Run owner record.

## Core Features

1. A proven mismatch is told apart. The owner proof returns a distinct error
   when the live start identity was read, is comparable and differs from the
   recorded one; read failures and uncomparable tokens stay unreadable, and an
   invalid or own PID stays unproven without a mismatch. The mismatch keeps
   today's message bytes and still counts as an unproven owner for every
   caller that refuses one.
2. The Run owner record candidate. For a terminal Run with a proven mismatch,
   the report lists one `runOwnerRecord` entry in `processCandidates` with the
   owner PID, an empty process list, the proof and the action, prints it as
   `Run owner record candidate: Run <run-id> owner PID <pid>`, and leaves
   `preservedCandidates` without it. The tree of the process now holding the
   PID is never enumerated.
3. The signal-free release. `--apply` clears the Run's recorded owner PID and
   identity in one compare-and-set write on the inspected values, only while
   the Run is terminal, journals a Run Event holding the released PID,
   identity and reason, and counts the release as `runOwnerRecordsReleased`
   (`run-owner-records-released` in text). It never asks the process controller
   to terminate anything.
4. Documentation. The `reconcile` guide and the Roundfix `reconcile`
   reference describe the Run owner record and its release; the glossary defines
   **Stale Run Owner Record** and revises **Reconcile Command**.

## Non-Goals / Out of Scope

- Changing Force Stop, the Proven Signal operations, the Delivery Queue owner
  proof or its `deliver status` line.
- Changing `doctor`'s residue check; it stops reading a released Run's lineage
  because the Run records no owner any more.
- Deleting a Run, changing its outcome or timestamps, or changing the Run
  Database schema.
- Releasing an Active Run's owner record, or one whose identity cannot be
  read.

## Success Metrics

1. Success Metric: in a temporary Roundfix Home, a terminal Run whose owner
   PID the fake process table answers with a proven mismatch is listed by
   `roundfix reconcile` as exactly one `runOwnerRecord` process candidate with
   the action `would release the Run owner record with --apply` and no preserved
   candidate, in text and JSON, and the Run Database bytes are unchanged.
2. Success Metric: `roundfix reconcile --apply` releases it: the Run's owner
   PID and identity are cleared, its outcome is unchanged, `runOwnerRecordsReleased`
   is 1, the process table was asked to terminate nothing, and a second
   `--apply` reports no candidate and no preserved entry for the Run.
3. Success Metric: an identity read failure, an uncomparable token, an
   unproven owner without a mismatch, a platform without process control and a
   Run recorded without an owner identity each stay in `preservedCandidates`
   with today's refusal reason through a dry run and `--apply`, and the record
   is kept.
4. Success Metric: with a fake process table, the owner proof reports a
   comparable different identity as the mismatch error, also unproven and not
   unreadable, with today's message bytes; the own PID and PID 0 are unproven
   without a mismatch; read failures and foreign tokens are unreadable without
   a mismatch; and Force Stop's termination over a mismatch returns the
   refusal and signals nothing.
5. Success Metric: a live process started by the test, holding the recorded
   PID under its own start identity, is classified as the mismatch by the
   real process table and is still running afterwards.
6. Success Metric: the release writes nothing for a different PID or identity
   or for an Active Run, writes exactly one journal event when it releases,
   and writes nothing when called again.
7. Success Metric: after release, `roundfix reconcile` on the maintainer's
   machine lists the four Runs named in Sources as Run owner record candidates,
   `roundfix reconcile --apply` releases them, and a later `roundfix
   reconcile` reports none of them.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-22.

- The maintainer's own `roundfix reconcile` output, recorded on 2026-10-21 in
  the adopted Backlog Entry before this Spec existed: the refusal reason
  `owner process identity cannot be proven: live process start identity does
  not match the recorded owner identity`, which the mismatch error must keep
  byte for byte.
- POSIX.1-2024, Base Definitions 4.17 "Process ID Reuse"
  (https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap04.html):
  "A process ID shall not be reused by the system until the process lifetime
  ends. In addition, if there exists a process group whose process group ID
  is equal to that process ID, the process ID shall not be reused by the
  system until the process group lifetime ends."
- proc_pid_stat(5) (https://man7.org/linux/man-pages/man5/proc_pid_stat.5.html):
  "(22) starttime %llu The time the process started after system boot", a
  field the page does not mark `[PT]`, the marker of fields shown as 0 when
  the ptrace access check denies them.

## Unreachable Acceptance

- criterion: Success Metric 7, after release, `roundfix reconcile` on the
  maintainer's machine lists the four Runs named in Sources as Run owner record
  candidates, `roundfix reconcile --apply` releases them, and a later
  `roundfix reconcile` reports none of them.
  reason: only the maintainer's real Roundfix Home holds those Runs, and no
  Verification or QA row may read or write the real `~/.roundfix` or run
  `roundfix reconcile --apply` outside a temporary Roundfix Home.
  satisfied-by: the operator's `roundfix reconcile`, `roundfix reconcile
  --apply` and second `roundfix reconcile` in this repository's checkout after
  this Spec's release.

## Validation

The decisions were checked on 2026-10-22 against these sources.

- Characterization at `34b4e246` with two `go test -overlay` probes that
  wrote nothing to the repository. With a fake process table, `InspectTree`
  returned `ErrOwnerProcessIdentityUnproven` both for a comparable different
  identity and for the own PID and PID 0, and `ErrOwnerIdentityUnreadable` for
  a read failure and a foreign token; it enumerated and signalled nothing in
  every case. Through the command, with a live child process holding the
  recorded PID under its own identity and the real process table, `roundfix
  reconcile`, `--format=json` and `--apply --format=json` each exited 0 and
  listed the Run as a preserved `process` candidate with the Backlog Entry's
  refusal reason (`preserved=1`); the Run Database bytes were unchanged, the
  owner PID was still recorded after `--apply`, and the child was alive. The
  TechSpec records the code it read.
- Secondbrain, read only. Read `wiki/index.md`; ran `qmd query "PID reuse
  process start identity reconcile stale owner record" --all --files
  --min-score 0.3`, which returned the Roundfix mirror's decision that process
  residue is a readiness fact, this Backlog Entry's mirror and the concept
  `wiki/concepts/sinalizacao-segura-de-processos-por-pid.md`, which records
  that a PID must be paired with its start identity and that a different
  start time means the PID was reused ("se o processo atual tem start time
  diferente do gravado, o PID foi reutilizado — não sinalize"). They confirmed
  the rule this Spec applies; no prior decision on releasing the record
  exists. No Inbox Entry was written, because the briefing forbids writing
  there.
- Exa: POSIX.1-2024 4.17 and proc_pid_stat(5) above. Challenging: a
  permission denial on Linux could have shown a start time of 0 and faked a
  mismatch, but `starttime` is not among the `[PT]` fields; and the recorded
  owner's descendants could outlive it, but POSIX keeps the PID unused while
  its process group lives, so a reused PID means that group had ended.

## Glossary

- adds: **Stale Run Owner Record**
- changes: **Reconcile Command**

## Decisions

- 2026-10-22 (about 00:55 BRT): the cycle AO scope, all four Specs in the
  order "Z: portão e corretiva", "Espera de checks", "Retomar item
  arquivado", "Y: reconcile PID", one minor each, and the authorization "the
  same rules as cycle AN"; this Spec's design is the Backlog Entry's
  Expected, kept by the maintainer.
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").

Autonomous route: no person present; decisions taken from the operator's cycle AO briefing for Spec 0300 of 2026-10-22 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- "Releasing the stale process record" clears the Run's recorded owner PID
  and identity; the Run, its outcome and its timestamps stay, and a Run Event
  keeps what was cleared.
- The record is reported inside `processCandidates` with the new `kind`
  `runOwnerRecord`, so the existing envelope and counts keep their meaning.
- The own-PID case stays unproven without a mismatch: the proof refuses it
  before reading any identity, as today, and the next invocation runs under
  another PID.
- `doctor` is left unchanged: it is read-only, and after the release it has
  no lineage to read for the Run.
- One new ADR, ADR-0313, records the rule; the glossary gains **Stale Owner
  Record** and revises **Reconcile Command**.

## Open Questions

None.
