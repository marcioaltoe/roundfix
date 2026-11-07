---
task: task_01
spec: 0364-a-rerun-finishes-a-partial-npm-publication
status: completed
type: docs
complexity: low
---

# Task 01: The runbook, the glossary and ADR-0146 say a rerun finishes a partial publication

## Overview

Tasks 02 and 03 make the release workflow's Publication Preflight accept a
coordinate already published at the target version and let a tag run publish
only what is missing, once the published content equals its own build. This
Task writes that rule first: the amendment notice in ADR-0146, the
**Release Set** and **Publication Preflight** entries of `CONTEXT.md`, the
release runbook's recovery procedure and refusals, and the two Coverage
Reviews of the runbook's surfaces; it re-records the Behavior Surface Record
last.

## Requirements

1. MUST, in `docs/adr/0146-release-publication-is-all-or-nothing-with-a-bounded-token-fallback.md`,
   set `updated_at: 2026-11-06T23:30:00Z` and add, as a new paragraph right
   after the paragraph that ends "the identity gap cannot produce exactly the
   partial release the invariant forbids.", this paragraph, exactly as the
   block below holds it (one line in the file):

   ```markdown
   **Amended by [ADR-0374](0374-a-rerun-completes-a-release-set-from-the-same-content.md) (2026-11-06).** A coordinate already published at the target version no longer stops publication when its published tarball holds the same content as the tag run's build: the preflight reports it `published`, the tag run compares the decompressed tar streams after the build, and the publish step skips it, so a rerun completes a Release Set that stopped part way. A coordinate published with different content, or whose version was unpublished, still stops publication before any package publishes, and the bounded token fallback is unchanged.
   ```

   No other byte of ADR-0146 changes.
2. MUST, in `CONTEXT.md`, replace the definition line of **Release Set**
   and the definition line of **Publication Preflight** with these two lines,
   in that order of entries, exactly as the block below holds them:

   ```markdown
   The launcher package and five platform package coordinates that form one indivisible Roundfix publication unit. Every coordinate must be eligible, or already published with the same content as the tag run's build, before any one is published ([ADR-0374](docs/adr/0374-a-rerun-completes-a-release-set-from-the-same-content.md)).
   A read-only eligibility check of the entire Release Set against registry truth for an exact target version. It stops publication unless every coordinate is proven eligible or already published at that version, and it stops when all of them are already published; a tag run then compares each already published coordinate with its own build before it publishes the rest.
   ```

   The first line replaces the **Release Set** definition, the second the
   **Publication Preflight** definition; the term lines and both `_Avoid_`
   lines stay as they are.
3. MUST, in `docs/user-guide/release-runbook.md`:
   - add, as a new paragraph right after the paragraph of "Rehearsing the
     Publication Preflight" that ends "cross-compilation, npm publication, or
     GitHub Release creation.", the paragraph: "A coordinate already published at the version shows `published` and does not stop the rehearsal; only a tag run compares its content with a build.";
   - add, right before the heading `## Release workflow failure vocabulary`,
     a section `## Finishing a partial publication` with these two
     paragraphs:
     "When `npm publish` fails after some coordinates of the Release Set have published, finish the release by rerunning the failed job of the same workflow run ("Re-run failed jobs"); do not cut a new version for it. The Publication Preflight reports each coordinate already published at the target version as `published`. After the build, the step "Compare published packages" packs this run's build of each such coordinate, fetches the published tarball with `npm pack <coordinate>@<version>` and compares the two decompressed tar streams. When they are byte-identical it prints `same content`, and the publish step prints `skipping <coordinate>@<version>: already published with the same content` and publishes only the missing coordinates, platform packages first and the launcher last."
     and
     "The step compares the decompressed tar streams because the tarball integrity changes with the gzip output of the Node build that packs it, even when every packed byte is the same. A rerun rebuilds the binaries with the Go release that `actions/setup-go` resolves at that moment: if a newer Go patch release appeared after the first attempt, the rebuild differs and the step stops with `different content`. The version then cannot be completed, and the next patch version republishes the same changes, as it does for a version that was removed with `npm unpublish`.";
   - replace the two vocabulary rows that begin "| `registry:` |
     Publication Preflight |" and "| `undetermined:` | Publication
     Preflight |" with these three rows:
     "| `registry:` | Publication Preflight | The named coordinate cannot accept the target version because the version was published and then unpublished (`already used`), remains in the printed post-unpublish cooldown, or is absent; or every coordinate of the Release Set is already published at the version. | Follow the named reason: choose an unused version, wait until the printed cooldown time, or restore the absent coordinate; then rerun the rehearsal. A Release Set already published at the version has nothing left to publish. |",
     "| `registry:` | Compare published packages | The named coordinate is already published at the target version with content that differs from this run's build (`different content`); the line prints both tarball integrities. This run has published nothing. | Release the next patch version; a published version can never be replaced. |",
     "| `undetermined:` | Publication Preflight, Compare published packages | Registry eligibility could not be proven because the Release Set, transport, HTTP response, registry body, or cooldown timestamp could not be read, or because this run's build or the published tarball of a coordinate could not be packed, fetched or decompressed. This run has published nothing. | Inspect the named read or classification failure and rerun after registry access is reliable. |".
4. MUST add to `docs/references/skill-coverage.json`, on the surfaces
   `guide docs/user-guide/release-runbook.md` and `command release plan`,
   after their `sources` list, the Coverage Review:
   `"review": "2026-11-06 — Spec 0364 task_01: the runbook gains the procedure that finishes a partial npm publication and the refusals of the release workflow's Compare published packages step; the release reference describes the release plan and the tag-triggered workflow without naming its publication steps, so it needs no text change."`
   No skill changes.
5. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's LAST step, after the last edit of the runbook. No record is
   hand-edited.
6. MUST NOT change any other paragraph of these files, any other ADR (ADR-0374
   is already in the tree), any skill, any Go file or
   `.github/workflows/release.yml`.

## Subtasks

- [ ] Add the amendment notice to ADR-0146 and bump its `updated_at`.
- [ ] Update the two glossary entries.
- [ ] Add the rehearsal paragraph, the recovery section and the vocabulary rows to the runbook.
- [ ] Add the two Coverage Reviews, then re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] ADR-0146 names ADR-0374 as its amendment and keeps the fallback rule.
- [ ] The glossary and the runbook say a coordinate published with the same content counts for the Release Set and how a rerun finishes the release.
- [ ] The documentation contracts pass, including the glossary and coverage checks.

## Context

- instruction: `docs/adr/0374-a-rerun-completes-a-release-set-from-the-same-content.md`
- interface: `docs/adr/0146-release-publication-is-all-or-nothing-with-a-bounded-token-fallback.md`
- interface: `CONTEXT.md`
- interface: `docs/user-guide/release-runbook.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The glossary, the runbook and the coverage records — consumers: `./internal/docscontract`

## Verification

- `tr -s '[:space:]' ' ' < docs/adr/0146-release-publication-is-all-or-nothing-with-a-bounded-token-fallback.md | grep -qF -- 'Amended by [ADR-0374](0374-a-rerun-completes-a-release-set-from-the-same-content.md) (2026-11-06).' || exit 1; tr -s '[:space:]' ' ' < docs/adr/0146-release-publication-is-all-or-nothing-with-a-bounded-token-fallback.md | grep -qF -- 'updated_at: 2026-11-06T23:30:00Z' || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'or already published with the same content as the tag run' || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'compares each already published coordinate with its own build' || exit 1; grep -qF -- '## Finishing a partial publication' docs/user-guide/release-runbook.md || exit 1; tr -s '[:space:]' ' ' < docs/user-guide/release-runbook.md | grep -qF -- 'skipping <coordinate>@<version>: already published with the same content' || exit 1; grep -qF -- 'Spec 0364 task_01' docs/references/skill-coverage.json || exit 1; grep -F -- '**Release Set**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Release Set\*\*(:| —( |$))' || exit 1; grep -F -- '**Publication Preflight**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Publication Preflight\*\*(:| —( |$))' || exit 1; go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract -run '^(TestTheSkillCoverageMapIsCurrent|TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryStatesEachRuleOnce)$'` — expected: exit 0. Before this Task none of the files carries the new text, so the first check fails; afterwards the text is in place and the glossary, ADR-link and coverage contracts pass.

## References

- `_prd.md` → Goals; User Story 3; Core Feature 4
- `_techspec.md` → System Architecture; Build Order 1
- ADR-0374

## Result

Skills loaded: tech-writer

- ADR-0146 now carries `updated_at: 2026-11-06T23:30:00Z` and, right after the paragraph ending "the identity gap cannot produce exactly the partial release the invariant forbids.", the one-line amendment notice naming ADR-0374 and keeping the bounded token fallback unchanged; no other byte of the ADR changed.
- `CONTEXT.md` defines **Release Set** as eligible or already published with the same content as the tag run's build (citing ADR-0374), and **Publication Preflight** as stopping unless every coordinate is eligible or already published, still stopping when all are published, and leaving the content comparison to a tag run. Term lines, `_Avoid_` lines and the other entries are untouched.
- The runbook's "Rehearsing the Publication Preflight" gains the `published`/no-stop sentence; a new "Finishing a partial publication" section explains "Re-run failed jobs", the `published` classification, the "Compare published packages" step (pack, `npm pack <coordinate>@<version>`, decompressed tar-stream comparison, `same content`, `skipping …`, platform packages first and launcher last) and the `different content` / Go patch-release case. The vocabulary table's two old rows are replaced by the three required rows (Preflight `registry:`, Compare published packages `registry:`, and the widened `undetermined:`).
- Both Coverage Reviews (`command release plan`, `guide docs/user-guide/release-runbook.md`) carry the 2026-11-06 Spec 0364 task_01 review line; no skill changed.
- Behavior Surface Record re-recorded last, after the final runbook edit, with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` (ok, 0.477s); the diff touches only the runbook surface hash in `behavior-surfaces.json`.

### Focused checks (this session)

- `go test -count=1 -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs)$'` → ok (0.327s).
- Whitespace-normalized `grep -cF` probes for the ADR amendment sentence, the `COMPARES each already published coordinate` glossary phrase, the `## Finishing a partial publication` heading and the `skipping <coordinate>@<version>: already published with the same content` runbook text all return 1.
- `python3 -c "json.load(...)"` on `skill-coverage.json` parses; `git status --short` shows only the five Spec-listed files plus this task file.
- The declared `## Verification` commands were not run (Daemon-owned).

## Carry-forward provenance

- Source Run: `run_20261107T002718Z_f5b7975ba1530c9f`
- Source commit: `90c3672dfa27de1a1b060a6f5da2669eed5b8e0e`
