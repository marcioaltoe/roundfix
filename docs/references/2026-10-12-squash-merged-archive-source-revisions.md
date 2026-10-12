# Squash-merged Archive Record revisions (2026-10-12)

A measurement made on 2026-10-12, between 10:18 and 10:20 UTC, while Spec 0256
was authored. It asks two questions about this repository's own Archive
Records: is a record's `source_revision` reachable from the default branch
after a squash merge, and can a clone that lacks it recover it from the
forge without writing a ref?

## Setup

An empty repository was created with `git init` in a scratch directory
outside the repository, with `GIT_TERMINAL_PROMPT=0`. It fetched only
`refs/heads/main` of `https://github.com/marcioaltoe/roundfix.git` into
`refs/remotes/origin/main`, which resolved to
`7a8b7f70f464acfbe155540a583b6f509f930849`, the default branch after
pull request #490.

## Results

| record | `source_revision` | delivery commit on `main` | present after the `main` fetch | ancestor of `main` | present after the pull request fetch |
| --- | --- | --- | --- | --- | --- |
| `0252-one-gate-per-run-and-a-smaller-baseline` | `4d3e9059fa3d6704a8c5c0628cca369adbd7c2f7` | `e4dd6a68` "feat: deliver 0252-one-gate-per-run-and-a-smaller-baseline (#481)" | no | no | yes, from `refs/pull/481/head` |
| `0253-authoring-rules-that-stop-qa-reruns` | `1740b715c99db7a1ae7bd78e6bcd3a0b71f51ac7` | `1db414e2` "feat: deliver 0253-authoring-rules-that-stop-qa-reruns (#484)" | no | no | yes, from `refs/pull/484/head` |
| `0254-a-faster-make-test` | `c74b91f92628c31480650b05d4b812685afe2909` | `7b82943b` "feat: deliver 0254-a-faster-make-test (#487)" | no | no | yes, from `refs/pull/487/head` |

The delivery commit is the first commit on `main` that added the record,
found with `git log -1 --diff-filter=A --format='%h %s' refs/remotes/origin/main -- docs/history/specs/<slug>.md`.
None of the three records carries a `pull_request` field. Each delivery
commit's subject ends with the pull request number in the form `(#N)`.

Each recovery ran
`git fetch -q --no-tags --no-write-fetch-head https://github.com/marcioaltoe/roundfix.git refs/pull/<N>/head`
and exited 0. `git cat-file -e <source_revision>^{commit}` then succeeded,
and `git cat-file -e <source_revision>:docs/specs/<slug>/_prd.md` found each
Spec's PRD. For 0254, `git show c74b91f9…:docs/specs/0254-a-faster-make-test/_prd.md`
printed the PRD front matter (`spec: 0254-a-faster-make-test`,
`status: active`). After all three fetches `git for-each-ref` listed only
`refs/remotes/origin/main`: a fetch with a source and no destination wrote no
ref. `git merge-base --is-ancestor c74b91f9… refs/remotes/origin/main` exited
1.

## What it shows

- A squash merge leaves an Archive Record's `source_revision` outside the
  default branch's ancestry. A clone of the default branch cannot read the
  removed Spec folder, so every reader that opens the folder at
  `source_revision` reports it unavailable.
- The forge keeps the revision under the pull request's head ref, and a
  fetch of that ref recovers it without creating a local ref.
- The pull request number is known after the merge from the delivery
  commit's subject, even for records that do not store it.

## Published sources

- GitHub Docs, "Checking out pull requests locally"
  (<https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/checking-out-pull-requests-locally>):
  "Commits in a pull request are available in a repository even before the
  pull request is merged." The same page states: "The remote `refs/pull/`
  namespace is read-only."
- Git documentation, `git fetch` (<https://git-scm.com/docs/git-fetch>):
  "Passing `--no-write-fetch-head` from the command line tells Git not to
  write the file."
