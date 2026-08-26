---
status: done
absorbed_by: 2026-08-06-rollup-baseline-and-derived-tooling.md
created_at: 2026-08-07
updated_at: 2026-08-26
kind: finding
---

# Two HTTP Contract defaults, and only one is read (2026-08-07)

Changing the HTTP Contract Decision default from `Post-only` to `REST` on
2026-08-07 surfaced a second declaration of the same default that nothing reads.

## What was observed

The default exists in two assets:

- `internal/baseline/assets/decisions.json` — the catalog decision declaration,
  `"default": {"mode": ...}`. `baselineDecisionDefinition` resolves the
  interactive prompt's default from `catalog.Decision(id)`, so this is the value
  a maintainer sees and accepts with Enter.
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json` — the
  profile's `httpContract` block, `"default": "Post-only"`.

No Go code reads the profile's `httpContract.default`. `normalizeHTTPContract`
consults the declaration's `modes` for validation and never its `default`. The
contract vocabulary in `internal/baseline/assets/contract-v1.json` declares an
error code `profile.httpContract.default.invalid`, and that code has no
implementation anywhere in Go.

The catalog default is now `REST`; the profile default still reads `Post-only`.
Nothing breaks, because the profile value is inert.

## Why it matters

An inert duplicate of a live value is a trap rather than a bug. The next person
changing this default will grep, find two, and either change the wrong one — and
observe no effect — or change both and assume both mattered. The declared but
unimplemented error code makes the profile field look load-bearing, which is
worse than an undocumented field.

The repository already carries a rule for this shape: an assertion reads the
constant it means, and a duplicated value changes at every occurrence in the
same commit. The rule assumes every occurrence is read; here one is not.

## Route

Not fixed here. The 2026-08-07 authorization for the REST default is bounded to
`internal/baseline/assets/decisions.json` and explicitly does not reach
`internal/baseline/assets/profiles/**`, so the divergence is recorded rather
than silently edited.

Whoever picks this up must first decide which of three the profile field is,
because the fix differs for each:

- a field that should be read, and the missing validation is the real defect;
- a field that duplicates the catalog and should be deleted along with its
  unimplemented error code;
- a deliberate per-profile override that the resolution path forgot to consult.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
