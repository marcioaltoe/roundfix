---
type: fix
status: open
created: 2026-09-08
spec: null
reason: null
---

# Capture-time Git instructions conflict with Secondbrain synchronization ownership

## Symptom

The generated guide requires a commit at capture time. The originating
session reports an autosync job that commits and pushes the shared checkout,
leaving individual sessions responsible for overlapping durability operations.

## Where

`docs/agents/secondbrain.md` and
`internal/baseline/assets/modules/secondbrain.json` still require the immediate
commit. Secondbrain `scripts/autosync.sh` stages all changes, commits, rebases,
and pushes. Its `inbox/README.md` still requires capture-time commits.

## Expected

Define one durable capture responsibility and make the canonical guide and
Secondbrain door contract agree. Preserve other sessions' work. Select the
policy against the verified synchronization mechanism; a configured job alone
does not prove a remote push. Route to provisional P1.

## Evidence

Secondbrain `inbox/roundfix/_triaged/2026-08-27-o-guia-manda-a-sessao-commitar-num-brain-que-ja-se-commita-e-se-empurra-sozinho.md` records the original
observation. Triage on 2026-09-08 checked the current local sources named above.
The Inbox body remains the observation's provenance; this entry records intent.

On 2026-09-08, `launchctl print gui/501/com.marcioaltoe.secondbrain-autosync`
found the configured scheduler with state `not running`, 33 recorded runs,
last exit code 0, and interval 7200 seconds. This confirms a configured job and
its recorded last outcome. It does not prove that any new capture or this
triage has reached the remote repository.
