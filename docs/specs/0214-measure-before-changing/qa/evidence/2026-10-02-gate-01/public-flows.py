import sys,tempfile,sqlite3
from pathlib import Path
sys.path.insert(0,str(Path(__file__).parent))
from gate import *
with tempfile.TemporaryDirectory(prefix='roundfix-0214-qa-',dir='/private/tmp') as tmp:
 base=Path(tmp);repo=base/'repo';home=base/'home';outside=base/'outside'
 for p in (repo,home,outside):p.mkdir()
 subprocess.run(['git','init','--initial-branch=main',str(repo)],check=True,capture_output=True)
 (repo/'bin').mkdir();(repo/'bin/roundfix').symlink_to(ROOT/'bin/roundfix')
 (outside/'bin').mkdir();(outside/'bin/roundfix').symlink_to(ROOT/'bin/roundfix')
 d=repo/'docs/specs/0300-example';d.mkdir(parents=True)
 d.joinpath('_tasks.md').write_text('---\nschema: spec-tasks/v1\nqa: task_03\ngraph:\n  nodes:\n    - {id: task_01, file: task_01.md}\n    - {id: task_02, file: task_02.md}\n    - {id: task_03, file: task_03.md}\n    - {id: task_04, file: task_04.md}\n---\n')
 d.joinpath('task_04.md').write_text('# Corrective Task\n\n## Overview\nPre-PR review found that the golden output was not re-recorded\n')
 dbpath=home/'.roundfix/roundfix.db';dbpath.parent.mkdir();db=sqlite3.connect(dbpath)
 source=(ROOT/'internal/store/store.go').read_text()
 for table in ('runs','run_events'):
  sql=source.split('`CREATE TABLE IF NOT EXISTS '+table+' (',1)[1].split('`',1)[0];db.execute('CREATE TABLE '+table+' ('+sql)
 db.execute('PRAGMA user_version = 22')
 rid='run_20261001T120000Z_0000000000000001';at='2026-10-01T12:00:00Z';art=home/'artifacts';log=art/'runs'/rid/'verification';log.mkdir(parents=True)
 db.execute('INSERT INTO runs (id,kind,state,git_root,repository_root,local_branch,spec_slug,artifact_dir,created_at,updated_at,completed_at) VALUES (?,?,?,?,?,?,?,?,?,?,?)',(rid,'implement','Clean',str(repo),str(repo),'main','0300-example',str(art),at,at,at))
 cur=0
 def event(kind,payload):
  global cur
  cur+=1;db.execute('INSERT INTO run_events(run_id,cursor,source,kind,created_at,payload) VALUES(?,?,?,?,?,?)',(rid,cur,'daemon',kind,at,json.dumps(payload)))
 for task in ('task_01','task_02','task_03','task_04'):event('daemon.task',{'task':task,'phase':'started'})
 for i,(task,attempt,cmd,diag) in enumerate([('task_01',1,'go test ./internal/store','database is locked'),('task_02',1,'go test ./internal/cli','--- FAIL: TestCausesFixture'),('task_02',2,'make verify','no matching words')]):
  path=log/f'{i}.log';path.write_text(diag);event('daemon.verification',{'task':task,'attempt':attempt,'phase':'failed','command':cmd,'diagnostic_path':str(path)});event('daemon.verification',{'task':task,'attempt':attempt,'phase':'verdict','verdict':'failed'})
 event('daemon.task',{'task':'task_01','phase':'verification_feedback'});event('daemon.verification',{'task':'task_01','phase':'verdict','attempt':2,'verdict':'passed'})
 db.commit();db.close()
 env=ENV.copy();env['HOME']=str(home)
 def hashes():return {str(p.relative_to(home)):hashlib.sha256(p.read_bytes()).hexdigest() for p in home.rglob('*') if p.is_file()}
 before=hashes();checks=[]
 def invoke(n,args,label,cwd=repo):
  r=run(n,['./bin/roundfix']+args,label,cwd,env);after=hashes();assert after==before,(label,'home changed');checks.append({'command':args,'before':before,'after':after});return r
 ts=(ROOT/'docs/specs/0214-measure-before-changing/_techspec.md').read_text().split('### Surface Transcripts',1)[1]
 blocks=[__import__('textwrap').dedent(b) for b in re.findall(r'```transcript\n(.*?)```',ts,re.S)]
 for n,args,block,cwd in [(4,['runs','causes','--since','2026-10-01'],blocks[0],repo),(5,['runs','causes','--until','2026-09-01'],blocks[1],repo),(6,['runs','causes','--since','2026-13-01'],blocks[2],repo),(7,['runs','causes'],blocks[3],outside)]:
  r=invoke(n,args,f'row-{n}',cwd);wantout=block.split('stdout:\n')[1].split('stderr:\n')[0];wanterr=block.split('stderr:\n')[1].split('exit:')[0];code=int(block.split('exit:')[1].strip());digest=hashlib.sha256((ROOT/'internal/runcause/signatures.json').read_bytes()).hexdigest()[:12]
  wantout=wantout.replace('<digest>',digest);ok=(r.stdout==wantout and r.stderr==wanterr and r.returncode==code)
  print(n,ok,r.returncode,r.stderr[:250]);settle(n,'pass' if ok else 'fail',f'Built binary transcript matches stdout, stderr and exit {code} exactly: [capture](evidence/2026-10-02-gate-01/row-{n}.txt). Fresh SHA-256 of every disposable-home file equals its pre-command hash. Literal ./bin/roundfix is a symlink to the binary built in this Run Worktree.')
 r=invoke(8,['runs','causes','--since','2026-10-01','--until','2026-10-02','--format','json'],'public-json');j=json.loads(r.stdout);assert r.returncode==0 and not r.stderr
 assert set(j)=={'schema','window','signatures_sha256','runs','items','tasks','summary','specs_not_found'}
 assert len(j['items'])==4 and len(j['tasks'])==4 and j['summary']['classified']==3
 assert j['items'][2]['signature'] is None and j['items'][2]['evidence'] is None and j['items'][3]['attempt'] is None
 for item in j['items']:assert set(item)=={'kind','spec','task','run_id','attempt','check','class','signature','evidence'}
 for task in j['tasks']:assert set(task)=={'spec','task','qa','corrective','verdicts_passed','verdicts_failed','feedback_rounds','runs'}
 synopsis='roundfix runs causes [--since <YYYY-MM-DD>] [--until <YYYY-MM-DD>] [--format <text|json>]'
 for k,args in enumerate([['runs','causes','--help'],['runs','--help'],['--help']]):
  r=invoke(8,args,f'public-help-{k}');assert r.returncode==0 and not r.stderr and synopsis in r.stdout
 for k,args in enumerate([['--format','yaml'],['--unknown'],['--since','2026-10-01','--until','2026-10-01']]):
  r=invoke(8,['runs','causes']+args,f'public-invalid-{k}');assert r.returncode==2 and not r.stdout
 r=invoke(8,['runs','causes','--since','2026-10-01'],'public-repeat');assert r.stdout==(E/'row-4.txt').read_text().split('\n',1)[1].split('\nSTDERR:\n')[0]
 settle(8,'pass','JSON includes every top-level/item/Task/nullable field and agrees with text counts; all three help commands contain API Contract 3 synopsis. Invalid format, unknown flag and inverted window exit 2. Fresh-process repeat equals transcript 1. [JSON](evidence/2026-10-02-gate-01/public-json.txt), help-0/1/2 and invalid-0/1/2 captures; [per-command SHA-256](evidence/2026-10-02-gate-01/home-hashes.json).')
 (E/'home-hashes.json').write_text(json.dumps(checks,indent=2))
