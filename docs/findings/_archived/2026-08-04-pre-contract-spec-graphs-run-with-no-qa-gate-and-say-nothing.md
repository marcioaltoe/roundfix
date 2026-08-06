---
status: done
created_at: 2026-08-04
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-qa-gates-and-verification-evidence.md
---

# 2026-08-04 — Pre-contract Spec graphs run with no QA gate and say nothing

status: pending

## What was observed

Roundfix 0.3.1 removed the `--qa` flag and moved the gate into the Task Graph
as an authored `qa: task_NN` declaration. In a consumer repository (oraculum)
with eleven queued Specs authored under the old flag, **none of the eleven
declares a gate**:

```
0015-health-e-identidade-operacional     *** NO qa: DECLARATION ***
0016-contrato-de-perimetro-de-filiais    *** NO qa: DECLARATION ***
…
0025-consultor-financeiro-mcp            *** NO qa: DECLARATION ***
```

`roundfix implement` accepts every one of those graphs. Preflight does not
refuse them, no stderr line mentions the absence, and the Run ends `Clean` with
its Tasks completed and **no gate ever executed**. The outcome is
indistinguishable from a Spec that passed QA.

The gap is not theoretical. When the gate was authored by hand for Spec 0015
and finally ran, it returned `fail` twice and found two `Trust-Damage` defects
that every focused Task Verification had passed:

- the assembled stdio composition dropped `ambiente_presumido`, because the bit
  was computed as `ambienteFonte === undefined && configuracao.ambientePresumido`
  and the real composition supplies `ambiente` through tool contexts;
- the compiled `/health` could never reach `ok`, because the ERP probe threw
  without an operational Organization and the entry point never supplied one —
  reporting a dependency as *down* that was never checked.

Both were reachable only by driving the compiled process through a real user
seat. Both would have shipped silently on any of the other ten Specs.

## Root cause

Two rules that are individually defensible compose into a silent skip:

1. The `write-tasks` skill instructs authors to leave legacy graphs alone —
   *"An absent declaration remains valid only for a legacy graph proven to
   predate the contract; leave that graph byte-identical."*
2. The Daemon treats an absent declaration as "no gate authored" and proceeds
   to a normal `Clean`.

So the skill says *don't touch it*, the Daemon says *nothing to run*, and the
Supervisor gets a green Run. Nothing in the pipeline states that the Spec shipped
without QA. A contract migration that silently downgrades verification for every
pre-existing Spec is the one migration shape that cannot be left to notice.

The blast radius scales with adoption: every repository that queued Specs before
0.3.1 has the same latent set, and the more Specs were written ahead of time —
exactly the pattern a "write the PRD, let the agent build" workflow encourages —
the more of them ship ungated.

## What would settle it

- Make the absence loud rather than silent. At minimum one stderr line at Run
  start and one line in the deterministic report: `qa not declared — no gate
  will run for this Spec`. An operator reading a `Clean` outcome must be able to
  tell "passed QA" from "had no QA".
- Consider refusing at Preflight for graphs whose Spec folder was created after
  the contract shipped, keeping the byte-identical exemption only for graphs
  that provably predate it. The current rule trusts authoring date without
  checking it.
- Ship a migration path: `roundfix spec qa-declare <slug>` (or a `doctor` line)
  that reports which Specs under `specs.root` carry no declaration, so the gap
  is discoverable in one command instead of by grepping eleven folders.
- Reconsider the "leave legacy graphs byte-identical" instruction for graphs
  whose Tasks are all still `pending`. Nothing is in flight, nothing is
  invalidated, and authoring the gate costs one node.

## Evidence

- Consumer repository: oraculum at `20fc93b`, eleven Specs under `docs/specs/`.
- Gate authored by hand for Spec 0015 as `task_05`; reports
  `qa/qa-report-2026-08-04.md` (`verdict: fail`, F-001/F-002) and
  `qa/qa-report-2026-08-04-02.md` (`verdict: pass`, 10 rows, zero findings)
  are archived under
  `docs/specs/_archived/0015-health-e-identidade-operacional/qa/`.
