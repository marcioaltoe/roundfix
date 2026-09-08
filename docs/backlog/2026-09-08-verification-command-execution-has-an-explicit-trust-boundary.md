---
type: fix
status: open
created: 2026-09-08
spec: null
reason: null
---

# Verification execution has no stated trust boundary for another author's Spec

## Symptom

`roundfix spec check --run-verification` executes authored shell commands,
and the authoring skills require that form. The contract does not say whether
a Spec supplied by another author is assumed trusted or requires a separate
execution authorization.

## Where

`internal/cli/spec_check.go`, `internal/daemon/verification_probe.go`, and
the `write-prd`, `write-techspec`, and `write-tasks` skill contracts. Spec 0116
increased routine use without changing the shared prober's execution model.

## Expected

State the trust and authorization contract at the execution boundary and in
Spec authorship. Integrate it with the user-requested Spec-contained authority
record. Provisional P1 owns the decision; this entry does not prescribe a new
prompt or sandbox and does not claim a demonstrated exploit.

## Evidence

Secondbrain `inbox/roundfix/_triaged/2026-08-30-run-verification-executa-comandos-de-spec-de-terceiros.md` records the original
observation. Triage on 2026-09-08 checked the current local sources named above.
The Inbox body remains the observation's provenance; this entry records intent.
