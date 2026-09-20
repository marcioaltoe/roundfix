# Real Spec 0128 supersession journey

Outside evidence came from repository content this Spec did not author:
`docs/specs/0128-release-planning-with-bare-stable-tags/` and
`docs/history/specs/0147-a-planner-that-reads-both-tag-spellings/`.
Both were copied unchanged into the disposable Git root
`/private/tmp/roundfix-qa0151-real.ynHhZQ`.

1. `./bin/roundfix archive 0128-release-planning-with-bare-stable-tags`
   exited `2`. The diagnostic named the absent `_tasks.md` and instructed the
   operator to run the write-tasks workflow. It reported no Run, Agent, commit,
   or push side effect.
2. `./bin/roundfix supersede --spec 0128-release-planning-with-bare-stable-tags --by 0147-a-planner-that-reads-both-tag-spellings --reason 'Spec 0147 delivered the complete content of Spec 0128.'`
   exited `0` and printed the two slugs.
3. A fresh copy of the amended active Spec was saved outside the repository.
4. A second archive invocation exited `0` and moved the Spec to
   `docs/history/specs/0128-release-planning-with-bare-stable-tags`.
5. `diff -qr` between the saved amended tree and the archived destination
   exited `0` with no differences. A fresh `find` confirmed the archived tree
   contains `_supersession.md` and every pre-existing PRD, TechSpec,
   authorization, and reference file.

The two CLI processes and the independent destination read prove persistence;
the archived folder remained byte-identical to the amended source.

