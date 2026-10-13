from pathlib import Path
import json, os, shutil, subprocess, tempfile
repo_source=Path.cwd().resolve()
evidence=Path(__file__).resolve().parent
binary=repo_source/'bin/roundfix'
log=[]
def run(args,cwd,env,allow=0):
 p=subprocess.run(args,cwd=cwd,env=env,text=True,capture_output=True)
 log.append('COMMAND: '+repr(args)+'\nCWD: '+str(cwd)+'\nEXIT: '+str(p.returncode)+'\nSTDOUT:\n'+p.stdout+'STDERR:\n'+p.stderr)
 if p.returncode!=allow: raise RuntimeError(log[-1])
 return p.stdout
try:
 cache=subprocess.check_output(['go','env','GOPATH','GOCACHE'],text=True).splitlines()
 with tempfile.TemporaryDirectory(prefix='rf262-binary-',dir='/private/tmp') as scratch:
  root=Path(scratch); home=root/'home'; home.mkdir(); repo=root/'repo'; repo.mkdir(); real=root/'real'; real.mkdir(); link=root/'linked'; link.symlink_to(real,target_is_directory=True)
  env={k:v for k,v in os.environ.items() if not k.startswith(('GIT_','ROUNDFIX_')) and k!='NODE_OPTIONS'}
  env.update(HOME=str(home),GIT_CONFIG_NOSYSTEM='1',GIT_CONFIG_GLOBAL='/dev/null',GIT_TERMINAL_PROMPT='0',GOPATH=cache[0],GOCACHE=cache[1],GOPROXY='off',GOSUMDB='off',GOTOOLCHAIN='local')
  run(['git','init','-b','main'],repo,env)
  for key,value in [('user.name','QA fixture'),('user.email','qa@example.invalid'),('commit.gpgsign','false'),('core.hooksPath','/dev/null')]: run(['git','config',key,value],repo,env)
  (repo/'.gitignore').write_text('bin/\n')
  (repo/'.roundfixrc.yml').write_text('worktree:\n  location: '+json.dumps(str(link))+'\n')
  s=repo/'docs/specs/qa-symlinked-root'; s.mkdir(parents=True)
  (s/'_prd.md').write_text('---\nspec: qa-symlinked-root\nstatus: active\n---\n\n# Disposable QA fixture\n')
  (s/'task_01.md').write_text('---\ntask: task_01\nspec: qa-symlinked-root\nstatus: completed\ntype: backend\n---\n\n# Delivered fixture\n')
  (s/'_tasks.md').write_text('---\nschema: spec-tasks/v1\nspec: qa-symlinked-root\ngraph:\n  nodes:\n    - id: task_01\n      file: task_01.md\n      needs: []\n---\n')
  run(['git','add','.'],repo,env); run(['git','commit','-m','docs: seed disposable QA fixture'],repo,env)
  helper=root/'helper'; helper.mkdir(); (helper/'main.go').write_bytes((evidence/'fixture.go.txt').read_bytes()); shutil.copy(repo_source/'go.sum',helper/'go.sum')
  run(['go','mod','init','roundfix/qa'],helper,env)
  run(['go','mod','edit','-go=1.26','-require=roundfix@v0.0.0','-replace=roundfix='+str(repo_source)],helper,env)
  data=json.loads(run(['go','run','-mod=mod','.','seed',str(home),str(repo),str(link)],helper,env))
  (evidence/'R3-fixture.json').write_text(json.dumps(data,indent=2)+'\n')
  (repo/'bin').mkdir(); (repo/'bin/roundfix').symlink_to(binary)
  assert (repo/'bin/roundfix').resolve()==binary.resolve()
  version=run(['./bin/roundfix','--version'],repo,env)
  dry=json.loads(run(['./bin/roundfix','reconcile','--format','json'],repo,env)); (evidence/'R3-dry-run.json').write_text(json.dumps(dry,indent=2)+'\n')
  row=next(r for r in dry['results'] if r['runId']==data['run_id'])
  assert row['classification']=='safe',row
  assert row['worktree']==data['worktree'] and row['refusalReason']=='',row
  assert data['worktree']!=data['resolved'] and Path(data['worktree']).is_dir()
  listing=run(['git','worktree','list','--porcelain'],repo,env); assert 'worktree '+data['resolved']+'\n' in listing
  run(['git','show-ref','--verify','refs/heads/'+data['branch']],repo,env)
  applied=json.loads(run(['./bin/roundfix','reconcile','--apply','--format','json'],repo,env)); (evidence/'R3-apply.json').write_text(json.dumps(applied,indent=2)+'\n')
  row=next(r for r in applied['results'] if r['runId']==data['run_id']); assert row['classification']=='safe' and row['action']=='released',row
  assert not Path(data['worktree']).exists()
  listing=run(['git','worktree','list','--porcelain'],repo,env); assert data['resolved'] not in listing and data['branch'] not in listing
  run(['git','show-ref','--verify','refs/heads/'+data['branch']],repo,env,allow=128)
  fresh=json.loads(run(['./bin/roundfix','reconcile',data['run_id'],'--format','json'],repo,env)); (evidence/'R3-fresh-process.json').write_text(json.dumps(fresh,indent=2)+'\n')
  assert fresh['results'][0]['classification']=='released',fresh
  log.append('ASSERTIONS: safe dry-run retains recorded spelling, directory, registration and branch; apply releases all three; a fresh process reports released.\nBINARY: '+version)
 print('R3 binary journey passed')
finally:
 (evidence/'R3-binary.log').write_text('\n'.join(log))
