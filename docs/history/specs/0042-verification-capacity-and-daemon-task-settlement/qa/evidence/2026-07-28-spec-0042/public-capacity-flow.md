# Built-CLI capacity flow

Run: `run_20260728T140033Z_867f7808dae52080`

Entry point: the built `bin/roundfix implement --spec 0001-capacity-flow`
against a disposable Git repository on `ma/qa-public-capacity`.

Parity: the external ACP Task Agent boundary was a local fake. The built
Roundfix binary, Config loading, model proof protocol, Run Database, Run and
Task Worktrees, Git commits, real `sh -c` Verification commands, Run Event
Journal, event replay, Attach, integration, and settlement were production
paths.

Observed:

- the Live Run View reported `Task Capacity: 2` and
  `Verification Capacity: 1`;
- both Agent processes crossed a two-party readiness barrier, proving overlap;
- the shell log was serialized:

```text
start task_02
end task_02
start task_01
end task_01
```

- task_02 waited at cursor 10 and started at cursor 11; task_01 waited at
  cursor 12 and did not start until cursor 15, after task_02's passed verdict;
- both Tasks settled `completed`, the Run reached `Clean`, and the user
  checkout had two Daemon-created Task commits;
- a fresh `roundfix events` replay and two fresh Attach invocations preserved
  both capacities, ordering, terminal Task files, and Clean outcome.

The fake Agent deliberately attempted `status: completed` for task_01 and
`status: failed` for task_02. The final user-checkout Task files both contain
`status: completed`, proving the Agent values did not bypass the real shell
gate or control settlement.
