# 2026-08-04 — Branch Integrity Preflight prescribes a remedy that reintroduces superseded work

status: pending

## What was observed

A QA gate failed a Spec because the Supervisor had committed protected tooling
(a Baseline profile) onto the Spec branch. The correct remediation, named by the
gate itself, was *"remove the out-of-scope profile change from this Spec's
ancestry"*. The branch was unpushed, so it was rebuilt: the offending file was
dropped, the platform work moved to its own PR, and every Spec commit replayed.

The rebuild left three terminal Run Branches whose commits carried the *old*
ancestry. Their content survived in the rebuilt branch under new SHAs, so Git
ancestry reported them unintegrated. `roundfix watch` then refused to start:

```
Branch Integrity Preflight refused pending Run Branch work for PR Head Branch "…".
- branch=roundfix/run-…33453c4b ahead_commits=3  integration_command="git merge --ff-only roundfix/run-…33453c4b"
- branch=roundfix/run-…d0934d88 ahead_commits=10 integration_command="git merge --ff-only roundfix/run-…d0934d88"
- branch=roundfix/run-…74adf985 ahead_commits=1  integration_command="git merge --ff-only roundfix/run-…74adf985"
Next action: inspect each pending Run Worktree, then run the listed integration command …
```

**Running the prescribed command would have undone the remediation.** Those
branches contain the very commit the gate demanded be removed (`docs: alinhar o
contrato de backend com o design do Axis`, already landed separately as its own
PR) plus two superseded `qa report … (fail)` commits that a later `(pass)`
report replaced. The printed next action pointed straight back at the state the
gate had just rejected.

`roundfix reconcile` could not help: it classified the three as `unintegrated`
and `dirty`, action `preserve`, with no force path — correct conservatism, but it
left no supported way forward. The only exit was `--skip-branch-integrity`,
which is documented as a bypass for *ignoring* guardrails, not as the normal
route out of a legitimate remediation.

Verification that nothing would be lost had to be done by hand, comparing commit
subjects between each stranded head and the rebuilt branch:

```
### commits nas branches presas          ### commits na minha branch
docs: autorar o portao de QA …           docs: autorar o portao de QA …
feat: falhar alto em ambiente …          feat: falhar alto em ambiente …
docs: alinhar o contrato … com o Axis    (deliberadamente removido → PR próprio)
docs: qa report … (fail)                 docs: qa report … (pass)
…                                        + 3 commits novos
```

## Root cause

Branch Integrity Preflight decides "pending work" from Git ancestry alone. It
cannot distinguish:

- a Run Branch carrying work the target branch genuinely lacks, from
- a Run Branch whose content the target branch already has under different SHAs
  because history was legitimately rebuilt.

Both look identical to `merge-base --is-ancestor`. The command then prints a
fast-forward as the *next action* with no hedge, even though it is a mutation
whose correctness the tool has not established. Reconcile shares the blind spot
from the other side: it preserves, but never recognises supersession by content.

This matters more for an autonomous loop than for a human. A Supervisor agent
following the printed `Next action:` — which is precisely what a deterministic
next-action field invites — would silently reintroduce the work a QA gate had
just ordered removed, and the branch would fail the same gate again with no
indication why. The instruction to *"verify a stranded Run Branch before
discarding it"* exists in the skill; the tool's own output argues against it.

## What would settle it

- Soften the prescription. When a pending branch's ancestry diverges from the
  target rather than trailing it, print the branches and stop at "inspect", not
  at a `git merge --ff-only` next action. A remedy the tool cannot prove safe
  should not be phrased as the deterministic next step.
- Give `reconcile` a supersession check: compare patch-ids or commit subjects
  between the stranded head and the target. Content already present under a
  different SHA is releasable, and saying so turns a dead end into a supported
  path.
- Recognise history rebuild as a first-class outcome of a failed governance
  gate. The gate can order a commit removed from a Spec's ancestry; when it
  does, the branch-integrity machinery should expect rewritten history rather
  than treat it as pending work.
- Failing all that, document `--skip-branch-integrity` as the sanctioned exit
  for this specific case, with the manual verification it requires. Today it
  reads as a guardrail bypass, so using it correctly looks like cutting a
  corner.

## Evidence

- Consumer repository: oraculum, PR #38, branch
  `ma/spec-0015-health-identidade-operacional`.
- Gate finding that ordered the removal: F-001 in
  `docs/specs/_archived/0015-health-e-identidade-operacional/qa/qa-report-2026-08-04.md`.
- Platform work relanded separately as oraculum PR #37.
- Related: [a Spec cycle leaves branches and worktrees nobody audits](2026-08-02-a-spec-cycle-leaves-branches-and-worktrees-nobody-audits.md).
