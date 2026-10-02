from pathlib import Path
import subprocess,sys,os,re
root=Path(__file__).resolve().parents[6]
e=Path(__file__).parent
report=e.parents[1]/'qa-report-2026-10-02-01.md'
def settle(row,status,detail):
 s=report.read_text(); s=re.sub(r'(?m)^\| '+row+r' \| pending \|', '| '+row+' | '+status+' |',s)
 a=s.index('### '+row+'\n'); b=s.find('\n### ',a+1); b=len(s) if b<0 else b
 part=s[a:b].replace('Pending execution.',detail); report.write_text(s[:a]+part+s[b:])
def run(name,cmd,env=None):
 p=subprocess.run(cmd,cwd=root,env=env,text=True,capture_output=True)
 (e/name).write_text('command: '+repr(cmd)+'\nexit: '+str(p.returncode)+'\nstdout:\n'+p.stdout+'\nstderr:\n'+p.stderr)
 print(name,p.returncode); return p
if __name__=='__main__':
 env=dict(os.environ,GOCACHE='/private/tmp/roundfix-0211-qa-gocache')
 row=sys.argv[1]
 cmds={
 '02':[('02-run-start.txt',['go','test','-count=1','-v','./internal/cli','-run','RunStart|ArchivedRetry|OperatorArchiveHistory']),('02-operator-archive.txt',['go','test','-count=1','-v','./internal/delivery','-run','OperatorArchive|ArchivedHead|EnvironmentOnlyPartial'])],
 '03':[('03-delivery.txt',['go','test','-count=1','-v','./internal/delivery'])],
 '04':[('04-retry-tests.txt',['go','test','-count=1','-v','./internal/cli','-run','DeliverRetry'])],
 '05':[('05-agent-tests.txt',['go','test','-count=1','-v','./internal/agent','-run','NodeOptions|MissingPreload|ACPXCommandEnv|CommandEnvDefaults']),('05-doctor-tests.txt',['go','test','-count=1','-v','./internal/cli','-run','Doctor.*Preload|Setup.*Preload'])]}
 results=[run(n,c,env) for n,c in cmds[row]]
 assert all(p.returncode==0 for p in results)
