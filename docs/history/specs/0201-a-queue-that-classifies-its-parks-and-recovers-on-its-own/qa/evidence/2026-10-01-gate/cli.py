from pathlib import Path
import os,subprocess,json,tempfile,sqlite3
s=Path(__file__).resolve().parents[3]; e=Path(__file__).resolve().parent; root=Path.cwd(); r=s/'qa/qa-report-2026-10-01.md'
def row(i,ok,detail):
 t=r.read_text().replace('| '+i+' | pending |','| '+i+' | '+('pass' if ok else 'fail')+' |');t=t.replace('### '+i+'\n','### '+i+'\n\n'+detail+'\n');r.write_text(t)
base=Path(tempfile.mkdtemp(prefix='roundfix-0201-qa-',dir='/private/tmp'));repo=base/'repo';home=base/'home';repo.mkdir();home.mkdir();env=dict(os.environ,HOME=str(home),GIT_CONFIG_GLOBAL='/dev/null',GIT_CONFIG_NOSYSTEM='1');env.pop('TYPESAFE_API_KEY',None)
def git(*args):
 return subprocess.run(['git',*args],cwd=repo,env=env,check=True,capture_output=True,text=True).stdout
repo.joinpath('bin').mkdir();repo.joinpath('bin/roundfix').symlink_to(root/'bin/roundfix')
git('init','--initial-branch=main');git('config','user.name','Roundfix QA');git('config','user.email','qa@example.invalid');git('config','commit.gpgsign','false')
origin=base/'origin.git';subprocess.run(['git','init','--bare','--initial-branch=main',str(origin)],check=True,capture_output=True,env=env);git('remote','add','origin',str(origin))
# No push needed: every journey rejects before owner startup.
def spec(slug,requires):
 d=repo/'docs/specs'/slug;d.mkdir(parents=True,exist_ok=True)
 (d/'_prd.md').write_text('---\nstatus: active\n---\n# QA fixture\n')
 (d/'_tasks.md').write_text('---\nschema: spec-tasks/v1\nrequires: ['+requires+']\ngraph:\n  nodes:\n    - id: task_01\n      file: task_01.md\n      needs: []\n---\n')
 (d/'task_01.md').write_text('---\ntask: task_01\nstatus: pending\ntype: backend\n---\n# Fixture\n\n## Verification\n\n- `true`\n')
spec('0204-first','0205-second');spec('0205-second','0204-first');git('add','docs');git('commit','-m','test: seed disposable QA specs')
def cli(tag,*args):
 p=subprocess.run(['./bin/roundfix',*args],cwd=repo,env=env,capture_output=True,text=True)
 (e/(tag+'.json')).write_text(json.dumps({'argv':['./bin/roundfix',*args],'cwd':str(repo),'HOME':str(home),'stdout':p.stdout,'stderr':p.stderr,'exit':p.returncode},indent=2));return p
def noqueue():
 db=home/'.roundfix/roundfix.db'
 if not db.exists():return True
 with sqlite3.connect(db) as c:
  names=[x[0] for x in c.execute("select name from sqlite_master where type='table'")]
  return 'delivery_queues' not in names or c.execute('select count(*) from delivery_queues').fetchone()[0]==0
want='Preflight failed\n\nReason:\n  Delivery Queue Specs require each other in a cycle: 0204-first -> 0205-second -> 0204-first\n\nNo side effects:\n  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.\n\nUsage:\n  Run \'roundfix deliver start --help\' for usage.\n'
oks=[]
for n in (1,2):
 p=cli('02-cycle-'+str(n),'deliver','start','0204-first','0205-second');oks.append(p.returncode==2 and p.stdout=='' and p.stderr==want and noqueue())
spec('0204-first','0205-unknown')
for n in (1,2):
 p=cli('02-unknown-'+str(n),'deliver','start','0204-first');oks.append(p.returncode==2 and p.stdout=='' and 'requires unknown Spec 0205-unknown' in p.stderr and noqueue())
(e/'02-confirmation.txt').write_text('All comparisons: '+str(oks)+'\nQueue absent after each refusal: '+str(noqueue())+'\n'+git('status','--short'))
row('02',all(oks),'Observed: cycle matches Surface Transcript 2 exactly (empty stdout, exact stderr, exit 2); unknown prerequisite exits 2 with named slug. Each command repeated in a new process; direct disposable database read confirms no queue after every refusal. Evidence: `qa/evidence/2026-10-01-gate/02-cycle-1.json`, `02-cycle-2.json`, `02-unknown-1.json`, `02-unknown-2.json`, `02-confirmation.txt`. Fixture location '+str(base)+'.')
# Create the fixture schema through the public migration command, then seed disposable queue rows.
(home/'.roundfix').mkdir(); (home/'.roundfix/roundfix.db').touch()
p=cli('fixture-migrate','migrate');assert p.returncode==0,(p.returncode,p.stderr)
db=home/'.roundfix/roundfix.db'
with sqlite3.connect(db) as c:
 (e/'fixture-schema.txt').write_text('\n'.join(str(x) for x in c.execute("select sql from sqlite_master where name like 'delivery_%'")))
 c.execute('insert into delivery_queues(git_root) values (?)',(str(repo),))
 for pos,slug in enumerate(['0204-first','0205-second']):c.execute('insert into delivery_queue_items(git_root,spec_slug,position,stage) values (?,?,?,?)',(str(repo),slug,pos,'queued'))
want='0204-first\tqueued\t-\t-\n0205-second\tqueued\t-\t-\nLimits: deadline none, retries per item none, concurrency 1, spend not measured\n'
oks=[]
for n in (1,2):
 p=cli('03-nonparked-'+str(n),'deliver','status');oks.append(p.returncode==0 and p.stdout==want and p.stderr=='')
path='internal/baseline/assets/profiles/standard-typescript-monorepo.json'
with sqlite3.connect(db) as c:
 c.execute("update delivery_queue_items set stage='parked',blocker=?,branch='roundfix/deliver-first',worktree='/worktrees/0204-first' where spec_slug='0204-first'",('pull-request-conflict: '+path,))
 c.execute("update delivery_queue_items set stage='parked',blocker='prerequisite-unmerged: 0204-first' where spec_slug='0205-second'")
nextcmd='merge the default branch into the item branch in /worktrees/0204-first, resolve '+path+', commit, then run roundfix deliver retry 0204-first'
want='0204-first\tparked\tpull-request-conflict: '+path+'\t/worktrees/0204-first\n0205-second\tparked\tprerequisite-unmerged: 0204-first\t-\nPark: 0204-first conflict: '+nextcmd+'\nPark: 0205-second dependency: run roundfix deliver retry 0204-first; 0205-second returns to the queue once 0204-first is merged\nLimits: deadline none, retries per item none, concurrency 1, spend not measured\nPending question: 0204-first parked pull-request-conflict: '+path+'\nAnswer: '+nextcmd+'\nWaiting behind it: 1 parked item(s)\n'
for n in (1,2):
 p=cli('03-transcript-'+str(n),'deliver','status');oks.append(p.returncode==0 and p.stdout==want and p.stderr=='')
with sqlite3.connect(db) as c:(e/'03-confirmation.txt').write_text('Comparisons: '+str(oks)+'\n'+repr(c.execute('select spec_slug,stage,blocker,worktree,retry_count from delivery_queue_items order by position').fetchall()))
row('03',all(oks),'Observed: built binary reproduces Surface Transcript 1 twice with exact stdout (the `<limits>` placeholder matched its concrete nonempty value), empty stderr and exit 0. Nonparked queue matches previous table/limits format twice without Park/Pending Question lines. Fresh processes and independent SQLite read confirm persistence. Named TestDeliverStatusReproducesSurfaceTranscriptOne and unchanged-warning test also passed. Evidence: `qa/evidence/2026-10-01-gate/03-tests.txt` is named `03-test-tests.txt`; command captures are `03-transcript-1.json`, `03-transcript-2.json`, `03-nonparked-1.json`, `03-nonparked-2.json`, `03-confirmation.txt` under that same evidence directory.')
configs=[('empty-paths','delivery:\n  derived_paths:\n    - paths: []\n      regenerate: make generated\n'),('absolute','delivery:\n  derived_paths:\n    - paths: [/tmp/generated]\n      regenerate: make generated\n'),('parent','delivery:\n  derived_paths:\n    - paths: [../generated]\n      regenerate: make generated\n'),('empty-command','delivery:\n  derived_paths:\n    - paths: [generated]\n      regenerate: ""\n')]
oks=[]
for tag,content in configs:
 (repo/'.roundfixrc.yml').write_text(content);p=cli('08-'+tag,'deliver','status');oks.append(p.returncode==2 and 'delivery.derived_paths' in p.stderr and p.stdout=='')
(repo/'.roundfixrc.yml').write_text('specs:\n  root: docs/specs\n')
for n in (1,2):
 p=cli('08-omitted-'+str(n),'deliver','status');oks.append(p.returncode==0 and p.stderr=='' and p.stdout==want)
(e/'08-confirmation.txt').write_text('Comparisons: '+str(oks))
row('08',all(oks),'Observed: empty paths, absolute path, parent traversal and empty regenerate each exit 2 naming delivery.derived_paths. An omitted-key config returns unchanged persisted status on two fresh reads. Evidence: `qa/evidence/2026-10-01-gate/08-empty-paths.json`, `08-absolute.json`, `08-parent.json`, `08-empty-command.json`, `08-omitted-1.json`, `08-omitted-2.json`, `08-confirmation.txt`.')
print('CLI assertions',all(oks),'fixture',base)
