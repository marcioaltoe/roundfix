---
type: fix
status: open
created: 2026-11-05
spec: null
reason: null
---

# The release workflow survives a transparency-log 409 on the roundfix package

Twice the npm publication stopped half way. v0.130.0 (2026-11-03) published `@roundfix/cli-darwin-arm64` and then the transparency log answered `409` for `@roundfix/cli-darwin-x64`. v0.132.1 (2026-11-05, run 42056741468, attempt 1, step "Publish to npm") published the five `@roundfix/cli-*` platform packages and then answered `npm error error creating tlog entry - (409) an equivalent entry already exists in the transparency log` for `roundfix`. In both cases the rerun of the failed job refused in "Publication preflight" with `registry: <package>@<version> is already used`, so the version could never be completed and the next patch republished the same content (0.130.1, and 0.133.0 for 0.132.1).

Intent: a release whose publication stops after some packages can be finished, not abandoned. Candidates to evaluate: the preflight accepts a platform package already published at the release version when its tarball integrity matches the one this run built, so a rerun publishes only what is missing; and the publish step does not retry `npm publish --provenance` for the same tarball after a provenance entry was created (an inferred cause of "an equivalent entry already exists"; not verified).

Evidence: `gh run view 42056741468 --attempt 1 --log-failed`; `npm view roundfix versions` (no 0.130.0, no 0.132.1); the CHANGELOG entries for 0.130.1 and 0.133.0.
