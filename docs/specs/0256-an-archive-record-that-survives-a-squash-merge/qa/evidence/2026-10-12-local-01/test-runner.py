from pathlib import Path
import os,subprocess,json,re,sys
root=Path.cwd().resolve();e=root/'docs/specs/0256-an-archive-record-that-survives-a-squash-merge/qa/evidence/2026-10-12-local-01';selected=json.loads((e/'test-selection.json').read_text());results=[]
for pkg,names in selected.items():
 cmd=['go','test','-json','-count=1']
 if pkg=='internal/docscontract':cmd+=['-tags','docscontract']
 cmd+=['-run','^('+'|'.join(names)+')$','./'+pkg]
 env=os.environ.copy();env.update(GOPROXY='off',GOSUMDB='off')
 print('RUN',pkg,len(names),flush=True)
 # Keep all evidence writers idle until this package and its suiteguard exit.
 r=subprocess.run(cmd,cwd=root,env=env,capture_output=True,text=True,timeout=300)
 log=e/(pkg.replace('/','-')+'.jsonl');log.write_text(r.stdout);(e/(pkg.replace('/','-')+'.stderr.txt')).write_text(r.stderr)
 events=[json.loads(x) for x in r.stdout.splitlines() if x.startswith('{')]
 passed={x.get('Test') for x in events if x.get('Action')=='pass'};failed=[x.get('Test','package') for x in events if x.get('Action')=='fail'];missing=sorted(set(names)-passed)
 item={'package':pkg,'command':cmd,'exit':r.returncode,'expected':names,'missing':missing,'failures':failed,'passed_top_level':len(set(names)&passed),'events':[x for x in events if x.get('Action') in ('pass','fail','skip')]};results.append(item);(e/'test-results.json').write_text(json.dumps(results,indent=2));print('DONE',pkg,'exit',r.returncode,'passed',item['passed_top_level'],'missing',missing,'failed',failed,flush=True)
print('TOTAL',sum(x['passed_top_level'] for x in results),flush=True)
sys.exit(0 if all(x['exit']==0 and not x['missing'] and not x['failures'] for x in results) else 1)
