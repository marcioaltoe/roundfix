---
spec: 0053-qa-gate-reachability-and-verdict-semantics
prd: _prd.md
created: 2026-07-28
---

# QA gate reachability and verdict semantics — Technical Spec

## Executive Summary

Three seams change: the QA report contract, the reconcile classifier, and
Branch Integrity Preflight. The primary trade-off is where the
environment-blocked distinction lives. Roundfix today parses only the QA
report's frontmatter `verdict` scalar — there is no row parser — and building
a markdown-table parser would pin an agent-authored table shape forever. We
instead extend the frontmatter with additive typed counts and keep the
verdict *rule* in the qa-gate Skill per ADR-0080, with `internal/spec`
enforcing consistency (a `pass` claiming finding-blocked rows is refused as
unreadable). The second trade-off is supersession proof: the classifier
derives it entirely from Git evidence it already reaches — the Daemon's
deterministic QA-commit message contract and the report-recency comparator —
keeping ADR-0053's proof-based rule and adding no store schema. The report
naming contract narrows to numeric suffixes only, because the existing
recency comparator ranks a non-numeric `-<scope-or-build>` suffix *below*
same-day siblings, which would misselect the newest report.

## Project Constraints

- Identifier strategy: not applicable — Run IDs, Git heads, QA report
  paths, and commit-message contracts are reused; no project-owned Internal
  Identifier is created. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — Pull Request observation goes
  through the existing `gh` boundary, read-only; no new authentication
  provider, credential policy, or HTTP route. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0015 keeps QA a phase inside the
  Spec Run; ADR-0051 keeps QA in its own Agent Session; ADR-0052 protects
  terminal completion; ADR-0053 keeps reconciliation explicit and
  proof-based; ADR-0057 keeps Task status Daemon-owned; ADR-0036 preserves
  the separate review-artifact commit; ADR-0080 defines the
  environment-blocked verdict rule. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-28, the maintainer expressly
  authorizes changes to exactly `.agents/skills/qa-gate/SKILL.md`,
  `skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, and
  `skills/roundfix/SKILL.md`, plus the deterministic Skill-digest fallout in
  exactly `internal/baseline/assets/setups/go-cli.json`,
  `internal/baseline/assets/setups/rust-cli.json`,
  `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`,
  and `internal/baseline/testdata/parity-corpus/v1/manifest.json`. No other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## System Architecture

Everything extends existing modules:

- `internal/spec` (`qa.go`) — verdict parsing gains the typed blocked
  counts and the consistency rule; the recency comparator is already
  correct for numeric suffixes and is unchanged.
- `internal/agent` (`spec_prompt.go`) — the QA prompt contract gains the
  Pull Request fact and the collision-safe report filename instruction.
- `internal/daemon` (`task_engine.go`) — the QA plan resolves the open Pull
  Request for the target branch, best-effort, before building the prompt.
- `internal/worktree` + `internal/cli` (`reconcile.go`, branch-integrity in
  `cli.go`) — the `superseded` classification, its `--apply` release, and
  the QA-report-only guard on automatic fast-forward integration.
- The qa-gate and roundfix Skill pairs plus user-guide docs — the verdict
  rule, naming contract, and reconcile vocabulary.

No new packages. The QA-report-only probe is one shared helper in
`internal/worktree` consumed by both the classifier and the preflight.

## Implementation Design

### Interfaces

QA report frontmatter (additive; unknown keys are already tolerated):

```yaml
verdict: pass            # pass | fail | partial (unchanged vocabulary)
rows_blocked_environment: 3   # optional, default 0
rows_blocked_finding: 0       # optional, default 0
```

`internal/spec` enforcement in `QAVerdict`:

```go
// verdict pass + rows_blocked_finding > 0        -> QAReportError (unreadable)
// verdict pass + rows_blocked_environment > 0    -> valid only with the
//   report's evidence obligation satisfied by the Skill's rule (ADR-0080);
//   roundfix trusts the verdict, validates count consistency only.
// verdict partial/fail                            -> unchanged handling.
```

The Daemon's `settleQAVerdict`, the Run outcome mapping (`pass` integrates,
everything else settles Unresolved), and `roundfix archive`'s
`verdict == pass` gate are all unchanged — the reachability fix is that the
Skill's verdict rule can now legitimately produce `pass`, with the archive
test matrix gaining the previously missing `partial` case and the new
consistency refusals.

QA prompt (`QAPromptRequest` gains one field, rendered as a labeled fact):

```go
type QAPromptRequest struct {
    // existing: SpecSlug, SpecDir, PRDPath, RunBranch, TargetBranch, UserCheckout
    PullRequest string // "#40 (owner/repo)" or "" when none is open
}
```

The Daemon resolves it best-effort from the target branch through the
existing `gh` boundary before `BuildQAPrompt`; when empty, the prompt states
that no open Pull Request exists and Pull Request journeys are
environment-blocked. The `qaGateContract` filename instruction changes to
the collision-safe scheme; the substring-pinned prompt tests update in the
same change.

Report naming contract (one scheme, both surfaces):
`qa-report-YYYY-MM-DD.md` for the first run of a day, `-NN` numeric suffixes
for same-day reruns. The Skill's `-<scope-or-build>` alternative is removed
because `parseQAReportName` ranks a non-numeric suffix below every sequenced
same-day sibling, which would misselect the newest report; the comparator
itself is untouched.

QA-report-only probe (shared helper in `internal/worktree`):

```go
// qaReportOnlyBranch reports whether every commit in target..run matches the
// Daemon QA-commit contract for run.SpecSlug and every changed path lies
// under the Spec's qa/ directory.
func qaReportOnlyBranch(ctx, git gitRunner, gitRoot, targetHead, runHead, slug string) (bool, error)
```

It parses `git log --format` output against `QACommitMessage`'s
deterministic shape (subject `docs: qa report for <slug> (<verdict>)`,
`Roundfix-Spec:` trailer) and `git diff --name-only` against
`docs/specs/<slug>/qa/` (and the archived path).

### Data Models

- No Run Database schema change. `IntegrationReconciliation` records the
  new classification string; `terminalRunSnapshot` mirrors any new field so
  the apply-time freshness proof still compares complete state.
- Reconcile classification set gains `superseded` with its reason string;
  `roundfix-reconcile/v1` output adds the classification value and a
  `superseded` summary counter — additive, schema version unchanged.

### API Contracts

- `inspectTerminalRun`: at the existing ancestry-miss branch, before
  settling `unintegrated`, classify `superseded` when (i) the Run is a
  terminal Implement Run, (ii) `qaReportOnlyBranch` proves the branch holds
  only QA-report commits for the Run's Spec, and (iii) the target branch's
  newest QA report for that Spec — compared with the existing recency
  comparator via `git ls-tree` of the Spec's qa/ directory at the target
  head — is newer than the branch's report. Any unprovable step falls back
  to `unintegrated`/`preserve`.
- `--apply` releases `superseded` Runs: the apply revalidation's
  `safe`/`released` allowance extends to `superseded`, re-proving the
  classification against fresh heads before cleanup, and the release reason
  records the superseding report path.
- Branch Integrity Preflight: pending Run Branches that prove QA-report-only
  are excluded from automatic fast-forward integration (today they would be
  silently merged onto the user's head branch, dragging a failing report
  in); the refusal lists them separately with
  `superseded QA report — release with: roundfix reconcile --apply` instead
  of a `git merge --ff-only` command. Genuinely pending task-work branches
  keep the current listing and integration behavior.

## Coverage Map

- Goal 1 / Stories 1–2 → frontmatter counts + `QAVerdict` consistency,
  Daemon Pull Request fact, qa-gate Skill verdict rule and read-only
  journey instructions (ADR-0080).
- Goal 2 / Story 5 → typed counts in report + Skill report template.
- Goal 3 / Stories 3–4 → `superseded` classification, `--apply` release,
  supersession-aware next actions, QA-report-only guard in preflight
  (Features 4–7).
- Goal 4 / Story 6 → `qaGateContract` filename instruction + Skill naming
  section (Feature 8).

## Integration Points

- GitHub — read-only `gh` calls to resolve the open Pull Request for the
  target branch and for the Skill-guided observation journeys; no mutation.
- Git — all supersession proof through the classifier's existing
  `gitRunner`; no new git surfaces.

## Testing Approach

Existing seams throughout: `internal/spec/qa_test.go` (verdict matrix gains
consistency refusals and the environment-blocked `pass`),
`internal/cli/archive_test.go` (adds the missing `partial` case and the new
`pass`-with-environment-counts case),
`internal/agent/spec_prompt_test.go` (contract substrings and the Pull
Request fact), `internal/worktree/worktree_test.go` (superseded
classification table beside the existing inspect tests, plus the shared
probe), `internal/cli/cli_test.go` reconcile and branch-integrity tables
(superseded release, refusal wording, no auto-ff for QA-report-only
branches), and `internal/daemon/task_engine_test.go` (Pull Request
resolution is best-effort and failure-tolerant). The Skill contract tests
gate the qa-gate pair edit.

## Build Order

1. Typed blocked counts and consistency enforcement in `internal/spec`,
   with archive-gate test coverage including `partial`.
2. QA prompt: Pull Request fact plumbed from the Daemon plan, and the
   collision-safe filename instruction (depends on: 1 for the report
   contract vocabulary).
3. Shared QA-report-only probe and the `superseded` classification with
   `--apply` release (no dependency on 1–2).
4. Branch Integrity Preflight: exclude QA-report-only branches from
   auto-integration and render supersession guidance (depends on: 3).
5. Docs, qa-gate and roundfix Skill pairs (verdict rule, journeys, naming,
   reconcile vocabulary), and the authorized derived digest pins (depends
   on: 1–4).

## Risks & Considerations

- **Verdict trust boundary.** Roundfix validates count consistency, not
  evidence quality — the evidence obligation stays with the Skill-governed
  Agent and the human reading the report, as it does today for `pass`
  itself. ADR-0080 makes that explicit rather than new.
- **Supersession false negatives are safe.** Any unparseable commit,
  missing target report, or older target report falls back to `preserve`;
  the failure mode is the status quo, never deletion.
- **Auto-integration behavior change.** Excluding QA-report-only branches
  from fast-forward integration changes a silent behavior users may have
  relied on; the release note and Skill docs state it, and the branches
  remain releasable explicitly.
- **Prompt contract tests.** The `qaGateContract` string is
  substring-pinned by four tests; the wording change and tests move in one
  commit to keep the gate deterministic.

## Decisions

- The blocked-cause distinction lives in frontmatter scalars; no markdown
  row parser is built. See ADR-0080 for the verdict rule.
- The naming contract narrows to numeric suffixes; the recency comparator
  is left untouched.
- Supersession is proven from the QA-commit message contract plus the
  recency comparator over target-side reports — Git evidence only, per
  ADR-0053; no store schema change.
- The qa-gate Skill edit in step 5 also folds in the one-pass
  authorization-audit reporting that Spec 0054's choreography work asks of
  the gate, since this Spec owns the only authorized qa-gate mutation.
