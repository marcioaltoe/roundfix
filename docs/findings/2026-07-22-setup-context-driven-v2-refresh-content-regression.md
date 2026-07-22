---
status: done
created_at: 2026-07-22
updated_at: 2026-07-22
---

# Setup Context-Driven v2 refresh — replaced rule bodies instead of extending them, plus a rule-migration defect (2026-07-22)

This investigation compared, in `/Users/marcio/dev/fluxus` (branch
`ma/skills-update`, uncommitted working tree, decisions confirmed 2026-07-22):

- the **v2 generated corpus** (`AGENTS.md` + `docs/agents/*.md`, all setup
  markers now `version=2`/`version=3`) and the **v2 asset catalog**
  (`.agents/skills/setup-context-driven/assets/`);
- the **v1 generated corpus and catalog** (commit `84da566` / `HEAD`);
- the **pre-setup repository-authored corpus** (commit `c96a896`: 575-line root
  `AGENTS.md` plus repo-authored `docs/agents/{spec-routing,autonomous-work,
  triage-labels,...}.md`);
- the reference sample
  `/Users/marcio/dev/skills/samples/AGENTS.typescript-bun.md` (573 lines).

It extends
`2026-07-21-setup-context-driven-coverage-audit-and-skill-restoration.md`
(routed to Spec 0043). The v2 refresh clearly implemented part of that
finding's §1 — `docs/agents/agent-instructions.md` v2 gained portable rules
(`rule.core.verification-selected`, `rule.core.verification-integrity`,
`rule.core.git-delivery`, `rule.core.security-configuration`,
`rule.core.dependency-discipline`, `rule.core.english-guidance`,
`rule.core.research-authority`, `rule.core.skill-dispatch`), and
`docs/agents/skill-dispatch.md` and `docs/agents/monorepo.md` now exist. Those
gains should be kept. But the v1→v2 architecture change — moving guide bodies
out of template prose (`assets/templates/guides/*.md`, now reduced to
`{{artifact.rules}}`) into single-bullet `guidance` strings in
`assets/modules/*.json` — **replaced** prose without a migration-fidelity
check, and one guide lost its content entirely.

## 1. `guide.spec-routing` lost all routing guidance — prose→rule migration miss

- **Symptom / evidence**: generated `docs/agents/spec-routing.md` v2 and
  `docs/agents/issue-tracker.md` v2 have byte-identical rule bodies ("Keep each
  Spec under `docs/specs/<feature-slug>/`. Dependencies live only in
  `_tasks.md`…"). Nothing in the generated corpus now says how to **route** a
  change, while the root block still promises "Route … local Specs through
  `docs/agents/spec-routing.md`".
- **Root cause (template layer)**: the v2 refresh deleted the prose of every
  `assets/templates/guides/*.md` in favor of `{{artifact.rules}}`. For every
  other guide the prose was migrated into a module rule `guidance` string; for
  `spec-routing` it was not. The deleted template prose held the only copy of
  the routing rules (large initiative → idea/PRD/techspec; ordinary feature →
  PRD or techspec before tasks when decisions are needed; refactor/bug fix →
  focused techspec when behavior changes; trivial change may skip a Spec only
  when intent, acceptance criteria, and verification are obvious).
- **Root cause (module layer)**: in `assets/modules/spec-workflow.json` v2,
  both `guide.spec-routing` and `guide.issue-tracker` list the same single
  rule `rule.spec.artifacts`; no routing rule exists in the v2 catalog.
- **Root cause (test layer)**: `tests/test_spec_triage_decisions.py` asserts
  only that `docs/agents/spec-routing.md` exists or is absent per decision — no
  test asserts guide content, and `assets/contract-v1.json`'s module contract
  checks rule-ID **uniqueness** but nothing forbids two guides rendering the
  same rule set.
- **Action / suggestion**: add a `rule.spec.routing` carrying the entry-point
  matrix and the trivial-change exemption, map `guide.spec-routing` to it, and
  keep `rule.spec.artifacts` on `guide.issue-tracker` only. Add catalog-level
  regression checks: (a) no two supporting guides in a module may render an
  identical rule list; (b) a prose→rule migration test comparing each v1
  template's normative clauses against the v2 rule set (see §2).

## 2. v2 prose→rule compression dropped normative clauses without an audit trail

Comparing each deleted v1 template's prose with the v2 `guidance` string that
replaced it, these obligations vanished (setup-owned regressions, not
sample-vs-setup gaps):

- **`guide.docs-layout` v2→v3** (`rule.context.docs-layout`): the findings
  lifecycle was reduced to a status enum plus `done` semantics. Lost: the
  required frontmatter template, the definitions of `partial` and `deferred`
  (both required a **recorded reason**, and `partial` a link to the covering
  Spec), "append evidence and routing links instead of rewriting
  observations", and "update `updated_at` on every status/evidence change".
  These are the contract the finding→Spec routing workflow parses.
- **`guide.agent-instructions`**: lost "record the commands that prove each
  acceptance criterion" and "keep follow-up work out of the current slice"
  (the slice rule now survives only inside `rule.spec.artifacts`, i.e. only
  during Task execution).
- **`guide.autonomous-work`**: "The Supervisor **never** writes feature code
  or tests" was softened to "delegates implementation code and tests", and
  "normally through Roundfix" disappeared — the generated corpus no longer
  names the delegation channel at all.
- **`guide.frontend`**: lost the "unless a stable test id is the explicit
  contract" qualifier on the internals-coupling prohibition.
- **`guide.secondbrain`**: lost the exact query contract
  (`qmd query "<question>" --all --files --min-score 0.3`) and the concrete
  consult triggers (Vortex, Tax, Visio, Gesttione) — v2 keeps only the
  generalized categories and a bare "run `qmd query`".
- **Root `AGENTS.md` v2**: every root block became a pointer-only line. That
  is acceptable **only** where the target guide carries the rule body;
  combined with §1 the spec-workflow root pointer currently resolves to
  tracker hygiene, not routing.
- **Action / suggestion**: treat v1 template prose as the coverage baseline
  for the v2 rule set: for each deleted template, require the migration to
  show every normative clause either landed in a rule `guidance` (same or new
  rule ID) or was intentionally retired with a recorded reason. Restore the
  findings-lifecycle details, the acceptance-evidence and slice rules, the
  Supervisor "never implements" phrasing, and the Secondbrain query contract
  as rule text.

## 3. Hard rules from the sample still have no semantic equivalent in the v2 catalog

Continuation of the 2026-07-21 finding §1 — partially addressed, still
incomplete. Categories from `AGENTS.typescript-bun.md` (and from the pre-setup
Fluxus root, a customized instance of the same sample) with no v2 counterpart:

- **Zero-warnings gate**: v2 says "treat every failure as blocking" but never
  states that lint **warnings** are blocking failures (`bun run lint` treats
  warnings as errors in this profile).
- **Design-contract read gate**: sample: "ALWAYS READ `DESIGN.md` before
  writing any UI code". `guide.frontend` v2 only says to "follow the
  repository-owned design contract" — it does not require reading it before UI
  work, and the root frontend block does not name `DESIGN.md`.
- **Dependent-API check before writing tests** (sample HIGH PRIORITY): absent.
- **Never guess user-answerable decisions** (AskUserQuestion / stop-and-ask
  rule): absent.
- **External research protocol**: `context7` appears in dispatch, but the
  sample's web-research fallback (`exa-web-search`, multiple varied searches)
  and the explicit "never use external research tools on local code" phrasing
  are only weakly implied by `rule.core.research-authority`.
- **Commit/PR conventions**: `conventional-commits` and `github-pr-workflow`
  are installed in the repository skill set but appear in no module's
  `requiredSkills` or `skillDispatch`; the sample's Conventional Commits /
  `cog verify` gate has no equivalent.
- **Domain skill dispatch**: verified across all 16 `assets/modules/*.json` —
  no module requires or dispatches any domain skill; a sweep for
  `react|drizzle|hono|tanstack|shadcn|tailwind|vite|turborepo|onioncry|zod|
  conventional-commits` over the module catalog matches nothing. The
  technology modules dispatch only generic skills (`context7`,
  `testing-boss`, `systematic-debugging`, `no-workarounds`,
  `domain-modeling`), so the generated `guide.skill-dispatch` cannot route
  React/TanStack/Hono/Drizzle/Zod/shadcn/Tailwind/bun/turborepo/vite/docker/
  onioncry/logtape/external-api-adapters/data-sync-workflows/
  integration-contract-testing/observability-audit work. The sample's
  keyword → skill dispatch protocol and its anti-patterns/enforcement list
  are the highest-leverage part of the legacy instructions and remain
  uncovered.
- **Action / suggestion**: extend the technology and surface modules
  (`typescript`, `bun`, `frontend`, `backend`, `monorepo`) with
  `skillDispatch` entries for the domain skills their profiles install, so
  the generated dispatch guide derives from the same catalog that pins the
  skills. Add portable rules for the zero-warnings gate, the design-contract
  read gate, the ask-don't-guess rule, and commit-convention dispatch. Keep
  genuinely project-specific mandates (Fluxus `systems/<domain>/` layout,
  OnionCry specifics, POST-only API contract) out of the skill — see §5.

## 4. `guide.skill-dispatch` renders near-duplicate entries from overlapping module catalogs

- **Symptom / evidence**: the generated `docs/agents/skill-dispatch.md` lists
  `context7` 4×, `testing-boss` 3×, `systematic-debugging` 3×, `tech-writer`
  3×, and repeats all 11 workflow skills twice, as one flat list with no
  module attribution.
- **Root cause**: `render_skill_dispatch` (`scripts/context_setup.py:5272`)
  already deduplicates exact `(skill_name, when)` pairs — the duplication is
  in the **catalog**, not the renderer. `context-workflow.json` and
  `spec-workflow.json` each declare the identical 11-skill dispatch list with
  *differently phrased* `when` clauses (e.g. `brainstorming`: "Starting
  creative feature or behavior design before implementation." vs "Clarifying a
  creative Spec direction before authoring requirements."), and
  `spec-workflow` `dependsOn` `context-workflow`, so both always co-activate
  and every near-duplicate survives dedupe. The generic-skill entries
  (`context7`, `testing-boss`, …) repeat the same way across `core`,
  `typescript`, `bun`, `backend`, and `frontend`.
- **Action / suggestion**: two options, in preference order: (a) deduplicate
  at the catalog level — a module that depends on another must not re-declare
  a dispatch entry for the same skill unless the trigger is genuinely
  distinct, enforced by a contract check for near-identical `when` clauses;
  (b) group rendering by skill ID (one bullet per skill, distinct triggers
  merged as sub-clauses, optionally attributed per module). Add a rendering
  test asserting one bullet per skill ID in the generated guide.

## 5. Repository delegation now dangles — the baseline silently poses as a full replacement

- **Symptom / evidence**: `packages/backend/AGENTS.md` (repository-authored,
  22 KB, still present) opens by delegating "skill enforcement, verification
  gate, git restrictions, coding style, commit guidelines" to the **root
  `AGENTS.md`**. The rich 574-line root was deleted in Fluxus commit `497ce45`
  (before the first apply), so the setup's preserve-outside-markers guarantee
  never saw those bytes; the generated root never restored the delegated
  categories (coding style, commit guidelines, and the domain skill matrix are
  still absent per §3). The delegation has been dangling since `84da566` and
  the v2 refresh did not detect it.
- **Related**: `docs/agents/triage-labels.md` remains absent (already recorded
  as observation §3 of the 2026-07-21 finding; not re-analyzed here).
- **Action / suggestion**: two complementary fixes in the skill:
  1. **Adoption/audit signal**: during apply and audit, scan repository
     `AGENTS.md`/`CLAUDE.md` files (root and packages) for references to
     setup-managed paths or to the root instructions, and warn when a
     referenced category has no rule in the active catalog — at minimum emit
     an info finding stating that the generated baseline is a floor, not a
     replacement, listing the repository-authored documents that still claim
     delegated coverage.
  2. **Repository-owned extension scaffold**: offer (behind a decision) a
     scaffolded, unmarked `docs/agents/repository.md` (or root section outside
     markers) where project-specific hard rules live, and have the root
     template link to it. That gives migrating repositories an explicit,
     preserved home for content like the Fluxus Secondbrain mirror contract
     (`.secondbrain-project` / `.secondbrain-export`, read-only mirror), the
     pt-BR domain-vocabulary language policy, the frontend `systems/<domain>/`
     mandates, and the backend OnionCry/controller rules — none of which exist
     anywhere in the working tree today.

## Scope boundary — extract portable invariants, do not copy the sample

Unchanged from the 2026-07-21 finding: the sample is evidence of useful rule
categories, not a template. Its `.knowledge` workspace flow, Fable/Codex
`gpt-5.5`/Opus runtime defaults, `<Project name>` architecture tree, and
standard-REST method guidance are obsolete or conflict with confirmed Fluxus
decisions (`Codex with gpt-5.6-sol`, POST-only action contract). §1, §2 and §4
are setup-owned defects fixable in the module catalog, templates, and
contract checks; §3 is catalog coverage; §5 is apply/audit behavior plus a
Fluxus follow-up to restore its repository-authored content in a repo-owned
document.

## Routing

Pending. Candidate owners: extend
[Spec 0043 — Context-Driven Baseline coverage and Repository Skill Set restoration](../specs/0043-context-driven-baseline-coverage-and-skill-restoration/_prd.md)
with §1–§4 (they are regressions/gaps in the same v2 catalog that Spec 0043
produced), and route §5's Fluxus-side content restoration as a separate
repository task in Fluxus itself.

Routed 2026-07-22: Spec 0043 is archived, so §1–§4 plus §5's setup-owned
delegation audit signal and extension scaffold are owned by
[Spec 0044 — Context-Driven Baseline upgrade retention and formatter compatibility](../specs/0044-upgrade-retention-and-formatter-compatibility/_prd.md);
§5's Fluxus-side content restoration remains a Fluxus repository task.
