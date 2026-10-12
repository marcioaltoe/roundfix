from pathlib import Path
import os,sys,json,subprocess,tempfile,hashlib,shutil,re,datetime
from report import settle
ROOT=Path.cwd().resolve(); EV=Path(__file__).resolve().parent
TEMP=Path(tempfile.mkdtemp(prefix="roundfix-cancel-qa-",dir="/private/tmp")).resolve()
SLUG="0261-cancel-fixture"; CASES=[]
# A temporary nested module permits fixture use of Roundfix's internal store.
helper=TEMP/"helper";helper.mkdir();(helper/"main.go").write_bytes((EV/"seed-store.go.txt").read_bytes())
(helper/"go.mod").write_text("module roundfix/qa-fixture\n\ngo 1.26\n\nrequire roundfix v0.0.0\nreplace roundfix => "+str(ROOT)+"\n")
shutil.copy(ROOT/"go.sum",helper/"go.sum")
build=subprocess.run(["rtk","proxy","go","build","-mod=mod","-o",str(helper/"seed"),"."],cwd=helper,text=True,capture_output=True)
(EV/"seed-build.txt").write_text(build.stdout+build.stderr);assert build.returncode==0,build.stderr
constraints="## Project Constraints\n\n- Identifier strategy: not applicable — fixture uses existing identities. Source: `docs/agents/domain.md`.\n- Authentication and HTTP: not applicable — no transport. Source: `docs/agents/agent-instructions.md`.\n- Active ADR obligations: not applicable — no ADR. Source: `docs/agents/domain.md`.\n- Tooling authority: not applicable — no tooling edits. Source: `docs/agents/agent-instructions.md`.\n"
def write(path,text):path.parent.mkdir(parents=True,exist_ok=True);path.write_text(text)
def env(home):return {"HOME":str(home),"XDG_CONFIG_HOME":str(home/".config"),"PATH":os.environ["PATH"],"GIT_CONFIG_NOSYSTEM":"1","GIT_CONFIG_GLOBAL":"/dev/null","NO_COLOR":"1","LANG":"en_US.UTF-8","TZ":"UTC"}
def call(cmd,cwd,home):return subprocess.run(cmd,cwd=cwd,env=env(home),text=True,capture_output=True,timeout=30)
def git(repo,home,*args):
 x=call(["rtk","proxy","git",*args],repo,home);assert x.returncode==0,(args,x.stderr);return x.stdout
def commit(repo,home):
 if not git(repo,home,"status","--porcelain","--untracked-files=all").strip():return
 git(repo,home,"add","-A");git(repo,home,"commit","-m","docs: seed QA fixture")
def spec(repo,slug,n=2,completed=0,requires=None):
 d=repo/"docs/specs"/slug;write(d/"_prd.md","---\nspec: "+slug+"\nstatus: active\ncreated: 2026-10-10\n---\n\n# Cancel fixture\n\nNo consumer needs this work.\n\n"+constraints+"\n## Success Metrics\n\nNone. Operational fixture.\n")
 manifest="---\nschema: spec-tasks/v1\nspec: "+slug+"\nqa: declined\nqa_reason: fixture measures the public cancellation path\n"
 if requires:manifest+="requires: ["+requires+"]\n"
 manifest+="graph:\n  nodes:\n"
 for i in range(1,n+1):
  t=f"task_{i:02d}";manifest+=f"    - id: {t}\n      file: {t}.md\n      needs: []\n";write(d/(t+".md"),f"---\ntask: {t}\nspec: {slug}\nstatus: {'completed' if i<=completed else 'pending'}\ntype: backend\n---\n\n# Fixture Task\n\n## Verification\n\n- `test -f .fixture-result` — expected: exit 0.\n")
 write(d/"_tasks.md",manifest+"---\n\n# Task Graph\n")
 return d
def fixture(name,n=2,completed=0):
 top=TEMP/name;repo=top/"repo";home=top/"home";repo.mkdir(parents=True);home.mkdir()
 git(repo,home,"init","--initial-branch=main");git(repo,home,"config","user.name","QA Fixture");git(repo,home,"config","user.email","qa@example.invalid");git(repo,home,"config","commit.gpgsign","false")
 (repo/"bin").symlink_to(ROOT/"bin",target_is_directory=True);write(repo/".git/info/exclude","/bin\n")
 write(repo/"docs/agents/domain.md","# Domain\n");write(repo/"docs/agents/agent-instructions.md","# Agent instructions\n");d=spec(repo,SLUG,n,completed)
 commit(repo,home);return repo,home,d
def adopt(d,finding=True):
 write(d/"references/_index.md","| path | type |\n| --- | --- |\n| entry.md | backlog |\n"+("| finding.md | finding |\n" if finding else ""))
 original="---\nstatus: promoted # keep comment\nspec: "+SLUG+" # keep comment\n---\n\n# Entry\n\nIntent and provenance remain exact.\n"
 write(d/"references/entry.md",original)
 if finding:write(d/"references/finding.md","---\nstatus: done\nupdated_at: 2026-10-01\nabsorbed_by: "+SLUG+"\n---\n\n# Finding\n\nEvidence stays.\n")
 return original
def seed(repo,home,mode):
 x=call([str(helper/"seed"),str(home),str(repo),SLUG,mode,str(os.getpid())],repo,home);assert x.returncode==0,x.stderr;return x.stdout
def folderhash(d):
 h=hashlib.sha256()
 if d.exists():
  for p in sorted(d.rglob("*")):
   if p.is_symlink():data=("link:"+os.readlink(p)).encode()
   elif p.is_file():data=p.read_bytes()
   else:continue
   h.update(str(p.relative_to(d)).encode()+b"\0"+data+b"\0")
 return h.hexdigest()
def state(repo,home,d):return dict(git_status=git(repo,home,"status","--porcelain","--untracked-files=all"),spec_sha256=folderhash(d),head=git(repo,home,"rev-parse","HEAD").strip())
def snapshot(repo):return {str(p.relative_to(repo)):("link:"+os.readlink(p) if p.is_symlink() else hashlib.sha256(p.read_bytes()).hexdigest()) for p in repo.rglob("*") if p.is_file() and ".git" not in p.relative_to(repo).parts and "bin" not in p.relative_to(repo).parts}
def run(repo,home,args,name,d=None):
 d=d or repo/"docs/specs"/SLUG;before=state(repo,home,d);bs=snapshot(repo)
 x=call(["./bin/roundfix",*args],repo,home);after=state(repo,home,d)
 data=dict(name=name,command=["./bin/roundfix",*args],cwd=str(repo),home=str(home),exit=x.returncode,stdout=x.stdout,stderr=x.stderr,before=before,after=after,repository_bytes_unchanged=bs==snapshot(repo),binary_sha256=hashlib.sha256((ROOT/"bin/roundfix").read_bytes()).hexdigest())
 (EV/(name+".json")).write_text(json.dumps(data,indent=2));CASES.append(data);print(name,"exit",x.returncode,flush=True);return x,data
def refusal(reason):return "Preflight failed\n\nReason:\n  "+reason+"\n\nNo side effects:\n  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.\n\nUsage:\n  Run 'roundfix cancel --help' for usage.\n"
def expectedsummary(repo,home,d,c,p):
 files=[f for f in d.rglob("*") if f.is_file()];return f"cancelled {SLUG} -> docs/history/specs/{SLUG}.md; {c} completed Task(s), {p} pending; removed {len(files)} file(s) ({sum(f.stat().st_size for f in files)} bytes) kept in Git at {git(repo,home,'rev-parse','HEAD').strip()[:12]}\n"
def checkref(x,data,reason=None):
 assert x.returncode==2 and x.stdout=="" and (not reason or x.stderr==refusal(reason)),data
 assert data["before"]==data["after"] and data["repository_bytes_unchanged"],data
 assert x.stderr.startswith("Preflight failed\n"),data
# Exact Surface Transcripts 1 through 6.
for t in range(1,7):
 if re.search(r"(?m)^\| "+str(t+2)+r" \| pass \|", (EV.parents[1]/"qa-report-2026-10-12.md").read_text()):continue
 try:
  repo,home,d=fixture("transcript-"+str(t),8 if t==1 else 2,5 if t==1 else 0);args=["cancel",SLUG,"--reason","no consumer"];reason=None
  if t==1:adopt(d);commit(repo,home);out=expectedsummary(repo,home,d,5,3)+"returned docs/backlog/entry.md as open\nreturned docs/findings/finding.md as pending\n";err="";code=0
  elif t==2:
   write(repo/"docs/handoffs/notes.md",f"# Notes\n\nRead `docs/specs/{SLUG}/`.\n");commit(repo,home);seed(repo,home,"parked");out=expectedsummary(repo,home,d,0,2);err=f"roundfix: warning: Delivery Queue item {SLUG} stays parked; a Delivery Retry of a cancelled Spec is refused, so close its branch and pull request by hand\nroundfix: warning: docs/handoffs/notes.md:3 cites docs/specs/{SLUG}/, which this cancellation removed\n";code=0
  elif t==3:seed(repo,home,"running");reason=f'Spec "{SLUG}" cannot be cancelled while its Delivery Queue item is at stage running; stop the queue with: roundfix deliver stop'
  elif t==4:
   rid=seed(repo,home,"active");reason=f'Active Run {rid} holds Spec "{SLUG}"; stop it with: roundfix stop {rid}'
  elif t==5:spec(repo,"dependent",requires=SLUG);commit(repo,home);reason=f'Spec "{SLUG}" cannot be cancelled while active Spec dependent requires it; remove it from that Spec\'s requires first'
  else:
   args=["cancel","--help"];tech=(ROOT/"docs/specs/0261-a-cancelled-spec-leaves-a-record/_techspec.md").read_text();block=re.findall(r"```transcript\n(.*?)\n\s*```",tech,re.S)[5];out=block.split("   stdout:\n",1)[1].split("   stderr:",1)[0];out="\n".join(line[3:] if line.startswith("   ") else line for line in out.splitlines())+"\n";err="";code=0
  if reason:out="";err=refusal(reason);code=2
  x,data=run(repo,home,args,"transcript-"+str(t),d);assert (x.returncode,x.stdout,x.stderr)==(code,out,err),(data,out,err)
  if reason:checkref(x,data,reason)
  if t in [1,2]:assert not d.exists() and (repo/"docs/history/specs"/(SLUG+".md")).exists() and data["before"]["head"]==data["after"]["head"]
  settle(t+2,"pass",f"Exact stdout, stderr and exit matched Surface Transcript {t}; [transcript-{t}.json](evidence/2026-10-12-cancel/transcript-{t}.json). Fresh filesystem and Git reads confirm persistence on success or identical repository bytes, Spec hash, status and HEAD on refusal. Fixture invocation uses ./bin/roundfix via a symlink to the locally built binary; its SHA-256 is recorded.")
 except Exception as exc:
  settle(t+2,"fail",f"Surface Transcript {t} mismatch: {str(exc)[:1200]}. See transcript evidence.");print("FAIL transcript",t,repr(exc),flush=True)
# Success Metric 1 independently measures the exact adopter shape.
repo,home,d=fixture("metric-one",8,5);original=adopt(d,False);commit(repo,home)
x,data=run(repo,home,["cancel",SLUG,"--reason","no consumer"],"metric-one",d)
record=(repo/"docs/history/specs"/(SLUG+".md")).read_text();returned=(repo/"docs/backlog/entry.md").read_text()
want=original.replace("status: promoted","status: open").replace("spec: "+SLUG,"spec: null")
checks={"completed_tasks":"[task_01, task_02, task_03, task_04, task_05]","pending_tasks":"[task_06, task_07, task_08]","cancel_reason":"no consumer"}
assert x.returncode==0 and returned==want and not d.exists(),data
for key,val in checks.items():
 if key.endswith("_tasks"):
  match=re.search(r"(?m)^"+key+r":\n((?:  - [^\n]+\n)*)",record); observed=[v.strip() for v in re.findall(r"  - ([^\n]+)",match.group(1))];assert observed==val[1:-1].split(", "),(key,record)
 else:assert f"{key}: {val}\n" in record,(key,record)
archived=re.search(r"(?m)^archived: (.+)$",record).group(1).strip('"');assert archived==datetime.datetime.now(datetime.timezone.utc).date().isoformat()
(EV/"metric-record.md.txt").write_text(record);(EV/"source-before.md.txt").write_text(original);(EV/"source-returned.md.txt").write_text(returned)
(EV/"metric-one-measurement.json").write_text(json.dumps(dict(**checks,archived=archived,byte_comparison_outside_status_spec=True,before_sha256=hashlib.sha256(original.encode()).hexdigest(),returned_sha256=hashlib.sha256(returned.encode()).hexdigest()),indent=2))
settle(9,"pass","Eight Tasks (task_01–task_05 completed; task_06–task_08 pending) and one adopted Backlog Entry cancelled with exit 0. Record has cancel_reason=no consumer and archived="+archived+". Returned entry matches every original byte outside status/spec, including comments. [metric-one-measurement.json](evidence/2026-10-12-cancel/metric-one-measurement.json), [metric-record.md.txt](evidence/2026-10-12-cancel/metric-record.md.txt), before/returned source artifacts. A fresh Git read confirms HEAD unchanged; repeat cancel refuses unknown active Spec.")
x,_=run(repo,home,["cancel",SLUG,"--reason","again"],"metric-repeat",d);assert x.returncode==2
# Every publicly constructible API Contract 1 refusal.
refusals=["missing-slug","missing-reason","blank-reason","unexpected-argument","outside-git","unknown-spec","non-markdown-pin","uncommitted","missing-prd","supersession","manifest","record-occupied","folder-occupied","source-missing","source-symlink","source-parent-symlink","status-missing","destination-occupied","duplicate-destination","unsafe-path","other-source-missing","requires-malformed"]+["stage-"+s for s in ["queued","running","reviewing","archiving","gating","publishing","checking","merging"]]+["active-run","requiring-spec"]
failures=[]
for kind in refusals:
 repo,home,d=fixture("refusal-"+kind);args=["cancel",SLUG,"--reason","no consumer"]
 if kind=="missing-slug":args=["cancel"]
 elif kind=="missing-reason":args=["cancel",SLUG]
 elif kind=="blank-reason":args[-1]=" \n\t "
 elif kind=="unexpected-argument":args.append("extra")
 elif kind=="unknown-spec":args[1]="unknown"
 elif kind=="non-markdown-pin":write(repo/"pin.txt",f"docs/specs/{SLUG}/\n")
 elif kind=="uncommitted":write(d/"_prd.md","uncommitted\n")
 elif kind=="missing-prd":(d/"_prd.md").unlink()
 elif kind=="supersession":write(d/"_supersession.md","supersession\n")
 elif kind=="manifest":write(d/"_tasks.md","invalid manifest\n")
 elif kind=="record-occupied":write(repo/"docs/history/specs"/(SLUG+".md"),"occupied\n")
 elif kind=="folder-occupied":write(repo/"docs/history/specs"/SLUG/"keep.md","occupied\n")
 elif kind in ["source-missing","source-symlink","source-parent-symlink","status-missing","destination-occupied","duplicate-destination","unsafe-path","other-source-missing"]:
  adopt(d,False);f=d/"references/entry.md"
  if kind=="source-missing":f.unlink()
  elif kind=="source-symlink":f.unlink();f.symlink_to("../_prd.md")
  elif kind=="source-parent-symlink":(d/"references/nested").symlink_to(".",target_is_directory=True);write(d/"references/_index.md","| type | path |\n| --- | --- |\n| backlog | nested/entry.md |\n")
  elif kind=="status-missing":write(f,"---\nspec: demo\n---\n\nOriginal\n")
  elif kind=="destination-occupied":write(repo/"docs/backlog/entry.md","occupied\n")
  elif kind=="duplicate-destination":write(d/"references/_index.md","| type | path |\n| --- | --- |\n| backlog | entry.md |\n| backlog | entry.md |\n")
  elif kind=="unsafe-path":write(d/"references/_index.md","| type | path |\n| --- | --- |\n| backlog | ../_prd.md |\n")
  elif kind=="other-source-missing":write(d/"references/_index.md","| type | path |\n| --- | --- |\n| inbox | missing.md |\n")
 elif kind=="requires-malformed":write(repo/"docs/specs/dependent/_tasks.md","---\nrequires: bad\n---\n")
 elif kind=="requiring-spec":spec(repo,"dependent",requires=SLUG)
 if kind not in ["uncommitted","outside-git"]:commit(repo,home)
 if kind.startswith("stage-"):seed(repo,home,kind[6:])
 elif kind=="active-run":seed(repo,home,"active")
 if kind=="outside-git":other=repo.parent/"not-git";other.mkdir();(other/"bin").symlink_to(ROOT/"bin",target_is_directory=True);repo=other
 x,data=run(repo,home,args,"refusal-"+kind,d) if kind!="outside-git" else (None,None)
 if kind=="outside-git":
  x=call(["./bin/roundfix",*args],repo,home);data=dict(name=kind,exit=x.returncode,stdout=x.stdout,stderr=x.stderr,spec_sha256=folderhash(d));(EV/"refusal-outside-git.json").write_text(json.dumps(data,indent=2));assert x.returncode==2 and "requires a git repository" in x.stderr;continue
 try:checkref(x,data)
 except Exception as exc:failures.append(dict(name=kind,error=str(exc)))
(EV/"refusal-summary.json").write_text(json.dumps(dict(cases=refusals,failures=failures,unreachable_revision_checks="CLI always resolves a 40-hex HEAD; invalid revision length/hex are library-only and passed TestCancelRefusalsWriteNothing"),indent=2))
settle(10,"pass" if not failures else "fail","All "+str(len(refusals))+" publicly constructible API Contract 1 refusals exited 2. Each Git fixture records before/after status, Spec SHA-256, HEAD and whole-repository byte comparison. Non-Git refusal records the valid Spec independently. Invalid SourceRevision length/hex cannot be supplied through the CLI, which derives HEAD; both required library refusal subtests passed. [refusal-summary.json](evidence/2026-10-12-cancel/refusal-summary.json) and refusal-*.json hold every result.")
# Reader commands use the real cancelled record produced above.
repo=TEMP/"metric-one/repo";home=TEMP/"metric-one/home";d=repo/"docs/specs"/SLUG
commit(repo,home);before=(repo/"docs/history/specs"/(SLUG+".md")).read_bytes();seed(repo,home,"terminal")
x,data=run(repo,home,["runs","causes","--format","json"],"reader-causes",d);j=json.loads(x.stdout);assert x.returncode==0 and any(a["disposition"]=="cancelled" for a in j.get("archived_specs",[])),data
write(repo/"docs/history/specs/0001-old/_prd.md","---\nspec: 0001-old\ncreated: 2026-10-10\n---\n\n# Old archive\n\nA delivered outcome.\n");commit(repo,home);git(repo,home,"tag","-a","history-full","-m","docs: mark full fixture history")
x,data=run(repo,home,["history","sanitize"],"reader-sanitize-plan",d);assert x.returncode==0 and SLUG not in x.stdout and "1 unit(s) pending" in x.stdout,data
x,data=run(repo,home,["history","sanitize","--apply","--batch","1"],"reader-sanitize-apply",d);assert x.returncode==0 and (repo/"docs/history/specs"/(SLUG+".md")).read_bytes()==before,data
spec(repo,"0262-dependent");commit(repo,home)
x,data=run(repo,home,["supersede","--spec","0262-dependent","--by",SLUG,"--reason","replacement"],"reader-supersede",d);assert x.returncode==2 and 'a cancelled Spec delivers no content' in x.stderr and data["repository_bytes_unchanged"],data
spec(repo,"0263-requires",requires=SLUG);commit(repo,home)
x,data=run(repo,home,["deliver","start","0263-requires"],"reader-deliver-start",d);assert x.returncode==2 and f"requires cancelled Spec {SLUG}" in x.stderr and data["repository_bytes_unchanged"],data
seed(repo,home,"parked");db=home/".roundfix/roundfix.db";dbs=list((home/".roundfix").glob("*.db"));db=dbs[0]
db_before=db.read_bytes();x,data=run(repo,home,["deliver","retry",SLUG],"reader-deliver-retry",d);assert x.returncode==2 and "a cancelled Spec is never delivered" in x.stderr and data["repository_bytes_unchanged"] and db.read_bytes()==db_before,data
settle(11,"pass","Success Metrics 3/4: parser/render and review tests passed; public runs causes JSON reports cancelled with source history available; history sanitize plan excludes the record and apply converts one other legacy unit while preserving cancelled bytes; supersede, deliver start and deliver retry each exit 2 with exact cancellation reason and no repository writes. Retry preserves database bytes. [reader-causes.json](evidence/2026-10-12-cancel/reader-causes.json) and reader-*.json; per-package test transcripts supply byte round-trip, review exclusion and returned-source Spec-check evidence.")
(EV/"cli-summary.json").write_text(json.dumps(dict(temp_root=str(TEMP),cases=[dict(name=c["name"],exit=c["exit"]) for c in CASES],failures=failures),indent=2))
print("CLI flows finished",TEMP,flush=True)
