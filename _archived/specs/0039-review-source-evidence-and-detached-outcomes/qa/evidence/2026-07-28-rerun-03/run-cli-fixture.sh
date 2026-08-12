#!/bin/sh

set -eu

if [ "$#" -ne 1 ]; then
  printf 'usage: %s <watch-skip|watch-transient|watch-failure|watch-detached|fetch-zero>\n' "$0" >&2
  exit 2
fi

mode="$1"
fixture_dir="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
repo_root="$(CDPATH= cd -- "$fixture_dir/../../../../../../" && pwd)"
scratch_repo="/private/tmp/roundfix-qa0039-rerun-03-${mode}-repo"
artifact_dir="/private/tmp/roundfix-qa0039-rerun-03-${mode}-artifacts"
state_file="/private/tmp/roundfix-qa0039-rerun-03-${mode}-state"

if [ -e "$scratch_repo" ] || [ -e "$artifact_dir" ] || [ -e "$state_file" ]; then
  printf 'QA scratch path already exists for %s\n' "$mode" >&2
  exit 96
fi

rtk git clone --shared "$repo_root" "$scratch_repo"
rtk mkdir -p "$artifact_dir"

PATH="$fixture_dir:$PATH"
QA_EXPECTED_HEAD="6646b84845d32335f35c39df95a6c161a973c9ec"
QA_GH_LOG="/private/tmp/roundfix-qa0039-rerun-03-${mode}-gh.log"
QA_NOTIFY_LOG="/private/tmp/roundfix-qa0039-rerun-03-${mode}-notify.log"
export PATH QA_EXPECTED_HEAD QA_GH_LOG QA_NOTIFY_LOG

case "$mode" in
  watch-transient)
    QA_GH_MODE="transient_then_skip"
    QA_GH_STATE="$state_file"
    export QA_GH_MODE QA_GH_STATE
    ;;
  watch-failure)
    QA_GH_MODE="permanent_failure"
    export QA_GH_MODE
    ;;
  watch-skip|watch-detached|fetch-zero)
    ;;
  *)
    printf 'unknown mode: %s\n' "$mode" >&2
    exit 2
    ;;
esac

cd "$scratch_repo"

if [ "$mode" = "fetch-zero" ]; then
  exec "$repo_root/bin/roundfix" \
    fetch \
    --source coderabbit \
    --pr 123 \
    --base-repo owner/repo \
    --head-repo owner/repo \
    --head-branch roundfix/run-run_20260728T171509Z_8bf024d33e949d37 \
    --artifact-dir "$artifact_dir" \
    --no-input
fi

if [ "$mode" = "watch-detached" ]; then
  exec "$repo_root/bin/roundfix" \
    watch \
    --source coderabbit \
    --pr 123 \
    --base-repo owner/repo \
    --head-repo owner/repo \
    --head-branch roundfix/run-run_20260728T171509Z_8bf024d33e949d37 \
    --artifact-dir "$artifact_dir" \
    --until-clean \
    --max-rounds 1 \
    --no-input \
    --agent claude \
    --model claude-fable-5 \
    --reasoning-effort xhigh \
    --detach
fi

exec "$repo_root/bin/roundfix" \
  watch \
  --source coderabbit \
  --pr 123 \
  --base-repo owner/repo \
  --head-repo owner/repo \
  --head-branch roundfix/run-run_20260728T171509Z_8bf024d33e949d37 \
  --artifact-dir "$artifact_dir" \
  --until-clean \
  --max-rounds 1 \
  --no-input \
  --agent claude \
  --model claude-fable-5 \
  --reasoning-effort xhigh
