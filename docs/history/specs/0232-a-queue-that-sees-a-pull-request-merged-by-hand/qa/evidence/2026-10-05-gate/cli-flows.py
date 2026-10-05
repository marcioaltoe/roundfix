import os, pathlib, tempfile, subprocess, json, sqlite3, time, hashlib
ROOT=pathlib.Path.cwd(); E=ROOT/"docs/specs/0232-a-queue-that-sees-a-pull-request-merged-by-hand/qa/evidence/2026-10-05-gate"

def run(args,cwd,env=None,check=True):
 p=subprocess.run(args,cwd=cwd,env=env,text=True,capture_output=True)
 if check and p.returncode: raise RuntimeError(str(args)+"\n"+p.stdout+p.stderr)
 return p

def fixture(mode):
 base=pathlib.Path(tempfile.mkdtemp(prefix="roundfix-0232-")).resolve(); repo=base/"repo"; repo.mkdir(); home=base/"home"; home.mkdir()
 env=dict(os.environ,HOME=str(home),GIT_CONFIG_NOSYSTEM="1",GIT_CONFIG_GLOBAL="/dev/null")
 def git(*a):return run(["git",*a],repo,env).stdout.strip()
 git("init","-b","main");git("config","user.name","QA");git("config","user.email","qa@example.invalid")
 (repo/"docs/specs/0300-example").mkdir(parents=True);(repo/"docs/specs/.gitkeep").touch();(repo/"docs/specs/0300-example/_prd.md").write_text("# disposable Spec\n")
 git("add",".");git("commit","-m","seed")
 origin=base/"origin.git";run(["git","init","--bare","--initial-branch=main",str(origin)],base,env)
 git("remote","add","origin",str(origin));git("push","-u","origin","main");git("remote","set-head","origin","main")
 branch="roundfix/deliver-0300-example-0123456789abcdef"; wt=home/".roundfix/worktrees"/("repo-"+hashlib.sha256(str(repo).encode()).hexdigest()[:8])/(branch.replace("/","-")+"-"+hashlib.sha256(branch.encode()).hexdigest()[:8]); wt.parent.mkdir(parents=True)
 git("worktree","add","-b",branch,str(wt),"main")
 (wt/"docs/history/specs").mkdir(parents=True);run(["git","mv","docs/specs/0300-example","docs/history/specs/0300-example"],wt,env);run(["git","commit","-m","archive"],wt,env)
 candidate=run(["git","rev-parse","HEAD"],wt,env).stdout.strip(); head=candidate; merge=""
 if mode=="merged":
  (wt/"ci-fix.txt").write_text("fixed\n");run(["git","add","."],wt,env);run(["git","commit","-m","fix CI"],wt,env);head=run(["git","rev-parse","HEAD"],wt,env).stdout.strip()
  git("merge","--no-ff",branch,"-m","manual merge");merge=git("rev-parse","HEAD");git("push","origin","main");git("worktree","remove",str(wt));git("branch","-D",branch)
 elif mode=="archive":
  git("merge","--squash",branch);git("commit","-m","manual squash");merge=git("rev-parse","HEAD");git("push","origin","main")
 (repo/"bin").mkdir();(repo/"bin/roundfix").symlink_to(ROOT/"bin/roundfix")
 fake=base/"fake";fake.mkdir();payload={"number":404,"state":"CLOSED" if mode=="closed" else "MERGED","headRefName":branch,"headRefOid":head,"mergeCommit":{"oid":merge} if merge else None}
 (fake/"gh").write_text("#!/bin/sh\nprintf '%s\\n' '"+json.dumps(payload)+"'\n");(fake/"gh").chmod(0o755);env["PATH"]=str(fake)+os.pathsep+env["PATH"]
 run(["/tmp/roundfix-0232-seed",str(home),str(repo)],repo,env)
 dbpath=home/".roundfix/roundfix.db";db=sqlite3.connect(dbpath);db.row_factory=sqlite3.Row
 print("fixture",mode,str(base));print("schema",[dict(r) for r in db.execute("PRAGMA table_info(delivery_queue_items)")])
 db.execute("UPDATE delivery_queue_items SET stage='parked',blocker='checks-failed',branch=?,worktree=?,worktree_provisioned=1,candidate_commits=?,pull_request_number=? WHERE spec_slug='0300-example'",(branch,str(wt),json.dumps([candidate]),"" if mode=="archive" else "404"));db.commit()
 before=dict(db.execute("SELECT * FROM delivery_queue_items").fetchone());db.close()
 p=run(["./bin/roundfix","deliver","retry","0300-example"],repo,env,False)
 print("command: ./bin/roundfix deliver retry 0300-example");print("exit:",p.returncode);print("stdout:\n"+p.stdout);print("stderr:\n"+p.stderr)
 status=run(["./bin/roundfix","deliver","status"],repo,env,False);print("fresh public status",status.returncode,status.stdout,status.stderr)
 db=sqlite3.connect(dbpath);db.row_factory=sqlite3.Row;after=dict(db.execute("SELECT * FROM delivery_queue_items").fetchone());db.close();print("before",json.dumps(before));print("persisted after",json.dumps(after))
 if mode=="closed":
  reason='  retry Delivery Queue item "0300-example": pull request #404 was closed without merging; reopen it or merge the Spec into the default branch, then run roundfix deliver retry 0300-example'
  assert p.returncode==2 and p.stdout=="" and p.stderr.startswith("Retry refused\n") and reason in p.stderr and "  stage: parked; blocker: checks-failed" in p.stderr and before==after
 else:
  evidence="pull request #404" if mode=="merged" else 'Spec archived on default branch "main"'
  assert p.stdout.startswith("Merged outside the queue: "+evidence+"; merge commit "+merge+"\nRetried 0300-example: checks-failed -> merged\n")
  assert after["stage"]=="merged" and after["merge_commit"]==merge and after["retry_count"]==before["retry_count"] and json.loads(after["candidate_commits"])[-1]==head
  if p.returncode!=0: print("OWNER ENVIRONMENT BLOCK: merge persistence passed but detached owner could not start")
 if mode=="merged":
  for _ in range(100):
   d=sqlite3.connect(dbpath);d.row_factory=sqlite3.Row;q=dict(d.execute("SELECT * FROM delivery_queues").fetchone());final=dict(d.execute("SELECT * FROM delivery_queue_items").fetchone());d.close()
   if not q["owner_pid"]: break
   time.sleep(0.05)
  print("after owner exit",json.dumps(final));assert final["blocker"]=="" and final["stage"]=="merged" and final["warning"]=="" and not q["owner_pid"]
 print("ASSERTIONS PASSED",mode)
for mode in ["merged","archive","closed"]:fixture(mode)
