---
status: accepted
created_at: 2026-07-05T22:17:04Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# A parsed prompt result outranks the acpx exit code

When the acpx NDJSON stream has already delivered a valid `session/prompt` response for a Batch, a subsequent nonzero acpx exit is classified as teardown noise: journaled loudly with the stderr tail, while the Batch proceeds to the Daemon's verbatim verification — which remains the only gate for settling and committing (ADR-0014). Without a parsed result, a nonzero exit stays a Batch failure exactly as before. Motivated by two dogfood Runs in one day where acpx's 10 MiB message buffer killed finished turns and a false failure discarded completed, verifiable work.
