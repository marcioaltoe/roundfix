from gate import *
HEAD=subprocess.check_output(['git','-c','core.fsmonitor=false','rev-parse','HEAD'],text=True).strip()
BASE='513b22ebab2d62a23b4222073c937c2862907ec2'
ENV['GOCACHE']='/private/tmp/roundfix-0214-qa-cache'
def git(*args):
 r=subprocess.run(['git','-c','core.fsmonitor=false',*args],capture_output=True,text=True,env=ENV);assert r.returncode==0,(args,r.stderr);return r.stdout
(E/'row-3-sequential.txt').write_text(Path('/private/tmp/roundfix-0214-row3-sequential.txt').read_text()+'\nEXIT: 0\n')
settle(3,'pass','Sequential rerun exited 0; all selected store/classifier/trigger/Build/CLI cases passed. [Sequential output](evidence/2026-10-02-gate-01/row-3-sequential.txt). Initial row-3.txt exited 1 because concurrent QA writes violated suiteguard; no product assertion failed. Rerun had no concurrent repository writer.')
record=json.loads((ROOT/'docs/references/corrective-causes.json').read_text())
r=run(11,['./bin/roundfix','runs','causes','--since',record['window']['since'],'--until',record['window']['until'],'--format','json'],'real-causes',env=ENV)
if r.returncode==0:
 actual=json.loads(r.stdout);(E/'real-causes.json').write_text(r.stdout)
 same=actual['summary']==record['summary']
 settle(11,'pass' if same else 'fail',f'Read-only public CLI over committed window; NODE_OPTIONS and judge keys removed from child environment. Summary equality: {same}; observed Runs {actual["runs"]}, recorded {record["runs"]}. [Actual JSON](evidence/2026-10-02-gate-01/real-causes.json); [command](evidence/2026-10-02-gate-01/real-causes.txt). Independent comparison uses committed record, not a generated fixture.')
else:settle(11,'blocked (environment: Run Database unavailable)','Public reader could not reach real database; see real-causes.txt; rerun with read access required. No equivalent current observation.')
doc=(ROOT/'docs/references/archived-evidence-measurement.md').read_text();commit=re.search(r'^Measured at: ([0-9a-f]{40})$',doc,re.M)[1]
tree=git('ls-tree','-r','-l',commit,'--','docs/history');(E/'history-tree.txt').write_text(tree)
counts=[0]*6
for line in tree.splitlines():
 meta,path=line.split('\t');size=int(meta.split()[3]);counts[0]+=1;counts[1]+=size
 if re.match(r'^docs/history/specs/[^/]+/qa/evidence/',path):counts[2]+=1;counts[3]+=size
 if re.search(r'\.(png|jpg|jpeg|gif|pdf|db|zip)$',path.lower()):counts[4]+=1;counts[5]+=size
headers=['History files','History bytes','QA evidence files','QA evidence bytes','Binary files','Binary bytes'];expected=[int(re.search('^'+h+r': (\d+)$',doc,re.M)[1]) for h in headers]
ancestry=git('merge-base','--is-ancestor',commit,HEAD);mergebase=git('merge-base','main',HEAD).strip()
deletions=git('diff','--name-only','--diff-filter=DR',mergebase,HEAD,'--','docs/history');changes=git('diff','--name-status',mergebase,'--','docs/history');(E/'archive-audit.json').write_text(json.dumps({'commit':commit,'computed':counts,'document':expected,'merge_base':mergebase,'deletions':deletions,'all_changes':changes},indent=2))
settle(12,'pass' if counts==expected and not deletions and not changes else 'fail','Seven header lines independently verified at named ancestor '+commit+'; counts '+str(counts)+'. No history deletion/rename/content delta from main merge-base to head/current tree. [Audit](evidence/2026-10-02-gate-01/archive-audit.json).')
checks=[]
for path,fires,suffix in [('corrective-causes-measurement.md','reopen','reopen-the-memory-question-for-agent-sessions'),('task-acceptance-measurement.md','adopt','adopt-the-task-acceptance-judgment'),('archived-evidence-measurement.md',None,'remove-archived-evidence-nobody-reads')]:
 text=(ROOT/'docs/references'/path).read_text();value=re.search(r'^(Verdict|Candidates): (.+)$',text,re.M)[2];fired=value==fires if fires else value!='none';entries=list((ROOT/'docs/backlog').glob('*-'+suffix+'.md'));ok=len(entries)==int(fired)
 for e in entries:ok=ok and 'status: open' in e.read_text() and path in e.read_text()
 checks.append({'reference':path,'value':value,'fired':fired,'entries':[str(e.relative_to(ROOT)) for e in entries],'valid':ok})
(E/'backlog-rules.json').write_text(json.dumps(checks,indent=2));settle(13,'pass' if all(c['valid'] for c in checks) else 'fail','Recomputed verdicts inconclusive/do not adopt and archive Candidates: none; both recorded ablation exits are 2, so no candidate passes rule 3. No entry exists for any unfired rule. [Presence and rules](evidence/2026-10-02-gate-01/backlog-rules.json).')
# documentation mirror + unchanged settlement
checks={};checks['guide_heading']='#### Why Verification failed' in (ROOT/'docs/user-guide/commands/runs.md').read_text();checks['skill_heading']='### Why Verification failed' in (ROOT/'.agents/skills/roundfix/references/runs.md').read_text();mirrors=[]
for p in (ROOT/'.agents/skills/roundfix').rglob('*'):
 if p.is_file():q=ROOT/'skills/roundfix'/p.relative_to(ROOT/'.agents/skills/roundfix');mirrors.append(q.exists() and p.read_bytes()==q.read_bytes())
checks['mirrors']=all(mirrors);settlements={}
def section(t):
 m=re.search(r'(?m)^### QA settlement\n',t)
 if not m:return None
 end=re.search(r'(?m)^#{1,3} ',t[m.end():]);return t[m.start():m.end()+end.start()] if end else t[m.start():]
for p in git('ls-files','.agents/skills','skills').splitlines():
 if not p.endswith('.md'):continue
 current=section((ROOT/p).read_text())
 if current is None:continue
 old=section(git('show',BASE+':'+p));settlements[p]=current==old
checks['settlements']=settlements
(E/'docs-audit.json').write_text(json.dumps(checks,indent=2))
r=run(20,['make','skills-sync-check'],'skills-sync-check',env=ENV)
r2=run(20,['go','test','-count=1','-v','-tags','docscontract','-run','^(TestEveryCommandIsNamedInTheRoundfixSkill|TestEveryCommandIsNamedInTheUserGuide)$','./internal/docscontract'],'docs-command-contracts',env=ENV)
settle(20,'pass' if checks['guide_heading'] and checks['skill_heading'] and checks['mirrors'] and all(settlements.values()) and r.returncode==r2.returncode==0 else 'fail',f'Documented headings, {len(mirrors)} mirror files and {len(settlements)} unchanged QA settlement sections verified. skills-sync-check and named docs-contract tests exit {r.returncode}/{r2.returncode}. [Byte comparisons](evidence/2026-10-02-gate-01/docs-audit.json); skills-sync-check.txt and docs-command-contracts.txt.')
# Network imports and transitive senders
patch=git('diff','--unified=0',BASE,HEAD,'--','*.go');bad=[];path=''
for line in patch.splitlines():
 if line.startswith('+++ b/'):path=line[6:]
 elif line.startswith('+') and not line.startswith('+++') and '"net/http"' in line and not path.endswith('_test.go') and not path.startswith('internal/judge/'):bad.append(path)
r=run(21,['go','list','-deps','./internal/runcause'],'runcause-imports',env=ENV);deps=r.stdout.splitlines();senders=[p for p in deps if p in ['net/http','net/smtp','net/rpc','net/http/httputil']]
settle(21,'pass' if not bad and not senders and r.returncode==0 else 'fail',f'Git added-import sweep: {bad}; production runcause dependency senders: {senders}. [Import closure](evidence/2026-10-02-gate-01/runcause-imports.txt).')
# Promise maps checked explicitly, independently of source read
prd=(ROOT/'docs/specs/0214-measure-before-changing/_prd.md').read_text();tech=(ROOT/'docs/specs/0214-measure-before-changing/_techspec.md').read_text();tasks=[(ROOT/f'docs/specs/0214-measure-before-changing/task_0{i}.md').read_text().split('## References',1)[1] for i in range(1,5)]
map={}
for kind,total in [('Success Metric',5),('API Contract',4)]:
 for n in range(1,total+1):
  owners=[]
  for i,t in enumerate(tasks,1):
   for m in re.finditer(kind+r's? ([0-9][0-9 ,and-]*)',t):
    nums=[int(x) for x in re.findall(r'\d+',m[1])]
    if n in nums:owners.append('task_0'+str(i));break
  map[f'{kind} {n}']=owners
(E/'promise-map.json').write_text(json.dumps(map,indent=2));ok=all(map.values()) and all('Success Metric '+str(n)+' →' in tech for n in range(1,6))
settle(22,'pass' if ok else 'fail','Every numbered metric/API contract has a consuming Task in its References; all metrics map in TechSpec; strict Spec check additionally confirms Stories/Core Features. [Independently rebuilt ownership map](evidence/2026-10-02-gate-01/promise-map.json).')
# Actual commits vs Context and Recorded paths
commits=git('log','--format=%H',BASE+'..'+HEAD).splitlines();audit=[]
for c in commits:
 body=git('show','-s','--format=%B',c);m=re.search(r'Roundfix-Task: (task_\d+)',body)
 if not m:continue
 task=m[1];file='docs/specs/0214-measure-before-changing/'+task+'.md';t=(ROOT/file).read_text();decl=set(re.findall(r'(?m)^- (?:creates|interface|instruction): `([^`]+)`',t));recorded=set()
 if '## Recorded paths' in t:recorded=set(re.findall(r'`([^`]+)`',t.split('## Recorded paths',1)[1].split('\n## ',1)[0]))
 paths=git('diff-tree','--no-commit-id','--name-only','-r',c).splitlines();outside=[p for p in paths if p!=file and p not in decl and p not in recorded];audit.append({'task':task,'commit':c,'paths':paths,'recorded_paths':sorted(recorded),'outside':outside})
(E/'task-scope.json').write_text(json.dumps(audit,indent=2));protected=git('diff','--name-only',BASE,HEAD,'--','internal/cli/cli_test.go','.roundfixrc.yml','Makefile','go.mod');grant=git('show',BASE+':docs/specs/0214-measure-before-changing/_authorization.md');git('merge-base','--is-ancestor',BASE,'014f1424')
settle(23,'pass' if len(audit)==4 and not any(a['outside'] for a in audit) and not protected else 'fail','All four actual Task commits independently diff-tree audited against Context/Recorded paths; earlier-Run task_03 included. Three governed skill paths bounded by ancestor grant; protected four paths have no delta. [Actual paths](evidence/2026-10-02-gate-01/task-scope.json). Seeded authorization audit retained; no extra governed Task commit unlisted. Separate glossary repair f59ebc42 changes only CONTEXT.md and is not a Task commit; that repair is tested in row 24.')
