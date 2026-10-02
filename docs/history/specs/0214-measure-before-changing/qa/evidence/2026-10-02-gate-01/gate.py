from pathlib import Path
import subprocess,json,os,re,hashlib
ROOT=Path.cwd()
REPORT=ROOT/'docs/specs/0214-measure-before-changing/qa/qa-report-2026-10-02-01.md'
E=REPORT.parent/'evidence/2026-10-02-gate-01'
def settle(n,status,detail):
 s=REPORT.read_text();a=s.index('## Results');b=s.index('## Mechanical skips',a);s=s[:a]+re.sub(rf'(?m)^\| {n} \| [^|]+ \|',f'| {n} | {status} |',s[a:b])+s[b:]
 a=s.index(f'### {n}\n');b=s.find('\n### ',a+1);b=len(s) if b<0 else b
 block=s[a:b];block=re.sub(r'Evidence: .*',lambda m:'Evidence: '+detail,block,flags=re.S)
 REPORT.write_text(s[:a]+block+s[b:])
def run(n,args,label=None,cwd=ROOT,env=None):
 r=subprocess.run(args,cwd=cwd,env=env,capture_output=True,text=True)
 name=label or f'row-{n}';(E/(name+'.txt')).write_text('$ '+ ' '.join(args)+'\n'+r.stdout+'\nSTDERR:\n'+r.stderr+f'\nEXIT: {r.returncode}\n')
 return r
EXCLUDED={'NODE_OPTIONS','ROUNDFIX_OPENROUTER_API_KEY','ROUNDFIX_TYPESAFE_API_KEY','OPENROUTER_API_KEY','TYPESAFE_API_KEY'}
ENV={k:os.environ[k] for k in os.environ if k not in EXCLUDED}
