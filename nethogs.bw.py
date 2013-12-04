#!/usr/bin/python

import time,sys

def watch():
  while True:
    try: new = raw_input()
    except (EOFError, KeyboardInterrupt, SystemExit): yield(False)
    if new: yield(new)
    else:
      try: time.sleep(0.5)
      except (KeyboardInterrupt, SystemExit): yield(False)

class NH:
  def __init__(self, timeout=60, showprogress=False):
    self.data = {}
    self.timeout = time.time() + int(timeout)
    self.count = 0
    self.showprogress = showprogress
    self.runcomplete = False

  def parseline(self,line):
    ld = {}
    l = line.split()
    if len(l) != 3: return False
    d = l[0].split('/')
    if len(d) < 3: return False
    self.count += 1
    ld['sent'] = float(l[1])
    ld['recv'] = float(l[2])
    ld['prog'] = '/'.join(d[0:-2])
    ld['pid']  = d[-2]
    ld['uid']  = d[-1]
    return ld

  def run(self):
    for l in watch():
      if l == False: break
      line = parseline(l)
      if line == False: break
      if not self.data.has_key(line['prog']):
        self.data[line['prog']] = (0.0,0.0)
      zipped = zip(self.data[line['prog']], (line['sent'],line['recv']))
      self.data[line['prog']] = tuple(sum(x) for x in zipped)
      if self.showprogress: self.progress()
      if time.time() > self.timeout: break
    self.runcomplete = True

  def progress(self):
    print '\r>> %d seconds remaining; %d lines processed' % (timeout - time.time(), count),
    sys.stdout.flush()

  def get(self):
    return (self.runcomplete, self.data)

  def fmtitem(self, item):
    item = str(item)
    if len(item) < 7:
      item = item + "\t"
    item = item + "\t"
    return item

  def pp(self):
    if not self.runcomplete:
      print "call method run to populate data"
      return
    else:
      for prog in self.data:
        sent = self.data['prog'][0]
        recv = self.data['prog'][1]
        print self.fmtitem(prog),
        print self.fmtitem(sent),
        print self.fmtitem(recv)

if __name__ == "__main__":
  d = NH(timeout=60, showprogress=True)
  d.run()
  data = d.get()

