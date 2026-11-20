---
status: accepted
knowledge: true
created_at: 2026-11-18T21:00:00Z
updated_at: 2026-11-18T21:00:00Z
deprecated_at: null
superseded_by: null
---

# A renamed Spec keeps its Spec Number and is one row of the Spec queue

The status view of the Follow-up Command counts a phase entry of an
Authoring Phase for a slug only while no other slug of the same Spec Number
holds that Spec. It amends ADR-0410, under which a phase entry made its slug
a Spec being authored whenever the checkout held no PRD and no Archive Record
for that slug.

## Context

On 2026-11-18 the status view listed
`0401-the-phase-boundary-follows-only-the-release-of-a-spec` as
"authoring: tasks" after an amendment had renamed the Spec to
`0401-a-supervisor-session-ends-only-when-the-maintainer-asks` and the Spec
had been delivered. The Phase Record is append-only, so the entries of the
old slug stay for seven days, and the status view joined them to the
checkout by slug alone. The `write-prd` skill fixes the identity a rename
keeps: "Numbers are never reused and travel with the spec when archived."

## Decision

A phase entry for a slug whose folder is not under the Spec Root does not
make that slug a Spec being authored when another slug of the same Spec
Number has a folder under the Spec Root, an Archive Record, or a Delivery
Queue item of any stage, a merged one included. Among slugs that only phase
entries hold, the one whose counted entry is the latest is the Spec, and the
others are its earlier slugs. A folder under the Spec Root always keeps its
row, as before. A slug without a Spec Number is unaffected.

**What stays.** The view still writes nothing and reads no other source. A
phase entry still counts only for the checkout, its linked worktrees and the
working trees with the same `origin` remote, within seven days. The Phase
Record is not rewritten: the earlier slug's entries stay until Run Retention
removes their day file.

## Considered Options

- **Keep the row and mark it as a renamed slug.** Refused: the row describes
  no work. The table would show one Spec twice, and the maintainer would
  still have to learn that the mark means nothing to do.
- **Count only phase entries whose slug has a folder in the checkout.**
  Refused: an author often works in another clone of the repository, and
  ADR-0410 counts its entries before the plan reaches the checkout.
- **Rewrite the Phase Record when a Spec is renamed.** Refused: the record is
  append-only, and a rename is a Git change that no Roundfix command sees.
- **Read the rename from Git history.** Refused: the slug of a plan that
  never reached the default branch has no history in the checkout, and
  rename detection costs a `git log` per status call.

## Consequences

- One Spec Number is at most one row of a Spec being authored.
- A repository that reused a Spec Number against the `write-prd` rule sees
  only one of the two Specs being authored. The rule forbids that reuse, so
  the view relies on it.
- While the checkout still holds a renamed Spec's old folder, that folder is
  its row until the renamed plan reaches the checkout.
