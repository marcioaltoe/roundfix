### supersede

```bash
roundfix supersede --spec <slug> --by <slug> --reason <text>
```

Use this command when another Spec delivered the content of a Spec that was
never decomposed. It records the supersession in `_supersession.md`; do not
edit the Spec's status or write a note instead. `--spec` names the superseded
Spec, `--by` names the active or archived Spec that delivered its content, and
`--reason` records the explanation. The command writes only that amendment,
creates no Run, and never commits or pushes.

The command refuses an unknown superseded Spec, an unknown superseding Spec, a
self-supersession, or a Spec that already carries a supersession. Archive
accepts this recorded supersession in place of completed Tasks and a passing
gate when the Spec has no Task Graph; every other archive precondition still
applies, and a Spec with a Task Graph keeps its existing evidence rules.

Options:

- `--spec <slug>` — Superseded Spec slug under the configured Spec Root.
- `--by <slug>` — Active or archived Spec slug that delivered the content.
- `--reason <text>` — Explanation recorded in the supersession amendment.

Exit codes:

- `0` — supersession recorded.
- `1` — supersession write failed.
- `2` — Preflight Validation failed, including an unknown Spec, a
  self-supersession, or an existing supersession.

