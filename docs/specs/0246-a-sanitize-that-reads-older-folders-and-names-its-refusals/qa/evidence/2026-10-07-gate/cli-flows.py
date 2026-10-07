from pathlib import Path
import subprocess,os,json,tempfile,re,hashlib
ROOT=Path(__file__).resolve().parents[5]
# Explicit repository root passed for replay after archive.
ROOT=Path(os.environ.get('QA_REPOSITORY',Path.cwd()))
E=Path(__file__).resolve().parent
G="""---
schema: spec-tasks/v1
qa: task_qa
graph:
  nodes:
    - id: task_01
      file: task_01.md
    # - id: task_02
    #   file: task_02.md
    - id: task_qa
      file: task_qa.md
      needs: [task_01]
---

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
"""
def setup(shapes):
 d=Path(tempfile.mkdtemp(prefix='roundfix-0246-qa-'));repo=d/'repo';repo.mkdir();home=d/'home';home.mkdir()
 env={'PATH':os.environ['PATH'],'HOME':str(home),'XDG_CONFIG_HOME':str(home/'config'),'GIT_CONFIG_NOSYSTEM':'1','GIT_CONFIG_GLOBAL':'/dev/null'}
 def git(*a):
  p=subprocess.run(['git',*a],cwd=repo,env=env,text=True,capture_output=True);assert p.returncode==0,p.stderr;return p.stdout.strip()
 git('init','-b','main');git('config','user.name','QA');git('config','user.email','qa@example.invalid');git('config','commit.gpgsign','false');git('config','tag.gpgsign','false')
 (repo/'docs/specs').mkdir(parents=True);(repo/'README.md').write_text('# Synthetic QA\n')
 for slug,shape in shapes:
  f=repo/'docs/history/specs'/slug;f.mkdir(parents=True)
  (f/'_prd.md').write_text('invalid PRD' if shape=='refused' else f'---\nspec: {slug}\ncreated: 2026-09-01\n---\n\n# {slug}\n\nPreserved synthetic outcome.\n')
  if shape not in ['clean','refused']:
   rows='| task_01 | Implementation | '+('refactor' if shape=='retired' else 'backend')+' | low | - |\n'
   if shape=='outside':rows+='| task_02 | Omitted | backend | low | - |\n'
   if shape=='outside-many':rows+='| task_03 | Omitted | backend | low | - |\n| task_04 | Omitted | backend | low | - |\n'
   (f/'_tasks.md').write_text(G+rows+'| task_qa | QA | qa | low | task_01 |\n')
   for task in ['task_01','task_qa']:(f/(task+'.md')).write_text(f'---\ntask: {task}\nspec: {slug}\nstatus: completed\ntype: '+('qa' if task=='task_qa' else 'backend')+'\n---\n\n# Synthetic Task\n')
   if shape=='failed':
    (f/'qa').mkdir()
    for day in [1,2,3]:(f/'qa'/f'qa-report-2026-09-0{day}.md').write_text('---\nverdict: fail\n---\n\n# Synthetic failed QA\n')
 git('add','.');git('commit','-m','Synthetic fixtures (#42)');rev=git('rev-parse','HEAD');git('tag','-a','history-full','-m','Full synthetic history')
 (repo/'bin').symlink_to(ROOT/'bin',target_is_directory=True)
 (repo/'.git/info/exclude').write_text('bin\n')
 return repo,env,rev

def snapshot(repo):return {str(p.relative_to(repo)):hashlib.sha256(p.read_bytes()).hexdigest() for p in (repo/'docs').rglob('*') if p.is_file()}
def run(repo,env,name,*args):
 cmd=['./bin/roundfix','history','sanitize',*args];p=subprocess.run(cmd,cwd=repo,env=env,text=True,capture_output=True)
 data={'cwd':str(repo),'command':cmd,'exit':p.returncode,'stdout':p.stdout,'stderr':p.stderr};(E/(name+'.json')).write_text(json.dumps(data,indent=2));return p

def match(actual,lines):
 pats=[]
 for line in lines:
  if line=='...':pats.append(r'(?:[^\n]*\n)*');continue
  chunks=re.split(r'(<[^>]+>)',line);pats.append(''.join(r'[^\n]+' if c.startswith('<') and c.endswith('>') else re.escape(c) for c in chunks)+r'\n')
 assert re.fullmatch(''.join(pats),actual),(actual,lines)

repo,env,rev=setup([('aaa','outside'),('bbb','failed'),('ccc','refused')]);before=snapshot(repo);p=run(repo,env,'transcript-1');assert p.returncode==0 and p.stderr=='' and snapshot(repo)==before
match(p.stdout,['history sanitize plan: <u> unit(s) pending; <f> file(s) (<b> bytes) leave docs/history; 1 unit(s) refused','...','tolerates docs/history/specs/<slug>: projection row <id> names a Task outside the graph','...','folder docs/history/specs/<slug>: removes <n> file(s) (<b> bytes) and writes docs/history/specs/<slug>.md (<b> bytes, failed-qa, delivery <d>)','...','refused docs/history/specs/<slug>: <reason>','...','apply with: roundfix history sanitize --apply --batch <n> (needs the annotated tag history-full at or before HEAD)'])
repo,env,rev=setup([('aaa','refused'),('bbb','clean'),('ccc','clean')]);before=snapshot(repo);p=run(repo,env,'transcript-2','--apply','--batch','2');assert p.returncode==0 and p.stderr==''
match(p.stdout,['refused docs/history/specs/<slug>: <reason>','history sanitize applied 2 unit(s): wrote 2 Archive Record(s), reduced 0 file(s), removed <d> file(s) (<b> bytes) kept in Git at <rev> and tag history-full; promoted 0 file(s) to docs/references/; 1 unit(s) refused; <r> unit(s) remain'])
assert snapshot(repo)['docs/history/specs/aaa/_prd.md']==before['docs/history/specs/aaa/_prd.md']
assert all((repo/f'docs/history/specs/{s}.md').exists() and not (repo/f'docs/history/specs/{s}').exists() for s in ['bbb','ccc']);r=run(repo,env,'transcript-2-readback');assert r.returncode==0 and '1 unit(s) refused' in r.stdout and 'folder ' not in r.stdout
repo,env,rev=setup([('aaa','refused')]);before=snapshot(repo);p=run(repo,env,'transcript-3','--apply','--batch','1');assert p.returncode==2 and snapshot(repo)==before
match(p.stdout,['refused docs/history/specs/<slug>: <reason>']);match(p.stderr,['Preflight failed','...']);r=run(repo,env,'transcript-3-readback');assert r.returncode==0 and snapshot(repo)==before
repo,env,rev=setup([('aaa','outside'),('bbb','outside-many'),('ccc','retired'),('ddd','failed'),('eee','clean')]);before=snapshot(repo);p=run(repo,env,'four-shapes-plan');assert p.returncode==0 and p.stderr=='' and snapshot(repo)==before and '5 unit(s) pending' in p.stdout and 'failed-qa' in p.stdout
p=run(repo,env,'four-shapes-apply','--apply','--batch','5');assert p.returncode==0 and p.stderr=='' and 'wrote 5 Archive Record(s)' in p.stdout
records=[]
for slug in ['aaa','bbb','ccc','ddd','eee']:
 f=repo/f'docs/history/specs/{slug}.md';assert f.exists() and not (repo/f'docs/history/specs/{slug}').exists();text=f.read_text();assert f'source_revision: {rev}' in text and 'tolerates' not in text;records.append(str(f))
 if slug=='ddd':assert 'disposition: failed-qa' in text and 'qa_verdict: fail' in text and 'qa_report: qa-report-2026-09-03.md' in text and 'qa_override:' not in text
r=run(repo,env,'four-shapes-readback');assert r.returncode==0 and r.stderr=='' and 'nothing pending' in r.stdout
(E/'records.json').write_text(json.dumps(records,indent=2));print('Three transcript matches and five-shape conversion/readback PASS; records:',records)
