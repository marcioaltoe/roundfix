from pathlib import Path
import re,json,subprocess,tempfile,shutil,hashlib,os,tarfile,io
from report import settle
R=Path.cwd().resolve();E=Path(__file__).resolve().parent;S=R/"docs/specs/0261-a-cancelled-spec-leaves-a-record";BASE="7cee52134ce36025e9cf66a3bacdbde42aed2698";HEAD="b412a32d9e4e2f2c300c619cf088358802b9e9ab"
def git(*args):
 x=subprocess.run(["rtk","proxy","git",*args],capture_output=True,text=True);assert x.returncode==0,x.stderr;return x.stdout
commits={"task_02":"40e65f388772b97aae05d86216e6c3bf9fdef511","task_01":"a3ef66e1f372721ea70d5c3d220c439479f5df6f","task_04":"ec26061f0f15b144d86613741188e753b7c71ee3","task_03":"b412a32d9e4e2f2c300c619cf088358802b9e9ab"};rows=[];errors=[];allpaths=[]
for task,c in commits.items():
 paths=git("diff-tree","--no-commit-id","--name-only","-r",c).splitlines();allpaths+=paths;text=(S/(task+".md")).read_text();declared=set(re.findall(r"^- (?:instruction|interface|creates): `([^`]+)`",text,re.M));recorded=text.split("## Recorded paths",1)[1] if "## Recorded paths" in text else "";declared.update(re.findall(r"`([^`]+)`",recorded));declared.add("docs/specs/0261-a-cancelled-spec-leaves-a-record/"+task+".md")
 outside=[p for p in paths if p not in declared];rows.append(dict(task=task,commit=c,command="git diff-tree --no-commit-id --name-only -r "+c,actual_paths=paths,declared_paths=sorted(declared),outside=outside))
 if outside:errors.append(dict(task=task,outside=outside))
T=Path(tempfile.mkdtemp(prefix="roundfix-scope-qa-",dir="/private/tmp")).resolve();helper=T/"helper";helper.mkdir();shutil.copy(E/"scope-helper.go.txt",helper/"main.go");(helper/"go.mod").write_text("module roundfix/qa-fixture\n\ngo 1.26\n\nrequire roundfix v0.0.0\nreplace roundfix => "+str(R)+"\n");shutil.copy(R/"go.sum",helper/"go.sum")
x=subprocess.run(["rtk","proxy","go","build","-mod=mod","-o",str(helper/"scope"),"."],cwd=helper,text=True,capture_output=True);assert x.returncode==0,x.stderr
x=subprocess.run([str(helper/"scope"),*sorted(set(allpaths))],text=True,capture_output=True);assert x.returncode==0,x.stderr;governed=json.loads(x.stdout)
authority=(S/"_authorization.md").read_text();front=authority.split("---",2)[1];bounds=set(re.findall(r"^  - ([^\n]+)$",front.split("paths:",1)[1].split("operations:",1)[0],re.M));unbounded=[p for p,v in governed.items() if v and p not in bounds]
if unbounded:errors.append(dict(unbounded=unbounded))
changed=git("diff","--name-only",BASE,HEAD).splitlines();protected=[p for p in changed if p in ["go.mod","go.sum",".roundfixrc.yml","Makefile"] or p.startswith((".github/workflows/","docs/history/specs/","internal/baseline/assets/retention/")) or (p.startswith("internal/baseline/assets/modules/") and p!="internal/baseline/assets/modules/spec-workflow.json")]
if protected:errors.append(dict(protected=protected))
# The archive-root and source-baseline assets have no additional changes.
source_assets=[p for p in changed if re.search(r"(?:source-baseline|source_baseline|retention)/",p)]
if source_assets:errors.append(dict(source_baseline=source_assets))
# Every earlier clause byte is retained, and QA settlement is unchanged both across skills and from base.
before=json.loads(git("show",BASE+":internal/baseline/assets/modules/spec-workflow.json"));after=json.loads(Path("internal/baseline/assets/modules/spec-workflow.json").read_text())
def clauses(x):
 if isinstance(x,dict):
  if x.get("id")=="clause.spec.keep-artifacts-in-spec-folder":yield x
  for v in x.values():yield from clauses(v)
 elif isinstance(x,list):
  for v in x:yield from clauses(v)
b=list(clauses(before))[0];a=list(clauses(after))[0];assert a["guidance"].startswith(b["guidance"]);assert {k:v for k,v in b.items() if k!="guidance"}=={k:v for k,v in a.items() if k!="guidance"}
sections=[]
for directory in [".agents/skills","skills"]:
 for p in Path(directory).glob("*/SKILL.md"):
  text=p.read_text()
  if "### QA settlement\n" not in text:continue
  text=text.split("### QA settlement\n",1)[1];text=re.split(r"(?m)^#{1,3} ",text,maxsplit=1)[0]
  old=git("show",BASE+":"+str(p)).split("### QA settlement\n",1)[1];old=re.split(r"(?m)^#{1,3} ",old,maxsplit=1)[0]
  sections.append(dict(path=str(p),sha256=hashlib.sha256(text.encode()).hexdigest(),unchanged_from_base=text==old))
assert len({s["sha256"] for s in sections})==1 and all(s["unchanged_from_base"] for s in sections),sections
mirrors=[]
for p in ["archive-spec/SKILL.md","roundfix/SKILL.md","roundfix/references/archive.md","roundfix/references/deliver.md"]:
 equal=Path(".agents/skills",p).read_bytes()==Path("skills",p).read_bytes();mirrors.append(dict(path=p,byte_identical=equal));assert equal,p
# A scratch clone uses the same audited head, no upstream and isolated HOME.
clone=T/"repo";home=T/"home";home.mkdir()
x=subprocess.run(["rtk","proxy","git","clone","--shared","--no-hardlinks","--no-checkout",str(R),str(clone)],capture_output=True,text=True);assert x.returncode==0,x.stderr
x=subprocess.run(["rtk","proxy","git","checkout","--detach",HEAD],cwd=clone,capture_output=True,text=True);assert x.returncode==0,x.stderr
(clone/"bin").symlink_to(R/"bin",target_is_directory=True)
env={"HOME":str(home),"XDG_CONFIG_HOME":str(home/".config"),"PATH":os.environ["PATH"],"GIT_CONFIG_NOSYSTEM":"1","GIT_CONFIG_GLOBAL":"/dev/null","NO_COLOR":"1"}
plans=[]
for i in [1,2]:
 cmd=["./bin/roundfix","baseline","update","--repo",".","--no-skills","--format","json"];x=subprocess.run(cmd,cwd=clone,env=env,capture_output=True,text=True,timeout=45);(E/f"baseline-update-{i}.json").write_text(x.stdout);(E/f"baseline-update-{i}.stderr.txt").write_text(x.stderr);obj=json.loads(x.stdout);plans.append(dict(command=cmd,exit=x.returncode,state=obj.get("state"),cwd=str(clone)))
 if x.returncode!=0 or obj.get("state")!="current":errors.append(dict(baseline=plans[-1]))
# Authorization predates the change: the operative file exists at BASE, with no Task edit.
assert git("show",BASE+":docs/specs/0261-a-cancelled-spec-leaves-a-record/_authorization.md")==authority
for c in commits.values():assert "docs/specs/0261-a-cancelled-spec-leaves-a-record/_authorization.md" not in git("diff-tree","--no-commit-id","--name-only","-r",c).splitlines()
result=dict(commit_range=BASE+".."+HEAD,tasks=rows,governed=governed,bounds=sorted(bounds),unbounded=unbounded,protected_changes=protected,source_baseline_changes=source_assets,settlement_sections=sections,mirrors=mirrors,baseline_plans=plans,errors=errors)
(E/"scope-audit.json").write_text(json.dumps(result,indent=2));print("scope audit",len(rows),"Task commits",len(set(allpaths)),"paths",sum(governed.values()),"governed; errors",errors,flush=True)
settle(15,"pass" if not errors else "fail","Audited all four actual Task commits with diff-tree over "+BASE+".."+HEAD+". Every path is declared or Recorded paths; all 11 Governed Paths are within the operative pre-Task grant. Protected manifests, Project Config, Makefile, workflows, all other modules, Source Baseline assets and existing Spec records are unchanged. The clause preserves every earlier byte. Both isolated baseline update JSON runs return current. All six QA settlement sections are byte-identical and unchanged from base, and all four edited canonical/mirror pairs match. [scope-audit.json](evidence/2026-10-12-cancel/scope-audit.json), baseline-update-1/2.json. The seeded mechanical authorization audit already names task_02 at the same commit and grant; no duplicate Mechanical Refusal Code row was added.")
assert not errors,errors
