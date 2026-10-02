import os,pathlib,subprocess,time,signal,json
root=pathlib.Path.cwd(); out=pathlib.Path('/private/tmp/roundfix-0213-qa'); tmp=out/'interrupt-tmp'; tmp.mkdir(exist_ok=True);record=out/'interrupt-owner-pid'; binary=str(out/'cli.test');env=os.environ.copy();env['TMPDIR']=str(tmp);env['ROUNDFIX_DETACH_TEST_OWNER_RECORD']=str(record)
def snapshot():
 text=subprocess.check_output(['ps','-axo','pid=,ppid=,pgid=,command='],text=True)
 return [line.strip() for line in text.splitlines() if binary in line]
log=(out/'interrupt.log').open('w');cmd=[binary,'-test.run=^TestRunImplementDetachSurvivesCallerProcessGroupKill$','-test.timeout=180s','-test.v'];p=subprocess.Popen(cmd,cwd=root/'internal/cli',env=env,stdout=log,stderr=subprocess.STDOUT,start_new_session=True); result={'command':cmd,'cwd':str(root/'internal/cli'),'binary_pid':p.pid,'TMPDIR':str(tmp)}; owned=set([p.pid]);start=time.monotonic()
try:
 while not record.exists() or not list(tmp.rglob('prompt-started')):
  if p.poll() is not None:raise RuntimeError('binary ended before prompt: '+str(p.returncode))
  if time.monotonic()-start>90:raise RuntimeError('prompt rendezvous deadline')
  time.sleep(.05)
 result['owner_pid']=int(record.read_text());owned.add(result['owner_pid']);result['prompt_started']=[str(x) for x in tmp.rglob('prompt-started')];result['before_kill']=snapshot();owned.update(int(x.split()[0]) for x in result['before_kill']);os.kill(p.pid,signal.SIGKILL);result['binary_exit']=p.wait(timeout=10);killed=time.monotonic();observations=[]
 while True:
  current=snapshot();observations.append({'seconds':round(time.monotonic()-killed,3),'processes':current});owned.update(int(x.split()[0]) for x in current)
  if not current or time.monotonic()-killed>=60:break
  time.sleep(.1)
 result['observations']=observations;result['remaining']=current;result['pass']=not current;result['cleanup_pids']=[]
finally:
 if p.poll() is None:os.kill(p.pid,signal.SIGKILL);p.wait()
 # Reclaim only PID records or exact compiled-binary command matches observed by this probe.
 for line in snapshot():
  pid=int(line.split()[0])
  if pid in owned:
   try:os.kill(pid,signal.SIGKILL);result.setdefault('cleanup_pids',[]).append(pid)
   except ProcessLookupError:pass
 log.close();(out/'interrupt-summary.json').write_text(json.dumps(result,indent=2))
print(json.dumps(result,indent=2));raise SystemExit(0 if result.get('pass') else 1)
