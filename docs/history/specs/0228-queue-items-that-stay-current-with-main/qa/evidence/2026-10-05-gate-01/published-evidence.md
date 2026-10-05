# Published evidence, observed 2026-10-05

The local curl-cffi attempt on the Git page could not establish the proxy tunnel (curl 7, CONNECT response 403). It was not retried. The web read tool successfully reached all three primary sources below; no outside-source row is blocked by this transport limit.

## Q09 — Git

Source: [git merge documentation](https://git-scm.com/docs/git-merge), HOW CONFLICTS ARE PRESENTED and HOW TO RESOLVE CONFLICTS, lines 492–544 as retrieved by the web reader.

Observed: “The part before the `=======` is typically your side”. The document identifies the other side after the separator, distinguishes diff3's base side, and describes aborting an unresolved merge. It still supports the resolver's explicit two-sided format and default-branch side selection. This is independent support for the protocol; Q03's actual merge and the negative Git fixtures establish implementation behavior.

## Q10 — Changesets

Source: [Changesets decisions](https://github.com/changesets/changesets/blob/main/docs/decisions.md), How changesets are combined, lines 197–208 as retrieved by the web reader.

Observed: “we flatten the version bumps into one single bump”. The document selects the combined bump when consuming accumulated changes. It still supports combination-time version choice as an analogy. The page explicitly marks itself outdated and points to https://changesets.dev/guide/technical-decisions; the cited combination principle remains on the named source. This is not proof of Roundfix's implementation; Q02/Q03 supply that proof.

## Q11 — Gerrit

Source: [Gerrit review labels](https://gerrit-review.googlesource.com/Documentation/config-labels.html), label.Label-Name.copyCondition and changekind:NO_CODE_CHANGE, lines 182–204 as retrieved by the web reader.

Observed: “only the commit message may be different”. Approval copying depends on an explicit query and the NO_CODE_CHANGE predicate preserves the parent tree and code diff. This still supports constrained reuse as an analogy, not arbitrary documentation changes or automatic approval. Roundfix instead proves a bounded archive correction and requests round 2, independently exercised by Q04's correction/lineage tests.
