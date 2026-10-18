---
status: accepted
created_at: 2026-10-18T18:00:00Z
updated_at: 2026-10-18T18:00:00Z
deprecated_at: null
superseded_by: null
---

# QA evidence holds text and images, and the QA Report commit refuses any other binary

On 2026-10-18 the QA gate of a Spec in the Conexus repository ran four times.
The `qa-gate` skill tells the gate to store screenshots under
`qa/evidence/`, the Daemon committed the PNGs with each QA Report, and the next
gate's mechanical stage refused every one of them as `QA-EVIDENCE-PATH` "is a
built binary", because any blob with a NUL byte in its first 8000 bytes counts
as a binary and every PNG has one. One full QA run was lost to evidence the
flow itself created. A text-only requirement in the QA Task did not help: the
fourth run's report commit still versioned 26 PNGs that the gate's browser
script wrote by itself, and the mechanical stage, which reads tracked files
before the QA Agent's turn, could not see them. On 2026-10-18 the maintainer
chose "Aceitar imagens (Recomendado)": "O detector aceita PNG, JPEG e SVG como
evidência legível e recusa só binário de build ou executável. O Daemon recusa
no commit do relatório qualquer binário que não seja imagem."

## Decision

**Readable QA evidence is text, PNG or JPEG content.** One classifier decides,
from the first 8000 bytes of a file:

- content with no NUL byte is text, and an SVG image is text;
- content that begins with the PNG signature `89 50 4E 47 0D 0A 1A 0A` is a
  PNG image;
- content that begins with `FF D8 FF`, the JPEG Start Of Image marker followed
  by the next marker, is a JPEG image;
- anything else is a binary that is not an image.

The content decides, never the file name, so a renamed executable is still a
binary and a screenshot saved without an extension is still an image.

**The mechanical stage accepts images.** A tracked file under the QA Report's
`evidence/` directory raises `QA-EVIDENCE-PATH` only when it is a gitlink or a
binary that is not an image. The stage reads tracked files before the QA
Agent's turn, as before.

**The QA Report commit refuses what the next mechanical stage would refuse.**
After the Secret Scan of the paths the QA Report commit would stage finds
nothing, the Daemon reads each of those paths under the report's `evidence/`
directory with the same classifier. A binary that is not an image settles the
QA Task `failed` with the reason `QA evidence carries a non-image binary:
<path>; <path>`, publishes one QA Run Event that names the paths, and creates
no QA Report commit; the report and its evidence stay uncommitted in the Run
Worktree, exactly as a Secret Match leaves them (ADR-0273). The defect is
reported by the gate that created it, not by the next one.

**No size bound is added.** The code has no limit a QA evidence file could
inherit: the pre-PR review omits `qa/evidence/` from its diff, and the archive
removes the Spec folder from the tree. The remote's own file limits still
apply at the push.

## Considered Options

- Forbidding binaries in `qa/evidence/` and teaching text-only visual evidence
  was rejected. The `qa-gate` skill's frontend surface asks for key-state and
  failure screenshots, and Conexus showed that a text-only requirement in the
  QA Task does not stop a browser script that saves captures by itself.
- Deciding by file extension was rejected: a name is not content, and a binary
  saved as `.png` would pass.
- Accepting every binary a viewer can open, such as PDF, video or a trace
  archive, was rejected. The maintainer chose images, and this repository's
  own history shows what else a gate commits: Spec 0090's QA evidence held
  SQLite Run Databases from temporary homes, which is the scratch state the
  check exists to refuse.
- Dropping a non-image binary from the QA Report commit and settling the gate
  was rejected: the report may cite the file, and a silently missing citation
  is worse than a named refusal.
- Requiring an XML prefix for SVG was rejected as redundant: content with no
  NUL byte is already accepted as text, and an SVG that held one would not be a
  valid SVG.

## Consequences

- Screenshots a gate saves under `qa/evidence/` are committed and accepted by
  every later gate of the same Spec.
- A gate that saves a video, a trace archive, a PDF or a database under
  `qa/evidence/` fails at its own report commit, and the reason names the
  path. The gate records such evidence as text or keeps it outside the Spec.
- Images committed with QA evidence stay in the item branch's history after
  the archive removes the Spec folder.
- This refines ADR-0273's QA Report commit, which already refuses a Secret
  Match, with a second refusal at the same boundary.
