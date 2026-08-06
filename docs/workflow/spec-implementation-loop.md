# Spec implementation loop — supervisor instructions

Working notes for running Specs end to end through Roundfix. Everything here is
supervisor discipline that needs no product change; the product-side gaps live
in [the QA cycle finding](../findings/_archived/2026-07-29-qa-cycle-latency-and-detector-placement.md).
This file is a staging area: adjust it as the loop teaches, and promote the
parts that stabilize into a skill or a command.

## The loop

For each Spec, in backlog order, without pausing between Specs:

1. Sync `main`, branch `ma/<spec-slug>`.
2. Author whatever the Spec is missing — TechSpec, then Task Graph. Invocation
   is the authorization; present the breakdown as a report, never as a
   blocking question.
3. Commit and push the authoring work.
4. `roundfix implement --spec <slug> --qa --detach`, then monitor by polling
   the Console Log for a terminal signature.
5. On QA pass: archive on the branch, push, open the Pull Request.
6. `roundfix watch --source coderabbit --pr <n> --until-clean`.
7. On Clean or the configured round cap: squash merge, delete the branch, sync
   `main`, `roundfix reconcile`, rebuild the binary.
8. Next Spec.

Stop only for a decision that is genuinely the maintainer's: an unforeseen
tooling authorization, or a scope change. Not for approvals the loop already
covers.

## Writing Verifications

**Verify the class, not the case.** A Task that removes a leak from one call
site should assert that no file in that area can leak, not that the one file
stopped leaking. Case assertions pass while the next instance of the same
defect ships.

**Put the detector where the risk is.** If a Task changes a contract with an
external surface — an ACP adapter, a Review Source, a package registry — its
Verification must probe that surface, not only compile and unit-test. The
adapter defect in Spec 0052 was invisible to every fixture in the repository
and visible in thirty seconds to
`roundfix profiles validate --category <c> --json`. A twenty-minute gate is
the wrong place to learn something a thirty-second probe answers.

**Suspect any fixture that models an external response.** A fixture encodes
what you believed the surface returns. The adapter echo defect existed because
every fixture modelled a list no adapter actually sends. When a Task depends on
an external response shape, capture that shape from the real surface first and
build the fixture from the capture.

**Prove the effect, not the absence of regression.** A Verification that only
runs broad suites proves nothing about the Task. Include at least one
executable check that fails before the Task and passes after.

## Writing Specs

**Specify the hazard, not just the objective.** Three of the four Spec 0052
findings and all three chained Vortex failures came from objectives written
without their hazard: "time cap" without "cancel, never abandon"; "technical
detail only in the log" without "the log is a destination someone reads".
Whenever a requirement names where data goes, name who reads that destination.

**Do not dismiss a risk you can cheaply test.** The Spec 0052 TechSpec listed
the exact collision that later broke it, called it hypothetical, and moved on.
If a risk is written down and a probe would settle it, run the probe while
writing the Spec.

**Bound the tooling authorization to the real fallout.** Editing one Skill
cascades into up to seven derived digest files; a Baseline module clause
cascades further, into the source-baseline manifest, formatter goldens, and the
profile golden digest. An authorization that undercounts the fallout blocks its
own Task.

## Recovering from a QA failure

**Read every finding before fixing any of them.** Fixes are cheap; cycles cost
twenty minutes. One commit that closes four findings costs one cycle; four
commits cost four.

**Never touch protected tooling during recovery.** The gate audits the
supervisor exactly as it audits Agents. A comment-only edit to Project Config
cost a full cycle in Spec 0052 because that path was outside the Spec's bounded
authorization. If a fix needs a path the Spec never authorized, either drop the
fix or route it to the next Spec — reordering history to insert an
authorization ahead of the commit it authorizes costs more than it saves.

**Prefer fixing the source of a brittle assertion.** A test pinning a literal
that a Spec legitimately changes should read the constant instead. Bumping the
literal buys one green cycle; reading the constant ends the class.

**Verify a stranded branch before discarding it.** Failed QA cycles leave Run
Branches that Branch Integrity Preflight then refuses. Confirm each branch's
unique content is already present or genuinely superseded before removing it —
the suggested `git merge --ff-only` can be actively wrong, because integrating
a superseded QA report overwrites a passing one at the same dated path.

## Concurrency

`worktree.concurrency` may be raised whenever the Task Graph is genuinely
file-disjoint. Per
[ADR-0026](../adr/0026-task-integration-is-a-serialized-cherry-pick-queue.md)
the serialized cherry-pick is deliberate: a conflict means the graph declared
two Tasks independent when they were not, and the Task settles `failed` with
its Worktree kept for inspection. That is a decomposition defect to fix, not a
reason to serialize every Spec forever. Raise the value; treat conflicts as
graph feedback.

## Reporting

Report at milestones — Run started, terminal outcome, QA verdict, merge — not
per Agent event. Per-event narration consumes context and time without
improving the decision the maintainer has to make.
