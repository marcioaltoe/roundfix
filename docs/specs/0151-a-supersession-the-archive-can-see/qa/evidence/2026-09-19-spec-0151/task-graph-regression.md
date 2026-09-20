# Task Graph archive regression

The public CLI ran from disposable Git root
`/private/tmp/roundfix-qa0151-taskgraph.h5Txpz` against a copy of Spec 0151.

1. `./bin/roundfix archive 0151-a-supersession-the-archive-can-see` exited `2`
   with `Task "task_04" is "pending"; archive requires every Task to be
   "completed"`. A fresh `git status --porcelain` was empty.
2. The public `supersede` command added a valid `_supersession.md` and exited
   `0`.
3. Repeating archive exited `2` with the exact same incomplete-Task diagnostic.
4. `git diff --exit-code` over the Spec while excluding the new amendment
   exited `0`; `git status --porcelain` named only `_supersession.md`.

The amendment therefore did not bypass or alter archive's Task Graph evidence
path.

