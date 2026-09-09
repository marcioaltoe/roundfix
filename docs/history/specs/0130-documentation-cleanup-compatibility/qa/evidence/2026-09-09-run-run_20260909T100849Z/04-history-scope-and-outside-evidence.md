# History, scope and outside evidence

## Outside evidence

Read-only GitHub observation of Actions job `102229198572` in run
`34276050740` reported workflow `CI | Verify`, branch
`chore/remove-obsolete-workflow-docs`, head `795bcfa834d95265fe0d036ac5f3deb8f4d60e15`
and conclusion `failure`. Its failed `make verify` log named:

- missing `docs/workflow/authorizations/2026-08-06-proof-cost.md` from
  `TestOutputsForCommand`; and
- stale catalog, formatter and plan-characterization artifacts, each directing
  the operator to `make baseline-digests`.

Source: <https://github.com/marcioaltoe/roundfix/actions/runs/34276050740/job/102229198572>.

## History and tooling scope

- `rtk git cat-file -e '81a6afb48f4a3683d0e5fad52f3919cf1bdfbbf4^{commit}'`
  exited 0. Listing `docs/workflow/authorizations` at that object produced 42
  paths.
- Both historical Task object checks exited 128 with `Not a valid object name`:
  `419a4661ac769ff7ee6ce5423bd795185c859d01` and
  `c80e1266658929f68e8046af82f88e13392dc56d`.
- `rtk git merge-base --is-ancestor 54614682 c6a5ea15` exited 0. The authorization
  therefore precedes the implementation commit.
- `rtk git diff-tree --no-commit-id --name-only -r c6a5ea15` listed only
  `task_01.md` and the five paths bounded by `_authorization.md`.
- The authorization commit changed only the six Spec authoring artifacts. There
  is no post-Task descendant commit in the audited range, so no prerequisite or
  consequent fix was folded into or ordered around the tooling change.
- `rtk git diff --exit-code 54614682..c6a5ea15 -- CONTEXT.md` and the equivalent
  `go.mod go.sum` check exited 0. No glossary term or dependency changed.


## 2026-09-09 — complete historical-object inventory

An independent `git cat-file -e <sha>^{commit}` reread in the clean audited clone
returned exit 128 for each of the four configured historical Task objects:

| Task object | Result | Dependent historical case |
| --- | --- | --- |
| `419a4661ac769ff7ee6ce5423bd795185c859d01` | unavailable | grouped authorization audit |
| `65c51ebf2e19220ff50d25fe03be809fcdf353f0` | unavailable | grouped authorization audit |
| `28acf39cc193ad490646cb5a1d23500e0c08c273` | unavailable | grouped authorization audit |
| `c80e1266658929f68e8046af82f88e13392dc56d` | unavailable | regeneration audit |

The two historical subtests remain skipped, with the first missing object named
for each case. This inventory does not manufacture four executed subtests.
The 42 recovered real grants and controlled accepted/refused Git-change cases
remain the executed coverage recorded above.
