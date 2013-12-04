#!/usr/bin/python

import time,sys

Data = {}

timeout = time.time() + int(sys.argv[1])

def watch():
	while True:
		try:
			new = raw_input()
		except (EOFError,KeyboardInterrupt,SystemExit):
			yield(False)
		if new:
			yield(new)
		else:
			try:
				time.sleep(0.5)
			except (KeyboardInterrupt, SystemExit):
				yield(False)

count = 0
for l in watch():
	if l == False:
		break
	line = l.split()
	if len(line) != 3:
		continue
	d = line[0].split('/')
	if len(d) < 3:
		continue
	sent = float(line[1])
	recv = float(line[2])
	prog = '/'.join(d[0:-2])
	pid  = d[-2]
	uid  = d[-1]
	count += 1
	print '\r>> %d seconds remaining; %d lines processed' % (timeout - time.time(), count),
	sys.stdout.flush()
	#print pid, uid, prog, recv, sent
	if not Data.has_key(prog):
		Data[prog] = (0.0,0.0)
	Data[prog] = tuple(sum(x) for x in zip(Data[prog],(sent,recv)))
	if time.time() > timeout:
		break

print
print Data
