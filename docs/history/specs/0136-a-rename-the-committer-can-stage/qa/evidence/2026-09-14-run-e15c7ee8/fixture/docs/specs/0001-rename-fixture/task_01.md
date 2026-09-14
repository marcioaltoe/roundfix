---
task: task_01
spec: 0001-rename-fixture
status: pending
type: backend
---

# Exercise the Daemon commit path

## Verification

- `case "$ROUNDFIX_QA_MODE" in staged_rename|unstaged_rename) test -f after.txt && test ! -e before.txt ;; delete) test ! -e before.txt ;; governed_rename|unavailable_revision|unresolvable_spec_root|governed_push) test -f notes.txt && test ! -e Makefile ;; ordinary_push) test "$(sed -n '1p' ordinary.txt)" = changed ;; governed_modify_push) grep -q '^qa:' Makefile ;; *) exit 1 ;; esac` — expected: passes only after the fixture Agent performs the selected mutation.
