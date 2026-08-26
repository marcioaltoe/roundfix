---
status: superseded # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-10T00:00:00Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null
superseded_by: ADR-0147
---

# The refusal that fired first is the refusal

Spec 0091 set out to make membership decide the verdict: Roundfix would read
what a runtime advertises, refuse an unadvertised model itself, and name the
advertised set, rather than inherit whatever the adapter chose to say. Its QA
gate then measured the built binary against all three live runtimes and found
the premise does not hold in execution: Codex, Claude and OpenCode each reach
requested-model application and return `selection_rejected` from the adapter
before Roundfix's membership check can speak. The membership verdict does fire
against fake harnesses, where no adapter competes for the refusal.

Chasing the promise means pre-empting every adapter — proving membership before
any application attempt on every runtime and version, forever, against
adapters Roundfix does not own. That is a race Roundfix loses by construction,
and winning it would buy a better-worded message for a selection that was
already going to be refused.

Roundfix therefore accepts whichever refusal fires first. The adapter's message
reaches the operator unchanged when the adapter refuses. The membership verdict
remains as the net beneath it, for the case that motivated the Spec: a runtime
that *declines* to refuse, which Spec 0091's characterization corpus recorded
for `claude` proving a model it never advertised. Silence from an adapter is
still not proof that a selection is sound.

The consequence is deliberate and is what the maintainer chose on 2026-08-10:
an unadvertised model fails, and the error travels to whoever is orchestrating
— the supervising agent or the person at the terminal — who corrects the
profile. Roundfix does not translate that failure into a better one.

Requiring the refusal to always name the advertised set was rejected with this
decision. It remains true where Roundfix owns the verdict, and it is not
promised where the adapter owns it: a message Roundfix cannot author is not a
contract it can keep.
