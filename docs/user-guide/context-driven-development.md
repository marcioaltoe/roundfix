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

## Configure or audit the baseline

The `setup-context-driven` skill manages the portable Context-Driven Baseline:
declared root blocks, setup-owned guides, and
`docs/agents/setup-context.json`. A Source Baseline is the immutable input
corpus and identity. Its independent Normative Clause Manifest accounts for
every Normative Clause, recommendation, and Operational Contract. The Setup
Manifest records the selected profile, modules, Repository Capabilities,
decisions, and managed artifacts.

Setup preserves unmarked repository content, does not generate
project-specific architecture, and never removes extra installed skills.
Repository-authored policy belongs in a recognized typed document or in
Repository-Specific Normative Rules, not in the portable Source Baseline.

### Standard TypeScript Monorepo Profile

Select `standard-typescript-monorepo` only when the repository uses both
required workspaces:

- `packages/frontend`
- `packages/backend`

The profile requires TypeScript, Bun, Turborepo, Vite, React, Hono, Drizzle,
Zod, Tailwind, shadcn, TanStack Query, TanStack Router, Better Auth,
PostgreSQL, LogTape, Oxlint, Oxfmt, and Vitest. Setup proves these Repository
Capabilities from local package, workspace, executable, installed-skill, and
repository-contract evidence. It does not install the stack, connect to
PostgreSQL, or execute repository scripts.

Frontend feature code uses systems with a public system boundary; code inside
one system imports its internal modules directly. Backend code uses `domain`,
`application`, and `infrastructure` layers. HTTP handlers stay thin, use cases
stay HTTP-independent, and Drizzle owns persistence. The profile does not
create generic `modules` or a `services` layer. Inngest and Docker are
optional modules.

Each repository must persist an HTTP Contract Decision with one mode:

- `REST`
- `Post-only`

The decision has no default. Each exception is ordered and records `scope`,
`methods`, `owner`, and `reason`. Setup can reuse a supported typed repository
document and its digest; otherwise the maintainer supplies the decision in a
decision file. Hono capability never determines HTTP policy.

Each Repository Capability is required, recommended, or optional. Required
Repository Capabilities block readiness when absent. A present capability with
no usable version remains an unresolved decision.
Recommended-capability gaps emit one non-blocking
`capability.recommended.missing` warning with an explanation and next action.
Optional modules activate only when local evidence or an explicit repository
contract selects them.

Context7 and Exa are required Repository Capabilities. Firecrawl, `rtk`, and
`rg` are recommended. Search local repository code and documentation first,
using `rg` when available. Use Context7 for authoritative current library,
framework, runtime, and tool documentation. If that cannot answer an external
question, use Exa for three to seven varied searches, then verify conclusions
against primary sources. Firecrawl can collect structured external material,
but external research never substitutes for local code search.

The profile installs exact Skill Activation bundles:

| Trigger ID | Bundle ID | Required skills |
| --- | --- | --- |
| `trigger.production-code` | `bundle.production-code` | `coding-guidelines`, `clean-code`, `solid` |
| `trigger.frontend.react-feature` | `bundle.frontend-react` | `react`, `react-best-practices`, `react-composition-patterns` |
| `trigger.frontend.ui-quality` | `bundle.frontend-ui-quality` | `frontend-design`, `interaction-design`, `interface-design`, `fixing-accessibility`, `wcag-audit-patterns`, `web-design-guidelines` |
| `trigger.hono.endpoint` | `bundle.hono-endpoint` | `hono-api-best-practices`, `hono`, `zod` |
| `trigger.hono.endpoint-persistence` | `bundle.hono-endpoint-persistence` | `hono-api-best-practices`, `hono`, `zod`, `drizzle-orm` |
| `trigger.testing` | `bundle.testing` | `testing-boss`, `tdd`, `vitest` |
| `trigger.debugging` | `bundle.debugging` | `systematic-debugging`, `diagnosing-bugs`, `no-workarounds` |
| `trigger.security` | `bundle.security` | `security-best-practices`, `security-threat-model` |
| `trigger.qa` | `bundle.qa` | `qa-gate`, `evidence-gate` |
| `trigger.delivery` | `bundle.delivery` | `conventional-commits`, `github-pr-workflow` |

Start with a local, read-only audit:

```bash
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py audit --repo <repo> --profile standard-typescript-monorepo --format json
```

The response uses `setup-context-driven/audit-v1`. Read the selected profile,
ordered modules, canonical setup, `capabilities`, `sourceBaseline`,
`sourceEntries`, findings, and complete `plannedChanges`. Each planned change
includes `action`, `path`, `managedId`, `state`, `reason`, `beforeDigest`, and
`afterDigest`, plus `condition`, `fromPath`, or `referenceEdits` when
applicable. Findings name `code`, `managedId`, `path`, `message`, `action`, and
structured `remediation` when available.

Resolve each `decision.required` finding one at a time. Scalar answers can use
repeatable `--decision <id=value>`. Structured decisions use a repeatable
`--decision-file <path>`. A decision file has this strict shape:

```json
{
  "schemaVersion": "setup-context-driven/decisions/0.0.1",
  "version": "0.0.1",
  "decisions": [
    {
      "id": "verification.gate",
      "value": "<verification-command>"
    },
    {
      "id": "http.contract",
      "value": {
        "mode": "REST",
        "exceptions": [],
        "source": {
          "path": "docs/architecture/http-contract.json",
          "digest": "<lowercase-sha256>"
        }
      }
    }
  ]
}
```

The ordered `decisions` array contains exact `{id, value}` records. An HTTP
Contract Decision contains only `mode`, ordered `exceptions`, and `source`.
`source.path` is repository-relative and `source.digest` binds the current
typed source. Use `Post-only` instead of `REST` when that is the repository's
policy. Conflicting scalar answers across `--decision` and decision files are
invalid input.

`verification.gate` is required for every profile and records the repository's
authoritative Verification command. Rerun audit with the complete decision
set:

```bash
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py audit --repo <repo> --profile standard-typescript-monorepo --format json --decision-file <decision-file>
```

Only a fully resolved Decision Plan has an authorizable `planDigest`. Review
the complete Change Plan, including capabilities, source entries, dispositions,
Repository-Specific Normative Rules bytes, retention accounting, and every
planned file change. Preview apply without mutation:

```bash
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py apply --repo <repo> --profile standard-typescript-monorepo --format json --decision-file <decision-file>
```

A non-empty preview emits `plan.confirmation.required`, exits `3`, and writes
nothing. Confirm only the exact returned digest:

```bash
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py apply --repo <repo> --profile standard-typescript-monorepo --format json --decision-file <decision-file> --confirm-plan <planDigest>
```

Apply recomputes the plan. If repository bytes, decisions, profile, or catalog
content changed, it emits `plan.confirmation.stale`, exits `3`, and writes
nothing. Review the recomputed plan and confirm its new digest. An empty plan
exits `0` without confirmation. Audit and apply use bundled assets only and
remain local and network-free.

After apply, run the repository's selected formatter and Verification. Then
rerun the same resolved audit and an unconfirmed reapply preview:

```bash
<formatter-command>
<verification-command>
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py audit --repo <repo> --profile standard-typescript-monorepo --format json --decision-file <decision-file>
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py apply --repo <repo> --profile standard-typescript-monorepo --format json --decision-file <decision-file>
```

The final audit must exit `0`; the reapply preview must report no changes and
exit `0` without confirmation. This audit → decisions → apply preview → exact
confirmation → apply → formatter → Verification → audit → reapply sequence is
the complete operating workflow. A clean audit proves setup-owned semantic
coverage, typed references, the authorized managed tree, and installed
Repository Skill Set digests. It does not prove repository-authored
architecture or policy completeness.

### Baseline Readoption

An older or incompatible Setup Manifest starts Baseline Readoption. Setup does
not interpret it as a clean install or permission to overwrite. The initial
audit inventories bounded root and nested `AGENTS.md`, `CLAUDE.md`, and
`docs/agents/` carriers. Every nonblank byte enters exactly one structural
Source Baseline Entry with an ID, byte range, kind, digest, and opaque bytes.
The scanner does not infer whether an entry is normative.

Copy `sourceBaseline.id`, `sourceBaseline.digest`, and every ordered
`sourceEntries` record from the fresh audit into the optional `readoption`
object in the same decision file:

```json
{
  "schemaVersion": "setup-context-driven/decisions/0.0.1",
  "version": "0.0.1",
  "decisions": [],
  "readoption": {
    "sourceBaseline": {
      "id": "<source-baseline-id>",
      "digest": "<lowercase-sha256>"
    },
    "dispositions": [
      {
        "entryId": "<source-entry-id>",
        "entryDigest": "<lowercase-sha256>",
        "classification": "normative-clause",
        "disposition": "managed-entry",
        "destination": {
          "managedId": "<managed-entry-id>"
        },
        "reason": ""
      }
    ]
  }
}
```

Map every source entry exactly once. `classification` is one of
`normative-clause`, `recommendation`, `operational-contract`, or
`non-governed`. `disposition` is one of:

- `managed-entry`: `destination` contains exactly `managedId`.
- `repository-document`: `destination` contains a supported `documentType`,
  safe repository-relative `path`, and current `digest`.
- `repository-rules`: `destination` identifies
  `docs/agents/repository-rules.md`, declares `documentType` as
  `repository-rules`, and contains canonical base64 `proposedBytes` plus their
  lowercase SHA-256 `digest`.
- `rejected`: `destination` is `null` and `reason` is non-empty.

A `non-governed` classification also requires its own reason. Missing,
duplicate, unknown, stale, unsafe, or structurally invalid records block
without writes. Source identity, ordered entries, classifications,
destinations, reasons, proposed Repository-Specific Normative Rules bytes, and
the complete Change Plan all contribute to `planDigest`.

The first confirmed `repository-rules` write creates an unmarked,
repository-owned file. Later setup runs require the typed target to exist but
never compare, format, restore, rewrite, or remove its bytes. Complete
Baseline Readoption through the same preview, exact confirmation, apply,
formatter, Verification, audit, and reapply sequence above.

### Upgrade Retention Contract

A baseline version transition is authorizable only after the Upgrade Retention Contract accounts for every previously managed mandatory clause. The Change Plan adds ordered `retentionAccounting` when the repository has a recognized source-baseline transition. Each entry reports `fromClause`, `enforcement`, a `retained`, `moved`, `replaced`, or `rejected` disposition, `targets`, and a recorded `reason`. Accepted targets must exist in the selected future artifact graph and preserve the clause's enforcement strength. This normative accounting contributes to `planDigest`, so changing only a transition mapping still invalidates an earlier confirmation.

An unknown source baseline, unaccounted clause, unreachable target, or enforcement mismatch is a blocking finding: audit and apply exit `1` without writes. Use the finding code to correct the Setup Manifest identity or the canonical transition ledger, then rerun audit and review and confirm the complete new Change Plan. Do not infer an unknown baseline or patch generated Markdown around the failure. Decisions and missing or stale confirmation remain exit `3` and use the existing confirmation flow shown above.

### The 0.0.1 version boundary

`0.0.1` is the shared identity for Roundfix-owned release and setup surfaces:

- the Roundfix CLI and its npm launcher and platform packages;
- the Context-Driven Baseline, Source Baseline assets, schemas, manifests,
  profiles, modules, decisions, templates, setup snapshots, managed artifacts,
  formatter provenance, compatibility fixtures, and managed markers;
- all Roundfix-owned distributed skills, including their canonical and embedded
  copies;
- the Release Plan JSON schema and restarted changelog.

The reset does not erase or renumber operational, upstream, or third-party
contracts:

- User Config, Project Config, Runs, Run Database rows, and existing
  operational state remain intact;
- the Run Database `PRAGMA user_version` and external `skills-lock.json` schema
  retain their operational meanings;
- upstream-managed skill content and metadata versions remain owned by their
  upstream source;
- third-party protocol versions remain unchanged;
- Git history, Specs, and accepted or partially superseded ADRs remain
  authoritative.

Historical tags and GitHub Releases are neither setup state nor permission to
mutate remote history. Their separate read-only reset workflow is documented
in the [release runbook](release-runbook.md#planning-the-001-release-history-reset).

### Repository ownership and delegation floor

The `repository.extension.enabled=true` decision can include one initial Repository-Owned Extension creation in the confirmed Change Plan. The created file has no setup markers and never enters `managedArtifacts`. Its bytes remain repository-authored and outside setup management: audit, reapply, and profile transitions do not compare, rewrite, format, or remove them. Setup manages only the typed reference from generated guidance. If a previously selected extension is missing, `reference.repository.missing` blocks rather than recreating project content.

The bounded delegation scan reads root and nested agent-instruction documents without making them setup-owned. A `delegation.baseline-floor` finding means a repository-authored document delegates to a setup-managed category absent from the active catalog. It is informational, does not affect the Change Plan or exit status, and does not block apply; it states that the generated Context-Driven Baseline is a floor, not a replacement for project policy.

### Skill dispatch and Formatter-Stable Output

Selected module contracts drive both installed skills and generated dispatch. `docs/agents/skill-dispatch.md` renders each installed skill once; one skill entry can contain multiple distinct trigger lines. A profile cannot hide an installed skill from dispatch or declare the same skill through multiple owning modules.

Generated managed Markdown is Formatter-Stable Output for the selected profile. Formatter proof is pinned and profile-specific: the TypeScript/Bun profile declares an exact Oxfmt version and golden-corpus digest, while profiles without a selected Markdown formatter declare that explicitly. Ordinary Verification is hermetic: it compares generated bytes with the checked-in pinned corpus and neither downloads nor executes the formatter. Final QA runs the real pinned formatter in a disposable TypeScript/Bun fixture, followed by that fixture's selected Verification, a fresh audit, and a second apply. Audit and apply themselves remain offline and never run a formatter or the target repository's Verification.

## Restore required external skills

Restoration is never a side effect of audit or documentation apply. Missing or drifted external skills carry structured audit remediation with an immutable GitHub source, exact commit, source path, expected complete-tree digest, and preview argv. Run restoration only as an explicit maintainer action; do not substitute a generic skill refresh.

Preview selected drift without mutating the repository:

```bash
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py restore-skills --repo <repo> --profile <profile-id> --skill <name> --format json
```

Omit `--skill` to select all missing or drifted external skills required by the profile. The preview uses schema `setup-context-driven/restore-v1`, acquires the declared exact commit, verifies each source subtree, and reports all created, refreshed, and removed files plus each `skills-lock.json` edit. It exits `3` with `plan.confirmation.required` for a non-empty plan.

Review `skills`, `acquisitions`, `plannedChanges`, and `planDigest`. Confirm that exact digest by rerunning the same command:

```bash
rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py restore-skills --repo <repo> --profile <profile-id> --skill <name> --format json --confirm-plan <planDigest>
```

Use `--source-dir <git-source>` on both preview and confirmation to use an offline Git checkout or bare object store containing the declared commit. Restoration owns Git acquisition; audit and apply never fetch. The command has no branch or default-revision fallback, preserves unrelated lock entries and extra skills, verifies the restored complete-tree digest, and rolls back all targets on an apply failure.

Exit categories are stable across the operator workflow:

- exit `0`: clean audit, successful apply or restoration, or an already-current plan;
- exit `1`: blocking audit or retention findings, or a restoration source, proof, safety, lock-adapter, or write failure;
- exit `2`: invalid arguments, decisions, confirmation digest, skill selection, or lock input;
- exit `3`: unresolved decisions, confirmation required, or stale confirmation.

Text and JSON results go to stdout; diagnostics go to stderr. The `audit` and `apply` commands never use the network. `restore-skills` can acquire only the declared immutable Git source and runs solely after an explicit maintainer request.

After restoration, rerun the resolved audit. Spec 0036 and the Doctor Command own Repository Skill Set readiness and lock-hash compatibility; this setup workflow relies on that compatibility gate and does not add Doctor behavior.

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
