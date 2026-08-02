---
status: open
created_at: 2026-07-30
updated_at: 2026-07-30
---

# Failed QA runs accumulate Run Branches nothing can release (2026-07-30)

Spec 0053 gives `roundfix reconcile` a `superseded` classification so a Run
Branch holding nothing but an obsolete QA report stops being preserved forever.
Implementing that Spec produced the case its design does not cover: four QA
runs, none passing, and four Run Branches that neither `reconcile --apply` nor
Branch Integrity Preflight can clear.

## Symptom / evidence

After four QA cycles on Spec 0053, `roundfix watch --source coderabbit --pr 51`
refused:

```text
Branch Integrity Preflight refused pending Run Branch work for PR Head Branch
"ma/spec-0053-qa-gate-reachability".
- branch=roundfix/run-run_20260730T141915Z_93576ce770ae59e2 ahead_commits=1 …
- branch=roundfix/run-run_20260730T144822Z_084a2f73f6a0c7f8 ahead_commits=1 …
- branch=roundfix/run-run_20260730T152956Z_3f8ddc9020907985 ahead_commits=1 …
- branch=roundfix/run-run_20260730T160806Z_2a8bab139ef9089e ahead_commits=1 …
```

`roundfix reconcile` classified all four `unintegrated`:

```text
Summary: total=135 safe=0 superseded=0 unintegrated=4 dirty=0 unknown=0
         released=131 applied=0 preserved=4 operational-failures=0
```

Each branch holds exactly one commit: a QA report with a failing verdict.

## Root cause

Supersession requires a **newer QA report on the target branch**, and a QA
report reaches the target branch only when the Daemon integrates a passing Run.
A failing Run settles Unresolved and leaves its report on the Run Branch. So
the release path is:

- one failing report, then a passing one → the failing branch is superseded and
  releasable. This is the case Spec 0053 designed for and it works.
- N failing reports, none passing → no report ever reaches the target, nothing
  is superseded, and every branch is preserved indefinitely.

The second case is the normal shape of a Spec under active QA repair, which is
exactly when the operator most needs the watch to run.

The preflight's fallback does not help either. Its printed
`git merge --ff-only <branch>` stops working as soon as the target branch gains
a commit after the Run — which corrective work always does — so by the time the
branches are noticed, the suggested command fails too.

## Action / suggestion

Give reconcile a classification for a QA-report-only branch whose report is
**superseded by a newer report anywhere for that Spec**, not only by one on the
target branch. The Run Database already records each Run's Spec and QA verdict,
and the QA-report-only probe already proves the branch carries nothing else, so
the evidence exists without a schema change.

A narrower alternative: when several QA-report-only branches exist for one
Spec, every branch except the one holding the newest report is superseded by a
sibling. That covers the observed case using only Git evidence and keeps
ADR-0053's proof-based rule.

Whichever shape is chosen, the printed next action must be a command that can
actually run. Today the operator is told to fast-forward a branch that cannot
fast-forward, and the only remaining exits are `--skip-branch-integrity` or
deleting branches by hand outside the tool.

## What worked — keep

- The refusal itself is right. Four branches holding unmerged QA evidence is
  exactly the state a user should be told about before a Run touches the
  branch.
- The classification is honest: with no newer report on the target, these
  branches genuinely are `unintegrated`, and the fix that made preflight agree
  with the classifier (Spec 0053) is what made this gap visible instead of
  letting preflight offer a release the classifier would refuse.

## Routing — 2026-08-01

Routed to [Spec 0066](../specs/0066-run-teardown-reclaims-what-it-created/_prd.md) on 2026-08-01.
