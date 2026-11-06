---
type: fix
status: open
created: 2026-11-06
spec: null
reason: null
---

# The repository build cache stays within a size bound

On 2026-11-05 at about 19:07 BRT the Mac mini's system volume filled (161 MiB free). The Delivery Queue parked three items with `delivery-error: ... build item binary: verification command "make build" failed` (`compile: writing output: ... no space left on device`), and a fourth Task ended `Verification unknown`. Two Go build caches held most of the space: the user cache `~/Library/Caches/go-build` (28 GB, cleared with `go clean -cache` to recover) and the Makefile's `GOCACHE`, `$TMPDIR/roundfix-go-build-<uid>` (25 GB then, 26 GB the next morning), which every Run's Verification on this repository shares and nothing trims beyond Go's own five-day eviction.

Intent: a busy day of Runs cannot fill the disk through the build cache. Candidates to evaluate: a size bound for the Makefile cache that a target enforces before building (trim or clear above it), a doctor or Delivery Queue warning when free space on the cache's volume drops under a floor, and a park reason that names a full disk instead of a generic `delivery-error`.

Evidence: `/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/delivery/0354-an-external-skill-is-checked-only-for-presence/item-binary-build.log`; `df -h /` before and after `go clean -cache`; `du -sh $TMPDIR/roundfix-go-build-501`.
