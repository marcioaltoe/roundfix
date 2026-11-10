---
type: fix
status: open
created: 2026-11-10
spec: null
reason: null
---

# The Catalog Skill Bump leaves a stale catalog digest in the Setup Manifest

## Symptom

When one change edits a Baseline module and runs `make skills-third-party-bump`, the bump refreshes `docs/agents/setup-context.json` with a binary it built before `make baseline-digests` rewrote the profile's golden digest. The catalog digest it records is then stale and `TestTheSetupManifestIsCurrent` fails until someone runs the Managed Refresh by hand.

## Where

The `skills-third-party-bump` target of the `Makefile` and the order it shares with `make baseline-digests`.

## Expected

After a Catalog Skill Bump and `make baseline-digests` in the same change, the Setup Manifest is current with no extra command, whichever of the two ran first, or the bump says which command must follow it.

## Evidence

Spec 0370's TechSpec prescribed the manual step as its step 5: `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text` run twice, "the second reporting `File changes: 0`", because "The bump refreshes the Setup Manifest with a binary it built before `make baseline-digests` rewrote the profile's golden digest, so without this step the recorded catalog digest is stale and `TestTheSetupManifestIsCurrent` fails." The operator has not reproduced it outside that Spec. Recorded by the operator on 2026-11-10 from the Session Handoff of 2026-11-10.
