# Local workflow probes

Build: `171f6a378c9e640a8a10c9382e28b501b21ff5a0`

- `rtk ruby .../current_workflow_probe.rb` — exit 0. The probe parsed the
  current workflow; exercised the npm 11.5.1 boundary; tag, valid dispatch,
  invalid semver, and mismatched-version validation; all classifier fixtures;
  eligible, used, cooldown, malformed, absent, HTTP 503, transport, and mixed
  preflight runs; all-OIDC, open-fallback, and closed-fallback publication; the
  six-coordinate order; job-scoped permissions; and the review repair's
  expression-free generated shell.
- `rtk ruby .../failure_attribution_probe.rb` — exit 0. A network timeout
  emitted `publish:` with the npm diagnostic, exited 1, and was not retried as
  identity.
- `rtk ruby .../secret_boundary_probe.rb` — exit 1. It observed
  `SECRET_EXPRESSION_SCOPE=publish-step`, one `NODE_AUTH_TOKEN` assignment,
  zero secret expressions in generated shell, and one `$NPM_TOKEN` fallback
  reference. This is QA-005: the credential remains available to the entire
  step even though only the fallback command consumes it.

No sentinel token appeared in captured stdout, stderr, or summaries.
