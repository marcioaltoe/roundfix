---
status: accepted
created_at: 2026-08-08T00:00:00Z
updated_at: 2026-08-08T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Profile readiness covers every configured Agent Work Category

`CONTEXT.md` defines Agent Selection Profile Readiness as requiring Exact Agent
Selection Proof for every distinct tuple, and the Doctor Command resolved only
the five required Agent Work Categories, so a profile configured for an optional
category was parsed, registered, and never proven. Measured on 2026-08-08: with
a `data` profile whose preferred Agent Selection was broken, `roundfix doctor`
printed `profiles: ok (5 distinct tuples; 10 category references)` and listed
only `claude` and `codex` under adapter readiness, while `roundfix profiles
validate --category data` reported that same profile with `source: project` and
failed it. The gap cost a session of misdiagnosis, because a configured profile
that a readiness command never mentions is indistinguishable from one that did
not register. Profile readiness therefore resolves every Agent Work Category the
effective configuration defines — the required five plus every optional category
present in the merged profiles — and adapter readiness enumerates the runtimes
those tuples reference. Categories that merely inherit `general` contribute no
distinct tuple and are not enumerated, so the change widens coverage without
inventing proofs. Resolving all ten categories unconditionally was rejected
because inherited categories add references without adding a tuple to prove.
