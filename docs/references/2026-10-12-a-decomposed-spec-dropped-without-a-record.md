# A decomposed Spec dropped without a record (2026-10-12)

This reference keeps the outside evidence behind `roundfix cancel`, so a
reader can check it without the network or another repository.

## The adopter event

On 2026-10-10 the maintainer of an adopter repository reviewed its Spec plan
and dropped eight open Specs that already had Task Graphs. One of them had 5
Tasks `completed` and 3 `pending`. The other seven had every Task `pending`.
The reasons were a missing consumer, a decision taken in another project, and
a redesigned connection contract.

The adopter found no Roundfix command that could record the drop:

- the Supersede Command needs a Spec that delivered the content, and the
  Archive Command accepts a supersession only for a Spec without a Task Graph;
- the QA Archive Override still requires every non-QA Task to be `completed`;
- Task `status` belongs to the Daemon, so the pending Tasks cannot be closed
  by hand.

The adopter deleted the eight folders in one commit, handed the adopted
sources back to its backlog or to the Spec that now owned them, and wrote the
reasons into a planning document. No Archive Record remains. What the Spec
with 5 completed Tasks delivered has no history entry.

The adopter's plan table recorded two outcomes for the eight Specs: some were
"Substituída" (replaced by another Spec that took over a source) and one was
"Deprecada; saúde S2S já entregue" (dropped, while what it had delivered
stays).

Provenance: the Secondbrain Inbox Entry
`inbox/roundfix/_triaged/2026-10-10-nao-ha-caminho-para-cancelar-uma-spec-ja-decomposta.md`
(origin: the adopter, created 2026-10-10, resolved to the Backlog Entry this
reference supports), and the adopter's mirrored planning document
`projects/vortex/mirror/docs/references/2026-10-10-plano-de-specs-ate-o-release.md`,
section "Specs fora do plano". No adopter file is copied here.

## The current Roundfix contract, measured on `7a8b7f70`

- `roundfix supersede --help`: "Records that another active or archived Spec
  delivered this Spec's content."
- `roundfix archive --help`: "The QA override requires an approval source and
  reason, still requires every non-QA Task completed".
- The record parser accepts only the dispositions `pass`, `partial`,
  `qa-override`, `failed-qa`, `superseded` and `no-qa`.

## Published practice

- PEP 1, "PEP Purpose and Guidelines" (<https://peps.python.org/pep-0001/>):
  "A PEP can also be “Rejected”. Perhaps after all is said and done it was not
  a good idea. It is still important to have a record of this fact. The
  “Withdrawn” status is similar". PEP 1 keeps "Withdrawn" and "Superseded" as
  separate states, and a superseded PEP carries `Superseded-By`.
- GitHub Changelog, "The new GitHub Issues – May 19th update"
  (<https://github.blog/changelog/2022-05-19-the-new-github-issues-may-19th-update/>):
  "When closing, you can opt to close as `complete` (default option) or
  `not-planned`."

Both separate "done" from "dropped" and keep a record of the drop, which is
what the `cancelled` disposition does. Both keep "replaced by another
proposal" as its own state, which is why a cancellation does not name a
delivering Spec.
