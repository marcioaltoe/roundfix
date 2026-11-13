---
schema: roundfix/archive-record/v1
spec: 0384-the-makefile-follows-the-go-build-cache-set-with-go-env
title: The Makefile follows the Go build cache set with go env
status: archived
created: "2026-11-12"
archived: "2026-11-13"
disposition: partial
source: docs/specs/0384-the-makefile-follows-the-go-build-cache-set-with-go-env
source_revision: fef2100c91e545787d539b07d30e4559a5cef8e1
qa_task: task_03
qa_report: qa-report-2026-11-13.md
qa_verdict: partial
unproven:
  - the operator's run of `make build` in the main checkout of the Mac mini after the merge, with `find` over the three directories for entries newer than a marker file and `du` of each, recorded session evidence.
adrs:
  - ADR-0394
  - ADR-0370
sources:
  - 2026-11-12-gesttione-drive-the-roundfix-go-build-cache-fills-the-mini-internal-disk.md
  - 2026-11-12-the-writer-of-the-default-go-build-cache-on-the-mini.md
regeneration:
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/setup-context.json
promoted: []
type: fix
---

# The Makefile follows the Go build cache set with go env

The Makefile now follows a usable cache configured with go env under the existing 10240 MiB bound and warns once before using the 2048 MiB temporary fallback when that cache is unusable, without creating a missing mount path. The seven cache cases, baseline comparison, six sabotage probes, documentation checks and exact Task scope audit support this behavior; the real post-merge Mac mini observation remains unproven and the pre-PR control remains blocked. Offline archive planning completed without judge contact; its candidates are the two adopted source records, whose reusable cache policy is already recorded in ADR-0394, the glossary and repository guide, so no additional promotion is needed.
