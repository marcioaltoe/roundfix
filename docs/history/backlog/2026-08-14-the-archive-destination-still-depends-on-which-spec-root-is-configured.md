---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-14
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The archive destination still depends on which Spec Root is configured

## Symptom

`roundfix archive` sends a completed Spec to `docs/history/specs/<slug>` when the
Spec Root is the built-in `docs/specs`, and to `<spec-root>/_archived/<slug>` when
it is anything else. Spec 0094 moved the built-in destination and left the
asymmetry: the default root is still the only configuration whose archive lives
outside it, and a non-default root still archives under the pre-0094 name.

Two consequences observed in `conexus` on 2026-08-13, archiving
`0002-acesso-do-dashboard-aos-dados` and `0003-experiencia-financeira`:

- **The archive splits in two.** The repository already kept
  `docs/specs/_archived/` with Spec 0001 and sixteen legacy `phase-*` Specs. The
  two new ones went elsewhere, so `ls docs/specs/_archived/` stopped being the
  complete record of what was built — which is the property the `archive-spec`
  skill states it wants: "one `ls <spec-root>/` separates live work from history".
- **No configuration reaches the wanted path.** Asking for the archive inside the
  Spec Root requires `docs/specs` to be treated as a non-default root, and it is
  precisely that name which triggers the special case.

A second defect rides along and is independent of the destination: relative links
authored inside a Spec break on archival, because the move adds a directory level
and the contract forbids rewriting an archived Spec. A link `../../../findings/x.md`,
correct while the Spec is active, resolves one level short afterwards. The depth
changes identically wherever the archive lives, so no destination fixes it; what
would fix it is the archive rewriting link destinations during the move, the only
edit compatible with preserving the observations byte for byte.

## Where

`ArchiveSpecRoot` in `internal/spec/archive.go`, which branches on whether the
Spec Root is the built-in one, and the Archive Command's move, which relocates
files without touching the links inside them.

## Expected

One rule for every Spec Root, so a repository's archive location is a property of
its configuration rather than of whether that configuration happens to be the
default. Whichever rule wins, adopters holding the other convention need a stated
migration — the Baseline already migrates layouts and could carry this one.

The link rewrite is its own decision: rewriting Markdown link destinations during
the move is the one change that keeps a Spec's prose byte-identical while keeping
its links resolvable.

## Evidence

Minted from the Inbox Entry
`inbox/roundfix/2026-08-13-arquivo-de-spec-cai-fora-da-raiz-configurada.md` in the
Secondbrain. The sibling entry
`inbox/roundfix/2026-08-12-archive-e-o-guia-do-baseline-discordam-sobre-onde-a-spec-arquiva.md`
recorded the same split in `fiscus`, where Specs 0001 to 0006 and Spec 0007 ended
up under two conventions; the guide half of that entry is closed, because Spec
0094 brought `docs/agents/docs-layout.md` into agreement with the command.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
