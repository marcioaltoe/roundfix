### attach

```bash
roundfix attach [<run-id>]
```

Read-only. With a Run ID it replays the Run Event Journal and follows live
events — never creating Runs, fetching, starting Agents, committing, pushing,
stopping, or resolving threads. Without a Run ID at an interactive terminal it
opens the Run Browser; in non-interactive mode it exits `2` naming
`roundfix runs list`. The Live Run View shows a `WORK QUEUE` pane next to a
`SESSION.TIMELINE` pane grouping Run Events by Batch; raw payloads never
render inline and full content stays in the Detail Modal. For spec Runs its
header reports `Task Capacity` and `Verification Capacity`, and each Task row
uses `Agent working`, `Waiting for Verification`, or `Verifying` before the
Daemon settles it. Verification Feedback returns that Task to `Agent working`;
meaning never depends on color.

