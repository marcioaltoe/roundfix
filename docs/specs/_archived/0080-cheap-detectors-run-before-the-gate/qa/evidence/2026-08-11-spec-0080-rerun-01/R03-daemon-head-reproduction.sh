#!/bin/sh
set -eu

GOCACHE="${TMPDIR:-/tmp}/roundfix-spec0080-rerun-gocache" \
  go test ./internal/daemon \
    -run '^TestQAGatePromptUsesTaskContextBuilderAndPreviousReportIdentity$' \
    -count=1 \
    -v
