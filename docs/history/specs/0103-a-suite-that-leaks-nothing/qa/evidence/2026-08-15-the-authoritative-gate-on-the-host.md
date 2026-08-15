---
spec: 0103-a-suite-that-leaks-nothing
date: 2026-08-15
build: 4482c01907014f7bbdbd10c5217c28a7873007fb
kind: equivalent-evidence
row: R01, R02
---

# The authoritative gate on the host, at the gate's own build

`make verify` was run from the repository root on the operator's machine at
the same commit the QA gate reported failing, outside any Agent Session.

```text
$ git rev-parse HEAD
4482c01907014f7bbdbd10c5217c28a7873007fb
$ make verify; echo "exit=$?"
ok  	roundfix/skills	7.948s
go test -count=1 ./skills -run 'TestNoPythonBaselineRuntime|TestThinSetupSkill|TestCheckRejectsExecutableSetupEngineArtifacts|TestRecommendedSkillsMatchLock'
ok  	roundfix/skills	0.209s
go run -buildvcs=false ./cmd/roundfix skills check
Roundfix skill check passed: roundfix, write-idea, write-prd, write-techspec, write-tasks, setup-context-driven, implement-task, implement-spec, brainstorming, council, business-analyst, archive-spec, qa-gate, evidence-gate
go build -buildvcs=false -ldflags "-X 'roundfix/internal/app.BuildCommit=4482c019' -X 'roundfix/internal/app.BuildTime=2026-08-14 20:41:27 -0300'" -o bin/roundfix ./cmd/roundfix
exit=0
```

Machine: 10 CPUs, 16 GB RAM. The suite runs `-parallel 16`.

The gate's F-01 reports three ACPX fixture journeys dying under repository-wide
load — a fixture killed before its milestone, an unproven adapter lineage, and a
session close returning exit code -1. The same trio passes 5/5 in isolation and
the whole `internal/agent` package passes at `-parallel 16`. It fails only when
the authoritative suite runs inside an Agent Session, which is itself a dense
process tree layered on a suite already oversubscribed sixteen ways across ten
cores.

This is the spawn-density family the Spec was pulled forward to address, observed
at the one place the Spec cannot reach: the gate's own execution context. The
saturation is recorded as backlog.
