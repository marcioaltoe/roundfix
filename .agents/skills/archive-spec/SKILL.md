---
name: archive-spec
description: Archive a completed spec — verify every task completed, QA passed, and indexed references are self-contained, then stamp the archive metadata and move docs/specs/<slug>/ to docs/specs/_archived/<slug>/. Runs automatically at the end of the implement-spec loop after a QA pass, or whenever the user asks to archive a spec.
argument-hint: "<spec slug> [--release <tag or PR URL>]"
metadata:
  category: delivery
  tags: [workflow, documentation, process]
  version: 0.0.2
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Archive Spec

Move a completed spec out of the active set: `docs/specs/<slug>/` → `docs/specs/_archived/<slug>/`, with the completion stamped in its frontmatter. Archived means _implemented, verified, and self-contained_ — every task done, QA passed, and every indexed reference owned by the Spec — after this, one `ls docs/specs/` separates live work from history, and the archive stays greppable as the record of what was built and why.

The trigger is spec completion, not publication: run this automatically at the end of the `implement-spec` loop once the QA gate passes, or whenever the user asks. Merge and release are separate, user-driven steps — the archive commit simply travels with the branch and ships inside the feature's own PR.

## Preconditions — verify, don't trust

Check all three with fresh command evidence before touching anything:

1. **Every task completed.** Read each `task_NN.md` listed in `_tasks.md`; every
   `status` must be `completed`.

   **Command:** run `grep -n '^status:' <each-task-file-listed-in-_tasks.md>` and
   retain its output. Any value other than `status: completed` blocks the
   archive and names the Task file.

2. **QA passed.** The newest report in `qa/` must pass the repository's QA
   verifier. Do not substitute a line grep for structured validation. In
   Roundfix, use the same `internal/spec.QAVerdict` contract as the Archive
   Command: select the newest report by the `qa-report-YYYY-MM-DD[-NN].md`
   filename contract, parse its YAML frontmatter, require a supported
   `verdict`, require both blocked-row fields to be non-negative integers when
   present, and reject `verdict: pass` when `rows_blocked_finding` is nonzero.
   Retain the verifier's report path and result as evidence. A missing `qa/`
   directory, malformed newest report, or non-passing verdict blocks the
   archive; proceed only if the user explicitly says "archive anyway", and
   record that override in the stamped frontmatter (`qa_override: true`).

3. **The Spec is self-contained.** Apply this precondition when
   `docs/specs/<slug>/references/_index.md` exists or is a symbolic link; a
   symbolic link, including a broken one, is invalid index state rather than a
   legacy Spec. Only a Spec where that path neither exists nor is a symbolic
   link predates this contract and passes without retrofitting historical
   artifacts. For an indexed Spec, every indexed `path` must exist relative to
   `_index.md`, every never-updated `source` path must be absent, and no
   Markdown link destination inside the Spec may point into `docs/_inbox/` or
   `docs/findings/`.

   **Commands:** first run this link-destination check; its syntax deliberately
   matches inline or reference-style Markdown links, not prose that merely
   names either tree:

   ```bash
   spec_dir=docs/specs/<slug>
   link_hits=$(grep --include='*.md' -RInE '(\]\([^)]*(docs/)?(_inbox|findings)/[^)]*\)|^[[:space:]]*\[[^]]+\]:[[:space:]]*<?[^[:space:]>]*(docs/)?(_inbox|findings)/)' "$spec_dir")
   link_status=$?
   if test "$link_status" -eq 0; then
     printf '%s\n' "$link_hits"
     echo "self-containment failed: rewrite each listed link at adoption step 8"
     exit 1
   fi
   test "$link_status" -eq 1 || exit "$link_status"
   ```

   Then parse and validate every data row. The index belongs to the current
   Spec: `owner` must equal its four-digit prefix, `type` must be `inbox` or
   `finding`, and each `source` and `path` must appear only once. A `path` must
   be one basename relative to `_index.md`; reject absolute paths, `.`, `..`,
   path separators, and symbolic links instead of allowing traversal or a link
   outside `references/`. Run the following from the repository root and
   retain the normalized rows plus any diagnostic as evidence:

   ```bash
   slug=<slug>
   index="docs/specs/$slug/references/_index.md"
   expected_owner=${slug%%-*}
   parsed_index=$(mktemp) || exit 1
   trap 'rm -f "$parsed_index"' EXIT HUP INT TERM

   if test -L "$(dirname "$index")"; then
     printf 'self-containment failed: references/ must not be a symbolic link\n' >&2
     exit 1
   fi

   if test -L "$index" || test ! -f "$index"; then
     printf 'self-containment failed: references/_index.md must be a regular, non-symbolic-link file\n' >&2
     exit 1
   fi

   awk -F '|' -v expected_owner="$expected_owner" '
   function trim(value) {
     gsub(/^[[:space:]]+|[[:space:]]+$/, "", value)
     return value
   }
   function reject(message) {
     print "self-containment failed: " message > "/dev/stderr"
     invalid = 1
   }
   function separator(value) {
     value = trim(value)
     return value ~ /^:?-{3,}:?$/
   }
   /^[[:space:]]*$/ || /^[[:space:]]*#/ { next }
   !header {
     if (NF != 7 || trim($2) != "source" || trim($3) != "type" ||
         trim($4) != "owner" || trim($5) != "adopted date" ||
         trim($6) != "path") {
       reject("_index.md must use the fixed source | type | owner | adopted date | path header")
     } else {
       header = 1
     }
     next
   }
   !divider {
     if (NF != 7 || !separator($2) || !separator($3) || !separator($4) ||
         !separator($5) || !separator($6)) {
       reject("_index.md has an invalid table separator")
     } else {
       divider = 1
     }
     next
   }
   {
     source = trim($2)
     type = trim($3)
     owner = trim($4)
     adopted = trim($5)
     path = trim($6)
     if (NF != 7 || source == "" || type == "" || owner == "" ||
         adopted == "" || path == "") {
       reject("invalid index row at line " NR ": " $0)
       next
     }
     if (type != "inbox" && type != "finding") {
       reject("type must be `inbox` or `finding` at line " NR ": " type)
     }
     if (owner != expected_owner) {
       reject("owner must be " expected_owner " at line " NR ": " owner)
     }
     if (seen_source[source]++) {
       reject("duplicate source at line " NR ": " source)
     }
     if (seen_path[path]++) {
       reject("duplicate path at line " NR ": " path)
     }
     if (adopted !~ /^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]$/) {
       reject("adopted date must be YYYY-MM-DD at line " NR ": " adopted)
     }
     if ((type == "inbox" && source !~ /^docs\/_inbox\/[^\/]+\.md$/) ||
         (type == "finding" && source !~ /^docs\/findings\/[^\/]+\.md$/)) {
       reject("source does not match type at line " NR ": " source)
     }
     print source "|" type "|" owner "|" adopted "|" path
   }
   END {
     if (!header || !divider) {
       reject("_index.md is missing its fixed table header")
     }
     if (invalid) {
       exit 1
     }
   }
   ' "$index" > "$parsed_index" || exit $?

   while IFS='|' read -r source type owner adopted path; do
     case "$path" in
       ""|.|..|/*|*/*|*\\*)
         printf 'self-containment failed: path must be one basename relative to `_index.md`: %s\n' "$path" >&2
         exit 1
         ;;
     esac
     source_basename=${source##*/}
     if test "$path" != "$source_basename"; then
       printf 'self-containment failed: path must equal source basename for %s: %s != %s\n' "$source" "$path" "$source_basename" >&2
       exit 1
     fi
     current="$(dirname "$index")/$path"
     if test -L "$current" || test ! -f "$current"; then
       printf 'self-containment failed: invalid or missing path %s; finish adoption step 7\n' "$path" >&2
       exit 1
     fi
     if test -e "$source" || test -L "$source"; then
       printf 'self-containment failed: source still exists at %s; finish adoption step 7\n' "$source" >&2
       exit 1
     fi
   done < "$parsed_index"
   cat "$parsed_index"
   ```

`qa_override: true` overrides only failed or missing QA evidence in precondition
2. It never overrides self-containment: verification can be overridden by the
maintainer, but self-containment is a property of the artifact and must be
repaired by finishing adoption.

A merged PR or release tag is **not** a precondition. If the user passes `--release`, or a merged PR/tag is already known, stamp it as metadata — but never block the archive waiting for one.

If any check fails, stop and report the offending Task, report, source, or link
and the adoption step that fixes a self-containment failure — the Spec stays
active.

## Steps

1. **Stamp** `_prd.md` frontmatter:

   ```yaml
   status: archived
   archived: YYYY-MM-DD
   qa_override: true # only when archiving despite failed/missing QA
   release: <tag or PR URL> # only when known — from --release or an already-merged PR/tag
   ```

2. **Move** with history preserved:

   ```bash
   mkdir -p docs/specs/_archived
   git mv docs/specs/<slug> docs/specs/_archived/<slug>
   ```

3. **Commit** — `chore(specs): archive <slug>` (Conventional Commits). Do not push unless asked.

4. **Report** — the new path, the release reference when one was stamped, and anything carried over (open follow-ups from task `## Result` sections belong in new specs, not in the archive).

5. **Suggest the publish step** — when the work isn't merged yet, close by suggesting the PR (via `github-pr-workflow`). This is where that suggestion lives in the workflow — the implement loop ends at the archive and doesn't offer it. Suggest only: opening the PR is the user's call.

## Unarchive

Rare, explicit, reversed: `git mv` back, set `status: active`, remove `release`/`archived`. Reopening usually means new work — prefer a fresh spec that references the archived one.

## Anti-patterns

- Archiving with failing or missing QA silently — the override must be the user's word, on the record.
- Blocking the archive on a merged PR or release tag — completion (tasks + QA) is the gate; publishing is a separate, user-driven step.
- Leaving a completed spec in the active set "until the PR merges" — the active folder is for live work only.
- Editing an archived spec — it is a record; new requirements get a new spec that links back.
- Deleting instead of archiving — the graveyard is where "didn't we already try this?" gets answered.
