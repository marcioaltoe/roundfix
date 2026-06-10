# Resolve uses downloaded review issues

Roundfix treats `roundfix resolve --pr <n>` as work over all downloaded Unresolved Review Issues for an Open Pull Request by default, not as a Review Source fetch. Without an explicit Round, `resolve` processes all Compatible Artifact Rounds instead of choosing only the latest Round. A Round is only an optional filter. This keeps `fetch`, `resolve`, and `watch` on one durable history while preserving a clear command boundary: `fetch` downloads, `resolve` resolves downloaded issues, and `watch` automates both across review rounds.

A new `resolve` Run reuses Compatible Artifacts already present in the Artifact Directory. Compatible Artifacts match the requested Head Repository, PR Head Branch, and pull request number. If a Round is provided, they must also match that Round.

When `resolve` includes multiple Compatible Artifact Rounds, Roundfix deduplicates repeated unresolved Review Issues before batching. Deduplication uses the Review Issue Fingerprint, preferring `source_ref` when present and otherwise using a provider-specific fingerprint such as `review_hash`. Only the newest occurrence is assigned to an Agent; older occurrences are associated to the newest issue.

This means Roundfix does not need `fetch` to upsert Review Issue artifacts in place. Repeated fetches may preserve repeated Review Source findings in newer Round directories, and `resolve` remains the boundary that decides which occurrence should create Agent work.

After the assigned newest occurrence reaches `resolved` or `invalid`, the daemon marks older duplicate occurrences with terminal `status: duplicated` and sets `duplicate_of` to the newest occurrence. The Agent does not set `duplicated`, because duplicated status is a daemon-owned artifact bookkeeping result, not direct Agent work.

Duplicated older occurrences are local-only and do not resolve Review Source threads separately.

Newest occurrence is deterministic: higher Round wins first; if duplicate occurrences are in the same Round, later `source_review_submitted_at` wins; if that is missing or tied, later `round_created_at` wins. If the newest occurrence is still ambiguous, `resolve` fails during Preflight Validation instead of choosing nondeterministically.

If no artifacts match the requested pull request and Round scope, `resolve` fails during Preflight Validation before creating a Run and tells the user to run `fetch` or `watch`.
