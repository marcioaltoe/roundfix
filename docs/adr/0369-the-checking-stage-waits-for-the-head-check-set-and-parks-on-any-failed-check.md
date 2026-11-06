---
status: accepted
created_at: 2026-11-06T12:00:00Z
updated_at: 2026-11-06T12:00:00Z
deprecated_at: null
superseded_by: null
---

# The checking stage waits for the Head Check Set and parks on any failed check

A Delivery Queue item in `checking` moves to `merging` only when every member
of its Head Check Set is completed and none failed: every check run and commit
status the pull request head lists, the latest GitHub Actions workflow run per
workflow file and event of the head commit, and every Expected Workflow of the
head's own workflow files. A merge state of `CLEAN` alone never moves it.

## Context

On 2026-11-05 the Delivery Queue of Fluxus merged pull requests #466 and #469
with the `CI | Backend` job about to fail, and #477 before its CI ran. The
queue read every check `gh pr checks` listed, required or not, and parked on a
listed failure, but the three workflow runs of each head existed seconds after
the pull request was opened while their job check runs appeared more than a
minute later; the list held one passing check, and the merge state was
`CLEAN` because the base branch requires no check. The maintainer chose:
"na etapa de checks, a fila espera todos os workflows que o head disparou
terminarem, até o prazo de sempre. Se algum falhar, mesmo sem ser obrigatório,
o item para com o nome do check e o link."

## Decision

**The set.** At each poll the queue reads the head's check runs and commit
statuses as before, the workflow runs of the head commit through `gh api
repos/{owner}/{repo}/actions/runs?head_sha=<head>`, keeping the latest run per
workflow file and event, and the head's workflow files through local `git`.
A check run that belongs to a listed run counts once. Check suites are not
read: a suite with no check run never completes.

**Expected Workflows close the first seconds.** A workflow file whose trigger
makes GitHub run it for every pull request into the base branch is expected:
`pull_request` or `pull_request_target` without a path filter or
`branches-ignore`, whose `branches` names the base branch literally without
a `!` pattern and whose `types` name `opened` and `synchronize`, when present;
or a `push` trigger with no filter. Until such a file has a run for the head,
the set is not complete. A filtered workflow is never expected, but its run is
a member once GitHub creates it.

**Conclusions.** `success`, `neutral` and `skipped` pass. Any other
conclusion fails: `failure`, `timed_out`, `cancelled`, `action_required`,
which never completes without a person, `startup_failure` and `stale`, which
did not test the head. An empty set keeps waiting.

**Parks keep their names.** A failed check still goes through the stale and
outside-change re-runs of ADR-0236, and a run the queue re-ran is pending
until a newer attempt reports. Every other failure parks `checks-failed:
<check> (<conclusion>) <link>`, each failed member joined by `; `. The deadline
stays `delivery.check_timeout` and parks `checks-timeout: current-head checks
not settled within delivery.check_timeout <duration>; pending: <check>;
<check>`. Both keep their Park Classes.

**The wait is visible.** The console log names the pending members when they
change and the member count when the set completes. The Queue Owner records
the pending members on the item's Queue Turn, and `roundfix deliver status`
prints `Turn: <slug> advancing at checking; waiting for <n> check(s): <check>;
<check>`.

## Considered Options

- A settle delay before trusting the list was rejected: it adds a fixed wait
  to every merge and still guesses how long GitHub takes.
- Evaluating `paths`, branch globs and `branches-ignore` to expect filtered
  workflows was rejected: it re-implements GitHub's filter semantics, and a
  wrong guess waits for a run that never comes until the Check Wait ends.
- Waiting for every check suite was rejected: an app can create a suite that
  never gets a check run.
- Requiring only the branch's required checks, or warning about failed
  optional checks while still merging, was rejected by the maintainer's
  answer.
- Reading `UNSTABLE` from the merge state was rejected: GitHub sets it only
  for checks that already reported.

## Consequences

A pull request now merges no earlier than every workflow its head triggered,
which costs nothing when CI already takes minutes. A repository whose only
workflows are filtered still relies on GitHub creating their runs within the
first poll interval. An operator whose `gh` token cannot read Actions sees a
logged read error and a `checks-timeout` park. A blocker recorded before this
change keeps its bare text; status and retry read both forms.
