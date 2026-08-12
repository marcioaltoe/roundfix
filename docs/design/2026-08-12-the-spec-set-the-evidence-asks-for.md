# The Spec set the evidence asks for (2026-08-12)

This is a routing document, not a Spec. It reads every live Finding and Backlog
Entry in this repository, plus the nine Inbox Entries addressed to Roundfix that
were triaged into `docs/findings/` and `docs/backlog/` on 2026-08-12, and groups
them into eighteen proposed Specs ordered by measured cost, plus one priority
correction the maintainer directed on 2026-08-12 that ships before the rest.

The coverage audit near the end of this document is the honest answer to "does
this set cover everything": it did not, on the first pass. Nine items were
uncovered, two Specs were added and four had their scope widened to close them.

Nothing here is a commitment. Each proposal names the evidence it rests on, the
scope it would take, and the route from `docs/agents/spec-routing.md` that fits.
The next step for any one of them is `write-prd` or `write-techspec`, which is
where the decisions this document deliberately leaves open get settled.

## What is already closed

Verified against the tree at `366e664` before proposing anything. These records
are stale and should move through their lifecycle rather than into a Spec.

| Record | Closed by |
| --- | --- |
| `2026-08-07-a-sixty-four-value-bound-locks-out-the-opencode-runtime.md` | Spec 0088; `maxRetainedCapabilityValues` now bounds retention, `maxAdvertisedCapabilityValues` bounds plausibility |
| `2026-08-08-a-runtime-that-advertises-a-catalog-is-not-unusable.md` | same |
| `2026-08-07-claude-agent-selections-are-never-proven.md` | Spec 0091 (ADR-0119); `internal/agent/catalog.go` corrected |
| `2026-08-06-three-gigabytes-of-event-journal-inside-the-retention-window.md` | Spec 0081 |
| `2026-08-07-a-live-contract-lives-inside-an-archived-spec.md` | the coverage record moved to `docs/references/coverage-record.json` |
| `2026-08-06-a-promoted-backlog-entry-has-nowhere-valid-to-go.md` | its own 2026-08-07 addendum; `SC-BACKLOG-UNMOVED` enforces it |

Two more are partially closed and their remainder is folded into a proposal
below rather than left standing:

- `2026-08-07-the-only-gate-reports-green-on-a-red-suite.md` — the Makefile now
  sets `GO := go` and keeps the wrapper on `GO_HUMAN`, so the authoritative gate
  no longer routes through it. The cache half of the same false-green class is
  still live and is Spec 0104 below.
- `2026-08-11-a-git-worktree-that-fails-only-under-load.md` — production worktree
  creation now sets `core.fsmonitor=false` (`internal/worktree/worktree.go:1917`).
  The remaining ask, a Roundfix sentence around the raw git fatal, is in Spec
  0097.

## The outside-evidence rule

`docs/agents/spec-routing.md` requires each Spec to rest at least one acceptance
row on evidence originating outside its own artifacts. Most of this set has that
for free: the friction was measured in `fluxus`, `vortex`, `oraculum` and
`fiscus` — repositories these Specs did not build — and those measurements are
now recorded in `docs/findings/2026-08-12-*`. Each proposal names its outside
source. Where none exists, the proposal says so.

---

## Wave 0 — the priority correction

Maintainer direction on 2026-08-12: this Spec ships, and a fix release follows
it, before any other implementation in this set starts.


### 0094 — one history root under `docs/`

**Problem.** Retired documentation is stored in three places and none is right.
Spec 0085 built the single archive root the maintainer asked for on 2026-08-09 —
"uma estrutura com `_archived/specs|findings|adr|backlog`" — and anchored it at the
repository root, so an archive of documentation now sits beside `cmd/`,
`internal/` and `skills/`. Every other adopted repository still holds the pre-0085
per-tree layout. Retired decision records never moved at all. And Review Artifacts
have no terminal home: the orphan case has accumulated to fifty folders here.

**Two facts found while designing, both of which shrink the work.** A review whose
owning Spec is known already writes to `<spec>/reviews/` — no underscore — and
already travels with its Spec into history; only the orphan case needs a home.
And Spec 0085 already gave every decision record a lifecycle status, so measured
2026-08-12 there are zero retiring records in `docs/adr/` and one declined intent
entry: that family moves one file today and gains its rule for the future.

**One fact found while designing that grows it.** `!**/_reviews/**` matches by
leading underscore, so `docs/specs/<slug>/reviews/` is **not** excluded from
review. It escapes only because all twelve Spec-owned review directories currently
sit inside the archive. The first active Spec that gains one reintroduces the
self-review deadlock recorded on 2026-08-06. Path-anchoring the exclusion closes
it.

**Scope.** The resolver's answers move to `docs/history/{specs,findings,adr,backlog,reviews}`
and it stays the single place that knows. The orphan Review Artifact root loses its
underscore and gains one refusal: it never resolves into history, so a new Round for
an already-retired Spec's Pull Request writes to the live root. Retirement
classification for decision records and intent entries, narrowed so a pending
proposal is not filed as history. Review liveness from local Git. Layout discovery
across every legacy shape. The relocation ledger in the Plan and its transaction.
Then this repository's own bytes, the carriers that name the old location, and the
review configuration.

**The fleet question inverts.** Roundfix is the only repository carrying Spec
0085's root archive, so nothing outside it migrates away from that — the local move
is a rename and a commit. Every other adopted repository holds the pre-0085
per-tree layout, and those are what the Baseline migration serves. Measured
2026-08-12: four candidate repositories hold 173 to 735 archived files each.

**The design question, and why it is not the file move.** An orphan review is
finished when its Pull Request is, and that fact lives outside the repository. The
accurate source — the provider's Pull Request state — cannot be a requirement,
because reading it needs a credential and a network the migration must run without
(`gh` answered `HTTP 401` in the authoring session). ADR-0123 decides on local Git
reachability instead and accepts that the answer can lag: a squash merge leaves no
ancestor, so an undecidable head stays live. That is the safe direction — keeping a
finished artifact costs a directory, retiring a live one breaks a Round in flight.

**The tension recorded rather than hidden.** Spec 0085's Goal 2 was "retired
material leaves the directories an Agent loads by default", and `docs/` is such a
directory in the loosest reading. ADR-0120 supersedes that goal explicitly, on the
grounds that an Agent reaches `docs/agents/` through `CLAUDE.md` rather than
loading `docs/` wholesale. Goals 1 and 3 — one root, one filter — are preserved.

**Evidence.** `docs/specs/0094-one-history-root-under-docs/references/2026-08-12-the-archive-root-sits-beside-docs-instead-of-inside-it.md`,
adopted from `docs/backlog/` on 2026-08-12;
`docs/workflow/authorizations/2026-08-09-what-an-agent-reads-before-it-decides.md`,
which records the original request verbatim;
`_archived/specs/0085-what-an-agent-reads-before-it-decides/task_04.md`.

**Route.** `write-prd` → `write-techspec` → `write-tasks`. Complete: PRD, TechSpec,
and a nine-Task graph in six waves, all checker-clean.

**Authorization — granted 2026-08-12.**
`docs/workflow/authorizations/2026-08-12-the-archive-root-under-docs.md`.

**Then a fix release**, per the maintainer's direction, before Wave 1 starts.

---

## Wave 1 — the loop's honesty

The largest measured cost in the whole corpus. Of 201 failed Tasks since
2026-08-03, 123 are the QA gate returning a verdict rather than code breaking,
and 7 are a Verification that was vacuous
(`docs/backlog/2026-08-10-the-loop-is-measured-and-the-gate-is-where-it-costs.md`).
Two independent repositories name the same lever as their top priority.

### 0095 — a Verification that ran before anyone believed it

**Problem.** The Verification contract is checked for form and never for
execution. Six defects in one `oraculum` night were pure shell semantics; the
`fiscus` 0004 graph came within one runner's exit code of being entirely vacuous.

**Scope.** `spec check --run-verification`, executing each `## Verification` line
in a disposable checkout and reporting the exit code. Two new detectors:
`SC-VERIFY-INVERTED-EXIT` for known reversed forms (`grep -c`, `grep -v | wc -l`,
`test $(…)` without a comparison) and `SC-VERIFY-NON-HERMETIC` for undeclared
environment variables, `test -n "$VAR" &&` guards, and dependence on state
outside the repository. Re-enable `SC-VERIFY-VACUOUS-COMMAND`, registered and
staged out at `internal/speccheck/coherence.go:59`. Accept a `creates:` entry in
`## Context` so a Task can declare its own output without `SC-REF-UNRESOLVED`
firing. One sentence in the authoring contract: a Verification command passes
only by exiting zero, with the three working forms as the worked answer.

**Evidence.** `2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md`
finding 5; `2026-08-12-five-unresolved-runs-to-deliver-one-spec.md` finding 1;
`2026-08-12-three-consecutive-specs-measure-the-loop.md` finding 3;
`2026-08-09-a-verification-command-passes-only-by-exiting-zero.md`;
`2026-08-10-the-loop-is-measured-and-the-gate-is-where-it-costs.md`.

**Route.** `write-prd` → `write-techspec` → `write-tasks`. The `write-tasks`
skill edit needs express tooling authorization with bounded files.

**Open decision the PRD must settle.** Whether `--run-verification` runs on the
pre-work tree (proving a command *can* fail) or on the post-work tree (proving it
*can* pass), or both as distinct modes. The vacuity rule and the exit-zero rule
pull in opposite directions here.

### 0096 — a failure the Agent can read

**Problem.** A failing Verification can hand the Agent Session nothing at all,
and the same failure can repeat across Runs without anyone saying it repeated.

**Scope.** When a Verification fails with empty stdout and stderr, say exactly
that in the Verification Feedback instead of sending an empty message. When a
verdict and diagnostic signature match a prior failure of the same Work Item,
report it in the report and in the `task-status` Run Event. Name the sanctioned
exits when the two-corrective-Task ceiling is reached — amend the TechSpec and
recut, or promote the excess to its own Spec — so the excess becomes a decision
inside the loop's authority instead of a policy stop. Say which surface a task
file was read from: `settle` already prints `Settle surface:` and does not state
that the task file comes from there, so fixing an unsatisfiable Verification
means editing the file in two places, discovered by trial. Settle whether
`budget.max_run_duration` is evaluated at Work Item boundaries and say so in its
documentation, since a Run measured at 2h34m against a configured 2h has no
second case and no established cause.

**Evidence.** `2026-08-12-a-redirected-verification-hands-the-agent-an-empty-diagnostic.md`;
`2026-08-12-five-unresolved-runs-to-deliver-one-spec.md` findings 4, 5 and 8;
`2026-08-12-three-consecutive-specs-measure-the-loop.md` finding 7.

**Route.** `write-techspec` → `write-tasks` for the first two; the ceiling exits
are authoring-contract work and travel with 0106.

**Outside evidence.** Run `run_20260808T153649Z_78746d4b80d08fc7` and the two
`fiscus` task_07 diagnostics.

### 0097 — a wave that cannot collide

**Problem.** Raising `worktree.concurrency` made visible three failures that were
already latent: bootstrap collides on the shared `.git`, Tasks declared
independent edit the same file and die at cherry-pick, and Task Worktree creation
fails under load with a raw git errno.

**Scope.** Serialize Worktree Bootstrap even when Tasks run in parallel. Detect
file-set intersection between Tasks of the same wave before dispatch and either
serialize or refuse, naming the paths. Wrap a Task Worktree creation failure in a
Roundfix sentence carrying the Run, the Task and the concurrency level, with the
git text as evidence rather than as the message.

**Evidence.** `2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md`
findings 2 and 3; `2026-08-11-a-git-worktree-that-fails-only-under-load.md`;
`2026-08-12-three-consecutive-specs-measure-the-loop.md` finding 8 rule 5.

**Route.** `write-techspec` → `write-tasks`.

**Outside evidence.** The `oraculum` Specs 0032 and 0035 Runs, which lost all
Agent work at cherry-pick.

### 0098 — a hook that cannot outrank the gate

**Problem.** The Daemon is the verification authority and then commits through a
`pre-commit` hook that can refuse. Three Runs died with correct, already-verified
work staged in the Task Worktree, and the command designed for recovery refused
because the status was `completed` rather than `failed`.

**Scope.** Settle the design conflict: either the Daemon commits as the authority
it already is, or a hook refusal is a repairable class that spends a repair
round, or at minimum it is detected and names the recovery. Cover the state
"verified, settled, not committed" in recovery. Write the invariant — a commit
hook may never be stricter than the authoritative Verification — into the
Baseline module that owns it.

**Evidence.** `2026-08-12-a-hook-failure-kills-a-run-that-already-verified-its-work.md`.

**Route.** `write-prd` → `write-techspec` → `write-tasks`; the choice between the
three shapes is a product decision, not an implementation detail.

**Check first.** Spec 0092 ("a run that can hand back its work") may already close
the `settle` half. Confirm before authoring, and narrow the Spec if so.

---

## Wave 2 — the tool meets the repository as it is

Two divergences of one shape, measured in the same `fluxus` session: Roundfix
assumes a convention the target repository does not follow and fails closed
without offering to accommodate the real one.

### 0099 — a tool that meets the repository's conventions

**Problem.** `archive` refuses a graph that declined the QA gate, refuses a
declared-only `partial` that its own `--help` promises to accept, and
`release plan` demands a `v`-prefixed tag in a repository whose release workflow
uses the absent prefix to decide that a release is stable.

**Scope.** `archive` accepts `qa: declined` with a non-empty `qa_reason` as an
equivalent verdict, recording the reason in the archival artifact. The Daemon
settles the QA Task as `completed` when a `partial` is declared-only — the same
check `archive` already performs. `release plan` detects the tag pattern from the
repository's existing tags instead of prescribing one.

**Evidence.** `2026-08-12-archive-refuses-a-graph-that-declined-the-qa-gate.md`;
`2026-08-12-release-plan-requires-a-v-prefixed-tag.md`;
`2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md` finding 1.

**Route.** `write-techspec` → `write-tasks`.

**Outside evidence.** `fluxus` Spec 0023, whose Run reached `Clean` with 4 of 4
Tasks and could not be archived, and the eight unprefixed `fluxus` release tags.

### 0100 — a review the loop always asks for

**Problem.** With automatic review disabled, the request is published only for
heads the loop itself pushes. A head pushed by hand, and a Pull Request with no
findings at all, both wait thirty minutes for evidence nobody asked for. Two
Runs on PR #153 spent sixty minutes of wall clock and delivered nothing against
twelve issues that were fetchable the whole time.

**Scope.** Request the review on the head the Run starts from when the configured
Review Source has no accepted evidence for it, using the existing
`ReviewRequestMarker` for idempotence instead of comparing against
`currentHeadSHA` (`internal/watch/watch.go:555`). Fail fast on a refusal Roundfix
can already read: `success` plus "automatic reviews are disabled" plus
`request_review: true` is a known state with a known remedy, reachable in the
first minute. Do not discard the wait on a transient `gh` failure —
`ghCommandError` already carries `temporary` and nothing consumes it on that
path. Add the coherence Preflight for a ruleset requiring thread resolution
against `include_nitpicks: false`.

**Evidence.** `2026-08-10-a-head-the-loop-did-not-push-is-a-head-nobody-reviews.md`;
`2026-08-12-three-consecutive-specs-measure-the-loop.md` findings 4 and 5;
`2026-08-12-a-hook-failure-kills-a-run-that-already-verified-its-work.md` cites
the same request gap on a clean Pull Request.

**Route.** `write-techspec` → `write-tasks`.

**Do not regress.** Twice, Roundfix refused to read a green CodeRabbit check as
proof of review — once with `Review skipped: automatic reviews are disabled`,
once with `Review limit reached`. That distinction prevented a merge on a false
review signal and must survive this change.

### 0101 — a terminal branch with one disposition

**Problem.** Terminal Runs leave commits on Run Branches that only a human can
dispose of. Ancestry proof fails after a rebase even when the content is already
integrated; a superseded branch gets a recovery suggestion that would regress the
repository; a worktree of the same repository gets its own `repo-id`, so each
side sees the other's Run Branches.

**Scope.** Prove integration by content — identical tree, or absence of an
exclusive file — not only by ancestry. Classify a Run Branch as releasable when
its terminal Run's target branch no longer exists, or offer a way to record
supersession with a reason for `reconcile` to act on. Resolve one repository
identity across a repository and its worktrees.

**Evidence.** `2026-08-12-three-consecutive-specs-measure-the-loop.md` finding 6;
`2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md` finding 8;
`2026-08-12-a-hook-failure-kills-a-run-that-already-verified-its-work.md` records
the manual integration ladder; `2026-08-06-rollup-run-lifecycle-and-branch-integrity.md`
is the standing rollup.

**Route.** `write-prd` → `write-techspec` → `write-tasks`. Disposition semantics
are product contract.

**Why it matters beyond convenience.** In two sessions the only exit was
`--skip-branch-integrity`, which disables the guard for legitimately pending work
too. A guard whose normal path is impassable teaches its own bypass.

---

## Wave 3 — preflight and the gate

### 0102 — a preflight that proves what the Run needs

**Problem.** Preflight proves things the Run will not use and misses things it
will. Every configured fallback tuple must prove, so an intermittent adapter in a
category that would never be reached refuses Runs whose preferred Selection is
perfect. Meanwhile a declared `worktree.copy` source that does not exist is only
discovered per worktree, after the Run exists.

**Scope.** Validate each `worktree.copy` source in Preflight and refuse naming
the absent ones. Always prove the preferred Selection; prove fallbacks lazily or
tolerantly, or allow an opt-out. A cleanup failure that only reports a Session
absent because the setup error explains why no Session exists must not be joined
into the surfaced error, and no cleanup error should outrank the setup error in
the printed `next:`.

**Evidence.** `2026-08-12-five-unresolved-runs-to-deliver-one-spec.md` finding 3;
`2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md` finding 4;
`2026-08-08-a-failed-proof-appends-a-cleanup-error-the-maintainer-cannot-act-on.md`.

**Route.** `write-techspec` → `write-tasks`.

### 0103 — a suite that leaks nothing

**Problem.** The test suite writes into the repository it is reading, leaves
processes running for days, and carries a flake family whose common condition is
spawn density. None of it is visible from any Roundfix command.

**Scope.** Establish and prove the invariant that a test may not write inside the
repository root — a harness snapshotting `git status --porcelain` around the
suite names the offending case directly. Terminate the detached child in the
detach tests rather than only releasing it, move the release sentinel out of
`t.TempDir()`, and bound the fake ACPX wait. Replace shell-script fakes with
compiled test binaries re-executing `os.Args[0]`, removing the write-then-exec
window. Give Roundfix a way to see its own residue — `runs list --orphans` or a
`doctor` line reporting detached processes with no live Run record. Prove the
process tree on Force Stop, not the registered owner alone. Keep a gate's scratch
repositories and built binaries out of the Spec's evidence directory.

**Evidence.** `2026-08-06-a-test-mutates-the-repository-another-test-is-reading.md`;
`2026-08-06-the-detach-tests-leak-the-process-they-prove-survives.md`;
`2026-08-10-a-fake-adapter-goes-silent-under-a-dense-start.md`;
`2026-08-09-the-qa-gate-commits-scratch-repositories-and-binaries.md`;
`2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md` finding 7.

**Route.** `write-techspec` → `write-tasks`.

**Measured cost.** Four orphaned processes aged up to 3d15h consumed 2h40m of CPU
doing nothing, and a lost race in the gate can both report a false failure and
delete tracked files from the operator's working tree.

### 0104 — a gate that cannot certify its own cache

**Problem.** `docs/agents/specific-repository.md` tells a reader to run
`go clean -testcache` before trusting `make verify`. The Makefile exports
`GOCACHE ?= $(CURDIR)/.gocache`, so a bare `go clean -testcache` clears a
different cache and leaves the gate's own exactly as stale. Measured during Spec
0088: a real regression reported `make verify exit=0` with
`ok roundfix/internal/spec (cached)`.

**Scope.** A Makefile target that cleans the cache the gate actually uses, with
the guidance pointing at it. Settle whether `make verify` should refuse to report
success when its Go suite is entirely cached, since a cached pass proves the
cache and not the tree. Establish whether the `rtk` wrapper can lose a
package-level failure in a multi-package run, or record that it is unowned: the
Makefile now sets `GO := go` so the gate no longer routes through it, but
`GO_HUMAN` still does and the mechanism was never reproduced.

**Evidence.** `2026-08-08-go-clean-testcache-clears-a-cache-the-gate-does-not-use.md`;
the residue of `2026-08-07-the-only-gate-reports-green-on-a-red-suite.md`.

**Route.** `write-techspec` → `write-tasks`.

**Blocking precondition.** Both changes touch protected tooling. The PRD and
TechSpec must record express maintainer authorization with the exact bounded
files before any Task may run, and the authorization lands as its own commit.

### 0105 — the gate's own economics

**Problem.** The expensive tail of the loop is one family: the Spec assumes the
world behaves like its fakes. Spec 0091 proved it end to end — a design premise
survived authoring, implementation and every unit test, and died at the QA gate
four Runs later. Separately, the Pull Request row is unreachable by design in
every Spec and cost six of eight gate executions in one `vortex` Spec.

**Scope.** When a Spec crosses an external surface — an ACP adapter, an HTTP
contract, a database — its characterization Task records what the real thing
does. Apply the equivalent-evidence path to the Pull Request row by default in
`qa-gate`, requiring the evidence to be recorded. Make `## Verification`
immutable for the review-resolution Agent, or surface an edit to it as a contract
change in the Round report. Refuse to let the review Agent resolve a finding that
requires test infrastructure the repository does not have — choosing a test
substrate is Spec scope, not a correction round — and treat a newly introduced
environment variable as the cheap trigger for that refusal or for human review.
Roundfix owns the QA Task's Verification instead of
letting an author re-assert a derived verdict. Report independent static findings
together so a governance failure does not blind the whole flow matrix. Teach the
citation parser a conjunction separator and a bare ADR number, or name the
recognised form in the failure message.

**Evidence.** `2026-08-10-the-loop-is-measured-and-the-gate-is-where-it-costs.md`;
`2026-08-12-three-consecutive-specs-measure-the-loop.md` findings 1, 2 and 9;
`2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md` finding 6;
`2026-08-12-five-unresolved-runs-to-deliver-one-spec.md` finding 2;
`2026-08-12-a-hook-failure-kills-a-run-that-already-verified-its-work.md` records
the governance blindness case.

**Route.** `write-prd` → `write-techspec` → `write-tasks`. This is large enough to
split; the PRD should decide whether the review-Agent guardrails travel
separately.

**Read against.** Eleven non-`pass` verdicts in one session failed on contract,
not on business logic — and the same gate found four real defects no suite would
have caught. Any proposal that loosens the gate must answer to both numbers.

---

## Wave 4 — the Baseline and the method

### 0106 — a decision that reaches every artifact

**Problem.** The Baseline's preservation path diverges from its greenfield path.
Excluding an artifact removes it from the manifest and leaves its bytes on disk
with the setup markers intact, so the repository states a rule the Baseline no
longer recognises in a file whose markers claim setup owns it. Greenfield
adoption is unreachable for any repository with stale managed carriers. Changing
the HTTP Contract discards every recorded exception silently. Preflight warns on
the Baseline's own fixtures on every plan and apply, teaching its reader to skip
the one warning class that also carries the real signal.

**Scope.** Exclusion plans a deletion the way readoption already does, and strips
an orphaned managed region from a shared carrier such as `AGENTS.md`. Greenfield
either accounts for stale managed carriers without classifying them, or refuses
up front naming Preservation. The `http.contract` "Change" branch carries
exceptions and provenance forward. The inert `httpContract.default` in the
profile assets is settled — read, deleted, or made a real override — along with
its unimplemented error code. Carrier discovery excludes the Baseline's own
embedded assets. `SC-ADR-RELATED` scopes to active Specs so the archived count is
structurally zero. Give the derivation chain one mechanical owner: a declared
source change must identify every derived artifact and one sanctioned command
must regenerate the complete consequence set, with the semantic retention proof
that the Baseline lifecycle still lacks. Part of that chain is the Setup
Manifest's recorded catalog digest, which nothing compares against the embedded
catalog: a catalog edit leaves the manifest naming a catalog that is gone, and
the gate passes. Spec 0094 meets this directly — it edits catalog modules under
an authorization that does not reach the manifest, so it must leave the drift in
place rather than repair it opportunistically.

**Evidence.** `2026-08-10-an-excluded-artifact-outlives-the-decision-that-excluded-it.md`;
`2026-08-06-rollup-baseline-and-derived-tooling.md` live edge;
`2026-08-10-the-preflight-reads-the-catalog-fixtures-as-repository-carriers.md`;
`2026-08-07-greenfield-adoption-cannot-satisfy-its-own-gate.md`;
`2026-08-07-changing-the-http-contract-discards-its-exceptions.md`;
`2026-08-07-two-http-contract-defaults-and-only-one-is-read.md`;
`2026-08-06-minting-an-adr-opens-gaps-no-one-can-ever-close.md`.

**Route.** `write-prd` → `write-techspec` → `write-tasks`. Deleting repository
bytes on a decision flip is exactly the kind of change the Baseline otherwise
asks about, so the product decision comes first.

**Carries a bounded authorization.** The profile-asset half is bounded to
`internal/baseline/assets/profiles/**`, which the 2026-08-07 REST authorization
explicitly did not reach.

### 0107 — the authoring rules the guides do not carry

**Problem.** Three fleet repositories independently derived the same class of
rule from measured failures, and every one of them belongs to a guide Roundfix
owns — so it cannot be written where it was learned without being overwritten by
the next `baseline update`.

**Scope.** Carry into the Baseline modules: a Verification is hermetic; a
requirement describes the property, not the magnitude of the data; a use case
test is born against ports and doubles while a persistence proof is born in
infra; every Task that changes an already recorded contract updates the record in
the same slice; commit scope is per Spec. Add the tooling-authority chronology
rule — the commit that authorizes a tooling mutation does not contain it.
Document `## Unreachable Acceptance`: where it lives, the entry format, and what
`satisfied-by` means. Give `write-tasks` the rule that same-wave Tasks do not
share an edit target. Give `autonomous-work.md` an owner for environment
preparation and a rule for an identically repeated Work Item failure. Have
`write-techspec` ask for tooling authority by class rather than by path. Give the
Secondbrain guide a production clause with a trigger. Close the traceability
chain the authoring rollup still stands on: a finding or ADR reaches a Task, and
that Task's evidence proves the authored consequence shipped — a citation check
proves an obligation was named, not obeyed.

**Evidence.** `2026-08-12-three-consecutive-specs-measure-the-loop.md` finding 8;
`2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md` finding 9;
`2026-08-12-five-unresolved-runs-to-deliver-one-spec.md` findings 6 and 7;
`2026-08-12-the-secondbrain-guide-permits-capture-without-triggering-it.md`;
`2026-08-06-rollup-spec-authoring-and-contract-enforcement.md`.

**Route.** `write-techspec` → `write-tasks`.

**The rule this Spec must obey itself.** A finding filed on 2026-08-06 predicted
Spec 0090's F-001 exactly, and nothing stopped it until the rule became a
detector. Every clause added here should name which check decides it, or state
that none does.

### 0108 — what an Agent loads to answer one question

**Problem.** 103 KB of catalog JSON renders 1,450 lines of prose into thirteen
files that also hold 41 local rules, and most of a `baseline update` plan's output
is the retention evidence proving it did not eat them. Separately,
`.agents/skills/roundfix/SKILL.md` is 2,172 lines across 29 sections — six times
the next largest owned skill — and an Agent resolving a review comment loads
release planning, retention policy and the whole configuration schema with it.

**Scope.** Divide by enforcement rather than by canonicity: a clause a mechanical
check already decides is documentation and can live wherever it is cheapest to
read; a clause that governs every task and has no check stays always-on and
minimal; a clause that governs one kind of work belongs to that skill;
repository personalisation belongs in `AGENTS.md`. Split the roundfix skill into
`references/` with the skill keeping what routes.

**Evidence.** `2026-08-09-the-canonical-method-lives-in-rendered-guides-instead-of-skills.md`;
`2026-08-09-the-roundfix-skill-carries-nine-commands-in-one-file.md`.

**Route.** `write-prd` → `write-techspec` → `write-tasks`, and it either extends
or supersedes part of Spec 0085.

**The risk that decides the shape.** A guide is always-on and a skill is on
demand. Mandatory dispatch is itself prose in `skill-dispatch.md` and nothing
verifies it, so the move trades a weak guarantee for a weaker one unless dispatch
becomes checkable. Measure how much of the 2,172 lines a typical dispatch
actually uses before assuming the cost is real.

### 0109 — what a Session consumed

**Problem.** Roundfix persists the effective Agent Selection for every Task and
records nothing about what that Session consumed. Concurrency and reasoning
effort are different knobs, and a provider's daily total cannot separate them.

**Scope.** Attach consumption to each Agent Session alongside the identity
already persisted: input and output tokens with cache reads separated when the
adapter reports them, the owning Task, Task Type and Batch, whether the Session
ran a Preferred or a Fallback Selection, and whether it was the initial turn or a
Verification retry. A missing measurement stays observably absent rather than
becoming a zero.

**Evidence.** `2026-08-08-record-usage-per-agent-session.md`.

**Route.** `write-prd` → `write-techspec` → `write-tasks`.

**Settle first.** Whether `codex-acp` and `claude-agent-acp` expose usage at all.
The event shape cannot be fixed before that answer.

### 0110 — a refresh that does not re-interview, and a prompt that ends

**Problem.** Spec 0082 closed the manifest-driven update and explicitly left
three observations from the same finding live and unowned.

**Scope.** The sealed prompt outliving `SealedPromptTimeout` — a run left alone
was still alive past ten minutes against a declared five-minute bound, and it
affects first adoption and every other sealed-prompt caller.
`internal/agent/sealed.go` now passes `--timeout` and `--ttl` to the subprocess
in addition to the context deadline, so the first Task is to establish whether
the observation still reproduces before designing anything. Which binary the
sealed prompt actually spawns for runtime `codex`, given that `doctor` reports
`~/.local/bin/codex` while `which -a` resolves a cmux shim first — the finding is
explicit that nobody should act on the shim as a cause without proving it.
Selecting a replacement profile mid-interview re-asks decisions already answered,
including ones both profiles share.

**Evidence.** `2026-08-07-the-setup-refresh-interviews-a-repository-that-already-answered.md`,
whose Route section names all three as outside Spec 0082;
`2026-08-06-rollup-baseline-and-derived-tooling.md`.

**Route.** `write-techspec` → `write-tasks`, with a reproduction Task before any
design Task. Two of the three have `unknown` root causes on the record.

### 0111 — one terminal audit across a Run's surfaces

**Problem.** Two standing rollups converge on the same missing thing. Run
lifecycle spreads across process trees, Run Branches, Task and Run Worktrees,
refs, artifacts, notifications and database storage, and no single terminal audit
covers ownership and classification across them. Agent selection spreads across
accepted configuration, runtime launch and Task filesystem access, and nothing
proves those three describe the same environment.

**Scope.** One audit that reports the terminal disposition of every resource a
Run created, across repositories, preserving the evidence needed to check the
cleanup. One execution-policy proof that exercises the same selection and access
contract the Task Session will receive, so an accepted `agent_full_access` or a
documented sandbox escape cannot validate and then fail at the first real Task.

**Evidence.** `2026-08-06-rollup-run-lifecycle-and-branch-integrity.md` live edge;
`2026-08-06-rollup-agent-selection-and-execution-environments.md` live edge.

**Route.** `write-prd` → `write-techspec` → `write-tasks`.

**Depends on.** 0097, 0101, 0102 and 0103, each of which closes one surface this
audit would otherwise have to model on its own. Authoring it before them would
build the audit against seams that are about to move.

---

## Coverage audit

Every live Finding, every numbered observation inside the three findings minted
on 2026-08-12, every open Backlog Entry, and every rollup live edge, mapped to
the Spec that owns it. Run on 2026-08-12 against this document's first pass; the
nine uncovered items it found are marked and were closed by adding Specs 0110 and
0111 and by widening 0096, 0104, 0106 and 0107.

### Findings — one record per row

| Record | Owner |
| --- | --- |
| a-test-mutates-the-repository-another-test-is-reading | 0103 |
| minting-an-adr-opens-gaps-no-one-can-ever-close | 0106 |
| the-detach-tests-leak-the-process-they-prove-survives | 0103 |
| the-loop-cannot-fix-comments-about-its-own-artifacts | direct change, stated below |
| changing-the-http-contract-discards-its-exceptions | 0106 |
| greenfield-adoption-cannot-satisfy-its-own-gate | 0106 |
| two-http-contract-defaults-and-only-one-is-read | 0106 |
| the-only-gate-reports-green-on-a-red-suite | 0104 — **was uncovered** for the wrapper half |
| the-setup-refresh-interviews-a-repository-that-already-answered | 0110 — **was uncovered** |
| a-fake-adapter-goes-silent-under-a-dense-start | 0103 |
| a-head-the-loop-did-not-push-is-a-head-nobody-reviews | 0100 |
| a-git-worktree-that-fails-only-under-load | 0097 |

Seven further findings are closed and listed at the top of this document.

### The 2026-08-12 findings — one row per numbered observation

| Observation | Owner |
| --- | --- |
| vortex 1 — Pull Request row in the gate | 0105 |
| vortex 2 — resolve Agent rewrites `## Verification` | 0105 |
| vortex 3 — non-hermetic Verification, `creates:` | 0095 |
| vortex 4 — `include_nitpicks` × ruleset | 0100 |
| vortex 5 — transient `gh` failure | 0100 |
| vortex 6 — superseded Run Branch | 0101 |
| vortex 7 — `settle` reads the task file from the surface | 0096 — **was uncovered** |
| vortex 8 — five authoring rules | 0107 |
| vortex 9 — Heavy lift induces invented infrastructure | 0105 — **was cited as evidence with no matching scope** |
| oraculum 1 — declared-only `partial` | 0099 |
| oraculum 2 — parallel bootstrap collides on `.git` | 0097 |
| oraculum 3 — integration conflict found late | 0097 |
| oraculum 4 — preflight proves fallbacks | 0102 |
| oraculum 5 — `spec check --run-verification` | 0095 |
| oraculum 6 — Roundfix owns the QA Verification | 0105 |
| oraculum 7 — force stop leaves an orphan grandchild | 0103 |
| oraculum 8 — `reconcile` scoped by repo-id | 0101 |
| oraculum 9 — baseline instruction gaps | 0107 |
| fiscus 1 — `--run-verification` corroboration | 0095 |
| fiscus 2 — citation parser, Portuguese list and bare ADR | 0105 |
| fiscus 3 — missing `worktree.copy` source | 0102 |
| fiscus 4 — identical repeated failure | 0096 |
| fiscus 5 — two-corrective-Task ceiling | 0096 |
| fiscus 6 — tooling authority by class | 0107 |
| fiscus 7 — autonomous-work guide gaps | 0107 |
| fiscus 8 — Run above the configured budget | 0096 — **was uncovered** |

### Backlog — one record per row

| Record | Owner |
| --- | --- |
| a-failed-proof-appends-a-cleanup-error | 0102 |
| go-clean-testcache-clears-a-cache-the-gate-does-not-use | 0104 |
| record-usage-per-agent-session | 0109 |
| a-verification-command-passes-only-by-exiting-zero | 0095, with the diagnostic half in 0096 |
| the-canonical-method-lives-in-rendered-guides | 0108 |
| the-qa-gate-commits-scratch-repositories-and-binaries | 0103 |
| the-roundfix-skill-carries-nine-commands-in-one-file | 0108 |
| an-excluded-artifact-outlives-the-decision | 0106 |
| the-loop-is-measured-and-the-gate-is-where-it-costs | 0105 |
| the-preflight-reads-the-catalog-fixtures | 0106 |
| the-gate-accepts-a-manifest-that-names-a-catalog-that-is-gone | 0106 |
| a-hook-failure-kills-a-run (2026-08-12) | 0098 |
| a-redirected-verification-hands-the-agent-an-empty-diagnostic | 0096 |
| archive-refuses-a-graph-that-declined-the-qa-gate | 0099 |
| release-plan-requires-a-v-prefixed-tag | 0099 |
| the-secondbrain-guide-permits-capture-without-triggering-it | 0107 |
| the-archive-root-sits-beside-docs-instead-of-inside-it | 0094 — adopted, now under its `references/` |
| atomic-inbox-capture-helper | not in the set, stated below |
| speak-acp-directly-instead-of-through-acpx | not in the set, stated below |

`one-reader-in-cli-still-couples-verify-to-the-docs-tree` is `declined` with its
reason recorded and its resolution measured — 61s to 15.6s — so it is terminal,
not uncovered. `a-runtime-that-advertises-a-catalog-is-not-unusable` is closed
and listed at the top.

`the-gate-accepts-a-manifest-that-names-a-catalog-that-is-gone` was listed as
closed in this document's first pass and is not. Commit `bf9a702` corrected the
one stale value and created that backlog entry; it added no gate, and nothing in
the tree compares the Setup Manifest's recorded catalog digest against the
embedded catalog. The entry is open and owned by 0106.

### Rollup live edges

| Rollup | Live edge | Owner |
| --- | --- | --- |
| qa-gates-and-verification-evidence | gate authoring needs mechanical discovery and contract checks | 0095, 0105 |
| review-and-delivery-convergence | every non-decision interruption needs a typed recovery path | 0096, 0098, 0100, 0101 |
| run-lifecycle-and-branch-integrity | no single terminal audit covers cross-surface ownership | 0111 — **was uncovered** |
| agent-selection-and-execution-environments | config, launch and filesystem access must describe one environment | 0111 — **was uncovered** |
| spec-authoring-and-contract-enforcement | traceability from finding and ADR through Task to shipped evidence | 0107 — **was uncovered** |
| baseline-and-derived-tooling | one mechanically owned derivation path and semantic retention proof | 0106 — **was uncovered** |

Review Issue identity across Rounds, a member defect of the convergence rollup,
is closed rather than uncovered: `ReviewIssueFingerprint` and the duplicate
grouping in `internal/rounds/rounds.go` give an Issue a stable identity.

### What the audit does not claim

It maps every recorded observation to an owner. It does not prove the owner's
scope is sufficient to close it — that is what each PRD decides — and it does not
cover defects nobody has written down. Six of the mappings rest on a partial
reading of a rollup's live edge rather than on an enumerated member list, because
several members live under `_archived/findings/` and were read only through the
rollup's own summary.

## Deliberately not in the set

- `docs/backlog/2026-08-06-atomic-inbox-capture-helper.md` — small, and its
  trigger is 0107's Secondbrain clause. Revisit after that lands.
- `docs/backlog/2026-08-12-speak-acp-directly-instead-of-through-acpx.md` — large,
  unscoped, and Spec 0091 already answered its sharpest cost. Route through
  `write-idea` before any PRD.
- `docs/findings/2026-08-06-the-loop-cannot-fix-comments-about-its-own-artifacts.md` —
  the path filter closed the deadlock. The remaining half, emitting Markdown that
  passes the lint Roundfix asks other repositories to pass, is a few lines in the
  composer in `internal/rounds/rounds.go` and a test. Direct change, no Spec.
- The six standing rollups under `docs/findings/2026-08-06-rollup-*.md` — each
  should be reviewed for closure once its wave above lands, through the Finding
  lifecycle contract rather than as work of its own. Their live edges are not
  left implicit: every one is mapped in the coverage audit below.

## Suggested order

Waves are sequential in value, not in scheduling. Within them:

0. **0094 before everything**, by maintainer direction on 2026-08-12, with a fix
   release after it and before Wave 1 starts.
1. **0095** first and alone — it changes what every later Spec's Tasks are allowed
   to assert, and it is the only item two independent repositories named as their
   top priority.
2. **0096, 0098, 0099** next — each is small, each removes a class of Run that
   delivers nothing.
3. **0097 and 0103** together — both are about parallelism, and 0103's harness is
   what proves 0097 did not make things worse.
4. **0100, 0101, 0102** — the delivery surfaces; independent of each other.
5. **0104** early if a tooling authorization is convenient, since every claim made
   through `make verify` inherits its answer.
6. **0105, 0106, 0107, 0108, 0109** — the largest and the least urgent, and 0107
   should follow whichever of 0095–0103 lands first so its clauses cite shipped
   behaviour rather than intent.
7. **0110** anywhere after 0095; its first Task is a reproduction, so it can be
   cheap or can dissolve entirely if the timeout observation no longer holds.
8. **0111 last.** It depends on 0097, 0101, 0102 and 0103, and authoring it
   earlier would model seams that are about to move.
