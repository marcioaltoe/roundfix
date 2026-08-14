---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-14
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The vacuous-Verification event lists the commands that passed

## Symptom

The vacuity check is right and valuable: a Verification command that passes
before the work exists cannot tell work from nothing. What costs Runs is reading
the event it emits.

```json
{"category":"verification","work_item":"task_01","phase":"failed",
 "classification":"verification_vacuous",
 "commands":["bunx turbo run test --filter=@fiscus/backend -- tests/…",
             "bunx turbo run typecheck --filter=@fiscus/backend"]}
```

The natural reading of `commands` is "the commands it ran". It actually means
"the commands that passed against the unchanged tree" — the offenders. The Task
declared four commands; the two absent from the list were the ones that failed,
which is to say the correct ones.

That reading produced two wrong conclusions in sequence: first that Roundfix was
discarding `grep` commands, then that the rule was about the sequence rather
than about each command. Only the probe log settled it —
`<artifact_dir>/runs/<run-id>/verification/batch-NNN-probe-1.log` showed the new
command running and failing correctly, which proved the problem lay in the
others.

## Where

The `verification` Run Event with `classification: verification_vacuous`, and
the console line derived from it.

## Expected

The field says what it holds. Naming it for the offenders — or reporting each
command with its own verdict, so a reader sees which passed and which failed
without opening a probe log — turns a two-Run diagnosis into a one-line read.

Worth settling in the same work: whether the probe log path belongs in the event,
since it is what actually resolved the confusion both times.

## Evidence

Minted from the Inbox Entry
`inbox/roundfix/2026-08-12-o-evento-de-verificacao-vacua-lista-os-comandos-que-passaram.md`
in the Secondbrain, captured from a `fiscus` session. The same misreading slowed
diagnosis in `roundfix` on 2026-08-12 while implementing Spec 0094, where eleven
authored commands were refused across five probe cycles.
