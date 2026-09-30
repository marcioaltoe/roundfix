---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A conflict confined to declared derived paths is resolved by regeneration

On 2026-09-30 Pull Request #298 reported `CONFLICTING` after two other items
merged. The queue owner waited for checks that GitHub never starts on a
conflicting Pull Request and parked the item with `checks-timeout`. Every
conflicted file but one was derived: a digest pin, catalog snapshots, plan
goldens and the Setup Manifest, all rewritten by a sanctioned regeneration
command. The one source conflict was a Baseline profile, and resolving it with
one side's bytes dropped another Spec's decision.

Now the owner reads GitHub's `mergeable` state while it waits for checks, and a
`CONFLICTING` Pull Request stops the wait at once. The owner merges the
refreshed default branch into the item branch, never rebases it. When every
conflicted path matches a derived-path declaration in Project Config, the owner
takes the default branch's bytes for those paths, runs each matched
declaration's regeneration command, and commits the merge only when the
regeneration changed no path outside the declarations. The new head then
passes the repository gate, is pushed and is checked again. Any other
conflicted path aborts the merge, and the item parks with the list of source
paths.

The owner does not repeat the Pre-PR Review for that merge. The merge adds only
commits already on the default branch and outputs of a declared command, and
the reviewed diff of the item's own change does not move. A merge the operator
resolves by hand touches source, so a Delivery Retry sends it through the
review again.

## Consequences

The declaration lives in Project Config because merge recovery is a repository
policy, like the Verification and bootstrap commands the owner already runs
from there. It grants no authority: a declared command runs only inside the
owner's merge, and the audit of a Task commit still reads the grant under
[ADR-0149](0149-one-regeneration-declaration-the-grant-names-the-command-the-tree-names-outputs.md)
and [ADR-0178](0178-a-task-commit-is-authorized-by-the-grant-it-ran-under.md).
A path declared derived that also carries hand-edited source is resolved
wrongly; the repository gate that follows is the net, and a failed gate parks
the item with the merge still unpushed. A repository without a declaration
parks every conflict, as a source conflict.

Rebasing was rejected. It rewrites a pushed branch, needs a forced push and
discards the head the Pull Request, the review and the checks were recorded
against.
