---
status: accepted
created_at: 2026-10-12T00:00:00Z
updated_at: 2026-10-12T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The Better Auth base path is a repository decision that a Baseline update records

Better Auth mounts its routes under a `basePath` that each application
chooses; the library's default is `/api/auth`. The Standard TypeScript profile
already records the provider's route exception, scope included, in the
`auth.provider` decision (ADR-0076), and the backend guide renders it. Two
adopters have removed the `/api` prefix: one recorded in a repository-authored
note that the setup-owned guide now contradicts its own decision, and the
other queued the same removal for 2026-10. Neither could change the scope at
`roundfix baseline update`, because the update never changes a recorded
decision, and changing the provider value through `baseline plan` and
`baseline apply` appends a second Better Auth exception to `http.contract`
instead of moving the first one. On 2026-10-11 the maintainer chose "Decisão
no Setup Manifest".

**The decision.** `auth.base_path` is an optional string decision of the two
TypeScript profiles that select `auth.provider`. Its default is `/api/auth`.
When it is recorded, the provider's route exception covers `<base path>/*`,
the matching `http.contract` exception moves with it instead of being
duplicated, and the backend guide states the recorded path. When it is
absent, nothing changes: the recorded provider scope stays as it is, and the
guide bytes stay identical. A value is an absolute path with at least one
segment and no trailing slash, wildcard, whitespace or marker-shaped text;
`auth.base_path` without `auth.provider` is refused.

**The update records it.** `roundfix baseline update --decision <id=value>`
records or changes a decision the repository's profile selects, before the
plan is built, so the Plan Digest covers the new Setup Manifest. It parses the
value as `baseline plan --decision` does, refuses a decision the profile does
not select or fixes, and answers a new decision without `--adopt-suggested`.
The update stays the Managed Refresh: it never changes the profile or the
preservation mode.

## Considered options

Making the existing `routeException.scope` the only input was rejected: the
repository would have to restate the owner and the reason to move a path, and
the scope is a derived projection of where the application mounts Better
Auth, which is one string. Making `auth.base_path` a required decision was
rejected because every adopter's next update would stop at exit 3 for a value
nobody needs to change (ADR-0118 accepts that cost only when the answer is
not already recorded). Keeping the automation path alone (`baseline plan`
with every recorded decision, then `baseline apply`) was rejected because it
asks an adopter to restate every recorded decision to change one.

## Consequences

Current adopters keep `/api/auth/*` and byte-identical backend guides; their
next update rewrites only the Setup Manifest's recorded digests. An adopter that
mounts Better Auth elsewhere runs one confirmed update. The decision catalog
gains one optional decision and the backend template one token that renders
nothing while the decision is absent.
