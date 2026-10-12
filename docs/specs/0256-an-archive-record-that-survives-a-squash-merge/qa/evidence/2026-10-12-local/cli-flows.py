import os,subprocess,tempfile,json,re,hashlib
from pathlib import Path
root=Path.cwd().resolve();slug='0300-qa-archive';ev=root/'docs/specs/0256-an-archive-record-that-survives-a-squash-merge/qa/evidence/2026-10-12-local'
base=Path(tempfile.mkdtemp(prefix='roundfix-0256-flows-',dir='/tmp'));env=os.environ.copy()
for k in list(env):
 if k.startswith('ROUNDFIX_') or k.startswith('GIT_'):env.pop(k,None)
env.update(HOME=str(base/'home'),XDG_CONFIG_HOME=str(base/'config'),GIT_CONFIG_NOSYSTEM='1',GIT_CONFIG_GLOBAL='/dev/null',GIT_TERMINAL_PROMPT='0',NO_COLOR='1',GOPROXY='off',GOSUMDB='off')
(base/'home').mkdir(); records=[]
def run(repo,args,ok=True):
 r=subprocess.run(args,cwd=repo,env=env,capture_output=True,text=True,timeout=60)
 if ok and r.returncode:raise RuntimeError((args,r.returncode,r.stdout,r.stderr))
 return r
def git(repo,*args,ok=True):return run(repo,['git',*args],ok).stdout.strip()
def write(p,t):p.parent.mkdir(parents=True,exist_ok=True);p.write_text(t)
def fixture(n):
 repo=base/('repo'+str(n));repo.mkdir();git(repo,'init','-b','main');git(repo,'config','user.name','QA');git(repo,'config','user.email','qa@example.invalid');git(repo,'config','commit.gpgsign','false')
 write(repo/'CONTEXT.md','# Context\n');write(repo/'AGENTS.md','# Instructions\n')
 for name in ['domain','agent-instructions']:write(repo/f'docs/agents/{name}.md','# Instructions\n')
 d=repo/'docs/specs'/slug
 write(d/'_prd.md',f'---\nspec: {slug}\nstatus: active\ncreated: 2026-10-12\n---\n\n# QA archive fixture\n')
 write(d/'_tasks.md','---\nschema: spec-tasks/v1\nqa: task_02\ngraph:\n  nodes:\n    - {id: task_01, file: task_01.md, needs: []}\n    - {id: task_02, file: task_02.md, needs: [task_01]}\n---\n')
 for ident,typ in [('task_01','backend'),('task_02','qa')]:write(d/(ident+'.md'),f'---\ntask: {ident}\nspec: {slug}\nstatus: completed\ntype: {typ}\ncomplexity: low\n---\n\n# {ident}\n\n## Overview\nHermetic fixture.\n\n## Requirements\n1. MUST exercise archive.\n\n## Acceptance Criteria\n- [x] Fixture complete.\n\n## Verification\n- `true`\n')
 git(repo,'add','.');git(repo,'commit','-m','docs: base fixture')
 forge=base/('forge'+str(n)+'.git');git(base,'init','--bare','--initial-branch=main',str(forge));git(repo,'remote','add','origin',str(forge));git(repo,'push','origin','main');git(repo,'symbolic-ref','refs/remotes/origin/HEAD','refs/remotes/origin/main');git(repo,'checkout','-b','fix/qa-fixture')
 write(d/'qa/qa-report-2026-10-12.md','---\nverdict: pass\n---\n\n# Fixture QA\n');git(repo,'add','.');git(repo,'commit','-m','docs: completed fixture')
 (repo/'bin').mkdir();(repo/'bin/roundfix').symlink_to(root/'bin/roundfix')
 return repo,git(repo,'rev-parse','HEAD')
def snap(repo):return {str(p.relative_to(repo)):hashlib.sha256(p.read_bytes()).hexdigest() for p in repo.rglob('*') if p.is_file() and '.git' not in p.parts and 'bin' not in p.parts}
def match(pattern,value):
 lines=pattern.splitlines();out=value.splitlines()
 def rec(i,j):
  if i==len(lines):return j==len(out)
  if lines[i]=='...':return any(rec(i+1,k) for k in range(j,len(out)+1))
  p=''.join('.+' if x.startswith('<') else re.escape(x) for x in re.split(r'(<[^>]+>)',lines[i]))
  return j<len(out) and re.fullmatch(p,out[j]) is not None and rec(i+1,j+1)
 return rec(0,0)
def check(n,repo,args,out,err,code):
 r=run(repo,['./bin/roundfix',*args],False);x={'transcript':n,'repository':str(repo),'command':['./bin/roundfix',*args],'stdout':r.stdout,'stderr':r.stderr,'exit':r.returncode,'expected_stdout':out,'expected_stderr':err,'expected_exit':code};x['matches']=r.returncode==code and match(out,r.stdout) and match(err,r.stderr);records.append(x);(ev/'transcripts.json').write_text(json.dumps(records,indent=2));assert x['matches'],x
 return r
confirmation=f'archived {slug} -> docs/history/specs/{slug}.md; removed <count> file(s) (<bytes> bytes) kept in Git at <revision>\n'
r,rev=fixture(1);check(1,r,['archive',slug],confirmation,"roundfix: warning: source_revision <revision> is not an ancestor of origin/main; a squash merge keeps it only in the pull request's head ref, so record the pull request with --pull-request once it opens\n",0);assert not (r/'docs/specs'/slug).exists();record=r/f'docs/history/specs/{slug}.md';assert rev in record.read_text();git(r,'show',rev+f':docs/specs/{slug}/_prd.md');git(r,'add','docs');git(r,'commit','-m','docs: archive fixture')
before=snap(r);check(3,r,['archive',slug,'--pull-request','7'],f'recorded pull request 7 in docs/history/specs/{slug}.md\n','',0);after=snap(r);assert [p for p in before if before[p]!=after[p]]==[f'docs/history/specs/{slug}.md'];assert 'pull_request: "7"' in record.read_text();data=record.read_bytes();replay=run(r,['./bin/roundfix','archive',slug,'--pull-request','7']);assert replay.stdout==f'docs/history/specs/{slug}.md already records pull request 7\n' and not replay.stderr and record.read_bytes()==data
records.append({'probe':'record-only replay','stdout':replay.stdout,'exit':0,'byte_identical':True})
git(r,'push','origin','HEAD:main');reason=f'Archive Record docs/history/specs/{slug}.md is already on origin/main; an archived record stays byte-identical'
stderr=f"Preflight failed\n\nReason:\n  {reason}\n\nNo side effects:\n  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.\n\nUsage:\n  Run 'roundfix archive --help' for usage.\n"
check(4,r,['archive',slug,'--pull-request','7'],'',stderr,2);assert record.read_bytes()==data
r2,rev2=fixture(2);check(2,r2,['archive',slug,'--pull-request','7'],confirmation.rstrip('\n')+'; pull request 7\n','roundfix: warning: source_revision <revision> is not an ancestor of origin/main; readers recover it from refs/pull/7/head\n',0);assert 'pull_request: "7"' in (r2/f'docs/history/specs/{slug}.md').read_text()
check(5,r2,['archive','--help'],'Usage:\n  roundfix archive <slug>\n  roundfix archive <slug> --qa-override --approval <source> --reason <text>\n  roundfix archive <slug> --plan\n  roundfix archive <slug> --pull-request <number>\n...\n  --pull-request Store the delivery pull request number in the Archive Record\n...\n','',0)
r3,rev3=fixture(3);git(r3,'push','origin','HEAD:main');git(r3,'checkout','-B','main');rr=run(r3,['./bin/roundfix','archive',slug]);assert not rr.stderr and match(confirmation,rr.stdout);records.append({'probe':'archive on default branch','stdout':rr.stdout,'stderr':rr.stderr,'exit':rr.returncode,'warning_absent':True})
(ev/'transcripts.json').write_text(json.dumps(records,indent=2));(ev/'fixture-location.txt').write_text(str(base)+'\n');print('Five transcripts match stdout/stderr/exit; replay, record-only scope, Git persistence and default branch no-warning checks pass.',flush=True)
