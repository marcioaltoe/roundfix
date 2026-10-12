import os,subprocess,tempfile,json,shutil,hashlib
from pathlib import Path
root=Path.cwd().resolve();ev=root/'docs/specs/0256-an-archive-record-that-survives-a-squash-merge/qa/evidence/2026-10-12-local';old=Path((ev/'fixture-location.txt').read_text().strip());base=Path(tempfile.mkdtemp(prefix='roundfix-0256-recovery-',dir='/tmp'));slug='0300-qa-archive';seed=Path((ev/'store-seed-module.txt').read_text().strip())/'seed'
results=[]
for mode in ['record','subject','unreachable']:
 env=os.environ.copy()
 for k in list(env):
  if k.startswith('ROUNDFIX_') or k.startswith('GIT_'):env.pop(k,None)
 home=base/(mode+'-home');home.mkdir();env.update(HOME=str(home),XDG_CONFIG_HOME=str(home/'config'),GIT_CONFIG_NOSYSTEM='1',GIT_CONFIG_GLOBAL='/dev/null',GIT_TERMINAL_PROMPT='0',NO_COLOR='1')
 def run(repo,args,ok=True):
  r=subprocess.run(args,cwd=repo,env=env,capture_output=True,text=True,timeout=60)
  if ok and r.returncode:raise RuntimeError((args,r.returncode,r.stdout,r.stderr))
  return r
 def git(repo,*args,ok=True):return run(repo,['git',*args],ok).stdout.strip()
 producer=base/(mode+'-producer');git(base,'clone','--no-local',str(old/'repo2'),str(producer));git(producer,'config','user.name','QA');git(producer,'config','user.email','qa@example.invalid');git(producer,'config','commit.gpgsign','false')
 recordpath=f'docs/history/specs/{slug}.md';record=(old/'repo2'/recordpath).read_text();revision=record.split('source_revision: ')[1].splitlines()[0].strip('"')
 if mode=='subject':record='\n'.join(l for l in record.split('\n') if not l.startswith('pull_request:'))
 shutil.rmtree(producer/'docs/specs'/slug);(producer/'docs/history/specs').mkdir(parents=True);(producer/recordpath).write_text(record);git(producer,'add','docs');git(producer,'commit','-m','docs: archive fixture');prhead=git(producer,'rev-parse','HEAD')
 git(producer,'checkout','-B','main','origin/main');git(producer,'merge','--squash',prhead);git(producer,'commit','-m','feat: deliver QA fixture (#7)')
 clone=base/(mode+'-clone');git(base,'clone','--no-local','--single-branch','--branch','main',str(producer),str(clone));clone=clone.resolve();forge=base/(mode+'-forge.git');git(base,'init','--bare','--initial-branch=main',str(forge));git(forge,'fetch','--no-tags','--no-write-fetch-head',str(producer),prhead+':refs/pull/7/head');git(forge,'fetch','--no-tags','--no-write-fetch-head',str(producer),'main:refs/heads/main');git(clone,'remote','set-url','origin',str(forge if mode!='unreachable' else base/'absent.git'))
 assert run(clone,['git','cat-file','-e',revision+'^{commit}'],False).returncode!=0
 assert not (clone/'.git/FETCH_HEAD').exists()
 (clone/'bin').mkdir();(clone/'bin/roundfix').symlink_to(root/'bin/roundfix');run(clone,[str(seed),str(home),str(clone),slug]);before=git(clone,'for-each-ref');refsfile=ev/('refs-'+mode+'-before.txt');refsfile.write_text(before+'\n')
 r=run(clone,['./bin/roundfix','runs','causes','--format','json']);report=json.loads(r.stdout);item=report['archived_specs'][0];assert item['source_available']==(mode!='unreachable'),report
 after=git(clone,'for-each-ref');(ev/('refs-'+mode+'-after.txt')).write_text(after+'\n');assert before==after;assert not (clone/'.git/FETCH_HEAD').exists();rr=run(clone,['./bin/roundfix','runs','causes','--format','json']);assert json.loads(rr.stdout)['archived_specs'][0]['source_available']==(mode!='unreachable')
 if mode!='unreachable':git(clone,'show',revision+f':docs/specs/{slug}/_prd.md')
 (ev/('causes-'+mode+'.json')).write_text(r.stdout)
 results.append({'mode':mode,'command':['./bin/roundfix','runs','causes','--format','json'],'exit':r.returncode,'stderr':r.stderr,'source_revision':revision,'source_missing_before':True,'source_available':item['source_available'],'refs_identical':before==after,'FETCH_HEAD_absent':True,'fresh_process_confirmation':True,'clone':str(clone),'remote':str(forge),'record_sha256':hashlib.sha256((clone/recordpath).read_bytes()).hexdigest()});(ev/'recovery-results.json').write_text(json.dumps(results,indent=2));print(mode,results[-1],flush=True)
