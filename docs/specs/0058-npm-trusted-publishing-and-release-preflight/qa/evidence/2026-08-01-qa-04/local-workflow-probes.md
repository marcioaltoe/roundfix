# Local workflow probes

Build: `e45dd37d2f2ced6dcaa3533fcea939a867b3ea6c`

## Prior QA-005 reproduction

`rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-08-01-qa-03/secret_boundary_probe.rb` exited 0. The original detector now reports `SECRET_EXPRESSION_SCOPE=other`, one `NODE_AUTH_TOKEN` assignment, zero script secret expressions, and `PASS: retained token is unavailable outside the fallback branch`.

## Dynamic secret-boundary and publish-flow probe

`rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-08-01-qa-04/secret_boundary_probe.rb` exited 0 with every assertion passing.

- YAML resolves the folded `NPM_FALLBACK_TOKEN` environment mapping to the exact GitHub secret expression, while the parsed `run:` script contains no secret expression and the workflow contains one retained-token reference.
- The script copies the exported value into a non-exported shell variable and unsets `NPM_FALLBACK_TOKEN` before the first `npm publish`.
- Six successful OIDC attempts each observed `NODE_AUTH_TOKEN`, `NPM_FALLBACK_TOKEN`, and `npm_fallback_token` absent from their process environment.
- An open-window authentication failure made one token-free OIDC call, one token-bearing retry, then five later token-free OIDC calls. The summary recorded only the coordinate.
- The closed-window case emitted `identity:`, made one token-free call, and performed no retry.
- The network case emitted `publish:`, surfaced the underlying timeout, made one token-free call, and performed no retry.
- The sentinel appeared in no stdout, stderr, job summary, fallback record, or generated temporary artifact.
- The call log preserved all five platform packages before the launcher.

These are command-boundary substitutes for irreversible npm publication; the real OIDC exchange remains a separate live-release row.

## Full workflow probe

`rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-08-01-qa-04/full_workflow_probe.rb` exited 0. It exercised the runtime floor, tag and dispatch validation, fixture classifier, full stubbed preflight response matrix, multi-coordinate accumulation, six-coordinate publication, open/closed fallback, and platform-before-launcher order. It also independently compared the current workflow with `21bc4bf^` and found the `Verify gate`, `Cross-compile and stage`, and `GitHub Release` scripts byte-identical; GitHub Release remains after npm publication.

`rtk git diff --name-only 21bc4bf^..HEAD -- dist/npm cmd internal` exited 0 with no output, confirming the Spec changed no package manifest, command implementation, or Upgrade Command path.
