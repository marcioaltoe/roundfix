from pathlib import Path
import re,subprocess,os
base=Path(__file__).parent
sets={'skills':['skills/owned_skill_versions_test.go','skills/owned_skill_version_raise_test.go'],'internal/config':['internal/config/delivery_derived_lines_test.go','internal/config/verification_tools_test.go','internal/config/delivery_derived_paths_test.go'],'internal/cli':['internal/cli/deliver_conflict_test.go','internal/cli/deliver_derived_lines_test.go','internal/cli/deliver_derived_skill_layout_test.go','internal/cli/deliver_review_correction_test.go','internal/cli/review_lineage_test.go'],'internal/delivery':['internal/delivery/corrective_spec_test.go','internal/delivery/archived_head_retry_test.go','internal/delivery/review_correction_retry_test.go']}
procs=[]
for pkg,files in sets.items():
 names=[]
 for f in files:
  if Path(f).exists(): names+=re.findall(r'^func (Test\w+)\(',Path(f).read_text(),re.M)
 cmd=['go','test','-count=1','-v','./'+pkg,'-run','^('+'|'.join(names)+')$']
 log=Path('/private/tmp/roundfix-0228-qa01')/(pkg.replace('/','-')+'-tests.txt');h=log.open('w');h.write('GOCACHE='+os.environ['GOCACHE']+'\nCommand: '+repr(cmd)+'\n');h.flush()
 procs.append((pkg,subprocess.Popen(cmd,stdout=h,stderr=subprocess.STDOUT),h,log,names))
for pkg,proc,h,log,names in procs:
 code=proc.wait();h.write('\nExit: '+str(code)+'\n');h.close();text=log.read_text();missing=[n for n in names if '--- PASS: '+n+' ' not in text];print(pkg,code,'missing passes:',missing,flush=True)

import shutil
for pkg,proc,h,log,names in procs:shutil.copyfile(log,base/log.name)
