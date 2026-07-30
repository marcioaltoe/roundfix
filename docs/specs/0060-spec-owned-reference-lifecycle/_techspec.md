---
spec: 0060-spec-owned-reference-lifecycle
prd: _prd.md
created: 2026-07-30
---

# Spec-owned reference lifecycle — Technical Spec

## Executive Summary

This Spec ships no Go code. Every behavior is a contract written into the
authorial Skills and the repository-authored sections of two agent guides, and
its enforcement is the Skills' own gates. That makes the central design question
not "where does the logic live" but "which stage owns each transition, and what
evidence does its gate read".

Three seams change. `write-prd` gains the adoption transition — the inventory,
the `git mv`, the finding status flip, the reference index, and the link
rewrite, all in one reviewable change. `archive-spec` gains a self-containment
precondition it can verify with `git` and `grep` alone, alongside the two
preconditions it already verifies. `write-idea`, `write-techspec`, and
`write-tasks` gain only the referencing convention, so a downstream artifact
never links a source at its pre-adoption path.

The primary trade-off is enforcement strength. A Skill gate is an instruction an
agent follows, not a compiled check — weaker than `internal/spec`. Moving the
validation into Go would make it binding, but the lifecycle is a
documentation-authoring workflow that runs before and after any Roundfix Run, so
a Go gate would only see the result and could not perform the move. We keep
enforcement in the Skills and make the gate mechanical instead: each check is a
named command whose output is the evidence, so a skipped check is visible in the
report rather than merely claimed.

The second trade-off is the index format. A machine-readable index would let a
future Go gate validate provenance, but nothing reads it today and an unread
schema drifts. We specify a Markdown table with fixed columns — greppable,
diffable, and cheap to promote to a parsed format later if a consumer appears.

## Project Constraints

- Identifier strategy: not applicable — documents keep their basenames and Git
  history; no project-owned Internal Identifier is created. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — documentation lifecycle only.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0083 makes adopted sources move to
  their one owning Spec with Git history, new promotions only. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the PRD authorizes exactly the five authorial
  Skill pairs (`write-idea`, `write-prd`, `write-techspec`, `write-tasks`,
  `archive-spec`, each in `.agents/skills/` and `skills/`) plus the Skill-digest
  fallout in exactly `internal/baseline/assets/setups/go-cli.json`,
  `internal/baseline/assets/setups/rust-cli.json`,
  `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`. The
  repository-authored sections of `docs/agents/docs-layout.md` and
  `docs/agents/spec-routing.md` are documentation, not protected tooling. No
  other protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.
- Durable knowledge flows upstream only: applicable — this Spec's own
  `_prd.md` links two findings, so it is its own first subject; the adoption it
  specifies applies to new promotions only and this Spec is not retrofitted.
  Source: `docs/agents/specific-repository.md`.

## System Architecture

- `.agents/skills/write-prd/SKILL.md` + `skills/write-prd/SKILL.md` — the
  adoption transition, inserted as a new step between "Record decisions" and
  "Write", because the move must be part of the same change that creates the
  Spec folder.
- `.agents/skills/write-prd/references/prd-template.md` — the template gains
  nothing structural; adopted sources are referenced from the body as they are
  today, now at their post-move paths.
- `.agents/skills/archive-spec/SKILL.md` + embedded pair — a third
  precondition, its commands, and the failure wording.
- `.agents/skills/write-idea/SKILL.md`, `write-techspec`, `write-tasks` + pairs
  — one referencing rule each.
- `docs/agents/docs-layout.md` — `docs/specs/<slug>/references/` enters the
  layout; the inbox and findings entries state that an adopted document leaves
  them.
- `docs/agents/spec-routing.md` — routing states that committing to
  implementation transfers ownership.
- `CONTEXT.md` — the `Spec` glossary entry's artifact set gains the references
  directory and its index.

`make skills-sync` propagates each `.agents/skills/` edit to its embedded
counterpart, and `make baseline-digests` owns the setup-snapshot and
parity-fixture fallout. No file outside the authorized list is edited by hand.

## Implementation Design

### The adoption transition (`write-prd`, new step)

Runs after scope is accepted and before the Spec folder is written, so the move
and the folder arrive in one change:

1. **Inventory.** List every inbox note and finding the PRD relies on. A
   document merely cited for background is not adopted; a document whose content
   the PRD depends on is.
2. **Classify.** A raw inbox note that is source material moves directly. A note
   that records observed behavior becomes a finding first, then moves as a
   finding. Field evidence never enters a Spec without passing through the
   findings tree, because the findings tree is where "didn't we already see
   this?" is answered.
3. **Claim ownership.** For each adopted source, check whether another active or
   archived Spec already owns it. The first Spec that commits to implementation
   is the primary owner; a later Spec links the owner's copy and adopts nothing.
4. **Flip then move.** A finding records `status: done` and its owning-Spec link
   **before** the move, in the same change, so the status trail is legible at the
   old path in Git history.
5. **Move.** `git mv <source> docs/specs/<slug>/references/<basename>` — one
   move, never a copy, never a stub. The basename and bytes are preserved; the
   observations inside are not rewritten.
6. **Index.** Append a row to
   `docs/specs/<slug>/references/_index.md`.
7. **Rewrite links.** Repository-wide, repoint every link to the old path at the
   new one, relative to the linking file.
8. **Gate.** Fail, and do not report completion, while any adopted source still
   exists at its original path or any link to it is unresolved.

### The reference index

`docs/specs/<slug>/references/_index.md`:

```markdown
# Adopted sources

| source                                    | type    | owner | adopted    | path                    |
| ----------------------------------------- | ------- | ----- | ---------- | ----------------------- |
| docs/findings/2026-07-25-spec-owned-…md   | finding | 0060  | 2026-07-30 | 2026-07-25-spec-owned-…md |
| docs/_inbox/2026-07-26-triage-notes.md    | inbox   | 0060  | 2026-07-30 | 2026-07-26-triage-notes.md |
```

`source` is the pre-adoption path and is never updated — it is the provenance
record. `path` is relative to the index, so it survives the archive move
unchanged. `owner` is the four-digit Spec number, so a secondary Spec's index
row names the primary owner rather than itself.

### The archive precondition (`archive-spec`, third check)

Added beside "every task completed" and "QA passed", with the same
verify-don't-trust framing and named commands:

```bash
# 1. no reference resolves outside the Spec
grep -rn 'docs/_inbox/\|docs/findings/' docs/specs/<slug>/

# 2. every indexed source exists at its recorded relative path
# 3. no adopted source remains at its pre-adoption path
```

Check 1 must not fire on a legitimate mention: a link into those trees fails the
gate, while prose naming a tree does not. The instruction states the
distinction and the gate reads link syntax, not bare words.

The failure names the offending source or link and the transition that fixes it
— which of the eight adoption steps was skipped — rather than only reporting
that the archive is blocked.

`qa_override: true` does not override this check. The QA override is the
maintainer's word about verification; self-containment is a property of the
artifact and is repaired by finishing the move, not by a decision.

### Migration boundary

New promotions only. An existing archived Spec is not retrofitted, and a finding
that already routed to a shipped Spec keeps its `done` status and links exactly
as they are. The archive gate therefore applies only to Specs whose PRD carries
a `references/_index.md`; a Spec without one is a pre-contract Spec and passes
the check trivially. This is what keeps the gate from turning every historical
archive red.

## Coverage Map

- Goal 1 / Story 1 → adoption steps 1–6 in `write-prd` (Features 1, 7).
- Goal 1 / Story 4 → adoption step 3, ownership claim (Feature 5).
- Goal 2 / Stories 2–3 → the index and the archive precondition
  (Features 2, 4).
- Goal 3 / Story 5 → adoption steps 4 and 7 plus the authoring gate
  (Feature 3).
- Feature 6 → the migration boundary.

## Integration Points

- Git — `git mv` for history preservation and `git log --follow` for the trail
  at the old path. No other tooling.

## Testing Approach

There is no Go seam, so verification is the Skill contract tests plus one
end-to-end rehearsal:

- `roundfix skills check` and `make skills-sync-check` prove the edited pairs
  stay in sync and the embedded bundle matches canonical.
- The existing Skill documentation contract tests in `./skills` gate the
  authorized pairs; any test pinning `write-prd` or `archive-spec` step
  structure updates in the same change as the wording.
- **Rehearsal, on a scratch branch, discarded afterward**: create a throwaway
  finding and inbox note, run the adoption transition by hand exactly as written,
  and confirm the move preserves history (`git log --follow`), the index
  populates, no link points at the old path
  (`grep -rn <old-basename> --include='*.md' .` returns only the index's
  provenance column and Git history), and the archive gate passes. Then inject a
  stale link and confirm the gate fails naming that link. This is the only way to
  prove an instruction-level gate actually catches its case, and its result is
  recorded in the Task's `## Result` rather than committed.
- `make verify` — exit 0.

## Build Order

1. `write-prd` adoption transition and the index format, including the authoring
   gate (no dependency).
2. `archive-spec` self-containment precondition, its commands, its failure
   wording, and the migration boundary (depends on: 1 for the index the check
   reads).
3. Referencing convention in `write-idea`, `write-techspec`, `write-tasks`
   (no dependency on 1–2, but lands after so all five pairs move once).
4. `docs/agents/docs-layout.md`, `docs/agents/spec-routing.md`, the `Spec`
   glossary entry in `CONTEXT.md`, the rehearsal, `make skills-sync`, and the
   authorized derived digest fallout (depends on: 1–3).

## Risks & Considerations

- **A gate written as prose can be skipped silently.** This is the Spec's main
  weakness and it cannot be fully closed here. The mitigation is that each check
  is a named command whose output goes in the report, so a skipped check leaves a
  hole a reader can see. If skipping recurs, the follow-up is a Go gate reading
  the index at archive time — deliberately out of scope until there is evidence
  it is needed.
- **The link rewrite is the step most likely to be done partially.** Fixing one
  of three occurrences is this repository's most repeated defect. The gate's
  grep is repository-wide for exactly that reason, and the instruction says to
  count occurrences before and after.
- **`git mv` of a findings file breaks any external bookmark.** Accepted: the
  provenance column records where it came from, and Git history resolves the old
  path. Nothing outside the repository is promised a stable path.
- **The archive gate could fire on prose.** A guide that discusses
  `docs/findings/` as a concept must not fail the check. The gate reads link
  syntax; the rehearsal includes a prose mention to prove it.
- **Shared sources invite two owners.** Step 3 makes ownership a check against
  both the active and archived Spec trees, not an assumption; two Specs authored
  in parallel are the case to watch, and the first to commit to implementation
  wins.

## Decisions

- `write-prd` owns the move, `archive-spec` owns the validation — per the PRD.
- The index is a Markdown table with a never-updated provenance column and a
  relative current path, so it survives archiving without edits and can be
  promoted to a parsed format when a consumer appears.
- `qa_override` does not override self-containment: the override speaks to
  verification, not to whether the artifact is whole.
- The gate applies only to Specs carrying an index, which is what keeps the
  migration boundary from reddening every historical archive.
