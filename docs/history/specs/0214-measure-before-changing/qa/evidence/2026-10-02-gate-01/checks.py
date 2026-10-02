from gate import *
HEAD=subprocess.check_output(['git','-c','core.fsmonitor=false','rev-parse','HEAD'],text=True).strip()
BASE='513b22ebab2d62a23b4222073c937c2862907ec2'
ENV['GOCACHE']='/private/tmp/roundfix-0214-qa-cache'
def check(n,args,detail,label=None):
 r=run(n,args,label,env=ENV)
 settle(n,'pass' if r.returncode==0 else 'fail',detail+f' Exit {r.returncode}; [capture](evidence/2026-10-02-gate-01/{label or "row-"+str(n)}.txt).')
 print(n,r.returncode,flush=True)
 return r
settle(1,'pass',f'Daemon-provided make verify-changed: pass, exit 0; diagnostics removed on success. Recorded at audited head {HEAD}; not rerun.')
r=run(1,['./bin/roundfix','--version'],'user-flow-version',env=ENV)
r=run(1,['./bin/roundfix','spec','check','0214-measure-before-changing','--strict'],'spec-check',env=ENV)
assert r.returncode==0
tech=(ROOT/'docs/specs/0214-measure-before-changing/_techspec.md').read_text()
checks={}
for section,path in [('The signature table','internal/runcause/signatures.json'),('The Task acceptance question','internal/judge/testdata/task-acceptance-question.json')]:
 block=tech.split('### '+section,1)[1].split('```json\n',1)[1].split('```',1)[0]
 checks[path]={'equal':block.encode()==(ROOT/path).read_bytes(),'sha256':hashlib.sha256((ROOT/path).read_bytes()).hexdigest()}
(E/'embedded-contracts.json').write_text(json.dumps(checks,indent=2))
settle(2,'pass' if all(x['equal'] for x in checks.values()) else 'fail','Both JSON blocks independently extracted and byte-compared; [digests and equalities](evidence/2026-10-02-gate-01/embedded-contracts.json).')
check(3,['go','test','-count=1','-v','-run','Test(RunEventsOfKinds|LoadRefuses|Classify|Trigger|Build|RunsCauses)','./internal/store','./internal/runcause','./internal/cli'],'All named store, classifier, trigger, Build and command boundary cases; includes unsafe logs, corrective/QA/missing graphs and home SHA-256 preservation.')
check(9,['go','test','-count=1','-v','-run','Test(TaskAcceptance|AUROC|ClusterBootstrap|MeasureTaskAcceptance)','./internal/judge'],'Fake transport proves Result/source exclusion, endpoint key isolation, generic-key no-request block, pin refusal, per-request log and ceiling. Live harness must SKIP without its flag.')
check(10,['go','test','-count=1','-v','-run','^TestCausesRecordIsConsistent$','./internal/runcause','-args','-causes-record='+str(ROOT/'docs/references/corrective-causes.json'),'-causes-document='+str(ROOT/'docs/references/corrective-causes-measurement.md')],'Cause record recomputed from items; document verdict inconclusive (100/132 unclassified).','causes-consistency')
check(10,['go','test','-count=1','-v','-run','^TestTaskAcceptanceRecordIsConsistent$','./internal/judge','-args','-task-acceptance-record='+str(ROOT/'docs/references/task-acceptance-remeasurement.json'),'-task-acceptance-document='+str(ROOT/'docs/references/task-acceptance-measurement.md'),'-task-acceptance-labels='+str(ROOT/'docs/references/corrective-causes.json')],'Both committed-record consistency tests recomputed figures/digests/state hashes and passed. Verdicts: inconclusive; do not adopt (237 answered, 31 repaired, AUROC 0.466881, 95% interval 0.358034–0.574594). Cause capture: causes-consistency.txt.','acceptance-consistency')
