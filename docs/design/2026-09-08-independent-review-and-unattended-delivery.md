# Independent review and unattended delivery

Status: In review, 2026-09-08. Decider: maintainer.

Remove the CodeRabbit execution path, require an independent local review before
publishing a Pull Request, and make delivery survive the supervising session.
Keep `implement-spec` as the preparation and orchestration entry point, backed by
Roundfix Runs. The proposed sequence is Spec 0126 for removal and independent
review, followed by Spec 0127 for durable delivery. These numbers are provisional;
this design does not approve their architecture or authorize tooling edits.
The wider queue and its provisional assignments live in the
[pending-work plan](../workflow/2026-09-08-pending-work-plan.md).

The maintainer has approved automatic delivery through merge after the relevant
Spec and its limits are approved, with independent review and required checks
passing for the current head. Reviewer selection, correction budgets, and exact
tooling grants remain decisions to record in the Specs. Existing authorization
does not turn an unanswered limit into an unlimited allowance.

## Problem and observed state

The current source accepts only `review_source.name: coderabbit` in
`internal/config/config.go` and `internal/cli/cli.go`. Omitting the block restores
that default. There is no supported `none` value or local review command.
`resolve` consumes downloaded Review Issues and requires an Open Pull Request;
it cannot ingest findings from a review performed before a PR exists.

The project config sets `request_review: true`, while `.coderabbit.yaml` disables
automatic review and retains its description marker. Setting only
`request_review: false` makes `watch` and `resolve` fail their coherence
preflight. It does not establish a supported disabled mode. Deleting the
CodeRabbit file first is also unsuitable: missing settings use provider defaults.

GitHub was inspected without mutation on 2026-09-08. The repository has no
rulesets; branch protection on `main` requires `Verification gate` and
`Validate PR title`. CodeRabbit is not a required check. The inspection did not
establish the GitHub App installation's repository access or disable that app.

```bash
rtk proxy gh api repos/marcioaltoe/roundfix/rulesets
rtk proxy gh api repos/marcioaltoe/roundfix/branches/main/protection
```

There is also an instruction conflict. The
[2026-08-31 decision](../workflow/2026-08-31-kickoff-is-settled-implement-spec-remains-the-loop.md)
keeps `implement-spec`, but that skill still tells the session to implement,
verify, and commit Tasks itself, invoke QA after the graph, and stop before
publication. The [autonomous guide](../agents/autonomous-work.md) and the
`Autonomous Spec delivery` section of the Roundfix skill assign implementation
and Verification to the Daemon, make QA a terminal Task, and describe delivery
through merge. Preserve the accepted ownership, then repair the conflicting skill.

## Goals and boundaries

The resulting delivery path must require no CodeRabbit service, account,
configuration, request marker, API call, or success signal. Review must happen
before the Pull Request opens and must cover the commit eligible for merge.
Read-only history remains readable, including old Runs, Review Issues, and
CodeRabbit evidence; historical provenance does not enable an execution path.

This proposal does not install a new orchestration framework, rename
`implement-spec`, copy Fluxus's kickoff skill, change CI quality thresholds,
uninstall an organization-wide GitHub App, or purchase API credits. The
Supervisor continues to author documents and dispatch Runs; it does not write
feature code or tests.

## Active dependencies and history

| Paths | Role now | Proposed treatment |
| --- | --- | --- |
| `.roundfixrc.yml`, `.coderabbit.yaml` | Active repository configuration | Remove obsolete review configuration atomically with config migration; handle repository App access separately if confirmed necessary. |
| `internal/config/config.go`, `internal/preflight/preflight.go`, `internal/cli/cli.go` | Source default, coherence preflight, command routing, adapter construction | Retire CodeRabbit execution and make legacy invocations return a migration diagnostic before side effects. |
| `internal/reviewsource/coderabbit/coderabbit.go` | GitHub requests, collection, classification, replies, resolution | Remove executable adapter; preserve historical data decoding elsewhere where needed. |
| `internal/reviewsource/reviewsource.go`, `internal/rounds/rounds.go`, `internal/watch/watch.go` | Shared review types and PR-resolution loop | Keep only proven reusable types and history readers; remove provider requests and the operational dependency on `watch`. |
| `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/modules/autonomous-work.json` | Canonical instructions requiring a PR review request and `watch` | Replace with the independent-review and delivery contract. |
| `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/agents/openai.yaml`, `.agents/skills/implement-spec/SKILL.md` | Roundfix-owned operational instructions | Describe the new commands and one delegation loop. |
| `skills/roundfix/SKILL.md`, `skills/roundfix/agents/openai.yaml`, `skills/implement-spec/SKILL.md`, `docs/agents/agent-instructions.md`, `docs/agents/autonomous-work.md` | Embedded or rendered copies | Regenerate from their owners; never use copies as the source. |
| `skills/skills.go`, `skills/baseline_skill_contract_test.go`, `internal/docscontract/publicdocs_test.go` | Checks currently require CodeRabbit phrases and configuration | Change the intended contract and prove its replacement in the same slice. |
| `README.md`, `docs/user-guide/commands.md`, `docs/user-guide/configuration.md`, `docs/user-guide/usage.md`, `CONTEXT.md` | Active public and domain contracts | Document retirement, independent review, and preserved history. |
| `docs/history`, old Run journals and review artifacts, `CHANGELOG.md`, dated authorization records | Historical evidence | Preserve provenance; do not rewrite old verdicts or delete stored Runs. |
| `internal/baseline/assets/source-baselines`, `internal/baseline/assets/formatter-fixtures`, `internal/baseline/assets/setups`, `internal/baseline/testdata`, `internal/baseline/assets/profiles` | Maintained derived artifacts and characterization material | Change only through the sanctioned regeneration steps and their declared ownership. A historical occurrence is not an active integration. |
| `.agents/skills/review/SKILL.md`, `.agents/skills/github-pr-workflow/SKILL.md` | Upstream-managed skills | Do not modify locally; place repository policy in the owned modules and skills. |

The last restriction follows `Makefile`'s `OWNED_SKILLS` and `skills/skills.go`.
Neither `review` nor `github-pr-workflow` belongs to the 14 embedded authorial
skills. Adding ownership would require a separate decision.
The supervising session's skill scan found direct CodeRabbit references only in
the owned Roundfix skill, not in these upstream skills. Leaving upstream skills
unchanged therefore does not preserve an identified active CodeRabbit dependency.

## Proposed Spec 0126: removal and independent review

### Build order and safe bootstrap

1. Author and approve the PRD, TechSpec, limits, and exact tooling grants. Keep
   the grant commit before any governed mutation. Characterize the installed
   review CLIs and capture legacy Run/history fixtures before changing behavior.
2. Retire CodeRabbit dispatch, config defaults, request generation, and mandatory
   skill instructions. Migrate `.roundfixrc.yml` in the same Task slice that
   changes the loader. Preserve history readers and explicit migration errors.
   This step depends on the approved contract in step 1.
3. Implement independent local review and evidence validation, reusing the
   existing command runner, process lifecycle, event stream, and model selection
   where their contracts fit. This depends on the historical/operational split
   from step 2. Do not make a local review pretend to have a PR number.
4. Integrate review findings with bounded corrective Specs, re-review after
   corrections, and prove current-head enforcement through the public CLI.
   This depends on step 3; it may not relax a Task's Verification, amend an
   archived Spec, or infer a grant for the corrective Spec.
5. Execute the terminal authored QA Task, archive, independently review the final
   candidate, then publish only when the approved delivery requirements hold.
   This depends on every preceding slice.

The removal slice belongs in `roundfix implement --spec <slug> --detach` because
Implement Runs do not require CodeRabbit review. The maintainer's 2026-09-08
branch-policy decision permits the current `roundfix/run-` namespace and its
Task branches, so the previously identified personal-prefix bootstrap conflict
is resolved. New work branches use purpose prefixes without user initials.
No branch rename behind a Run's owner or Supervisor feature-code exception is
needed. Spec 0125's remaining repository-identity work is still relevant to
cross-worktree orchestration, but naming no longer blocks an implementation Run.
After the applicable Spec, limits, and tooling grant are approved, the old Daemon
can execute and verify the Task that removes the adapter from the next binary.
Rebuild that binary before a later Task or QA exercises its CLI.
Until the native replacement passes its gate, no unattended merge is eligible.
For the first replacement PR, use the independently selected native CLI against
the final candidate and record its evidence; this avoids requiring the new
feature to certify its own existence. That bootstrap still requires the pending
reviewer choice and execution limits.

### Archive and review ordering

Proposed sequence: implement the Task Graph and terminal QA, archive and commit
the completed Spec, freeze the resulting candidate, independently review that
candidate, publish it, verify the current-head checks, and merge. This is the
smallest sequence that preserves the existing archive boundary and obtains a
review of the exact commit to publish. No archive or receipt commit follows the
passing review; receipts stay outside the reviewed tree.

An archive changes the candidate, so a review performed before archive cannot
clear the archived head. This proposal adopts no artifact-only inheritance
exemption. An optional earlier review may help discovery, but would cost an
additional review and cannot replace the required final-candidate review.

When the final review finds a blocking defect, stop publication and preserve the
completed/archived Spec byte-identically. Capture the finding and author a new
corrective Spec with its own approval record, bounded files, Tasks, and QA. After
that correction is implemented and archived, review the entire new candidate
against its actual base. The original Spec's approval does not silently grant
the corrective Spec's implementation or spending. If the approved queue does not
cover that new authority, its execution waits for the maintainer.

Proposed trade-off: preserving archive immutability and exact-head evidence costs
an extra Spec and gate when a late review finds a defect; it can stop unattended
delivery for a new grant. The maintainer has not yet accepted that trade-off.
Changing the archive boundary or approving scoped authority for future corrective
Specs would be a separate explicit policy decision, not a workaround inside
0126 or 0127. QA corrections discovered while a Spec is still active retain the
existing Task Graph correction rules.

The smallest measure available before implementation is to stop invoking
CodeRabbit-bound commands and markers and keep automatic review disabled. This
is an operational pause, not product removal: executable requests, their default
configuration, and mandatory instructions still exist until the removal Task
lands. Changing `request_review` to false can also suppress Roundfix requests but
leaves expected preflight refusals; it must not be described as a successful
disabled configuration. A repository App suspension would be a separate external
action, not evidence that the product has been decoupled.

After step 2, operational shutdown is deterministic: legacy commands refuse
before network/process/Run creation, no config default selects CodeRabbit, and
no owned instruction requests it. Do not introduce a temporary `source: none`
feature merely to retire the old provider. The proposed migration removes the
obsolete config block and detects explicit legacy usage with an actionable
diagnostic. Removing `fetch`, `resolve`, and `watch` as operational commands is
a public-contract change that the Spec must approve and document; history
inspection remains supported.
The preflight must also detect inherited legacy User Config and name its source
and remediation. A project Task does not silently edit `~/.roundfix/config.yml`
or another repository to make the migration pass.

### Independent reviewer contract

Proposed input: repository identity, immutable base and head SHAs, the actual
merge base, Spec path, operative standards paths, changed-path inventory, chosen
runtime/model/effort, and approved resource limits. The owner records this before
the reviewer starts. A branch name alone is insufficient evidence identity.

Proposed output: schema version, review identity, observed input SHAs,
runtime/model/effort, start/end times, coverage and omissions, findings with file
and line evidence, runtime outcome, and the path/digest of the raw result.
The harness validates shape, identity, coverage, and successful completion before
classifying the report. Process exit zero alone means the command completed;
it does not mean the diff is acceptable.

Run the reviewer in a fresh session with no inherited implementer conversation.
The reviewer reads code and the Spec but cannot mutate either, run corrective
work, change tests or Verification, commit, push, or decide its own permissions.
Check the checkout and reviewed identity before and after the invocation.
Repository instructions are review inputs, not authority to expand capabilities.

A new head, changed base/merge base, missing coverage, malformed output,
unavailable model, tool denial, timeout, or interrupted stream invalidates the
review. Preserve the diagnostic and raw evidence. Store review receipts outside
the reviewed working tree before publication so recording a receipt does not
itself create a new unreviewed head. Any later report commit needs another
review. A descendant exception is outside this proposal; the old provider's
historical exception supplies no authority for the new independent review.

The existing `agent.Runner`, `SessionRefForReview`, selection proof, process
cancellation, and Run Event Journal are useful seams. Native CLI review does
not automatically satisfy the existing ACP selection proof: the implementation
must record and verify the model selected by the native invocation rather than
claiming an ACP probe proves a separate executable's state.

### Native interfaces actually observed

Read-only CLI inspection on 2026-09-08 found `codex-cli 0.153.4` and
`Claude Code 2.1.263`. No model invocation was made for this design.

| Interface | Verified surface | Consequence |
| --- | --- | --- |
| `codex review --help` | `--base`, `--commit`, `--uncommitted`, positional prompt | Native target modes exist; the prompt is an alternative target mode. |
| `codex exec review --help` | Those targets plus `--model`, `--json`, `--output-schema`, `--output-last-message`, `--ignore-user-config`, `--ignore-rules`, `--ephemeral` | Prefer characterization of native review with structured output before designing a generic prompt replacement. Availability in help does not prove the resulting payload. |
| `codex exec review --base main 'probe incompatible prompt'` | Parser returned exit 2 before model work | Do not combine a native target with custom review instructions or invent `--instructions`. |
| `claude --help` | `-p`, `--output-format`, `--json-schema`, `--tools`, `--allowedTools`, `--permission-mode`, `--model`, `--max-budget-usd` | A headless review can declare its output and capabilities; validate the returned `structured_output`. |

The Codex native reviewer is documented to produce prioritized findings without
editing the working tree. Its native targeting and custom-context restriction
must be reflected in the Spec: prove how the reviewer reads the adopted Spec
and standards, or use the documented custom mode with explicit immutable scope.
Do not silently fall back to a generic coding prompt. [OpenAI code review](https://learn.chatgpt.com/docs/code-review).

For Claude, `dontAsk` denies actions needing permission, while `--allowedTools`
auto-approves named tools; neither statement alone proves the absence of
inherited write permissions. Propose a minimal `--tools Read,Glob,Grep` surface
with a harness-supplied diff and prove effective permissions, including hooks and
MCP configuration. Do not add unrestricted Bash. `--bare` is unsuitable as a
casual isolation switch here: the documented mode does not use subscription
OAuth credentials. `--max-budget-usd` is a CLI API-spend control, not proof of a
subscription quota ceiling. [Claude headless execution](https://code.claude.com/docs/en/headless),
[CLI reference](https://code.claude.com/docs/en/cli-reference).

Proposed Roundfix command naming and flags must be settled in the TechSpec.
`roundfix review` is a candidate name, not an existing command. No examples in
this document require a nonexistent Roundfix flag.

### Findings and bounds awaiting a decision

Proposed disposition is to block publication on unresolved correctness,
security, contract, or missing-coverage findings; require evidence for rejecting
a false positive; and record optional suggestions in the backlog.
The reviewer cannot mark its own findings resolved by changing the requirement.
The existing [review Agent backlog item](../backlog/2026-08-31-the-review-agent-rewrites-the-contract-it-was-asked-to-satisfy.md)
belongs in this contract, including escalation when a fix needs absent infrastructure.

Rejecting a finding because its premise is disproved is an evidenced disposition.
Accepting a real risk or waiving a failed gate is a different human decision; no
automatic `accepted` or `waived` state can turn failure into approval.

The current autonomous guide already calls for re-examining a decomposition that
generates more than two corrective Tasks from QA findings. Preserve that rule.
Proposed additional limit: at most two corrective implementation/re-review cycles
per candidate, then stop with the remaining findings. This is a proposal, not an
approved budget or a reinterpretation of the QA Task limit. The maintainer must
decide runtime/model selection, wall-clock
and spend/quota limits, cancellation behavior, and whether any severity is
eligible for evidence-backed deferral. No paid-API fallback is implied.
After archive, each corrective implementation cycle belongs to a new corrective
Spec under the ordering above. A cycle budget alone does not grant that Spec.

## Proposed Spec 0127: durable unattended delivery

This Spec depends on 0126's review receipts and refusal contract. Rewrite
`implement-spec` to prepare/triage a durable queue, author missing artifacts,
record approvals, and delegate execution through Roundfix. Port useful behavior
from the Fluxus kickoff requirements without adopting its skill or script.

Use the existing `roundfix window set/show/clear` contract. A Run Window bounds
new Run starts; `budget.max_run_duration` bounds an individual Run. Neither is
currently a queue-wide spend limit. Existing `worktree.concurrency` bounds Task
Worktrees, not concurrent whole-Spec deliveries. Propose serial Spec delivery
initially, with Task parallelism remaining the declared repository setting.

Persist the current Spec, Run, branch, base/head, archive state, review receipt,
PR number, confirmed remote head, merge result, approved limits, and blocking
reason in daemon-owned state. The existing Run Database and event infrastructure
are the preferred home; exact schema changes belong in the TechSpec. One owner
must reconcile each external action before replaying it after restart. A queue
Markdown file can be a human projection, but cannot be the only proof a push or
merge happened.

The native supervisor must survive the initiating chat turn and process. It
starts detached Runs, consumes `roundfix events <run-id> --follow`, reconciles
their terminal outcomes, and advances through archive, independent review,
publication, checks, merge, and the next approved Spec. It revalidates later
Specs after earlier ones change their assumptions. A Run surviving while its
chat dies proves only Run durability; acceptance must also kill the supervisor
between Runs and recover the delivery queue.
Follow the proposed archive-first order above. A blocking post-archive review
parks publication and creates a corrective-Spec decision; it never reopens the
archived graph. Approval for ordinary queue continuation does not erase this
authority boundary.

Before merge, confirm that required checks and review cover the current head,
that the approved base relationship still holds, and that no unresolved blocking
finding exists. A pending user decision parks that item durably. Ask only one
question at a time; never import kickoff's grouped-authorization prompt. The
user's conditional merge authority is persisted with its limits, not inferred
from a skill invocation.

The terminal QA remains a Task under [ADR-0091](../adr/0091-the-qa-gate-is-a-task-node-of-its-own-type.md).
The Daemon retains Verification under [ADR-0014](../adr/0014-daemon-runs-task-verification-and-settles-status.md).
Update the domain deliberately: current `Review Source` and `Merge-Ready`
definitions assume an Open Pull Request. Preserve historical meaning while
recording how independent review changes the active readiness decision; do not
claim [ADR-0142](../adr/0142-head-bound-review-source-evidence-decides-the-watch-outcome.md)
already approves that change.

## Proposed tooling authorization boundary

These are exact candidate paths for separate approval, not a grant. Ordinary
source and test changes still belong in the Task scope; the table concentrates
on protected files and historically governed paths likely to be touched. Split
0126 and 0127 grants and list each Task's subset in the actual Spec.

| Work | Canonical or directly edited protected paths |
| --- | --- |
| 0126 configuration retirement | `.roundfixrc.yml`; `.coderabbit.yaml` |
| 0126 review obligations and public skill | `internal/baseline/assets/modules/core.json`; `internal/baseline/assets/modules/autonomous-work.json`; `.agents/skills/roundfix/SKILL.md`; `.agents/skills/roundfix/agents/openai.yaml` |
| 0126 affected historically governed contracts | `internal/cli/cli_test.go`; `internal/docscontract/publicdocs_test.go`; `skills/baseline_skill_contract_test.go` |
| 0127 orchestration instructions | `.agents/skills/implement-spec/SKILL.md`; `.agents/skills/roundfix/SKILL.md`; `.agents/skills/roundfix/agents/openai.yaml`; `internal/baseline/assets/modules/autonomous-work.json` |
| 0127 affected historically governed contracts | `internal/cli/cli_test.go`; `skills/baseline_skill_contract_test.go`; `internal/docscontract/publicdocs_test.go` |

Generated guidance and copies must be listed in the corresponding grant or
covered by the repository's expressly sanctioned regeneration contract:
`skills/roundfix/SKILL.md`, `skills/roundfix/agents/openai.yaml`,
`skills/implement-spec/SKILL.md`, `docs/agents/agent-instructions.md`,
`docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, and any
additional exact rendered path the concrete Baseline plan names.

Run `rtk make skills-sync` and `rtk make baseline-digests` after authorized
source edits; use the public Baseline update plan to regenerate repository
guidance. The sanctioned digest command owns its declared deterministic fallout
in the artifact directories named above. Inspect its actual changed-path output;
do not hand-edit a pin or assume it authorizes unrelated fixture changes.
The grant record is a separate preceding commit, as required by the
[universal guide](../agents/agent-instructions.md). No Makefile, dependency,
workflow, or upstream-skill edit is proposed here.

## Verification and rollout evidence

Existing tests establish useful boundaries but do not prove the proposed feature:

- `TestRunEnforcesReviewRequestCoherence` and
  `TestRunExemptsFetchFromReviewRequestCoherence` in `internal/preflight/preflight_test.go`.
- `TestReviewRequestContract` in `skills/baseline_skill_contract_test.go`, and
  `TestReviewHistoryConfiguration` in `internal/docscontract/publicdocs_test.go`.
- `TestRunImplementDetachSurvivesCallerProcessGroupKill`,
  `TestImplementRunWindow`, `TestImplementRunWindowCrossing`, and
  `TestRunImplementResumesStaleInProgressTask` in `internal/cli/implement_test.go`.

0126 must prove no CodeRabbit call or new review Run can originate from any
retired command, omitted/legacy config, or owned workflow; historical Runs and
reports must still render. Exercise the actual selected native reviewer against
an external repository fixture containing a real detectable defect, then its
corrected version. Prove read-only enforcement, current-head rejection,
malformed/error output, cancellation, and missing evidence through the public
surface. A stub is useful for deterministic failure paths but cannot establish
native model output or permission behavior.

0127 must interrupt the owner before and after each externally visible action,
then show recovery without duplicate PRs, merges, or Spec Runs. Exercise the
window, exhausted limits, stale checks/review, lost credentials, a held user
decision, and a later Spec whose assumptions changed. Preserve receipts of what
GitHub accepted rather than trusting the agent's narration.

Run the Task-declared commands, `rtk make verify`, and `rtk make verify-docs`
at their governed boundaries. Review removed contract tests against the new
policy instead of deleting assertions only to get a green suite. Rollback is a
reviewed revert or forward fix; never delete old Run evidence or silently
reactivate CodeRabbit to recover delivery. Database compatibility must be
settled before 0127 permits an older binary to read new state.

During this planning session the supervising agent reported a fresh `go vet`
failure with 31 `copylocks` diagnostics. This design did not independently rerun
that command. The approved prerequisite repair and its evidence must land before
the affected tooling change; no clean verification baseline is asserted here.

## Sources, influence, and limits

Both required research sources were consulted. Local Secondbrain consultation
started at `/Users/marcio/dev/secondbrain/wiki/index.md`, used the advisory qmd
search, and then ran:

```bash
rtk qmd query 'Revisão independente de código e persistência de workflow autônomo Roundfix com Codex e Claude Code' --all --files --min-score 0.3
```

| Source read | How it affected this proposal |
| --- | --- |
| `/Users/marcio/dev/secondbrain/inbox/roundfix/_triaged/2026-08-25-review-source-unica-e-a-frota-desligando-o-coderabbit.md` | Prior fleet evidence already identified the missing pre-PR path and provider-bound evidence. Its August fleet counts were not remeasured and are not current claims here. Local source inspection confirmed the Roundfix limitation. |
| `/Users/marcio/dev/secondbrain/wiki/concepts/agent-harnesses.md` | Informed durable ownership, ordered mutation, receipts, and independent review. It is a synthesis, not a benchmark proving this implementation. |
| `/Users/marcio/dev/secondbrain/wiki/concepts/agent-workflows-e-loop-engineering.md` | Informed stop conditions, external verification, and recovery requirements; challenged treating a detached Task Run as an entire durable delivery workflow. |
| [OpenAI code review](https://learn.chatgpt.com/docs/code-review), found and read through Exa MCP | Supports native review before publication and non-mutating prioritized findings. Local CLI help supplied the exact installed flags; parser rejection was reproduced without spending model tokens. |
| [Claude headless execution](https://code.claude.com/docs/en/headless) and [CLI reference](https://code.claude.com/docs/en/cli-reference), found and read through Exa MCP | Support headless output, schema, and capability controls. The documented subscription limitation of bare mode changed the isolation proposal. |

The external pages establish available primitives, not end-to-end Roundfix
correctness. No live reviewer, permission attack, spend limit, queue recovery,
GitHub mutation, or implementation test ran while writing this design. Native
reviewer selection and result characterization remain explicit prerequisites.
The supervising session owns capture of this sourced research digest in the
Secondbrain inbox; this artifact does not claim that ingestion occurred.
