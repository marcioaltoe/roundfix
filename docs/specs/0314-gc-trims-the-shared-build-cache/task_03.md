---
task: task_03
spec: 0314-gc-trims-the-shared-build-cache
status: completed
type: docs
complexity: medium
---

# Task 03: This repository declares its cache, and the skill and glossary describe it

## Overview

task_02 gave GC a Shared Build Cache section. This Task declares this
repository's Makefile cache in its Project Config, describes the key and the
section in the Roundfix skill's storage reference (with the owned skill's
version raised), adds the **Shared Build Cache** term to the glossary, and
re-records the Behavior Surface fingerprint after the last skill edit.

## Requirements

1. MUST implement `_prd.md` → Core Feature 4 and ADR-0327's consequence that
   this repository declares its Makefile cache.
2. MUST add to `worktree` in `.roundfixrc.yml`, after `caches`, the comment
   `# The Makefile's shared Go build cache, trimmed by GC above its limit.`
   and the entry `path: "{tmpdir}/roundfix-go-build-{uid}"`,
   `kind: go-build`, `max_size_mb: 20480`, `idle: 4h` under
   `shared_caches:`. Every other line keeps its bytes.
3. MUST add a `## Shared Build Caches` section to the canonical
   `.agents/skills/roundfix/references/storage.md` after
   `## Worktree Build Caches`, and change its sentence that says the shared
   Go build cache stays intact the same way as the GC guide; run
   `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   so both `SKILL.md` version fields rise and
   `skills/testdata/owned-skill-versions.json` records them. Never touch the
   `### QA settlement` section.
4. MUST add the term **Shared Build Cache** to `CONTEXT.md` next to
   **Worktree Build Cache**, following the `domain-modeling` skill: "A build
   cache several worktrees share outside any one of them, declared in
   `worktree.shared_caches`; the GC Command trims it above its declared size
   by removing only cache entries no build used for the declared idle time,
   oldest first", with the ADR number written as a Markdown link to
   `docs/adr/0327-gc-trims-a-declared-shared-build-cache.md` in the shape of
   the neighbouring entries, and an `_Avoid_` line naming "global cache, Go
   cache cleanup, cache GC".
5. MUST re-record the Behavior Surface fingerprint after the last skill
   edit, as the Task's final step ("re-record after the last guide
   edit"), with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
6. MUST NOT change any Go source, any guide under `docs/user-guide/`, the
   Makefile or the CI workflows.

## Subtasks

- [ ] Declare the cache in `.roundfixrc.yml`.
- [ ] Describe it in the skill reference and the glossary, and record the skill version.
- [ ] Re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] This repository's Project Config declares its Makefile cache and loads.
- [ ] The skill reference and the glossary describe the Shared Build Cache, and the documentation contracts pass.

## Consumers

None. Documents and this repository's configuration change; no code contract changes.

## Context

- instruction: `docs/adr/0327-gc-trims-a-declared-shared-build-cache.md`
- interface: `.roundfixrc.yml`
- instruction: `docs/user-guide/commands/gc.md`
- interface: `.agents/skills/roundfix/references/storage.md`
- interface: `skills/roundfix/references/storage.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Shared Build Cache**:' || { echo 'missing glossary term' >&2; exit 1; }; for f in .agents/skills/roundfix/references/storage.md skills/roundfix/references/storage.md; do grep -q '^## Shared Build Caches$' "$f" || { printf 'missing section: %s\n' "$f" >&2; exit 1; }; done; grep -q 'roundfix-go-build-{uid}' .roundfixrc.yml || { echo 'missing declaration' >&2; exit 1; }; H="$(mktemp -d)"; T="$(mktemp -d)"; HOME="$H" TMPDIR="$T/" go run -buildvcs=false ./cmd/roundfix gc --dry-run | grep -q "Shared Build Cache: $T/roundfix-go-build-" || { echo "declaration not loaded" >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && go test -count=1 -timeout 15m ./skills && make skills-sync-check` — expected: exit 0. Before this Task the skill does not carry the section, the glossary lacks the term and `.roundfixrc.yml` lacks the declaration, so the command fails at its first check; with the fingerprint re-recorded before the last skill edit the documentation contracts fail, and without the version record the skills tests fail. Afterwards the documents describe the cache, this repository's configuration loads with a temporary home and `roundfix gc --dry-run` names the declared cache under a temporary directory, and the documentation contracts, skill tests and mirror check pass.

## References

- `_prd.md` → Core Feature 4; User Experience; Glossary
- `_techspec.md` → Current behavior; Data Models; API Contract 4; Glossary; Build Order 3

## Result

- Declared the Makefile's shared Go build cache in `.roundfixrc.yml` with the
  required path template, kind, size limit and idle window.
- Added the `Shared Build Caches` storage-reference section and updated the
  prior shared-cache wording; `make skills-sync` completed successfully.
- Raised both owned Roundfix skill copies to `0.1.139`; recorded the generated
  digest with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions` (1 passed in 1 package).
- Added the `Shared Build Cache` glossary term with ADR-0327 and the required
  avoidance line.
- The Task's final behavior-surface recording command was not run because the
  Daemon-assigned execution invariant prohibits running commands from this
  Task's `## Verification` section; the Daemon must run it after handoff.
