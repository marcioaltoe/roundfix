---
status: pending
type: backend
---

# Task: Handle Deleted File Staging

Ensure deleted files stage correctly.

## Work
- Verify `git add --all` handles deletions
- Test with a Task that deletes a file
- Confirm deletion commits cleanly

## Verification
- `git_test=$(mktemp -d) && cd "$git_test" && git init -q . && git config user.email t@t && git config user.name t && echo hi > file.txt && git add file.txt && git commit -qm init && rm file.txt && git add --all && git diff --cached --name-status | grep -q "^D" && cd - >/dev/null`


## References

- Core Feature 3: Deleted File Handling

## Result
Deleted files stage without errors.
