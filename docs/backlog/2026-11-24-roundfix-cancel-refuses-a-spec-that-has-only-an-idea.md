---
type: fix
status: open
created: 2026-11-24
spec: null
reason: null
---

# `roundfix cancel` refuses a Spec that has only an `_idea.md`

## Symptom

The Quaestor project had Spec `0125-quaestor-sem-google-drive` with only `_idea.md` (an approved `write-idea`, `status: proposed`) and moved it to another repository. To close it here it ran `roundfix cancel 0125-quaestor-sem-google-drive --reason "..."` on Roundfix 0.148.1 and got:

```text
Preflight failed — read cancellation PRD: open .../docs/specs/0125-quaestor-sem-google-drive/_prd.md: no such file or directory
```

`docs-layout.md` says to cancel through the runtime and never delete the folder by hand, but the cancel does not cover a Spec that holds only an Idea. The workaround was to keep the folder with `status: parked` and a note naming the destination, and to swap the links to retired backlog entries for Git Locators.

## Expected

Either `roundfix cancel` accepts a Spec with only `_idea.md` (reading the title from the Idea), or `write-idea` and `docs-layout` say how to close or move an Idea that will not become a PRD in this repository.

## Source

Cross-project message `01a1f699-52e1-7a22-99b0-9e0c74f1b4dd` from the Quaestor project, 2026-11-19 (dated 2026-11-20 in the inbox). Nature: measured (command output on Roundfix 0.148.1). No urgency stated. Recorded by triage on 2026-11-24.
