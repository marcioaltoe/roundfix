# Live tagged release

Fresh `gh run list --workflow release.yml` and `gh release list` inspection found no post-implementation tagged release.

- The newest tagged release run is `30477190935`, commit `755e9fba2c8ca03ec674c0d31d6dc7a056cc18c7`, tag `v0.0.2`, completed successfully on 2026-07-29.
- The first implementation Task is commit `21bc4bf1e2d4659556795ac2f59a7867b23fbdd2` on 2026-07-31, so run `30477190935` predates the OIDC/preflight implementation.
- The only published GitHub Release is `v0.0.2`, also dated 2026-07-29.

Blocked cause: no maintainer-authorized fresh version/tag/release exists after the implementation. Unblocking requires all six live trusted-publisher bindings plus an authorized new tag whose completed workflow shows six OIDC publications and a persisted empty fallback record.

Local command-boundary probes and the live publish-free dispatch prove the surrounding Release Set, ordering, failure, and mutation barriers, but they are not equivalent proof of a real GitHub OIDC exchange or an empty persisted fallback record.
