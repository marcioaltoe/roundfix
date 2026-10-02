import concurrent.futures as cf, json, os, pathlib, subprocess, time
root=pathlib.Path.cwd(); out=pathlib.Path('/private/tmp/roundfix-0213-qa'); env=os.environ.copy(); env['GOCACHE']='/private/tmp/roundfix-0213-qa-cache'
packages=['store','daemon','cli','testfixture','baseline']
for pkg in packages:
 with (out/(pkg+'-compile.log')).open('w') as log:
  rc=subprocess.call(['go','test','-c','-o',str(out/(pkg+'.test')),'./internal/'+pkg],env=env,stdout=log,stderr=subprocess.STDOUT)
 print('compile',pkg,rc,flush=True)
 if rc: raise SystemExit(rc)
def run(label,pkg,pattern,count=1,cpu='1,4',timeout='600s'):
 path=out/(label+'.log'); start=time.monotonic()
 command=[str(out/(pkg+'.test')),'-test.run='+pattern,'-test.count='+str(count),'-test.cpu='+cpu,'-test.timeout='+timeout,'-test.v']
 with path.open('w') as log: rc=subprocess.call(command,cwd=root/'internal'/pkg,env=env,stdout=log,stderr=subprocess.STDOUT)
 data=path.read_text(); result={'label':label,'command':command,'cwd':str(root/'internal'/pkg),'exit':rc,'seconds':round(time.monotonic()-start,3),'passes':data.count('--- PASS:'),'failures':data.count('--- FAIL:')}
 (out/(label+'.json')).write_text(json.dumps(result,indent=2));return result
# Each test command starts its suite guard only after report planning is complete.
with cf.ThreadPoolExecutor(max_workers=4) as pool:
 tasks=[pool.submit(run,'fixtures','cli','^TestScriptFixtureIsTheCompiledTestBinary$',1),pool.submit(run,'residue','testfixture','^TestNoTestWritesAnExecutableOutsideTheResidue$',1),pool.submit(run,'death','cli','^(TestImplementDetachChildEndsWhenItsTestBinaryDies|TestDetachSurvivorEndsWhenItsTestBinaryDies|TestRunImplementDetachSurvivesCallerProcessGroupKill|TestDetachedChildIsTerminatedAtTeardown)$',1),pool.submit(run,'assets','baseline','^(TestAssetsSync|TestBaselineAssetsSync)',3)]
 for task in cf.as_completed(tasks):print(json.dumps(task.result()),flush=True)
def stress_worker(worker,pkg,pattern,rounds,count):
 return [run(pkg+'-w'+str(worker)+'-r'+str(r),pkg,pattern,count) for r in range(rounds)]
with cf.ThreadPoolExecutor(max_workers=8) as pool:
 results=list(pool.map(lambda w:stress_worker(w,'store','^TestBatchClosesOnCountLingerAndImmediate$/^linger_closes_a_quiet_batch$',10,50),range(8)))
(out/'linger-summary.json').write_text(json.dumps(results,indent=2));print('linger commands',80,'failed',sum(x['exit']!=0 for xs in results for x in xs),flush=True)
# Ten burners self-terminate at ten minutes even if the coordinator fails.
burners=[subprocess.Popen(['python3','-c','import time; end=time.monotonic()+600; x=1\nwhile time.monotonic()<end: x=(x*1664525+1013904223)%4294967296'],stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL) for _ in range(10)]
(out/'burner-pids.json').write_text(json.dumps([p.pid for p in burners]))
try:
 with cf.ThreadPoolExecutor(max_workers=8) as pool:
  results=list(pool.map(lambda w:stress_worker(w,'daemon','^(TestTaskBudgetReasonNamesTheSettlementThatRenewedIt|TestTaskBudgetCancelsAStalledTaskOneAllowanceAfterTheLastSettlement)$',6,20),range(8)))
 (out/'budget-summary.json').write_text(json.dumps(results,indent=2));print('budget commands',48,'failed',sum(x['exit']!=0 for xs in results for x in xs),flush=True)
finally:
 for p in burners:
  if p.poll() is None: p.terminate()
  p.wait()
print('ALL TEST PROBES FINISHED',flush=True)
