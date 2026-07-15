# Review artifacts are committed in a separate docs commit

Review Issue artifacts become part of the repository's history: after a clean
integration, `resolve` and `watch` create one Daemon-owned docs commit on the
user's checkout (`docs: review round NNN for pr <n>`) carrying every dirty
path under the Run's resolved review artifact root, then run Final Push from
the checkout so the commit rides it. (The original clause fast-forwarding the
Run Branch over the docs commit became moot when ADR-0042 removed Run
Branches from review Runs; the docs commit lands directly on the PR Head
Branch.) The commit stays
separate from Batch fix commits so fix diffs stay clean, `fetch` still never
commits, `auto_commit: false` disables it with everything else, and roots
outside the repository — an explicit external Artifact Directory, an external
Spec Root, or a path crossing a symbolic link — are never staged (ADR-0035
semantics), leaving those artifacts to the owning repository. Supersedes
ADR-0029's never-commit clause; its location hierarchy is unchanged.
