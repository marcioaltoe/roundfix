---
type: fix
status: open
created: 2026-10-03
spec: null
reason: null
---

# A reopened QA gate has no path for a report written before row inputs

## Symptom

The `qa-gate` skill 0.0.6, as Roundfix 0.26.0 installs it, requires every executed row to declare non-empty `inputs:`, and the same section forbids adding inputs after execution (`.agents/skills/qa-gate/SKILL.md` around lines 344–349). A gate reopened on a report written before 0.0.6 has no valid move. fluxus's pre-PR review flagged the contradiction on 2026-10-02. It could only dismiss the finding, because the skill is vendored from Roundfix.

## Expected

The skill states that a row re-executed in a new pass declares its inputs for that pass, which is not adding inputs after the fact. A row carried from a report without `inputs:` is never carriable: it re-runs. The Daemon's carry logic already treats a row without inputs as `no inputs`.

## Source

Secondbrain `inbox/roundfix/2026-10-02-qa-gate-0-0-6-relatorio-antigo-sem-inputs.md` (fluxus).
