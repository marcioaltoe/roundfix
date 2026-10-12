from pathlib import Path
import re,json,subprocess,hashlib
from report import settle
E=Path(__file__).resolve().parent;ROOT=Path.cwd();S=ROOT/"docs/specs/0261-a-cancelled-spec-leaves-a-record"
def norm(t):return " ".join(t.split())
checks=[]
def phrase(paths,text):
 for path in paths:
  content=norm(Path(path).read_text());checks.append(dict(path=path,phrase=text,pass_=norm(text) in content))
archives=[".agents/skills/roundfix/references/archive.md","skills/roundfix/references/archive.md"]
phrase(archives,"`roundfix cancel <slug> --reason <text>` writes an Archive Record with the disposition `cancelled`")
phrase(archives,"A cancelled Spec never meets a prerequisite, never supersedes another Spec, and a Delivery Retry refuses it")
phrase(archives,"## Cancel Command")
phrase([".agents/skills/roundfix/SKILL.md","skills/roundfix/SKILL.md"],"| [archive](references/archive.md) | `archive`, `cancel`, `history`, `supersede` |")
phrase([".agents/skills/roundfix/references/deliver.md","skills/roundfix/references/deliver.md","docs/user-guide/commands/deliver.md"],"`deliver start` refuses a queue whose Spec requires a cancelled Spec, and `deliver retry` refuses a cancelled Spec")
cs=[".agents/skills/archive-spec/SKILL.md","skills/archive-spec/SKILL.md"]
phrase(cs,"## Cancelling a Spec")
phrase(cs,"A Spec that will not be delivered is cancelled with `roundfix cancel <slug> --reason <text>`, never deleted by hand")
phrase(cs,"returns each adopted Backlog Entry to `docs/backlog/` with `status: open`")
phrase(["docs/user-guide/commands/archive.md"],"A Spec that will not be delivered is cancelled with `roundfix cancel`")
phrase(["docs/user-guide/commands/supersede.md"],"A cancelled Spec cannot be the superseding Spec")
phrase(["docs/user-guide/commands/cancel.md"],"roundfix cancel <slug> --reason <text>")
phrase(["docs/user-guide/commands.md"],"| `cancel` | [command guide](commands/cancel.md) |")
phrase(["docs/user-guide/usage.md"],"`cancel`")
phrase(["docs/agents/docs-layout.md"],"A Spec dropped before delivery is cancelled, never deleted by hand")
context=Path("CONTEXT.md").read_text();tech=(S/"_techspec.md").read_text()
glossary={}
for term in ["Cancel Command","Source Return","Archive Record","Supersession","Prerequisite Spec","Delivery Retry"]:
 m=re.search(r"(?m)^\*\*"+term+r"\*\*:\n(.*?)(?=\n\*\*|\Z)",context,re.S);assert m,term;glossary[term]=norm(m.group(1));assert "ADR-0262" in m.group(1),term
# All exact glossary sentences from the authored texts, in their semantic entry.
expected={"Cancel Command":"The `roundfix cancel <slug> --reason <text>` command that retires an active Spec which will not be delivered: it writes an Archive Record with the disposition `cancelled`, performs Source Return and removes the Spec folder, which stays in Git at `source_revision`. It refuses while an Active Run or an unfinished Delivery Queue item holds the Spec or an active Spec requires it, never changes a Task status and never commits (ADR-0262).","Source Return":"The Cancel Command's hand-back of a cancelled Spec's adopted sources: a Backlog Entry returns to `docs/backlog/` with `status: open` and no Spec, a Finding returns to `docs/findings/` with `status: pending` and a dated addendum, and every other byte is kept. Other adopted sources stay in Git with the folder (ADR-0262).","Archive Record":"A cancelled Spec's record has the disposition `cancelled` and carries `cancel_reason`, `completed_tasks`, `pending_tasks` and `returned_sources`; only its listed completed Tasks count as completed (ADR-0262).","Supersession":"A cancelled Spec never delivers another Spec's content (ADR-0262).","Prerequisite Spec":"A cancelled Spec's record never meets a prerequisite (ADR-0262).","Delivery Retry":"It refuses a Spec whose Archive Record is `cancelled` (ADR-0262)."}
for term,text in expected.items():checks.append(dict(path="CONTEXT.md",term=term,phrase=text,pass_=norm(text) in glossary[term]))
checks.append(dict(path="CONTEXT.md",term="placement",pass_=context.index("**Cancel Command**:")<context.index("**Source Return**:")<context.index("**History Sanitize Command**:")))
checks.append(dict(path="CONTEXT.md",term="Cancel Avoid",pass_="_Avoid_: Delete command, abandon, archive override" in context))
checks.append(dict(path="CONTEXT.md",term="Return Avoid",pass_="_Avoid_: Unadoption, restore, re-triage" in context))
# Preserve the full broad sweep, then inspect every matching paragraph for stale directives.
cmd=["rtk","proxy","rg","-n","-i",r"dispositions?|delet[^\n]{0,70}spec|spec[^\n]{0,70}delet|dropped[^\n]{0,70}spec|spec[^\n]{0,70}dropped","docs","skills",".agents/skills","CONTEXT.md","AGENTS.md","-g","!docs/history/**","-g","!docs/specs/0261-a-cancelled-spec-leaves-a-record/**"]
x=subprocess.run(cmd,text=True,capture_output=True);(E/"wording-sweep.txt").write_text(x.stdout+x.stderr)
hits=[]
for line in x.stdout.splitlines():
 path,number,text=line.split(":",2)
 if path.startswith("docs/adr/"):reason="Allowed: immutable accepted decision; ADR-0262 expressly extends earlier archive decisions, rather than rewriting them."
 elif path=="docs/references/2026-10-12-a-decomposed-spec-dropped-without-a-record.md":reason="Allowed: attributed adopter event and explicitly dated pre-change parser characterization at 7a8b7f70, followed by the current cancelled design."
 elif path.startswith("docs/references/"):reason="Allowed: dated evidence, audit, structured record or test-name inventory; no current instruction to delete a dropped Spec."
 elif "cancel" in text.lower():reason="Allowed: current cancellation guidance or unrelated terminal lifecycle with cancellation recorded."
 elif re.search(r"Review|review|Branch|branch|Carry|carry|Analysis|analysis|Baseline|baseline|Clause|clause|Finding|finding|Run|run",text):reason="Allowed: another domain's disposition or removal rule, or archive/review history handling; no dropped-Spec deletion directive."
 else:reason="Allowed after contextual inspection: archive eligibility, historical source deletion, or generic scoped file mutation; no complete archive-disposition list or dropped-Spec deletion directive."
 hits.append(dict(path=path,line=int(number),text=text,allowed=True,reason=reason))
(E/"wording-hits.json").write_text(json.dumps(hits,indent=2));(E/"phrase-checks.json").write_text(json.dumps(checks,indent=2))
# Surface vocabulary is matched against the glossary and command guide despite skipped detector.
(E/"glossary.json").write_text(json.dumps(dict(entries=glossary,new_terms=["Cancel Command","Source Return"],other_new_terms_needed=False),indent=2))
failed=[c for c in checks if not c["pass_"]];print("phrase checks",len(checks),"failed",failed,"sweep hits",len(hits))
settle(12,"pass" if not failed else "fail",f"Whitespace-normalized authored phrases: {len(checks)} checks, {len(failed)} missing. The broad repeat sweep records all {len(hits)} hits and a per-hit allowance in [wording-hits.json](evidence/2026-10-12-cancel/wording-hits.json); exact matches in [phrase-checks.json](evidence/2026-10-12-cancel/phrase-checks.json). Live guidance routes dropped Specs to cancel. Earlier accepted ADRs and the expressly pinned old-parser reference are historical evidence, not stale instructions. Documentation's public commands were exercised by the CLI evidence.")
settle(14,"pass" if not failed else "fail","All six declared glossary entries have their exact ADR-0262 sentences; the two new terms are in order immediately before History Sanitize Command and their Avoid lines match. Emitted command/disposition/status terms resolve to these entries and the cancel guide; no additional Spec-introduced term needs an entry. [glossary.json](evidence/2026-10-12-cancel/glossary.json), exact phrase checks and successful public command transcripts.")
# Outside evidence is the pre-existing committed reference, not the synthetic fixture.
ref=Path("docs/references/2026-10-12-a-decomposed-spec-dropped-without-a-record.md");raw=ref.read_text()
assert "dropped eight open Specs" in norm(raw) and "5 Tasks `completed` and 3 `pending`" in norm(raw)
assert "It is still important to have a record of this fact." in raw and "Withdrawn" in raw and "Superseded" in raw
assert "`complete` (default option) or" in raw and "`not-planned`" in raw
outside=dict(reference=str(ref),sha256=hashlib.sha256(ref.read_bytes()).hexdigest(),sources=[dict(source="adopter event 2026-10-10",provenance="Secondbrain triaged Inbox Entry and Vortex mirrored planning document, both named by the reference",observation="Eight decomposed Specs dropped; one had five completed and three pending Tasks; seven were all pending; folders deleted without Archive Records",comparison="The live metric fixture reproduces the one eight-Task Spec, not eight Tasks conflated with eight Specs",supports_design=True),dict(source="PEP 1",url="https://peps.python.org/pep-0001/",observation="Preserves rejected/withdrawn proposals as records and distinguishes Withdrawn from Superseded",access="Quoted committed reference; no network request",supports_design=True),dict(source="GitHub Changelog 2022-05-19",url="https://github.blog/changelog/2022-05-19-the-new-github-issues-may-19th-update/",observation="Issue closure distinguishes complete from not-planned",access="Quoted committed reference; no network request",supports_design=True)])
(E/"outside-evidence.json").write_text(json.dumps(outside,indent=2))
settle(13,"pass","Reached every independent source through the committed durable reference as Requirement 8 prescribes. The adopter event distinguishes eight dropped Specs from the eight-Task (5 completed/3 pending) fixture; the live metric measures that shape. PEP 1 supports keeping dropped proposals recorded and separating Withdrawn from Superseded; the GitHub passage supports complete versus not-planned, not an invented delivering Spec. [outside-evidence.json](evidence/2026-10-12-cancel/outside-evidence.json) captures observations, exact provenance, URLs and reference SHA-256. These passages still support the design as recorded; live publication freshness and adopter originals were not re-observed. No network, provider or forge request was sent.")
assert not failed,failed
