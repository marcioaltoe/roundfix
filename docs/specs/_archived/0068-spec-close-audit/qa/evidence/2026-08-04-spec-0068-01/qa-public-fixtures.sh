#!/bin/sh
set -eu

qa_binary="/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260805T013130Z_31398370e8ba8670/bin/roundfix"
qa_source="/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260805T013130Z_31398370e8ba8670"
qa_root="$(mktemp -d /private/tmp/roundfix-qa0068-rerun.XXXXXX)"
qa_repo="$qa_root/repository"
qa_remote="$qa_root/origin.git"
qa_worktree="$qa_root/scratch-worktree"
qa_slug="0068-spec-close-audit"
qa_branch="ma/spec-close-scratch-rerun"

rtk git init --initial-branch=main "$qa_repo"
rtk git -C "$qa_repo" config user.name "Roundfix QA"
rtk git -C "$qa_repo" config user.email "roundfix-qa@example.com"
rtk git -C "$qa_repo" config commit.gpgsign false
rtk mkdir -p "$qa_repo/docs/specs/$qa_slug"
rtk cp "$qa_source/docs/specs/$qa_slug/_prd.md" "$qa_repo/docs/specs/$qa_slug/_prd.md"
rtk cp "$qa_source/README.md" "$qa_repo/README.md"
rtk git -C "$qa_repo" add -A
rtk git -C "$qa_repo" commit -m "docs: seed QA fixture"

rtk git init --bare "$qa_remote"
rtk git -C "$qa_repo" remote add origin "$qa_remote"
rtk git -C "$qa_repo" push -u origin main
rtk git -C "$qa_repo" worktree add -b "$qa_branch" "$qa_worktree" main
rtk cp "$qa_source/README.md" "$qa_worktree/scratch.txt"
rtk git -C "$qa_worktree" add scratch.txt
rtk git -C "$qa_worktree" commit -m "feat: add scratch fixture" -m "Roundfix-Spec: $qa_slug"
rtk git -C "$qa_repo" push -u origin "$qa_branch"
rtk git -C "$qa_repo" merge --squash "$qa_branch"
rtk git -C "$qa_repo" commit -m "feat: merge scratch fixture"
rtk git -C "$qa_repo" push origin main

qa_refs_before="$(rtk git -C "$qa_repo" show-ref)"
qa_worktrees_before="$(rtk git -C "$qa_repo" worktree list --porcelain)"

set +e
qa_text_output="$(cd "$qa_repo" && rtk proxy "$qa_binary" spec audit "$qa_slug" 2>&1)"
qa_text_exit=$?
qa_json_output="$(cd "$qa_repo" && rtk proxy "$qa_binary" spec audit "$qa_slug" --format json 2>&1)"
qa_json_exit=$?
qa_repeat_output="$(cd "$qa_repo" && rtk proxy "$qa_binary" spec audit "$qa_slug" --format json 2>&1)"
qa_repeat_exit=$?
set -e

qa_refs_after="$(rtk git -C "$qa_repo" show-ref)"
qa_worktrees_after="$(rtk git -C "$qa_repo" worktree list --porcelain)"
test "$qa_refs_before" = "$qa_refs_after"
test "$qa_worktrees_before" = "$qa_worktrees_after"
test "$qa_text_exit" -eq 1
test "$qa_json_exit" -eq 1
test "$qa_repeat_exit" -eq 1
test "$(rtk proxy jq -r '.schema' <<EOF
$qa_json_output
EOF
)" = "roundfix-specaudit/v1"

qa_reclaim="$(rtk proxy jq -r --arg path "$qa_worktree" '.survivors[] | select(.name == $path and .is_worktree == true and .kind == "residue") | .reclaim' <<EOF
$qa_json_output
EOF
)"
test -n "$qa_reclaim"
case "$qa_reclaim" in
  *"git worktree remove -- '$qa_worktree' && git branch -D -- '$qa_branch'"*) ;;
  *) exit 41 ;;
esac

qa_branch_reclaim="$(rtk proxy jq -r --arg branch "$qa_branch" '.survivors[] | select(.name == $branch and .is_worktree == false) | .reclaim' <<EOF
$qa_json_output
EOF
)"
test "$qa_branch_reclaim" = "$qa_reclaim"

cd "$qa_repo"
rtk proxy sh -c "$qa_reclaim"
test ! -e "$qa_worktree"
test -z "$(rtk git branch --list "$qa_branch")"

qa_clean_root="$(mktemp -d /private/tmp/roundfix-qa0068-clean-rerun.XXXXXX)"
qa_clean_repo="$qa_clean_root/repository"
rtk git init --initial-branch=main "$qa_clean_repo"
rtk git -C "$qa_clean_repo" config user.name "Roundfix QA"
rtk git -C "$qa_clean_repo" config user.email "roundfix-qa@example.com"
rtk git -C "$qa_clean_repo" config commit.gpgsign false
rtk mkdir -p "$qa_clean_repo/docs/specs/$qa_slug"
rtk cp "$qa_source/docs/specs/$qa_slug/_prd.md" "$qa_clean_repo/docs/specs/$qa_slug/_prd.md"
rtk git -C "$qa_clean_repo" add -A
rtk git -C "$qa_clean_repo" commit -m "docs: seed clean QA fixture"
rtk mkdir -p "$qa_clean_repo/docs/specs/_archived"
rtk git -C "$qa_clean_repo" mv "docs/specs/$qa_slug" "docs/specs/_archived/$qa_slug"
rtk git -C "$qa_clean_repo" commit -m "docs: archive clean QA fixture"

set +e
qa_clean_text="$(cd "$qa_clean_repo" && rtk proxy "$qa_binary" spec audit "$qa_slug" 2>&1)"
qa_clean_text_exit=$?
qa_clean_json="$(cd "$qa_clean_repo" && rtk proxy "$qa_binary" spec audit --format json "$qa_slug" 2>&1)"
qa_clean_json_exit=$?
qa_missing="$(cd "$qa_clean_repo" && rtk proxy "$qa_binary" spec audit 2>&1)"
qa_missing_exit=$?
qa_unknown="$(cd "$qa_clean_repo" && rtk proxy "$qa_binary" spec audit no-such-slug --format json 2>&1)"
qa_unknown_exit=$?
qa_format="$(cd "$qa_clean_repo" && rtk proxy "$qa_binary" spec audit "$qa_slug" --format yaml 2>&1)"
qa_format_exit=$?
set -e

test "$qa_clean_text_exit" -eq 0
test "$qa_clean_json_exit" -eq 0
test "$qa_missing_exit" -eq 2
test "$qa_unknown_exit" -eq 2
test "$qa_format_exit" -eq 2
test "$(rtk proxy jq -r '.schema' <<EOF
$qa_clean_json
EOF
)" = "roundfix-specaudit/v1"
test "$(rtk proxy jq '.survivors | length' <<EOF
$qa_clean_json
EOF
)" -eq 0
test "$(rtk proxy jq '.undelivered | length' <<EOF
$qa_clean_json
EOF
)" -eq 0
test -z "$(rtk git -C "$qa_clean_repo" status --short)"

rtk proxy printf '%s\n' "FIXTURE_ROOT=$qa_root"
rtk proxy printf '%s\n' "SCRATCH_TEXT_EXIT=$qa_text_exit"
rtk proxy printf '%s\n' "$qa_text_output"
rtk proxy printf '%s\n' "SCRATCH_JSON_EXIT=$qa_json_exit"
rtk proxy printf '%s\n' "$qa_json_output"
rtk proxy printf '%s\n' "REPEAT_JSON_EXIT=$qa_repeat_exit"
rtk proxy printf '%s\n' "STATE_UNCHANGED_BEFORE_RECLAIM=yes"
rtk proxy printf '%s\n' "EMITTED_RECLAIM=$qa_reclaim"
rtk proxy printf '%s\n' "EMITTED_RECLAIM_EXECUTED=yes"
rtk proxy printf '%s\n' "WORKTREE_REMOVED=yes"
rtk proxy printf '%s\n' "LOCAL_BRANCH_REMOVED=yes"
rtk proxy printf '%s\n' "CLEAN_FIXTURE_ROOT=$qa_clean_root"
rtk proxy printf '%s\n' "CLEAN_TEXT_EXIT=$qa_clean_text_exit"
rtk proxy printf '%s\n' "$qa_clean_text"
rtk proxy printf '%s\n' "CLEAN_JSON_EXIT=$qa_clean_json_exit"
rtk proxy printf '%s\n' "$qa_clean_json"
rtk proxy printf '%s\n' "MISSING_EXIT=$qa_missing_exit"
rtk proxy printf '%s\n' "$qa_missing"
rtk proxy printf '%s\n' "UNKNOWN_EXIT=$qa_unknown_exit"
rtk proxy printf '%s\n' "$qa_unknown"
rtk proxy printf '%s\n' "INVALID_FORMAT_EXIT=$qa_format_exit"
rtk proxy printf '%s\n' "$qa_format"
