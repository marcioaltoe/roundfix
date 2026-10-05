from pathlib import Path
import os,subprocess,tempfile,json,re,shutil
root=Path.cwd(); evidence=Path(__file__).resolve().parent
scratch=Path(tempfile.mkdtemp(prefix="roundfix-0228-qa-",dir="/private/tmp")); repo=scratch/'item'; remote=scratch/'origin.git'
env=os.environ|{'GOCACHE':'/private/tmp/roundfix-0228-task04-gocache','GIT_CONFIG_COUNT':'2','GIT_CONFIG_KEY_0':'core.fsmonitor','GIT_CONFIG_VALUE_0':'false','GIT_CONFIG_KEY_1':'core.hooksPath','GIT_CONFIG_VALUE_1':'/dev/null'}
log=open(scratch/'reproduce.txt','w')
def run(cmd,cwd=repo):
 log.write('COMMAND '+repr(cmd)+'\n');log.flush();p=subprocess.run(cmd,cwd=cwd,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True);log.write(p.stdout+'\nExit: '+str(p.returncode)+'\n');log.flush()
 if p.returncode: raise RuntimeError(p.stdout)
 return p.stdout.strip()
try:
 run(['git','clone','--no-hardlinks',str(root),str(repo)],root)
 run(['git','config','commit.gpgsign','false']);run(['git','config','user.name','QA']);run(['git','config','user.email','qa@example.invalid'])
 recordpath=repo/'skills/testdata/owned-skill-versions.json'; original=json.loads(recordpath.read_text())
 def record(): run(['go','test','-count=1','-v','./skills','-run','^TestEveryOwnedSkillVersionIsRecorded$','-record-skill-versions'])
 def edit(file,text):
  for base in ['.agents/skills','skills']:
   p=repo/base/'roundfix'/file;p.write_text(p.read_text()+text)
 edit('references/deliver.md','\nQA item documentation addition.\n');record()
 new=json.loads(recordpath.read_text());last=new['skills']['roundfix'][-1]['version'];v=list(map(int,original['skills']['roundfix'][-1]['version'].split('.')));v[2]+=1;want='.'.join(map(str,v))
 assert last==want
 for name,entries in original['skills'].items():
  assert new['skills'][name][:len(entries)]==entries
  assert len(new['skills'][name])==len(entries)+(name=='roundfix')
 for base in ['.agents/skills','skills']:
  data=(repo/base/'roundfix/SKILL.md').read_text().split('---')[1];assert re.findall(r'^ *version: (.*)$',data,re.M)==[want,want]
 run(['go','test','-count=1','-v','./skills','-run','^TestEveryOwnedSkillVersionIsRecorded$'])
 recordbytes=recordpath.read_bytes();record();assert recordpath.read_bytes()==recordbytes
 log.write('Q02 PASS: both canonical/mirror version fields '+want+'; exactly one entry added; historical digests unchanged; fresh non-record check and idempotent recording pass.\n')
 # Disposable fixture reset, never the Run checkout.
 run(['git','reset','--hard','HEAD']);run(['git','checkout','-b','test/qa-collision-main'])
 run(['git','init','--bare',str(remote)],root);run(['git','remote','set-url','origin',str(remote)])
 base=run(['git','rev-parse','HEAD'])
 edit('references/review.md','\nQA default documentation addition.\n');record()
 edit('references/review.md','\nQA second default documentation addition.\n');record()
 default=json.loads(recordpath.read_text())['skills']['roundfix'][-1]['version']
 run(['git','add','.']);run(['git','commit','-m','test: default skill edits recorded']);run(['git','push','origin','test/qa-collision-main']);run(['git','symbolic-ref','refs/remotes/origin/HEAD','refs/remotes/origin/test/qa-collision-main']);run(['git','symbolic-ref','HEAD','refs/heads/test/qa-collision-main'],remote)
 run(['git','checkout','-b','test/qa-collision-item',base]);edit('references/deliver.md','\nQA item documentation addition.\n');record();run(['git','add','.']);run(['git','commit','-m','test: item skill edit recorded'])
 probe=scratch/'collision_probe_test.go';probe.write_text((evidence/'collision-probe.go.txt').read_text());overlay=scratch/'overlay.json';overlay.write_text(json.dumps({'Replace':{str(root/'internal/cli/qa_collision_probe_test.go'):str(probe)}}));env['QA_COLLISION_REPO']=str(repo)
 run(['go','test','-count=1','-v','-overlay',str(overlay),'./internal/cli','-run','^TestQAActualSkillCollision$'],root)
 merged=json.loads(recordpath.read_text())['skills']['roundfix'][-1]['version'];v=list(map(int,default.split('.')));v[2]+=1;assert merged=='.'.join(map(str,v))
 for b in ['.agents/skills','skills']:
  assert re.findall(r'^ *version: (.*)$',(repo/b/'roundfix/SKILL.md').read_text().split('---')[1],re.M)==[merged,merged]
 run(['go','test','-count=1','-v','./skills','-run','^TestEveryOwnedSkillVersionIsRecorded$']);assert run(['git','status','--porcelain'])==''
 log.write('Q03 PASS: default='+default+' merged='+merged+'; actual skill edits plus record collision; trailer asserted; clean merged tree record check passes.\n')
finally:
 log.close();shutil.copyfile(scratch/'reproduce.txt',evidence/'reproduce.txt');print('scratch:',scratch,flush=True)
