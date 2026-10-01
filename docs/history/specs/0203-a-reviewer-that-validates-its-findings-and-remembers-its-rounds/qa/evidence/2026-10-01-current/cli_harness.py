import os, subprocess, tempfile, json, pathlib, re
ROOT=pathlib.Path.cwd(); E=pathlib.Path(__file__).resolve().parent
scratch=pathlib.Path(tempfile.mkdtemp(prefix="qa0203-",dir="/private/tmp")); fake=scratch/"fake";fake.mkdir()
fakecode = r'''#!/usr/bin/env python3
import sys,os,json,pathlib
args=sys.argv[1:]; state=pathlib.Path(os.environ["QA_FAKE_STATE"])
with (state/"calls.jsonl").open("a") as f:f.write(json.dumps(args)+"\n")
if "--version" in args:
 print("0.19.3" if pathlib.Path(sys.argv[0]).name=="acpx" else "@agentclientprotocol/codex-acp 2.0.1");sys.exit()
model=args[args.index("--model")+1] if "--model" in args else "gpt-5.6-sol"
mf=state/"models.json"; models=json.loads(mf.read_text()) if mf.exists() else {}
if "ensure" in args:
 name=args[args.index("--name")+1];models[name]=model;mf.write_text(json.dumps(models))
if "show" in args:model=models.get(args[-1],model)
if "-s" in args:model=models.get(args[args.index("-s")+1],model)
opts=[{"id":"model","category":"model","type":"select","currentValue":model,"options":[{"value":"gpt-5.6-sol"},{"value":"gpt-5.5"}]},{"id":"reasoning_effort","type":"select","currentValue":"high","options":[{"value":"high"}]}]
if "prompt" in args:
 text=sys.stdin.read(); (state/("prompt-"+str(len(list(state.glob("prompt-*"))))+".txt")).write_text(text)
 answer=(state/"answer.txt").read_text()
 if "--deny-all" in args:answer=json.dumps({"schema":"roundfix/review-validation/v1","findings":[{"id":"F1","verdict":"dismiss","rule":"convention:C2","reason":"restates the Daemon settlement"}]})
 print(json.dumps({"jsonrpc":"2.0","method":"session/update","params":{"sessionId":"fake-acp-session","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":answer}}}}))
 print(json.dumps({"jsonrpc":"2.0","id":1,"result":{"stopReason":"end_turn"}}));sys.exit()
if "set" in args:
 i=args.index("set");print(json.dumps({"action":"config_set","configId":args[i+1],"value":args[i+2],"configOptions":opts}));sys.exit()
print(json.dumps({"schema":"acpx.session.v1","acpx":{"current_model_id":model,"config_options":opts},"configOptions":opts}))
'''
for name in ["acpx","npx"]:
 f=fake/name;f.write_text(fakecode);f.chmod(0o755)
env=dict(os.environ,PATH=str(fake)+":"+os.environ["PATH"])
def fixture(name, task=False):
 r=scratch/name;r.mkdir();a=scratch/(name+"-artifacts");a.mkdir();state=scratch/(name+"-state");state.mkdir()
 def git(*args):return subprocess.check_output(["git","-c","core.fsmonitor=false",*args],cwd=r,text=True,stderr=subprocess.DEVNULL).strip()
 git("init","-b","main");git("config","commit.gpgsign","false");git("config","core.hooksPath","/dev/null");git("config","user.name","QA");git("config","user.email","qa@example.invalid")
 (r/"review.txt").write_text("before\n");git("add",".");git("commit","-m","test: base");base=git("rev-parse","HEAD")
 (r/".roundfixrc.yml").write_text(f"defaults:\n  artifact_dir: {json.dumps(str(a))}\npre_pr_review:\n  provider: codex\nprofiles:\n  review:\n    preferred:\n      runtime: codex\n      model: gpt-5.6-sol\n      reasoning_effort: high\n    fallbacks:\n      - runtime: codex\n        model: gpt-5.5\n        reasoning_effort: high\n")
 (r/"review.txt").write_text("before\ncandidate\n")
 if task:
  t=r/"docs/specs/0001-example/task_01.md";t.parent.mkdir(parents=True);t.write_text("---\ntask: task_01\nspec: 0001-example\nstatus: completed\n---\n\n## Requirements\nDo the work.\n\n## Result\nStatus is Daemon-owned.\n\n## Next\nOther work.\n")
 git("add",".");git("commit","-m","test: candidate")
 return r,a,state,git,base
observations={}
def run(n,fx,answer,args=None):
 r,a,s,g,b=fx;(s/"answer.txt").write_text(answer);en=dict(env,QA_FAKE_STATE=str(s))
 cmd=[str(ROOT/"bin/roundfix"),"review"]+(args if args is not None else ["--base",b])
 out=subprocess.run(cmd,cwd=r,env=en,text=True,capture_output=True)
 (E/f"transcript-{n}.stdout").write_text(out.stdout);(E/f"transcript-{n}.stderr").write_text(out.stderr)
 observations[str(n)]={"command":"./bin/roundfix "+" ".join(cmd[1:]),"cwd":str(r),"exit":out.returncode,"stdout":out.stdout,"stderr":out.stderr}
 return out
f=fixture("c2",True); run(1,f,"Findings:\n- `docs/specs/0001-example/task_01.md:4` — The Task is `completed` while its Result says status is Daemon-owned. Failure: an invalid Task state transition.")
run(7,f,"",["dispose","F1","--dismiss","--evidence","restates the Daemon settlement"])
f=fixture("mixed");run(2,f,"Findings:\n- internal/untouched.go:10 Failure: outside.\n- review.txt:2 Failure: candidate.")
f=fixture("no-anchor");run(3,f,"Findings:\n- A finding with no anchor.")
f=fixture("continued");run("4-setup",f,"Findings:\n- review.txt:2 Failure: candidate.");r,a,s,g,b=f
(r/"review.txt").write_text("before\nfixed\n");g("add",".");g("commit","-m","test: fix");fix=g("rev-parse","HEAD")
run("4-dispose",f,"",["dispose","F1","--fixed-by",fix]);run(4,f,"No findings")
f=fixture("ceiling");run("5-setup1",f,"Findings:\n- review.txt:2 Failure: candidate.");r,a,s,g,b=f
(r/"review.txt").write_text("before\nround two\n");g("add",".");g("commit","-m","test: round two");run("5-setup2",f,"Findings:\n- review.txt:2 Failure: still broken.")
record=list(a.glob("pre-pr-review/*/pre-pr-review.json"))[0];old=record.read_bytes();calls=(s/"calls.jsonl").read_text().count('"prompt"')
(r/"review.txt").write_text("before\nround three\n");g("add",".");g("commit","-m","test: third candidate");fix=g("rev-parse","HEAD")
run(5,f,"No findings");assert old==record.read_bytes();assert calls==(s/"calls.jsonl").read_text().count('"prompt"')
run("6-dispose",f,"",["dispose","F1","--fixed-by",fix]);run(6,f,"No findings");assert calls==(s/"calls.jsonl").read_text().count('"prompt"')
# Exact transcript matching, with only the TechSpec's placeholders variable.
spec=(ROOT/"docs/specs/0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds/_techspec.md").read_text()
blocks=re.findall(r"```transcript\n(.*?)```",spec,re.S)
for i,block in enumerate(blocks,1):
 expected=block.split("stdout:\n",1)[1];stdout,rest=expected.split("stderr:\n",1);stderr,exitcode=rest.rsplit("exit: ",1)
 def pattern(t):
  t="\n".join(line[3:] if line.startswith("   ") else line for line in t.splitlines()).strip("\n")
  parts=re.split(r"(<[^>]+>)",t);return "".join("[^\n]+" if x.startswith("<") and x.endswith(">") else re.escape(x) for x in parts)
 o=observations[str(i)];o["matches"]=o["exit"]==int(exitcode.strip()) and bool(re.fullmatch(pattern(stdout),o["stdout"].rstrip("\n"))) and bool(re.fullmatch(pattern(stderr),o["stderr"].rstrip("\n")))
 (E/"transcripts.json").write_text(json.dumps(observations,indent=2))
 print(i,"PASS" if o["matches"] else "FAIL",o["exit"],o["stderr"].strip())
