import os, subprocess, tempfile, shutil, json, re, hashlib
from pathlib import Path
root=Path.cwd(); spec="0191-claims-with-receipts-and-contracts-as-they-ship"
evidence=root/"docs/specs"/spec/"qa/evidence/2026-10-01-c5340e8e"
scratch=Path(tempfile.mkdtemp(prefix="roundfix-0191-qa-")); clone=scratch/"candidate"
def git(*args):
    p=subprocess.run(["rtk","proxy","git","-c","core.fsmonitor=false",*args],cwd=clone if clone.exists() else root,text=True,capture_output=True)
    if p.returncode: raise RuntimeError(p.stderr)
    return p.stdout
subprocess.run(["rtk","proxy","git","clone","--no-hardlinks","--quiet",str(root),str(clone)],check=True)
git("config","user.name","QA Fixture");git("config","user.email","qa@example.invalid");git("config","commit.gpgsign","false");git("config","core.hooksPath","/dev/null")
(clone/"bin").mkdir(exist_ok=True);shutil.copy2(root/"bin/roundfix",clone/"bin/roundfix")
# HOME is scoped to child CLI processes solely to honor the authored disposable-home requirement.
# The parent shell environment and real user home are unchanged.
qa_env=dict(os.environ);qa_env["HOME"]=str(scratch/"home");Path(qa_env["HOME"]).mkdir()
qa_env["GIT_CONFIG_GLOBAL"]="/dev/null";qa_env["GIT_CONFIG_NOSYSTEM"]="1"
(clone/".roundfixrc.yml").write_text("specs:\n  root: docs/specs\n")
def write(rel,txt):
    p=clone/rel;p.parent.mkdir(parents=True,exist_ok=True);p.write_text(txt)
def commit(msg):
    git("add","docs/specs/0200-example");git("commit","--quiet","-m",msg)
transcript_text=(root/"docs/specs"/spec/"_techspec.md").read_text().split("### Surface Transcripts\n",1)[1].split("### Skills, guides",1)[0]
blocks=re.findall(r"```transcript\n(.*?)\n   ```",transcript_text,re.S)
blocks=["\n".join(line[3:] if line.startswith("   ") else line for line in b.splitlines()) for b in blocks]
def matches(expected,actual):
    lines=expected.splitlines(); out=actual.splitlines()
    def walk(i,j):
        if i==len(lines):return j==len(out)
        if lines[i]=="...":return any(walk(i+1,k) for k in range(j,len(out)+1))
        if j==len(out):return False
        parts=re.split(r"(<[^>]+>)",lines[i]);pattern="".join(".+" if q.startswith("<") and q.endswith(">") else re.escape(q) for q in parts)
        return bool(re.fullmatch(pattern,out[j])) and walk(i+1,j+1)
    return walk(0,0)
def run(label,args,transcript=None):
    cmd=["./bin/roundfix",*args];p=subprocess.run(cmd,cwd=clone,env=qa_env,text=True,capture_output=True)
    result={"cwd":str(clone),"command":cmd,"exit":p.returncode,"stdout":p.stdout,"stderr":p.stderr,"candidate":git("rev-parse","HEAD").strip()}
    if transcript:
        b=blocks[transcript-1];expout=b.split("stdout:\n",1)[1].split("\nstderr:",1)[0];experr=b.split("stderr:\n",1)[1].rsplit("exit:",1)[0].rstrip("\n");exitcode=int(b.rsplit("exit:",1)[1])
        result["transcript_match"]={"stdout":matches(expout,p.stdout),"stderr":matches(experr,p.stderr),"exit":p.returncode==exitcode}
    (evidence/(label+".json")).write_text(json.dumps(result,indent=2)+"\n")
    print(label,"exit",p.returncode,"match",result.get("transcript_match",{}),flush=True)
    return result
run("identity",["--version"])
run("precondition",["spec","check",spec,"--strict"])
prd="docs/specs/0200-example/_prd.md";tech="docs/specs/0200-example/_techspec.md"
constraints="## Project Constraints\n\n"+"".join(f"- {name}: not applicable — isolated QA artifact; Source: `docs/agents/spec-routing.md`.\n" for name in ["Identifier strategy","Authentication and HTTP","Active ADR obligations","Tooling authority"])+"\n"
header="---\nspec: 0200-example\nstatus: active\n---\n# Example\n\n"
header+=constraints+"## Success Metrics\n\nNone. Isolated diagnostic reproduction.\n\n## Example claim\n\n"
write(prd,header+'ADR-0116: "reads the cited records"\n');commit("fixture: changed word")
run("transcript-1",["spec","check","0200-example","--stage","prd"],1)
run("transcript-1-confirm",["spec","check","0200-example","--stage","prd","--format","json"])
write(prd,header+"ADR-0116 requires the check to read the cited record.\n");commit("fixture: missing receipt")
run("transcript-2",["spec","check","0200-example","--stage","prd"],2)
run("transcript-2-strict",["spec","check","0200-example","--stage","prd","--strict"])
write(prd,header+'ADR-0116 requires the check to read the cited record. ADR-0116: "reads the cited record"\n')
run("transcript-2-repaired",["spec","check","0200-example","--stage","prd","--strict"])
run("transcript-2-repaired-confirm",["spec","check","0200-example","--stage","prd","--strict","--format","json"])
write(prd,header)
block="```transcript\n$ never-execute-this-command\nstdout:\nhello\nstderr:\nexit: 0\n```\n"
write(tech,"# TechSpec\n## Surface Transcripts\n\n1. Surface Transcript: example\n\n"+block.replace("exit: 0\n","")+"\n2. Surface Transcript: second\n\n"+block)
write("docs/specs/0200-example/_tasks.md","---\nschema: spec-tasks/v1\nspec: 0200-example\nqa: task_02\ngraph:\n  nodes:\n    - id: task_01\n      file: task_01.md\n      needs: []\n    - id: task_02\n      file: task_02.md\n      needs: [task_01]\n---\n# Tasks\n")
for id,kind,req,refs in [("task_01","backend","Implement the example.","Surface Transcripts 1-2"),("task_02","qa","MUST reproduce Surface Transcript 1.","")]:
    write("docs/specs/0200-example/"+id+".md",f"---\ntask: {id}\nspec: 0200-example\nstatus: pending\ntype: {kind}\ncomplexity: low\n---\n# Task\n\n## Requirements\n\n1. {req}\n\n## Verification\n\n- `true`\n\n## References\n\n{refs}\n")
commit("fixture: malformed and ungated transcripts")
run("transcript-3",["spec","check","0200-example"],3)
run("transcript-3-confirm",["spec","check","0200-example","--format","json"])
# Transcript 4 exercises absence of a declaration; retaining the malformed block would still raise an error in an unheld Spec.
write(tech,"# TechSpec\n"+constraints+"## API Contracts\n\nNone. No application API in this diagnostic fixture.\n");guide=clone/".agents/skills/write-techspec/references/concrete-contracts.md";guide.unlink()
run("transcript-4",["spec","check","0200-example","--stage","techspec"],4)
run("transcript-4-confirm",["spec","check","0200-example","--stage","techspec","--format","json"])
shutil.copy2(root/".agents/skills/write-techspec/references/concrete-contracts.md",guide)
for f in ["_prd.md","_techspec.md"]:
    txt=(root/"docs/specs"/spec/f).read_text().replace(spec,"0201-own-rules")
    write("docs/specs/0201-own-rules/"+f,txt)
git("add","docs/specs/0201-own-rules");git("commit","--quiet","-m","fixture: own rules held after guide")
run("own-rules",["spec","check","0201-own-rules","--stage","techspec","--strict"])
run("own-rules-confirm",["spec","check","0201-own-rules","--stage","techspec","--strict","--format","json"])
old="0182-delivery-that-reviews-and-retries-from-where-the-item-stands"
for f in ["_prd.md","_techspec.md"]:
    shutil.copy2(root/"docs/history/specs"/old/f,clone/"docs/specs/0200-example"/f)
commit("fixture: archived Spec 0182 replay")
run("outside-replay",["spec","check","0200-example","--stage","techspec"])
run("outside-replay-confirm",["spec","check","0200-example","--stage","techspec","--format","json"])
(evidence/"scratch.json").write_text(json.dumps({"scratch":str(scratch),"clone":str(clone),"original_head":subprocess.check_output(["git","rev-parse","HEAD"],text=True).strip(),"binary_sha256":hashlib.sha256((root/"bin/roundfix").read_bytes()).hexdigest()},indent=2)+"\n")
