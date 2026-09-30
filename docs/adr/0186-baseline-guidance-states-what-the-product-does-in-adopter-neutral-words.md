---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Baseline guidance states what the product does, in adopter-neutral words

The Context-Driven Baseline ships its clauses to every adopter as mandatory
guidance. An audit on 2026-09-30 found clauses that described an older product:
an archive-then-review order the Delivery Queue does not follow, a refusal the
audit no longer makes, and an approval record no code reads. Other clauses
cited this repository's own Spec and ADR numbers, which name different records
in an adopter's repository. One module rendered a paragraph twice.

A Baseline clause now states only behavior the shipped product has. When a
clause asks for a control the product does not enforce, it says who enforces it
instead of describing a mechanism. A clause names Roundfix commands and the
review providers Roundfix configures, because Roundfix delivers the Baseline and
runs the loop. It never cites a Spec number or an ADR number: the rule is
written out where it applies. Every clause keeps its enforcement level when its
wording changes.

Three mechanical checks hold this in place. No two clauses share text. No
shipped guidance cites a Spec or ADR number. The loop order the guide declares
equals the order the Delivery Queue runs.

## Consequences

An adopter's guide no longer points at records it cannot open, and this
repository's guides lose their ADR links in the same clauses. The decisions
stay in `docs/adr/`. A change to the Delivery Queue's stage order now fails a
test until the clause changes with it. Wording drift that no check can see
still needs the release step that Spec 0195 adds.
