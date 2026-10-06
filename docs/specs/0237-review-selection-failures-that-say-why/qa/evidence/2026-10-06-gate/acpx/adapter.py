import json,os,sys
for line in sys.stdin:
 r=json.loads(line)
 if 'id' not in r: continue
 m=r.get('method'); response={'jsonrpc':'2.0','id':r['id']}
 fail=open(os.environ['QA_FAIL_FILE']).read().strip()
 if fail=='disconnect' and m=='session/prompt': sys.exit(7)
 if m==fail: response['error']={'code':-32603,'message':'QA adapter refused '+m,'data':{'marker':'not-diagnostic'}}
 elif m=='initialize': response['result']={'protocolVersion':1,'agentCapabilities':{'loadSession':True},'agentInfo':{'name':'qa-fake','version':'1'},'authMethods':[]}
 elif m in ('session/new','session/load'): response['result']={'sessionId':'qa-session','models':{'currentModelId':'default','availableModels':[{'modelId':'default','name':'Default'},{'modelId':'qa-model','name':'QA'}]}}
 else: response['result']={}
 print(json.dumps(response),flush=True)
