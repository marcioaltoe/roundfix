from pathlib import Path
import subprocess,re,json
s=Path(__file__).resolve().parents[3]; e=Path(__file__).resolve().parent;r=s/'qa/qa-report-2026-10-01.md'
def row(i,status,detail):
 t=r.read_text().replace('| '+i+' | pending |','| '+i+' | '+status+' |');t=t.replace('### '+i+'\n','### '+i+'\n\n'+detail+'\n');r.write_text(t)
def run(i,paths,names):
 regex='^('+'|'.join(sorted(set(names)))+')$';cmd=['rtk','proxy','go','test','-count=1','-v','-run',regex]+paths
 p=subprocess.run(cmd,capture_output=True,text=True);output=p.stdout+p.stderr
 missing=[n for n in set(names) if '--- PASS: '+n not in output]
 log=e/(i+'-tests.txt');log.write_text('command: '+json.dumps(cmd)+'\n'+output+'\nexit: '+str(p.returncode)+'\nmissing passes: '+str(missing)+'\n')
 row(i,'pass' if p.returncode==0 and not missing else 'fail',f'Observed: {len(set(names))} named tests executed; exit {p.returncode}; missing passes {missing}. Evidence: `qa/evidence/2026-10-01-gate/{log.name}`. Tests assert observable boundary calls and persisted state, with real temporary Git repositories where applicable. No repository Verification rerun.')
 print(i,p.returncode,'named:',len(set(names)),'missing:',missing,flush=True)
for i,n in [('04',1),('05',2),('06',3),('07',4)]:
 text=(s/f'task_0{n}.md').read_text().split('## Verification',1)[1].split('## References',1)[0]
 names=re.findall(r'Test[A-Za-z0-9_]+',text)
 if i=='04':
  names+=['TestGitHubCLIBuildFailuresAreUnattributable','TestGitHubCLISetupFailuresAreUnattributable','TestGitHubCLILogsWithoutGoPackagesAreUnattributable','TestMultipleFailedJobsInOneRunAreRerunTogetherOnce','TestAnOldCheckAttemptCannotKeepResettingTheTimeout']
 if i=='07':
  names+=[m for f in Path('internal/cli').glob('deliver_conflict_test.go') for m in re.findall(r'func (Test[A-Za-z0-9_]+)',f.read_text())]
 run(i,['./internal/spec','./internal/config','./internal/delivery','./internal/cli','./skills'],names)
run('03-test',['./internal/cli'],['TestDeliverStatusReproducesSurfaceTranscriptOne','TestDeliverStatusPrintsNoWarningLineWithoutAWarning'])
