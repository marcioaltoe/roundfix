# Real tagged release reachability

Read-only `rtk gh run list --workflow release.yml --event push` and
`rtk gh release list` observations show the newest tagged release is `v0.0.2`,
run `30477190935`, created on 2026-07-29 at commit
`755e9fba2c8ca03ec674c0d31d6dc7a056cc18c7`. It predates Task 01 commit
`21bc4bf` on 2026-07-31 and therefore cannot prove this Spec's OIDC path or
fallback summary.

Environment block: QA has no maintainer-authorized fresh version/tag/release,
and creating one is irreversible publication outside this gate's mutation
authority. Unblocking requires all six trusted-publisher bindings plus a new
maintainer-authorized tag whose completed run publishes all six coordinates
with an empty fallback record; registry reads must then confirm the exact
version on all six coordinates.

Equivalent evidence is incomplete. The local publish-boundary harness proves
OIDC-first command shape, fallback behavior, ordering, and non-leakage; live
dispatch `30703974453` proves the exact remote runtime, Verification, registry
classification, and mutation barrier. Neither performs npm's OIDC exchange or
persists an empty fallback record.
