---
status: done
created_at: 2026-07-21
updated_at: 2026-07-21
---

# Setup Context-Driven — coverage, audit, and skill restoration gaps (2026-07-21)

This investigation followed a real `setup-context-driven` 0.8.1 apply in
`/Users/marcio/dev/fluxus` with the `typescript-bun-monorepo` profile. The
generated instructions were compared semantically with
`/Users/marcio/dev/skills/samples/AGENTS.typescript-bun.md`, then the generated
links, manifest, required-skill restoration path, final audit, and Fluxus gate
were checked. The sample is evidence of useful rules, not a template to copy
verbatim.

## 1. The TypeScript/Bun baseline does not replace the useful legacy instruction coverage

- **Symptom / evidence**: the sample has 573 lines and explicit sections for
  language policy, skill enforcement and dispatch, universal verification,
  Git restrictions, local-versus-external search, dependency installation,
  frontend architecture, design-system use, naming, commit/PR conventions,
  and security/configuration. The generated Fluxus corpus consists of the
  90-line root `AGENTS.md` and ten guides totaling 180 lines. Line count is not
  a quality metric by itself, but a semantic search confirmed that the
  generated corpus has no equivalent rule bodies for those categories.
- **Symptom / evidence**: the catalog knows many required skills, but generated
  instructions do not expose a dispatch contract. For example,
  `assets/modules/core.json` requires six workflow skills while defining only
  three compact rules; `frontend.json` requires `context7`, `testing-boss`, and
  `systematic-debugging` while its only generated rule concerns observable UI
  behavior. In Fluxus, `packages/backend/AGENTS.md:3` explicitly delegates
  shared skill enforcement, the verification gate, Git restrictions, coding
  style, and commit guidelines to the root `AGENTS.md`; the generated root no
  longer contains that shared contract.
- **Symptom / evidence**: `verification.gate=rtk make verify` is rendered in
  `docs/agents/autonomous-work.md`, so it applies explicitly to autonomous
  Tasks. The universal guide requires fresh evidence but does not require the
  selected repository gate for every completion claim, zero warnings, or
  preservation of lint, formatter, architecture, and test configuration.
- **Root cause**: the setup audit proves byte-level conformance to the selected
  asset catalog, not semantic completeness. The current universal,
  TypeScript/Bun, frontend, and backend modules intentionally encode a very
  small rule set, and no catalog contract states which portable instruction
  categories a profile must cover.
- **Action / suggestion**: keep the root compact, but add setup-owned supporting
  modules/guides for (1) universal safety, selected-gate enforcement, and
  verification-config integrity; (2) skill dispatch derived from module
  `requiredSkills`; (3) language, research-source, dependency, Git/delivery,
  and security rules; and (4) frontend architecture and a repository-owned
  `DESIGN.md` pointer when that surface is enabled. Define these as rule IDs and
  test profile coverage semantically instead of asserting a target line count.
  Clarify which rules remain repository-authored so adoption cannot imply that
  the generated baseline is a complete replacement when it is not.

## 2. A clean single-context audit accepts a broken managed link

- **Symptom / evidence**: generated Fluxus `AGENTS.md:35` says to read
  `docs/agents/monorepo.md`, but that file does not exist. The manifest contains
  `root.monorepo` and no `guide.monorepo`. A fresh audit returned exit `0` with:

  ```json
  {
    "ok": true,
    "summary": {"errors": 0, "decisions": 0, "warnings": 0, "info": 0},
    "findings": [],
    "plannedChanges": []
  }
  ```

- **Symptom / evidence**: `assets/modules/monorepo.json` declares both the root
  block and `docs/agents/monorepo.md`. However,
  `test_domain_layout_selects_distinct_guidance_and_audits_clean` explicitly
  expects that guide to be absent for `domain.layout=single-context` while the
  unconditional root template still links to it.
- **Root cause**: decision-driven artifact selection removes the supporting
  guide for a single-context repository without conditionally changing the
  root pointer. Audit validates managed digests and expected artifact presence,
  but not relative Markdown links among generated artifacts.
- **Action / suggestion**: either generate the monorepo guide whenever the
  monorepo root block is present or render a single-context-safe pointer. Add an
  audit pass for setup-owned relative Markdown links and regression cases for
  every decision combination that adds or removes a referenced guide.

## 3. The pre-apply preview did not account for every observed removal

- **Symptom / evidence**: before apply, Fluxus had an untracked
  `docs/agents/triage-labels.md`. With `triage.external=false`, the file
  disappeared during the setup session even though the pre-apply managed-change
  summary did not name that removal. Because the file was untracked, commit
  `84da566` cannot preserve its preimage; this finding records the session
  observation, not a reconstructed root cause.
- **Root cause**: unknown. The current evidence does not distinguish a
  conditional managed cleanup from a legacy-marker or adoption transition, but
  the missing preview entry conflicts with the skill's requirement to present
  every create, refresh, and remove operation before apply.
- **Action / suggestion**: add transition fixtures with pre-existing marked and
  unmarked conditional guides. For each fixture, compare the complete file-tree
  delta after apply with `plannedChanges`; every removal must have a path,
  managed ID, state, condition, and reason before confirmation. Preserve an
  unmarked file unless the user explicitly approves adoption or removal.

## 4. Required-skill drift has no portable, deterministic restoration path

- **Symptom / evidence**: after the docs apply, audit blocked on
  `skills.required.drift` for `agent-browser`, `baseline-ui`,
  `frontend-design`, `interface-design`, and `turborepo`. The emitted action is
  only “Refresh ... from the typescript-bun canonical skill setup.” The bundled
  setup identifies its revision as `pinned-2026-07-15`, but these five entries
  have an empty `source` object and only a path plus content digest.
- **Symptom / evidence**: restoring from the local canonical checkout produced
  the expected content but rewrote `skills-lock.json` with machine-local source
  paths. Restoring from `marcioaltoe/skills` used the current default branch and
  reintroduced drift. Supplying the observed immutable commit to `bunx skills
  add` displayed the ref but installed default-branch content and omitted the
  ref from the lock entry. The session therefore required manual reconciliation
  to content verified at commit
  `5d09cbf0bb47b325d7fb093af82371342741191c`, followed by manual portable
  GitHub provenance and hash repair.
- **Symptom / evidence**: Fluxus commit `84da566` records changes across all five
  skill directories and deletion of
  `.agents/skills/interface-design/agents/openai.yaml`. That mutation was part
  of the out-of-band drift remediation, not the setup apply preview. After the
  manual reconciliation, setup audit exited `0`, canonical directory
  comparisons passed, and `rtk make verify` passed.
- **Root cause**: bundled snapshots are sufficient to detect content drift but
  not to reproduce the required content. A date label plus digest is not an
  immutable, installer-consumable source, and setup currently depends on
  external installer behavior that did not preserve the requested revision or
  portable provenance in this run.
- **Action / suggestion**: make every external required-skill snapshot
  restorable from an immutable `source`, `ref`, and `path`, and include those
  values in the finding's machine-readable remediation. Add a disposable-repo
  restoration test that starts from drift, executes the supported command,
  verifies the complete directory digest, proves that removed files were
  previewed, and rejects absolute local paths in `skills-lock.json`. If the
  external installer cannot honor immutable refs, bundle a deterministic
  restore mechanism or block with an explicit upstream dependency instead of
  suggesting a generic refresh.

## Scope boundary — preserve the compact architecture, not the sample verbatim

Several sample rules are obsolete or conflict with current Fluxus contracts:

- the `.knowledge` sparse-checkout and separate documentation commit flow is
  retired; generated Secondbrain guidance is correctly read-only;
- Fable, Codex `gpt-5.5`, and Claude Opus 4.8 defaults differ from the confirmed
  `Codex with gpt-5.6-sol` decisions;
- `<Project name>`, PostgreSQL 18, and Better Auth are template- or
  project-specific assumptions;
- standard REST method guidance conflicts with Fluxus's explicit POST-only
  action contract and current unauthenticated, internal-only backend.

The corrective work should extract portable invariants from the sample and
leave project-specific architecture in repository-owned guides. Existing Specs
0036 and 0040 overlap with skill-source correctness and Context-Driven Baseline
evolution, respectively, but neither currently maps all observations above.
Keep this finding `pending` until a Spec explicitly owns the selected scope.

## What worked — keep

- Compact root blocks with detailed supporting guides are easier to scan and
  remain the right document architecture.
- Decision rendering correctly preserved the selected profile, single-context
  domain layout, local Specs, autonomous runtime values, verification command,
  and read-only Secondbrain guidance.
- Managed-content audit, required-skill digest detection, final idempotent
  audit, and the Fluxus repository gate all worked once the content and
  provenance were reconciled.

## Routing

All four observations are routed to
[Spec 0043 — Context-Driven Baseline coverage and Repository Skill Set restoration](../specs/0043-context-driven-baseline-coverage-and-skill-restoration/_prd.md).
Specs 0036 and 0040 retain their existing scopes for Doctor Skill Readiness and
mandatory Baseline ADRs, respectively.
