---
spec: 0363-a-reused-owner-pid-means-the-owner-is-gone
status: active
created: 2026-11-06
surfaces: [cli, docs]
type: fix
---

# A reused owner PID means the owner is gone

`roundfix doctor` 0.136.0 ends with `residue: partial` in every repository once the recorded owner of a terminal Run has died and the system has given its PID to another process: the proof of the owner's identity reports a start-identity mismatch, and Doctor files it as a process table it could not read. No command clears the warning. The same reused PID also keeps a dead owner present for three readers that check only whether the PID is alive: the orphan reclaim of an Active Run's lock, the Delivery Queue's Run-slot holders and `roundfix deliver start`. After this Spec a proven mismatch means the owner is gone for all four, `partial` again means an inspection that could not finish, and every place that signals or that cannot prove the mismatch keeps today's behavior.

## Sources

One ai-memory message, recorded as a Backlog Entry adopted at birth (see
`references/_index.md`):

1. [`references/2026-11-06-doctor-residue-stays-partial-when-an-owner-pid-is-reused.md`](references/2026-11-06-doctor-residue-stays-partial-when-an-owner-pid-is-reused.md),
   reported by argus, `type: fix`.

At the authoring base `docs/backlog/` holds no open entry and `docs/findings/`
no unresolved Finding that shares this context.

## Prerequisites

None in `requires`. Authored against `9a285d77` (v0.136.1). The queue
delivers 0360, then 0362, then this Spec, then 0361. Spec 0360 also changes
`internal/cli/deliver_workflow.go` (not `RunSlots`),
`docs/user-guide/commands/deliver.md`, `docs/user-guide/commands/doctor.md`,
`.agents/skills/roundfix/references/deliver.md` and
`.agents/skills/roundfix/references/setup.md`, the Roundfix skill's
`SKILL.md` pair, `skills/testdata/owned-skill-versions.json` and
`docs/references/behavior-surfaces.json`; Spec 0362 changes the Roundfix
skill's `SKILL.md` pair, the version record and the Behavior Surface Record.
This Spec edits only the function `RunSlots` in `deliver_workflow.go`, one
paragraph of each guide and reference it names, and regenerates every record
by its command, so a rebase reruns the command instead of merging bytes.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; Runs, queues and owner records keep their identities and
  columns. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the change reads the local process table and the local Run
  Database. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0373 (new) records the rule.
  ADR-0313: "A proven owner identity mismatch releases a terminal Run's owner
  record without a signal", so a mismatch is already proof the owner ended.
  ADR-0044: "Orphaned Run locks are reclaimed only on proven owner death",
  which the reclaim keeps by counting a mismatch as that proof. ADR-0127:
  "Process residue is a readiness fact, not a Run record", so Doctor still
  writes nothing. ADR-0279: "Status reports a Queue Owner proven gone", whose
  status, retry and resume rules stay unchanged. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AZ batch 4 rules answer of 2026-11-06 quoted in the authorization
  record, which covers this Spec as a new Spec of the batch; no Makefile,
  workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0363-a-reused-owner-pid-means-the-owner-is-gone/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/setup.md`, `skills/roundfix/SKILL.md`.

## Goals

- `roundfix doctor` reports `residue: ok` when the only obstacle was a
  terminal Run whose owner PID now belongs to another process, and never
  reports that process's tree as Roundfix residue.
- `residue: partial` stays for a process table that cannot be read and for a
  Run whose recorded identity is unproven.
- The orphan reclaim of an Active Run's lock, the Delivery Queue's Run-slot
  holders and `roundfix deliver start` treat a proven mismatch as an owner
  that is gone, exactly as they treat an absent PID.
- Nothing that sends a signal, and nothing that cannot prove the mismatch,
  changes its behavior.

## User Stories

1. As a maintainer running `roundfix doctor` in any repository, I want a
   reused owner PID to stop showing as a process-table failure, so that
   `partial` means something I can act on.
2. As a maintainer whose Implement Run's owner died and whose PID was reused,
   I want the next Run for that Spec to reclaim the stale lock, so that I do
   not have to stop a Run that no longer exists.
3. As an operator of the Delivery Queue, I want a dead queue owner or Run
   slot holder whose PID was reused to stop blocking `deliver start` and the
   slot wait, so that the queue moves without a manual release.

## Core Features

1. One owner-gone function in the store: absent, or alive with a proven
   different start identity, and whether the PID was reused.
2. Doctor's residue check counts a mismatch as an owner that is gone and
   reports neither residue nor a read failure for that Run.
3. The orphan reclaim of an Active Run's lock (CLI preflight and the store's
   own guard), `RunSlots` and `deliver start` use the function; the reclaim
   reason names the reuse.
4. Documentation: the doctor and deliver guides and the Roundfix skill's
   setup and deliver references state the rule; the Roundfix skill
   records its next version.

## User Experience

On argus, `roundfix doctor` prints `residue: ok (no process residue found)`
instead of three `partial` lines. A new `roundfix implement` for a Spec whose
earlier Run's owner PID was reused prints the existing reclaim warning with
the reason `owner process <pid> gone (pid reused); lock reclaimed`. Nothing
else in any output changes.

## Non-Goals / Out of Scope

- Releasing owner records from Doctor; ADR-0313: "another read command was
  rejected", and `reconcile --apply` keeps the release.
- Force Stop, `deliver stop`, any terminate or signal path, `reconcile`,
  `deliver status`, `retry` and `resume`, `runs list` and its TUI, `gc`,
  worktree disposal and the carry-forward staging sweep, whose own identity
  rule is left as it is.
- Legacy Runs with no recorded identity, and the `runs.max_active` count in
  the store's Run creation, which counts Active Runs without reading PIDs.
- The Makefile, the CI workflows, `go.mod`, `go.sum` and the lint or
  formatter configuration.

## Success Metrics

1. Success Metric: Doctor over a terminal Run whose owner proof returns the
   mismatch error reports `residue: ok (no process residue found)`; at
   `9a285d77` it reports `residue: partial` with the mismatch text.
2. Success Metric: Doctor over an unreadable identity, an unproven recorded
   identity, or a process-detail failure keeps its `partial` line byte for
   byte.
3. Success Metric: the owner-gone function answers gone and reused for a
   live PID with a different comparable identity, gone and not reused for an
   absent PID, and present for a matching identity, an unreadable identity, a
   foreign-platform token, an empty recorded identity on a live PID and the
   caller's own PID.
4. Success Metric: with a recorded owner PID that is the test's live parent
   process and a recorded identity that differs from its real one, the
   Active Run lock is reclaimed with the reuse reason, `RunSlots` excludes the
   holder and `deliver start` releases the dead queue owner; with the
   matching identity none of them does. At `9a285d77` all three keep the dead
   owner present.
5. Success Metric: no test sends a signal to a real process, and every
   existing Force Stop, `deliver stop`, reconcile and orphan-reclaim test
   passes unchanged.

## Acceptance evidence

The outside-evidence row rests on a source this Spec did not write:
POSIX.1-2017 (The Open Group Base Specifications Issue 7), Base Definitions,
section 4.14 Process ID Reuse
(https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap04.html,
read through Exa on 2026-11-06): "A process ID shall not be reused by the
system until the process lifetime ends." A live process at the recorded PID
with a different start identity therefore proves the recorded owner's
lifetime ended.

## Unreachable Acceptance

None. Every criterion is checked with fake process controllers, the test's
own parent process read without a signal, temporary Run Databases and
temporary repositories.

## Validation

The decisions were checked on 2026-11-06 against these sources.

- Characterization at `9a285d77`: a scratch test gave
  `doctorResidueResults` a terminal Run whose recorded owner was the test's
  live parent PID with the identity `<GOOS>:0.0`. `InspectTreeProcesses`
  returned `ErrOwnerProcessIdentityMismatch`, and Doctor printed
  `residue: partial (could not read process table for Run run_reused: owner
  process <pid> failed step "prove owner process identity": owner process
  identity cannot be proven: live process start identity does not match the
  recorded owner identity)`, the text argus reported.
- Reading at `9a285d77`: `Residue` (`internal/cli/doctor.go`) files every
  `ReadLineage` error as a read failure; `reclaimOrphanedActiveRun`
  (`internal/cli/cli.go`), `Store.ReclaimOrphanedRun`
  (`internal/store/store.go`), `RunSlots`
  (`internal/cli/deliver_workflow.go`) and `runDeliverStart`
  (`internal/cli/deliver.go`) call `store.ProcessAlive` on the PID alone;
  `forceStopRun` refuses a mismatch, `inspectReconcileProcesses` releases it,
  and `deliver status`, `retry` and `resume` read it as gone.
- Pending ai-memory messages for this project: the three messages of this
  batch, which Specs 0361 to 0363 adopt; none was consumed.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "PID reuse process start time identity doctor residue" --all --files --min-score 0.3`;
  no page decides how a reused PID is classified, so the repository's ADRs
  govern.
- Exa: the POSIX passage quoted above.

## Glossary

None.

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AZ batch 4 brief for Spec 0363 of 2026-11-06 and the maintainer's answers below.

The maintainer asked no separate question for this report, a clear defect
with the reporter's own correction; the batch's rules cover it as a new Spec
(AskUserQuestion, 2026-11-06 ~12:00 BRT), verbatim:

- Rules: "Mesmas do AZ (Recommended)" — option text: "Sem supervisão, merge
  pela fila, conserto e retry quando parar. Release assim que cada Spec ficar
  pronta; resposta ao OnionCry, ao sql-library e ao Jarvis."
- The maintainer's standing instruction of 2026-11-05: "Quero que faça os
  releases hoje assim que estiverem prontos. E se existe no baseline regra
  para um release por dia por algum motivo deve ser removida."
- The argus message's own correction: "quando o start identity do PID vivo
  não bate, classifique o dono como morto (`owner gone: pid reused`) e não
  registre falha de leitura. Assim o residue fica `ok`, a não ser que haja
  outro resíduo de verdade. Guarde o `partial` para quando a tabela de
  processos realmente não pode ser lida."

Taken from those and the reading, with no further interview:

- Doctor prints no extra line for a reused owner: the owner record is
  harmless and `reconcile --apply` releases it, and an extra line would make
  `ok` read like a warning.
- The three PID-only readers join the rule because the same reused PID keeps
  a dead owner present in each; the store function makes the rule one
  implementation.
- Every signal path, every reader that already treats a mismatch as gone,
  and every answer short of a proven mismatch stay as they are (ADR-0373
  lists them).
- ADR-0373 records the rule; no older ADR is amended, because each already
  treats a proven mismatch as the end of the owner or does not speak to it.
- The release is the operator's patch right after the merge; no Task
  releases.

## Open Questions

None.
