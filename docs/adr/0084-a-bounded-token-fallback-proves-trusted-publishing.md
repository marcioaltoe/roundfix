---
status: accepted
created_at: 2026-07-31T00:00:00Z
updated_at: 2026-07-31T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A bounded token fallback proves trusted publishing without risking a partial release

npm does not validate a trusted publisher configuration when it is saved and
exposes no read-only endpoint, CLI command, or dry-run that proves OIDC
identity before a publish; configuration is per-package, so a correct launcher
says nothing about the fifth platform package. A publication preflight can
therefore establish registry-state eligibility — a used version, a
post-unpublish cooldown — but never identity, and an `ENEEDAUTH` on the third
of six coordinates would leave exactly the partial release ADR-0082 forbids.
During a bounded rollback window the release workflow therefore retries a
coordinate that fails OIDC authentication with the existing `NPM_TOKEN`,
records every coordinate that needed the fallback, and completes the release
set; the token and the fallback are removed together after one release that
publishes all six coordinates with the fallback untouched. Rehearsing on a
prerelease tag was rejected because it burns a version per attempt and still
proves nothing about the next release's configuration, and accepting the gap
was rejected because it makes ADR-0082's guarantee conditional on a
configuration npm will not let us verify.
