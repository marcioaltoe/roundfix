---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An owned skill's version names its content, and its minimum is the bundle

Roundfix compares the version an installed owned skill declares with a minimum
the binary holds. The minimum was one literal, `0.0.2`, for all 14 owned
skills, while the bundle already carried `qa-gate` and `write-tasks` at `0.0.3`
and `implement-spec` at `0.1.0`. A repository with an older copy passed Doctor,
and the `baseline update` preview called it current. Nothing stopped a skill's
content from changing under one version, so equal versions did not mean equal
content. The `go-cli` and `rust-cli` setups did not list the Roundfix skill, and
no release step re-read skills or guides.

The minimum version of an owned skill is now the version the binary carries
for it. It is read from the embedded `SKILL.md`; there is no second list. An
installed copy below it fails Doctor's `skills:` line, and the `baseline
update` preview lists it.

A version names one content. A record in the repository holds, for each owned
skill, every version it has shipped and the digest of that version's folder. A
test refuses a skill whose digest differs from the one recorded for its
version, and a version that is not recorded. Recording only adds a new, higher
version; it never replaces a recorded digest.

The `minimumVersion` in a setup snapshot stays the catalog's declared lower
bound. It does not follow each skill release, because a skill edit must leave
the catalog byte-identical.

Which owned skills a setup lists is Roundfix's decision. The asset sync keeps
an owned entry the snapshot records when the upstream list omits it, and still
drops an external one.

Every release runs a skills-and-guides check before its Pull Request. The step
rests on checks that need no model and no network.

## Consequences

After a binary upgrade Doctor fails until the repository refreshes its owned
skills, which one command does. Every edit to an owned skill now needs a
version change and a line in the record, in the same change. The preview
reports only an installed skill that is older; a repository with no installed
owned skill is still Doctor's to report, because counting it changed the
result of repositories that never installed skills. The setup snapshots keep a
minimum that is lower than the binary's, and the two are no longer expected to
agree.
