from pathlib import Path
import re,subprocess,os,json
s=Path('docs/specs/0229-a-jev-router-that-checks-credit-and-names-its-model');e=s/'qa/evidence/2026-10-05-jev-credit'
env=os.environ.copy();env['GOCACHE']='/private/tmp/roundfix-0229-qa-gocache';env['TMPDIR']='/private/tmp/roundfix-0229-qa-tmp';Path(env['TMPDIR']).mkdir(exist_ok=True)
for name in ['ROUNDFIX_OPENROUTER_API_KEY','ROUNDFIX_TYPESAFE_API_KEY','TYPESAFE_API_KEY']:env.pop(name,None)
for i in (2,3,4):
 text=(s/f'task_0{i}.md').read_text();match=re.search(r'go test -count=1 -v -run "([^"]+)" ((?:\./[\w/]+\s*)+)',text);assert match
 cmd=['go','test','-count=1','-v','-run',match[1]]+match[2].strip().split()
 p=subprocess.run(cmd,env=env,text=True,capture_output=True)
 (e/f'task-0{i}-tests-full.txt').write_text('COMMAND: '+repr(cmd)+'\n'+p.stdout+p.stderr+'\nexit: '+str(p.returncode)+'\n')
 names=match[1].strip('^()$').split('|');missing=[n for n in names if '--- PASS: '+n not in p.stdout]
 print('task',i,'exit',p.returncode,'required',len(names),'missing',missing,flush=True)
 print((p.stdout+p.stderr)[-400:],flush=True)
 r=s/'qa/qa-report-2026-10-05.md';t=r.read_text();t=re.sub(r'\| '+str(i)+r' \| (?:pending|fail|pass|blocked[^|]+) \|','| '+str(i)+' | '+('pass' if p.returncode==0 and not missing else 'fail')+' |',t);r.write_text(t)
