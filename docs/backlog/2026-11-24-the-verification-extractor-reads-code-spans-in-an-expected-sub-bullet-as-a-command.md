---
type: fix
status: open
created: 2026-11-24
spec: null
reason: null
---

# The Verification extractor reads code spans in an `Esperado` sub-bullet as a command, and `spec check` does not flag it

## Symptom

In corrective `task_03` of the Quaestor project's Spec `0006-codigo-em-ingles`, the `## Verification` section had a command bullet followed by an expectation sub-bullet whose names were between backticks (`  - Esperado: o padrão traz `chave` e não repete `key`. ...`). The Run's Pre-work Verification treated `chave` as a command. The log `batch-001-probe-2.log` says `sh: chave: command not found`, and the Task failed with `Pre-work Verification unknown: runner could not observe a verdict for command "chave": command could not be executed (exit 127)`. Removing the backticks from the expectation line fixed it. `roundfix spec check <slug> --strict` passed with no finding; only the Run found it.

## Expected

The command extractor reads only the code of the first-level bullet, or `spec check` flags an expectation sub-bullet whose code spans would become commands. The sender also suggests that `spec check` report "Pre-work Verification vacuous" when it already runs with `--run-verification`; that part is added to [the open entry on the vacuous verdict](2026-11-20-spec-check-called-honest-a-verification-the-daemon-refused-as-vacuous.md).

## Source

Cross-project message `01a1fa32-305f-7be2-9c17-bc200f117002` from the Quaestor project, 2026-11-20. Evidence cited by the sender: `~/.roundfix/artifacts/85a8e4c6a0c58989/runs/run_20261120T173226Z_56fed74aac57867d/verification/batch-001-probe-2.log`; commits `f171a38` and `a8d1a75` of gesttione-solutions/quaestor (PR #28). Nature: inferred. Recorded by triage on 2026-11-24.
