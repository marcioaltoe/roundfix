# Public Stop Request while queued

Run: `run_20260728T141723Z_eca0841d68e61067`

Entry points:

```text
roundfix implement --spec 0001-capacity-flow ... --detach
roundfix events run_20260728T141723Z_eca0841d68e61067 --filter task-status,verification,outcome
roundfix stop run_20260728T141723Z_eca0841d68e61067
roundfix events run_20260728T141723Z_eca0841d68e61067 --follow --filter outcome
```

The built detached Run had Task Capacity 2, Verification Capacity 1, and two
30-second real shell Verification commands. Before Stop:

- task_01 waiting cursor 10, started cursor 11;
- task_02 waiting cursor 12 and no task_02 shell marker.

The public Stop Command recorded a Stop Request while the Run state was
`Verifying`. Actual final evidence:

- task_01 passed at cursors 13–14;
- task_02 started at cursor 15 after the Stop Request;
- the shell log contains `start task_02` and `end task_02`;
- task_02 passed at cursors 21–22 and settled `completed` at cursor 23;
- the Run reached `Stopped` only at cursor 27, about 60 seconds after start;
- fresh Attach shows both Tasks `completed`.

This contradicts US8/Core Feature 9: the Task waiting for Verification
Capacity did not honor the Stop Request, started a child command, and did not
remain resumable.

The internal direct-context-cancellation integration check passes. The defect
is specifically the real public Stop Request path, which that test does not
exercise.
