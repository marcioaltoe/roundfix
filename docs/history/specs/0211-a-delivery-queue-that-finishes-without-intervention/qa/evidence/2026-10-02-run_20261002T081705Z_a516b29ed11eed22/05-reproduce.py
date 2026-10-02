from pathlib import Path
import os, subprocess, tempfile, shutil, sys, json
root=Path(sys.argv[1]).resolve()
node=shutil.which('node')
assert node, 'real Node required'
scratch=Path(tempfile.mkdtemp(prefix='roundfix-0211-doctor-repro-',dir='/private/tmp'))
(scratch/'home').mkdir(); (scratch/'bin').mkdir(); (scratch/'repo').mkdir()
subprocess.run(['git','init','-q',str(scratch/'repo')],check=True)
(scratch/'repo/bin').mkdir()
(scratch/'repo/bin/roundfix').symlink_to(root/'bin/roundfix')
# acpx is a Node program; record options before invoking real Node.
acpx=scratch/'bin/acpx'
acpx.write_text('#!/bin/sh\nprintf "%s\\n" "$NODE_OPTIONS" >> "$QA_NODE_RECORD"\ncase "$1" in --version) exec '+node+' -e \'console.log("0.12.0")\' 2>> "$QA_NODE_ERRORS";; *) exit 1;; esac\n')
acpx.chmod(0o700)
for name in ['node','codex','claude','npm','npx']:
 p=scratch/'bin'/name
 p.write_text('#!/bin/sh\ncase "$1" in --version) printf "v24.0.0\\n";; *) exit 1;; esac\n'); p.chmod(0o700)
env=dict(os.environ,HOME=str(scratch/'home'),PATH=str(scratch/'bin')+':/usr/bin:/bin',NODE_OPTIONS='--require='+str(scratch/'missing.cjs')+' --max-old-space-size=4096',QA_NODE_RECORD=str(scratch/'received.txt'),QA_NODE_ERRORS=str(scratch/'errors.txt'))
p=subprocess.run(['./bin/roundfix','doctor'],cwd=scratch/'repo',env=env,text=True,capture_output=True,timeout=40)
received=(scratch/'received.txt').read_text().splitlines()
print(json.dumps({'binary':str(root/'bin/roundfix'),'cwd':str(scratch/'repo'),'node':node,'exit':p.returncode,'received_NODE_OPTIONS':received,'real_node_stderr':(scratch/'errors.txt').read_text(),'stdout':p.stdout,'stderr':p.stderr},indent=2))
# Contract: every acpx invocation receives only kept options.
assert received and all(v=='--max-old-space-size=4096' for v in received), 'Doctor starts its first acpx with a dead preload'
