---
type: fix
status: promoted
created: 2026-10-26
spec: 0323-the-status-the-park-table-and-write-tasks-say-what-to-do-next
reason: null
---

# The Park Class table joins the secret rewrite to every finding park

## Symptom

The Park Class table of the `deliver` guide and of both Roundfix skill
reference copies ends its `finding` row with
`secret-found`: in `<worktree>`, rewrite the item branch from the earliest
commit the blocker names so that each named line holds a value no Secret
Pattern matches … then run `roundfix deliver retry <slug>`. The colon after
the last blocker reads as the action for the whole row, so a reader of a
`run-unresolved`, `gate-failed` or `review-findings` park is told to
rewrite history for a secret that does not exist.

## Where

`docs/user-guide/commands/deliver.md` (section on parked items, the
`| Class | Blockers |` table), `.agents/skills/roundfix/references/deliver.md`
and its mirror `skills/roundfix/references/deliver.md` (the
`| Class | Blockers or action |` table). The queue itself is right:
`internal/delivery/park_class.go` `ClassifyPark` gives the rewrite only to
`secret-found`.

## Expected

The table maps each class to its blockers only. The `secret-found` action
is stated for that blocker alone, in its own sentence, as the queue prints
it, the way the `archived-item-unproven` action already is.

## Evidence

Pre-PR review finding F1 of Spec 0299, dismissed on 2026-10-22 at 20:05 BRT
as outside that Spec's scope, and recorded as item 4 of the cycle AO triage
inputs: "Guia deliver.md, tabela de Park Class: a instrução de reescrever
segredo depois de 'secret-found:' parece valer para a linha 'finding' inteira
… Separar a ação por blocker (guia + as duas referências)." The row is
unchanged at `aa8b3db7`.
