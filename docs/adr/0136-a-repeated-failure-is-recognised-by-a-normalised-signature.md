---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-16T09:00:30Z
updated_at: 2026-08-16T09:00:30Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# A repeated failure is recognised by a normalised diagnostic signature

A Task that fails twice with the same assertion is reported as new both times, and
a whole Run was spent on 2026-08-08 reproducing a diagnostic the loop had already
seen. Recognising the repetition needs a comparable form, and raw equality is not
it: two identical failures differ by timestamps, temporary paths, durations and
process identifiers, so a byte comparison would call every repetition new. The
signature is therefore a normalised form of the captured diagnostic — volatile
spans replaced before hashing — paired with the failing command and the Work Item.
Normalisation can collide, calling two different failures the same; that is the
accepted direction, because a false repetition costs a Supervisor one comparison
while a missed repetition costs a Run.
