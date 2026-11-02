---
spec: 0338-the-release-decision-travels-with-the-spec
prd: _prd.md
created: 2026-11-01
---

# The release decision travels with the Spec — Technical Spec

## Executive Summary

Two release facts decided at planning time are lost before the release. The
Delivery Queue's `pullRequestTitle` hardcodes `feat: deliver`, and the
Release Plan Command asks for approval of every minor increment even when the
delivered Spec's authorization record already granted `release`. The design
adds the Delivery Type: an optional PRD `type` (`feat`, `fix`, `perf`,
`refactor`) that `spec check` validates, the archive copies into the Archive
Record and the queue puts in the title. It adds the Release Ceiling: two
optional authorization keys, `release_ceiling` and `release_paths`, that the
authorization reader validates, and an optional Git-source capability through
which the release plan reads the records of the Specs a range delivers and
answers `approval_required` when every condition of ADR-0349 holds. Backlog
Entries become planning evidence so plan commits stop forcing manual
classification. The Baseline's release and authorization clauses, the
authoring skills and the Roundfix skill state the rules.

The trade-off accepted: a range is answered only when every delivered Spec
consents with a sufficient ceiling and every impact-setting change is such a
delivery, so a single direct `feat:` commit makes the plan ask as before.

## Project Constraints

- Identifier strategy: not applicable — no identifier or slug is minted; the
  new keys are optional and additive. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git reads only. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0143 keeps the plan read-only,
  ADR-0143: "A generic release request must never silently authorize a
  version decision"; ADR-0309 owns the title shape and is amended; ADR-0302
  (cited by ADR-0309) keeps the numbered form; ADR-0256 keeps Lagging
  Surfaces blocking, ADR-0256: "A lagging surface or an unreadable map makes
  the new"; the plan reads no `source_revision`, which after a squash
  survives only on the pull request head; ADR-0340 gives
  `release_paths` its path rules, ADR-0340: "A Spec's `_authorization.md` may
  carry `executables:`, a list of exact repository-relative file paths.";
  ADR-0343 cites ADR-0340 and is untouched; ADR-0349 is new. ADR
  Alternatives below relates each. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AU answers of 2026-11-01 quoted in the authorization record; the Governed
  Paths are the skill files, Baseline modules, derived Baseline files and
  `internal/speccheck/constraints.go` the Tasks change; no Makefile,
  workflow, module or lint configuration changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0338-the-release-decision-travels-with-the-spec/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/archive.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/release.md`, `.agents/skills/roundfix/references/spec.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/prd-template.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `docs/agents/agent-instructions.md`, `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/speccheck/constraints.go`, `skills/roundfix/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## Current behavior

Read on 2026-11-01 at `dd3b2b62`; `_prd.md` → Validation records the runs.

1. `pullRequestTitle(specSlug, specTitle)` (`internal/cli/deliver_workflow.go`)
   returns `"feat: deliver " + specSlug` within 91 UTF-16 units, otherwise
   `feat: deliver <number>: <title>`; `Publication` calls it with a closure
   that reads the Archive Record's `title` only for a long slug.
2. `BuildArchiveRecord` (`internal/spec/archive_record.go`) reads `_prd.md`
   frontmatter into a struct without `type`; `archiveRecordMapping` writes
   optional keys such as `pull_request` only when set; `ParseArchiveRecord`
   reads keys from a map and ignores unknown ones.
3. `speccheck.Check` (`internal/speccheck/constraints.go`) runs the PRD
   detectors, `detectClarificationRecord` among them; nothing reads the PRD's
   `type`.
4. `parseAuthorizationFrontmatter` (`internal/authorization/authorization.go`)
   ignores unknown keys; `release` is a valid operation that no caller
   consumes.
5. `releaseplan.Build` resolves the range, classifies commits with
   `ClassifyChanges` and calls `CalculateProposal`; a minor or major impact
   gives `approval_required`. `isPlanningEvidencePath` lists `docs/specs/`,
   `docs/adr/`, `docs/findings/` and `docs/handoffs/`, so a commit that also
   touches `docs/backlog/` is a blocking commit.
6. `releasePlanGitSource` (`internal/cli/releaseplan_git_source.go`) loads
   commits with `rev-list --reverse` and `diff-tree --name-only`; the command
   prints `Approval required:` and the next action, then the checks.
   `follow_up_status.go` builds the same plan for the follow-up mode.
7. `clause.core.plan-the-release-first` (`core.json`) says a generic release
   request authorizes only a conclusive patch plan;
   `clause.spec.project-constraints-02-tooling-authorization`
   (`spec-workflow.json`) lists the record's keys without a ceiling.

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/spec`: `DeliveryType` and its parser; the Archive Record's `type`.
- `internal/speccheck`: `SC-PRD-TYPE`.
- `internal/cli/deliver_workflow.go`: the title's type.
- `internal/authorization`: the two record keys.
- `internal/releaseplan`: `ReleaseGrant`, `ReleaseGrantSource`,
  `AnswerApproval`, the planning-evidence path.
- `internal/cli/releaseplan_*`: the grant lookup and its rendering, the help.
- Baseline modules `core` and `spec-workflow` and their derived files.
- The Roundfix, `write-prd`, `write-techspec` and `write-tasks` skills, the
  user guides, `CONTEXT.md`, ADR-0309 and ADR-0349.

## Implementation Design

### Interfaces

```go
// internal/spec/delivery_type.go

// DeliveryType is the Conventional Commits type a Spec delivers as.
type DeliveryType string

const (
	DeliveryTypeFeat     DeliveryType = "feat"
	DeliveryTypeFix      DeliveryType = "fix"
	DeliveryTypePerf     DeliveryType = "perf"
	DeliveryTypeRefactor DeliveryType = "refactor"
)

// DeliveryTypes returns the supported set in this order.
func DeliveryTypes() []DeliveryType

// ParsePRDDeliveryType reads the optional `type` key of a PRD's frontmatter.
// present is false when the key is absent; an unsupported or non-string
// value returns UnsupportedDeliveryTypeError.
func ParsePRDDeliveryType(prd []byte) (deliveryType DeliveryType, present bool, err error)

type UnsupportedDeliveryTypeError struct{ Value string }

// internal/spec/archive_record.go
type ArchiveRecord struct {
	// ...existing fields unchanged...
	Type string // optional `type` key, written only when the PRD declares it
}

// internal/cli/deliver_workflow.go
func pullRequestTitle(specSlug string, deliveryType spec.DeliveryType, specTitle func() (string, error)) (string, error)

// internal/authorization/authorization.go
type AuthorizationRecord struct {
	// ...existing fields unchanged...
	ReleaseCeiling string   // "" | patch | minor | major
	ReleasePaths   []string
}

// internal/releaseplan/release_grant.go

// ReleaseGrant is the release authorization of one Spec a range delivers.
type ReleaseGrant struct {
	Spec      string // the delivered Spec's slug
	CommitSHA string // the delivering commit
	Record    string // repository-relative record path
	Revision  string // the resolved first parent the record was read at
	Granted   bool   // an operative record lists release and a ceiling
	Ceiling   Impact // ImpactNone when not granted
	Paths     []string
	Detail    string // why the record grants no ceiling; empty when Granted
}

// ReleaseGrantSource is an optional GitSource capability.
type ReleaseGrantSource interface {
	ReleaseGrants(context.Context, Range, []Commit) ([]ReleaseGrant, error)
}

// ReleaseAuthorization is what the plan read and decided.
type ReleaseAuthorization struct {
	Answered bool
	Reason   string
	Grants   []ReleaseGrant
}

// AnswerApproval applies ADR-0349 to a plan built from changes.
func AnswerApproval(plan Plan, grants []ReleaseGrant) Plan

// internal/releaseplan/model.go
type Plan struct {
	// ...existing fields unchanged...
	ReleaseAuthorization *ReleaseAuthorization // nil: no delivered record carries a ceiling
}
```

### Data Models

- `_prd.md` frontmatter: optional `type: feat | fix | perf | refactor`.
- Archive Record: optional `type` key, written after `delivery_commit` and
  only when the PRD declares one; older records parse unchanged.
- `_authorization.md` frontmatter: optional `release_ceiling` (string) and
  `release_paths` (list of exact repository-relative paths).
- Release plan JSON (`roundfix.release-plan/0.0.1`, unchanged version):
  optional `releaseAuthorization` object `{answered, reason?, grants: [{spec,
  commitSHA, record, revision, granted, ceiling?, paths, detail?}]}`, omitted
  when `Plan.ReleaseAuthorization` is nil.

### API Contracts

1. API Contract: the Delivery Type. `ParsePRDDeliveryType` reads the
   frontmatter's `type`; absent → (`""`, false, nil); one of the four →
   (value, true, nil); anything else, `null` included →
   `UnsupportedDeliveryTypeError` whose message is
   `PRD type "<value>" is not one of feat, fix, perf, refactor`.
2. API Contract: `SC-PRD-TYPE`. `speccheck.Check` calls a new
   `detectPRDDeliveryType(result, repoRoot, prdPath)` right after
   `detectClarificationRecord`; an unsupported value adds one error finding,
   code `SC-PRD-TYPE`, summary the error message, location `_prd.md`, fix
   "Set the PRD's type to feat, fix, perf or refactor, or remove it to
   deliver as feat." An absent key or a supported value adds nothing.
3. API Contract: the Archive Record. `BuildArchiveRecord` calls
   `ParsePRDDeliveryType` on the PRD it already reads; an error refuses the
   archive with that error; a present value sets `Type`.
   `archiveRecordMapping` writes `type` after `delivery_commit` when `Type` is
   non-empty; `ParseArchiveRecord` reads `type` when present. A PRD without
   `type` gives the same record bytes as before.
4. API Contract: the title. `pullRequestTitle` uses `<type>: deliver ` for
   every prefix where it used `feat: deliver ` (`<type>` = `feat` when
   `deliveryType` is empty), with the same 91-unit budget measured on the new
   prefix, the same cut rules and the same lazy title read. `Publication`
   reads the Archive Record once through the existing `archivePaths` and
   `ReadArchivedSpec` before calling it, passes `Record.Type`, and returns
   `plan publication: read Spec type: <err>` when the record cannot be read
   and `plan publication: <UnsupportedDeliveryTypeError>` for an unsupported
   value; the title closure reuses the record read.
5. API Contract: the record keys. `parseAuthorizationFrontmatter` reads
   `release_ceiling` with `authorizationScalar` and `release_paths` with the
   same exact-path rules as `executables`. `classifyAuthorizationRecord`, for
   an approved record, refuses with code `contradictory` and field
   `release_ceiling` when a ceiling is present and `release` is not in
   `operations`; with code `malformed_record` and field `release_ceiling` when
   the value is not `patch`, `minor` or `major`; with code `malformed_record`
   and field `release_paths` when a ceiling has no paths or paths have no
   ceiling. A record without either key classifies as before.
6. API Contract: finding the grants. `releasePlanGitSource.ReleaseGrants`
   runs `git diff-tree --no-commit-id --name-status --no-renames -r <sha>`
   for each commit; a commit delivers `<slug>` when it has a `D` line for
   `<dir>/<slug>/_authorization.md` and an `A` line for a path whose base name
   is `<slug>.md`. For each delivery it calls `spec.ReadAuthorization` with
   `RecordPath` `<dir>/<slug>/_authorization.md`, `Revision` `<sha>^`, role
   `spec` and `AskingSpec` `<slug>`. `Granted` is true when the outcome is
   granted, `release` is permitted and the ceiling is set; otherwise `Detail`
   is `record <outcome>: <reason code>` for a refusal, `record grants no
   release` or `record grants release without a ceiling`. Commits without a
   delivery produce no grant.
7. API Contract: answering. `Build` calls `ReleaseGrants` only when the
   proposal's state is `approval_required` and the source implements
   `ReleaseGrantSource`, then returns `AnswerApproval(plan, grants)`.
   `AnswerApproval` returns the plan unchanged when no grant is `Granted`.
   Otherwise it sets `ReleaseAuthorization{Grants: grants}` and checks, in
   this order, recording the first failure as `Reason`:
   1. `Classification.Breaking` → `the plan is breaking`;
   2. `Classification.ManualReason != ""` → `the plan was classified by hand`;
   3. a grant not `Granted` → `<spec> grants no release ceiling`;
   4. a ceiling below `Classification.Impact` →
      `impact <impact> is above the <ceiling> ceiling of <spec>`;
   5. a change whose `AutomaticImpact` equals `Classification.Impact`, that
      crosses the maintenance boundary and is not a grant's commit →
      `commit <sha12> sets the impact outside a delivered Spec`.
   With no failure, `Answered` is true, `State` becomes `ready`, and
   `Approval` becomes `{Required: false, Increment, ProposedVersion}` with no
   question; `ProposedVersion` and `Classification` keep their values.
8. API Contract: rendering. When `ReleaseAuthorization` is nil, text and JSON
   keep their bytes. Otherwise, after the `Next action:` line and before the
   checks, text prints `Release authorization: answered` or
   `Release authorization: not answered — <reason>`, then one line per grant
   `- <spec> ceiling=<ceiling|none> record=<record>@<revision12> paths=<comma-joined paths|none>`;
   JSON adds `releaseAuthorization`. The exit code follows the state, so an
   answered plan exits 0 unless a skill-coverage check blocks it.
9. API Contract: planning evidence. `isPlanningEvidencePath` adds
   `docs/backlog/`.
10. API Contract: the help. `commandUsage("release plan")` adds, after the
    decision states, the paragraph: "approval_required is answered when every
    Spec the range delivers has an approved authorization record whose
    release_ceiling is at or above the impact, nothing is breaking or
    classified by hand, and every change that sets the impact is such a
    delivery; the decision is then ready and the plan names the records."
11. API Contract: the Baseline sentences.
    - `clause.core.plan-the-release-first` appends: "A Release Ceiling the
      maintainer recorded in the `_authorization.md` of every Spec the
      release delivers answers a minor or major plan at or below that
      ceiling, never a breaking or manually classified one." Every existing
      sentence keeps its bytes. `rule.core.git-delivery` rises by one.
    - `clause.spec.project-constraints-02-tooling-authorization` inserts,
      after the operations list "(`implement`, `commit`, `push`,
      `pull_request`, `merge`, and `release`),", the words "the impact
      ceiling a `release` grant answers without a second approval
      (`release_ceiling`: `patch`, `minor` or `major`) with the version-only
      paths that release may change (`release_paths`),".
      `rule.spec.project-constraints` rises by one.
12. API Contract: the skills and guides.
    - `write-prd`: the template's frontmatter gains
      `type: feat # feat | perf for a feature Spec; the Delivery Type`, and a
      section `## Delivery Type and release decision` says the PRD carries
      `type`, and that when the maintainer asks for a release with the Spec,
      the skill asks the release question (ceiling and release paths) with
      the structured question tool and records `release`, `release_ceiling`
      and `release_paths` in `_authorization.md` with the verbatim answer.
    - `write-techspec`: the minimal PRD frontmatter adds `type` (`fix` for a
      bug fix, `refactor` for a refactor), and the same release section.
    - `write-tasks`: a section `## Release decision` says a Task never
      performs the release, the QA scope row reads `release_ceiling` and
      `release_paths` as grants of the later release, not of a Task, and a
      record whose ceiling is missing its paths is a decomposition refusal.
    - The Roundfix skill: `references/deliver.md` replaces `feat: deliver`
      with `<type>: deliver` and names the Delivery Type; `references/spec.md`
      lists `SC-PRD-TYPE`; `references/archive.md` names the record's `type`;
      `references/release.md` and the `SKILL.md` `## Release planning`
      section state when records answer, the `Release authorization:` lines,
      `releaseAuthorization`, Backlog Entries as planning evidence, and that a
      `refactor:` delivery is a commit to classify by hand.
    - The guides `docs/user-guide/commands/deliver.md`, `spec.md`,
      `archive.md`, `docs/user-guide/release-runbook.md` and `usage.md` say
      the same. `CONTEXT.md` adds **Delivery Type** and **Release Ceiling**.
      ADR-0309 gains an amendment naming the Delivery Type.

### Surface Transcripts

None. The changed lines are asserted by the named tests; no whole-command
transcript is declared.

### Invariants

```text
1. A PRD without type archives to the same record bytes and delivers with the same title bytes as before.
2. Every queue title is within 91 UTF-16 units and classifies as its Delivery Type in the release plan.
3. A record without release_ceiling or release_paths classifies exactly as before.
4. A plan whose range delivers no record with a ceiling has the same state, text, JSON and exit code as before.
5. A breaking, manually classified, above-ceiling or outside-change plan stays approval_required.
6. The release plan never writes and never reads Roundfix configuration.
```

### Declared breaks

- `Publication` reads the Archive Record before building the title, so the
  missing-record error of `TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug`
  (`internal/cli/deliver_queue_subject_test.go`) changes from
  `plan publication: read Spec title:` to `plan publication: read Spec type:`;
  the assertion is updated, the test name kept.
- `pullRequestTitle` gains a parameter; the two title tests pass `""` and keep
  their cases.
- A commit changing only planning paths and `docs/backlog/` is now
  maintenance-only; `TestMaintenanceOnly` gains the `docs/backlog/` case.
- A range whose records answer reports `ready` and exits 0 where it reported
  `approval_required` and exited 3; no record carries a ceiling today.

## Coverage Map

- Goal 1, Goal 2, Core Features 1 and 2, Success Metrics 1-3 →
  `ParsePRDDeliveryType`, `SC-PRD-TYPE`, the Archive Record and
  `pullRequestTitle` (API Contracts 1-4), Invariants 1 and 2.
- Goal 3, Core Feature 3, Success Metric 4 → the record keys (API Contract 5),
  Invariant 3.
- Goal 3, Goal 4, Core Feature 4, Success Metrics 5 and 6 →
  `ReleaseGrants`, `AnswerApproval`, the rendering, the planning-evidence path
  and the help (API Contracts 6-10), Invariants 4-6.
- Goal 5, Core Feature 5, Success Metrics 7 and 8 → the Baseline sentences and
  the skills and guides (API Contracts 11 and 12).

Deviations from the PRD: None.

## Absence Evidence

1. A Delivery Type reader — command: `grep -rn 'DeliveryType' internal/spec internal/cli/deliver_workflow.go`; result: no match; the title prefix is the constant `"feat: deliver "`.
2. A PRD type finding — command: `grep -rn 'SC-PRD-TYPE' internal/speccheck`; result: no match.
3. Release ceiling keys — command: `grep -rn 'release_ceiling\|release_paths' internal docs/agents`; result: no match.
4. A release grant reader — command: `grep -rn 'ReleaseGrant\|AnswerApproval' internal/releaseplan internal/cli`; result: no match; `Build` ends with `CalculateProposal`.
5. Backlog as planning evidence — command: `grep -n 'docs/backlog/' internal/releaseplan/classifier.go`; result: no match.

## ADR Alternatives

- ADR-0143: followed; its rejected alternative, a plan that mutates or a
  generic request that decides a version, stays rejected: the plan still
  writes nothing and only a recorded, dated, per-Spec decision answers.
- ADR-0309: amended; its rejected alternatives (always numbering, cutting the
  slug, reserving ten digits, renaming the pull request after it opens) stay
  rejected: only the type before `: deliver` changes.
- ADR-0302: unaffected; the numbered form it introduced keeps its rules.
- ADR-0256: followed; a Lagging Surface still blocks an answered plan, and its
  rejected alternative of a release without the skill check stays rejected.
- ADR-0340: followed; `release_paths` reuses its exact-path rules and its
  rejected glob form stays rejected.
- ADR-0343: unaffected; it cites ADR-0340 for settle's executable check,
  and its rejected count-based acceptance stays rejected.
- ADR-0349: written by this Spec.

## Integration Points

None external. The plan and the queue read local Git objects.

## Testing Approach

All tests are hermetic: temporary Git repositories, `t.TempDir()` homes, no
network, no forge.

- `internal/spec/delivery_type_test.go`:
  `TestPRDDeliveryTypeReadsTheFrontmatter` (absent, each of the four, `docs`,
  `null`, a list) and `TestArchiveRecordCarriesThePRDDeliveryType` (a PRD with
  `type: fix` builds and renders `type: fix` after `delivery_commit` and
  parses back; a PRD without it renders no `type` key; `type: docs` refuses).
- `internal/speccheck/delivery_type_test.go`: `TestPRDTypeFindingNamesTheValue`
  (`docs` gives one `SC-PRD-TYPE` error; each supported value and an absent
  key give none).
- `internal/cli/deliver_queue_subject_test.go`: the existing cases pass `""`;
  `TestQueuePullRequestTitleFollowsTheDeliveryType` checks each type for a
  fitting slug and the numbered form within 91 units for a long slug, and
  that `releaseplan.ClassifyCommit` of each title plus ` (#575)` reports that
  type; `TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug` gains a
  record with `Type: "fix"` giving `fix: deliver 0296: …`.
- `internal/authorization/release_ceiling_test.go`:
  `TestReleaseCeilingIsReadAndValidated` (well-formed pair granted with both
  values; the four refusals with their code and field; a record without the
  keys unchanged).
- `internal/releaseplan/release_grant_test.go`:
  `TestAnswerApprovalFollowsEveryCeiling` with in-memory plans for each
  condition of API Contract 7 and its reason text, and the nil case.
- `internal/releaseplan/classifier_test.go`: `TestMaintenanceOnly` gains a
  `docs/backlog/` path.
- `internal/cli/releaseplan_release_ceiling_test.go`:
  `TestReleasePlanAnswersApprovalFromTheReleaseCeiling` builds a repository
  with tag `v0.1.0`, a plan commit adding `docs/specs/0001-x/_prd.md`,
  `_authorization.md` and moving a Backlog Entry, then a `feat: deliver
  0001-x (#2)` commit deleting the Spec folder, adding
  `docs/history/specs/0001-x.md` and changing `internal/x.go`. Subtests:
  ceiling `minor` (exit 0, `Decision: ready`, `Proposed version: v0.2.0`,
  `Release authorization: answered` and the grant line; JSON
  `releaseAuthorization.answered` true); ceiling `patch` (exit 3, reason
  `impact minor is above the patch ceiling of 0001-x`); no ceiling (the
  record lists `release` without `release_ceiling`: exit 3, and the text and
  JSON equal, byte for byte, those of the same repository whose record omits
  `release`, with no `Release authorization` line and no
  `releaseAuthorization` key); a direct `feat:`
  commit (exit 3, `sets the impact outside a delivered Spec`); a breaking
  footer (exit 3, `the plan is breaking`); `--impact minor --reason x` (exit
  3, `classified by hand`).
- `internal/baseline/release_clause_test.go`:
  `TestTheReleaseClauseLetsACeilingAnswerTheApproval` asserts the new core
  sentence in the embedded module and the formatter golden, and the new
  spec-workflow words in the module and the `spec-routing.md` golden.
- Existing suites: `./internal/releaseplan`, the `ReleasePlan` tests of
  `./internal/cli`, `./internal/authorization`, `./internal/spec`,
  `./internal/speccheck`, `./internal/baseline`, `./skills`,
  `-tags docscontract ./internal/docscontract`.
- Sabotages, each run on a temporary copy:
  - `pullRequestTitle` ignores `deliveryType`:
    `TestQueuePullRequestTitleFollowsTheDeliveryType` fails on `fix`;
  - `archiveRecordMapping` drops `type`:
    `TestArchiveRecordCarriesThePRDDeliveryType` fails;
  - the `release` operation check removed from the record classifier:
    `TestReleaseCeilingIsReadAndValidated` fails on the contradictory case;
  - the ceiling comparison inverted: `TestAnswerApprovalFollowsEveryCeiling`
    fails on the patch-ceiling case and the CLI subtest "ceiling patch"
    reports `ready`;
  - condition 5 removed: the CLI subtest "a direct feat commit" reports
    `ready`;
  - `ReleaseAuthorization` set for a plan without a ceiling: the CLI subtest
    "no ceiling" fails on the byte comparison;
  - `docs/backlog/` removed from the planning paths: the CLI subtest "ceiling
    minor" reports `manual_classification_required`.

## Build Order

1. The Delivery Type: parser, `SC-PRD-TYPE`, Archive Record `type`, title,
   tests, the deliver, spec and archive guides and references, `CONTEXT.md`,
   the ADR-0309 amendment, the Roundfix skill version and fingerprints last
   (task_01).
2. The Release Ceiling and the plan: record keys, grants, `AnswerApproval`,
   rendering, help, planning evidence, tests, the release guides and
   references, the Roundfix `SKILL.md` release section, `CONTEXT.md`,
   ADR-0349, the Roundfix skill version and fingerprints last (task_02)
   (depends on: 1).
3. The Baseline sentences: modules, module versions, `make baseline-digests`,
   the Managed Refresh twice, the clause test (task_03) (depends on: 2).
4. The authoring skills: `write-prd`, `write-techspec`, `write-tasks`, their
   versions (task_04) (depends on: 3).
5. QA gate (task_05) (depends on: 4).

The Tasks run in series: task_01 and task_02 both edit the Roundfix skill,
its version record, `CONTEXT.md` and the Behavior Surface Record; task_03
names in its clauses the behavior task_02 builds; task_04 teaches the keys
task_02 validates and edits the owned-skill version record.

Files Specs 0339 and 0340, authored in parallel and delivered after this one,
may also change, each kept to a narrow edit here so they rebase cleanly:
`skills/testdata/owned-skill-versions.json`, `.agents/skills/roundfix/SKILL.md`
and its mirror (version fields), `docs/references/behavior-surfaces.json`,
`CONTEXT.md`, `internal/baseline/module-versions.json`, the catalog files, the
Standard TypeScript profile, goldens and plan characterization goldens,
`docs/agents/setup-context.json` and the managed guides, and the
`write-tasks` skill. This Spec's functions are `pullRequestTitle`,
`Publication`, `BuildArchiveRecord`, `archiveRecordMapping`,
`ParseArchiveRecord`, `parseAuthorizationFrontmatter`,
`classifyAuthorizationRecord`, `releaseplan.Build`, `isPlanningEvidencePath`,
`printReleasePlanText` and `releasePlanJSONFromPlan`; a later Spec that
touches one rebases on it.

## Risks & Considerations

- A `refactor:` delivery that changes released code is a commit the plan asks
  to classify by hand, so a refactor Spec's release still asks once; the
  release reference says to rerun with `--impact patch --reason`.
- A repository that merges deliveries without squash has no commit that both
  deletes the record and adds the Archive Record on the first-parent line;
  its plan finds no grant and asks as before.
- A record refused at the parent revision (for example a path that became a
  symlink) gives no grant and the plan asks; the grant line names the reason.
- `releasePlanGitSource` gains one `diff-tree` per commit for plans that would
  ask; ready and no-release plans run no extra Git command.

## Glossary

- adds: **Delivery Type** — the Conventional Commits type a Spec delivers as
- adds: **Release Ceiling** — the highest impact a Spec's release grant answers
- not a term: **Delivery Queue** — existing glossary term, unchanged
- not a term: **Archive Record** — existing glossary term, unchanged

## Decisions

- The PRD `type` is validated at authoring by `SC-PRD-TYPE`, and through the
  Archive Record it reaches the title, so `deliver start`, which runs the
  strict check, refuses a bad value before any Run.
- The Release Ceiling is two flat record keys beside `operations`, validated
  by the authorization reader, refused as a whole when malformed.
- The release plan reads the records through an optional source capability,
  only for plans that would ask, at the first parent of each delivering
  commit.
- `docs/backlog/` becomes planning evidence (ADR-0349).
- Secondbrain: `wiki/index.md` read; the `qmd query` named in `_prd.md` →
  Validation found only mirrors and a release-cadence page, none deciding
  either rule. Exa: the Conventional
  Commits specification and the Cocogitto commit guide, quoted in `_prd.md` →
  Acceptance evidence.
- Four implementation Tasks in series and a QA gate.
