---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-14
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# A Spec archives before continuous integration answers

## Symptom

The mandatory order per Spec is: implement the graph with its authored gate,
archive, open the Pull Request, watch until Clean, and merge. Archiving therefore
happens **before** any answer from continuous integration, so a CI failure
arrives when the Spec is already history.

Measured in `conexus` on 2026-08-13 with Spec `0003-experiencia-financeira`. The
authored gate closed `pass`, `make verify` and `make smoke` passed locally, the
Spec archived, and the Pull Request opened. CI then failed a frontend test the
local suite approves: it fixes an absolute last-synchronised instant and asserts
the relative-age text the screen derives from the real clock, so the result
changes with the day and the timezone — passing at UTC−3 and failing at UTC.

Adding the corrective Task required un-archiving: moving the folder back to the
active directory, changing `status: archived` to `active` in the PRD, and
removing the archive date stamp. Editing the archived Spec instead would be
worse, because the contract requires it to stay byte-identical.

## Where

The autonomous loop's mandatory order, and the Archive Command's position within
it.

## Expected

Either archiving waits for the evidence that can still refuse the work, or
un-archiving is a supported act with its own command rather than three manual
edits against a contract that forbids touching archived Specs.

The first reading is that CI is part of the evidence a Spec's completion rests
on, and archiving before it is archiving on an incomplete claim. The second is
that archiving records the Spec's own gate and CI belongs to delivery, in which
case the loop needs the reverse operation it currently lacks.

## Evidence

Minted from the Inbox Entry
`inbox/roundfix/2026-08-13-a-spec-arquiva-antes-de-a-integracao-continua-responder.md`
in the Secondbrain. The same gap was met in `roundfix` on 2026-08-13 from the
other side: Spec 0094 archived, its Pull Request opened, and CI failed twice on
the documented spawn-density flake family before passing on a third run. Nothing
had to be un-archived only because the failures were flakes rather than defects.
