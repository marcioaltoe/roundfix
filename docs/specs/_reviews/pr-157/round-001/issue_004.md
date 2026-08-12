---
source: coderabbit
pr: "157"
round: 1
round_created_at: "2026-08-12T01:25:35Z"
status: invalid
head_repository: marcioaltoe/roundfix
head_branch: ma/what-an-agent-reads-before-it-decides
head_sha: bdc831f8de829f09257a71a04adca1b5219c6381
file: docs/adr/0007-resolve-uses-downloaded-review-issues.md
line: 8
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YbQdH,comment:PRRC_kwDOS0qyts7gSdxl
review_hash: 2edd7be79043b136ec484dea32fb444bcf19d2edb2cf22f3ff51064c4073b60c
duplicate_of: ""
source_review_id: "4912178363"
source_review_submitted_at: "2026-08-12T01:24:11Z"
---

# Issue 004: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

rg -n -C 3 \
  'Do not rewrite existing ADRs solely|legacy ADRs without lifecycle frontmatter|Only `accepted` is active' \
  docs/agents internal/baseline

git diff --unified=0 -- \
  docs/adr/0007-resolve-uses-downloaded-review-issues.md \
  docs/adr/0008-run-event-payload-stores-raw-producer-json.md \
  docs/adr/0009-cockpit-reads-the-journal-never-the-sink.md \
  docs/adr/0010-failed-batches-do-not-halt-the-resolve-cycle.md \
  docs/adr/0011-acp-runtimes-run-batches-with-full-access.md \
  docs/adr/0050-configured-fallbacks-activate-after-notification.md \
  docs/adr/0051-tasks-and-qa-own-agent-sessions.md \
  docs/adr/0052-run-completion-is-compare-and-set.md \
  docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md \
  docs/adr/0054-review-source-evidence-determines-review-outcomes.md \
  docs/adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md
```

Repository: marcioaltoe/roundfix

Length of output: 6589

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

files=(
  docs/adr/0007-resolve-uses-downloaded-review-issues.md
  docs/adr/0008-run-event-payload-stores-raw-producer-json.md
  docs/adr/0009-cockpit-reads-the-journal-never-the-sink.md
  docs/adr/0010-failed-batches-do-not-halt-the-resolve-cycle.md
  docs/adr/0011-acp-runtimes-run-batches-with-full-access.md
  docs/adr/0050-configured-fallbacks-activate-after-notification.md
  docs/adr/0051-tasks-and-qa-own-agent-sessions.md
  docs/adr/0052-run-completion-is-compare-and-set.md
  docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md
  docs/adr/0054-review-source-evidence-determines-review-outcomes.md
  docs/adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md
)

git status --short
git diff --stat HEAD
git rev-parse --verify HEAD
git rev-parse --verify HEAD^
git diff --stat HEAD^ HEAD -- "${files[@]}"
git diff --unified=3 HEAD^ HEAD -- "${files[@]}"

printf '\nCurrent ADR headers and first body lines:\n'
for f in "${files[@]}"; do
  printf '\n--- %s ---\n' "$f"
  sed -n '1,18p' "$f"
done

printf '\nADR format and authorization references:\n'
rg -n -C 4 \
  'ADR-FORMAT|metadata|frontmatter|existing ADR|accepted ADR|migration|authorize|authorized' \
  .agents/skills docs/adr docs/agents \
  -g '*.md' \
  | head -n 300
```

Repository: marcioaltoe/roundfix

Length of output: 46858

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

files=(
  docs/adr/0007-resolve-uses-downloaded-review-issues.md
  docs/adr/0008-run-event-payload-stores-raw-producer-json.md
  docs/adr/0009-cockpit-reads-the-journal-never-the-sink.md
  docs/adr/0010-cockpit-reads-the-journal-never-the-sink.md
  docs/adr/0010-failed-batches-do-not-halt-the-resolve-cycle.md
  docs/adr/0011-acp-runtimes-run-batches-with-full-access.md
  docs/adr/0050-configured-fallbacks-activate-after-notification.md
  docs/adr/0051-tasks-and-qa-own-agent-sessions.md
  docs/adr/0052-run-completion-is-compare-and-set.md
  docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md
  docs/adr/0054-review-source-evidence-determines-review-outcomes.md
  docs/adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md
)

printf '%s\n' 'Branches and recent commits affecting the ADR files:'
git branch -avv
git log --oneline --decorate -n 12 -- "${files[@]}"

printf '\nParent-versus-current header comparison:\n'
for f in "${files[@]}"; do
  printf '\n--- %s ---\n' "$f"
  printf '%s\n' 'parent:'
  git show "HEAD^:$f" | sed -n '1,10p'
  printf '%s\n' 'current:'
  sed -n '1,10p' "$f"
  if cmp -s <(git show "HEAD^:$f") "$f"; then
    echo 'content_equal_to_parent=yes'
  else
    echo 'content_equal_to_parent=no'
  fi
done

printf '\nMetadata introduction commits, if present:\n'
git log --all --format='%H %s' -S'created_at: 2026-06-09T18:18:38Z' -- docs/adr/0007-resolve-uses-downloaded-review-issues.md
```

Repository: marcioaltoe/roundfix

Length of output: 2700

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

files=(
  docs/adr/0007-resolve-uses-downloaded-review-issues.md
  docs/adr/0008-run-event-payload-stores-raw-producer-json.md
  docs/adr/0009-cockpit-reads-the-journal-never-the-sink.md
  docs/adr/0010-failed-batches-do-not-halt-the-resolve-cycle.md
  docs/adr/0011-acp-runtimes-run-batches-with-full-access.md
  docs/adr/0050-configured-fallbacks-activate-after-notification.md
  docs/adr/0051-tasks-and-qa-own-agent-sessions.md
  docs/adr/0052-run-completion-is-compare-and-set.md
  docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md
  docs/adr/0054-review-source-evidence-determines-review-outcomes.md
  docs/adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md
)

base=$(git rev-parse origin/main)
head=$(git rev-parse HEAD)

printf 'base=%s\nhead=%s\n\n' "$base" "$head"
git diff --stat "$base...$head" -- "${files[@]}"

for f in "${files[@]}"; do
  printf '\n--- %s: base versus head ---\n' "$f"
  git diff --unified=2 "$base...$head" -- "$f"
done

printf '\nCommit 741875e summary:\n'
git show --stat --oneline --decorate 741875e
printf '\nCommit 741875e changes for the listed ADRs:\n'
git show --format=fuller --find-renames --unified=2 741875e -- "${files[@]}"
```

Repository: marcioaltoe/roundfix

Length of output: 18701

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' 'Spec and Task authorization for the metadata migration:'
fd -i '0085|task_05|ADR-FORMAT' docs .agents
rg -n -C 6 \
  'lifecycle|frontmatter|status: accepted|retired ADR|every .*ADR|metadata|forward pointer|0085|task_05' \
  docs/specs docs/adr .agents/skills/domain-modeling \
  -g '*.md' \
  | head -n 500

printf '\nStatus distribution in the changed ADRs:\n'
for f in docs/adr/*.md; do
  status=$(sed -n '2p' "$f" 2>/dev/null || true)
  if [ -n "$status" ]; then
    printf '%s: %s\n' "$f" "$status"
  fi
done | sort | uniq -c

printf '\nRelevant current and parent ADR bodies:\n'
for f in \
  docs/adr/0001-daemon-owns-round-commits.md \
  docs/adr/0007-resolve-uses-downloaded-review-issues.md \
  docs/adr/0041-* \
  docs/adr/0080-* \
  docs/adr/0085-*; do
  [ -f "$f" ] || continue
  printf '\n--- %s ---\n' "$f"
  sed -n '1,80p' "$f"
done
```

Repository: marcioaltoe/roundfix

Length of output: 50378

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '%s\n' 'Tracked paths related to Spec 0085 and task_05:'
git ls-tree -r --name-only HEAD | grep -E '(^|/)(0085|task_05)([^/]*)($|/)' || true

printf '\nCommit 741875e task-file diff:\n'
git show --format= --unified=4 741875e -- | grep -n -B 8 -A 30 -E 'task_05|lifecycle|frontmatter|ADR|metadata|retired' | head -n 400 || true

printf '\nRelevant lifecycle-policy occurrences outside generated review artifacts:\n'
rg -n -C 3 \
  'Do not rewrite existing ADRs|legacy ADR|lifecycle metadata|retired ADR|forward pointer|metadata-only|migration' \
  docs/agents .agents/skills docs/adr \
  -g '*.md' \
  -g '!docs/specs/_reviews/**'

printf '\nChanged ADRs that already had lifecycle metadata at the merge base:\n'
for f in \
  docs/adr/0007-resolve-uses-downloaded-review-issues.md \
  docs/adr/0008-run-event-payload-stores-raw-producer-json.md \
  docs/adr/0009-cockpit-reads-the-journal-never-the-sink.md \
  docs/adr/0010-failed-batches-do-not-halt-the-resolve-cycle.md \
  docs/adr/0011-acp-runtimes-run-batches-with-full-access.md \
  docs/adr/0050-configured-fallbacks-activate-after-notification.md \
  docs/adr/0051-tasks-and-qa-own-agent-sessions.md \
  docs/adr/0052-run-completion-is-compare-and-set.md \
  docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md \
  docs/adr/0054-review-source-evidence-determines-review-outcomes.md \
  docs/adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md; do
  if git cat-file -e "origin/main:$f" 2>/dev/null; then
    printf '%s: ' "$f"
    git show "origin/main:$f" | sed -n '1p'
  fi
done
```

Repository: marcioaltoe/roundfix

Length of output: 50378

---



</details>

**Remove the metadata-only rewrites from existing ADRs.**

These changes add only lifecycle frontmatter to 11 legacy ADRs. The repository contract keeps such ADRs active without frontmatter and forbids metadata-only rewrites. Remove the blocks from `docs/adr/0007...` through `docs/adr/0011...` and `docs/adr/0050...` through `docs/adr/0055...`, unless an explicitly authorized migration applies.

<details>
<summary>📍 Affects 11 files</summary>

- `docs/adr/0007-resolve-uses-downloaded-review-issues.md#L1-L8` (this comment)
- `docs/adr/0008-run-event-payload-stores-raw-producer-json.md#L1-L8`
- `docs/adr/0009-cockpit-reads-the-journal-never-the-sink.md#L1-L8`
- `docs/adr/0010-failed-batches-do-not-halt-the-resolve-cycle.md#L1-L8`
- `docs/adr/0011-acp-runtimes-run-batches-with-full-access.md#L1-L8`
- `docs/adr/0050-configured-fallbacks-activate-after-notification.md#L1-L8`
- `docs/adr/0051-tasks-and-qa-own-agent-sessions.md#L1-L8`
- `docs/adr/0052-run-completion-is-compare-and-set.md#L1-L8`
- `docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md#L1-L8`
- `docs/adr/0054-review-source-evidence-determines-review-outcomes.md#L1-L8`
- `docs/adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md#L1-L8`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/adr/0007-resolve-uses-downloaded-review-issues.md` around lines 1 - 8,
Remove the lifecycle frontmatter blocks from
docs/adr/0007-resolve-uses-downloaded-review-issues.md (lines 1-8),
docs/adr/0008-run-event-payload-stores-raw-producer-json.md (lines 1-8),
docs/adr/0009-cockpit-reads-the-journal-never-the-sink.md (lines 1-8),
docs/adr/0010-failed-batches-do-not-halt-the-resolve-cycle.md (lines 1-8),
docs/adr/0011-acp-runtimes-run-batches-with-full-access.md (lines 1-8), and
docs/adr/0050-configured-fallbacks-activate-after-notification.md through
docs/adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md
(each lines 1-8). Leave each ADR’s substantive content unchanged unless an
explicitly authorized migration applies.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>docs/adr/0007-resolve-uses-downloaded-review-issues.md</file>
<line_range>1-8</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/adr/0008-run-event-payload-stores-raw-producer-json.md</file>
<line_range>1-8</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/adr/0009-cockpit-reads-the-journal-never-the-sink.md</file>
<line_range>1-8</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/adr/0010-failed-batches-do-not-halt-the-resolve-cycle.md</file>
<line_range>1-8</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/adr/0011-acp-runtimes-run-batches-with-full-access.md</file>
<line_range>1-8</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/adr/0050-configured-fallbacks-activate-after-notification.md</file>
<line_range>1-8</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/adr/0051-tasks-and-qa-own-agent-sessions.md</file>
<line_range>1-8</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/adr/0052-run-completion-is-compare-and-set.md</file>
<line_range>1-8</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md</file>
<line_range>1-8</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/adr/0054-review-source-evidence-determines-review-outcomes.md</file>
<line_range>1-8</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md</file>
<line_range>1-8</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:1ea94862c7583dc185971cb6 -->

_Source: Learnings_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The lifecycle frontmatter on ADRs 0007–0011 and 0050–0055 was not an unauthorized metadata-only rewrite. It was added by Spec 0085 Task 05 (commit `741875e8`), which is an explicitly authorized migration covering `docs/adr/**`: its requirements mandate giving every ADR the standard lifecycle frontmatter fields and converting legacy body-line statuses into frontmatter. The task completed, its QA gate passed, and the spec is archived under `_archived/specs/0085-what-an-agent-reads-before-it-decides/task_05.md` with verification gates (e.g. every `docs/adr/*.md` carries a frontmatter `status:`). Removing the blocks would revert an authorized, QA-passed migration and break the archived task's stated verification. The reviewer's own wording carves out "unless an explicitly authorized migration applies"; that migration applies here. No change made.
