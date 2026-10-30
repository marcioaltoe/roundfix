---
spec: 0335-follow-up-takes-every-documented-type-and-an-absent-spec-root-is-empty
prd: _prd.md
created: 2026-10-30
---

# The follow-up command takes every documented message type, an absent Spec Root counts as empty, and the follow-up guides name the JSON fields — Technical Spec

## Executive Summary

Three defects, two seams. The Spec Root resolver in `internal/config` refuses
a root that does not exist, and every command that only enumerates Specs
inherits that refusal, although `spec.ListActiveDetailed` already reads a
missing root as empty. The resolver gains a variant that reports an absent
in-repository root as a value instead of an error, the existing resolver
keeps its bytes by delegating to it, and only the enumerating readers switch
to the variant. In `internal/cli`, the `--messages` parser hard-codes four
accepted types and the type fallback of the Follow-up Priority hard-codes its
own; both read one table, the error lists it, and a documentation contract
test in the existing `internal/cli` docscontract file holds both follow-up
guides to that table and to the JSON tags of the report types.

The trade-off accepted: two resolver entry points instead of one, so that more
than twenty named-Spec callers keep their precise refusal untouched, at the
cost of each enumerating caller choosing the variant explicitly (ADR-0348).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme is minted or
  changed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added; no Verification or QA row reaches the network.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0348 and ADR-0319 govern the
  design, as `_prd.md` → Project Constraints states and ADR Alternatives
  below relates. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AS decisions of 2026-10-30 quoted in the authorization record; no
  Makefile, workflow, module or lint configuration changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0335-follow-up-takes-every-documented-type-and-an-absent-spec-root-is-empty/_authorization.md`;
  bounded files:
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/follow-up.md`,
  `.agents/skills/roundfix/references/modes.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-10-30 at `35de7f74`; `_prd.md` → Validation records the runs.

1. `config.ResolveSpecsRoot` (`internal/config/config.go`) returns
   `specs.root resolved to %q, which does not exist; create the directory or
   update specs.root` when `os.Stat` of the resolved root reports not-exist,
   before it classifies `External` and `BuiltInRoot`.
2. Callers, by kind (sweep: `grep -rn 'ResolveSpecsRoot(' internal cmd`
   outside tests, 28 call sites):
   - Enumerating, refused today: `readFollowUpSpecs`
     (`internal/cli/follow_up_inventory.go`, list and status views),
     `runSpecCheckCommand` without a slug (`internal/cli/spec_check.go`),
     `runDeliverPlan` without a slug (`internal/cli/deliver_plan.go`),
     `loadDeliveryCommand` (`internal/cli/deliver.go`), which `deliver
     status`, `stop`, `resume` and the queue owner call without using the
     root, `runHistoryCommand` and `printHistoryPlan`
     (`internal/cli/history.go`), and `planBaselineUpdateHistory`
     (`internal/cli/baseline_update_history.go`), which reports the history
     plan `blocked`.
   - Named Spec, unchanged: `runSpecCheckCommand` and `runDeliverPlan` with
     slugs, `runDeliverStart` and `runDeliverRetry` (through
     `loadDeliveryCommand`), `runSpecAuditCommand`, `runSpecJudgeCommand`,
     `implement`, `settle`, `archive` (two sites and `archive_reader.go`),
     `cancel`, `reopen`, `supersede`, `deliver_revalidate.go`,
     `deliver_workflow.go` (two sites), `specaudit.Audit`, and `reconcile`
     (its carry-forward needs Specs; without carry-forward it already ignores
     the error).
   - Tolerant already, unchanged: `runs causes` (stats first), `attach`
     (falls back), the review-artifact sites in `cli.go` (non-default roots
     only), `upgrade`'s history notice (skips silently) and the interactive
     `implement` picker in `buildInteractiveInputRequest`, which refuses with
     "no active Specs to implement" for an empty root anyway.
3. `followUpMessages` accepts `type` empty (stored as `message`) or `feat`,
   `fix`, `perf`, `refactor`, and otherwise fails with `invalid type %q`.
4. `followUpOrder` (`internal/cli/follow_up.go`), without a judge answer,
   gives class `health` to `fix`, `perf` and `Kind == "finding"`, `feature`
   to `feat`, and `other` to the rest.
5. The JSON tags of `followUpDocument`, `followUpSource`, `followUpJudge`,
   `followUpItem` and `followUpStatusDocument` are the report's fields; the
   skill reference calls item fields "class", "reference", "requester",
   "blocking project" and the judge's `cost_usd` "cost", and the user guide
   lists six of the thirteen item fields.

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/config/config.go`: `SpecsRoot.Absent`, `ResolveSpecsRootOrAbsent`
  and `SpecsRoot.RequirePresent` (task_01).
- `internal/cli`: the enumerating readers of Current behavior 2 call the
  variant; `runSpecCheckCommand` and `runDeliverPlan` with slugs and
  `runDeliverStart` and `runDeliverRetry` call `RequirePresent` (task_01).
- `internal/cli/follow_up_inventory.go`, `internal/cli/follow_up.go`: the
  message type table, its parser use, its error and its ordering use
  (task_02).
- `internal/cli/live_documents_docscontract_test.go`: the guide contract
  (task_02).
- Guides and the Roundfix skill: the configuration guide, four skill
  references and `CONTEXT.md` describe the absent root (task_01); the
  follow-up guide and reference name the fields and types, and the modes
  reference tells the agent to write only documented types (task_02).

## Implementation Design

### Interfaces

```go
// internal/config/config.go
type SpecsRoot struct {
	Path        string
	External    bool
	BuiltInRoot bool
	Absent      bool // no filesystem entry; nearest existing ancestor inside the repository
}
func ResolveSpecsRootOrAbsent(loaded Loaded, repoRoot string) (SpecsRoot, error)
func (root SpecsRoot) RequirePresent() error // nil unless Absent
func ResolveSpecsRoot(loaded Loaded, repoRoot string) (SpecsRoot, error) // = OrAbsent + RequirePresent

// internal/cli/follow_up_inventory.go
var followUpMessageTypes = []struct{ Name, Class string }{ /* sorted by Name */ }
func followUpTypeClass(typ, kind string) int // 2 health, 3 feature, 4 other
```

### Data Models

None. The `roundfix/follow-up/v1` and `roundfix/follow-up-messages/v1`
schemas keep their fields; only the accepted `type` values grow.

### API Contracts

1. API Contract: absent Spec Root resolution. `ResolveSpecsRootOrAbsent`
   resolves the path as `ResolveSpecsRoot` does today. When `os.Lstat` of
   the resolved path reports not-exist, it walks to the nearest ancestor
   that `os.Lstat` finds, evaluates its symbolic links, and returns
   `SpecsRoot{Path: resolved, BuiltInRoot: <lexical rule of today>,
   Absent: true}` with a nil error when that ancestor is inside or equal to
   the evaluated repository root; otherwise it returns today's
   `which does not exist` error. A path whose `os.Lstat` finds an entry
   keeps every current branch: a symbolic link whose target is missing, a
   file and a stat error are refused with today's text. `RequirePresent`
   returns, for an absent root, the error `ResolveSpecsRoot` returns today,
   byte for byte; `ResolveSpecsRoot` returns that error, so its callers see
   no change.
2. API Contract: enumerating readers. For an absent root,
   `readFollowUpSpecs` returns no Spec and no error, so the list and status
   views report `specs` `checked` with count `0`; `spec check` and `deliver
   plan` without a slug check no Spec; `deliver status`, `stop`, `resume`
   and the queue owner load as with an empty root; `history sanitize`
   (plan and apply) and `planBaselineUpdateHistory` plan with the root's
   path and its archive root as with an empty root. For each, standard
   output, standard error and exit code equal those of the same repository
   with an empty Spec Root.
3. API Contract: named Specs. `spec check <slug>…`, `deliver plan <slug>…`,
   `deliver start` and `deliver retry` call `RequirePresent` before reading a
   Spec, so an absent root gives today's refusal text and exit code; every
   other caller of Current behavior 2 keeps `ResolveSpecsRoot`.
4. API Contract: message types. `followUpMessageTypes` holds, sorted by name:
   `feat` → `feature`, `feature` → `feature`, `finding` → `health`, `fix` →
   `health`, `information` → `other`, `other` → `other`, `perf` → `health`,
   `refactor` → `other`. The parser accepts an empty `type` (stored as
   `message`, as today) or a name from the table, and keeps the name as the
   item's `type`. Any other value fails with `invalid type %q; use one of
   <names joined by ", ">, or omit type`, built from the table, wrapped in
   today's `messages file %q: ` prefix, exit 2 and the usage line.
5. API Contract: the type fallback. Without a judge answer, `followUpOrder`
   takes an item's class from `followUpTypeClass`: the table's class for a
   listed type, `health` for `Kind == "finding"`, and `other` for any other
   type (`message`, `note`, `spec`, an unlisted Backlog type). Every item the
   current code classifies keeps its class.
6. API Contract: the documented fields. Both follow-up guides carry these
   sentences, each followed by its backticked names in the JSON tag order of
   the named type and ending with a period: "The list report carries these
   JSON fields:" (`followUpDocument`), "The status report carries these JSON
   fields:" (`followUpStatusDocument`), "Each source carries these JSON
   fields:" (`followUpSource`), "The judge carries these JSON fields:"
   (`followUpJudge`), "Each item carries these JSON fields:"
   (`followUpItem`), and "A message type is empty or one of" followed by the
   table's names in order. The skill reference no longer says "class",
   "reference", "requester" or "blocking project" for item fields.

### Surface Transcripts

None. No command gains a flag, an argument or an exit code; the new error
line is fixed in API Contract 4 and asserted whole by
`TestFollowUpRefusesAnUnknownMessageTypeNamingTheAcceptedOnes`, and API
Contract 2 is asserted as equality with the empty-root run rather than as a
fixed transcript.

### Invariants

```text
1. ResolveSpecsRoot returns the same value or the same error bytes as before for every input.
2. An absent Spec Root is never reported as unavailable, refused or blocked by an enumerating reader.
3. A dangling Spec Root link, a file, an unreadable root and a missing root outside the repository stay refused.
4. archive and cancel write nothing in the Spec Root beyond what they write today.
5. The parser, its error, the ordering fallback and both follow-up guides read or state one type list.
6. Every JSON field name a follow-up guide states is a JSON tag of the type its sentence names, in order, and every tag is stated.
```

## Coverage Map

- Goal 1, Core Feature 1, Success Metrics 1–2 → API Contracts 4 and 5:
  the `followUpMessageTypes` table, the `followUpMessages` parser and its
  error, and `followUpTypeClass` in `followUpOrder`.
- Goal 2, Core Features 2–3, Success Metrics 3–4 → API Contracts 1 to 3:
  `ResolveSpecsRootOrAbsent` and `SpecsRoot.RequirePresent` in
  `internal/config`; the enumerating readers `readFollowUpSpecs`,
  `runSpecCheckCommand` and `runDeliverPlan` without slugs,
  `loadDeliveryCommand`, `runHistoryCommand`, `printHistoryPlan` and
  `planBaselineUpdateHistory`; and the `RequirePresent` guards of the
  named-Spec readers.
- Goal 3, Core Feature 4, Success Metric 5 → API Contract 6: the field and
  type sentences of both follow-up guides, held by
  `TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes`.
- Core Feature 5 → the modes reference (task_02).
- Success Metric 6 → declared unreachable; the guides tell adopters the
  accepted types.

Deviations from the PRD: None.

## Absence Evidence

1. An absent Spec Root value — command: `grep -rn 'Absent\|OrAbsent\|RequirePresent' internal/config/config.go`; result: no match; the only not-exist branch returns an error.
2. A message type table — command: `grep -rn 'feat.*fix.*perf.*refactor\|item.Type == .feat.' internal/cli/follow_up_inventory.go internal/cli/follow_up.go`; result: one literal type switch in each file, no shared list.
3. A guide field contract — command: `grep -rn 'follow-up.md' internal/cli/*docscontract_test.go internal/docscontract`; result: no test reads either follow-up guide.

## ADR Alternatives

- ADR-0348: implemented as recorded, including its rejection of an archive placeholder, of a tolerant resolver for every caller and of reading a missing external root as empty.
- ADR-0319: refined; the four questions, thresholds and classes stay, and only the type fallback reads the shared table, giving every previously accepted type its old class.

## Integration Points

None. Git, `gh` and the judge are reached as today; tests use temporary
repositories, the follow-up fixture's fake `gh` and fake judge.

## Testing Approach

All tests are hermetic: temporary repositories and homes, the follow-up
fixture's fakes, no network, no real `~/.roundfix`.

- task_01, `internal/config/specs_root_absent_test.go`:
  `TestResolveSpecsRootOrAbsentReadsAnAbsentRootInsideTheRepository`
  (built-in and relative roots, `docs` itself missing: `Absent` true, no
  error, `BuiltInRoot` as today; `ResolveSpecsRoot` on the same input
  returns today's error text),
  `TestResolveSpecsRootOrAbsentStillRefusesAMissingExternalRoot` (an
  absolute root under another temporary directory),
  `TestResolveSpecsRootOrAbsentStillRefusesADanglingLink` (`docs/specs` a
  symbolic link to a missing directory), and
  `TestResolveSpecsRootOrAbsentStillRefusesAFile`.
  `internal/cli/absent_spec_root_test.go`:
  `TestAnAbsentSpecRootReadsLikeAnEmptyOne` (subtests `follow-up`,
  `follow-up --status`, `spec check`, `deliver plan`, `deliver status`,
  `history sanitize`: run without `docs/specs`, create it empty, run again,
  compare stdout, stderr and exit, and assert the first run never prints
  `which does not exist`),
  `TestBaselineUpdateHistoryPlansAnAbsentSpecRootLikeAnEmptyOne`
  (`planBaselineUpdateHistory` reports the same status and message for
  both), and `TestANamedSpecInAnAbsentSpecRootKeepsTheRefusal` (`spec check
  <slug>`, `deliver plan <slug>`, `deliver start <slug>`: exit 2 and the
  `which does not exist; create the directory or update specs.root` text).
  The existing `TestResolveSpecsRoot*`, `TestFollowUp*`, `TestDeliverPlan*`,
  `TestHistorySanitize*` and `TestSpecCheckWithoutStageIsUnchanged` consume
  the changed resolver and run in the same Verification.
- task_02, `internal/cli/follow_up_test.go`:
  `TestFollowUpAcceptsEveryDocumentedMessageType` (one message per table
  name plus one without `type`, judge unavailable: exit 0, each item's type
  kept and its class as the PRD's Success Metric 1 states, asserted with
  literal expectations) and
  `TestFollowUpRefusesAnUnknownMessageTypeNamingTheAcceptedOnes` (type
  `chore`: stdout empty, exit 2, stderr exactly the error of API Contract 4
  built from the table, then the usage line);
  `internal/cli/live_documents_docscontract_test.go`:
  `TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes` (reads
  `docs/user-guide/commands/follow-up.md`,
  `.agents/skills/roundfix/references/follow-up.md` and its mirror, extracts
  the names after each sentence of API Contract 6 and compares them in order
  with the JSON tags by reflection and with the table).
  No test is removed or renamed, so the coverage record is not re-recorded.

## Build Order

1. The absent Spec Root resolution, its enumerating and named-Spec callers,
   their tests, the configuration guide, the four skill references and the
   glossary (task_01).
2. The message type table, its parser, error and ordering use, the
   follow-up guide and reference fields, the modes reference and the guide
   contract test (task_02) (depends on: 1 — both edit
   `internal/cli/follow_up_inventory.go`, the follow-up reference, the
   Roundfix `SKILL.md` pair, the owned-skill record and the Behavior Surface
   record).
3. QA gate (task_03) (depends on: 1, 2).

## Risks & Considerations

- A wrongly widened absence could hide a broken knowledge-workspace link:
  the `Lstat`-first rule keeps a dangling link refused, and a test pins it.
- `history sanitize` and the Baseline history plan use the Spec Root's
  relative path as a prefix when they look for live citations and legacy
  folders (`planHistoryUnit`, `planLinksUnit`, `planHistoryRewrite`); an
  absent root must give what an empty one gives, and the equality tests of
  task_01 are what prove it for the plans; `--apply` reads the same plan.
- Spec 0333 also bumps the Roundfix owned skill; the record commands, not
  hand edits, resolve the versions after the rebase.
- The skills repository vendors a copy of the Roundfix follow-up reference;
  it is upstream content and follows the release, not this Spec.

## Glossary

- changes: **Spec Root**

## Decisions

- Absent means no entry at the resolved path with an in-repository nearest
  ancestor; only enumerating readers take it. See ADR-0348.
- Two entry points keep the named-Spec callers byte-identical. See ADR-0348.
- One type table drives the parser, its error and the ordering fallback; the
  docscontract test lives in the existing `//verify:always` docscontract
  file of `internal/cli`, so it reads the unexported report types and adds no
  new contract file.
