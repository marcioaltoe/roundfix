# CONTEXT-driven development

CONTEXT-driven development is how this repository plans and builds: every change
flows through a pipeline of planning artifacts — idea, PRD, tech spec, tasks —
that live as local markdown under `docs/specs/<slug>/`, grounded in a shared
vocabulary (`CONTEXT.md`) and recorded decisions (`docs/adr/`). Downstream stages
read the artifacts, not the conversation, so any fresh agent session continues
from the files alone.

The method is adapted from Matt Pocock's **skills** work; see
[Source and attribution](#source-and-attribution).

## Why "CONTEXT-driven"

The name points at the artifact at the center of the method: `CONTEXT.md`, the
project glossary. Code, docs, prompts, task names, and TUI copy all draw their
vocabulary from it, so an agent and a human describe the same thing the same way.
Shared, written-down language is what keeps artifacts consistent across sessions
and models — it is the context that drives the work, rather than a single long
conversation an agent cannot reload.

Three artifacts carry that context:

- **`CONTEXT.md`** — the domain vocabulary. New terms are added as Specs
  introduce them.
- **`docs/adr/`** — Architecture Decision Records. Product and technical
  decisions are recorded here in the same authoring pass, never only in chat.
- **`docs/specs/<slug>/`** — the planning artifacts for one change, the only home
  of that work's plan.

## The pipeline

```text
write-idea -> write-prd -> write-techspec -> write-tasks -> implement -> qa-gate -> roundfix archive
```

Each stage reads and writes `docs/specs/<slug>/` and produces one artifact:

| Stage | Produces | Altitude |
| --- | --- | --- |
| `write-idea` | `_idea.md` — the opportunity, scored and debated | Why |
| `write-prd` | `_prd.md` — user stories and acceptance criteria | What and why |
| `write-techspec` | `_techspec.md` — contracts, data models, failure modes | How |
| `write-tasks` | `_tasks.md` + `task_NN.md` — a dependency-ordered Task Graph | Work units |
| `implement` | Completed Tasks, each with verification evidence | Execution |
| `qa-gate` | `qa/` report validating the feature against the PRD | Verdict |
| `roundfix archive` | The Spec stamped and moved to `docs/specs/_archived/` | Record |

Not every change runs the full pipeline. The entry point depends on the change —
a large initiative starts at `write-idea`, a standard feature at `write-prd`, a
refactor or bug fix at `write-techspec`, and a trivial fix skips the spec folder
entirely. The routing rules live in
[`docs/agents/spec-routing.md`](agents/spec-routing.md); every route converges on
`write-tasks`, so implementation always executes from a Task Graph rather than an
ad-hoc plan.

## Adopt or update the Context-Driven Baseline

The Baseline Command is the sole public authority for adopting, updating, and
verifying a Context-Driven Baseline. It manages declared root blocks,
setup-owned guides, immutable root-instruction backups, and the Setup Manifest
at `docs/agents/setup-context.json`. It preserves repository-authored content
outside managed boundaries.

One selected Baseline Profile composes Repository Capabilities. Each selected
skill is a Skill Activation. During Baseline Readoption, the Source Baseline's
independent Normative Clause Manifest accounts for every Normative Clause,
recommendation, and Operational Contract. The interactive questions and strict
Decision Documents form the Decision Plan; the projected repository mutation
is the Change Plan. Repository policy such as an HTTP Contract Decision remains
repository-owned input; implementation evidence cannot answer it.

### Instruction hierarchy

The generated root index applies only the active tiers, in this order:

1. Universal instructions.
2. Context and documentation.
3. Spec workflow.
4. Autonomous work, when enabled.
5. Stack guidance.
6. Surface guidance.
7. Optional knowledge sources.
8. Repository-Specific Normative Rules, only when accepted residual rules
   exist.

A narrower guide may add constraints for its concern but cannot weaken a
universal Normative Clause or confirmed project decision. Each operative rule
has one semantic owner, and each active guide has one root pointer. Read the
universal tier first, then continue through the active tiers that match the
work; do not use file order or a narrower guide to override earlier policy.

Run the interactive workflow from a terminal:

```bash
roundfix baseline --repo . --format text
```

The same command handles first adoption and later updates. It detects the
current state, asks numbered questions, resolves exactly one Baseline Profile,
shows one consolidated Change Plan with file changes first, and asks once
whether to apply the displayed Plan Digest. It writes nothing before that final
confirmation.

Each decision prompt and the contextual preservation prompt mark one visible
default. Press Enter to accept it. A still-valid value from the existing Setup
Manifest takes precedence, even when a changed Profile Digest requires
adoption again. Without a stored value, the embedded catalog suggests:

| Decision | Suggested value |
| --- | --- |
| Generated language | `English` |
| Verification gate | `rtk make verify` |
| HTTP contract | `Post-only` |
| Spec artifacts | Yes |
| Domain layout | `single-context` |
| External triage | No |
| Autonomous work | Yes |
| Backend runtime | `codex gpt-5.6-sol` |
| Design runtime | `claude fable high` |
| Secondbrain | Yes |
| Repository-Specific Normative Rules carrier | Permitted when non-empty |

Existing root instructions make Preservation the default; an empty instruction
inventory makes Greenfield the default. A recoverable existing profile is the
profile default. Classification, plan approval, and apply never accept an empty
answer.

The interactive command refuses redirected or absent terminal input. Scripts,
CI jobs, and Agents must use the non-interactive `baseline plan` and `baseline
apply` commands described under [Automation](#automation).

### Greenfield composition and update redistribution

Greenfield composition selects exactly one Baseline Profile and renders only
the semantic guides and root pointers selected by its modules. It does not
generate a generic repository guide. With no accepted residual rules, it also
creates no `docs/agents/specific-repository.md` carrier and no root pointer to
that path.

Update and Baseline Readoption account for every existing repository rule
before mutation. Roundfix segments the current carriers without changing their
exact source bytes, then proposes one reviewed disposition for each segment:

- move the exact bytes into a repository-owned block in the active semantic
  owner, such as the domain, backend, frontend, Spec, or documentation guide;
- retain genuinely ambiguous or unmodeled rules as Repository-Specific
  Normative Rules in `docs/agents/specific-repository.md`;
- retain a recognized typed repository document at its typed destination; or
- reject or classify non-governed evidence with an individual reason.

The Change Plan exposes each source identity and digest, proposed semantic
owner, residual disposition, root-pointer effect, before and after identity,
and Upgrade Retention Contract entry. Review that complete ledger before
confirmation. No source rule can disappear or be weakened without an explicit,
reasoned disposition.

Confirmed redistribution removes recognized legacy carriers after their exact
rules have destinations. If no residual remains, the plan removes the empty
repository-specific carrier and its root pointer. If residuals remain, the
carrier contains only those accepted exact rules. Arbitrary nested instruction
carriers stay repository-owned, byte-identical, and warning-only.

### First adoption

An unconfigured or incompatible repository enters adoption. Choose one
instruction mode:

- **Greenfield** creates content-addressed backups of safe root instruction
  carriers and imports none of their rules.
- **Preservation** creates the same backups, inventories each root Source
  Baseline Entry, and opens one consolidated classification review. Every entry
  must be retained as a managed entry, a recognized repository document, or a
  Repository-Specific Normative Rule, or rejected with its own reason.

Repository-Specific Normative Rules have one canonical carrier:
`docs/agents/specific-repository.md`. Baseline creates it and adds its managed
root pointer only when approved non-empty rules exist. Greenfield therefore
creates neither an empty carrier nor a dangling pointer. During adoption,
Roundfix safely migrates the exact non-empty bytes from either legacy
`docs/agents/repository.md` or `docs/agents/repository-rules.md` and removes an
empty legacy scaffold. Divergent non-empty legacy carriers block planning for
manual reconciliation; Roundfix never guesses which rules win.
Existing `0.0.1` Decision Documents that name either legacy destination are
normalized to the canonical path before the Plan Digest is computed.

Nested instruction carriers remain repository-owned and unchanged. Conflicts
found in them appear as warnings. Unsafe root carriers—external or escaping
links, cycles, unreadable paths, and special files—block planning and are never
followed.

If semantic classification is unavailable, invalid, too large, or incomplete,
the interactive workflow keeps the immutable snapshot and continues with the
same structured manual review. ACP output is a proposal only; deterministic
validation and the maintainer's confirmation remain authoritative.

### Update, profile change, and rejected plans

A compatible Setup Manifest enters update with its current Baseline Profile and
stored decisions. An incompatible Profile Digest enters adoption but retains
still-valid stored decisions and a resolvable profile as interactive defaults.
The workflow asks whether to keep or change the profile. Changing it produces a
new complete Change Plan and Plan Digest; profiles are never combined.

Rejecting a final plan does not authorize a partial write. Select the decision
area to revisit:

- Baseline Profile;
- Repository-Specific Normative Rules;
- repository divergences and decisions;
- projected files.

You can change the structured decision directly or submit a free-form
suggestion within Baseline scope. A valid suggestion is translated into
proposed decisions only. Roundfix recalculates the complete plan, requires the
digest to change, displays the new plan, and asks for confirmation again.

### ADR and Findings lifecycle

New ADRs prepend the repository-owned lifecycle overlay to the body contract in
`.agents/skills/domain-modeling/ADR-FORMAT.md`. Only `accepted` is active;
`proposed`, `rejected`, `deprecated`, and `superseded` are inactive. A legacy
ADR without lifecycle frontmatter stays active unless its body explicitly
marks it inactive, and Baseline never rewrites an existing ADR only to add this
metadata.

<!-- baseline-adr-lifecycle-template:start -->
```markdown
---
status: proposed # proposed | accepted | rejected | deprecated | superseded
created_at: YYYY-MM-DDTHH:MM:SSZ
updated_at: YYYY-MM-DDTHH:MM:SSZ
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# <Short title of the decision>

<One to three sentences describing the context, decision, and reason.>
```
<!-- baseline-adr-lifecycle-template:end -->

Findings use `pending`, `partial`, `deferred`, and `done`. Set `done` when the
implementation Spec is created and linked, not when implementation or release
finishes. Preserve the original investigation and append later evidence,
routing, and status context as dated addenda.

<!-- baseline-findings-template:start -->
```markdown
---
status: pending # pending | partial | deferred | done
created_at: YYYY-MM-DD
updated_at: YYYY-MM-DD
---

# <Area> — <short title> (YYYY-MM-DD)

<Two to four sentences describing the session or investigation, the attempted outcome, and links to adjacent evidence.>

## 1. <Finding title — symptom, not hypothesis>

- Symptom / evidence: <observed behavior, command output, identifiers, and paths needed to reproduce>
- Root cause: <proven cause, or `unknown` with what was ruled out>
- Action / suggestion: <fix, mitigation, or route to a Spec, direct change, or upstream report>

## 2. <Next finding>

## What worked — keep

<Optional evidence about behavior worth preserving.>

## Addendum — YYYY-MM-DD — <short title>

<Append new evidence, root-cause proof, status context, or routing links without rewriting the original observation.>
```
<!-- baseline-findings-template:end -->

Use `pending` for a new finding without an implementation Spec. Use `partial`
when the linked Spec covers only the selected implementation scope and record
why the rest is unnecessary. Use `deferred` only with a recorded reason not to
implement. Update `updated_at` for every lifecycle change or addendum and keep
`created_at` unchanged.

### Profile alignment and adaptation

The interactive workflow audits the selected Baseline Profile before
instruction classification. It prints each blocking and advisory divergence.
For profile-specific blockers, choose one explicit path:

1. **Change Baseline Profile** and audit the replacement.
2. Review a **repository-owned Profile adaptation**.
3. **Decline without writing**.

An adaptation proposal lists every module and profile-specific Repository
Capability it would remove. Review those removals, choose a valid repository
Profile ID, and let Roundfix validate and re-audit the in-memory draft. Only a
`ready` re-audit proceeds to instruction classification and the consolidated
Change Plan. The Profile file appears in that plan at
`.roundfix/baseline/profiles/<id>.json`; planning does not write it, and apply
writes it only after approval of the final Plan Digest.

Universal required capabilities cannot be removed or waived. Follow the exact
remediation reported by alignment. For example, preview a missing Context7
skill restoration:

```bash
roundfix baseline skills restore --repo . --profile standard-typescript-monorepo --skill context7 --format json
```

Review the returned restoration Plan Digest, then confirm the same current
preview:

```bash
roundfix baseline skills restore --repo . --profile standard-typescript-monorepo --skill context7 --confirm-plan <digest> --format json
```

Rerun the interactive workflow after remediation. Profile alignment must be
`ready` before classification or the final Change Plan.

### Decision Documents

Use repeatable `--decision id=value` flags for scalar answers and repeatable
`--decision-file <path>` flags for strict structured input. A Decision Document
uses `setup-context-driven/decisions/0.0.1`, rejects duplicate keys and unknown
fields, and contains an ordered array of exact `{id, value}` records.

This complete greenfield input is parser-valid. Replace the decisions with the
repository owner's answers:

<!-- baseline-decision-document:start -->
```json
{
  "schemaVersion": "setup-context-driven/decisions/0.0.1",
  "version": "0.0.1",
  "decisions": [
    {
      "id": "autonomous.enabled",
      "value": false
    },
    {
      "id": "domain.layout",
      "value": "single-context"
    },
    {
      "id": "language.generated",
      "value": "English"
    },
    {
      "id": "preservation.mode",
      "value": "greenfield"
    },
    {
      "id": "repository.extension.enabled",
      "value": false
    },
    {
      "id": "secondbrain.enabled",
      "value": false
    },
    {
      "id": "spec.scaffold",
      "value": true
    },
    {
      "id": "triage.external",
      "value": false
    },
    {
      "id": "verification.gate",
      "value": "make verify"
    }
  ]
}
```
<!-- baseline-decision-document:end -->

Preservation adds one `readoption` object containing the current
`sourceBaseline` identity and an ordered `dispositions` array. Each disposition
names the exact Source Baseline Entry and digest, one classification, one typed
destination or `null`, and an individual reason when rejected or classified as
non-governed. Start from the current workflow's emitted classification
skeleton; never copy Source Baseline identities or entry digests from another
repository state.

### Automation

Planning is read-only and non-interactive:

```bash
roundfix baseline plan --repo . --profile go-cli-tui --decision-file baseline-decisions.json --format json
```

Automation can provide the same reviewed repository-owned adaptation as an
explicit strict Profile document:

```bash
roundfix baseline plan --repo . --profile-file team-backend.json --decision-file baseline-decisions.json --format json
```

`--profile-file` and `--profile` are mutually exclusive. The draft must use the
strict `roundfix/custom-baseline-profile/v1` schema, bind the embedded catalog,
and be a valid adaptation of one built-in Profile. Roundfix resolves it in
memory and includes its canonical repository path and exact bytes in the
portable Plan; it does not write the Profile file during planning.

Write its stdout to `baseline-plan.json` using the calling shell or process.
Exit `0` emits one complete `roundfix/baseline-plan/v1` document. Exit `3`
emits a `roundfix/baseline-result/v1` next action and no partial plan. Invalid
or mutually exclusive Profile input exits `2`.

Review the complete document, especially `fileChanges`, `managedEntries`,
`retention`, `warnings`, and `planDigest`. Apply only the reviewed artifact:

```bash
roundfix baseline apply --repo . --plan baseline-plan.json --confirm-plan sha256:<64-lowercase-hex> --format json
```

`--confirm-plan` must equal the exact `planDigest` inside the supplied file.
Apply strictly parses the file, rejects duplicate or unknown fields, verifies
the document digest, validates repository lineage and every bounded preimage,
and applies only the supplied postimages. It never recalculates or substitutes
a newer plan.

Run the reported formatter and repository Verification recommendations outside
Baseline. Repair any failure outside managed guidance. Generate a fresh plan
with the same `--profile-file` and decisions; require no file changes, or an
exact idempotent reapply at exit `0`, before reporting the adaptation current.

Plans use repository-relative paths and clone-stable Git root lineage. A plan
can move to another clone only when its Git object format, root commits,
catalog identity, selected profile, and every bounded preimage still match:

```bash
roundfix baseline apply --repo <matching-clone> --plan baseline-plan.json --confirm-plan sha256:<64-lowercase-hex> --format json
```

Treat a plan as a sensitive review artifact: it contains repository policy,
bounded preimage evidence, and exact generated bytes. Transport it through the
same access-controlled channel used for source changes.

### Automation schemas

`roundfix/baseline-plan/v1` is the portable approval artifact:

| Field | Contract |
| --- | --- |
| `schemaVersion` | Exact value `roundfix/baseline-plan/v1` |
| `repository` | Clone-stable Git identity |
| `catalog` | Embedded catalog schema and digest |
| `profile` | One resolved Baseline Profile |
| `decisions` | Normalized owner decisions |
| `retention` | Complete Upgrade Retention Contract ledger |
| `preimages` / `postimages` | Exact bounded states approved for replacement |
| `warnings` | Non-blocking repository findings |
| `setupManifest` | Proposed Setup Manifest |
| `managedEntries` | Canonical ordered change ledger |
| `fileChanges` | Derived one-row-per-path projection |
| `planDigest` | Digest binding the canonical plan and exact postimages |

`roundfix/baseline-result/v1` is returned when no complete plan exists and by
apply:

| Field | Contract |
| --- | --- |
| `schemaVersion` | Exact value `roundfix/baseline-result/v1` |
| `operation` / `state` | Operation and stable outcome state |
| `category` | Failure or next-action category when applicable |
| `message` / `nextAction` | Human-readable outcome and deterministic recovery action |
| `planDigest` | Relevant digest when one exists |
| `verifiedPostimages` | Exact verified writes; always an array |
| `warnings` | Structured findings; always an array |
| `recommendations` | Commands reported but not executed; always an array |

Requested text or JSON goes to stdout. Diagnostics and progress go to stderr.
JSON mode returns a result document on stdout once argument parsing has selected
the Baseline contract.

| Exit | Category |
| --- | --- |
| `0` | Plan emitted, apply verified, profile operation succeeded, no-op reapply, restoration complete, or assets current |
| `1` | Execution, apply, output, Baseline verification, recovery, asset refresh, or incomplete rollback failed |
| `2` | Invalid arguments or schema, unsafe repository or target, or Preflight Validation failed |
| `3` | A decision, manual classification, renewed approval, current confirmation, or fresh preimage is required |
| `130` | Operation canceled |

### Profiles

Roundfix ships `go-cli-tui`, `rust-cli`, and
`standard-typescript-monorepo`. Inspect a built-in or repository-owned Baseline
Profile before selecting it:

```bash
roundfix baseline profile show go-cli-tui --format json
roundfix baseline profile validate go-cli-tui --format text
```

Create a repository-owned profile from one built-in profile:

```bash
roundfix baseline profile init --id team-go-cli --from go-cli-tui
roundfix baseline profile validate team-go-cli --format text
```

Repository-owned profiles live only at
`.roundfix/baseline/profiles/<id>.json`. They can reference only embedded
catalog entry IDs; they cannot compose profiles, load a user-scoped catalog,
declare executable content, or fetch remote content. Review and version the
profile with the repository.

The output keeps three concepts separate:

- A **profile expectation** is a portable role or Repository Capability the
  selected profile expects. It is not proof that the repository implements it.
- A **repository command** is executable only when Roundfix found and
  digest-bound a supported local declaration. The Baseline Command reports it
  but does not run it.
- A **recommendation** is a next command for the maintainer. It is not an
  executed action or readiness claim.

### Repository Skill Set restoration

Baseline adoption reports missing or drifted external skills but never restores
them as a side effect. Preview the explicit operation:

```bash
roundfix baseline skills restore --repo . --profile go-cli-tui --skill <skill-name> --format json
```

A non-empty preview exits `3` and returns an exact Plan Digest. Confirm that
same current preview:

```bash
roundfix baseline skills restore --repo . --profile go-cli-tui --skill <skill-name> --confirm-plan <64-lowercase-hex> --format json
```

Use `--source-dir <path>` for an offline Git checkout or bare object store that
contains the declared immutable commit. Restoration validates provider,
repository, commit, source tree, complete-tree digest, target safety,
`skills-lock.json` compatibility, and the complete preimage. It updates only
the selected skill trees and lock records through the recoverable transaction.

### Canonical asset synchronization

Asset synchronization is a maintainer operation, not part of repository
adoption. Check an explicit immutable canonical setups checkout:

```bash
roundfix baseline assets sync --source-dir <canonical-setups> --check --format json
```

Without `--check`, Roundfix validates the generated catalog in memory and
updates only `internal/baseline/assets/setups` through the recoverable
transaction. It never installs skills, writes to the canonical source, or reads
an installed setup skill as runtime authority.

### Recovery and troubleshooting

| Outcome | Recovery |
| --- | --- |
| Missing decision or manual classification | Follow `nextAction`, complete the strict Decision Document or structured human review, then rerun planning. No partial plan exists. |
| Confirmation mismatch | Review the supplied plan and copy its exact `planDigest`; do not edit the plan or shorten the digest. |
| Stale plan | Repository bytes, catalog identity, profile, or decisions changed. Run `baseline plan` again, review the complete replacement, and approve its new digest. |
| Cross-clone or lineage mismatch | Use a clone with matching root commits and object format, or generate and approve a plan in the target repository. |
| Unsafe carrier | Repair or remove the reported escaping, external, cyclic, unreadable, or special-file carrier, then rerun planning. Roundfix never follows it. |
| Interrupted transaction | Rerun the same public apply operation. Before validating a new write, Roundfix locks and recovers the Git-private transaction journal. Do not delete recovery state. |
| Incomplete rollback | Stop further Baseline writes, preserve the checkout and Git-private recovery state, inspect every reported path, restore trustworthy bytes from version control or the recorded preimages, then generate a new plan. Do not bypass the journal or confirmation. |
| Repository formatter or Verification fails | Repair the repository outside Baseline, run the command again, then generate a fresh plan or perform an exact empty reapply as appropriate. |

An exact reapply against an already verified postimage is an idempotent exit
`0`. Baseline verification covers only Baseline-owned outcomes: managed
postimages, immutable backups, carrier relationships, Setup Manifest identity,
retention accounting, and resolved Baseline state.

### Security and execution limits

Audit and planning are local, read-only, and network-free. Baseline never:

- installs application or Agent dependencies;
- connects to PostgreSQL or live infrastructure;
- executes repository formatter, build, test, lint, migration, or Verification
  commands;
- follows unsafe links or mutates nested instruction carriers;
- infers Normative Clauses, HTTP policy, exception ownership, or exception
  rationale from implementation evidence;
- lets ACP output authorize repository mutation;
- combines Baseline Profiles or loads user-scoped or remote profiles.

Semantic classification and rejected-plan revision use a fresh sealed ACP
Session with no tools or terminal, denied permissions, one turn, and a
two-minute deadline. Canonical input is limited to 2 MiB and 256 entries;
accepted output is limited to 512 KiB. Limit, proof, schema, tool-use, timeout,
or cleanup failure discards the proposal and returns to deterministic manual
review.

After apply, run reported formatter and Verification recommendations yourself.
Then generate a fresh plan. No file changes means the Baseline is current; do
not treat an advisory profile expectation as an executed command.

### Migrate from the script-backed setup skill

The public Baseline Command replaces direct setup-script operation. Existing
strict `setup-context-driven/decisions/0.0.1` Decision Documents remain useful
input, but old audit or preview output is not a portable Baseline Plan.

1. Install a Roundfix version that exposes `roundfix baseline --help`.
2. Stop invoking setup scripts or treating skill assets as runtime authority.
3. Run `roundfix baseline` for the interactive migration, or generate a new
   `roundfix/baseline-plan/v1` document with `baseline plan`.
4. Review the file projection, complete managed-entry and retention ledgers,
   warnings, repository identity, and Plan Digest.
5. Apply only that exact document and digest through `baseline apply`.
6. Run the reported repository commands yourself, then require a fresh empty
   plan or exact idempotent reapply.

A compatible Setup Manifest enters update. An incompatible manifest enters
Baseline Readoption and requires explicit preservation or greenfield handling;
it is never treated as permission to overwrite. There is no independent skill
engine or script fallback. If the public command refuses, follow its
`nextAction` or repair the reported state.

## How Roundfix executes it

Roundfix is the runtime for the implementation half of the pipeline. `roundfix
implement --spec <slug>` executes the Task Graph as one Run — Tasks in dependency
order, each gated by its own Verification commands — and `roundfix archive`
closes the loop after `qa-gate` passes. The operational flow is documented in the
[usage guide](usage.md); the autonomous role split (an orchestrator authors
Specs, an ACP Runtime implements them) is in
[`docs/agents/autonomous-work.md`](agents/autonomous-work.md).

The result is that planning and execution share one artifact contract: the same
`docs/specs/<slug>/` files a human reads are what the agent implements from and
what QA checks against.

## Source and attribution

This method is adapted from **Matt Pocock's skills** —
<https://github.com/mattpocock/skills> — described there as "Skills for Real
Engineers. Straight from my .claude directory." Pocock's approach organizes
agent work as composable skills built on software-engineering fundamentals:

- **Alignment before execution** — resolve gaps between developer and agent
  before writing code.
- **Domain-driven language** — build shared terminology through `CONTEXT.md` and
  ADRs to keep agents aligned. This is the origin of the "CONTEXT" in
  CONTEXT-driven development.
- **Tight feedback loops** — test-driven development and fast verification to
  validate quality.

Roundfix adapts those ideas into this repository's pipeline (`write-idea`,
`write-prd`, `write-techspec`, `write-tasks`, `implement-spec`/`implement-task`,
`qa-gate`, `archive-spec`) and ships the authorial skill bundle in the binary.
The skills are maintained as an adaptation at
[`marcioaltoe/skills`](https://github.com/marcioaltoe/skills) and pinned through
`skills-lock.json`; the wording, task contract, and Roundfix integration are this
repository's own. Credit for the underlying method belongs to Matt Pocock.
