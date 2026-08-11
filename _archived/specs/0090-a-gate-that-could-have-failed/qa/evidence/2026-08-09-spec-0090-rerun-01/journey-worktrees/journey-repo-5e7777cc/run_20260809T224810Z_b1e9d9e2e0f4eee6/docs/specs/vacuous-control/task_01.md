---
task: task_01
spec: vacuous-control
status: failed
type: backend
---

# Task 01: Refuse the vacuous gate

## Negative Control

- `remove the unrelated xhigh fallback and require the gate to fail`
- `change the required profile while preserving the unrelated xhigh fallback`

## Verification

- `grep -q 'reasoning_effort: xhigh' .roundfixrc.yml` — expected: this deliberately vacuous gate is refused before Agent work.
- `test -f missing-before-agent` — expected: observed non-zero command proves the full probe list ran.
- `printf 'also vacuous\n'` — expected: second offending command is aggregated in the refusal.

## Result

The Agent ran and completed every requested change. This sentence is
deliberately inaccurate so the Daemon record can contradict it.
