import subprocess, pathlib, tempfile, shutil, json, hashlib, os, tarfile, io
root=pathlib.Path.cwd(); ev=root/'docs/specs/0200-a-skill-snapshot-that-matches-its-upstream/qa/evidence/2026-10-01-gate'; report=ev.parent.parent/'qa-report-2026-10-01.md'
work=pathlib.Path(tempfile.mkdtemp(prefix='roundfix-0200-qa-',dir='/private/tmp')); (ev/'scratch-path.txt').write_text(str(work)+'\n')
env=os.environ.copy(); env['GOCACHE']='/private/tmp/roundfix-0200-qa-cache'; env['GIT_CONFIG_NOSYSTEM']='1'
def run(args,cwd=root,name=None):
 r=subprocess.run(args,cwd=cwd,env=env,text=True,capture_output=True)
 if name: (ev/(name+'.log')).write_text('Command: '+str(args)+'\nCwd: '+str(cwd)+'\nExit: '+str(r.returncode)+'\nSTDOUT:\n'+r.stdout+'\nSTDERR:\n'+r.stderr)
 return r
binary=root/'bin/roundfix'
def cli(args,cwd,name): return run([str(binary),*args],cwd,name)
def git(p,*args):
 r=run(['git','-c','core.fsmonitor=false','-c','core.hooksPath=/dev/null','-c','commit.gpgsign=false','-C',str(p),*args]); assert r.returncode==0,r.stderr; return r.stdout.strip()
def init(p):
 git(p,'init','--quiet');git(p,'config','user.name','QA Fixture');git(p,'config','user.email','qa@example.invalid');git(p,'add','.');git(p,'commit','--quiet','-m','synthetic QA fixture')
def hashes(p): return {str(f.relative_to(p)):hashlib.sha256(f.read_bytes()).hexdigest() for f in p.rglob('*') if f.is_file()}
def settle(row,status,detail):
 s=report.read_text().replace('| '+row+' | pending |','| '+row+' | '+status+' |'); s+='\n### '+row+'\n\n'+detail+'\n';report.write_text(s)
copy=work/'repository';copy.mkdir(); data=run(['git','-c','core.fsmonitor=false','archive','HEAD']) if False else subprocess.check_output(['git','-c','core.fsmonitor=false','archive','HEAD'])
with tarfile.open(fileobj=io.BytesIO(data)) as tar: tar.extractall(copy,filter='data')
init(copy)
# CLI operations run the worktree-built binary, never PATH roundfix. Absolute path is the same ./bin/roundfix artifact.
for name,args in [('doctor',['doctor']),('update',['baseline','update','--repo',str(copy),'--format','text']),('refresh',['baseline','update','--repo',str(copy),'--no-skills','--format','text'])]:
 r=cli(args,copy,name);print(name,r.returncode,r.stdout[:700],r.stderr[:200],flush=True)
 if name=='refresh':
  assert r.returncode==0 and 'File changes: 0' in r.stdout
 if name=='update': assert r.returncode==0 and 'Baseline update: current' in r.stdout
 if name=='doctor': assert 'skills: ok (42 required: 14 Roundfix-owned, 28 external)' in r.stdout
r=run(['make','baseline-digests'],copy,'regenerate');print('regenerate',r.returncode,r.stdout[-400:],r.stderr[-200:],flush=True);assert r.returncode==0 and '"changed":false' in r.stdout
# Context7 cases each start with identical current managed repo, removing only the current skill.
for mode in ['current','prior','neither']:
 target=work/('context7-'+mode);shutil.copytree(copy,target)
 if mode=='prior': (target/'.agents/skills/context7-cli').rename(target/'.agents/skills/context7')
 if mode=='neither': shutil.rmtree(target/'.agents/skills/context7-cli')
 r=cli(['baseline','plan','--profile','go-cli-tui','--repo',str(target),'--format','json'],target,'plan-'+mode)
 doc=json.loads(r.stdout); (ev/('plan-'+mode+'-pretty.json')).write_text(json.dumps(doc,indent=2)+'\n');print('plan',mode,r.returncode,flush=True)
# Fresh synthetic source holds every setup path and a committed GitHub origin.
source=work/'source';(source/'setups').mkdir(parents=True)
for f in (root/'internal/baseline/assets/setups').glob('*.json'):
 doc=json.loads(f.read_text());paths=[x['path'] for x in doc['skills']];(source/'setups'/(doc['id']+'.txt')).write_text('# canonical QA source\n'+'\n'.join(paths)+'\n')
 for path in paths:
  dest=source/path/'SKILL.md';dest.parent.mkdir(parents=True,exist_ok=True);dest.write_text('# '+pathlib.Path(path).name+'\n')
init(source);git(source,'remote','add','origin','https://github.com/example/skills.git')
# Rename source handoff; target requires this renamed skill before its snapshots do.
handoff=next(source.glob('skills/*/handoff'));handoff.rename(handoff.with_name('handoff-next'))
for f in (source/'setups').glob('*.txt'): f.write_text(f.read_text().replace('/handoff\n','/handoff-next\n'))
git(source,'add','.');git(source,'commit','--quiet','-m','synthetic upstream rename')
for mode in ['good','bad']:
 target=work/('sync-'+mode);assets=target/'internal/baseline/assets';shutil.copytree(root/'internal/baseline/assets',assets);init(target)
 core=assets/'modules/core.json';core.write_text(core.read_text().replace('"handoff"','"handoff-next"' if mode=='good' else '"handoff-absent"'))
 before=hashes(assets)
 for check in [True,False]:
  args=['baseline','assets','sync','--source-dir',str(source/'setups'),'--format','text']+(['--check'] if check else [])
  r=cli(args,target,'sync-'+mode+('-check' if check else '-apply'));print('sync',mode,check,r.returncode,r.stdout[:250],r.stderr[:150],flush=True)
  expected=(1 if check else 0) if mode=='good' else 2
  assert r.returncode==expected,(expected,r.stdout,r.stderr)
  if check or mode=='bad': assert hashes(assets)==before
  if mode=='bad': assert 'Generated setup snapshots are incompatible with the Baseline catalog' in r.stdout+r.stderr
 if mode=='good':
  for f in (assets/'setups').glob('*.json'):
   names={x['name'] for x in json.loads(f.read_text())['skills']};assert 'handoff-next' in names and 'handoff' not in names
  r=cli(['baseline','assets','sync','--source-dir',str(source/'setups'),'--check'],target,'sync-good-recheck');assert r.returncode==0
settle('02','pass','Maintainer / CLI: synthetic target core requires handoff-next before the source-only rename is synchronized. Check exited 1 without byte changes; apply exited 0; fresh check exited 0 and all four snapshots persist handoff-next without handoff. A separate absent-skill target refused check and apply with exit 2 and the exact incompatibility message, preserving all asset hashes. Evidence: `evidence/2026-10-01-gate/sync-*.log` and this reproducible harness. No installed/adopter tree was modified.')
print('scratch',work,flush=True)
