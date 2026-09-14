#!/bin/sh
set -eu

real_git=/opt/homebrew/bin/git

measure() {
  mode=$1
  repo=$(mktemp -d "/private/tmp/roundfix-qa-0136-git-${mode}.XXXXXX")
  "$real_git" -C "$repo" init --initial-branch=main >/dev/null
  "$real_git" -C "$repo" config user.name 'Roundfix QA'
  "$real_git" -C "$repo" config user.email 'roundfix-qa@example.com'
  "$real_git" -C "$repo" config commit.gpgsign false
  printf '%s\n' 'tracked' > "$repo/before.txt"
  "$real_git" -C "$repo" add before.txt
  "$real_git" -C "$repo" commit -m initial >/dev/null

  if [ "$mode" = staged ]; then
    "$real_git" -C "$repo" mv before.txt after.txt
  else
    mv "$repo/before.txt" "$repo/after.txt"
  fi

  stderr_path="$repo/git-add.stderr"
  set +e
  "$real_git" -C "$repo" add -f -- before.txt 2>"$stderr_path"
  code=$?
  set -e

  printf 'mode=%s\n' "$mode"
  printf 'git_add_source_exit=%s\n' "$code"
  printf '%s\n' 'stderr:'
  sed -n '1,10p' "$stderr_path"
  printf '%s\n' 'status:'
  "$real_git" -C "$repo" status --porcelain=v1
}

"$real_git" --version
measure staged
measure unstaged
