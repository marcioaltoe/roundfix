---
type: fix
status: promoted
created: 2026-10-12
spec: 0257-the-baseline-respects-the-repository
---

# The Baseline hardcodes the Better Auth base path

## Symptom

The Standard TypeScript profile's backend guide states that Better Auth owns
`GET` and `POST` for `/api/auth/*`. Its `setup-context.json` repeats
`"scope": "/api/auth/*"`. Both sit inside Baseline markers, so an adopter
cannot edit them. The literal also appears in
`internal/baseline/assets/decisions.json`, the formatter golden fixture and
the `setup-context-driven` skill (`SKILL.md`, around line 198).

Better Auth's `basePath` is chosen by each repository. Vortex Spec 0106
(queued on 2026-10-10) removes the whole `/api` prefix, and Better Auth then
answers only on `/auth/*`. From that point the guide describes a route that no
longer exists. The rule is still right. Only the literal is wrong.

## Where

Standard TypeScript profile: `guide.backend`, `docs/agents/setup-context.json`, `internal/baseline/assets/decisions.json`, the `setup-context-driven` skill.

## Expected

Maintainer decision, 2026-10-11 ("Decisão no Setup Manifest"): a recorded
repository decision (for example `auth.base_path`) in the Setup Manifest. The
guide and `setup-context.json` render that value. Current adopters keep
`/api/auth` as the default, and the vortex declares `/auth` at its next
`baseline update`.

## Evidence

- Secondbrain `inbox/roundfix/_triaged/2026-08-24-o-guia-de-backend-fixa-api-auth-como-escopo-do-better-auth.md`
  (vortex, 2026-08-24, rerouted from `skills` on 2026-10-09).
- Secondbrain `inbox/roundfix/_triaged/2026-10-10-o-vortex-remove-api-auth-e-a-baseline-fixa-esse-caminho.md`
  (vortex, 2026-10-10).
