---
status: accepted
created_at: 2026-11-01T22:20:00Z
updated_at: 2026-11-01T22:20:00Z
deprecated_at: null
superseded_by: null
---

# A QA Report settles only when its local links resolve in the committed tree

On 2026-11-01 Argus reported a QA Report that settled `completed`, a
qualifying partial in a Clean Run, while two of its evidence links pointed at
`.log` files the repository's `.gitignore` excludes. The QA Report commit
never contained them, the Run Worktree that held them was removed, and after
integration `roundfix spec check --strict` reported `SC-HISTORY-RETAINED` for
both links. Nothing between the gate's turn and the commit read the report's
links, so the loss surfaced one Spec later, in a check that could only point
at it. The maintainer chose "Recusar no settle (Recommended)": "Antes do
settle da QA, cada link local de evidência precisa existir na árvore que vai
ser commitada; link ignorado ou ausente falha com o caminho e o motivo, e o
Agente tem um turno de reparo. A skill qa-gate manda salvar evidência com
extensão que o repositório não ignora (.txt)."

## Decision

**A local link resolves when the committed tree would hold its target.** A
local link is a Markdown link of the QA Report, inline or a reference
definition, outside code, whose destination is a relative path; URLs,
anchors and absolute paths are not local links. Its target is resolved
against the report's directory. The target resolves when it exists in the
working tree and Git tracks it or would add it: `git ls-files --cached
--others --exclude-standard` lists the path, or a file under it for a
directory. Otherwise the link is an **Unresolved Evidence Link** with one
reason:

- `ignored`: the target exists, but Git ignores it and does not track it;
- `missing`: the target does not exist, or it lies outside the repository.

Outside a Git repository, a target resolves when it exists.

**One extraction serves this rule and `SC-HISTORY-RETAINED`.** Both read the
links of a Markdown file through `spec.RelativeLinks`, the scanner that
already excludes code and URLs, so a link the QA settlement accepts is a link
the later history check reads the same way.

**Every QA settlement path applies it.** `roundfix qa-report accept`, the
Settle Command for a `qa` Task and the Daemon's QA settlement refuse a report
with an Unresolved Evidence Link, naming each link as written, its line and
its reason. The CLI paths apply it after the eligibility decision, so an
ineligible report keeps its current refusal.

**The Daemon gives the gate one repair turn.** When the QA turn ends with an
Unresolved Evidence Link in the newest report, the Daemon returns the same
Agent Session once with each link and reason and the instruction to keep the
evidence under a name the repository does not ignore or to stop citing a file
it did not keep, without changing a verdict or a row the evidence supports.
At the QA Report commit boundary the Daemon checks again: a remaining link
settles the QA Task `failed` with `QA Report evidence links do not resolve in
the committed tree: <link> (line <n>, <reason>); …`, publishes one QA Run
Event that names the links, and creates no QA Report commit, so the report
and its evidence stay in the Run Worktree as a Secret Match or a non-image
binary leaves them.

**`roundfix archive` is unchanged.** It judges a committed report, whose
missing links `SC-HISTORY-RETAINED` already reports.

## Considered Options

- Stage ignored evidence with `git add -f`. The repository's ignore rules
  say what must not be versioned, build logs and local state among them; a
  forced add would commit what the project excluded, and could carry
  content the Secret Scan and binary classifier were never designed to bound.
- Rename ignored evidence in the Daemon. The report's links and the
  evidence would then differ from what the gate wrote, and the Daemon would
  edit Agent-owned report text.
- Warn and settle. The Argus report already settled with a later check
  pointing at the loss; a warning keeps the loss and moves it nowhere.
- Leave the check to `spec check`. It runs after integration, when the Run
  Worktree that held the evidence is gone.
- No repair turn. The gate that wrote the link is the only actor that knows
  whether the evidence can be kept under another name, and a refusal without
  a turn costs a whole QA Run for a rename.
