#!/bin/sh
set -eu

rtk proxy env GOCACHE=/private/tmp/roundfix-spec0080-repro-gocache \
  rtk go test ./internal/cli \
  -run '^TestRunImplementQAVerdictMatrix/pass$' \
  -count=1 -v
