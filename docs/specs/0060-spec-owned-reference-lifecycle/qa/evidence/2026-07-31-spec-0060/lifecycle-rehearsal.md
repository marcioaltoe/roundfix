# Lifecycle rehearsal

Build: `00ca18ee7c0fa2bbc31f00b98c41c4208170cf5f`

Scratch checkout: `/tmp/roundfix-qa0060-flow.NlbV9k/repo` on local throwaway
branch `ma/qa-0060-rehearsal`. No commit was created. After execution,
`rtk rm -rf /tmp/roundfix-qa0060-flow.NlbV9k` removed the checkout and
`rtk proxy /bin/test ! -e /tmp/roundfix-qa0060-flow.NlbV9k` exited 0.

## Adoption

The rehearsal used one tracked finding and one throwaway inbox source. It
followed `write-prd` in order: inventory, classification, active/archive owner
search, finding `status: done` and owner-link update, two `rtk git mv` calls,
fixed five-column index creation, repository-wide link rewrite, and gate checks.

Observed evidence:

```text
rtk git status --short
RM docs/findings/2026-07-29-qa-gate-round-economics.md -> docs/specs/9999-qa-reference-lifecycle/references/2026-07-29-qa-gate-round-economics.md
A  docs/specs/9999-qa-reference-lifecycle/references/2099-01-01-qa-reference-note.md

rtk grep -n '^status: done$' <moved-finding>
2:status: done

rtk grep -n 'Spec 9999' <moved-finding>
132:[Spec 9999 — QA reference lifecycle](../_prd.md).
```

Both original-path absence checks and both indexed-current-path checks exited
0. Repository-wide fixed-string searches for each old path, excluding
`_index.md`, exited 1 with no matches. The index retained both old paths only
as provenance.

The moved finding diff changed only lifecycle metadata (`status`, `updated_at`)
and appended the owner route; its observations were unchanged.

## Authoring negative probes

Recreating the finding at its original path produced exit 1 and:

```text
authoring failed: docs/findings/2026-07-29-qa-gate-round-economics.md remains; repeat adoption step 5 (Move)
```

Changing the PRD link to a missing destination produced exit 1 and:

```text
authoring failed: _prd.md link references/missing-2099-01-01-qa-reference-note.md is unresolved; repeat adoption step 7 (Rewrite links)
```

After each repair, the same absence/resolution checks exited 0 from a fresh
shell.

## Archive preconditions and probes

The documented Task and QA commands reported `status: completed` and
`verdict: pass`. The exact `archive-spec` link-destination regular expression
exited through its success branch on the self-contained fixture.

With stale inline and reference-style links injected together, the command
exited 1, printed both offending lines, and printed:

```text
self-containment failed: rewrite each listed link at adoption step 7
```

After repair it exited 0 and reported that the existing prose mention of
`docs/findings/` did not match.

Deleting an indexed current file produced exit 1 and named the missing path
with step-5 guidance. Recreating an original inbox source produced exit 1 and
named the surviving `source` with the one-move step-5 guidance. Both recovery
checks then exited 0.

Adding `qa_override: true` while injecting a stale findings link still produced
exit 1 and named the link. The current Spec 0060 has no index; its explicit
no-index check exited 0 as the new-promotions-only migration boundary.

## Archive portability and shared ownership

After `rtk git mv` moved the whole fixture Spec under `docs/specs/_archived/`,
both index-relative `path` values still resolved. A later secondary Spec linked
the archived primary owner's finding, had no `references/_index.md`, and a
repository file search returned exactly one physical finding basename. The
primary index continued to name owner `9999`.

## Git-history policy block and equivalent evidence

This QA contract prohibits commits. On the staged 98% rename from
`docs/findings/` to the archived Spec reference, `git log --follow` at the
uncommitted destination therefore exited 0 with no commits; at the committed
source path it reached `fe345ab`. This is the same Git precondition recorded in
Task 03's Result.

Equivalent observed evidence for the move mechanism is available without a QA
commit:

- `rtk git diff --cached --summary` recognized the finding as a 98% rename and
  the content diff showed only lifecycle metadata plus the appended owner route.
- `rtk git log --follow` at the committed source reached its pre-adoption
  history.
- On the repository's already committed active-to-archive directory move,
  `rtk git log --follow --oneline --
  docs/specs/_archived/0053-qa-gate-reachability-and-verdict-semantics/_prd.md`
  reached both `d930d3e` and pre-archive commit `397227f`.

The only unreachable observation is `git log --follow` starting at a newly
committed moved-finding destination. Unblocking it requires authority to make a
throwaway adoption commit in a disposable clone.
