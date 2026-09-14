#!/bin/sh
set -eu

mode=${1:?mode is required}
binary=${2:?rebuilt roundfix binary is required}
evidence_root=${3:?evidence output directory is required}
script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
scratch=$(mktemp -d "/private/tmp/roundfix-qa-0136-${mode}.XXXXXX")
repo=$scratch/repo
home_dir=$scratch/home
remote=$scratch/remote.git
branch="feat/qa-${mode}"
real_git=/opt/homebrew/bin/git

mkdir -p "$repo" "$home_dir" "$evidence_root/$mode"
cp -R "$script_dir/fixture/." "$repo/"

case "$mode" in
  governed_rename)
    cp "$script_dir/authorizations/commit-only.md" "$repo/docs/specs/0001-rename-fixture/_authorization.md"
    ;;
  *)
    cp "$script_dir/authorizations/implement-commit.md" "$repo/docs/specs/0001-rename-fixture/_authorization.md"
    ;;
esac

case "$mode" in
  ordinary_push|governed_push|governed_modify_push)
    cp "$script_dir/config/auto-push.yml" "$repo/.roundfixrc.yml"
    ;;
  *)
    cp "$script_dir/config/no-push.yml" "$repo/.roundfixrc.yml"
    ;;
esac

"$real_git" -C "$repo" init --initial-branch=main >/dev/null
"$real_git" -C "$repo" config user.name 'Roundfix QA'
"$real_git" -C "$repo" config user.email 'roundfix-qa@example.com'
"$real_git" -C "$repo" config commit.gpgsign false
"$real_git" -C "$repo" add -A
"$real_git" -C "$repo" commit -m 'seed QA fixture' >/dev/null
"$real_git" -C "$repo" switch -c "$branch" >/dev/null

case "$mode" in
  ordinary_push|governed_push|governed_modify_push)
    "$real_git" init --bare "$remote" >/dev/null
    "$real_git" -C "$repo" remote add origin "$remote"
    "$real_git" -C "$repo" push -u origin "HEAD:$branch" >/dev/null
    ;;
esac

git_failure=""
case "$mode" in
  unavailable_revision|unresolvable_spec_root)
    git_failure=$mode
    ;;
esac

stdout_path="$scratch/stdout.txt"
stderr_path="$scratch/stderr.txt"
set +e
(
  cd "$repo"
  PATH="$script_dir/fake-bin:$PATH" \
  HOME="$home_dir" \
  GIT_CONFIG_GLOBAL=/dev/null \
  GIT_CONFIG_SYSTEM=/dev/null \
  ROUNDFIX_TUI=never \
  ROUNDFIX_QA_MODE="$mode" \
  ROUNDFIX_QA_GIT_FAILURE="$git_failure" \
  "$binary" implement --spec 0001-rename-fixture --no-input \
    --agent codex --model gpt-5.6-sol --reasoning-effort high \
    --agent-command "$script_dir/fake-bin/codex-acp" \
    >"$stdout_path" 2>"$stderr_path"
)
exit_code=$?
set -e

cp "$stdout_path" "$evidence_root/$mode/stdout.txt"
cp "$stderr_path" "$evidence_root/$mode/stderr.txt"
"$real_git" -C "$repo" log -1 --format='%H%n%s%n%B' >"$evidence_root/$mode/head.txt"
"$real_git" -C "$repo" diff-tree --no-commit-id --name-status --no-renames -r HEAD >"$evidence_root/$mode/changed.txt"
"$real_git" -C "$repo" status --porcelain=v1 >"$evidence_root/$mode/status.txt"
printf '%s\n' "$exit_code" >"$evidence_root/$mode/exit-code.txt"
printf '%s\n' "$scratch" >"$evidence_root/$mode/scratch-path.txt"

if [ -d "$remote" ]; then
  "$real_git" --git-dir="$remote" rev-parse "$branch" >"$evidence_root/$mode/remote-head.txt"
fi

printf 'mode=%s\nexit=%s\nscratch=%s\n' "$mode" "$exit_code" "$scratch"
printf '%s\n' 'stdout:'
sed -n '1,80p' "$stdout_path"
printf '%s\n' 'stderr tail:'
tail -n 40 "$stderr_path"
printf '%s\n' 'head:'
sed -n '1,20p' "$evidence_root/$mode/head.txt"
printf '%s\n' 'changed:'
sed -n '1,20p' "$evidence_root/$mode/changed.txt"
printf '%s\n' 'status:'
sed -n '1,20p' "$evidence_root/$mode/status.txt"
if [ -f "$evidence_root/$mode/remote-head.txt" ]; then
  printf '%s\n' 'remote head:'
  sed -n '1,5p' "$evidence_root/$mode/remote-head.txt"
fi
