from pathlib import Path
import re, subprocess,json
from report import settle
root=Path.cwd();ev=Path(__file__).resolve().parent
created={"./internal/spec":["internal/spec/cancel_test.go","internal/spec/archive_record_cancelled_test.go"],"./internal/cli":["internal/cli/cancel_test.go","internal/cli/cancelled_record_readers_test.go"],"./internal/speccheck":["internal/speccheck/cancelled_record_test.go"],"./internal/baseline":["internal/baseline/cancel_clause_test.go"]}
kept=["TestArchiveCommandLeavesTheArchiveRecord","TestSupersedeAcceptsAnArchiveRecord","TestPrerequisitesCountAnArchiveRecord","TestReviewDiffOmitsTheRemovedSpecFolder"]
results=[]
for pkg,files in created.items():
 names=[n for f in files for n in re.findall(r"^func (Test\w+)\(",Path(f).read_text(),re.M)]
 if pkg=="./internal/cli":names+=kept
 cmd=["rtk","proxy","go","test","-count=1","-json","-run","^("+"|".join(sorted(set(names)))+")$",pkg]
 run=subprocess.run(cmd,text=True,capture_output=True);(ev/(pkg.rsplit("/",1)[1]+"-tests.jsonl")).write_text(run.stdout);(ev/(pkg.rsplit("/",1)[1]+"-tests.stderr.txt")).write_text(run.stderr)
 passed={x.get("Test") for line in run.stdout.splitlines() if (x:=json.loads(line)).get("Action")=="pass"}
 missing=set(names)-passed;results.append(dict(package=pkg,command=cmd,exit=run.returncode,required=names,missing=sorted(missing)))
 print(pkg,"exit",run.returncode,"top-level",len(names),"missing",sorted(missing),flush=True)
for pkg,names,tags in [("./internal/docscontract",["TestTheSkillCoverageMapIsCurrent"],["-tags","docscontract"]),("./skills",["TestEveryOwnedSkillVersionIsRecorded"],[])]:
 cmd=["rtk","proxy","go","test","-count=1","-json"]+tags+["-run","^("+"|".join(names)+")$",pkg]
 run=subprocess.run(cmd,text=True,capture_output=True);(ev/(pkg.rsplit("/",1)[1]+"-tests.jsonl")).write_text(run.stdout)
 passed={x.get("Test") for line in run.stdout.splitlines() if (x:=json.loads(line)).get("Action")=="pass"};missing=set(names)-passed
 results.append(dict(package=pkg,command=cmd,exit=run.returncode,required=names,missing=sorted(missing)));print(pkg,run.returncode,sorted(missing),flush=True)
(ev/"test-summary.json").write_text(json.dumps(results,indent=2))
settle(2,"pass" if all(x["exit"]==0 and not x["missing"] for x in results) else "fail","Focused tests executed with -count=1 and exact named-pass checks: [test-summary.json](evidence/2026-10-12-cancel/test-summary.json), plus per-package JSON transcripts. Every created test and all four kept tests are required. Phrase checks are recorded below by the documentation audit, which must also pass for this row.")
