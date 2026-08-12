# R03 — two-tier verification

- The first cold `make verify` attempt was terminated by the managed network
  sandbox when an integration path reached `api.github.com`; this was an
  environment boundary, not a product verdict.
- The authorized identical rerun, after `go clean -testcache`, exited 2 after
  119.23 seconds. `internal/agent`, `internal/app`, `internal/baseline`, and
  `internal/baselineacp` passed; `internal/cli` failed its QA verdict, auto-push,
  attach, external-spec-root, QA-only, target-branch, and interactive-input
  journeys.
- `make verify-incremental` exited 2 after 57.44 seconds with the same CLI
  failure class. It met the under-60-second timing bound but did not validate
  the current change.
- The Makefile comparison shows the original `verify:` declaration is
  byte-identical and `verify-incremental:` is additive. The complete target
  still accepts `VERIFY_TEST_TARGET=test-budget`; the local target fixes the
  cached `test` tier.

Both tiers are red on this build, so the success metric fails.
