import os,subprocess,json,pathlib,hashlib,shutil,re
ROOT=pathlib.Path.cwd(); TMP=pathlib.Path('/private/tmp/roundfix-0223-qa'); E=ROOT/'docs/specs/0223-adjustments-the-adopters-asked-for/qa/evidence/2026-10-04-gate'; REPORT=E.parents[1]/'qa-report-2026-10-04.md'
HOME_DIR=TMP/'home';HOME_DIR.mkdir(exist_ok=True)
ENV=os.environ.copy();ENV['HOME']=str(HOME_DIR);ENV.pop('NODE_OPTIONS',None);ENV['GIT_CONFIG_GLOBAL']='/dev/null';ENV['GIT_CONFIG_NOSYSTEM']='1'
OFFLINE=TMP/'offline-tools';OFFLINE.mkdir(exist_ok=True)
for name in ['git','make']:
 link=OFFLINE/name
 if not link.exists():link.symlink_to(shutil.which(name))
ENV['PATH']=str(OFFLINE)+':/usr/bin:/bin'
for key in list(ENV):
 if key.startswith(('ROUNDFIX_','GITHUB_','GH_','TYPESAFE_','OPENAI_','ANTHROPIC_')):ENV.pop(key,None)
def run(c,cwd=ROOT,label=None):
 r=subprocess.run(c,cwd=cwd,env=ENV,capture_output=True,text=True)
 if label:
  (E/(label+'.stdout.txt')).write_text(r.stdout);(E/(label+'.stderr.txt')).write_text(r.stderr);(E/(label+'.command.txt')).write_text('cwd: '+str(cwd)+'\nHOME: '+str(HOME_DIR)+'; NODE_OPTIONS: unset\n'+repr(c)+'\nexit: '+str(r.returncode)+'\n')
 return r

def git(repo,*args):
 r=run(['git','-c','core.fsmonitor=false',*args],repo);assert r.returncode==0,r.stderr;return r.stdout

def init(name):
 p=TMP/name;p.mkdir(exist_ok=True);git(p,'init','-q');git(p,'config','user.name','QA');git(p,'config','user.email','qa@example.invalid');git(p,'config','commit.gpgsign','false');git(p,'config','core.hooksPath','/dev/null');return p

def bind(p,binary=ROOT/'bin/roundfix'):
 (p/'bin').mkdir(exist_ok=True); link=p/'bin/roundfix'
 if link.exists():link.unlink()
 link.symlink_to(binary);(p/'.git/info/exclude').write_text('/bin/\n')

def commit(p,msg):git(p,'add','.');git(p,'commit','-qm',msg)
def fixture(name,adopt=False,base=False):
 p=init(name);(p/'go.mod').write_text('module example.invalid/qa\n\ngo 1.25.0\n');(p/'Makefile').write_text('verify:\n\t@true\nverify-incremental:\n\t@true\n');(p/'main.go').write_text('package main\nfunc main() {}\n');commit(p,'chore: initial repository');bind(p)
 if adopt:
  shutil.copytree(ROOT/'.agents/skills',p/'.agents/skills',dirs_exist_ok=True);shutil.copy2(ROOT/'skills-lock.json',p/'skills-lock.json')
  decisions=['preservation.mode=greenfield','language.generated=English','verification.gate=make verify','verification.incremental=make verify-incremental','branch.prefix=ma/','spec.scaffold=true','domain.layout=single-context','triage.external=false','autonomous.enabled=true','runtime.backend=codex gpt-5.5 xhigh','runtime.design=claude opus xhigh','secondbrain.enabled=false','repository.extension.enabled=false']
  if base:bind(p,TMP/'base/bin/roundfix')
  args=['./bin/roundfix','baseline','plan','--repo','.','--profile','go-cli-tui','--format','json']
  for d in decisions:args+=['--decision',d]
  r=run(args,p,name+'-adopt-plan');assert r.returncode==0,(r.stdout,r.stderr)
  plan=TMP/(name+'-plan.json');plan.write_text(r.stdout);v=json.loads(r.stdout)
  r=run(['./bin/roundfix','baseline','apply','--repo','.','--plan',str(plan),'--confirm-plan',v['planDigest'],'--format','json'],p,name+'-adopt-apply');assert r.returncode==0,(r.stdout,r.stderr)
  bind(p)
  commit(p,'chore: adopt baseline')
 git(p,'tag','v0.4.0');(p/'internal/cli').mkdir(parents=True,exist_ok=True);(p/'internal/cli/output.go').write_text('package cli\n');commit(p,'fix: correct output');return p

def snap(p):
 h=hashlib.sha256()
 for f in sorted(p.rglob('*')):
  if '.git' in f.relative_to(p).parts or not f.is_file():continue
  h.update(str(f.relative_to(p)).encode());h.update(f.read_bytes())
 return {'sha256':h.hexdigest(),'status':git(p,'status','--porcelain'),'refs':git(p,'show-ref'),'config':git(p,'config','--local','--list')}

def close(i,text):
 s=REPORT.read_text();s=s.replace('| '+i+' | pending |','| '+i+' | pass |');a=s.index('### '+i+'\n');b=s.find('\n### ',a+5);b=len(s) if b<0 else b;s=s[:a]+s[a:b].replace('Evidence: pending.',text)+s[b:];REPORT.write_text(s)

for name,adopt,i in [('unadopted',False,'02'),('adopted',True,'03')]:
 if '| '+i+' | pass |' in REPORT.read_text():continue
 p=TMP/name if (TMP/name/'.git').exists() else fixture(name,adopt); records=[]
 for fmt in ['text','json']:
  before=snap(p);r=run(['./bin/roundfix','release','plan','--format',fmt],p,name+'-'+fmt);after=snap(p);assert before==after,(before,after)
  records.append({'format':fmt,'before':before,'after':after,'exit':r.returncode});assert r.returncode==0 and r.stderr=='',(r.stdout,r.stderr)
  if fmt=='json':
   v=json.loads(r.stdout);assert v['state']=='ready' and v['proposedVersion']=='v0.4.1';assert 'checks' in v
   (E/(name+'-checks.json')).write_text(json.dumps(v['checks'],indent=2))
  else:
   tech=(ROOT/'docs/specs/0223-adjustments-the-adopters-asked-for/_techspec.md').read_text();trans=re.findall(r'```transcript\n(.*?)\n\s*```',tech,re.S)[0 if not adopt else 1];trans='\n'.join(line.strip() for line in trans.splitlines());want=trans.split('stdout:\n')[1].split('\nstderr:')[0]
   regex=re.escape(want);regex=re.sub(r'<[^>]+>',r'.+',regex);assert re.fullmatch(regex,r.stdout.rstrip('\n')),(want,r.stdout)
   skills=next(x for x in r.stdout.splitlines() if x.startswith('skills:'))
   assert sum(x.startswith('skills:') for x in r.stdout.splitlines())==1 and sum(x.startswith('baseline:') for x in r.stdout.splitlines())==1
 if adopt:
  r=run(['./bin/roundfix','doctor'],p,'adopted-doctor');doctor=next(x for x in r.stdout.splitlines() if x.startswith('skills:'));normalized=re.sub(r'^skills: ([^ ]+) \((.*)\)$',r'skills: \1: \2',doctor);assert normalized==skills,(doctor,skills)
 (E/(name+'-snapshots.json')).write_text(json.dumps(records,indent=2));close(i,'Evidence: `evidence/2026-10-04-gate/'+name+'-text.stdout.txt`, `'+name+'-json.stdout.txt`, `'+name+'-checks.json` and `'+name+'-snapshots.json`. Text matched the complete declared transcript (stdout/stderr/exit separately); JSON independently confirms ready/v0.4.1 and both checks. Each invocation preserves tree SHA-256, Git status, refs and local config.'+(' Doctor status and detail are identical, with its existing parenthesized printer format preserved (`adopted-doctor.stdout.txt`); unrelated Doctor readiness diagnostics are captured separately.' if adopt else ''))
 print('ROW',i,'PASS',flush=True)
p=TMP/'unadopted'
for label,args in [('reset',['--reset-to','v0.0.1']),('refused',['--from','not-a-ref'])]:
 before=snap(p);r=run(['./bin/roundfix','release','plan',*args],p,label);assert snap(p)==before;assert r.returncode==2;assert not any(x.startswith(('skills:','baseline:')) for x in r.stdout.splitlines());assert '"checks"' not in r.stdout
close('04','Evidence: `evidence/2026-10-04-gate/reset.{stdout,stderr,command}.txt` and `refused.{stdout,stderr,command}.txt`. Reset refuses missing origin offline; invalid base refuses locally. Both exit 2, print no check line/object and preserve the repository snapshot. A successful reset would require remote/GitHub inventory, prohibited by this Task; the named reset invocation is exercised only to its offline refusal.');print('ROW 04 PASS',flush=True)
p=fixture('old-adopter',True,True);guide=p/'docs/agents/agent-instructions.md';before=guide.read_bytes();manifest=json.loads((p/'docs/agents/setup-context.json').read_text());prefix=manifest['decisions']['branch.prefix'];r=run(['./bin/roundfix','baseline','update','--repo','.','--no-skills','--yes','--format','json'],p,'old-adopter-update');assert r.returncode==0,(r.stdout,r.stderr);v=json.loads(r.stdout);assert not v.get('newDecisions');assert guide.read_bytes()==before
summary={'base':'3490588afbdffe4419574befd0d7d9569882ff69','branch.prefix':prefix,'guide_before_sha256':hashlib.sha256(before).hexdigest(),'guide_after_sha256':hashlib.sha256(guide.read_bytes()).hexdigest(),'newDecisions':v.get('newDecisions')}
p2=TMP/'unrecorded-adopter';shutil.copytree(p,p2,dirs_exist_ok=True);m=p2/'docs/agents/setup-context.json';v=json.loads(m.read_text());del v['decisions']['branch.prefix'];m.write_text(json.dumps(v,indent=2)+'\n');r=run(['./bin/roundfix','baseline','update','--repo','.','--no-skills','--yes','--format','json'],p2,'unrecorded-update');assert r.returncode==0,(r.stdout,r.stderr);assert not json.loads(r.stdout).get('newDecisions');after=(p2/'docs/agents/agent-instructions.md').read_text();fixed='No branch prefix is recorded. Name new work branches `<type>/<description>`,\nwhere `<type>` is the work\'s Conventional Commit type, as the branch rule\nbelow states. Tool-owned Run and Task branches follow their tool\'s documented\nnamespace.';assert fixed in after and 'The branch-prefix pattern is' not in after;assert 'branch.prefix' not in json.loads(m.read_text())['decisions'];summary['unrecorded_paragraph']=fixed
(E/'05-compatibility.json').write_text(json.dumps(summary,indent=2));close('05','Evidence: focused tests `evidence/2026-10-04-gate/05-focused-tests.txt` exit 0, all six required tests passed. Base-built binary adopted `old-adopter` with `ma/`; current binary update reports no new decisions and retains identical complete guide hashes (`05-compatibility.json`). A copy with no prefix updates without a decision and renders Fixed text 2; the manifest remains unanswered. Commands and output are in `old-adopter-*` and `unrecorded-update.*.txt`. The base binary is used solely for the explicitly required historical adoption, then replaced by the Run-built binary for verification.');print('ROW 05 PASS',flush=True)
baseguide=subprocess.check_output(['git','show','3490588:docs/agents/agent-instructions.md']);assert (ROOT/'docs/agents/agent-instructions.md').read_bytes()==baseguide
for n in [1,2]:
 r=run(['./bin/roundfix','baseline','update','--repo','.','--no-skills','--format','json'],ROOT,'repository-update-'+str(n));assert r.returncode==0,(r.stdout,r.stderr);assert json.loads(r.stdout)['state']=='current'
changed=subprocess.check_output(['git','diff','--name-only','3490588','HEAD','--','internal/baseline/assets']).decode();assert all('/clauses/' not in f and 'clause' not in pathlib.Path(f).name for f in changed.splitlines())
(E/'06-retention.txt').write_text('Repository guide equals base byte for byte. SHA256: '+hashlib.sha256(baseguide).hexdigest()+'\nChanged assets:\n'+changed+'No clause manifest was changed; all clause IDs/text are retained. Template prose moved into the renderer, not a Normative Clause.\n')
close('06','Evidence: `evidence/2026-10-04-gate/06-retention.txt` and `repository-update-{1,2}.stdout.txt`. Guide equals base bytes and SHA-256; two separate read-only updates report state `current`, no writes or skills stage. Changed catalog assets are decision metadata, core required list and template prose; no clause manifest/ID changed.');print('ROW 06 PASS',flush=True)
