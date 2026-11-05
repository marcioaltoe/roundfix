### attach

```bash
roundfix attach [<run-id>]
```

`<run-id>` may also be the eight-digit short id that the
`runs list` table shows, or a longer prefix of the id's random suffix, when
it names one Run; a prefix that names several Runs exits `2` and lists
them ([ADR-0363](../../adr/0363-runs-list-prints-a-table-at-a-terminal-and-a-short-run-id-names-a-run.md)).

`attach` is a Read-Only Run Command: it runs from any directory. If
`worktree.location` would be refused there, it prints that refusal once on
stderr as a warning, naming the detected Git root and the `.git` it found,
and continues.

Read-only. With a Run ID it replays the Run Event Journal and follows live
events — never creating Runs, fetching, starting Agents, committing, pushing,
stopping, or resolving threads. Without a Run ID at an interactive terminal it
opens the Run Browser; in non-interactive mode it exits `2` naming
`roundfix runs list`. When the Run Browser closes without a selection it
prints `roundfix: Run Browser closed (<key>); no Run selected.` on stderr and
exits `0`. The Live Run View shows a `WORK QUEUE` pane next to a
`SESSION.TIMELINE` pane grouping Run Events by Batch; raw payloads never
render inline and full content stays in the Detail Modal. For spec Runs its
header reports `Task Capacity` and `Verification Capacity`, and each Task row
uses `Agent working`, `Waiting for Verification`, or `Verifying` before the
Daemon settles it. Verification Feedback returns that Task to `Agent working`;
meaning never depends on color.

